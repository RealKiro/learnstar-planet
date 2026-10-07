// 第三方扫码登录 / 绑定 / 通讯录导入的服务层测试。
//
// 全部外部平台调用都指向 httptest 假上游（企微 / 钉钉 / 飞书），**不访问真实外网**；
// 用例覆盖：企微 gettoken 成功/失败、token 缓存命中不重复请求、code 换 userid、
// 部门列表 + 成员 + fetchContacts 去重与 department_names 映射；
// 钉钉 / 飞书的 getUserByCode 与 fetchContacts（含 authUrl 字符串逐字断言）；
// 平台选项三种形态；auth-url（未配置平台 / 无学校 / redirect_uri 透传）；
// 微信 / QQ / 人人通两条分支（已绑定无 token、未绑定 need_binding + temp_token 落库与 10 分钟过期）；
// 企微免注册建号；第三方平台登录四分支；bind-after-scan；绑定 / 解绑；通讯录导入统计与去重。
package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwtauth "github.com/RealKiro/learnstar-planet/backend-go/internal/auth"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// 假上游工具
// ============================================================

// recordingServer 启动假上游并记录每个请求的「方法 + 路径」，返回服务器与读取函数。
func recordingServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, func() []string) {
	t.Helper()
	var mu sync.Mutex
	seen := []string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen = append(seen, r.Method+" "+r.URL.Path)
		mu.Unlock()
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv, func() []string {
		mu.Lock()
		defer mu.Unlock()
		out := make([]string, len(seen))
		copy(out, seen)
		return out
	}
}

// jsonBody 写出 JSON 响应。
func jsonBody(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}

// newThirdPartyService 建第三方服务 + JWT 管理器（假上游由调用方注入）。
func newThirdPartyService(t *testing.T, db *gorm.DB) *services.ThirdParty {
	t.Helper()
	return services.NewThirdPartyService(db, jwtauth.New("test-secret", 1))
}

// ptr 返回字符串指针（可空字段用）。
func ptr(v string) *string { return &v }

// optionKeys 提取平台选项的 key 列表。
func optionKeys(options []services.PlatformOption) []string {
	keys := make([]string, 0, len(options))
	for _, option := range options {
		keys = append(keys, option.Key)
	}
	return keys
}

// marshalJSON 把值序列化为 JSON 文本（用于断言响应键集）。
func marshalJSON(t *testing.T, v any) string {
	t.Helper()
	raw, err := json.Marshal(v)
	require.NoError(t, err)
	return string(raw)
}

