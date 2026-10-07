// 企业微信（企微）调用服务：access_token 缓存 + OAuth code 换 userid + 部门 / 成员通讯录。
//
// 忠实移植自 Laravel App\Services\WechatWorkService 的
// getAccessToken / getUserIdByCode / getDepartmentList / getDepartmentUsers / fetchContacts；
// 配置项名与 Laravel config/wechat-work.php 完全一致：
//
//	WECHAT_WORK_CORPID / WECHAT_WORK_AGENTID / WECHAT_WORK_SECRET /
//	WECHAT_WORK_TOKEN / WECHAT_WORK_ENCODING_AES_KEY
//
// 另加 WECHAT_WORK_API_BASE（见下方说明）。
//
// 与 Laravel 的有意差异：
//  1. **access_token 缓存落表**：Laravel 用 `Cache::put("wecom_at:<schoolId>", $token, $ttl)`
//     （默认 file 缓存）；Go 端无 Cache 层，改用表 `wechat_work_tokens(school_id 主键, token,
//     expires_at)` —— 多实例共享同一 token（比进程内全局缓存更正确：不会各实例各请求一次
//     gettoken 触发企微限频）。TTL 口径不变：`max(expires_in - 300, 60)` 秒。
//  2. **接口地址可注入**：Laravel 硬编码 `https://qyapi.weixin.qq.com`；Go 端默认即官方地址，
//     但支持环境变量 `WECHAT_WORK_API_BASE` 覆盖 —— **仅为测试 / 自建代理留的注入点，
//     生产部署恒为官方地址**（不配置即官方地址）。
//  3. **HTTP 客户端可注入**（`SetHTTPClient`，同 services.AIService 的模式）：测试用
//     `httptest` 假上游，绝不访问真实外网。
//  4. `!isset($access_token)` 的边界：Laravel 对 `access_token: null` 会走到 PHP TypeError；
//     Go 端把「缺失 / null / 空串」统一视为失败并抛「token失败」。
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
	"github.com/RealKiro/learnstar-planet/backend-go/internal/util"
	"gorm.io/gorm"
)

// wechatWorkAPITimeout 单次请求超时（同 Laravel `Http::timeout(10)`）。
const wechatWorkAPITimeout = 10 * time.Second

// DefaultWechatWorkAPIBase 企微接口官方根地址（生产恒为此地址）。
const DefaultWechatWorkAPIBase = "https://qyapi.weixin.qq.com"

// WechatWorkService 企业微信调用服务。
type WechatWorkService struct {
	db     *gorm.DB
	client HTTPDoer

	// CorpID / AgentID / Secret / Token / EncodingAESKey 对应 Laravel config/wechat-work.php 的
	// 同名字段（env: WECHAT_WORK_CORPID / _AGENTID / _SECRET / _TOKEN / _ENCODING_AES_KEY）。
	// 导出字段便于测试直接覆盖（同 AdminOps.UploadRoot 的做法）。
	CorpID         string
	AgentID        int
	Secret         string
	Token          string
	EncodingAESKey string
	// APIBase 企微接口根地址：env WECHAT_WORK_API_BASE，默认官方地址（测试 / 自建代理注入点）。
	APIBase string
}

// NewWechatWorkService 创建企微服务（构造时读环境变量，同 AdminOps 读 UPLOAD_DIR 的做法）。
func NewWechatWorkService(db *gorm.DB) *WechatWorkService {
	agentID, _ := strconv.Atoi(envOrDefault("WECHAT_WORK_AGENTID", "0"))
	return &WechatWorkService{
		db:             db,
		client:         &http.Client{},
		CorpID:         envOrDefault("WECHAT_WORK_CORPID", ""),
		AgentID:        agentID,
		Secret:         envOrDefault("WECHAT_WORK_SECRET", ""),
		Token:          envOrDefault("WECHAT_WORK_TOKEN", ""),
		EncodingAESKey: envOrDefault("WECHAT_WORK_ENCODING_AES_KEY", ""),
		APIBase:        strings.TrimRight(envOrDefault("WECHAT_WORK_API_BASE", DefaultWechatWorkAPIBase), "/"),
	}
}

