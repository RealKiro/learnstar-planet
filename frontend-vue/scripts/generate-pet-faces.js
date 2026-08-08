// ===== 学趣星球 · 宠物脸部差异化生成 =====
// 用户反馈：所有角色大圆脸趋同。此脚本为每个角色生成差异化脸谱（脸型轮廓+五官组合），
// 批量替换各角色 6 阶 SVG 中的 <g id="xxx-face">...</g> 模板块。
//
// 用法：node scripts/generate-pet-faces.js
// 说明：只替换 <g id="<species>-face"> 或 <g id="<alias>-face"> 的模板块；
//       一/二阶若未用模板（内联手绘脸），本脚本不改（后续人工）。
// 手/脚比例：脚本同时提供 HAND_FOOT 参数提示，主体比例调整需人工逐张，此处先修脸部。

import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const PET_DIR = path.join(__dirname, '..', 'public', 'pets')

// ---------- 脸谱参数 ----------
// 每角色：face 脸型轮廓 / eye 眼型 / brow 眉型 / nose 鼻型 / mouth 唇型 / skin 肤色 / mark 面部特征
// 脸型：oval 瓜子脸 / square 国字方脸 / long 长脸 / round 圆脸 / baby 娃娃脸 / pointed 尖长脸
// 眼型：phoenix 凤目 / apricot 杏眼 / peach 桃花眼 / triangle 三角眼 / narrow 细长眼 / round 圆眼 / eagle 鹰目 / fox 狐目
// 眉型：sword 剑眉 / straight 一字眉 / willow 柳叶眉 / thick 浓眉 / thin 细眉 / angry 怒眉
// 鼻型：high 高挺 / garlic 蒜头 / hooked 鹰钩 / flat 塌鼻 / pointed 尖鼻
// 唇型：thin 薄唇 / full 厚唇 / pressed 抿唇 / smile 微笑 / cold 冷峭 / cruel 狞笑
const FACE_SPECS = {
  // ===== 东方神话 =====
  sun_wukong:     { face: 'pointed', eye: 'phoenix', brow: 'angry', nose: 'pointed', mouth: 'pressed', skin: '#F2DEC2', mark: '雷公嘴·金瞳' },
  nezha:          { face: 'baby',    eye: 'apricot', brow: 'sword', nose: 'pointed', mouth: 'smile',   skin: '#F6E2CC', mark: '娃娃脸·倔强' },
  zhong_kui:      { face: 'square',  eye: 'eagle',   brow: 'thick', nose: 'high',    mouth: 'cold',    skin: '#E8C8A8', mark: '豹眼·虬髯' },
  jiang_ziya:     { face: 'long',    eye: 'phoenix', brow: 'willow',nose: 'high',    mouth: 'smile',   skin: '#F0D8BC', mark: '白眉·白须' },
  yang_jian:      { face: 'square',  eye: 'eagle',   brow: 'sword', nose: 'high',    mouth: 'pressed', skin: '#F0DCC6', mark: '天眼' },
  lei_zhenzi:     { face: 'pointed', eye: 'round',   brow: 'straight',nose:'pointed',mouth: 'pressed', skin: '#8EACFF', mark: '蓝靛面·鸟喙' },
  huang_tianhua:  { face: 'oval',    eye: 'phoenix', brow: 'sword', nose: 'high',    mouth: 'pressed', skin: '#F2DEC2', mark: '少年·剑眉' },
  tu_xingsun:     { face: 'round',   eye: 'narrow',  brow: 'thin',  nose: 'garlic',  mouth: 'smile',   skin: '#F2DEC2', mark: '胖腮·蒜头鼻' },
  yang_ren:       { face: 'square',  eye: 'phoenix', brow: 'willow',nose: 'high',    mouth: 'pressed', skin: '#F0D8BC', mark: '白眉·掌目' },
  wei_hu:         { face: 'square',  eye: 'eagle',   brow: 'sword', nose: 'high',    mouth: 'pressed', skin: '#F2DEC2', mark: '怒目金刚' },
  daji:           { face: 'oval',    eye: 'fox',     brow: 'willow',nose:'pointed',  mouth: 'smile',   skin: '#F6E2CC', mark: '狐目·朱唇' },
  shen_gongbao:   { face: 'long',    eye: 'eagle',   brow: 'sword', nose: 'hooked',  mouth: 'cold',    skin: '#F0DCC6', mark: '鹰目·薄唇' },
  lv_dongbin:     { face: 'oval',    eye: 'phoenix', brow: 'sword', nose: 'high',    mouth: 'smile',   skin: '#F2DEC2', mark: '剑眉星目' },
  he_xiangu:      { face: 'oval',    eye: 'apricot', brow: 'willow',nose:'pointed',  mouth: 'smile',   skin: '#F6E2CC', mark: '柳眉·杏眼' },
  zhang_guolao:   { face: 'round',   eye: 'narrow',  brow: 'willow',nose:'garlic',   mouth: 'smile',   skin: '#F0D8BC', mark: '白眉长垂' },
  tie_guaili:     { face: 'long',    eye: 'round',   brow: 'thick', nose: 'high',    mouth: 'smile',   skin: '#E0C8AC', mark: '蓬头·眼亮' },
  han_zhongli:    { face: 'square',  eye: 'eagle',   brow: 'thick', nose: 'high',    mouth: 'pressed', skin: '#F0D0B0', mark: '赤发·虬髯' },
  lan_caihe:      { face: 'round',   eye: 'apricot', brow: 'willow',nose:'pointed',  mouth: 'smile',   skin: '#F2DEC2', mark: '少年·弯眉眼' },
  cao_guojiu:     { face: 'oval',    eye: 'peach',   brow: 'willow',nose:'high',     mouth: 'smile',   skin: '#F0D8BC', mark: '温润·桃花眼' },
  taishang_laojun:{ face: 'round',   eye: 'phoenix', brow: 'willow',nose:'high',     mouth: 'smile',   skin: '#F2E2CC', mark: '鹤发童颜' },

  // ===== 星座圣斗士 =====
  aries:          { face: 'oval',    eye: 'apricot', brow: 'sword', nose: 'high',    mouth: 'smile',   skin: '#F2DEC2', mark: '温柔·贤者' },
  taurus:         { face: 'square',  eye: 'eagle',   brow: 'thick', nose: 'garlic',  mouth: 'pressed', skin: '#E8C8A8', mark: '国字·宽鼻' },
  gemini:         { face: 'long',    eye: 'triangle',brow: 'sword', nose: 'hooked',  mouth: 'cold',    skin: '#F2DEC2', mark: '善恶双相' },
  cancer:         { face: 'long',    eye: 'narrow',  brow: 'thin',  nose: 'hooked',  mouth: 'cruel',   skin: '#E8D4BE', mark: '阴柔·狞笑' },
  leo:            { face: 'square',  eye: 'eagle',   brow: 'sword', nose: 'high',    mouth: 'pressed', skin: '#E8C8A8', mark: '虎目·正义' },
  virgo:          { face: 'oval',    eye: 'phoenix', brow: 'willow',nose:'high',     mouth: 'smile',   skin: '#F0E0CC', mark: '闭目·禅相' },
  libra:          { face: 'round',   eye: 'eagle',   brow: 'thick', nose: 'high',    mouth: 'pressed', skin: '#E8C8A8', mark: '少年·虎目' },
  scorpio:        { face: 'square',  eye: 'eagle',   brow: 'sword', nose: 'high',    mouth: 'pressed', skin: '#E8C8A8', mark: '忠义·剑眉' },
  sagittarius:    { face: 'oval',    eye: 'apricot', brow: 'willow',nose:'high',     mouth: 'smile',   skin: '#F0E0CC', mark: '仁厚·星目' },
  capricorn:      { face: 'long',    eye: 'phoenix', brow: 'sword', nose: 'high',    mouth: 'pressed', skin: '#E8C8A8', mark: '剑士·锐利' },
  aquarius:       { face: 'oval',    eye: 'phoenix', brow: 'willow',nose:'high',     mouth: 'pressed', skin: '#F0E0CC', mark: '清冷·凤目' },
  pisces:         { face: 'oval',    eye: 'peach',   brow: 'willow',nose:'pointed',  mouth: 'smile',   skin: '#F6E2CC', mark: '绝美·桃花眼' },
}

