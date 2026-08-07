#!/usr/bin/env node
/**
 * 生成全部宠物物种 × 6 阶段的 AI 生图提示词 CSV。
 * 数据契约：docs/宠物系统全案.md 6.4 节 + docs/pet-image-manifest.md
 *
 * 用法：node scripts/generate-pet-prompts.js [输出路径]（默认 frontend-vue/pet-prompts.csv）
 */
import fs from 'fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const srcDir = path.join(__dirname, '../src/utils')
const src = fs.readFileSync(path.join(srcDir, 'petData.ts'), 'utf8')
const ext = fs.readFileSync(path.join(srcDir, 'petDataExtended.ts'), 'utf8')
const all = src + '\n' + ext

// —— 六阶段视觉基调（参考《系列角色六阶段文案策划.txt》）——
const STAGES = [
  { key: 'egg', zh: '灵胎初醒', main: '#4A3A2A', sub: '#C4A672', kw: '渺小、独行、破败背景', el: '孤灯/破庙', startLv: 1 },
  { key: 'baby', zh: '凡尘砺心', main: '#5B6B7A', sub: '#8A9BA8', kw: '行路、风雨、远山', el: '雨/风', startLv: 3 },
  { key: 'growing', zh: '道法初成', main: '#C44B2A', sub: '#E8B84B', kw: '光芒初绽、首胜、昂首', el: '光芒/法器', startLv: 5 },
  { key: 'mature', zh: '大劫淬炼', main: '#1A0A0A', sub: '#8B1A1A', kw: '碎裂、坠落、血与火', el: '火焰/裂痕', startLv: 7 },
  { key: 'legendary', zh: '封神登天', main: '#5B2A8A', sub: '#F5D742', kw: '冠冕、宝座、万众仰视', el: '冠冕/神位', startLv: 9 },
  { key: 'transcendent', zh: '归真永恒', main: '#F5F0E8', sub: '#8A9BA8', kw: '清风、青山、无我之相', el: '山水/清风', startLv: 11 },
]

// —— 解析系列 id → 名称 ——
const seriesName = {}
for (const m of all.matchAll(/id: '([a-z_]+)', name: '([^']+)'/g)) seriesName[m[1]] = m[2]

// —— 解析物种：id / name / seriesId / 各级等级名 ——
function parseSpecies(block, sid) {
  const out = []
  // 基础格式：{ level: N, name: 'X', ... , seriesId: 'S' } 前是物种头 id: 'x', name: 'y'
  const speciesRe = /id: '([a-z_]+)', name: '([^']+)', seriesId: '([a-z_]+)',\s*levels: \[([\s\S]*?)\n(\s*)\],/g
  let m
  while ((m = speciesRe.exec(block))) {
    const [, id, name, seriesId, levelsStr] = m
    const lv = {}
    for (const lm of levelsStr.matchAll(/\{ level: (\d+), name: '([^']+)'/g)) lv[+lm[1]] = lm[2]
    for (const lm of levelsStr.matchAll(/Lv\((\d+), '([^']+)'/g)) lv[+lm[1]] = lm[2]
    out.push({ id, name, seriesId, lv })
  }
  return out
}
const species = [
  ...parseSpecies(src, ''),
  ...parseSpecies(ext, ''),
]
// 去重（EXTRA 合并）
const seen = new Set()
const unique = species.filter(s => (seen.has(s.id) ? false : (seen.add(s.id), true)))

// —— 生成 CSV ——
const rows = [['speciesId', 'name', 'stage', 'stageZh', 'prompt']]
for (const s of unique) {
  for (const st of STAGES) {
    const lvName = s.lv[st.startLv] || s.lv[st.startLv + 1] || '未知形态'
    const series = seriesName[s.seriesId] || s.seriesId
    const prompt =
      `${s.name} 的 ${st.zh} 阶段。` +
      `外形：${lvName}。` +
      `系列：${series}。` +
      `阶段基调：${st.zh}，${st.main} 主色、${st.sub} 辅色，${st.kw}，必有${st.el}。` +
      `风格：儿童友好的 2.5D 卡通渲染，柔和渐变背景，居中全身像，正面 3/4 视角，无文字。` +
      `画布：800×1000 竖版，角色主体占 60% 居中偏下。`
    rows.push([s.id, s.name, st.key, st.zh, prompt])
  }
}
const csv = rows.map(r => r.map(c => `"${c.replace(/"/g, '""')}"`).join(',')).join('\n')
const outPath = process.argv[2] || path.join(__dirname, '../pet-prompts.csv')
fs.writeFileSync(outPath, '﻿' + csv, 'utf8')
console.log(`已生成 ${unique.length} 物种 × 6 阶段 = ${unique.length * 6} 条提示词 → ${outPath}`)
console.log('示例：', rows[1][4] || '')
