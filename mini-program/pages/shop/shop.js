const app = getApp();
Page({
  data: { items: [], redemptions: [], loading: true, tab: 'shop' },
  onShow() { this.loadData(); },
  loadData() {
    this.setData({ loading: true });
    Promise.all([
      app.apiRequest({ path: '/api/v1/teacher/shop/items' }),
      app.apiRequest({ path: '/api/v1/teacher/shop/redemptions' }),
    ]).then(([items, redemptions]) => {
      this.setData({ items: items.data || [], redemptions: redemptions.data || [], loading: false });
    }).catch(() => this.setData({ loading: false }));
  },
  switchTab(e) { this.setData({ tab: e.currentTarget.dataset.tab }); },
  approve(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({ title: '确认', content: '确认通过此兑换申请？', success: (r) => {
      if (r.confirm) app.apiRequest({ path: `/api/v1/teacher/shop/redemptions/${id}/approve`, method: 'PUT' })
        .then(() => { wx.showToast({ title: '已通过', icon: 'success' }); this.loadData(); });
    }});
  }
});