// ---------- 脸谱别名映射：模板基础名 → species ----------
// 模板 id 形如 "<基础名>-face" 或 "<基础名>face"；此处建立基础名 → speciesId 映射
const BASE_ALIAS = {
  sun_wukong: 'sun_wukong',
  sunwukong: 'sun_wukong',
  nezha: 'nezha',
  zhong_kui: 'zhong_kui',
  zhongkui: 'zhong_kui',
  jiang_ziya: 'jiang_ziya',
  jiangziya: 'jiang_ziya',
  yang_jian: 'yang_jian',
  yangjian: 'yang_jian',
  lei_zhenzi: 'lei_zhenzi',
  leizhenzi: 'lei_zhenzi',
  huang_tianhua: 'huang_tianhua',
  huangtianhua: 'huang_tianhua',
  tu_xingsun: 'tu_xingsun',
  tuxingsun: 'tu_xingsun',
  yang_ren: 'yang_ren',
  yangren: 'yang_ren',
  wei_hu: 'wei_hu',
  weihu: 'wei_hu',
  daji: 'daji',
  shen: 'shen_gongbao',
  shen_gongbao: 'shen_gongbao',
  shengongbao: 'shen_gongbao',
  lv_dongbin: 'lv_dongbin',
  lvdongbin: 'lv_dongbin',
  he_xiangu: 'he_xiangu',
  hexiangu: 'he_xiangu',
  zhang_guolao: 'zhang_guolao',
  zhangguolao: 'zhang_guolao',
  tie_guaili: 'tie_guaili',
  tieguaili: 'tie_guaili',
  han_zhongli: 'han_zhongli',
  hanzhongli: 'han_zhongli',
  lan_caihe: 'lan_caihe',
  lancaihe: 'lan_caihe',
  cao_guojiu: 'cao_guojiu',
  caoguojiu: 'cao_guojiu',
  taishang_laojun: 'taishang_laojun',
  laojun: 'taishang_laojun',
  aries: 'aries',
  taurus: 'taurus',
  gemini: 'gemini',
  cancer: 'cancer',
  leo: 'leo',
  virgo: 'virgo',
  libra: 'libra',
  scorpio: 'scorpio',
  sag: 'sagittarius',
  sagittarius: 'sagittarius',
  capricorn: 'capricorn',
  aquarius: 'aquarius',
  pisces: 'pisces',
  hongmao: 'hongmao',
}

