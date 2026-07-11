<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// 如果已登录，直接从首页开始
if (authStore.isLoggedIn) {
  if (authStore.isAdmin) router.replace({ name: 'admin-dashboard' })
  else if (authStore.isTeacher) router.replace({ name: 'teacher-dashboard' })
  else if (authStore.isParent) router.replace({ name: 'parent-home' })
}
</script>

<template>
  <div style="min-height: 100vh;">
    <!-- 导航栏 -->
    <nav style="position:fixed;top:0;left:0;right:0;padding:16px 32px;display:flex;align-items:center;justify-content:space-between;z-index:50;">
      <div style="font-size:20px;font-weight:700;color:#F1F5F9;display:flex;align-items:center;gap:8px;">
        <span>🌌</span> 学趣星球
      </div>
      <div style="display:flex;gap:16px;">
        <a href="#" style="color:#CBD5E1;font-size:14px;">功能介绍</a>
        <a href="#" style="color:#CBD5E1;font-size:14px;">联系我们</a>
        <router-link to="/login" class="btn btn-primary btn-sm" style="background:var(--gradient-primary);color:white;">登录</router-link>
      </div>
    </nav>

    <!-- Hero 区 -->
    <section style="min-height:600px;background:var(--gradient-hero);display:flex;align-items:center;padding:48px 32px;position:relative;overflow:hidden;">
      <div style="position:absolute;width:800px;height:800px;background:radial-gradient(circle,rgba(79,70,229,0.12) 0%,transparent 70%);top:-300px;right:-200px;"></div>
      <div style="position:relative;z-index:1;max-width:560px;">
        <h1 style="font-size:48px;font-weight:900;color:#F1F5F9;margin-bottom:16px;line-height:1.2;">
          让学习充满<span style="color:var(--color-primary-light);">乐趣</span>的成长平台
        </h1>
        <p style="font-size:18px;color:#94A3B8;margin-bottom:32px;">
          积分激励、宠物陪伴、排行竞技 —— 全免费，零门槛，让每个学生的努力都被看见
        </p>
        <div style="display:flex;gap:16px;">
          <router-link to="/login" class="btn" style="background:var(--gradient-primary);color:white;padding:14px 32px;border-radius:var(--radius-lg);font-size:16px;font-weight:600;box-shadow:0 0 20px rgba(79,70,229,0.3);">
            立即使用
          </router-link>
          <button class="btn" style="background:rgba(255,255,255,0.08);color:#CBD5E1;padding:14px 32px;border-radius:var(--radius-lg);font-size:16px;font-weight:600;border:1px solid rgba(255,255,255,0.15);">
            了解更多
          </button>
        </div>
        <div style="display:flex;gap:32px;margin-top:48px;">
          <div style="text-align:center;"><div style="font-size:28px;font-weight:900;color:#F1F5F9;">50,000+</div><div style="font-size:13px;color:#94A3B8;">活跃班级</div></div>
          <div style="text-align:center;"><div style="font-size:28px;font-weight:900;color:#F1F5F9;">200万+</div><div style="font-size:13px;color:#94A3B8;">学生用户</div></div>
          <div style="text-align:center;"><div style="font-size:28px;font-weight:900;color:#F1F5F9;">4.9⭐</div><div style="font-size:13px;color:#94A3B8;">教师好评</div></div>
        </div>
      </div>
    </section>

    <!-- 功能介绍 -->
    <section style="padding:48px 32px;">
      <h2 style="text-align:center;font-size:28px;font-weight:700;margin-bottom:8px;">核心功能</h2>
      <p style="text-align:center;color:var(--color-text-secondary);margin-bottom:32px;">12大模块，覆盖班级管理全场景</p>
      <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:24px;max-width:960px;margin:0 auto;">
        <div v-for="feat in features" :key="feat.title" class="card" style="text-align:center;">
          <div style="font-size:36px;margin-bottom:16px;">{{ feat.icon }}</div>
          <div style="font-weight:600;font-size:16px;margin-bottom:8px;">{{ feat.title }}</div>
          <div style="font-size:13px;color:var(--color-text-secondary);">{{ feat.desc }}</div>
        </div>
      </div>
    </section>

    <!-- 进化体系展示 -->
    <section style="padding:32px;">
      <h2 style="text-align:center;margin-bottom:16px;">宠物10阶进化体系</h2>
      <div style="display:flex;align-items:center;justify-content:center;gap:4px;flex-wrap:wrap;padding:24px 0;">
        <template v-for="(stage, i) in evolutionStages" :key="stage.name">
          <div v-if="i > 0" style="font-size:16px;color:var(--color-text-secondary);">→</div>
          <div style="display:flex;flex-direction:column;align-items:center;gap:4px;padding:8px;border-radius:var(--radius-md);cursor:pointer;transition:transform 0.2s;">
            <span style="font-size:32px;">{{ stage.emoji }}</span>
            <span style="font-size:11px;font-weight:500;">{{ stage.name }}</span>
          </div>
        </template>
      </div>
    </section>

    <!-- 底部 -->
    <footer style="text-align:center;padding:32px;color:var(--color-text-secondary);font-size:13px;border-top:1px solid var(--color-border);">
      MIT 开源许可证 · 学趣星球 · 让学习充满乐趣
    </footer>
  </div>
</template>

<script lang="ts">
const features = [
  { icon: '⭐', title: '积分激励', desc: '自定义规则，实时加减分，进步看得见' },
  { icon: '🌟', title: '宠物进化', desc: '10阶宇宙进化体系，积分变经验驱动成长' },
  { icon: '🏆', title: '排行竞技', desc: '总积分/周进步/宠物等级三大排行' },
  { icon: '📢', title: '班级通知', desc: '一键发布，实时推送家长端' },
  { icon: '📊', title: '成绩管理', desc: '成绩录入分析，班级对比，趋势可视化' },
  { icon: '🤖', title: 'AI助教', desc: '生成班级反馈、学生分析、家校沟通建议' },
  { icon: '✅', title: '智能考勤', desc: '一键点名签到，到课/请假/迟到状态实时统计' },
  { icon: '📱', title: '扫码收作业', desc: '生成专属二维码，学生扫码提交，自动汇总' },
  { icon: '🛍️', title: '积分商城', desc: '学生兑换实物/特权奖励，教师审批发放' },
  { icon: '📡', title: '实时广播', desc: '消息直达教室桌面端，支持文字/语音/横幅/全屏' },
  { icon: '📝', title: '在线答题', desc: '题库管理+课堂即时检测，自动判分统计' },
  { icon: '🔗', title: '多端登录', desc: '微信/企微/QQ/人人通扫码，账号密码双通道' },
]

const evolutionStages = [
  { emoji: '🌟', name: '星尘' },
  { emoji: '🌙', name: '月芽' },
  { emoji: '🌱', name: '灵苗' },
  { emoji: '🌿', name: '青藤' },
  { emoji: '🌳', name: '慧树' },
  { emoji: '🦋', name: '蝶灵' },
  { emoji: '🦅', name: '鹰慧' },
  { emoji: '🦁', name: '狮睿' },
  { emoji: '🦄', name: '灵角' },
  { emoji: '✨', name: '星耀' },
  { emoji: '🌌', name: '银河' },
]
</script>
