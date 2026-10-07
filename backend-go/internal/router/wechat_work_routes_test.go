// 公开字典（积分分类）与企微回调（验签 / AES 解密 / 事件接收）的路由契约与端到端行为。
//
// 覆盖要点：3 条新路由注册；`GET /common/score-categories` 无需登录、字段与 CategoryLabels 完全一致、
// sort 从 1 递增、响应不含 message；`GET /wechat-work/callback` 用**标准库自造的 msg_signature + AES-256-CBC
// 密文**验证「成功回显 echostr / 签名错 / 缺参 / 密文篡改」；`POST` 的 `<Encrypt>` 报文解密 → 事件解析 →
// 建请假记录 + 考勤置 leave（source=wechat_work），非审批事件、SpStatus 缺失、school_id 缺失与坏报文
// 仍返回 `{"errcode":0,"errmsg":"ok"}`。
//
// 假上游：企微接口地址由 WECHAT_WORK_API_BASE 指向 httptest.NewServer（服务构造时读环境变量，
// 故 env 必须在 router.New 之前设置）；**不访问真实外网**。
package router_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/config"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/database"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/router"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// wechatCallbackRoutes 本批 3 条公开路由（路径逐字对齐 Laravel routes/api.php 第 326、331-332 行）。
var wechatCallbackRoutes = []string{
	"GET /api/v1/common/score-categories",
	"GET /api/v1/wechat-work/callback",
	"POST /api/v1/wechat-work/callback",
}

// TestWechatCallbackRoutesRegistered 3 条路由注册齐全且无需 token。
func TestWechatCallbackRoutesRegistered(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	registered := map[string]bool{}
	for _, rt := range engine.Routes() {
		registered[rt.Method+" "+rt.Path] = true
	}
	for _, key := range wechatCallbackRoutes {
		assert.True(t, registered[key], "缺少路由 %s", key)
	}

	// 公开：不带 token 不应被鉴权中间件拦成 401。
	for _, key := range wechatCallbackRoutes {
		method, path := splitRouteKey(t, key)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(method, path, nil))
		assert.NotEqual(t, http.StatusUnauthorized, w.Code, "%s 是公开路由，不应 401", key)
	}
}

// TestScoreCategoriesEndpoint 分类字典：字段与 services.CategoryLabels 一致、sort 从 1 递增、无 message。
func TestScoreCategoriesEndpoint(t *testing.T) {
	_, engine, _ := newTestEngine(t)

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/common/score-categories", nil))
	require.Equal(t, http.StatusOK, w.Code)

	body := map[string]any{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	_, hasMessage := body["message"]
	assert.False(t, hasMessage, "Laravel 该接口不带 message")

	items, ok := body["data"].([]any)
	require.True(t, ok)
	require.Len(t, items, len(services.CategoryLabels), "分类条数应与唯一真源一致")

	seen := map[string]bool{}
	for i, raw := range items {
		item, ok := raw.(map[string]any)
		require.True(t, ok)
		id, _ := item["id"].(string)
		assert.Equal(t, float64(i+1), item["sort"], "sort 应从 1 递增")
		assert.Equal(t, services.CategoryLabels[id], item["name"], "%s 的名称应取 CategoryLabels", id)
		assert.NotEmpty(t, item["icon"], "%s 应有图标", id)
		seen[id] = true
	}
	assert.Equal(t, map[string]bool{
		"classroom": true, "homework": true, "behavior": true,
		"literacy": true, "daily": true, "academic": true, "custom": true,
	}, seen, "键集必须与 CategoryLabels 完全一致")
}

// TestWechatWorkLeaveRecordTable 新表随 AutoMigrate 建好，且列名与 Laravel 迁移一致。
func TestWechatWorkLeaveRecordTable(t *testing.T) {
	db, _, _ := newTestEngine(t)

	require.True(t, db.Migrator().HasTable("wechat_work_leave_records"))
	columns, err := db.Migrator().ColumnTypes(&models.WechatWorkLeaveRecord{})
	require.NoError(t, err)
	names := make([]string, 0, len(columns))
	for _, column := range columns {
		names = append(names, column.Name())
	}
	for _, want := range []string{
		"school_id", "class_id", "student_id", "parent_wework_userid", "student_name_from_wework",
		"sp_no", "leave_start_date", "leave_end_date", "leave_type", "reason", "approve_status",
		"approved_at", "raw_data", "synced_at",
	} {
		assert.Contains(t, names, want)
	}
}

// ============================================================
// 企微回调：加密链路（全部用标准库自造报文）
// ============================================================

// wechatTestAESKey 生成 43 字符的 EncodingAESKey（企微后台给出的就是这种形态：
// 32 字节 key 的 base64 去掉尾部 '='）。
func wechatTestAESKey() string {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte('a' + i%26)
	}
	return strings.TrimRight(base64.StdEncoding.EncodeToString(key), "=")
}

