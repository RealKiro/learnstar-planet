const app = getApp();
Page({
  data: { scores: [], loading: true, studentId: null },
  onLoad(params) { this.setData({ studentId: params.student_id }); this.loadData(); },
  loadData() {
    if (!this.data.studentId) return;
    this.setData({ loading: true });
    app.apiRequest({ path: `/api/v1/parent/scores/detail?student_id=${this.data.studentId}` })
      .then((res) => { this.setData({ scores: res.data || [], loading: false }); })
      .catch(() => this.setData({ loading: false }));
  }
});
