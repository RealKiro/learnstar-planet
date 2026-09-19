<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BroadcastPage from './BroadcastPage.vue'
import NoticesPage from './NoticesPage.vue'

type Tab = 'broadcast' | 'notice'

const route = useRoute()
const router = useRouter()

/** tab 以 URL 为唯一真相：/teacher/communication?tab=notice 可深链、可分享，
    也让原独立路由 /teacher/broadcast、/teacher/notices 的重定向有的放矢。 */
const activeTab = ref<Tab>(route.query.tab === 'notice' ? 'notice' : 'broadcast')

function selectTab(tab: Tab) {
  if (tab === activeTab.value) return
  activeTab.value = tab
  router.replace({ query: { ...route.query, tab } })
}
</script>

<template>
  <div>
    <div class="page-head">
      <div>
        <p class="page-eyebrow">消息中心</p>
        <h2 class="page-title">📢 消息中心</h2>
      </div>
    </div>

    <!-- 标签导航 -->
    <div class="tab-bar">
      <button :class="['tab-btn', { active: activeTab === 'broadcast' }]" @click="selectTab('broadcast')">
        📡 实时广播
      </button>
      <button :class="['tab-btn', { active: activeTab === 'notice' }]" @click="selectTab('notice')">
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
