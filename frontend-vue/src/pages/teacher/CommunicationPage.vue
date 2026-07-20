<script setup lang="ts">
import { ref } from 'vue'
import BroadcastPage from './BroadcastPage.vue'
import NoticesPage from './NoticesPage.vue'

const activeTab = ref<'broadcast' | 'notice'>('broadcast')
</script>

<template>
  <div>
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:24px;">
      <div>
        <p style="font-size:13px;color:var(--color-text-secondary);margin-bottom:4px;">消息中心</p>
        <h2 style="font-size:24px;font-weight:700;">📢 消息中心</h2>
      </div>
    </div>

    <!-- 标签导航 -->
    <div class="tab-bar">
      <button :class="['tab-btn', { active: activeTab === 'broadcast' }]" @click="activeTab = 'broadcast'">
        📡 实时广播
      </button>
      <button :class="['tab-btn', { active: activeTab === 'notice' }]" @click="activeTab = 'notice'">
        📋 班级通知
      </button>
    </div>

    <BroadcastPage v-if="activeTab === 'broadcast'" />
    <NoticesPage v-if="activeTab === 'notice'" />
  </div>
</template>

<style scoped>
.tab-bar { display: flex; gap: 4px; margin-bottom: 20px; background: var(--color-bg); border-radius: 12px; padding: 4px; }
.tab-btn { flex: 1; padding: 10px 12px; border: none; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer; background: transparent; color: var(--color-text-secondary); transition: all 0.2s; }
.tab-btn:hover { background: rgba(124,58,237,0.06); color: var(--color-text); }
.tab-btn.active { background: #7c3aed; color: #fff; box-shadow: 0 2px 8px rgba(124,58,237,0.25); }
</style>
