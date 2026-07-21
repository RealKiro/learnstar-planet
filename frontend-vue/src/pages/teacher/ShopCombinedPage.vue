<script setup lang="ts">
import { ref } from 'vue'
import ShopPage from './ShopPage.vue'
import ExchangeCenterPage from './ExchangeCenterPage.vue'

const activeTab = ref<'shop' | 'exchange'>('shop')
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">班级经济</p>
        <h2 style="font-size:24px;font-weight:700;">🛍️ 积分商城</h2>
      </div>
    </div>

    <div class="tab-bar">
      <button :class="['tab-btn', { active: activeTab === 'shop' }]" @click="activeTab = 'shop'">
        📦 商品管理
      </button>
      <button :class="['tab-btn', { active: activeTab === 'exchange' }]" @click="activeTab = 'exchange'">
        🔄 兑换管理
      </button>
    </div>

    <ShopPage v-if="activeTab === 'shop'" />
    <ExchangeCenterPage v-if="activeTab === 'exchange'" />
  </div>
</template>

<style scoped>
.tab-bar { display: flex; gap: 4px; margin-bottom: 20px; background: var(--color-bg); border-radius: 12px; padding: 4px; }
.tab-btn { flex: 1; padding: 10px 12px; border: none; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer; background: transparent; color: var(--color-text-secondary); transition: all 0.2s; }
.tab-btn:hover { background: rgba(124,58,237,0.06); color: var(--color-text); }
.tab-btn.active { background: #7c3aed; color: #fff; box-shadow: 0 2px 8px rgba(124,58,237,0.25); }
</style>
