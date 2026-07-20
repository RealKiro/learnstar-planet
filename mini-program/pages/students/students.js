const app = getApp();
Page({
  data: { students: [], loading: true, search: '' },
  onShow() { this.loadData(); },
  loadData() {
    this.setData({ loading: true });
    app.apiRequest({ path: '/api/v1/teacher/students?per_page=200' }).then((res) => {
      this.setData({ students: res.data || [], loading: false });
    }).catch(() => this.setData({ loading: false }));
  },
  onSearch(e) {
    const val = e.detail.value.toLowerCase();
    const all = this.data._all || this.data.students;
    if (!val) return this.setData({ students: all });
    const filtered = all.filter(s => s.name.toLowerCase().includes(val) || (s.student_no || '').includes(val));
    this.setData({ students: filtered, _all: all });
  },
  viewDetail(e) {
    const id = e.currentTarget.dataset.id;
    wx.navigateTo({ url: `/pages/pets/pets?student_id=${id}` });
  }
});
