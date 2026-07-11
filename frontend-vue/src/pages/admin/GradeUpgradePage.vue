<script setup lang="ts">
import { ref } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, GradeUpgradePreview } from '@/types'

const toast = useToastStore()

const preview = ref<GradeUpgradePreview | null>(null)
const loading = ref(false)
const executing = ref(false)
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
  if (!confirm(`确认执行年级升级？\n\n将升级 ${upgradeCount} 个班级，毕业 ${graduateCount} 个班级。\n此操作不可撤销！`)) return
  executing.value = true
  try {
    await apiPost('/api/v1/admin/grade-upgrade/execute')
    toast.show('年级升级已执行完成', 'success')
    // 重新拉取预览
    await loadPreview()
  } catch { /* handled */ }
  finally { executing.value = false }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;flex-wrap:wrap;gap:12px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">学年管理</p>
        <h2 style="font-size:24px;font-weight:700;">年级升级</h2>
      </div>
      <div style="display:flex;gap:8px;">
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="loadPreview" :disabled="loading">
          {{ loading ? '加载中...' : '🔄 重新预览' }}
        </button>
        <button v-if="preview" class="btn btn-sm btn-primary" :disabled="executing" @click="executeUpgrade">
          {{ executing ? '执行中...' : '▶ 执行升级' }}
        </button>
      </div>
    </div>

    <!-- 提示卡片 -->
    <div class="card" style="padding:20px 24px;margin-bottom:24px;background:rgba(79,70,229,0.04);border:1px solid rgba(79,70,229,0.15);">
      <div style="display:flex;align-items:flex-start;gap:12px;">
        <span style="font-size:24px;">ℹ️</span>
        <div>
          <p style="font-weight:600;margin-bottom:4px;">关于年级升级</p>
          <p style="font-size:13px;color:var(--color-text-secondary);line-height:1.6;">
            执行升级后，所有班级将整体上升一个年级（如三年级 → 四年级），六年级班级学生将标记为毕业并归档。
            请先点击「重新预览」查看变更详情，确认无误后再执行。
          </p>
        </div>
      </div>
    </div>

    <!-- 未预览 -->
    <div v-if="!hasPreviewed && !loading" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📊</div>
      <p style="margin-bottom:16px;">点击下方按钮预览年级升级方案</p>
      <button class="btn btn-primary" @click="loadPreview">🔍 查看升级预览</button>
    </div>

    <!-- 加载中 -->
    <div v-else-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <!-- 预览数据 -->
    <template v-else-if="preview">
      <!-- 汇总卡片 -->
      <div class="stats-grid" style="margin-bottom:24px;">
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

      <div v-if="preview.summary.note" class="card" style="padding:12px 20px;margin-bottom:16px;background:rgba(245,158,11,0.06);border:1px solid rgba(245,158,11,0.2);">
        <p style="font-size:13px;color:#B45309;">⚠️ {{ preview.summary.note }}</p>
      </div>

      <!-- 待升级班级 -->
      <div class="card" style="margin-bottom:24px;">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">⬆️ 待升级班级（{{ preview.upgrade_classes.length }}）</h3>
        <div v-if="preview.upgrade_classes.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">无待升级班级</div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>当前班级</th><th>当前年级</th><th>→</th><th>升级后名称</th><th>升级后年级</th><th>学生数</th></tr></thead>
            <tbody>
              <tr v-for="(c, i) in preview.upgrade_classes" :key="i">
                <td style="font-weight:600;">{{ c.class_name }}</td>
                <td><span style="display:inline-block;padding:2px 10px;border-radius:20px;font-size:12px;background:rgba(107,114,128,0.1);color:#6B7280;">{{ c.old_grade }}</span></td>
                <td style="color:var(--color-text-secondary);">→</td>
                <td style="font-weight:600;">{{ c.new_name }}</td>
                <td><span style="display:inline-block;padding:2px 10px;border-radius:20px;font-size:12px;font-weight:600;background:rgba(79,70,229,0.08);color:var(--color-primary);">{{ c.new_grade }}</span></td>
                <td>{{ c.student_count }} 人</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 待毕业班级 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">🎓 待毕业班级（{{ preview.graduate_classes.length }}）</h3>
        <div v-if="preview.graduate_classes.length === 0" style="text-align:center;padding:24px;color:var(--color-text-secondary);">无待毕业班级</div>
        <div v-else class="data-table">
          <table>
            <thead><tr><th>班级名称</th><th>学生数</th></tr></thead>
            <tbody>
              <tr v-for="(c, i) in preview.graduate_classes" :key="i">
                <td style="font-weight:600;">{{ c.class_name }}</td>
                <td>{{ c.student_count }} 人</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- 预览为空 -->
    <div v-else class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">📭</div>
      <p>暂无升级预览数据</p>
    </div>
  </div>
</template>
