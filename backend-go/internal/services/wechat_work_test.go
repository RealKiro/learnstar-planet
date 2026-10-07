// 企微审批查询（getapprovalinfo / getapprovaldetail / parse）的服务层测试。
//
// 假上游全部走 httptest（**不访问真实外网**）：覆盖 getapprovalinfo 的分页续页与 errcode 中断、
// 请求体口径（starttime/endtime/cursor/size/filters）、getApprovalDetail 的 nil 分支（errcode 非 0、
// info 缺失/空）与成功分支，以及 parse 的逐字段取值（学生姓名 / 请假类型 / 事由 / 起止时间 /
// sp_status 缺省 / approve_time 的严格 === 2 语义 / controls[0].text 兜底）。
package services_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newWechatWorkApprovalService 建「配置齐全 + 指向假上游」的企微服务（token 缓存落内存库）。
func newWechatWorkApprovalService(t *testing.T, handler http.HandlerFunc) (*services.WechatWorkService, func() []string) {
	t.Helper()
	srv, seen := recordingServer(t, handler)
	svc := services.NewWechatWorkService(setupDB(t))
	svc.CorpID, svc.Secret, svc.APIBase = "corp-1", "sec-1", srv.URL
	svc.SetHTTPClient(srv.Client())
	return svc, seen
}

// TestWechatWorkGetLeaveApprovalSpNosPaging 分页续页 + 空串丢弃 + 请求体口径。
func TestWechatWorkGetLeaveApprovalSpNosPaging(t *testing.T) {
	bodies := []map[string]any{}
	svc, seen := newWechatWorkApprovalService(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"tok-1","expires_in":7200}`)
		case "/cgi-bin/oa/getapprovalinfo":
			raw, _ := io.ReadAll(r.Body)
			payload := map[string]any{}
			_ = json.Unmarshal(raw, &payload)
			bodies = append(bodies, payload)
			if payload["cursor"] == float64(0) {
				// 空串会被丢弃；next_cursor > 0 触发续页。
				jsonBody(w, `{"errcode":0,"sp_no_list":["SP-1","","SP-2"],"next_cursor":2}`)
				return
			}
			jsonBody(w, `{"errcode":0,"sp_no_list":["SP-3"],"next_cursor":0}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	spNos, err := svc.GetLeaveApprovalSpNos(7, 1700000000, 1700086399)
	require.NoError(t, err)
	assert.Equal(t, []string{"SP-1", "SP-2", "SP-3"}, spNos)

	assert.Equal(t, []string{
		"GET /cgi-bin/gettoken",
		"POST /cgi-bin/oa/getapprovalinfo",
		"POST /cgi-bin/oa/getapprovalinfo",
	}, seen())

	require.Len(t, bodies, 2)
	assert.Equal(t, float64(1700000000), bodies[0]["starttime"])
	assert.Equal(t, float64(1700086399), bodies[0]["endtime"])
	assert.Equal(t, float64(0), bodies[0]["cursor"])
	assert.Equal(t, float64(100), bodies[0]["size"])
	assert.Equal(t, map[string]any{"key": "sp_status", "value": "2"}, bodies[0]["filters"])
	assert.Equal(t, float64(2), bodies[1]["cursor"], "next_cursor 应作为下一页的 cursor")
}

