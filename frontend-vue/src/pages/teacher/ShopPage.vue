<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiGet, apiPost, apiDelete } from '@/utils/api'
import { useToastStore } from '@/stores/toast'
import type { ApiResponse, ShopItem } from '@/types'

const toast = useToastStore()

const items = ref<ShopItem[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiGet<ApiResponse<ShopItem[]>>('/api/v1/teacher/shop/items')
    items.value = res.data || []
  } catch { items.value = [] }
  finally { loading.value = false }
})

async function deleteItem(item: ShopItem) {
  if (!confirm(`确定删除商品「${item.name}」？`)) return
  try {
    await apiDelete(`/api/v1/teacher/shop/items/${item.id}`)
    items.value = items.value.filter(i => i.id !== item.id)
    toast.show(`已删除 ${item.name}`, 'success')
  } catch { /* handled */ }
}

const showAddForm = ref(false)
const newItem = ref({ name: '', cost: 10, stock: 0 })

async function addItem() {
  if (!newItem.value.name || newItem.value.cost < 1) {
    toast.show('请填写商品名称和价格', 'error')
    return
  }
  try {
    await apiPost('/api/v1/teacher/shop/items', newItem.value)
    toast.show('商品已添加', 'success')
    showAddForm.value = false
    newItem.value = { name: '', cost: 10, stock: 0 }
    // Reload
    const res = await apiGet<ApiResponse<ShopItem[]>>('/api/v1/teacher/shop/items')
    items.value = res.data || []
  } catch { /* handled */ }
}
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <h2 style="font-size:24px;font-weight:700;">积分商城</h2>
      <button class="btn btn-sm btn-primary" @click="showAddForm = !showAddForm">
        {{ showAddForm ? '取消' : '添加奖品' }}
      </button>
    </div>

    <div v-if="showAddForm" class="card" style="margin-bottom:24px;">
      <h3 style="font-size:16px;font-weight:600;margin-bottom:16px;">添加新奖品</h3>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;">
        <div class="form-group">
          <label>商品名称</label>
          <input v-model="newItem.name" class="form-input" placeholder="如：课外图书">
        </div>
        <div class="form-group">
          <label>所需积分</label>
          <input v-model.number="newItem.cost" type="number" min="1" class="form-input">
        </div>
        <div class="form-group">
          <label>库存（0=无限）</label>
          <input v-model.number="newItem.stock" type="number" min="0" class="form-input">
        </div>
      </div>
      <button class="btn btn-primary" style="margin-top:16px;width:auto;" @click="addItem">添加</button>
    </div>

    <div v-if="loading" style="text-align:center;padding:48px;color:var(--color-text-secondary);">加载中...</div>

    <div v-else-if="items.length === 0" class="card" style="text-align:center;padding:48px;color:var(--color-text-secondary);">
      <div style="font-size:48px;margin-bottom:8px;">🛍️</div>
      <p>暂无商品，点击「添加奖品」开始</p>
    </div>

    <div v-else style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:16px;">
      <div v-for="item in items" :key="item.id" class="card"
        style="text-align:center;padding:24px;cursor:pointer;position:relative;transition:transform 0.3s;"
        @mouseenter="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = 'translateY(-4px)'"
        @mouseleave="(e: MouseEvent) => (e.currentTarget as HTMLElement).style.transform = ''">
        <div style="position:absolute;top:8px;right:12px;cursor:pointer;font-size:14px;color:var(--color-text-secondary);"
          @click.stop="deleteItem(item)">✕</div>
        <div style="font-size:48px;margin-bottom:12px;">
          {{ { '课外图书': '📖', '美术用品': '🎨', '免作业1天': '🎮', '班级之星徽章': '🏆', '冰淇淋奖励': '🍦', '电影票': '🎬' }[item.name] || '🎁' }}
        </div>
        <div style="font-weight:600;font-size:16px;margin-bottom:8px;">{{ item.name }}</div>
        <div style="color:var(--color-secondary);font-weight:700;font-size:16px;">
          ⭐ {{ item.cost }} 积分
        </div>
        <div v-if="item.stock > 0" style="font-size:12px;color:var(--color-text-secondary);margin-top:4px;">
          库存: {{ item.stock }}
        </div>
      </div>
    </div>
  </div>
</template>