// mustHash 生成 bcrypt 哈希（建带已知密码的账号用）。
func mustHash(t *testing.T, plain string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

// ============================================================
// 企微服务
// ============================================================

// TestWechatWorkGetAccessToken 企微 gettoken：成功、缓存命中不重复请求、两种失败、未配置。
func TestWechatWorkGetAccessToken(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)

	calls := 0
	srv, seen := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			calls++
			assert.Equal(t, "corp-1", r.URL.Query().Get("corpid"))
			assert.Equal(t, "sec-1", r.URL.Query().Get("corpsecret"))
			jsonBody(w, `{"errcode":0,"errmsg":"ok","access_token":"AT-1","expires_in":7200}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	svc := services.NewWechatWorkService(db)
	svc.CorpID = "corp-1"
	svc.Secret = "sec-1"
	svc.APIBase = srv.URL
	svc.SetHTTPClient(srv.Client())

	token, err := svc.GetAccessToken(school.ID)
	require.NoError(t, err)
	assert.Equal(t, "AT-1", token)
	assert.Equal(t, 1, calls)

	// 缓存命中：不再请求上游。
	var cached models.WechatWorkToken
	require.NoError(t, db.First(&cached, school.ID).Error)
	assert.Equal(t, "AT-1", cached.Token)
	assert.True(t, cached.ExpiresAt.After(time.Now()), "TTL = expires_in - 300 秒，应仍未过期")

	token, err = svc.GetAccessToken(school.ID)
	require.NoError(t, err)
	assert.Equal(t, "AT-1", token)
	assert.Equal(t, 1, calls, "缓存命中不应重复请求 gettoken")
	assert.Len(t, seen(), 1)

	// 过期后重新取。
	require.NoError(t, db.Model(&models.WechatWorkToken{}).Where("school_id = ?", school.ID).
		Update("expires_at", time.Now().Add(-time.Minute)).Error)
	token, err = svc.GetAccessToken(school.ID)
	require.NoError(t, err)
	assert.Equal(t, "AT-1", token)
	assert.Equal(t, 2, calls)

	// 未配置 → 「企微未配置」。
	unconfigured := services.NewWechatWorkService(db)
	unconfigured.APIBase = srv.URL
	unconfigured.SetHTTPClient(srv.Client())
	_, err = unconfigured.GetAccessToken(school.ID + 99)
	require.Error(t, err)
	assert.Equal(t, "企微未配置", err.Error())

	// 无 access_token → 「token失败」。
	noToken, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonBody(w, `{"errcode":40013,"errmsg":"invalid corpid"}`)
	})
	bad := services.NewWechatWorkService(db)
	bad.CorpID = "corp-1"
	bad.Secret = "sec-1"
	bad.APIBase = noToken.URL
	bad.SetHTTPClient(noToken.Client())
	_, err = bad.GetAccessToken(school.ID + 100)
	require.Error(t, err)
	assert.Equal(t, "token失败", err.Error())

	// errcode != 0（带 access_token 字段也不认，同 Laravel `!isset` 之前不判 errcode 的部分）：
	// Laravel 只判 `!isset($r['access_token'])`，故 errcode 非 0 但有 token 时仍会返回该 token。
	withCode, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonBody(w, `{"errcode":40001,"errmsg":"bad","access_token":"AT-X","expires_in":600}`)
	})
	codeSvc := services.NewWechatWorkService(db)
	codeSvc.CorpID = "corp-1"
	codeSvc.Secret = "sec-1"
	codeSvc.APIBase = withCode.URL
	codeSvc.SetHTTPClient(withCode.Client())
	token, err = codeSvc.GetAccessToken(school.ID + 101)
	require.NoError(t, err)
	assert.Equal(t, "AT-X", token)
}

// TestWechatWorkGetUserIDByCode code 换 userid：UserId 优先、userid 兜底、errcode 非 0 → 空串。
func TestWechatWorkGetUserIDByCode(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)

	srv, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cgi-bin/gettoken" {
			jsonBody(w, `{"errcode":0,"access_token":"AT-1","expires_in":7200}`)
			return
		}
		switch r.URL.Query().Get("code") {
		case "code-upper":
			jsonBody(w, `{"errcode":0,"UserId":"zhangsan"}`)
		case "code-lower":
			jsonBody(w, `{"errcode":0,"userid":"lisi"}`)
		default:
			jsonBody(w, `{"errcode":40029,"errmsg":"invalid code"}`)
		}
	})

	svc := services.NewWechatWorkService(db)
	svc.CorpID, svc.Secret, svc.APIBase = "corp-1", "sec-1", srv.URL
	svc.SetHTTPClient(srv.Client())

	userid, err := svc.GetUserIDByCode(school.ID, "code-upper")
	require.NoError(t, err)
	assert.Equal(t, "zhangsan", userid)

	userid, err = svc.GetUserIDByCode(school.ID, "code-lower")
	require.NoError(t, err)
	assert.Equal(t, "lisi", userid)

	userid, err = svc.GetUserIDByCode(school.ID, "code-bad")
	require.NoError(t, err)
	assert.Empty(t, userid, "errcode 非 0 时返回空串（不抛错）")
}

// TestWechatWorkFetchContacts 部门 + 成员：只遍历顶级部门、按 userid 去重、department_names 映射。
func TestWechatWorkFetchContacts(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)

	srv, seen := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"AT-1","expires_in":7200}`)
		case "/cgi-bin/department/list":
			jsonBody(w, `{"errcode":0,"department":[
				{"id":1,"name":"总部","parentid":0},
				{"id":2,"name":"一年级组","parentid":1}
			]}`)
		case "/cgi-bin/user/list":
			assert.Equal(t, "1", r.URL.Query().Get("fetch_child"), "顶级部门必须 fetch_child=1")
			jsonBody(w, `{"errcode":0,"userlist":[
				{"userid":"u1","name":"张三","mobile":"138","email":"a@b.c","position":"教师","department":[1,2]},
				{"userid":"u1","name":"张三（重复）","department":[2]},
				{"name":"无 userid 的成员"},
				{"userid":"u2","name":"李四","department":[2,9]}
			]}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	svc := services.NewWechatWorkService(db)
	svc.CorpID, svc.Secret, svc.APIBase = "corp-1", "sec-1", srv.URL
	svc.SetHTTPClient(srv.Client())

	contacts, err := svc.FetchContacts(school.ID)
	require.NoError(t, err)
	require.Len(t, contacts.Departments, 2)
	require.Len(t, contacts.Members, 2, "userid 去重 + 丢弃无 userid 的成员")

	first := contacts.Members[0]
	assert.Equal(t, "u1", first.UserID)
	assert.Equal(t, "张三", first.Name)
	assert.Equal(t, "138", first.Mobile)
	assert.Equal(t, "教师", first.Position)
	assert.Equal(t, []string{"总部", "一年级组"}, first.DepartmentNames)

	second := contacts.Members[1]
	assert.Equal(t, "u2", second.UserID)
	assert.Equal(t, []string{"一年级组"}, second.DepartmentNames, "未知部门 ID 不产生空名条目")

	// 只遍历顶级部门（id=2 的 parentid != 0 不再单独查询）。
	userListCalls := 0
	for _, req := range seen() {
		if req == "GET /cgi-bin/user/list" {
			userListCalls++
		}
	}
	assert.Equal(t, 1, userListCalls)

	// 部门接口报错 → 「获取企微部门失败：xxx」。
	failing, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cgi-bin/gettoken" {
			jsonBody(w, `{"errcode":0,"access_token":"AT-2","expires_in":7200}`)
			return
		}
		jsonBody(w, `{"errcode":60011,"errmsg":"no privilege"}`)
	})
	bad := services.NewWechatWorkService(db)
	bad.CorpID, bad.Secret, bad.APIBase = "corp-1", "sec-1", failing.URL
	bad.SetHTTPClient(failing.Client())
	_, err = bad.FetchContacts(school.ID)
	require.Error(t, err)
	assert.Equal(t, "获取企微部门失败：no privilege", err.Error())
}

// ============================================================
// provider 管理器 / authUrl / 钉钉 / 飞书
// ============================================================

// TestThirdPartyManagerProviderFor 未配置 / 未知平台报错文案同 Laravel。
func TestThirdPartyManagerProviderFor(t *testing.T) {
	db := setupDB(t)
	svc := newThirdPartyService(t, db)

	_, err := svc.Manager().ProviderFor("")
	require.Error(t, err)
	assert.Equal(t, "未配置或未知的第三方平台", err.Error())

	_, err = svc.Manager().ProviderFor("weibo")
	require.Error(t, err)

	for platform, wantKey := range map[string]string{
		"wechat_work": "wechat_work",
		"dingtalk":    "dingtalk",
		"feishu":      "feishu",
	} {
		provider, err := svc.Manager().ProviderFor(platform)
		require.NoError(t, err)
		assert.Equal(t, wantKey, provider.Key())
	}
}

// TestProviderAuthURLs authUrl 字符串逐字断言（含 urlencode 与 state=schoolId）。
func TestProviderAuthURLs(t *testing.T) {
	db := setupDB(t)
	svc := newThirdPartyService(t, db)
	svc.WechatWork().CorpID = "ww1234"
	svc.WechatWork().AgentID = 1000014
	svc.Manager().DingTalkAppKey = "ding-app"
	svc.Manager().LarkAppID = "cli_app"

	wechatWork, err := svc.Manager().ProviderFor("wechat_work")
	require.NoError(t, err)
	assert.Equal(t,
		"https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=ww1234&agentid=1000014"+
			"&redirect_uri=https%3A%2F%2Fexample.com%2Flogin&state=7",
		wechatWork.AuthURL(7, "https://example.com/login"))

	dingtalk, err := svc.Manager().ProviderFor("dingtalk")
	require.NoError(t, err)
	assert.Equal(t,
		"https://login.dingtalk.com/oauth2/auth?redirect_uri=https%3A%2F%2Fexample.com%2Flogin"+
			"&response_type=code&client_id=ding-app&scope=openid&state=7",
		dingtalk.AuthURL(7, "https://example.com/login"))

	feishu, err := svc.Manager().ProviderFor("feishu")
	require.NoError(t, err)
	assert.Equal(t,
		"https://accounts.feishu.cn/open-apis/authen/v1/authorize?client_id=cli_app"+
			"&redirect_uri=https%3A%2F%2Fexample.com%2Flogin&response_type=code&state=7",
		feishu.AuthURL(7, "https://example.com/login"))
}

// TestDingTalkProvider 钉钉 getUserByCode / fetchContacts（含失败分支）。
func TestDingTalkProvider(t *testing.T) {
	db := setupDB(t)
	svc := newThirdPartyService(t, db)

	var mu sync.Mutex
	meToken := ""
	srv, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1.0/oauth2/userAccessToken":
			jsonBody(w, `{"accessToken":"dt-uat","expireIn":7200}`)
		case "/v1.0/contact/users/me":
			mu.Lock()
			meToken = r.Header.Get("x-acs-dingtalk-access-token")
			mu.Unlock()
			jsonBody(w, `{"unionId":"dt-union-1","userId":"dt-user","nick":"钉钉老师","mobile":"137","email":"dt@x.c","avatarUrl":"https://a/1.png"}`)
		case "/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"dt-internal"}`)
		case "/topapi/v2/department/listsub":
			jsonBody(w, `{"errcode":0,"result":{"list":[
				{"dept_id":1,"name":"总部","parent_id":0},
				{"dept_id":2,"name":"一年级组","parent_id":1}
			]}}`)
		case "/topapi/v2/user/list":
			jsonBody(w, `{"errcode":0,"result":{"list":[
				{"userid":"dt-u1","name":"钉钉张三","mobile":"136","email":"a@b.c","title":"语文老师","dept_id_list":[1,2]},
				{"userid":"dt-u1","name":"重复"},
				{"userid":"dt-u2","name":"钉钉李四","dept_id_list":[2]}
			]}}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	manager := svc.Manager()
	manager.SetHTTPClient(srv.Client())
	manager.DingTalkAppKey = "ding-app"
	manager.DingTalkAppSecret = "ding-secret"
	manager.DingTalkOAuthBase = srv.URL
	manager.DingTalkAPIBase = srv.URL

	provider, err := manager.ProviderFor("dingtalk")
	require.NoError(t, err)

	info, err := provider.GetUserByCode(1, "code-1")
	require.NoError(t, err)
	assert.Equal(t, "dt-union-1", info.PlatformID, "unionId 优先")
	assert.Equal(t, "钉钉老师", info.Name)
	assert.Equal(t, "137", info.Mobile)
	assert.Equal(t, "https://a/1.png", info.Avatar)
	mu.Lock()
	assert.Equal(t, "dt-uat", meToken)
	mu.Unlock()

	contacts, err := provider.FetchContacts(1)
	require.NoError(t, err)
	require.Len(t, contacts.Members, 2, "按 userid 去重")
	assert.Equal(t, "钉钉张三", contacts.Members[0].Name)
	assert.Equal(t, "语文老师", contacts.Members[0].Position, "position 取 title")
	assert.Equal(t, []string{"总部", "一年级组"}, contacts.Members[0].DepartmentNames)
	assert.Equal(t, []string{"一年级组"}, contacts.Members[1].DepartmentNames)

	// 拿不到企业内部 token → 业务错误文案。
	failing, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonBody(w, `{"errcode":40001,"errmsg":"bad appkey"}`)
	})
	failManager := services.NewThirdPartyManager(nil)
	failManager.SetHTTPClient(failing.Client())
	failManager.DingTalkAPIBase = failing.URL
	failManager.DingTalkOAuthBase = failing.URL
	failProvider, err := failManager.ProviderFor("dingtalk")
	require.NoError(t, err)
	_, err = failProvider.FetchContacts(1)
	require.Error(t, err)
	assert.Equal(t, "钉钉未配置或获取企业内部 token 失败", err.Error())

	// 登录失败：无 accessToken。
	noToken, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonBody(w, `{"message":"code 无效"}`)
	})
	failManager.SetHTTPClient(noToken.Client())
	failManager.DingTalkOAuthBase = noToken.URL
	failProvider, err = failManager.ProviderFor("dingtalk")
	require.NoError(t, err)
	_, err = failProvider.GetUserByCode(1, "bad")
	require.Error(t, err)
	assert.Equal(t, "钉钉登录失败：code 无效", err.Error())
}

// TestLarkProvider 飞书 getUserByCode / fetchContacts（含失败分支）。
func TestLarkProvider(t *testing.T) {
	db := setupDB(t)
	svc := newThirdPartyService(t, db)

	srv, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/v3/token":
			jsonBody(w, `{"access_token":"lark-uat"}`)
		case "/open-apis/authen/v1/user_info":
			assert.Equal(t, "Bearer lark-uat", r.Header.Get("Authorization"))
			jsonBody(w, `{"data":{"union_id":"lark-union-1","name":"飞书老师","mobile":"135","email":"fs@x.c","avatar_url":"https://a/2.png"}}`)
		case "/open-apis/auth/v3/tenant_access_token/internal":
			jsonBody(w, `{"code":0,"tenant_access_token":"lark-tenant"}`)
		case "/open-apis/contact/v3/departments":
			assert.Equal(t, "true", r.URL.Query().Get("fetch_child"))
			jsonBody(w, `{"code":0,"data":{"items":[
				{"department_id":"od-1","name":"总部","parent_department_id":"0"},
				{"department_id":"od-2","name":"一年级组"}
			]}}`)
		case "/open-apis/contact/v3/users/find_by_department":
			jsonBody(w, `{"code":0,"data":{"items":[
				{"union_id":"lark-u1","name":"飞书张三","mobile":"134","email":"x@y.z","job_title":"数学老师","department_ids":["od-1","od-2"]},
				{"union_id":"lark-u1","name":"重复"},
				{"user_id":"lark-u2","name":"飞书李四","department_ids":["od-2"]}
			]}}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	manager := svc.Manager()
	manager.SetHTTPClient(srv.Client())
	manager.LarkAppID = "cli-1"
	manager.LarkAppSecret = "cli-secret"
	manager.LarkAccountsBase = srv.URL
	manager.LarkAPIBase = srv.URL

	provider, err := manager.ProviderFor("feishu")
	require.NoError(t, err)

	info, err := provider.GetUserByCode(1, "code-1")
	require.NoError(t, err)
	assert.Equal(t, "lark-union-1", info.PlatformID)
	assert.Equal(t, "飞书老师", info.Name)
	assert.Equal(t, "135", info.Mobile)
	assert.Equal(t, "https://a/2.png", info.Avatar)

	contacts, err := provider.FetchContacts(1)
	require.NoError(t, err)
	require.Len(t, contacts.Members, 2, "union_id 去重 + user_id 兜底")
	assert.Equal(t, "飞书张三", contacts.Members[0].Name)
	assert.Equal(t, "数学老师", contacts.Members[0].Position)
	assert.Equal(t, []string{"总部", "一年级组"}, contacts.Members[0].DepartmentNames)
	assert.Equal(t, "lark-u2", contacts.Members[1].UserID)

	// tenant_access_token 缺失 → 业务错误文案。
	failing, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonBody(w, `{"code":10003,"msg":"invalid app_secret"}`)
	})
	failManager := services.NewThirdPartyManager(nil)
	failManager.SetHTTPClient(failing.Client())
	failManager.LarkAccountsBase = failing.URL
	failManager.LarkAPIBase = failing.URL
	failProvider, err := failManager.ProviderFor("feishu")
	require.NoError(t, err)
	_, err = failProvider.FetchContacts(1)
	require.Error(t, err)
	assert.Equal(t, "飞书未配置或获取 tenant_access_token 失败", err.Error())

	// 登录失败：无 access_token → 飞书登录失败：<error_description>。
	noToken, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonBody(w, `{"error":"invalid_grant","error_description":"code 已失效"}`)
	})
	failManager.SetHTTPClient(noToken.Client())
	failManager.LarkAccountsBase = noToken.URL
	failProvider, err = failManager.ProviderFor("feishu")
	require.NoError(t, err)
	_, err = failProvider.GetUserByCode(1, "bad")
	require.Error(t, err)
	assert.Equal(t, "飞书登录失败：code 已失效", err.Error())
}