// ---------- 脸型轮廓 ----------
function faceShape(face) {
  switch (face) {
    case 'oval':   // 瓜子脸：上部圆润、下巴略尖
      return `<ellipse cx="0" cy="0" rx="10.5" ry="11" fill="url(#skin)" stroke="#C8A882" stroke-width="0.8"/>
      <path d="M -8 8 Q 0 12 8 8 Q 5 12 0 12.5 Q -5 12 -8 8 Z" fill="url(#skin)" opacity="0.6"/>`
    case 'square': // 国字方脸：宽额、方下颌
      return `<ellipse cx="0" cy="0" rx="11" ry="10.5" fill="url(#skin)" stroke="#C8A882" stroke-width="0.8"/>
      <path d="M -10 4 L 10 4 L 8 10 L -8 10 Z" fill="url(#skin)" stroke="#C8A882" stroke-width="0.5"/>`
    case 'long':   // 长脸：窄长
      return `<ellipse cx="0" cy="0" rx="9.5" ry="12" fill="url(#skin)" stroke="#C8A882" stroke-width="0.8"/>`
    case 'round':  // 圆脸：宽圆
      return `<ellipse cx="0" cy="0" rx="11" ry="11.5" fill="url(#skin)" stroke="#C8A882" stroke-width="0.8"/>`
    case 'baby':   // 娃娃脸：圆、下巴更圆
      return `<ellipse cx="0" cy="0" rx="10.5" ry="11.5" fill="url(#skin)" stroke="#C8A882" stroke-width="0.8"/>
      <path d="M -7 7 Q 0 12 7 7 Q 4 13 0 13 Q -4 13 -7 7 Z" fill="url(#skin)" opacity="0.6"/>`
    case 'pointed': // 尖长脸：上宽下尖
      return `<ellipse cx="0" cy="0" rx="10" ry="11.5" fill="url(#skin)" stroke="#C8A882" stroke-width="0.8"/>
      <path d="M -7 7 Q 0 13 7 7 Q 4 14 0 14 Q -4 14 -7 7 Z" fill="url(#skin)" opacity="0.5"/>`
    default:
      return `<ellipse cx="0" cy="0" rx="10.5" ry="11" fill="url(#skin)" stroke="#C8A882" stroke-width="0.8"/>`
  }
}

