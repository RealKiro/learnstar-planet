// ===== 学趣星球 · 宠物角色档案（形态 / 习性） =====
// 按物种 id 索引；缺失时返回兜底，保证所有宠物图鉴可展示。
// 专属诗文见 petHandbookData.getPoems，进化台词见 getEvoLines。

export interface PetProfile {
  /** 形态 */
  form: string
  /** 习性 */
  habit: string
}

export const PET_PROFILES: Record<string, PetProfile> = {
  // ===== 山海经 =====
  zhulong:       { form: '人面蛇身，通体赤红，口衔火烛。', habit: '视为昼，瞑为夜，吹为冬，呼为夏，独掌昼夜与四季。' },
  yinglong:      { form: '身披金鳞，背生双翼的神龙。', habit: '主云雨，曾助黄帝斩蚩尤，性刚烈忠勇，风雨随行。' },
  nine_tail_fox: { form: '通体雪白，身后拖九条长尾。', habit: '居于青丘之山，善魅惑亦为祥瑞，月光之下最是活跃。' },
  kunpeng:       { form: '北冥之巨鱼，可化身为鸟。', habit: '怒而飞则翼若垂天之云，扶摇而上九万里，遨游海空。' },
  fenghuang:     { form: '五彩华羽，尾羽如虹。', habit: '非梧桐不栖，非竹实不食，非醴泉不饮，性高洁至傲。' },
  qilin:         { form: '麋身牛尾，马蹄而一角。', habit: '不践生草，不履生虫，乃仁兽，现世则天下太平。' },
  qiongqi:       { form: '状如虎而生双翼，遍体凶煞。', habit: '喜食善人而恶人予饲，颠倒黑白，为四凶之一。' },
  bifang:        { form: '青色单足之鸟，白喙赤足。', habit: '鸣声如鹤，见则其邑有讹火，乃火灾之兆。' },
  pixiu:         { form: '龙头马身麟脚，背生双翼。', habit: '只进不出，聚财纳福，镇宅辟邪，人皆喜爱。' },
  jingwei:       { form: '形如乌，赤首白喙赤足。', habit: '炎帝之女所化，常衔西山木石以填东海，至死不渝。' },
  xiangliu:      { form: '九首蛇身，自环其尾。', habit: '食于九土，所过之处皆化为泽水，乃共工之臣。' },
  xiezhi:        { form: '似羊而独角，目光清正。', habit: '能辨曲直，见人争斗则触其不直者，司掌公正。' },

  // ===== 数码宝贝 =====
  mecha_dragon:  { form: '橙色小恐龙，背甲坚挺。', habit: '贪吃好动，最重友情，一遇危险便喷出小型火焰。' },
  cyber_cat:     { form: '白色圣猫，耳戴神圣之环。', habit: '外表高傲实则忠诚，爪击迅捷，身怀光明之力。' },
  space_mecha:   { form: '奶油色圆身，扑腾小翅膀。', habit: '天真乐观，爱在空中翻滚，是希望的化身。' },
  quantum_beast: { form: '身披蓝毛的小狼，额生独角。', habit: '腼腆胆小，一旦信任便倾尽所有，友情之力深藏。' },
  digital_phoenix:{ form: '粉色小鸟，头顶凤冠。', habit: '最爱唱歌，歌声悦耳能抚慰人心，纯真动人。' },
  mecha_shark:   { form: '白色海狮，圆滚滚胖乎乎。', habit: '爱撒娇爱卖萌，率真直爽，诚实可靠。' },

  // ===== 龙系 =====
  wyvern:        { form: '赤色飞龙，蝠翼如膜，尾带尖刺。', habit: '翱翔于魔界火山，凶猛善战，以喷吐烈焰著称。' },
  chinese_alligator: { form: '低伏的宽吻鳄龙，身披硬鳞。', habit: '潜伏于江河泥潭，静若处子，动若雷霆。' },
  xuanwu:        { form: '龟蛇合体，背负硬壳。', habit: '镇守北方，性沉稳坚韧，遇水则灵。' },
  t_rex:         { form: '直立霸王龙，头大颚壮，前肢短小。', habit: '称霸白垩纪的顶级掠食者，吼声震天。' },
  spinosaurus:   { form: '背生巨大帆状脊，形似鳄龙。', habit: '栖于水边，捕鱼为生，帆可调节体温。' },
  ankylosaurus:  { form: '身披重甲，尾端生骨锤。', habit: '行动缓慢但防御无双，一锤可碎石。' },
  diplodocus:    { form: '长颈细尾的巨型蜥脚龙。', habit: '以高树之叶为食，性情温和，群居而游。' },
}

export function getPetProfile(speciesId: string): PetProfile {
  return PET_PROFILES[speciesId] || { form: '形态尚待探明的神秘宠物。', habit: '习性未知，正等待与小主人一起探索。' }
}