// ============================================================
// 平台选项 / auth-url
// ============================================================

// TestThirdPartyOptions 三种形态：默认三平台 / 学校开关交集 / 颜色表。
func TestThirdPartyOptions(t *testing.T) {
	db := setupDB(t)
	svc := newThirdPartyService(t, db)

	// 无学校 → 默认三平台（Laravel: count($platforms) === 0 → 默认）。
	options, err := svc.Options()
	require.NoError(t, err)
	require.Len(t, options, 3)
	assert.Equal(t, []string{"wechat_work", "wechat", "qq"}, optionKeys(options))
	assert.Equal(t, "#2B7CE9", options[0].Color)
	assert.Equal(t, "企业微信", options[0].Label)
	assert.Equal(t, "🏢", options[0].Icon)
	assert.Equal(t, "#07C160", options[1].Color)
	assert.Equal(t, "#12B7F5", options[2].Color)

	school := seedSchool(t, db)

	// 学校开关为交集（顺序按 platforms()）。
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("settings", `{"enabled_third_party_platforms":["wechat","renren","feishu"]}`).Error)
	options, err = svc.Options()
	require.NoError(t, err)
	assert.Equal(t, []string{"wechat", "renren", "feishu"}, optionKeys(options))
	assert.Equal(t, "#FF6A00", options[1].Color)
	assert.Equal(t, "#3370FF", options[2].Color)

	// 开关里全是未知平台 → 交集为空 → 回落默认三平台。
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("settings", `{"enabled_third_party_platforms":["weibo"]}`).Error)
	options, err = svc.Options()
	require.NoError(t, err)
	assert.Equal(t, []string{"wechat_work", "wechat", "qq"}, optionKeys(options))

	// 空数组 → 默认三平台。
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("settings", `{"enabled_third_party_platforms":[]}`).Error)
	options, err = svc.Options()
	require.NoError(t, err)
	assert.Equal(t, []string{"wechat_work", "wechat", "qq"}, optionKeys(options))

	// 颜色与标签表：逐个平台单独启用后断言（dingtalk / feishu 也走同一 match 表）。
	for platform, want := range map[string]struct{ color, label, icon string }{
		"wechat_work": {"#2B7CE9", "企业微信", "🏢"},
		"dingtalk":    {"#0089FF", "钉钉", "🔷"},
		"feishu":      {"#3370FF", "飞书", "🪶"},
		"wechat":      {"#07C160", "微信", "💬"},
		"qq":          {"#12B7F5", "QQ", "🐧"},
		"renren":      {"#FF6A00", "人人通空间", "🌐"},
	} {
		require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
			Update("settings", `{"enabled_third_party_platforms":["`+platform+`"]}`).Error)
		options, err := svc.Options()
		require.NoError(t, err)
		require.Len(t, options, 1)
		assert.Equal(t, want.color, options[0].Color, "平台 %s 的品牌色", platform)
		assert.Equal(t, want.label, options[0].Label, "平台 %s 的显示名", platform)
		assert.Equal(t, want.icon, options[0].Icon, "平台 %s 的图标", platform)
	}
}