// SetHTTPClient 注入 HTTP 客户端（测试用 httptest 假上游）。
func (s *WechatWorkService) SetHTTPClient(c HTTPDoer) {
	if c != nil {
		s.client = c
	}
}

// GetAccessToken 取企微 access_token（命中未过期的表缓存则直接返回）。
func (s *WechatWorkService) GetAccessToken(schoolID uint) (string, error) {
	var cached models.WechatWorkToken
	err := s.db.First(&cached, schoolID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if err == nil && cached.Token != "" && cached.ExpiresAt.After(util.Now()) {
		return cached.Token, nil
	}

	if s.CorpID == "" || s.Secret == "" {
		return "", errors.New("企微未配置")
	}

	r, reqErr := s.getJSON("/cgi-bin/gettoken", url.Values{
		"corpid":     {s.CorpID},
		"corpsecret": {s.Secret},
	})
	if reqErr != nil {
		return "", reqErr
	}
	token := jsonString(r["access_token"])
	if token == "" {
		// Laravel 在此 Log::error('企微token失败', $r) 后抛 RuntimeException。
		return "", errors.New("token失败")
	}

	ttl := jsonIntDefault(r["expires_in"], 7200) - 300
	if ttl < 60 {
		ttl = 60
	}
	expiresAt := util.Now().Add(time.Duration(ttl) * time.Second)

	if delErr := s.db.Where("school_id = ?", schoolID).Delete(&models.WechatWorkToken{}).Error; delErr != nil {
		return "", delErr
	}
	if createErr := s.db.Create(&models.WechatWorkToken{
		SchoolID:  schoolID,
		Token:     token,
		ExpiresAt: expiresAt,
	}).Error; createErr != nil {
		return "", createErr
	}
	return token, nil
}

// ============================================================
// 审批（请假）查询：getapprovalinfo 列表 + getapprovaldetail 详情
// 移植自 Laravel WechatWorkService::getLeaveApprovalSpNos / getApprovalDetail / parse
// ============================================================

// ApprovalDetail 企微审批详情（等价 Laravel `parse()` 返回的关联数组）。
// 键名与 HandleWebhookCallback / SyncForSchool 用到的字段逐一对齐
// （sp_no / submitter_userid / student_name / leave_start / leave_end / leave_type /
// reason / approve_status / approved_at），可直接 marshal 成 raw_data 存库。
// `student_name` / `leave_type` / `reason` / `approved_at` 可为 null（同 Laravel 的可空语义）。
type ApprovalDetail struct {
	SpNo            string  `json:"sp_no"`
	SubmitterUserID string  `json:"submitter_userid"`
	StudentName     *string `json:"student_name"`
	LeaveType       *string `json:"leave_type"`
	// LeaveStart / LeaveEnd 请假起止（Unix 秒；缺失为 0，调用方按 Laravel 的 `?:` 语义回退）。
	LeaveStart int64   `json:"leave_start"`
	LeaveEnd   int64   `json:"leave_end"`
	Reason     *string `json:"reason"`
	// ApproveStatus 审批状态（同 Laravel `sp_status`，缺省 1）。
	ApproveStatus int `json:"approve_status"`
	// ApprovedAt 通过时间（Unix 秒，可为 null）。
	ApprovedAt *int64 `json:"approved_at"`
}

// GetLeaveApprovalSpNos 拉取指定时间段内「已通过」（sp_status=2）的请假审批单号。
//
// 口径逐条对齐 Laravel：分页游标从 0 开始，`errcode != 0` 立即 break（**不抛错**，返回已收集的
// 单号），每页取 `sp_no_list` 中的非空字符串，`next_cursor > 0` 时续页。
func (s *WechatWorkService) GetLeaveApprovalSpNos(schoolID uint, start, end int64) ([]string, error) {
	token, err := s.GetAccessToken(schoolID)
	if err != nil {
		return nil, err
	}

	spNos := []string{}
	cursor := 0
	for {
		r, err := s.postJSON("/cgi-bin/oa/getapprovalinfo", url.Values{"access_token": {token}}, map[string]any{
			"starttime": start,
			"endtime":   end,
			"cursor":    cursor,
			"size":      100,
			"filters":   map[string]any{"key": "sp_status", "value": "2"},
		})
		if err != nil {
			return nil, err
		}
		if jsonErrcode(r) != 0 {
			break
		}
		for _, item := range jsonList(r["sp_no_list"]) {
			if no, ok := item.(string); ok && no != "" {
				spNos = append(spNos, no)
			}
		}
		cursor = jsonIntDefault(r["next_cursor"], 0)
		if cursor <= 0 {
			break
		}
	}
	return spNos, nil
}

// GetApprovalDetail 取单张审批单详情并解析为 ApprovalDetail。
//
// 返回 (nil, nil) 表示「上游明确没有该单」：`errcode != 0` 或 `info` 缺失/非对象/为空
// （同 Laravel 返回 null）。access_token 获取失败、传输层异常向上抛错（Laravel 亦不捕获）。
func (s *WechatWorkService) GetApprovalDetail(schoolID uint, spNo string) (*ApprovalDetail, error) {
	token, err := s.GetAccessToken(schoolID)
	if err != nil {
		return nil, err
	}
	r, err := s.postJSON("/cgi-bin/oa/getapprovaldetail", url.Values{"access_token": {token}}, map[string]any{"sp_no": spNo})
	if err != nil {
		return nil, err
	}
	if jsonErrcode(r) != 0 {
		return nil, nil
	}
	info, ok := r["info"].(map[string]any)
	if !ok || len(info) == 0 {
		return nil, nil
	}
	return parseApprovalInfo(info), nil
}

// parseApprovalInfo 逐条对齐 Laravel `WechatWorkService::parse($info)`：
//   - sp_no / applyer.userid / sp_status（缺省 1）；
//   - approved_at：sp_record 中**严格等于 2**（`=== 2`，字符串 "2" 不算）且存在 approve_time 的
//     最后一条的时间戳；
//   - apply_data.contents 逐条按 title 关键词归类（顺序同 Laravel 的 if/elseif 链，
//     后匹配的会覆盖先前值）：学生姓名 / 请假类型 / 请假事由 / 请假时间（new_begin|new_start + new_end）。
func parseApprovalInfo(info map[string]any) *ApprovalDetail {
	d := &ApprovalDetail{
		SpNo:          jsonString(info["sp_no"]),
		ApprovedAt:    nil,
		ApproveStatus: jsonIntDefault(info["sp_status"], 1),
	}
	if applyer, ok := info["applyer"].(map[string]any); ok {
		d.SubmitterUserID = jsonString(applyer["userid"])
	}

	for _, item := range jsonList(info["sp_record"]) {
		rec, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if !jsonStrictIntEquals(rec["sp_status"], 2) {
			continue
		}
		if raw, ok := rec["approve_time"]; ok && raw != nil {
			at := int64(jsonIntDefault(raw, 0))
			d.ApprovedAt = &at
		}
	}

	applyData, _ := info["apply_data"].(map[string]any)
	for _, item := range jsonList(applyData["contents"]) {
		c, ok := item.(map[string]any)
		if !ok {
			continue
		}
		title := jsonString(c["title"])
		value := c["value"]
		text := approvalContentText(value)

		switch {
		case matchAny(title, []string{"学生姓名", "学生", "姓名"}):
			if text != "" {
				v := text
				d.StudentName = &v
			} else {
				d.StudentName = approvalControlsText(value)
			}
		case matchAny(title, []string{"请假类型", "类型"}):
			v := text
			d.LeaveType = &v
		case matchAny(title, []string{"请假事由", "事由", "原因", "说明"}):
			v := text
			d.Reason = &v
		case matchAny(title, []string{"请假时间", "起止时间", "开始时间"}):
			if m, ok := value.(map[string]any); ok {
				if raw, ok := m["new_begin"]; ok {
					d.LeaveStart = int64(jsonIntDefault(raw, 0))
				} else {
					d.LeaveStart = int64(jsonIntDefault(m["new_start"], 0))
				}
				d.LeaveEnd = int64(jsonIntDefault(m["new_end"], 0))
			}
		}
	}

	return d
}

// approvalContentText 等价 Laravel parse 里的 `$vt`：
// `is_array($v) ? (is_string($v['text']) ? $v['text'] : ”) : (string) $v`。
func approvalContentText(value any) string {
	if m, ok := value.(map[string]any); ok {
		return jsonString(m["text"])
	}
	return phpScalarString(value)
}

// approvalControlsText 等价 Laravel parse 的姓名兜底
// `is_array($v) && isset($v['controls'][0]['text']) ? $v['controls'][0]['text'] : null`；
// `$vt ?: 兜底` 的 PHP 假值判断在调用方（text 为空串时走兜底）。
func approvalControlsText(value any) *string {
	m, ok := value.(map[string]any)
	if !ok {
		return nil
	}
	list := jsonList(m["controls"])
	if len(list) == 0 {
		return nil
	}
	first, ok := list[0].(map[string]any)
	if !ok {
		return nil
	}
	raw, ok := first["text"]
	if !ok || raw == nil {
		return nil
	}
	v := phpScalarString(raw)
	return &v
}

// matchAny 等价 Laravel `mf($title, $keywords)`（mb_strpos 子串包含）。
func matchAny(title string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(title, kw) {
			return true
		}
	}
	return false
}

