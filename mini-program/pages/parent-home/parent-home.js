const app = getApp();
Page({
  data: { children: [], loading: true },
  onShow() { this.loadData(); },
  loadData() {
    this.setData({ loading: true });
    app.apiRequest({ path: '/api/v1/parent/home' }).then((res) => {
      this.setData({ children: res.data?.children || [], loading: false });
    }).catch(() => this.setData({ loading: false }));
  }
});
