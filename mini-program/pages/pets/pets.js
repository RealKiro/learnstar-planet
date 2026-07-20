const app = getApp();
Page({
  data: { pets: [], loading: true, stats: { total: 0, avg_level: 0, max_level: 0 } },
  onShow() { this.loadData(); },
  loadData() {
    this.setData({ loading: true });
    app.apiRequest({ path: '/api/v1/teacher/pets/overview' }).then((res) => {
      const list = (res.data || []).map(p => ({
        student_name: p.student_name, pet_name: p.pet_name || '未孵化',
        level: p.level || 0, stage: p.stage_name || '未孵化', emoji: p.emoji || '🥚',
      }));
      const total = list.length;
      const avg_level = total > 0 ? (list.reduce((s, p) => s + p.level, 0) / total).toFixed(1) : 0;
      const max_level = total > 0 ? Math.max(...list.map(p => p.level)) : 0;
      this.setData({ pets: list, stats: { total, avg_level, max_level }, loading: false });
    }).catch(() => this.setData({ loading: false }));
  }
});
