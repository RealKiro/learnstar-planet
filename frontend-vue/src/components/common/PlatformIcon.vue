<script setup lang="ts">
// 第三方平台官方品牌图标（简化 SVG，颜色取自各平台官方品牌色）
const props = defineProps<{
  platform: string
  size?: number
}>()
const size = props.size || 24

// 各平台官方品牌色
const colors: Record<string, string> = {
  wechat_work: '#0082EF',
  dingtalk: '#0089FF',
  feishu: '#3370FF',
  renren: '#F36C1D',
  wechat: '#07C160',
  qq: '#12B7F5',
}
const color = colors[props.platform] || '#7c3aed'
</script>

<template>
  <!-- 企业微信：蓝底白字"企" -->
  <svg v-if="platform === 'wechat_work'" :width="size" :height="size" viewBox="0 0 32 32" fill="none">
    <rect width="32" height="32" rx="6" fill="#0082EF"/>
    <path d="M10 8h12v4h-8v2h6v4h-6v6h-4V8z" fill="#fff"/>
  </svg>

  <!-- 钉钉：蓝色"钉"形 -->
  <svg v-else-if="platform === 'dingtalk'" :width="size" :height="size" viewBox="0 0 32 32" fill="none">
    <rect width="32" height="32" rx="6" fill="#0089FF"/>
    <path d="M23 6l-9 9 5-1-8 12 8-9-4 1 11-13h-3z" fill="#fff"/>
  </svg>

  <!-- 飞书：蓝白色羽毛 -->
  <svg v-else-if="platform === 'feishu'" :width="size" :height="size" viewBox="0 0 32 32" fill="none">
    <rect width="32" height="32" rx="6" fill="#3370FF"/>
    <path d="M8 24c3-6 8-10 16-12-1 6-4 10-9 13l-3-4h-4v3z" fill="#fff" opacity="0.95"/>
    <circle cx="9" cy="23" r="2" fill="#3370FF"/>
  </svg>

  <!-- 人人通空间：橙色渐变圆 -->
  <svg v-else-if="platform === 'renren'" :width="size" :height="size" viewBox="0 0 32 32" fill="none">
    <defs>
      <linearGradient id="renren-g" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0" stop-color="#F36C1D"/>
        <stop offset="1" stop-color="#F99A2E"/>
      </linearGradient>
    </defs>
    <circle cx="16" cy="16" r="14" fill="url(#renren-g)"/>
    <path d="M11 13c0-3 2-5 5-5s5 2 5 5c0 1.5-.5 2.8-1.5 3.8L16 24l-3.5-7.2C11.5 15.8 11 14.5 11 13z" fill="#fff"/>
  </svg>

  <!-- 微信：绿色对话气泡 -->
  <svg v-else-if="platform === 'wechat'" :width="size" :height="size" viewBox="0 0 32 32" fill="none">
    <rect width="32" height="32" rx="6" fill="#07C160"/>
    <path d="M23 12.5c0-4.1-3.4-7.5-7.5-7.5S8 8.4 8 12.5s3.4 7.5 7.5 7.5c.6 0 1.2-.1 1.8-.2l3.2 1.7-.9-2.7c1.9-1.2 3.4-3.2 3.4-5.3z" fill="#fff" opacity="0.9"/>
    <circle cx="13" cy="12" r="1.2" fill="#07C160"/>
    <circle cx="18" cy="12" r="1.2" fill="#07C160"/>
  </svg>

  <!-- QQ：企鹅 -->
  <svg v-else-if="platform === 'qq'" :width="size" :height="size" viewBox="0 0 32 32" fill="none">
    <rect width="32" height="32" rx="6" fill="#12B7F5"/>
    <path d="M16 5c3.5 0 6.5 2.5 7 6.5.3 2.5-.3 4.8-1.5 6.8.3 3-1 5.5-2.5 7-.8-.5-1.3-1.2-1.8-2-2 .7-4 .7-6 0-.5.8-1 1.5-1.8 2-1.5-1.5-2.8-4-2.5-7-1.2-2-1.8-4.3-1.5-6.8C9.5 7.5 12.5 5 16 5z" fill="#fff"/>
    <circle cx="12.5" cy="14" r="1" fill="#12B7F5"/>
    <circle cx="19.5" cy="14" r="1" fill="#12B7F5"/>
  </svg>

  <!-- 兜底：蓝色方块 + 平台首字母 -->
  <svg v-else :width="size" :height="size" viewBox="0 0 32 32" fill="none">
    <rect width="32" height="32" rx="6" fill="#7c3aed"/>
    <text x="16" y="22" text-anchor="middle" font-size="16" font-weight="bold" fill="#fff" font-family="sans-serif">{{ (platform || '?')[0] }}</text>
  </svg>
</template>
