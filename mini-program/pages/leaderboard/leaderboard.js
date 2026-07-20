const app = getApp();
Page({
  data: { tab: 'total', totalList: [], weeklyList: [], petList: [], loading: true },
  onShow() { this.loadData(); },
  loadData() {
    this.setData({ loading: true });
    Promise.all([
      app.apiRequest({ path: '/api/v1/teacher/leaderboard/total' }),
      app.apiRequest({ path: '/api/v1/teacher/leaderboard/weekly' }),
      app.apiRequest({ path: '/api/v1/teacher/leaderboard/pet-level' }),
    ]).then(([total, weekly, pet]) => {
      this.setData({
        totalList: (total.data || []).slice(0, 20),
        weeklyList: (weekly.data || []).slice(0, 20),
        petList: (pet.data || []).slice(0, 20),
        loading: false,
      });
    }).catch(() => this.setData({ loading: false }));
  },
  switchTab(e) { this.setData({ tab: e.currentTarget.dataset.tab }); },
  medal(r) {
    if (r === 1) return '🥇'; if (r === 2) return '🥈'; if (r === 3) return '🥉';
    return r;
  }
});