// TestThirdPartyAuthURLOfSchool auth-url：按 school_id 取学校、未配置平台报错、redirect_uri 透传。
func TestThirdPartyAuthURLOfSchool(t *testing.T) {
	db := setupDB(t)
	svc := newThirdPartyService(t, db)

	// 无学校 → (nil, nil)，由控制器返回 400「系统尚未初始化」。
	school, err := svc.SchoolForAuthURL(0)
	require.NoError(t, err)
	assert.Nil(t, school)

	first := seedSchool(t, db)
	second := models.School{Name: "第二学校", Code: "second-school", Status: "active"}
	require.NoError(t, db.Create(&second).Error)

	school, err = svc.SchoolForAuthURL(0)
	require.NoError(t, err)
	require.NotNil(t, school)
	assert.Equal(t, first.ID, school.ID, "未传 school_id 取第一所学校")

	school, err = svc.SchoolForAuthURL(second.ID)
	require.NoError(t, err)
	require.NotNil(t, school)
	assert.Equal(t, second.ID, school.ID)

	// 不存在的 school_id → 回退第一所学校（同 Laravel `School::find($id)` 为 null 时…）
	// 注：Laravel 是 `$schoolId > 0 ? School::find($schoolId) : School::first()`，
	// find 返回 null 时不会回退，直接走 400。Go 端保持一致：返回 nil。
	school, err = svc.SchoolForAuthURL(999999)
	require.NoError(t, err)
	assert.Nil(t, school, "指定了不存在的 school_id → 400「系统尚未初始化」（与 Laravel 一致，不回退）")

	// 未配置平台 → 报错，由控制器转成 400 固定文案。
	_, err = svc.AuthURLOfSchool(&first, "https://example.com/login")
	require.Error(t, err)
	assert.Equal(t, "未配置或未知的第三方平台", err.Error())

	require.NoError(t, db.Model(&models.School{}).Where("id = ?", first.ID).
		Update("settings", `{"third_party_platform":"dingtalk"}`).Error)
	svc.Manager().DingTalkAppKey = "ding-app"

	require.NoError(t, db.First(&first, first.ID).Error)
	result, err := svc.AuthURLOfSchool(&first, "https://example.com/redirect")
	require.NoError(t, err)
	assert.Equal(t, "dingtalk", result.Platform)
	assert.Contains(t, result.AuthURL, "redirect_uri=https%3A%2F%2Fexample.com%2Fredirect")
	assert.Contains(t, result.AuthURL, "client_id=ding-app")

	// state 解析：正确 ID → 命中；非法/不存在 → 回退第一所学校。
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", second.ID).
		Update("settings", `{"third_party_platform":"feishu"}`).Error)
	resolved, err := svc.ResolveSchoolFromState(itoa(second.ID))
	require.NoError(t, err)
	require.NotNil(t, resolved)
	assert.Equal(t, second.ID, resolved.ID)

	resolved, err = svc.ResolveSchoolFromState("abc")
	require.NoError(t, err)
	require.NotNil(t, resolved)
	assert.Equal(t, first.ID, resolved.ID, "非法 state 回退第一所学校")

	resolved, err = svc.ResolveSchoolFromState("999999")
	require.NoError(t, err)
	require.NotNil(t, resolved)
	assert.Equal(t, first.ID, resolved.ID, "不存在的 state 回退第一所学校")

	resolved, err = svc.ResolveSchoolFromState("2abc")
	require.NoError(t, err)
	require.NotNil(t, resolved)
	assert.Equal(t, second.ID, resolved.ID, "PHP `(int)` 语义：前导数字生效")
}

// TestLoginURL APP_URL 决定 url('/login') 的等价物。
func TestLoginURL(t *testing.T) {
	t.Setenv("APP_URL", "")
	assert.Equal(t, "http://localhost/login", services.LoginURL())
	t.Setenv("APP_URL", "https://learnstar.example.com")
	assert.Equal(t, "https://learnstar.example.com/login", services.LoginURL())
}

// ============================================================
// 微信 / QQ / 人人通扫码登录
// ============================================================

// TestLoginWithWechatBranches 已绑定（无 token）/ unionid 命中 / 未绑定 need_binding。
func TestLoginWithWechatBranches(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)

	// 未绑定 → need_binding + temp_token 落库（10 分钟）。
	result, err := svc.LoginWithWechat("wx-open-1", nil, ptr("微信昵称"), ptr("https://a/1.png"))
	require.NoError(t, err)
	need, ok := result.(services.WechatNeedBindingResult)
	require.True(t, ok, "未绑定应返回 WechatNeedBindingResult，实际 %T", result)
	assert.Equal(t, "need_binding", need.Status)
	assert.Equal(t, "wx-open-1", need.OpenID)
	assert.Nil(t, need.UnionID, "unionid 未传时为 null")
	require.NotEmpty(t, need.TempToken)

	var row models.TempBindingContext
	require.NoError(t, db.First(&row, "temp_token = ?", need.TempToken).Error)
	assert.WithinDuration(t, time.Now().Add(10*time.Minute), row.ExpiresAt, 30*time.Second)
	assert.Contains(t, row.Context, `"platform":"wechat"`)
	assert.Contains(t, row.Context, `"platform_id":"wx-open-1"`)

	// 已绑定（openid）→ logged_in，**不带 token**。
	teacher := seedTeacher(t, db, school.ID, "wechat-teacher")
	require.NoError(t, db.Create(&models.ThirdPartyBinding{
		UserID: teacher.ID, Platform: "wechat", PlatformID: "wx-open-1",
	}).Error)

	result, err = svc.LoginWithWechat("wx-open-1", nil, nil, nil)
	require.NoError(t, err)
	loggedIn, ok := result.(services.LoggedInResult)
	require.True(t, ok, "已绑定应返回 LoggedInResult，实际 %T", result)
	assert.Equal(t, "logged_in", loggedIn.Status)
	require.NotNil(t, loggedIn.User)
	assert.Equal(t, teacher.ID, loggedIn.User.ID)

	var reloaded models.User
	require.NoError(t, db.First(&reloaded, teacher.ID).Error)
	assert.NotNil(t, reloaded.LastLoginAt, "已绑定登录应刷新 last_login_at")

	// unionid 优先（即使 openid 未绑定）。
	require.NoError(t, db.Create(&models.ThirdPartyBinding{
		UserID: teacher.ID, Platform: "wechat", PlatformID: "wx-open-2", PlatformUnionID: "wx-union-1",
	}).Error)
	result, err = svc.LoginWithWechat("wx-open-3", ptr("wx-union-1"), nil, nil)
	require.NoError(t, err)
	loggedIn, ok = result.(services.LoggedInResult)
	require.True(t, ok)
	assert.Equal(t, teacher.ID, loggedIn.User.ID, "unionid 命中同一账号")

	// 序列化断言：logged_in 分支 JSON 不含 token；need_binding 含 openid/unionid。
	payload := marshalJSON(t, loggedIn)
	assert.NotContains(t, payload, "token")
	assert.Contains(t, payload, `"status":"logged_in"`)
}

