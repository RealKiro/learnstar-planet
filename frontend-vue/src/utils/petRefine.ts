// ===== 学趣星球 · 宠物六阶段精修层（神态 / 动作 / 衣着 / 梳造） =====
// 在 v2 提示词（形态+意象+场景+风格+色彩）基础上，为每一阶段补充四个角色维度：
//   神态 spirit —— 表情与眼神（随阶段心境变化）
//   动作 action —— 姿态与标志动作（{m} 替换为该物种 movement 演化）
//   衣着 attire —— 服饰 / 圣衣 / 毛皮 递进
//   梳造 hair   —— 发型 / 头饰 / 角羽 递进
// 数据源：六阶段基座（按大类）+ 重点角色手写覆盖（REFINE_SPECIAL）。
// 生成脚本：scripts/generate-pet-prompts-v2.js。

export type RefineCategory =
  | 'animal'    // 兽类（宝可梦/国宝/史前哺乳/白虎等）：毛皮鳞甲演化
  | 'dragon'    // 神龙蛇类（烛龙/青龙/玄武/扬子鳄/沧龙/飞龙）：鳞角升腾
  | 'mythic'    // 神兽瑞兽（凤凰/麒麟/貔貅/独角兽/狮鹫等）：灵光化形
  | 'saint'     // 圣衣战士（星座守护 12）：星屑→青铜→黄金
  | 'sword'     // 剑客（虹猫蓝兔 9）：剑意素衣→侠装
  | 'immortal'  // 仙侠人物（东方神话 20）：布衣→道袍→仙袍
  | 'digi'      // 数码兽（数码宝贝 6）：数码蛋→幼年→究极
  | 'spirit'    // 灵体法师（巫师/魔杖猫/龙骑士/精灵）：法袍星芒
  | 'plant'     // 植物（妙蛙种子/树人/花精灵）：芽→花→树
  | 'food'      // 食物拟人（节日 11）：原料→成品
  | 'artifact'  // 器物（花灯/贤者石/灯芯/魔像）：素坯→灵光
  | 'prehistoric' // 远古爬行（恐龙/巨齿鲨/沧龙等卵生）：幼体→霸主

export interface StageRefine {
  spirit: string
  action: string
  attire: string
  hair: string
}
export type RefineBase = Record<RefineCategory, StageRefine[]>
export type RefineSpecial = Record<string, { attire: string; hair: string }[]>

/** 每物种所属大类（决定用哪套六阶段基座） */
export const SPECIES_CATEGORY: Record<string, RefineCategory> = {
  // 山海经
  zhulong: 'dragon', yinglong: 'dragon', nine_tail_fox: 'mythic', kunpeng: 'mythic',
  fenghuang: 'mythic', qilin: 'mythic', qiongqi: 'mythic', bifang: 'mythic',
  pixiu: 'mythic', jingwei: 'mythic', xiangliu: 'dragon', xiezhi: 'mythic',
  qinglong: 'dragon', baihu: 'animal', zhuque: 'mythic', xuanwu: 'dragon',
  taotie: 'mythic', baize: 'mythic',
  // 东方神话
  jiang_ziya: 'immortal', nezha: 'immortal', yang_jian: 'immortal', lei_zhenzi: 'immortal',
  huang_tianhua: 'immortal', tu_xingsun: 'immortal', yang_ren: 'immortal', wei_hu: 'immortal',
  daji: 'immortal', shen_gongbao: 'immortal', sun_wukong: 'immortal', lv_dongbin: 'immortal',
  he_xiangu: 'immortal', zhang_guolao: 'immortal', tie_guaili: 'immortal', han_zhongli: 'immortal',
  lan_caihe: 'immortal', cao_guojiu: 'immortal', taishang_laojun: 'immortal', zhong_kui: 'immortal',
  // 宝可梦
  charmander: 'animal', bulbasaur: 'plant', squirtle: 'animal', eevee: 'animal',
  pikachu: 'animal', riolu: 'animal', ice_fox: 'animal', rock_rhino: 'animal',
  wind_falcon: 'animal', light_deer: 'animal', dark_panther: 'animal', steel_armadillo: 'animal',
  // 数码宝贝
  mecha_dragon: 'digi', cyber_cat: 'digi', space_mecha: 'digi', quantum_beast: 'digi',
  digital_phoenix: 'digi', mecha_shark: 'digi',
  // 国宝
  panda: 'animal', golden_monkey: 'animal', red_crowned_crane: 'animal', south_china_tiger: 'animal',
  chinese_alligator: 'dragon', crested_ibis: 'animal', tibetan_antelope: 'animal',
  snow_leopard: 'animal', milu_deer: 'animal', siberian_tiger: 'animal',
  red_panda: 'animal', finless_porpoise: 'animal',
  // 魔法奇幻
  unicorn: 'mythic', wyvern: 'dragon', fairy: 'plant', treant: 'plant', griffin: 'mythic',
  mermaid: 'mythic', grey_wizard: 'spirit', wand_cat: 'spirit', dragon_knight: 'spirit',
  alchemy_golem: 'artifact', nightmare_horse: 'mythic', lamp_spirit: 'artifact',
  // 史前生物
  t_rex: 'prehistoric', triceratops: 'prehistoric', pterosaur: 'prehistoric', mammoth: 'animal',
  sabertooth: 'animal', mosasaur: 'dragon', spinosaurus: 'prehistoric', ankylosaurus: 'prehistoric',
  diplodocus: 'prehistoric', megalodon: 'prehistoric', ground_sloth: 'animal', woolly_rhino: 'animal',
  // 星座守护
  aries: 'saint', taurus: 'saint', gemini: 'saint', cancer: 'saint', leo: 'saint', virgo: 'saint',
  libra: 'saint', scorpio: 'saint', sagittarius: 'saint', capricorn: 'saint', aquarius: 'saint', pisces: 'saint',
  // 传统节日
  zongzi: 'food', tangyuan: 'food', mooncake: 'food', qingtuan: 'food', chongyang_cake: 'food',
  niangao: 'food', laba_porridge: 'food', spring_pancake: 'food', tanghulu: 'food',
  osmanthus_cake: 'food', wonton: 'food', festival_lantern: 'artifact',
  // 虹猫蓝兔
  hongmao: 'sword', lantu: 'sword', doudou: 'sword', dabeng: 'sword', tiaotiao: 'sword',
  shali: 'sword', dada: 'sword', heixinhu: 'sword', heixiaohu: 'sword',
}

