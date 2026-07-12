<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

if (authStore.isLoggedIn) {
  if (authStore.isAdmin) router.replace({ name: 'admin-dashboard' })
  else if (authStore.isTeacher) router.replace({ name: 'teacher-dashboard' })
  else if (authStore.isParent) router.replace({ name: 'parent-home' })
}
</script>

<template>
  <div style="min-height:100vh;background:#0F172A;">
    <!-- 导航栏 -->
    <nav style="position:fixed;top:0;left:0;right:0;padding:16px 32px;display:flex;align-items:center;justify-content:space-between;z-index:50;background:rgba(15,23,42,0.8);backdrop-filter:blur(20px);border-bottom:1px solid rgba(255,255,255,0.06);">
      <div style="font-size:20px;font-weight:700;color:#F1F5F9;display:flex;align-items:center;gap:8px;">
        <span>🌌</span> 学趣星球
      </div>
      <div style="display:flex;align-items:center;gap:20px;">
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="color:#94A3B8;font-size:14px;display:flex;align-items:center;gap:6px;transition:color 0.2s;" onmouseover="this.style.color='#F1F5F9'" onmouseout="this.style.color='#94A3B8'">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
          GitHub
        </a>
        <router-link to="/login" style="background:rgba(79,70,229,0.9);color:#F1F5F9;padding:8px 20px;border-radius:8px;font-size:14px;font-weight:500;transition:all 0.2s;" onmouseover="this.style.background='rgba(67,56,202,0.95)'" onmouseout="this.style.background='rgba(79,70,229,0.9)'">登录</router-link>
      </div>
    </nav>

    <!-- Hero 区 -->
    <section style="min-height:100vh;display:flex;align-items:center;justify-content:center;padding:80px 32px 48px;position:relative;overflow:hidden;">
      <!-- 背景光效 -->
      <div style="position:absolute;width:700px;height:700px;background:radial-gradient(circle,rgba(79,70,229,0.15) 0%,transparent 70%);top:-200px;right:-100px;animation:floatOrb 12s ease-in-out infinite;"></div>
      <div style="position:absolute;width:500px;height:500px;background:radial-gradient(circle,rgba(16,185,129,0.08) 0%,transparent 70%);bottom:-150px;left:-100px;animation:floatOrb 15s ease-in-out infinite reverse;"></div>

      <div style="position:relative;z-index:1;max-width:680px;text-align:center;">
        <!-- 开源徽章 -->
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="display:inline-flex;align-items:center;gap:8px;background:rgba(255,255,255,0.06);border:1px solid rgba(255,255,255,0.12);border-radius:100px;padding:6px 16px;font-size:13px;color:#94A3B8;margin-bottom:24px;text-decoration:none;transition:all 0.2s;" onmouseover="this.style.background='rgba(255,255,255,0.1)'" onmouseout="this.style.background='rgba(255,255,255,0.06)'">
          <span style="color:#10B981;">●</span>
          MIT 开源 · 完全免费 · 自托管
        </a>

        <h1 style="font-size:52px;font-weight:900;color:#F1F5F9;margin-bottom:20px;line-height:1.15;letter-spacing:-1px;">
          让每个孩子的努力<br>
          <span style="background:linear-gradient(135deg,#818CF8,#6EE7B7);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text;">都被看见</span>
        </h1>

        <p style="font-size:18px;color:#94A3B8;margin-bottom:36px;line-height:1.6;">
          积分激励 · 宠物养成 · AI 助教 · 多端同步<br>
          开源班级管理系统，Docker 一键部署，数据完全自主掌控
        </p>

        <div style="display:flex;gap:16px;justify-content:center;margin-bottom:48px;">
          <router-link to="/login" style="background:linear-gradient(135deg,#6366F1,#818CF8);color:white;padding:14px 36px;border-radius:12px;font-size:16px;font-weight:600;box-shadow:0 8px 24px rgba(79,70,229,0.3);transition:all 0.25s;text-decoration:none;" onmouseover="this.style.transform='translateY(-2px)';this.style.boxShadow='0 12px 32px rgba(79,70,229,0.4)'" onmouseout="this.style.transform='';this.style.boxShadow='0 8px 24px rgba(79,70,229,0.3)'">
            立即使用
          </router-link>
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="background:rgba(255,255,255,0.06);color:#E2E8F0;padding:14px 36px;border-radius:12px;font-size:16px;font-weight:600;border:1px solid rgba(255,255,255,0.12);transition:all 0.25s;text-decoration:none;display:flex;align-items:center;gap:8px;" onmouseover="this.style.background='rgba(255,255,255,0.1)'" onmouseout="this.style.background='rgba(255,255,255,0.06)'">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
            源码
          </a>
        </div>

        <!-- 项目特点 -->
        <div style="display:flex;gap:40px;justify-content:center;flex-wrap:wrap;">
          <div v-for="item in highlights" :key="item.label" style="text-align:center;">
            <div style="font-size:28px;margin-bottom:4px;">{{ item.icon }}</div>
            <div style="font-size:13px;color:#64748B;">{{ item.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 功能介绍 -->
    <section style="padding:80px 32px;background:rgba(255,255,255,0.02);border-top:1px solid rgba(255,255,255,0.06);border-bottom:1px solid rgba(255,255,255,0.06);">
      <div style="max-width:960px; ">
        <h2 style="text-align:center;font-size:32px;font-weight:800;color:#F1F5F9;margin-bottom:12px;">12 大功能模块</h2>
        <p style="text-align:center;color:#64748B;margin-bottom:48px;font-size:15px;">覆盖班级管理全场景，全部免费，无任何付费功能</p>
        <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:16px;">
          <div v-for="feat in features" :key="feat.title"
            style="background:rgba(255,255,255,0.03);border:1px solid rgba(255,255,255,0.06);border-radius:16px;padding:24px;transition:all 0.25s;"
            onmouseover="this.style.background='rgba(255,255,255,0.06)';this.style.borderColor='rgba(79,70,229,0.3)';this.style.transform='translateY(-3px)'"
            onmouseout="this.style.background='rgba(255,255,255,0.03)';this.style.borderColor='rgba(255,255,255,0.06)';this.style.transform=''">
            <div style="font-size:28px;margin-bottom:12px;">{{ feat.icon }}</div>
            <div style="font-weight:600;font-size:15px;color:#E2E8F0;margin-bottom:6px;">{{ feat.title }}</div>
            <div style="font-size:13px;color:#64748B;line-height:1.5;">{{ feat.desc }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 宠物进化展示 -->
    <section style="padding:64px 32px;text-align:center;">
      <h2 style="font-size:28px;font-weight:800;color:#F1F5F9;margin-bottom:8px;">宠物 11 阶进化体系</h2>
      <p style="color:#64748B;margin-bottom:32px;font-size:14px;">积分变经验，驱动宠物从星尘进化到银河</p>
      <div style="display:flex;align-items:center;justify-content:center;gap:2px;flex-wrap:wrap;max-width:800px; ">
        <template v-for="(stage, i) in evolutionStages" :key="stage.name">
          <div v-if="i > 0" style="font-size:14px;color:#334155;margin:0 2px;">→</div>
          <div style="display:flex;flex-direction:column;align-items:center;gap:4px;padding:10px 8px;border-radius:12px;transition:all 0.2s;"
            onmouseover="this.style.background='rgba(79,70,229,0.1)';this.style.transform='scale(1.15)'"
            onmouseout="this.style.background='transparent';this.style.transform=''">
            <span style="font-size:28px;">{{ stage.emoji }}</span>
            <span style="font-size:11px;font-weight:500;color:#64748B;">{{ stage.name }}</span>
          </div>
        </template>
      </div>
    </section>

    <!-- 开源号召 -->
    <section style="padding:64px 32px;text-align:center;">
      <div style="max-width:600px; background:linear-gradient(135deg,rgba(79,70,229,0.1),rgba(16,185,129,0.05));border:1px solid rgba(79,70,229,0.15);border-radius:20px;padding:48px 32px;">
        <div style="font-size:40px;margin-bottom:16px;">🚀</div>
        <h2 style="font-size:26px;font-weight:800;color:#F1F5F9;margin-bottom:12px;">加入开源社区</h2>
        <p style="color:#94A3B8;margin-bottom:24px;font-size:15px;line-height:1.6;">
          项目完全开源，欢迎 Star ⭐、Fork、提交 Issue 和 PR<br>
          如果你正在用它管理班级，你就是这个社区的一员
        </p>
        <div style="display:flex;gap:12px;justify-content:center;flex-wrap:wrap;">
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="background:rgba(255,255,255,0.08);color:#E2E8F0;padding:12px 28px;border-radius:10px;font-size:14px;font-weight:500;border:1px solid rgba(255,255,255,0.12);text-decoration:none;transition:all 0.2s;" onmouseover="this.style.background='rgba(255,255,255,0.12)'" onmouseout="this.style.background='rgba(255,255,255,0.08)'">⭐ 给个 Star</a>
          <a href="https://github.com/RealKiro/learnstar-planet/issues" target="_blank" style="background:transparent;color:#94A3B8;padding:12px 28px;border-radius:10px;font-size:14px;font-weight:500;border:1px solid rgba(255,255,255,0.08);text-decoration:none;transition:all 0.2s;" onmouseover="this.style.color='#E2E8F0'" onmouseout="this.style.color='#94A3B8'">提交 Issue</a>
        </div>
      </div>
    </section>

    <!-- 底部 -->
    <footer style="padding:32px;text-align:center;border-top:1px solid rgba(255,255,255,0.06);">
      <div style="color:#475569;font-size:13px;margin-bottom:8px;">
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="color:#475569;text-decoration:none;">学趣星球</a>
        · MIT 开源许可证 · 自托管 · 完全免费
      </div>
      <div style="color:#334155;font-size:12px;">Made with 💜 by open source community</div>
    </footer>
  </div>
</template>

<style scoped>
@keyframes floatOrb {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(40px, -30px); }
}
</style>

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
  { icon: '⭐', title: '积分激励', desc: '自定义规则，实时加减分，进步看得见' },
  { icon: '🌟', title: '宠物进化', desc: '11 阶宇宙进化体系，积分变经验驱动成长' },
  { icon: '🏆', title: '排行竞技', desc: '总积分/周进步/宠物等级三大排行' },
  { icon: '📢', title: '班级通知', desc: '一键发布，实时推送家长端' },
  { icon: '📊', title: '成绩管理', desc: '成绩录入分析，班级对比，趋势可视化' },
  { icon: '🤖', title: 'AI 助教', desc: '生成班级反馈、学生分析、家校沟通建议' },
  { icon: '✅', title: '智能考勤', desc: '一键点名签到，到课/请假/迟到状态统计' },
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