// TestTempBindingContextExpiry 过期上下文视为不存在（绑定走请求体参数），且写入时惰性清理过期行。
func TestTempBindingContextExpiry(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)

	teacher := seedCredentialUser(t, db, school.ID, "teacher", "过期老师", "ls123456", "active")

	scanned, err := svc.LoginWithWechat("wx-open-1", nil, ptr("昵称"), nil)
	require.NoError(t, err)
	need, ok := scanned.(services.WechatNeedBindingResult)
	require.True(t, ok)
	require.NotEmpty(t, need.TempToken)

	// 手动把上下文改成已过期。
	require.NoError(t, db.Model(&models.TempBindingContext{}).Where("temp_token = ?", need.TempToken).
		Update("expires_at", time.Now().Add(-time.Second)).Error)

	// 过期 → 视为不存在：绑定与昵称覆盖全部走请求体参数（platform=renren / nick=请求体昵称）。
	bound, err := svc.BindAfterScan(
		need.TempToken, "过期老师", "ls123456", "renren", "rr-expired", nil, ptr("过期后的昵称"), nil,
	)
	require.NoError(t, err)
	assert.Equal(t, "bound", bound.Status)
	require.NotNil(t, bound.User)

	var binding models.ThirdPartyBinding
	require.NoError(t, db.Where("user_id = ? AND platform = ?", teacher.ID, "renren").First(&binding).Error)
	assert.Equal(t, "rr-expired", binding.PlatformID, "上下文过期 → 使用请求体里的 platform/platform_id")

	var wechatBindings int64
	require.NoError(t, db.Model(&models.ThirdPartyBinding{}).
		Where("user_id = ? AND platform = ?", teacher.ID, "wechat").Count(&wechatBindings).Error)
	assert.Equal(t, int64(0), wechatBindings, "上下文过期 → 不按 wechat 记绑定")

	// 过期行由写入路径惰性清理。
	stale := models.TempBindingContext{
		TempToken: "stale-token", Context: "{}", ExpiresAt: time.Now().Add(-time.Hour),
	}
	require.NoError(t, db.Create(&stale).Error)
	_, err = svc.LoginWithQQ("qq-open-2", nil, nil)
	require.NoError(t, err)

	var count int64
	require.NoError(t, db.Model(&models.TempBindingContext{}).
		Where("temp_token = ?", "stale-token").Count(&count).Error)
	assert.Equal(t, int64(0), count, "写入时惰性清理过期行")
}

// TestLoginWithQQAndRenrenBranches QQ / 人人通：未绑定 need_binding（键集不同），已绑定无 token。
func TestLoginWithQQAndRenrenBranches(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)
	teacher := seedTeacher(t, db, school.ID, "qq-teacher")

	// QQ 未绑定：{status, temp_token, openid}
	result, err := svc.LoginWithQQ("qq-open-1", ptr("QQ昵称"), nil)
	require.NoError(t, err)
	qqNeed, ok := result.(services.QQNeedBindingResult)
	require.True(t, ok, "实际 %T", result)
	assert.Equal(t, "need_binding", qqNeed.Status)
	assert.Equal(t, "qq-open-1", qqNeed.OpenID)
	require.NotEmpty(t, qqNeed.TempToken)
	payload := marshalJSON(t, qqNeed)
	assert.NotContains(t, payload, "unionid", "QQ 分支不带 unionid 键")
	assert.NotContains(t, payload, "platform_id", "QQ 分支不带 platform_id 键")

	require.NoError(t, db.Create(&models.ThirdPartyBinding{
		UserID: teacher.ID, Platform: "qq", PlatformID: "qq-open-1",
	}).Error)
	result, err = svc.LoginWithQQ("qq-open-1", nil, nil)
	require.NoError(t, err)
	qqLoggedIn, ok := result.(services.LoggedInResult)
	require.True(t, ok, "实际 %T", result)
	assert.Equal(t, teacher.ID, qqLoggedIn.User.ID)
	assert.NotContains(t, marshalJSON(t, qqLoggedIn), "token", "已绑定分支不带 token（Laravel 原样行为）")

	// 人人通未绑定：{status, temp_token, platform_id}
	result, err = svc.LoginWithRenren("rr-user-1", ptr("人人通昵称"), nil)
	require.NoError(t, err)
	rrNeed, ok := result.(services.RenrenNeedBindingResult)
	require.True(t, ok, "实际 %T", result)
	assert.Equal(t, "need_binding", rrNeed.Status)
	assert.Equal(t, "rr-user-1", rrNeed.PlatformID)
	require.NotEmpty(t, rrNeed.TempToken)
	assert.NotContains(t, marshalJSON(t, rrNeed), "openid", "人人通分支不带 openid 键")

	require.NoError(t, db.Create(&models.ThirdPartyBinding{
		UserID: teacher.ID, Platform: "renren", PlatformID: "rr-user-1",
	}).Error)
	result, err = svc.LoginWithRenren("rr-user-1", nil, nil)
	require.NoError(t, err)
	rrLoggedIn, ok := result.(services.LoggedInResult)
	require.True(t, ok)
	assert.Equal(t, teacher.ID, rrLoggedIn.User.ID)
}

// ============================================================
// 企微冒烟登录（免注册建号）
// ============================================================