// wechatDecodeKey 同实现口径：base64(aes_key + "=") → 32 字节。
func wechatDecodeKey(t *testing.T, aesKey string) []byte {
	t.Helper()
	key, err := base64.StdEncoding.DecodeString(aesKey + "=")
	require.NoError(t, err)
	require.Len(t, key, 32)
	return key
}

// wechatEncrypt 按企微规格加密：明文 = 16 字节随机串 + 4 字节大端长度 + 正文，
// 再按 32 字节块做 PKCS7 填充（填充值 1..32），最后 AES-256-CBC（无额外填充）+ base64。
func wechatEncrypt(t *testing.T, aesKey, payload string) string {
	t.Helper()
	key := wechatDecodeKey(t, aesKey)

	buf := make([]byte, 0, len(payload)+64)
	buf = append(buf, bytes.Repeat([]byte{0x41}, 16)...)
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(payload)))
	buf = append(buf, length...)
	buf = append(buf, payload...)

	pad := 32 - len(buf)%32
	buf = append(buf, bytes.Repeat([]byte{byte(pad)}, pad)...)

	block, err := aes.NewCipher(key)
	require.NoError(t, err)
	out := make([]byte, len(buf))
	cipher.NewCBCEncrypter(block, key[:aes.BlockSize]).CryptBlocks(out, buf)
	return base64.StdEncoding.EncodeToString(out)
}

// wechatSignature 同实现口径：sha1(排序后拼接 token/timestamp/nonce/encrypt)。
func wechatSignature(token, timestamp, nonce, encrypt string) string {
	parts := []string{token, timestamp, nonce, encrypt}
	sort.Strings(parts)
	sum := sha1.Sum([]byte(strings.Join(parts, "")))
	return hex.EncodeToString(sum[:])
}

// wechatUpstream 起一个假企微上游（handler 为 nil 时一律 404）。
func wechatUpstream(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	if handler == nil {
		handler = func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }
	}
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

// newWechatCallbackEngine 建内存库 + 假企微上游 + 完整路由（env 必须在 router.New 之前注入）。
func newWechatCallbackEngine(t *testing.T, handler http.HandlerFunc) (*gorm.DB, *gin.Engine, models.School, models.Student) {
	t.Helper()

	upstream := wechatUpstream(t, handler)
	t.Setenv("WECHAT_WORK_API_BASE", upstream.URL)
	t.Setenv("WECHAT_WORK_CORPID", "corp-1")
	t.Setenv("WECHAT_WORK_SECRET", "sec-1")
	t.Setenv("WECHAT_WORK_TOKEN", "test-wechat-token")
	t.Setenv("WECHAT_WORK_ENCODING_AES_KEY", wechatTestAESKey())

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, database.Migrate(db))

	school := models.School{Name: "回调测试学校", Code: "ww-callback-school", Status: "active"}
	require.NoError(t, db.Create(&school).Error)
	teacher := models.User{
		SchoolID: school.ID, Role: "teacher", Username: "ww-callback-teacher", Name: "李老师", Status: "active",
	}
	require.NoError(t, db.Create(&teacher).Error)
	class := models.ClassRoom{SchoolID: school.ID, Name: "三年级（2）班", TeacherID: &teacher.ID, Status: "active"}
	require.NoError(t, db.Create(&class).Error)
	student := models.Student{ClassID: class.ID, Name: "小明", StudentNo: "0001", Status: "active"}
	require.NoError(t, db.Create(&student).Error)
	parent := models.User{
		SchoolID: school.ID, Role: "teacher", Username: "ww-callback-parent", Name: "家长", Status: "active",
	}
	require.NoError(t, db.Create(&parent).Error)
	require.NoError(t, db.Model(&models.Student{}).Where("id = ?", student.ID).
		Update("parent_id", parent.ID).Error)
	require.NoError(t, db.Create(&models.ThirdPartyBinding{
		UserID: parent.ID, Platform: "wechat_work", PlatformID: "parent-callback-1",
	}).Error)

	engine := router.New(db, &config.Config{JWTSecret: "test-secret", JWTExpHours: 1})
	return db, engine, school, student
}