// ---------- 眼型 ----------
function eyeShape(eye) {
  switch (eye) {
    case 'phoenix': // 凤目：细长上挑
      return `<path d="M -8 0 Q -5 -1.8 -2 -0.9 L -2 0.9 Q -5 0 -8 1.2 Z" fill="#FFFDF2" stroke="#1E1008" stroke-width="0.7"/>
      <path d="M 8 0 Q 5 -1.8 2 -0.9 L 2 0.9 Q 5 0 8 1.2 Z" fill="#FFFDF2" stroke="#1E1008" stroke-width="0.7"/>
      <circle cx="-4.8" cy="-0.3" r="1.3" fill="#1E1008"/>
      <circle cx="4.8" cy="-0.3" r="1.3" fill="#1E1008"/>
      <circle cx="-5.3" cy="-0.9" r="0.5" fill="#FFF"/>
      <circle cx="4.3" cy="-0.9" r="0.5" fill="#FFF"/>`
    case 'apricot': // 杏眼：圆润
      return `<ellipse cx="-4.6" cy="0.4" rx="2.4" ry="2.2" fill="#FFFDF2" stroke="#1E1008" stroke-width="0.7"/>
      <ellipse cx="4.6" cy="0.4" rx="2.4" ry="2.2" fill="#FFFDF2" stroke="#1E1008" stroke-width="0.7"/>
      <circle cx="-4.6" cy="0.4" r="1.4" fill="#1E1008"/>
      <circle cx="4.6" cy="0.4" r="1.4" fill="#1E1008"/>
      <circle cx="-5.1" cy="-0.4" r="0.5" fill="#FFF"/>
      <circle cx="4.1" cy="-0.4" r="0.5" fill="#FFF"/>`
    case 'peach':   // 桃花眼：弯月含情
      return `<path d="M -7.5 0.4 Q -4.5 -1.4 -1.8 -0.4 L -1.8 1 Q -4.5 0.2 -7.5 1.8 Z" fill="#1E1008"/>
      <path d="M 7.5 0.4 Q 4.5 -1.4 1.8 -0.4 L 1.8 1 Q 4.5 0.2 7.5 1.8 Z" fill="#1E1008"/>
      <circle cx="-4.5" cy="0.3" r="1.5" fill="#1E1008"/>
      <circle cx="4.5" cy="0.3" r="1.5" fill="#1E1008"/>
      <circle cx="-5" cy="-0.5" r="0.5" fill="#FFF"/>
      <circle cx="4" cy="-0.5" r="0.5" fill="#FFF"/>`
    case 'triangle': // 三角眼：阴鸷
      return `<path d="M -7.5 0.4 Q -4.5 -1.8 -1.8 -0.6 L -1.8 0.9 Q -4.5 0.1 -7.5 1.8 Z" fill="#1E1008"/>
      <path d="M 7.5 0.4 Q 4.5 -1.8 1.8 -0.6 L 1.8 0.9 Q 4.5 0.1 7.5 1.8 Z" fill="#1E1008"/>
      <circle cx="-4.5" cy="0" r="1.6" fill="#A8282E"/>
      <circle cx="4.5" cy="0" r="1.6" fill="#A8282E"/>`
    case 'narrow':  // 细长眼：眯眼
      return `<path d="M -7 0.6 Q -4.5 -1.4 -1.8 -0.4 L -1.8 1.4 Q -4.5 0.6 -7 2.2 Z" fill="#1E1008"/>
      <path d="M 7 0.6 Q 4.5 -1.4 1.8 -0.4 L 1.8 1.4 Q 4.5 0.6 7 2.2 Z" fill="#1E1008"/>`
    case 'round':   // 圆眼：大而正
      return `<circle cx="-4.6" cy="0.4" r="2.8" fill="#1E1008"/>
      <circle cx="4.6" cy="0.4" r="2.8" fill="#1E1008"/>
      <circle cx="-5.3" cy="-0.5" r="1" fill="#FFF"/>
      <circle cx="4" cy="-0.5" r="1" fill="#FFF"/>`
    case 'eagle':   // 鹰目：锐利
      return `<ellipse cx="-4.6" cy="0.4" rx="2.4" ry="2" fill="#1E1008"/>
      <ellipse cx="4.6" cy="0.4" rx="2.4" ry="2" fill="#1E1008"/>
      <circle cx="-5.1" cy="-0.4" r="0.7" fill="#FFE066"/>
      <circle cx="4.1" cy="-0.4" r="0.7" fill="#FFE066"/>`
    case 'fox':     // 狐目：上挑勾魂
      return `<path d="M -7.5 0.4 Q -4.5 -1.6 -1.8 -0.5 L -1.8 0.9 Q -4.5 0.1 -7.5 1.9 Z" fill="#1E1008"/>
      <path d="M 7.5 0.4 Q 4.5 -1.6 1.8 -0.5 L 1.8 0.9 Q 4.5 0.1 7.5 1.9 Z" fill="#1E1008"/>
      <circle cx="-4.5" cy="0.2" r="1.5" fill="#1E1008"/>
      <circle cx="4.5" cy="0.2" r="1.5" fill="#1E1008"/>
      <circle cx="-5" cy="-0.6" r="0.5" fill="#FFF"/>
      <circle cx="4" cy="-0.6" r="0.5" fill="#FFF"/>`
    default:
      return `<ellipse cx="-4.6" cy="0.4" rx="2.4" ry="2.2" fill="#1E1008"/>
      <ellipse cx="4.6" cy="0.4" rx="2.4" ry="2.2" fill="#1E1008"/>
      <circle cx="-5.1" cy="-0.4" r="0.8" fill="#FFF"/>
      <circle cx="4.1" cy="-0.4" r="0.8" fill="#FFF"/>`
  }
}

