<script setup lang="ts">
import { ref, computed, watch } from 'vue'

export interface ReasonGroup {
  title: string
  emoji: string
  color: string
  items: string[]
}

const props = defineProps<{
  show: boolean
  type: 'add' | 'sub'
  /** 单学生弹窗的学生名 */
  studentName?: string
  /** 批量弹窗：选中学生数 */
  selectedCount?: number
  stepValue: number
  busy: boolean
  error: string
  reasonGroups: ReasonGroup[]
}>()

const emit = defineEmits<{
  close: []
  confirm: [reasons: string[]]
  'go-login': []
}>()

// ===== 理由多选（弹窗内自持状态，打开时重置） =====
const selectedReasons = ref<string[]>([])
const customReason = ref('')
const showCustom = ref(false)

watch(() => props.show, (v) => { if (v) resetReasonPick() })

function toggleReason(r: string) {
  const i = selectedReasons.value.indexOf(r)
  if (i >= 0) selectedReasons.value.splice(i, 1)
  else selectedReasons.value.push(r)
}
function resetReasonPick() { selectedReasons.value = []; customReason.value = ''; showCustom.value = false }
function selectedList(): string[] {
  const list = [...selectedReasons.value]
  if (showCustom.value && customReason.value.trim()) list.push('✍️ ' + customReason.value.trim())
  return list
}

// ===== 权限预检：单次 ±30 需教师登录（内联提示 + 动态按钮） =====
const needsLogin = computed(() => Math.abs(props.stepValue * selectedList().length) >= 30)

function confirm() {
  if (needsLogin.value || props.busy) return
  emit('confirm', selectedList())
}
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="show" class="reason-overlay" @click.self="emit('close')">
        <div class="reason-box">
          <div class="reason-head">
            <h3>{{ type === 'add' ? '🌟 选择加分原因' : '⚠️ 选择减分原因' }}</h3>
            <span class="count-chip">已选: <strong>{{ selectedList().length }}</strong> 项</span>
          </div>
          <p class="reason-sub">
            <template v-if="selectedCount !== undefined">
              为 <strong>{{ selectedCount }}</strong> 名学生选择原因 · 每人每项 <strong class="gold">{{ stepValue }}</strong> 分，可多选<span class="cap-hint">· 每人单次上限 30 分</span>
            </template>
            <template v-else>
              为 <strong>{{ studentName }}</strong> 选择原因 · 每项 <strong class="gold">{{ stepValue }}</strong> 分，可多选<span class="cap-hint">· 单次上限 30 分</span>
            </template>
          </p>

          <div class="reason-groups">
            <div v-for="g in reasonGroups" :key="g.title">
              <div class="rg-title">
                <span>{{ g.emoji }} {{ g.title }}</span>
                <span class="rg-line"></span>
              </div>
              <div class="rg-btns">
                <button
                  v-for="r in g.items"
                  :key="r"
                  @click="toggleReason(r)"
                  class="reason-pick"
                  :style="selectedReasons.includes(r) ? { borderColor: g.color, background: g.color + '1A', color: g.color, fontWeight: 700 } : {}"
                >
                  {{ r }}<span v-if="selectedReasons.includes(r)"> ✓</span>
                </button>
              </div>
            </div>
          </div>

          <div v-if="showCustom" class="custom-wrap">
            <input v-model="customReason" type="text" placeholder="输入自定义原因..." maxlength="20" class="custom-input" />
          </div>
          <button @click="showCustom = !showCustom" class="custom-toggle">➕ 自定义原因</button>

          <!-- 权限受限：单次 ±30 需教师账号登录（内联提示 + 动态按钮，非弹窗） -->
          <div v-if="needsLogin" class="perm-prompt">
            <span class="perm-icon">🔒</span>
            <div class="perm-text">
              <div class="perm-title">单次{{ type === 'add' ? '加分' : '减分' }}超过 30 分</div>
              <div class="perm-desc">教室端为避免误操作设了单次上限，请登录教师账号后操作</div>
            </div>
            <button class="perm-btn" @click="emit('go-login')">🔑 登录教师账号</button>
          </div>
          <div v-else-if="error" class="reason-error">{{ error }}</div>

          <div class="reason-actions">
            <button @click="emit('close')" class="act-cancel">取消</button>
            <button
              v-if="!needsLogin"
              @click="confirm"
              :disabled="busy || !selectedList().length"
              class="act-confirm"
              :style="busy || !selectedList().length ? 'opacity:0.5;cursor:not-allowed;' : ''"
            >
              ✅ 确认{{ type === 'add' ? '加分' : '减分' }} ({{ selectedList().length }}项)
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.reason-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 300;
}
.reason-box {
  background: var(--color-bg-card);
  border: 1px solid var(--tint-3);
  border-radius: var(--md-radius);
  padding: 24px 28px;
  max-width: 480px;
  width: 90%;
  max-height: 82vh;
  overflow-y: auto;
  box-shadow: var(--md-elevation);
  animation: popIn 0.25s ease;
}
@keyframes popIn { from { transform: scale(0.92); opacity: 0; } to { transform: scale(1); opacity: 1; } }
.reason-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 4px; }
.reason-head h3 { font-size: 19px; font-weight: 700; margin: 0; }
.count-chip {
  font-size: 13px;
  font-weight: 700;
  color: var(--md-primary);
  background: rgba(167,139,250,0.12);
  padding: 3px 12px;
  border-radius: 20px;
}
.reason-sub { font-size: 13px; color: var(--md-text-secondary); margin-bottom: 16px; }
.reason-sub strong { color: var(--color-text); }
.reason-sub .gold { color: var(--md-gold); }
.cap-hint { color: var(--md-gold); opacity: 0.7; margin-left: 6px; }

