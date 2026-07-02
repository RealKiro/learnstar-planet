// 学趣星球 - 登录页（教师登录 / 管理员登录 分开）
Page({
  data: {
    loginType: 'teacher',  // 默认教师登录
    username: '',
    password: '',
    adminUsername: '',
    adminPassword: '',
    loading: false,
    showBinding: false,   // 扫码后需要绑定教师账号
    tempToken: '',
    bindUsername: '',
    bindPassword: '',
    teacherAttempts: 0,
    bindAttempts: 0,
  },

  onLoad() {
    // 已登录则直接跳转
    const app = getApp();
    if (app.globalData.token) {
      this.navigateByRole();
    }
  },

  // 切换登录类型
  switchLoginType(e) {
    const type = e.currentTarget.dataset.type;
    this.setData({ loginType: type, showBinding: false });
  },

  // 教师输入绑定
  bindUsernameInput(e) { this.setData({ username: e.detail.value }); },
  bindPasswordInput(e) { this.setData({ password: e.detail.value }); },

  // 管理员输入绑定
  bindAdminUsernameInput(e) { this.setData({ adminUsername: e.detail.value }); },
  bindAdminPasswordInput(e) { this.setData({ adminPassword: e.detail.value }); },

  // 绑定表单输入
  bindBindUsernameInput(e) { this.setData({ bindUsername: e.detail.value }); },
  bindBindPasswordInput(e) { this.setData({ bindPassword: e.detail.value }); },

  // 取消绑定，返回登录
  cancelBinding() { this.setData({ showBinding: false }); },

  // ===== 教师账号密码登录 =====
  handleTeacherLogin() {
    if (!this.data.username || !this.data.password) {
      wx.showToast({ title: '请输入教师账号和密码', icon: 'none' });
      return;
    }

    this.setData({ loading: true });
    const app = getApp();

    app.loginWithCredentials(this.data.username, this.data.password)
      .then(() => {
        this.setData({ loading: false, teacherAttempts: 0 });
        if (app.globalData.role === 'teacher') {
          wx.switchTab({ url: '/pages/dashboard/dashboard' });
        } else {
          wx.showToast({ title: '此账号非教师账号，请切换管理员登录', icon: 'none' });
          this.setData({ loading: false });
        }
      })
      .catch((err) => {
        const attempts = this.data.teacherAttempts + 1;
        const maxAttempts = 3;
        this.setData({ loading: false, teacherAttempts: attempts });
        if (attempts >= maxAttempts) {
          wx.showToast({ title: '密码错误次数过多，请联系管理员修改密码', icon: 'none' });
        } else {
          wx.showToast({ title: (err.message || '账号或密码错误') + '，请重试（还剩 ' + (maxAttempts - attempts) + ' 次）', icon: 'none' });
        }
      });
  },

  // ===== 管理员账号密码登录 =====
  handleAdminLogin() {
    if (!this.data.adminUsername || !this.data.adminPassword) {
      wx.showToast({ title: '请输入管理员账号和密码', icon: 'none' });
      return;
    }

    this.setData({ loading: true });
    const app = getApp();

    app.loginWithCredentials(this.data.adminUsername, this.data.adminPassword)
      .then(() => {
        this.setData({ loading: false });
        if (app.globalData.role === 'school_admin') {
          wx.redirectTo({ url: '/pages/dashboard/dashboard' });
        } else {
          wx.showToast({ title: '此账号非管理员账号，请切换教师登录', icon: 'none' });
          this.setData({ loading: false });
        }
      })
      .catch((err) => {
        this.setData({ loading: false });
        wx.showToast({ title: err.message || '管理员登录失败', icon: 'none' });
      });
  },

  // ===== 微信一键登录（仅教师可用） =====
  handleWxLogin() {
    this.setData({ loading: true });
    const app = getApp();

    app.wxSilentLogin()
      .then((res) => {
        this.setData({ loading: false });
        if (res.status === 'logged_in') {
          this.navigateByRole();
        } else if (res.status === 'need_binding') {
          this.setData({
            showBinding: true,
            tempToken: res.temp_token,
          });
          wx.showToast({ title: '请输入教师账号密码完成绑定', icon: 'none' });
        }
      })
      .catch(() => {
        this.setData({ loading: false });
        wx.showToast({ title: '微信登录失败', icon: 'none' });
      });
  },

  // 扫码后绑定教师账号
  handleBind() {
    if (!this.data.bindUsername || !this.data.bindPassword) {
      wx.showToast({ title: '请输入教师账号和密码', icon: 'none' });
      return;
    }

    const app = getApp();
    app.bindAfterScan(this.data.tempToken, this.data.bindUsername, this.data.bindPassword)
      .then(() => {
        this.setData({ showBinding: false, bindAttempts: 0 });
        wx.switchTab({ url: '/pages/dashboard/dashboard' });
      })
      .catch((err) => {
        const attempts = this.data.bindAttempts + 1;
        const maxAttempts = 3;
        this.setData({ bindAttempts: attempts });
        if (attempts >= maxAttempts) {
          wx.showToast({ title: '密码错误次数过多，请联系管理员修改密码', icon: 'none' });
        } else {
          wx.showToast({ title: (err.message || '账号或密码错误') + '，请重试（还剩 ' + (maxAttempts - attempts) + ' 次）', icon: 'none' });
        }
      });
  },

  navigateByRole() {
    const app = getApp();
    if (app.globalData.role === 'teacher') {
      wx.switchTab({ url: '/pages/dashboard/dashboard' });
    } else if (app.globalData.role === 'parent') {
      wx.redirectTo({ url: '/pages/parent-home/parent-home' });
    } else if (app.globalData.role === 'school_admin') {
      wx.redirectTo({ url: '/pages/dashboard/dashboard' });
    }
  }
});
