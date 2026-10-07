// 第三方办公平台适配层（每校配置一个平台）：provider 工厂 + 企业微信 / 钉钉 / 飞书三个 provider。
//
// 忠实移植自 Laravel App\Services\ThirdParty\ThirdPartyManager、ThirdPartyProvider 与
// WeChatWorkProvider / DingTalkProvider / LarkProvider：接口地址、参数名、取值路径与错误文案逐条对齐。
//
// 与 Laravel 的有意差异：
//  1. 钉钉 / 飞书的接口地址在 Laravel 里是硬编码常量；Go 端默认同样硬编码官方地址，但把
//     `DingTalkOAuthBase/DingTalkAPIBase/LarkAccountsBase/LarkAPIBase` 导出为**测试用注入点**
//     （同 AI 官方账单批次保留 `SetEndpointOverrides` 的做法）——生产不配置即官方地址，
//     测试全部指向 httptest 假上游，不访问真实外网。
//  2. 应用凭证从环境变量读取（同 Laravel config/dingtalk.php、config/feishu.php 的键名：
//     DINGTALK_APP_KEY / DINGTALK_APP_SECRET / FEISHU_APP_ID / FEISHU_APP_SECRET）。
//  3. HTTP 客户端可注入（`SetHTTPClient`）；Laravel 用 `Http::timeout(10)`，Go 端用 10 秒超时
//     （钉钉 / 飞书 provider 沿用与企微一致的超时口径）。
//  4. providerFor 对未配置 / 未知平台抛的错误文案与 Laravel 逐字一致（`未配置或未知的第三方平台`）。
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ============================================================
// 共享类型
// ============================================================

// ThirdPartyMember 通讯录成员（键名逐字同 Laravel provider 返回的 members 元素）。
type ThirdPartyMember struct {
	UserID          string   `json:"userid"`
	Name            string   `json:"name"`
	Mobile          string   `json:"mobile"`
	Email           string   `json:"email"`
	Position        string   `json:"position"`
	DepartmentNames []string `json:"department_names"`
}

// ThirdPartyContacts 通讯录（部门原始数组 + 归一化成员）。
// Platform 仅由 `GET /admin/third-party/contacts` 追加（同 Laravel `$contacts['platform'] = $provider->key()`）。
type ThirdPartyContacts struct {
	Departments []map[string]any   `json:"departments"`
	Members     []ThirdPartyMember `json:"members"`
	Platform    string             `json:"platform,omitempty"`
}

// ThirdPartyUserInfo provider 用授权 code 换回的用户身份
// （键名同 PHP 数组 platform_id / name / mobile / email / avatar）。
type ThirdPartyUserInfo struct {
	PlatformID string
	Name       string
	Mobile     string
	Email      string
	Avatar     string
}

// ThirdPartyProvider 第三方办公平台适配接口（同 Laravel ThirdPartyProvider）。
type ThirdPartyProvider interface {
	// Key 平台标识：wechat_work / dingtalk / feishu。
	Key() string
	// AuthURL 扫码授权 URL（前端生成二维码跳转）。
	AuthURL(schoolID uint, redirectURI string) string
	// GetUserByCode 用授权 code 换取用户身份。
	GetUserByCode(schoolID uint, code string) (ThirdPartyUserInfo, error)
	// FetchContacts 拉取通讯录（部门 + 成员，按 userid 去重）。
	FetchContacts(schoolID uint) (*ThirdPartyContacts, error)
}

// ============================================================
// 管理器
// ============================================================

// 官方端点（生产恒为此值；导出字段仅测试 / 自建代理可覆盖）。
const (
	DefaultDingTalkOAuthBase = "https://api.dingtalk.com"
	DefaultDingTalkAPIBase   = "https://oapi.dingtalk.com"
	DefaultLarkAccountsBase  = "https://accounts.feishu.cn"
	DefaultLarkAPIBase       = "https://open.feishu.cn"
)

