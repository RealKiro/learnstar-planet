<script setup lang="ts">
import LeftIntro from './LeftIntro.vue'
import LoginPanel from './LoginPanel.vue'
</script>

<template>
  <div class="landing-wrapper">
    <div class="layout">
      <div class="left">
        <LeftIntro />
      </div>
      <div class="right">
        <LoginPanel />
      </div>
    </div>
  </div>
</template>

<style scoped>
* { box-sizing: border-box }
.landing-wrapper { display: flex; justify-content: center; align-items: center; min-height: 100vh; background: #fbfbfd }
.layout { display: flex; width: 100%; max-width: 1100px; height: 90vh; max-height: 760px; background: #fff; color: #1d1d1f; font-family: 'Noto Sans SC',-apple-system,BlinkMacSystemFont,'SF Pro Display',sans-serif; -webkit-font-smoothing: antialiased; border-radius: 24px; box-shadow: 0 8px 60px rgba(0,0,0,.06); overflow: hidden }
.left { flex: 1; height: 100%; position: relative; background: #fbfbfd; overflow: hidden; border-radius: 24px 0 0 24px }
.right { flex: 0 0 400px; display: flex; align-items: center; justify-content: center; padding: 32px; background: #fff; border-left: 1px solid #f0f0f3; border-radius: 0 24px 24px 0; height: 100% }
@media (max-width: 768px) {
  .landing-wrapper { align-items: stretch }
  .layout { flex-direction: column; height: auto; max-height: none; max-width: 100%; border-radius: 0; box-shadow: none }
  .left { height: 600px; border-radius: 0 }
  .right { flex: none; height: auto; border-left: none; border-top: 1px solid #f0f0f3; border-radius: 0; padding: 32px 24px }
}
</style>
