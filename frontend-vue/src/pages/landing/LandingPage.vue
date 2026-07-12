<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

if (authStore.isLoggedIn) {
  if (authStore.isAdmin) router.replace({ name: 'admin-dashboard' })
  else if (authStore.isTeacher) router.replace({ name: 'teacher-dashboard' })
  else if (authStore.isParent) router.replace({ name: 'parent-home' })
}

// ── 滚动揭示动画 ──
const observedEls = ref<Set<HTMLElement>>(new Set())
let observer: IntersectionObserver | null = null

function setRevealRef(el: HTMLElement | null) {
  if (el && !observedEls.value.has(el) && observer) {
    observer.observe(el)
    observedEls.value.add(el)
  }
}

onMounted(() => {
  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          entry.target.classList.add('revealed')
          observer!.unobserve(entry.target)
          observedEls.value.delete(entry.target as HTMLElement)
        }
      }
    },
    { threshold: 0.12, rootMargin: '0px 0px -30px 0px' },
  )
})

onUnmounted(() => observer?.disconnect())
</script>

<template>
  <div class="page">
    <!-- ===== 导航栏 — 磨砂玻璃 ===== -->
    <header class="nav">
      <div class="nav-inner">
        <a href="/welcome" class="nav-brand" aria-label="学趣星球">
          🌌&nbsp;学趣星球
        </a>
        <nav class="nav-links">
          <a href="#features">功能</a>
          <a href="#evolution">宠物进化</a>
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="nav-github">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
            GitHub
          </a>
          <router-link to="/login" class="nav-cta">登录</router-link>
        </nav>
      </div>
    </header>

    <!-- ===== Hero ===== -->
    <section class="hero">
      <div class="hero-inner reveal" :ref="setRevealRef">
        <span class="badge">
          <span class="badge-dot" /> MIT 开源 · 完全免费 · 自托管
        </span>

        <h1 class="hero-title">
          让每个孩子的努力<br />
          <span class="hero-title-grad">都被看见</span>
        </h1>

        <p class="hero-desc">
          积分激励 · 宠物养成 · AI 助教 · 多端同步<br />
          开源班级管理系统，Docker 一键部署，数据完全自主掌控
        </p>

        <div class="hero-actions">
          <router-link to="/login" class="btn-primary">立即使用</router-link>
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="btn-secondary">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
            源码
          </a>
        </div>

        <div class="hero-tags">
          <div v-for="item in highlights" :key="item.label" class="tag">
            <span class="tag-icon">{{ item.icon }}</span>
            <span class="tag-label">{{ item.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== 功能模块 — Bento 卡片网格 ===== -->
    <section id="features" class="section-gray">
      <div class="container">
        <div class="section-head reveal" :ref="setRevealRef">
          <span class="section-eyebrow">功能模块</span>
          <h2 class="section-title">12 大模块，覆盖全场景</h2>
          <p class="section-desc">全部免费，无任何付费功能，Docker 一条命令即可部署</p>
        </div>

        <div class="bento">
          <div
            v-for="feat in features"
            :key="feat.title"
            class="bento-card reveal"
            :class="feat.size"
            :ref="setRevealRef"
          >
            <span class="bento-icon">{{ feat.icon }}</span>
            <h3 class="bento-title">{{ feat.title }}</h3>
            <p class="bento-desc">{{ feat.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== 宠物进化 — 水平时间线 ===== -->
    <section id="evolution" class="section-white">
      <div class="container">
        <div class="section-head reveal" :ref="setRevealRef">
          <span class="section-eyebrow">宠物系统</span>
          <h2 class="section-title">11 阶进化，积分驱动成长</h2>
          <p class="section-desc">从星尘到银河，每一次进步都有回响</p>
        </div>

        <div class="evo-track reveal" :ref="setRevealRef">
          <template v-for="(stage, i) in evolutionStages" :key="stage.name">
            <span v-if="i > 0" class="evo-arrow">→</span>
            <div class="evo-node">
              <span class="evo-emoji">{{ stage.emoji }}</span>
              <span class="evo-name">{{ stage.name }}</span>
            </div>
          </template>
        </div>
      </div>
    </section>

    <!-- ===== CTA ===== -->
    <section class="section-gray">
      <div class="container">
        <div class="cta reveal" :ref="setRevealRef">
          <div class="cta-icon">🚀</div>
          <h2 class="cta-title">加入开源社区</h2>
          <p class="cta-desc">
            项目完全开源，欢迎 Star ⭐、Fork、提交 Issue 和 PR<br />
            如果你正在用它管理班级，你就是这个社区的一员
          </p>
          <div class="cta-btns">
            <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="btn-outline">⭐ 给个 Star</a>
            <a href="https://github.com/RealKiro/learnstar-planet/issues" target="_blank" class="btn-ghost">提交 Issue</a>
          </div>
        </div>
      </div>
    </section>

    <!-- ===== 底部 ===== -->
    <footer class="footer">
      <p class="footer-line">
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" class="footer-link">学趣星球</a>
        &nbsp;· MIT 开源许可证 · 自托管 · 完全免费
      </p>
      <p class="footer-made">Made with 💜 by open source community</p>
    </footer>
  </div>
</template>

<script lang="ts">
const highlights = [
  { icon: '🆓', label: '完全免费' },
  { icon: '🔓', label: 'MIT 开源' },
  { icon: '🏠', label: '自托管部署' },
  { icon: '🗄️', label: '4 种数据库' },
  { icon: '🐳', label: 'Docker 一键部署' },
  { icon: '📱', label: '多端适配' },
]

const features = [
  { icon: '⭐', title: '积分激励', desc: '自定义规则，实时加减分，进步看得见', size: 'bento-lg' },
  { icon: '🌟', title: '宠物进化', desc: '11 阶宇宙进化体系，积分变经验驱动成长', size: '' },
  { icon: '🏆', title: '排行竞技', desc: '总积分 / 周进步 / 宠物等级三大排行', size: '' },
  { icon: '📢', title: '班级通知', desc: '一键发布，实时推送家长端', size: 'bento-wide' },
  { icon: '📊', title: '成绩管理', desc: '成绩录入分析，班级对比，趋势可视化', size: '' },
  { icon: '🤖', title: 'AI 助教', desc: '生成班级反馈、学生分析、家校沟通建议', size: 'bento-tall' },
  { icon: '✅', title: '智能考勤', desc: '一键点名签到，到课 / 请假 / 迟到状态统计', size: '' },
  { icon: '📱', title: '扫码收作业', desc: '生成专属二维码，学生扫码提交，自动汇总', size: 'bento-wide' },
  { icon: '🛍️', title: '积分商城', desc: '学生兑换实物 / 特权奖励，教师审批发放', size: '' },
  { icon: '📡', title: '实时广播', desc: '消息直达教室桌面端，支持文字 / 语音 / 横幅 / 全屏', size: '' },
  { icon: '📝', title: '在线答题', desc: '题库管理 + 课堂即时检测，自动判分统计', size: '' },
  { icon: '🔗', title: '多端登录', desc: '微信 / 企微 / QQ / 人人通扫码，账号密码双通道', size: 'bento-lg' },
]

const evolutionStages = [
  { emoji: '🌟', name: '星尘' }, { emoji: '🌙', name: '月芽' },
  { emoji: '🌱', name: '灵苗' }, { emoji: '🌿', name: '青藤' },
  { emoji: '🌳', name: '慧树' }, { emoji: '🦋', name: '蝶灵' },
  { emoji: '🦅', name: '鹰慧' }, { emoji: '🦁', name: '狮睿' },
  { emoji: '🦄', name: '灵角' }, { emoji: '✨', name: '星耀' },
  { emoji: '🌌', name: '银河' },
]
</script>

<style scoped>
/* ================================================================
   LANDING PAGE — Apple 浅色风格
   白色卡片 + 微妙阴影 + 磨砂导航 + Bento 网格
   ================================================================ */

/* ── 根 ── */
.page {
  min-height: 100vh;
  background: #FFFFFF;
  color: #1D1D1F;
  font-family: 'Noto Sans SC', -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'Segoe UI', sans-serif;
  -webkit-font-smoothing: antialiased;
}

/* ── 导航栏 — 磨砂玻璃 ── */
.nav {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: saturate(180%) blur(20px);
  -webkit-backdrop-filter: saturate(180%) blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}
.nav-inner {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0 24px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.nav-brand {
  font-size: 18px;
  font-weight: 700;
  color: #1D1D1F;
  text-decoration: none;
  letter-spacing: -0.3px;
}
.nav-links {
  display: flex;
  align-items: center;
  gap: 20px;
}
.nav-links a:not(.nav-cta) {
  font-size: 13px;
  font-weight: 500;
  color: #6E6E73;
  text-decoration: none;
  transition: color 0.2s;
}
.nav-links a:not(.nav-cta):hover { color: #1D1D1F; }
.nav-github {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}
.nav-cta {
  display: inline-block;
  padding: 7px 18px;
  border-radius: 9999px;
  background: #1D1D1F;
  color: #FFFFFF !important;
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s;
}
.nav-cta:hover {
  background: #333;
  transform: scale(1.03);
}

/* ── Hero ── */
.hero {
  padding: 140px 20px 80px;
  background: linear-gradient(180deg, #FBFBFD 0%, #FFFFFF 100%);
}
.hero-inner {
  max-width: 660px;
  margin: 0 auto;
  text-align: center;
}
.badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #F5F5F7;
  border: 1px solid #E5E5EA;
  border-radius: 9999px;
  padding: 6px 18px;
  font-size: 13px;
  color: #6E6E73;
  margin-bottom: 32px;
}
.badge-dot {
  width: 7px; height: 7px;
  background: #34C759;
  border-radius: 50%;
  display: inline-block;
  box-shadow: 0 0 6px rgba(52, 199, 89, 0.35);
  animation: pulse-dot 2.5s ease-in-out infinite;
}
@keyframes pulse-dot {
  0%, 100% { box-shadow: 0 0 4px rgba(52, 199, 89, 0.3); }
  50%      { box-shadow: 0 0 12px rgba(52, 199, 89, 0.55); }
}
.hero-title {
  font-size: clamp(40px, 6vw, 60px);
  font-weight: 800;
  line-height: 1.12;
  letter-spacing: -1.5px;
  color: #1D1D1F;
  margin-bottom: 24px;
}
.hero-title-grad {
  background: linear-gradient(135deg, #5E5CE6 0%, #FF375F 50%, #FF9F0A 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.hero-desc {
  font-size: 17px;
  color: #6E6E73;
  line-height: 1.7;
  margin-bottom: 36px;
}

/* ── 按钮 ── */
.hero-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 56px;
}
.btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 14px 36px;
  border-radius: 9999px;
  background: linear-gradient(135deg, #5E5CE6, #7D7AFF);
  color: #FFFFFF;
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
  box-shadow: 0 4px 20px rgba(94, 92, 230, 0.25);
  transition: all 0.3s ease;
}
.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 28px rgba(94, 92, 230, 0.35);
}
.btn-secondary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 14px 36px;
  border-radius: 9999px;
  background: #FFFFFF;
  color: #1D1D1F;
  font-size: 16px;
  font-weight: 600;
  border: 1px solid #D2D2D7;
  text-decoration: none;
  transition: all 0.3s ease;
}
.btn-secondary:hover {
  background: #F5F5F7;
  border-color: #C7C7CC;
}
.btn-outline {
  display: inline-flex;
  align-items: center;
  padding: 12px 28px;
  border-radius: 9999px;
  background: #FFFFFF;
  color: #1D1D1F;
  font-size: 14px;
  font-weight: 600;
  border: 1px solid #D2D2D7;
  text-decoration: none;
  transition: all 0.2s;
}
.btn-outline:hover { background: #F5F5F7; border-color: #C7C7CC; }
.btn-ghost {
  display: inline-flex;
  align-items: center;
  padding: 12px 28px;
  border-radius: 9999px;
  background: transparent;
  color: #6E6E73;
  font-size: 14px;
  font-weight: 500;
  border: 1px solid #E5E5EA;
  text-decoration: none;
  transition: all 0.2s;
}
.btn-ghost:hover { color: #1D1D1F; border-color: #D2D2D7; }

/* ── Hero 标签 ── */
.hero-tags {
  display: flex;
  gap: 24px;
  justify-content: center;
  flex-wrap: wrap;
}
.tag {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  background: #FFFFFF;
  border: 1px solid #F0F0F3;
  border-radius: 14px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
  transition: all 0.3s ease;
}
.tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.06);
}
.tag-icon { font-size: 18px; }
.tag-label { font-size: 13px; color: #86868B; font-weight: 500; }

/* ── 分区交替背景 ── */
.section-gray  { background: #F5F5F7; padding: 80px 20px; }
.section-white { background: #FFFFFF; padding: 80px 20px; }
.container { max-width: 1100px; margin: 0 auto; }

/* ── 分区标题 ── */
.section-head { text-align: center; margin-bottom: 52px; }
.section-eyebrow {
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 2.5px;
  color: #5E5CE6;
  margin-bottom: 12px;
}
.section-title {
  font-size: clamp(26px, 4vw, 34px);
  font-weight: 800;
  color: #1D1D1F;
  margin-bottom: 12px;
  letter-spacing: -0.4px;
}
.section-desc {
  font-size: 16px;
  color: #86868B;
}

/* ── Bento 卡片网格 ── */
.bento {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}
.bento-card {
  position: relative;
  background: #FFFFFF;
  border: 1px solid #F0F0F3;
  border-radius: 20px;
  padding: 28px;
  transition: all 0.35s ease;
  cursor: default;
  box-shadow: 0 1px 3px rgba(0,0,0,0.03);
}
.bento-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 36px rgba(0,0,0,0.08);
  border-color: #E5E5EA;
}
.bento-icon  { display: block; font-size: 32px; margin-bottom: 14px; }
.bento-title { font-size: 15px; font-weight: 700; color: #1D1D1F; margin-bottom: 6px; }
.bento-desc  { font-size: 13px; color: #86868B; line-height: 1.55; }

/* Bento 尺寸变体 */
.bento-wide { grid-column: span 2; }
.bento-lg   { grid-column: span 2; }
.bento-tall { grid-row: span 2; }

/* ── 宠物进化 ── */
.evo-track {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
  flex-wrap: wrap;
  padding: 20px 0;
}
.evo-arrow { color: #C7C7CC; font-size: 14px; font-weight: 300; user-select: none; }
.evo-node {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  border-radius: 16px;
  transition: all 0.25s ease;
  cursor: default;
}
.evo-node:hover {
  background: #F5F5F7;
  transform: scale(1.1);
}
.evo-emoji { font-size: 30px; }
.evo-name  { font-size: 12px; color: #86868B; font-weight: 500; }

/* ── CTA ── */
.cta {
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
  background: #FFFFFF;
  border: 1px solid #F0F0F3;
  border-radius: 24px;
  padding: 52px 40px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.cta-icon  { font-size: 44px; margin-bottom: 18px; }
.cta-title { font-size: 28px; font-weight: 800; color: #1D1D1F; margin-bottom: 16px; }
.cta-desc  { font-size: 15px; color: #86868B; line-height: 1.7; margin-bottom: 32px; }
.cta-btns  { display: flex; gap: 12px; justify-content: center; flex-wrap: wrap; }

/* ── Footer ── */
.footer {
  padding: 32px 20px;
  text-align: center;
  background: #F5F5F7;
  border-top: 1px solid #E5E5EA;
}
.footer-line { color: #86868B; font-size: 13px; margin-bottom: 6px; }
.footer-link { color: #515154; text-decoration: none; transition: color 0.2s; }
.footer-link:hover { color: #1D1D1F; }
.footer-made { color: #AEAEB2; font-size: 12px; }

/* ── 滚动揭示 ── */
.reveal {
  opacity: 0;
  transform: translateY(36px);
  transition: opacity 0.7s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.7s cubic-bezier(0.16, 1, 0.3, 1);
}
.revealed {
  opacity: 1;
  transform: translateY(0);
}

/* ── Bento 交错延迟 ── */
.bento-card:nth-child(1)  { transition-delay: 0.00s, 0.00s; }
.bento-card:nth-child(2)  { transition-delay: 0.05s, 0.05s; }
.bento-card:nth-child(3)  { transition-delay: 0.10s, 0.10s; }
.bento-card:nth-child(4)  { transition-delay: 0.15s, 0.15s; }
.bento-card:nth-child(5)  { transition-delay: 0.20s, 0.20s; }
.bento-card:nth-child(6)  { transition-delay: 0.05s, 0.05s; }
.bento-card:nth-child(7)  { transition-delay: 0.10s, 0.10s; }
.bento-card:nth-child(8)  { transition-delay: 0.15s, 0.15s; }
.bento-card:nth-child(9)  { transition-delay: 0.20s, 0.20s; }
.bento-card:nth-child(10) { transition-delay: 0.05s, 0.05s; }
.bento-card:nth-child(11) { transition-delay: 0.10s, 0.10s; }
.bento-card:nth-child(12) { transition-delay: 0.15s, 0.15s; }

/* ── 响应式 ── */
@media (max-width: 1024px) {
  .bento { grid-template-columns: repeat(3, 1fr); }
  .bento-wide, .bento-lg { grid-column: span 2; }
}
@media (max-width: 768px) {
  .nav-inner { padding: 0 16px; height: 48px; }
  .nav-links a:not(.nav-cta) { display: none; }
  .hero { padding: 100px 16px 60px; }
  .hero-title { font-size: clamp(28px, 7vw, 42px); }
  .hero-actions { flex-direction: column; align-items: center; }
  .hero-tags { gap: 10px; }
  .bento { grid-template-columns: 1fr; }
  .bento-wide, .bento-lg, .bento-tall { grid-column: span 1; grid-row: span 1; }
  .section-gray, .section-white { padding: 48px 16px; }
  .section-head { margin-bottom: 32px; }
  .evo-track { gap: 1px; }
  .evo-node  { padding: 6px 8px; }
  .evo-emoji { font-size: 22px; }
  .evo-name  { font-size: 10px; }
  .cta { padding: 32px 20px; }
}
@media (max-width: 480px) {
  .evo-arrow { font-size: 10px; }
}
</style>
