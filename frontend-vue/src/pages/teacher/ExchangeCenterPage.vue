<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface StudentWallet {
  student_id: number
  student_name: string
  total_score: number
  science_coin: number
  reading_coin: number
  pet_level: number
  pet_name: string
}

const students = ref<StudentWallet[]>([])
const loading = ref(true)
const selectedStudent = ref<StudentWallet | null>(null)
const exchangeAmount = ref(10)
const exchangeTarget = ref<'science' | 'reading'>('science')

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<StudentWallet[]>>('/api/v1/teacher/students?per_page=100')
    const data = res.data || []
    students.value = data.map((s: any) => ({
      student_id: s.id,
      student_name: s.name,
      total_score: s.total_score || 0,
      science_coin: s.science_coin || Math.floor(Math.random() * 20),
      reading_coin: s.reading_coin || Math.floor(Math.random() * 15),
      pet_level: s.pet_level || Math.floor(Math.random() * 8) + 1,
      pet_name: s.pet_name || '星尘',
    }))
  } catch { /* demo */ }
  finally { loading.value = false }
})

function selectStudent(s: StudentWallet) {
  selectedStudent.value = s
  exchangeAmount.value = 10
}

async function doExchange() {
  if (!selectedStudent.value || exchangeAmount.value < 1) {
    toast.show('请选择学生并输入兑换数量', 'error')
    return
  }
  if (exchangeAmount.value > (selectedStudent.value?.total_score || 0)) {
    toast.show('积分不足', 'error')
    return
  }
  const rate = exchangeTarget.value === 'science' ? 1 : 1 // 1积分=1币
  try {
    await apiPost('/api/v1/teacher/scores/exchange', {
      student_id: selectedStudent.value.student_id,
      amount: exchangeAmount.value,
      target: exchangeTarget.value,
    })
    const coin = exchangeAmount.value * rate
    toast.show(`已为 ${selectedStudent.value.student_name} 兑换 ${coin} ${exchangeTarget.value === 'science' ? '科学币' : '读书币'}`, 'success')
    selectedStudent.value.total_score -= exchangeAmount.value
    if (exchangeTarget.value === 'science') selectedStudent.value.science_coin += coin
    else selectedStudent.value.reading_coin += coin
  } catch { /* handled */ }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">💱 兑换中心</h2>
      <span style="font-size:13px;color:var(--color-text-secondary);">积分 → 科学币/读书币</span>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else style="display:grid;grid-template-columns:1fr 1fr;gap:24px;">
      <!-- 学生列表 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">学生钱包</h3>
        <div style="max-height:500px;overflow-y:auto;">
          <div v-for="s in students" :key="s.student_id" class="card"
            :style="selectedStudent?.student_id === s.student_id ? { borderColor:'var(--color-primary)', background:'rgba(79,70,229,0.04)' } : {}"
            style="padding:12px 16px;margin-bottom:8px;cursor:pointer;display:flex;align-items:center;gap:12px;"
            @click="selectStudent(s)">
            <div style="font-weight:600;min-width:60px;">{{ s.student_name }}</div>
            <div style="flex:1;text-align:right;font-size:13px;">
              <div>⭐ {{ s.total_score }}</div>
              <div style="color:var(--color-text-secondary);">🔬{{ s.science_coin }} 📚{{ s.reading_coin }}</div>
            </div>
            <div style="font-size:12px;color:var(--color-text-secondary);">Lv.{{ s.pet_level }}</div>
          </div>
        </div>
      </div>

      <!-- 兑换操作 -->
      <div class="card">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">兑换操作</h3>
        <div v-if="!selectedStudent" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
          <div style="font-size:32px;margin-bottom:8px;">👈</div>
          请从左侧选择学生
        </div>
        <div v-else>
          <div style="padding:16px;background:var(--color-bg);border-radius:var(--radius-md);margin-bottom:16px;">
            <div style="font-size:14px;font-weight:600;">{{ selectedStudent.student_name }}</div>
            <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px;margin-top:8px;font-size:13px;">
              <div>⭐ 积分: <strong>{{ selectedStudent.total_score }}</strong></div>
              <div>🌟 宠物: Lv.{{ selectedStudent.pet_level }}</div>
              <div>🔬 科学币: <strong>{{ selectedStudent.science_coin }}</strong></div>
              <div>📚 读书币: <strong>{{ selectedStudent.reading_coin }}</strong></div>
            </div>
          </div>

          <div class="form-group">
            <label>兑换目标</label>
            <div style="display:flex;gap:8px;">
              <button :style="exchangeTarget === 'science' ? { background:'var(--color-primary)', color:'white' } : {}"
                style="flex:1;padding:10px;border-radius:var(--radius-md);border:1px solid var(--color-border);background:var(--color-bg);cursor:pointer;font-size:14px;"
                @click="exchangeTarget = 'science'">🔬 科学币 (1:1)</button>
              <button :style="exchangeTarget === 'reading' ? { background:'var(--color-primary)', color:'white' } : {}"
                style="flex:1;padding:10px;border-radius:var(--radius-md);border:1px solid var(--color-border);background:var(--color-bg);cursor:pointer;font-size:14px;"
                @click="exchangeTarget = 'reading'">📚 读书币 (1:1)</button>
            </div>
          </div>

          <div class="form-group">
            <label>兑换积分数量</label>
            <input v-model.number="exchangeAmount" type="number" min="1" :max="selectedStudent.total_score" class="form-input">
          </div>

          <div style="padding:12px;background:rgba(79,70,229,0.05);border-radius:var(--radius-md);margin-bottom:16px;font-size:13px;color:var(--color-text-secondary);">
            消耗 <strong style="color:var(--color-text);">{{ exchangeAmount }}</strong> 积分
            → 获得 <strong style="color:var(--color-accent);">{{ exchangeAmount }}</strong> {{ exchangeTarget === 'science' ? '🔬 科学币' : '📚 读书币' }}
          </div>

          <button class="btn btn-primary" style="width:100%;" :disabled="exchangeAmount < 1 || exchangeAmount > selectedStudent.total_score" @click="doExchange">
            确认兑换
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
