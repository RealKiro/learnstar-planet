// 抽取虹猫角色卡（一次性展示脚本）
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'
import ts from 'typescript'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const frontendDir = path.resolve(__dirname, '..')
const utilsDir = path.join(frontendDir, 'src/utils')
const tmpDir = path.resolve(frontendDir, '..', '.pet-audit-tmp')

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

const ID = process.argv[2] || 'hongmao'
const sp = getSpeciesById(ID)
const story = PET_LIFE_STORIES[ID]
const traits = PET_TRAITS[ID]
const profile = PET_PROFILES[ID]
const series = PET_SERIES.find(s => s.species.some(x => x.id === ID))
const labels = getStageLabelMap(series?.id)
const stageNames = ['卵生', '幼年', '成长', '成熟', '传说', '道果']
const stageLevels = [1, 3, 5, 7, 9, 11]
const stageLevelEnds = [2, 4, 6, 8, 10, 12]

const out = []
out.push(`# 🎴 角色卡：${sp.name}（${ID}）`)
out.push(`系列：${series.name} ｜ 系列 emoji：${series.emoji} ｜ 基础 emoji：${SPECIES_EMOJI[ID]} ｜ 画风绑定：${SERIES_ART_STYLE[series.id]}（${series.id === 'qixia' ? '2D 国漫武侠' : ''}）`)
out.push(`> ${story?.theme || ''}`)
out.push('')

out.push('## 十二级进化链')
out.push('| 等级 | 名称 | 描述 | 累计积分 | 阶段 emoji |')
out.push('|---|---|---|---|---|')
for (const l of sp.levels || []) {
  out.push(`| L${l.level} | ${l.name} | ${l.description} | ${l.requiredScore} | ${getStageEmoji(ID, l.level)} |`)
}
out.push('')

out.push('## 六阶段 · 剧情与特质')
for (let i = 0; i < 6; i++) {
  const st = story?.stages?.[i]
  if (!st) continue
  const lv = stageLevels[i]
  const lvEnd = stageLevelEnds[i]
  out.push(`### ${stageNames[i]}（L${lv}-${lvEnd}）· ${labels[Object.keys(labels).length ? ['egg', 'baby', 'growing', 'mature', 'legendary', 'transcendent'][i] : 'egg']}`)
  out.push(`- **阶段名**：${st.name}（${st.age}）｜ 关键词：${st.keyword}`)
  out.push(`- **品性**：${st.character}`)
  out.push(`- **互动/冲突**：${st.interaction || '（待撰写）'}`)
  out.push(`- **动作**：${st.action}`)
  out.push(`- **服饰**：${st.attire}${st.attireDetail ? '；' + st.attireDetail : ''}`)
  out.push(`- **功法**：${st.technique}｜ **心境**：${st.mood}`)
  out.push(`- **台词**：「${st.line}」`)
  out.push(`- **诗文**：${st.poem}`)
  out.push(`- **特质**：${traits?.personality?.[i] ?? '—'} ｜ **技能**：${traits?.ability?.[i] ?? '—'}`)
  out.push(`- **画面**：${st.svg}`)
  if (st.body || st.face || st.attireDetail || st.hairStyle || st.powerEffect) {
    out.push(`- **视觉规格**：${[st.body, st.face, st.hairStyle, st.attireDetail, st.powerEffect].filter(Boolean).join(' / ')}`)
  }
  out.push(`- **AI 提示词（动漫风）**：${composeStageAiPrompt(ID, lv, 'anime')}`)
  out.push('')
}

out.push('## 九维档案（当前回退全局）')
if (profile) {
  out.push(`形：${profile.form} ｜ 习性：${profile.habit} ｜ 出身：${profile.origin}`)
  out.push(`雅号：${profile.epithet} ｜ 动作：${profile.movement} ｜ 象征：${profile.symbol} ｜ 主题：${profile.theme}`)
  if (profile.stages) {
    for (const [k, v] of Object.entries(profile.stages)) {
      out.push(`- ${k}：${v.form} ｜ ${v.theme}`)
    }
  }
}

fs.writeFileSync(path.join(frontendDir, 'hongmao-card.md'), out.join('\n'))
console.log(out.join('\n'))