// ThirdPartyManager 按学校配置的平台返回对应 provider（每校一平台）。
type ThirdPartyManager struct {
	wechatWork *WechatWorkService
	httpClient HTTPDoer

	// 应用凭证（env 同 Laravel config/dingtalk.php 与 config/feishu.php）。
	DingTalkAppKey    string
	DingTalkAppSecret string
	LarkAppID         string
	LarkAppSecret     string

	// 端点（默认官方地址；**仅测试 / 自建代理注入点**，生产勿改）。
	DingTalkOAuthBase string
	DingTalkAPIBase   string
	LarkAccountsBase  string
	LarkAPIBase       string
}

// NewThirdPartyManager 创建管理器（构造时读环境变量，同 AdminOps 读 UPLOAD_DIR 的做法）。
func NewThirdPartyManager(wechatWork *WechatWorkService) *ThirdPartyManager {
	return &ThirdPartyManager{
		wechatWork:        wechatWork,
		httpClient:        &http.Client{},
		DingTalkAppKey:    envOrDefault("DINGTALK_APP_KEY", ""),
		DingTalkAppSecret: envOrDefault("DINGTALK_APP_SECRET", ""),
		LarkAppID:         envOrDefault("FEISHU_APP_ID", ""),
		LarkAppSecret:     envOrDefault("FEISHU_APP_SECRET", ""),
		DingTalkOAuthBase: DefaultDingTalkOAuthBase,
		DingTalkAPIBase:   DefaultDingTalkAPIBase,
		LarkAccountsBase:  DefaultLarkAccountsBase,
		LarkAPIBase:       DefaultLarkAPIBase,
	}
}

// SetHTTPClient 注入 HTTP 客户端（同时注入企微服务；测试用 httptest 假上游）。
func (m *ThirdPartyManager) SetHTTPClient(c HTTPDoer) {
	if c == nil {
		return
	}
	m.httpClient = c
	if m.wechatWork != nil {
		m.wechatWork.SetHTTPClient(c)
	}
}

// WechatWork 返回底层企微服务（供 OAuth code 换 userid 等复用）。
func (m *ThirdPartyManager) WechatWork() *WechatWorkService { return m.wechatWork }

// ProviderFor 按平台标识返回 provider；未配置（空串）或未知平台抛错（文案同 Laravel）。
func (m *ThirdPartyManager) ProviderFor(platform string) (ThirdPartyProvider, error) {
	switch platform {
	case "wechat_work":
		return &wechatWorkProvider{manager: m}, nil
	case "dingtalk":
		return &dingTalkProvider{manager: m}, nil
	case "feishu":
		return &larkProvider{manager: m}, nil
	default:
		return nil, errors.New("未配置或未知的第三方平台")
	}
}

// ============================================================
// HTTP 小工具（钉钉 / 飞书 provider 共用）
// ============================================================

// client 返回实际使用的 HTTP 客户端。
func (m *ThirdPartyManager) client() HTTPDoer {
	if m.httpClient != nil {
		return m.httpClient
	}
	return http.DefaultClient
}

// JSON 请求（body 为 nil 时不带请求体；headers 追加请求头），响应解析为键值表。
// 解析失败返回空表（Laravel `->json()` 解析失败返回 null，后续 `?? ` 取值等价于「没有该键」）。
func (m *ThirdPartyManager) requestJSON(method, endpoint string, body any, headers map[string]string) (map[string]any, error) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(payload)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := m.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{}, nil
	}
	return out, nil
}

// jsonNestedMap 取嵌套对象（缺失 / 非对象返回空表）。
func jsonNestedMap(m map[string]any, key string) map[string]any {
	if v, ok := m[key].(map[string]any); ok {
		return v
	}
	return map[string]any{}
}

// ============================================================
// 企业微信 provider
// ============================================================

type wechatWorkProvider struct{ manager *ThirdPartyManager }

func (p *wechatWorkProvider) Key() string { return "wechat_work" }

// AuthURL 企微扫码授权 URL（拼接口径逐字同 Laravel WeChatWorkProvider::authUrl）。
func (p *wechatWorkProvider) AuthURL(schoolID uint, redirectURI string) string {
	corpID, agentID := "", 0
	if p.manager != nil && p.manager.wechatWork != nil {
		corpID = p.manager.wechatWork.CorpID
		agentID = p.manager.wechatWork.AgentID
	}
	return "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=" + url.QueryEscape(corpID) +
		"&agentid=" + strconv.Itoa(agentID) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&state=" + strconv.FormatUint(uint64(schoolID), 10)
}