// TestWechatWorkGetLeaveApprovalSpNosStopsOnErrcode errcode != 0 立即中断（已收集的单号保留，该页列表丢弃）。
func TestWechatWorkGetLeaveApprovalSpNosStopsOnErrcode(t *testing.T) {
	approvalCalls := 0
	svc, _ := newWechatWorkApprovalService(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/gettoken":
			jsonBody(w, `{"errcode":0,"access_token":"tok-1"}`)
		case "/cgi-bin/oa/getapprovalinfo":
			approvalCalls++
			if approvalCalls == 1 {
				jsonBody(w, `{"errcode":0,"sp_no_list":["SP-1"],"next_cursor":5}`)
				return
			}
			jsonBody(w, `{"errcode":40001,"errmsg":"invalid credential","sp_no_list":["SP-9"],"next_cursor":9}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	spNos, err := svc.GetLeaveApprovalSpNos(1, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, []string{"SP-1"}, spNos, "errcode 非 0 时应中断且丢弃该页列表")
	assert.Equal(t, 2, approvalCalls, "errcode 非 0 后不应再续页")
}

// TestWechatWorkGetApprovalDetailNil 上游 errcode 非 0 / info 缺失 / info 为空 / 响应非 JSON → (nil, nil)。
func TestWechatWorkGetApprovalDetailNil(t *testing.T) {
	cases := map[string]string{
		"errcode 非 0": `{"errcode":40001,"errmsg":"invalid sp_no"}`,
		"info 缺失":     `{"errcode":0}`,
		"info 为空对象":   `{"errcode":0,"info":{}}`,
		"info 非对象":    `{"errcode":0,"info":[]}`,
		"响应非 JSON":    `not-json-at-all`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			svc, _ := newWechatWorkApprovalService(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/cgi-bin/gettoken" {
					jsonBody(w, `{"errcode":0,"access_token":"tok-1"}`)
					return
				}
				_, _ = w.Write([]byte(body))
			})

			detail, err := svc.GetApprovalDetail(1, "SP-1")
			require.NoError(t, err)
			assert.Nil(t, detail)
		})
	}
}

// TestWechatWorkGetApprovalDetailParse parse 的逐字段映射（含严格 === 2 与 controls 兜底）。
func TestWechatWorkGetApprovalDetailParse(t *testing.T) {
	const info = `{"errcode":0,"info":{
		"sp_no":"SP-1",
		"applyer":{"userid":"parent-1"},
		"sp_status":2,
		"sp_record":[
			{"sp_status":1,"approve_time":111},
			{"sp_status":2,"approve_time":222},
			{"sp_status":"2","approve_time":333}
		],
		"apply_data":{"contents":[
			{"title":"学生姓名","value":{"text":"小明"}},
			{"title":"请假类型","value":{"text":"病假"}},
			{"title":"请假事由","value":"发烧需要休息"},
			{"title":"请假时间","value":{"new_begin":1700000000,"new_end":1700086400}}
		]}
	}}`

	svc, _ := newWechatWorkApprovalService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cgi-bin/gettoken" {
			jsonBody(w, `{"errcode":0,"access_token":"tok-1"}`)
			return
		}
		assert.Equal(t, "/cgi-bin/oa/getapprovaldetail", r.URL.Path)
		raw, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"sp_no":"SP-1"}`, string(raw))
		jsonBody(w, info)
	})

	detail, err := svc.GetApprovalDetail(1, "SP-1")
	require.NoError(t, err)
	require.NotNil(t, detail)

	assert.Equal(t, "SP-1", detail.SpNo)
	assert.Equal(t, "parent-1", detail.SubmitterUserID)
	assert.Equal(t, 2, detail.ApproveStatus)
	require.NotNil(t, detail.StudentName)
	assert.Equal(t, "小明", *detail.StudentName)
	require.NotNil(t, detail.LeaveType)
	assert.Equal(t, "病假", *detail.LeaveType)
	require.NotNil(t, detail.Reason)
	assert.Equal(t, "发烧需要休息", *detail.Reason)
	assert.Equal(t, int64(1700000000), detail.LeaveStart)
	assert.Equal(t, int64(1700086400), detail.LeaveEnd)
	require.NotNil(t, detail.ApprovedAt)
	assert.Equal(t, int64(222), *detail.ApprovedAt, "字符串 \"2\" 不算通过（PHP === 严格比较）")
}

// TestWechatWorkGetApprovalDetailParseFallbacks 姓名走 controls[0].text 兜底、sp_status 缺省为 1、
// 起止时间缺失为 0、approved_at 缺失为 null。
func TestWechatWorkGetApprovalDetailParseFallbacks(t *testing.T) {
	const info = `{"errcode":0,"info":{
		"sp_no":"SP-2",
		"applyer":{},
		"apply_data":{"contents":[
			{"title":"学生","value":{"controls":[{"text":"小红"}]}},
			{"title":"起止时间","value":{"new_start":1700000000}}
		]}
	}}`

	svc, _ := newWechatWorkApprovalService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cgi-bin/gettoken" {
			jsonBody(w, `{"errcode":0,"access_token":"tok-1"}`)
			return
		}
		jsonBody(w, info)
	})

	detail, err := svc.GetApprovalDetail(1, "SP-2")
	require.NoError(t, err)
	require.NotNil(t, detail)

	assert.Equal(t, "SP-2", detail.SpNo)
	assert.Equal(t, "", detail.SubmitterUserID)
	assert.Equal(t, 1, detail.ApproveStatus, "sp_status 缺失时按 1")
	require.NotNil(t, detail.StudentName)
	assert.Equal(t, "小红", *detail.StudentName, "text 为空时用 controls[0].text")
	assert.Nil(t, detail.LeaveType)
	assert.Nil(t, detail.Reason)
	assert.Equal(t, int64(1700000000), detail.LeaveStart)
	assert.Equal(t, int64(0), detail.LeaveEnd)
	assert.Nil(t, detail.ApprovedAt)
}
