<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { apiGet, apiPost, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse } from '@/types'

const toast = useToastStore()

interface ShopItemExt {
  id: number
  name: string
  cost: number
  stock: number
  category: string
  description?: string
}

const items = ref<ShopItemExt[]>([])
const loading = ref(true)
const filterCategory = ref<string>('')
const showRedeemModal = ref(false)
const selectedItem = ref<ShopItemExt | null>(null)
const selectedStudentId = ref<number | null>(null)

const students = ref<Array<{ id: number; name: string; total_score: number; pet_level?: number }>>([])

const categories = ref([
  { key: '', label: '全部' },
  { key: 'points', label: '⭐ 积分充值' },
  { key: 'stationery', label: '✏️ 文具用品' },
  { key: 'food', label: '🍎 零食饮料' },
  { key: 'privilege', label: '🎁 特权奖励' },
  { key: 'activity', label: '🎪 活动参与' },
])
const showAddForm = ref(false)
const newItem = ref({ name: '', cost: 10, stock: 0, category: 'stationery', description: '' })
const itemErrors = reactive<Record<string, string>>({})
function iClr(f: string) { delete itemErrors[f] }
function iVld(field: string): boolean {
  if (field === 'name' && !newItem.value.name.trim()) { itemErrors.name = '请填写商品名称，如"铅笔"'; return false }
  if (field === 'cost' && (!newItem.value.cost || newItem.value.cost < 1)) { itemErrors.cost = '积分至少为 1'; return false }
  delete itemErrors[field]; return true
}

const addStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const redeemStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const deleteStatusMap = ref<Record<number, 'idle' | 'loading' | 'success' | 'error'>>({})

const filteredItems = computed(() => {
  if (!filterCategory.value) return items.value
  return items.value.filter(i => (i as any).category === filterCategory.value)
})

onMounted(async () => {
  try {
    const [itemRes, studentRes] = await Promise.all([
      apiGet<ApiResponse<ShopItemExt[]>>('/api/v1/teacher/shop/items'),
      apiGet<ApiResponse<Array<{ id: number; name: string; total_score: number }>>>('/api/v1/teacher/students?per_page=100'),
    ])
    items.value = (itemRes.data || []).map(i => ({
      ...i,
      category: (i as any).category || 'reward',
      description: (i as any).description || '',
    }))
    students.value = (studentRes.data || []).map(s => ({
      ...s,
      pet_level: (s as any).pet_level || Math.floor(Math.random() * 8) + 1,
    }))
  } catch {
    items.value = demoItems()
    students.value = demoStudents()
  }
  finally { loading.value = false }
})

function demoItems(): ShopItemExt[] {
  return [
    { id: 1, name: '课外图书', cost: 50, stock: 20, category: 'reading', description: '精选课外读物一本' },
    { id: 2, name: '科学实验套装', cost: 80, stock: 10, category: 'science', description: '趣味科学实验器材' },
    { id: 3, name: '免作业1天', cost: 30, stock: 0, category: 'reward', description: '特权：免写一次作业' },
    { id: 4, name: '班级之星徽章', cost: 100, stock: 5, category: 'reward', description: '荣誉徽章+表扬信' },
    { id: 5, name: '冰淇淋奖励', cost: 40, stock: 0, category: 'reward', description: '课间领取冰淇淋' },
    { id: 6, name: '迟到抵消券', cost: 60, stock: 15, category: 'penalty', description: '抵消一次迟到记录' },
  ]
}
function demoStudents() {
  return [
    { id: 1, name: '小明', total_score: 385, pet_level: 5 },
    { id: 2, name: '小红', total_score: 365, pet_level: 4 },
    { id: 3, name: '小刚', total_score: 350, pet_level: 3 },
    { id: 4, name: '小丽', total_score: 395, pet_level: 6 },
    { id: 5, name: '小华', total_score: 420, pet_level: 7 },
    { id: 6, name: '小强', total_score: 280, pet_level: 2 },
  ]
}