// TestLoginWithWechatWorkAutoRegister 未绑定 → 免注册建号（明文密码 + 绑定行 + token）；
// 已绑定 → 只返回 {status, user}（不带 token）。
func TestLoginWithWechatWorkAutoRegister(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)

	// nick 未传 → 姓名兜底为 userid。
	result, err := svc.LoginWithWechatWork("ww-user-1", nil, nil, &school)
	require.NoError(t, err)
	assert.Equal(t, "logged_in", result.Status)
	require.NotEmpty(t, result.Token, "未绑定分支必须返回 token")
	require.NotNil(t, result.User)
	assert.Equal(t, "ww-user-1", result.User.Name)
	assert.Equal(t, "ww-user-1", result.User.Username)
	assert.Equal(t, "ww-user-1", result.User.Nickname)
	assert.Equal(t, "teacher", result.User.Role)
	assert.Equal(t, "active", result.User.Status)
	require.NotNil(t, result.User.LastLoginAt)

	var stored models.User
	require.NoError(t, db.First(&stored, result.User.ID).Error)
	assert.Equal(t, services.DefaultTeacherPassword, stored.PlainPassword, "明文密码与默认密码一致")
	assert.NotEmpty(t, stored.PasswordHash)

	var binding models.ThirdPartyBinding
	require.NoError(t, db.Where("platform = ? AND platform_id = ?", "wechat_work", "ww-user-1").
		First(&binding).Error)
	assert.Equal(t, stored.ID, binding.UserID)
	require.NotNil(t, binding.VerifiedAt)

	// token 可用（同 JWT 校验路径）。
	claims, err := jwtauth.New("test-secret", 1).Parse(result.Token)
	require.NoError(t, err)
	assert.Equal(t, stored.ID, claims.UserID)

	// 已绑定 → 无 token。
	again, err := svc.LoginWithWechatWork("ww-user-1", nil, nil, &school)
	require.NoError(t, err)
	assert.Equal(t, "logged_in", again.Status)
	assert.Empty(t, again.Token, "已绑定分支不带 token（Laravel 原样行为）")
	require.NotNil(t, again.User)
	assert.Equal(t, stored.ID, again.User.ID)
	var userCount int64
	require.NoError(t, db.Model(&models.User{}).Where("role = ?", "teacher").Count(&userCount).Error)
	assert.Equal(t, int64(1), userCount, "不应重复建号")

	// nick 传入 → 姓名取 nick；同名再登（不同 userid）→ 用户名/昵称 _2 去重。
	withNick, err := svc.LoginWithWechatWork("ww-user-2", ptr("王老师"), nil, &school)
	require.NoError(t, err)
	require.NotNil(t, withNick.User)
	assert.Equal(t, "王老师", withNick.User.Name)
	assert.Equal(t, "王老师", withNick.User.Username)

	sameName, err := svc.LoginWithWechatWork("ww-user-3", ptr("王老师"), nil, &school)
	require.NoError(t, err)
	require.NotNil(t, sameName.User)
	assert.Equal(t, "王老师_2", sameName.User.Username)
	assert.Equal(t, "王老师_2", sameName.User.Nickname)

	// 无学校 → status=error + 文案。
	emptyDB := setupDB(t)
	emptySvc := newThirdPartyService(t, emptyDB)
	noSchool, err := emptySvc.LoginWithWechatWork("ww-user-9", nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "error", noSchool.Status)
	assert.Equal(t, "系统尚未初始化，请先联系管理员", noSchool.Message)
}

// ============================================================
// 第三方平台登录（钉钉 / 飞书 / 企微 provider 统一入口）
// ============================================================

// TestLoginWithThirdPartyBranches 四个分支：已绑定 / 手机号匹配本地账号 / 新建账号 / 无学校。
func TestLoginWithThirdPartyBranches(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)

	// ① 新建账号（姓名实名 + 默认密码 + 手机号/邮箱落库 + 绑定）。
	info := services.ThirdPartyUserInfo{
		PlatformID: "dt-1", Name: "钉钉老师", Mobile: "13800000001", Email: "dt@x.c", Avatar: "https://a/1.png",
	}
	result, err := svc.LoginWithThirdParty("dingtalk", info, &school)
	require.NoError(t, err)
	assert.Equal(t, "logged_in", result.Status)
	require.NotEmpty(t, result.Token)
	require.NotNil(t, result.User)
	assert.Equal(t, "钉钉老师", result.User.Username)
	assert.Equal(t, "13800000001", result.User.Phone)
	assert.Equal(t, "dt@x.c", result.User.Email)

	var binding models.ThirdPartyBinding
	require.NoError(t, db.Where("platform = ? AND platform_id = ?", "dingtalk", "dt-1").First(&binding).Error)
	assert.Equal(t, result.User.ID, binding.UserID)
	assert.Equal(t, "钉钉老师", binding.PlatformNick)

	// ② 已绑定 → 直接登录（无 token，用户相同）。
	again, err := svc.LoginWithThirdParty("dingtalk", info, &school)
	require.NoError(t, err)
	assert.Equal(t, "logged_in", again.Status)
	assert.Empty(t, again.Token, "已绑定分支不带 token")
	require.NotNil(t, again.User)
	assert.Equal(t, result.User.ID, again.User.ID)

	// ③ 手机号匹配本地已有账号 → 复用账号 + 补建绑定 + 返回 token。
	existing := models.User{
		SchoolID: school.ID, Role: "teacher", Username: "李老师", Name: "李老师",
		Phone: "13800000002", Status: "active",
	}
	require.NoError(t, db.Create(&existing).Error)

	reuseInfo := services.ThirdPartyUserInfo{
		PlatformID: "dt-2", Name: "李老师", Mobile: "13800000002", Avatar: "https://a/2.png",
	}
	reused, err := svc.LoginWithThirdParty("dingtalk", reuseInfo, &school)
	require.NoError(t, err)
	require.NotNil(t, reused.User)
	assert.Equal(t, existing.ID, reused.User.ID, "按手机号复用本地账号，不重复建号")
	require.NotEmpty(t, reused.Token)
	var reuseBinding models.ThirdPartyBinding
	require.NoError(t, db.Where("platform = ? AND platform_id = ?", "dingtalk", "dt-2").First(&reuseBinding).Error)
	assert.Equal(t, existing.ID, reuseBinding.UserID)
	var reloaded models.User
	require.NoError(t, db.First(&reloaded, existing.ID).Error)
	assert.Equal(t, "https://a/2.png", reloaded.AvatarPath, "本地头像为空时同步第三方头像")

	// ④ 无学校 → status=error（不建号）。
	emptyDB := setupDB(t)
	emptySvc := newThirdPartyService(t, emptyDB)
	noSchool, err := emptySvc.LoginWithThirdParty("dingtalk", info, nil)
	require.NoError(t, err)
	assert.Equal(t, "error", noSchool.Status)
	assert.Equal(t, "系统尚未初始化，请先联系管理员", noSchool.Message)
}

// ============================================================
// 扫码后绑定 / 主动绑定 / 解绑
// ============================================================