func (p *wechatWorkProvider) GetUserByCode(schoolID uint, code string) (ThirdPartyUserInfo, error) {
	if p.manager == nil || p.manager.wechatWork == nil {
		return ThirdPartyUserInfo{}, errors.New("企业微信登录失败")
	}
	userid, err := p.manager.wechatWork.GetUserIDByCode(schoolID, code)
	if err != nil {
		return ThirdPartyUserInfo{}, err
	}
	if userid == "" {
		return ThirdPartyUserInfo{}, errors.New("企业微信登录失败")
	}
	detail := p.userDetail(schoolID, userid)
	name := detail["name"]
	if name == "" {
		name = userid
	}
	return ThirdPartyUserInfo{
		PlatformID: userid,
		Name:       name,
		Mobile:     detail["mobile"],
		Email:      detail["email"],
		Avatar:     detail["avatar"],
	}, nil
}

func (p *wechatWorkProvider) FetchContacts(schoolID uint) (*ThirdPartyContacts, error) {
	if p.manager == nil || p.manager.wechatWork == nil {
		return nil, errors.New("企微未配置")
	}
	return p.manager.wechatWork.FetchContacts(schoolID)
}

// userDetail 读企微成员详情（errcode 非 0 时返回空详情，同 Laravel）。
func (p *wechatWorkProvider) userDetail(schoolID uint, userid string) map[string]string {
	empty := map[string]string{"name": "", "mobile": "", "email": "", "avatar": ""}
	token, err := p.manager.wechatWork.GetAccessToken(schoolID)
	if err != nil {
		return empty
	}
	apiBase := p.manager.wechatWork.APIBase
	r, err := p.manager.requestJSON(http.MethodGet,
		apiBase+"/cgi-bin/user/get?access_token="+url.QueryEscape(token)+"&userid="+url.QueryEscape(userid),
		nil, nil)
	if err != nil || jsonErrcode(r) != 0 {
		return empty
	}
	return map[string]string{
		"name":   jsonString(r["name"]),
		"mobile": jsonString(r["mobile"]),
		"email":  jsonString(r["email"]),
		"avatar": jsonString(r["avatar"]),
	}
}

// ============================================================
// 钉钉 provider
// ============================================================

type dingTalkProvider struct{ manager *ThirdPartyManager }

func (p *dingTalkProvider) Key() string { return "dingtalk" }

// AuthURL 钉钉扫码授权 URL（拼接口径逐字同 Laravel DingTalkProvider::authUrl）。
func (p *dingTalkProvider) AuthURL(schoolID uint, redirectURI string) string {
	return "https://login.dingtalk.com/oauth2/auth?redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code&client_id=" + url.QueryEscape(p.manager.DingTalkAppKey) +
		"&scope=openid&state=" + strconv.FormatUint(uint64(schoolID), 10)
}

func (p *dingTalkProvider) GetUserByCode(schoolID uint, code string) (ThirdPartyUserInfo, error) {
	r, err := p.manager.requestJSON(http.MethodPost, p.manager.DingTalkOAuthBase+"/v1.0/oauth2/userAccessToken",
		map[string]any{
			"clientId":     p.manager.DingTalkAppKey,
			"clientSecret": p.manager.DingTalkAppSecret,
			"code":         code,
			"grantType":    "authorization_code",
		}, nil)
	if err != nil {
		return ThirdPartyUserInfo{}, err
	}
	token := jsonString(r["accessToken"])
	if token == "" {
		return ThirdPartyUserInfo{}, fmt.Errorf("钉钉登录失败：%s", jsonString(r["message"]))
	}

	me, err := p.manager.requestJSON(http.MethodGet, p.manager.DingTalkOAuthBase+"/v1.0/contact/users/me",
		nil, map[string]string{"x-acs-dingtalk-access-token": token})
	if err != nil {
		return ThirdPartyUserInfo{}, err
	}

	platformID := jsonString(me["unionId"])
	if platformID == "" {
		platformID = jsonString(me["userId"])
	}
	if platformID == "" {
		platformID = code
	}
	name := jsonString(me["nick"])
	if name == "" {
		name = jsonString(me["name"])
	}
	return ThirdPartyUserInfo{
		PlatformID: platformID,
		Name:       name,
		Mobile:     jsonString(me["mobile"]),
		Email:      jsonString(me["email"]),
		Avatar:     jsonString(me["avatarUrl"]),
	}, nil
}

