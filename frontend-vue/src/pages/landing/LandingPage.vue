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
  <div style="background:#0F172A;color:#F1F5F9;font-family:system-ui,-apple-system,sans-serif;overflow-x:hidden;">
    <!-- 导航栏 -->
    <nav style="position:fixed;top:0;left:0;right:0;z-index:50;background:rgba(15,23,42,0.85);backdrop-filter:blur(20px);border-bottom:1px solid rgba(255,255,255,0.06);">
      <div style="max-width:960px;margin:0 auto;padding:16px 24px;display:flex;align-items:center;justify-content:space-between;">
        <div style="font-size:18px;font-weight:700;display:flex;align-items:center;gap:8px;">
          <span>🌌</span> 学趣星球
        </div>
        <div style="display:flex;align-items:center;gap:16px;">
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="color:#94A3B8;font-size:13px;display:flex;align-items:center;gap:6px;text-decoration:none;transition:color 0.2s;" onmouseover="this.style.color='#E2E8F0'" onmouseout="this.style.color='#94A3B8'">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
            GitHub
          </a>
          <router-link to="/login" style="background:linear-gradient(135deg,#6366F1,#818CF8);color:#fff;padding:8px 20px;border-radius:8px;font-size:13px;font-weight:500;text-decoration:none;transition:opacity 0.2s;" onmouseover="this.style.opacity='0.9'" onmouseout="this.style.opacity='1'">登录</router-link>
        </div>
      </div>
    </nav>

    <!-- Hero -->
    <section style="min-height:100vh;display:flex;flex-direction:column;align-items:center;justify-content:center;padding:100px 0 60px;position:relative;overflow:hidden;width:100%;">
      <div style="position:absolute;width:600px;height:600px;background:radial-gradient(circle,rgba(99,102,241,0.12) 0%,transparent 70%);top:10%;left:50%;transform:translate(-50%,-50%);pointer-events:none;"></div>
      <div style="position:relative;z-index:1;max-width:600px;width:100%;padding:0 20px;text-align:center;margin:0 auto;">
        <div style="display:inline-flex;align-items:center;gap:6px;background:rgba(255,255,255,0.05);border:1px solid rgba(255,255,255,0.1);border-radius:100px;padding:5px 14px;font-size:12px;color:#94A3B8;margin-bottom:28px;">
          <span style="width:6px;height:6px;background:#10B981;border-radius:50%;display:inline-block;"></span>
          MIT 开源 · 完全免费 · 自托管
        </div>

        <h1 style="font-size:clamp(36px,5vw,52px);font-weight:900;line-height:1.2;letter-spacing:-0.5px;margin-bottom:20px;">
          让每个孩子的努力<br>
          <span style="background:linear-gradient(135deg,#818CF8,#6EE7B7);-webkit-background-clip:text;-webkit-text-fill-color:transparent;">都被看见</span>
        </h1>

        <p style="font-size:16px;color:#94A3B8;line-height:1.7;margin-bottom:36px;">
          积分激励 · 宠物养成 · AI 助教 · 多端同步<br>开源班级管理系统，Docker 一键部署，数据完全自主掌控
        </p>

        <div style="display:flex;gap:12px;justify-content:center;flex-wrap:wrap;margin-bottom:56px;">
          <router-link to="/login" style="background:linear-gradient(135deg,#6366F1,#818CF8);color:#fff;padding:13px 32px;border-radius:10px;font-size:15px;font-weight:600;text-decoration:none;box-shadow:0 4px 20px rgba(99,102,241,0.35);transition:transform 0.2s,box-shadow 0.2s;" onmouseover="this.style.transform='translateY(-1px)';this.style.boxShadow='0 6px 28px rgba(99,102,241,0.5)'" onmouseout="this.style.transform='';this.style.boxShadow='0 4px 20px rgba(99,102,241,0.35)'">立即使用</router-link>
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="background:rgba(255,255,255,0.06);color:#E2E8F0;padding:13px 32px;border-radius:10px;font-size:15px;font-weight:600;border:1px solid rgba(255,255,255,0.1);text-decoration:none;display:flex;align-items:center;gap:8px;transition:background 0.2s;" onmouseover="this.style.background='rgba(255,255,255,0.1)'" onmouseout="this.style.background='rgba(255,255,255,0.06)'">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
            源码
          </a>
        </div>

        <div style="display:flex;gap:32px;justify-content:center;flex-wrap:wrap;">
          <div v-for="item in highlights" :key="item.label" style="text-align:center;">
            <div style="font-size:24px;margin-bottom:4px;">{{ item.icon }}</div>
            <div style="font-size:12px;color:#64748B;">{{ item.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 功能模块 -->
    <section style="padding:80px 20px;border-top:1px solid rgba(255,255,255,0.06);border-bottom:1px solid rgba(255,255,255,0.06);width:100%;">
      <div style="max-width:1000px;margin:0 auto;">
        <h2 style="text-align:center;font-size:28px;font-weight:800;margin-bottom:8px;">12 大功能模块</h2>
        <p style="text-align:center;color:#64748B;margin-bottom:48px;font-size:14px;">覆盖班级管理全场景，全部免费，无任何付费功能</p>
        <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:14px;">
          <div v-for="feat in features" :key="feat.title"
            style="background:rgba(255,255,255,0.03);border:1px solid rgba(255,255,255,0.06);border-radius:14px;padding:22px;transition:all 0.25s;"
            onmouseover="this.style.background='rgba(99,102,241,0.08)';this.style.borderColor='rgba(99,102,241,0.25)';this.style.transform='translateY(-3px)'"
            onmouseout="this.style.background='rgba(255,255,255,0.03)';this.style.borderColor='rgba(255,255,255,0.06)';this.style.transform=''">
            <div style="font-size:26px;margin-bottom:10px;">{{ feat.icon }}</div>
            <div style="font-weight:600;font-size:14px;margin-bottom:5px;">{{ feat.title }}</div>
            <div style="font-size:12px;color:#64748B;line-height:1.5;">{{ feat.desc }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 宠物进化 -->
    <section style="padding:64px 20px;width:100%;">
      <div style="max-width:900px;margin:0 auto;text-align:center;">
        <h2 style="font-size:26px;font-weight:800;margin-bottom:6px;">宠物 11 阶进化体系</h2>
        <p style="color:#64748B;margin-bottom:32px;font-size:13px;">积分变经验，驱动宠物从星尘进化到银河</p>
        <div style="display:flex;align-items:center;justify-content:center;gap:0;flex-wrap:wrap;">
          <template v-for="(stage, i) in evolutionStages" :key="stage.name">
            <span v-if="i > 0" style="color:#334155;font-size:13px;margin:0 2px;">→</span>
            <div style="display:flex;flex-direction:column;align-items:center;gap:4px;padding:8px 10px;border-radius:10px;transition:all 0.2s;cursor:default;"
              onmouseover="this.style.background='rgba(99,102,241,0.1)';this.style.transform='scale(1.12)'"
              onmouseout="this.style.background='transparent';this.style.transform=''">
              <span style="font-size:26px;">{{ stage.emoji }}</span>
              <span style="font-size:11px;color:#64748B;">{{ stage.name }}</span>
            </div>
          </template>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section style="padding:64px 20px;width:100%;">
      <div style="max-width:560px;margin:0 auto;background:linear-gradient(135deg,rgba(99,102,241,0.08),rgba(16,185,129,0.04));border:1px solid rgba(99,102,241,0.12);border-radius:18px;padding:44px 32px;text-align:center;">
        <div style="font-size:36px;margin-bottom:14px;">🚀</div>
        <h2 style="font-size:24px;font-weight:800;margin-bottom:10px;">加入开源社区</h2>
        <p style="color:#94A3B8;margin-bottom:24px;font-size:14px;line-height:1.6;">
          项目完全开源，欢迎 Star ⭐、Fork、提交 Issue 和 PR<br>如果你正在用它管理班级，你就是这个社区的一员
        </p>
        <div style="display:flex;gap:10px;justify-content:center;flex-wrap:wrap;">
          <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="background:rgba(255,255,255,0.08);color:#E2E8F0;padding:11px 26px;border-radius:10px;font-size:14px;font-weight:500;border:1px solid rgba(255,255,255,0.1);text-decoration:none;transition:all 0.2s;" onmouseover="this.style.background='rgba(255,255,255,0.14)'" onmouseout="this.style.background='rgba(255,255,255,0.08)'">⭐ 给个 Star</a>
          <a href="https://github.com/RealKiro/learnstar-planet/issues" target="_blank" style="color:#94A3B8;padding:11px 26px;border-radius:10px;font-size:14px;font-weight:500;border:1px solid rgba(255,255,255,0.06);text-decoration:none;transition:all 0.2s;" onmouseover="this.style.color='#E2E8F0';this.style.borderColor='rgba(255,255,255,0.14)'" onmouseout="this.style.color='#94A3B8';this.style.borderColor='rgba(255,255,255,0.06)'">提交 Issue</a>
        </div>
      </div>
    </section>

    <!-- 底部 -->
    <footer style="padding:28px 20px;text-align:center;border-top:1px solid rgba(255,255,255,0.05);width:100%;">
      <div style="color:#475569;font-size:12px;margin-bottom:6px;">
        <a href="https://github.com/RealKiro/learnstar-planet" target="_blank" style="color:#475569;text-decoration:none;">学趣星球</a>
        · MIT 开源许可证 · 自托管 · 完全免费
      </div>
      <div style="color:#334155;font-size:11px;">Made with 💜 by open source community</div>
    </footer>
  </div>
</template>

<style scoped>
@keyframes floatOrb {
  0%, 100% { transform: translate(-50%, -50%) translate(0, 0); }
  50% { transform: translate(-50%, -50%) translate(30px, -20px); }
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
                                                                                                                                                                                                                                                                                                                               