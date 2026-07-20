const fs = require('fs');
let c = fs.readFileSync('frontend-vue/src/pages/admin/SchoolSettingsPage.vue', 'utf8');

// Add demo functions before </script>
const demoFuncs = `
const demoLoading = ref(false)

async function seedDemo() {
  demoLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/demo/seed', { method: 'POST', headers: { 'Authorization': 'Bearer ' + localStorage.getItem('auth_token'), 'Content-Type': 'application/json' } })
    const data = await res.json()
    toast.show(data.message || '演示数据已生成', res.ok ? 'success' : 'error')
  } catch { toast.show('操作失败', 'error') }
  finally { demoLoading.value = false }
}
async function cleanDemo() {
  demoLoading.value = true
  try {
    const res = await fetch('/api/v1/admin/demo/clean', { method: 'POST', headers: { 'Authorization': 'Bearer ' + localStorage.getItem('auth_token'), 'Content-Type': 'application/json' } })
    const data = await res.json()
    toast.show(data.message || '演示数据已清除', res.ok ? 'success' : 'error')
  } catch { toast.show('操作失败', 'error') }
  finally { demoLoading.value = false }
}
`;

c = c.replace('</script>', demoFuncs + '\n</script>');

// Add demo section before </template>
const demoSect = `
    <!-- 演示数据管理 -->
    <div class="card" style="max-width:640px;padding:32px;margin-top:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:4px;">🧪 演示数据管理</h3>
      <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:16px;">
        管理员: demo_admin / demo123 · 教师: demo_t1~t4 / demo123 · 班级码: DEMO00
      </p>
      <div style="display:flex;gap:12px;">
        <button class="btn btn-sm btn-primary" :disabled="demoLoading" @click="seedDemo">{{ demoLoading ? '处理中...' : '\u{1F4E5} 生成演示数据' }}</button>
        <button class="btn btn-sm" :disabled="demoLoading" @click="cleanDemo"
          style="background:var(--color-bg-card);color:var(--color-danger);border:1px solid rgba(239,68,68,0.2);">
          {{ demoLoading ? '处理中...' : '\u{1F5D1}️ 清除演示数据' }}
        </button>
      </div>
    </div>
`;

c = c.replace('</template>', demoSect + '\n</template>');

fs.writeFileSync('frontend-vue/src/pages/admin/SchoolSettingsPage.vue', c);
console.log('Fixed SchoolSettingsPage');
