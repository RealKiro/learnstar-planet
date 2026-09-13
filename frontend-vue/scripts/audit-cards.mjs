// 角色卡专项审计：契合度人工核对表 + 唯一性全量检测
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'
import ts from 'typescript'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const frontendDir = path.resolve(__dirname, '..')
const repoRoot = path.resolve(frontendDir, '..')
const utilsDir = path.join(frontendDir, 'src/utils')
const tmpDir = path.resolve(repoRoot, '..', '.pet-audit-tmp')

const FILES = ['petDataExtended', 'petData', 'petLifeStories', 'petTraits', 'stageEmoji']
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

// 1b. 拷贝系列数据 JSON（TS wrapper 经 require('../data/pets/*.json') 引用）
const jsonSrcDir = path.join(utilsDir, '..', 'data', 'pets')
fs.mkdirSync(path.join(tmpDir, 'data', 'pets'), { recursive: true })
for (const f of fs.readdirSync(jsonSrcDir)) fs.copyFileSync(path.join(jsonSrcDir, f), path.join(tmpDir, 'data', 'pets', f))
const require2 = createRequire(path.join(tmpDir, 'utils', 'petData.js'))
const { PET_SERIES, SPECIES_EMOJI } = require2('./petData.js')
const { getStageEmoji } = require2('./stageEmoji.js')
const { PET_TRAITS } = require2('./petTraits.js')

const allSpecies = PET_SERIES.flatMap(s => s.species)
const stageNames = ['卵', '幼年', '成长', '成熟', '传说', '道果']

// ① 契合度人工核对表（输出到文件供人工审阅）
const fitLines = []
for (const s of PET_SERIES) {
  fitLines.push(`## ${s.name}（${s.id}）`)
  for (const sp of s.species) {
    fitLines.push(
      `${sp.name} | 基础 ${SPECIES_EMOJI[sp.id]} | 传说 ${getStageEmoji(sp.id, 9)} | 道果 ${getStageEmoji(sp.id, 12)} | L1 ${sp.levels?.[0]?.name ?? ''} / L12 ${(sp.levels || []).slice(-1)[0]?.name ?? ''}`
    )
  }
}
fs.writeFileSync(path.join(frontendDir, 'card-fit-review.txt'), fitLines.join('\n'))

// ② 各阶段 emoji 跨物种重复（L1 卵生统一 🥚 为设计，跳过）
console.log('==== ① 阶段 emoji 跨物种重复（排除卵生） ====')
for (const lv of [3, 5, 7, 9, 12]) {
  const seen = {}
  for (const sp of allSpecies) {
    const e = getStageEmoji(sp.id, lv)
    ;(seen[e] = seen[e] || []).push(sp.name)
  }
  for (const [e, names] of Object.entries(seen)) {
    if (names.length > 1) console.log(`L${lv} ${e} → ${names.join(', ')}`)
  }
}

// ③ 12 级名称 / 描述跨物种重复
function checkDup(getter, label) {
  const seen = {}
  for (const sp of allSpecies) {
    for (const l of sp.levels || []) {
      const v = (getter(l) || '').trim()
      if (!v) continue
      ;(seen[v] = seen[v] || []).push(`${sp.name}·L${l.level}`)
    }
  }
  const dups = Object.entries(seen).filter(([, v]) => v.length > 1)
  console.log(`\n==== ${label}：${dups.length} 组重复 ====`)
  for (const [v, where] of dups.slice(0, 40)) console.log(`  「${v}」→ ${where.join(' / ')}`)
  if (dups.length > 40) console.log(`  ...共 ${dups.length} 组`)
}
checkDup(l => l.name, '等级名称跨物种重复')
checkDup(l => l.description, '等级描述跨物种重复')

// ④ 各阶段 personality / ability 跨物种重复
console.log('\n==== 各阶段特质跨物种重复 ====')
let dupCount = 0
for (const kind of ['personality', 'ability']) {
  for (let si = 0; si < 6; si++) {
    const seen = {}
    for (const [id, t] of Object.entries(PET_TRAITS)) {
      const v = (t[kind] || [])[si]
      if (!v) continue
      ;(seen[v] = seen[v] || []).push(id)
    }
    for (const [v, ids] of Object.entries(seen)) {
      if (ids.length > 1) { dupCount++; console.log(`${kind}#${stageNames[si]} 「${v}」→ ${ids.join(', ')}`) }
    }
  }
}
console.log(`特质跨物种重复合计：${dupCount} 组`)