.reason-groups { display: flex; flex-direction: column; gap: 12px; margin-bottom: 12px; }
.rg-title { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 600; color: var(--md-text-secondary); margin-bottom: 8px; }
.rg-line { flex: 1; height: 1px; background: var(--tint-3); }
.rg-btns { display: flex; flex-wrap: wrap; gap: 8px; }
.reason-pick {
  min-height: 44px;
  padding: 8px 16px;
  border-radius: 20px;
  border: 1px solid var(--tint-3);
  background: var(--tint-1);
  color: var(--color-text);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}
.reason-pick:hover { background: var(--tint-3); transform: translateY(-1px); }

.custom-wrap { margin-bottom: 10px; }
.custom-input {
  width: 100%;
  padding: 11px 16px;
  border-radius: 20px;
  border: 1px solid var(--tint-3);
  background: var(--tint-1);
  color: var(--color-text);
  font-size: 14px;
  outline: none;
  font-family: inherit;
}
.custom-toggle {
  padding: 6px 14px;
  border-radius: 20px;
  border: 1px dashed var(--tint-4);
  background: transparent;
  color: var(--md-text-secondary);
  font-size: 13px;
  cursor: pointer;
  font-family: inherit;
  margin-bottom: 6px;
}

/* 权限受限：单次 ±30 需教师登录（内联提示 + 动态按钮，非弹窗） */
.perm-prompt {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  margin-bottom: 12px;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.25);
  border-radius: 12px;
  flex-wrap: wrap;
  animation: popIn 0.2s ease;
}
.perm-icon { font-size: 20px; flex-shrink: 0; }
.perm-text { flex: 1; min-width: 0; }
.perm-title { font-size: 13px; font-weight: 700; color: var(--color-danger-text); }
.perm-desc { font-size: 12px; color: var(--md-text-secondary); margin-top: 2px; }
.perm-btn {
  padding: 8px 16px;
  border-radius: 20px;
  border: none;
  background: linear-gradient(135deg, var(--md-primary), var(--md-secondary));
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  transition: 0.15s;
  white-space: nowrap;
}
.perm-btn:hover { transform: translateY(-1px); box-shadow: 0 4px 14px rgba(0,0,0,0.2); }
.reason-error {
  margin-bottom: 12px;
  padding: 10px 12px;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.2);
  border-radius: 10px;
  color: var(--color-danger-text);
  font-size: 13px;
  text-align: center;
}

.reason-actions { display: flex; gap: 10px; margin-top: 12px; }
.act-cancel {
  flex: 1;
  padding: 12px;
  border-radius: 12px;
  border: 1px solid var(--tint-3);
  background: transparent;
  color: var(--md-text-secondary);
  font-size: 14px;
  cursor: pointer;
  font-family: inherit;
}
.act-confirm {
  flex: 1.4;
  padding: 12px;
  border-radius: 12px;
  border: none;
  background: var(--md-primary);
  color: #fff;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