// TestWechatWorkCallbackVerify GET 验签：成功回显 echostr 明文；签名错 / 缺参 / 密文篡改返回空串。
func TestWechatWorkCallbackVerify(t *testing.T) {
	_, engine, _, _ := newWechatCallbackEngine(t, nil)
	aesKey := wechatTestAESKey()
	token := "test-wechat-token"
	const timestamp, nonce = "1700000000", "nonce-1"

	encrypted := wechatEncrypt(t, aesKey, "<echostr>1293712937</echostr>")
	validSignature := wechatSignature(token, timestamp, nonce, encrypted)

	get := func(query url.Values) (int, string) {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/wechat-work/callback?"+query.Encode(), nil))
		return w.Code, w.Body.String()
	}

	t.Run("签名正确回显明文 echostr", func(t *testing.T) {
		status, body := get(url.Values{
			"msg_signature": {validSignature},
			"timestamp":     {timestamp},
			"nonce":         {nonce},
			"echostr":       {encrypted},
		})
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "1293712937", body)
	})

	t.Run("签名错误返回空串", func(t *testing.T) {
		status, body := get(url.Values{
			"msg_signature": {wechatSignature("wrong-token", timestamp, nonce, encrypted)},
			"timestamp":     {timestamp},
			"nonce":         {nonce},
			"echostr":       {encrypted},
		})
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "", body)
	})

	t.Run("缺参返回空串", func(t *testing.T) {
		status, body := get(url.Values{"timestamp": {timestamp}, "nonce": {nonce}})
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "", body)
	})

	t.Run("密文不可解返回空串", func(t *testing.T) {
		status, body := get(url.Values{
			"msg_signature": {wechatSignature(token, timestamp, nonce, "AAAA")},
			"timestamp":     {timestamp},
			"nonce":         {nonce},
			"echostr":       {"AAAA"},
		})
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "", body)
	})
}

// wechatEventPayload 拼一封企微事件明文（可控制 Event / SpNo / SpStatus，spStatus 为空表示元素缺失）。
func wechatEventPayload(event, spNo, spStatus string) string {
	spStatusXML := ""
	if spStatus != "" {
		spStatusXML = "<SpStatus>" + spStatus + "</SpStatus>"
	}
	return "<xml><ToUserName><![CDATA[corp-1]]></ToUserName>" +
		"<MsgType><![CDATA[event]]></MsgType>" +
		"<Event><![CDATA[" + event + "]]></Event>" +
		"<ApprovalInfo><SpNo><![CDATA[" + spNo + "]]></SpNo>" + spStatusXML + "</ApprovalInfo></xml>"
}

