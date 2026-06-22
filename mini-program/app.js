// 学趣星球 - 微信小程序 App
App({
  globalData: {
    apiBase: 'https://your-domain.com/api',
    userInfo: null,
    role: null,       // teacher / parent / admin
    schoolId: null,
    classId: null,
    openid: null,
    unionid: null,
    token: null,
  },

  onLaunch() {
    // 检查本地缓存的登录状态
    const token = wx.getStorageSync('token');
    if (token) {
      this.globalData.token = token;
      this.globalData.role = wx.getStorageSync('role');
      this.globalData.userInfo = wx.getStorageSync('userInfo');
    }
  },

  // 微信静默登录（获取openid）
  wxSilentLogin() {
    return new Promise((resolve, reject) => {
      wx.login({
        success: (res) => {
          if (res.code) {
            // 将code发送到后端换取openid
            wx.request({
              url: this.globalData.apiBase + '/auth/login/wechat',
              method: 'POST',
              data: { code: res.code },
              success: (apiRes) => {
                if (apiRes.data.status === 'logged_in') {
                  // 已绑定账号，直接登录成功
                  this.handleLoginSuccess(apiRes.data);
                  resolve(apiRes.data);
                } else if (apiRes.data.status === 'need_binding') {
                  // 未绑定账号，需要输入账号密码绑定
                  resolve(apiRes.data);
                }
              },
              fail: reject
            });
          }
        },
        fail: reject
      });
    });
  },

  // 账号密码登录
  loginWithCredentials(username, password) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: this.globalData.apiBase + '/auth/login',
        method: 'POST',
        data: { username, password },
        success: (res) => {
          if (res.data.token) {
            this.handleLoginSuccess(res.data);
            resolve(res.data);
          } else {
            reject(new Error(res.data.message || '登录失败'));
          }
        },
        fail: reject
      });
    });
  },

  // 扫码后绑定账号
  bindAfterScan(tempToken, username, password) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: this.globalData.apiBase + '/auth/bind-after-scan',
        method: 'POST',
        data: {
          temp_token: tempToken,
          username,
          password,
          platform: 'wechat',
          platform_id: this.globalData.openid,
        },
        success: (res) => {
          if (res.data.status === 'bound') {
            this.handleLoginSuccess(res.data);
            resolve(res.data);
          } else {
            reject(new Error(res.data.message || '绑定失败'));
          }
        },
        fail: reject
      });
    });
  },

  // 登录成功后处理
  handleLoginSuccess(data) {
    this.globalData.token = data.token || data.user.token;
    this.globalData.role = data.user.role;
    this.globalData.userInfo = data.user;
    this.globalData.schoolId = data.user.school_id;

    wx.setStorageSync('token', this.globalData.token);
    wx.setStorageSync('role', this.globalData.role);
    wx.setStorageSync('userInfo', this.globalData.userInfo);

    // 根据角色跳转不同首页
    if (this.globalData.role === 'teacher') {
      wx.switchTab({ url: '/pages/dashboard/dashboard' });
    } else if (this.globalData.role === 'parent') {
      wx.redirectTo({ url: '/pages/parent-home/parent-home' });
    }
  },

  // 统一API请求
  apiRequest(options) {
    const token = this.globalData.token;
    return new Promise((resolve, reject) => {
      wx.request({
        url: this.globalData.apiBase + options.path,
        method: options.method || 'GET',
        data: options.data || {},
        header: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        success: (res) => {
          if (res.statusCode === 401) {
            // token过期，清除登录状态
            this.clearLogin();
            wx.redirectTo({ url: '/pages/login/login' });
            reject(new Error('登录过期'));
          } else {
            resolve(res.data);
          }
        },
        fail: reject
      });
    });
  },

  clearLogin() {
    this.globalData.token = null;
    this.globalData.role = null;
    this.globalData.userInfo = null;
    wx.removeStorageSync('token');
    wx.removeStorageSync('role');
    wx.removeStorageSync('userInfo');
  }
});
