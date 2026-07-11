<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, ScoreRule, Student } from '@/types'

const toast = useToastStore()

const students = ref<Student[]>([])
const rules = ref<ScoreRule[]>([])
const scoreSummary = ref({ total: 0, today: 0, this_week: 0 })
const loading = ref(true)

const selectedStudentId = ref<number | null>(null)
const customPoints = ref(5)
const customReason = ref('')

onMounted(async () => {
  try {
    const [studentsRes, rulesRes, summaryRes] = await Promise.all([
      apiGet<ApiResponse<Student[]>>('/api/v1/teacher/students?per_page=100'),
      apiGet<ApiResponse<ScoreRule[]>>('/api/v1/teacher/scores/rules'),
      apiGet<ApiResponse<{ total: number; today: number; this_week: number }>>('/api/v1/teacher/scores/summary'),
    ])
    students.value = studentsRes.data || []
    rules.value = rulesRes.data || []
    scoreSummary.value = summaryRes.data || { total: 0, today: 0, this_week: 0 }
  } catch { /* handled */ }
  finally { loading.value = false }
})

async function handleGiveScore(student: Student, points: number, reason: string) {
  if (!points || !reason.trim()) {
    toast.show('请选择积分值和填写理由', 'error')
    return
  }
  try {
    const res = await apiPost<ApiResponse<{ new_score: number }>>('/api/v1/teacher/scores/give', {
      student_id: student.id,
      points,
      reason: reason.trim(),
    })
    toast.show(`${student.name} ${points > 0 ? '+' : ''}${points}分 — ${reason}`, 'success')
    const data = res as unknown as { data: { new_score: number } }
    if (data.data?.new_score != null) {
      student.total_score = data.data.new_score
    } else {
      student.total_score += points
    }
  } catch { /* handled */ }
}

async function handleRuleScore(student: Student, rule: ScoreRule) {
  try {
    await apiPost(`/api/v1/teacher/scores/give-by-rule/${rule.id}`, {
      student_id: student.id,
    })
    toast.show(`${student.name} 按规则「${rule.name}」${rule.is_penalty ? '扣' : '加'}${Math.abs(rule.points)}分`, 'success')
    student.total_score += rule.points
  } catch { /* handled */ }
}

// 快捷规则批量应用到全班：一次请求而非 N 次
async function handleBatchRuleScore(rule: ScoreRule) {
  if (students.value.length === 0) return
  try {
    await apiPost('/api/v1/teacher/scores/batch-give', {
      student_ids: students.value.map(s => s.id),
      points: rule.points,
      reason: rule.name,
    })
    toast.show(`全班按规则「${rule.name}」${rule.is_penalty ? '扣' : '加'}${Math.abs(rule.points)}分（${students.value.length}人）`, 'success')
    students.value.forEach(s => { s.total_score += rule.points })
  } catch { /* handled */ }
}

const positiveRules = computed(() => rules.value.filter(r => !r.is_penalty))
const negativeRules = computed(() => rules.value.filter(r => r.is_penalty))
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">积分管理</h2>
      <div style="display:flex;gap:8px;">
        <button class="btn btn-sm btn-primary">批量加分</button>
      </div>
    </div>

    <!-- 积分统计 -->
    <div class="stats-grid">
      <div class="stat-card stat-card--primary">
        <span class="stat-card__icon">⭐</span>
        <div class="stat-card__value">{{ scoreSummary.total.toLocaleString() }}</div>
        <div class="stat-card__label">累计积分</div>
      </div>
      <div class="stat-card stat-card--accent">
        <span class="stat-card__icon">📅</span>
        <div class="stat-card__value">{{ scoreSummary.today.toLocaleString() }}</div>
        <div class="stat-card__label">今日积分</div>
      </div>
      <div class="stat-card stat-card--secondary">
        <span class="stat-card__icon">📈</span>
        <div class="stat-card__value">{{ scoreSummary.this_week.toLocaleString() }}</div>
        <div class="stat-card__label">本周积分</div>
      </div>
      <div class="stat-card stat-card--info">
        <span class="stat-card__icon">📋</span>
        <div class="stat-card__value">{{ rules.length }}</div>
        <div class="stat-card__label">积分规则</div>
      </div>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else style="display:grid;grid-template-columns:1fr 1fr;gap:24px;">
      <!-- 快捷规则区 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">快捷规则</h3>

        <div v-if="positiveRules.length">
          <h4 style="font-size:13px;color:var(--color-text-secondary);margin-bottom:8px;">加分规则</h4>
          <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px;margin-bottom:16px;">
            <div v-for="rule in positiveRules" :key="rule.id" class="card"
              style="padding:10px 12px;border-color:rgba(16,185,129,0.3);cursor:pointer;text-align:center;font-size:13px;"
              @click="handleBatchRuleScore(rule)">
              <span style="font-weight:700;color:var(--color-accent);">+{{ rule.points }}</span>
              <span style="margin-left:6px;">{{ rule.name }}</span>
            </div>
          </div>
        </div>

        <div v-if="negativeRules.length">
          <h4 style="font-size:13px;color:var(--color-text-secondary);margin-bottom:8px;">减分规则</h4>
          <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px;">
            <div v-for="rule in negativeRules" :key="rule.id" class="card"
              style="padding:10px 12px;border-color:rgba(239,68,68,0.3);cursor:pointer;text-align:center;font-size:13px;"
              @click="handleBatchRuleScore(rule)">
              <span style="font-weight:700;color:var(--color-danger);">{{ rule.points }}</span>
              <span style="margin-left:6px;">{{ rule.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 手动加减分 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">手动加减分</h3>
        <div class="form-group">
          <label>选择学生</label>
          <select v-model.number="selectedStudentId" class="form-select">
            <option :value="null">请选择学生</option>
            <option v-for="s in students" :key="s.id" :value="s.id">
              {{ s.name }} — {{ s.total_score }}分
            </option>
          </select>
        </div>
        <div class="form-group">
          <label>积分值（正数加分，负数减分）</label>
          <input v-model.number="customPoints" type="number" class="form-input" placeholder="如：5 或 -3">
        </div>
        <div class="form-group">
          <label>理由</label>
          <input v-model="customReason" class="form-input" placeholder="如：作业完成" @keydown.enter="() => {
            const s = students.find(x => x.id === selectedStudentId)
            if (s) handleGiveScore(s, customPoints, customReason)
          }">
        </div>
        <button class="btn btn-primary" style="width:auto;" :disabled="!selectedStudentId"
          @click="() => {
            const s = students.find(x => x.id === selectedStudentId)
            if (s) handleGiveScore(s, customPoints, customReason)
          }">
          确认发放
        </button>
      </div>
    </div>
  </div>
</template>