function openRedeem(item: ShopItemExt) {
  selectedItem.value = item
  showRedeemModal.value = true
  selectedStudentId.value = null
}

async function submitRedeem() {
  if (!selectedStudentId.value || !selectedItem.value) {
    toast.show('请选择学生和商品', 'error', { position: 'center', duration: 2000 })
    return
  }
  redeemStatus.value = 'loading'
  try {
    await apiPost('/api/v1/teacher/shop/redemptions', {
      student_id: selectedStudentId.value,
      item_id: selectedItem.value.id,
    })
    redeemStatus.value = 'success'
    setTimeout(() => { redeemStatus.value = 'idle'; showRedeemModal.value = false }, 1500)
  } catch {
    redeemStatus.value = 'error'
    setTimeout(() => { redeemStatus.value = 'idle' }, 3000)
  }
}

async function deleteItem(item: ShopItemExt) {
  if (!confirm(`确定删除商品「${item.name}」？`)) return
  deleteStatusMap.value[item.id] = 'loading'
  try {
    await apiDelete(`/api/v1/teacher/shop/items/${item.id}`)
    deleteStatusMap.value[item.id] = 'success'
    items.value = items.value.filter(i => i.id !== item.id)
    setTimeout(() => { deleteStatusMap.value[item.id] = 'idle' }, 1500)
  } catch {
    deleteStatusMap.value[item.id] = 'error'
    setTimeout(() => { deleteStatusMap.value[item.id] = 'idle' }, 3000)
  }
}

async function addItem() {
  if (!iVld('name') | !iVld('cost')) return
  addStatus.value = 'loading'
  try {
    await apiPost('/api/v1/teacher/shop/items', newItem.value)
    addStatus.value = 'success'
    newItem.value = { name: '', cost: 10, stock: 0, category: 'reward', description: '' }
    const res = await apiGet<ApiResponse<ShopItemExt[]>>('/api/v1/teacher/shop/items')
    items.value = (res.data || []).map(i => ({
      ...i,
      category: (i as any).category || 'reward',
      description: (i as any).description || '',
    }))
    setTimeout(() => { addStatus.value = 'idle'; showAddForm.value = false }, 1500)
  } catch {
    addStatus.value = 'error'
    setTimeout(() => { addStatus.value = 'idle' }, 3000)
  }
}

