const app = getApp();
Page({
  data: { notices: [], loading: true, studentId: null },
  onLoad(params) { this.setData({ studentId: params.student_id }); this.loadData(); },
  loadData() {
    this.setData({ loading: true });
    app.apiRequest({ path: '/api/v1/parent/notices' })
      .then((res) => { this.setData({ notices: res.data || [], loading: false }); })
      .catch(() => this.setData({ loading: false }));
  },
  read(e) {
    const id = e.currentTarget.dataset.id;
    app.apiRequest({ path: `/api/v1/parent/notices/${id}`, method: 'GET' });
  }
});