// ---------- 眉型 ----------
function browShape(brow) {
  switch (brow) {
    case 'sword':   // 剑眉：上扬
      return `<path d="M -8 -1.8 Q -5 -3.4 -2 -2.6 M 8 -1.8 Q 5 -3.4 2 -2.6" stroke="#1E1008" stroke-width="1.4" stroke-linecap="round" fill="none"/>`
    case 'straight': // 一字眉
      return `<path d="M -8 -2 L 8 -2" stroke="#1E1008" stroke-width="1.4" stroke-linecap="round" fill="none"/>`
    case 'willow':  // 柳叶眉：弯细
      return `<path d="M -8 -1.5 Q -5 -3 0 -2.6 Q 5 -3 8 -1.5" stroke="#1E1008" stroke-width="1.1" stroke-linecap="round" fill="none"/>`
    case 'thick':   // 浓眉
      return `<path d="M -8 -1.5 Q -5 -3.2 -2 -2.4 M 8 -1.5 Q 5 -3.2 2 -2.4" stroke="#2E1A10" stroke-width="1.7" stroke-linecap="round" fill="none"/>`
    case 'thin':    // 细眉
      return `<path d="M -8 -1.3 Q -5 -2.4 -2 -2 M 8 -1.3 Q 5 -2.4 2 -2" stroke="#1E1008" stroke-width="0.9" stroke-linecap="round" fill="none"/>`
    case 'angry':   // 怒眉：倒竖
      return `<path d="M -8 -2.6 Q -5 -4 -2 -2.8 M 8 -2.6 Q 5 -4 2 -2.8" stroke="#1E1008" stroke-width="1.5" stroke-linecap="round" fill="none"/>`
    default:
      return `<path d="M -8 -1.5 Q -5 -3.2 -2 -2.4 M 8 -1.5 Q 5 -3.2 2 -2.4" stroke="#1E1008" stroke-width="1.3" stroke-linecap="round" fill="none"/>`
  }
}