// TestBindAfterScan 错密码（status=error）/ 成功（bound）/ 上下文一次性 / nickname 与 avatar 覆盖条件。
func TestBindAfterScan(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)

	hash := mustHash(t, "ls123456")
	teacher := models.User{
		SchoolID: school.ID, Role: "teacher", Username: "绑定老师", Name: "绑定老师",
		Nickname: "绑定老师", PasswordHash: hash, Status: "active",
	}
	require.NoError(t, db.Create(&teacher).Error)

	// 错密码 → status=error（HTTP 200 由路由层断言）。
	failed, err := svc.BindAfterScan(
		"not-exist-token", "绑定老师", "wrong-password", "wechat", "wx-open-9", nil, nil, nil,
	)
	require.NoError(t, err)
	assert.Equal(t, "error", failed.Status)
	assert.Equal(t, "账号或密码错误，请核对后重试", failed.Message)
	assert.Nil(t, failed.User)

	// 先扫码拿 temp_token（nick/avatar 存入上下文）。
	scanned, err := svc.LoginWithWechat("wx-open-9", nil, ptr("微信昵称"), ptr("https://a/wx.png"))
	require.NoError(t, err)
	need, ok := scanned.(services.WechatNeedBindingResult)
	require.True(t, ok)

	// 成功绑定：上下文覆盖 platform/platform_id/nick/avatar，且临时上下文被消费。
	bound, err := svc.BindAfterScan(
		need.TempToken, "绑定老师", "ls123456", "wechat", "ignored-id", nil, nil, nil,
	)
	require.NoError(t, err)
	assert.Equal(t, "bound", bound.Status)
	require.NotNil(t, bound.User)
	assert.Equal(t, "微信昵称", bound.User.Nickname, "本地昵称是默认值（= 姓名）→ 被第三方昵称覆盖")
	assert.Equal(t, "https://a/wx.png", bound.User.AvatarPath, "本地头像为空 → 被覆盖")

	var binding models.ThirdPartyBinding
	require.NoError(t, db.Where("user_id = ? AND platform = ?", teacher.ID, "wechat").First(&binding).Error)
	assert.Equal(t, "wx-open-9", binding.PlatformID, "平台 ID 取上下文（覆盖请求体中的 ignored-id）")
	assert.Equal(t, "微信昵称", binding.PlatformNick)

	var leftover int64
	require.NoError(t, db.Model(&models.TempBindingContext{}).
		Where("temp_token = ?", need.TempToken).Count(&leftover).Error)
	assert.Equal(t, int64(0), leftover, "临时上下文一次性消费")

	// 自定义昵称 / 已有头像不被覆盖。
	custom := models.User{
		SchoolID: school.ID, Role: "teacher", Username: "自定义老师", Name: "自定义老师",
		Nickname: "我的自定义昵称", AvatarPath: "https://a/keep.png", PasswordHash: hash, Status: "active",
	}
	require.NoError(t, db.Create(&custom).Error)

	scanned2, err := svc.LoginWithQQ("qq-open-77", ptr("QQ新昵称"), ptr("https://a/qq.png"))
	require.NoError(t, err)
	need2, ok := scanned2.(services.QQNeedBindingResult)
	require.True(t, ok)

	bound2, err := svc.BindAfterScan(
		need2.TempToken, "自定义老师", "ls123456", "qq", "qq-open-77", nil, nil, nil,
	)
	require.NoError(t, err)
	assert.Equal(t, "bound", bound2.Status)
	require.NotNil(t, bound2.User)
	assert.Equal(t, "我的自定义昵称", bound2.User.Nickname, "自定义昵称不被覆盖")
	assert.Equal(t, "https://a/keep.png", bound2.User.AvatarPath, "已有头像不被覆盖")

	// 无 temp_token（已过期/不存在）时用请求体里的 platform/platform_id。
	bound3, err := svc.BindAfterScan(
		"missing-token", "绑定老师", "ls123456", "renren", "rr-1", nil, ptr("直传昵称"), nil,
	)
	require.NoError(t, err)
	assert.Equal(t, "bound", bound3.Status)
	var rrBinding models.ThirdPartyBinding
	require.NoError(t, db.Where("user_id = ? AND platform = ?", teacher.ID, "renren").First(&rrBinding).Error)
	assert.Equal(t, "rr-1", rrBinding.PlatformID)

	// 同平台再次扫码绑定 → 复用已有绑定行更新（Laravel bindThirdParty 的 update 分支）。
	bound4, err := svc.BindAfterScan(
		"missing-token", "绑定老师", "ls123456", "renren", "rr-2", ptr("union-2"), nil, nil,
	)
	require.NoError(t, err)
	assert.Equal(t, "bound", bound4.Status)
	var rrRows []models.ThirdPartyBinding
	require.NoError(t, db.Where("user_id = ? AND platform = ?", teacher.ID, "renren").Find(&rrRows).Error)
	require.Len(t, rrRows, 1, "同平台只保留一条绑定")
	assert.Equal(t, "rr-2", rrRows[0].PlatformID)
	assert.Equal(t, "union-2", rrRows[0].PlatformUnionID)

	adminHash := mustHash(t, "admin-pass")
	admin := models.User{
		SchoolID: school.ID, Role: "school_admin", Username: "绑定管理员", Name: "管理员",
		PasswordHash: adminHash, Status: "active",
	}
	require.NoError(t, db.Create(&admin).Error)
	adminResult, err := svc.BindAfterScan(
		"missing-token", "绑定管理员", "admin-pass", "wechat", "wx-admin", nil, nil, nil,
	)
	require.NoError(t, err)
	assert.Equal(t, "error", adminResult.Status, "管理员不支持第三方登录")
}

// TestBindForUserAndUnbind 主动绑定：重复 422、成功；解绑幂等。
func TestBindForUserAndUnbind(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)

	teacher := seedTeacher(t, db, school.ID, "bind-teacher")
	other := seedTeacher(t, db, school.ID, "bind-other")

	require.NoError(t, svc.BindForUser(&teacher, "wechat", "wx-1", ptr("昵称"), ptr("https://a/1.png")))
	var binding models.ThirdPartyBinding
	require.NoError(t, db.Where("user_id = ? AND platform = ?", teacher.ID, "wechat").First(&binding).Error)
	assert.Equal(t, "wx-1", binding.PlatformID)
	assert.Equal(t, "昵称", binding.PlatformNick)
	require.NotNil(t, binding.VerifiedAt)

	// 任意用户已绑定同一 platform_id → 422。
	err := svc.BindForUser(&other, "wechat", "wx-1", nil, nil)
	require.Error(t, err)
	appErr, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnprocessableEntity, appErr.Status)
	assert.Equal(t, "该第三方账号已被其他用户绑定", appErr.Message)

	// 解绑：删掉当前用户该平台绑定；再解绑不报错。
	require.NoError(t, svc.UnbindForUser(&teacher, "wechat"))
	var count int64
	require.NoError(t, db.Model(&models.ThirdPartyBinding{}).
		Where("user_id = ? AND platform = ?", teacher.ID, "wechat").Count(&count).Error)
	assert.Equal(t, int64(0), count)
	require.NoError(t, svc.UnbindForUser(&teacher, "wechat"), "无绑定时解绑也应成功")
}

// ============================================================
// 管理端通讯录
// ============================================================

// TestWechatWorkContactsService 企微通讯录：成功 / 上游异常 → 400 / 学校缺失 → 404。
func TestWechatWorkContactsService(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)

	srv, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"AT-1","expires_in":7200}`)
		case "/cgi-bin/department/list":
			jsonBody(w, `{"errcode":0,"department":[{"id":1,"name":"总部","parentid":0}]}`)
		case "/cgi-bin/user/list":
			jsonBody(w, `{"errcode":0,"userlist":[{"userid":"u1","name":"张三","department":[1]}]}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	svc.WechatWork().CorpID = "corp-1"
	svc.WechatWork().Secret = "sec-1"
	svc.WechatWork().APIBase = srv.URL
	svc.Manager().SetHTTPClient(srv.Client())

	contacts, err := svc.WechatWorkContacts(school.ID)
	require.NoError(t, err)
	require.Len(t, contacts.Members, 1)
	assert.Equal(t, "张三", contacts.Members[0].Name)
	assert.Empty(t, contacts.Platform, "企微通讯录不附带 platform 字段")

	// 学校不存在 → 404「未找到学校」。
	_, err = svc.WechatWorkContacts(school.ID + 999)
	appErr, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appErr.Status)
	assert.Equal(t, "未找到学校", appErr.Message)

	// 上游报错 → 400 + 原始文案。
	require.NoError(t, db.Model(&models.WechatWorkToken{}).Where("school_id = ?", school.ID).
		Update("expires_at", time.Now().Add(-time.Minute)).Error)
	broken, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonBody(w, `{"errcode":40013,"errmsg":"invalid corpid"}`)
	})
	svc.WechatWork().APIBase = broken.URL
	svc.Manager().SetHTTPClient(broken.Client())
	_, err = svc.WechatWorkContacts(school.ID)
	appErr, ok = services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, appErr.Status)
	assert.Equal(t, "token失败", appErr.Message)
}

