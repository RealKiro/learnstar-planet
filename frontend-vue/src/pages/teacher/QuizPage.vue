<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, Quiz, QuestionBank } from '@/types'

const toast = useToastStore()

const quizzes = ref<Quiz[]>([])
const banks = ref<QuestionBank[]>([])
const loading = ref(true)

const showQuizModal = ref(false)
const showBankModal = ref(false)
const sendingQuiz = ref(false)
const sendingBank = ref(false)

const quizForm = ref({ title: '', subject: '', question_count: 5, time_limit: 30 })
const bankForm = ref({ name: '', subject: '' })

onMounted(async () => {
  try {
    const [quizRes, bankRes] = await Promise.all([
      apiGet<ApiResponse<Quiz[]>>('/api/v1/teacher/quizzes'),
      apiGet<ApiResponse<QuestionBank[]>>('/api/v1/teacher/question-banks'),
    ])
    quizzes.value = quizRes.data || []
    banks.value = bankRes.data || []
  } catch { /* handled */ }
  finally { loading.value = false }
})

function openQuizModal() {
  quizForm.value = { title: '', subject: '', question_count: 5, time_limit: 30 }
  showQuizModal.value = true
}

function openBankModal() {
  bankForm.value = { name: '', subject: '' }
  showBankModal.value = true
}

async function handleCreateQuiz() {
  if (!quizForm.value.title.trim()) {
    toast.show('请填写测验标题', 'error')
    return
  }
  sendingQuiz.value = true
  try {
    await apiPost('/api/v1/teacher/quizzes', {
      title: quizForm.value.title.trim(),
      subject: quizForm.value.subject.trim(),
      question_count: quizForm.value.question_count,
      time_limit: quizForm.value.time_limit,
    })
    toast.show('测验已创建', 'success')
    showQuizModal.value = false
    const res = await apiGet<ApiResponse<Quiz[]>>('/api/v1/teacher/quizzes')
    quizzes.value = res.data || []
  } catch { /* handled */ }
  finally { sendingQuiz.value = false }
}

async function handleCreateBank() {
  if (!bankForm.value.name.trim()) {
    toast.show('请填写题库名称', 'error')
    return
  }
  sendingBank.value = true
  try {
    await apiPost('/api/v1/teacher/question-banks', {
      name: bankForm.value.name.trim(),
      subject: bankForm.value.subject.trim(),
    })
    toast.show('题库已创建', 'success')
    showBankModal.value = false
    const res = await apiGet<ApiResponse<QuestionBank[]>>('/api/v1/teacher/question-banks')
    banks.value = res.data || []
  } catch { /* handled */ }
  finally { sendingBank.value = false }
}

async function toggleQuiz(quiz: Quiz) {
  try {
    await apiPost(`/api/v1/teacher/quizzes/${quiz.id}/${quiz.is_active ? 'stop' : 'start'}`)
    quiz.is_active = !quiz.is_active
    toast.show(quiz.is_active ? '测验已开始' : '测验已停止', 'success')
  } catch { /* handled */ }
}