func (p *dingTalkProvider) FetchContacts(schoolID uint) (*ThirdPartyContacts, error) {
	token := p.internalToken()
	if token == "" {
		return nil, errors.New("钉钉未配置或获取企业内部 token 失败")
	}

	departments := p.allDepartments(token)
	nameByID := map[int]string{}
	for _, d := range departments {
		nameByID[jsonIntDefault(d["dept_id"], 0)] = jsonString(d["name"])
	}

	members := make([]ThirdPartyMember, 0, len(departments))
	seen := map[string]bool{}
	for _, d := range departments {
		deptID := jsonIntDefault(d["dept_id"], 0)
		if deptID <= 0 {
			continue
		}
		for _, u := range p.deptUsers(token, deptID) {
			uid := jsonString(u["userid"])
			if uid == "" || seen[uid] {
				continue
			}
			seen[uid] = true

			deptNames := make([]string, 0, 1)
			for _, raw := range jsonList(u["dept_id_list"]) {
				if name, ok := nameByID[jsonIntDefault(raw, 0)]; ok {
					deptNames = append(deptNames, name)
				}
			}
			members = append(members, ThirdPartyMember{
				UserID:          uid,
				Name:            jsonString(u["name"]),
				Mobile:          jsonString(u["mobile"]),
				Email:           jsonString(u["email"]),
				Position:        jsonString(u["title"]),
				DepartmentNames: deptNames,
			})
		}
	}

	return &ThirdPartyContacts{Departments: departments, Members: members}, nil
}

// internalToken 企业内部应用 access_token（缺失返回空串，同 Laravel）。
func (p *dingTalkProvider) internalToken() string {
	r, err := p.manager.requestJSON(http.MethodGet, p.manager.DingTalkAPIBase+"/gettoken", nil, nil)
	if err != nil {
		// Laravel 未捕获传输层异常；此处按「拿不到 token」处理，由调用方抛业务错误文案。
		return ""
	}
	return jsonString(r["access_token"])
}

func (p *dingTalkProvider) allDepartments(token string) []map[string]any {
	r, err := p.manager.requestJSON(http.MethodPost,
		p.manager.DingTalkAPIBase+"/topapi/v2/department/listsub?access_token="+token,
		map[string]any{"dept_id": 1}, nil)
	if err != nil {
		return nil
	}
	return jsonObjectList(jsonNestedMap(r, "result")["list"])
}

func (p *dingTalkProvider) deptUsers(token string, deptID int) []map[string]any {
	r, err := p.manager.requestJSON(http.MethodPost,
		p.manager.DingTalkAPIBase+"/topapi/v2/user/list?access_token="+token,
		map[string]any{"dept_id": deptID, "cursor": 0, "size": 100}, nil)
	if err != nil {
		return nil
	}
	return jsonObjectList(jsonNestedMap(r, "result")["list"])
}

// ============================================================
// 飞书（Lark）provider
// ============================================================

type larkProvider struct{ manager *ThirdPartyManager }

func (p *larkProvider) Key() string { return "feishu" }

// AuthURL 飞书扫码授权 URL（拼接口径逐字同 Laravel LarkProvider::authUrl）。
func (p *larkProvider) AuthURL(schoolID uint, redirectURI string) string {
	return "https://accounts.feishu.cn/open-apis/authen/v1/authorize?client_id=" + url.QueryEscape(p.manager.LarkAppID) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code&state=" + strconv.FormatUint(uint64(schoolID), 10)
}

