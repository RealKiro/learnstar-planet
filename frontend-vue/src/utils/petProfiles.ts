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

  // ===== 宝可梦风 =====
  charmander:    { form: '橙色小蜥蜴，尾尖燃着火焰。', habit: '活泼好斗，尾巴的火焰随心情起伏，越兴奋烧得越旺。' },
  bulbasaur:     { form: '蓝绿蛙状小兽，背驮种子苞。', habit: '安静沉着，爱晒太阳，背苞靠阳光积蓄生长的力量。' },
  squirtle:      { form: '蓝色小龟，背驮硬壳。', habit: '机灵顽皮，爱搞恶作剧，龟壳既是铠甲也是家。' },
  eevee:         { form: '棕毛狐狸样小兽，颈毛蓬松。', habit: '适应力极强，可进化出多种形态，温顺而亲人。' },
  pikachu:       { form: '黄色小鼠，双颊红晕。', habit: '活泼灵动，放电时脸颊噼啪作响，是快乐的发电站。' },
  riolu:         { form: '蓝黑相间的小犬，胸前黑纹。', habit: '能感知波导之力，善良正义，总为弱者挺身而出。' },

  // ===== 国宝守护 =====
  panda:         { form: '黑白圆滚，黑耳黑眼圈。', habit: '爱啃竹子，慵懒温和，是行走的国宝萌物。' },
  golden_monkey: { form: '金色长毛，蓝脸朝天鼻。', habit: '活泼好动，群居深山，攀援跳跃如履平地。' },
  south_china_tiger: { form: '橙底黑纹，身形矫健。', habit: '独行山林的王者，喜夜间巡猎，虎啸生风。' },
  red_crowned_crane: { form: '白羽黑颈，头顶一点朱红。', habit: '高洁优雅，屹立如松，鸣声嘹亮穿云。' },
  crested_ibis:  { form: '白羽粉翼，头冠蓬松。', habit: '珍稀祥瑞，栖于水田浅滩，以小鱼虾为食。' },
  tibetan_antelope: { form: '白褐相间，长角如剑。', habit: '高原上的奔跑健将，角细长而直，如剑指天。' },
  snow_leopard:  { form: '灰白底黑斑，长尾粗壮。', habit: '雪线之上的独行侠，攀岩走壁如履平地。' },
  milu_deer:     { form: '角似鹿非鹿，蹄似牛非牛。', habit: '人称四不像，喜水泽之地，性情温和从容。' },
  siberian_tiger: { form: '硕大虎躯，厚毛如披风。', habit: '森林之王，喜寒惧热，领地意识极强。' },
  red_panda:     { form: '红褐毛皮，环纹大尾。', habit: '慵懒可爱，爱趴在树枝上打盹，以竹叶为食。' },

  // ===== 魔法奇幻 =====
  unicorn:       { form: '纯白骏马，额生螺旋独角。', habit: '纯洁神圣，只亲近心灵纯净之人，角光能驱邪。' },
  griffin:       { form: '鹰首狮身，金色双翼。', habit: '天空与陆地的双重霸主，忠勇无双，翱翔万里。' },
  wand_cat:      { form: '紫毛小猫，爪握魔法杖。', habit: '神秘机敏，挥杖即施咒，是巫师的灵宠。' },
  nightmare_horse: { form: '紫焰骏马，四蹄踏火。', habit: '奔跑时烈焰环绕，夜行千里，桀骜难驯。' },
  fairy:         { form: '金翼小仙子，身绕光点。', habit: '善良灵动，在花间飞舞，将希望播撒四方。' },
  treant:        { form: '深绿树人，枝干为臂。', habit: '守护森林的智者，沉稳慈悲，与古木同寿。' },
  mermaid:       { form: '碧鳞人鱼，长发如波。', habit: '歌声惑人，遨游深海，却向往人间的灯火。' },
  grey_wizard:   { form: '灰袍巫师，手持法杖。', habit: '深藏不露，精通奥术，只在关键时刻出手。' },
  dragon_knight: { form: '红铠骑士，身侧有龙影。', habit: '既屠龙亦御龙，忠诚勇敢，身经百战。' },
  alchemy_golem: { form: '金甲魔像，符文刻身。', habit: '由炼金术驱动，力大无穷，唯主人之命是从。' },
  lamp_spirit:   { form: '金灯灵体，灯火摇曳。', habit: '夜夜为人点灯引路，灯芯不灭则初心不改。' },

  // ===== 史前生物 =====
  triceratops:   { form: '大颈盾配三尖角，身强体壮。', habit: '群居食草，角盾护群，憨厚却不可欺。' },
  pterosaur:     { form: '展翼无齿，颈长身轻。', habit: '在远古天空滑翔，俯冲湖面捕鱼为生。' },
  mammoth:       { form: '长毛巨象，弯长象牙。', habit: '冰河时代的庞然大物，群居迁徙，力可撼树。' },
  sabertooth:    { form: '上颚獠牙如剑，肌肉虬结。', habit: '凶猛的猎手，一击毙敌，獠牙是其信条。' },
  megalodon:     { form: '巨牙利齿的远古巨鲨。', habit: '海洋霸主，猎尽海中巨物，深海也为其让路。' },
  ground_sloth:  { form: '巨爪大懒兽，毛长体壮。', habit: '看似笨拙实则力大，巨爪可掘土开洞。' },
  woolly_rhino:  { form: '长毛犀牛，一双弯角。', habit: '冰原巨兽，厚毛御寒，以苔草为食。' },

  // ===== 星座守护 =====
  aries:         { form: '金色羊首星灵，弯角如月。', habit: '朝气蓬勃，冲劲十足，是黄道十二宫的第一宫。' },
  taurus:        { form: '橙黄牛首星灵，角坚如铁。', habit: '沉稳务实，认定目标便锲而不舍，绝不回头。' },
  gemini:        { form: '翠绿双面星灵，一体双相。', habit: '思维敏捷，兴趣多变，两个心灵活成一个意志。' },
  cancer:        { form: '银白蟹形星灵，双螯护身。', habit: '恋家念旧，外壳坚硬，内里却是最柔软的心。' },
  leo:           { form: '金鬃狮首星灵，威光赫赫。', habit: '天生的王者，热情慷慨，誓死守护自己的领地。' },
  virgo:         { form: '纯白星灵，手执麦穗。', habit: '追求完美，细致入微，是智慧与纯洁的化身。' },
  libra:         { form: '粉紫星灵，身前悬天秤。', habit: '讲求公平，权衡利弊，一生追寻和谐之美。' },
  scorpio:       { form: '深红蝎形星灵，尾刺带毒。', habit: '深沉敏锐，蛰伏待机，重情亦记仇。' },
  sagittarius:   { form: '紫光半人马，背负弓箭。', habit: '自由奔放，箭指远方，为追逐理想永不停蹄。' },
  capricorn:     { form: '蓝鳞羊尾，踏浪登峰。', habit: '坚韧不拔，目标长远，一步一个脚印向上攀。' },
  aquarius:      { form: '天蓝星灵，手捧宝瓶。', habit: '思想超前，乐于奉献，瓶中智慧之水长流。' },
  pisces:        { form: '碧蓝双鱼，首尾相衔。', habit: '温柔多梦，共情力强，游弋于现实与幻想之间。' },

  // ===== 民间传说 =====
  nian:          { form: '红毛独角，獠牙外露。', habit: '岁末而出，畏爆竹与红色，后成驱邪纳吉之兽。' },
  taotie:        { form: '巨口凶兽，双角狰狞。', habit: '贪食无度，吞天噬地，名列四大凶兽。' },
  baize:         { form: '白身独角，通体智慧纹。', habit: '知天下万物之名，能言善解，是智慧的祥兽。' },
  lion_dance:    { form: '红金狮头，绣球相伴。', habit: '逢年过节起舞助兴，驱邪纳吉，热闹喜庆。' },
  magpie:        { form: '青黑白腹，长尾如扇。', habit: '报喜之鸟，爱在枝头叽喳，最讨人欢喜。' },
  firework_spirit: { form: '粉色灵光，烟花形态。', habit: '逢节绽放，璀璨一瞬却也灿烂一瞬。' },
  zongzi:        { form: '翠绿粽形，系着绳结。', habit: '端午登场，糯米裹着思念，清香四溢。' },
  tangyuan:      { form: '雪白圆润，浮于甜汤。', habit: '元宵团圆之食，软糯甜糯，寓意和美。' },
  mooncake:      { form: '棕金饼形，压花纹章。', habit: '中秋月圆时登场，内藏豆沙蛋黄，甜而不腻。' },
  qingtuan:      { form: '碧绿糯米团，艾草清香。', habit: '清明时节江南味，裹着豆沙馅，软糯可口。' },
  chongyang_cake:{ form: '金色米糕，层叠如塔。', habit: '重阳登高时相随，松软香甜，寓意步步高。' },
  niangao:       { form: '白糯方糕，软糯弹牙。', habit: '年节必备，谐音年年高，蒸煮皆宜。' },
  laba_porridge: { form: '五谷杂粮熬成的浓粥。', habit: '腊八之日一碗暖，百味杂陈也是团圆味。' },
  spring_pancake:{ form: '薄嫩面饼，卷着春蔬。', habit: '立春咬春，一卷新绿，咬出满口春意。' },
  tanghulu:      { form: '冰糖山楂串，红亮晶莹。', habit: '街头的酸甜童年，外层脆甜内里酸。' },
  osmanthus_cake:{ form: '米白方糕，点缀桂花瓣。', habit: '秋日桂香入糕，清甜温润，满是金秋气息。' },
  wonton:        { form: '薄皮小馄饨，汤鲜味美。', habit: '街巷里的一碗热汤，皮薄馅嫩，暖心暖胃。' },
  festival_lantern: { form: '红彩花灯，烛光摇曳。', habit: '元宵提灯夜游，照亮团圆的路。' },

  // ===== 七侠剑客 =====
  hongmao:       { form: '红毛小猫，手持长虹剑。', habit: '机智勇敢，侠肝义胆，长虹剑法出神入化。' },
  lantu:         { form: '蓝毛白兔，冰魄剑在手。', habit: '冰雪聪明，温柔而坚定，剑出便是一道冰魄寒光。' },
  doudou:        { form: '白色小狗，爱玩爱闹。', habit: '调皮捣蛋却心地纯善，是七侠中的开心果。' },
  dabeng:        { form: '棕毛大熊，天生神力。', habit: '豪爽憨厚，力大无穷，奔雷剑势大力沉。' },
  tiaotiao:      { form: '绿毛小猴，身轻如燕。', habit: '古灵精怪，跳跃如飞，青光剑快若闪电。' },
  shali:         { form: '紫毛小猫，温柔坚韧。', habit: '外柔内刚，剑法细腻，紫云剑名动江湖。' },
  dada:          { form: '绿面熊猫，沉稳大气。', habit: '文雅温和，太极剑法以柔克刚，四两拨千斤。' },
  qilin_sacred:  { form: '金鳞瑞兽，祥云环绕。', habit: '威严神圣，是正义的化身，默默守护七侠。' },
  heixinhu:      { form: '黑毛猛虎，双目如炬。', habit: '争强好胜，亦正亦邪，虎亦有虎的尊严。' },
  zhuzhijie:     { form: '粉色小猫，温柔爱美。', habit: '细腻体贴，是队伍里的柔情担当。' },
  niuxuanfeng:   { form: '灰毛壮牛，弯角如刀。', habit: '忠厚勤劳，力大沉稳，一诺千金。' },

  // ===== 水生 =====
  mosasaur:      { form: '深海巨蜥龙，鳍状四肢。', habit: '远古海洋的顶级掠食者，游弋于苍茫汪洋。' },
  finless_porpoise: { form: '圆头无背鳍，灰蓝光滑。', habit: '性情温顺爱嬉水，嘴角常挂微笑，是江中的微笑天使。' },
  dragon_boat:   { form: '龙首舟身，彩旗飘飘。', habit: '端午竞渡的主角，锣鼓一响便破浪争先。' },

  // ===== 鸟 =====
  zhuque:        { form: '赤红神鸟，周身烈焰。', habit: '四象之一镇守南方，浴火重生，永不言败。' },
  wind_falcon:   { form: '翠绿猛隼，翼疾如风。', habit: '在高空盘旋，俯冲捕猎，是天空的风之子。' },
  lingge:        { form: '淡蓝灵鸟，翎羽如歌。', habit: '歌声婉转能抚慰伤痛，是七侠的传信使者。' },

  // ===== 机甲系 =====
  lightsaber_warrior: { form: '紫铠机甲，手持光剑。', habit: '崇尚荣耀，剑术无双，为正义而战。' },
  fission_giant: { form: '蓝甲巨人，能量核心。', habit: '力能裂石，遇强则强，是移动的堡垒。' },
  nano_swarm:    { form: '绿光虫群，聚散无常。', habit: '成群结队，纳米重组，无孔不入。' },
  storm_jet:     { form: '流线战机，蓝翼破空。', habit: '超音速巡航，与风暴为伴，疾如闪电。' },
  bio_armor:     { form: '翠绿生化铠甲。', habit: '与宿主共生，能自我修复，越战越坚。' },
  starship_core: { form: '紫光核心，星轨环绕。', habit: '星舰的心脏，能量无穷，静待启航。' },

  // ===== 人形·神灵 =====
  fu_star:       { form: '红衣福童，怀抱福字。', habit: '天官赐福，所到之处福气满满。' },
  shou_star:     { form: '白须仙翁，手拄寿杖。', habit: '寿比南山，慈祥仁厚，赠人长寿。' },
  god_of_wealth: { form: '金袍财神，怀抱金元宝。', habit: '招财进宝，逢年过节人皆奉之，喜气盈门。' },
  door_god:      { form: '赤面武将，手持门戟。', habit: '镇守门户，驱邪避鬼，威严刚正不阿。' },
  kitchen_god:   { form: '橙袍灶君，端坐灶前。', habit: '掌管一家烟火，岁末上天言一家好事。' },

  // ===== 补齐 =====
  lantern:       { form: '白纸灯笼，竹骨为身。', habit: '夜里为人照路，烛火摇曳间传递人间温暖。' },
  dumpling:      { form: '白皮元宝形，褶边如花。', habit: '冬至除夕的团圆味，馅料藏满家的牵挂。' },
  qinglong:      { form: '青鳞神龙，鹿角长须。', habit: '四象之首镇东方，行云布雨，佑一方风调雨顺。' },
  baihu:         { form: '白毛巨虎，金纹隐现。', habit: '四象之一镇西方，威震百兽，为战神之象。' },
  ice_fox:       { form: '冰蓝狐狸，尾带霜华。', habit: '生于极寒雪原，行动如风，所过之处凝霜成路。' },
  rock_rhino:    { form: '岩石甲身，额生尖角。', habit: '如磐石般沉稳，角可撞碎山岩，一步一个坑。' },
  light_deer:    { form: '金光神鹿，鹿角如炬。', habit: '在光影中穿行，能驱散黑暗，是光的信使。' },
  dark_panther:  { form: '墨黑豹影，身姿流线。', habit: '潜伏于暗夜，出击快如闪电，无声无息。' },
  steel_armadillo: { form: '钢甲犰狳，甲带如鳞。', habit: '遇敌便蜷成铁球，刀枪不入，坚不可摧。' },
}

export function getPetProfile(speciesId: string): PetProfile {
  return PET_PROFILES[speciesId] || { form: '形态尚待探明的神秘宠物。', habit: '习性未知，正等待与小主人一起探索。' }
}
