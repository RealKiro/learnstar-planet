// 宠物数据全量审计脚本：检查特质/档案/诗词/emoji/视觉规格的完整性与雷同度
// 用法：cd frontend-vue && node scripts/analyze-pets.mjs（报告输出到 docs/pet-audit-report.json）
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'
import ts from 'typescript'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const frontendDir = path.resolve(__dirname, '..')
const repoRoot = path.resolve(frontendDir, '..')
const utilsDir = path.join(frontendDir, 'src/utils')
const docsDir = path.resolve(repoRoot, 'docs')
const tmpDir = path.resolve(repoRoot, '..', '.pet-audit-tmp')
fs.writeFileSync(path.join(tmpDir, 'package.json'), JSON.stringify({ type: 'commonjs' }))

// 1. 转译 TS → CJS
const FILES = ['petDataExtended', 'petData', 'petLifeStories', 'petProfiles', 'petTraits', 'stageEmoji']
fs.rmSync(tmpDir, { recursive: true, force: true })
fs.mkdirSync(path.join(tmpDir, 'utils'), { recursive: true })
for (const name of FILES) {
  const src = fs.readFileSync(path.join(utilsDir, `${name}.ts`), 'utf8')
  const js = ts.transpileModule(src, {
    compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020, esModuleInterop: true },
  }).outputText
  fs.writeFileSync(path.join(tmpDir, 'utils', `${name}.js`), js)
}

// 2. 载入
const require2 = createRequire(path.join(tmpDir, 'utils', 'petData.js'))
const { PET_SERIES, SPECIES_EMOJI, getLevelData } = require2('./petData.js')
const { PET_LIFE_STORIES } = require2('./petLifeStories.js')
const { PET_PROFILES } = require2('./petProfiles.js')
const { PET_TRAITS } = require2('./petTraits.js')
const stageEmojiMod = require2('./stageEmoji.js')

// 3. 审计
const allSpecies = PET_SERIES.flatMap(s => s.species)
const speciesIds = allSpecies.map(s => s.id)
const report = { total: allSpecies.length, series: [], emojiDup: [], missingStory: [], storyIssues: [], traitsIssues: [], profileIssues: [], levelIssues: [], stageEmojiDup: {} }

// 系列分布
for (const s of PET_SERIES) {
  report.series.push({ id: s.id, name: s.name, count: s.species.length })
}

// emoji 重复（同形态：基础 emoji 或各阶段 emoji 相同的物种对）
const byEmoji = {}
for (const [id, e] of Object.entries(SPECIES_EMOJI)) {
  ;(byEmoji[e] = byEmoji[e] || []).push(id)
}
for (const [e, ids] of Object.entries(byEmoji)) {
  if (ids.length > 1) report.emojiDup.push({ emoji: e, ids })
}

// stageEmoji：阶段去重情况（该模块提供 getStageEmoji(speciesId, level)）
if (typeof stageEmojiMod.getStageEmoji === 'function') {
  const stageSeen = {}
  for (const id of speciesIds) {
    for (const lv of [1, 3, 5, 7, 9, 12]) {
      const e = stageEmojiMod.getStageEmoji(id, lv)
      const key = `L${lv}:${e}`
      ;(stageSeen[key] = stageSeen[key] || []).push(id)
    }
  }
  for (const [key, ids] of Object.entries(stageSeen)) {
    if (ids.length > 3) report.stageEmojiDup[key] = ids.length // 同阶段同 emoji 的物种数（>3 视为雷同）
  }
}

// 生灵档案
const VISUAL_FIELDS = ['body', 'face', 'attireDetail', 'hairStyle', 'powerEffect', 'actionSeq']
let interactionCovered = 0, interactionTotal = 0
for (const id of speciesIds) {
  const st = PET_LIFE_STORIES[id]
  if (!st) { report.missingStory.push(id); continue }
  if ((st.stages || []).length !== 6) {
    report.storyIssues.push({ id, issue: `stages=${(st.stages || []).length}` })
    continue
  }
  st.stages.forEach((s, i) => {
    const stageName = ['卵生', '幼年', '成长', '成熟', '传说', '道果'][i]
    if (!s.poem || s.poem.length < 10) report.storyIssues.push({ id, issue: `${stageName}诗词缺失/过短(${(s.poem || '').length}字)` })
    if (!s.line || s.line.length < 4) report.storyIssues.push({ id, issue: `${stageName}台词缺失` })
    if (!s.character || s.character.length < 4) report.storyIssues.push({ id, issue: `${stageName}品性缺失` })
    interactionTotal++
    if (s.interaction && s.interaction.length > 4) interactionCovered++
    const missingVisual = VISUAL_FIELDS.filter(f => !s[f] || String(s[f]).length < 6)
    if (missingVisual.length >= 3) report.storyIssues.push({ id, issue: `${stageName}视觉规格缺${missingVisual.length}项(${missingVisual.join(',')})` })
  })
}
report.interactionCoverage = { covered: interactionCovered, total: interactionTotal }

