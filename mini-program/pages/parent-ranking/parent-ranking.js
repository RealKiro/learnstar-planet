const app = getApp();
Page({
  data: { rankings: [], myRank: null, loading: true, studentId: null },
  onLoad(params) { this.setData({ studentId: params.student_id }); this.loadData(); },
  loadData() {
    if (!this.data.studentId) return;
    this.setData({ loading: true });
    app.apiRequest({ path: `/api/v1/parent/ranking?student_id=${this.data.studentId}` })
      .then((res) => {
        const list = (res.data || []).slice(0, 30);
        const mine = list.find(r => r.is_mine);
        this.setData({ rankings: list, myRank: mine || null, loading: false });
      }).catch(() => this.setData({ loading: false }));
  }
});
