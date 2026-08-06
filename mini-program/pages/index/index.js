// ============================================================
// 班宠星球 - 微信小程序首页
// ============================================================

const app = getApp();

Page({
  data: {
    isTeacher: false,
    classRooms: [],
    currentClassId: null,
    freeBadge: true
  },

  onLoad() {
    const role = app.globalData.role;
    this.setData({
      isTeacher: role === 'teacher',
      freeBadge: true
    });

    if (role === 'teacher') {
      this.loadTeacherData();
    }
  },

  /**
   * 教师端：加载班级列表
   */
  loadTeacherData() {
    app.apiRequest('/teacher/classes', 'GET', {}, (data) => {
      this.setData({
        classRooms: data,
        currentClassId: data.length > 0 ? data[0].id : null
      });
      if (data.length > 0) {
        wx.switchTab({ url: '/pages/dashboard/dashboard' });
      }
    });
  },

  /**
   * 选择班级
   */
  selectClass(e) {
    const classId = e.currentTarget.dataset.id;
    this.setData({ currentClassId: classId });
    app.globalData.classRoomId = classId;
    wx.switchTab({ url: '/pages/dashboard/dashboard' });
  },

  /**
   * 切换角色（用于演示，家长端已移除）
   */
  switchRole() {
    app.globalData.role = 'teacher';
    this.setData({ isTeacher: true });
  }
});