// jsonStrictIntEquals 等价 PHP `$v === 2`：仅接受 JSON 数值（float64/int），
// 字符串 "2" 不算（同 Laravel `===` 严格比较）。
func jsonStrictIntEquals(v any, want int) bool {
	switch n := v.(type) {
	case float64:
		return n == float64(want)
	case int:
		return n == want
	case int64:
		return n == int64(want)
	}
	return false
}

// phpScalarString 等价 PHP `(string) $scalar` 的常见分支（null → ""；true → "1"；false → ""）。
func phpScalarString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		if t {
			return "1"
		}
		return ""
	case float64:
		if t == math.Trunc(t) && math.Abs(t) < 1e15 {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'g', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	}
	return ""
}

// GetUserIDByCode 用 OAuth code 换取企微 UserId。
//
// 口径同 Laravel getUserIdByCode：errcode 非 0 时记 warning 并返回空串（不抛错）；
// 取 `UserId`，其次 `userid`。access_token 获取失败（企微未配置 / token失败 / 网络异常）会向上抛错
// —— Laravel 同样不捕获（控制器里表现为 500）。
func (s *WechatWorkService) GetUserIDByCode(schoolID uint, code string) (string, error) {
	token, err := s.GetAccessToken(schoolID)
	if err != nil {
		return "", err
	}
	r, err := s.getJSON("/cgi-bin/user/getuserinfo", url.Values{
		"access_token": {token},
		"code":         {code},
	})
	if err != nil {
		return "", err
	}
	if jsonErrcode(r) != 0 {
		return "", nil
	}
	if v := jsonString(r["UserId"]); v != "" {
		return v, nil
	}
	return jsonString(r["userid"]), nil
}