// ---------- 鼻型 ----------
function noseShape(nose) {
  switch (nose) {
    case 'high':    // 高挺鼻
      return `<path d="M 0 1.5 L 0 3.6" stroke="#C0A884" stroke-width="1.1" stroke-linecap="round"/>`
    case 'garlic':  // 蒜头鼻
      return `<circle cx="0" cy="3.4" r="2" fill="#C8A880"/>
      <path d="M -2 4.4 Q 0 4.8 2 4.4" stroke="#B8A080" stroke-width="0.8" fill="none"/>`
    case 'hooked':  // 鹰钩鼻
      return `<path d="M 0 1.5 L 0 2.6 Q 1.6 4 0 5 Q -1.6 4 0 2.6" fill="none" stroke="#C0A884" stroke-width="1.1" stroke-linecap="round"/>`
    case 'flat':    // 塌鼻
      return `<circle cx="0" cy="3.4" r="1.8" fill="#C8A880"/>`
    case 'pointed': // 尖鼻
      return `<path d="M 0 1.5 L 0 3 Q 1.2 4.4 0 5.4 Q -1.2 4.4 0 3" fill="none" stroke="#C0A884" stroke-width="1" stroke-linecap="round"/>`
    default:
      return `<path d="M 0 1.5 L 0 3.4" stroke="#C0A884" stroke-width="1" stroke-linecap="round"/>`
  }
}

// ---------- 唇型 ----------
function mouthShape(mouth) {
  switch (mouth) {
    case 'thin':    // 薄唇
      return `<path d="M -3.2 6 Q 0 6.8 3.2 6" fill="none" stroke="#8E5A3E" stroke-width="1" stroke-linecap="round"/>`
    case 'full':    // 厚唇
      return `<path d="M -3.4 5.6 Q 0 8 3.4 5.6 Q 2 7 0 7 Q -2 7 -3.4 5.6 Z" fill="#B06A4A"/>`
    case 'pressed': // 抿唇
      return `<path d="M -3.2 6.2 Q 0 7 3.2 6.2" fill="none" stroke="#8E5A3E" stroke-width="1.2" stroke-linecap="round"/>`
    case 'smile':   // 微笑
      return `<path d="M -3 5.6 Q 0 7.6 3 5.6" fill="none" stroke="#B06A4A" stroke-width="1.1" stroke-linecap="round"/>`
    case 'cold':    // 冷峭
      return `<path d="M -3.2 6.4 Q 0 6.8 3.2 6.4" fill="none" stroke="#6E4A3E" stroke-width="1.1" stroke-linecap="round"/>`
    case 'cruel':   // 狞笑
      return `<path d="M -3.4 5.8 Q 0 8 3.4 5.8" fill="none" stroke="#8E5A3E" stroke-width="1.3" stroke-linecap="round"/>`
    default:
      return `<path d="M -3 5.4 Q 0 7 3 5.4" fill="none" stroke="#8E5A3E" stroke-width="1.1" stroke-linecap="round"/>`
  }
}

// ---------- 面部特征装饰（白眉/白须/血痕/妖纹等，随 mark 附加） ----------
function markShape(mark, spec) {
  const parts = []
  if (mark.includes('白眉')) {
    parts.push(`<path d="M -8 -1.5 Q -5 -3 0 -2.6 Q 5 -3 8 -1.5" stroke="#E8E4DC" stroke-width="1.3" stroke-linecap="round" fill="none"/>`)
    parts.push(`<path d="M -3.5 6.4 Q 0 8.4 3.5 6.4" fill="none" stroke="#E8E4DC" stroke-width="1.2" stroke-linecap="round"/>`)
  }
  if (mark.includes('虬髯')) {
    parts.push(`<path d="M -4.5 7 Q -2 9 0 7.4 M 0 7.4 Q 2 9 4.5 7 M -2 9.2 Q 0 11 2 9.2" stroke="#2E1A10" stroke-width="1.4" stroke-linecap="round" fill="none"/>`)
  }
  if (mark.includes('血痕') || mark.includes('妖纹')) {
    parts.push(`<path d="M -7 -3 Q -6 -1 -7 1" stroke="#B8432E" stroke-width="1.2" stroke-linecap="round" opacity="0.8"/>`)
  }
  if (mark.includes('金瞳')) {
    parts.push(`<circle cx="-4.8" cy="-0.3" r="0.6" fill="#FFE066"/><circle cx="4.8" cy="-0.3" r="0.6" fill="#FFE066"/>`)
  }
  return parts.join('\n      ')
}