func (p *larkProvider) GetUserByCode(schoolID uint, code string) (ThirdPartyUserInfo, error) {
	r, err := p.manager.requestJSON(http.MethodPost, p.manager.LarkAccountsBase+"/oauth/v3/token",
		map[string]any{
			"grant_type":    "authorization_code",
			"client_id":     p.manager.LarkAppID,
			"client_secret": p.manager.LarkAppSecret,
			"code":          code,
		}, nil)
	if err != nil {
		return ThirdPartyUserInfo{}, err
	}
	token := jsonString(r["access_token"])
	if token == "" {
		return ThirdPartyUserInfo{}, fmt.Errorf("飞书登录失败：%s", jsonString(r["error_description"]))
	}

	me, err := p.manager.requestJSON(http.MethodGet, p.manager.LarkAPIBase+"/open-apis/authen/v1/user_info",
		nil, map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		return ThirdPartyUserInfo{}, err
	}
	data := jsonNestedMap(me, "data")

	platformID := jsonString(data["union_id"])
	if platformID == "" {
		platformID = jsonString(data["open_id"])
	}
	if platformID == "" {
		platformID = code
	}
	return ThirdPartyUserInfo{
		PlatformID: platformID,
		Name:       jsonString(data["name"]),
		Mobile:     jsonString(data["mobile"]),
		Email:      jsonString(data["email"]),
		Avatar:     jsonString(data["avatar_url"]),
	}, nil
}

func (p *larkProvider) FetchContacts(schoolID uint) (*ThirdPartyContacts, error) {
	token := p.tenantToken()
	if token == "" {
		return nil, errors.New("飞书未配置或获取 tenant_access_token 失败")
	}

	departments := p.allDepartments(token)
	// 飞书成员上只带 department_ids，前端按「部门名」匹配班级 —— 这里翻译成名称。
	deptNameByID := map[string]string{}
	for _, d := range departments {
		did := jsonString(d["department_id"])
		if did != "" {
			deptNameByID[did] = jsonString(d["name"])
		}
	}

	members := make([]ThirdPartyMember, 0, len(departments))
	seen := map[string]bool{}
	for _, d := range departments {
		deptID := jsonString(d["department_id"])
		if deptID == "" {
			continue
		}
		for _, u := range p.deptUsers(token, deptID) {
			uid := jsonString(u["union_id"])
			if uid == "" {
				uid = jsonString(u["user_id"])
			}
			if uid == "" || seen[uid] {
				continue
			}
			seen[uid] = true

			deptNames := make([]string, 0, 1)
			for _, raw := range jsonList(u["department_ids"]) {
				name := deptNameByID[jsonString(raw)]
				if name == "" || containsString(deptNames, name) {
					continue
				}
				deptNames = append(deptNames, name)
			}
			members = append(members, ThirdPartyMember{
				UserID:          uid,
				Name:            jsonString(u["name"]),
				Mobile:          jsonString(u["mobile"]),
				Email:           jsonString(u["email"]),
				Position:        jsonString(u["job_title"]),
				DepartmentNames: deptNames,
			})
		}
	}

	return &ThirdPartyContacts{Departments: departments, Members: members}, nil
}

// tenantToken 飞书 tenant_access_token（缺失返回空串，同 Laravel）。
func (p *larkProvider) tenantToken() string {
	r, err := p.manager.requestJSON(http.MethodPost,
		p.manager.LarkAPIBase+"/open-apis/auth/v3/tenant_access_token/internal",
		map[string]any{"app_id": p.manager.LarkAppID, "app_secret": p.manager.LarkAppSecret}, nil)
	if err != nil {
		return ""
	}
	return jsonString(r["tenant_access_token"])
}

func (p *larkProvider) allDepartments(token string) []map[string]any {
	r, err := p.manager.requestJSON(http.MethodGet,
		p.manager.LarkAPIBase+"/open-apis/contact/v3/departments?fetch_child=true&page_size=100",
		nil, map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		return nil
	}
	return jsonObjectList(jsonNestedMap(r, "data")["items"])
}

func (p *larkProvider) deptUsers(token, deptID string) []map[string]any {
	r, err := p.manager.requestJSON(http.MethodGet,
		p.manager.LarkAPIBase+"/open-apis/contact/v3/users/find_by_department?department_id="+
			url.QueryEscape(deptID)+"&page_size=100",
		nil, map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		return nil
	}
	return jsonObjectList(jsonNestedMap(r, "data")["items"])
}