// ============================================================
// 六阶段基座模板（{m} = 该物种标志动作 movement，脚本替换）
// 阶段序：灵胎初醒 / 凡尘砺心 / 道法初成 / 大劫淬炼 / 封神登天 / 归真永恒
// ============================================================
export const REFINE_BASE: RefineBase = {
  // —— 兽类 ——
  animal: [
    { spirit: '眼帘低垂，呼吸轻浅，沉在初生的微光里', action: '蜷卧于光中，静默沉睡', attire: '细嫩绒毛/初生软鳞，未成形', hair: '新生短毛，耳未立/角未出' },
    { spirit: '眸光清澈，好奇打量，怯生生又藏不住欢喜', action: '蹒跚学步，试探着触碰四周', attire: '绒毛渐丰，隐约透出本色', hair: '耳尾渐立，茸毛蓬松' },
    { spirit: '眸光晶亮，意气初显，跃跃欲试', action: '初展身手，初露锋芒', attire: '毛色分明，姿态挺拔', hair: '鬃毛/羽角初生，精神抖擞' },
    { spirit: '眸光深沉，眉峰微蹙，带着淬炼的坚毅', action: '豁尽全力，势如破竹', attire: '成体毛皮/鳞甲，或有伤痕', hair: '角/鬃/尾冠成型，凌乱而烈' },
    { spirit: '目光如炬，不怒自威，威仪自生', action: '绝技大成，昂首傲立山巅', attire: '金色祥纹缀身，王者的气象', hair: '金鬃/长羽/圣角，威风凛凛' },
    { spirit: '王者归位，万灵俯首', action: '啸震九霄，天地同应', attire: '神光化铠，日月护身', hair: '王者之冕，诸兽共尊' },
  ],
  // —— 神龙蛇类 ——
  dragon: [
    { spirit: '龙息轻吐，沉睡于混沌灵光', action: '盘蜷于光，沉眠未醒', attire: '幼嫩鳞胚，泛着初光', hair: '无角，须影初现' },
    { spirit: '竖瞳初睁，窥探云水', action: '初探云水，游弋学步', attire: '鳞甲渐密，颜色初显', hair: '角芽微露，须丝轻扬' },
    { spirit: '眸光锐亮，潜龙欲腾', action: '腾空而起，初显神威', attire: '鳞甲生辉，腹光流转', hair: '双角初成，须如流云' },
    { spirit: '龙威炽烈，怒目电光', action: '全力施为，风雷随身', attire: '战损鳞甲，雷火纹显', hair: '角芒凌厉，须张如戟' },
    { spirit: '龙目洞彻，神威赫赫', action: '绝技大成，行云布雨', attire: '金鳞覆身，祥光万道', hair: '龙角如珊瑚，须垂百丈' },
    { spirit: '万龙之首，睥睨三界', action: '一怒天地色变，号令风雨', attire: '龙身化岳，鳞甲映日月', hair: '龙角擎天，须断星河' },
  ],
  // —— 神兽瑞兽 ——
  mythic: [
    { spirit: '灵光中沉睡，兽性未醒的宁静', action: '蜷于灵光，微息起伏', attire: '灵光虚影，形尚未凝', hair: '灵毫光点，未成形相' },
    { spirit: '初踏山川，懵懂而灵慧', action: '蹒跚踏云，好奇嗅闻', attire: '半实灵体，羽毛/鳞纹初显', hair: '幼羽/灵尾，泛着微光' },
    { spirit: '眼神古老而专注，神通初显', action: '展翅昂首，神姿初现', attire: '华羽灵纹，瑞气氤氲', hair: '长羽/灵角渐生' },
    { spirit: '眸光如电，威严中带着坚韧', action: '全力施为，祥光与煞气并舞', attire: '神纹满身，羽鳞如铠', hair: '羽冠/灵角，光华流转' },
    { spirit: '神兽威严，目光洞彻九幽', action: '绝技大成，百兽来朝', attire: '神光加身，五色祥云', hair: '圣羽垂天，瑞角冲霄' },
    { spirit: '瑞气冲霄，万灵俯首', action: '真身镇世，祥瑞永驻', attire: '祥光化雨，福泽苍生', hair: '瑞光冠冕，天地共仰' },
  ],
  // —— 圣衣战士 ——
  saint: [
    { spirit: '命星微光中沉睡，小宇宙初醒而不自知', action: '星屑环绕，静立若定', attire: '星屑虚影，圣衣未凝', hair: '命星点光，束发未冠' },
    { spirit: '修行中的专注，汗珠映着星辉', action: '挥拳踢腿，苦练招式', attire: '素白练功服，圣衣雏形初覆', hair: '束发紧绷，汗湿鬓角' },
    { spirit: '为守护而战的炽热，第七感觉醒', action: '绝技初现，拳风呼啸', attire: '青铜色圣衣，护肩护腕初成', hair: '青铜头冠微启' },
    { spirit: '绝境中燃烧小宇宙，眸光不屈', action: '全力施展，血战不退', attire: '白银圣衣，甲胄染尘', hair: '头冠残破，银发凌乱' },
    { spirit: '金眸光华内敛，不怒自威', action: '绝技大成，镇守宫门', attire: '黄金圣衣全装，辉光闪耀', hair: '星座金盔，庄严生辉' },
    { spirit: '第八感圆满，战意与神性合一', action: '绝技巅峰，星碎神裂', attire: '神圣衣全开，星辉圣甲', hair: '神冠星盔，威临宇宙' },
  ],
  // —— 剑客 ——
  sword: [
    { spirit: '剑意初凝，稚气未脱的认真', action: '剑意化种，静候萌发', attire: '剑意虚影，素衣未备', hair: '发丝初束，剑穗微晃' },
    { spirit: '初握剑的紧张与欢喜', action: '笨拙挥剑，破绽百出却认真', attire: '素色短打，布带缠手', hair: '束发利落，额前碎发' },
    { spirit: '剑气初成，眸光锐利', action: '剑意初现，剑光如虹', attire: '练功劲装，护腕已备', hair: '发带束起，英气渐显' },
    { spirit: '浴血后的沉静，杀气内敛', action: '全力挥剑，剑影纵横', attire: '侠客劲装，血迹未干', hair: '发乱而不颓，眼神如刃' },
    { spirit: '目光如剑，豪气干云', action: '剑法大成，一剑定乾坤', attire: '锦袍披风，名剑在腰', hair: '侠客冠束，风采卓然' },
    { spirit: '剑道圆满，锋芒内敛却慑人', action: '人剑合一，剑光裂天', attire: '剑者长袍，剑气绕身', hair: '高冠束发，剑眉入鬓' },
  ],
  // —— 仙侠人物 ——
  immortal: [
    { spirit: '灵光中沉睡，仙缘初定的宁静', action: '灵光化形，气息未定', attire: '灵光虚影，衣形未凝', hair: '灵光为发，未成形相' },
    { spirit: '入世初见的纯澈', action: '初踏红尘，好奇四顾', attire: '粗布素衣，简朴无饰', hair: '少年总角/披发，质朴' },
    { spirit: '悟道中的专注，眸光渐亮', action: '功法初显，招式渐成', attire: '道袍初成，法器中藏', hair: '束发簪冠，渐有仙姿' },
    { spirit: '渡劫时的坚毅，眼神无畏', action: '全力抗劫，天劫加身', attire: '法衣破损，法力激荡', hair: '发冠凌乱，眸若星辰' },
    { spirit: '仙光内蕴，目光深邃', action: '功法大成，威临天下', attire: '仙袍法冠，法宝随身', hair: '仙髻金冠，风姿卓然' },
    { spirit: '道果大成，神光自照的庄严', action: '法相全开，气吞八荒', attire: '功德金光，道果法衣', hair: '万道归冠，无上威严' },
  ],
  // —— 数码兽 ——
  digi: [
    { spirit: '数据流中沉睡的懵懂', action: '悬浮微光，数据环绕', attire: '数码蛋，数据纹流转', hair: '数据环光，无固定形' },
    { spirit: '探索数码世界的雀跃', action: '蹦跳探索，好奇触碰', attire: '幼年兽体，毛茸茸', hair: '幼毛/呆毛，圆润可爱' },
    { spirit: '进化前的坚定，眸光发亮', action: '绝技初现，数据粒子升腾', attire: '成熟体装甲，渐渐成形', hair: '头甲初现，眼神明亮' },
    { spirit: '直面黑暗的勇毅', action: '全力施为，与深渊对峙', attire: '完全体铠甲，伤痕累累', hair: '战甲破损，眸光如焰' },
    { spirit: '究极体的神圣威仪', action: '绝技大成，数据凝成圣甲', attire: '究极体圣甲，辉光万丈', hair: '圣盔金角，威仪堂堂' },
    { spirit: '数据升维，神明降临', action: '绝技轰天，数码洪流归一', attire: '神体数据化，圣光铸形', hair: '数据光环，如日当空' },
  ],
  // —— 灵体法师 ——
  spirit: [
    { spirit: '星芒中蜷缩，似梦似醒', action: '悬于法阵微光，静候唤醒', attire: '法袍虚影，星光为衣', hair: '灵光发丝，半透明' },
    { spirit: '好奇试探的灵动', action: '吟唱初学，魔力火花迸溅', attire: '学徒法袍，布料朴素', hair: '发丝扬起，缀着星尘' },
    { spirit: '光华渐盛，神采焕发', action: '魔力初现，法阵脚下亮起', attire: '法师长袍，绣星纹', hair: '魔法帽/发冠初戴' },
    { spirit: '暗夜中坚持，眼神倔强', action: '全力施法，魔力风暴', attire: '法袍破损，魔力燃烧', hair: '发丝凌乱，眸中带焰' },
    { spirit: '威能内蕴，光晕庄严', action: '魔法大成，法阵漫天', attire: '贤者白袍/秘袍，权杖在手', hair: '王冠/法冠，威仪自生' },
    { spirit: '万法归心，与元素同源', action: '法阵通天，咒文如潮', attire: '元素法身，光尘化袍', hair: '星冕加身，魔力如渊' },
  ],
  // —— 植物 ——
  plant: [
    { spirit: '沉睡的种子里，藏着破土的心跳', action: '静卧泥土，待春雨而萌', attire: '种子/种壳，纹理细腻', hair: '顶端嫩芽微顶' },
    { spirit: '初生的好奇，向着光的方向', action: '嫩芽破土，努力伸展', attire: '幼芽新叶，娇嫩欲滴', hair: '顶芽鲜嫩，两片子叶' },
    { spirit: '生机勃发，神采奕奕', action: '生机初现，枝叶舒展', attire: '茎叶渐盛，花苞初成', hair: '花苞/嫩叶环生' },
    { spirit: '风雨中的坚韧，眸光不折', action: '全力生长，根系深扎', attire: '枝叶繁茂，带风霜痕迹', hair: '花叶翻卷，仍自挺立' },
    { spirit: '华彩绽放，生机盎然的威仪', action: '生机大成，繁花满枝', attire: '花开满枝/树冠如云', hair: '花冠/树冠，光华流转' },
    { spirit: '万灵朝拜，生命之尊', action: '花开万里，果实垂天', attire: '参天古木，藤蔓缠霄', hair: '万叶化冠，春永驻世' },
  ],
  // —— 食物 ——
  food: [
    { spirit: '原料的宁静，蕴着期待', action: '静置案板，等待巧手', attire: '米面/食材原料，朴实', hair: '无，原料形态' },
    { spirit: '制作中的专注', action: '揉捏塑形，蒸汽初升', attire: '半成品，形态渐显', hair: '裹叶/印模初成' },
    { spirit: '成形时的欢喜', action: '摆上蒸笼，静待火候', attire: '形态完整，色泽初显', hair: '印花/装饰点缀' },
    { spirit: '火候淬炼的专注', action: '经受蒸煮/煎炸，坚韧定型', attire: '色泽加深，香气酝酿', hair: '糖衣/油光初亮' },
    { spirit: '出锅的骄傲，暖意融融', action: '热气腾腾，登堂亮相', attire: '成品佳肴，色香味全', hair: '装饰华美，名品之姿' },
    { spirit: '名品圆满，色香俱全的骄傲', action: '登堂入室，香飘满座', attire: '华美名品，灵光透色', hair: '宝光流转，镇席之品' },
  ],
  // —— 器物 ——
  artifact: [
    { spirit: '沉睡的灵光，物灵未醒', action: '静置无声，灵光内蕴', attire: '素坯/原石，未成形', hair: '无，器物胚形' },
    { spirit: '初生灵智的好奇', action: '微光颤动，器灵初醒', attire: '初雕成型，轮廓渐显', hair: '雕纹/铭文初刻' },
    { spirit: '灵光渐盛，灵动自生', action: '器光初现，灵光流转', attire: '成形精工，纹饰渐繁', hair: '纹饰/嵌饰增辉' },
    { spirit: '淬炼中的忍耐', action: '全力受炼，火炼淬洗', attire: '历经淬炼，温润内敛', hair: '包浆/裂纹，岁月痕' },
    { spirit: '灵光大成的庄严', action: '器灵大成，光华夺目', attire: '华彩流光，灵气充盈', hair: '祥纹满饰，宝光外放' },
    { spirit: '器灵化神，灵性通明', action: '器道通天，光华耀世', attire: '神纹流转，镇世之宝', hair: '灵光化形，器魂永驻' },
  ],
  // —— 远古爬行 ——
  prehistoric: [
    { spirit: '化石般的沉睡，原始生命力潜伏', action: '静卧蛋中，尾/爪微动', attire: '蛋壳斑驳，远古纹路', hair: '无，蛋中初形' },
    { spirit: '初生的笨拙好奇', action: '破壳蹒跚，懵懂张望', attire: '稚嫩皮甲，柔软未坚', hair: '头冠/角芽未显' },
    { spirit: '捕猎时的专注凶狠', action: '低伏潜行，初试锋芒', attire: '鳞甲渐厚，色泽加深', hair: '头冠/帆脊渐起' },
    { spirit: '生存竞争中的冷酷坚毅', action: '全力搏杀，与天地争食', attire: '成体鳞甲，伤痕累累', hair: '角/帆/鬃威猛' },
    { spirit: '霸主之威，眼神睥睨', action: '猎技大成，万兽辟易', attire: '霸主之躯，王纹隐现', hair: '顶冠/长角，王者相' },
    { spirit: '洪荒之巅，唯我独尊', action: '踏碎山河，万兽臣服', attire: '远古神躯，鳞甲映日', hair: '骨冠擎天，威压万古' },
  ],
}