// TestThirdPartyContactsService 第三方通讯录：附带 platform 标识 / 未配置平台 400。
func TestThirdPartyContactsService(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)

	// 未配置平台 → 400「未配置或未知的第三方平台」。
	_, err := svc.ThirdPartyContacts(school.ID)
	appErr, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, appErr.Status)
	assert.Equal(t, "未配置或未知的第三方平台", appErr.Message)

	// 配置为钉钉 → 成功后附带 platform=dingtalk。
	require.NoError(t, db.Model(&models.School{}).Where("id = ?", school.ID).
		Update("settings", `{"third_party_platform":"dingtalk"}`).Error)
	srv, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"dt-internal"}`)
		case "/topapi/v2/department/listsub":
			jsonBody(w, `{"errcode":0,"result":{"list":[{"dept_id":1,"name":"总部"}]}}`)
		case "/topapi/v2/user/list":
			jsonBody(w, `{"errcode":0,"result":{"list":[{"userid":"dt-u1","name":"钉钉张三","dept_id_list":[1]}]}}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	svc.Manager().SetHTTPClient(srv.Client())
	svc.Manager().DingTalkAPIBase = srv.URL
	svc.Manager().DingTalkOAuthBase = srv.URL

	contacts, err := svc.ThirdPartyContacts(school.ID)
	require.NoError(t, err)
	assert.Equal(t, "dingtalk", contacts.Platform)
	require.Len(t, contacts.Members, 1)
	assert.Equal(t, "钉钉张三", contacts.Members[0].Name)
}

// TestImportContacts 导入：教师查重（请求内 + 库内）、学生同班同名跳过、统计与文案。
func TestImportContacts(t *testing.T) {
	db := setupDB(t)
	school := seedSchool(t, db)
	svc := newThirdPartyService(t, db)
	class := seedClass(t, db, school.ID)

	// 库内已存在的教师（按手机号 / 按用户名各一位）。
	require.NoError(t, db.Create(&models.User{
		SchoolID: school.ID, Role: "teacher", Username: "已存在老师", Name: "已存在老师",
		Phone: "13800000000", Status: "active",
	}).Error)
	require.NoError(t, db.Create(&models.User{
		SchoolID: school.ID, Role: "teacher", Username: "同名老师", Name: "同名老师", Status: "active",
	}).Error)

	// 同班已有学生，用于触发「该班级已有同名学生」。
	require.NoError(t, db.Create(&models.Student{
		ClassID: class.ID, Name: "已有学生", Status: "active",
	}).Error)

	result, err := svc.ImportContacts(school.ID, services.ContactsImportRequest{
		Teachers: []services.ContactsImportTeacher{
			// 手机号带空白与连字符 → 归一化后与库内已存在手机号一致 → 跳过。
			{Name: "重复手机老师", Mobile: ptr("138 0000-0000")},
			// 同名（= username）已存在 → 跳过。
			{Name: "同名老师"},
			// 请求内同手机号重复 → 只保留第一条。
			{Name: "新老师甲", Mobile: ptr("13900000001")},
			{Name: "新老师甲重复", Mobile: ptr("13900000001")},
			// 正常创建（带邮箱）。
			{Name: "新老师乙", Phone: ptr("13900000002"), Email: ptr("new@x.c")},
		},
		Students: []services.ContactsImportStudent{
			{Name: "新同学", ClassID: intPtr(int(class.ID)), Gender: ptr("男生")},
			{Name: "已有学生", ClassID: intPtr(int(class.ID))},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 2, result.Data.CreatedTeachers)
	assert.Equal(t, 1, result.Data.CreatedStudents)
	assert.Equal(t, "已导入 2 名教师、1 名学生，跳过已存在教师 2 名，跳过同名学生 1 名", result.Message)
	require.Len(t, result.Data.SkippedTeachers, 2)
	assert.Equal(t, "手机号已存在", result.Data.SkippedTeachers[0].Reason)
	assert.Equal(t, "同名账号已存在", result.Data.SkippedTeachers[1].Reason)
	require.Len(t, result.Data.SkippedStudents, 1)
	assert.Equal(t, "该班级已有同名学生", result.Data.SkippedStudents[0].Reason)

	require.Len(t, result.Data.TeacherAccounts, 2)
	for _, account := range result.Data.TeacherAccounts {
		assert.Equal(t, services.DefaultTeacherPassword, account.InitialPassword)
		assert.NotEmpty(t, account.Username)
	}
	assert.Equal(t, "新老师甲", result.Data.TeacherAccounts[0].Name)
	assert.Equal(t, "新老师甲", result.Data.TeacherAccounts[0].Username)
	assert.Equal(t, "新老师乙", result.Data.TeacherAccounts[1].Name)

	// 落库校验：手机号已归一化、邮箱写入、学生自动分配宠物。
	var created models.User
	require.NoError(t, db.Where("username = ?", "新老师甲").First(&created).Error)
	assert.Equal(t, "13900000001", created.Phone)
	assert.Equal(t, services.DefaultTeacherPassword, created.PlainPassword)

	var createdB models.User
	require.NoError(t, db.Where("username = ?", "新老师乙").First(&createdB).Error)
	assert.Equal(t, "new@x.c", createdB.Email)

	var student models.Student
	require.NoError(t, db.Where("class_id = ? AND name = ?", class.ID, "新同学").First(&student).Error)
	assert.Equal(t, "男生", student.Gender)
	var petCount int64
	require.NoError(t, db.Model(&models.Pet{}).Where("student_id = ?", student.ID).Count(&petCount).Error)
	assert.Equal(t, int64(1), petCount, "导入学生应自动分配默认宠物")

	// 校验失败 → 422 + 点号键（同 Laravel validator 键名）。
	_, err = svc.ImportContacts(school.ID, services.ContactsImportRequest{
		Teachers: []services.ContactsImportTeacher{{Name: ""}},
		Students: []services.ContactsImportStudent{{Name: "缺班级"}},
	})
	validationErr, ok := services.AsValidationError(err)
	require.True(t, ok, "应为 ValidationError，实际 %v", err)
	assert.Contains(t, validationErr.Errors, "teachers.0.name")
	assert.Contains(t, validationErr.Errors, "students.0.class_id")

	// 学校不存在 → 404「未找到学校」。
	_, err = svc.ImportContacts(school.ID+999, services.ContactsImportRequest{})
	appErr, ok := services.AsAppError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appErr.Status)
	assert.Equal(t, "未找到学校", appErr.Message)

	// 空请求体（无 teachers / students）→ 全部统计为 0，且数组为空数组而非 null。
	empty, err := svc.ImportContacts(school.ID, services.ContactsImportRequest{})
	require.NoError(t, err)
	assert.Equal(t, "已导入 0 名教师、0 名学生", empty.Message)
	assert.NotNil(t, empty.Data.SkippedTeachers)
	assert.NotNil(t, empty.Data.TeacherAccounts)
}