// ---------- 生成脸谱模板（保留原始模板 id，只换块内内容） ----------
function generateFaceBlock(templateId, speciesId) {
  const spec = FACE_SPECS[speciesId]
  if (!spec) return null
  return `<g id="${templateId}">
      ${faceShape(spec.face)}
      <!-- ${spec.mark} -->
      ${markShape(spec.mark, spec)}
      <!-- 眉 -->
      ${browShape(spec.brow)}
      <!-- 眼 -->
      ${eyeShape(spec.eye)}
      <!-- 鼻 -->
      ${noseShape(spec.nose)}
      <!-- 唇 -->
      ${mouthShape(spec.mouth)}
    </g>`
}

// ---------- 批量替换 ----------
let replaced = 0, injected = 0, files = 0, errors = 0

// 物种形态角色（脸部本就是动物特征，不注入人形脸谱）
const SPECIES_LIKE = new Set(['daji','zhulong','charmander','hongmao','lantu','doudou','dabeng','tiaotiao','shali','dada','heixinhu','heixiaohu','ma_sanniang','lei_zhenzi'])

function processFile(filePath) {
  const raw = fs.readFileSync(filePath, 'utf8')
  let out = raw
  // 模式 A：文件已有 <g id="...face"> 模板 → 替换块内容（保留模板 id）
  const idMatch = /<g id="([a-zA-Z_\-]+face)">([\s\S]*?)<\/g>/.exec(out)
  if (idMatch) {
    const templateId = idMatch[1]
    const base = templateId.replace(/([-_])?face$/, '')
    const speciesId = BASE_ALIAS[base] || base
    if (!FACE_SPECS[speciesId]) return
    const newBlock = generateFaceBlock(templateId, speciesId)
    if (!newBlock) return
    out = out.replace(/<g id="[a-zA-Z_\-]+face">[\s\S]*?<\/g>/, newBlock)
    if (out !== raw) { fs.writeFileSync(filePath, out); replaced++ }
    return
  }
  // 模式 B：无模板文件（baby/egg 等内联手绘脸）→ 注入差异化脸谱
  // 识别角色名：从文件名 species-stage.svg
  const fname = path.basename(filePath).replace('.svg', '')
  const dash = fname.lastIndexOf('-')
  const speciesId = fname.slice(0, dash)
  if (!FACE_SPECS[speciesId]) return
  if (SPECIES_LIKE.has(speciesId)) return // 物种形态保留
  // 识别头部圆（fill=url(#skin) 的最大 circle）
  const circles = [...out.matchAll(/<circle cx="([\d.]+)" cy="([\d.]+)" r="([\d.]+)" fill="url\(#skin\)"/g)]
    .map(m => ({ x:+m[1], y:+m[2], r:+m[3] }))
  if (!circles.length) return
  const head = circles.sort((a,b) => b.r - a.r)[0]
  if (head.r < 7) return // 头太小（如远处小人）不注入
  // 生成模板块并注入 <defs>
  const faceId = speciesId + '-face'
  const faceBlock = generateFaceBlock(faceId, speciesId)
  // 找到 <defs> 末尾插入
  const defsMatch = /(<defs>[\s\S]*?)(<\/defs>)/.exec(out)
  if (!defsMatch) return
  const defsInsert = `\n    <!-- 差异化脸谱（脚本注入） -->\n    ${faceBlock}\n  `
  out = out.replace(/(<defs>[\s\S]*?)(<\/defs>)/, (m, d1, d2) => d1 + defsInsert + d2)
  // ===== 注入方案：末尾覆盖 + 头发上移 =====
  // 1) 从原位置移除头部圆及周边五官元素；2) 抽出头发元素（深色系/发注释块）；
  // 3) 在 </svg> 前插入 <use>（模板脸覆盖，最后绘制）；4) 把头发元素插到 use 之后（发型在模板之上）。
  // 头部中心附近半径
  const RADIUS = head.r + 5
  // 头发深棕色（发和眼共用，靠"发注释块"识别保留）
  const HAIR_DARK = ['1E1008','2E2620','1E1814','3A2E24','2A1E18','1E1412','2A1A0E','3A2A1E','B8432E','B83A2E','E84E3A','E8E4DC','D8D0C8','C8784E','8E7252','5E4E8E']
  const lines = out.split('\n')
  const hairLines = [] // 抽出的头发元素
  const kept = []
  let pendingHairComment = false // 刚见到头发注释，下一元素若是头顶深色则保留
  for (const line of lines) {
    // 头部圆 → 删除（替换为 use，稍后插入）
    if (new RegExp('<circle cx="'+head.x+'" cy="'+head.y+'" r="'+head.r+'"[^>]*/>').test(line)) continue
    // 头发注释：注释含 发/髻/冠/簪/辫/毛 且不含"头（"（头注释是整体，非头发）
    if (/<!--/.test(line) && /发|髻|冠|簪|辫|毛/.test(line) && !/头（/.test(line)) {
      pendingHairComment = true
      // 注释本身不输出到 kept（避免残留），但记录待抽
      continue
    }
    // 元素行坐标
    let cx, cy
    const c = /<circle cx="([-\d.]+)" cy="([-\d.]+)"/.exec(line)
    const e = /<ellipse cx="([-\d.]+)" cy="([-\d.]+)"/.exec(line)
    const p = /<path d="M ([-\d.]+)[ ,]([-\d.]+)/.exec(line)
    if (c) { cx=+c[1]; cy=+c[2] }
    else if (e) { cx=+e[1]; cy=+e[2] }
    else if (p) { cx=+p[1]; cy=+p[2] }
    const dist = cx!==undefined ? Math.hypot(cx-head.x, cy-head.y) : Infinity
    // 头发：深色 **fill** path（覆盖头顶），五官是 stroke 或小圆
    const isHairFill = /fill="#(1E1008|2E2620|1E1814|3A2E24|2A1E18|1E1412|2A1A0E|3A2A1E|B8432E|B83A2E|E84E3A|E8E4DC|D8D0C8|C8784E|8E7252|5E4E8E)"/.test(line)
    const hasStrokeDark = /stroke="#(1E1008|2E2620)"/.test(line)
    const isHairElement = isHairFill && !hasStrokeDark && cx!==undefined && cy <= head.y + 1 && dist <= RADIUS + 8
    if (pendingHairComment && isHairElement) {
      hairLines.push(line)
      pendingHairComment = false
      continue
    }
    // 若在头发注释后但元素非头发（如眼睛）→ 不抽，重置
    pendingHairComment = false
    // 头顶深色 fill 元素（无注释引导，但明显是发冠）→ 也抽
    if (isHairElement && /<path/.test(line) && dist <= RADIUS + 6) {
      hairLines.push(line)
      continue
    }
    // 头部半径内的其他元素（五官/脸部）→ 删除（被模板覆盖）
    if (cx!==undefined && dist <= RADIUS) {
      // 身体/手臂粗线 → 保留（可能伸入头部下方）
      if (/(stroke-width="[3-9]|stroke-width="1[0-9])/.test(line)) { kept.push(line); continue }
      continue // 其余头部附近元素删除
    }
    kept.push(line)
  }
  // 在 </svg> 前插入 use + 头发
  const useLine = `  <use href="#${faceId}" transform="translate(${head.x},${head.y}) scale(1.05)"/>`
  const hairInsert = hairLines.join('\n')
  const insertBlock = `\n  <!-- 差异化脸谱覆盖（脚本注入） -->\n${useLine}\n${hairInsert ? '\n  <!-- 原头发（保留在模板之上） -->\n' + hairInsert : ''}\n`
  out = kept.join('\n').replace('</svg>', insertBlock + '</svg>')
  if (out !== raw) { fs.writeFileSync(filePath, out); injected++ }
}

function walk(dir) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, entry.name)
    if (entry.isDirectory()) walk(p)
    else if (entry.name.endsWith('.svg')) {
      files++
      try { processFile(p) } catch (e) { errors++; console.error('ERR', p, e.message) }
    }
  }
}

walk(PET_DIR)
console.log(`扫描 SVG: ${files} 个，替换模板: ${replaced} 个，注入脸谱: ${injected} 个，错误: ${errors}`)

// 列出所有角色脸谱规格（供核对）
console.log('\n=== 脸谱分配 ===')
for (const [id, s] of Object.entries(FACE_SPECS)) {
  console.log(`${id.padEnd(16)} ${s.face.padEnd(8)} ${s.eye.padEnd(9)} ${s.brow.padEnd(8)} ${s.nose.padEnd(8)} ${s.mouth.padEnd(8)} ${s.mark}`)
}