// ============================================================
// 重点角色手写精修（衣着 / 梳造 六阶段；覆盖基座）
// 按 speciesId 索引，数组长度 6 = 六阶段
// ============================================================
export const REFINE_SPECIAL: RefineSpecial = {
  // ---------- 东方神话 20 ----------
  nezha: [
    { attire: '灵珠之光中的孩童轮廓，红肚兜若隐若现', hair: '总角发髻，呆毛微翘' },
    { attire: '红肚兜赤足孩童，混天绫缠在腕间', hair: '双总角发髻，额前齐刘海' },
    { attire: '莲花战甲初成，乾坤圈套在颈间', hair: '发髻渐紧，束发带飞扬' },
    { attire: '三头六臂初现，火尖枪在手，战甲燃焰', hair: '六臂之姿，发丝燃着火光' },
    { attire: '三头六臂全开，莲花金甲，风火轮踏焰', hair: '金冠束发，英姿勃发' },
    { attire: '三头六臂莲花战神，金甲全开', hair: '金冠束发，六臂威仪' },
  ],
  sun_wukong: [
    { attire: '灵石中的石猴轮廓，身裹混沌仙气', hair: '石纹未褪，通体金毛初生' },
    { attire: '花果山的小石猴，藤叶为衣', hair: '金色猴毛，红冠未成' },
    { attire: '虎皮裙初穿，金箍棒在手', hair: '金毛如焰，凤翅紫金冠初戴' },
    { attire: '大闹天宫的战甲，锁子黄金甲', hair: '紫金冠斜戴，战意凛然' },
    { attire: '齐天大圣冠冕，金甲红袍', hair: '凤翅紫金冠，金毛猎猎' },
    { attire: '斗战胜佛袈裟，佛光内蕴', hair: '金毛归于平静，项间念珠' },
  ],
  yang_jian: [
    { attire: '天眼灵光中的银甲虚影', hair: '束发未冠，额间天目微闭' },
    { attire: '灌江口少年，布衣素服', hair: '束发利落，英气初显' },
    { attire: '银甲初着，三尖两刃刀在手', hair: '银冠束发，天目渐睁' },
    { attire: '八九玄功激荡，战袍猎猎', hair: '发丝凌乱，天目怒张' },
    { attire: '二郎真君银甲金冠，哮天犬随侧', hair: '银盔高束，天目如电' },
    { attire: '清源妙道真君，玄功大成', hair: '束发道冠，天目内敛' },
  ],
  lei_zhenzi: [
    { attire: '雷卵灵光中的雏鸟轮廓', hair: '无，蛋中初形' },
    { attire: '幼年孩童，青布短衣', hair: '总角发髻，天真' },
    { attire: '食杏果后鸟面初显，身生细羽', hair: '发间生出细羽，青面初现' },
    { attire: '鸟面雷公之相，背生双翅，黄金棍在手', hair: '风雷双翅，发如闪电' },
    { attire: '雷公真身，风雷战甲，金棍引雷', hair: '雷霆之发，翼展千里' },
    { attire: '风雷显圣，忠义神将之姿', hair: '雷冠神相，威而不怒' },
  ],
  huang_tianhua: [
    { attire: '将门灵种，虎子之气的灵光', hair: '总角发髻，稚气' },
    { attire: '少年将军，素色劲装', hair: '束发带冠，少年英气' },
    { attire: '青布战袍，莫邪剑初佩', hair: '束发簪冠，眉目英武' },
    { attire: '玉麒麟踏阵，莫邪剑出鞘', hair: '战盔斜戴，血战之姿' },
    { attire: '少年将军金甲，玉麒麟为骑', hair: '金冠束发，英姿飒爽' },
    { attire: '封神之姿，少年英雄的永恒', hair: '神光化冠，风华长存' },
  ],
  tu_xingsun: [
    { attire: '地灵种子，土色灵光', hair: '矮小身形，发丝粗短' },
    { attire: '矮小少年，粗布短衣', hair: '平头短发，机灵' },
    { attire: '镔铁棍在手，土遁初成', hair: '发间沾土，狡黠' },
    { attire: '土遁如飞，捆仙绳藏袖', hair: '灰头土脸，笑意狡黠' },
    { attire: '地行仙之姿，镔铁棍如臂使指', hair: '短须微蓄，仙气内藏' },
    { attire: '地行仙圆满，镔铁棍破土裂山', hair: '短须如戟，仙威内藏' },
  ],
  yang_ren: [
    { attire: '忠臣灵光，朝服虚影', hair: '束发官帽，端肃' },
    { attire: '商朝上大夫，绯色官袍', hair: '冠冕齐整，正直' },
    { attire: '被剜双目，白巾覆眼', hair: '发丝散乱，悲怆' },
    { attire: '掌中双目初睁，飞电枪在手', hair: '道袍布冠，掌目生辉' },
    { attire: '掌中目仙，云霞兽为骑', hair: '道冠仙髻，目光如炬' },
    { attire: '道德真君之姿，掌目洞彻天地', hair: '仙光化发，慈严并具' },
  ],
  wei_hu: [
    { attire: '金刚灵种，佛光初凝', hair: '短发，僧俗未定' },
    { attire: '护法少年，素衣', hair: '束发，端正' },
    { attire: '金甲初着，降魔杵在手', hair: '发冠，威仪渐显' },
    { attire: '护法金刚相，金甲怒目', hair: '怒目圆睁，法相庄严' },
    { attire: '韦陀护法，降魔杵镇邪', hair: '金冠宝相，慈悲怒目' },
    { attire: '三教护法，金刚不坏', hair: '神光化髻，万邪辟易' },
  ],
  daji: [
    { attire: '妖狐之卵，绯红幽光', hair: '狐尾虚影蜷绕' },
    { attire: '幼年狐妖，素衫', hair: '青丝如墨，狐尾初藏' },
    { attire: '入宫华服，宫装初着', hair: '云髻初盘，珠钗点缀' },
    { attire: '九尾华裳，狐火隐现', hair: '高髻金步摇，媚眼如丝' },
    { attire: '九尾祸世，长裙曳地', hair: '云髻高绾，狐尾尽展' },
    { attire: '九尾齐天，祸乱天下', hair: '九尾尽展，魅惑万世' },
  ],
  shen_gongbao: [
    { attire: '邪道灵光，黑袍虚影', hair: '束发，面带谄笑' },
    { attire: '阐教弟子，青白道袍', hair: '道髻初束，眼神闪烁' },
    { attire: '叛教下山，玄黑道袍', hair: '道髻斜歪，黑虎随行' },
    { attire: '黑虎为骑，双剑在手', hair: '乱发披肩，三寸舌如刀' },
    { attire: '说客之相，五岳三山皆访', hair: '长须垂胸，眼神阴鸷' },
    { attire: '翻江倒海，三寸舌搅动乾坤', hair: '长须飞扬，舌灿莲花' },
  ],
  jiang_ziya: [
    { attire: '灵光种子中的垂钓老影', hair: '白发，斗笠未戴' },
    { attire: '渭水边布衣老翁，直钩垂钓', hair: '白发苍苍，束发布巾' },
    { attire: '杏黄旗初展，道袍素净', hair: '道髻高束，眉目慈和' },
    { attire: '封神榜启，杏黄旗挥', hair: '道冠峨立，仙风道骨' },
    { attire: '姜太公临，金甲法衣', hair: '仙髻法冠，打神鞭在手' },
    { attire: '封神执榜人，坐于云端', hair: '白发如银，慈眉垂目' },
  ],
  lv_dongbin: [
    { attire: '黄粱书生，青衫负剑', hair: '书生巾，鬓发如裁' },
    { attire: '儒生装束，弃儒从道', hair: '束发，眉目清俊' },
    { attire: '灰白道袍，纯阳剑在手', hair: '道髻，逍遥巾束发' },
    { attire: '五雷天心，道袍猎猎', hair: '长发飞扬，剑眉入鬓' },
    { attire: '纯阳帝君，白衣佩剑', hair: '仙髻金冠，剑仙风姿' },
    { attire: '纯阳剑仙，飞剑千里斩妖', hair: '仙髻金冠，剑意冲霄' },
  ],
  he_xiangu: [
    { attire: '采药女虚影，竹篓在背', hair: '青丝，荆钗布裙' },
    { attire: '岭南采药女，素裙', hair: '发辫垂肩，灵秀' },
    { attire: '八仙之女，云母为饰', hair: '云髻初盘，莲花簪' },
    { attire: '百花阵启，仙裙生辉', hair: '云髻高绾，花钿点额' },
    { attire: '八仙过海，素衣仙袍', hair: '仙髻凤钗，容光摄人' },
    { attire: '瑶池仙姑，莲花在手', hair: '素发仙髻，不染凡尘' },
  ],
  zhang_guolao: [
    { attire: '樵夫虚影，柴担在肩', hair: '白发，粗布巾' },
    { attire: '恒山樵夫，麻布短褐', hair: '白发苍然，腰间酒葫芦' },
    { attire: '倒骑白驴，鱼鼓在手', hair: '白发束起，仙翁之气' },
    { attire: '呼风唤雨，道袍翻飞', hair: '白发飘逸，仙骨道姿' },
    { attire: '鱼鼓天音，混元仙翁', hair: '仙髻道冠，银须飘飘' },
    { attire: '混元仙翁，鱼鼓天音震三界', hair: '银发仙髻，倒骑驴巡天' },
  ],
  tie_guaili: [
    { attire: '巴国隐士虚影，铁拐顿地', hair: '乱发，破帽' },
    { attire: '弃儒从道，布衣', hair: '束发，清癯' },
    { attire: '附身乞丐，破衣跛足', hair: '乱发蓬面，铁拐为杖' },
    { attire: '混元丹术，破袍仙光', hair: '发须皆乱，目藏精光' },
    { attire: '铁拐仙人，宝葫芦纳天地', hair: '乱发仙髻，破中藏道' },
    { attire: '铁拐大仙，破衣而入道', hair: '形残道全，宝葫芦常挂' },
  ],
  han_zhongli: [
    { attire: '武将虚影，红袍未褪', hair: '束发，鬓角霜' },
    { attire: '燕台武将，玄甲红袍', hair: '束发盔缨，战意未消' },
    { attire: '弃武修道，道袍初着', hair: '束发，芭蕉扇藏袖' },
    { attire: '赤阳神功，红袍大氅', hair: '赤发虬髯，目若朗星' },
    { attire: '正阳帝君，红袍金冠', hair: '金冠束发，赤面威严' },
    { attire: '正阳帝君，大日金轮镇山河', hair: '赤发金冠，仙光圆满' },
  ],
  lan_caihe: [
    { attire: '乐童虚影，破绿衫', hair: '赤足，散发' },
    { attire: '市井歌者，绿衫赤足', hair: '散发，笑意洒脱' },
    { attire: '踏歌而行，玉板在手', hair: '发丝飞扬，踏歌自若' },
    { attire: '踏歌成阵，歌声破空', hair: '散发狂歌，仙气外溢' },
    { attire: '九天玄音，歌仙之姿', hair: '发髻松散，逍遥不羁' },
    { attire: '濠梁仙翁，赤足高歌', hair: '白发作歌，笑看红尘' },
  ],
  cao_guojiu: [
    { attire: '公子虚影，锦袍玉带', hair: '束发金冠，华贵' },
    { attire: '锦绣公子，锦袍', hair: '束发玉冠，温润' },
    { attire: '青布道袍，弃官入道', hair: '道髻，眉目清雅' },
    { attire: '礼乐大道，玉磬在手', hair: '道冠，音律为魂' },
    { attire: '九天神音，鹤氅白袍', hair: '莲花冠，玉磬清响' },
    { attire: '国舅真仙，白鹤为伴', hair: '素发仙髻，超然物外' },
  ],
  taishang_laojun: [
    { attire: '上古真人灵光，紫气东来', hair: '白发童子面，未成相' },
    { attire: '青牛老者，紫金道袍初现', hair: '白发白须，面如童子' },
    { attire: '八卦炉前，道袍烟云', hair: '发髻高束，拂尘在手' },
    { attire: '老君炉炼，紫气万丈', hair: '白发怒张，炉火映面' },
    { attire: '道德天尊，紫金仙袍', hair: '仙髻金冠，紫气绕身' },
    { attire: '道祖至尊，三清合一紫气东来', hair: '紫金道冠，万道之源' },
  ],
  zhong_kui: [
    { attire: '寒窗书生灵光，布巾束发', hair: '布巾束发，书生清癯' },
    { attire: '青衫落拓，腰悬旧笔', hair: '发丝散乱，眉宇郁结' },
    { attire: '大红官袍初着，乌纱帽', hair: '乌纱帽微歪，髭须初张' },
    { attire: '驱鬼之王，红袍与剑', hair: '虬髯怒张，豹眼圆睁' },
    { attire: '玄黑蟒袍金纹，紫金冠', hair: '紫金冠，虬髯如戟' },
    { attire: '镇宅捉鬼，判官笔指鬼青锋剑镇邪', hair: '虬髯怒张，紫金冠生威' },
  ],
  // ---------- 星座守护 12 ----------
  aries: [
    { attire: '星屑虚影，白羊命星微光', hair: '束发，眉心一点星芒' },
    { attire: '素白练功服，圣衣未成', hair: '束发，汗珠映星' },
    { attire: '青铜圣衣初覆肩甲', hair: '羊角头冠微启' },
    { attire: '白银圣衣，水晶墙光', hair: '银盔，眉宇沉静' },
    { attire: '黄金圣衣全装，白羊宫主', hair: '白羊金盔，宝相庄严' },
    { attire: '神圣衣星辉，星光灭绝之姿', hair: '星光化冠，超脱凡尘' },
  ],
  taurus: [
    { attire: '星辉灵光，金牛星芒沉稳', hair: '粗发，魁伟身形' },
    { attire: '练功服，赤膊练力', hair: '短发如刺，汗如雨下' },
    { attire: '青铜圣衣，力士之姿', hair: '牛角头冠初现' },
    { attire: '白银圣衣，巨型号角', hair: '双角盔，怒目圆睁' },
    { attire: '黄金圣衣，金牛宫主', hair: '金牛金盔，力撼山岳' },
    { attire: '黄金之角，仁厚巨人', hair: '金盔生辉，憨厚威仪' },
  ],
  gemini: [
    { attire: '双子星辉，一明一暗', hair: '双生虚影，善恶交织' },
    { attire: '练功服，温雅少年', hair: '束发，眼神却晦暗' },
    { attire: '青铜圣衣，善恶分形', hair: '半面面具，半面光明' },
    { attire: '异次元空间，白银圣衣', hair: '黑发狂舞，邪气外露' },
    { attire: '黄金圣衣，双子宫主', hair: '双面金盔，善恶并存' },
    { attire: '银河星爆，神的化身', hair: '神光化发，一念神魔' },
  ],
  cancer: [
    { attire: '冥辉灵光，幽光沉沉', hair: '发丝带冥火' },
    { attire: '练功服，邪气未显', hair: '束发，嘴角冷笑' },
    { attire: '青铜圣衣，冥气缠绕', hair: '蟹甲头冠，诡异' },
    { attire: '积尸气冥界波，冥光冲霄', hair: '冥火为发，阴笑森然' },
    { attire: '黄金圣衣，巨蟹宫主', hair: '金盔蟹钳，邪中带狂' },
    { attire: '冥界之王，身陷幽暗', hair: '冥光吞没，亦正亦邪' },
  ],
  leo: [
    { attire: '星辉灵光，狮鬃般的金焰', hair: '金发，烈如狮鬃' },
    { attire: '练功服，热血修行', hair: '束发，汗意蒸腾' },
    { attire: '青铜圣衣，拳光初现', hair: '狮盔微启，眸光如电' },
    { attire: '闪电光速拳，战意熊熊', hair: '金发飞扬，怒目如狮' },
    { attire: '黄金圣衣，狮子宫主', hair: '黄金狮盔，威风凛凛' },
    { attire: '等离子光速拳，正义之光', hair: '金光化冠，浩然正气' },
  ],
  virgo: [
    { attire: '星辉莲花，禅意初凝', hair: '束发，眉目安详' },
    { attire: '练功服，静坐冥想', hair: '束发，闭目禅坐' },
    { attire: '青铜圣衣，禅相庄严', hair: '发髻微束，宝相初显' },
    { attire: '天魔降伏，神光内敛', hair: '闭目怒相，禅威并具' },
    { attire: '黄金圣衣，处女宫主', hair: '金莲花冠，佛陀相' },
    { attire: '第八感觉醒，最接近神', hair: '圣光化莲，眉目悲悯' },
  ],
  libra: [
    { attire: '星辉龙珠，天秤虚影', hair: '束发，苍老之相' },
    { attire: '五老峰练功，布衣', hair: '须发皆白，佝偻' },
    { attire: '青铜圣衣，庐山升龙霸', hair: '白发，龙吟初响' },
    { attire: '返老还童，白银圣衣', hair: '银发转黑，青年之相' },
    { attire: '黄金圣衣，天秤宫主', hair: '金盔，十二件武器负背' },
    { attire: '老师坐镇，两百年守望', hair: '白发金盔，目光如渊' },
  ],
  scorpio: [
    { attire: '星辉灵光，猩红针芒微露', hair: '束发，眼神犀利' },
    { attire: '练功服，忠义修行', hair: '束发，沉稳' },
    { attire: '青铜圣衣，毒针初芒', hair: '发间一点猩红' },
    { attire: '深红毒针，蝎尾流光', hair: '红发如焰，凛然' },
    { attire: '黄金圣衣，天蝎宫主', hair: '蝎尾金盔，锐利' },
    { attire: '天蝎之怒，深红针下留情', hair: '金盔生辉，义字当先' },
  ],
  sagittarius: [
    { attire: '星辉之箭，弓芒初凝', hair: '束发，俊朗' },
    { attire: '练功服，弓弦初张', hair: '束发，目光如鹰' },
    { attire: '青铜圣衣，黄金之箭雏形', hair: '发丝随风，弓在手' },
    { attire: '黄金之箭，守护而战', hair: '金发飞扬，英烈之气' },
    { attire: '黄金圣衣，射手宫主', hair: '金盔，弓满月' },
    { attire: '英魂长存，星空中守望', hair: '星光化发，永恒之姿' },
  ],
  capricorn: [
    { attire: '星辉之剑，剑芒初凝', hair: '束发，冷峻' },
    { attire: '练功服，剑士修行', hair: '束发，沉静如渊' },
    { attire: '青铜圣衣，圣剑雏形', hair: '发间剑光隐现' },
    { attire: '圣剑Excalibur，剑影纵横', hair: '黑发如墨，眸光如剑' },
    { attire: '黄金圣衣，摩羯宫主', hair: '金盔，双臂为剑' },
    { attire: '觉醒之剑，忠诚所向', hair: '剑光化发，一生为剑' },
  ],
  aquarius: [
    { attire: '星辉冰晶，寒气初凝', hair: '束发，清冷' },
    { attire: '练功服，冰道修行', hair: '束发，凝霜为汗' },
    { attire: '青铜圣衣，冰晶环雏形', hair: '发间凝霜，眼如冰湖' },
    { attire: '绝对零度之道，冰封千里', hair: '白发冰霜，无情之冷' },
    { attire: '黄金圣衣，水瓶宫主', hair: '冰蓝金盔，严师之姿' },
    { attire: '冰河之师，严中带暖', hair: '冰发微融，目光温柔' },
  ],
  pisces: [
    { attire: '星辉花苞，绯红暗香', hair: '束发，俊美如画' },
    { attire: '练功服，花之修行', hair: '束发，眉目如诗' },
    { attire: '青铜圣衣，恶魔玫瑰雏形', hair: '发间玫瑰初绽' },
    { attire: '食人鱼玫瑰，血玫飞射', hair: '金发染血，凄美' },
    { attire: '黄金圣衣，双鱼宫主', hair: '玫瑰金盔，美中藏毒' },
    { attire: '凄美绽放，玫瑰成刃', hair: '花瓣化发，绝美之姿' },
  ],
}
