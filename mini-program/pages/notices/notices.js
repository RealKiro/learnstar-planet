const app = getApp();
Page({
  data: { notices: [], loading: true },
  onShow() { this.loadData(); },
  loadData() {
    this.setData({ loading: true });
    app.apiRequest({ path: '/api/v1/teacher/notices' }).then((res) => {
      this.setData({ notices: res.data || [], loading: false });
    }).catch(() => this.setData({ loading: false }));
  },
  publish(e) {
    const id = e.currentTarget.dataset.id;
    app.apiRequest({ path: `/api/v1/teacher/notices/${id}/publish`, method: 'PUT' })
      .then(() => { wx.showToast({ title: '已发布', icon: 'success' }); this.loadData(); });
  }
});