// GetDepartmentList 企微部门列表（原始数组，同 Laravel `$r['department'] ?? []`）。
func (s *WechatWorkService) GetDepartmentList(schoolID uint) ([]map[string]any, error) {
	token, err := s.GetAccessToken(schoolID)
	if err != nil {
		return nil, err
	}
	r, err := s.getJSON("/cgi-bin/department/list", url.Values{"access_token": {token}})
	if err != nil {
		return nil, err
	}
	if jsonErrcode(r) != 0 {
		return nil, fmt.Errorf("获取企微部门失败：%s", jsonString(r["errmsg"]))
	}
	return jsonObjectList(r["department"]), nil
}

// GetDepartmentUsers 部门成员（`fetch_child=1` 会带上所有子部门成员）。
func (s *WechatWorkService) GetDepartmentUsers(schoolID uint, departmentID int, fetchChild bool) ([]map[string]any, error) {
	token, err := s.GetAccessToken(schoolID)
	if err != nil {
		return nil, err
	}
	child := "0"
	if fetchChild {
		child = "1"
	}
	r, err := s.getJSON("/cgi-bin/user/list", url.Values{
		"access_token":  {token},
		"department_id": {strconv.Itoa(departmentID)},
		"fetch_child":   {child},
	})
	if err != nil {
		return nil, err
	}
	if jsonErrcode(r) != 0 {
		return nil, fmt.Errorf("获取企微成员失败：%s", jsonString(r["errmsg"]))
	}
	return jsonObjectList(r["userlist"]), nil
}