const catLabels: Record<string, string> = { points: '⭐ 积分充值', stationery: '✏️ 文具用品', food: '🍎 零食饮料', privilege: '🎁 特权奖励', activity: '🎪 活动参与', experience: '🎨 体验活动' }
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">积分商城</h2>
      <div style="display:flex;gap:8px;">
        <button class="btn btn-sm" style="background:var(--color-bg-card);color:var(--color-text);border:1px solid var(--color-border);" @click="showAddForm = !showAddForm">
          {{ showAddForm ? '取消' : '添加奖品' }}
        </button>
      </div>
    </div>

    <!-- 宠物成长绑定说明 -->
    <div class="card" style="margin-bottom:24px;padding:16px 24px;background:rgba(79,70,229,0.04);border-color:rgba(79,70,229,0.15);">
      <div style="display:flex;align-items:center;gap:12px;">
        <span style="font-size:28px;">🌟</span>
        <div>
          <div style="font-weight:600;font-size:14px;color:var(--color-primary);">宠物成长与积分联动</div>
          <div style="font-size:13px;color:var(--color-text-secondary);">
            班主任和科任老师增减积分时，学生宠物的成长值同步变化。积分越高，宠物成长越快，逐步进化为更高级形态。
          </div>
        </div>
      </div>
    </div>

    <!-- 分类筛选 -->
    <div style="display:flex;gap:8px;margin-bottom:16px;">
      <button v-for="c in categories" :key="c.key"
        :style="filterCategory === c.key ? { background:'var(--color-primary)', color:'#fff', borderColor:'var(--color-primary)' } : { color: 'var(--color-text)' }"
        style="padding:8px 16px;border-radius:20px;font-size:13px;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);"
        @click="filterCategory = c.key">
        {{ c.label }}
      </button>
    </div>

    <div v-if="showAddForm" class="card" style="margin-bottom:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">添加新奖品</h3>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;">
        <div class="form-group"><label>商品名称 <span style="color:#f87171;">*</span></label><input v-model="newItem.name" class="form-input" placeholder="如：铅笔" :style="{ borderColor: itemErrors.name ? '#f87171' : '' }" @blur="iVld('name')" @input="iClr('name')"><div v-if="itemErrors.name" style="color:#f87171;font-size:11px;margin-top:2px;">{{ itemErrors.name }}</div></div>
        <div class="form-group"><label>所需积分 <span style="color:#f87171;">*</span></label><input v-model.number="newItem.cost" type="number" min="1" class="form-input" :style="{ borderColor: itemErrors.cost ? '#f87171' : '' }" @blur="iVld('cost')" @input="iClr('cost')"><div v-if="itemErrors.cost" style="color:#f87171;font-size:11px;margin-top:2px;">{{ itemErrors.cost }}</div></div>
        <div class="form-group"><label>库存（0=无限）</label><input v-model.number="newItem.stock" type="number" min="0" class="form-input"></div>
        <div class="form-group">
          <label>类别</label>
          <select v-model="newItem.category" class="form-select">
            <option value="points">⭐ 积分充值</option>
            <option value="stationery">✏️ 文具用品</option>
            <option value="food">🍎 零食饮料</option>
            <option value="privilege">🎁 特权奖励</option>
            <option value="activity">🎪 活动参与</option>
          </select>
        </div>
      </div>
      <div class="form-group" style="margin-top:8px;">
        <label>描述</label>
        <input v-model="newItem.description" class="form-input" placeholder="简短描述">
      </div>
      <button class="btn btn-primary" style="margin-top:16px;width:auto;" :disabled="addStatus !== 'idle'"
        :style="{ background: addStatus === 'loading' ? '#f59e0b' : addStatus === 'success' ? '#10b981' : addStatus === 'error' ? '#ef4444' : '#7c3aed' }"
        @click="addItem">
        <template v-if="addStatus === 'idle'">添加</template>
        <template v-else-if="addStatus === 'loading'">添加中...</template>
        <template v-else-if="addStatus === 'success'">✅ 已添加</template>
        <template v-else-if="addStatus === 'error'">❌ 失败</template>
      </button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="filteredItems.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">🛍️</div>
      <p>{{ filterCategory ? '该类别暂无商品' : '暂无商品，点击「添加奖品」开始' }}</p>
    </div>

    <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:16px;">
      <div v-for="item in filteredItems" :key="item.id" class="card"
        style="text-align:center;padding:24px;cursor:pointer;position:relative;transition:transform 0.3s;"
        @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = 'translateY(-4px)'"
        @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = ''">
        <div style="position:absolute;top:8px;right:12px;font-size:14px;"
          :style="{
            cursor: (!deleteStatusMap[item.id] || deleteStatusMap[item.id] === 'idle') ? 'pointer' : 'default',
            color: deleteStatusMap[item.id] === 'loading' ? '#f59e0b' : deleteStatusMap[item.id] === 'success' ? '#10b981' : deleteStatusMap[item.id] === 'error' ? '#ef4444' : 'var(--color-text-secondary)'
          }"
          @click.stop="(!deleteStatusMap[item.id] || deleteStatusMap[item.id] === 'idle') && deleteItem(item)">
          <template v-if="!deleteStatusMap[item.id] || deleteStatusMap[item.id] === 'idle'">✕</template>
          <template v-else-if="deleteStatusMap[item.id] === 'loading'">⏳</template>
          <template v-else-if="deleteStatusMap[item.id] === 'success'">✅</template>
          <template v-else-if="deleteStatusMap[item.id] === 'error'">❌</template>
        </div>
        <div style="font-size:48px;margin-bottom:12px;">
          {{ { '课外图书': '📖', '科学实验套装': '🔬', '免作业1天': '🎮', '班级之星徽章': '🏆', '冰淇淋奖励': '🍦', '迟到抵消券': '⏰' }[item.name] || '🎁' }}
        </div>
        <div style="font-weight:600;font-size:16px;margin-bottom:4px;">{{ item.name }}</div>
        <div style="font-size:12px;color:var(--color-text-secondary);margin-bottom:8px;" v-if="item.description">{{ item.description }}</div>
        <div style="display:inline-flex;align-items:center;gap:4px;padding:2px 10px;border-radius:12px;font-size:12px;margin-bottom:8px;background:rgba(79,70,229,0.08);color:var(--color-primary);">
          {{ (catLabels as any)[item.category] || '🎁 奖励' }}
        </div>
        <div style="color:var(--color-secondary);font-weight:700;font-size:16px;">⭐ {{ item.cost }} 积分</div>
        <div v-if="item.stock > 0" style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">库存: {{ item.stock }}</div>
        <button class="btn btn-sm btn-primary" style="margin-top:12px;width:100%;" @click.stop="openRedeem(item)">立即兑换</button>
      </div>
    </div>

    <!-- 兑换弹窗 -->
    <div v-if="showRedeemModal" style="position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:200;">
      <div class="card" style="width:420px;max-width:90vw;max-height:80vh;overflow-y:auto;">
        <h3 style="font-size:18px;font-weight:600;margin-bottom:16px;">兑换 {{ selectedItem?.name }}</h3>
        <div class="form-group">
          <label>选择学生</label>
          <select v-model.number="selectedStudentId" class="form-select">
            <option :value="null">请选择学生</option>
            <option v-for="s in students" :key="s.id" :value="s.id" :disabled="s.total_score < (selectedItem?.cost || 0)">
              {{ s.name }} — {{ s.total_score }}分 {{ s.total_score < (selectedItem?.cost || 0) ? '(积分不足)' : '' }}
            </option>
          </select>
        </div>
        <div v-if="selectedStudentId" style="margin-bottom:16px;padding:12px;background:var(--color-bg);border-radius:var(--radius-md);">
          <div style="font-size:13px;">
            当前积分: <strong>{{ students.find(s => s.id === selectedStudentId)?.total_score || 0 }}</strong>
            → 兑换后: <strong style="color:var(--color-accent);">{{ (students.find(s => s.id === selectedStudentId)?.total_score || 0) - (selectedItem?.cost || 0) }}</strong>
          </div>
          <div style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">
            宠物等级: Lv.{{ students.find(s => s.id === selectedStudentId)?.pet_level || 0 }}
          </div>
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end;">
          <button class="btn btn-ghost" @click="showRedeemModal = false">取消</button>
          <button class="btn btn-primary" :disabled="redeemStatus !== 'idle' || !selectedStudentId"
            :style="{ background: redeemStatus === 'loading' ? '#f59e0b' : redeemStatus === 'success' ? '#10b981' : redeemStatus === 'error' ? '#ef4444' : '#7c3aed' }"
            @click="submitRedeem">
            <template v-if="redeemStatus === 'idle'">确认兑换</template>
            <template v-else-if="redeemStatus === 'loading'">兑换中...</template>
            <template v-else-if="redeemStatus === 'success'">✅ 已兑换</template>
            <template v-else-if="redeemStatus === 'error'">❌ 失败</template>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cat-btn { padding:8px 16px;border-radius:20px;font-size:13px;cursor:pointer;background:var(--color-bg);border:1px solid var(--color-border);color:var(--color-text);font-family:inherit;transition:all 0.15s; }
.cat-btn:hover { border-color:var(--color-primary);color:var(--color-primary); }
.cat-btn.active { background:var(--color-primary);color:#fff;border-color:var(--color-primary); }
</style>
