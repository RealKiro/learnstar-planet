// ============================================================
// 班宠星球 - 教师端数据看板（微信小程序）
// ============================================================

const app = getApp();

Page({
  data: {
    classRoomId: null,
    stats: { net_total: 0, today_add: 0, score_count: 0 },
    topStudents: [],
    recentScores: [],
    petDistribution: {},
    scoreTrend: [],
    loading: true,
    freeBadge: true
  },

  onShow() {
    this.setData({ classRoomId: app.globalData.classRoomId });
    if (this.data.classRoomId) {
      this.loadDashboard();
    }
  },

  onPullDownRefresh() {
    this.loadDashboard();
    wx.stopPullDownRefresh();
  },

  loadDashboard() {
    this.setData({ loading: true });

    // 加载统计数据
    app.apiRequest(`/teacher/scores/stats/${this.data.classRoomId}`, 'GET', {}, (data) => {
      this.setData({ stats: data });
    });

    // 加载排行榜
    app.apiRequest(`/teacher/leaderboard/${this.data.classRoomId}`, 'GET', {}, (data) => {
      this.setData({ topStudents: data.slice(0, 10), loading: false });
    });

    // 加载积分趋势
    app.apiRequest(`/teacher/scores/trend/${this.data.classRoomId}`, 'GET', {}, (data) => {
      this.setData({ scoreTrend: data });
    });
  },

  /**
   * 快捷加分
   */
  quickAddScore(e) {
    const { studentId, ruleCode } = e.currentTarget.dataset;
    app.apiRequest('/teacher/scores/quick', 'POST', {
      student_id: studentId,
      rule_code: ruleCode
    }, () => {
      wx.showToast({ title: '加分成功', icon: 'success' });
      this.loadDashboard();
    });
  },

  /**
   * 查看学生详情
   */
  viewStudent(e) {
    const studentId = e.currentTarget.dataset.id;
    wx.navigateTo({ url: `/pages/students/students?id=${studentId}` });
  },

  /**
   * 前往积分管理
   */
  goScores() {
    wx.switchTab({ url: '/pages/scores/scores' });
  }
});