function progressOf(q: Quiz): number {
  if (!q.total_students) return 0
  return Math.round(q.submission_count / q.total_students * 100)
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">在线答题</h2>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <template v-else>
      <!-- 测验列表 -->
      <div style="margin-bottom:32px;">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px;">
          <h3 style="font-size:16px;font-weight:600;">测验列表</h3>
          <button class="btn btn-sm btn-primary" @click="openQuizModal">创建测验</button>
        </div>

        <div v-if="quizzes.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
          <div style="font-size:48px;margin-bottom:12px;">🎯</div>
          <p style="margin-bottom:16px;">在线答题功能即将上线</p>
          <button class="btn btn-primary" style="width:auto;" @click="openQuizModal">创建测验</button>
        </div>

        <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(320px,1fr));gap:16px;">
          <div v-for="q in quizzes" :key="q.id" class="card" style="padding:16px;">
            <div style="display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:12px;">
              <div>
                <h4 style="font-size:15px;font-weight:600;">{{ q.title }}</h4>
                <span style="font-size:12px;color:var(--color-text-secondary);">{{ q.subject }}</span>
              </div>
              <span :style="{
                fontSize: '12px',
                padding: '2px 8px',
                borderRadius: 'var(--radius-sm)',
                background: q.is_active ? 'rgba(16,185,129,0.1)' : 'var(--color-bg)',
                color: q.is_active ? 'var(--color-accent)' : 'var(--color-text-secondary)',
              }">{{ q.is_active ? '进行中' : '未开始' }}</span>
            </div>

            <div style="display:flex;gap:16px;font-size:13px;color:var(--color-text-secondary);margin-bottom:12px;">
              <span>题目：{{ q.question_count }}</span>
              <span v-if="q.time_limit">限时：{{ q.time_limit }}分钟</span>
            </div>

            <div style="margin-bottom:12px;">
              <div style="display:flex;justify-content:space-between;font-size:12px;margin-bottom:4px;">
                <span style="color:var(--color-text-secondary);">提交进度</span>
                <span style="font-weight:600;">{{ q.submission_count }}/{{ q.total_students }}</span>
              </div>
              <div style="height:8px;background:var(--color-bg);border-radius:4px;overflow:hidden;">
                <div :style="{ width: `${progressOf(q)}%`, height: '100%', background: 'var(--color-primary)' }"></div>
              </div>
            </div>

            <button class="btn btn-sm" :class="q.is_active ? 'btn-ghost' : 'btn-primary'"
              style="width:auto;" @click="toggleQuiz(q)">
              {{ q.is_active ? '停止' : '开始' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 题库列表 -->
      <div>
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px;">
          <h3 style="font-size:16px;font-weight:600;">题库列表</h3>
          <button class="btn btn-sm btn-primary" @click="openBankModal">创建题库</button>
        </div>

        <div v-if="banks.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
          <div style="font-size:48px;margin-bottom:12px;">📚</div>
          <p>暂无题库，点击「创建题库」添加</p>
        </div>

        <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:16px;">
          <div v-for="b in banks" :key="b.id" class="card" style="padding:16px;">
            <div style="display:flex;align-items:center;gap:12px;margin-bottom:8px;">
              <div style="width:40px;height:40px;border-radius:var(--radius-sm);background:var(--gradient-primary);display:flex;align-items:center;justify-content:center;font-size:20px;">📚</div>
              <div style="flex:1;min-width:0;">
                <h4 style="font-size:14px;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">{{ b.name }}</h4>
                <span style="font-size:12px;color:var(--color-text-secondary);">{{ b.subject }}</span>
              </div>
            </div>
            <div style="display:flex;gap:16px;font-size:12px;color:var(--color-text-secondary);">
              <span>题目：{{ b.question_count }}</span>
              <span>使用：{{ b.usage_count }}次</span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 创建测验弹窗 -->
    <div v-if="showQuizModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:100;" @click.self="showQuizModal = false">
      <div class="card" style="width:440px;max-width:90vw;padding:24px;">
        <h3 style="font-size:18px;font-weight:600;margin-bottom:20px;">创建测验</h3>
        <div class="form-group">
          <label>测验标题</label>
          <input v-model="quizForm.title" class="form-input" placeholder="如：第一单元测验">
        </div>
        <div class="form-group">
          <label>科目</label>
          <input v-model="quizForm.subject" class="form-input" placeholder="如：数学">
        </div>
        <div class="form-group">
          <label>题目数量</label>
          <input v-model.number="quizForm.question_count" type="number" min="1" class="form-input">
        </div>
        <div class="form-group">
          <label>限时（分钟）</label>
          <input v-model.number="quizForm.time_limit" type="number" min="1" class="form-input">
        </div>
        <div style="display:flex;justify-content:flex-end;gap:8px;margin-top:8px;">
          <button class="btn btn-ghost" @click="showQuizModal = false">取消</button>
          <button class="btn btn-primary" style="width:auto;" :disabled="sendingQuiz" @click="handleCreateQuiz">
            {{ sendingQuiz ? '创建中...' : '创建' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 创建题库弹窗 -->
    <div v-if="showBankModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:100;" @click.self="showBankModal = false">
      <div class="card" style="width:400px;max-width:90vw;padding:24px;">
        <h3 style="font-size:18px;font-weight:600;margin-bottom:20px;">创建题库</h3>
        <div class="form-group">
          <label>题库名称</label>
          <input v-model="bankForm.name" class="form-input" placeholder="如：数学应用题库">
        </div>
        <div class="form-group">
          <label>科目</label>
          <input v-model="bankForm.subject" class="form-input" placeholder="如：数学">
        </div>
        <div style="display:flex;justify-content:flex-end;gap:8px;margin-top:8px;">
          <button class="btn btn-ghost" @click="showBankModal = false">取消</button>
          <button class="btn btn-primary" style="width:auto;" :disabled="sendingBank" @click="handleCreateBank">
            {{ sendingBank ? '创建中...' : '创建' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
