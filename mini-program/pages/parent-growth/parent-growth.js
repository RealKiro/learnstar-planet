const app = getApp();
Page({
  data: { pet: null, loading: true, studentId: null },
  onLoad(params) { this.setData({ studentId: params.student_id }); this.loadData(); },
  loadData() {
    if (!this.data.studentId) return;
    this.setData({ loading: true });
    app.apiRequest({ path: `/api/v1/parent/pet?student_id=${this.data.studentId}` })
      .then((res) => { this.setData({ pet: res.data || null, loading: false }); })
      .catch(() => this.setData({ loading: false }));
  },
  feed() {
    if (!this.data.studentId) return;
    app.apiRequest({ path: `/api/v1/parent/pet/feed?student_id=${this.data.studentId}`, method: 'POST' })
      .then(() => { wx.showToast({ title: '喂养成功', icon: 'success' }); this.loadData(); });
  }
});