// TestWechatWorkCallbackReceive POST 接收：解密 → 事件解析 → 建请假记录 + 考勤置 leave。
func TestWechatWorkCallbackReceive(t *testing.T) {
	aesKey := wechatTestAESKey()
	token := "test-wechat-token"
	const timestamp, nonce = "1700000000", "nonce-1"

	detailCalls := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			_, _ = w.Write([]byte(`{"errcode":0,"access_token":"tok-1","expires_in":7200}`))
		case "/cgi-bin/oa/getapprovaldetail":
			detailCalls++
			_, _ = w.Write([]byte(`{"errcode":0,"info":{"sp_no":"SP-WEB-1","sp_status":2,` +
				`"applyer":{"userid":"parent-callback-1"},` +
				`"apply_data":{"contents":[{"title":"学生姓名","value":{"text":"小明"}},` +
				`{"title":"请假事由","value":{"text":"发烧"}}]}}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	db, engine, school, student := newWechatCallbackEngine(t, handler)
	schoolIDParam := strconv.FormatUint(uint64(school.ID), 10)

	post := func(payload, schoolID string) (int, string) {
		encrypted := wechatEncrypt(t, aesKey, payload)
		query := url.Values{
			"msg_signature": {wechatSignature(token, timestamp, nonce, encrypted)},
			"timestamp":     {timestamp},
			"nonce":         {nonce},
			"school_id":     {schoolID},
		}
		body := "<xml><Encrypt><![CDATA[" + encrypted + "]]></Encrypt></xml>"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wechat-work/callback?"+query.Encode(), strings.NewReader(body))
		req.Header.Set("Content-Type", "text/xml")
		engine.ServeHTTP(w, req)
		return w.Code, w.Body.String()
	}

	expectNoRecord := func(t *testing.T, spNo string) {
		t.Helper()
		var count int64
		require.NoError(t, db.Model(&models.WechatWorkLeaveRecord{}).Where("sp_no = ?", spNo).Count(&count).Error)
		assert.Equal(t, int64(0), count)
	}

	t.Run("审批事件建请假记录并落考勤", func(t *testing.T) {
		status, body := post(wechatEventPayload("sys_approval_change", "SP-WEB-1", "2"), schoolIDParam)
		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"errcode":0,"errmsg":"ok"}`, body)
		assert.Equal(t, 1, detailCalls)

		var rec models.WechatWorkLeaveRecord
		require.NoError(t, db.Where("sp_no = ?", "SP-WEB-1").First(&rec).Error)
		assert.Equal(t, "approved", rec.ApproveStatus)
		require.NotNil(t, rec.StudentID)
		assert.Equal(t, student.ID, *rec.StudentID)
		assert.Equal(t, "parent-callback-1", rec.ParentWeworkUserid)
		assert.Equal(t, "发烧", rec.Reason)
		require.NotNil(t, rec.SyncedAt)

		var att models.Attendance
		require.NoError(t, db.Where("student_id = ?", student.ID).First(&att).Error)
		assert.Equal(t, "leave", att.Status)
		assert.Equal(t, "wechat_work", att.Source)
		assert.Equal(t, "发烧", att.Remark)

		// 幂等：重复推送不新增记录、不再查详情。
		status, body = post(wechatEventPayload("sys_approval_change", "SP-WEB-1", "2"), schoolIDParam)
		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"errcode":0,"errmsg":"ok"}`, body)
		assert.Equal(t, 1, detailCalls)

		var count int64
		require.NoError(t, db.Model(&models.WechatWorkLeaveRecord{}).Count(&count).Error)
		assert.Equal(t, int64(1), count)
	})

	t.Run("非审批事件不落库", func(t *testing.T) {
		status, body := post(wechatEventPayload("click", "SP-WEB-2", "2"), schoolIDParam)
		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"errcode":0,"errmsg":"ok"}`, body)
		expectNoRecord(t, "SP-WEB-2")
	})

	t.Run("SpStatus 缺失按 1 处理（不落库）", func(t *testing.T) {
		status, body := post(wechatEventPayload("sys_approval_change", "SP-WEB-3", ""), schoolIDParam)
		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"errcode":0,"errmsg":"ok"}`, body)
		expectNoRecord(t, "SP-WEB-3")
	})

	t.Run("school_id 缺失时不处理", func(t *testing.T) {
		status, body := post(wechatEventPayload("sys_approval_change", "SP-WEB-4", "2"), "")
		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"errcode":0,"errmsg":"ok"}`, body)
		expectNoRecord(t, "SP-WEB-4")
	})

	t.Run("坏报文仍返回 ok", func(t *testing.T) {
		badBodies := []string{
			"",
			"not-xml-at-all",
			"<xml><Foo>1</Foo></xml>",
			"<xml><Encrypt>@@@</Encrypt></xml>",
			"<xml><Encrypt>" + wechatEncrypt(t, aesKey, "<xml><Foo/></xml>") + "</Encrypt></xml>",
		}
		for _, body := range badBodies {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/wechat-work/callback?school_id="+schoolIDParam+"&msg_signature=x&timestamp=1&nonce=1",
				strings.NewReader(body),
			)
			engine.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code, "报文 %q", body)
			assert.JSONEq(t, `{"errcode":0,"errmsg":"ok"}`, w.Body.String(), "报文 %q", body)
		}
	})
}
