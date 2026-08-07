// ===== 学趣星球 · 宠物初始形态设计表（灵胎初醒 / Lv1） =====
// 设计原则（用户定稿）：
//   1. 初始形态符合物种特点 —— 植物=种子 / 胎生哺乳=胚胎 / 卵生=蛋卵 / 数码兽=数码蛋 / 神兽=灵种灵卵 / 食物=原料 / 器物=物件
//   2. 每个角色初始形态有属性区别（木火土金水 + 雷冰风光暗钢秘 十二属性）
//   3. 形态名建议同步 petData 的 Lv1 等级名（保持一致后生图/UI/诗文三方对齐）
// 字段：cat 分类 / elem 五行属性 / name 初始形态名 / desc 初始形态描述（提示词 egg 阶段用）
// 缺失物种时返回通用兜底。

export interface PetInitial {
  cat: 'seed' | 'embryo' | 'egg' | 'digi' | 'spirit' | 'food' | 'artifact' | 'sword' | 'person' | 'mythic'
  elem: '木' | '火' | '土' | '金' | '水' | '雷' | '冰' | '风' | '光' | '暗' | '钢' | '秘'
  name: string
  desc: string
}

/** 十二属性视觉特征（供 desc 点缀 / 生成 SVG 时取色） */
export const INITIAL_ELEMENT_GLOSS: Record<PetInitial['elem'], { accent: string; glow: string; motif: string }> = {
  木: { accent: '#4ADE80', glow: '#DCFCE7', motif: '青色芽纹、嫩绿生机' },
  火: { accent: '#F97316', glow: '#FFEDD5', motif: '赤焰纹路、暖金光晕' },
  土: { accent: '#D97706', glow: '#FEF3C7', motif: '大地岩纹、赭黄微光' },
  金: { accent: '#EAB308', glow: '#FEF9C3', motif: '鎏金纹路、金属冷光' },
  水: { accent: '#38BDF8', glow: '#E0F2FE', motif: '水波鳞纹、潮蓝微光' },
  雷: { accent: '#A78BFA', glow: '#EDE9FE', motif: '细碎电纹、噼啪电光' },
  冰: { accent: '#7DD3FC', glow: '#E0F2FE', motif: '霜花晶纹、寒蓝微光' },
  风: { accent: '#93C5FD', glow: '#EFF6FF', motif: '流线旋纹、青白微光' },
  光: { accent: '#FDE68A', glow: '#FEFCE8', motif: '芒星光晕、圣白微光' },
  暗: { accent: '#A78BFA', glow: '#2E1065', motif: '幽暗光晕、紫黑微光' },
  钢: { accent: '#94A3B8', glow: '#F1F5F9', motif: '钢铁装甲纹、银灰冷光' },
  秘: { accent: '#C4B5FD', glow: '#F5F3FF', motif: '奥术符文、星芒法阵' },
}

