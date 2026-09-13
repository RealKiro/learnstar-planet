// README 内链锚点校验（一次性）：解析标题 slug 与全部内链，报告断链
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const p = path.resolve(__dirname, '../../README.md')
const t = fs.readFileSync(p, 'utf8')

const slugify = s =>
  s.trim().toLowerCase()
    .replace(/[\u{1F300}-\u{1FAFF}\u{2600}-\u{27BF}\u{FE0F}\u{200D}]/gu, '') // 去 emoji
    .replace(/[^\p{L}\p{N}\s-]/gu, '') // 去标点（保留 CJK/字母数字/空格/连字符）
    .replace(/\s/g, '-') // 每个空格各转一个连字符（GitHub 规则："备份 / 恢复" → "备份--恢复"）

const headings = new Set()
for (const m of t.matchAll(/^(#{1,6})\s+(.+)$/gm)) headings.add(slugify(m[2]))

const links = []
for (const m of t.matchAll(/\]\(#([^)]+)\)/g)) links.push(decodeURIComponent(m[1]))
for (const m of t.matchAll(/href="#([^"]+)"/g)) links.push(decodeURIComponent(m[1]))

// 上面的正则提取出的锚点不含 '#'，直接比对即可。
// ⚠️ 旧版这里误写了 l.startsWith('#')（恒为假），导致断链永远报不出来。
const broken = [...new Set(links)].filter(l => !headings.has(l))
console.log('标题数:', headings.size, '| 内链数:', links.length)
console.log(broken.length ? '断链:\n' + broken.map(b => '  #' + b).join('\n') : '内链锚点全部有效 ✓')
