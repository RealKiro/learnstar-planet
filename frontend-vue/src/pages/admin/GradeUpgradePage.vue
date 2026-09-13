<script setup lang="ts">
import { ref } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { openConfirm } from '@/components/common/ConfirmDialog.vue'
import type { ApiResponse, GradeUpgradePreview } from '@/types'

const preview = ref<GradeUpgradePreview | null>(null)
const loading = ref(false)
const executing = ref(false)
const executeStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const hasPreviewed = ref(false)

async function loadPreview() {
  loading.value = true
  hasPreviewed.value = true
  try {
    const res = await apiGet<ApiResponse<GradeUpgradePreview>>('/api/v1/admin/grade-upgrade/preview')
    preview.value = (res as unknown as { data: GradeUpgradePreview }).data
  } catch { preview.value = null }
  finally { loading.value = false }
}

async function executeUpgrade() {
  if (!preview.value) return
  const upgradeCount = preview.value.summary.upgrade_class_count
  const graduateCount = preview.value.summary.graduate_class_count
  const ok = await openConfirm({
    title: '执行年级升级',
    message: `确认执行年级升级？\n\n将升级 ${upgradeCount} 个班级，毕业 ${graduateCount} 个班级。\n此操作不可撤销！`,
    danger: true,
    confirmText: '确认执行',
  })
  if (!ok) return
  executing.value = true
  executeStatus.value = 'loading'
  try {
    await apiPost('/api/v1/admin/grade-upgrade/execute')
    executeStatus.value = 'success'
    setTimeout(() => { executeStatus.value = 'idle' }, 2000)
    await loadPreview()
  } catch {
    executeStatus.value = 'error'
    setTimeout(() => { executeStatus.value = 'idle' }, 3000)
  }
  finally { executing.value = false }
}
</script>