export const PET_INITIAL: Record<string, PetInitial> = {
  // ===== 山海经（神兽灵种 / 灵卵） =====
  zhulong:       { cat: 'mythic', elem: '火', name: '烛火微光', desc: '一粒烛火灵种，火光如豆却照彻幽冥，烛龙始祖的气息在火苗中蜷曲，昼夜之神初醒' },
  yinglong:      { cat: 'egg', elem: '水', name: '虺卵', desc: '一枚虺卵，青灰色蛋壳覆着细密金鳞纹，隐隐透出风雨雷光，是应龙化形前的沉睡之卵' },
  nine_tail_fox: { cat: 'mythic', elem: '水', name: '灵狐蛋', desc: '一枚灵狐蛋，月白蛋壳泛着清辉，尾尖般的银色光毫在壳上流转，九尾之灵伏卧其中' },
  kunpeng:       { cat: 'mythic', elem: '水', name: '鲲苗', desc: '一尾鲲苗，通体青蓝如墨玉，在混沌灵光中游弋，鳍翼初具而鹏形未显' },
  fenghuang:     { cat: 'egg', elem: '火', name: '火羽蛋', desc: '一枚火羽蛋，赤色蛋壳缀满五色羽纹，暖光从壳内透出，凤鸣在蛋中回响' },
  qilin:         { cat: 'mythic', elem: '土', name: '瑞兽玉', desc: '一枚温润瑞兽玉，玉身内蕴祥云纹，麟角初形若隐若现，仁兽之息如春风化雨' },
  qiongqi:       { cat: 'egg', elem: '土', name: '凶兽卵', desc: '一枚凶兽卵，暗色蛋壳覆着虎纹与翼影，煞气凝成黑气在壳面游走' },
  bifang:        { cat: 'egg', elem: '火', name: '毕方卵', desc: '一枚毕方卵，青白蛋壳上燃着一点赤红火纹，单足立火的雏形在壳中孕育' },
  pixiu:         { cat: 'mythic', elem: '金', name: '貔貅玉', desc: '一枚鎏金貔貅玉，玉兽张嘴吞吐财气，金光明灭间纳福之相初成' },
  jingwei:       { cat: 'egg', elem: '木', name: '精卫卵', desc: '一枚精卫卵，浅赤蛋壳缀着白羽纹，衔枝填海的执念在壳中轻轻鼓动' },
  xiangliu:      { cat: 'egg', elem: '水', name: '九首卵', desc: '一枚九首卵，蛇纹盘绕的暗青蛋壳，隐隐裂出九道微光，毒水之气在壳下涌动' },
  xiezhi:        { cat: 'egg', elem: '金', name: '獬豸卵', desc: '一枚獬豸卵，雪白蛋壳顶生一点金角凸起，断狱明理的法光在壳面流转' },
  qinglong:      { cat: 'mythic', elem: '木', name: '木灵珠', desc: '一粒木灵珠，青碧如春水凝成，珠身缠着鹿角须影，东方之木生机勃发' },
  baihu:         { cat: 'egg', elem: '金', name: '金灵卵', desc: '一枚金灵卵，雪白金纹蛋壳，虎啸之声在金属冷光中隐隐回响' },
  zhuque:        { cat: 'egg', elem: '火', name: '火灵卵', desc: '一枚火灵卵，赤羽纹蛋壳透出烈焰光晕，离火之精在壳中振翅欲鸣' },
  xuanwu:        { cat: 'mythic', elem: '水', name: '水灵珠', desc: '一粒水灵珠，玄青珠身内蕴龟蛇之影，北冥之水在其间缓缓流转' },
  taotie:        { cat: 'egg', elem: '土', name: '食欲之卵', desc: '一枚饕餮卵，暗金蛋壳覆着贪婪纹，巨口虚影在壳面时隐时现，永不知足' },
  baize:         { cat: 'mythic', elem: '金', name: '智慧之卵', desc: '一枚智慧之卵，素白蛋壳上浮现问号与书页纹，万灵之名在壳中等待被念出' },

  // ===== 东方神话（人物灵光化形 / 修行起点） =====
  jiang_ziya:    { cat: 'person', elem: '水', name: '灵光种子', desc: '一缕灵光凝成垂钓老者的剪影，渭水之畔愿者上钩，封神之榜在其身侧浮沉' },
  nezha:         { cat: 'person', elem: '火', name: '灵珠种子', desc: '一团赤红灵珠之光，火尖枪的锋芒在其中跳跃，莲花初生的轮廓若隐若现' },
  yang_jian:     { cat: 'person', elem: '金', name: '天眼灵光', desc: '一道天眼灵光，银甲虚影与额间竖瞳的星芒若隐若现，二郎神威初显' },
  lei_zhenzi:    { cat: 'person', elem: '雷', name: '雷卵种子', desc: '一枚雷纹灵卵，紫电纹路在蛋壳上蜿蜒，风雷双翅的雏形在雷光中孕育' },
  huang_tianhua: { cat: 'person', elem: '土', name: '将门灵种', desc: '一粒将门灵种，玉麒麟的蹄影与莫邪剑的寒光在土色灵光中交叠' },
  tu_xingsun:    { cat: 'person', elem: '土', name: '地灵种子', desc: '一粒地灵种子，黄土灵光凝成遁形之影，镔铁棍的轮廓埋在地脉之中' },
  yang_ren:      { cat: 'person', elem: '木', name: '忠臣灵光', desc: '一缕忠臣灵光，双掌之目在光影中睁开，掌中星芒是剖目重见天日的光' },
  wei_hu:        { cat: 'person', elem: '金', name: '金刚灵种', desc: '一粒金刚灵种，金甲虚影与降魔杵的宝光交错，护法威严初凝' },
  daji:          { cat: 'person', elem: '暗', name: '妖狐之卵', desc: '一枚妖狐之卵，绯红蛋壳泛着魅惑幽光，九尾虚影在狐火中摇曳' },
  shen_gongbao:  { cat: 'person', elem: '风', name: '邪道灵光', desc: '一缕邪道灵光，双剑虚影与黑虎之息交织，三寸舌的巧言在风中回响' },
  sun_wukong:    { cat: 'person', elem: '土', name: '灵根育孕', desc: '一粒灵石灵根，花果山仙气凝成石卵，石中孕育着金毛猴王的轮廓' },
  lv_dongbin:    { cat: 'person', elem: '金', name: '黄粱书生', desc: '一位黄粱书生剪影，青衫负剑立于道袍灵光中，纯阳剑气如露未凝' },
  he_xiangu:     { cat: 'person', elem: '木', name: '岭南采药女', desc: '一位采药女虚影，竹篓药香凝成青雾，莲花在足下缓缓绽放' },
  zhang_guolao:  { cat: 'person', elem: '土', name: '恒山樵夫', desc: '一位恒山樵夫虚影，背负柴薪倒骑白驴，鱼鼓之声在灵光中轻响' },
  tie_guaili:    { cat: 'person', elem: '水', name: '巴国隐士', desc: '一位巴国隐士虚影，铁拐顿地破衣飘摇，葫芦灵光纳着一方天地' },
  han_zhongli:   { cat: 'person', elem: '火', name: '燕台武将', desc: '一位燕台武将虚影，红袍赤面战意未消，芭蕉扇上焰光初燃' },
  lan_caihe:     { cat: 'person', elem: '风', name: '濠梁乐童', desc: '一位踏歌乐童虚影，破绿衫赤足迎风，玉板轻拍溅起歌谣般的灵光' },
  cao_guojiu:    { cat: 'person', elem: '金', name: '锦绣公子', desc: '一位锦绣公子虚影，锦袍玉带随灵光浮动，玉磬之音如金玉轻鸣' },
  taishang_laojun: { cat: 'person', elem: '火', name: '上古真人', desc: '一缕上古真人灵光，紫气自东而来，拂尘虚影与八卦炉的焰光交织' },
  zhong_kui:     { cat: 'person', elem: '暗', name: '终南寒窗', desc: '一缕书生灵光伏案夜读，寒窗孤灯下字迹成链，隐约凝成赤面虬髯的轮廓' },

  // ===== 宝可梦（属性蛋，系列孵蛋设定） =====
  charmander:    { cat: 'egg', elem: '火', name: '火纹蛋', desc: '一枚火纹蛋，橙色蛋壳燃着一点尾焰，暖光从壳缝漏出，火之生命初醒' },
  bulbasaur:     { cat: 'seed', elem: '木', name: '种子蛋', desc: '一粒种子蛋，青白种壳顶着一枚嫩芽，翠绿叶脉蔓延壳面，像会发芽的小世界' },
  squirtle:      { cat: 'egg', elem: '水', name: '水纹蛋', desc: '一枚水纹蛋，蓝色蛋壳泛着水波鳞光，浪纹在壳面流动，是水之灵胎' },
  eevee:         { cat: 'egg', elem: '光', name: '伊布蛋', desc: '一枚伊布蛋，奶棕蛋壳泛着多变的光晕，进化之力在壳中随时光流转' },
  pikachu:       { cat: 'egg', elem: '雷', name: '电气蛋', desc: '一枚电气蛋，暖黄蛋壳泛着微光，细碎电流纹在壳面爬动，仿佛蛰伏的雷种' },
  riolu:         { cat: 'egg', elem: '钢', name: '波导蛋', desc: '一枚波导蛋，蓝黑蛋壳凝着波导光环，钢铁意志的波纹在壳面一圈圈荡开' },
  ice_fox:       { cat: 'egg', elem: '冰', name: '冰晶蛋', desc: '一枚冰晶蛋，剔透蛋壳凝着霜花，寒蓝微光从内部透出，是冰原的灵胎' },
  rock_rhino:    { cat: 'egg', elem: '土', name: '岩卵', desc: '一枚岩卵，赭石蛋壳覆着龟裂岩纹，沉重如磐石，山岳之力在内中沉睡' },
  wind_falcon:   { cat: 'egg', elem: '风', name: '风卵', desc: '一枚风卵，青白蛋壳绕着流线旋纹，轻如一片浮云，风之羽在壳中蓄势' },
  light_deer:    { cat: 'egg', elem: '光', name: '光卵', desc: '一枚光卵，莹白蛋壳透着晨曦微光，鹿角般的金色芒纹在壳面浮现' },
  dark_panther:  { cat: 'egg', elem: '暗', name: '暗卵', desc: '一枚暗卵，墨黑蛋壳融入夜色，只有紫色幽光在壳缝间若隐若现' },
  steel_armadillo: { cat: 'egg', elem: '钢', name: '钢卵', desc: '一枚钢卵，银灰蛋壳覆着装甲鳞纹，金属冷光中透着坚不可摧的气息' },

  // ===== 数码宝贝（数码蛋） =====
  mecha_dragon:  { cat: 'digi', elem: '火', name: '数码蛋', desc: '一枚悬浮的数码蛋，橙色数据流环绕成环，蛋面裂纹透出勇气之火的光芒' },
  cyber_cat:     { cat: 'digi', elem: '光', name: '数码蛋', desc: '一枚悬浮的数码蛋，圣白数据流结成光环，神圣之息在蛋中凝成爪影' },
  space_mecha:   { cat: 'digi', elem: '光', name: '数码蛋', desc: '一枚悬浮的数码蛋，金色数据流如羽翼舒展，希望之光在蛋中盘旋' },
  quantum_beast: { cat: 'digi', elem: '冰', name: '数码蛋', desc: '一枚悬浮的数码蛋，蓝色数据流凝成霜纹，友情之冰在蛋中静静凝固' },
  digital_phoenix:{ cat: 'digi', elem: '木', name: '数码蛋', desc: '一枚悬浮的数码蛋，翠绿数据流如藤蔓缠绕，歌声的萌动在蛋中孕育' },
  mecha_shark:   { cat: 'digi', elem: '水', name: '数码蛋', desc: '一枚悬浮的数码蛋，海蓝数据流荡出涟漪，诚实的水波在蛋中轻摇' },

  // ===== 国宝（按真实生物：胎生=胚胎 / 卵生=蛋） =====
  panda:         { cat: 'embryo', elem: '土', name: '粉红团子', desc: '一团粉红胚胎，肉嘟嘟蜷成团子，黑白斑纹在茸毛下若隐若现，散发温润的大地生机' },
  golden_monkey: { cat: 'embryo', elem: '金', name: '金丝灵胎', desc: '一团金丝灵胎，淡金绒毛初生，蓝脸的轮廓在暖光中隐约，攀山越涧之性已在血脉' },
  red_crowned_crane: { cat: 'egg', elem: '冰', name: '鹤卵', desc: '一枚鹤卵，青灰蛋壳缀着云纹，一点朱红在壳顶隐现，鹤唳之音在蛋中回荡' },
  south_china_tiger: { cat: 'embryo', elem: '金', name: '虎纹胚胎', desc: '一团虎纹胚胎，橙黑斑纹在微光中初显，额间"王"字的虚影若隐若现' },
  chinese_alligator: { cat: 'egg', elem: '水', name: '鳄卵', desc: '一枚鳄卵，米白蛋壳覆着细密鳞纹，半截小鳄的剪影在壳中蜷伏' },
  crested_ibis:  { cat: 'egg', elem: '火', name: '朱卵', desc: '一枚朱卵，米白蛋壳晕着一抹绯红，东方宝石的祥光在壳面流转' },
  tibetan_antelope: { cat: 'embryo', elem: '冰', name: '雪原灵胎', desc: '一团雪原灵胎，白褐绒毛初生，剑形长角的一点雏形在寒光中显现' },
  snow_leopard:  { cat: 'embryo', elem: '冰', name: '斑雪灵胎', desc: '一团斑雪灵胎，灰白绒毛缀着黑斑，雪山之王的冷冽气息在胚胎中沉淀' },
  milu_deer:     { cat: 'embryo', elem: '水', name: '泽畔灵胎', desc: '一团泽畔灵胎，浅褐绒毛沾着水汽，四不像的角影在湿地灵光中若隐若现' },
  siberian_tiger:{ cat: 'embryo', elem: '冰', name: '雪虎灵胎', desc: '一团雪虎灵胎，厚绒毛泛着冰晶微光，林海雪原的威压已在胚胎中酝酿' },
  red_panda:     { cat: 'embryo', elem: '木', name: '红绒灵胎', desc: '一团红绒灵胎，火红绒毛裹着竹香，环纹尾巴的雏形在暖光中蜷起' },
  finless_porpoise: { cat: 'embryo', elem: '水', name: '江流灵胎', desc: '一团江流灵胎，圆头灰蓝泛着水光，嘴角的微笑弧度在江涛中若隐若现' },

  // ===== 魔法奇幻（灵种 / 蛋 / 器物） =====
  unicorn:       { cat: 'mythic', elem: '光', name: '星光种子', desc: '一粒星光种子，银白灵光凝成独角雏形，圣洁芒星在种壳上闪烁' },
  wyvern:        { cat: 'egg', elem: '火', name: '龙蛋', desc: '一枚龙蛋，赤色蛋壳覆着岩浆裂纹，龙翼的虚影在火光中鼓动' },
  fairy:         { cat: 'seed', elem: '木', name: '花苞', desc: '一枚花苞灵种，金色花蕊在晨露中初绽，精灵翅影藏在花瓣之间' },
  treant:        { cat: 'seed', elem: '木', name: '种子', desc: '一粒橡树种子，深褐种壳裂出嫩绿新芽，森林的脉搏在根须间跳动' },
  griffin:       { cat: 'egg', elem: '金', name: '狮鹫蛋', desc: '一枚狮鹫蛋，金白蛋壳生着细羽纹，鹰喙与狮爪的雏形在蛋中蓄势' },
  mermaid:       { cat: 'egg', elem: '水', name: '珍珠', desc: '一枚深海珍珠，碧光流转如潮水，人鱼歌谣的涟漪在珠身内漾开' },
  grey_wizard:   { cat: 'spirit', elem: '秘', name: '魔杖种子', desc: '一粒魔杖种子，星芒符文绕种盘旋，灰袍贤者的灵光在奥秘中浮现' },
  wand_cat:      { cat: 'spirit', elem: '秘', name: '魔法绒球', desc: '一团魔法绒球，紫毛间缀着星点，魔杖的微光在绒球中央轻轻闪动' },
  dragon_knight: { cat: 'egg', elem: '火', name: '契约之卵', desc: '一枚契约之卵，暗金蛋壳刻着誓约纹，龙影与骑士的羁绊在卵中相连' },
  alchemy_golem: { cat: 'artifact', elem: '金', name: '贤者之石', desc: '一块贤者之石，鎏金棱面内蕴符文，炼金之火在其中凝成一粒跳动的心' },
  nightmare_horse:{ cat: 'egg', elem: '暗', name: '月影卵', desc: '一枚月影卵，紫黑蛋壳晕着月辉，梦魇蹄影在夜色灵光中踏空' },
  lamp_spirit:   { cat: 'artifact', elem: '火', name: '灯芯种子', desc: '一粒灯芯种子，灯油灵光托着一朵火苗，千愿的浮影在火中明灭' },

  // ===== 史前生物（恐龙爬行=蛋 / 哺乳=胚胎） =====
  t_rex:         { cat: 'egg', elem: '火', name: '化石蛋', desc: '一枚化石蛋，斑驳壳面泛着远古纹路，内部透出熔岩般的赤光，像沉睡的史前之心' },
  triceratops:   { cat: 'egg', elem: '土', name: '角蛋', desc: '一枚角蛋，沙棕蛋壳顶生三处角突，颈盾的雏形在蛋壳边缘隆起' },
  pterosaur:     { cat: 'egg', elem: '风', name: '飞卵', desc: '一枚飞卵，灰白蛋壳生着翼膜纹，气流在蛋身周围盘旋，似欲乘风而起' },
  mammoth:       { cat: 'embryo', elem: '冰', name: '长毛灵胎', desc: '一团长毛灵胎，厚实绒毛覆着霜晶，弯长象牙的一点雏形在寒光中显现' },
  sabertooth:    { cat: 'embryo', elem: '金', name: '剑齿灵胎', desc: '一团剑齿灵胎，肌肉初成、绒毛未褪，上颚獠牙的锋利雏形若隐若现' },
  mosasaur:      { cat: 'egg', elem: '水', name: '水卵', desc: '一枚水卵，深蓝蛋壳泛着鳞光，鳍肢的剪影在蛋中游弋，苍茫海意在壳内涌动' },
  spinosaurus:   { cat: 'egg', elem: '水', name: '帆纹蛋', desc: '一枚帆纹蛋，河畔泥色蛋壳隆起一道帆脊纹，河水的湿润气息包裹着它' },
  ankylosaurus:  { cat: 'egg', elem: '土', name: '锤尾蛋', desc: '一枚锤尾蛋，岩甲色蛋壳覆着板甲纹，尾端锤形的凸起在壳下隐约' },
  diplodocus:    { cat: 'egg', elem: '木', name: '长颈蛋', desc: '一枚长颈蛋，青灰蛋壳有着细长的颈纹，温柔巨兽的轮廓在壳中延伸' },
  megalodon:     { cat: 'egg', elem: '水', name: '巨牙卵', desc: '一枚巨牙卵，深灰卵壳生着锯齿纹，巨齿的锋芒在深蓝灵光中初现' },
  ground_sloth:  { cat: 'embryo', elem: '土', name: '巨爪灵胎', desc: '一团巨爪灵胎，长毛初生的笨拙轮廓，巨型爪钩的雏形蜷在胸前' },
  woolly_rhino:  { cat: 'embryo', elem: '冰', name: '冰角灵胎', desc: '一团冰角灵胎，厚毛覆着苔原霜气，双角雏形在冰原灵光中缓缓凝成' },

  // ===== 星座守护（星屑 / 星辉灵光，属性按黄道三性） =====
  aries:         { cat: 'spirit', elem: '火', name: '星屑之种', desc: '一团星屑灵光，白羊命星在其中闪烁，暖金色小宇宙初醒，羊角般的星光微凝' },
  taurus:        { cat: 'spirit', elem: '土', name: '星辉之种', desc: '一团星辉灵光，金牛星芒沉稳厚重，大地般的小宇宙在星辉中缓缓脉动' },
  gemini:        { cat: 'spirit', elem: '风', name: '星辉双子', desc: '两团交缠的星辉灵光，一明一暗宛如双子，善恶之争在风中无声展开' },
  cancer:        { cat: 'spirit', elem: '水', name: '冥辉之种', desc: '一团冥辉灵光，蟹钳般的紫芒时隐时现，冥界幽光在蛋形光晕中流转' },
  leo:           { cat: 'spirit', elem: '火', name: '星辉之种', desc: '一团星辉灵光，狮鬃般的金色芒焰喷张，战意如火的小宇宙炽热燃烧' },
  virgo:         { cat: 'spirit', elem: '土', name: '星辉莲花', desc: '一朵星辉莲花，白金花瓣层层合拢，闭目禅相在莲心流转，最接近神的光' },
  libra:         { cat: 'spirit', elem: '风', name: '星辉龙珠', desc: '一粒星辉龙珠，天秤虚影与十二道金光浮沉，守望者的龙息在珠中沉眠' },
  scorpio:       { cat: 'spirit', elem: '水', name: '星辉之种', desc: '一团星辉灵光，蝎尾般的猩红针芒微露，深红毒针的锐意在星辉中蛰伏' },
  sagittarius:   { cat: 'spirit', elem: '火', name: '星辉之箭', desc: '一支星辉之箭，黄金弓芒凝成箭形，射穿苍穹的意志在星光中引而不发' },
  capricorn:     { cat: 'spirit', elem: '土', name: '星辉之剑', desc: '一柄星辉之剑，银白剑芒凝成圣剑雏形，Excalibur的锋锐在土色星光中成形' },
  aquarius:      { cat: 'spirit', elem: '冰', name: '星辉冰晶', desc: '一粒星辉冰晶，冰蓝灵光凝霜成晶，绝对零度的气息在寒光中沉淀' },
  pisces:        { cat: 'spirit', elem: '水', name: '星辉花苞', desc: '一朵星辉花苞，绯红花瓣裹着双鱼游影，玫瑰之毒的暗香在灵光中浮动' },

  // ===== 传统节日（食物原料 / 器物） =====
  zongzi:        { cat: 'food', elem: '木', name: '箬叶籽', desc: '一粒箬叶包裹的糯米灵种，青绿叶片系着绳结，米香与端午的艾草气萦绕' },
  tangyuan:      { cat: 'food', elem: '水', name: '糯米团', desc: '一团雪白糯米团，软糯圆润浮在甜汤灵光里，流心芝麻的甜意在团中酝酿' },
  mooncake:      { cat: 'food', elem: '金', name: '面团胚', desc: '一块面团胚，金黄饼形压着花印纹，满月般的圆满寓意在胚中成形' },
  qingtuan:      { cat: 'food', elem: '木', name: '艾草芽', desc: '一枚艾草芽，碧绿糯米团裹着青碧草色，清明踏青的清香在芽尖浮动' },
  chongyang_cake:{ cat: 'food', elem: '金', name: '米粉种', desc: '一粒米粉种，金黄米塔初具层级，重阳登高的茱萸香气在粉粒间萦绕' },
  niangao:       { cat: 'food', elem: '金', name: '糯米砖', desc: '一块糯米砖，莹白软糯的方砖雏形，年年高升的红纸光晕在砖面浮动' },
  laba_porridge: { cat: 'food', elem: '土', name: '谷粒籽', desc: '一把谷粒籽，红豆桂圆莲子红枣五谷杂陈，腊八的暖意在百味中酝酿' },
  spring_pancake:{ cat: 'food', elem: '木', name: '面糊团', desc: '一团面糊，薄嫩面胚泛着麦香，春蔬的嫩绿在饼皮中透出咬春之意' },
  tanghulu:      { cat: 'food', elem: '火', name: '山楂果', desc: '一粒山楂果，红亮圆润裹着初凝的冰糖衣，街头的酸甜在脆壳下待绽' },
  osmanthus_cake:{ cat: 'food', elem: '金', name: '桂子籽', desc: '一粒桂子籽，米白糕胚缀着金桂花瓣，金秋的桂香在甜意中酝酿' },
  wonton:        { cat: 'food', elem: '水', name: '面皮种', desc: '一张面皮种，薄如云翳的面胚裹着鲜肉初形，热汤的暖意在褶皱中藏起' },
  festival_lantern: { cat: 'artifact', elem: '火', name: '竹骨灯', desc: '一盏竹骨灯，素绢灯面初糊成胚，烛火微光在骨架上跳动，灯谜的期待在绢面浮现' },

  // ===== 虹猫蓝兔七侠传（剑意灵光） =====
  hongmao:       { cat: 'sword', elem: '火', name: '剑意种子', desc: '一缕赤虹剑意凝成的种子，剑光如丝缠绕，隐约映出红衣小猫的虚影' },
  lantu:         { cat: 'sword', elem: '冰', name: '冰魄剑意', desc: '一缕冰魄剑意凝成的种子，寒光凝霜成种，蓝衣玉兔的轮廓在冰魄中浮现' },
  doudou:        { cat: 'sword', elem: '水', name: '雨花剑意', desc: '一缕雨花剑意凝成的种子，剑光如雨丝纷扬，白绒小狗与药囊的虚影若隐若现' },
  dabeng:        { cat: 'sword', elem: '雷', name: '奔雷剑意', desc: '一缕奔雷剑意凝成的种子，雷光在剑种中轰鸣，棕熊神力的轮廓在雷芒中显现' },
  tiaotiao:      { cat: 'sword', elem: '风', name: '青光剑意', desc: '一缕青光剑意凝成的种子，剑光快如飞燕掠影，绿衣小猴的虚影在风中跳跃' },
  shali:         { cat: 'sword', elem: '木', name: '紫云剑意', desc: '一缕紫云剑意凝成的种子，紫云翻涌如雾，紫衣剑客的温柔轮廓在其中隐现' },
  dada:          { cat: 'sword', elem: '风', name: '旋风剑意', desc: '一缕旋风剑意凝成的种子，剑气旋成环舞，竹林熊猫的沉稳虚影在风中端坐' },
  heixinhu:      { cat: 'sword', elem: '暗', name: '暗影虎意', desc: '一缕暗影虎意凝成的种子，黑风裹着威压，猛虎如炬的双目在黑暗中睁开' },
  heixiaohu:     { cat: 'sword', elem: '暗', name: '暗影虎卵', desc: '一枚暗影虎卵，墨黑灵光裹着虎影，利爪的锋芒在执念般的暗光中初现' },
}

export function getPetInitial(speciesId: string): PetInitial {
  return PET_INITIAL[speciesId] || {
    cat: 'mythic',
    elem: '秘',
    name: '未知灵种',
    desc: '一粒尚未明辨属性的神秘灵种，在未知的光芒中安静蛰伏。',
  }
}
