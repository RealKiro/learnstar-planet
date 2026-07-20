const app = getApp();
Page({
  data: { userInfo: null, role: '' },
  onShow() {
    this.setData({
      userInfo: app.globalData.userInfo,
      role: app.globalData.role,
    });
  },
  logout() {
    wx.showModal({
      title: '退出登录',
      content: '确定要退出当前账号吗？',
      success: (res) => {
        if (res.confirm) {
          app.clearLogin();
          wx.redirectTo({ url: '/pages/login/login' });
        }
      }
    });
  },
  copyApi() {
    wx.setClipboardData({ data: app.globalData.apiBase });
  }
});