// FetchContacts 拉取全部通讯录（部门 + 成员），按 userid 去重。
//
// 口径逐条对齐 Laravel fetchContacts：**只遍历顶级部门**（`parentid == 0`），
// 每个顶级部门用 `fetch_child=1` 一次性带出其所有子部门成员；
// 成员的 `department`（部门 ID 列表）经部门表映射为 `department_names`。
func (s *WechatWorkService) FetchContacts(schoolID uint) (*ThirdPartyContacts, error) {
	departments, err := s.GetDepartmentList(schoolID)
	if err != nil {
		return nil, err
	}

	nameByID := map[int]string{}
	for _, d := range departments {
		nameByID[jsonIntDefault(d["id"], 0)] = jsonString(d["name"])
	}

	members := make([]ThirdPartyMember, 0, len(departments))
	seen := map[string]bool{}
	for _, d := range departments {
		if jsonIntDefault(d["parentid"], -1) != 0 {
			continue
		}
		deptID := jsonIntDefault(d["id"], 0)
		if deptID <= 0 {
			continue
		}
		users, err := s.GetDepartmentUsers(schoolID, deptID, true)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			uid := jsonString(u["userid"])
			if uid == "" || seen[uid] {
				continue
			}
			seen[uid] = true

			deptNames := make([]string, 0, 1)
			for _, raw := range jsonList(u["department"]) {
				if name, ok := nameByID[jsonIntDefault(raw, 0)]; ok {
					deptNames = append(deptNames, name)
				}
			}

			members = append(members, ThirdPartyMember{
				UserID:          uid,
				Name:            jsonString(u["name"]),
				Mobile:          jsonString(u["mobile"]),
				Email:           jsonString(u["email"]),
				Position:        jsonString(u["position"]),
				DepartmentNames: deptNames,
			})
		}
	}

	return &ThirdPartyContacts{Departments: departments, Members: members}, nil
}

// getJSON 发 GET 请求并把响应体解析为键值表。
// 解析失败时返回空表（Laravel `->json()` 解析失败返回 null，后续 `?? ` 取值等价于「没有该键」）。
func (s *WechatWorkService) getJSON(path string, query url.Values) (map[string]any, error) {
	requestURL := s.APIBase + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	ctx, cancel := context.WithTimeout(context.Background(), wechatWorkAPITimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := json.Unmarshal(body, &out); err != nil {
		return map[string]any{}, nil
	}
	return out, nil
}

// postJSON 发 POST 请求（JSON body）并把响应体解析为键值表。
// 解析失败时返回空表（同 getJSON；Laravel `->json()` 解析失败返回 null，后续 `?? ` 取值等价于「没有该键」）。
func (s *WechatWorkService) postJSON(path string, query url.Values, payload any) (map[string]any, error) {
	requestURL := s.APIBase + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wechatWorkAPITimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := json.Unmarshal(body, &out); err != nil {
		return map[string]any{}, nil
	}
	return out, nil
}

// ============================================================
// JSON 取值小工具（PHP 数组访问语义）
// ============================================================

// jsonErrcode 读取 errcode：缺省或非数值返回 -1（同 Laravel `($r['errcode'] ?? -1) !== 0`）。
func jsonErrcode(m map[string]any) int {
	v, ok := m["errcode"]
	if !ok {
		return -1
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	}
	return -1
}

// jsonString 取字符串值，非字符串返回空串。
func jsonString(v any) string {
	s, _ := v.(string)
	return s
}

// jsonIntDefault 取整数值，缺失 / 非数值返回 def。
func jsonIntDefault(v any, def int) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i
		}
	}
	return def
}

// jsonList 把值转为任意元素数组，非数组返回 nil。
func jsonList(v any) []any {
	list, _ := v.([]any)
	return list
}

// jsonObjectList 把值转为对象数组，元素非对象时以空对象占位（保留数组长度）。
func jsonObjectList(v any) []map[string]any {
	list := jsonList(v)
	if list == nil {
		return nil
	}
	out := make([]map[string]any, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
			continue
		}
		out = append(out, map[string]any{})
	}
	return out
}