// 互动叙事模板化检测：道果阶段以「把X传给/留给/教给」收尾视为套模板
let formulaicFinale = 0
const formulaicIds = []
for (const id of speciesIds) {
  const st = PET_LIFE_STORIES[id]
  const last = st?.stages?.[5]?.interaction
  if (last && /^(它|他|她)(把|教|受封.*留)/.test(last) && /(传给|留给|教给|教他们|教幼|教人|教孩子|教采药)/.test(last)) {
    formulaicFinale++
    formulaicIds.push(id)
  }
}
report.formulaicFinale = { count: formulaicFinale, ids: formulaicIds }

// 特质：数组长度 + 跨物种重复
const traitDup = {}
for (const id of speciesIds) {
  const t = PET_TRAITS[id]
  if (!t) { report.traitsIssues.push({ id, issue: '缺失' }); continue }
  if (t.personality.length < 6) report.traitsIssues.push({ id, issue: `personality=${t.personality.length}项(第6阶回退)` })
  if (t.ability.length < 6) report.traitsIssues.push({ id, issue: `ability=${t.ability.length}项(第6阶回退)` })
  const set = new Set([...t.personality, ...t.ability])
  if (set.size < t.personality.length + t.ability.length - 2) report.traitsIssues.push({ id, issue: '阶段内重复特质' })
  for (const p of t.personality) {
    const k = 'P:' + p
    ;(traitDup[k] = traitDup[k] || []).push(id)
  }
}
report.traitDupAcrossSpecies = Object.entries(traitDup).filter(([, ids]) => ids.length > 1).length

// 九维档案
for (const id of speciesIds) {
  const p = PET_PROFILES[id]
  if (!p) { report.profileIssues.push({ id, issue: '缺失' }); continue }
  const stages = p.stages || {}
  const covered = Object.keys(stages).length
  if (covered < 6) report.profileIssues.push({ id, issue: `stages=${covered}/6` })
}

// 12 级数据完整性
for (const sp of allSpecies) {
  const lv = (sp.levels || []).length
  if (lv !== 12) report.levelIssues.push({ id: sp.id, levels: lv })
  for (let i = 1; i <= 12; i++) {
    if (!getLevelData(sp.id, i)) report.levelIssues.push({ id: sp.id, missingLevel: i })
  }
}

fs.writeFileSync(path.join(docsDir, 'pet-audit-report.json'), JSON.stringify(report, null, 2))
console.log('=== 宠物数据审计 ===')
console.log('物种总数:', report.total)
console.log('系列分布:', report.series.map(s => `${s.name}:${s.count}`).join(' / '))
console.log('基础 emoji 重复对:', report.emojiDup.length, JSON.stringify(report.emojiDup).slice(0, 200))
console.log('缺生灵档案:', report.missingStory.length ? report.missingStory.join(',') : '无')
console.log('档案阶段数异常:', report.storyIssues.filter(i => i.issue.startsWith('stages')).length)
console.log('诗词缺失/过短:', report.storyIssues.filter(i => i.issue.includes('诗词')).length)
console.log('台词缺失:', report.storyIssues.filter(i => i.issue.includes('台词')).length)
console.log('品性缺失:', report.storyIssues.filter(i => i.issue.includes('品性')).length)
console.log('剧情互动覆盖:', report.interactionCoverage.covered + '/' + report.interactionCoverage.total)
console.log('道果收官模板化:', report.formulaicFinale.count, report.formulaicFinale.ids.join(','))
console.log('视觉规格缺≥3项:', report.storyIssues.filter(i => i.issue.includes('视觉规格')).length)
console.log('特质缺失/不足:', report.traitsIssues.length, JSON.stringify(report.traitsIssues.slice(0, 10)))
console.log('跨物种重复特质组:', report.traitDupAcrossSpecies)
console.log('九维档案问题:', report.profileIssues.length, JSON.stringify(report.profileIssues.slice(0, 10)))
console.log('12级数据问题:', report.levelIssues.length, JSON.stringify(report.levelIssues.slice(0, 5)))
console.log('同阶段雷同 emoji 组(>3物种):', Object.keys(report.stageEmojiDup).length, JSON.stringify(report.stageEmojiDup).slice(0, 300))