<template>
  <div>
    <div class="gup-head">
      <div>
        <p class="page-eyebrow">学年管理</p>
        <h2 class="page-title">年级升级</h2>
      </div>
      <div class="gup-row-8">
        <button class="btn btn-sm gup-btn-plain" @click="loadPreview" :disabled="loading">
          {{ loading ? '加载中...' : '🔄 重新预览' }}
        </button>
        <button v-if="preview" class="btn btn-sm gup-btn-solid" :disabled="executeStatus === 'loading'" @click="executeUpgrade"
          :style="{ background: executeStatus === 'loading' ? '#f59e0b' : executeStatus === 'success' ? '#10b981' : executeStatus === 'error' ? '#ef4444' : '#7c3aed' }">
          <span v-if="executeStatus === 'idle'">▶ 执行升级</span>
          <span v-else-if="executeStatus === 'loading'">⏳ 执行中...</span>
          <span v-else-if="executeStatus === 'success'">✅ 升级完成</span>
          <span v-else>❌ 执行失败</span>
        </button>
      </div>
    </div>

    <!-- 提示卡片 -->
    <div class="card gup-info-card">
      <div class="gup-info-row">
        <span class="gup-emoji-24">ℹ️</span>
        <div>
          <p class="gup-info-title">关于年级升级</p>
          <p class="gup-info-desc">
            执行升级后，所有班级将整体上升一个年级（如三年级 → 四年级），六年级班级学生将标记为毕业并归档。
            请先点击「重新预览」查看变更详情，确认无误后再执行。
          </p>
        </div>
      </div>
    </div>

    <!-- 未预览 -->
    <div v-if="!hasPreviewed && !loading" class="card empty-state">
      <div class="empty-state__icon">📊</div>
      <p class="gup-mb-16">点击下方按钮预览年级升级方案</p>
      <button class="btn btn-primary" @click="loadPreview">🔍 查看升级预览</button>
    </div>

    <!-- 加载中 -->
    <div v-else-if="loading" class="empty-state">加载中...</div>

    <!-- 预览数据 -->
    <template v-else-if="preview">
      <!-- 汇总卡片 -->
      <div class="stats-grid gup-mb-24">
        <div class="stat-card stat-card--primary">
          <span class="stat-card__icon">⬆️</span>
          <div class="stat-card__value">{{ preview.summary.upgrade_class_count }}</div>
          <div class="stat-card__label">升级班级</div>
        </div>
        <div class="stat-card stat-card--secondary">
          <span class="stat-card__icon">👨‍🎓</span>
          <div class="stat-card__value">{{ preview.summary.upgrade_student_count }}</div>
          <div class="stat-card__label">升级学生</div>
        </div>
        <div class="stat-card stat-card--accent">
          <span class="stat-card__icon">🎓</span>
          <div class="stat-card__value">{{ preview.summary.graduate_class_count }}</div>
          <div class="stat-card__label">毕业班级</div>
        </div>
        <div class="stat-card stat-card--info">
          <span class="stat-card__icon">📜</span>
          <div class="stat-card__value">{{ preview.summary.graduate_student_count }}</div>
          <div class="stat-card__label">毕业学生</div>
        </div>
      </div>

      <div v-if="preview.summary.note" class="card gup-warn-card">
        <p class="gup-warn-text">⚠️ {{ preview.summary.note }}</p>
      </div>

      <!-- 待升级班级 -->
      <div class="card gup-mb-24">
        <h3 class="section-title">⬆️ 待升级班级（{{ preview.upgrade_classes.length }}）</h3>
        <div v-if="preview.upgrade_classes.length === 0" class="muted-center">无待升级班级</div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>当前班级</th><th>当前年级</th><th>→</th><th>升级后名称</th><th>升级后年级</th><th>学生数</th></tr></thead>
            <tbody>
              <tr v-for="(c, i) in preview.upgrade_classes" :key="i">
                <td class="gup-fw-600">{{ c.class_name }}</td>
                <td><span class="gup-pill-gray">{{ c.old_grade }}</span></td>
                <td class="gup-secondary">→</td>
                <td class="gup-fw-600">{{ c.new_name }}</td>
                <td><span class="gup-pill-primary">{{ c.new_grade }}</span></td>
                <td>{{ c.student_count }} 人</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 待毕业班级 -->
      <div class="card">
        <h3 class="section-title">🎓 待毕业班级（{{ preview.graduate_classes.length }}）</h3>
        <div v-if="preview.graduate_classes.length === 0" class="muted-center">无待毕业班级</div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>班级名称</th><th>学生数</th></tr></thead>
            <tbody>
              <tr v-for="(c, i) in preview.graduate_classes" :key="i">
                <td class="gup-fw-600">{{ c.class_name }}</td>
                <td>{{ c.student_count }} 人</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- 预览为空 -->
    <div v-else class="card empty-state">
      <div class="empty-state__icon">📭</div>
      <p>暂无升级预览数据</p>
    </div>
  </div>
</template>

<style scoped>
/* ===== P1 内联样式收口（声明逐字保留以保渲染等价） ===== */
.gup-fw-600 { font-weight:600; }
.gup-mb-24 { margin-bottom:24px; }
.gup-head { display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;flex-wrap:wrap;gap:12px; }
.gup-row-8 { display:flex;gap:8px; }
.gup-btn-plain { background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border); }
.gup-btn-solid { transition:all 0.3s ease;border:none;color:#fff; }
.gup-info-card { padding:20px 24px;margin-bottom:24px;background:rgba(79,70,229,0.04);border:1px solid rgba(79,70,229,0.15); }
.gup-info-row { display:flex;align-items:flex-start;gap:12px; }
.gup-emoji-24 { font-size:24px; }
.gup-info-title { font-weight:600;margin-bottom:4px; }
.gup-info-desc { font-size:13px;color:var(--color-text-secondary);line-height:1.6; }
.gup-mb-16 { margin-bottom:16px; }
.gup-warn-card { padding:12px 20px;margin-bottom:16px;background:rgba(245,158,11,0.06);border:1px solid rgba(245,158,11,0.2); }
.gup-warn-text { font-size:13px;color: var(--c-amber-deeper); }
.gup-pill-gray { display:inline-block;padding:2px 10px;border-radius:20px;font-size:12px;background:rgba(107,114,128,0.1);color:var(--color-text-secondary); }
.gup-secondary { color:var(--color-text-secondary); }
.gup-pill-primary { display:inline-block;padding:2px 10px;border-radius:20px;font-size:12px;font-weight:600;background:rgba(79,70,229,0.08);color:var(--color-primary); }
</style>
