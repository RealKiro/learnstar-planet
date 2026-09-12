// 全量角色卡生成器：126 物种 → docs/角色卡/{系列名}/{speciesId}.md
// 用法：node scripts/generate-cards.mjs [speciesId]（缺省生成全部）
// 每张卡含：基本信息 / 12 级进化链 / 六阶段（剧情正文 + 互动冲突 + 台词诗文 + 特质 + 视觉 + AI 提示词）/ 九维档案
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'
import ts from 'typescript'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const frontendDir = path.resolve(__dirname, '..')
const repoRoot = path.resolve(frontendDir, '..')
const docsDir = path.join(repoRoot, 'docs')
const utilsDir = path.join(frontendDir, 'src/utils')
const tmpDir = path.resolve(repoRoot, '..', '.pet-audit-tmp')

const FILES = ['petDataExtended', 'petData', 'petLifeStories', 'petProfiles', 'petTraits', 'stageEmoji', 'seriesStages', 'petArtStyles', 'stageAiPrompt']
fs.rmSync(tmpDir, { recursive: true, force: true })
fs.mkdirSync(path.join(tmpDir, 'utils'), { recursive: true })
fs.writeFileSync(path.join(tmpDir, 'package.json'), JSON.stringify({ type: 'commonjs' }))
for (const name of FILES) {
  const src = fs.readFileSync(path.join(utilsDir, `${name}.ts`), 'utf8')
  const js = ts.transpileModule(src, {
    compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020, esModuleInterop: true },
  }).outputText
  fs.writeFileSync(path.join(tmpDir, 'utils', `${name}.js`), js)
}
const require2 = createRequire(path.join(tmpDir, 'utils', 'petData.js'))
const { PET_SERIES, SPECIES_EMOJI, getSpeciesById } = require2('./petData.js')
const { PET_LIFE_STORIES } = require2('./petLifeStories.js')
const { PET_PROFILES } = require2('./petProfiles.js')
const { PET_TRAITS } = require2('./petTraits.js')
const { getStageEmoji } = require2('./stageEmoji.js')
const { getStageLabelMap, SERIES_ART_STYLE } = require2('./seriesStages.js')
const { composeStageAiPrompt } = require2('./stageAiPrompt.js')

const allSpecies = PET_SERIES.flatMap(s => s.species.map(sp => ({ ...sp, seriesName: s.name })))
const only = process.argv[2] || ''
const targets = only ? allSpecies.filter(sp => sp.id === only) : allSpecies
if (!targets.length) { console.error('未找到物种：' + only); process.exit(1) }

const stageNames = ['卵生', '幼年', '成长', '成熟', '传说', '道果']
const stageLevels = [1, 3, 5, 7, 9, 11]
const stageLevelEnds = [2, 4, 6, 8, 10, 12]
const stageKeys = ['egg', 'baby', 'growing', 'mature', 'legendary', 'transcendent']

// 剧情正文组装：把结构化字段织成一段连贯的叙事
function sent(s) {
  const v = (s || '').trim()
  if (!v) return ''
  return /[。！？…」』]$/.test(v) ? v : v + '。'
}
function composeNarrative(st, spName) {
  const parts = [
    sent(`${st.age}，${spName}——${st.name}，${st.keyword}`),
    sent(st.character),
    sent(st.action),
    st.interaction ? sent('※ ' + st.interaction) : '',
    st.technique ? sent('所习：' + st.technique) : '',
    sent('心境：' + st.mood),
  ]
  return parts.filter(Boolean).join('')
}

function buildCard(sp) {
  const series = PET_SERIES.find(s => s.species.some(x => x.id === sp.id))
  const story = PET_LIFE_STORIES[sp.id]
  const traits = PET_TRAITS[sp.id]
  const profile = PET_PROFILES[sp.id]
  const labels = getStageLabelMap(series?.id)
  const artStyle = SERIES_ART_STYLE[series?.id] || 'anime'
  const L = []
  L.push(`# 🎴 角色卡：${sp.name}（${sp.id}）`)
  L.push(``)
  L.push(`> 系列：${series?.name} ｜ 基础 emoji：${SPECIES_EMOJI[sp.id] || '—'} ｜ 画风绑定：${artStyle} ｜ 主题：${story?.theme || '—'}`)
  L.push(``)
  L.push(`## 1. 十二级进化链`)
  L.push(``)
  L.push(`| 等级 | 名称 | 描述 | 累计积分 | 阶段 emoji |`)
  L.push(`|---|---|---|---|---|`)
  for (const l of sp.levels || []) {
    L.push(`| L${l.level} | ${l.name} | ${l.description} | ${l.requiredScore} | ${getStageEmoji(sp.id, l.level)} |`)
  }
  L.push(``)
  L.push(`## 2. 六阶段 · 剧情正文与角色特质`)
  for (let i = 0; i < 6; i++) {
    const st = story?.stages?.[i]
    if (!st) continue
    const key = stageKeys[i]
    L.push(``)
    L.push(`### 2.${i + 1} ${stageNames[i]}（L${stageLevels[i]}-${stageLevelEnds[i]}）· ${st.name}（${st.age}）`)
    L.push(``)
    L.push(`> 关键词：${st.keyword} ｜ 心境：${st.mood}`)
    L.push(``)
    L.push(`**剧情正文**：${composeNarrative(st, sp.name)}`)
    L.push(``)
    L.push(`**互动 / 冲突**：${st.interaction || '（待撰写——本阶段与其他角色的互动/冲突剧情）'}`)
    L.push(``)
    if (st.line) L.push(`**台词**：「${st.line}」`)
    if (st.poem) L.push(`**诗文**：${st.poem}`)
    L.push(`**特质**：${traits?.personality?.[i] ?? '—'} ｜ **技能**：${traits?.ability?.[i] ?? '—'}`)
    L.push(``)
    L.push(`**画面描述**：${st.svg || '—'}`)
    if (st.body || st.face || st.hairStyle || st.attireDetail || st.powerEffect) {
      L.push(`**视觉规格**：${[st.body, st.face, st.hairStyle, st.attireDetail, st.powerEffect].filter(Boolean).join(' / ')}`)
    }
    L.push(`**AI 提示词（${artStyle}）**：${composeStageAiPrompt(sp.id, stageLevels[i], artStyle)}`)
  }
  if (profile) {
    L.push(``)
    L.push(`## 3. 九维档案`)
    L.push(``)
    L.push(`形：${profile.form} ｜ 习性：${profile.habit} ｜ 出身：${profile.origin}`)
    L.push(`雅号：${profile.epithet} ｜ 动作：${profile.movement} ｜ 象征：${profile.symbol} ｜ 主题：${profile.theme}`)
    if (profile.stages) {
      for (const [k, v] of Object.entries(profile.stages)) {
        L.push(`- **${k}**：${v.form} ｜ ${v.theme}`)
      }
    }
  }
  return L.join('\n')
}

let ok = 0
for (const sp of targets) {
  const series = PET_SERIES.find(s => s.species.some(x => x.id === sp.id))
  const dir = path.join(docsDir, '角色卡', series?.name || '未分组')
  fs.mkdirSync(dir, { recursive: true })
  fs.writeFileSync(path.join(dir, `${sp.id}.md`), buildCard(sp))
  ok++
}
console.log(`✓ 已生成 ${ok} 张角色卡 → docs/角色卡/{系列名}/`)
