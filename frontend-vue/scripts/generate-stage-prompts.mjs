// 全量阶段 AI 提示词导出：按系列绑定的画风（126 物种 × 6 阶段）→ docs/stage-ai-prompts/{styleId}.md
// 用法：cd frontend-vue && node scripts/generate-stage-prompts.mjs
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'
import ts from 'typescript'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const frontendDir = path.resolve(__dirname, '..')
const repoRoot = path.resolve(frontendDir, '..')
const utilsDir = path.join(frontendDir, 'src/utils')
const outDir = path.join(repoRoot, 'docs/stage-ai-prompts')

const FILES = ['petDataExtended', 'petData', 'petLifeStories', 'petTraits', 'stageEmoji', 'petArtStyles', 'seriesStages', 'stageAiPrompt']
const tmpDir = path.resolve(repoRoot, '..', '.pet-audit-tmp')
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
const { PET_SERIES } = require2('./petData.js')
const { composeStageAiPrompt } = require2('./stageAiPrompt.js')
const { SERIES_ART_STYLE } = require2('./seriesStages.js')
const { PET_ART_STYLES } = require2('./petArtStyles.js')

const allSpecies = PET_SERIES.flatMap(s => s.species.map(sp => ({ ...sp, seriesName: s.name })))
const stageNames = ['卵生', '幼年', '成长', '成熟', '传说', '道果']

// 按画风分组（画风由系列绑定）
const byStyle = new Map()
for (const s of PET_SERIES) {
  const styleId = SERIES_ART_STYLE[s.id] || 'anime'
  if (!byStyle.has(styleId)) byStyle.set(styleId, { styleName: PET_ART_STYLES.find(x => x.id === styleId)?.name || styleId, series: [] })
  byStyle.get(styleId).series.push(s)
}

fs.mkdirSync(outDir, { recursive: true })
let total = 0
for (const [styleId, { styleName, series }] of byStyle) {
  const styleMeta = PET_ART_STYLES.find(x => x.id === styleId)
  const lines = [
    `# 学宠星球 · 阶段 AI 生图提示词（画风：${styleName}）`,
    '',
    `> 本画风覆盖系列：${series.map(s => s.name).join(' / ')}`,
    `> 画风行：${styleMeta?.cn ?? ''} / ${styleMeta?.en ?? ''}`,
    '> 统一约束已拼入每条提示词：1024×1024 PNG、透明背景、单主体居中 3/4 视角、全身完整、无文字水印。',
    `> 同系列角色的六个阶段使用同一画风，保证角色一致性。`,
    '',
  ]
  let count = 0
  for (const s of series) {
    lines.push(`## 系列：${s.name}（${s.id}）`)
    for (const sp of s.species) {
      lines.push(`### ${sp.name}（${sp.id}）`)
      ;[1, 3, 5, 7, 9, 12].forEach((lv, i) => {
        lines.push(`- **L${lv} ${stageNames[i]}**：${composeStageAiPrompt(sp.id, lv, styleId)}`)
      })
      lines.push('')
      count += 6
      total += 6
    }
  }
  const out = path.join(outDir, `${styleId}.md`)
  fs.writeFileSync(out, lines.join('\n'))
  console.log(`✓ ${out}（${count} 条，系列：${series.map(s => s.name).join('/')}）`)
}
console.log(`合计 ${total} 条提示词，输出目录 docs/stage-ai-prompts/`)
