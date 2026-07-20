const app = getApp();
Page({
  data: { students: [], rules: [], loading: true, selectedRule: '', step: 1 },
  onShow() { this.loadData(); },
  loadData() {
    this.setData({ loading: true });
    Promise.all([
      app.apiRequest({ path: `/api/v1/teacher/students?per_page=100` }),
      app.apiRequest({ path: `/api/v1/teacher/scores/rules` }),
    ]).then(([students, rules]) => {
      const ruleList = (rules.data || []).map(r => ({ name: r.name, points: r.points, is_penalty: r.is_penalty }));
      this.setData({ students: students.data || [], rules: ruleList, loading: false, selectedRule: ruleList[0]?.name || '' });
    }).catch(() => this.setData({ loading: false }));
  },
  giveScore(e) {
    const { id, name } = e.currentTarget.dataset;
    const rule = this.data.rules.find(r => r.name === this.data.selectedRule);
    if (!rule) return wx.showToast({ title: '请选择积分规则', icon: 'none' });
    wx.showModal({
      title: `${rule.is_penalty ? '扣' : '加'}分确认`,
      content: `对 ${name} ${rule.is_penalty ? '扣' : '加'} ${Math.abs(rule.points)} 分（${rule.name}）`,
      success: (res) => {
        if (res.confirm) {
          app.apiRequest({
            path: '/api/v1/teacher/scores/give',
            method: 'POST',
            data: { student_id: id, points: rule.points, reason: rule.name },
          }).then(() => {
            wx.showToast({ title: '操作成功', icon: 'success' });
            this.loadData();
          });
        }
      }
    });
  },
  onRuleChange(e) { this.setData({ selectedRule: this.data.rules[e.detail.value]?.name || '' }); }
});
