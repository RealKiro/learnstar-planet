<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { PET_SERIES, getSpeciesEmoji } from '@/utils/petData'

const currentSlide = ref(0)
let timer: ReturnType<typeof setInterval>

function goToSlide(idx: number) {
  currentSlide.value = idx
  clearInterval(timer)
  timer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % PET_SERIES.length }, 5000)
}

onMounted(() => {
  timer = setInterval(() => { currentSlide.value = (currentSlide.value + 1) % PET_SERIES.length }, 5000)
})
onUnmounted(() => clearInterval(timer))
</script>

<template>
  <div>
    <h2 style="font-size:24px;font-weight:700;margin-bottom:20px;">📚 宠物图鉴</h2>

    <div style="position:relative;background:rgba(255,255,255,0.02);border-radius:var(--md-radius);padding:24px;border:1px solid rgba(255,255,255,0.04);overflow:hidden;min-height:400px;">
      <div v-for="(series, si) in PET_SERIES" :key="series.id"
        :style="{ display: currentSlide === si ? 'block' : 'none', animation: 'fadeIn 0.5s ease' }">
        <div style="display:flex;align-items:center;gap:16px;margin-bottom:16px;">
          <span style="font-size:48px;">{{ series.emoji }}</span>
          <span style="font-size:28px;font-weight:700;">{{ series.name }}</span>
        </div>
        <div style="color:var(--md-text-secondary);font-size:16px;margin-bottom:20px;padding:12px 16px;background:rgba(255,255,255,0.03);border-radius:12px;border-left:3px solid var(--md-primary);">
          {{ series.name }}系列的宠物们
        </div>
        <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:16px;">
          <div v-for="sp in series.species" :key="sp.id"
            style="background:rgba(255,255,255,0.03);border-radius:12px;padding:14px 12px;text-align:center;border:1px solid rgba(255,255,255,0.04);transition:0.2s;">
            <span style="font-size:36px;display:block;margin-bottom:4px;">{{ getSpeciesEmoji(sp.id) }}</span>
            <div style="font-weight:600;font-size:16px;">{{ sp.name }}</div>
            <div style="font-size:13px;color:var(--md-text-secondary);margin-top:4px;line-height:1.4;">
              {{ sp.levels[0]?.description || '' }}
            </div>
            <div style="display:flex;justify-content:center;gap:2px;margin-top:8px;">
              <span v-for="lvl in 12" :key="lvl"
                :style="{
                  display:'inline-block', width:'8px', height:'8px', borderRadius:'50%',
                  background: lvl === 12 ? '#f472b6' : (lvl <= 3 ? 'rgba(252,211,77,0.3)' : 'rgba(255,255,255,0.1)'),
                  boxShadow: lvl === 12 ? '0 0 8px rgba(244,114,182,0.3)' : 'none',
                  transition:'0.3s',
                }"></span>
            </div>
          </div>
        </div>
      </div>

      <div style="display:flex;justify-content:center;gap:8px;margin-top:24px;">
        <button v-for="(_, i) in PET_SERIES" :key="i"
          :style="{
            width: currentSlide === i ? '32px' : '12px', height:'12px',
            borderRadius: currentSlide === i ? '6px' : '50%',
            border: 'none', cursor:'pointer',
            background: currentSlide === i ? 'var(--md-primary)' : 'rgba(255,255,255,0.1)',
            boxShadow: currentSlide === i ? '0 0 12px rgba(167,139,250,0.3)' : 'none',
            transition:'0.2s',
          }"
          @click="goToSlide(i)"></button>
      </div>
    </div>

    <style>
      @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
    </style>
  </div>
</template>
