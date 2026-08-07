# 宠物 AI 生图提示词 · 全 125 物种 × 6 阶段（v2）

> 生成：2026-08-07 · 由 `scripts/generate-pet-prompts-v2.js` 自动生成，**不要手改**（改数据源后重跑）
> 取代 v1 `pet-prompts.csv` 通用模板。生成图片后按 `docs/pet-image-manifest.md` 流程登记（MANIFEST + 拷贝到 backend/public/pets/）。

## 一、使用说明

- **中文主句**（每段首行）：复制给即梦 / 通义万相 / Stable Diffusion（中文模型）。
- **EN 行**：复制给 Midjourney / DALL·E / Stable Diffusion（英文模型）。
- **六阶等级制**：灵胎初醒(Lv1) → 凡尘砺心(Lv3) → 道法初成(Lv5) → 大劫淬炼(Lv7) → 封神登天(Lv9) → 道果圆满(Lv11)。
- **精修四维度**：每阶段含 神态 / 动作 / 衣着 / 梳造 四维描述（重点角色手写，其余按大类基座），见 `frontend-vue/src/utils/petRefine.ts`。
- **人生档案版**：东方神话 姜子牙 / 杨戬 / 雷震子（有六阶段人生档案）的提示词为「人生叙事」版（品性 / 姿态 / 服饰 / 功法 / 画面 / 台词 / 诗词 / 主题句），源自 `frontend-vue/src/utils/petLifeStories.ts`，与手绘 SVG 完全对齐；其余角色走抽象道行版（petRefine）。
- **画布**：800×1000 竖版，角色主体居中偏下占约 60%，正面 3/4 视角，无文字无水印。
- **命名规范**：产物存 `frontend-vue/public/pets/{seriesId}/{speciesId}-{stage}.webp`。
- **版权提醒**：宝可梦 / 数码宝贝 / 星座圣斗士 / 虹猫蓝兔为受保护 IP，本提示词按项目「版权直名、非商用致敬」政策直接使用角色名，请勿用于商用。

## 二、初始形态分类规范（灵胎初醒 / Lv1）

- **设计原则**：初始形态符合物种特点 —— 植物类=一粒种子 / 胎生哺乳=一团胚胎 / 卵生=一枚蛋卵 / 数码兽=数码蛋 / 神兽瑞兽=灵种灵卵 / 人物=灵光化形 / 食物=原料 / 器物=物件 / 剑客=剑意灵光
- **属性区别**：每个角色初始形态带十二属性（木火土金水 + 雷冰风光暗钢秘），由 `frontend-vue/src/utils/petInitial.ts` 定义
- **分布**：seed 种子 4 / embryo 胚胎 13 / egg 蛋卵 37 / digi 数码蛋 6 / spirit 灵光 14 / person 人物 20 / food 原料 11 / artifact 器物 3 / sword 剑意 9 / mythic 神兽灵种 10
- **形态名**建议同步 petData 的 Lv1 等级名（保持一致后 生图 / UI / 诗文 三方对齐）

## 三、系列视觉规范总表

| 系列 | 物种数 | 风格（中文） | 风格（英文） |
|------|--------|--------------|--------------|
| 山海经 | 18 | 新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘 | Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture |
| 东方神话 | 20 | 新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长 | xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura |
| 宝可梦 | 12 | 经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈 | classic Japanese anime adventure style, cel-shaded, vibrant and wholesome |
| 数码宝贝 | 6 | 日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻 | Japanese digital evolution anime, cel shading with digital particles, techno glow |
| 国宝 | 12 | 皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖 | Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop |
| 魔法奇幻 | 12 | 欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗 | western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere |
| 史前生物 | 12 | BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫 | paleoart documentary style with cartoon charm, realistic scale/skin texture |
| 星座守护 | 12 | 华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高 | golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs |
| 传统节日 | 12 | 民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈 | Chinese folk festival mascot illustration, nianhua palette, round and festive |
| 虹猫蓝兔七侠传 | 10 | 2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇 | 2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop |

## 四、分系列完整提示词

### 1. 山海经（18 物种）

> **风格**：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。**阶段演绎**：
> - 灵胎初醒：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘（混沌玄青/青铜褐）
> - 凡尘砺心：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则（石青灰/赭石棕）
> - 道法初成：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量（朱砂红/鎏金）
> - 大劫淬炼：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变（玄墨黑/血赤）
> - 封神登天：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方（鎏金/朱紫）
> - 道果圆满：图腾显圣，山海为印，洪荒万古唯此一尊（云白/淡金）

#### 烛龙（`zhulong`）

**灵胎初醒 · 烛火微光**
- 烛龙，灵胎初醒阶段·烛火微光。初始形态：一粒烛火灵种，火光如豆却照彻幽冥，烛龙始祖的气息在火苗中蜷曲，昼夜之神初醒。火属性灵光微微环绕。神态：龙息轻吐，沉睡于混沌灵光。动作：盘蜷于光，沉眠未醒。衣着：幼嫩鳞胚，泛着初光。梳造：无角，须影初现。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 龙息轻吐，沉睡于混沌灵光; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 衔烛之苗**
- 烛龙，凡尘砺心阶段·衔烛之苗。形象：人面蛇身，通体赤红，口衔火烛。 核心意象：烛火、昼夜、钟山。神态：竖瞳初睁，窥探云水。动作：初探云水，游弋学步。衣着：鳞甲渐密，颜色初显。梳造：角芽微露，须丝轻扬。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 竖瞳初睁，窥探云水; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 衔烛而行**
- 烛龙，道法初成阶段·衔烛而行。形象：人面蛇身，通体赤红，口衔火烛。 核心意象：烛火、昼夜、钟山。神态：眸光锐亮，潜龙欲腾。动作：腾空而起，初显神威，开眼为昼，闭眼为夜，呼吸之间风雷自生。衣着：鳞甲生辉，腹光流转。梳造：双角初成，须如流云。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眸光锐亮，潜龙欲腾; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 昼夜使者**
- 烛龙，大劫淬炼阶段·昼夜使者。形象：人面蛇身，通体赤红，口衔火烛。 核心意象：烛火、昼夜、钟山。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 烛龙巡天**
- 烛龙，封神登天阶段·烛龙巡天。形象：人面蛇身，通体赤红，口衔火烛。 核心意象：烛火、昼夜、钟山。神态：受封龙君，神威赫赫。动作：登天行云，布雨泽四方。衣着：金鳞覆身，受冕祥光。梳造：龙角如珊瑚，加冕为尊。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封龙君，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 开目为昼**
- 烛龙，道果圆满阶段·开目为昼。形象：人面蛇身，通体赤红，口衔火烛。 核心意象：烛火、昼夜、钟山。神态：万龙之源，睥睨三界。动作：真身化岳，日月为伴，开眼为昼，闭眼为夜，呼吸之间风雷自生。衣着：龙身映日月，鳞甲如星辰。梳造：龙角擎天，道纹绕体。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 万龙之源，睥睨三界; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 应龙（`yinglong`）

**灵胎初醒 · 虺卵**
- 应龙，灵胎初醒阶段·虺卵。初始形态：一枚虺卵，青灰色蛋壳覆着细密金鳞纹，隐隐透出风雨雷光，是应龙化形前的沉睡之卵。水属性灵光微微环绕。神态：龙息轻吐，沉睡于混沌灵光。动作：盘蜷于光，沉眠未醒。衣着：幼嫩鳞胚，泛着初光。梳造：无角，须影初现。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 龙息轻吐，沉睡于混沌灵光; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 虺行**
- 应龙，凡尘砺心阶段·虺行。形象：身披金鳞，背生双翼的神龙。 核心意象：金鳞、双翼、风雷云雨。神态：竖瞳初睁，窥探云水。动作：初探云水，游弋学步。衣着：鳞甲渐密，颜色初显。梳造：角芽微露，须丝轻扬。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 竖瞳初睁，窥探云水; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 蛟龙**
- 应龙，道法初成阶段·蛟龙。形象：身披金鳞，背生双翼的神龙。 核心意象：金鳞、双翼、风雷云雨。神态：眸光锐亮，潜龙欲腾。动作：腾空而起，初显神威，双翼一展，云腾雨至，江河为之让路。衣着：鳞甲生辉，腹光流转。梳造：双角初成，须如流云。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眸光锐亮，潜龙欲腾; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 应龙雏形**
- 应龙，大劫淬炼阶段·应龙雏形。形象：身披金鳞，背生双翼的神龙。 核心意象：金鳞、双翼、风雷云雨。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 战龙**
- 应龙，封神登天阶段·战龙。形象：身披金鳞，背生双翼的神龙。 核心意象：金鳞、双翼、风雷云雨。神态：受封龙君，神威赫赫。动作：登天行云，布雨泽四方。衣着：金鳞覆身，受冕祥光。梳造：龙角如珊瑚，加冕为尊。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封龙君，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 神龙降世**
- 应龙，道果圆满阶段·神龙降世。形象：身披金鳞，背生双翼的神龙。 核心意象：金鳞、双翼、风雷云雨。神态：万龙之源，睥睨三界。动作：真身化岳，日月为伴，双翼一展，云腾雨至，江河为之让路。衣着：龙身映日月，鳞甲如星辰。梳造：龙角擎天，道纹绕体。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 万龙之源，睥睨三界; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 九尾狐（`nine_tail_fox`）

**灵胎初醒 · 灵狐蛋**
- 九尾狐，灵胎初醒阶段·灵狐蛋。初始形态：一枚灵狐蛋，月白蛋壳泛着清辉，尾尖般的银色光毫在壳上流转，九尾之灵伏卧其中。水属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 白狐**
- 九尾狐，凡尘砺心阶段·白狐。形象：通体雪白，身后拖九条长尾。 核心意象：九尾、青丘山、月光。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 三尾灵狐**
- 九尾狐，道法初成阶段·三尾灵狐。形象：通体雪白，身后拖九条长尾。 核心意象：九尾、青丘山、月光。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，九尾摇曳如月华流泻，一顾倾人再顾倾国。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 五尾**
- 九尾狐，大劫淬炼阶段·五尾。形象：通体雪白，身后拖九条长尾。 核心意象：九尾、青丘山、月光。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 七尾**
- 九尾狐，封神登天阶段·七尾。形象：通体雪白，身后拖九条长尾。 核心意象：九尾、青丘山、月光。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 九尾天狐**
- 九尾狐，道果圆满阶段·九尾天狐。形象：通体雪白，身后拖九条长尾。 核心意象：九尾、青丘山、月光。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，九尾摇曳如月华流泻，一顾倾人再顾倾国。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 鲲鹏（`kunpeng`）

**灵胎初醒 · 鲲苗**
- 鲲鹏，灵胎初醒阶段·鲲苗。初始形态：一尾鲲苗，通体青蓝如墨玉，在混沌灵光中游弋，鳍翼初具而鹏形未显。水属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 巨鲲**
- 鲲鹏，凡尘砺心阶段·巨鲲。形象：北冥之巨鱼，可化身为鸟。 核心意象：北冥、垂天之云、九万里长空。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 化鲲**
- 鲲鹏，道法初成阶段·化鲲。形象：北冥之巨鱼，可化身为鸟。 核心意象：北冥、垂天之云、九万里长空。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，振翅扶摇九万里，水击三千里，绝云气负青天。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 鹏雏**
- 鲲鹏，大劫淬炼阶段·鹏雏。形象：北冥之巨鱼，可化身为鸟。 核心意象：北冥、垂天之云、九万里长空。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 鲲鹏**
- 鲲鹏，封神登天阶段·鲲鹏。形象：北冥之巨鱼，可化身为鸟。 核心意象：北冥、垂天之云、九万里长空。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 北冥之主**
- 鲲鹏，道果圆满阶段·北冥之主。形象：北冥之巨鱼，可化身为鸟。 核心意象：北冥、垂天之云、九万里长空。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，振翅扶摇九万里，水击三千里，绝云气负青天。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 凤凰（`fenghuang`）

**灵胎初醒 · 火羽蛋**
- 凤凰，灵胎初醒阶段·火羽蛋。初始形态：一枚火羽蛋，赤色蛋壳缀满五色羽纹，暖光从壳内透出，凤鸣在蛋中回响。火属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼凤**
- 凤凰，凡尘砺心阶段·幼凤。形象：五彩华羽，尾羽如虹。 核心意象：梧桐、竹实、醴泉。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 赤凤**
- 凤凰，道法初成阶段·赤凤。形象：五彩华羽，尾羽如虹。 核心意象：梧桐、竹实、醴泉。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，五色并举，振翅则百鸟来朝，一鸣则天下太平。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 凤雏**
- 凤凰，大劫淬炼阶段·凤雏。形象：五彩华羽，尾羽如虹。 核心意象：梧桐、竹实、醴泉。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 神凤**
- 凤凰，封神登天阶段·神凤。形象：五彩华羽，尾羽如虹。 核心意象：梧桐、竹实、醴泉。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 涅槃之凤**
- 凤凰，道果圆满阶段·涅槃之凤。形象：五彩华羽，尾羽如虹。 核心意象：梧桐、竹实、醴泉。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，五色并举，振翅则百鸟来朝，一鸣则天下太平。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 麒麟（`qilin`）

**灵胎初醒 · 瑞兽玉**
- 麒麟，灵胎初醒阶段·瑞兽玉。初始形态：一枚温润瑞兽玉，玉身内蕴祥云纹，麟角初形若隐若现，仁兽之息如春风化雨。土属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼麟**
- 麒麟，凡尘砺心阶段·幼麟。形象：麋身牛尾，马蹄而一角。 核心意象：独角、麋身、太平之兆。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 仁兽**
- 麒麟，道法初成阶段·仁兽。形象：麋身牛尾，马蹄而一角。 核心意象：独角、麋身、太平之兆。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，不踏生草不履生虫，步生祥云，角有瑞光。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 踏火麒麟**
- 麒麟，大劫淬炼阶段·踏火麒麟。形象：麋身牛尾，马蹄而一角。 核心意象：独角、麋身、太平之兆。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 天麟**
- 麒麟，封神登天阶段·天麟。形象：麋身牛尾，马蹄而一角。 核心意象：独角、麋身、太平之兆。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 万瑞之祖**
- 麒麟，道果圆满阶段·万瑞之祖。形象：麋身牛尾，马蹄而一角。 核心意象：独角、麋身、太平之兆。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，不踏生草不履生虫，步生祥云，角有瑞光。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 穷奇（`qiongqi`）

**灵胎初醒 · 凶兽卵**
- 穷奇，灵胎初醒阶段·凶兽卵。初始形态：一枚凶兽卵，暗色蛋壳覆着虎纹与翼影，煞气凝成黑气在壳面游走。土属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 虎翼兽**
- 穷奇，凡尘砺心阶段·虎翼兽。形象：状如虎而生双翼，遍体凶煞。 核心意象：虎身双翼、凶煞之气、颠倒黑白。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 翼虎**
- 穷奇，道法初成阶段·翼虎。形象：状如虎而生双翼，遍体凶煞。 核心意象：虎身双翼、凶煞之气、颠倒黑白。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，振翅扑食，专噬忠善，风随其身。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 穷奇·残暴**
- 穷奇，大劫淬炼阶段·穷奇·残暴。形象：状如虎而生双翼，遍体凶煞。 核心意象：虎身双翼、凶煞之气、颠倒黑白。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 穷奇·噬恶**
- 穷奇，封神登天阶段·穷奇·噬恶。形象：状如虎而生双翼，遍体凶煞。 核心意象：虎身双翼、凶煞之气、颠倒黑白。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 四凶之穷奇**
- 穷奇，道果圆满阶段·四凶之穷奇。形象：状如虎而生双翼，遍体凶煞。 核心意象：虎身双翼、凶煞之气、颠倒黑白。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，振翅扑食，专噬忠善，风随其身。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 毕方（`bifang`）

**灵胎初醒 · 毕方卵**
- 毕方，灵胎初醒阶段·毕方卵。初始形态：一枚毕方卵，青白蛋壳上燃着一点赤红火纹，单足立火的雏形在壳中孕育。火属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 一足鸟**
- 毕方，凡尘砺心阶段·一足鸟。形象：青色单足之鸟，白喙赤足。 核心意象：单足、白喙、青色火焰。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 青鸟**
- 毕方，道法初成阶段·青鸟。形象：青色单足之鸟，白喙赤足。 核心意象：单足、白喙、青色火焰。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，一足独立，长鸣如鹤，过处即有讹火。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 毕方·衔火**
- 毕方，大劫淬炼阶段·毕方·衔火。形象：青色单足之鸟，白喙赤足。 核心意象：单足、白喙、青色火焰。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 毕方·焚风**
- 毕方，封神登天阶段·毕方·焚风。形象：青色单足之鸟，白喙赤足。 核心意象：单足、白喙、青色火焰。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 毕方神鸟**
- 毕方，道果圆满阶段·毕方神鸟。形象：青色单足之鸟，白喙赤足。 核心意象：单足、白喙、青色火焰。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，一足独立，长鸣如鹤，过处即有讹火。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 貔貅（`pixiu`）

**灵胎初醒 · 貔貅玉**
- 貔貅，灵胎初醒阶段·貔貅玉。初始形态：一枚鎏金貔貅玉，玉兽张嘴吞吐财气，金光明灭间纳福之相初成。金属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 貔貅雏形**
- 貔貅，凡尘砺心阶段·貔貅雏形。形象：龙头马身麟脚，背生双翼。 核心意象：龙头麟脚、双翼、招财纳福。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 貔貅·纳财**
- 貔貅，道法初成阶段·貔貅·纳财。形象：龙头马身麟脚，背生双翼。 核心意象：龙头麟脚、双翼、招财纳福。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，张口吸纳四方财气，只进不出，聚而不散。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 貔貅·镇宅**
- 貔貅，大劫淬炼阶段·貔貅·镇宅。形象：龙头马身麟脚，背生双翼。 核心意象：龙头麟脚、双翼、招财纳福。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 貔貅·吞金**
- 貔貅，封神登天阶段·貔貅·吞金。形象：龙头马身麟脚，背生双翼。 核心意象：龙头麟脚、双翼、招财纳福。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 貔貅至尊**
- 貔貅，道果圆满阶段·貔貅至尊。形象：龙头马身麟脚，背生双翼。 核心意象：龙头麟脚、双翼、招财纳福。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，张口吸纳四方财气，只进不出，聚而不散。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 精卫（`jingwei`）

**灵胎初醒 · 精卫卵**
- 精卫，灵胎初醒阶段·精卫卵。初始形态：一枚精卫卵，浅赤蛋壳缀着白羽纹，衔枝填海的执念在壳中轻轻鼓动。木属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 女娃鸟**
- 精卫，凡尘砺心阶段·女娃鸟。形象：形如乌，赤首白喙赤足。 核心意象：木石、东海、不灭之志。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 衔木鸟**
- 精卫，道法初成阶段·衔木鸟。形象：形如乌，赤首白喙赤足。 核心意象：木石、东海、不灭之志。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，衔西山之木石，日夜往复，以填东海。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 精卫·填海**
- 精卫，大劫淬炼阶段·精卫·填海。形象：形如乌，赤首白喙赤足。 核心意象：木石、东海、不灭之志。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 精卫·不屈**
- 精卫，封神登天阶段·精卫·不屈。形象：形如乌，赤首白喙赤足。 核心意象：木石、东海、不灭之志。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 精卫神鸟**
- 精卫，道果圆满阶段·精卫神鸟。形象：形如乌，赤首白喙赤足。 核心意象：木石、东海、不灭之志。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，衔西山之木石，日夜往复，以填东海。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 相柳（`xiangliu`）

**灵胎初醒 · 九首卵**
- 相柳，灵胎初醒阶段·九首卵。初始形态：一枚九首卵，蛇纹盘绕的暗青蛋壳，隐隐裂出九道微光，毒水之气在壳下涌动。水属性灵光微微环绕。神态：龙息轻吐，沉睡于混沌灵光。动作：盘蜷于光，沉眠未醒。衣着：幼嫩鳞胚，泛着初光。梳造：无角，须影初现。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 龙息轻吐，沉睡于混沌灵光; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 双首**
- 相柳，凡尘砺心阶段·双首。形象：九首蛇身，自环其尾。 核心意象：九首、蛇身、水泽之祸。神态：竖瞳初睁，窥探云水。动作：初探云水，游弋学步。衣着：鳞甲渐密，颜色初显。梳造：角芽微露，须丝轻扬。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 竖瞳初睁，窥探云水; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 多首蛇**
- 相柳，道法初成阶段·多首蛇。形象：九首蛇身，自环其尾。 核心意象：九首、蛇身、水泽之祸。神态：眸光锐亮，潜龙欲腾。动作：腾空而起，初显神威，九首齐张，所过之地化为泽国，毒水横流。衣着：鳞甲生辉，腹光流转。梳造：双角初成，须如流云。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眸光锐亮，潜龙欲腾; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 相柳**
- 相柳，大劫淬炼阶段·相柳。形象：九首蛇身，自环其尾。 核心意象：九首、蛇身、水泽之祸。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 凶兽相柳**
- 相柳，封神登天阶段·凶兽相柳。形象：九首蛇身，自环其尾。 核心意象：九首、蛇身、水泽之祸。神态：受封龙君，神威赫赫。动作：登天行云，布雨泽四方。衣着：金鳞覆身，受冕祥光。梳造：龙角如珊瑚，加冕为尊。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封龙君，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 相柳·湮世**
- 相柳，道果圆满阶段·相柳·湮世。形象：九首蛇身，自环其尾。 核心意象：九首、蛇身、水泽之祸。神态：万龙之源，睥睨三界。动作：真身化岳，日月为伴，九首齐张，所过之地化为泽国，毒水横流。衣着：龙身映日月，鳞甲如星辰。梳造：龙角擎天，道纹绕体。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 万龙之源，睥睨三界; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 獬豸（`xiezhi`）

**灵胎初醒 · 獬豸卵**
- 獬豸，灵胎初醒阶段·獬豸卵。初始形态：一枚獬豸卵，雪白蛋壳顶生一点金角凸起，断狱明理的法光在壳面流转。金属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 独角羊**
- 獬豸，凡尘砺心阶段·独角羊。形象：似羊而独角，目光清正。 核心意象：独角、公正、断狱之威。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 神羊**
- 獬豸，道法初成阶段·神羊。形象：似羊而独角，目光清正。 核心意象：独角、公正、断狱之威。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，见人争斗，以独角触其不直者，刚正不阿。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 獬豸·司法**
- 獬豸，大劫淬炼阶段·獬豸·司法。形象：似羊而独角，目光清正。 核心意象：独角、公正、断狱之威。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 獬豸·执法**
- 獬豸，封神登天阶段·獬豸·执法。形象：似羊而独角，目光清正。 核心意象：独角、公正、断狱之威。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 司法神兽**
- 獬豸，道果圆满阶段·司法神兽。形象：似羊而独角，目光清正。 核心意象：独角、公正、断狱之威。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，见人争斗，以独角触其不直者，刚正不阿。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 青龙（`qinglong`）

**灵胎初醒 · 木灵珠**
- 青龙，灵胎初醒阶段·木灵珠。初始形态：一粒木灵珠，青碧如春水凝成，珠身缠着鹿角须影，东方之木生机勃发。木属性灵光微微环绕。神态：龙息轻吐，沉睡于混沌灵光。动作：盘蜷于光，沉眠未醒。衣着：幼嫩鳞胚，泛着初光。梳造：无角，须影初现。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 龙息轻吐，沉睡于混沌灵光; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小青龙**
- 青龙，凡尘砺心阶段·小青龙。形象：青鳞神龙，鹿角长须。 核心意象：青鳞鹿角、东方之位、风调雨顺。神态：竖瞳初睁，窥探云水。动作：初探云水，游弋学步。衣着：鳞甲渐密，颜色初显。梳造：角芽微露，须丝轻扬。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 竖瞳初睁，窥探云水; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 苍龙**
- 青龙，道法初成阶段·苍龙。形象：青鳞神龙，鹿角长须。 核心意象：青鳞鹿角、东方之位、风调雨顺。神态：眸光锐亮，潜龙欲腾。动作：腾空而起，初显神威，腾云驾雾而起，行云布雨泽润四方。衣着：鳞甲生辉，腹光流转。梳造：双角初成，须如流云。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眸光锐亮，潜龙欲腾; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 腾云龙**
- 青龙，大劫淬炼阶段·腾云龙。形象：青鳞神龙，鹿角长须。 核心意象：青鳞鹿角、东方之位、风调雨顺。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 星宿青龙**
- 青龙，封神登天阶段·星宿青龙。形象：青鳞神龙，鹿角长须。 核心意象：青鳞鹿角、东方之位、风调雨顺。神态：受封龙君，神威赫赫。动作：登天行云，布雨泽四方。衣着：金鳞覆身，受冕祥光。梳造：龙角如珊瑚，加冕为尊。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封龙君，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 青龙尊者**
- 青龙，道果圆满阶段·青龙尊者。形象：青鳞神龙，鹿角长须。 核心意象：青鳞鹿角、东方之位、风调雨顺。神态：万龙之源，睥睨三界。动作：真身化岳，日月为伴，腾云驾雾而起，行云布雨泽润四方。衣着：龙身映日月，鳞甲如星辰。梳造：龙角擎天，道纹绕体。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 万龙之源，睥睨三界; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 白虎（`baihu`）

**灵胎初醒 · 金灵卵**
- 白虎，灵胎初醒阶段·金灵卵。初始形态：一枚金灵卵，雪白金纹蛋壳，虎啸之声在金属冷光中隐隐回响。金属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小白虎**
- 白虎，凡尘砺心阶段·小白虎。形象：白毛巨虎，金纹隐现。 核心意象：白毛金纹、西方之位、战神之象。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 啸风虎**
- 白虎，道法初成阶段·啸风虎。形象：白毛巨虎，金纹隐现。 核心意象：白毛金纹、西方之位、战神之象。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，金纹一闪，虎啸裂空威震百兽。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眸光晶亮，意气初显，跃跃欲试; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 啸林虎**
- 白虎，大劫淬炼阶段·啸林虎。形象：白毛巨虎，金纹隐现。 核心意象：白毛金纹、西方之位、战神之象。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 镇西神虎**
- 白虎，封神登天阶段·镇西神虎。形象：白毛巨虎，金纹隐现。 核心意象：白毛金纹、西方之位、战神之象。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 白虎尊者**
- 白虎，道果圆满阶段·白虎尊者。形象：白毛巨虎，金纹隐现。 核心意象：白毛金纹、西方之位、战神之象。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，金纹一闪，虎啸裂空威震百兽。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 神光自照，与天地同尊; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 朱雀（`zhuque`）

**灵胎初醒 · 火灵卵**
- 朱雀，灵胎初醒阶段·火灵卵。初始形态：一枚火灵卵，赤羽纹蛋壳透出烈焰光晕，离火之精在壳中振翅欲鸣。火属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小朱雀**
- 朱雀，凡尘砺心阶段·小朱雀。形象：赤红神鸟，周身烈焰。 核心意象：赤羽烈焰、南方之位、浴火重生。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 离火鸟**
- 朱雀，道法初成阶段·离火鸟。形象：赤红神鸟，周身烈焰。 核心意象：赤羽烈焰、南方之位、浴火重生。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，振翅浴火而起，烈焰燎空化为凤凰。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 焚风朱雀**
- 朱雀，大劫淬炼阶段·焚风朱雀。形象：赤红神鸟，周身烈焰。 核心意象：赤羽烈焰、南方之位、浴火重生。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 浴火朱雀**
- 朱雀，封神登天阶段·浴火朱雀。形象：赤红神鸟，周身烈焰。 核心意象：赤羽烈焰、南方之位、浴火重生。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 朱雀尊者**
- 朱雀，道果圆满阶段·朱雀尊者。形象：赤红神鸟，周身烈焰。 核心意象：赤羽烈焰、南方之位、浴火重生。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，振翅浴火而起，烈焰燎空化为凤凰。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 玄武（`xuanwu`）

**灵胎初醒 · 水灵珠**
- 玄武，灵胎初醒阶段·水灵珠。初始形态：一粒水灵珠，玄青珠身内蕴龟蛇之影，北冥之水在其间缓缓流转。水属性灵光微微环绕。神态：龙息轻吐，沉睡于混沌灵光。动作：盘蜷于光，沉眠未醒。衣着：幼嫩鳞胚，泛着初光。梳造：无角，须影初现。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 龙息轻吐，沉睡于混沌灵光; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小玄龟**
- 玄武，凡尘砺心阶段·小玄龟。形象：龟蛇合体，背负硬壳。 核心意象：龟蛇合体、北冥之水、玄武星宿。神态：竖瞳初睁，窥探云水。动作：初探云水，游弋学步。衣着：鳞甲渐密，颜色初显。梳造：角芽微露，须丝轻扬。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 竖瞳初睁，窥探云水; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 负水龟**
- 玄武，道法初成阶段·负水龟。形象：龟蛇合体，背负硬壳。 核心意象：龟蛇合体、北冥之水、玄武星宿。神态：眸光锐亮，潜龙欲腾。动作：腾空而起，初显神威，龟蛇盘踞，身周四象水光流转，一尾横扫千军。衣着：鳞甲生辉，腹光流转。梳造：双角初成，须如流云。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眸光锐亮，潜龙欲腾; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 驮山龟**
- 玄武，大劫淬炼阶段·驮山龟。形象：龟蛇合体，背负硬壳。 核心意象：龟蛇合体、北冥之水、玄武星宿。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 镇水玄武**
- 玄武，封神登天阶段·镇水玄武。形象：龟蛇合体，背负硬壳。 核心意象：龟蛇合体、北冥之水、玄武星宿。神态：受封龙君，神威赫赫。动作：登天行云，布雨泽四方。衣着：金鳞覆身，受冕祥光。梳造：龙角如珊瑚，加冕为尊。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封龙君，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 玄武尊者**
- 玄武，道果圆满阶段·玄武尊者。形象：龟蛇合体，背负硬壳。 核心意象：龟蛇合体、北冥之水、玄武星宿。神态：万龙之源，睥睨三界。动作：真身化岳，日月为伴，龟蛇盘踞，身周四象水光流转，一尾横扫千军。衣着：龙身映日月，鳞甲如星辰。梳造：龙角擎天，道纹绕体。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 万龙之源，睥睨三界; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 饕餮（`taotie`）

**灵胎初醒 · 食欲之卵**
- 饕餮，灵胎初醒阶段·食欲之卵。初始形态：一枚饕餮卵，暗金蛋壳覆着贪婪纹，巨口虚影在壳面时隐时现，永不知足。土属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小饕餮**
- 饕餮，凡尘砺心阶段·小饕餮。形象：巨口凶兽，双角狰狞。 核心意象：巨口、狰狞双角、贪食之纹。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 大胃兽**
- 饕餮，道法初成阶段·大胃兽。形象：巨口凶兽，双角狰狞。 核心意象：巨口、狰狞双角、贪食之纹。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，巨口一张，吞尽万物永无餍足。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 无餍之兽**
- 饕餮，大劫淬炼阶段·无餍之兽。形象：巨口凶兽，双角狰狞。 核心意象：巨口、狰狞双角、贪食之纹。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 贪噬之相**
- 饕餮，封神登天阶段·贪噬之相。形象：巨口凶兽，双角狰狞。 核心意象：巨口、狰狞双角、贪食之纹。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 饕餮尊者**
- 饕餮，道果圆满阶段·饕餮尊者。形象：巨口凶兽，双角狰狞。 核心意象：巨口、狰狞双角、贪食之纹。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，巨口一张，吞尽万物永无餍足。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 白泽（`baize`）

**灵胎初醒 · 智慧之卵**
- 白泽，灵胎初醒阶段·智慧之卵。初始形态：一枚智慧之卵，素白蛋壳上浮现问号与书页纹，万灵之名在壳中等待被念出。金属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：混沌初开的洪荒氛围，天地未分，异兽灵胎裹着混沌宝光，渺小而神秘。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：混沌玄青（#2A2A35）主调 + 青铜褐（#8B7E6A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; nascent chaos; 灵光中沉睡，兽性未醒的宁静; palette #2A2A35 with #8B7E6A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小问号兽**
- 白泽，凡尘砺心阶段·小问号兽。形象：白身独角，通体智慧纹。 核心意象：白身独角、智慧纹、通晓万物。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：幼兽行走于上古山川之间，穿行风雨与原始密林，初窥天地法则。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：石青灰（#4A5568）主调 + 赭石棕（#A0855B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; wandering ancient wilds; 初踏山川，懵懂而灵慧; palette #4A5568 with #A0855B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 百晓兽**
- 白泽，道法初成阶段·百晓兽。形象：白身独角，通体智慧纹。 核心意象：白身独角、智慧纹、通晓万物。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，开口便能道出天下妖魅鬼怪之名。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：异兽神通初显，与山川灵气共鸣，身上浮现神秘纹路，散发古老力量。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：朱砂红（#8B2E2E）主调 + 鎏金（#D4AF37）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; awakening power; 眼神古老而专注，神通初显; palette #8B2E2E with #D4AF37 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 万卷兽**
- 白泽，大劫淬炼阶段·万卷兽。形象：白身独角，通体智慧纹。 核心意象：白身独角、智慧纹、通晓万物。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 万物之书**
- 白泽，封神登天阶段·万物之书。形象：白身独角，通体智慧纹。 核心意象：白身独角、智慧纹、通晓万物。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 受封瑞兽，祥云拱卫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 白泽尊者**
- 白泽，道果圆满阶段·白泽尊者。形象：白身独角，通体智慧纹。 核心意象：白身独角、智慧纹、通晓万物。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，开口便能道出天下妖魅鬼怪之名。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气化道，福泽天地; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 2. 东方神话（20 物种）

> **风格**：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。**阶段演绎**：
> - 灵胎初醒：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘（玄青/仙金）
> - 凡尘砺心：离开洞府行走人间，红尘历练，磨砺道心（黛青/灰白）
> - 道法初成：道法初成，法宝显威，法相庄严初现（朱砂/明黄）
> - 大劫淬炼：天劫淬炼、心魔试炼，仙体历经大劫而不灭（墨褐/血绛）
> - 封神登天：功行圆满封神登天，位列仙班，金光紫气环绕（紫金/御金）
> - 道果圆满：道果圆满，紫气东来，万法皆通证道果（素白/淡紫）

#### 姜子牙（`jiang_ziya`） · 人生档案版

**灵胎初醒 · 灵光种子**
- 姜子牙，灵胎初醒阶段·朝歌市井（32岁，朝歌卖面·市井隐志）。形象：姜子牙，白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。品性：忠厚老实，不善钻营，朝歌城中卖面粉为生。心中有济世之志却无门路，常被人笑"痴人说梦"。不怨天尤人，每日早起晚归，卖完面粉便在城墙根下读书。。姿态：拂晓即起，挑着面粉担子穿街过巷；收摊后蹲在墙角，以指为笔在地上演算卦理；见到乞丐会分半个饼。。服饰：粗布褐衣，肩头有补丁；发髻用木簪别着，鬓角已有几缕白发（三十二岁的人看起来像四十岁）；脚蹬草鞋，鞋底磨穿了一半。。体型：身高约6头身，中年男子身形（32岁却显老），微佝偻，清瘦。。衣物细节：粗布褐衣带补丁，草鞋磨穿一半。。发型妆造：发髻木簪别着，鬓角几缕白发。。脸型五官：国字脸，眉间刻着风霜，眼神沉静，鼻梁挺，薄唇抿紧，下巴方正，鬓角斑白。。武器招式：无兵器，以指为笔在地上演卦。。功法：无功法，不习武。唯有一肚子书——自幼熟读《周易》《尚书》，能演八卦，会推星象，但无人信他。。功法表现：无神力，唯手边书卷与卦纹。。画面：构图：破旧城墙根下，一布衣男子蹲坐，以指画地。背景是灰黄的市井街巷，远处有卖菜挑担的虚影。色调：土褐+灰黄，晨光从左侧斜照，在地面拉出长长的影子。氛围：市井、朴实、隐忍。。台词："圣人云'五十知天命'，我今年三十二，却连命在哪里都不知道。也罢——卖我的面，读我的书，天若不弃我，终有我一展抱负之日。"。动作帧（动图）：①拂晓挑面担穿街 ②收摊蹲墙角以指画地 ③演算卦理 ④见乞丐分半个饼。诗词：朝歌市井卖面郎，日暮城墙演卦章。人笑痴人说梦话，不知胸中有汪洋。。主题句：三十二岁不知命，七十二岁知天命，一百岁后——命就是我，我就是命。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：土褐+灰黄，晨光从左侧斜照，在地面拉出长长的影子。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 朝歌市井, age 32岁, 朝歌卖面·市井隐志; scene: 拂晓即起，挑着面粉担子穿街过巷；收摊后蹲在墙角，以指为笔在地上演算卦理；见到乞丐会分半个饼。; 朝歌市井卖面郎，日暮城墙演卦章。人笑痴人说梦话，不知胸中有汪洋。; palette: 土褐+灰黄，晨光从左侧斜照，在地面拉出长长的影子; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 垂钓渭水**
- 姜子牙，凡尘砺心阶段·昆仑拜师（35岁，昆仑修道·挑水砍柴）。形象：姜子牙，白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。品性：三十二岁上昆仑，面见元始天尊时双膝跪地，不哭不求，只说了八个字："弟子愿学，望师收录。"从此每日砍柴挑水，整整三年不曾闻道，却无一句怨言。其心之诚，令玉虚宫众弟子动容。。姿态：清晨第一个起床扫地；扫地时扫帚不停，口中默念《道德经》；劈柴时不看柴看纹路——他说"木有木的纹理，跟卦象一个道理"；独自坐于后山石上，面朝东海，静看日出。。服饰：玉虚宫青布道袍，布巾束发，腰间系一条麻绳。眉目间已无朝歌时的愁苦，代之以沉静的光。。体型：身高约6头身，青袍道者，背脊渐挺，沉静。。衣物细节：玉虚宫青布道袍，腰系麻绳。。发型妆造：布巾束发，眉目沉静。。脸型五官：国字脸，眉目舒展，眼神内敛有光，鼻梁挺，嘴角平和，下巴方正，静穆。。武器招式：无兵器，执扫帚/柴刀。。功法：始修吐纳导引术（元始天尊所授基础功课）；粗通道法基础（观星、望气、掐指演算）；三年砍柴挑水已磨出沉稳心性。。功法表现：吐纳导引，晨光笼身。。画面：构图：昆仑山玉虚宫后山，一青袍男子面东而坐，身前云海翻涌，身后古松苍劲。色调：青灰+墨绿+云白，晨光照亮其背影轮廓。氛围：静修、沉淀、决心。。台词："师父曾说'大器晚成'。弟子在朝歌蹉跎三十二年，今日方知——那三十二年不是蹉跎，是等待。弟子不惧晚，只怕不成。"。动作帧（动图）：①清晨扫地 ②扫帚不停默念道德经 ③劈柴看纹路 ④坐后山石上面朝东海看日出。诗词：三十二岁上昆仑，三年挑水未闻真。心中一片玉壶月，不惧光阴不惧尘。。主题句：三十二岁不知命，七十二岁知天命，一百岁后——命就是我，我就是命。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青灰+墨绿+云白，晨光照亮其背影轮廓。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 昆仑拜师, age 35岁, 昆仑修道·挑水砍柴; scene: 清晨第一个起床扫地；扫地时扫帚不停，口中默念《道德经》；劈柴时不看柴看纹路——他说"木有木的纹理，跟卦象一个道理"；独自坐于后山石上，面朝东海，静看日出。; 三十二岁上昆仑，三年挑水未闻真。心中一片玉壶月，不惧光阴不惧尘。; palette: 青灰+墨绿+云白，晨光照亮其背影轮廓; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 悟道下山**
- 姜子牙，道法初成阶段·渭水垂钓（72岁，渭水垂钓·愿者上钩）。形象：姜子牙，白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。品性：七十二岁下山，元始天尊授其封神榜、杏黄旗、打神鞭。下山后不急于求用，径往渭水之滨，以直钩垂钓。旁人笑他"直钩钓鱼，愿者上钩"，他只微微一笑——等的不是鱼，是文王。。姿态：坐于渭水边青石上，手持竹竿，竿头直钩离水三尺；眼不看水面，看的是天际流云；口中时而低吟《周易》卦辞；身侧放着一卷竹简，写满了三十六路兵法。。服饰：白衣鹤氅，腰系黄绦，头戴逍遥巾。与昆仑时不同——不再是弟子打扮，已有"师者"气度。竹竿为杖，三尺直钩系着红绳。。体型：身高约6头身，白衣鹤氅老者（72岁），仙风道骨。。衣物细节：白衣鹤氅，腰系黄绦，头戴逍遥巾，直钩系红绳。。发型妆造：白发白须，逍遥巾束发。。脸型五官：国字脸，长眉入鬓，眼神悠远，鼻梁挺，嘴角含笑，白须长垂，下巴方正。。武器招式：竹竿为杖；封神榜/杏黄旗/打神鞭悬身侧。。功法：封神榜（元始天尊所授，三百六十五路正神待封）；杏黄旗（一展可定乾坤）；打神鞭（专打不正之神）；六韬兵法（下山前日夜精研，已烂熟于心）。。功法表现：杏黄旗明黄为视觉中心，封神榜黄绸包裹。。画面：构图：渭水边青石上，一白衣老者持竿垂钓，直钩悬于水面之上。背景是苍茫山川与长空流云。色调：青绿+灰白+淡金（夕阳），直钩上一点红绳如朱砂。氛围：等待、从容、胸有成竹。。台词："直钩钓鱼，愿者上钩。我不钓文王——文王会自己来。若他来了，我便以这身本事，还他一个八百年的天下。"。动作帧（动图）：①坐渭水青石 ②手持竹竿直钩悬水三尺 ③眼望天际流云 ④低吟周易卦辞。诗词：渭水垂纶直钩悬，离水三尺待流年。旁人笑我痴人梦，我笑旁人看不穿。。主题句：三十二岁不知命，七十二岁知天命，一百岁后——命就是我，我就是命。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青绿+灰白+淡金（夕阳），直钩上一点红绳如朱砂。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 渭水垂钓, age 72岁, 渭水垂钓·愿者上钩; scene: 坐于渭水边青石上，手持竹竿，竿头直钩离水三尺；眼不看水面，看的是天际流云；口中时而低吟《周易》卦辞；身侧放着一卷竹简，写满了三十六路兵法。; 渭水垂纶直钩悬，离水三尺待流年。旁人笑我痴人梦，我笑旁人看不穿。; palette: 青绿+灰白+淡金（夕阳），直钩上一点红绳如朱砂; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 杏黄旗展**
- 姜子牙，大劫淬炼阶段·封神大战（80-90岁，西岐拜相·封神大战）。形象：姜子牙，白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。品性：西岐拜相，统率三军。百战之中从不亲临前线，却算无遗策——对手的每一步都在他的卦里。但他不是冷酷之人，每有将士阵亡，必亲自祭奠。面对截教众仙，他始终记得师父的话："封神榜上三百六十五人，不杀一人，如何封神？"。姿态：阵前一手持杏黄旗，一手掐指演算；遇强敌时后退三步，打神鞭一指；战后独自登台，望天良久，在竹简上记下阵亡将士的名字。。服饰：金甲战袍（西岐所赐），肩披玄色披风，腰悬打神鞭，背后负着封神榜（以黄绸包裹），杏黄旗插于阵前。。体型：身高约6头身，金甲战袍，身姿威重。。衣物细节：金甲战袍，玄色披风，腰悬打神鞭，背负封神榜（黄绸包裹）。。发型妆造：金盔束发，白发白须。。脸型五官：国字脸，浓眉，眼神如炬，鼻梁挺，唇线坚毅，下巴方正，金盔下白发。。武器招式：打神鞭通灵可大可小；杏黄旗展收定风雨。。功法：打神鞭已通灵（可大可小，专打仙人）；杏黄旗可展可收（一展定乾坤，一收藏风雨）；封神榜已展开三分之一（随着战事推进，榜上名字逐一亮起）；六韬兵法化为实战。。功法表现：封神榜展开如云，榜上名字逐一亮起。。画面：构图：阵前高台之上，姜子牙一手持杏黄旗，一手掐诀，身后封神榜展开如云。背景是战火与雷电交加的天际。色调：暗红+金+墨黑，杏黄旗的明黄是画面的视觉中心。氛围：战争、决断、天命的执行者。。台词："师父说'杀一人易，封一神难'。今日我杀一人，明日他上封神榜——百年之后，他是天神，我是凡人。可我不悔。这天下的秩序，总要有人来定。"。动作帧（动图）：①一手持杏黄旗 ②一手掐指演算 ③遇强敌退三步打神鞭一指 ④战后登台记阵亡将士名。诗词：杏黄旗展定乾坤，打神鞭指众仙魂。封神榜上三百六，一笔一划血写成。。主题句：三十二岁不知命，七十二岁知天命，一百岁后——命就是我，我就是命。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：暗红+金+墨黑，杏黄旗的明黄是画面的视觉中心。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 封神大战, age 80-90岁, 西岐拜相·封神大战; scene: 阵前一手持杏黄旗，一手掐指演算；遇强敌时后退三步，打神鞭一指；战后独自登台，望天良久，在竹简上记下阵亡将士的名字。; 杏黄旗展定乾坤，打神鞭指众仙魂。封神榜上三百六，一笔一划血写成。; palette: 暗红+金+墨黑，杏黄旗的明黄是画面的视觉中心; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 封神榜启**
- 姜子牙，封神登天阶段·封神台（93岁，封神台·册封诸神）。形象：姜子牙，白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。品性：牧野之战后，登封神台册封三百六十五位正神。每念一个名字，便想起那个人的生平——有的曾是他的对手，有的曾是他的战友。念完最后一个名字，他合上封神榜，回头看了一眼人间，轻轻叹了口气："都走了。"。姿态：登台时每一步都走得极稳；宣读封神榜时声音不高，却字字清晰，天地可闻；封神完毕，合上榜单，对着空无一人的封神台说了一句："各位，好走。"。服饰：紫金道袍，头戴莲花冠，手持玉笏，封神榜已合拢负于背后。面如冠玉，白发白须，已是标准的"神仙相"。。体型：身高约7头身，紫金道袍，白发白须神仙相。。衣物细节：紫金道袍，莲花冠，手持玉笏，封神榜合拢负于背后。。发型妆造：白发白须，面如冠玉。。脸型五官：国字脸，白眉，眼神庄严，鼻梁挺，口含玉笏之气，下巴方正，面如冠玉。。武器招式：玉笏在手；封神榜/打神鞭/杏黄旗皆收敛。。功法：封神榜已合拢（三百六十五神皆已归位）；打神鞭已收敛锋芒；杏黄旗已收起；道法大成（已可夜观天象而知三界事）。。功法表现：台下云雾三百六十五道神光依次亮起。。画面：构图：封神台上，姜子牙手持榜单宣读，台下云雾中有三百六十五道神光依次亮起。色调：紫金+明黄+圣白，神光如繁星点点。氛围：庄严、圆满、告别的时刻。。台词："三百六十五位正神，今日归位。姜子牙一生，只做了一件事——就是这件事。够了。"。动作帧（动图）：①登封神台 ②宣读封神榜 ③合上榜单 ④回头望人间轻叹"都走了"。诗词：封神台上榜文宣，三百六十五路仙。念罢诸神归位去，回看人间已千年。。主题句：三十二岁不知命，七十二岁知天命，一百岁后——命就是我，我就是命。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金+明黄+圣白，神光如繁星点点。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 封神台, age 93岁, 封神台·册封诸神; scene: 登台时每一步都走得极稳；宣读封神榜时声音不高，却字字清晰，天地可闻；封神完毕，合上榜单，对着空无一人的封神台说了一句："各位，好走。"; 封神台上榜文宣，三百六十五路仙。念罢诸神归位去，回看人间已千年。; palette: 紫金+明黄+圣白，神光如繁星点点; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 姜太公临**
- 姜子牙，道果圆满阶段·渭水归去（100+岁，渭水归去·无钩垂钓）。形象：姜子牙，白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。品性：封神之后，周武王拜他为"师尚父"，他却辞官不受，只身回到渭水边。还是那块青石，还是那根竹竿——只是直钩变成了无钩。他坐在水边，什么都不等，什么都不钓。有人问他"太公何不留在朝中享福？"他笑了笑："我在渭水边，比在朝堂上舒服。"。姿态：独坐渭水青石上，竹竿横膝，不垂不钓；偶尔低头看看水中的倒影——白发苍苍的老头儿，和当年卖面粉的年轻人，好像是同一个人，又好像不是；闭上眼，风从耳边过，像三十岁那年朝歌城头的风。。服饰：布衣鹤氅，已无任何法器随身。白发披散，未束未冠。脚下是草鞋，手边是竹杖。与凡间老者无异，唯有一双眼睛——清澈如初上山时。。体型：身高约6头身，布衣老者独坐，佝偻却清朗。。衣物细节：布衣鹤氅，无任何法器随身，草鞋竹杖。。发型妆造：白发披散，未束未冠。。脸型五官：国字脸，眉淡目清，眼神澄澈如初，鼻梁挺，嘴角含一抹淡笑，下巴圆润，清癯。。武器招式：无兵器，竹竿横膝。。功法：所有法器已还归玉虚宫（封神榜归天、打神鞭归位、杏黄旗归宗）；已无法力外显，唯有天眼通（可观三界事，但不看）；道心圆满（已无需修行，行住坐卧皆是道）。。功法表现：无法力外显，唯余天眼通（可观不视）。。画面：构图：渭水边青石上，一白发老者独坐，竹竿横膝，水波不兴。背景是远山与暮色，画面空阔。色调：水墨灰白+淡青+夕阳暖黄，水的颜色几乎与天融为一体。氛围：回归、平静、无我——一蓑烟雨任平生。。台词："三十二岁上昆仑，七十二岁下山，九十三岁封神。活了一百多岁——其实我最喜欢的，还是三十二岁那年，在朝歌城墙上演卦的自己。那时候什么都不会，可那会儿的风，真自由啊。"。动作帧（动图）：①独坐渭水青石 ②竹竿横膝不垂不钓 ③低头看水中倒影 ④闭眼听风。诗词：封神事了返渭滨，青石竹竿不钓鳞。白发老翁归何处？朝歌城里卖面人。。主题句：三十二岁不知命，七十二岁知天命，一百岁后——命就是我，我就是命。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：水墨灰白+淡青+夕阳暖黄，水的颜色几乎与天融为一体。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 渭水归去, age 100+岁, 渭水归去·无钩垂钓; scene: 独坐渭水青石上，竹竿横膝，不垂不钓；偶尔低头看看水中的倒影——白发苍苍的老头儿，和当年卖面粉的年轻人，好像是同一个人，又好像不是；闭上眼，风从耳边过，像三十岁那年朝歌城头的风。; 封神事了返渭滨，青石竹竿不钓鳞。白发老翁归何处？朝歌城里卖面人。; palette: 水墨灰白+淡青+夕阳暖黄，水的颜色几乎与天融为一体; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 哪吒（`nezha`） · 人生档案版

**灵胎初醒 · 灵珠种子**
- 哪吒，灵胎初醒阶段·灵珠降世（出生，陈塘关李府·灵珠子转世）。形象：哪吒，莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。品性：陈塘关总兵李靖第三子，灵珠子转世，生来就是一团烈火般的少年。。姿态：三岁随兄在校场玩耍，夺枪便耍，天生神力，枪法自通。。服饰：红衣肚兜，光着脚丫，颈戴乾坤圈，腕系混天绫。。体型：身高约3头身，奶胖幼童，圆润白嫩，红衣赤足。。衣物细节：红衣肚兜朱红，颈戴乾坤圈，腕系混天绫红绸。。发型妆造：总角双髻，刘海齐眉，眉心一点灵珠印记。。脸型五官：娃娃圆脸，眉弯如月，大眼圆睁，鼻梁小巧，小嘴粉润，下巴圆，眉心一点灵珠印。。武器招式：夺来的木枪随手耍；乾坤圈可掷可收。。功法：天生神力，灵珠之体；乾坤圈可掷可收，混天绫可缠可缚。。功法表现：乾坤圈绕身金光；混天绫红绸飘飞。。画面：构图：陈塘关李府后园，红衣赤足幼童持枪而舞，乾坤圈绕身，混天绫飘红，背景城墙海天。色调：朱红+石青+暖金，童真明朗。氛围：少年意气、天不怕地不怕。。台词："我才不学那些规矩。我是灵珠子转世，生来就要闹他个天翻地覆。"。动作帧（动图）：①校场夺枪 ②天生神力舞枪 ③乾坤圈绕腕转 ④咧嘴得意叉腰。诗词：灵珠一点降陈塘，红衣赤足小儿郎。乾坤圈转混天绫，天生便要闹一场。。主题句：剔骨还父、割肉还母，莲花重生的少年——一生快意，从不低头。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱红+石青+暖金，童真明朗。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 灵珠降世, age 出生, 陈塘关李府·灵珠子转世; scene: 三岁随兄在校场玩耍，夺枪便耍，天生神力，枪法自通。; 灵珠一点降陈塘，红衣赤足小儿郎。乾坤圈转混天绫，天生便要闹一场。; palette: 朱红+石青+暖金，童真明朗; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 陈塘哪吒**
- 哪吒，凡尘砺心阶段·闹海抽龙（7岁，东海玩水·抽龙筋惹祸）。形象：哪吒，莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。品性：盛夏去东海边玩水，见海涛汹汹，童心大起，跃入海中搅动乾坤圈，搅得东海龙宫晃动。。姿态：龙王三太子敖丙出水喝问，他反手一记乾坤圈，抽了龙筋缠在手上。。服饰：红衣半湿，混天绫缠腕，乾坤圈在手，脚踩浪头。。体型：身高约3头身，红衣湿透，踏浪而立。。衣物细节：红衣半湿，混天绫缠腕，乾坤圈在手。。发型妆造：总角沾水，湿发贴额，眉间灵珠光。。脸型五官：娃娃圆脸，眉微扬，大眼明亮，鼻尖微翘，咧嘴笑，下巴圆，湿发贴额。。武器招式：乾坤圈一击碎甲；混天绫翻海。。功法：乾坤圈一击碎甲；混天绫翻海覆浪，搅动东海。。功法表现：搅动海波如山；抽得龙筋金光一条。。画面：构图：东海碧波之上，红衣幼童踏浪而立，手中拎着长长的龙筋，海水翻腾，背景龙宫隐约。色调：海蓝+朱红+白浪。氛围：闯祸的嚣张、少年无惧。。台词："你就是龙王三太子？正好，借你的龙筋给我系混天绫！"。动作帧（动图）：①跃入东海 ②搅动乾坤圈 ③龙王三太子出水 ④一记乾坤圈抽龙筋。诗词：东海嬉水起风波，乾坤一掷碎鳞多。龙王三子俯首去，抽得龙筋缚天河。。主题句：剔骨还父、割肉还母，莲花重生的少年——一生快意，从不低头。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：海蓝+朱红+白浪。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 闹海抽龙, age 7岁, 东海玩水·抽龙筋惹祸; scene: 龙王三太子敖丙出水喝问，他反手一记乾坤圈，抽了龙筋缠在手上。; 东海嬉水起风波，乾坤一掷碎鳞多。龙王三子俯首去，抽得龙筋缚天河。; palette: 海蓝+朱红+白浪; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 闹海屠龙**
- 哪吒，道法初成阶段·剔骨还父（7岁，自刎谢罪·剔骨还父割肉还母）。形象：哪吒，莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。品性：龙王水淹陈塘关，逼李靖交出儿子。他不愿连累父母与全城百姓，提剑自刎——剔骨还父，割肉还母。。姿态：东海水淹城关，他登城头，仰天大笑，挥剑自刎；魂魄飘至师父太乙真人座前。。服饰：单衣立于城头，乾坤圈混天绫尽数卸下，赤手空拳。。体型：身高约3头身，单衣立于城头，身形单薄却挺直。。衣物细节：单衣素白，乾坤圈混天绫尽卸。。发型妆造：散发，面无血色。。脸型五官：娃娃脸煞白，眉紧蹙，大眼含泪又坚定，鼻梁挺，唇紧抿，下巴微收。。武器招式：无兵器，以剑自刎谢罪。。功法：以一己之命担一城之责；魂魄不灭，得师父接引。。功法表现：魂魄离体化作金光，飘向师父莲池。。画面：构图：陈塘关城头，少年单衣立于雨中，海潮滔天，身后是颤抖的父母，他挥剑自刎的一瞬，漫天水光。色调：墨青+血绛+惨白，悲壮肃穆。氛围：大悲、决绝、以一命抵全城。。台词："爹爹娘亲，这身骨肉是你们给的，今日一并还了。若有来世，我哪吒，仍是你们儿子！"。动作帧（动图）：①登城头 ②望海潮大笑 ③提剑 ④挥剑自刎，漫天水光。诗词：四海龙王逼城关，为全忠孝自刎还。剔骨还父血还母，莲花座下再登天。。主题句：剔骨还父、割肉还母，莲花重生的少年——一生快意，从不低头。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨青+血绛+惨白，悲壮肃穆。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 剔骨还父, age 7岁, 自刎谢罪·剔骨还父割肉还母; scene: 东海水淹城关，他登城头，仰天大笑，挥剑自刎；魂魄飘至师父太乙真人座前。; 四海龙王逼城关，为全忠孝自刎还。剔骨还父血还母，莲花座下再登天。; palette: 墨青+血绛+惨白，悲壮肃穆; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 莲花化身**
- 哪吒，大劫淬炼阶段·莲花化身（7岁，太乙真人·莲池重生）。形象：哪吒，莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。品性：太乙真人惜其魂魄，以莲池中莲花莲叶为骨肉，金丹为魂，重塑肉身——从此不染凡尘，百邪不侵。。姿态：在师父莲池中醒来，见自己立于莲心，身轻如燕；试握火尖枪，一枪刺出莲花火焰冲天。。服饰：莲花化身的金肌玉骨，红肚兜外罩金缕衣，脚下踏出两团风火轮。。体型：身高约3头身，莲花化身的金肌玉骨，身轻如燕。。衣物细节：红肚兜外罩金缕衣，脚下风火轮，火尖枪/乾坤圈/混天绫/风火轮四宝齐备。。发型妆造：束发金环，眉目英气。。脸型五官：莲花化身娃娃脸，眉目英气，大眼有神，鼻梁挺直，嘴角含笑，下巴圆润，金光面庞。。武器招式：火尖枪一枪燎原；风火轮踏焰。。功法：莲花化身（不坏之身）；三头六臂初现；火尖枪、风火轮、乾坤圈、混天绫四宝齐备。。功法表现：莲花圣光护体；枪尖莲花火焰。。画面：构图：太乙真人莲池之中，少年立于巨大金莲之上，莲花圣光环绕，火尖枪/风火轮/乾坤圈/混天绫四宝飞旋，背景云海仙山。色调：莲粉+鎏金+圣白，重生夺目。氛围：脱胎换骨、神光初绽。。台词："师父，这具莲花身子，比原来的还趁手！从今往后，我哪吒没了凡骨，只有一颗赤子之心。"。动作帧（动图）：①莲心醒来 ②起身立于莲上 ③试握火尖枪 ④一枪刺出莲花火冲天。诗词：太乙莲池起金光，莲花化骨魄还阳。三头六臂初开眼，火尖枪挑日月长。。主题句：剔骨还父、割肉还母，莲花重生的少年——一生快意，从不低头。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：莲粉+鎏金+圣白，重生夺目。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 莲花化身, age 7岁, 太乙真人·莲池重生; scene: 在师父莲池中醒来，见自己立于莲心，身轻如燕；试握火尖枪，一枪刺出莲花火焰冲天。; 太乙莲池起金光，莲花化骨魄还阳。三头六臂初开眼，火尖枪挑日月长。; palette: 莲粉+鎏金+圣白，重生夺目; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 三太子**
- 哪吒，封神登天阶段·封神先锋（少年，西岐先锋·三头六臂大战）。形象：哪吒，莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。品性：随姜子牙伐纣，为西岐先锋。三头六臂全开，风火轮踏遍沙场，一人当千，战无不胜。。姿态：阵前现出三头六臂，火尖枪横扫，乾坤圈掷出索命，混天绫卷起千重红浪；战后仍是那副少年得意模样。。服饰：金甲红绫，三头六臂，各持火尖枪/乾坤圈/混天绫/风火轮，莲花圣光护体。。体型：身高约4头身，三头六臂全开，金甲少年战神。。衣物细节：金甲红绫，三头六臂各持火尖枪/乾坤圈/混天绫/风火轮。。发型妆造：束发金冠，眉间灵珠光盛。。脸型五官：少年娃娃脸，剑眉，大眼凌厉，鼻梁挺，咧嘴笑，下巴微抬，灵珠光盛。。武器招式：火尖枪挑、乾坤圈掷、混天绫卷、风火轮驰。。功法：三头六臂全开；火尖枪挑、乾坤圈掷、混天绫卷、风火轮驰；莲花化身百邪不侵。。功法表现：莲花圣光护体；乾坤圈划破长空。。画面：构图：沙场之上，少年三头六臂立于风火轮上，火尖枪斜指，乾坤圈掷出划破长空，身后西岐旌旗，背景战火连天。色调：朱红+鎏金+墨黑，热血炽烈。氛围：少年战神、一人成军。。台词："姜丞相，我打架不需要排兵布阵——我一个人，就是一支军队！"。动作帧（动图）：①现三头六臂 ②火尖枪横扫 ③乾坤圈掷出 ④混天绫卷千重红浪。诗词：西岐帐下少年郎，三头六臂破八荒。风火轮上乾坤定，莲花光里斩天狼。。主题句：剔骨还父、割肉还母，莲花重生的少年——一生快意，从不低头。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱红+鎏金+墨黑，热血炽烈。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 封神先锋, age 少年, 西岐先锋·三头六臂大战; scene: 阵前现出三头六臂，火尖枪横扫，乾坤圈掷出索命，混天绫卷起千重红浪；战后仍是那副少年得意模样。; 西岐帐下少年郎，三头六臂破八荒。风火轮上乾坤定，莲花光里斩天狼。; palette: 朱红+鎏金+墨黑，热血炽烈; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 肉身成圣**
- 哪吒，道果圆满阶段·三坛海会（封神后，中坛元帅·莲心护世）。形象：哪吒，莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。品性：封神后受封中坛元帅、三坛海会大神，镇守中坛，护佑三界。三头六臂敛去，仍是那个红衣赤足、天不怕地不怕的少年——只是多了几分护佑众生的慈。。姿态：巡视三界，护佑孩童与忠良；见人间孩童夜啼，会悄悄在窗边放一朵莲花；依旧会为不平事冲冠一怒。。服饰：红绫金甲，莲花圣光内敛；三头六臂收作一身，乾坤圈混天绫仍随身。。体型：身高约4头身，负手云端，三头六臂敛作一身，仍少年。。衣物细节：红绫金甲，莲花圣光内敛，乾坤圈混天绫随身。。发型妆造：束发，眉目不改少年。。脸型五官：少年娃娃脸，眉目含慈，大眼清澈，鼻梁挺，嘴角带笑，下巴圆润，不改稚气。。武器招式：乾坤圈混天绫静悬身侧，一念即出。。功法：三坛海会大神神位；莲花化身不坏；三头六臂一念即出；护佑三界孩童。。功法表现：莲花圣光化作六瓣金莲于身后。。画面：构图：云端之上，红衣金甲少年负手而立，身后莲花圣光化作文六瓣金莲，脚下祥云，俯瞰人间城郭与万家灯火，身侧乾坤圈混天绫静静悬浮。色调：朱红+鎏金+云白+莲粉，庄重而不失少年气。氛围：护世、赤诚、终成中坛元帅。。台词："做了神，我还是那个陈塘关的哪吒。三界不平，我就要管；孩子受委屈，我就要护。"。动作帧（动图）：①云端负手 ②巡视三界 ③见孩童夜啼窗边放莲花 ④为不平事冲冠一怒。诗词：三坛海会镇中坛，莲心一颗照人间。历经剔骨重生后，赤子之心更少年。。主题句：剔骨还父、割肉还母，莲花重生的少年——一生快意，从不低头。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱红+鎏金+云白+莲粉，庄重而不失少年气。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 三坛海会, age 封神后, 中坛元帅·莲心护世; scene: 巡视三界，护佑孩童与忠良；见人间孩童夜啼，会悄悄在窗边放一朵莲花；依旧会为不平事冲冠一怒。; 三坛海会镇中坛，莲心一颗照人间。历经剔骨重生后，赤子之心更少年。; palette: 朱红+鎏金+云白+莲粉，庄重而不失少年气; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 杨戬（`yang_jian`） · 人生档案版

**灵胎初醒 · 天眼灵光**
- 杨戬，灵胎初醒阶段·桃山遗孤（12岁，母被压山·天眼初开）。形象：杨戬，银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。品性：瑶姬仙子之子，母亲因私婚凡间被天庭压于桃山之下。自幼流落人间，沉默寡言、心思极重，却刚毅不折。凡人不识他，只当他是个孤僻的少年。。姿态：常独坐高坡望桃山方向发呆；以打猎为生；夜里常梦见母亲在山底唤他名字，醒来枕巾尽湿。。服饰：粗布短褐，眉目清秀，额间有一道淡淡的天眼印记，时隐时现。。体型：身高约5头身，少年身形，清瘦，眉目清秀。。衣物细节：粗布短褐，无饰。。发型妆造：乱发束起，额间淡淡天眼印记。。脸型五官：鹅蛋脸，剑眉微蹙，凤目坚毅，鼻梁高挺，薄唇抿紧，下巴微尖，额间天眼印记。。武器招式：无兵器，以打猎为生（弓/石）。。功法：无师自通的天眼初现（能看穿山川，隐隐感知母亲所在）；体质非凡（天神后裔，力大目明）。。功法表现：天眼初现，透出一点金光。。画面：构图：桃山远景巍峨如囚牢，少年杨戬立于山脚仰头望山，额间天眼透出一点微光。背景苍茫旷野。色调：灰褐+青灰+一点天眼金光。氛围：孤苦、刚毅、望母。。台词："娘被压在山下，我要把她救出来。他们都笑我痴，可我额头上这眼睛，是真看得见那山的。"。动作帧（动图）：①独坐高坡望桃山 ②发呆良久 ③握紧拳头 ④额间天眼微光一闪。诗词：桃山之下母被囚，襁褓流离少年愁。天眼初睁窥万里，誓劈此山救娘亲。。主题句：先是我杨戬，才是二郎神。一柄刀，一口天眼，护一人间烟火。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：灰褐+青灰+一点天眼金光。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 桃山遗孤, age 12岁, 母被压山·天眼初开; scene: 常独坐高坡望桃山方向发呆；以打猎为生；夜里常梦见母亲在山底唤他名字，醒来枕巾尽湿。; 桃山之下母被囚，襁褓流离少年愁。天眼初睁窥万里，誓劈此山救娘亲。; palette: 灰褐+青灰+一点天眼金光; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年二郎**
- 杨戬，凡尘砺心阶段·拜师玉鼎（15岁，玉泉学艺·八九玄功）。形象：杨戬，银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。品性：经玉泉山金霞洞玉鼎真人收为弟子。天资卓绝却最肯下苦功，学八九玄功时日夜不辍。真人赞他"此子非池中物，将来必成大器"。。姿态：每日黎明即起，对着玉泉山瀑布演练刀法；课余练八九玄功（变化之术）至深夜；闲暇时逗弄师父养的白犬——后来成了他的哮天犬。。服饰：玉泉山青布道袍，束发利落，额间天眼渐睁，已能看清妖邪本相。。体型：身高约5头身，少年道人，身形渐健。。衣物细节：玉泉山青布道袍，束发利落。。发型妆造：束发利落，额间天眼渐睁。。脸型五官：鹅蛋脸，剑眉舒展，凤目专注，鼻梁高挺，嘴角沉静，下巴微尖，天眼渐睁。。武器招式：三尖两刃刀初学，对瀑劈砍。。功法：始修八九玄功（七十二变）；初学三尖两刃刀；天眼可辨妖邪，观星望气。。功法表现：天眼可辨妖邪；刀风带起水雾。。画面：构图：玉泉山瀑布前，少年道人杨戬舞三尖两刃刀，水花飞溅，一只白犬蹲坐瀑布边望着他。背景青山飞瀑。色调：青绿+墨青+云白。氛围：苦修、专注、与犬为伴。。台词："师父说九转玄功练成可移山倒海，我且练它三年，换娘亲出山。白犬，你也好好长大，将来陪我下山。"。动作帧（动图）：①黎明对瀑布练刀 ②刀光水花飞溅 ③白犬蹲坐旁观 ④深夜练八九玄功。诗词：玉泉山上拜师门，八九玄功昼夜勤。瀑前练得三尖刃，为救慈亲不惜身。。主题句：先是我杨戬，才是二郎神。一柄刀，一口天眼，护一人间烟火。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青绿+墨青+云白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 拜师玉鼎, age 15岁, 玉泉学艺·八九玄功; scene: 每日黎明即起，对着玉泉山瀑布演练刀法；课余练八九玄功（变化之术）至深夜；闲暇时逗弄师父养的白犬——后来成了他的哮天犬。; 玉泉山上拜师门，八九玄功昼夜勤。瀑前练得三尖刃，为救慈亲不惜身。; palette: 青绿+墨青+云白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 八九玄功**
- 杨戬，道法初成阶段·劈山救母（18岁，一刀劈山·大孝动天）。形象：杨戬，银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。品性：八九玄功大成，握开山之能。不拜玉帝（其舅），不认天条，只认一个"孝"字。桃山前长跪三日，第四日挥刃劈山，救出母亲。。姿态：手持三尖两刃刀，一刀裂开桃山；母亲得见天日；面对天庭追责，他选择——自刎还母（后被救，玉帝终与他和好，封他二郎神）。。服饰：银甲初着，三尖两刃刀在手，额间天眼怒张放出神光。。体型：身高约6头身，银甲初着，身形英挺。。衣物细节：银甲初着，三尖两刃刀在手。。发型妆造：束发，额间天眼怒张。。脸型五官：鹅蛋脸，剑眉怒竖，凤目含光，鼻梁高挺，唇线决绝，下巴坚毅，天眼怒张。。武器招式：三尖两刃刀开山裂石。。功法：八九玄功大成（移山倒海）；三尖两刃刀可开山裂石；天眼照彻三界。。功法表现：天眼放神光，山裂金光。。画面：构图：桃山从中裂开，金光迸射，杨戬银甲持刀立于山前，天眼放出神光，母亲身影在光中得见天日。背景天光破云。色调：墨褐+金光+血绛。氛围：决绝、赤诚、劈天救母。。台词："天条是舅舅定的，孝道是我娘的。我杨戬，只认孝，不认天。"。动作帧（动图）：①桃山前长跪 ②起身拔刀 ③一刀劈下 ④山裂金光迸射。诗词：桃山高万丈，天眼裂层云。一刀劈天地，救母出沉沦。。主题句：先是我杨戬，才是二郎神。一柄刀，一口天眼，护一人间烟火。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐+金光+血绛。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 劈山救母, age 18岁, 一刀劈山·大孝动天; scene: 手持三尖两刃刀，一刀裂开桃山；母亲得见天日；面对天庭追责，他选择——自刎还母（后被救，玉帝终与他和好，封他二郎神）。; 桃山高万丈，天眼裂层云。一刀劈天地，救母出沉沦。; palette: 墨褐+金光+血绛; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 天眼洞开**
- 杨戬，大劫淬炼阶段·封神大战（20余岁，西岐先锋·肉身成圣）。形象：杨戬，银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。品性：奉师命下山助西岐伐纣。身为先锋屡立战功，擒魔家四将、斗梅山七怪、劈金毛犼，战无不胜。但本性淡泊，不贪功名，只做该做之事。。姿态：阵前独骑，三尖两刃刀横扫千军；哮天犬如电随行，专咬妖邪；战后不居功，静立一旁；识破土行孙地行之术等无数奇谋。。服饰：银甲金冠，玄色披风，哮天犬相随，三尖两刃刀在手，额间天眼如电。。体型：身高约7头身，银甲金冠，战将英姿。。衣物细节：银甲金冠，玄色披风，哮天犬相随。。发型妆造：束发金冠，天眼如电。。脸型五官：鹅蛋脸，剑眉入鬓，凤目如电，鼻梁高挺，嘴角淡然，下巴英挺，天眼如电。。武器招式：三尖两刃刀横扫千军。。功法：八九玄功化身万千；天眼照破妖邪原形；哮天犬助阵；三尖两刃刀通神。。功法表现：天眼照破妖邪原形；哮天犬助阵。。画面：构图：沙场之上，杨戬银甲独骑，三尖两刃刀横扫，哮天犬如电随行，天眼照破漫天妖雾。背景战火连天。色调：暗红+银白+金+墨黑。氛围：战功赫赫、淡泊、自在。。台词："封神榜上有我名，也不过是排班轮值。我更愿意在这红尘里，做个自在的守护者。"。动作帧（动图）：①阵前独骑 ②三尖两刃刀横扫 ③哮天犬如电随行 ④战后静立不居功。诗词：三尖刃下妖邪伏，哮天犬随护圣躯。封神榜上名虽列，心向灌江不向天。。主题句：先是我杨戬，才是二郎神。一柄刀，一口天眼，护一人间烟火。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：暗红+银白+金+墨黑。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 封神大战, age 20余岁, 西岐先锋·肉身成圣; scene: 阵前独骑，三尖两刃刀横扫千军；哮天犬如电随行，专咬妖邪；战后不居功，静立一旁；识破土行孙地行之术等无数奇谋。; 三尖刃下妖邪伏，哮天犬随护圣躯。封神榜上名虽列，心向灌江不向天。; palette: 暗红+银白+金+墨黑; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 二郎真君**
- 杨戬，封神登天阶段·二郎真君（封神后，灌江口·听调不听宣）。形象：杨戬，银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。品性：封神后肉身成圣，受封灌江口二郎显圣真君。不居天庭（听调不听宣），在灌江口设庙护佑一方。梅山六圣、一千二百草头神相随，哮天犬不离左右。。姿态：端坐庙堂或巡游人间，护佑风调雨顺；哮天犬蹲伏案前；遇民间疾苦则显圣相助；闲时在江畔负手观水。。服饰：银甲金冠，玄色披风，腰间玉带，额间天眼如电（已大成）。。体型：身高约7头身，银甲端坐庙堂，威严。。衣物细节：银甲金冠，玄色披风，腰悬玉带。。发型妆造：束发金冠，天眼如电（大成）。。脸型五官：鹅蛋脸，剑眉，凤目威严，鼻梁高挺，唇微抿，下巴英挺，天眼含光。。武器招式：三尖两刃刀侧立。。功法：八九玄功圆满；天眼照三界；哮天犬、梅山六圣护法；法力深不可测。。功法表现：身后梅山六圣虚影，庙外香火袅袅。。画面：构图：灌江口二郎庙，杨戬银甲端坐主位，天眼含光，哮天犬伏于足侧，身后梅山六圣虚影浮现，庙外香火袅袅。色调：紫金+明黄+圣白。氛围：逍遥、守护、威严。。台词："做神有什么好？不过是换个地方守护。我在灌江口，护的是人间的烟火气。"。动作帧（动图）：①端坐主位 ②天眼含光 ③哮天犬伏足侧 ④巡视灌江口。诗词：灌江口上立真君，梅山六圣护清名。听调不听宣，人间烟火渡长生。。主题句：先是我杨戬，才是二郎神。一柄刀，一口天眼，护一人间烟火。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金+明黄+圣白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 二郎真君, age 封神后, 灌江口·听调不听宣; scene: 端坐庙堂或巡游人间，护佑风调雨顺；哮天犬蹲伏案前；遇民间疾苦则显圣相助；闲时在江畔负手观水。; 灌江口上立真君，梅山六圣护清名。听调不听宣，人间烟火渡长生。; palette: 紫金+明黄+圣白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 玄功大成**
- 杨戬，道果圆满阶段·清源妙道（千年后，灌江垂钓·本心不渝）。形象：杨戬，银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。品性：千百年守护灌江口，已是人间信仰。但本心仍如少年——喜欢在江边看云，喜欢听人间百姓的家长里短。有人说他是神，他只笑笑："我先是杨戬，才是二郎神。"。姿态：青石江畔负手而立，天眼半阖看云卷云舒；偶尔陪哮天犬在江边撒欢；遇人间孩童呼他"杨二哥"，会笑着递过一颗糖。。服饰：银甲已卸，布衣素服，唯有额间天眼隐隐含光。。体型：身高约7头身，布衣素服，负手而立，本真。。衣物细节：布衣素服，银甲已卸。。发型妆造：束发道髻，天眼微阖含光。。脸型五官：鹅蛋脸，剑眉舒展，凤目半阖温柔，鼻梁高挺，嘴角含笑，下巴柔和，天眼微阖。。武器招式：无兵器，负手而立。。功法：法力尽敛，唯余天眼（可观，但不轻用）；哮天犬相伴；一念可护一方。。功法表现：天眼微阖一点柔光；法力尽敛。。画面：构图：灌江畔青石上，布衣杨戬负手而立，天眼微阖，哮天犬卧于足边打盹。背景暮色江水、远山如黛。色调：水墨灰白+淡青+夕阳暖黄。氛围：逍遥、释然、本真。。台词："做了几千年的神，最怀念的还是拜师那年，在玉泉山对着瀑布劈刀的日子。那时候心里只有一个念头——救娘亲。现在娘亲早不在了，可那个念头，还亮着。"。动作帧（动图）：①青石江畔负手 ②天眼半阖看云卷云舒 ③陪哮天犬江边撒欢 ④见孩童递一颗糖。诗词：灌江碧水映斜阳，真君垂钓不钓王。千载功名都放下，只留天眼照故乡。。主题句：先是我杨戬，才是二郎神。一柄刀，一口天眼，护一人间烟火。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：水墨灰白+淡青+夕阳暖黄。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 清源妙道, age 千年后, 灌江垂钓·本心不渝; scene: 青石江畔负手而立，天眼半阖看云卷云舒；偶尔陪哮天犬在江边撒欢；遇人间孩童呼他"杨二哥"，会笑着递过一颗糖。; 灌江碧水映斜阳，真君垂钓不钓王。千载功名都放下，只留天眼照故乡。; palette: 水墨灰白+淡青+夕阳暖黄; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 雷震子（`lei_zhenzi`） · 人生档案版

**灵胎初醒 · 雷卵种子**
- 雷震子，灵胎初醒阶段·燕山雷婴（出生，燕山雷婴·雷雨降生）。形象：雷震子，鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。品性：生于燕山古洞的雷雨之中，凡胎一个，周身却缠绕未散的雷光。不哭不闹，睁着乌溜溜的眼睛看人，仿佛知道自己在这世间是"被拣来"的孩子。。姿态：蜷卧在古洞青石上，裹着粗布襁褓；文王循着啼哭寻来时，他伸出小手抓住文王的胡须，咧嘴笑了；洞外雷声渐远，雨丝斜打。。服饰：粗布襁褓裹得严实；额前一撮柔软的胎发；眉心一点淡淡的雷纹胎记，似有若无——这是老天爷盖的戳。。体型：身高约1头身（婴儿），蜷卧襁褓中，圆润软糯。。衣物细节：粗布襁褓米褐麻布，裹得严实，仅露小脸。。发型妆造：额前一撮柔软胎发，眉心一点淡金雷纹胎记。。脸型五官：婴儿圆脸，眉淡如线，眼圆黑亮，鼻小嘴嘟，下巴圆润，眉心一点淡金雷纹。。武器招式：无武器，婴儿手脚绵软。。功法：无。凡胎一个，不懂神通。只是出生那夜雷声震天，旁人说他"是雷公送到人间的孩子"。。功法表现：眉心雷纹淡金微闪，周身未散雷光如金线游走。。画面：构图：燕山古洞，青石上一裹粗布襁褓的婴儿蜷卧，眉心一点淡金雷纹；洞口外雨雾渐歇，文王俯身欲抱的背影。色调：大地灰褐+淡金微光，未散的雷光如金线游走。氛围：命定、孤苦、被接住。。台词："（文王将他抱进怀里）'雷雨之中降生的孩子……就叫你雷震子吧。'（婴儿却咧嘴笑了，仿佛真听懂了这个名字。）"。动作帧（动图）：①静卧襁褓 ②小手轻抬抓住文王胡须 ③咧嘴一笑 ④眼帘合拢入睡。诗词：燕山雷雨裂长空，古洞青石襁褓中。文王驻马抱将起，百子添来一稚童。。主题句：文王从雷雨里接回一个孩子，这个孩子用一生，把每一次雷雨都炼成了护人的翅膀。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：大地灰褐+淡金微光，未散的雷光如金线游走。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 燕山雷婴, age 出生, 燕山雷婴·雷雨降生; scene: 蜷卧在古洞青石上，裹着粗布襁褓；文王循着啼哭寻来时，他伸出小手抓住文王的胡须，咧嘴笑了；洞外雷声渐远，雨丝斜打。; 燕山雷雨裂长空，古洞青石襁褓中。文王驻马抱将起，百子添来一稚童。; palette: 大地灰褐+淡金微光，未散的雷光如金线游走; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年雷仔**
- 雷震子，凡尘砺心阶段·西岐王孙（少年，西岐王孙·文王百子）。形象：雷震子，鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。品性：在西岐长成少年。文王百子，他排行最小，最得父王疼爱。憨直，力大，心热。不知自己身世，只当自己是姬家第一百个孩子，因此格外恋家。。姿态：清晨在校场举石锁练力；午后溜去城郊追野兔，傍晚满载而归分给兄弟；夜里缠着父王讲"那年路过燕山"的故事——文王每讲到"雷雨"二字，他的眼睛就亮一下。。服饰：素色布衣劲装，头扎总角，眉目清朗；眉心雷纹胎记淡淡，似有若无。。体型：身高约4头身，少年圆润敦实，总角之龄。。衣物细节：素色布衣劲装青灰，腰间无饰，布鞋。。发型妆造：总角双髻，黑发油亮，眉心雷纹淡淡。。脸型五官：少年圆脸，浓眉大眼，鼻梁挺直，嘴角上扬，下巴圆润，眉心雷纹淡淡。。武器招式：无兵器，徒手举石锁练力。。功法：天生神力，能徒手搬起城门石墩；习武略有根底；云中子路过西岐远远望见，捋须叹道："此子不凡，他日当振翅九天。"。功法表现：天生神力，搬动石墩尘土微扬。。画面：构图：西岐城郊晨光，一总角少年布衣蹲身，双手抱起石锁；身后城门石狮为伴，远处宫苑檐角。色调：青灰+暖黄晨光，天光通透。氛围：健康、憨直、被爱着。。台词："（雷震子仰头）'父王，他们说我是雷雨里捡来的，当真？'（文王捋须）'是老天爷把你放进雷雨里，等我捡你回来。'（他挠头）'那我将来接好风雷，也去捡别人。'"。动作帧（动图）：①蹲身抱石锁 ②发力挺腰举起 ③放下石锁咧嘴笑 ④跑向城郊追野兔。诗词：王孙不羡锦衣光，校场举石力最强。夜听父王谈燕山，心驰雷雨向八方。。主题句：文王从雷雨里接回一个孩子，这个孩子用一生，把每一次雷雨都炼成了护人的翅膀。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青灰+暖黄晨光，天光通透。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 西岐王孙, age 少年, 西岐王孙·文王百子; scene: 清晨在校场举石锁练力；午后溜去城郊追野兔，傍晚满载而归分给兄弟；夜里缠着父王讲"那年路过燕山"的故事——文王每讲到"雷雨"二字，他的眼睛就亮一下。; 王孙不羡锦衣光，校场举石力最强。夜听父王谈燕山，心驰雷雨向八方。; palette: 青灰+暖黄晨光，天光通透; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 食杏生翅**
- 雷震子，道法初成阶段·玉柱变形（15岁，玉柱变形·仙杏生翅）。形象：雷震子，鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。品性：云中子收其为徒，带回终南山玉柱洞。本性馋而直，师父留两枚仙杏嘱"待我归来再食"，他偷吃一枚——正欲再食，忽觉浑身燥热，背后"噗"地生出一对风雷双翅，面皮化作蓝靛，发色朱红，口生獠牙。对水一照，嚎啕大哭："这副模样，还怎么去见父王！"。姿态：云中子归来抚掌而笑："天赐神通，何必自弃？"手把手教他控翅、御风、运雷；他惊惶地对着水面一遍遍练习振翅，从镜中那张蓝靛面孔里，渐渐看出一分英武来。。服饰：蓝靛面，朱砂发束成髻；青布道袍；背后双翅初展，左翅"风"纹、右翅"雷"纹金光隐现。。体型：身高约5头身，少年身形，背生双翅（初展较小）。。衣物细节：青布道袍，翅根风雷金圈，左翅"风"右翅"雷"金光隐现。。发型妆造：朱砂发束成髻，蓝靛面，金瞳雷纹，鸟喙獠牙初显。。脸型五官：鸟面初显，蓝靛面，朱砂眉，金瞳竖线，鸟喙微凸，下巴尖削，獠牙初露。。武器招式：黄金棍拄地，初得未精；风雷双翅初展。。功法：风雷双翅（日行万里、御雷而行）；黄金棍初得于云中子；初习风雷遁法。。功法表现：翅缘风雷金光流转；仙杏落地两枚金杏。。画面：构图：终南山玉柱洞前云海，一青袍少年背生双翅，翅缘风雷金光流转；他回望水面倒影，惊愕渐化为坚定；身后古松苍苍。色调：暖朱赤金+宝蓝翠绿，仙杏赤金与风雷青蓝交织。氛围：惊变、接纳、初得神通。。台词："（他指着自己的脸）'师父，这张脸还能变回去么？'（云中子摇头）'变不回。但你记好——你之异相，正是你之神通。丑与不丑，看心不看脸。'"。动作帧（动图）：①惊觉背上异动 ②双翅"噗"地展开 ③回望水面倒影惊愕 ④握黄金棍渐定。诗词：终南山上仙杏香，一翅生风一翅雷。对水照来惊变后，方知天赐此形骸。。主题句：文王从雷雨里接回一个孩子，这个孩子用一生，把每一次雷雨都炼成了护人的翅膀。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：暖朱赤金+宝蓝翠绿，仙杏赤金与风雷青蓝交织。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 玉柱变形, age 15岁, 玉柱变形·仙杏生翅; scene: 云中子归来抚掌而笑："天赐神通，何必自弃？"手把手教他控翅、御风、运雷；他惊惶地对着水面一遍遍练习振翅，从镜中那张蓝靛面孔里，渐渐看出一分英武来。; 终南山上仙杏香，一翅生风一翅雷。对水照来惊变后，方知天赐此形骸。; palette: 暖朱赤金+宝蓝翠绿，仙杏赤金与风雷青蓝交织; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 黄金棍**
- 雷震子，大劫淬炼阶段·潼关救主（18岁，潼关救主·雷翼护主）。形象：雷震子，鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。品性：闻父王被囚羕里七年、归途又遭截杀，连夜展翅北飞，昼夜不停。潼关之下追兵如潮，他一翅遮天，黄金棍横扫千军，负起父王飞越城关。落地时文王已白发苍苍，他单膝跪地："父王，儿子来接您回家了。"。姿态：双翅卷起风雷逼退追兵；黄金棍荡开箭雨；背负文王低飞，怕惊着老人特意放缓翅速；沿途以雷光开路。。服饰：蓝靛面，朱砂发披散，银蓝战袍；风雷双翅尽展，左"风"右"雷"金光灼灼；手持黄金棍。。体型：身高约6头身，青年身形，双翅全开展翅飞翔。。衣物细节：银蓝战袍，双翅全开"风""雷"金光灼灼，黄金棍在手。。发型妆造：朱砂发披散，蓝靛面，金瞳怒张。。脸型五官：鸟面雷公，蓝靛面，朱砂粗眉，金瞳怒张，鸟喙外凸，下巴尖削，獠牙显露。。武器招式：黄金棍斜指城下，荡开箭雨。。功法：风雷遁法大成（负人飞行如履平地）；黄金棍使得出神入化；双翅可生风雷护体。。功法表现：双翅卷起风雷逼退追兵，身侧雷光如蛇。。画面：构图：潼关城楼之上雷云翻滚，雷震子双翅全开展翅而飞，黄金棍斜指城下，身侧雷光如蛇；身后是黑压压的追兵虚影与城墙剪影。色调：玄黑深红+闪电白蓝，黄金棍一点明金。氛围：救主、雷霆、至孝。。台词："（他跪在白发文王面前）'父王莫怕，儿子翅膀大，飞得稳。当年您从雷雨里把我接回家，今日我披着雷雨送您回去。'"。动作帧（动图）：①展翅冲天 ②黄金棍横扫追兵 ③侧身负父于背 ④放缓翅速护父低飞。诗词：一翅风来一翅雷，潼关城下救父归。黄金棍扫千军散，负得慈亲出险围。。主题句：文王从雷雨里接回一个孩子，这个孩子用一生，把每一次雷雨都炼成了护人的翅膀。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄黑深红+闪电白蓝，黄金棍一点明金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 潼关救主, age 18岁, 潼关救主·雷翼护主; scene: 双翅卷起风雷逼退追兵；黄金棍荡开箭雨；背负文王低飞，怕惊着老人特意放缓翅速；沿途以雷光开路。; 一翅风来一翅雷，潼关城下救父归。黄金棍扫千军散，负得慈亲出险围。; palette: 玄黑深红+闪电白蓝，黄金棍一点明金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 救父出关**
- 雷震子，封神登天阶段·封神先锋（20余岁，封神先锋·雷部先锋）。形象：雷震子，鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。品性：随姜子牙兴周伐纣，为西岐雷部先锋。作战骁勇，从不贪功，护在武王与姜丞相左右。他记着师父那句话"你之异相，正是你之神通"——当年嫌丑的这张脸，如今让敌军闻风丧胆。。姿态：阵前展翅掠阵，黄金棍点破敌阵；雷声一响，敌军先自惊乱；战后蹲在营帐外为阵亡将士收殓，不发一言；武王劝他回朝受封，他摇头："我飞惯了，住不惯金殿。"。服饰：蓝靛面，朱砂发束金冠，银蓝战甲；风雷双翅金光大盛；黄金棍在手。。体型：身高约6头身，悬空展翅，战甲加身，威仪初成。。衣物细节：银蓝战甲，束发金冠，双翅全开金光大盛，身后周旗。。发型妆造：朱砂发束金冠，蓝靛面，金瞳如电。。脸型五官：鸟面威仪，蓝靛面，朱砂眉如刀，金瞳如电，鸟喙收紧，下巴坚毅，獠牙内敛。。武器招式：黄金棍通神，横扫/点刺/劈砸三连。。功法：风雷双翅全开（可催动漫天雷雨）；黄金棍通神；雷部先锋，可调雷雨开路。。功法表现：雷光自翅尖迸发，可催动漫天雷雨。。画面：构图：战场天穹之上，雷震子悬空展翅，风雷漫天，黄金棍横空；身后西岐旌旗与云海，下方山河在望。色调：紫金明黄宝蓝+纯白圣光，雷电自翅尖迸发。氛围：威临、骁勇、护苍生。。台词："（面对武王封赏）'末将不要封地。只要一样——风雷所至，护得住这天下老百姓的太平，便够了。'"。动作帧（动图）：①双翅大展掠阵 ②黄金棍点破敌阵 ③雷声一震敌军惊乱 ④战后蹲营外沉默。诗词：西岐帐下雷部英，一棍一翅万军惊。雷声最怕何处响？响处正是我刀兵。。主题句：文王从雷雨里接回一个孩子，这个孩子用一生，把每一次雷雨都炼成了护人的翅膀。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金明黄宝蓝+纯白圣光，雷电自翅尖迸发。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 封神先锋, age 20余岁, 封神先锋·雷部先锋; scene: 阵前展翅掠阵，黄金棍点破敌阵；雷声一响，敌军先自惊乱；战后蹲在营帐外为阵亡将士收殓，不发一言；武王劝他回朝受封，他摇头："我飞惯了，住不惯金殿。"; 西岐帐下雷部英，一棍一翅万军惊。雷声最怕何处响？响处正是我刀兵。; palette: 紫金明黄宝蓝+纯白圣光，雷电自翅尖迸发; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 风雷显圣**
- 雷震子，道果圆满阶段·肉身成圣（封神后，肉身成圣·雷部正神）。形象：雷震子，鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。品性：封神事毕，李靖、哪吒、杨戬、雷震子等七人肉身成圣，位列其中。不居天位，不争香火，奉天庭意旨行雷部神职——布雷诛邪、赏善罚恶。他记得自己的边界：行云布雨是龙王与雨师的差事，他只管打雷——龙王布雨时，他在云端为其开路，雷声隆隆，从不多下一滴雨，也从不少劈一道邪。守在人间与仙界的云隙之间：该劈的邪祟，劈它个明明白白；该护的人，一个都不许漏。他最喜欢的，还是化作一阵风掠过西岐城郊，听一听当年那片练武校场里的笑声。。姿态：左手擎雷令法旨（奉天庭意旨），右手黄金棍引金雷；数道雷光钉锁云隙妖邪（诛邪），城郊放风筝的孩童在另一侧安然无恙（护生）；龙王布雨时于云端布雷开路，雷声隆隆为其开道（雷司其职，不越雨司之权）；路见人间孩童放风筝，俯身吹出一口气，让风筝飞得更高些。。服饰：蓝靛面，朱砂发披散，紫金圣甲；双翅威仪展开，翅上"风""雷"二字金光灼灼；肩披风雷环；左手擎雷令法旨，右手执黄金棍。。体型：身高约7头身，云端立定，法相威仪，紫金圣甲。。衣物细节：紫金圣甲，肩披风雷环，双翅威仪展开"风""雷"灼灼，破云金光法相。。发型妆造：朱砂发披散，蓝靛面，紫金束发冠，金瞳。。脸型五官：鸟面法相，蓝靛面，朱砂眉入鬓，金瞳微阖含光，鸟喙微阖，下巴端方，不怒自威。。武器招式：左手雷令法旨（奉天庭意旨），右手黄金棍引天雷诛邪。。功法：风雷双翅圆满（一念可至天涯）；黄金棍引天雷诛邪；雷令法旨（奉天庭意旨布雷）；雷司其职、赏善罚恶；肉身成圣，不坏不死。。功法表现：数道金雷钉锁妖邪；俯身吹一口金光托住城郊孩童风筝。。画面：构图：破云金光法相之中，雷部正神紫金圣甲、双翅威仪展开、肩披风雷环，左手擎雷令法旨、右手黄金棍引金雷；数道雷光劈向云隙间被钉锁的妖邪黑气（诛邪），右下城郊安然、放风筝孩童被气流托起（护生）。色调：云白淡青+紫金法相+青白雷光+墨紫妖气。氛围：奉天庭意旨布雷诛邪、雷司其职不越界的神职威仪。。台词："（像对天上说，也像对老去的文王冢说）'有人说我是老天爷用雷雨捏出来的孩子。对，我是。所以我这辈子，就做雷该做的事——妖邪该劈时劈得明明白白，龙王布雨该开路时，我把雷声打在最前头。至于人间要护的人，一个都不许漏。'"。动作帧（动图）：①立于法云 ②左臂擎雷令法旨 ③右手黄金棍引金雷 ④金雷劈向云隙妖邪。诗词：肉身成圣返云端，风雷二字隐翅间。人间若问谁是雷，便是当年捡来童。。主题句：文王从雷雨里接回一个孩子，这个孩子用一生，把每一次雷雨都炼成了护人的翅膀。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：云白淡青+紫金法相+青白雷光+墨紫妖气。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 肉身成圣, age 封神后, 肉身成圣·雷部正神; scene: 左手擎雷令法旨（奉天庭意旨），右手黄金棍引金雷；数道雷光钉锁云隙妖邪（诛邪），城郊放风筝的孩童在另一侧安然无恙（护生）；龙王布雨时于云端布雷开路，雷声隆隆为其开道（雷司其职，不越雨司之权）；路见人间孩童放风筝，俯身吹出一口气，让风筝飞得更高些。; 肉身成圣返云端，风雷二字隐翅间。人间若问谁是雷，便是当年捡来童。; palette: 云白淡青+紫金法相+青白雷光+墨紫妖气; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 黄天化（`huang_tianhua`） · 人生档案版

**灵胎初醒 · 将门灵种**
- 黄天化，灵胎初醒阶段·将门灵童（出生，武成王府·天生瑞气）。形象：黄天化，青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。品性：武成王黄飞虎长子，将门之后。出生时满堂瑞气，天生一股虎胆豪气。。姿态：五岁便敢骑烈马，抓周抓的是父亲案上的虎符；听得金鼓之声，便不哭不闹。。服饰：总角小儿，穿小号武将劲装，腰悬玩具木剑。。体型：身高约3头身，总角幼童，敦实虎气。。衣物细节：小号武将劲装，腰悬玩具木剑。。发型妆造：总角双髻，虎头虎脑。。脸型五官：圆脸虎气，剑眉，大眼有神，鼻梁挺，咧嘴笑，下巴圆，额头光洁。。武器招式：木剑随手比划。。功法：将门家传武艺启蒙；天生神力，骑射之才。。功法表现：无神力，将门虎气。。画面：构图：武成王府演武场，总角孩童骑在小马驹上，手持木剑，身后父亲黄飞虎含笑而立，背景红墙金瓦。色调：暖金+朱红+石青。氛围：将门虎子、意气初萌。。台词："我爹是武成王，将来我要比他还能打！"。动作帧（动图）：①骑小马驹 ②手持木剑比划 ③抓周抓虎符 ④咧嘴笑得志得意满。诗词：将门一出满堂香，抓周虎符握中央。五岁敢骑烈马走，天生虎胆小儿郎。。主题句：少年热血报父恩，纵马持剑伐无道——英雄不问年岁，只问肝胆。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：暖金+朱红+石青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 将门灵童, age 出生, 武成王府·天生瑞气; scene: 五岁便敢骑烈马，抓周抓的是父亲案上的虎符；听得金鼓之声，便不哭不闹。; 将门一出满堂香，抓周虎符握中央。五岁敢骑烈马走，天生虎胆小儿郎。; palette: 暖金+朱红+石青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年天化**
- 黄天化，凡尘砺心阶段·紫阳学道（幼年，青峰山·清虚道德真君收徒）。形象：黄天化，青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。品性：自幼被清虚道德真君带上青峰山紫阳洞学道。虽是道童，火性不改，常把丹房闹得鸡飞狗跳。。姿态：清晨随师父采药练气，午后偷溜下山比试武艺；师父罚他抄《道德经》，他抄着抄着画起兵器图。。服饰：青布道袍，束发小髻，腰间别着师父赐的短剑。。体型：身高约4头身，小道童，精瘦好动。。衣物细节：青布道袍，腰别短剑。。发型妆造：束发小髻，眉目英气。。脸型五官：少年圆脸，剑眉，眼神明亮，鼻梁挺，嘴角微扬，下巴圆，英气。。武器招式：短剑初学，剑法日进。。功法：吐纳导引，道术筑基；暗习家传武艺，剑法日渐精进。。功法表现：道术筑基，吐纳导引。。画面：构图：青峰山紫阳洞前，小道童持剑而舞，道袍翻飞，身后古松与仙云，山下一片苍茫。色调：黛青+云白+一点朱红。氛围：山中苦修、少年火性。。台词："师父，我学道法，也学剑法。等我下山，要把天下不平事都砍平了。"。动作帧（动图）：①采药练气 ②偷溜下山 ③比试武艺 ④抄道德经画兵器图。诗词：青峰山上紫阳开，道童偏偏爱剑来。丹房偷炼三分火，只待下山动九垓。。主题句：少年热血报父恩，纵马持剑伐无道——英雄不问年岁，只问肝胆。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青+云白+一点朱红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 紫阳学道, age 幼年, 青峰山·清虚道德真君收徒; scene: 清晨随师父采药练气，午后偷溜下山比试武艺；师父罚他抄《道德经》，他抄着抄着画起兵器图。; 青峰山上紫阳开，道童偏偏爱剑来。丹房偷炼三分火，只待下山动九垓。; palette: 黛青+云白+一点朱红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 火龙标**
- 黄天化，道法初成阶段·下山救父（少年，奉师命下山·玉麒麟莫邪剑）。形象：黄天化，青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。品性：闻父亲黄飞虎被纣王逼反、困于途，奉师命下山救父。师父赐莫邪宝剑、攒心钉，又赠玉麒麟为坐骑。。姿态：乘玉麒麟昼夜兼程，路上斩关夺隘；见父亲时翻身下骑，单膝跪地——"孩儿来了"。。服饰：青布战袍，莫邪剑初佩，胯下玉麒麟，眉目英武。。体型：身高约5头身，少年将军，跨玉麒麟。。衣物细节：青布战袍，莫邪剑初佩，玉麒麟坐骑。。发型妆造：束发簪冠，眉目英武。。脸型五官：少年方脸，剑眉入鬓，凤目英武，鼻梁挺，唇线坚毅，下巴方正。。武器招式：莫邪剑出鞘如虹；攒心钉百发百中。。功法：莫邪宝剑出鞘如虹；攒心钉百发百中；玉麒麟日行千里。。功法表现：玉麒麟日行千里，踏云生风。。画面：构图：关隘之前，少年将军胯下玉麒麟昂首，莫邪剑出鞘寒光如虹，身后是黄飞虎的残军；背景雄关与烽烟。色调：青布+剑光白+战火绛。氛围：救父心切、少年破阵。。台词："父亲莫慌，孩儿已得师父真传。这天下，谁也伤不了您！"。动作帧（动图）：①乘玉麒麟昼夜兼程 ②斩关夺隘 ③翻身下骑 ④单膝跪地"孩儿来了"。诗词：师命下山救父归，麒麟踏碎五关围。莫邪一剑寒光起，攒心钉下鬼神悲。。主题句：少年热血报父恩，纵马持剑伐无道——英雄不问年岁，只问肝胆。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青布+剑光白+战火绛。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 下山救父, age 少年, 奉师命下山·玉麒麟莫邪剑; scene: 乘玉麒麟昼夜兼程，路上斩关夺隘；见父亲时翻身下骑，单膝跪地——"孩儿来了"。; 师命下山救父归，麒麟踏碎五关围。莫邪一剑寒光起，攒心钉下鬼神悲。; palette: 青布+剑光白+战火绛; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 下山辅周**
- 黄天化，大劫淬炼阶段·伐纣先锋（少年，金鸡岭大战·攒心钉扬威）。形象：黄天化，青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。品性：随父伐纣，为西岐前锋。骁勇善战，攒心钉钉敌将如探囊取物，玉麒麟冲阵无人能挡。。姿态：阵前纵麒麟直取敌将，攒心钉连钉数员大将；战后不居功，只问"父王安好？"。服饰：少年将军金甲，玉麒麟为骑，莫邪剑在手，英姿飒爽。。体型：身高约6头身，少年金甲将军，英姿飒爽。。衣物细节：少年将军金甲，玉麒麟为骑，莫邪剑在手。。发型妆造：金冠束发，英姿飒爽。。脸型五官：少年方脸，剑眉，凤目如虹，鼻梁挺，嘴角带战意，下巴方正，英姿。。武器招式：莫邪剑通神；攒心钉例无虚发；玉麒麟冲阵。。功法：莫邪剑通神；攒心钉例无虚发；玉麒麟冲阵踏营。。功法表现：攒心钉金光脱手；麒麟冲阵踏营。。画面：构图：金鸡岭战场，少年将军乘玉麒麟跃阵而起，莫邪剑高举，攒心钉金光脱手；背景旌旗与烽火。色调：金甲明黄+战火红+墨黑。氛围：少年悍勇、一往无前。。台词："父亲在前，我就在前。敌军再多，我黄天化也要杀出一条血路！"。动作帧（动图）：①纵麒麟直取敌将 ②攒心钉连钉数将 ③莫邪剑横扫 ④战后问"父王安好？"。诗词：西岐帐下一先锋，攒心钉落将星空。麒麟踏破千营帐，少年血性贯长虹。。主题句：少年热血报父恩，纵马持剑伐无道——英雄不问年岁，只问肝胆。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：金甲明黄+战火红+墨黑。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 伐纣先锋, age 少年, 金鸡岭大战·攒心钉扬威; scene: 阵前纵麒麟直取敌将，攒心钉连钉数员大将；战后不居功，只问"父王安好？"; 西岐帐下一先锋，攒心钉落将星空。麒麟踏破千营帐，少年血性贯长虹。; palette: 金甲明黄+战火红+墨黑; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 将门虎子**
- 黄天化，封神登天阶段·捐躯金鸡岭（少年，战殁沙场·三山正神）。形象：黄天化，青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。品性：金鸡岭大战中浴血苦战，为破敌阵孤身陷围，身负重伤仍力战不退——少年将军，血洒疆场。。姿态：敌阵中枪失前蹄，他翻身而起，持剑再战；气绝时仍面向敌阵，虎目圆睁。。服饰：金甲染血，莫邪剑拄地，玉麒麟伏于身侧低鸣。。体型：身高约6头身，金甲染血，倚剑而立。。衣物细节：金甲染血，莫邪剑拄地，玉麒麟伏于身侧。。发型妆造：散发染血，英气不灭。。脸型五官：少年方脸，剑眉染血，凤目圆睁不屈，鼻梁挺，唇紧抿，下巴方正，英气不灭。。武器招式：莫邪剑力竭仍战。。功法：力竭仍战；一身肝胆，至死不退。。功法表现：忠魂不灭，天边残阳如血。。画面：构图：金鸡岭战场暮色，少年将军金甲染血倚剑而立，玉麒麟伏地，远方战旗半卷；天边一线残阳如血。色调：血绛+墨黑+残金，悲壮。氛围：壮烈、陨落、英魂不灭。。台词："父亲……孩儿先走一步。这战场，我守到了最后一刻……"。动作帧（动图）：①敌阵中枪失前蹄 ②翻身而起持剑再战 ③力竭仍战 ④气绝时虎目圆睁望敌阵。诗词：金鸡岭上少年殇，一腔热血洒疆场。马革裹尸何所惧，英魂直上封神堂。。主题句：少年热血报父恩，纵马持剑伐无道——英雄不问年岁，只问肝胆。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：血绛+墨黑+残金，悲壮。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 捐躯金鸡岭, age 少年, 战殁沙场·三山正神; scene: 敌阵中枪失前蹄，他翻身而起，持剑再战；气绝时仍面向敌阵，虎目圆睁。; 金鸡岭上少年殇，一腔热血洒疆场。马革裹尸何所惧，英魂直上封神堂。; palette: 血绛+墨黑+残金，悲壮; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 莫邪封神**
- 黄天化，道果圆满阶段·炳灵公（封神后，三山正神·烈火长明）。形象：黄天化，青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。品性：封神后受封三山正神炳灵公，执掌三山烟火。少年将军虽殁，一身烈火长明不熄，护佑山川与行旅。。姿态：巡守三山，护佑樵夫行旅；山火起时，以神火烧尽妖氛；依旧是那副少年将军的模样。。服饰：炳灵公神袍，周身烈火神光，莫邪剑化作山间虹气，玉麒麟常伴。。体型：身高约7头身，炳灵公神袍，烈火环绕。。衣物细节：炳灵公神袍，烈火神光，莫邪剑化山间虹气。。发型妆造：束发，少年面容，烈火为光。。脸型五官：炳灵公少年方脸，剑眉，凤目含笑含威，鼻梁挺，嘴角从容，下巴方正，神光。。武器招式：莫邪剑化作山间虹气护山。。功法：三山正神神位；烈火长明不熄；护山护人，佑一方平安。。功法表现：烈火长明不熄，火凤盘旋。。画面：构图：三山云海之上，少年将军神袍烈火环绕，周身火凤盘旋，莫邪剑气化作虹光，俯瞰苍翠群山与人间烟火。色调：烈火赤金+山青+神光白，炽烈而安宁。氛围：英魂长明、护佑山川。。台词："生为将，殁为神，我黄天化，天生就是护三山护百姓的命。"。动作帧（动图）：①巡守三山 ②护佑樵夫行旅 ③山火起以神火焚妖氛 ④仍是少年将军模样。诗词：三山正神炳灵公，烈火长明照碧空。少年一殁忠魂在，化作风雷镇万峰。。主题句：少年热血报父恩，纵马持剑伐无道——英雄不问年岁，只问肝胆。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：烈火赤金+山青+神光白，炽烈而安宁。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 炳灵公, age 封神后, 三山正神·烈火长明; scene: 巡守三山，护佑樵夫行旅；山火起时，以神火烧尽妖氛；依旧是那副少年将军的模样。; 三山正神炳灵公，烈火长明照碧空。少年一殁忠魂在，化作风雷镇万峰。; palette: 烈火赤金+山青+神光白，炽烈而安宁; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 土行孙（`tu_xingsun`） · 人生档案版

**灵胎初醒 · 地灵种子**
- 土行孙，灵胎初醒阶段·夹龙山童（出生，惧留孙门下·生得矮小）。形象：土行孙，矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。品性：夹龙山飞云洞惧留孙门下小徒弟，生得矮小，却机灵鬼马，一肚子主意。。姿态：个子矮够不着桌案，就发明了踩高跷扫地；师父炼丹，他偷偷给炉子扇风想抢功。。服饰：灰布短打，扎小髻，腰间别着一根短棍。。体型：身高约2.5头身，矮小圆胖道童。。衣物细节：灰布短打，腰别短棍。。发型妆造：扎小髻，机灵鬼马。。脸型五官：圆脸胖腮，细眉眯眼，蒜头鼻，咧嘴笑，下巴圆，机灵相。。武器招式：短棍随手耍。。功法：机灵过人；一肚子土遁奇门的心思。。功法表现：无神力，机灵过人。。画面：构图：夹龙山飞云洞前，矮小道童踩着高跷扫地，衣摆翻飞，背景云雾山峦。色调：灰褐+土黄+青灰，俏皮。氛围：机灵、乐观、身矮志高。。台词："师父，个矮怎么了？山不在高，有仙则名——我土行孙，照样是仙！"。动作帧（动图）：①踩高跷扫地 ②够不着桌案垫脚 ③给炉子扇风 ④得意叉腰。诗词：夹龙山上一小童，身矮偏有凌云胸。踩跷扫地也堪乐，笑问天公谁最灵。。主题句：矮小身躯藏神技，土中自有真天地——不争高，只争快。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：灰褐+土黄+青灰，俏皮。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 夹龙山童, age 出生, 惧留孙门下·生得矮小; scene: 个子矮够不着桌案，就发明了踩高跷扫地；师父炼丹，他偷偷给炉子扇风想抢功。; 夹龙山上一小童，身矮偏有凌云胸。踩跷扫地也堪乐，笑问天公谁最灵。; palette: 灰褐+土黄+青灰，俏皮; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年土行**
- 土行孙，凡尘砺心阶段·偷学土遁（幼年，地行术·遁地千里）。形象：土行孙，矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。品性：见师父地行术神奇，趁师父闭关偷学口诀。第一次遁地就迷了路，从茅厕里钻出来，被罚扫山门。。姿态：偷学地行术，整日在山间钻来钻去；练成后能在土中如鱼得水，片刻穿山越岭。。服饰：灰布短打沾满泥土，腰间短棍改成镔铁棍。。体型：身高约2.5头身，灰衣沾泥。。衣物细节：灰布短打沾满泥土，腰间镔铁棍。。发型妆造：乱发沾土，笑嘻嘻。。脸型五官：圆脸沾土，细眉扬起，眯眼贼笑，蒜头鼻，下巴圆。。武器招式：镔铁棍初练。。功法：地行术大成（日行千里，土中穿行）；镔铁棍初练。。功法表现：地行术大成，土中如鱼得水。。画面：构图：夹龙山间，矮小少年从土里探出半个身子，泥土飞溅，一脸得意，身后尘土成丘。色调：土黄+赭石+草青。氛围：顽皮、得意、土中天地。。台词："师父您看，这土里的天地，比地上还大呢！我土行孙，天生就是土里长出来的！"。动作帧（动图）：①偷学口诀 ②钻入土中 ③从茅厕探出 ④被罚扫山门仍笑嘻嘻。诗词：偷得地行一遁法，山间钻透九重沙。笑从土里探出头，师父气笑手叉腰。。主题句：矮小身躯藏神技，土中自有真天地——不争高，只争快。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：土黄+赭石+草青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 偷学土遁, age 幼年, 地行术·遁地千里; scene: 偷学地行术，整日在山间钻来钻去；练成后能在土中如鱼得水，片刻穿山越岭。; 偷得地行一遁法，山间钻透九重沙。笑从土里探出头，师父气笑手叉腰。; palette: 土黄+赭石+草青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 土遁神通**
- 土行孙，道法初成阶段·被撺掇下山（少年，申公豹撺掇·助纣被擒）。形象：土行孙，矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。品性：被申公豹三寸舌说动，下山助商。捆仙绳绑了西岐几员大将，风光一时，后被杨戬设计擒获。。姿态：凭捆仙绳与地行术屡建"功"，却又被识破；被擒后得知是非，低头认错。。服饰：战甲短打，镔铁棍与捆仙绳在手。。体型：身高约2.5头身，矮小战将，捆仙绳在手。。衣物细节：战甲短打，镔铁棍与捆仙绳。。发型妆造：束发，垂头丧气。。脸型五官：圆脸，细眉，眯眼精明，蒜头鼻，撇嘴，下巴圆，垂头丧气。。武器招式：捆仙绳（仙魔难逃）。。功法：捆仙绳（惧留孙所授，仙魔难逃）；地行术穿营劫阵。。功法表现：捆仙绳金光脱手。。画面：构图：军帐之中，矮小战将手持捆仙绳欲捆敌将，身后申公豹黑影鼓舌，帐外旌旗；背景商周两军对峙。色调：玄黑+捆绳金+土黄。氛围：误入歧途、机灵反被聪明误。。台词："哎，我这是被那张嘴坑了！捆仙绳捆得住别人，捆不住我这颗糊涂心。"。动作帧（动图）：①下山助商 ②捆仙绳绑敌将 ③被杨戬设计擒获 ④低头认错。诗词：下山误入云与雾，三寸舌前尽糊涂。捆仙绳长难自缚，回头才见正路途。。主题句：矮小身躯藏神技，土中自有真天地——不争高，只争快。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄黑+捆绳金+土黄。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 被撺掇下山, age 少年, 申公豹撺掇·助纣被擒; scene: 凭捆仙绳与地行术屡建"功"，却又被识破；被擒后得知是非，低头认错。; 下山误入云与雾，三寸舌前尽糊涂。捆仙绳长难自缚，回头才见正路途。; palette: 玄黑+捆绳金+土黄; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 地行如飞**
- 土行孙，大劫淬炼阶段·归周立功（少年，将功赎罪·地行千里建功）。形象：土行孙，矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。品性：归周后将功赎罪，地行术屡立奇功——钻地断粮道、潜行盗兵符、地底擒敌将，成了姜子牙手里的"土遁奇兵"。。姿态：夜半遁地潜入敌营，盗取兵符；大战时自地下突然杀出，打敌一个措手不及。。服饰：西岐军士短打，镔铁棍与捆仙绳随身。。体型：身高约2.5头身，西岐军士短打。。衣物细节：西岐军士短打，镔铁棍捆仙绳。。发型妆造：束发，精神抖擞。。脸型五官：圆脸，细眉，眯眼机灵，蒜头鼻，咧嘴笑，下巴圆。。武器招式：镔铁棍横扫；捆仙绳擒敌。。功法：地行术出神入化；捆仙绳擒敌；镔铁棍近战。。功法表现：地行术神出鬼没。。画面：构图：夜营沙场，矮小军士自地下破土而出，镔铁棍横扫，捆仙绳金光脱手；背景月光下的营帐。色调：夜蓝+土金+明黄。氛围：神出鬼没、土遁奇功。。台词："姜丞相，打仗我不跟人比身高，比谁先到！这地下，是我土行孙的主场！"。动作帧（动图）：①夜半遁地 ②潜入敌营盗兵符 ③大战自地下杀出 ④打敌措手不及。诗词：归周洗心再立功，地行千里若游龙。盗符断道穿营去，矮小身形立奇功。。主题句：矮小身躯藏神技，土中自有真天地——不争高，只争快。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：夜蓝+土金+明黄。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 归周立功, age 少年, 将功赎罪·地行千里建功; scene: 夜半遁地潜入敌营，盗取兵符；大战时自地下突然杀出，打敌一个措手不及。; 归周洗心再立功，地行千里若游龙。盗符断道穿营去，矮小身形立奇功。; palette: 夜蓝+土金+明黄; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 入地为营**
- 土行孙，封神登天阶段·土府星君（封神后，封神榜·执掌土行）。形象：土行孙，矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。品性：封神后受封土府星君，执掌大地土行。矮小身躯里，装着一方大地的神通。。姿态：掌土府星位，巡行山川大地；勘山测土、护田固基，人间土木之事皆归其管。。服饰：星君官袍，镔铁棍化作执笏，捆仙绳化作腰绦。。体型：身高约3头身，星君官袍，端方。。衣物细节：星君官袍，镔铁棍化执笏，捆仙绳化腰绦。。发型妆造：官帽束发。。脸型五官：圆脸，细眉，眯眼端方，蒜头鼻，嘴角含笑，下巴圆，官威。。武器招式：执笏镇土。。功法：土府星君神位；地行术掌大地之脉；可缩地成寸。。功法表现：地行术掌大地之脉，缩地成寸。。画面：构图：云端土府星君殿前，矮小星君执笏而立，脚下山河大地如棋盘，周身土黄星光环绕；背景大地山川。色调：土金+星青+云白，厚重又精神。氛围：执掌大地、小身担大任。。台词："封我土府星君，我认！大地厚德载物，我土行孙，就做这大地的看门人。"。动作帧（动图）：①掌土府星位 ②巡行山川 ③勘山测土 ④护田固基。诗词：土府星君镇厚坤，地行千里掌三坟。矮身偏担大地任，德载万物一寸真。。主题句：矮小身躯藏神技，土中自有真天地——不争高，只争快。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：土金+星青+云白，厚重又精神。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 土府星君, age 封神后, 封神榜·执掌土行; scene: 掌土府星位，巡行山川大地；勘山测土、护田固基，人间土木之事皆归其管。; 土府星君镇厚坤，地行千里掌三坟。矮身偏担大地任，德载万物一寸真。; palette: 土金+星青+云白，厚重又精神; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 地行仙归**
- 土行孙，道果圆满阶段·地行无极（千年后，与地同行·快过日月）。形象：土行孙，矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。品性：千百年执掌土行，地行术已臻化境——一念遁地，千里可至；与大地同息，与山川同行。仍是那副矮矮的、笑嘻嘻的模样。。姿态：偶尔化作一阵土雾掠过人间，替迷路的旅人指路；地下过客，无不闻其名而称奇。。服饰：星君便装，镔铁棍斜挎，脚不沾地，周身土灵。。体型：身高约3头身，便装星君，脚不沾地。。衣物细节：星君便装，镔铁棍斜挎。。发型妆造：束发，笑意盈盈。。脸型五官：圆脸，细眉，眯眼含笑，蒜头鼻，咧嘴笑，下巴圆，逍遥。。武器招式：镔铁棍；捆仙绳一念化万绳。。功法：地行无极（一念千里）；与地同息、山川为路；捆仙绳一念化万绳。。功法表现：一念千里，土金流光。。画面：构图：黄昏大地，矮小身影化作一道土金流光掠过大川长河，身后大地如锦；背景暮色与群山。色调：土金+暮橙+黛青。氛围：逍遥、神速、与地同游。。台词："有人问我，一辈子在地底下钻，闷不闷？我笑——大地的热闹，你们看不见。"。动作帧（动图）：①化作土雾掠过人间 ②替迷路旅人指路 ③一念遁地 ④笑嘻嘻现身。诗词：地行无极念即程，山川为路日月轻。矮身一遁三千里，笑看人间岁月更。。主题句：矮小身躯藏神技，土中自有真天地——不争高，只争快。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：土金+暮橙+黛青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 地行无极, age 千年后, 与地同行·快过日月; scene: 偶尔化作一阵土雾掠过人间，替迷路的旅人指路；地下过客，无不闻其名而称奇。; 地行无极念即程，山川为路日月轻。矮身一遁三千里，笑看人间岁月更。; palette: 土金+暮橙+黛青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 杨任（`yang_ren`） · 人生档案版

**灵胎初醒 · 忠臣灵光**
- 杨任，灵胎初醒阶段·直言上大夫（壮年，商朝上大夫·犯颜直谏）。形象：杨任，双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。品性：商朝上大夫，忠直耿介。朝堂上人皆唯唯，独他敢对纣王直言相谏。。姿态：朝会之上，据理力争，句句掷地有声；退朝后仍忧国事，秉烛写谏章。。服饰：朝服冠冕，腰佩玉带，眉目端方。。体型：身高约7头身，朝服文臣，端方。。衣物细节：朝服冠冕，腰佩玉带。。发型妆造：束发戴冠，眉目端方。。脸型五官：国字脸，浓眉，虎目清正，鼻梁挺，唇线刚直，下巴方正。。武器招式：无兵器，手持谏章。。功法：通晓治国之策；一腔浩然正气。。功法表现：浩然正气。。画面：构图：商朝朝堂，一位朝服文臣手持谏章立于阶前，气势凛然，身后群臣唯唯诺诺的虚影。色调：玄黑+朝服绛红+烛火金。氛围：忠直、孤勇、犯颜直谏。。台词："臣食君禄，当言君过。大王纵要杀我，这谏章臣也要写完！"。动作帧（动图）：①朝会据理力争 ②句句掷地有声 ③退朝秉烛写谏章 ④执谏章立于阶前。诗词：朝堂唯唯皆无骨，独有杨公敢诤言。一纸谏章肝与胆，忠魂直上九重天。。主题句：剜目不改忠臣志，掌中生目更清明——看得越透，走得越正。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄黑+朝服绛红+烛火金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 直言上大夫, age 壮年, 商朝上大夫·犯颜直谏; scene: 朝会之上，据理力争，句句掷地有声；退朝后仍忧国事，秉烛写谏章。; 朝堂唯唯皆无骨，独有杨公敢诤言。一纸谏章肝与胆，忠魂直上九重天。; palette: 玄黑+朝服绛红+烛火金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年杨任**
- 杨任，凡尘砺心阶段·剜目之祸（壮年，触怒纣王·双目被剜）。形象：杨任，双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。品性：纣王怒其直谏，施以炮烙般酷刑——剜其双目。他至死不悔，昏死过去。。姿态：被施刑时睁着双眼直视纣王，至死不低头；血溅朝堂，被弃于荒野。。服饰：朝服染血，双目殷红，被弃于路旁。。体型：身高约7头身，朝服染血，被弃于野。。衣物细节：朝服染血，双目殷红。。发型妆造：散发，面无血色。。脸型五官：国字脸，眉紧蹙，双目殷红，鼻梁挺，唇紧抿，下巴方正。。武器招式：无兵器，攥着谏章。。功法：浩然正气不灭；一缕忠魂不散。。功法表现：浩然正气不灭。。画面：构图：残阳荒郊，被剜目的文臣倒在路边，血染朝服，手中仍攥着那卷谏章；远处朝堂飞檐如笼。色调：血绛+玄黑+惨白，凄厉悲壮。氛围：冤屈、孤忠、不屈。。台词："天日昭昭，忠臣之心，剜得掉眼睛，剜不掉正气！"。动作帧（动图）：①被施刑仍直视纣王 ②血溅朝堂 ③昏死被弃 ④手中仍攥谏章。诗词：谏章触怒帝王心，剜目血溅紫宸深。忠骨何曾惧酷烈，荒郊野老抱孤贞。。主题句：剜目不改忠臣志，掌中生目更清明——看得越透，走得越正。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：血绛+玄黑+惨白，凄厉悲壮。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 剜目之祸, age 壮年, 触怒纣王·双目被剜; scene: 被施刑时睁着双眼直视纣王，至死不低头；血溅朝堂，被弃于荒野。; 谏章触怒帝王心，剜目血溅紫宸深。忠骨何曾惧酷烈，荒郊野老抱孤贞。; palette: 血绛+玄黑+惨白，凄厉悲壮; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 被挖双目**
- 杨任，道法初成阶段·掌中生目（壮年，道德真君救治·掌目重光）。形象：杨任，双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。品性：清虚道德真君闻其忠义，以仙药救其残躯，令其掌中生目——从此一目在手，照彻阴阳，重见天日。。姿态：掌中双目光华流转，第一次"看"见山色时，他老泪纵横；双手一抬，掌目扫视天地。。服饰：青布道袍，双手掌心各生一目，光华内敛。。体型：身高约7头身，青袍道者，掌中双目。。衣物细节：青布道袍，双手掌心各生一目。。发型妆造：束发，掌目放光。。脸型五官：国字脸，眉目（掌目）清明，鼻梁挺，嘴角舒展，下巴方正，老泪纵横后释然。。武器招式：无兵器，掌目为眼。。功法：掌中天目（能视地三尺、照破虚妄）；得仙家真传。。功法表现：掌中天目照彻虚妄，光华照云海。。画面：构图：仙山洞府晨光，青袍文臣双手抬起，掌中双目放光，光华照彻云海；背景青山如洗。色调：青袍+掌目金光+云白。氛围：劫后重生、看破世情。。台词："夺我双目者，是天；还我光明者，是道。这一双眼，我要用来看清人间的是非。"。动作帧（动图）：①仙药救残躯 ②掌中生目 ③第一次"看"见山色老泪纵横 ④双手一抬掌目扫视天地。诗词：仙家救取忠臣魂，掌中天目照乾坤。剜目何曾消正气，重光更看世事真。。主题句：剜目不改忠臣志，掌中生目更清明——看得越透，走得越正。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青袍+掌目金光+云白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 掌中生目, age 壮年, 道德真君救治·掌目重光; scene: 掌中双目光华流转，第一次"看"见山色时，他老泪纵横；双手一抬，掌目扫视天地。; 仙家救取忠臣魂，掌中天目照乾坤。剜目何曾消正气，重光更看世事真。; palette: 青袍+掌目金光+云白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 眼中生手**
- 杨任，大劫淬炼阶段·辅周伐纣（壮年，飞电枪·云霞兽辅周）。形象：杨任，双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。品性：怀报国之志辅周伐纣。骑云霞兽、使飞电枪，运筹帷幄，以掌中目识破敌军虚实，屡建奇功。。姿态：阵前掌目一照，敌营虚实尽在掌握；云霞兽踏云突袭，飞电枪直取敌将。。服饰：道袍外罩战甲，掌中目威光凛凛，胯下云霞兽。。体型：身高约7头身，道袍外罩战甲，骑云霞兽。。衣物细节：道袍外罩战甲，胯下云霞兽。。发型妆造：束发，掌目威光凛凛。。脸型五官：国字脸，浓眉，掌目放光，鼻梁挺，嘴角坚毅，下巴方正。。武器招式：飞电枪如电。。功法：掌中天目照破虚实；飞电枪如电；云霞兽腾云。。功法表现：掌中天目照破虚实。。画面：构图：两军阵前，青袍将军骑云霞兽立于云端，掌中双目照出漫天光华，飞电枪指向前方敌阵；背景烽火云海。色调：道青+电光蓝+云白。氛围：运筹帷幄、以目破敌。。台词："纣王以为剜了我的眼，我便瞎了。我偏要睁着掌中这双眼，看着他的江山倾覆。"。动作帧（动图）：①阵前掌目一照 ②敌营虚实尽现 ③云霞兽踏云突袭 ④飞电枪直取敌将。诗词：辅周伐纣抱孤忠，掌目明察万里戎。云霞兽踏飞电起，兵行诡道立奇功。。主题句：剜目不改忠臣志，掌中生目更清明——看得越透，走得越正。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道青+电光蓝+云白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 辅周伐纣, age 壮年, 飞电枪·云霞兽辅周; scene: 阵前掌目一照，敌营虚实尽在掌握；云霞兽踏云突袭，飞电枪直取敌将。; 辅周伐纣抱孤忠，掌目明察万里戎。云霞兽踏飞电起，兵行诡道立奇功。; palette: 道青+电光蓝+云白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 飞电枪出**
- 杨任，封神登天阶段·甲子太岁（封神后，受封甲子·执掌岁月）。形象：杨任，双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。品性：封神后受封甲子太岁，执掌岁月轮转。忠臣之目，看尽人间四季更替、朝代兴衰。。姿态：掌甲子轮回，司人间岁时；以掌中目勘岁月之痕，护耕读有序。。服饰：太岁神袍，掌中天目温润，岁月神光内敛。。体型：身高约7头身，太岁神袍，端坐。。衣物细节：太岁神袍，掌目温润。。发型妆造：束发，岁月神光内敛。。脸型五官：国字脸，白眉，掌目温润，鼻梁挺，嘴角肃穆，下巴方正。。武器招式：无兵器，掌目为法。。功法：甲子太岁神位；执掌岁月；掌中目察兴衰。。功法表现：二十四节气光轮环绕。。画面：构图：太岁神殿，神袍文臣端坐，掌中目光化作二十四节气的光轮环绕，身后四季景象流转；背景浩瀚岁月。色调：岁月金+昼夜蓝+云白。氛围：执掌岁月、清明肃穆。。台词："我管岁月，也管公道。天地四时有序，人心也该有个清浊分明。"。动作帧（动图）：①掌甲子轮回 ②司人间岁时 ③以掌目勘岁月 ④护耕读有序。诗词：甲子星高掌岁轮，忠臣化神照尘尘。掌中一眼观兴替，四时有序太岁尊。。主题句：剜目不改忠臣志，掌中生目更清明——看得越透，走得越正。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：岁月金+昼夜蓝+云白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 甲子太岁, age 封神后, 受封甲子·执掌岁月; scene: 掌甲子轮回，司人间岁时；以掌中目勘岁月之痕，护耕读有序。; 甲子星高掌岁轮，忠臣化神照尘尘。掌中一眼观兴替，四时有序太岁尊。; palette: 岁月金+昼夜蓝+云白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 道德真君**
- 杨任，道果圆满阶段·明照幽明（神位长存，洞视阴阳·掌目永明）。形象：杨任，双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。品性：千百年执掌岁月，掌中天目已能照彻幽冥阳世。忠臣之心不灭，目光所及，善恶皆明。。姿态：以掌目照察人间善恶、幽冥冤屈；夜半为迷途游魂指路，白日为困顿之人指明路。。服饰：太岁便服，掌中双目神光流转，周身浩然正气。。体型：身高约7头身，太岁便服，浩然正气。。衣物细节：太岁便服，掌目神光流转。。发型妆造：束发，双目（掌目）明澈。。脸型五官：国字脸，眉目明澈，掌目神光，鼻梁挺，嘴角安和，下巴方正。。武器招式：无兵器。。功法：洞视幽明；甲子岁月之力；掌目照善恶。。功法表现：掌目双光一道照人间一道照幽冥。。画面：构图：黄昏阴阳交汇处，神袍文臣掌中双目放出两道明光，一道照人间城郭、一道照幽冥微光，身周岁月流转；背景暮色苍茫。色调：明光金+幽冥蓝+暮紫。氛围：洞视幽明、正气长存。。台词："我这双眼，天生就是用来分是非的。天能剜它一次，剜不掉我照尽天下的念头。"。动作帧（动图）：①掌目照察人间善恶 ②照幽冥冤屈 ③为迷途游魂指路 ④为困顿人指明路。诗词：幽明两界掌中看，善恶分明照不偏。岁月如轮身是轴，忠魂不灭自安然。。主题句：剜目不改忠臣志，掌中生目更清明——看得越透，走得越正。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：明光金+幽冥蓝+暮紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 明照幽明, age 神位长存, 洞视阴阳·掌目永明; scene: 以掌目照察人间善恶、幽冥冤屈；夜半为迷途游魂指路，白日为困顿之人指明路。; 幽明两界掌中看，善恶分明照不偏。岁月如轮身是轴，忠魂不灭自安然。; palette: 明光金+幽冥蓝+暮紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 韦护（`wei_hu`） · 人生档案版

**灵胎初醒 · 金刚灵种**
- 韦护，灵胎初醒阶段·山中少年（出生，道行天尊门下·道心初立）。形象：韦护，身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。品性：道行天尊弟子，山中修道的少年。性憨厚，心向善，一心想修成护法神通。。姿态：清晨随师兄采药，午后在山溪边练吐纳；见山中樵夫被野兽所伤，背其下山医治。。服饰：青布道袍，束发，眉目憨厚沉稳。。体型：身高约4头身，青袍少年道童，憨厚沉稳。。衣物细节：青布道袍，束发，无饰。。发型妆造：束发，眉目憨厚。。脸型五官：圆脸憨厚，眉目柔和，眼大而正，鼻梁直，嘴角带憨笑，下巴圆。。武器招式：无兵器。。功法：吐纳导引，道心初立。。功法表现：吐纳导引，道心初立。。画面：构图：云雾山间，青袍少年道童立于溪边吐纳练气，身后古松飞瀑，山鸟相伴。色调：山青+云白+褐木。氛围：纯朴、向善、道心初萌。。台词："师父，我想修一身护人的本领。山里的老虎不伤人，世间的恶念才伤人。"。动作帧（动图）：①采药 ②溪边练吐纳 ③背樵夫下山医治 ④俯身问伤。诗词：山中修道少年郎，憨厚心肠日月光。采药溪边闻道法，他年一杵镇八方。。主题句：一杵降魔，护法三教——金刚怒目，亦是慈悲。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：山青+云白+褐木。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 山中少年, age 出生, 道行天尊门下·道心初立; scene: 清晨随师兄采药，午后在山溪边练吐纳；见山中樵夫被野兽所伤，背其下山医治。; 山中修道少年郎，憨厚心肠日月光。采药溪边闻道法，他年一杵镇八方。; palette: 山青+云白+褐木; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年韦护**
- 韦护，凡尘砺心阶段·金光护体（幼年，玄门勤修·初窥神通）。形象：韦护，身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。品性：勤修苦练，玄门功夫日渐精深，周身渐生护体金光，越发笃定护法之志。。姿态：日夜修持不辍，山石前练掌，金光随掌而起；下山云游，路遇不平，以金光护身相援。。服饰：青布道袍已洗得发白，周身若有若无的金光。。体型：身高约4头身，青袍泛金光。。衣物细节：青布道袍洗得发白，周身若有金光。。发型妆造：束发。。脸型五官：圆脸，眉目沉静，眼含金光，鼻梁直，嘴角笃定，下巴圆。。武器招式：无兵器，掌中金光。。功法：护体金光初成；道法筑基。。功法表现：护体金光初成。。画面：构图：山间道场，少年道人一掌推出，周身泛起淡淡金光，衣袍翻飞；背景晨光山影。色调：道青+护体金+晨光白。氛围：勤修、初露神通。。台词："这金光护得住我，更要护得住别人。修炼不是为自己，是为能伸出手。"。动作帧（动图）：①山石前练掌 ②金光随掌而起 ③云游 ④路遇不平以金光相援。诗词：寒暑不辍玄门功，护体金光初照空。云游但见不平事，一掌相援风雨中。。主题句：一杵降魔，护法三教——金刚怒目，亦是慈悲。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道青+护体金+晨光白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 金光护体, age 幼年, 玄门勤修·初窥神通; scene: 日夜修持不辍，山石前练掌，金光随掌而起；下山云游，路遇不平，以金光护身相援。; 寒暑不辍玄门功，护体金光初照空。云游但见不平事，一掌相援风雨中。; palette: 道青+护体金+晨光白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 降魔杵**
- 韦护，道法初成阶段·降魔杵成（少年，三教至宝·降魔杵出）。形象：韦护，身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。品性：道行天尊授其三教至宝降魔杵。神杵在手，万邪辟易，他终于修成梦寐以求的护法神通。。姿态：持降魔杵在山中一试，一杵砸下，山石崩裂，妖氛尽散；从此以杵护道。。服饰：道袍外罩护法战甲，手持降魔杵，法相庄严。。体型：身高约6头身，道袍外罩护法战甲。。衣物细节：护法战甲，手持降魔杵，法相庄严。。发型妆造：束发。。脸型五官：方脸庄严，剑眉，眼含威光，鼻梁挺，唇紧抿，下巴方正。。武器招式：降魔杵（三教至宝，万邪辟易）。。功法：降魔杵（三教至宝，万邪辟易）；护体金光大成。。功法表现：一杵砸出金光迸溅。。画面：构图：山中试杵，少年道人金甲持杵立于高崖，一杵砸出金光迸溅，妖氛四散；背景山川。色调：金甲+神杵金+山青。氛围：神通初成、护法威严。。台词："一杵在手，不是用来逞凶，是用来挡在众生面前，挡在邪恶面前。"。动作帧（动图）：①持降魔杵 ②山中一杵 ③山石崩裂妖氛尽散 ④从此以杵护道。诗词：三教至宝降魔杵，一挥万邪尽低头。山中一试山河动，护道之心自此稠。。主题句：一杵降魔，护法三教——金刚怒目，亦是慈悲。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：金甲+神杵金+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 降魔杵成, age 少年, 三教至宝·降魔杵出; scene: 持降魔杵在山中一试，一杵砸下，山石崩裂，妖氛尽散；从此以杵护道。; 三教至宝降魔杵，一挥万邪尽低头。山中一试山河动，护道之心自此稠。; palette: 金甲+神杵金+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 金刚怒目**
- 韦护，大劫淬炼阶段·护法先锋（少年，随军伐纣·千军护法）。形象：韦护，身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。品性：随姜子牙兴周伐纣，为西岐护法先锋。千军万马之中，降魔杵所向，护住军心与主帅。。姿态：阵前降魔杵一杵扫开敌军妖法；万军之中护住姜子牙法台，岿然不动。。服饰：护法金甲，手持降魔杵，身后金光如幢。。体型：身高约7头身，金甲护法，手持降魔杵。。衣物细节：护法金甲，降魔杵，身后金光如幢。。发型妆造：束发，怒目。。脸型五官：方脸，剑眉怒目，眼如金刚，鼻梁挺，唇线坚毅，下巴方正。。武器招式：降魔杵破邪降妖。。功法：降魔杵破邪降妖；护体金光罩军；金刚不坏。。功法表现：护体金光罩军，金刚不坏。。画面：构图：两军阵前，金甲护法手持降魔杵立于高台，金光如幢笼罩军阵，敌军妖法被金光挡住四散；背景战火旌旗。色调：护法金+战火红+玄黑。氛围：巍然、护法、一夫当关。。台词："冲锋陷阵有先锋，挡灾挡难有我。姜丞相，你的法台，我韦护守定了。"。动作帧（动图）：①阵前杵扫妖法 ②万军之中护姜子牙法台 ③金光如幢 ④岿然不动。诗词：西岐军帐护法身，千军万马一杵镇。金光如幢遮风雨，妖法邪术尽低头。。主题句：一杵降魔，护法三教——金刚怒目，亦是慈悲。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：护法金+战火红+玄黑。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 护法先锋, age 少年, 随军伐纣·千军护法; scene: 阵前降魔杵一杵扫开敌军妖法；万军之中护住姜子牙法台，岿然不动。; 西岐军帐护法身，千军万马一杵镇。金光如幢遮风雨，妖法邪术尽低头。; palette: 护法金+战火红+玄黑; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 护法大神**
- 韦护，封神登天阶段·肉身成圣（封神后，七人成圣·护法韦陀）。形象：韦护，身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。品性：伐纣功成，肉身成圣，位列护法韦陀。降魔杵化作无上护法，守护佛门与三界道场。。姿态：守护道场山门，降魔杵斜立身侧；凡有妖邪犯境，一杵即镇。。服饰：韦陀金身，降魔杵在手，宝相庄严，怒目含威。。体型：身高约7头身，韦陀金身，降魔杵在手。。衣物细节：韦陀金身，降魔杵，怒目含威。。发型妆造：束发，宝相庄严。。脸型五官：韦陀方脸，剑眉，怒目含威，鼻梁挺，唇抿，下巴方正，宝相。。武器招式：降魔杵镇一切邪。。功法：护法韦陀神位；降魔杵镇一切邪；金刚不坏身。。功法表现：佛光与幢幡于身后。。画面：构图：巍峨山门之前，护法韦陀金身持杵而立，降魔杵杵尖向下，身后佛光与幢幡；背景云山道场。色调：韦陀金+佛光白+山青。氛围：庄严、护法、宝相慈悲。。台词："成圣不是终点，护法才是。这世界总有黑夜，我便做黑夜里那道杵影。"。动作帧（动图）：①守护道场山门 ②降魔杵斜立 ③妖邪犯境一杵即镇 ④宝相庄严。诗词：肉身成圣列韦陀，降魔宝杵镇山河。护法三界常不寐，怒目金刚亦佛陀。。主题句：一杵降魔，护法三教——金刚怒目，亦是慈悲。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：韦陀金+佛光白+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 肉身成圣, age 封神后, 七人成圣·护法韦陀; scene: 守护道场山门，降魔杵斜立身侧；凡有妖邪犯境，一杵即镇。; 肉身成圣列韦陀，降魔宝杵镇山河。护法三界常不寐，怒目金刚亦佛陀。; palette: 韦陀金+佛光白+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 韦护尊者**
- 韦护，道果圆满阶段·杵化慈悲（神位长存，金刚怒目·降魔即慈悲）。形象：韦护，身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。品性：千百年护法，终于懂得：降魔杵降的从来不是魔，是魔性——护的从来是人心。金刚怒目之下，是一颗度化众生的慈悲心。。姿态：见妖邪被降后虔诚忏悔，会将其度化而非镇压；深夜入梦，护佑担惊受怕的孩子安眠。。服饰：韦陀便装，降魔杵化作一枚护身玉杵挂于胸前，金光内敛。。体型：身高约7头身，韦陀便装，金光内敛。。衣物细节：韦陀便装，降魔杵化作胸前玉杵，金光柔和。。发型妆造：束发。。脸型五官：韦陀方脸，剑眉舒展，怒目含慈，鼻梁挺，嘴角慈悲，下巴柔和。。武器招式：降魔即度化，一杵化作接引桥。。功法：降魔即度化；护法神通尽收，一念可护一方。。功法表现：金光柔和如月，照亮阶前安睡小兽。。画面：构图：暮色山门，护法韦陀背对落日而立，降魔杵已化作胸前温润玉杵，金光柔和如月，照亮阶前一只安睡的小兽；背景苍山如黛。色调：金身暖光+暮紫+山黛。氛围：慈悲、度化、怒目即菩萨。。台词："我一杵降的是恶念，不是生灵。魔若回头，我这一杵，便化作接引的桥。"。动作帧（动图）：①见妖邪忏悔度化 ②深夜入梦护孩童安眠 ③杵化玉杵挂胸前 ④暮色中温润如月。诗词：千年护法悟真章，怒目原来菩萨肠。杵落不是夺生念，度得回头即是光。。主题句：一杵降魔，护法三教——金刚怒目，亦是慈悲。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：金身暖光+暮紫+山黛。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 杵化慈悲, age 神位长存, 金刚怒目·降魔即慈悲; scene: 见妖邪被降后虔诚忏悔，会将其度化而非镇压；深夜入梦，护佑担惊受怕的孩子安眠。; 千年护法悟真章，怒目原来菩萨肠。杵落不是夺生念，度得回头即是光。; palette: 金身暖光+暮紫+山黛; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 妲己（`daji`） · 人生档案版

**灵胎初醒 · 妖狐之卵**
- 妲己，灵胎初醒阶段·轩辕坟狐（千年，轩辕坟中·九尾修炼）。形象：妲己，千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。品性：轩辕坟中修行千年的九尾狐，参得天地玄机，只差一步便可飞升。狐性本媚，却守着一片坟冢清修。。姿态：月下吐纳狐火，尾分九条，灵光缭绕；偶化人形，在坟前石阶上望月。。服饰：狐身时雪白九尾，月华披身；化人时白衣素裙，清丽如月。。体型：狐身体长如犬，雪白九尾如纱。。衣物细节：狐身雪白无饰，月华披身；化人白衣素裙。。发型妆造：狐身白毛；化人青丝如瀑。。脸型五官：狐脸清丽，细目含媚，鼻尖挺秀，嘴含微笑，下巴尖俏（化人时）；狐身细目温顺。。武器招式：无兵器，爪/尾。。功法：千年修为；九尾神通；狐火幻术。。功法表现：九尾流光，狐火幽幽。。画面：构图：轩辕坟古冢月夜，一只雪白九尾狐蹲坐坟头望月，九尾如纱，狐火幽幽；背景荒坟古柏。色调：月白+狐火幽蓝+墓灰。氛围：千年孤修、命定将变。。台词："修行千年，我本要踏月飞升。可那一声法旨，把我从天上，拉回了人间——也拉进了劫里。"。动作帧（动图）：①月下吐纳狐火 ②尾分九条 ③偶化人形望月 ④在坟前石阶蹲坐。诗词：轩辕坟外月华明，九尾千年修道成。只待一朝脱劫去，却闻法旨入红尘。。主题句：红颜祸国六百年，狐火焚尽成汤天下——妖也罢，命也罢，终究是一场劫。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：月白+狐火幽蓝+墓灰。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 轩辕坟狐, age 千年, 轩辕坟中·九尾修炼; scene: 月下吐纳狐火，尾分九条，灵光缭绕；偶化人形，在坟前石阶上望月。; 轩辕坟外月华明，九尾千年修道成。只待一朝脱劫去，却闻法旨入红尘。; palette: 月白+狐火幽蓝+墓灰; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年妖狐**
- 妲己，凡尘砺心阶段·附身入宫（千年，奉女娲命·附身妲己）。形象：妲己，千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。品性：奉女娲之命入商，附身冀州苏护之女妲己。本要祸乱商纣、助周伐纣，却在万人之上慢慢忘了来路。。姿态：附身之日，苏妲己眉眼一转，满城失色；入宫之后，一颦一笑，纣王再不上朝。。服饰：宫装云髻，蛾眉螓首，举手投足皆是倾国媚态。。体型：身高约7头身，宫装绝色，体态袅娜。。衣物细节：宫装云髻，蛾眉螓首。。发型妆造：云髻高绾，珠翠满头。。脸型五官：鹅蛋脸，眉目一转含媚，狐目勾魂，鼻梁挺，朱唇，下巴尖俏。。武器招式：无兵器，长袖轻舞。。功法：附身之术；魅惑之瞳；九尾幻术初施。。功法表现：魅惑之瞳，九尾幻术初施。。画面：构图：商宫大殿，绝色宫装女子回眸一笑，纣王神迷，满殿烛火摇曳；背景金碧宫阙。色调：宫红+金箔+妖媚紫。氛围：魅惑、权谋初起、迷失来路。。台词："我本奉命来断这殷商的气数。可这人间富贵，竟比狐火还烫——烫得我忘了自己是谁。"。动作帧（动图）：①附身之日眉眼一转 ②入宫 ③一颦一笑 ④纣王再不上朝。诗词：女娲法旨入宫墙，附得娇躯惑帝王。一颦一笑山河乱，忘了本是我狐妆。。主题句：红颜祸国六百年，狐火焚尽成汤天下——妖也罢，命也罢，终究是一场劫。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：宫红+金箔+妖媚紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 附身入宫, age 千年, 奉女娲命·附身妲己; scene: 附身之日，苏妲己眉眼一转，满城失色；入宫之后，一颦一笑，纣王再不上朝。; 女娲法旨入宫墙，附得娇躯惑帝王。一颦一笑山河乱，忘了本是我狐妆。; palette: 宫红+金箔+妖媚紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 入宫惑主**
- 妲己，道法初成阶段·魅惑君王（数载，权倾朝野·酒池肉林）。形象：妲己，千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。品性：自此权倾朝野，纣王言听计从。建酒池肉林、筑摘星楼，贪欲如狐火，烧了半个殷商。。姿态：摘星楼上长袖轻舞，九尾隐现；命人建酒池肉林，她笑看人间沉沦。。服饰：凤冠霞帔，宫装云髻，狐尾在裙下时隐时现。。体型：身高约7头身，凤冠霞帔，裙下九尾隐现。。衣物细节：凤冠霞帔，宫装云髻，狐尾裙下时隐时现。。发型妆造：云髻高绾，步摇轻颤。。脸型五官：鹅蛋脸，细眉入鬓，狐目迷离，鼻梁挺，朱唇带笑，下巴尖俏，眉间妖气。。武器招式：无兵器，长袖/狐尾。。功法：九尾幻术大盛；魅惑摄魂；狐火随心。。功法表现：九尾幻术大盛，狐火随心。。画面：构图：摘星楼上，盛装妖姬凭栏而立，裙下九尾隐现，脚下酒池肉林与朝歌万家灯火；背景夜空。色调：奢靡金+妖紫+夜黑。氛围：奢靡、倾覆、沉沦。。台词："大王要醉，我陪他醉；大王要荒，我陪他荒。这江山，我替你烧着玩——反正，是劫数。"。动作帧（动图）：①摘星楼长袖轻舞 ②九尾隐现 ③命建酒池肉林 ④笑看人间沉沦。诗词：摘星楼高接紫宸，酒池肉林夜夜春。狐火一烧山河沸，岂知身是劫中人。。主题句：红颜祸国六百年，狐火焚尽成汤天下——妖也罢，命也罢，终究是一场劫。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：奢靡金+妖紫+夜黑。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 魅惑君王, age 数载, 权倾朝野·酒池肉林; scene: 摘星楼上长袖轻舞，九尾隐现；命人建酒池肉林，她笑看人间沉沦。; 摘星楼高接紫宸，酒池肉林夜夜春。狐火一烧山河沸，岂知身是劫中人。; palette: 奢靡金+妖紫+夜黑; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 炮烙忠臣**
- 妲己，大劫淬炼阶段·罪孽滔天（数载，炮烙虿盆·残害忠良）。形象：妲己，千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。品性：设炮烙、筑虿盆，残害忠良无数。明知自己早已不是"奉旨祸商"那样简单，却已停不下来。。姿态：铜柱炮烙之下，忠臣化作焦骨；她袖手看那青烟，心底那点千年道心，忽明忽灭。。服饰：华服尽染妖气，眉间一丝阴鸷，狐性已占上风。。体型：身高约7头身，华服尽染妖气，眉间阴鸷。。衣物细节：华服妖气浓重，眉间一丝阴鸷。。发型妆造：云髻，眼角妖媚渐厉。。脸型五官：鹅蛋脸，细眉微蹙，狐目含煞，鼻梁挺，朱唇染血，下巴尖俏，阴鸷。。武器招式：无兵器，摄魂之眸。。功法：摄魂之眸摄人心魄；炮烙虿盆之酷刑由她构设。。功法表现：狐火化妖氛，炮烙虿盆阴气。。画面：构图：朝歌刑场，铜柱炮烙升腾青烟，妖姬华服立于高台，裙下九尾张扬，脚边虿盆毒蛇涌动；背景血色宫阙。色调：血绛+刑铁黑+妖紫。氛围：罪孽、沉沦、回不了头。。台词："我一边杀，一边也在问自己：这狐火，烧的是殷商，还是我自己？"。动作帧（动图）：①铜柱炮烙 ②袖手看青烟 ③虿盆毒蛇涌动 ④狐心忽明忽灭。诗词：炮烙虿盆血染袍，忠良白骨怨声高。狐心已入膏肓处，回头无岸恨滔滔。。主题句：红颜祸国六百年，狐火焚尽成汤天下——妖也罢，命也罢，终究是一场劫。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：血绛+刑铁黑+妖紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 罪孽滔天, age 数载, 炮烙虿盆·残害忠良; scene: 铜柱炮烙之下，忠臣化作焦骨；她袖手看那青烟，心底那点千年道心，忽明忽灭。; 炮烙虿盆血染袍，忠良白骨怨声高。狐心已入膏肓处，回头无岸恨滔滔。; palette: 血绛+刑铁黑+妖紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 九尾祸世**
- 妲己，封神登天阶段·城破显形（朝歌城破，原形毕露·姜尚斩狐）。形象：妲己，千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。品性：牧野之战，朝歌城破。姜子牙登台，斩狐于台下——千年修行，终作一场空。。姿态：大军压城，她褪去华服，现出九尾本相；被缚至法台，回头望了一眼那倾颓的宫阙。。服饰：华服尽褪，现出雪白九尾狐身，狐火明灭。。体型：狐身尽显，雪白九尾，被缚法台。。衣物细节：华服尽褪，雪白九尾狐身。。发型妆造：狐身白毛，狐火明灭。。脸型五官：狐脸渐显，媚态散尽，狐目圆睁，鼻尖微尖，嘴角颤抖，下巴尖，妖气散。。武器招式：无兵器，狐火护身。。功法：九尾本相尽出，狐火护身，却终究挡不住天命法旨。。功法表现：千年修为散尽，狐火渐熄。。画面：构图：朝歌城破之夜，法台上雪白九尾狐被缚，狐身低垂，远处宫阙倾颓、火光冲天；背景夜幕与军旗。色调：狐白+夜黑+火绛。氛围：穷途、宿命、一切成空。。台词："千年修行，一朝尽散。我不怨姜子牙——怨我接了那道法旨，却没守住那颗狐心。"。动作帧（动图）：①大军压城 ②褪华服现九尾本相 ③被缚法台 ④回望倾颓宫阙。诗词：牧野战鼓震朝歌，九尾狐身对法科。千年修为随风散，回首宫阙泪滂沱。。主题句：红颜祸国六百年，狐火焚尽成汤天下——妖也罢，命也罢，终究是一场劫。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：狐白+夜黑+火绛。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 城破显形, age 朝歌城破, 原形毕露·姜尚斩狐; scene: 大军压城，她褪去华服，现出九尾本相；被缚至法台，回头望了一眼那倾颓的宫阙。; 牧野战鼓震朝歌，九尾狐身对法科。千年修为随风散，回首宫阙泪滂沱。; palette: 狐白+夜黑+火绛; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 妖狐末路**
- 妲己，道果圆满阶段·狐心归山（劫后，妖气散尽·一缕魂归轩辕）。形象：妲己，千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。品性：一缕残魂飘回轩辕坟。祸福皆历，才知那千年的山野清修，才是唯一自在——妖也罢，命也罢，终究一场劫，如今劫散。。姿态：残魂落于轩辕坟头，化作一只再无法力的白狐，蜷在月下老柏根旁，听风声入眠。。服饰：再无法力的雪白小狐，无九尾，无狐火，唯余月下一身雪。。体型：小小白狐，无九尾无狐火。。衣物细节：再无法力的雪白小狐，一身雪。。发型妆造：一身白毛。。脸型五官：小狐脸，细目温顺，鼻尖湿润，嘴角安和，下巴圆，纯澈。。武器招式：无兵器。。功法：妖力尽散；唯余一丝清明，守着坟冢与故地。。功法表现：妖力尽散，唯月下一点清光。。画面：构图：轩辕坟月夜，一只小小的白狐蜷在古柏根下打盹，无尾无光，月光如水；背景荒坟与枯枝，安宁静谧。色调：月白+墓灰+老柏青。氛围：劫后、归寂、安宁。。台词："六百年，我做了一场大梦。梦醒，我还是轩辕坟外那只望月的狐狸——只是这次，我不再盼飞升了。"。动作帧（动图）：①残魂飘回轩辕坟 ②蜷在古柏根下 ③望月 ④听风声入眠。诗词：轩辕坟外又逢秋，一缕残魂作雪狐。祸福历尽心归处，月照荒坟万事休。。主题句：红颜祸国六百年，狐火焚尽成汤天下——妖也罢，命也罢，终究是一场劫。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：月白+墓灰+老柏青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 狐心归山, age 劫后, 妖气散尽·一缕魂归轩辕; scene: 残魂落于轩辕坟头，化作一只再无法力的白狐，蜷在月下老柏根旁，听风声入眠。; 轩辕坟外又逢秋，一缕残魂作雪狐。祸福历尽心归处，月照荒坟万事休。; palette: 月白+墓灰+老柏青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 申公豹（`shen_gongbao`） · 人生档案版

**灵胎初醒 · 邪道灵光**
- 申公豹，灵胎初醒阶段·玉虚弟子（少年，昆仑山·阐教高徒）。形象：申公豹，黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。品性：昆仑山玉虚宫元始天尊弟子，与姜子牙同门。资质出众，口才无双，却心高气傲，总觉得自己该是封神的主事人。。姿态：玉虚宫中论道，无人说得过他三寸舌；见姜子牙得师门看重，心中暗暗不平。。服饰：昆仑道袍，眉目英挺，腰间悬剑，一副得道高士模样。。体型：身高约7头身，昆仑道袍，英挺。。衣物细节：昆仑道袍，腰间悬剑。。发型妆造：束发，眉目英挺。。脸型五官：长脸，剑眉微挑，目如鹰隼，鼻梁高挺，薄唇善辩，下巴尖，傲气。。武器招式：腰间佩剑未出鞘。。功法：阐教道法精深；辩才无碍，舌灿莲花。。功法表现：谈经论道，舌灿莲花。。画面：构图：昆仑山玉虚宫论道场，意气风发的道人立于众仙之中，衣袖生风，侃侃而谈，身后云海仙山。色调：道青+云白+玉虚金。氛围：傲气、口才、暗生不平。。台词："姜师兄，你我同门，凭什么执封神榜的是你？这天下的英雄，谁不认得我申公豹这三寸舌？"。动作帧（动图）：①论道舌灿莲花 ②众仙叹服 ③见姜子牙得看重暗生不平 ④拂袖。诗词：昆仑高处羽衣轻，谈经论道舌灿星。自诩天下英雄主，哪知身在劫中行。。主题句：道友请留步——一张嘴说尽天下豪杰，却也误了一身道行。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道青+云白+玉虚金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 玉虚弟子, age 少年, 昆仑山·阐教高徒; scene: 玉虚宫中论道，无人说得过他三寸舌；见姜子牙得师门看重，心中暗暗不平。; 昆仑高处羽衣轻，谈经论道舌灿星。自诩天下英雄主，哪知身在劫中行。; palette: 道青+云白+玉虚金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 阐教入门**
- 申公豹，凡尘砺心阶段·被逐教门（青年，心术不正·逐出阐教）。形象：申公豹，黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。品性：屡次离间同门、卖弄心机，元始天尊察其心术不正，将其逐出昆仑。他不思悔改，反生怨怼。。姿态：被逐出山门时一步三回头，眼底怨毒渐浓；下山后誓言要让姜子牙难堪。。服饰：道袍已旧，眉宇间多了一丝阴鸷。。体型：身高约7头身，道袍已旧，眉宇阴鸷。。衣物细节：道袍已旧，剑在手。。发型妆造：束发，眉间阴鸷。。脸型五官：长脸，剑眉含怨，目带阴鸷，鼻梁高挺，薄唇紧抿，下巴尖。。武器招式：剑未出鞘。。功法：道法仍在，心术已歪；游说之术愈发老练。。功法表现：怨气渐生。。画面：构图：昆仑山门之下，被逐的道人抱剑回望，云海在身后翻涌，山门紧闭；背景苍茫群山。色调：道青+暮云灰+一点暗红。氛围：怨怼、执念、堕落前兆。。台词："逐我出教？好！我申公豹没了师门，还有这张嘴——我偏要与你们斗到底！"。动作帧（动图）：①被逐出山门 ②一步三回头 ③眼底怨毒渐浓 ④下山誓言。诗词：一纸逐书断昆仑，怨根从此入心门。三寸舌横天下走，誓教同门负怨深。。主题句：道友请留步——一张嘴说尽天下豪杰，却也误了一身道行。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道青+暮云灰+一点暗红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 被逐教门, age 青年, 心术不正·逐出阐教; scene: 被逐出山门时一步三回头，眼底怨毒渐浓；下山后誓言要让姜子牙难堪。; 一纸逐书断昆仑，怨根从此入心门。三寸舌横天下走，誓教同门负怨深。; palette: 道青+暮云灰+一点暗红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 叛教下山**
- 申公豹，道法初成阶段·游说万仙（青年，转投截教·道友请留步）。形象：申公豹，黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。品性：转投截教，凭三寸不烂之舌四处游说，一句"道友请留步"，说动无数截教仙人下山助纣。。姿态：云游四海，逢仙便挽袖攀谈；所到之处，必有人被他鼓动得热血上头。。服饰：截教黑道袍，身负双剑，笑意盈盈。。体型：身高约7头身，截教黑道袍。。衣物细节：截教黑道袍，身负双剑。。发型妆造：束发，笑意盈盈。。脸型五官：长脸，剑眉，目带狡黠，鼻梁高挺，薄唇含笑，下巴尖，笑意盈盈。。武器招式：双剑在背。。功法：游说之术登峰造极；三昧真火初窥；飞头术。。功法表现：三寸舌蛊惑，妖风随行。。画面：构图：云雾山头，黑袍道人拦住一位仙人笑谈，衣袖翻飞，舌灿莲花，身后各路仙影赴约而来；背景缥缈云海。色调：截教黑+诡紫+山青。氛围：蛊惑、游说、风云际会。。台词："道友请留步！你这一身本事，替那昏君卖命，岂不屈才？不如随我去，共襄"大业"！"。动作帧（动图）：①云游四海 ②逢仙挽袖攀谈 ③"道友请留步" ④鼓动仙人下山。诗词：道友留步一声呼，说动群仙出洞府。三寸舌间山河动，半为豪气半糊涂。。主题句：道友请留步——一张嘴说尽天下豪杰，却也误了一身道行。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：截教黑+诡紫+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 游说万仙, age 青年, 转投截教·道友请留步; scene: 云游四海，逢仙便挽袖攀谈；所到之处，必有人被他鼓动得热血上头。; 道友留步一声呼，说动群仙出洞府。三寸舌间山河动，半为豪气半糊涂。; palette: 截教黑+诡紫+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 请道友下山**
- 申公豹，大劫淬炼阶段·助纣为虐（中年，屡斗姜尚·逆天而行）。形象：申公豹，黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。品性：助商伐周，屡与姜子牙斗法。机关算尽，却一次次落败，越斗越深，越陷越不能回头。。姿态：阵前鼓动截教设下十绝阵、九曲黄河阵；败后又去搬援兵，从不认输。。服饰：黑袍猎猎，双剑在手，面容阴鸷。。体型：身高约7头身，黑袍猎猎。。衣物细节：黑袍猎猎，双剑在手。。发型妆造：束发，面容阴鸷。。脸型五官：长脸，剑眉倒竖，目含戾气，鼻梁高挺，薄唇冷峭，下巴尖。。武器招式：飞头术/三昧真火。。功法：飞头术、三昧真火；截教旁门道术；调动万仙。。功法表现：十绝阵煞气腾腾。。画面：构图：两军阵前，黑袍道人负手立于云端，身后截教众仙云集，脚下十绝阵煞气腾腾；背景战云密布。色调：黑袍+阵煞紫黑+战云绛。氛围：逆天、执迷、众仙聚祸。。台词："姜子牙，你赢我一局，赢不了我的嘴！只要我还能说话，这天下就没有我请不来的援兵！"。动作帧（动图）：①阵前鼓动截教 ②设十绝阵 ③败后去搬援兵 ④从不认输。诗词：十绝阵起黄河横，屡败屡战心不惊。逆天而行浑不觉，只向同门较输赢。。主题句：道友请留步——一张嘴说尽天下豪杰，却也误了一身道行。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黑袍+阵煞紫黑+战云绛。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 助纣为虐, age 中年, 屡斗姜尚·逆天而行; scene: 阵前鼓动截教设下十绝阵、九曲黄河阵；败后又去搬援兵，从不认输。; 十绝阵起黄河横，屡败屡战心不惊。逆天而行浑不觉，只向同门较输赢。; palette: 黑袍+阵煞紫黑+战云绛; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 五岳三山**
- 申公豹，封神登天阶段·镇压北海眼（封神前，被擒·塞入北海眼）。形象：申公豹，黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。品性：最终被元始天尊擒住，塞入北海眼堵海。一身道行，落到如此下场，才算醒悟——却已迟了。。姿态：被擒时仍不服，大声叫骂；被压入海眼刹那，望着昆仑方向，忽然沉默。。服饰：道袍零落，双剑尽失，被压于北海眼中。。体型：身高约7头身，道袍零落。。衣物细节：道袍零落，双剑尽失。。发型妆造：散发。。脸型五官：长脸，眉目颓然，目含悔恨，鼻梁高挺，薄唇颤抖，下巴尖。。武器招式：无兵器，道行尽封。。功法：道行尽封；唯余一颗悔恨之心。。功法表现：海眼漩涡吞身。。画面：构图：北海波涛之中，海眼漩涡吞噬着黑袍道人的残影，他回头望向昆仑方向，眼神悔恨；背景漆黑海天。色调：海渊黑+浪花白+一点悔恨金。氛围：穷途、悔悟、迟来的清醒。。台词："我这一生，都在跟人斗嘴、斗法、斗气。斗到最后，把自己斗进了海眼里。姜师兄，我错了……"。动作帧（动图）：①被擒仍叫骂 ②压入海眼 ③望昆仑方向 ④忽然沉默。诗词：北海眼中镇此身，滔天浊浪锁星辰。三寸舌终成铁索，悔不当初做真人。。主题句：道友请留步——一张嘴说尽天下豪杰，却也误了一身道行。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：海渊黑+浪花白+一点悔恨金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 镇压北海眼, age 封神前, 被擒·塞入北海眼; scene: 被擒时仍不服，大声叫骂；被压入海眼刹那，望着昆仑方向，忽然沉默。; 北海眼中镇此身，滔天浊浪锁星辰。三寸舌终成铁索，悔不当初做真人。; palette: 海渊黑+浪花白+一点悔恨金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 封神台末**
- 申公豹，道果圆满阶段·分水将军（封神后，执掌海眼·以罪镇海）。形象：申公豹，黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。品性：封神后受封分水将军，镇守北海眼。那张说尽天下英雄的嘴，从此只用来巡查海疆、安抚水族。。姿态：巡行北海，止波平浪；路见水族受困，会淡淡说一句"留步，我送你们一程"。。服饰：分水将军神袍，周身水灵环绕，眉宇间戾气尽消，只剩看透后的淡泊。。体型：身高约7头身，分水将军神袍，戾气尽消。。衣物细节：分水将军神袍，周身水灵。。发型妆造：束发，眉宇淡泊。。脸型五官：长脸，眉目淡泊，目含平静，鼻梁高挺，薄唇温和，下巴圆润，戾气尽消。。武器招式：无兵器，以法镇海。。功法：分水将军神位；执掌海眼；以罪身护海疆。。功法表现：浪头归平，波光银白。。画面：构图：北海之上，分水将军神袍立于浪头，衣袖轻拂，浪涛归平，身后水族影影绰绰相送；背景海天澄澈。色调：海蓝+神袍青+波光银。氛围：赎罪、镇海、看透归淡。。台词："从前我请人"留步"，是请他们入劫。如今我请水族"留步"，是请他们归家。一句留步，两种因果。"。动作帧（动图）：①巡行北海 ②止波平浪 ③路见水族受困 ④淡淡说"留步，送你们一程"。诗词：分水将军镇海门，涛声日夜洗尘魂。当年舌上风云乱，今作巡波一叩询。。主题句：道友请留步——一张嘴说尽天下豪杰，却也误了一身道行。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：海蓝+神袍青+波光银。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 分水将军, age 封神后, 执掌海眼·以罪镇海; scene: 巡行北海，止波平浪；路见水族受困，会淡淡说一句"留步，我送你们一程"。; 分水将军镇海门，涛声日夜洗尘魂。当年舌上风云乱，今作巡波一叩询。; palette: 海蓝+神袍青+波光银; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 孙悟空（`sun_wukong`） · 人生档案版

**灵胎初醒 · 灵根育孕**
- 孙悟空，灵胎初醒阶段·石猴出世（出生，花果山仙石·目运金光）。形象：孙悟空，通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。品性：花果山顶仙石孕育，受日月精华，一朝迸裂，化作石猴。出世即目运金光，射冲斗府，惊动天庭。。姿态：在花果山称王，率群猴嬉游；胆大包天，闯水帘洞，被尊为美猴王。。服饰：无衣无饰，通体金毛，一双火眼初露精光。。体型：身高约3头身，瘦小却筋骨强健，通体金毛，蜷坐跃起之势。。衣物细节：无衣，通体金黄短毛，前胸肚腹淡金，毛色在光下泛细碎金光。。发型妆造：头顶一撮乱毛自然蓬起，双耳尖立，脸侧金色颊毛。。脸型五官：脸型尖桃，眉如倒竖，金瞳圆睁，雷公嘴略凸，下巴尖瘦，颊毛蓬松。。武器招式：无武器，十指如钩，抓、挠、攀、跃皆利。。功法：天生灵猴，筋骨非凡；目运金光，惊动三界。。功法表现：双目睁开金光射斗，周身祥云隐隐随行。。画面：构图：花果山巅，仙石迸裂，石猴跃出，金光四射，身后飞瀑流泉、群猴朝拜；背景碧海青山。色调：石青+金毛亮金+花果红。氛围：天地初生、意气昂扬。。台词："天地生我孙悟空，可不是让我做个山大王就算了！我要学长生，要跳出这天地去！"。动作帧（动图）：①仙石迸裂、金光四射 ②石猴从中跃出 ③落地连翻两个筋斗 ④抓耳挠腮环顾四野。诗词：仙石迸裂日月惊，金光射斗动天庭。水帘洞中称大圣，跳出乾坤第一程。。主题句：跳出三界外，不在五行中——只做自己的齐天大圣。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：石青+金毛亮金+花果红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 石猴出世, age 出生, 花果山仙石·目运金光; scene: 在花果山称王，率群猴嬉游；胆大包天，闯水帘洞，被尊为美猴王。; 仙石迸裂日月惊，金光射斗动天庭。水帘洞中称大圣，跳出乾坤第一程。; palette: 石青+金毛亮金+花果红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 花果山王**
- 孙悟空，凡尘砺心阶段·拜师学艺（少年，灵台方寸山·菩提学艺）。形象：孙悟空，通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。品性：漂洋过海，遍访仙山，终在灵台方寸山拜菩提祖师。虽顽劣，却最肯下苦功，暗中学得一身本领。。姿态：在祖师处砍柴挑水七年，听闻大道便记在心上；夜半偷学七十二变，被祖师嗔怪又爱惜。。服饰：粗布道袍，毛脸雷公嘴，眼神灵动。。体型：身高约4头身，瘦长灵活，驼背微弓（猴态），四肢修长。。衣物细节：粗布青道袍，腰束麻绳，袖口挽起，衣料素净无纹。。发型妆造：金毛束成单髻，以木簪别住，脸侧颊毛仍在。。脸型五官：脸型尖桃，细长眉微挑，火眼初开，雷公嘴微收，下巴尖，颊毛半掩。。武器招式：无兵器，学拳脚与变化，翻腾跳跃如飞。。功法：七十二变；筋斗云（一个筋斗十万八千里）；长生妙法。。功法表现：变化之术初成，化作飞鸟/树木；一个筋斗翻上半空。。画面：构图：灵台方寸山月夜，毛脸猴道的少年在三星洞前演练变化，身形在云中腾挪，身后祖师洞府古松；背景月光云海。色调：道青+月白+金毛。氛围：苦学、机灵、初得神通。。台词："祖师，学不会变我就不走了！这猴子生来就是要跳出三界，长生不死！"。动作帧（动图）：①三星洞前打坐吐纳 ②起身舞袖演变化 ③摇身化作一棵大树 ④变回原形挠头憨笑。诗词：灵台山下学神通，七十二变指掌中。一个筋斗十万八，偷天换日不居功。。主题句：跳出三界外，不在五行中——只做自己的齐天大圣。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道青+月白+金毛。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 拜师学艺, age 少年, 灵台方寸山·菩提学艺; scene: 在祖师处砍柴挑水七年，听闻大道便记在心上；夜半偷学七十二变，被祖师嗔怪又爱惜。; 灵台山下学神通，七十二变指掌中。一个筋斗十万八，偷天换日不居功。; palette: 道青+月白+金毛; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 筋斗云**
- 孙悟空，道法初成阶段·龙宫取宝（青年，金箍棒·勾销生死簿）。形象：孙悟空，通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。品性：回花果山后，去东海龙宫寻得如意金箍棒，又大闹地府勾销猴类生死簿——"我命由我不由天"。。姿态：东海龙宫一试金箍棒，重一万三千五百斤，随手如意；地府之中，一笔勾去生死簿上猴属姓名。。服饰：金箍棒在手，虎皮裙初披，凤翅紫金冠戴起。。体型：身高约5头身，精壮起来，毛色金亮，昂首挺胸。。衣物细节：虎皮裙黄底黑纹，金箍棒在手，周身海气缭绕。。发型妆造：凤翅紫金冠初戴，金毛向后梳，气势初显。。脸型五官：脸型尖桃，眉峰扬起，金瞳精光，雷公嘴含笑，下巴微抬，颊毛金亮。。武器招式：金箍棒：两端金箍中间红，可大可小；起手式抡棒横扫。。功法：金箍棒如意变化；七十二变愈精；地府勾名，跳出轮回。。功法表现：棒起海波翻涌；地府勾名时笔锋带金光。。画面：构图：东海龙宫水波荡漾，石猴握着通体金光的金箍棒立于殿前，身后虾兵蟹将退避；背景珊瑚宝殿。色调：海蓝+棒光金+珊瑚红。氛围：得宝、豪气、睥睨四海。。台词："这棒子称手！从此我孙悟空，天不收、地不管，生死簿上再无我名！"。动作帧（动图）：①劈浪入海前行 ②龙宫前立定 ③金箍棒入手掂量 ④随手一旋，棒影如轮。诗词：东海龙宫取棒来，一万三千五斤开。地府勾销生死簿，我命从今不天裁。。主题句：跳出三界外，不在五行中——只做自己的齐天大圣。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：海蓝+棒光金+珊瑚红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 龙宫取宝, age 青年, 金箍棒·勾销生死簿; scene: 东海龙宫一试金箍棒，重一万三千五百斤，随手如意；地府之中，一笔勾去生死簿上猴属姓名。; 东海龙宫取棒来，一万三千五斤开。地府勾销生死簿，我命从今不天裁。; palette: 海蓝+棒光金+珊瑚红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 大闹天宫**
- 孙悟空，大劫淬炼阶段·大闹天宫（青年，齐天大圣·压五行山）。形象：孙悟空，通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。品性：天庭招安封他"弼马温"，他嫌官小反下天；自封齐天大圣，搅了蟠桃会，大闹天宫，十万天兵无人能挡。。姿态：蟠桃会上吃尽仙桃仙酒，酒醉闯凌霄殿；被如来以五行山镇压，一压五百年。。服饰：凤翅紫金冠、锁子黄金甲、藕丝步云履，金箍棒在手，威震天庭。。体型：身高约6头身，全盛之姿，金毛猎猎，昂然睥睨。。衣物细节：锁子黄金甲金鳞细密，虎皮裙黄底黑纹，藕丝步云履，红披风猎猎。。发型妆造：凤翅紫金冠，两根雉翎高挑，金毛披散，眉间红光。。脸型五官：脸型尖桃，眉飞入鬓，金瞳如电，雷公嘴龇出獠牙，下巴坚毅，眉眼红晕。。武器招式：金箍棒舞成棒影如墙，横扫千军、劈山裂地。。功法：七十二变、筋斗云、金箍棒；一身神通撼天动地。。功法表现：神通撼天动地，云雾随棒翻涌，双目金光如电。。画面：构图：凌霄殿前，金甲猴王挥舞金箍棒，天兵天将溃散，金冠猎猎；背景天宫金阙云海翻涌。色调：锁金甲+天宫金+战云紫。氛围：惊天动地、狂放不羁。。台词："皇帝轮流做，明年到我家！这天宫我来得、闹得，也坐得！"。动作帧（动图）：①拔棒冲天而起 ②棒指凌霄殿 ③横扫天兵如卷席 ④仰天长啸。诗词：齐天大圣反天庭，十万天兵扫若尘。玉帝也惧猴王勇，五行山下压猴身。。主题句：跳出三界外，不在五行中——只做自己的齐天大圣。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：锁金甲+天宫金+战云紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 大闹天宫, age 青年, 齐天大圣·压五行山; scene: 蟠桃会上吃尽仙桃仙酒，酒醉闯凌霄殿；被如来以五行山镇压，一压五百年。; 齐天大圣反天庭，十万天兵扫若尘。玉帝也惧猴王勇，五行山下压猴身。; palette: 锁金甲+天宫金+战云紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 五百年劫**
- 孙悟空，封神登天阶段·西天取经（五百年后，保唐僧·齐天大圣护师）。形象：孙悟空，通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。品性：被如来点化，保唐僧西天取经。戴上金箍，收了性子，一路降妖除魔、护师西行——从"齐天"到"护师"，心性渐定。。姿态：一路打杀妖魔鬼怪，火眼金睛识破百般变化；遇妖怪伤师，一棒便是雷霆。。服饰：锦布直裰，虎皮裙，金箍棒在手，火眼金睛如电。。体型：身高约6头身，身形矫健，略见风霜，脊背挺直护师而行。。衣物细节：锦布直裰，虎皮裙仍在，腰系紧带，风尘仆仆。。发型妆造：金毛束起，戴紧箍，脸侧颊毛，火眼金睛。。脸型五官：脸型尖桃，眉间略敛，火眼金睛，雷公嘴紧抿，下巴沉稳，风霜几缕。。武器招式：金箍棒：起手式藏于肘后，扫、点、劈三连，棒随身走。。功法：七十二变、筋斗云、金箍棒；火眼金睛识妖。。功法表现：火眼金睛识妖气，一道金光自眼而出；棒落时金光迸溅。。画面：构图：苍山古道，猴王持棒在前开路，火眼金睛放出神光，身后唐僧白马隐约，前方妖气隐隐；背景西行山川。色调：行者黄+袈裟红+山青。氛围：护师西行、历经艰险。。台词："俺老孙保师父取经，不是为了成佛，是为了一诺千金——说了护你到西天，少一步都不算数。"。动作帧（动图）：①火眼金睛扫视前路 ②金箍棒横在肩头 ③闻妖风腾身而起 ④一棒打散妖雾。诗词：五行山下悟前尘，保得唐僧踏劫轮。一棒扫开妖雾散，火眼金睛辨假真。。主题句：跳出三界外，不在五行中——只做自己的齐天大圣。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：行者黄+袈裟红+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 西天取经, age 五百年后, 保唐僧·齐天大圣护师; scene: 一路打杀妖魔鬼怪，火眼金睛识破百般变化；遇妖怪伤师，一棒便是雷霆。; 五行山下悟前尘，保得唐僧踏劫轮。一棒扫开妖雾散，火眼金睛辨假真。; palette: 行者黄+袈裟红+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 斗战胜佛**
- 孙悟空，道果圆满阶段·斗战胜佛（成佛后，一念降魔·持棒护三界）。形象：孙悟空，通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。品性：九九八十一难历尽，受封斗战胜佛。既佛也斗——金箍棒化金刚降魔，一念护持三界；跳脱如昨，仍是那只只做自己的猴。。姿态：云端持棒，降魔印随手而结；见人间不平，金箍棒一点即平；闲时回花果山看猴儿们闹。。服饰：锦斓袈裟、毗卢帽，金瞳威仪，金箍棒竖立身侧。。体型：身高约6头身，端坐莲台法相庄严，身形凝实，威而不怒。。衣物细节：锦斓袈裟红底金格纹，毗卢帽金线，金箍棒竖立身侧。。发型妆造：毗卢帽覆顶，金毛隐于帽下，金瞳威仪。。脸型五官：脸型端方，眉目慈悲中带威，金瞳内敛，雷公嘴微阖，下巴圆满，颊毛隐于帽下。。武器招式：金箍棒竖立拄地，化作降魔金刚杵；左手结降魔印。。功法：斗战胜佛神位；一念降魔；金箍棒化金刚杵；不坏金刚身。。功法表现：佛光法相自背后放出，金箍棒顶佛光流转，一念降魔。。画面：构图：破云金光法相之中，斗战胜佛端坐莲花，左手结降魔印，右手拄金箍棒，金瞳威仪，身后佛光双环；背景赤金法相。色调：赤金+袈裟红+佛光金。氛围：既佛既斗、降魔护世。。台词："做了佛，俺老孙还是那只猴。妖该劈时劈得明明白白，人该护时一个都不许漏——这就是俺的道。"。动作帧（动图）：①盘膝莲台 ②左掌结降魔印 ③右手拄金箍棒 ④金瞳微睁，佛光遍洒。诗词：斗战胜佛登莲台，金箍棒化降魔来。跳出三界心仍在，只做自己那猴怀。。主题句：跳出三界外，不在五行中——只做自己的齐天大圣。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：赤金+袈裟红+佛光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 斗战胜佛, age 成佛后, 一念降魔·持棒护三界; scene: 云端持棒，降魔印随手而结；见人间不平，金箍棒一点即平；闲时回花果山看猴儿们闹。; 斗战胜佛登莲台，金箍棒化降魔来。跳出三界心仍在，只做自己那猴怀。; palette: 赤金+袈裟红+佛光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 吕洞宾（`lv_dongbin`） · 人生档案版

**灵胎初醒 · 黄粱书生**
- 吕洞宾，灵胎初醒阶段·儒生进士（唐末，河中府·饱读诗书）。形象：吕洞宾，白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。品性：河中府儒生，唐末进士，才华横溢，抱负满怀，一心要济世安邦。。姿态：晨起读书，暮写策论；赴京赶考，一路意气风发。。服饰：青衫儒巾，佩书箱，眉目俊逸。。体型：身高约7头身，青衫儒生，清俊修长。。衣物细节：青衫儒巾，腰佩玉佩，书箱随身。。发型妆造：儒巾束发，发髻齐整。。脸型五官：俊逸面容，剑眉星目，鼻梁挺，嘴角含笑，下巴微尖。。武器招式：无兵器，笔墨为剑。。功法：满腹经纶，善诗文书画。。功法表现：无神力，文气盈身。。画面：构图：河中府书房，青衫儒生伏案苦读，烛火摇曳，窗外月色，背景书卷满架。色调：青衫+暖烛金+月白。氛围：意气风发、书生意气。。台词："十年寒窗，就为这金榜题名。我吕岩，要凭一肚子学问，换一个朗朗乾坤。"。动作帧（动图）：①伏案读书 ②掩卷沉思 ③挥笔写策论 ④负笈出门赴考。诗词：河中才子负笈行，一举成名天下惊。青衫仗剑书生气，欲把文章换太平。。主题句：黄粱一梦醒功名，仗剑天地任逍遥——度人先度己，成仙先成人。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青衫+暖烛金+月白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 儒生进士, age 唐末, 河中府·饱读诗书; scene: 晨起读书，暮写策论；赴京赶考，一路意气风发。; 河中才子负笈行，一举成名天下惊。青衫仗剑书生气，欲把文章换太平。; palette: 青衫+暖烛金+月白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 弃儒从道**
- 吕洞宾，凡尘砺心阶段·黄粱一梦（中年，邯郸客店·功名如梦）。形象：吕洞宾，白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。品性：邯郸客店，汉钟离以黄粱一梦点化——梦中历尽金榜题名、高官厚禄、宦海浮沉、潦倒落魄，醒来黄粱未熟。。姿态：梦中从状元到宰相再到阶下囚，荣华与落魄尽历；梦醒惊觉，功名不过黄粱一场。。服饰：梦中换过官袍囚衣，醒来仍是青衫。。体型：身高约7头身，倚枕卧于客店，梦中身姿变换。。衣物细节：青衫，枕边书箱；梦中官袍囚衣虚影。。发型妆造：儒巾散落，梦中束发变化。。脸型五官：梦中神色从狂喜到悲凉，醒时目光澄明。。武器招式：无兵器。。功法：勘破功名之幻；道心初动。。功法表现：黄粱一梦浮影，镜花水月。。画面：构图：邯郸客店，青衫儒生倚枕而卧，头顶浮起半生梦幻（状元/宰相/囚徒三幕虚影），炉上黄粱正熟。色调：梦幻金+客店褐+梦魇青。氛围：大梦初醒、勘破浮华。。台词："一枕黄粱，历尽半生。原来这功名富贵，不过是锅里的黄粱，一熟即散。"。动作帧（动图）：①卧枕入梦 ②梦中冠盖如云 ③梦中潦倒落魄 ④惊醒时黄粱未熟。诗词：邯郸道上梦黄粱，半世荣华一枕凉。醒来米熟身犹在，始信浮生梦一场。。主题句：黄粱一梦醒功名，仗剑天地任逍遥——度人先度己，成仙先成人。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：梦幻金+客店褐+梦魇青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 黄粱一梦, age 中年, 邯郸客店·功名如梦; scene: 梦中从状元到宰相再到阶下囚，荣华与落魄尽历；梦醒惊觉，功名不过黄粱一场。; 邯郸道上梦黄粱，半世荣华一枕凉。醒来米熟身犹在，始信浮生梦一场。; palette: 梦幻金+客店褐+梦魇青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 御剑飞行**
- 吕洞宾，道法初成阶段·拜师学道（中年，汉钟离收徒·纯阳剑法）。形象：吕洞宾，白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。品性：弃儒从道，拜汉钟离为师，学道习剑。师父授纯阳剑法与金丹大道。。姿态：山中练剑，剑光如虹；习吐纳导引，炼丹养气。。服饰：青布道袍，腰悬酒葫芦，背负纯阳剑。。体型：身高约7头身，青袍道者，身形矫健。。衣物细节：青布道袍，腰悬酒葫芦，背负纯阳剑。。发型妆造：道髻簪木，眉目俊逸。。脸型五官：俊逸面庞，眉目疏朗，眼神含剑意。。武器招式：纯阳剑（初学）：起手撩剑、横扫、回旋。。功法：纯阳剑法初成；吐纳导引，金丹大道筑基。。功法表现：剑气如虹，带起落叶。。画面：构图：终南道观，青袍道者持剑练于松前，剑光如虹，酒葫芦腰悬，身后古松云海。色调：道青+剑光白+葫芦棕。氛围：弃儒从道、剑心初成。。台词："师父，剑要快，心要慢。我且以这柄纯阳剑，斩尽天下不平事。"。动作帧（动图）：①松前起剑 ②剑光如虹横扫 ③收剑饮酒 ④负剑立定望云海。诗词：弃儒从道入玄门，纯阳剑起斩红尘。葫芦盛酒丹炉火，一念逍遥一念真。。主题句：黄粱一梦醒功名，仗剑天地任逍遥——度人先度己，成仙先成人。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道青+剑光白+葫芦棕。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 拜师学道, age 中年, 汉钟离收徒·纯阳剑法; scene: 山中练剑，剑光如虹；习吐纳导引，炼丹养气。; 弃儒从道入玄门，纯阳剑起斩红尘。葫芦盛酒丹炉火，一念逍遥一念真。; palette: 道青+剑光白+葫芦棕; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 五雷天心**
- 吕洞宾，大劫淬炼阶段·仗剑行侠（中年，十试淬炼·游戏人间）。形象：吕洞宾，白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。品性：剑术大成，仗剑行天下，游戏人间度有缘。历经十试，剑心淬炼，越发纯粹。。姿态：酒葫芦一倾，纯阳剑出鞘，飞剑千里斩妖；遇不平事拔剑相助，遇有缘人点化度之。。服饰：白衣剑仙，逍遥巾束发，腰悬酒葫芦，纯阳剑在手。。体型：身高约7头身，白衣剑仙，御剑飞行身姿飘逸。。衣物细节：白衣剑仙，逍遥巾，酒葫芦，纯阳剑。。发型妆造：逍遥巾束发，衣袂飘飘。。脸型五官：俊逸面容，眉目含笑，眼神温润又带剑意。。武器招式：纯阳剑：御剑飞行、飞剑斩妖、剑气纵横。。功法：御剑飞行；飞剑千里；纯阳剑法大成；点石成金。。功法表现：剑气如虹，剑光千里。。画面：构图：江湖山水间，白衣剑仙御剑飞行，身后剑气如虹，下方妖雾被斩开，背景大好河山。色调：白衣+剑光银+江山青。氛围：快意恩仇、游戏人间。。台词："这世道，妖要斩，人也要度。我吕岩这柄剑，斩妖亦斩心魔。"。动作帧（动图）：①御剑腾空 ②剑指妖雾 ③一剑劈开妖氛 ④收剑立于剑上饮酒。诗词：白衣仗剑走天涯，酒葫芦里盛晚霞。一剑斩妖千里外，闲来点石化金砂。。主题句：黄粱一梦醒功名，仗剑天地任逍遥——度人先度己，成仙先成人。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：白衣+剑光银+江山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 仗剑行侠, age 中年, 十试淬炼·游戏人间; scene: 酒葫芦一倾，纯阳剑出鞘，飞剑千里斩妖；遇不平事拔剑相助，遇有缘人点化度之。; 白衣仗剑走天涯，酒葫芦里盛晚霞。一剑斩妖千里外，闲来点石化金砂。; palette: 白衣+剑光银+江山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 点石成金**
- 吕洞宾，封神登天阶段·八仙之首（中年，八仙过海·纯阳吕祖）。形象：吕洞宾，白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。品性：八仙过海各显神通，他持纯阳剑踏浪而行，万剑归宗，位列八仙之首。。姿态：海面踏剑而行，剑气化作万千剑影；为众仙开路，谈笑间已渡沧海。。服饰：白衣道袍，纯阳剑在手，酒葫芦挂腰，八仙之首气度。。体型：身高约7头身，白衣道袍，踏剑凌波。。衣物细节：白衣道袍，纯阳剑，酒葫芦，八仙之首气度。。发型妆造：道髻逍遥巾，衣袂猎猎。。脸型五官：俊逸面容，眉目含笑，气度从容。。武器招式：万剑归宗，剑作轻舟。。功法：万剑归宗；纯阳剑法登峰；八仙之首神通。。功法表现：万千剑影随行，剑气如虹。。画面：构图：沧海之上，白衣剑仙踏剑而行，身后万千剑影如虹，八仙各显神通随行，背景海天相接。色调：沧溟蓝+剑光白+道袍青。氛围：各显神通、宗师气度。。台词："八仙过海，各显神通。我这一手，就一个"剑"字，足以渡海。"。动作帧（动图）：①踏剑入海 ②袖袍一挥 ③万千剑影如虹散开 ④谈笑渡海。诗词：八仙过海各称雄，剑作轻舟万剑同。纯阳一脉开宗祖，笑傲沧溟万里风。。主题句：黄粱一梦醒功名，仗剑天地任逍遥——度人先度己，成仙先成人。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：沧溟蓝+剑光白+道袍青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 八仙之首, age 中年, 八仙过海·纯阳吕祖; scene: 海面踏剑而行，剑气化作万千剑影；为众仙开路，谈笑间已渡沧海。; 八仙过海各称雄，剑作轻舟万剑同。纯阳一脉开宗祖，笑傲沧溟万里风。; palette: 沧溟蓝+剑光白+道袍青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 纯阳帝君**
- 吕洞宾，道果圆满阶段·纯阳剑仙（成仙后，点化众生·纯阳剑心）。形象：吕洞宾，白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。品性：八仙之首，纯阳吕祖，点化众生。一生度人无数，悟得"度人先度己，成仙先成人"。。姿态：化身书生、乞丐、老翁游戏人间，以游戏之法点化有缘；剑已化于心中，一念剑气可斩心魔。。服饰：白衣便装，纯阳剑化作腰间一枚温润剑符，酒葫芦仍在。。体型：身高约7头身，白衣便装，气质出尘。。衣物细节：白衣便装，腰间剑符，酒葫芦。。发型妆造：白发束起，仙气隐然。。脸型五官：俊逸面容，眉目含笑，眼神慈悲温润。。武器招式：剑已化心符，一念即剑气。。功法：点化众生；剑心通明（一念即剑）；度人先度己。。功法表现：剑符微光，点化之光。。画面：构图：人间市井，白衣仙翁扮作老翁与孩童讲古，腰间剑符微光，身后斜阳，背景烟火人间。色调：白衣+暮阳金+人间暖。氛围：点化众生、逍遥入世。。台词："我这一生，先度了自己，才渡得了别人。剑在心上，何须在手上。"。动作帧（动图）：①扮老翁入市 ②与孩童讲古 ③腰间剑符微光一闪 ④负手看斜阳。诗词：纯阳剑仙隐红尘，游戏人间渡有缘。剑化心符存一念，逍遥自在了无尘。。主题句：黄粱一梦醒功名，仗剑天地任逍遥——度人先度己，成仙先成人。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：白衣+暮阳金+人间暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 纯阳剑仙, age 成仙后, 点化众生·纯阳剑心; scene: 化身书生、乞丐、老翁游戏人间，以游戏之法点化有缘；剑已化于心中，一念剑气可斩心魔。; 纯阳剑仙隐红尘，游戏人间渡有缘。剑化心符存一念，逍遥自在了无尘。; palette: 白衣+暮阳金+人间暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 何仙姑（`he_xiangu`） · 人生档案版

**灵胎初醒 · 岭南采药女**
- 何仙姑，灵胎初醒阶段·岭南采药女（少女，增城·采药济人）。形象：何仙姑，素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。品性：增城采药女，孝顺勤劳，常采药济人，医者仁心。。姿态：清晨背竹篓上山采药；归来为村人煎药治病，分文不取。。服饰：素裙荆钗，背竹篓，清秀灵动。。体型：身高约6头身，素裙少女，身形灵巧。。衣物细节：素裙荆钗，竹篓背带。。发型妆造：青丝挽髻，荆钗别发。。脸型五官：鹅蛋脸清秀，眉眼弯弯带笑，鼻梁挺，嘴角温柔。。武器招式：无兵器，采药为生。。功法：识药辨草，医理精通。。功法表现：无神力，医者仁光。。画面：构图：岭南青山，素裙少女背竹篓采药，晨光洒落，身后村落炊烟，背景青山如黛。色调：素裙白+草药绿+晨光金。氛围：勤劳、济世、清秀。。台词："采药不为换钱，只为村人少受一分病痛。"。动作帧（动图）：①背篓上山 ②俯身采药 ③归家煎药 ④为村人送药。诗词：增城女儿采药忙，朝露沾裙晓风凉。荆钗素裙医者心，不问金银问安康。。主题句：采药人间医百病，得道仍怀济世心——仙姑手中莲花，渡人亦渡己。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素裙白+草药绿+晨光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 岭南采药女, age 少女, 增城·采药济人; scene: 清晨背竹篓上山采药；归来为村人煎药治病，分文不取。; 增城女儿采药忙，朝露沾裙晓风凉。荆钗素裙医者心，不问金银问安康。; palette: 素裙白+草药绿+晨光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 仙桃轻身**
- 何仙姑，凡尘砺心阶段·食云母轻身（少女，遇仙得药·身轻如燕）。形象：何仙姑，素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。品性：偶遇仙人，指点食云母粉。服后身轻如燕，步履如飞，渐通仙缘。。姿态：服云母后晨起试步，竟能踏草而行；山中采药，脚下生风。。服饰：素裙轻扬，步履如飞。。体型：身高约6头身，素裙少女，身姿轻灵。。衣物细节：素裙轻扬，竹篓仍在。。发型妆造：青丝，荆钗。。脸型五官：清秀面容，眉眼含笑，眼神明亮。。武器招式：无兵器。。功法：云母轻身法；身轻如燕，日行百里。。功法表现：身轻如燕，足尖生风。。画面：构图：岭南溪谷，素裙少女踏草而行，衣袂轻扬，草尖微弯不折，背景云雾山峦。色调：素裙+云母白+山青。氛围：轻身、仙缘初结。。台词："仙人指我食云母，从此轻身如燕。我且用这轻身术，采更多的药，救更多的人。"。动作帧（动图）：①服云母粉 ②起身试步 ③踏草而行 ④展臂迎风。诗词：云母一匙身渐轻，踏草穿云步履轻。采药更上山深处，济世之心愈发明。。主题句：采药人间医百病，得道仍怀济世心——仙姑手中莲花，渡人亦渡己。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素裙+云母白+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 食云母轻身, age 少女, 遇仙得药·身轻如燕; scene: 服云母后晨起试步，竟能踏草而行；山中采药，脚下生风。; 云母一匙身渐轻，踏草穿云步履轻。采药更上山深处，济世之心愈发明。; palette: 素裙+云母白+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 八仙之女**
- 何仙姑，道法初成阶段·吕祖点化（少女，得吕洞宾点化·手持莲花）。形象：何仙姑，素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。品性：吕洞宾游历岭南，见其心善根深，点化度之。授以莲花，传以仙法。。姿态：吕祖授如意莲花，她持莲花静听仙法；莲花在手，百花相和。。服饰：素衣仙姿，手持如意莲花。。体型：身高约6头身，素衣仙姿，手持莲花。。衣物细节：素衣仙姿，手持如意莲花。。发型妆造：青丝披肩，莲花簪。。脸型五官：清秀面容，眉目含仙气，嘴角带笑。。武器招式：如意莲花（百花相和）。。功法：得吕祖真传；如意莲花护体；仙道初成。。功法表现：莲花光晕，百花环绕。。画面：构图：山间溪畔，白衣仙姑手持如意莲花，吕祖身影于云中隐约，身后莲花满池。色调：莲花粉+仙白+溪青。氛围：得道点化、仙缘已定。。台词："吕祖教我，莲花出淤泥而不染。我做人，也要像这莲花——生在人间，心在云端。"。动作帧（动图）：①静立听法 ②双手接莲 ③莲花开合 ④拈花微笑。诗词：吕祖遥指一莲来，素手拈花笑靥开。淤泥不染清如水，仙缘从此入胸怀。。主题句：采药人间医百病，得道仍怀济世心——仙姑手中莲花，渡人亦渡己。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：莲花粉+仙白+溪青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 吕祖点化, age 少女, 得吕洞宾点化·手持莲花; scene: 吕祖授如意莲花，她持莲花静听仙法；莲花在手，百花相和。; 吕祖遥指一莲来，素手拈花笑靥开。淤泥不染清如水，仙缘从此入胸怀。; palette: 莲花粉+仙白+溪青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 百花阵**
- 何仙姑，大劫淬炼阶段·济世救人（少女，莲花护体·救苦济难）。形象：何仙姑，素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。品性：下山济世，莲花护体，救苦济难。采药、治病、镇灾，凡所至处，百姓得安。。姿态：以莲花化雨露救旱灾；以仙药解瘟疫；路见不平，玉手拈花抛起便是漫天仙光。。服饰：素衣仙裙，手持莲花，气质出尘。。体型：身高约6头身，素衣仙裙，飘然若仙。。衣物细节：素衣仙裙，手持莲花。。发型妆造：青丝，莲花簪。。脸型五官：清秀面容，眉目慈悲，嘴角带济世微笑。。武器招式：如意莲花：抛起化雨露/仙光。。功法：莲花护体；仙药救人；移花接木。。功法表现：莲花化甘霖，仙光护体。。画面：构图：村野灾后，仙姑持莲花洒下甘霖，枯苗回绿，百姓叩谢，背景村庄重焕生机。色调：莲花粉+甘霖蓝+新绿。氛围：济世、慈悲、救苦救难。。台词："采药医百病，莲花渡众生。我何仙姑，救人救到底。"。动作帧（动图）：①抛起莲花 ②莲花化雨露 ③枯苗回绿 ④百姓叩谢。诗词：莲花一瓣化甘霖，仙药回春救万民。素手拈花抛日月，人间何处不逢春。。主题句：采药人间医百病，得道仍怀济世心——仙姑手中莲花，渡人亦渡己。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：莲花粉+甘霖蓝+新绿。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 济世救人, age 少女, 莲花护体·救苦济难; scene: 以莲花化雨露救旱灾；以仙药解瘟疫；路见不平，玉手拈花抛起便是漫天仙光。; 莲花一瓣化甘霖，仙药回春救万民。素手拈花抛日月，人间何处不逢春。; palette: 莲花粉+甘霖蓝+新绿; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 八仙过海**
- 何仙姑，封神登天阶段·八仙过海（得道后，八仙过海·荷花仙子）。形象：何仙姑，素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。品性：八仙过海，她踏荷而行，荷花仙子，与群仙各显神通。。姿态：脚下荷叶生莲，踏浪而行；莲花开处，海波平静，鱼龙相贺。。服饰：素衣仙裙，手持莲花，荷花仙子气度。。体型：身高约6头身，素衣仙裙，踏荷凌波。。衣物细节：素衣仙裙，手持莲花。。发型妆造：青丝，莲花簪，仙气飘然。。脸型五官：清秀面容，眉目灵动，嘴角含笑。。武器招式：如意莲花（踏浪/镇波）。。功法：八仙过海；莲花神通大盛。。功法表现：莲花绽放，波平浪静。。画面：构图：沧海之上，仙姑踏荷而行，莲花绽放，八仙各显神通，海天开阔，背景明月沧波。色调：荷粉+沧溟蓝+月白。氛围：各显神通、凌波仙子。。台词："八仙过海，我这朵莲花，够我踏遍沧溟了。"。动作帧（动图）：①踏荷入海 ②莲花绽开 ③海波平静 ④凌波而行。诗词：八仙过海各神通，荷为舟楫踏浪行。莲花开处风波静，仙子凌波海月明。。主题句：采药人间医百病，得道仍怀济世心——仙姑手中莲花，渡人亦渡己。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：荷粉+沧溟蓝+月白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 八仙过海, age 得道后, 八仙过海·荷花仙子; scene: 脚下荷叶生莲，踏浪而行；莲花开处，海波平静，鱼龙相贺。; 八仙过海各神通，荷为舟楫踏浪行。莲花开处风波静，仙子凌波海月明。; palette: 荷粉+沧溟蓝+月白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 瑶池仙姑**
- 何仙姑，道果圆满阶段·花中真仙（成仙后，花开满山·泽被山水）。形象：何仙姑，素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。品性：花中真仙，天下百花皆为其化身。采药济世之心不改，得道仍怀济世之志。。姿态：所过之处，花开满山；以百花之力泽被山水，护佑一方风调雨顺。。服饰：素衣仙姿，莲花在侧，周身花气。。体型：身高约6头身，素衣仙姿，立于花海。。衣物细节：素衣仙姿，莲花在侧。。发型妆造：青丝，莲花簪，花气绕身。。脸型五官：清秀面容，眉目含笑，慈悲温婉。。武器招式：如意莲花（百花听命）。。功法：百花化身；花开满山；济世之心化道。。功法表现：满山花开，春意流转。。画面：构图：满山花海，素衣仙姑立于花间，一莲绽放百花生，远处村落祥和，背景春山如锦。色调：百花彩+莲粉+春山绿。氛围：花开满山、泽被人间。。台词："做了仙，我还是那个采药的姑娘。满山花开，就是我济世的心。"。动作帧（动图）：①立于花间 ②莲开百花生 ③挥手花随 ④负手望满山。诗词：花中真仙隐林泉，一莲开处百花连。采药之心终不泯，春风化雨满人间。。主题句：采药人间医百病，得道仍怀济世心——仙姑手中莲花，渡人亦渡己。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：百花彩+莲粉+春山绿。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 花中真仙, age 成仙后, 花开满山·泽被山水; scene: 所过之处，花开满山；以百花之力泽被山水，护佑一方风调雨顺。; 花中真仙隐林泉，一莲开处百花连。采药之心终不泯，春风化雨满人间。; palette: 百花彩+莲粉+春山绿; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 张果老（`zhang_guolao`） · 人生档案版

**灵胎初醒 · 恒山樵夫**
- 张果老，灵胎初醒阶段·恒山樵夫（壮年，恒山砍柴·憨厚樵夫）。形象：张果老，灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。品性：恒山砍柴的憨厚樵夫，日出而作，日落而息，与山林为伴。。姿态：清晨上山砍柴，挑担下山换米；山间唱山歌，自在逍遥。。服饰：粗布麻衣，柴担在肩，腰间别酒葫芦。。体型：身高约6头身，粗布樵夫，身形壮实。。衣物细节：粗布麻衣，柴担，腰间酒葫芦。。发型妆造：乱发扎巾，胡须粗犷。。脸型五官：圆脸憨厚，粗眉，眼大而诚，鼻梁宽，嘴角带笑，下巴圆。。武器招式：无兵器，柴刀/扁担。。功法：无道术，唯熟识山路林径。。功法表现：无神力，山野气息。。画面：构图：恒山古道，粗布樵夫挑柴而行，山歌悠悠，白云绕峰，背景青山叠嶂。色调：麻衣褐+山青+云白。氛围：憨厚、山野、自在。。台词："柴要劈得匀，山要认得熟。我张果老一辈子，就守着这恒山过日子。"。动作帧（动图）：①挥斧砍柴 ②捆柴上肩 ③挑担下山 ④歇脚唱山歌。诗词：恒山樵夫日日忙，柴担压肩唱山冈。白云深处人家少，自在山中岁月长。。主题句：世人匆匆向前赶，唯我倒骑白驴看云闲——回头看，才是向前。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：麻衣褐+山青+云白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 恒山樵夫, age 壮年, 恒山砍柴·憨厚樵夫; scene: 清晨上山砍柴，挑担下山换米；山间唱山歌，自在逍遥。; 恒山樵夫日日忙，柴担压肩唱山冈。白云深处人家少，自在山中岁月长。; palette: 麻衣褐+山青+云白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 轻身之能**
- 张果老，凡尘砺心阶段·食首乌返老（壮年，千年首乌·返老还童）。形象：张果老，灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。品性：山中偶得千年何首乌，食后白发转黑，皱纹渐消，竟返老还童，才知是仙缘。。姿态：服首乌后对溪自照，惊觉青春复来；从此通晓仙机。。服饰：麻衣依旧，容光焕发，倒像年轻人。。体型：身高约6头身，容光焕发，身姿矫健如青年。。衣物细节：麻衣，腰间酒葫芦。。发型妆造：黑发披肩，神采奕奕。。脸型五官：面容返青，粗眉，眼亮有神，嘴角含笑。。武器招式：无兵器。。功法：千年首乌之效；寿元大增。。功法表现：生机流转，寿元之气。。画面：构图：恒山溪畔，樵夫对水自照，水中映出少年面容，青丝如瀑，背景古树溪石。色调：麻衣+溪青+首乌棕。氛围：返老还童、奇遇仙缘。。台词："这一棵首乌，把我吃回了青春。山里的日子，原来藏着这样的宝贝。"。动作帧（动图）：①掘得首乌 ②食之 ③对溪自照惊觉 ④展臂身轻。诗词：千年首乌土中藏，一食白发转青苍。对溪照得青春面，方知此物是仙粮。。主题句：世人匆匆向前赶，唯我倒骑白驴看云闲——回头看，才是向前。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：麻衣+溪青+首乌棕。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 食首乌返老, age 壮年, 千年首乌·返老还童; scene: 服首乌后对溪自照，惊觉青春复来；从此通晓仙机。; 千年首乌土中藏，一食白发转青苍。对溪照得青春面，方知此物是仙粮。; palette: 麻衣+溪青+首乌棕; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 鱼鼓神通**
- 张果老，道法初成阶段·倒骑白驴（中年，倒骑驴·优游云海）。形象：张果老，灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。品性：得道后倒骑白驴，优游云海，看世间万象。世人皆向前，他偏倒着看——回头看，才是向前。。姿态：倒骑白驴悠然而行，鱼鼓轻敲；驴儿识途，他只看云看山。。服饰：灰白道袍，倒骑白驴，手持鱼鼓。。体型：身高约6头身，道翁身形，倒骑驴背。。衣物细节：灰白道袍，倒骑白驴，鱼鼓。。发型妆造：白发束起，白须垂胸。。脸型五官：圆脸慈祥，白眉长垂，眼含笑，鼻梁圆，嘴角悠然。。武器招式：鱼鼓（声定风波）。。功法：倒骑驴术（来去无踪）；鱼鼓定风波。。功法表现：鱼鼓一响，风波自定。。画面：构图：云山古道，白须道翁倒骑白驴，手持鱼鼓，白云绕身，前方道路蜿蜒，背景云海。色调：道袍灰白+云海白+山青。氛围：逍遥倒行、看云闲。。台词："世人匆匆向前赶，我偏倒骑驴看云。走错的路，回头再看，才看得明白。"。动作帧（动图）：①倒骑驴背 ②轻敲鱼鼓 ③驴儿识途前行 ④他只顾看云。诗词：倒骑白驴看云闲，鱼鼓轻敲世外天。回头方识来时路，不是寻常一半仙。。主题句：世人匆匆向前赶，唯我倒骑白驴看云闲——回头看，才是向前。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道袍灰白+云海白+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 倒骑白驴, age 中年, 倒骑驴·优游云海; scene: 倒骑白驴悠然而行，鱼鼓轻敲；驴儿识途，他只看云看山。; 倒骑白驴看云闲，鱼鼓轻敲世外天。回头方识来时路，不是寻常一半仙。; palette: 道袍灰白+云海白+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 呼风唤雨**
- 张果老，大劫淬炼阶段·屡召不朝（老年，唐帝屡召·装死脱身）。形象：张果老，灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。品性：唐玄宗闻其名屡次召见，他装死脱身；被识破又谈笑而遁，始终不肯入朝为官。。姿态：朝廷使者至，他倒头装死，气息全无；使者去，他又活转来倒骑驴大笑。。服饰：道袍，倒骑驴，一脸装死的安详。。体型：身高约6头身，道翁横卧，装死之姿。。衣物细节：道袍，倒骑驴。。发型妆造：白发，白须。。脸型五官：闭目安详，嘴角却含一丝捉狭笑意。。武器招式：无兵器，假死遁身。。功法：假死术；遁身法。。功法表现：假死气息尽敛。。画面：构图：客舍之中，使者面前道翁闭目横卧"气绝"，一旁白驴悠闲吃草，背景皇命诏书。色调：道袍灰白+诏书金+客舍褐。氛围：装死脱身、游戏君王。。台词："做官有什么好？不如我这倒骑驴，天不管地不管，自在逍遥。"。动作帧（动图）：①使者至 ②倒头装死 ③气息全无 ④使者去后翻身大笑。诗词：帝诏频来召老仙，闭目装死入黄泉。使者去后翻身起，笑倒驴儿云外天。。主题句：世人匆匆向前赶，唯我倒骑白驴看云闲——回头看，才是向前。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道袍灰白+诏书金+客舍褐。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 屡召不朝, age 老年, 唐帝屡召·装死脱身; scene: 朝廷使者至，他倒头装死，气息全无；使者去，他又活转来倒骑驴大笑。; 帝诏频来召老仙，闭目装死入黄泉。使者去后翻身起，笑倒驴儿云外天。; palette: 道袍灰白+诏书金+客舍褐; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 鱼鼓天音**
- 张果老，封神登天阶段·八仙过海（得道后，八仙过海·不老仙翁）。形象：张果老，灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。品性：八仙过海，他倒骑驴踏浪而行，鱼鼓一拍定风波，八仙中的不老仙翁。。姿态：白驴踏浪，倒骑而行；鱼鼓一拍，海波平息，众仙相顾而笑。。服饰：灰白道袍，倒骑驴，鱼鼓在手。。体型：身高约6头身，道翁倒骑驴，踏浪而行。。衣物细节：灰白道袍，倒骑白驴，鱼鼓。。发型妆造：白发，白须。。脸型五官：慈祥面容，白眉，眼含笑，嘴角悠然。。武器招式：鱼鼓（定风波）。。功法：鱼鼓天音；缩地成寸；辟谷神功。。功法表现：鱼鼓一拍，浪平风静。。画面：构图：沧海之上，白须道翁倒骑白驴踏浪，鱼鼓轻敲，海波自平，八仙随行，背景明月沧波。色调：道袍灰白+沧溟蓝+月白。氛围：不老仙翁、各显神通。。台词："八仙过海，我连驴都不肯正着骑——倒着骑，浪也认得我。"。动作帧（动图）：①倒骑驴入海 ②驴踏浪而行 ③鱼鼓一拍 ④海波平八仙笑。诗词：八仙过海各神通，倒骑驴儿踏浪行。鱼鼓一拍风波定，不老仙翁笑苍穹。。主题句：世人匆匆向前赶，唯我倒骑白驴看云闲——回头看，才是向前。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道袍灰白+沧溟蓝+月白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 八仙过海, age 得道后, 八仙过海·不老仙翁; scene: 白驴踏浪，倒骑而行；鱼鼓一拍，海波平息，众仙相顾而笑。; 八仙过海各神通，倒骑驴儿踏浪行。鱼鼓一拍风波定，不老仙翁笑苍穹。; palette: 道袍灰白+沧溟蓝+月白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 中条洞主**
- 张果老，道果圆满阶段·混元仙翁（成仙后，倒骑驴巡看人间）。形象：张果老，灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。品性：混元仙翁，倒骑白驴巡看人间。千岁之龄，一颗童心，看云看山看人间烟火。。姿态：云游人间，鱼鼓一响定风波；遇迷途人，指点一句回头路；依旧是那副倒骑驴的逍遥。。服饰：灰白道袍，倒骑驴，鱼鼓，仙翁气度。。体型：身高约6头身，道翁倒骑驴，云端而行。。衣物细节：灰白道袍，倒骑白驴，鱼鼓。。发型妆造：白发飘飘，白须垂胸。。脸型五官：慈祥面容，白眉长垂，眼含笑意，嘴角悠然。。武器招式：鱼鼓（震三界）。。功法：混元仙翁神位；鱼鼓天音震三界；倒骑驴巡天。。功法表现：鱼鼓天音，云海随行。。画面：构图：云端之上，白须仙翁倒骑驴，鱼鼓轻敲，身后云海与人间山河，背景夕阳金霞。色调：道袍灰白+金霞+云海白。氛围：混元逍遥、无求自在。。台词："做了一千年神仙，我还是最喜欢倒骑驴——因为只有回头看，才看得清自己走过的路。"。动作帧（动图）：①倒骑驴云游 ②鱼鼓定风波 ③指点迷途人 ④回望人间一笑。诗词：混元仙翁倒骑驴，千年云外看尘途。回头一笑人间事，鱼鼓声中万里疏。。主题句：世人匆匆向前赶，唯我倒骑白驴看云闲——回头看，才是向前。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道袍灰白+金霞+云海白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 混元仙翁, age 成仙后, 倒骑驴巡看人间; scene: 云游人间，鱼鼓一响定风波；遇迷途人，指点一句回头路；依旧是那副倒骑驴的逍遥。; 混元仙翁倒骑驴，千年云外看尘途。回头一笑人间事，鱼鼓声中万里疏。; palette: 道袍灰白+金霞+云海白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 铁拐李（`tie_guaili`） · 人生档案版

**灵胎初醒 · 巴国隐士**
- 铁拐李，灵胎初醒阶段·隐士李玄（中年，痴迷黄老·华岳修道）。形象：铁拐李，跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。品性：本名李玄，痴迷黄老之术的隐士，在华岳山中听道静修。。姿态：山中诵黄庭，吐纳导引；偶下山，讲道济人。。服饰：青布道袍，束发，清癯。。体型：身高约7头身，青袍隐士，清癯修长。。衣物细节：青布道袍，束发。。发型妆造：束发，清癯。。脸型五官：长脸清癯，眉目淡泊，眼含道意，鼻梁挺，嘴角平和。。武器招式：无兵器。。功法：黄老之术；吐纳导引。。功法表现：吐纳之气。。画面：构图：华岳山间草庐，青袍隐士闭目吐纳，身周松风，背景群山云海。色调：道青+山翠+云白。氛围：静修、慕道、清癯。。台词："黄老之术，讲的是清静无为。我且在这山中，把这一身浊气，修成清气。"。动作帧（动图）：①闭目吐纳 ②诵黄庭经 ③下山讲道 ④回山静修。诗词：华岳山中隐李玄，黄庭一卷日周旋。清癯道骨烟霞气，不向人间问利钱。。主题句：肉身成灰，一拐一拐走天涯——形可残，道不可缺。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：道青+山翠+云白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 隐士李玄, age 中年, 痴迷黄老·华岳修道; scene: 山中诵黄庭，吐纳导引；偶下山，讲道济人。; 华岳山中隐李玄，黄庭一卷日周旋。清癯道骨烟霞气，不向人间问利钱。; palette: 道青+山翠+云白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 弃儒从道**
- 铁拐李，凡尘砺心阶段·神游焚身（中年，元神出窍·肉身被焚）。形象：铁拐李，跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。品性：一日元神出窍神游太虚，弟子以为其坐化，将其肉身火化。归来时肉身已成灰烬。。姿态：元神游于太虚，忽见自家肉身被焚，大恸；魂魄无处依附，飘于荒野。。服饰：肉身成灰，魂魄飘零。。体型：魂影飘忽，无实体。。衣物细节：肉身成灰，魂影。。发型妆造：魂影淡发。。脸型五官：魂魄虚影，面目悲恸。。武器招式：无兵器。。功法：元神出窍；魂游太虚。。功法表现：魂游太虚，飘忽不定。。画面：构图：火化炉前青烟升腾，一缕魂魄立于荒野，望着余烬，身影飘忽，背景暮色苍茫。色调：火烬灰+暮色暗+魂影淡。氛围：大悲、肉身成灰、飘零。。台词："只道神游太虚去，归来肉身已成灰。这皮囊，终究是留不住。"。动作帧（动图）：①元神出窍 ②见肉身被焚 ③大恸 ④魂魄飘于荒野。诗词：元神出窍游太虚，归来肉身已焚墟。飘魂四顾无依处，泣向苍天问何如。。主题句：肉身成灰，一拐一拐走天涯——形可残，道不可缺。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：火烬灰+暮色暗+魂影淡。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 神游焚身, age 中年, 元神出窍·肉身被焚; scene: 元神游于太虚，忽见自家肉身被焚，大恸；魂魄无处依附，飘于荒野。; 元神出窍游太虚，归来肉身已焚墟。飘魂四顾无依处，泣向苍天问何如。; palette: 火烬灰+暮色暗+魂影淡; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 附身乞丐**
- 铁拐李，道法初成阶段·附身乞丐（中年，附身跛脚乞丐·铁拐为杖）。形象：铁拐李，跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。品性：汉钟离救之，令其附身于路旁一跛脚乞丐尸身。虽形残貌陋，道心却更坚。。姿态：附身后拄拐起身，跌跌撞撞走了几步，忽大笑——"有这双腿脚，够走天涯了"。。服饰：破衣烂衫，跛足，铁拐为杖。。体型：身高约6头身，跛足身形，拄铁拐。。衣物细节：破衣烂衫，铁拐为杖。。发型妆造：乱发蓬面。。脸型五官：蓬头乱发，眼却清亮有神，嘴角含笑。。武器招式：铁拐为兵。。功法：附身术；以铁拐为兵。。功法表现：附身之魂光。。画面：构图：荒野路旁，跛脚乞丐拄铁拐而立，破衣扬风，仰天大笑，汉钟离虚影于云中点头，背景旷野。色调：破衣褐+铁拐黑+暮光金。氛围：重获新生、形残道坚。。台词："形残了，道不残。我李玄，拄着这根铁拐，照样走遍三界。"。动作帧（动图）：①附身乞丐 ②拄拐起身 ③跌撞几步 ④仰天大笑。诗词：一拐一拐复还魂，破衣跛足笑红尘。皮囊虽破心尤亮，铁拐撑起又一身。。主题句：肉身成灰，一拐一拐走天涯——形可残，道不可缺。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衣褐+铁拐黑+暮光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 附身乞丐, age 中年, 附身跛脚乞丐·铁拐为杖; scene: 附身后拄拐起身，跌跌撞撞走了几步，忽大笑——"有这双腿脚，够走天涯了"。; 一拐一拐复还魂，破衣跛足笑红尘。皮囊虽破心尤亮，铁拐撑起又一身。; palette: 破衣褐+铁拐黑+暮光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 混元丹术**
- 铁拐李，大劫淬炼阶段·济世跛仙（中年，铁拐济世·看破皮囊）。形象：铁拐李，跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。品性：以最卑微之相行最慈悲之事。葫芦纳天地，铁拐济世，看破皮囊，形残心明。。姿态：葫芦一开，乾坤皆入其中；铁拐一顿，妖邪辟易；济人度世，从不以貌取人。。服饰：破衣铁拐，腰悬宝葫芦。。体型：身高约6头身，跛足拄拐，身形佝偻却精神。。衣物细节：破衣铁拐，宝葫芦。。发型妆造：乱发，蓬面。。脸型五官：蓬头之下，眼神慈悲明亮。。武器招式：铁拐（顿地辟邪）；葫芦（纳天地）。。功法：铁拐神通；宝葫芦纳天地；济世度人。。功法表现：葫芦纳星斗，铁拐辟邪。。画面：构图：市井街头，跛足仙翁拄拐而立，葫芦口倾出万千星斗，破衣扬风，众百姓敬拜，背景人间烟火。色调：破衣褐+葫芦金+星斗蓝。氛围：慈悲济世、形残心明。。台词："我这张破皮囊，正是度人的招牌——连我这副模样都肯度，谁不肯度？"。动作帧（动图）：①拄拐入市 ②葫芦一倾 ③星斗万千 ④百姓敬拜。诗词：跛仙铁拐走风尘，葫芦一倾纳乾坤。形残偏是慈悲相，度尽苍生不看人。。主题句：肉身成灰，一拐一拐走天涯——形可残，道不可缺。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衣褐+葫芦金+星斗蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 济世跛仙, age 中年, 铁拐济世·看破皮囊; scene: 葫芦一开，乾坤皆入其中；铁拐一顿，妖邪辟易；济人度世，从不以貌取人。; 跛仙铁拐走风尘，葫芦一倾纳乾坤。形残偏是慈悲相，度尽苍生不看人。; palette: 破衣褐+葫芦金+星斗蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 葫芦纳天地**
- 铁拐李，封神登天阶段·八仙过海（得道后，八仙过海·葫芦装海）。形象：铁拐李，跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。品性：八仙过海，他铁拐拄浪、葫芦装海，八仙中的老大哥，谈笑间渡海。。姿态：葫芦一倾，装下半海；铁拐一顿，踏浪如平。。服饰：破衣铁拐，宝葫芦，老大哥气度。。体型：身高约6头身，拄拐立浪，气度沉稳。。衣物细节：破衣铁拐，宝葫芦。。发型妆造：乱发。。脸型五官：蓬头，眼神豪迈带笑。。武器招式：铁拐拄浪；葫芦装海。。功法：葫芦装海；铁拐神通大盛。。功法表现：葫芦倾海，浪自平。。画面：构图：沧海之上，跛足仙翁铁拐拄浪而立，葫芦倾倒海水，八仙随行，背景海天辽阔。色调：破衣褐+沧海蓝+葫芦金。氛围：老大哥、装海神通。。台词："八仙过海，我这葫芦，先装一半海水，省得大家挤。"。动作帧（动图）：①铁拐拄浪 ②葫芦一倾装海 ③海水倒入海中 ④八仙同笑。诗词：八仙过海我先行，葫芦装海浪自平。铁拐一顿风云起，老大哥来笑一声。。主题句：肉身成灰，一拐一拐走天涯——形可残，道不可缺。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衣褐+沧海蓝+葫芦金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 八仙过海, age 得道后, 八仙过海·葫芦装海; scene: 葫芦一倾，装下半海；铁拐一顿，踏浪如平。; 八仙过海我先行，葫芦装海浪自平。铁拐一顿风云起，老大哥来笑一声。; palette: 破衣褐+沧海蓝+葫芦金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 铁拐大仙**
- 铁拐李，道果圆满阶段·乞丐道祖（成仙后，铁拐问天·乞丐成仙）。形象：铁拐李，跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。品性：以乞丐之相证得道祖之位。肉身成灰、附身跛丐，终成大觉——形可残，道不可缺。。姿态：行于三界，铁拐一顿山河动，葫芦一开纳星辰；依旧那副破衣，却已是道祖气象。。服饰：破衣铁拐，宝葫芦，周身道气冲霄。。体型：身高约6头身，拄拐云端，道祖气象。。衣物细节：破衣铁拐，宝葫芦，道气冲霄。。发型妆造：乱发，道气绕身。。脸型五官：蓬头之下，眼神澄澈如日月。。武器招式：铁拐化道；葫芦纳天地。。功法：乞丐道祖神位；铁拐化道；葫芦纳天地。。功法表现：星辰入葫芦，道气冲霄汉。。画面：构图：云端之上，跛足道祖拄拐而立，葫芦开处星斗漫天，破衣扬风而道气冲霄，背景星河云海。色调：破衣褐+星河蓝+道气金。氛围：乞丐道祖、形残道全。。台词："我这一生，皮囊烧过、换过、破过。唯有这颗道心，一路没丢过。"。动作帧（动图）：①拄拐立云 ②葫芦开纳星辰 ③铁拐一顿 ④俯视三界。诗词：乞丐道祖立乾坤，铁拐问天万法尊。皮囊虽破心如玉，一粒金丹渡世尘。。主题句：肉身成灰，一拐一拐走天涯——形可残，道不可缺。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衣褐+星河蓝+道气金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 乞丐道祖, age 成仙后, 铁拐问天·乞丐成仙; scene: 行于三界，铁拐一顿山河动，葫芦一开纳星辰；依旧那副破衣，却已是道祖气象。; 乞丐道祖立乾坤，铁拐问天万法尊。皮囊虽破心如玉，一粒金丹渡世尘。; palette: 破衣褐+星河蓝+道气金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 汉钟离（`han_zhongli`） · 人生档案版

**灵胎初醒 · 燕台武将**
- 汉钟离，灵胎初醒阶段·燕台武将（壮年，勇冠三军·汉朝武将）。形象：汉钟离，赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。品性：汉朝燕台武将钟离权，勇冠三军，沙场征伐，功勋赫赫。。姿态：挥戈跃马，阵前斩将；战后抚剑，看尽生死。。服饰：红袍战甲，赤发虬髯，腰悬宝刀。。体型：身高约7头身，红袍战甲，魁梧威猛。。衣物细节：红袍战甲，腰悬宝刀。。发型妆造：赤发披散，虬髯。。脸型五官：国字脸，赤发，浓眉，虎目含威，鼻梁挺，虬髯，下巴方正。。武器招式：长戈/宝刀，沙场杀伐。。功法：沙场武艺，万人难敌。。功法表现：无神力，煞气缠身。。画面：构图：沙场之上，红袍战将跃马挥戈，赤发猎猎，敌军溃散，背景烽烟战旗。色调：战甲红+烽烟黑+落日金。氛围：勇武、杀伐、武将本色。。台词："刀下亡魂千千万，功名簿上血染成。这刀，我握了半辈子，越握越沉。"。动作帧（动图）：①跃马提戈 ②阵前斩将 ③战后抚剑 ④望烽烟沉思。诗词：燕台猛将赤发髯，跃马挥戈破阵烟。功名血染刀头重，半生杀伐半生悬。。主题句：沙场猛将弃剑修，一扇火里见丹心——放下屠刀，自成大日。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：战甲红+烽烟黑+落日金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 燕台武将, age 壮年, 勇冠三军·汉朝武将; scene: 挥戈跃马，阵前斩将；战后抚剑，看尽生死。; 燕台猛将赤发髯，跃马挥戈破阵烟。功名血染刀头重，半生杀伐半生悬。; palette: 战甲红+烽烟黑+落日金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 弃武修道**
- 汉钟离，凡尘砺心阶段·兵败入山（中年，兵败入山·看破红尘）。形象：汉钟离，赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。品性：一战兵败，看破功名浮沉，弃武入山，寻道问仙。。姿态：散尽甲兵，独入深山；见云卷云舒，心渐澄明。。服饰：旧甲卸去，青布道衣。。体型：身高约7头身，青布道衣，身形依旧魁梧。。衣物细节：旧甲卸去，青布道衣。。发型妆造：赤发束起。。脸型五官：赤发束起，虎目渐平，眉宇含思。。武器招式：弃戈，无兵器。。功法：弃武修道；道心初萌。。功法表现：无神力，心渐澄明。。画面：构图：深山古径，卸甲的武将青布道衣独行，云雾绕身，回头望一眼来路烽烟，继续前行。色调：青布+山青+云白。氛围：看破、归隐、道心初萌。。台词："这一败，败得好。败尽了我的杀心，才装得下这颗道心。"。动作帧（动图）：①卸甲 ②独入深山 ③回头望烽烟 ④转身前行。诗词：一战兵消万事休，弃戈入山看云游。功名原是刀头血，洗净铅华道自谋。。主题句：沙场猛将弃剑修，一扇火里见丹心——放下屠刀，自成大日。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青布+山青+云白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 兵败入山, age 中年, 兵败入山·看破红尘; scene: 散尽甲兵，独入深山；见云卷云舒，心渐澄明。; 一战兵消万事休，弃戈入山看云游。功名原是刀头血，洗净铅华道自谋。; palette: 青布+山青+云白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 芭蕉扇**
- 汉钟离，道法初成阶段·东华传法（中年，东华帝君传法·赤阳真人）。形象：汉钟离，赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。品性：得东华帝君传法，授赤阳神功。悟得"放下屠刀，自成大日"，修成赤阳真人。。姿态：山中炼赤阳之气，掌出烈焰；芭蕉扇入手，一扇火海滔天。。服饰：红袍道衣，手持芭蕉扇，赤发虬髯。。体型：身高约7头身，红袍道衣，魁梧如塔。。衣物细节：红袍道衣，芭蕉扇。。发型妆造：赤发披散。。脸型五官：赤发，浓眉，虎目含火，鼻梁挺，虬髯，方正。。武器招式：芭蕉扇（扇风生火/清风送归）。。功法：赤阳神功；芭蕉扇（扇风生火）；三昧真火初成。。功法表现：赤阳真火，大日之气。。画面：构图：山间道场，红袍道人持芭蕉扇一扇，烈焰滔天，赤发飞舞，身后大日初升。色调：赤阳红+焰火橙+大日金。氛围：得道、赤阳、真人。。台词："放下屠刀，自成大日。我这一身杀伐气，炼成了赤阳真火。"。动作帧（动图）：①炼赤阳之气 ②芭蕉扇入手 ③一扇火海 ④赤发飞舞。诗词：东华传法授赤阳，一扇火海照天光。杀伐尽化丹炉火，自心即是大日堂。。主题句：沙场猛将弃剑修，一扇火里见丹心——放下屠刀，自成大日。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：赤阳红+焰火橙+大日金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 东华传法, age 中年, 东华帝君传法·赤阳真人; scene: 山中炼赤阳之气，掌出烈焰；芭蕉扇入手，一扇火海滔天。; 东华传法授赤阳，一扇火海照天光。杀伐尽化丹炉火，自心即是大日堂。; palette: 赤阳红+焰火橙+大日金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 炎龙召唤**
- 汉钟离，大劫淬炼阶段·度人点化（中年，黄粱点化·度吕洞宾）。形象：汉钟离，赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。品性：以黄粱一梦点化吕洞宾，又点化蓝采和、何仙姑等，八仙之师。。姿态：客店中借黄粱一梦点化吕岩；市井间点化歌者、医女，广结仙缘。。服饰：红袍道衣，芭蕉扇在手，仙师气度。。体型：身高约7头身，红袍道衣，仙师气度。。衣物细节：红袍道衣，芭蕉扇。。发型妆造：赤发披散。。脸型五官：赤发，浓眉，虎目含笑，虬髯，慈严并具。。武器招式：芭蕉扇；点化灵光。。功法：点化之术；黄粱一梦；赤阳神通。。功法表现：黄粱一梦，点化之光。。画面：构图：邯郸客店，红袍道人立于榻前，指尖一点灵光入吕洞宾梦中，黄粱正熟，背景客店烛影。色调：红袍+灵光金+客店褐。氛围：度人、仙师、点化。。台词："度人如点灯。我这一扇，扇醒一个吕洞宾，八仙便多一人。"。动作帧（动图）：①立于榻前 ②指尖一点灵光 ③吕洞宾梦中历劫 ④转身微笑。诗词：黄粱一枕度仙才，点化群贤入道来。赤阳一扇风波起，八仙从此次第开。。主题句：沙场猛将弃剑修，一扇火里见丹心——放下屠刀，自成大日。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：红袍+灵光金+客店褐。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 度人点化, age 中年, 黄粱点化·度吕洞宾; scene: 客店中借黄粱一梦点化吕岩；市井间点化歌者、医女，广结仙缘。; 黄粱一枕度仙才，点化群贤入道来。赤阳一扇风波起，八仙从此次第开。; palette: 红袍+灵光金+客店褐; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 赤阳神功**
- 汉钟离，封神登天阶段·正阳帝君（得道后，八仙武力之首·正阳帝君）。形象：汉钟离，赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。品性：位列八仙，武力之首，受封正阳帝君。赤阳神功大成，一扇火海镇三界。。姿态：八仙过海，芭蕉扇一扇，火海开道；受封正阳帝君，不改武将本色。。服饰：红袍大氅，正阳帝君冠，芭蕉扇在手。。体型：身高约7头身，红袍大氅，帝君威仪。。衣物细节：红袍大氅，正阳帝君冠，芭蕉扇。。发型妆造：赤发披散，帝君冠。。脸型五官：赤发，浓眉，虎目如日，虬髯，方正威严。。武器招式：芭蕉扇（一扇火海/清风送归）。。功法：正阳帝君神位；赤阳神功大成；芭蕉扇镇世。。功法表现：正阳大日，火凤绕身。。画面：构图：云端之上，红袍帝君持芭蕉扇而立，身后大日正阳，火凤绕身，八仙随行，背景赤霞云海。色调：正阳红+大日金+霞紫。氛围：正阳帝君、威震三界。。台词："做了帝君，我还是那个武将——只不过，从前斩人，如今斩邪。"。动作帧（动图）：①立云端 ②芭蕉扇一扇 ③火海开道 ④八仙随行。诗词：正阳帝君位列仙，一扇火海照三边。从前斩将今斩邪，赤阳不改旧时天。。主题句：沙场猛将弃剑修，一扇火里见丹心——放下屠刀，自成大日。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：正阳红+大日金+霞紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 正阳帝君, age 得道后, 八仙武力之首·正阳帝君; scene: 八仙过海，芭蕉扇一扇，火海开道；受封正阳帝君，不改武将本色。; 正阳帝君位列仙，一扇火海照三边。从前斩将今斩邪，赤阳不改旧时天。; palette: 正阳红+大日金+霞紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 正阳帝君**
- 汉钟离，道果圆满阶段·太阳真火（成仙后，太阳真火·御日行空）。形象：汉钟离，赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。品性：太阳真火圆满，正阳帝君御日行空，照临万物。放下屠刀者，自成大日。。姿态：御日巡天，一扇阳光遍洒；人间逢旱，扇来甘霖；逢夜，扇出暖阳。。服饰：红袍，周身太阳真火，芭蕉扇化作暖阳。。体型：身高约7头身，红袍，御日巡天。。衣物细节：红袍，周身太阳真火，芭蕉扇。。发型妆造：赤发披散，日光为冕。。脸型五官：赤发，浓眉，虎目含日光，虬髯，慈严。。武器招式：芭蕉扇（御日/送暖/甘霖）。。功法：太阳真火御日；芭蕉扇（呼风唤雨/送暖）；照临万物。。功法表现：太阳真火，暖阳遍洒。。画面：构图：天穹之上，红袍帝君御日而行，周身太阳真火，芭蕉扇一挥阳光遍洒，下方人间山川城池沐光，背景朝霞万里。色调：太阳金+真火红+朝霞紫。氛围：御日、道果、照临万物。。台词："我本是杀伐的武将，如今成了照人的大日。这世间冷暖，我都照得。"。动作帧（动图）：①御日巡天 ②芭蕉扇一挥 ③阳光遍洒 ④俯瞰人间。诗词：太阳真火御天行，一扇暖阳照万城。杀伐尽归丹心化，赤阳犹是旧时名。。主题句：沙场猛将弃剑修，一扇火里见丹心——放下屠刀，自成大日。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：太阳金+真火红+朝霞紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 太阳真火, age 成仙后, 太阳真火·御日行空; scene: 御日巡天，一扇阳光遍洒；人间逢旱，扇来甘霖；逢夜，扇出暖阳。; 太阳真火御天行，一扇暖阳照万城。杀伐尽归丹心化，赤阳犹是旧时名。; palette: 太阳金+真火红+朝霞紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 蓝采和（`lan_caihe`） · 人生档案版

**灵胎初醒 · 濠梁乐童**
- 蓝采和，灵胎初醒阶段·市井歌者（少年，市井放歌·乐童）。形象：蓝采和，破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。品性：唐代市井歌者，破绿衫赤足，醉歌于市，歌声清越。。姿态：手持玉板击节，踏歌而唱；唱至兴起，赤足起舞。。服饰：破绿衫，赤足，手持玉板，腰挂铜钱。。体型：身高约6头身，破绿衫赤足少年，清瘦灵动。。衣物细节：破绿衫，赤足，玉板，铜钱。。发型妆造：散乱短发，笑意盈盈。。脸型五官：少年面容，眉眼弯弯，眼含笑意，鼻梁挺，嘴角高歌。。武器招式：无兵器，玉板/歌声。。功法：清歌一曲，声入人心。。功法表现：歌声清越，入人心。。画面：构图：唐代街市，破绿衫赤足少年踏歌而行，玉板击节，行人驻足，背景市井烟火。色调：破衫绿+市井暖+铜钱黄。氛围：放歌、赤足、市井。。台词："世人忙忙碌碌，我自踏歌而行。这歌里的逍遥，你们听不听？"。动作帧（动图）：①玉板击节 ②踏歌而行 ③行人驻足 ④赤足起舞。诗词：市井歌童赤足行，破衫玉板唱天明。一声清越千人醉，笑看红尘几度荣。。主题句：一歌一笑一逍遥，赤足踏遍红尘道——天下熙攘，我只高歌。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衫绿+市井暖+铜钱黄。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 市井歌者, age 少年, 市井放歌·乐童; scene: 手持玉板击节，踏歌而唱；唱至兴起，赤足起舞。; 市井歌童赤足行，破衫玉板唱天明。一声清越千人醉，笑看红尘几度荣。; palette: 破衫绿+市井暖+铜钱黄; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 歌中有道**
- 蓝采和，凡尘砺心阶段·遇仙点化（少年，遇仙·放下执念）。形象：蓝采和，破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。品性：汉钟离见其灵性，点化之，令其放下对"名"的执念——歌为心声，不为唱红。。姿态：梦中历尽歌坛盛名与凄凉，醒来悟得"歌不是唱给别人听的，是唱给自己心的"。。服饰：破绿衫，赤足，玉板。。体型：身高约6头身，破衫赤足，立于月下。。衣物细节：破绿衫，赤足，玉板。。发型妆造：短发。。脸型五官：少年面容，眼神从迷惘到澄澈。。武器招式：无兵器。。功法：勘破名相；歌中见道。。功法表现：歌声入梦，点化之光。。画面：构图：梦中山台，歌童立于幻境中，眼前歌坛盛名如烟散尽，醒来时汉钟离立于月下，背景月夜客店。色调：破衫绿+梦影蓝+月白。氛围：放下、点化、歌者。。台词："从前唱歌，唱给万人听；如今唱歌，唱给自己心。心一净，歌就通了。"。动作帧（动图）：①梦游歌台 ②盛名如烟散 ③惊醒 ④望汉钟离而悟。诗词：点化一声醒梦魂，歌台盛名皆浮云。赤足踏歌心已净，此身原不属红尘。。主题句：一歌一笑一逍遥，赤足踏遍红尘道——天下熙攘，我只高歌。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衫绿+梦影蓝+月白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 遇仙点化, age 少年, 遇仙·放下执念; scene: 梦中历尽歌坛盛名与凄凉，醒来悟得"歌不是唱给别人听的，是唱给自己心的"。; 点化一声醒梦魂，歌台盛名皆浮云。赤足踏歌心已净，此身原不属红尘。; palette: 破衫绿+梦影蓝+月白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 踏歌仙人**
- 蓝采和，道法初成阶段·踏歌云游（少年，踏歌云游·一醉一歌）。形象：蓝采和，破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。品性：从此踏歌云游，一醉一歌，走到哪唱到哪，逍遥天地间。。姿态：玉板一拍，踏歌而行，一步一莲花；歌声所至，愁云散、欢颜开。。服饰：破绿衫，赤足，玉板，铜钱。。体型：身高约6头身，赤足踏歌，身形灵动。。衣物细节：破绿衫，赤足，玉板，铜钱。。发型妆造：短发飘动。。脸型五官：少年面容，眉眼弯弯，高歌含笑。。武器招式：玉板玄音；踏歌成阵。。功法：踏歌成阵（步步生莲）；玉板玄音。。功法表现：步步生莲，歌声清光。。画面：构图：云山古道，赤足少年踏歌而行，脚下步步生莲，歌声化作清光散开，背景青山云海。色调：破衫绿+莲白+山青。氛围：云游、踏歌、逍遥。。台词："天地是舞台，云霞是布景。我走到哪儿，歌就唱到哪儿。"。动作帧（动图）：①玉板一拍 ②踏歌而行 ③步步生莲 ④歌声散作清光。诗词：踏歌云游少年行，一步一莲天地清。玉板敲开愁千丈，歌声过处万山青。。主题句：一歌一笑一逍遥，赤足踏遍红尘道——天下熙攘，我只高歌。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衫绿+莲白+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 踏歌云游, age 少年, 踏歌云游·一醉一歌; scene: 玉板一拍，踏歌而行，一步一莲花；歌声所至，愁云散、欢颜开。; 踏歌云游少年行，一步一莲天地清。玉板敲开愁千丈，歌声过处万山青。; palette: 破衫绿+莲白+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 踏歌成阵**
- 蓝采和，大劫淬炼阶段·歌动人心（青年，歌动人心·化忧为喜）。形象：蓝采和，破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。品性：歌声能化忧为喜，破愁解怨。所到之处，哭者止泪，怒者息怒，哀者展颜。。姿态：一歌使满城愁云散；一歌使干戈化玉帛；散尽铜钱济贫者，歌声不为名利。。服饰：破绿衫，赤足，玉板。。体型：身高约6头身，赤足歌者，立于街中。。衣物细节：破绿衫，赤足，玉板，铜钱。。发型妆造：短发。。脸型五官：青年面容，眉眼含笑，目光温暖。。武器招式：玉板/歌声（言灵）。。功法：言灵之术（歌声化境）；金铜钱化雨散财。。功法表现：歌声光雨，化忧为喜。。画面：构图：城中愁巷，赤足歌者高歌，歌声化作光雨洒落，哭者展颜、怒者平息，铜钱化作光点飞向贫者，背景街巷。色调：破衫绿+歌光金+人间暖。氛围：歌动人心、化忧为喜。。台词："我唱的不是歌，是人心里的那口气。气顺了，什么坎都过得去。"。动作帧（动图）：①举玉板高歌 ②歌声化光雨 ③哭者展颜 ④铜钱化雨济贫。诗词：一歌能使万愁消，化忧为喜乐陶陶。铜钱化雨散寒士，歌到人间百病逃。。主题句：一歌一笑一逍遥，赤足踏遍红尘道——天下熙攘，我只高歌。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衫绿+歌光金+人间暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 歌动人心, age 青年, 歌动人心·化忧为喜; scene: 一歌使满城愁云散；一歌使干戈化玉帛；散尽铜钱济贫者，歌声不为名利。; 一歌能使万愁消，化忧为喜乐陶陶。铜钱化雨散寒士，歌到人间百病逃。; palette: 破衫绿+歌光金+人间暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 九天玄音**
- 蓝采和，封神登天阶段·八仙过海（青年，八仙过海·歌仙）。形象：蓝采和，破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。品性：八仙过海，他踏歌踏浪，玉板一拍海波开道，八仙中的歌仙。。姿态：赤足踏浪，歌声随行；玉板一拍，浪花化作音符铺路。。服饰：破绿衫，赤足，玉板，歌仙气度。。体型：身高约6头身，赤足踏浪，身形轻灵。。衣物细节：破绿衫，赤足，玉板。。发型妆造：短发。。脸型五官：青年面容，眉眼弯弯，高歌含笑。。武器招式：玉板（击浪成音）。。功法：八仙过海；歌中化境。。功法表现：浪花化音符，歌声随行。。画面：构图：沧海之上，赤足歌仙踏浪而歌，浪花化作音符铺路，八仙随行，背景明月沧波。色调：破衫绿+浪花白+月白。氛围：歌仙、踏浪、逍遥。。台词："八仙过海，各显神通。我的神通，就是让浪花跟着我的歌跳舞。"。动作帧（动图）：①赤足踏浪 ②玉板一拍 ③浪化音符铺路 ④踏歌而行。诗词：八仙过海我踏歌，浪花音符逐波波。玉板一拍千重浪，歌仙一唱海天和。。主题句：一歌一笑一逍遥，赤足踏遍红尘道——天下熙攘，我只高歌。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：破衫绿+浪花白+月白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 八仙过海, age 青年, 八仙过海·歌仙; scene: 赤足踏浪，歌声随行；玉板一拍，浪花化作音符铺路。; 八仙过海我踏歌，浪花音符逐波波。玉板一拍千重浪，歌仙一唱海天和。; palette: 破衫绿+浪花白+月白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 濠梁仙翁**
- 蓝采和，道果圆满阶段·大音希声（成仙后，大音希声·以歌化道）。形象：蓝采和，破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。品性：歌道圆满，大音希声——真正的大歌，无声胜有声。天下音乐皆为其化身。。姿态：不再高歌，只以心音传道；风过林梢、溪流过石，皆是他的歌。。服饰：素衫，赤足，玉板敛入袖中，周身音韵。。体型：身高约6头身，素衫赤足，负手山巅。。衣物细节：素衫，赤足，玉板藏袖。。发型妆造：短发。。脸型五官：青年面容，眉眼含笑，眼神澄澈。。武器招式：无兵器，心音为歌。。功法：大音希声；歌道圆满；天地为琴。。功法表现：风过成歌，天地为琴。。画面：构图：山巅风起，素衫歌者负手而立，风过松涛、溪流过石皆化作音符环绕，脚下云海，背景霞光。色调：素衫白+松青+霞光金。氛围：大音希声、以歌化道。。台词："年轻时唱歌，唱给人听；如今唱歌，唱给天地听——天地懂了，便是大音。"。动作帧（动图）：①负手山巅 ②风过松涛 ③音符环绕 ④静听天地。诗词：大音希声绕三界，风过松涛皆我歌。玉板已藏音韵在，无声处听大道和。。主题句：一歌一笑一逍遥，赤足踏遍红尘道——天下熙攘，我只高歌。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素衫白+松青+霞光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 大音希声, age 成仙后, 大音希声·以歌化道; scene: 不再高歌，只以心音传道；风过林梢、溪流过石，皆是他的歌。; 大音希声绕三界，风过松涛皆我歌。玉板已藏音韵在，无声处听大道和。; palette: 素衫白+松青+霞光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 曹国舅（`cao_guojiu`） · 人生档案版

**灵胎初醒 · 锦绣公子**
- 曹国舅，灵胎初醒阶段·皇亲国舅（青年，锦衣玉食·皇亲国舅）。形象：曹国舅，白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。品性：曹太后之弟曹景休，皇亲国戚，锦衣玉食，富贵满门。。姿态：华服锦袍，出入宫闱；赏花品茶，不知民间疾苦。。服饰：锦袍玉带，头戴莲花冠，富家公子气度。。体型：身高约7头身，锦袍玉带，雍容华贵。。衣物细节：锦袍玉带，莲花冠。。发型妆造：束发戴冠，富贵之气。。脸型五官：面容俊逸，眉目温润，眼含富贵气，鼻梁挺，嘴角含笑。。武器招式：无兵器。。功法：通音律，善礼乐；无道术。。功法表现：无神力，富贵之气。。画面：构图：富贵庭院，锦袍公子凭栏赏花，华服映日，侍从环绕，背景金楼玉阶。色调：锦袍金+富贵紫+华庭红。氛围：富贵、华服、不知疾苦。。台词："生来便是金枝玉叶，这世间，大约没有什么是我得不到的。"。动作帧（动图）：①凭栏赏花 ②品茶听曲 ③侍从环绕 ④低头沉思。诗词：皇亲国舅锦衣行，玉带金冠富贵荣。笙歌宴罢高楼月，不识人间有哭声。。主题句：生在富贵帝王家，偏寻云外清净地——放下荣华，方得大道。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：锦袍金+富贵紫+华庭红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 皇亲国舅, age 青年, 锦衣玉食·皇亲国舅; scene: 华服锦袍，出入宫闱；赏花品茶，不知民间疾苦。; 皇亲国舅锦衣行，玉带金冠富贵荣。笙歌宴罢高楼月，不识人间有哭声。; palette: 锦袍金+富贵紫+华庭红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 青布道袍**
- 曹国舅，凡尘砺心阶段·良心难安（青年，见弟为恶·良心难安）。形象：曹国舅，白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。品性：见胞弟仗势欺人、为恶乡里，锦衣玉食之下，良心渐生不安。。姿态：听百姓暗骂曹家，垂首无言；劝弟不听，辗转难眠。。服饰：锦袍依旧，眉间却多了忧色。。体型：身高约7头身，锦袍独坐，眉间忧色。。衣物细节：锦袍玉带。。发型妆造：束发，眉间忧色。。脸型五官：俊逸面容，眉微蹙，眼含愧意，嘴角抿紧。。武器招式：无兵器。。功法：无道术，唯有未泯的良知。。功法表现：无神力，良心之灼。。画面：构图：华庭月夜，锦袍公子独坐，烛火明灭，窗外隐约百姓指点的影子，神情愧怍，背景金屋生寒。色调：锦袍金+夜月白+愧影青。氛围：不安、愧怍、见弟为恶。。台词："这身荣华，原是踩着别人的苦换来的。我穿着它，越来越烫。"。动作帧（动图）：①听民怨 ②垂首无言 ③劝弟不听 ④夜中独坐难眠。诗词：锦衣公子夜难眠，胞弟为恶民怨连。玉食锦衣如芒刺，良心一寸怎安眠。。主题句：生在富贵帝王家，偏寻云外清净地——放下荣华，方得大道。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：锦袍金+夜月白+愧影青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 良心难安, age 青年, 见弟为恶·良心难安; scene: 听百姓暗骂曹家，垂首无言；劝弟不听，辗转难眠。; 锦衣公子夜难眠，胞弟为恶民怨连。玉食锦衣如芒刺，良心一寸怎安眠。; palette: 锦袍金+夜月白+愧影青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 玉磬仙人**
- 曹国舅，道法初成阶段·散家弃官（中年，散尽家财·弃官修道）。形象：曹国舅，白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。品性：看破富贵如浮云，散尽家财，弃官入道，寻访清净。。姿态：将家产尽散与贫民，脱下锦袍，青衫入山；宫门回望，释然而笑。。服饰：青衫素衣，弃冠散发。。体型：身高约7头身，青衫素衣，身形清减。。衣物细节：青衫素衣，弃冠散发。。发型妆造：散发，道气初现。。脸型五官：清俊面容，眉目释然，眼含宁静，嘴角浅笑。。武器招式：无兵器。。功法：弃富贵；道心初立。。功法表现：无神力，心静如水。。画面：构图：宫门之外，脱下锦袍的公子青衫散发，将金袋递给贫民，回望宫阙释然一笑，转身入山。色调：青衫+贫民暖+山青。氛围：放下、弃官、释然。。台词："这些金子，买不来一夜安眠。我放下它们，才睡得着觉。"。动作帧（动图）：①散尽家财 ②脱锦袍 ③回望宫阙 ④转身入山。诗词：散尽家财弃官袍，青衫一袭入云涛。回头宫阙金如土，不及山间一树桃。。主题句：生在富贵帝王家，偏寻云外清净地——放下荣华，方得大道。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：青衫+贫民暖+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 散家弃官, age 中年, 散尽家财·弃官修道; scene: 将家产尽散与贫民，脱下锦袍，青衫入山；宫门回望，释然而笑。; 散尽家财弃官袍，青衫一袭入云涛。回头宫阙金如土，不及山间一树桃。; palette: 青衫+贫民暖+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 礼乐大道**
- 曹国舅，大劫淬炼阶段·清修苦行（中年，清修苦行·礼乐自持）。形象：曹国舅，白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。品性：山中清修苦行，以礼乐自持。玉磬随身，一磬清音，涤荡尘心。。姿态：晨起击磬，清音绕山；日夜苦修，礼乐大道渐成。。服饰：青布道衣，手持玉磬，白衣鹤氅渐具。。体型：身高约7头身，白衣鹤氅，清瘦出尘。。衣物细节：青布道衣/白衣鹤氅，玉磬。。发型妆造：道髻，清癯。。脸型五官：清俊面容，眉目淡泊，眼含道意，嘴角宁静。。武器招式：玉磬（清音诀）。。功法：玉磬神音（清音诀）；礼乐大道。。功法表现：磬音涤尘，礼乐之道。。画面：构图：山中草庐，白衣道者晨击玉磬，清音化作音波散开，涤荡山间，背景晨曦松影。色调：白衣+磬光玉白+山青。氛围：清修、礼乐、苦行。。台词："富贵的时候，我耳朵里全是丝竹。如今清修，才听见磬声里的大道。"。动作帧（动图）：①晨起击磬 ②清音绕山 ③日夜苦修 ④礼乐入道。诗词：山中苦行礼乐勤，玉磬清音涤垢尘。富贵早随云去远，一磬一声见道真。。主题句：生在富贵帝王家，偏寻云外清净地——放下荣华，方得大道。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：白衣+磬光玉白+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 清修苦行, age 中年, 清修苦行·礼乐自持; scene: 晨起击磬，清音绕山；日夜苦修，礼乐大道渐成。; 山中苦行礼乐勤，玉磬清音涤垢尘。富贵早随云去远，一磬一声见道真。; palette: 白衣+磬光玉白+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 九天神音**
- 曹国舅，封神登天阶段·八仙过海（得道后，八仙过海·玉磬仙人）。形象：曹国舅，白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。品性：八仙过海，他玉磬一击，清音开道，八仙中的玉磬仙人。。姿态：磬声一起，海波为之让路；一磬清音，天地皆静。。服饰：白衣鹤氅，莲花冠，手持玉磬。。体型：身高约7头身，白衣鹤氅，击磬凌波。。衣物细节：白衣鹤氅，莲花冠，玉磬。。发型妆造：束发莲花冠。。脸型五官：清俊面容，眉目宁静，眼含雅意。。武器招式：玉磬（震东海）。。功法：磬震东海；玉磬神音大成。。功法表现：磬音开道，海波让路。。画面：构图：沧海之上，白衣道人击磬而行，磬音化作清波开道，八仙随行，背景海天澄净。色调：白衣+磬光玉白+沧溟蓝。氛围：玉磬、雅音、清净。。台词："八仙过海，我这玉磬一声，比什么神通都清净。"。动作帧（动图）：①举磬 ②一击 ③磬音开道 ④随八仙渡海。诗词：八仙过海磬声悠，一磬清音镇海流。玉磬仙人行处静，白云随我渡沧洲。。主题句：生在富贵帝王家，偏寻云外清净地——放下荣华，方得大道。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：白衣+磬光玉白+沧溟蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 八仙过海, age 得道后, 八仙过海·玉磬仙人; scene: 磬声一起，海波为之让路；一磬清音，天地皆静。; 八仙过海磬声悠，一磬清音镇海流。玉磬仙人行处静，白云随我渡沧洲。; palette: 白衣+磬光玉白+沧溟蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 国舅真仙**
- 曹国舅，道果圆满阶段·磬音清心（成仙后，磬音清心·以音律化人）。形象：曹国舅，白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。品性：磬音圆满，以音律化人。玉磬清音绕三界，众生闻之，嗔念自消。。姿态：云游三界，一磬清音，人间烦愁尽涤；玉磬化作清光，以礼乐度众生。。服饰：白衣鹤氅，玉磬化作清光随身，儒仙雅音。。体型：身高约7头身，白衣鹤氅，击磬云端。。衣物细节：白衣鹤氅，玉磬化清光。。发型妆造：束发莲花冠。。脸型五官：清俊面容，眉目慈悲，眼含温润。。武器招式：玉磬（化道清音）。。功法：磬音化道；以音律化人；儒仙雅音。。功法表现：磬音化清光，涤尽尘愁。。画面：构图：云海之上，白衣儒仙击磬，磬音化作漫天清光洒向人间，下方尘世愁云尽散，背景霞光万里。色调：白衣+磬光玉白+霞光金。氛围：磬音清心、以音化人。。台词："从前的丝竹是享乐，如今的磬音是渡人。这一声清音，愿众生皆闻。"。动作帧（动图）：①举磬云端 ②一击 ③磬音化清光 ④洒向人间。诗词：磬音清心绕三界，一声涤尽万般埃。儒仙化道音为渡，云外清音渡众来。。主题句：生在富贵帝王家，偏寻云外清净地——放下荣华，方得大道。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：白衣+磬光玉白+霞光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 磬音清心, age 成仙后, 磬音清心·以音律化人; scene: 云游三界，一磬清音，人间烦愁尽涤；玉磬化作清光，以礼乐度众生。; 磬音清心绕三界，一声涤尽万般埃。儒仙化道音为渡，云外清音渡众来。; palette: 白衣+磬光玉白+霞光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 太上老君（`taishang_laojun`） · 人生档案版

**灵胎初醒 · 上古真人**
- 太上老君，灵胎初醒阶段·一气化生（开天前，混沌初开·一气化生）。形象：太上老君，白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。品性：混沌初开之先的道祖，一气化生，分判阴阳，开天辟地。。姿态：身周紫气缭绕，一气化作清浊；指点间，天地初分。。服饰：紫金道袍，白发童颜，身绕紫气。。体型：身高约7头身，紫金道袍，白发童颜，身形伟岸。。衣物细节：紫金道袍，身绕紫气。。发型妆造：白发，白须。。脸型五官：童颜鹤发，眉目慈悲，眼含无量道意，鼻梁挺，嘴角含笑。。武器招式：无兵器，拂尘/拈诀。。功法：一气化生；开天辟地之力。。功法表现：一气化生，紫气东来。。画面：构图：混沌虚空，白发道祖立于紫气之中，一手拈诀，指尖清气上升浊气下沉，天地初分，背景浩瀚星海。色调：紫气+星海玄+初光金。氛围：开天、道祖、混沌初分。。台词："道生一，一生二，二生三，三生万物。这天地，是我一口气吹开的。"。动作帧（动图）：①立于混沌 ②拈诀 ③清气升浊气降 ④天地初分。诗词：混沌初开一气分，清升浊降定乾坤。紫气东来三万里，道生万物始见真。。主题句：道生一，一生二，二生三，三生万物——无为而无不为。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫气+星海玄+初光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 一气化生, age 开天前, 混沌初开·一气化生; scene: 身周紫气缭绕，一气化作清浊；指点间，天地初分。; 混沌初开一气分，清升浊降定乾坤。紫气东来三万里，道生万物始见真。; palette: 紫气+星海玄+初光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 紫气东来**
- 太上老君，凡尘砺心阶段·青牛出关（成道后，骑青牛出函谷·五千言化世）。形象：太上老君，白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。品性：骑青牛出函谷关，应关令尹喜之请，留下《道德经》五千言，道化众生。。姿态：青牛缓缓，西出函谷；尹喜求道，他留下五千言，随即骑牛远去。。服饰：紫金道袍，骑青牛，手持拂尘。。体型：身高约7头身，骑青牛，道祖气度。。衣物细节：紫金道袍，青牛，拂尘。。发型妆造：白发，白须。。脸型五官：童颜鹤发，眉目含笑，眼神深邃。。武器招式：无兵器，拂尘。。功法：五千言道；道德化世。。功法表现：五千言化金光，道满人间。。画面：构图：函谷关古道，白发道祖骑青牛缓缓而行，身后五千言化作金色文字洒向人间，关尹喜伏拜，背景落日关隘。色调：紫袍+青牛+经文化金。氛围：五千言、化世、骑牛。。台词："道可道，非常道。这五千字，够你们参三千年了。"。动作帧（动图）：①骑牛出关 ②尹喜求道 ③留下五千言 ④骑牛远去。诗词：青牛西出函谷关，五千言落道满天。尹喜求得真经在，老子骑牛云外还。。主题句：道生一，一生二，二生三，三生万物——无为而无不为。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫袍+青牛+经文化金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 青牛出关, age 成道后, 骑青牛出函谷·五千言化世; scene: 青牛缓缓，西出函谷；尹喜求道，他留下五千言，随即骑牛远去。; 青牛西出函谷关，五千言落道满天。尹喜求得真经在，老子骑牛云外还。; palette: 紫袍+青牛+经文化金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 八卦炉**
- 太上老君，道法初成阶段·太清炼丹（太清境，兜率宫·九转金丹）。形象：太上老君，白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。品性：居兜率宫（太清境），炼九转金丹，掌天地丹道。。姿态：八卦炉前，扇动炉火，九转金丹光华流转；偶开炉，丹光映天。。服饰：紫金道袍，白发，手持芭蕉扇（炼丹），八卦炉前。。体型：身高约7头身，紫金道袍，炉前而立。。衣物细节：紫金道袍，芭蕉扇，八卦炉。。发型妆造：白发，白须。。脸型五官：童颜鹤发，眉目含丹意，眼含炉火之光。。武器招式：芭蕉扇（炼丹/太清神雷）。。功法：九转金丹；太上丹道。。功法表现：九转金丹，丹光映天。。画面：构图：兜率宫八卦炉前，白发道祖执扇炼丹，炉火升腾，九转金丹丹光流转，背景宫阙云海。色调：炉火红+金丹金+太清紫。氛围：炼丹、太清、道祖。。台词："这炉子里的，不是丹药，是天地精气。炼到九转，便是一颗可化万物的道。"。动作帧（动图）：①执扇扇炉 ②炉火升腾 ③丹光流转 ④开炉丹光映天。诗词：兜率宫中炉火红，九转金丹照太清。一炉天地精气聚，万灵服此道常明。。主题句：道生一，一生二，二生三，三生万物——无为而无不为。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：炉火红+金丹金+太清紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 太清炼丹, age 太清境, 兜率宫·九转金丹; scene: 八卦炉前，扇动炉火，九转金丹光华流转；偶开炉，丹光映天。; 兜率宫中炉火红，九转金丹照太清。一炉天地精气聚，万灵服此道常明。; palette: 炉火红+金丹金+太清紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 西游护法**
- 太上老君，大劫淬炼阶段·八卦炉护法（封神，八卦炉前观劫·暗中护法）。形象：太上老君，白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。品性：封神劫起，他于八卦炉前观劫，暗中护法，点化有缘，为西岐之势暗中添柴。。姿态：炉中观三界劫数，一指点化玄机；授金刚琢、紫金红葫芦与门人，护法护道。。服饰：紫金道袍，白发，八卦炉前，法宝悬身。。体型：身高约7头身，紫金道袍，炉前观劫。。衣物细节：紫金道袍，八卦炉，金刚琢/紫金葫芦。。发型妆造：白发，白须。。脸型五官：童颜鹤发，眉目深邃，眼含天机。。武器招式：金刚琢（收物）；紫金红葫芦（纳妖）。。功法：金刚琢；紫金红葫芦；八卦炉观劫。。功法表现：炉中显劫，紫气护法。。画面：构图：兜率宫八卦炉前，道祖临炉观劫，炉中显出三界战火之象，身侧金刚琢紫金葫芦悬浮，背景紫气宫阙。色调：炉火红+紫气+战影金。氛围：观劫、护法、师祖。。台词："这场大劫，是天数，也是人心。我不下场，只把这炉火烧得旺些，好炼出该成的人。"。动作帧（动图）：①临炉观劫 ②炉中显战象 ③一指点化 ④法宝悬身。诗词：八卦炉前观劫尘，暗中添柴护道真。金刚琢落妖魔服，葫芦一收万法遵。。主题句：道生一，一生二，二生三，三生万物——无为而无不为。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：炉火红+紫气+战影金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 八卦炉护法, age 封神, 八卦炉前观劫·暗中护法; scene: 炉中观三界劫数，一指点化玄机；授金刚琢、紫金红葫芦与门人，护法护道。; 八卦炉前观劫尘，暗中添柴护道真。金刚琢落妖魔服，葫芦一收万法遵。; palette: 炉火红+紫气+战影金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 道德天尊**
- 太上老君，封神登天阶段·三清之首（封神后，三清之首·万法道祖）。形象：太上老君，白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。品性：三清之首，道德天尊，万法道祖。一炁化三清，道满天地。。姿态：端坐太清境，拂尘一挥，紫气东来三万里；万仙来朝，道法无边。。服饰：紫金道袍，白发童颜，拂尘在手，身绕紫气。。体型：身高约7头身，端坐太清境，道祖威仪。。衣物细节：紫金道袍，拂尘，紫气。。发型妆造：白发，白须。。脸型五官：童颜鹤发，眉目慈悲，眼含无量道。。武器招式：拂尘（挥紫气）；金刚琢/葫芦。。功法：一炁化三清；万法道祖；紫气东来。。功法表现：紫气东来，三清虚影。。画面：构图：太清境中，紫金道祖端坐，拂尘一挥紫气东来，身后三清虚影，万仙朝拜，背景云海紫霞。色调：紫金+紫气+圣白。氛围：三清之首、万法道祖。。台词："三清本是一气，万法本是一道。你们拜我，其实拜的是那个"道"。"。动作帧（动图）：①端坐 ②拂尘一挥 ③紫气东来 ④万仙朝拜。诗词：三清之首道德尊，一气化生道满门。紫气东来三万里，万法归宗本一根。。主题句：道生一，一生二，二生三，三生万物——无为而无不为。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金+紫气+圣白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 三清之首, age 封神后, 三清之首·万法道祖; scene: 端坐太清境，拂尘一挥，紫气东来三万里；万仙来朝，道法无边。; 三清之首道德尊，一气化生道满门。紫气东来三万里，万法归宗本一根。; palette: 紫金+紫气+圣白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 三清合一**
- 太上老君，道果圆满阶段·道化天地（道之化身，一炁化三清·道化万物）。形象：太上老君，白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。品性：已成为道本身——日升月落、春生秋杀皆是老君。无为而无不为，无处不是道。。姿态：一缕紫气化作日月星辰；一阵风即是他的道行。不显不隐，与万物同在。。服饰：道袍隐于万物，或现白发道祖本相，身化紫气。。体型：身高约7头身，道祖本相，身化紫气融入天地。。衣物细节：道袍隐于万物，本相时紫金道袍。。发型妆造：白发，白须。。脸型五官：童颜鹤发，眉目含笑，眼含整个天地。。武器招式：无兵器，道即是法。。功法：道化万物；无为而无不为；与天地同体。。功法表现：紫气化天地，无处不在。。画面：构图：天地之间，一缕紫气化作日月星辰、山川河流，白发道祖身影融入其间若隐若现，背景整个宇宙星海。色调：紫气+星海+山河青。氛围：道化天地、无处不在。。台词："你问我住哪？我不住天上，不住人间——我住在风里、水里、你抬头看的云里。道，无处不在。"。动作帧（动图）：①一缕紫气 ②化日月星辰 ③身影融天地 ④与万物同在。诗词：道化天地无形身，日月星辰皆是君。春生秋杀无非道，无为而治万物亲。。主题句：道生一，一生二，二生三，三生万物——无为而无不为。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫气+星海+山河青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 道化天地, age 道之化身, 一炁化三清·道化万物; scene: 一缕紫气化作日月星辰；一阵风即是他的道行。不显不隐，与万物同在。; 道化天地无形身，日月星辰皆是君。春生秋杀无非道，无为而治万物亲。; palette: 紫气+星海+山河青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 钟馗（`zhong_kui`） · 人生档案版

**灵胎初醒 · 终南寒窗**
- 钟馗，灵胎初醒阶段·终南寒窗（青年，寒门书生·饱读诗书）。形象：钟馗，赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。品性：终南山下寒门书生，貌虽丑，腹中锦绣。寒窗苦读十余载，胸有济世之志。。姿态：挑灯夜读，批注典籍；山间采药换纸墨，日复一日，笔耕不辍。。服饰：青布旧袍，布巾束发，面庞粗犷却目光清正。。体型：身高约7头身，粗犷书生，青布旧袍。。衣物细节：青布旧袍，布巾束发。。发型妆造：布巾束发，面庞粗犷目光清正。。脸型五官：国字脸粗犷，粗眉，豹眼清正，狮鼻阔口，络腮胡，下巴方正。。武器招式：无兵器，笔墨为伴。。功法：饱读经史，胸有锦绣；写得一手好文章。。功法表现：无神力，文气充盈。。画面：构图：终南山茅屋灯下，粗犷书生伏案夜读，烛火摇曳，窗外山月一轮；背景群山夜色。色调：灯暖黄+旧袍青+月白。氛围：寒窗、志气、身丑心正。。台词："皮囊丑些怕什么？我钟馗要凭这身学问，换一个干干净净的功名。"。动作帧（动图）：①挑灯夜读 ②批注典籍 ③采药换纸墨 ④伏案疾书。诗词：终南山下夜灯明，寒士挑灯读到更。貌丑何妨心锦绣，明朝金榜问功名。。主题句：生不让我中状元，死后我掌生死簿——以丑驱万鬼，以正镇诸邪。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：灯暖黄+旧袍青+月白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 灵胎初醒 终南寒窗, age 青年, 寒门书生·饱读诗书; scene: 挑灯夜读，批注典籍；山间采药换纸墨，日复一日，笔耕不辍。; 终南山下夜灯明，寒士挑灯读到更。貌丑何妨心锦绣，明朝金榜问功名。; palette: 灯暖黄+旧袍青+月白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 触阶而亡**
- 钟馗，凡尘砺心阶段·殿试被黜（青年，貌丑被黜·功名尽失）。形象：钟馗，赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。品性：进京殿试，文章夺魁，却因面貌丑陋被皇帝黜落。十年寒窗，一朝成空，他悲愤难平。。姿态：金榜无名，他在殿外立了整夜；见同榜者春风得意，他只攥紧了拳。。服饰：新制的青衫，在宫门前显得格外寒酸。。体型：身高约7头身，青衫在宫门前寒酸。。衣物细节：新制青衫。。发型妆造：束发，面色沉郁。。脸型五官：国字脸，粗眉紧锁，豹眼含愤，狮鼻，阔口抿紧，络腮胡，下巴方正。。武器招式：无兵器。。功法：满腹才学无处用；一腔悲愤难自抑。。功法表现：悲愤之气。。画面：构图：皇宫殿门之前，落第书生独自站在金榜前，人群散去，他背对着巍峨宫阙，身影落寞；背景夕阳宫墙。色调：宫墙红+旧袍青+夕照金。氛围：失意、愤懑、一腔孤愤。。台词："我以文章论忠奸，天下却以貌取人！这功名，不考也罢！"。动作帧（动图）：①金榜无名 ②殿外立整夜 ③见同榜春风得意 ④攥紧拳。诗词：殿试文章第一筹，只缘貌丑黜金瓯。十年寒窗皆付水，书生一愤动神州。。主题句：生不让我中状元，死后我掌生死簿——以丑驱万鬼，以正镇诸邪。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：宫墙红+旧袍青+夕照金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 凡尘砺心 殿试被黜, age 青年, 貌丑被黜·功名尽失; scene: 金榜无名，他在殿外立了整夜；见同榜者春风得意，他只攥紧了拳。; 殿试文章第一筹，只缘貌丑黜金瓯。十年寒窗皆付水，书生一愤动神州。; palette: 宫墙红+旧袍青+夕照金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 地府扬名**
- 钟馗，道法初成阶段·触阶含怨（青年，触阶而亡·怨气化鬼）。形象：钟馗，赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。品性：悲愤之下，触阶而亡。怨气不散，化为一缕厉鬼，徘徊不去，要问这世间何为公道。。姿态：在宫门前触阶，血溅丹墀；魂魄凝而不散，怒目向天。。服饰：青衫尽染，怒气满身，形貌愈发骇人。。体型：身高约7头身，青衫染血，魂魄凝形。。衣物细节：青衫尽染，怒气满身。。发型妆造：散发，面目骇人。。脸型五官：国字脸，粗眉怒竖，豹眼圆睁含怨，狮鼻，阔口龇牙，络腮胡怒张。。武器招式：无兵器。。功法：怨气凝形；一股执念不灭。。功法表现：怨气冲天，形貌愈骇。。画面：构图：宫门丹墀之上，触阶而亡的书生魂魄凝成虚影，青衫染血，怒目向天，宫阙在身后扭曲；背景昏天黑地。色调：血绛+墨黑+惨白。氛围：冤屈、愤恨、怨气冲天。。台词："苍天无眼，以貌取人！我钟馗生不能伸冤，死也要问一问这公道在哪！"。动作帧（动图）：①触阶 ②血溅丹墀 ③魂魄凝而不散 ④怒目向天。诗词：一怒触阶血未干，书生魂魄恨难安。以貌取人天不公，化作厉鬼问人间。。主题句：生不让我中状元，死后我掌生死簿——以丑驱万鬼，以正镇诸邪。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：血绛+墨黑+惨白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道法初成 触阶含怨, age 青年, 触阶而亡·怨气化鬼; scene: 在宫门前触阶，血溅丹墀；魂魄凝而不散，怒目向天。; 一怒触阶血未干，书生魂魄恨难安。以貌取人天不公，化作厉鬼问人间。; palette: 血绛+墨黑+惨白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 驱鬼之王**
- 钟馗，大劫淬炼阶段·地府授职（死后，阎王惜才·判官笔青锋剑）。形象：钟馗，赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。品性：魂魄至地府，阎王惜其刚正，授判官笔、青锋剑，命其捉拿天下恶鬼——生前屈才，死后得用。。姿态：领命之日，判官笔一落，生死簿上善恶分明；青锋剑出鞘，地府鬼魅俯首。。服饰：红袍乌纱，赤面虬髯，判官笔与青锋剑随身。。体型：身高约7头身，红袍乌纱，赤面虬髯。。衣物细节：红袍乌纱，判官笔与青锋剑。。发型妆造：乌纱帽，赤面虬髯。。脸型五官：国字脸赤面，粗眉，豹眼如电，狮鼻，阔口，虬髯如戟，下巴方正。。武器招式：判官笔定善恶；青锋剑斩恶鬼。。功法：判官笔（定善恶）；青锋剑（斩恶鬼）；镇鬼神通。。功法表现：判官笔金光，青锋剑寒光。。画面：构图：地府森罗殿前，红袍鬼官持判官笔立于殿阶，青锋剑横侧，身后众鬼俯首；背景幽蓝地府。色调：地府幽蓝+红袍+生死簿金。氛围：初掌权柄、刚正不阿。。台词："生前没让我中状元，死后倒让我掌了生死簿。也好——善恶报应，我来定！"。动作帧（动图）：①领命 ②判官笔一落善恶分明 ③青锋剑出鞘 ④地府鬼魅俯首。诗词：地府阎君惜此才，判官笔授青锋开。生前未遂凌云志，死后掌得生死台。。主题句：生不让我中状元，死后我掌生死簿——以丑驱万鬼，以正镇诸邪。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：地府幽蓝+红袍+生死簿金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 大劫淬炼 地府授职, age 死后, 阎王惜才·判官笔青锋剑; scene: 领命之日，判官笔一落，生死簿上善恶分明；青锋剑出鞘，地府鬼魅俯首。; 地府阎君惜此才，判官笔授青锋开。生前未遂凌云志，死后掌得生死台。; palette: 地府幽蓝+红袍+生死簿金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 五方鬼卒**
- 钟馗，封神登天阶段·万鬼皈依（驱魔年间，驱魔真君·五方鬼卒）。形象：钟馗，赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。品性：捉鬼无数，威名日盛。恶鬼闻其名而丧胆，愿归正的鬼反而俯首听命，麾下聚起五方鬼卒，护一方安宁。。姿态：夜巡街巷，判官笔点鬼、青锋剑斩邪；所过之处，阴气肃清，百姓夜不闭户。。服饰：紫金冠，红袍猎猎，判官笔与青锋剑在手，五方鬼卒随行。。体型：身高约7头身，紫金冠红袍，魁梧。。衣物细节：紫金冠，红袍猎猎，判官笔青锋剑，五方鬼卒。。发型妆造：紫金冠，赤面虬髯怒张。。脸型五官：国字脸赤面，粗眉，豹眼怒张，狮鼻，阔口，虬髯怒张，紫金冠下威严。。武器招式：判官笔点/青锋剑诛。。功法：驱魔真君神通；判官笔点善恶；青锋剑诛邪；号令五方鬼卒。。功法表现：驱邪金光，阴气肃清。。画面：构图：深夜街巷，红袍真君持判官笔立于屋脊，青锋剑出鞘，身后五方鬼卒黑影，前方阴气退散；背景月色瓦檐。色调：夜空墨+红袍+驱邪金光。氛围：镇宅、驱邪、威震幽冥。。台词："恶鬼我抓，善鬼我留。这生死簿上的善恶，我钟馗，一笔记一个准！"。动作帧（动图）：①夜巡街巷 ②判官笔点鬼 ③青锋剑斩邪 ④五方鬼卒随行。诗词：驱魔真君镇八方，判官笔落鬼魅慌。五方鬼卒听号令，夜巡街巷护一方。。主题句：生不让我中状元，死后我掌生死簿——以丑驱万鬼，以正镇诸邪。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：夜空墨+红袍+驱邪金光。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 封神登天 万鬼皈依, age 驱魔年间, 驱魔真君·五方鬼卒; scene: 夜巡街巷，判官笔点鬼、青锋剑斩邪；所过之处，阴气肃清，百姓夜不闭户。; 驱魔真君镇八方，判官笔落鬼魅慌。五方鬼卒听号令，夜巡街巷护一方。; palette: 夜空墨+红袍+驱邪金光; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 驱魔真君**
- 钟馗，道果圆满阶段·钟馗捉鬼（神位长存，判官笔镇宅·慈心驱鬼）。形象：钟馗，赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。品性：终成门神、驱魔真君，千家万户贴其像镇宅。他捉了一辈子鬼，最后明白：鬼由心生，捉鬼也是度心——怒目之下，是一颗慈悲心。。姿态：镇宅守门，邪祟不敢近；遇迷途恶鬼，判官笔一点，渡其往生；见人间孩童被噩梦惊扰，会轻轻哼一支镇魂调。。服饰：门神红袍，乌纱帽，判官笔与青锋剑，怒目而含慈悲。。体型：身高约7头身，门神红袍，怒目含笑。。衣物细节：门神红袍，乌纱帽，判官笔青锋剑。。发型妆造：乌纱帽，怒目含慈悲。。脸型五官：国字脸，粗眉舒展，豹眼含慈，狮鼻，阔口含笑，虬髯，怒目即菩萨。。武器招式：判官笔定善恶、渡鬼魂。。功法：驱魔真君神位；判官笔定善恶、渡鬼魂；青锋剑诛邪护生。。功法表现：镇宅金光，温暖护佑。。画面：构图：百姓门庭之上，门神红袍执笔而立，怒目含笑，脚下小鬼退避，身后万家灯火安宁；背景夜色与暖窗。色调：门神红+镇宅金+夜色墨。氛围：镇宅、慈悲、护佑人间。。台词："我这一笔，落下去是善恶，抬起来是慈悲。鬼也好，人也好，心正了，门就开。"。动作帧（动图）：①镇宅守门 ②遇迷途恶鬼判官笔度化 ③见孩童噩梦轻哼镇魂调 ④万家灯火安。诗词：门神一像镇千门，判官笔落善恶分。怒目原藏菩萨意，捉鬼亦是度心魂。。主题句：生不让我中状元，死后我掌生死簿——以丑驱万鬼，以正镇诸邪。。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：门神红+镇宅金+夜色墨。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; 道果圆满 钟馗捉鬼, age 神位长存, 判官笔镇宅·慈心驱鬼; scene: 镇宅守门，邪祟不敢近；遇迷途恶鬼，判官笔一点，渡其往生；见人间孩童被噩梦惊扰，会轻轻哼一支镇魂调。; 门神一像镇千门，判官笔落善恶分。怒目原藏菩萨意，捉鬼亦是度心魂。; palette: 门神红+镇宅金+夜色墨; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 3. 宝可梦（12 物种）

> **风格**：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。**阶段演绎**：
> - 灵胎初醒：蛋壳微光，新生命初醒，温暖亲近的孵化氛围（奶油/嫩粉）
> - 凡尘砺心：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮（晴蓝/暖橙）
> - 道法初成：初次对战或训练，火花四溅，能力初显的高光时刻（活力橙/电光蓝）
> - 大劫淬炼：进化之光笼罩，身体在光中蜕变，挣扎又坚定（深海蓝/战斗红）
> - 封神登天：完全进化形态，可靠的战斗伙伴，自信昂扬（伙伴金/烈焰橙）
> - 道果圆满：超进化或极致形态，圣光守护，散发传说级气场（圣光白/彩虹金）

#### 小火龙（`charmander`） · 人生档案版

**灵胎初醒 · 火纹蛋**
- 小火龙，灵胎初醒阶段·火苗幼蜥（出生，尾尖火焰·活泼好斗）。形象：小火龙，橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。品性：火系宝可梦，橙色小蜥蜴，尾尖燃着火焰。活泼好斗，尾巴的火焰随心情起伏，越兴奋烧得越旺。。姿态：追着蝴蝶跑，尾尖火焰欢快跳跃；看见石堆，鼓起腮帮一记"火花"，烤得石头滋滋响。。服饰：无衣，橙红鳞片，腹部米黄，尾尖火焰。。体型：身高约2头身，橙色小蜥蜴，身形活泼，尾长。。衣物细节：通体橙红鳞片，腹部米黄，尾尖燃橘红火焰。。发型妆造：头顶无发，背脊小小凸起。。脸型五官：圆眼灵动，口常张露小火苗，腮鼓。。武器招式：火花初窥，鼓起腮帮。。功法：火花初窥；尾焰随心情。。功法表现：尾尖火焰随心情起伏。。画面：构图：草地，橙色小蜥蜴追蝶，尾尖火焰欢快跳跃，烤石"滋滋"，背景晴空。色调：橙红+焰火橘+草青。氛围：活泼、好斗、火苗旺。。台词："我这尾尖的火，是生命之火——我越开心它越旺，我越有劲儿！谁也别想让它灭。"。动作帧（动图）：①追蝶 ②尾焰欢跳 ③鼓起腮帮喷火花 ④烤石滋滋。诗词：橙红蜥蜴尾焰扬，追蝶一记火花香。尾尖之火随心跳，快乐烧得愈汪洋。。主题句：尾尖之火不熄，是生命也是守护——从火苗到圣焰，燃尽一切黑暗。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：橙红+焰火橘+草青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 灵胎初醒 火苗幼蜥, age 出生, 尾尖火焰·活泼好斗; scene: 追着蝴蝶跑，尾尖火焰欢快跳跃；看见石堆，鼓起腮帮一记"火花"，烤得石头滋滋响。; 橙红蜥蜴尾焰扬，追蝶一记火花香。尾尖之火随心跳，快乐烧得愈汪洋。; palette: 橙红+焰火橘+草青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 蹒跚火苗**
- 小火龙，凡尘砺心阶段·火熄之惧（幼年，风暴之夜·火差点熄灭）。形象：小火龙，橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。品性：一场风暴，雨水浇得它尾尖之火摇摇欲熄。它拼命护着火，却第一次感到——"我的火，也会灭。"。姿态：风雨里，它蜷起身子用爪子护着尾焰，火苗在掌间忽明忽暗；等雨停，它望着重新燃起的火，眼眶发热。。服饰：橙红鳞片被雨打湿，尾焰微弱。。体型：身高约2头身，橙色小蜥蜴，蜷身护焰。。衣物细节：橙红鳞片淋湿，尾焰微弱。。发型妆造：无发，鳞片湿亮。。脸型五官：圆眼含忧，死死护着火。。武器招式：控火不足。。功法：控火不足；火熄之惧。。功法表现：尾焰微弱，风中摇曳。。画面：构图：风雨夜，橙色小蜥蜴蜷身护着尾焰，火苗在掌间忽明忽暗，背景狂风暴雨。色调：橙红+雨夜蓝+微弱焰火。氛围：火熄之惧、护焰、立志。。台词："我的火，是我的命。它要是灭了，我还算什么小火龙？——我要把它，练成谁也吹不熄的火。"。动作帧（动图）：①风雨浇尾 ②蜷身护焰 ③火苗忽暗 ④雨停望火立志。诗词：风暴之夜雨浇尾，火苗忽暗心头危。蜷身护焰等雨歇，立志燃火永难吹。。主题句：尾尖之火不熄，是生命也是守护——从火苗到圣焰，燃尽一切黑暗。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：橙红+雨夜蓝+微弱焰火。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 凡尘砺心 火熄之惧, age 幼年, 风暴之夜·火差点熄灭; scene: 风雨里，它蜷起身子用爪子护着尾焰，火苗在掌间忽明忽暗；等雨停，它望着重新燃起的火，眼眶发热。; 风暴之夜雨浇尾，火苗忽暗心头危。蜷身护焰等雨歇，立志燃火永难吹。; palette: 橙红+雨夜蓝+微弱焰火; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 烈火小龙**
- 小火龙，道法初成阶段·燃火苦修（少年，苦练控火·尾焰越旺）。形象：小火龙，橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。品性：把火熄之惧化作苦修之力，日夜练控火。从"火花"练到"烈焰"，尾尖火焰越烧越旺，风雨里也能稳如磐石。。姿态：对着山壁一遍遍喷火，从火花到火柱；风雨天也坚持练，让尾焰在雨里也能烧得旺旺的。。服饰：橙红鳞片更亮，尾焰橘红转赤。。体型：身高约2头身，橙色小蜥蜴，鳞片发亮。。衣物细节：橙红鳞片更亮，尾焰橘红转赤。。发型妆造：无发，背脊凸起渐显。。脸型五官：圆眼专注，腮帮鼓满火气。。武器招式：烈焰初成。。功法：控火渐精；烈焰初成。。功法表现：尾焰赤红，风雨不熄。。画面：构图：山壁前，橙色小蜥蜴火柱喷向山壁，风雨中尾焰仍旺，背景山岩。色调：橙红+焰火橘+山灰。氛围：苦修、控火、渐强。。台词："火，最怕灭，也最能炼。我且日日烧，烧到风雨也浇不熄，烧到能护住想护的人。"。动作帧（动图）：①对山壁喷火 ②火花到火柱 ③风雨天也练 ④尾焰越烧越旺。诗词：燃火苦修对山壁，火柱渐成风雨立。尾焰橘红转炽赤，只待他日焚夜黑。。主题句：尾尖之火不熄，是生命也是守护——从火苗到圣焰，燃尽一切黑暗。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：橙红+焰火橘+山灰。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 道法初成 燃火苦修, age 少年, 苦练控火·尾焰越旺; scene: 对着山壁一遍遍喷火，从火花到火柱；风雨天也坚持练，让尾焰在雨里也能烧得旺旺的。; 燃火苦修对山壁，火柱渐成风雨立。尾焰橘红转炽赤，只待他日焚夜黑。; palette: 橙红+焰火橘+山灰; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 烈火恐龙**
- 小火龙，大劫淬炼阶段·火焰护伴（少年，第一次火焰护伴·照夜路）。形象：小火龙，橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。品性：第一次以火焰护住同伴——夜路黑暗，它鼓起腮帮喷出火焰，照亮前路，也赶走了暗处的凶险。。姿态：同伴在夜里迷路遇险，它一记火焰轰开暗影，火光照亮夜路；看着同伴安然跟上来，它咧嘴笑，尾焰烧得更旺。。服饰：橙红鳞片，尾焰赤红，火光映身。。体型：身高约2头身，橙色小蜥蜴，尾焰赤红。。衣物细节：橙红鳞片，尾焰赤红。。发型妆造：无发，火光映身。。脸型五官：圆眼含笑，火光映腮。。武器招式：火焰初成。。功法：火焰护伴；照夜路。。功法表现：火焰照夜，驱散暗影。。画面：构图：夜路，橙色小蜥蜴火焰照亮前路，护住同伴，暗影被火光驱散，背景夜林。色调：橙红+焰火橘+夜蓝。氛围：小成、护伴、照夜。。台词："我的火，从前只用来烤石头。从今天起，它用来护人、照夜——火亮之处，谁也别想伤害我的同伴。"。动作帧（动图）：①同伴夜路遇险 ②一记火焰轰暗影 ③火光照路 ④咧嘴笑尾焰旺。诗词：火焰一记护同裳，照夜驱暗火光扬。小火龙儿初亮火，尾焰自此为伴燃。。主题句：尾尖之火不熄，是生命也是守护——从火苗到圣焰，燃尽一切黑暗。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：橙红+焰火橘+夜蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 大劫淬炼 火焰护伴, age 少年, 第一次火焰护伴·照夜路; scene: 同伴在夜里迷路遇险，它一记火焰轰开暗影，火光照亮夜路；看着同伴安然跟上来，它咧嘴笑，尾焰烧得更旺。; 火焰一记护同裳，照夜驱暗火光扬。小火龙儿初亮火，尾焰自此为伴燃。; palette: 橙红+焰火橘+夜蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 喷火龙**
- 小火龙，封神登天阶段·喷火龙（青年，完全进化·喷火龙形态）。形象：小火龙，橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。品性：完全进化成喷火龙——双翼展开遮天，尾焰如炬，口中喷出圣焰。它终于成了曾想成为的样子：能展翅守护的火之传说。。姿态：进化之光中双翼展开，一记圣焰冲天；它翱翔于天，俯瞰下方的同伴——"我的火，终于大到能护住一片天。"。服饰：喷火龙形态，橙红鳞甲，双翼遮天，尾焰如炬。。体型：身高约4头身，喷火龙形态，双翼遮天。。衣物细节：喷火龙形态，橙红鳞甲，双翼遮天，尾焰如炬。。发型妆造：无发，背鳍展开。。脸型五官：圆眼含威，口中圣焰。。武器招式：圣焰初成，烈焰冲霄。。功法：喷火龙形态；圣焰初成。。功法表现：圣焰冲天。。画面：构图：晴空，喷火龙形态双翼展开，圣焰冲天，下方同伴仰望，背景天穹。色调：橙红+圣焰金+天蓝。氛围：大成、展翅、圣焰。。台词："小时候我只会喷火花，如今我圣焰冲天。可我最珍惜的，不是这双翼与圣焰，是终于能展翅护住想护的人。"。动作帧（动图）：①进化之光 ②双翼展开 ③圣焰冲天 ④翱翔俯瞰同伴。诗词：完全进化喷火龙，双翼遮天圣焰雄。从兹不是喷花幼，护人之焰耀苍穹。。主题句：尾尖之火不熄，是生命也是守护——从火苗到圣焰，燃尽一切黑暗。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：橙红+圣焰金+天蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 封神登天 喷火龙, age 青年, 完全进化·喷火龙形态; scene: 进化之光中双翼展开，一记圣焰冲天；它翱翔于天，俯瞰下方的同伴——"我的火，终于大到能护住一片天。"; 完全进化喷火龙，双翼遮天圣焰雄。从兹不是喷花幼，护人之焰耀苍穹。; palette: 橙红+圣焰金+天蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 超级喷火龙X**
- 小火龙，道果圆满阶段·圣焰传说（终极，燃尽黑暗·火之传说）。形象：小火龙，橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。品性：圣焰圆满，成了燃尽黑暗的火之传说。它守在高天，为迷途者点亮归路——从尾尖一簇火苗，到照亮黑夜的圣焰，它护了一路。。姿态：黑夜，它盘旋于森林上空，圣焰为迷途的小家伙们引路；看着新生的小火龙对着它喷出第一簇火苗，它俯身，用翼尖轻轻拨了拨那簇小火。。服饰：喷火龙形态，圣焰如翼，周身火光照夜。。体型：身高约4头身，喷火龙形态，圣焰如翼。。衣物细节：喷火龙形态，圣焰如翼，火光绕身。。发型妆造：无发，背鳍如焰。。脸型五官：圆眼含笑，圣焰在口。。武器招式：圣焰圆满，燃尽黑暗。。功法：火之传说；燃尽黑暗。。功法表现：圣焰照夜，火之传说。。画面：构图：黑夜森林上空，喷火龙形态圣焰如翼盘旋引路，新生小火龙喷出第一簇火苗，背景星空与火光。色调：橙红+圣焰金+夜蓝。氛围：终极、燃夜、传说。。台词："从前我怕火会灭，如今我燃烧成传说。可我最珍惜的，不是这身圣焰，是每一个被我照亮归路的小家伙。"。动作帧（动图）：①盘旋森林上空 ②圣焰引路 ③见新生火苗 ④翼尖轻拨。诗词：圣焰传说照夜行，盘旋高天引归程。新生火苗翼尖拨，火名传世永长明。。主题句：尾尖之火不熄，是生命也是守护——从火苗到圣焰，燃尽一切黑暗。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：橙红+圣焰金+夜蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 道果圆满 圣焰传说, age 终极, 燃尽黑暗·火之传说; scene: 黑夜，它盘旋于森林上空，圣焰为迷途的小家伙们引路；看着新生的小火龙对着它喷出第一簇火苗，它俯身，用翼尖轻轻拨了拨那簇小火。; 圣焰传说照夜行，盘旋高天引归程。新生火苗翼尖拨，火名传世永长明。; palette: 橙红+圣焰金+夜蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 妙蛙种子（`bulbasaur`） · 人生档案版

**灵胎初醒 · 种子蛋**
- 妙蛙种子，灵胎初醒阶段·种子之芽（出生，背驮种子苞·爱晒太阳）。形象：妙蛙种子，蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。品性：草系宝可梦，蓝绿蛙状小兽，背驮种子苞。安静沉着，爱晒太阳，靠阳光积蓄生长的力量。。姿态：伏在暖阳下晒太阳，种子苞微微发亮；看见蔫了的花，会用藤鞭轻轻扶一扶。。服饰：无衣，蓝绿表皮，背驮种子苞，脚掌短。。体型：身高约2头身，蓝绿蛙状，背驮种子苞，身形圆润。。衣物细节：蓝绿表皮，背部深绿种子苞，脚掌带趾。。发型妆造：头顶小绿冠，种子苞含蕾。。脸型五官：圆眼清澈，颊有绿斑，神情安静。。武器招式：藤鞭初窥，叶芽。。功法：光合生长；藤鞭初窥。。功法表现：光合微光，种子苞发亮。。画面：构图：暖阳草地，蓝绿蛙状小兽伏地晒太阳，种子苞发亮，旁有蔫花被藤鞭轻扶，背景草原。色调：蓝绿+暖阳金+草青。氛围：安静、向阳、初醒。。台词："太阳一晒，我就有劲儿了。这背上的种子苞，会开出很漂亮的花吧？"。动作帧（动图）：①伏地晒太阳 ②种子苞发亮 ③见蔫花轻扶 ④安静入定。诗词：蓝绿蛙兽伏暖阳，种子苞儿微微光。静爱日光积蓄力，背苞含蕾待芬芳。。主题句：一粒种子向阳而生，背苞花开护绿意——藤鞭所至，皆是生机。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝绿+暖阳金+草青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 灵胎初醒 种子之芽, age 出生, 背驮种子苞·爱晒太阳; scene: 伏在暖阳下晒太阳，种子苞微微发亮；看见蔫了的花，会用藤鞭轻轻扶一扶。; 蓝绿蛙兽伏暖阳，种子苞儿微微光。静爱日光积蓄力，背苞含蕾待芬芳。; palette: 蓝绿+暖阳金+草青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼种子**
- 妙蛙种子，凡尘砺心阶段·烈日之灾（幼年，大旱·护不住花草）。形象：妙蛙种子，蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。品性：一场大旱，草木枯黄，河流干涸。它拼命想用藤鞭护住身边的花草，却看着它们一点点蔫死——"我的种子苞，还不会开花。"。姿态：烈日下，它用藤鞭给枯草挡阳、找水，却杯水车薪；看着一株株枯死，它垂下头，种子苞也蔫蔫的。。服饰：蓝绿表皮被晒得发干，种子苞无精打采。。体型：身高约2头身，蓝绿蛙状，有些蔫。。衣物细节：蓝绿表皮晒干，种子苞蔫。。发型妆造：绿冠垂着，种子苞无光。。脸型五官：圆眼含泪，神情执着。。武器招式：藤鞭无力。。功法：光合不足；无力。。功法表现：光合不足，日光灼灼。。画面：构图：龟裂大地，蓝绿蛙兽用藤鞭护着一株蔫草，种子苞垂着，背景烈日与枯黄草木。色调：蓝绿+枯黄+烈日白。氛围：无力、护绿、立志。。台词："我背上的种子苞，还不会开花。我护不住这些花草……可我不想放弃，我要让它们，都能活。"。动作帧（动图）：①烈日下护草 ②找水无果 ③看草木枯死 ④垂头立志。诗词：烈日大旱草木枯，藤鞭护草水难赊。种子苞儿蔫垂首，立志护绿盼雨苏。。主题句：一粒种子向阳而生，背苞花开护绿意——藤鞭所至，皆是生机。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝绿+枯黄+烈日白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 凡尘砺心 烈日之灾, age 幼年, 大旱·护不住花草; scene: 烈日下，它用藤鞭给枯草挡阳、找水，却杯水车薪；看着一株株枯死，它垂下头，种子苞也蔫蔫的。; 烈日大旱草木枯，藤鞭护草水难赊。种子苞儿蔫垂首，立志护绿盼雨苏。; palette: 蓝绿+枯黄+烈日白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 叶芽成长**
- 妙蛙种子，道法初成阶段·扎根苦修（少年，扎根大地·苦练藤鞭）。形象：妙蛙种子，蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。品性：把护不住花草的痛化作苦修之力，扎根大地，积蓄阳光，日夜练藤鞭——从细软练到坚韧，从探不到练到缠得住。。姿态：清晨将根扎进湿润的土里，晒够太阳，再对着枯木练藤鞭；藤鞭从细细软软，练到能卷起石块，种子苞也一天天鼓起来。。服饰：蓝绿表皮更绿，种子苞渐鼓，藤鞭粗韧。。体型：身高约2头身，蓝绿蛙状，扎根稳立。。衣物细节：蓝绿表皮更绿，种子苞鼓，藤鞭粗韧。。发型妆造：绿冠挺立，种子苞渐饱。。脸型五官：圆眼专注，神情坚毅。。武器招式：藤鞭渐韧。。功法：藤鞭渐韧；扎根蓄力。。功法表现：光合充盈，种子苞微光。。画面：构图：清晨草地，蓝绿蛙兽扎根土中晒阳，藤鞭卷起石块，种子苞渐鼓，背景晨光。色调：蓝绿+晨光金+草青。氛围：苦修、扎根、渐强。。台词："根扎得深，才吸得到水；藤鞭练得韧，才护得住花草。我且慢慢来，跟种子苞一起长大。"。动作帧（动图）：①根扎进土 ②晒足太阳 ③练藤鞭卷石块 ④种子苞渐鼓。诗词：扎根大地练藤鞭，朝沐暖阳夜蓄泉。细软渐成坚韧劲，种子苞儿渐次圆。。主题句：一粒种子向阳而生，背苞花开护绿意——藤鞭所至，皆是生机。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝绿+晨光金+草青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 道法初成 扎根苦修, age 少年, 扎根大地·苦练藤鞭; scene: 清晨将根扎进湿润的土里，晒够太阳，再对着枯木练藤鞭；藤鞭从细细软软，练到能卷起石块，种子苞也一天天鼓起来。; 扎根大地练藤鞭，朝沐暖阳夜蓄泉。细软渐成坚韧劲，种子苞儿渐次圆。; palette: 蓝绿+晨光金+草青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 花苞绽放**
- 妙蛙种子，大劫淬炼阶段·藤鞭护生（少年，第一次护住花草·生命复苏）。形象：妙蛙种子，蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。品性：第一次以藤鞭护住一片枯蔫的花草——它用积攒的阳光与水分，让它们重新活过来。看着绿意回来，它第一次懂了"种子苞的意义"。。姿态：一片干枯的花丛前，它藤鞭探入，将积蓄的养分与水渡给它们；看着花叶渐渐挺直、绽开，它笑弯了眼——"我做到了。"。服饰：蓝绿表皮，种子苞微绽，藤鞭舒展。。体型：身高约2头身，蓝绿蛙状，种子苞微绽。。衣物细节：蓝绿表皮，种子苞微绽，藤鞭舒展。。发型妆造：绿冠挺立，种子苞绽蕾。。脸型五官：圆眼含笑，眉眼弯弯。。武器招式：藤鞭护生。。功法：藤鞭护生；生命复苏。。功法表现：生命复苏，绿意回转。。画面：构图：干枯花丛前，蓝绿蛙兽藤鞭渡养分，花叶重绽，背景阳光与新绿。色调：蓝绿+新绿+暖阳金。氛围：小成、护绿、生命复苏。。台词："从前我护不住花草，如今我让它们活过来了。这背上的种子苞，终于有了意义——它开花，是为了护住更多的绿。"。动作帧（动图）：①到干枯花丛前 ②藤鞭渡养分 ③花叶重绽 ④笑弯了眼。诗词：藤鞭护草度生机，花叶挺直绽笑眉。种子苞儿初含意，护绿之心自此持。。主题句：一粒种子向阳而生，背苞花开护绿意——藤鞭所至，皆是生机。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝绿+新绿+暖阳金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 大劫淬炼 藤鞭护生, age 少年, 第一次护住花草·生命复苏; scene: 一片干枯的花丛前，它藤鞭探入，将积蓄的养分与水渡给它们；看着花叶渐渐挺直、绽开，它笑弯了眼——"我做到了。"; 藤鞭护草度生机，花叶挺直绽笑眉。种子苞儿初含意，护绿之心自此持。; palette: 蓝绿+新绿+暖阳金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 繁花之姿**
- 妙蛙种子，封神登天阶段·妙蛙花绽（青年，背苞开花·妙蛙花形态）。形象：妙蛙种子，蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。品性：背苞绽放成巨大的花，进化成妙蛙花形态——藤鞭如臂，花开护荫。它终于成了曾经想成为的样子：一棵能庇护大片绿意的树。。姿态：背苞在晨光中轰然绽放，花香弥漫；它立在林间，藤鞭舒展如臂，为底下的小草遮荫——"花开，是为了护住这一片。"。服饰：妙蛙花形态，背驮巨大的花，藤鞭粗壮，蓝绿表皮更健硕。。体型：身高约3头身，妙蛙花形态，背驮巨花。。衣物细节：妙蛙花形态，背花盛开，藤鞭粗壮。。发型妆造：背花如冠，绿叶舒展。。脸型五官：圆眼温和，背花盛开。。武器招式：藤鞭如臂，飞叶快刀。。功法：妙蛙花形态；花开护荫。。功法表现：花开护荫，生机四溢。。画面：构图：晨光林间，妙蛙花形态背花绽放，藤鞭舒展为小草遮荫，背景晨光与新绿。色调：妙蛙花红白+蓝绿+晨光金。氛围：大成、花开、护荫。。台词："我的种子苞，终于开花了。这朵花，是太阳晒出来的，也是那些枯死花草教我护出来的——花开之地，皆是生机。"。动作帧（动图）：①背苞晨光中绽放 ②花香弥漫 ③藤鞭为小草遮荫 ④立在林间。诗词：背苞绽放妙蛙花，藤鞭如臂护新芽。花开一树荫草处，护绿之心渐已华。。主题句：一粒种子向阳而生，背苞花开护绿意——藤鞭所至，皆是生机。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：妙蛙花红白+蓝绿+晨光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 封神登天 妙蛙花绽, age 青年, 背苞开花·妙蛙花形态; scene: 背苞在晨光中轰然绽放，花香弥漫；它立在林间，藤鞭舒展如臂，为底下的小草遮荫——"花开，是为了护住这一片。"; 背苞绽放妙蛙花，藤鞭如臂护新芽。花开一树荫草处，护绿之心渐已华。; palette: 妙蛙花红白+蓝绿+晨光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 花神**
- 妙蛙种子，道果圆满阶段·花开传说（终极，守护绿意·草系传说）。形象：妙蛙种子，蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。品性：花开圆满，成了守护整片森林的草系传说。它守着四季，看着枯草重生、花开又谢——它的花，护住了一代又一代的绿意。。姿态：春日，它在森林中央盛放，藤鞭为新生的小草引路；看着无数妙蛙种子伏在它花下，它垂下花冠，轻轻一颤——像在笑。。服饰：妙蛙花形态，背花如盖，藤鞭如林，周身绿意流转。。体型：身高约3头身，妙蛙花形态，背花如盖。。衣物细节：妙蛙花形态，背花如盖，藤鞭如林。。发型妆造：背花常开，绿叶常青。。脸型五官：圆眼温和，花冠轻垂含笑。。武器招式：藤鞭如林，花开护荫。。功法：草系传说；守护绿意。。功法表现：绿意流转，守护森林。。画面：构图：春日森林中央，妙蛙花形态盛放，无数妙蛙种子伏在花下，藤鞭为新生小草引路，背景春日林光。色调：妙蛙花红白+森林绿+春日金。氛围：终极、守护、传说。。台词："我的种子苞，从一粒到花开满树。可我最珍惜的，不是这身绿意，是每一个伏在我花下安然长大的小家伙。"。动作帧（动图）：①森林中央盛放 ②藤鞭为小草引路 ③妙蛙种子花下聚 ④花冠轻颤含笑。诗词：花开传说护森林，春来秋去绿常新。妙蛙种子花下聚，草系之名永长青。。主题句：一粒种子向阳而生，背苞花开护绿意——藤鞭所至，皆是生机。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：妙蛙花红白+森林绿+春日金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 道果圆满 花开传说, age 终极, 守护绿意·草系传说; scene: 春日，它在森林中央盛放，藤鞭为新生的小草引路；看着无数妙蛙种子伏在它花下，它垂下花冠，轻轻一颤——像在笑。; 花开传说护森林，春来秋去绿常新。妙蛙种子花下聚，草系之名永长青。; palette: 妙蛙花红白+森林绿+春日金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 杰尼龟（`squirtle`） · 人生档案版

**灵胎初醒 · 水纹蛋**
- 杰尼龟，灵胎初醒阶段·水泡幼龟（出生，蓝色小龟·机灵顽皮）。形象：杰尼龟，蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。品性：水系宝可梦，蓝色小龟，背驮硬壳。机灵顽皮，爱搞恶作剧，龟壳既是铠甲也是家。。姿态：在水边吐泡泡捉弄小鱼；被追时"唰"地缩进壳里，眨着眼看对方无从下手。。服饰：无衣，蓝色皮肤，背驮棕褐圆壳，腹部米黄。。体型：身高约2头身，蓝色小龟，圆壳，身形机灵。。衣物细节：通体蓝色皮肤，背驮棕褐圆壳，腹部米黄，尾卷。。发型妆造：头顶无发，壳缘花纹。。脸型五官：圆眼机灵，嘴角俏皮，腮帮圆。。武器招式：水枪初窥，吐水泡。。功法：水枪初窥；缩壳护身。。功法表现：口中水花，泡泡浮起。。画面：构图：水边，蓝色小龟吐泡泡，背驮圆壳，遇人缩进壳里眨眼，背景水泽。色调：蓝色+棕褐壳+水青。氛围：机灵、顽皮、壳护身。。台词："我这壳，是铠甲，也是家。谁也砸不破，谁也追不上——躲在里面，天塌下来也不怕。"。动作帧（动图）：①吐泡泡 ②捉弄小鱼 ③遇人缩壳 ④壳里眨眼。诗词：蓝色小龟吐水泡，机灵顽皮追鱼闹。遇事一缩壳里躲，壳下自有天地好。。主题句：缩进壳里不是怕，是等一记更大的水炮——龟壳之下，是护人的心。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝色+棕褐壳+水青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 灵胎初醒 水泡幼龟, age 出生, 蓝色小龟·机灵顽皮; scene: 在水边吐泡泡捉弄小鱼；被追时"唰"地缩进壳里，眨着眼看对方无从下手。; 蓝色小龟吐水泡，机灵顽皮追鱼闹。遇事一缩壳里躲，壳下自有天地好。; palette: 蓝色+棕褐壳+水青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼龟**
- 杰尼龟，凡尘砺心阶段·干涸之痛（幼年，水泽干涸·水还不够多）。形象：杰尼龟，蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。品性：一场大旱，水泽干涸，同伴渴得直喘。它拼命喷水，却只够解一瞬的渴——"我的水，还不够多。"。姿态：它用壳背水、用嘴喷水，却杯水车薪；看着同伴渴得趴下，它第一次恨自己"只会吐小水花"。。服饰：蓝色皮肤发干，龟壳晒得发烫。。体型：身高约2头身，蓝色小龟，有些干。。衣物细节：蓝色皮肤发干，龟壳发烫。。发型妆造：无发，壳缘花纹晒淡。。脸型五官：圆眼含忧，神情执着。。武器招式：水枪不足。。功法：水枪不足；干涸之痛。。功法表现：水花渐小。。画面：构图：干涸水泽，蓝色小龟用嘴喷水，同伴渴卧，龟壳被晒得发烫，背景龟裂大地。色调：蓝色+干裂棕+烈日白。氛围：无力、干渴、立志。。台词："我的水，只够解一瞬的渴。可我不想再看着大家渴下去了——我要练成，能救回一片水泽的水炮。"。动作帧（动图）：①水泽干涸 ②背水喷水 ③同伴渴卧 ④立志练水炮。诗词：大旱水泽渐干枯，同伴渴卧难支吾。背水喷嘴犹不足，立志练成水炮湖。。主题句：缩进壳里不是怕，是等一记更大的水炮——龟壳之下，是护人的心。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝色+干裂棕+烈日白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 凡尘砺心 干涸之痛, age 幼年, 水泽干涸·水还不够多; scene: 它用壳背水、用嘴喷水，却杯水车薪；看着同伴渴得趴下，它第一次恨自己"只会吐小水花"。; 大旱水泽渐干枯，同伴渴卧难支吾。背水喷嘴犹不足，立志练成水炮湖。; palette: 蓝色+干裂棕+烈日白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 水枪龟**
- 杰尼龟，道法初成阶段·蓄水苦修（少年，苦练水炮·龟壳磨亮）。形象：杰尼龟，蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。品性：把干涸之痛化作苦修之力，日夜练水枪。从水花练到水柱，从水柱练到水炮——龟壳也磨得发亮，存的水越来越多。。姿态：对着崖壁一遍遍喷水，从水花到水炮；清晨吸足水，把壳里存得满满的，练到月上枝头。。服饰：蓝色皮肤水润，龟壳磨得发亮。。体型：身高约2头身，蓝色小龟，龟壳发亮。。衣物细节：蓝色皮肤水润，龟壳磨亮。。发型妆造：无发，壳缘花纹清晰。。脸型五官：圆眼专注，腮帮鼓起蓄水。。武器招式：水炮渐成。。功法：水枪渐成水炮；蓄水之能。。功法表现：口中水炮，蓄势待发。。画面：构图：崖壁前，蓝色小龟水炮喷向崖壁，龟壳磨亮，背景水泽远影。色调：蓝色+水花白+崖灰。氛围：苦修、蓄水、渐强。。台词："水，不是天生的，是一口口存、一炮炮练的。我且把壳里存满，练到能救回一片水泽。"。动作帧（动图）：①对崖壁喷水 ②水花到水炮 ③壳里存水 ④练到月上枝头。诗词：蓄水苦修对崖鸣，水花渐成水炮横。龟壳磨亮存满水，只待他日泽又盈。。主题句：缩进壳里不是怕，是等一记更大的水炮——龟壳之下，是护人的心。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝色+水花白+崖灰。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 道法初成 蓄水苦修, age 少年, 苦练水炮·龟壳磨亮; scene: 对着崖壁一遍遍喷水，从水花到水炮；清晨吸足水，把壳里存得满满的，练到月上枝头。; 蓄水苦修对崖鸣，水花渐成水炮横。龟壳磨亮存满水，只待他日泽又盈。; palette: 蓝色+水花白+崖灰; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 水炮龟**
- 杰尼龟，大劫淬炼阶段·水炮护泽（少年，第一次水炮·救回水泽）。形象：杰尼龟，蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。品性：第一次以水炮救回干涸的水泽——水柱冲入河床，泉水重新涌出，同伴欢呼着扑进水里。它看着恢复生机的泽水，笑出了声。。姿态：一记水炮注入干涸河床，泉水应声涌出；同伴扑进水里的欢呼里，它缩进壳里偷着乐——"我的壳，终于存得住一片水泽了。"。服饰：蓝色皮肤水润，龟壳闪着水光。。体型：身高约2头身，蓝色小龟，龟壳水光。。衣物细节：蓝色皮肤水润，龟壳闪水光。。发型妆造：无发，壳缘水珠。。脸型五官：圆眼含笑，躲在壳里偷乐。。武器招式：水炮初成。。功法：水炮初成；救回水泽。。功法表现：水炮注河，泉涌泽回。。画面：构图：干涸河床，蓝色小龟水炮注入，泉水涌出，同伴扑水欢呼，背景水泽复苏。色调：蓝色+水花白+泽青。氛围：小成、护泽、救回。。台词："这一炮，我存了无数个日夜。从前我只够解一瞬的渴，如今我能救回一片水泽——我杰尼龟，终于能护住大家了。"。动作帧（动图）：①水炮注入河床 ②泉水涌出 ③同伴扑水 ④缩壳偷乐。诗词：水炮一记注干河，泉涌泽回众欢歌。壳存水泽终有报，蓝龟咧嘴笑呵呵。。主题句：缩进壳里不是怕，是等一记更大的水炮——龟壳之下，是护人的心。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝色+水花白+泽青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 大劫淬炼 水炮护泽, age 少年, 第一次水炮·救回水泽; scene: 一记水炮注入干涸河床，泉水应声涌出；同伴扑进水里的欢呼里，它缩进壳里偷着乐——"我的壳，终于存得住一片水泽了。"; 水炮一记注干河，泉涌泽回众欢歌。壳存水泽终有报，蓝龟咧嘴笑呵呵。; palette: 蓝色+水花白+泽青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 水箭龟**
- 杰尼龟，封神登天阶段·水箭龟（青年，双炮进化·水箭龟形态）。形象：杰尼龟，蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。品性：进化成水箭龟形态——背上伸出两门水炮，水柱如瀑。它终于成了曾想成为的样子：能守护水泽与同伴的水箭龟。。姿态：进化之光中，双炮伸展；一记双炮齐发，水柱冲霄，它望着恢复生机的水泽，眼里有光。。服饰：水箭龟形态，背驮双管水炮，蓝色皮肤更健硕。。体型：身高约3头身，水箭龟形态，背驮双炮。。衣物细节：水箭龟形态，背驮双管水炮。。发型妆造：无发，壳缘水纹。。脸型五官：圆眼沉稳，双炮在背。。武器招式：双炮齐发。。功法：水箭龟形态；双炮齐发。。功法表现：水柱冲霄。。画面：构图：水泽边，水箭龟形态背驮双炮，水柱冲霄，泽水波光，背景水泽山林。色调：蓝色+水花白+泽青。氛围：大成、双炮、护泽。。台词："小时候我只会吐水泡，如今我双炮齐发。可我最珍惜的，不是这身水炮，是终于能守护住这一片水泽和大家。"。动作帧（动图）：①进化之光 ②双炮伸展 ③水柱冲霄 ④望泽眼有光。诗词：双炮进化水箭龟，水柱冲霄泽更滋。从兹不是吐泡幼，护泽之心已堪持。。主题句：缩进壳里不是怕，是等一记更大的水炮——龟壳之下，是护人的心。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝色+水花白+泽青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 封神登天 水箭龟, age 青年, 双炮进化·水箭龟形态; scene: 进化之光中，双炮伸展；一记双炮齐发，水柱冲霄，它望着恢复生机的水泽，眼里有光。; 双炮进化水箭龟，水柱冲霄泽更滋。从兹不是吐泡幼，护泽之心已堪持。; palette: 蓝色+水花白+泽青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 海神龟**
- 杰尼龟，道果圆满阶段·泽水传说（终极，守护水泽·水之传说）。形象：杰尼龟，蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。品性：水炮圆满，成了守护水泽的传说。它守着四季的水，看枯泽回春、水草丰茂——它的壳里，永远存着一片能救急的水。。姿态：旱季，它守在泽边，一记水炮为枯泽续命；看着新生的水箭龟伏在它身边，它沉进水里，只留龟壳在波光中——像在点头。。服饰：水箭龟形态，双炮如瀑，周身水光流转。。体型：身高约3头身，水箭龟形态，双炮如瀑。。衣物细节：水箭龟形态，双炮如瀑，水光流转。。发型妆造：无发，壳缘水纹常新。。脸型五官：圆眼沉稳含笑，沉水留壳。。武器招式：水炮圆满，双炮如瀑。。功法：水之传说；守护水泽。。功法表现：守护水泽，水光长流。。画面：构图：泽畔，水箭龟形态双炮如瀑续泽，新生水箭龟伏旁，背景波光与水草。色调：蓝色+水光白+泽青。氛围：终极、守护、传说。。台词："从前我只够解一瞬的渴，如今我守护整片水泽。可我最珍惜的，是每一个伏在泽边安然喝水的小家伙。"。动作帧（动图）：①旱季守泽 ②一记水炮续命 ③新生水箭伏旁 ④沉水留壳点头。诗词：泽水传说守水乡，旱季一炮续枯荒。新生水箭伏泽畔，水之名传永流长。。主题句：缩进壳里不是怕，是等一记更大的水炮——龟壳之下，是护人的心。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：蓝色+水光白+泽青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 道果圆满 泽水传说, age 终极, 守护水泽·水之传说; scene: 旱季，它守在泽边，一记水炮为枯泽续命；看着新生的水箭龟伏在它身边，它沉进水里，只留龟壳在波光中——像在点头。; 泽水传说守水乡，旱季一炮续枯荒。新生水箭伏泽畔，水之名传永流长。; palette: 蓝色+水光白+泽青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 伊布（`eevee`）

**灵胎初醒 · 伊布蛋**
- 伊布，灵胎初醒阶段·伊布蛋。初始形态：一枚伊布蛋，奶棕蛋壳泛着多变的光晕，进化之力在壳中随时光流转。光属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 伊布**
- 伊布，凡尘砺心阶段·伊布。形象：棕毛狐狸样小兽，颈毛蓬松。 核心意象：蓬松颈毛、进化之力、多形态可能。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 进化预备**
- 伊布，道法初成阶段·进化预备。形象：棕毛狐狸样小兽，颈毛蓬松。 核心意象：蓬松颈毛、进化之力、多形态可能。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，蹭蹭主人的手，一记"撞击"轻快而出。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 元素显现**
- 伊布，大劫淬炼阶段·元素显现。形象：棕毛狐狸样小兽，颈毛蓬松。 核心意象：蓬松颈毛、进化之力、多形态可能。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 水伊布·潮汐**
- 伊布，封神登天阶段·水伊布·潮汐。形象：棕毛狐狸样小兽，颈毛蓬松。 核心意象：蓬松颈毛、进化之力、多形态可能。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 受封万兽之王，目光睥睨天地; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 水伊布·海神**
- 伊布，道果圆满阶段·水伊布·海神。形象：棕毛狐狸样小兽，颈毛蓬松。 核心意象：蓬松颈毛、进化之力、多形态可能。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，蹭蹭主人的手，一记"撞击"轻快而出。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 神光自照，与天地同尊; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 皮卡丘（`pikachu`） · 人生档案版

**灵胎初醒 · 电气蛋**
- 皮卡丘，灵胎初醒阶段·皮丘幼崽（出生，静电体质·胆小害羞）。形象：皮卡丘，黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。品性：电系宝可梦的幼崽，通体奶黄，双颊红晕。胆小害羞，怕打雷，一紧张就"滋滋"漏电。。姿态：躲在树后看别的宝可梦放电，满眼羡慕；自己一紧张，脸颊电光"滋滋"，吓自己一跳。。服饰：无衣，奶黄短毛，双颊红晕，尾尖闪电状。。体型：身高约2头身，圆滚滚奶黄小球，身形灵动。。衣物细节：通体奶黄短毛，背部两道深棕条纹，尾端闪电状，双颊电囊。。发型妆造：双耳尖长，尾如闪电，腮部电囊。。脸型五官：黑豆圆眼，双颊红晕两点，小鼻尖，咧嘴笑。。武器招式：静电微电，尾尖火花。。功法：静电微电（时灵时不灵）。。功法表现：静电环绕，双颊滋滋。。画面：构图：林间，奶黄小皮丘躲在树后，双颊"滋滋"漏电，远处雷云，背景草丛。色调：奶黄+红晕+林青。氛围：胆小、羡慕、初醒。。台词："为什么别人放电那么帅，我一紧张就只会滋滋响？我也想，有一天电光四射。"。动作帧（动图）：①躲树后 ②看别人放电满眼羡慕 ③一紧张滋滋漏电 ④吓自己一跳。诗词：奶黄皮丘躲树丛，红晕双颊怕雷公。一紧张来滋滋响，何时电光耀夜空。。主题句：从弱小的皮丘，到电光四射的雷丘——一记十万伏特，照亮所有黑夜。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶黄+红晕+林青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 灵胎初醒 皮丘幼崽, age 出生, 静电体质·胆小害羞; scene: 躲在树后看别的宝可梦放电，满眼羡慕；自己一紧张，脸颊电光"滋滋"，吓自己一跳。; 奶黄皮丘躲树丛，红晕双颊怕雷公。一紧张来滋滋响，何时电光耀夜空。; palette: 奶黄+红晕+林青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 皮丘萌动**
- 皮卡丘，凡尘砺心阶段·雷雨之惧（幼年，雷雨之夜·电还不够亮）。形象：皮卡丘，黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。品性：一次雷雨，它被同伴笑"电连萤火虫都不如"，又被雷声吓到瘫软。弱小之痛，第一次咬进心里。。姿态：雷雨中缩在角落，双颊电光只"滋滋"几下就暗了；同伴嘲笑，它攥紧小拳——"我要变得很强。"。服饰：奶黄短毛被雨打湿，双颊红晕黯淡。。体型：身高约2头身，缩成一团奶黄小球。。衣物细节：奶黄短毛淋湿，尾端闪电垂着。。发型妆造：双耳耷拉，尾如闪电却无光。。脸型五官：黑豆圆眼含泪，双颊红晕黯淡。。武器招式：静电微电不足。。功法：静电微电不足；受挫。。功法表现：电光黯淡，雨夜微光。。画面：构图：雷雨夜，奶黄皮丘缩在角落，双颊电光黯淡，远处雷云与同伴背影，背景雨夜。色调：奶黄+雨夜蓝+黯淡红晕。氛围：受挫、立志、弱小之痛。。台词："我的电，连萤火虫都不如……可我不想一直这样。总有一天，我要电到所有人都看得到。"。动作帧（动图）：①缩角落 ②双颊滋滋又暗 ③听同伴嘲笑 ④攥拳立志。诗词：雷雨夜中电光暗，同伴笑它不如萤。攥紧小拳立志起，他年电光照天明。。主题句：从弱小的皮丘，到电光四射的雷丘——一记十万伏特，照亮所有黑夜。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶黄+雨夜蓝+黯淡红晕。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 凡尘砺心 雷雨之惧, age 幼年, 雷雨之夜·电还不够亮; scene: 雷雨中缩在角落，双颊电光只"滋滋"几下就暗了；同伴嘲笑，它攥紧小拳——"我要变得很强。"; 雷雨夜中电光暗，同伴笑它不如萤。攥紧小拳立志起，他年电光照天明。; palette: 奶黄+雨夜蓝+黯淡红晕; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 蓄电皮卡丘**
- 皮卡丘，道法初成阶段·电系苦修（少年，苦练放电·脸颊噼啪）。形象：皮卡丘，黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。品性：把受挫之痛化作苦修之力，日日练放电。从"滋滋"练到"噼啪"，尾尖会冒火花，双颊电光越练越亮。。姿态：对着树桩一遍遍放电，电到双颊发麻；夜里望着雷云的方向，暗下决心——"总有一天，我也能召唤那样的闪电。"。服饰：奶黄短毛微焦，双颊电光噼啪。。体型：身高约2头身，奶黄小球，身形绷紧发力。。衣物细节：奶黄短毛微焦，背部深棕纹更明显。。发型妆造：双耳立起，尾如闪电带电光。。脸型五官：黑豆圆眼专注，双颊电光噼啪。。武器招式：电光初成，尾尖火花。。功法：电光初成；尾尖冒火花。。功法表现：双颊电光噼啪，夜中微亮。。画面：构图：夜林，奶黄皮丘对着树桩放电，双颊电光噼啪，尾尖火花，背景夜与雷云远影。色调：奶黄+电光金+夜蓝。氛围：苦修、立志、渐强。。台词："电，不是天生就强的，是练出来的。我且一天一点，练到它能照亮黑夜。"。动作帧（动图）：①对树桩放电 ②电到双颊发麻 ③尾尖火花 ④望雷云立志。诗词：苦练放电到夜阑，双颊噼啪电光寒。尾尖火花溅又起，皮丘立志照夜滩。。主题句：从弱小的皮丘，到电光四射的雷丘——一记十万伏特，照亮所有黑夜。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶黄+电光金+夜蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 道法初成 电系苦修, age 少年, 苦练放电·脸颊噼啪; scene: 对着树桩一遍遍放电，电到双颊发麻；夜里望着雷云的方向，暗下决心——"总有一天，我也能召唤那样的闪电。"; 苦练放电到夜阑，双颊噼啪电光寒。尾尖火花溅又起，皮丘立志照夜滩。; palette: 奶黄+电光金+夜蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 雷丘**
- 皮卡丘，大劫淬炼阶段·十万伏特（少年，第一次十万伏特·护住同伴）。形象：皮卡丘，黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。品性：第一次放出完整的十万伏特——那是它练习无数次的电光，那一刻照亮了黑夜，也护住了被欺的同伴。。姿态：同伴被欺负，它双颊电光迸发，一记十万伏特轰开坏蛋；看着同伴安然，它咧嘴笑了——"我的电，终于能护人了。"。服饰：奶黄短毛，双颊电光迸发，尾尖闪电铮亮。。体型：身高约2头身，奶黄小球，双颊电光迸发。。衣物细节：奶黄短毛，尾端闪电铮亮。。发型妆造：双耳立起，尾如闪电带电光。。脸型五官：黑豆圆眼亮起，双颊电光噼啪。。武器招式：十万伏特初成。。功法：十万伏特初成；护住同伴。。功法表现：电光照夜，护住同伴。。画面：构图：夜林，奶黄皮丘一记十万伏特轰开坏蛋，电光照亮黑夜，同伴安然，背景电光。色调：奶黄+电光金+夜蓝。氛围：小成、护人、电光初亮。。台词："这一记十万伏特，我练了无数个夜晚。从今天起，我的电，不是滋滋响，是能护人的光！"。动作帧（动图）：①同伴被欺 ②双颊电光迸发 ③十万伏特轰开坏蛋 ④咧嘴笑。诗词：十万伏特初轰隆，电光照夜护同裳。皮丘终成电光客，红晕双颊笑中明。。主题句：从弱小的皮丘，到电光四射的雷丘——一记十万伏特，照亮所有黑夜。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶黄+电光金+夜蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 大劫淬炼 十万伏特, age 少年, 第一次十万伏特·护住同伴; scene: 同伴被欺负，它双颊电光迸发，一记十万伏特轰开坏蛋；看着同伴安然，它咧嘴笑了——"我的电，终于能护人了。"; 十万伏特初轰隆，电光照夜护同裳。皮丘终成电光客，红晕双颊笑中明。; palette: 奶黄+电光金+夜蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 飞空皮卡丘**
- 皮卡丘，封神登天阶段·雷丘形态（青年，电光大成·雷丘姿态）。形象：皮卡丘，黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。品性：电光大成，进化成更强大的雷丘形态——身形更壮，电光更盛，尾如闪电，金色的电光缭绕周身。。姿态：雷云之下，它双颊电光迸发，金电缭绕；一记电光冲霄，震彻山林——"这就是我一直想成为的样子。"。服饰：雷丘形态，金色电光缭绕，双颊电囊鼓胀，尾如闪电。。体型：身高约3头身，雷丘形态，身形壮实。。衣物细节：雷丘形态，金色电光缭绕，背部深棕纹。。发型妆造：双耳挺立，尾如闪电金电流转。。脸型五官：黑豆圆眼含威，双颊电囊鼓胀，咧嘴自信。。武器招式：十万伏特大成，电光冲霄。。功法：雷丘形态；电光大成。。功法表现：金色电光缭绕，威震山林。。画面：构图：雷云之下，雷丘形态的金色皮卡丘电光缭绕，一记电光冲霄，背景山林与雷云。色调：金色+电光金+雷云。氛围：大成、电光冲霄、威震。。台词："小时候我只会滋滋响，如今电光冲霄。可我最珍惜的，不是这身电光，是终于能护住想护的人。"。动作帧（动图）：①雷云下蓄电 ②双颊电光迸发 ③金电缭绕 ④一记电光冲霄。诗词：雷丘形态电光豪，金电缭绕震九霄。从兹不是皮丘幼，护人之电照夜高。。主题句：从弱小的皮丘，到电光四射的雷丘——一记十万伏特，照亮所有黑夜。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：金色+电光金+雷云。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 封神登天 雷丘形态, age 青年, 电光大成·雷丘姿态; scene: 雷云之下，它双颊电光迸发，金电缭绕；一记电光冲霄，震彻山林——"这就是我一直想成为的样子。"; 雷丘形态电光豪，金电缭绕震九霄。从兹不是皮丘幼，护人之电照夜高。; palette: 金色+电光金+雷云; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 极巨化皮卡丘**
- 皮卡丘，道果圆满阶段·电光传说（终极，电光四射·守护森林）。形象：皮卡丘，黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。品性：电光圆满，成了守护森林的传说。皮神之名，随电光传遍——它依旧爱笑，依旧在雷雨夜守在森林上空，为每一个被吓到的小家伙照亮黑夜。。姿态：雷雨夜，它立于森林高处，电光四射为同伴照亮归路；看着新一代皮丘对着它满眼崇拜，它露出皮神那标志性的笑。。服饰：雷丘形态，金色电光如披风，双颊电囊如两颗小太阳。。体型：身高约3头身，雷丘形态，金电如披风。。衣物细节：雷丘形态，金色电光如披风。。发型妆造：双耳挺立，尾如闪电金电永转。。脸型五官：黑豆圆眼含笑，双颊电囊如小太阳。。武器招式：十万伏特圆满，电光四射。。功法：电光传说；守护森林。。功法表现：电光传说，照亮每一个夜路。。画面：构图：雷雨夜森林高处，雷丘形态金色电光四射，照亮下方无数归家的小宝可梦，新一代皮丘仰望，背景电光雷云。色调：金色+电光+夜蓝。氛围：终极、守护、传说。。台词："从前我怕雷雨，如今我在雷雨里守护大家。我的电，从滋滋响到电光传说——可我想做的，一直是照亮每一个夜路的人。"。动作帧（动图）：①立森林高处 ②电光四射 ③照亮归家路 ④对新一代皮丘笑。诗词：电光传说照夜明，皮神之名震林庭。护得众小归家路，十万伏特永传名。。主题句：从弱小的皮丘，到电光四射的雷丘——一记十万伏特，照亮所有黑夜。。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：金色+电光+夜蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; 道果圆满 电光传说, age 终极, 电光四射·守护森林; scene: 雷雨夜，它立于森林高处，电光四射为同伴照亮归路；看着新一代皮丘对着它满眼崇拜，它露出皮神那标志性的笑。; 电光传说照夜明，皮神之名震林庭。护得众小归家路，十万伏特永传名。; palette: 金色+电光+夜蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 利欧路（`riolu`）

**灵胎初醒 · 波导蛋**
- 利欧路，灵胎初醒阶段·波导蛋。初始形态：一枚波导蛋，蓝黑蛋壳凝着波导光环，钢铁意志的波纹在壳面一圈圈荡开。钢属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 利欧路**
- 利欧路，凡尘砺心阶段·利欧路。形象：蓝黑相间的小犬，胸前黑纹。 核心意象：胸前黑纹、波导之力、格斗之魂。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 练习生**
- 利欧路，道法初成阶段·练习生。形象：蓝黑相间的小犬，胸前黑纹。 核心意象：胸前黑纹、波导之力、格斗之魂。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，双掌凝聚波导弹，一发轰出破空而去。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 路卡利欧**
- 利欧路，大劫淬炼阶段·路卡利欧。形象：蓝黑相间的小犬，胸前黑纹。 核心意象：胸前黑纹、波导之力、格斗之魂。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 波导大师**
- 利欧路，封神登天阶段·波导大师。形象：蓝黑相间的小犬，胸前黑纹。 核心意象：胸前黑纹、波导之力、格斗之魂。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 受封万兽之王，目光睥睨天地; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 传说守护者**
- 利欧路，道果圆满阶段·传说守护者。形象：蓝黑相间的小犬，胸前黑纹。 核心意象：胸前黑纹、波导之力、格斗之魂。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，双掌凝聚波导弹，一发轰出破空而去。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 神光自照，与天地同尊; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 冰·霜翎狐（`ice_fox`）

**灵胎初醒 · 冰晶蛋**
- 冰·霜翎狐，灵胎初醒阶段·冰晶蛋。初始形态：一枚冰晶蛋，剔透蛋壳凝着霜花，寒蓝微光从内部透出，是冰原的灵胎。冰属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 霜尾幼狐**
- 冰·霜翎狐，凡尘砺心阶段·霜尾幼狐。形象：冰蓝狐狸，尾带霜华。 核心意象：冰蓝毛皮、霜华、极寒雪原。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 凌冰狐**
- 冰·霜翎狐，道法初成阶段·凌冰狐。形象：冰蓝狐狸，尾带霜华。 核心意象：冰蓝毛皮、霜华、极寒雪原。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，雪原上一纵而过，尾尖带起一路霜尘。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 凛冬狐**
- 冰·霜翎狐，大劫淬炼阶段·凛冬狐。形象：冰蓝狐狸，尾带霜华。 核心意象：冰蓝毛皮、霜华、极寒雪原。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 永冬狐**
- 冰·霜翎狐，封神登天阶段·永冬狐。形象：冰蓝狐狸，尾带霜华。 核心意象：冰蓝毛皮、霜华、极寒雪原。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 受封万兽之王，目光睥睨天地; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 极地尊者**
- 冰·霜翎狐，道果圆满阶段·极地尊者。形象：冰蓝狐狸，尾带霜华。 核心意象：冰蓝毛皮、霜华、极寒雪原。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，雪原上一纵而过，尾尖带起一路霜尘。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 神光自照，与天地同尊; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 岩·岩甲犀（`rock_rhino`）

**灵胎初醒 · 岩卵**
- 岩·岩甲犀，灵胎初醒阶段·岩卵。初始形态：一枚岩卵，赭石蛋壳覆着龟裂岩纹，沉重如磐石，山岳之力在内中沉睡。土属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小岩犀**
- 岩·岩甲犀，凡尘砺心阶段·小岩犀。形象：岩石甲身，额生尖角。 核心意象：岩石甲身、额尖角、磐石之重。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冲锋犀**
- 岩·岩甲犀，道法初成阶段·冲锋犀。形象：岩石甲身，额生尖角。 核心意象：岩石甲身、额尖角、磐石之重。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，低头一顶，尖角撞碎前方巨岩。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 铁岩犀**
- 岩·岩甲犀，大劫淬炼阶段·铁岩犀。形象：岩石甲身，额生尖角。 核心意象：岩石甲身、额尖角、磐石之重。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大地犀**
- 岩·岩甲犀，封神登天阶段·大地犀。形象：岩石甲身，额生尖角。 核心意象：岩石甲身、额尖角、磐石之重。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 受封万兽之王，目光睥睨天地; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 岩之尊者**
- 岩·岩甲犀，道果圆满阶段·岩之尊者。形象：岩石甲身，额生尖角。 核心意象：岩石甲身、额尖角、磐石之重。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，低头一顶，尖角撞碎前方巨岩。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 神光自照，与天地同尊; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 风·掠风隼（`wind_falcon`）

**灵胎初醒 · 风卵**
- 风·掠风隼，灵胎初醒阶段·风卵。初始形态：一枚风卵，青白蛋壳绕着流线旋纹，轻如一片浮云，风之羽在壳中蓄势。风属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小风隼**
- 风·掠风隼，凡尘砺心阶段·小风隼。形象：翠绿猛隼，翼疾如风。 核心意象：翠绿翼羽、疾风、万里长空。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 疾风隼**
- 风·掠风隼，道法初成阶段·疾风隼。形象：翠绿猛隼，翼疾如风。 核心意象：翠绿翼羽、疾风、万里长空。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，高空敛翅，一瞬俯冲如风掠地。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 破空隼**
- 风·掠风隼，大劫淬炼阶段·破空隼。形象：翠绿猛隼，翼疾如风。 核心意象：翠绿翼羽、疾风、万里长空。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 王者风隼**
- 风·掠风隼，封神登天阶段·王者风隼。形象：翠绿猛隼，翼疾如风。 核心意象：翠绿翼羽、疾风、万里长空。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 受封万兽之王，目光睥睨天地; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 天空尊者**
- 风·掠风隼，道果圆满阶段·天空尊者。形象：翠绿猛隼，翼疾如风。 核心意象：翠绿翼羽、疾风、万里长空。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，高空敛翅，一瞬俯冲如风掠地。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 神光自照，与天地同尊; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 光·晨曦鹿（`light_deer`）

**灵胎初醒 · 光卵**
- 光·晨曦鹿，灵胎初醒阶段·光卵。初始形态：一枚光卵，莹白蛋壳透着晨曦微光，鹿角般的金色芒纹在壳面浮现。光属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小光鹿**
- 光·晨曦鹿，凡尘砺心阶段·小光鹿。形象：金光神鹿，鹿角如炬。 核心意象：金鹿角、光束、驱暗之光。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 辉光鹿**
- 光·晨曦鹿，道法初成阶段·辉光鹿。形象：金光神鹿，鹿角如炬。 核心意象：金鹿角、光束、驱暗之光。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，鹿角燃起金光，踏光而行驱散阴霾。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 破晓鹿**
- 光·晨曦鹿，大劫淬炼阶段·破晓鹿。形象：金光神鹿，鹿角如炬。 核心意象：金鹿角、光束、驱暗之光。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 辉明鹿**
- 光·晨曦鹿，封神登天阶段·辉明鹿。形象：金光神鹿，鹿角如炬。 核心意象：金鹿角、光束、驱暗之光。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 受封万兽之王，目光睥睨天地; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 光辉尊者**
- 光·晨曦鹿，道果圆满阶段·光辉尊者。形象：金光神鹿，鹿角如炬。 核心意象：金鹿角、光束、驱暗之光。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，鹿角燃起金光，踏光而行驱散阴霾。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 神光自照，与天地同尊; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 暗·夜隐豹（`dark_panther`）

**灵胎初醒 · 暗卵**
- 暗·夜隐豹，灵胎初醒阶段·暗卵。初始形态：一枚暗卵，墨黑蛋壳融入夜色，只有紫色幽光在壳缝间若隐若现。暗属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 夜斑幼豹**
- 暗·夜隐豹，凡尘砺心阶段·夜斑幼豹。形象：墨黑豹影，身姿流线。 核心意象：墨黑皮毛、流线身姿、暗夜利爪。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 潜行豹**
- 暗·夜隐豹，道法初成阶段·潜行豹。形象：墨黑豹影，身姿流线。 核心意象：墨黑皮毛、流线身姿、暗夜利爪。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，身影没入夜色，骤然扑出快如闪电。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 幽冥豹**
- 暗·夜隐豹，大劫淬炼阶段·幽冥豹。形象：墨黑豹影，身姿流线。 核心意象：墨黑皮毛、流线身姿、暗夜利爪。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 暗夜猎豹**
- 暗·夜隐豹，封神登天阶段·暗夜猎豹。形象：墨黑豹影，身姿流线。 核心意象：墨黑皮毛、流线身姿、暗夜利爪。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 受封万兽之王，目光睥睨天地; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 暗夜尊者**
- 暗·夜隐豹，道果圆满阶段·暗夜尊者。形象：墨黑豹影，身姿流线。 核心意象：墨黑皮毛、流线身姿、暗夜利爪。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，身影没入夜色，骤然扑出快如闪电。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 神光自照，与天地同尊; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 钢·磁甲兽（`steel_armadillo`）

**灵胎初醒 · 钢卵**
- 钢·磁甲兽，灵胎初醒阶段·钢卵。初始形态：一枚钢卵，银灰蛋壳覆着装甲鳞纹，金属冷光中透着坚不可摧的气息。钢属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小磁兽**
- 钢·磁甲兽，凡尘砺心阶段·小磁兽。形象：钢甲犰狳，甲带如鳞。 核心意象：钢甲鳞带、铁球之姿、坚不可摧。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 铁甲兽**
- 钢·磁甲兽，道法初成阶段·铁甲兽。形象：钢甲犰狳，甲带如鳞。 核心意象：钢甲鳞带、铁球之姿、坚不可摧。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，蜷身一滚，化作铁球轰然撞出。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 磁暴兽**
- 钢·磁甲兽，大劫淬炼阶段·磁暴兽。形象：钢甲犰狳，甲带如鳞。 核心意象：钢甲鳞带、铁球之姿、坚不可摧。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 合金兽**
- 钢·磁甲兽，封神登天阶段·合金兽。形象：钢甲犰狳，甲带如鳞。 核心意象：钢甲鳞带、铁球之姿、坚不可摧。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 受封万兽之王，目光睥睨天地; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 钢之尊者**
- 钢·磁甲兽，道果圆满阶段·钢之尊者。形象：钢甲犰狳，甲带如鳞。 核心意象：钢甲鳞带、铁球之姿、坚不可摧。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，蜷身一滚，化作铁球轰然撞出。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 神光自照，与天地同尊; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 4. 数码宝贝（6 物种）

> **风格**：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。**阶段演绎**：
> - 灵胎初醒：数码蛋悬浮微光，数据流环绕，崭新世界初启（数据浅蓝/数码银）
> - 凡尘砺心：幼年期毛茸茸的小数码兽，在数码世界好奇探索（数码青/淡蓝）
> - 道法初成：成长进化，数据粒子环绕升腾，光芒凝聚（进化蓝/勇气橙）
> - 大劫淬炼：面对黑暗深渊的试炼，进化之光与黑暗对峙（深渊紫/战斗红）
> - 封神登天：究极进化，圣光加冕，数据粒子凝成铠甲（圣金/纯白）
> - 道果圆满：神圣形态，希望之光普照，数据洪流归于一念（圣光白/虹彩）

#### 亚古兽（`mecha_dragon`）

**灵胎初醒 · 数码蛋**
- 亚古兽，灵胎初醒阶段·数码蛋。初始形态：一枚悬浮的数码蛋，橙色数据流环绕成环，蛋面裂纹透出勇气之火的光芒。火属性灵光微微环绕。神态：数据流中沉睡的懵懂。动作：悬浮微光，数据环绕。衣着：数码蛋，数据纹流转。梳造：数据环光，无固定形。意境：数码蛋悬浮微光，数据流环绕，崭新世界初启。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数据浅蓝（#C9E4FF）主调 + 数码银（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; digital egg; 数据流中沉睡的懵懂; palette #C9E4FF with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 滚球兽**
- 亚古兽，凡尘砺心阶段·滚球兽。形象：橙色小恐龙，背甲坚挺。 核心意象：勇气徽章、橙色背甲、数码蛋。神态：探索数码世界的雀跃。动作：蹦跳探索，好奇触碰。衣着：幼年兽体，毛茸茸。梳造：幼毛/呆毛，圆润可爱。意境：幼年期毛茸茸的小数码兽，在数码世界好奇探索。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数码青（#7EC8E3）主调 + 淡蓝（#A8D8EA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; fresh/child era; 探索数码世界的雀跃; palette #7EC8E3 with #A8D8EA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 亚古兽·干劲**
- 亚古兽，道法初成阶段·亚古兽·干劲。形象：橙色小恐龙，背甲坚挺。 核心意象：勇气徽章、橙色背甲、数码蛋。神态：进化前的坚定，眸光发亮。动作：绝技初现，数据粒子升腾，口中聚焰，一记"小型火焰"喷向敌人。衣着：成熟体装甲，渐渐成形。梳造：头甲初现，眼神明亮。意境：成长进化，数据粒子环绕升腾，光芒凝聚。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：进化蓝（#3F8EFC）主调 + 勇气橙（#FF9E3D）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; evolution spark; 进化前的坚定，眸光发亮; palette #3F8EFC with #FF9E3D accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 暴龙兽·强攻**
- 亚古兽，大劫淬炼阶段·暴龙兽·强攻。形象：橙色小恐龙，背甲坚挺。 核心意象：勇气徽章、橙色背甲、数码蛋。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 机械暴龙兽·强化**
- 亚古兽，封神登天阶段·机械暴龙兽·强化。形象：橙色小恐龙，背甲坚挺。 核心意象：勇气徽章、橙色背甲、数码蛋。神态：究极体受封，神圣威仪。动作：绝技大成，数据凝圣甲。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体受封，神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 战斗暴龙兽·勇气**
- 亚古兽，道果圆满阶段·战斗暴龙兽·勇气。形象：橙色小恐龙，背甲坚挺。 核心意象：勇气徽章、橙色背甲、数码蛋。神态：数据升维，神明降临。动作：绝技轰天，数据化神体，口中聚焰，一记"小型火焰"喷向敌人。衣着：神体数据化，圣光铸形。梳造：数据光环，如日当空。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 数据升维，神明降临; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 迪路兽（`cyber_cat`）

**灵胎初醒 · 数码蛋**
- 迪路兽，灵胎初醒阶段·数码蛋。初始形态：一枚悬浮的数码蛋，圣白数据流结成光环，神圣之息在蛋中凝成爪影。光属性灵光微微环绕。神态：数据流中沉睡的懵懂。动作：悬浮微光，数据环绕。衣着：数码蛋，数据纹流转。梳造：数据环光，无固定形。意境：数码蛋悬浮微光，数据流环绕，崭新世界初启。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数据浅蓝（#C9E4FF）主调 + 数码银（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; digital egg; 数据流中沉睡的懵懂; palette #C9E4FF with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 咪罗兽·成长**
- 迪路兽，凡尘砺心阶段·咪罗兽·成长。形象：白色圣猫，耳戴神圣之环。 核心意象：神圣之环、光明徽章、圣洁之光。神态：探索数码世界的雀跃。动作：蹦跳探索，好奇触碰。衣着：幼年兽体，毛茸茸。梳造：幼毛/呆毛，圆润可爱。意境：幼年期毛茸茸的小数码兽，在数码世界好奇探索。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数码青（#7EC8E3）主调 + 淡蓝（#A8D8EA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; fresh/child era; 探索数码世界的雀跃; palette #7EC8E3 with #A8D8EA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 迪路兽·敏捷**
- 迪路兽，道法初成阶段·迪路兽·敏捷。形象：白色圣猫，耳戴神圣之环。 核心意象：神圣之环、光明徽章、圣洁之光。神态：进化前的坚定，眸光发亮。动作：绝技初现，数据粒子升腾，猫拳连击如电，一记"猫猫拳"撕裂黑暗。衣着：成熟体装甲，渐渐成形。梳造：头甲初现，眼神明亮。意境：成长进化，数据粒子环绕升腾，光芒凝聚。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：进化蓝（#3F8EFC）主调 + 勇气橙（#FF9E3D）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; evolution spark; 进化前的坚定，眸光发亮; palette #3F8EFC with #FF9E3D accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 猫人兽·雷霆**
- 迪路兽，大劫淬炼阶段·猫人兽·雷霆。形象：白色圣猫，耳戴神圣之环。 核心意象：神圣之环、光明徽章、圣洁之光。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 神圣天女兽·圣光**
- 迪路兽，封神登天阶段·神圣天女兽·圣光。形象：白色圣猫，耳戴神圣之环。 核心意象：神圣之环、光明徽章、圣洁之光。神态：究极体受封，神圣威仪。动作：绝技大成，数据凝圣甲。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体受封，神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 圣龙兽·光明**
- 迪路兽，道果圆满阶段·圣龙兽·光明。形象：白色圣猫，耳戴神圣之环。 核心意象：神圣之环、光明徽章、圣洁之光。神态：数据升维，神明降临。动作：绝技轰天，数据化神体，猫拳连击如电，一记"猫猫拳"撕裂黑暗。衣着：神体数据化，圣光铸形。梳造：数据光环，如日当空。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 数据升维，神明降临; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 巴达兽（`space_mecha`）

**灵胎初醒 · 数码蛋**
- 巴达兽，灵胎初醒阶段·数码蛋。初始形态：一枚悬浮的数码蛋，金色数据流如羽翼舒展，希望之光在蛋中盘旋。光属性灵光微微环绕。神态：数据流中沉睡的懵懂。动作：悬浮微光，数据环绕。衣着：数码蛋，数据纹流转。梳造：数据环光，无固定形。意境：数码蛋悬浮微光，数据流环绕，崭新世界初启。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数据浅蓝（#C9E4FF）主调 + 数码银（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; digital egg; 数据流中沉睡的懵懂; palette #C9E4FF with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 浮球兽·成长**
- 巴达兽，凡尘砺心阶段·浮球兽·成长。形象：奶油色圆身，扑腾小翅膀。 核心意象：希望徽章、奶油圆身、神圣羽翼。神态：探索数码世界的雀跃。动作：蹦跳探索，好奇触碰。衣着：幼年兽体，毛茸茸。梳造：幼毛/呆毛，圆润可爱。意境：幼年期毛茸茸的小数码兽，在数码世界好奇探索。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数码青（#7EC8E3）主调 + 淡蓝（#A8D8EA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; fresh/child era; 探索数码世界的雀跃; palette #7EC8E3 with #A8D8EA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 巴达兽·轻盈**
- 巴达兽，道法初成阶段·巴达兽·轻盈。形象：奶油色圆身，扑腾小翅膀。 核心意象：希望徽章、奶油圆身、神圣羽翼。神态：进化前的坚定，眸光发亮。动作：绝技初现，数据粒子升腾，扑扇小翅盘旋，一记"空气炮"轰向敌人。衣着：成熟体装甲，渐渐成形。梳造：头甲初现，眼神明亮。意境：成长进化，数据粒子环绕升腾，光芒凝聚。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：进化蓝（#3F8EFC）主调 + 勇气橙（#FF9E3D）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; evolution spark; 进化前的坚定，眸光发亮; palette #3F8EFC with #FF9E3D accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 天使兽·圣杖**
- 巴达兽，大劫淬炼阶段·天使兽·圣杖。形象：奶油色圆身，扑腾小翅膀。 核心意象：希望徽章、奶油圆身、神圣羽翼。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 神圣天使兽·审判**
- 巴达兽，封神登天阶段·神圣天使兽·审判。形象：奶油色圆身，扑腾小翅膀。 核心意象：希望徽章、奶油圆身、神圣羽翼。神态：究极体受封，神圣威仪。动作：绝技大成，数据凝圣甲。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体受封，神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 炽天使兽·希望**
- 巴达兽，道果圆满阶段·炽天使兽·希望。形象：奶油色圆身，扑腾小翅膀。 核心意象：希望徽章、奶油圆身、神圣羽翼。神态：数据升维，神明降临。动作：绝技轰天，数据化神体，扑扇小翅盘旋，一记"空气炮"轰向敌人。衣着：神体数据化，圣光铸形。梳造：数据光环，如日当空。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 数据升维，神明降临; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 加布兽（`quantum_beast`）

**灵胎初醒 · 数码蛋**
- 加布兽，灵胎初醒阶段·数码蛋。初始形态：一枚悬浮的数码蛋，蓝色数据流凝成霜纹，友情之冰在蛋中静静凝固。冰属性灵光微微环绕。神态：数据流中沉睡的懵懂。动作：悬浮微光，数据环绕。衣着：数码蛋，数据纹流转。梳造：数据环光，无固定形。意境：数码蛋悬浮微光，数据流环绕，崭新世界初启。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数据浅蓝（#C9E4FF）主调 + 数码银（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; digital egg; 数据流中沉睡的懵懂; palette #C9E4FF with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 独角兽·成长**
- 加布兽，凡尘砺心阶段·独角兽·成长。形象：身披蓝毛的小狼，额生独角。 核心意象：友情徽章、独角、蓝色毛皮。神态：探索数码世界的雀跃。动作：蹦跳探索，好奇触碰。衣着：幼年兽体，毛茸茸。梳造：幼毛/呆毛，圆润可爱。意境：幼年期毛茸茸的小数码兽，在数码世界好奇探索。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数码青（#7EC8E3）主调 + 淡蓝（#A8D8EA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; fresh/child era; 探索数码世界的雀跃; palette #7EC8E3 with #A8D8EA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 加布兽·机敏**
- 加布兽，道法初成阶段·加布兽·机敏。形象：身披蓝毛的小狼，额生独角。 核心意象：友情徽章、独角、蓝色毛皮。神态：进化前的坚定，眸光发亮。动作：绝技初现，数据粒子升腾，额角凝聚蓝光，一记"小狼爪"快如闪电。衣着：成熟体装甲，渐渐成形。梳造：头甲初现，眼神明亮。意境：成长进化，数据粒子环绕升腾，光芒凝聚。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：进化蓝（#3F8EFC）主调 + 勇气橙（#FF9E3D）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; evolution spark; 进化前的坚定，眸光发亮; palette #3F8EFC with #FF9E3D accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 加鲁鲁兽·狂岚**
- 加布兽，大劫淬炼阶段·加鲁鲁兽·狂岚。形象：身披蓝毛的小狼，额生独角。 核心意象：友情徽章、独角、蓝色毛皮。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 兽人加鲁鲁·月狼**
- 加布兽，封神登天阶段·兽人加鲁鲁·月狼。形象：身披蓝毛的小狼，额生独角。 核心意象：友情徽章、独角、蓝色毛皮。神态：究极体受封，神圣威仪。动作：绝技大成，数据凝圣甲。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体受封，神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 钢铁加鲁鲁兽·友情**
- 加布兽，道果圆满阶段·钢铁加鲁鲁兽·友情。形象：身披蓝毛的小狼，额生独角。 核心意象：友情徽章、独角、蓝色毛皮。神态：数据升维，神明降临。动作：绝技轰天，数据化神体，额角凝聚蓝光，一记"小狼爪"快如闪电。衣着：神体数据化，圣光铸形。梳造：数据光环，如日当空。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 数据升维，神明降临; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 比丘兽（`digital_phoenix`）

**灵胎初醒 · 数码蛋**
- 比丘兽，灵胎初醒阶段·数码蛋。初始形态：一枚悬浮的数码蛋，翠绿数据流如藤蔓缠绕，歌声的萌动在蛋中孕育。木属性灵光微微环绕。神态：数据流中沉睡的懵懂。动作：悬浮微光，数据环绕。衣着：数码蛋，数据纹流转。梳造：数据环光，无固定形。意境：数码蛋悬浮微光，数据流环绕，崭新世界初启。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数据浅蓝（#C9E4FF）主调 + 数码银（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; digital egg; 数据流中沉睡的懵懂; palette #C9E4FF with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 皮皮兽·成长**
- 比丘兽，凡尘砺心阶段·皮皮兽·成长。形象：粉色小鸟，头顶凤冠。 核心意象：纯真徽章、凤冠、藤蔓之花。神态：探索数码世界的雀跃。动作：蹦跳探索，好奇触碰。衣着：幼年兽体，毛茸茸。梳造：幼毛/呆毛，圆润可爱。意境：幼年期毛茸茸的小数码兽，在数码世界好奇探索。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数码青（#7EC8E3）主调 + 淡蓝（#A8D8EA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; fresh/child era; 探索数码世界的雀跃; palette #7EC8E3 with #A8D8EA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 比丘兽·鸣唱**
- 比丘兽，道法初成阶段·比丘兽·鸣唱。形象：粉色小鸟，头顶凤冠。 核心意象：纯真徽章、凤冠、藤蔓之花。神态：进化前的坚定，眸光发亮。动作：绝技初现，数据粒子升腾，凤冠一闪，一记"毒藤蔓"甩出缠绕敌手。衣着：成熟体装甲，渐渐成形。梳造：头甲初现，眼神明亮。意境：成长进化，数据粒子环绕升腾，光芒凝聚。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：进化蓝（#3F8EFC）主调 + 勇气橙（#FF9E3D）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; evolution spark; 进化前的坚定，眸光发亮; palette #3F8EFC with #FF9E3D accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 巴多拉兽·烈焰**
- 比丘兽，大劫淬炼阶段·巴多拉兽·烈焰。形象：粉色小鸟，头顶凤冠。 核心意象：纯真徽章、凤冠、藤蔓之花。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 伽楼达兽·疾风**
- 比丘兽，封神登天阶段·伽楼达兽·疾风。形象：粉色小鸟，头顶凤冠。 核心意象：纯真徽章、凤冠、藤蔓之花。神态：究极体受封，神圣威仪。动作：绝技大成，数据凝圣甲。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体受封，神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 凤凰兽·涅槃**
- 比丘兽，道果圆满阶段·凤凰兽·涅槃。形象：粉色小鸟，头顶凤冠。 核心意象：纯真徽章、凤冠、藤蔓之花。神态：数据升维，神明降临。动作：绝技轰天，数据化神体，凤冠一闪，一记"毒藤蔓"甩出缠绕敌手。衣着：神体数据化，圣光铸形。梳造：数据光环，如日当空。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 数据升维，神明降临; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 哥玛兽（`mecha_shark`）

**灵胎初醒 · 数码蛋**
- 哥玛兽，灵胎初醒阶段·数码蛋。初始形态：一枚悬浮的数码蛋，海蓝数据流荡出涟漪，诚实的水波在蛋中轻摇。水属性灵光微微环绕。神态：数据流中沉睡的懵懂。动作：悬浮微光，数据环绕。衣着：数码蛋，数据纹流转。梳造：数据环光，无固定形。意境：数码蛋悬浮微光，数据流环绕，崭新世界初启。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数据浅蓝（#C9E4FF）主调 + 数码银（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; digital egg; 数据流中沉睡的懵懂; palette #C9E4FF with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 噗噗兽·成长**
- 哥玛兽，凡尘砺心阶段·噗噗兽·成长。形象：白色海狮，圆滚滚胖乎乎。 核心意象：诚实徽章、圆滚滚、鱼群伙伴。神态：探索数码世界的雀跃。动作：蹦跳探索，好奇触碰。衣着：幼年兽体，毛茸茸。梳造：幼毛/呆毛，圆润可爱。意境：幼年期毛茸茸的小数码兽，在数码世界好奇探索。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：数码青（#7EC8E3）主调 + 淡蓝（#A8D8EA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; fresh/child era; 探索数码世界的雀跃; palette #7EC8E3 with #A8D8EA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 哥玛兽·潜水**
- 哥玛兽，道法初成阶段·哥玛兽·潜水。形象：白色海狮，圆滚滚胖乎乎。 核心意象：诚实徽章、圆滚滚、鱼群伙伴。神态：进化前的坚定，眸光发亮。动作：绝技初现，数据粒子升腾，尾巴一甩，召唤"鱼群大进军"扑向敌阵。衣着：成熟体装甲，渐渐成形。梳造：头甲初现，眼神明亮。意境：成长进化，数据粒子环绕升腾，光芒凝聚。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：进化蓝（#3F8EFC）主调 + 勇气橙（#FF9E3D）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; evolution spark; 进化前的坚定，眸光发亮; palette #3F8EFC with #FF9E3D accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 海狮兽·破浪**
- 哥玛兽，大劫淬炼阶段·海狮兽·破浪。形象：白色海狮，圆滚滚胖乎乎。 核心意象：诚实徽章、圆滚滚、鱼群伙伴。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 祖顿兽·重锤**
- 哥玛兽，封神登天阶段·祖顿兽·重锤。形象：白色海狮，圆滚滚胖乎乎。 核心意象：诚实徽章、圆滚滚、鱼群伙伴。神态：究极体受封，神圣威仪。动作：绝技大成，数据凝圣甲。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体受封，神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 维京兽·咆哮**
- 哥玛兽，道果圆满阶段·维京兽·咆哮。形象：白色海狮，圆滚滚胖乎乎。 核心意象：诚实徽章、圆滚滚、鱼群伙伴。神态：数据升维，神明降临。动作：绝技轰天，数据化神体，尾巴一甩，召唤"鱼群大进军"扑向敌阵。衣着：神体数据化，圣光铸形。梳造：数据光环，如日当空。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 数据升维，神明降临; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 5. 国宝（12 物种）

> **风格**：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。**阶段演绎**：
> - 灵胎初醒：新生命诞生于巢穴或洞穴，温暖微光，柔软新生（巢暖米/自然灰绿）
> - 凡尘砺心：幼崽蹒跚学步，跟随父母觅食，憨态可掬（幼崽棕/林绿）
> - 道法初成：青年期活力四射，嬉戏打闹，探索栖息地（竹林绿/暖阳黄）
> - 大劫淬炼：成年期沉稳担当，守护领地与族群，历经风雨（深林绿/熟褐）
> - 封神登天：成为族群的传奇，祥瑞光环，被守护与敬仰（祥瑞金/王绿）
> - 道果圆满：生态图腾显圣，祥瑞护世，生生不息（云白/生态金绿）

#### 大熊猫（`panda`）

**灵胎初醒 · 粉红团子**
- 大熊猫，灵胎初醒阶段·粉红团子。初始形态：一团粉红胚胎，肉嘟嘟蜷成团子，黑白斑纹在茸毛下若隐若现，散发温润的大地生机。土属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 蹒跚学步**
- 大熊猫，凡尘砺心阶段·蹒跚学步。形象：黑白圆滚，黑耳黑眼圈。 核心意象：黑白皮毛、竹笋、团团圆圆。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 好奇熊猫**
- 大熊猫，道法初成阶段·好奇熊猫。形象：黑白圆滚，黑耳黑眼圈。 核心意象：黑白皮毛、竹笋、团团圆圆。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，抱起竹笋慢悠悠地啃，一个翻身滚下坡。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 青涩少年**
- 大熊猫，大劫淬炼阶段·青涩少年。形象：黑白圆滚，黑耳黑眼圈。 核心意象：黑白皮毛、竹笋、团团圆圆。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 成年熊猫**
- 大熊猫，封神登天阶段·成年熊猫。形象：黑白圆滚，黑耳黑眼圈。 核心意象：黑白皮毛、竹笋、团团圆圆。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 太极熊猫**
- 大熊猫，道果圆满阶段·太极熊猫。形象：黑白圆滚，黑耳黑眼圈。 核心意象：黑白皮毛、竹笋、团团圆圆。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，抱起竹笋慢悠悠地啃，一个翻身滚下坡。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 金丝猴（`golden_monkey`）

**灵胎初醒 · 金丝灵胎**
- 金丝猴，灵胎初醒阶段·金丝灵胎。初始形态：一团金丝灵胎，淡金绒毛初生，蓝脸的轮廓在暖光中隐约，攀山越涧之性已在血脉。金属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 初显金毛**
- 金丝猴，凡尘砺心阶段·初显金毛。形象：金色长毛，蓝脸朝天鼻。 核心意象：金色长毛、蓝脸、高山密林。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 蓝面初显**
- 金丝猴，道法初成阶段·蓝面初显。形象：金色长毛，蓝脸朝天鼻。 核心意象：金色长毛、蓝脸、高山密林。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，抓住藤蔓荡起，一纵数米穿行林间。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 长尾猴**
- 金丝猴，大劫淬炼阶段·长尾猴。形象：金色长毛，蓝脸朝天鼻。 核心意象：金色长毛、蓝脸、高山密林。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金猴王**
- 金丝猴，封神登天阶段·金猴王。形象：金色长毛，蓝脸朝天鼻。 核心意象：金色长毛、蓝脸、高山密林。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 金丝大圣**
- 金丝猴，道果圆满阶段·金丝大圣。形象：金色长毛，蓝脸朝天鼻。 核心意象：金色长毛、蓝脸、高山密林。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，抓住藤蔓荡起，一纵数米穿行林间。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 丹顶鹤（`red_crowned_crane`）

**灵胎初醒 · 鹤卵**
- 丹顶鹤，灵胎初醒阶段·鹤卵。初始形态：一枚鹤卵，青灰蛋壳缀着云纹，一点朱红在壳顶隐现，鹤唳之音在蛋中回荡。冰属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 灰羽渐白**
- 丹顶鹤，凡尘砺心阶段·灰羽渐白。形象：白羽黑颈，头顶一点朱红。 核心意象：头顶朱红、白羽黑颈、松鹤延年。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 飞羽初展**
- 丹顶鹤，道法初成阶段·飞羽初展。形象：白羽黑颈，头顶一点朱红。 核心意象：头顶朱红、白羽黑颈、松鹤延年。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，引颈长鸣，展翅起舞，声闻九天。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 长寿鹤**
- 丹顶鹤，大劫淬炼阶段·长寿鹤。形象：白羽黑颈，头顶一点朱红。 核心意象：头顶朱红、白羽黑颈、松鹤延年。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 祥云伴鹤**
- 丹顶鹤，封神登天阶段·祥云伴鹤。形象：白羽黑颈，头顶一点朱红。 核心意象：头顶朱红、白羽黑颈、松鹤延年。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 羽化鹤**
- 丹顶鹤，道果圆满阶段·羽化鹤。形象：白羽黑颈，头顶一点朱红。 核心意象：头顶朱红、白羽黑颈、松鹤延年。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，引颈长鸣，展翅起舞，声闻九天。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 华南虎（`south_china_tiger`）

**灵胎初醒 · 虎纹胚胎**
- 华南虎，灵胎初醒阶段·虎纹胚胎。初始形态：一团虎纹胚胎，橙黑斑纹在微光中初显，额间"王"字的虚影若隐若现。金属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼虎**
- 华南虎，凡尘砺心阶段·幼虎。形象：橙底黑纹，身形矫健。 核心意象：橙底黑纹、额间王字、密林之王。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 斑纹加深**
- 华南虎，道法初成阶段·斑纹加深。形象：橙底黑纹，身形矫健。 核心意象：橙底黑纹、额间王字、密林之王。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，夜幕中悄然潜行，一声虎啸震慑山林。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 百兽之王**
- 华南虎，大劫淬炼阶段·百兽之王。形象：橙底黑纹，身形矫健。 核心意象：橙底黑纹、额间王字、密林之王。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金虎**
- 华南虎，封神登天阶段·金虎。形象：橙底黑纹，身形矫健。 核心意象：橙底黑纹、额间王字、密林之王。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 神虎**
- 华南虎，道果圆满阶段·神虎。形象：橙底黑纹，身形矫健。 核心意象：橙底黑纹、额间王字、密林之王。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，夜幕中悄然潜行，一声虎啸震慑山林。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 扬子鳄（`chinese_alligator`）

**灵胎初醒 · 鳄卵**
- 扬子鳄，灵胎初醒阶段·鳄卵。初始形态：一枚鳄卵，米白蛋壳覆着细密鳞纹，半截小鳄的剪影在壳中蜷伏。水属性灵光微微环绕。神态：龙息轻吐，沉睡于混沌灵光。动作：盘蜷于光，沉眠未醒。衣着：幼嫩鳞胚，泛着初光。梳造：无角，须影初现。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 龙息轻吐，沉睡于混沌灵光; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 鳞甲初生**
- 扬子鳄，凡尘砺心阶段·鳞甲初生。形象：低伏的宽吻鳄龙，身披硬鳞。 核心意象：硬鳞、宽吻、江河泥潭。神态：竖瞳初睁，窥探云水。动作：初探云水，游弋学步。衣着：鳞甲渐密，颜色初显。梳造：角芽微露，须丝轻扬。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 竖瞳初睁，窥探云水; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 水中霸主**
- 扬子鳄，道法初成阶段·水中霸主。形象：低伏的宽吻鳄龙，身披硬鳞。 核心意象：硬鳞、宽吻、江河泥潭。神态：眸光锐亮，潜龙欲腾。动作：腾空而起，初显神威，水中悄然潜行，猛然一扑咬定猎物。衣着：鳞甲生辉，腹光流转。梳造：双角初成，须如流云。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光锐亮，潜龙欲腾; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 铁甲鳄**
- 扬子鳄，大劫淬炼阶段·铁甲鳄。形象：低伏的宽吻鳄龙，身披硬鳞。 核心意象：硬鳞、宽吻、江河泥潭。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 龙威炽烈，怒目电光; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大鳄**
- 扬子鳄，封神登天阶段·大鳄。形象：低伏的宽吻鳄龙，身披硬鳞。 核心意象：硬鳞、宽吻、江河泥潭。神态：受封龙君，神威赫赫。动作：登天行云，布雨泽四方。衣着：金鳞覆身，受冕祥光。梳造：龙角如珊瑚，加冕为尊。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封龙君，神威赫赫; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 蛟鳄**
- 扬子鳄，道果圆满阶段·蛟鳄。形象：低伏的宽吻鳄龙，身披硬鳞。 核心意象：硬鳞、宽吻、江河泥潭。神态：万龙之源，睥睨三界。动作：真身化岳，日月为伴，水中悄然潜行，猛然一扑咬定猎物。衣着：龙身映日月，鳞甲如星辰。梳造：龙角擎天，道纹绕体。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 万龙之源，睥睨三界; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 朱鹮（`crested_ibis`）

**灵胎初醒 · 朱卵**
- 朱鹮，灵胎初醒阶段·朱卵。初始形态：一枚朱卵，米白蛋壳晕着一抹绯红，东方宝石的祥光在壳面流转。火属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 白羽渐生**
- 朱鹮，凡尘砺心阶段·白羽渐生。形象：白羽粉翼，头冠蓬松。 核心意象：粉红羽翼、长喙、水田浅滩。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 朱鹮**
- 朱鹮，道法初成阶段·朱鹮。形象：白羽粉翼，头冠蓬松。 核心意象：粉红羽翼、长喙、水田浅滩。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，长喙探入浅水，轻巧地叼起鱼虾。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 祥瑞之鸟**
- 朱鹮，大劫淬炼阶段·祥瑞之鸟。形象：白羽粉翼，头冠蓬松。 核心意象：粉红羽翼、长喙、水田浅滩。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 吉祥朱鹮**
- 朱鹮，封神登天阶段·吉祥朱鹮。形象：白羽粉翼，头冠蓬松。 核心意象：粉红羽翼、长喙、水田浅滩。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 凤凰近亲**
- 朱鹮，道果圆满阶段·凤凰近亲。形象：白羽粉翼，头冠蓬松。 核心意象：粉红羽翼、长喙、水田浅滩。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，长喙探入浅水，轻巧地叼起鱼虾。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 藏羚羊（`tibetan_antelope`）

**灵胎初醒 · 雪原灵胎**
- 藏羚羊，灵胎初醒阶段·雪原灵胎。初始形态：一团雪原灵胎，白褐绒毛初生，剑形长角的一点雏形在寒光中显现。冰属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小羚羊**
- 藏羚羊，凡尘砺心阶段·小羚羊。形象：白褐相间，长角如剑。 核心意象：剑形长角、白褐毛皮、可可西里。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 高原羚**
- 藏羚羊，道法初成阶段·高原羚。形象：白褐相间，长角如剑。 核心意象：剑形长角、白褐毛皮、可可西里。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，四蹄翻飞，在高原上一路狂奔如飞。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 奔风羚**
- 藏羚羊，大劫淬炼阶段·奔风羚。形象：白褐相间，长角如剑。 核心意象：剑形长角、白褐毛皮、可可西里。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 雪原之王**
- 藏羚羊，封神登天阶段·雪原之王。形象：白褐相间，长角如剑。 核心意象：剑形长角、白褐毛皮、可可西里。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 雪原尊者**
- 藏羚羊，道果圆满阶段·雪原尊者。形象：白褐相间，长角如剑。 核心意象：剑形长角、白褐毛皮、可可西里。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，四蹄翻飞，在高原上一路狂奔如飞。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 雪豹（`snow_leopard`）

**灵胎初醒 · 斑雪灵胎**
- 雪豹，灵胎初醒阶段·斑雪灵胎。初始形态：一团斑雪灵胎，灰白绒毛缀着黑斑，雪山之王的冷冽气息在胚胎中沉淀。冰属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小雪豹**
- 雪豹，凡尘砺心阶段·小雪豹。形象：灰白底黑斑，长尾粗壮。 核心意象：灰白斑纹、粗壮长尾、雪山之巅。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 雪山豹**
- 雪豹，道法初成阶段·雪山豹。形象：灰白底黑斑，长尾粗壮。 核心意象：灰白斑纹、粗壮长尾、雪山之巅。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，长尾平衡，纵身跃过悬崖峭壁。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 雪原豹**
- 雪豹，大劫淬炼阶段·雪原豹。形象：灰白底黑斑，长尾粗壮。 核心意象：灰白斑纹、粗壮长尾、雪山之巅。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 高山雪豹**
- 雪豹，封神登天阶段·高山雪豹。形象：灰白底黑斑，长尾粗壮。 核心意象：灰白斑纹、粗壮长尾、雪山之巅。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 雪山尊者**
- 雪豹，道果圆满阶段·雪山尊者。形象：灰白底黑斑，长尾粗壮。 核心意象：灰白斑纹、粗壮长尾、雪山之巅。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，长尾平衡，纵身跃过悬崖峭壁。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 麋鹿（`milu_deer`）

**灵胎初醒 · 泽畔灵胎**
- 麋鹿，灵胎初醒阶段·泽畔灵胎。初始形态：一团泽畔灵胎，浅褐绒毛沾着水汽，四不像的角影在湿地灵光中若隐若现。水属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小麋鹿**
- 麋鹿，凡尘砺心阶段·小麋鹿。形象：角似鹿非鹿，蹄似牛非牛。 核心意象：四不像之形、水泽湿地、角蹄。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 沼泽鹿**
- 麋鹿，道法初成阶段·沼泽鹿。形象：角似鹿非鹿，蹄似牛非牛。 核心意象：四不像之形、水泽湿地、角蹄。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，踏水而行，慢悠悠地在湿地踱步。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 泽畔鹿**
- 麋鹿，大劫淬炼阶段·泽畔鹿。形象：角似鹿非鹿，蹄似牛非牛。 核心意象：四不像之形、水泽湿地、角蹄。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 瑞泽鹿**
- 麋鹿，封神登天阶段·瑞泽鹿。形象：角似鹿非鹿，蹄似牛非牛。 核心意象：四不像之形、水泽湿地、角蹄。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 泽畔尊者**
- 麋鹿，道果圆满阶段·泽畔尊者。形象：角似鹿非鹿，蹄似牛非牛。 核心意象：四不像之形、水泽湿地、角蹄。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，踏水而行，慢悠悠地在湿地踱步。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 东北虎（`siberian_tiger`）

**灵胎初醒 · 雪虎灵胎**
- 东北虎，灵胎初醒阶段·雪虎灵胎。初始形态：一团雪虎灵胎，厚绒毛泛着冰晶微光，林海雪原的威压已在胚胎中酝酿。冰属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小东北虎**
- 东北虎，凡尘砺心阶段·小东北虎。形象：硕大虎躯，厚毛如披风。 核心意象：厚毛披风、硕大虎躯、林海雪原。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 雪原虎**
- 东北虎，道法初成阶段·雪原虎。形象：硕大虎躯，厚毛如披风。 核心意象：厚毛披风、硕大虎躯、林海雪原。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，一声长啸震彻林海，威压百兽。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 啸雪虎**
- 东北虎，大劫淬炼阶段·啸雪虎。形象：硕大虎躯，厚毛如披风。 核心意象：厚毛披风、硕大虎躯、林海雪原。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 雪原巨虎**
- 东北虎，封神登天阶段·雪原巨虎。形象：硕大虎躯，厚毛如披风。 核心意象：厚毛披风、硕大虎躯、林海雪原。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 雪原尊者**
- 东北虎，道果圆满阶段·雪原尊者。形象：硕大虎躯，厚毛如披风。 核心意象：厚毛披风、硕大虎躯、林海雪原。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，一声长啸震彻林海，威压百兽。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 小熊猫（`red_panda`）

**灵胎初醒 · 红绒灵胎**
- 小熊猫，灵胎初醒阶段·红绒灵胎。初始形态：一团红绒灵胎，火红绒毛裹着竹香，环纹尾巴的雏形在暖光中蜷起。木属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小熊猫**
- 小熊猫，凡尘砺心阶段·小熊猫。形象：红褐毛皮，环纹大尾。 核心意象：红褐毛皮、环纹大尾、竹叶。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 红棕熊**
- 小熊猫，道法初成阶段·红棕熊。形象：红褐毛皮，环纹大尾。 核心意象：红褐毛皮、环纹大尾、竹叶。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，竖起环纹大尾，抱着一捧竹叶细嚼。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 红锦猫熊**
- 小熊猫，大劫淬炼阶段·红锦猫熊。形象：红褐毛皮，环纹大尾。 核心意象：红褐毛皮、环纹大尾、竹叶。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 祥瑞小熊**
- 小熊猫，封神登天阶段·祥瑞小熊。形象：红褐毛皮，环纹大尾。 核心意象：红褐毛皮、环纹大尾、竹叶。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 竹林尊者**
- 小熊猫，道果圆满阶段·竹林尊者。形象：红褐毛皮，环纹大尾。 核心意象：红褐毛皮、环纹大尾、竹叶。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，竖起环纹大尾，抱着一捧竹叶细嚼。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 江豚（`finless_porpoise`）

**灵胎初醒 · 江流灵胎**
- 江豚，灵胎初醒阶段·江流灵胎。初始形态：一团江流灵胎，圆头灰蓝泛着水光，嘴角的微笑弧度在江涛中若隐若现。水属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：新生命诞生于巢穴或洞穴，温暖微光，柔软新生。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：巢暖米（#F5E9DC）主调 + 自然灰绿（#A8B8B0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; newborn glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E9DC with #A8B8B0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小江豚**
- 江豚，凡尘砺心阶段·小江豚。形象：圆头无背鳍，灰蓝光滑。 核心意象：无背鳍圆头、江涛、微笑。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼崽蹒跚学步，跟随父母觅食，憨态可掬。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：幼崽棕（#C9B79C）主调 + 林绿（#8FA8A0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; wobbly cub; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #C9B79C with #8FA8A0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 江中精灵**
- 江豚，道法初成阶段·江中精灵。形象：圆头无背鳍，灰蓝光滑。 核心意象：无背鳍圆头、江涛、微笑。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，跃出江面，嘴角的弧度像一抹微笑。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：青年期活力四射，嬉戏打闹，探索栖息地。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：竹林绿（#6B8E7A）主调 + 暖阳黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; playful youth; 眸光晶亮，意气初显，跃跃欲试; palette #6B8E7A with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 微笑天使**
- 江豚，大劫淬炼阶段·微笑天使。形象：圆头无背鳍，灰蓝光滑。 核心意象：无背鳍圆头、江涛、微笑。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 江中守护**
- 江豚，封神登天阶段·江中守护。形象：圆头无背鳍，灰蓝光滑。 核心意象：无背鳍圆头、江涛、微笑。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 受封万兽之王，目光睥睨天地; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 江流尊者**
- 江豚，道果圆满阶段·江流尊者。形象：圆头无背鳍，灰蓝光滑。 核心意象：无背鳍圆头、江涛、微笑。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，跃出江面，嘴角的弧度像一抹微笑。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 神光自照，与天地同尊; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 6. 魔法奇幻（12 物种）

> **风格**：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。**阶段演绎**：
> - 灵胎初醒：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生（奥秘紫/魔法薰衣草）
> - 凡尘砺心：初学魔法，好奇探索，小小的魔力火花（魔紫/星光蓝）
> - 道法初成：魔力觉醒，脚下浮现法阵，魔法光辉绽放（深魔紫/咒文金）
> - 大劫淬炼：黑暗试炼、魔力对决，在深渊边缘淬炼（暗夜紫/火焰红）
> - 封神登天：大法师或神话生物完全体，威震大陆，法阵与冠冕加身（传奇紫/王冠金）
> - 道果圆满：元素法相圆满，光与秘术铸就神格，威震大陆（永恒白/圣光紫）

#### 独角兽（`unicorn`）

**灵胎初醒 · 星光种子**
- 独角兽，灵胎初醒阶段·星光种子。初始形态：一粒星光种子，银白灵光凝成独角雏形，圣洁芒星在种壳上闪烁。光属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 灵光中沉睡，兽性未醒的宁静; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 白色小驹**
- 独角兽，凡尘砺心阶段·白色小驹。形象：纯白骏马，额生螺旋独角。 核心意象：螺旋独角、纯白毛皮、圣光。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 初踏山川，懵懂而灵慧; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 角光初现**
- 独角兽，道法初成阶段·角光初现。形象：纯白骏马，额生螺旋独角。 核心意象：螺旋独角、纯白毛皮、圣光。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，踏云而来，独角泛起一圈圣洁光环。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 眼神古老而专注，神通初显; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 彩虹鬃毛**
- 独角兽，大劫淬炼阶段·彩虹鬃毛。形象：纯白骏马，额生螺旋独角。 核心意象：螺旋独角、纯白毛皮、圣光。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 眸光如电，威严中带着坚韧; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 星夜独角兽**
- 独角兽，封神登天阶段·星夜独角兽。形象：纯白骏马，额生螺旋独角。 核心意象：螺旋独角、纯白毛皮、圣光。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封瑞兽，祥云拱卫; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 永恒独角兽**
- 独角兽，道果圆满阶段·永恒独角兽。形象：纯白骏马，额生螺旋独角。 核心意象：螺旋独角、纯白毛皮、圣光。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，踏云而来，独角泛起一圈圣洁光环。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 瑞气化道，福泽天地; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 飞龙（`wyvern`）

**灵胎初醒 · 龙蛋**
- 飞龙，灵胎初醒阶段·龙蛋。初始形态：一枚龙蛋，赤色蛋壳覆着岩浆裂纹，龙翼的虚影在火光中鼓动。火属性灵光微微环绕。神态：龙息轻吐，沉睡于混沌灵光。动作：盘蜷于光，沉眠未醒。衣着：幼嫩鳞胚，泛着初光。梳造：无角，须影初现。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 龙息轻吐，沉睡于混沌灵光; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 翼芽**
- 飞龙，凡尘砺心阶段·翼芽。形象：赤色飞龙，蝠翼如膜，尾带尖刺。 核心意象：蝠翼、尖刺长尾、火山烈焰。神态：竖瞳初睁，窥探云水。动作：初探云水，游弋学步。衣着：鳞甲渐密，颜色初显。梳造：角芽微露，须丝轻扬。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 竖瞳初睁，窥探云水; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 火焰龙**
- 飞龙，道法初成阶段·火焰龙。形象：赤色飞龙，蝠翼如膜，尾带尖刺。 核心意象：蝠翼、尖刺长尾、火山烈焰。神态：眸光锐亮，潜龙欲腾。动作：腾空而起，初显神威，双翼一振腾空，喷吐烈焰焚尽大地。衣着：鳞甲生辉，腹光流转。梳造：双角初成，须如流云。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 眸光锐亮，潜龙欲腾; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 雷龙**
- 飞龙，大劫淬炼阶段·雷龙。形象：赤色飞龙，蝠翼如膜，尾带尖刺。 核心意象：蝠翼、尖刺长尾、火山烈焰。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 龙威炽烈，怒目电光; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 古龙**
- 飞龙，封神登天阶段·古龙。形象：赤色飞龙，蝠翼如膜，尾带尖刺。 核心意象：蝠翼、尖刺长尾、火山烈焰。神态：受封龙君，神威赫赫。动作：登天行云，布雨泽四方。衣着：金鳞覆身，受冕祥光。梳造：龙角如珊瑚，加冕为尊。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封龙君，神威赫赫; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 太古龙**
- 飞龙，道果圆满阶段·太古龙。形象：赤色飞龙，蝠翼如膜，尾带尖刺。 核心意象：蝠翼、尖刺长尾、火山烈焰。神态：万龙之源，睥睨三界。动作：真身化岳，日月为伴，双翼一振腾空，喷吐烈焰焚尽大地。衣着：龙身映日月，鳞甲如星辰。梳造：龙角擎天，道纹绕体。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 万龙之源，睥睨三界; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 精灵（`fairy`）

**灵胎初醒 · 花苞**
- 精灵，灵胎初醒阶段·花苞。初始形态：一枚花苞灵种，金色花蕊在晨露中初绽，精灵翅影藏在花瓣之间。木属性灵光微微环绕。神态：沉睡的种子里，藏着破土的心跳。动作：静卧泥土，待春雨而萌。衣着：种子/种壳，纹理细腻。梳造：顶端嫩芽微顶。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 沉睡的种子里，藏着破土的心跳; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 叶精灵**
- 精灵，凡尘砺心阶段·叶精灵。形象：金翼小仙子，身绕光点。 核心意象：金翼、光点、花间飞舞。神态：初生的好奇，向着光的方向。动作：嫩芽破土，努力伸展。衣着：幼芽新叶，娇嫩欲滴。梳造：顶芽鲜嫩，两片子叶。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 初生的好奇，向着光的方向; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 暗精灵**
- 精灵，道法初成阶段·暗精灵。形象：金翼小仙子，身绕光点。 核心意象：金翼、光点、花间飞舞。神态：生机勃发，神采奕奕。动作：生机初现，枝叶舒展，金翼轻振，洒下点点星尘与花粉。衣着：茎叶渐盛，花苞初成。梳造：花苞/嫩叶环生。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 生机勃发，神采奕奕; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 火精灵**
- 精灵，大劫淬炼阶段·火精灵。形象：金翼小仙子，身绕光点。 核心意象：金翼、光点、花间飞舞。神态：风雨中的坚韧，眸光不折。动作：全力生长，根系深扎。衣着：枝叶繁茂，带风霜痕迹。梳造：花叶翻卷，仍自挺立。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 风雨中的坚韧，眸光不折; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大地精灵**
- 精灵，封神登天阶段·大地精灵。形象：金翼小仙子，身绕光点。 核心意象：金翼、光点、花间飞舞。神态：受封花王，生机威仪。动作：生机大成，繁花满枝。衣着：花开满枝，树冠如云。梳造：花冠受冕，光华流转。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封花王，生机威仪; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 森林之子**
- 精灵，道果圆满阶段·森林之子。形象：金翼小仙子，身绕光点。 核心意象：金翼、光点、花间飞舞。神态：万灵朝拜，生命之尊。动作：花开万里，果实垂天，金翼轻振，洒下点点星尘与花粉。衣着：参天古木，藤蔓缠霄。梳造：万叶化冠，春永驻世。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 万灵朝拜，生命之尊; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 树人（`treant`）

**灵胎初醒 · 种子**
- 树人，灵胎初醒阶段·种子。初始形态：一粒橡树种子，深褐种壳裂出嫩绿新芽，森林的脉搏在根须间跳动。木属性灵光微微环绕。神态：沉睡的种子里，藏着破土的心跳。动作：静卧泥土，待春雨而萌。衣着：种子/种壳，纹理细腻。梳造：顶端嫩芽微顶。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 沉睡的种子里，藏着破土的心跳; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小树苗**
- 树人，凡尘砺心阶段·小树苗。形象：深绿树人，枝干为臂。 核心意象：深绿树皮、年轮、森林之魂。神态：初生的好奇，向着光的方向。动作：嫩芽破土，努力伸展。衣着：幼芽新叶，娇嫩欲滴。梳造：顶芽鲜嫩，两片子叶。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 初生的好奇，向着光的方向; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 树人**
- 树人，道法初成阶段·树人。形象：深绿树人，枝干为臂。 核心意象：深绿树皮、年轮、森林之魂。神态：生机勃发，神采奕奕。动作：生机初现，枝叶舒展，大地震颤着迈步，枝条如臂挥展。衣着：茎叶渐盛，花苞初成。梳造：花苞/嫩叶环生。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 生机勃发，神采奕奕; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 古树人**
- 树人，大劫淬炼阶段·古树人。形象：深绿树人，枝干为臂。 核心意象：深绿树皮、年轮、森林之魂。神态：风雨中的坚韧，眸光不折。动作：全力生长，根系深扎。衣着：枝叶繁茂，带风霜痕迹。梳造：花叶翻卷，仍自挺立。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 风雨中的坚韧，眸光不折; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 森林卫士**
- 树人，封神登天阶段·森林卫士。形象：深绿树人，枝干为臂。 核心意象：深绿树皮、年轮、森林之魂。神态：受封花王，生机威仪。动作：生机大成，繁花满枝。衣着：花开满枝，树冠如云。梳造：花冠受冕，光华流转。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封花王，生机威仪; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 世界树幼苗**
- 树人，道果圆满阶段·世界树幼苗。形象：深绿树人，枝干为臂。 核心意象：深绿树皮、年轮、森林之魂。神态：万灵朝拜，生命之尊。动作：花开万里，果实垂天，大地震颤着迈步，枝条如臂挥展。衣着：参天古木，藤蔓缠霄。梳造：万叶化冠，春永驻世。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 万灵朝拜，生命之尊; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 狮鹫（`griffin`）

**灵胎初醒 · 狮鹫蛋**
- 狮鹫，灵胎初醒阶段·狮鹫蛋。初始形态：一枚狮鹫蛋，金白蛋壳生着细羽纹，鹰喙与狮爪的雏形在蛋中蓄势。金属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 灵光中沉睡，兽性未醒的宁静; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 羽翼初长**
- 狮鹫，凡尘砺心阶段·羽翼初长。形象：鹰首狮身，金色双翼。 核心意象：鹰首狮身、金色双翼、忠勇。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 初踏山川，懵懂而灵慧; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 金翼狮鹫**
- 狮鹫，道法初成阶段·金翼狮鹫。形象：鹰首狮身，金色双翼。 核心意象：鹰首狮身、金色双翼、忠勇。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，双翼一振冲天，俯冲时迅如流星。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 眼神古老而专注，神通初显; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 皇家狮鹫**
- 狮鹫，大劫淬炼阶段·皇家狮鹫。形象：鹰首狮身，金色双翼。 核心意象：鹰首狮身、金色双翼、忠勇。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 眸光如电，威严中带着坚韧; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 雷光狮鹫**
- 狮鹫，封神登天阶段·雷光狮鹫。形象：鹰首狮身，金色双翼。 核心意象：鹰首狮身、金色双翼、忠勇。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封瑞兽，祥云拱卫; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 神圣狮鹫**
- 狮鹫，道果圆满阶段·神圣狮鹫。形象：鹰首狮身，金色双翼。 核心意象：鹰首狮身、金色双翼、忠勇。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，双翼一振冲天，俯冲时迅如流星。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 瑞气化道，福泽天地; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 美人鱼（`mermaid`）

**灵胎初醒 · 珍珠**
- 美人鱼，灵胎初醒阶段·珍珠。初始形态：一枚深海珍珠，碧光流转如潮水，人鱼歌谣的涟漪在珠身内漾开。水属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 灵光中沉睡，兽性未醒的宁静; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 鳞光**
- 美人鱼，凡尘砺心阶段·鳞光。形象：碧鳞人鱼，长发如波。 核心意象：碧鳞鱼尾、如波长发、深海歌声。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 初踏山川，懵懂而灵慧; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 海妖**
- 美人鱼，道法初成阶段·海妖。形象：碧鳞人鱼，长发如波。 核心意象：碧鳞鱼尾、如波长发、深海歌声。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，鱼尾轻摆破浪，一曲歌声飘向海面。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 眼神古老而专注，神通初显; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 珊瑚公主**
- 美人鱼，大劫淬炼阶段·珊瑚公主。形象：碧鳞人鱼，长发如波。 核心意象：碧鳞鱼尾、如波长发、深海歌声。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 眸光如电，威严中带着坚韧; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 人鱼女王**
- 美人鱼，封神登天阶段·人鱼女王。形象：碧鳞人鱼，长发如波。 核心意象：碧鳞鱼尾、如波长发、深海歌声。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封瑞兽，祥云拱卫; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 大洋之魂**
- 美人鱼，道果圆满阶段·大洋之魂。形象：碧鳞人鱼，长发如波。 核心意象：碧鳞鱼尾、如波长发、深海歌声。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，鱼尾轻摆破浪，一曲歌声飘向海面。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 瑞气化道，福泽天地; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 灰袍智者（`grey_wizard`）

**灵胎初醒 · 魔杖种子**
- 灰袍智者，灵胎初醒阶段·魔杖种子。初始形态：一粒魔杖种子，星芒符文绕种盘旋，灰袍贤者的灵光在奥秘中浮现。秘属性灵光微微环绕。神态：星芒中蜷缩，似梦似醒。动作：悬于法阵微光，静候唤醒。衣着：法袍虚影，星光为衣。梳造：灵光发丝，半透明。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 星芒中蜷缩，似梦似醒; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小法师**
- 灰袍智者，凡尘砺心阶段·小法师。形象：灰袍巫师，手持法杖。 核心意象：灰袍、法杖、奥术符文。神态：好奇试探的灵动。动作：吟唱初学，魔力火花迸溅。衣着：学徒法袍，布料朴素。梳造：发丝扬起，缀着星尘。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 好奇试探的灵动; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 灰袍法师**
- 灰袍智者，道法初成阶段·灰袍法师。形象：灰袍巫师，手持法杖。 核心意象：灰袍、法杖、奥术符文。神态：光华渐盛，神采焕发。动作：魔力初现，法阵脚下亮起，法杖顿地，口中吟唱，星光自天穹应召而下。衣着：法师长袍，绣星纹。梳造：魔法帽/发冠初戴。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 光华渐盛，神采焕发; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 魔典持有者**
- 灰袍智者，大劫淬炼阶段·魔典持有者。形象：灰袍巫师，手持法杖。 核心意象：灰袍、法杖、奥术符文。神态：暗夜中坚持，眼神倔强。动作：全力施法，魔力风暴。衣着：法袍破损，魔力燃烧。梳造：发丝凌乱，眸中带焰。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 暗夜中坚持，眼神倔强; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 白袍法师**
- 灰袍智者，封神登天阶段·白袍法师。形象：灰袍巫师，手持法杖。 核心意象：灰袍、法杖、奥术符文。神态：受封贤者，光晕庄严。动作：魔法大成，法阵漫天。衣着：贤者白袍，权杖受命。梳造：法冠威仪，光晕庄严。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封贤者，光晕庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 白袍尊者**
- 灰袍智者，道果圆满阶段·白袍尊者。形象：灰袍巫师，手持法杖。 核心意象：灰袍、法杖、奥术符文。神态：与元素同源，法身自成。动作：法阵通天，咒文如潮，法杖顿地，口中吟唱，星光自天穹应召而下。衣着：元素法身，光尘化袍。梳造：星冕加身，魔力如渊。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 与元素同源，法身自成; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 魔杖猫（`wand_cat`）

**灵胎初醒 · 魔法绒球**
- 魔杖猫，灵胎初醒阶段·魔法绒球。初始形态：一团魔法绒球，紫毛间缀着星点，魔杖的微光在绒球中央轻轻闪动。秘属性灵光微微环绕。神态：星芒中蜷缩，似梦似醒。动作：悬于法阵微光，静候唤醒。衣着：法袍虚影，星光为衣。梳造：灵光发丝，半透明。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 星芒中蜷缩，似梦似醒; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 会魔法的猫**
- 魔杖猫，凡尘砺心阶段·会魔法的猫。形象：紫毛小猫，爪握魔法杖。 核心意象：紫色毛皮、魔法杖、闪烁咒光。神态：好奇试探的灵动。动作：吟唱初学，魔力火花迸溅。衣着：学徒法袍，布料朴素。梳造：发丝扬起，缀着星尘。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 好奇试探的灵动; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 魔法猫**
- 魔杖猫，道法初成阶段·魔法猫。形象：紫毛小猫，爪握魔法杖。 核心意象：紫色毛皮、魔法杖、闪烁咒光。神态：光华渐盛，神采焕发。动作：魔力初现，法阵脚下亮起，挥舞魔法杖，星光自杖尖倾泻而下。衣着：法师长袍，绣星纹。梳造：魔法帽/发冠初戴。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 光华渐盛，神采焕发; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 魔咒猫**
- 魔杖猫，大劫淬炼阶段·魔咒猫。形象：紫毛小猫，爪握魔法杖。 核心意象：紫色毛皮、魔法杖、闪烁咒光。神态：暗夜中坚持，眼神倔强。动作：全力施法，魔力风暴。衣着：法袍破损，魔力燃烧。梳造：发丝凌乱，眸中带焰。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 暗夜中坚持，眼神倔强; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大贤者猫**
- 魔杖猫，封神登天阶段·大贤者猫。形象：紫毛小猫，爪握魔法杖。 核心意象：紫色毛皮、魔法杖、闪烁咒光。神态：受封贤者，光晕庄严。动作：魔法大成，法阵漫天。衣着：贤者白袍，权杖受命。梳造：法冠威仪，光晕庄严。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封贤者，光晕庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 魔法尊者**
- 魔杖猫，道果圆满阶段·魔法尊者。形象：紫毛小猫，爪握魔法杖。 核心意象：紫色毛皮、魔法杖、闪烁咒光。神态：与元素同源，法身自成。动作：法阵通天，咒文如潮，挥舞魔法杖，星光自杖尖倾泻而下。衣着：元素法身，光尘化袍。梳造：星冕加身，魔力如渊。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 与元素同源，法身自成; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 龙骑士（`dragon_knight`）

**灵胎初醒 · 契约之卵**
- 龙骑士，灵胎初醒阶段·契约之卵。初始形态：一枚契约之卵，暗金蛋壳刻着誓约纹，龙影与骑士的羁绊在卵中相连。火属性灵光微微环绕。神态：星芒中蜷缩，似梦似醒。动作：悬于法阵微光，静候唤醒。衣着：法袍虚影，星光为衣。梳造：灵光发丝，半透明。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 星芒中蜷缩，似梦似醒; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 骑士候选**
- 龙骑士，凡尘砺心阶段·骑士候选。形象：红铠骑士，身侧有龙影。 核心意象：红铠、骑士长枪、龙之影。神态：好奇试探的灵动。动作：吟唱初学，魔力火花迸溅。衣着：学徒法袍，布料朴素。梳造：发丝扬起，缀着星尘。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 好奇试探的灵动; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 龙骑士**
- 龙骑士，道法初成阶段·龙骑士。形象：红铠骑士，身侧有龙影。 核心意象：红铠、骑士长枪、龙之影。神态：光华渐盛，神采焕发。动作：魔力初现，法阵脚下亮起，举枪跃上龙背，巨龙振翼腾空而起。衣着：法师长袍，绣星纹。梳造：魔法帽/发冠初戴。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 光华渐盛，神采焕发; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 骑士之誓**
- 龙骑士，大劫淬炼阶段·骑士之誓。形象：红铠骑士，身侧有龙影。 核心意象：红铠、骑士长枪、龙之影。神态：暗夜中坚持，眼神倔强。动作：全力施法，魔力风暴。衣着：法袍破损，魔力燃烧。梳造：发丝凌乱，眸中带焰。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 暗夜中坚持，眼神倔强; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 天空骑士**
- 龙骑士，封神登天阶段·天空骑士。形象：红铠骑士，身侧有龙影。 核心意象：红铠、骑士长枪、龙之影。神态：受封贤者，光晕庄严。动作：魔法大成，法阵漫天。衣着：贤者白袍，权杖受命。梳造：法冠威仪，光晕庄严。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封贤者，光晕庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 龙骑尊者**
- 龙骑士，道果圆满阶段·龙骑尊者。形象：红铠骑士，身侧有龙影。 核心意象：红铠、骑士长枪、龙之影。神态：与元素同源，法身自成。动作：法阵通天，咒文如潮，举枪跃上龙背，巨龙振翼腾空而起。衣着：元素法身，光尘化袍。梳造：星冕加身，魔力如渊。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 与元素同源，法身自成; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 炼金魔像（`alchemy_golem`）

**灵胎初醒 · 贤者之石**
- 炼金魔像，灵胎初醒阶段·贤者之石。初始形态：一块贤者之石，鎏金棱面内蕴符文，炼金之火在其中凝成一粒跳动的心。金属性灵光微微环绕。神态：沉睡的灵光，物灵未醒。动作：静置无声，灵光内蕴。衣着：素坯/原石，未成形。梳造：无，器物胚形。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 沉睡的灵光，物灵未醒; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 炼金仆从**
- 炼金魔像，凡尘砺心阶段·炼金仆从。形象：金甲魔像，符文刻身。 核心意象：金甲、符文刻痕、炼金秘术。神态：初生灵智的好奇。动作：微光颤动，器灵初醒。衣着：初雕成型，轮廓渐显。梳造：雕纹/铭文初刻。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 初生灵智的好奇; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 炼金术士**
- 炼金魔像，道法初成阶段·炼金术士。形象：金甲魔像，符文刻身。 核心意象：金甲、符文刻痕、炼金秘术。神态：灵光渐盛，灵动自生。动作：器光初现，灵光流转，符文亮起，双拳砸落震裂大地。衣着：成形精工，纹饰渐繁。梳造：纹饰/嵌饰增辉。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 灵光渐盛，灵动自生; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 点金手**
- 炼金魔像，大劫淬炼阶段·点金手。形象：金甲魔像，符文刻身。 核心意象：金甲、符文刻痕、炼金秘术。神态：淬炼中的忍耐。动作：全力受炼，火炼淬洗。衣着：历经淬炼，温润内敛。梳造：包浆/裂纹，岁月痕。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 淬炼中的忍耐; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 黄金石像**
- 炼金魔像，封神登天阶段·黄金石像。形象：金甲魔像，符文刻身。 核心意象：金甲、符文刻痕、炼金秘术。神态：受封灵宝，灵光大成的庄严。动作：器灵大成，光华夺目。衣着：华彩流光，灵气充盈。梳造：祥纹满饰，宝光外放。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封灵宝，灵光大成的庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 炼金尊者**
- 炼金魔像，道果圆满阶段·炼金尊者。形象：金甲魔像，符文刻身。 核心意象：金甲、符文刻痕、炼金秘术。神态：器灵化神，灵性通明。动作：器道通天，光华耀世，符文亮起，双拳砸落震裂大地。衣着：神纹流转，镇世之宝。梳造：灵光化形，器魂永驻。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 器灵化神，灵性通明; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 梦魇夜马（`nightmare_horse`）

**灵胎初醒 · 月影卵**
- 梦魇夜马，灵胎初醒阶段·月影卵。初始形态：一枚月影卵，紫黑蛋壳晕着月辉，梦魇蹄影在夜色灵光中踏空。暗属性灵光微微环绕。神态：灵光中沉睡，兽性未醒的宁静。动作：蜷于灵光，微息起伏。衣着：灵光虚影，形尚未凝。梳造：灵毫光点，未成形相。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 灵光中沉睡，兽性未醒的宁静; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 影鬃马**
- 梦魇夜马，凡尘砺心阶段·影鬃马。形象：紫焰骏马，四蹄踏火。 核心意象：紫焰、踏火四蹄、梦魇之力。神态：初踏山川，懵懂而灵慧。动作：蹒跚踏云，好奇嗅闻。衣着：半实灵体，羽毛/鳞纹初显。梳造：幼羽/灵尾，泛着微光。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 初踏山川，懵懂而灵慧; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 梦魇马**
- 梦魇夜马，道法初成阶段·梦魇马。形象：紫焰骏马，四蹄踏火。 核心意象：紫焰、踏火四蹄、梦魇之力。神态：眼神古老而专注，神通初显。动作：展翅昂首，神姿初现，四蹄踏火奔腾，紫焰在身后拖成长练。衣着：华羽灵纹，瑞气氤氲。梳造：长羽/灵角渐生。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 眼神古老而专注，神通初显; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 月光马**
- 梦魇夜马，大劫淬炼阶段·月光马。形象：紫焰骏马，四蹄踏火。 核心意象：紫焰、踏火四蹄、梦魇之力。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 眸光如电，威严中带着坚韧; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 夜之守护**
- 梦魇夜马，封神登天阶段·夜之守护。形象：紫焰骏马，四蹄踏火。 核心意象：紫焰、踏火四蹄、梦魇之力。神态：受封瑞兽，祥云拱卫。动作：绝技大成，百瑞来朝。衣着：五色神纹冠冕，祥光加身。梳造：圣羽垂天，瑞角冲霄。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封瑞兽，祥云拱卫; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 夜月尊者**
- 梦魇夜马，道果圆满阶段·夜月尊者。形象：紫焰骏马，四蹄踏火。 核心意象：紫焰、踏火四蹄、梦魇之力。神态：瑞气化道，福泽天地。动作：真身镇世，祥光化雨，四蹄踏火奔腾，紫焰在身后拖成长练。衣着：祥光铸身，福泽苍生。梳造：瑞光冠冕，天地共仰。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 瑞气化道，福泽天地; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 元素灯灵（`lamp_spirit`）

**灵胎初醒 · 灯芯种子**
- 元素灯灵，灵胎初醒阶段·灯芯种子。初始形态：一粒灯芯种子，灯油灵光托着一朵火苗，千愿的浮影在火中明灭。火属性灵光微微环绕。神态：沉睡的灵光，物灵未醒。动作：静置无声，灵光内蕴。衣着：素坯/原石，未成形。梳造：无，器物胚形。意境：魔法光球或种子，奥秘初现，星芒闪烁的神秘新生。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：奥秘紫（#7B68EE）主调 + 魔法薰衣草（#E6E6FA）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; arcane genesis; 沉睡的灵光，物灵未醒; palette #7B68EE with #E6E6FA accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小灯灵**
- 元素灯灵，凡尘砺心阶段·小灯灵。形象：金灯灵体，灯火摇曳。 核心意象：神灯、摇曳灯芯、许愿之光。神态：初生灵智的好奇。动作：微光颤动，器灵初醒。衣着：初雕成型，轮廓渐显。梳造：雕纹/铭文初刻。意境：初学魔法，好奇探索，小小的魔力火花。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：魔紫（#9370DB）主调 + 星光蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; curious novice; 初生灵智的好奇; palette #9370DB with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 灯之精灵**
- 元素灯灵，道法初成阶段·灯之精灵。形象：金灯灵体，灯火摇曳。 核心意象：神灯、摇曳灯芯、许愿之光。神态：灵光渐盛，灵动自生。动作：器光初现，灵光流转，灯芯一燃，应召而出，金光漫室。衣着：成形精工，纹饰渐繁。梳造：纹饰/嵌饰增辉。意境：魔力觉醒，脚下浮现法阵，魔法光辉绽放。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：深魔紫（#4B0082）主调 + 咒文金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; magic circle awakening; 灵光渐盛，灵动自生; palette #4B0082 with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 流光灯灵**
- 元素灯灵，大劫淬炼阶段·流光灯灵。形象：金灯灵体，灯火摇曳。 核心意象：神灯、摇曳灯芯、许愿之光。神态：淬炼中的忍耐。动作：全力受炼，火炼淬洗。衣着：历经淬炼，温润内敛。梳造：包浆/裂纹，岁月痕。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 淬炼中的忍耐; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 千愿灯灵**
- 元素灯灵，封神登天阶段·千愿灯灵。形象：金灯灵体，灯火摇曳。 核心意象：神灯、摇曳灯芯、许愿之光。神态：受封灵宝，灵光大成的庄严。动作：器灵大成，光华夺目。衣着：华彩流光，灵气充盈。梳造：祥纹满饰，宝光外放。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 受封灵宝，灵光大成的庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 灯灵尊者**
- 元素灯灵，道果圆满阶段·灯灵尊者。形象：金灯灵体，灯火摇曳。 核心意象：神灯、摇曳灯芯、许愿之光。神态：器灵化神，灵性通明。动作：器道通天，光华耀世，灯芯一燃，应召而出，金光漫室。衣着：神纹流转，镇世之宝。梳造：灵光化形，器魂永驻。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 器灵化神，灵性通明; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 7. 史前生物（12 物种）

> **风格**：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。**阶段演绎**：
> - 灵胎初醒：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息（化石米/石灰）
> - 凡尘砺心：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存（泥土棕/苔藓绿）
> - 道法初成：成长中的捕猎或迁徙，展露力量与速度（岩金/沙棕）
> - 大劫淬炼：严酷的冰河或火山环境下生存试炼，与天争命（火山黑/熔岩红）
> - 封神登天：成为领地霸主，王者的姿态与威严（霸主棕/王金）
> - 道果圆满：远古霸主图腾显圣，冰川之巅威压万古（冰川白/远古蓝）

#### 霸王龙（`t_rex`）

**灵胎初醒 · 化石蛋**
- 霸王龙，灵胎初醒阶段·化石蛋。初始形态：一枚化石蛋，斑驳壳面泛着远古纹路，内部透出熔岩般的赤光，像沉睡的史前之心。火属性灵光微微环绕。神态：化石般的沉睡，原始生命力潜伏。动作：静卧蛋中，尾/爪微动。衣着：蛋壳斑驳，远古纹路。梳造：无，蛋中初形。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 化石般的沉睡，原始生命力潜伏; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年暴龙**
- 霸王龙，凡尘砺心阶段·幼年暴龙。形象：直立霸王龙，头大颚壮，前肢短小。 核心意象：巨颚利齿、粗壮后肢、白垩纪。神态：初生的笨拙好奇。动作：破壳蹒跚，懵懂张望。衣着：稚嫩皮甲，柔软未坚。梳造：头冠/角芽未显。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 初生的笨拙好奇; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 追赶者**
- 霸王龙，道法初成阶段·追赶者。形象：直立霸王龙，头大颚壮，前肢短小。 核心意象：巨颚利齿、粗壮后肢、白垩纪。神态：捕猎时的专注凶狠。动作：低伏潜行，初试锋芒，巨颚一口咬碎骨骼，怒吼震彻山谷。衣着：鳞甲渐厚，色泽加深。梳造：头冠/帆脊渐起。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 捕猎时的专注凶狠; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 领地之王**
- 霸王龙，大劫淬炼阶段·领地之王。形象：直立霸王龙，头大颚壮，前肢短小。 核心意象：巨颚利齿、粗壮后肢、白垩纪。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 白垩纪领主**
- 霸王龙，封神登天阶段·白垩纪领主。形象：直立霸王龙，头大颚壮，前肢短小。 核心意象：巨颚利齿、粗壮后肢、白垩纪。神态：受封霸主，眼神睥睨。动作：猎技大成，万兽辟易。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封霸主，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 终极掠食者**
- 霸王龙，道果圆满阶段·终极掠食者。形象：直立霸王龙，头大颚壮，前肢短小。 核心意象：巨颚利齿、粗壮后肢、白垩纪。神态：洪荒之巅，唯我独尊。动作：踏碎山河，万兽臣服，巨颚一口咬碎骨骼，怒吼震彻山谷。衣着：远古神躯，鳞甲映日。梳造：骨冠擎天，威压万古。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 洪荒之巅，唯我独尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 三角龙（`triceratops`）

**灵胎初醒 · 角蛋**
- 三角龙，灵胎初醒阶段·角蛋。初始形态：一枚角蛋，沙棕蛋壳顶生三处角突，颈盾的雏形在蛋壳边缘隆起。土属性灵光微微环绕。神态：化石般的沉睡，原始生命力潜伏。动作：静卧蛋中，尾/爪微动。衣着：蛋壳斑驳，远古纹路。梳造：无，蛋中初形。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 化石般的沉睡，原始生命力潜伏; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 头盾初现**
- 三角龙，凡尘砺心阶段·头盾初现。形象：大颈盾配三尖角，身强体壮。 核心意象：三尖角、大颈盾、群居护群。神态：初生的笨拙好奇。动作：破壳蹒跚，懵懂张望。衣着：稚嫩皮甲，柔软未坚。梳造：头冠/角芽未显。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 初生的笨拙好奇; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 盾甲龙**
- 三角龙，道法初成阶段·盾甲龙。形象：大颈盾配三尖角，身强体壮。 核心意象：三尖角、大颈盾、群居护群。神态：捕猎时的专注凶狠。动作：低伏潜行，初试锋芒，低头亮出三尖角，怒吼着顶向来敌。衣着：鳞甲渐厚，色泽加深。梳造：头冠/帆脊渐起。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 捕猎时的专注凶狠; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 铁头龙**
- 三角龙，大劫淬炼阶段·铁头龙。形象：大颈盾配三尖角，身强体壮。 核心意象：三尖角、大颈盾、群居护群。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 群首领**
- 三角龙，封神登天阶段·群首领。形象：大颈盾配三尖角，身强体壮。 核心意象：三尖角、大颈盾、群居护群。神态：受封霸主，眼神睥睨。动作：猎技大成，万兽辟易。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封霸主，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 化石之王**
- 三角龙，道果圆满阶段·化石之王。形象：大颈盾配三尖角，身强体壮。 核心意象：三尖角、大颈盾、群居护群。神态：洪荒之巅，唯我独尊。动作：踏碎山河，万兽臣服，低头亮出三尖角，怒吼着顶向来敌。衣着：远古神躯，鳞甲映日。梳造：骨冠擎天，威压万古。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 洪荒之巅，唯我独尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 翼龙（`pterosaur`）

**灵胎初醒 · 飞卵**
- 翼龙，灵胎初醒阶段·飞卵。初始形态：一枚飞卵，灰白蛋壳生着翼膜纹，气流在蛋身周围盘旋，似欲乘风而起。风属性灵光微微环绕。神态：化石般的沉睡，原始生命力潜伏。动作：静卧蛋中，尾/爪微动。衣着：蛋壳斑驳，远古纹路。梳造：无，蛋中初形。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 化石般的沉睡，原始生命力潜伏; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 翼手龙**
- 翼龙，凡尘砺心阶段·翼手龙。形象：展翼无齿，颈长身轻。 核心意象：翼膜、长颈、远古天空。神态：初生的笨拙好奇。动作：破壳蹒跚，懵懂张望。衣着：稚嫩皮甲，柔软未坚。梳造：头冠/角芽未显。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 初生的笨拙好奇; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 风神翼龙**
- 翼龙，道法初成阶段·风神翼龙。形象：展翼无齿，颈长身轻。 核心意象：翼膜、长颈、远古天空。神态：捕猎时的专注凶狠。动作：低伏潜行，初试锋芒，翼膜一展，借风滑翔，俯冲掠水捕鱼。衣着：鳞甲渐厚，色泽加深。梳造：头冠/帆脊渐起。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 捕猎时的专注凶狠; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 火焰翼龙**
- 翼龙，大劫淬炼阶段·火焰翼龙。形象：展翼无齿，颈长身轻。 核心意象：翼膜、长颈、远古天空。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 王者翼龙**
- 翼龙，封神登天阶段·王者翼龙。形象：展翼无齿，颈长身轻。 核心意象：翼膜、长颈、远古天空。神态：受封霸主，眼神睥睨。动作：猎技大成，万兽辟易。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封霸主，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 太古翼龙**
- 翼龙，道果圆满阶段·太古翼龙。形象：展翼无齿，颈长身轻。 核心意象：翼膜、长颈、远古天空。神态：洪荒之巅，唯我独尊。动作：踏碎山河，万兽臣服，翼膜一展，借风滑翔，俯冲掠水捕鱼。衣着：远古神躯，鳞甲映日。梳造：骨冠擎天，威压万古。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 洪荒之巅，唯我独尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 猛犸象（`mammoth`）

**灵胎初醒 · 长毛灵胎**
- 猛犸象，灵胎初醒阶段·长毛灵胎。初始形态：一团长毛灵胎，厚实绒毛覆着霜晶，弯长象牙的一点雏形在寒光中显现。冰属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 象牙变长**
- 猛犸象，凡尘砺心阶段·象牙变长。形象：长毛巨象，弯长象牙。 核心意象：弯长象牙、厚毛、冰河雪原。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冰原猛犸**
- 猛犸象，道法初成阶段·冰原猛犸。形象：长毛巨象，弯长象牙。 核心意象：弯长象牙、厚毛、冰河雪原。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，长鼻卷起枯枝，弯长象牙挑开积雪。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 眸光晶亮，意气初显，跃跃欲试; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 象牙王**
- 猛犸象，大劫淬炼阶段·象牙王。形象：长毛巨象，弯长象牙。 核心意象：弯长象牙、厚毛、冰河雪原。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 冰川之主**
- 猛犸象，封神登天阶段·冰川之主。形象：长毛巨象，弯长象牙。 核心意象：弯长象牙、厚毛、冰河雪原。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封万兽之王，目光睥睨天地; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 远古巨兽**
- 猛犸象，道果圆满阶段·远古巨兽。形象：长毛巨象，弯长象牙。 核心意象：弯长象牙、厚毛、冰河雪原。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，长鼻卷起枯枝，弯长象牙挑开积雪。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 神光自照，与天地同尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 剑齿虎（`sabertooth`）

**灵胎初醒 · 剑齿灵胎**
- 剑齿虎，灵胎初醒阶段·剑齿灵胎。初始形态：一团剑齿灵胎，肌肉初成、绒毛未褪，上颚獠牙的锋利雏形若隐若现。金属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 长牙初成**
- 剑齿虎，凡尘砺心阶段·长牙初成。形象：上颚獠牙如剑，肌肉虬结。 核心意象：剑形獠牙、虬结肌肉、一击必杀。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 刃牙虎**
- 剑齿虎，道法初成阶段·刃牙虎。形象：上颚獠牙如剑，肌肉虬结。 核心意象：剑形獠牙、虬结肌肉、一击必杀。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，低伏潜行，骤然跃起一剑封喉。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 眸光晶亮，意气初显，跃跃欲试; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 冰河猎手**
- 剑齿虎，大劫淬炼阶段·冰河猎手。形象：上颚獠牙如剑，肌肉虬结。 核心意象：剑形獠牙、虬结肌肉、一击必杀。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金剑齿**
- 剑齿虎，封神登天阶段·金剑齿。形象：上颚獠牙如剑，肌肉虬结。 核心意象：剑形獠牙、虬结肌肉、一击必杀。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封万兽之王，目光睥睨天地; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 太古剑齿**
- 剑齿虎，道果圆满阶段·太古剑齿。形象：上颚獠牙如剑，肌肉虬结。 核心意象：剑形獠牙、虬结肌肉、一击必杀。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，低伏潜行，骤然跃起一剑封喉。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 神光自照，与天地同尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 沧龙（`mosasaur`）

**灵胎初醒 · 水卵**
- 沧龙，灵胎初醒阶段·水卵。初始形态：一枚水卵，深蓝蛋壳泛着鳞光，鳍肢的剪影在蛋中游弋，苍茫海意在壳内涌动。水属性灵光微微环绕。神态：龙息轻吐，沉睡于混沌灵光。动作：盘蜷于光，沉眠未醒。衣着：幼嫩鳞胚，泛着初光。梳造：无角，须影初现。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 龙息轻吐，沉睡于混沌灵光; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 鳞甲变厚**
- 沧龙，凡尘砺心阶段·鳞甲变厚。形象：深海巨蜥龙，鳍状四肢。 核心意象：鳍状四肢、巨颚、苍茫汪洋。神态：竖瞳初睁，窥探云水。动作：初探云水，游弋学步。衣着：鳞甲渐密，颜色初显。梳造：角芽微露，须丝轻扬。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 竖瞳初睁，窥探云水; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 海中霸主**
- 沧龙，道法初成阶段·海中霸主。形象：深海巨蜥龙，鳍状四肢。 核心意象：鳍状四肢、巨颚、苍茫汪洋。神态：眸光锐亮，潜龙欲腾。动作：腾空而起，初显神威，鳍肢一摆，巨颚咬碎深海猎物。衣着：鳞甲生辉，腹光流转。梳造：双角初成，须如流云。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 眸光锐亮，潜龙欲腾; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 王沧**
- 沧龙，大劫淬炼阶段·王沧。形象：深海巨蜥龙，鳍状四肢。 核心意象：鳍状四肢、巨颚、苍茫汪洋。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 龙威炽烈，怒目电光; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 古沧**
- 沧龙，封神登天阶段·古沧。形象：深海巨蜥龙，鳍状四肢。 核心意象：鳍状四肢、巨颚、苍茫汪洋。神态：受封龙君，神威赫赫。动作：登天行云，布雨泽四方。衣着：金鳞覆身，受冕祥光。梳造：龙角如珊瑚，加冕为尊。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封龙君，神威赫赫; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 神沧**
- 沧龙，道果圆满阶段·神沧。形象：深海巨蜥龙，鳍状四肢。 核心意象：鳍状四肢、巨颚、苍茫汪洋。神态：万龙之源，睥睨三界。动作：真身化岳，日月为伴，鳍肢一摆，巨颚咬碎深海猎物。衣着：龙身映日月，鳞甲如星辰。梳造：龙角擎天，道纹绕体。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 万龙之源，睥睨三界; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 棘龙（`spinosaurus`）

**灵胎初醒 · 帆纹蛋**
- 棘龙，灵胎初醒阶段·帆纹蛋。初始形态：一枚帆纹蛋，河畔泥色蛋壳隆起一道帆脊纹，河水的湿润气息包裹着它。水属性灵光微微环绕。神态：化石般的沉睡，原始生命力潜伏。动作：静卧蛋中，尾/爪微动。衣着：蛋壳斑驳，远古纹路。梳造：无，蛋中初形。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 化石般的沉睡，原始生命力潜伏; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 帆背龙**
- 棘龙，凡尘砺心阶段·帆背龙。形象：背生巨大帆状脊，形似鳄龙。 核心意象：巨大帆脊、长吻利齿、河畔。神态：初生的笨拙好奇。动作：破壳蹒跚，懵懂张望。衣着：稚嫩皮甲，柔软未坚。梳造：头冠/角芽未显。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 初生的笨拙好奇; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 巨帆龙**
- 棘龙，道法初成阶段·巨帆龙。形象：背生巨大帆状脊，形似鳄龙。 核心意象：巨大帆脊、长吻利齿、河畔。神态：捕猎时的专注凶狠。动作：低伏潜行，初试锋芒，帆脊竖起示威，长吻探入水中叼起巨鱼。衣着：鳞甲渐厚，色泽加深。梳造：头冠/帆脊渐起。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 捕猎时的专注凶狠; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 河畔霸主**
- 棘龙，大劫淬炼阶段·河畔霸主。形象：背生巨大帆状脊，形似鳄龙。 核心意象：巨大帆脊、长吻利齿、河畔。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 河畔之王**
- 棘龙，封神登天阶段·河畔之王。形象：背生巨大帆状脊，形似鳄龙。 核心意象：巨大帆脊、长吻利齿、河畔。神态：受封霸主，眼神睥睨。动作：猎技大成，万兽辟易。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封霸主，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 棘龙尊者**
- 棘龙，道果圆满阶段·棘龙尊者。形象：背生巨大帆状脊，形似鳄龙。 核心意象：巨大帆脊、长吻利齿、河畔。神态：洪荒之巅，唯我独尊。动作：踏碎山河，万兽臣服，帆脊竖起示威，长吻探入水中叼起巨鱼。衣着：远古神躯，鳞甲映日。梳造：骨冠擎天，威压万古。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 洪荒之巅，唯我独尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 甲龙（`ankylosaurus`）

**灵胎初醒 · 锤尾蛋**
- 甲龙，灵胎初醒阶段·锤尾蛋。初始形态：一枚锤尾蛋，岩甲色蛋壳覆着板甲纹，尾端锤形的凸起在壳下隐约。土属性灵光微微环绕。神态：化石般的沉睡，原始生命力潜伏。动作：静卧蛋中，尾/爪微动。衣着：蛋壳斑驳，远古纹路。梳造：无，蛋中初形。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 化石般的沉睡，原始生命力潜伏; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 板甲龙**
- 甲龙，凡尘砺心阶段·板甲龙。形象：身披重甲，尾端生骨锤。 核心意象：重甲骨板、尾端骨锤、不动如山。神态：初生的笨拙好奇。动作：破壳蹒跚，懵懂张望。衣着：稚嫩皮甲，柔软未坚。梳造：头冠/角芽未显。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 初生的笨拙好奇; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 锤尾龙**
- 甲龙，道法初成阶段·锤尾龙。形象：身披重甲，尾端生骨锤。 核心意象：重甲骨板、尾端骨锤、不动如山。神态：捕猎时的专注凶狠。动作：低伏潜行，初试锋芒，骨锤横扫，一击足以震退掠食者。衣着：鳞甲渐厚，色泽加深。梳造：头冠/帆脊渐起。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 捕猎时的专注凶狠; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 堡垒龙**
- 甲龙，大劫淬炼阶段·堡垒龙。形象：身披重甲，尾端生骨锤。 核心意象：重甲骨板、尾端骨锤、不动如山。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 要塞甲龙**
- 甲龙，封神登天阶段·要塞甲龙。形象：身披重甲，尾端生骨锤。 核心意象：重甲骨板、尾端骨锤、不动如山。神态：受封霸主，眼神睥睨。动作：猎技大成，万兽辟易。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封霸主，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 甲龙尊者**
- 甲龙，道果圆满阶段·甲龙尊者。形象：身披重甲，尾端生骨锤。 核心意象：重甲骨板、尾端骨锤、不动如山。神态：洪荒之巅，唯我独尊。动作：踏碎山河，万兽臣服，骨锤横扫，一击足以震退掠食者。衣着：远古神躯，鳞甲映日。梳造：骨冠擎天，威压万古。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 洪荒之巅，唯我独尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 梁龙（`diplodocus`）

**灵胎初醒 · 长颈蛋**
- 梁龙，灵胎初醒阶段·长颈蛋。初始形态：一枚长颈蛋，青灰蛋壳有着细长的颈纹，温柔巨兽的轮廓在壳中延伸。木属性灵光微微环绕。神态：化石般的沉睡，原始生命力潜伏。动作：静卧蛋中，尾/爪微动。衣着：蛋壳斑驳，远古纹路。梳造：无，蛋中初形。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 化石般的沉睡，原始生命力潜伏; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 长颈龙**
- 梁龙，凡尘砺心阶段·长颈龙。形象：长颈细尾的巨型蜥脚龙。 核心意象：超长脖颈、鞭状长尾、成群迁徙。神态：初生的笨拙好奇。动作：破壳蹒跚，懵懂张望。衣着：稚嫩皮甲，柔软未坚。梳造：头冠/角芽未显。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 初生的笨拙好奇; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 巨梁龙**
- 梁龙，道法初成阶段·巨梁龙。形象：长颈细尾的巨型蜥脚龙。 核心意象：超长脖颈、鞭状长尾、成群迁徙。神态：捕猎时的专注凶狠。动作：低伏潜行，初试锋芒，长颈探向高树，细尾如鞭甩动护身。衣着：鳞甲渐厚，色泽加深。梳造：头冠/帆脊渐起。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 捕猎时的专注凶狠; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 群居龙**
- 梁龙，大劫淬炼阶段·群居龙。形象：长颈细尾的巨型蜥脚龙。 核心意象：超长脖颈、鞭状长尾、成群迁徙。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 温柔巨兽**
- 梁龙，封神登天阶段·温柔巨兽。形象：长颈细尾的巨型蜥脚龙。 核心意象：超长脖颈、鞭状长尾、成群迁徙。神态：受封霸主，眼神睥睨。动作：猎技大成，万兽辟易。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封霸主，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 梁龙尊者**
- 梁龙，道果圆满阶段·梁龙尊者。形象：长颈细尾的巨型蜥脚龙。 核心意象：超长脖颈、鞭状长尾、成群迁徙。神态：洪荒之巅，唯我独尊。动作：踏碎山河，万兽臣服，长颈探向高树，细尾如鞭甩动护身。衣着：远古神躯，鳞甲映日。梳造：骨冠擎天，威压万古。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 洪荒之巅，唯我独尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 巨齿鲨（`megalodon`）

**灵胎初醒 · 巨牙卵**
- 巨齿鲨，灵胎初醒阶段·巨牙卵。初始形态：一枚巨牙卵，深灰卵壳生着锯齿纹，巨齿的锋芒在深蓝灵光中初现。水属性灵光微微环绕。神态：化石般的沉睡，原始生命力潜伏。动作：静卧蛋中，尾/爪微动。衣着：蛋壳斑驳，远古纹路。梳造：无，蛋中初形。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 化石般的沉睡，原始生命力潜伏; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 利齿鲨**
- 巨齿鲨，凡尘砺心阶段·利齿鲨。形象：巨牙利齿的远古巨鲨。 核心意象：巨牙、庞大身影、深海之渊。神态：初生的笨拙好奇。动作：破壳蹒跚，懵懂张望。衣着：稚嫩皮甲，柔软未坚。梳造：头冠/角芽未显。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 初生的笨拙好奇; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 深海鲨**
- 巨齿鲨，道法初成阶段·深海鲨。形象：巨牙利齿的远古巨鲨。 核心意象：巨牙、庞大身影、深海之渊。神态：捕猎时的专注凶狠。动作：低伏潜行，初试锋芒，巨口张开，一口吞下整条海兽。衣着：鳞甲渐厚，色泽加深。梳造：头冠/帆脊渐起。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 捕猎时的专注凶狠; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 深蓝鲨**
- 巨齿鲨，大劫淬炼阶段·深蓝鲨。形象：巨牙利齿的远古巨鲨。 核心意象：巨牙、庞大身影、深海之渊。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 巨齿霸主**
- 巨齿鲨，封神登天阶段·巨齿霸主。形象：巨牙利齿的远古巨鲨。 核心意象：巨牙、庞大身影、深海之渊。神态：受封霸主，眼神睥睨。动作：猎技大成，万兽辟易。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封霸主，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 巨齿尊者**
- 巨齿鲨，道果圆满阶段·巨齿尊者。形象：巨牙利齿的远古巨鲨。 核心意象：巨牙、庞大身影、深海之渊。神态：洪荒之巅，唯我独尊。动作：踏碎山河，万兽臣服，巨口张开，一口吞下整条海兽。衣着：远古神躯，鳞甲映日。梳造：骨冠擎天，威压万古。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 洪荒之巅，唯我独尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 大地懒（`ground_sloth`）

**灵胎初醒 · 巨爪灵胎**
- 大地懒，灵胎初醒阶段·巨爪灵胎。初始形态：一团巨爪灵胎，长毛初生的笨拙轮廓，巨型爪钩的雏形蜷在胸前。土属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 巨爪兽**
- 大地懒，凡尘砺心阶段·巨爪兽。形象：巨爪大懒兽，毛长体壮。 核心意象：巨型爪钩、长毛、缓慢而强大。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 大地懒**
- 大地懒，道法初成阶段·大地懒。形象：巨爪大懒兽，毛长体壮。 核心意象：巨型爪钩、长毛、缓慢而强大。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，巨爪刨开泥土，缓缓挪动庞然之躯。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 眸光晶亮，意气初显，跃跃欲试; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 耐心巨兽**
- 大地懒，大劫淬炼阶段·耐心巨兽。形象：巨爪大懒兽，毛长体壮。 核心意象：巨型爪钩、长毛、缓慢而强大。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大地巨人**
- 大地懒，封神登天阶段·大地巨人。形象：巨爪大懒兽，毛长体壮。 核心意象：巨型爪钩、长毛、缓慢而强大。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封万兽之王，目光睥睨天地; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 大地尊者**
- 大地懒，道果圆满阶段·大地尊者。形象：巨爪大懒兽，毛长体壮。 核心意象：巨型爪钩、长毛、缓慢而强大。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，巨爪刨开泥土，缓缓挪动庞然之躯。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 神光自照，与天地同尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 冰河犀牛（`woolly_rhino`）

**灵胎初醒 · 冰角灵胎**
- 冰河犀牛，灵胎初醒阶段·冰角灵胎。初始形态：一团冰角灵胎，厚毛覆着苔原霜气，双角雏形在冰原灵光中缓缓凝成。冰属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：化石在远古晨光中苏醒，蛋壳裂开，远古大地气息。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：化石米（#D2B48C）主调 + 石灰（#A9A9A9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; fossil dawn; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #D2B48C with #A9A9A9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 披毛犀**
- 冰河犀牛，凡尘砺心阶段·披毛犀。形象：长毛犀牛，一双弯角。 核心意象：双弯角、厚长毛、冰原冻土。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：幼兽在远古雨林或冰原蹒跚，跟随兽群学习生存。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：泥土棕（#8B7355）主调 + 苔藓绿（#6B8E7A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; wandering calf; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #8B7355 with #6B8E7A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冰角犀**
- 冰河犀牛，道法初成阶段·冰角犀。形象：长毛犀牛，一双弯角。 核心意象：双弯角、厚长毛、冰原冻土。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，低头拱开积雪，弯角翻出苔草为食。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：成长中的捕猎或迁徙，展露力量与速度。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：岩金（#B8860B）主调 + 沙棕（#CD853F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; growing hunter; 眸光晶亮，意气初显，跃跃欲试; palette #B8860B with #CD853F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 破冰犀**
- 冰河犀牛，大劫淬炼阶段·破冰犀。形象：长毛犀牛，一双弯角。 核心意象：双弯角、厚长毛、冰原冻土。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 冰河之主**
- 冰河犀牛，封神登天阶段·冰河之主。形象：长毛犀牛，一双弯角。 核心意象：双弯角、厚长毛、冰原冻土。神态：受封万兽之王，目光睥睨天地。动作：登顶山巅，受万兽朝拜。衣着：金色王冕祥纹，王者之姿。梳造：王冠受冕，金鬃猎猎。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 受封万兽之王，目光睥睨天地; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 冰河尊者**
- 冰河犀牛，道果圆满阶段·冰河尊者。形象：长毛犀牛，一双弯角。 核心意象：双弯角、厚长毛、冰原冻土。神态：神光自照，与天地同尊。动作：神形合一，啸震九霄，低头拱开积雪，弯角翻出苔草为食。衣着：神光铸体，日月随行。梳造：万灵共仰，神冕无上。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 神光自照，与天地同尊; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 8. 星座守护（12 物种）

> **风格**：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。**阶段演绎**：
> - 灵胎初醒：星屑微光中孕育，小宇宙初醒，命定的星芒（星屑银/白银）
> - 凡尘砺心：圣斗士修行岁月，拳法初练，热血汗水（青铜金/银灰）
> - 道法初成：第七感觉醒，星座之光闪耀，必杀技初现（黄金/星座蓝）
> - 大劫淬炼：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙（暗紫/圣战红）
> - 封神登天：身着黄金圣衣，镇守宫门，星座图腾全开（黄金圣衣/圣光白）
> - 道果圆满：第八感领悟，神之领域，小宇宙化作无尽星光（神光白/小宇宙金）

#### 白羊座·穆（`aries`） · 人生档案版

**灵胎初醒 · 星屑之种**
- 白羊座·穆，灵胎初醒阶段·星屑初醒（幼年，师从史昂·星屑微光）。形象：白羊座·穆，金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。品性：圣域白羊宫幼徒，师从教皇史昂。资质卓绝，小宇宙初醒，如星屑般明亮。。姿态：晨起在星空中冥想，小宇宙微光点点；跟着史昂学习圣衣修理之术。。服饰：青铜学徒服，少年身形，眉目温柔。。体型：身高约5头身，青铜学徒少年，身形清瘦。。衣物细节：青铜学徒服，腰间工具袋。。发型妆造：束发，眉目温柔。。脸型五官：少年面容，眉目温柔，眼含星屑之光，鼻梁挺，嘴角含笑。。武器招式：无兵器，念动力初窥。。功法：小宇宙初醒；念动力初窥。。功法表现：小宇宙星屑微光。。画面：构图：圣域星空下，青铜学徒少年闭目冥想，身周星屑微光点点，身后圣衣架与工具，背景夜穹群星。色调：星屑银+圣域白+夜空蓝。氛围：初醒、温柔、天赋。。台词："师父说，白羊座的心，要比星屑更细。我且慢慢练这双温柔的手。"。动作帧（动图）：①闭目冥想 ②小宇宙星屑浮现 ③修圣衣 ④抬头望星空。诗词：星屑微光照白羊，师从史昂理圣裳。一双妙手千般巧，心比星沙更细长。。主题句：温柔，也是圣域最锋利的守护。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银+圣域白+夜空蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 星屑初醒, age 幼年, 师从史昂·星屑微光; scene: 晨起在星空中冥想，小宇宙微光点点；跟着史昂学习圣衣修理之术。; 星屑微光照白羊，师从史昂理圣裳。一双妙手千般巧，心比星沙更细长。; palette: 星屑银+圣域白+夜空蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 念动初觉**
- 白羊座·穆，凡尘砺心阶段·圣衣修复（少年，圣衣修复天才）。形象：白羊座·穆，金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。品性：圣衣修复的天才，任何破损的圣衣在他手中都能焕然一新。温柔耐心，一丝不苟。。姿态：手持小锤与星光，一点点修复圣衣裂纹；修补时眼中映着星辰。。服饰：青铜学徒服，手持修复工具。。体型：身高约5头身，青铜学徒，俯身专注。。衣物细节：青铜学徒服，工具。。发型妆造：束发。。脸型五官：少年面容，眉目专注，眼含星光。。武器招式：无兵器，修复工具。。功法：圣衣修复术；念动力精细操控。。功法表现：星光缝合，圣衣焕新。。画面：构图：白羊宫工坊，少年俯身修复圣衣，星光化作丝线缝合裂纹，圣衣泛光，背景工具与典籍。色调：星屑银+圣衣金+工坊暖。氛围：匠心、温柔、专注。。台词："圣衣是会呼吸的伙伴。我修的，不是金属，是它们的心。"。动作帧（动图）：①持小锤 ②星光缝合裂纹 ③圣衣泛光 ④端详满意。诗词：星光小锤细细敲，圣衣裂纹补如初。温柔一双天工手，点石成金在此壶。。主题句：温柔，也是圣域最锋利的守护。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银+圣衣金+工坊暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 圣衣修复, age 少年, 圣衣修复天才; scene: 手持小锤与星光，一点点修复圣衣裂纹；修补时眼中映着星辰。; 星光小锤细细敲，圣衣裂纹补如初。温柔一双天工手，点石成金在此壶。; palette: 星屑银+圣衣金+工坊暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 水晶微光**
- 白羊座·穆，道法初成阶段·水晶墙（少年，第七感·星光灭绝初成）。形象：白羊座·穆，金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。品性：第七感觉醒，水晶墙一筑，星光灭绝初成。温柔中渐露锋芒。。姿态：双臂交叉，水晶墙如水幕立起；指尖一点，星光灭绝划破夜空。。服饰：青铜圣衣渐具白羊之形，角如新月。。体型：身高约6头身，青铜圣衣渐具，身形挺拔。。衣物细节：青铜圣衣，羊首角如新月。。发型妆造：束发。。脸型五官：少年面容，眉目坚定，眼含星光。。武器招式：水晶墙/星光灭绝。。功法：水晶墙（防御绝壁）；星光灭绝初成；念动力大成。。功法表现：水晶墙折射星光，灭绝如流星。。画面：构图：白羊宫前，少年双臂交叉，水晶墙如月幕立起，折射星光，指尖星光灭绝如流星，背景夜穹。色调：水晶透蓝+星光银+圣域白。氛围：初成、锋芒、贤者。。台词："师父，我的水晶墙，能为圣域挡住所有的箭了。"。动作帧（动图）：①双臂交叉 ②水晶墙立起 ③指尖星光灭绝 ④星光如流星。诗词：第七感觉水晶墙，星光灭绝指端藏。白羊之角初盈月，温柔渐露锋芒光。。主题句：温柔，也是圣域最锋利的守护。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：水晶透蓝+星光银+圣域白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 水晶墙, age 少年, 第七感·星光灭绝初成; scene: 双臂交叉，水晶墙如水幕立起；指尖一点，星光灭绝划破夜空。; 第七感觉水晶墙，星光灭绝指端藏。白羊之角初盈月，温柔渐露锋芒光。; palette: 水晶透蓝+星光银+圣域白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 水晶墙**
- 白羊座·穆，大劫淬炼阶段·圣战死斗（青年，十二宫之战·守白羊宫）。形象：白羊座·穆，金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。品性：圣战爆发，镇守白羊宫。温柔之下，是圣域最锋利的守护。。姿态：水晶墙连筑，挡下强敌无数；星光灭绝横扫，白羊宫前无人能越。。服饰：白银圣衣（战斗形态），白羊之角凛然。。体型：身高约6头身，白银圣衣，战斗姿态。。衣物细节：白银圣衣，白羊之角。。发型妆造：束发。。脸型五官：青年面容，眉目坚毅，眼含战意又温柔。。武器招式：水晶墙·星光灭绝。。功法：水晶墙·星光灭绝连击；念动力御敌。。功法表现：水晶墙如壁，灭绝如雨。。画面：构图：白羊宫前圣战，白银圣斗士双臂交叉，水晶墙挡住漫天攻击，星光灭绝横扫敌群，背景战火与宫阙。色调：银白+水晶蓝+战火金。氛围：守护、死斗、温柔即锋。。台词："我守护的，不只是白羊宫，是师父托付的整个圣域。"。动作帧（动图）：①双臂交叉 ②水晶墙挡攻击 ③星光灭绝横扫 ④立于宫前。诗词：白羊宫前战鼓催，水晶墙立万箭回。星光灭绝横空去，温柔化作破阵雷。。主题句：温柔，也是圣域最锋利的守护。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：银白+水晶蓝+战火金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 圣战死斗, age 青年, 十二宫之战·守白羊宫; scene: 水晶墙连筑，挡下强敌无数；星光灭绝横扫，白羊宫前无人能越。; 白羊宫前战鼓催，水晶墙立万箭回。星光灭绝横空去，温柔化作破阵雷。; palette: 银白+水晶蓝+战火金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 白羊宫主**
- 白羊座·穆，封神登天阶段·白羊宫主（黄金圣斗士，黄金白羊·圣域贤者）。形象：白羊座·穆，金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。品性：受封黄金圣斗士，白羊宫主。圣域最温柔的贤者，也是不可逾越的守护。。姿态：端坐白羊宫，念动间水晶墙四起；来人过宫，先过水晶墙再说。。服饰：黄金白羊圣衣，角如新月，星光流转。。体型：身高约7头身，黄金白羊圣衣，端坐威仪。。衣物细节：黄金白羊圣衣，角如新月。。发型妆造：束发。。脸型五官：青年面容，眉目温柔含笑，眼含星光。。武器招式：水晶墙/星光灭绝。。功法：黄金白羊之力；水晶墙·星光灭绝大成；圣衣修复宗师。。功法表现：水晶墙如穹顶，星光绕身。。画面：构图：白羊宫大殿，黄金白羊圣斗士端坐，水晶墙如穹顶环绕，星光流转，身后白羊宫星象，背景金色殿堂。色调：黄金+水晶透蓝+星屑银。氛围：黄金、贤者、镇守。。台词："黄金圣衣加身，我还是那个修圣衣的少年。温柔和锋利，本是一体。"。动作帧（动图）：①端坐 ②念动水晶墙 ③星光流转 ④俯瞰白羊宫。诗词：黄金白羊镇圣宫，水晶墙起万法通。温柔原是至坚处，星屑流光耀苍穹。。主题句：温柔，也是圣域最锋利的守护。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+水晶透蓝+星屑银。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 白羊宫主, age 黄金圣斗士, 黄金白羊·圣域贤者; scene: 端坐白羊宫，念动间水晶墙四起；来人过宫，先过水晶墙再说。; 黄金白羊镇圣宫，水晶墙起万法通。温柔原是至坚处，星屑流光耀苍穹。; palette: 黄金+水晶透蓝+星屑银; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 星光灭绝**
- 白羊座·穆，道果圆满阶段·万星朝圣（终极，星光灭绝·万星朝圣）。形象：白羊座·穆，金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。品性：第八感领悟，白羊贤者之魂圆满。星光灭绝化作万星朝圣——温柔了一辈子，守护了圣域一辈子。。姿态：双手一展，万星朝圣之光遍洒圣域；水晶墙化作星河，护住每一个角落。。服饰：神圣衣·白羊（或黄金终极态），星光为翼。。体型：身高约7头身，神圣衣白羊，星光为翼。。衣物细节：神圣衣·白羊，星光。。发型妆造：束发，星光绕。。脸型五官：面容温柔，眉目含光，眼如星海。。武器招式：星光灭绝·万星朝圣。。功法：万星朝圣；第八感；贤者之魂。。功法表现：万星朝圣，星河环绕。。画面：构图：白羊宫顶，白羊圣斗士立于星光之巅，双手一展万星朝圣，星光遍洒圣域，水晶墙化作星河环绕，背景漫天星辰。色调：神光白+万星金+圣域紫。氛围：终极、万星、贤者圆满。。台词："修了一辈子圣衣，也修了一辈子圣域。这万星朝圣的光，是给每一个守护者点的灯。"。动作帧（动图）：①立于星光之巅 ②双手一展 ③万星朝圣 ④星光遍洒。诗词：万星朝圣白羊峰，星光湮灭护圣宫。温柔守尽千般夜，终化星河照苍穹。。主题句：温柔，也是圣域最锋利的守护。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白+万星金+圣域紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 万星朝圣, age 终极, 星光灭绝·万星朝圣; scene: 双手一展，万星朝圣之光遍洒圣域；水晶墙化作星河，护住每一个角落。; 万星朝圣白羊峰，星光湮灭护圣宫。温柔守尽千般夜，终化星河照苍穹。; palette: 神光白+万星金+圣域紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 金牛座·阿鲁迪巴（`taurus`） · 人生档案版

**灵胎初醒 · 星辉之种**
- 金牛座·阿鲁迪巴，灵胎初醒阶段·大地之重（少年，初识大地的重量）。形象：金牛座·阿鲁迪巴，金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。品性：圣域金牛宫少年，天生神力，豪爽憨厚。初识大地的重量，也学会背负它。。姿态：扛起巨石绕着金牛宫走；与同门比力气，赢了就咧嘴憨笑。。服饰：青铜学徒服，体格魁梧。。体型：身高约6头身，魁梧少年，虎背熊腰。。衣物细节：青铜学徒服。。发型妆造：短发，粗犷。。脸型五官：圆脸憨厚，粗眉，大眼豪迈，鼻梁宽，咧嘴笑，下巴方正。。武器招式：无兵器，拳脚神力。。功法：天生神力；基础拳脚。。功法表现：无神力，天生巨力。。画面：构图：金牛宫外旷野，魁梧少年扛巨石而行，憨笑回头，尘土飞扬，背景远山大地。色调：青铜服+大地褐+天光金。氛围：憨厚、豪爽、初识大地。。台词："这块石头，我扛得起。圣域的分量，我也扛得起。"。动作帧（动图）：①扛起巨石 ②绕宫而行 ③与同门比力气 ④赢了咧嘴笑。诗词：金牛少年力拔山，巨石压肩若等闲。豪爽憨笑一声吼，初识大地背脊坚。。主题句：豪迈的巨人，用宽厚的背脊为圣域撑起屏障。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜服+大地褐+天光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 大地之重, age 少年, 初识大地的重量; scene: 扛起巨石绕着金牛宫走；与同门比力气，赢了就咧嘴憨笑。; 金牛少年力拔山，巨石压肩若等闲。豪爽憨笑一声吼，初识大地背脊坚。; palette: 青铜服+大地褐+天光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 蛮力渐长**
- 金牛座·阿鲁迪巴，凡尘砺心阶段·力拔山河（少年，力拔山河·练巨型号角）。形象：金牛座·阿鲁迪巴，金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。品性：日复一日锤炼巨型号角，一拳一吼，力拔山河。。姿态：对着山壁练吼，声震四野；以双臂撼动巨石，一吼山石俱裂。。服饰：青铜训练服，肌肉虬结。。体型：身高约6头身，肌肉虬结，魁梧。。衣物细节：青铜训练服。。发型妆造：短发。。脸型五官：圆脸，粗眉，虎目豪迈，咧嘴。。武器招式：巨型号角初练。。功法：巨型号角初练；神力渐成。。功法表现：拳风震山。。画面：构图：山壁前，魁梧少年一拳轰出，山石崩裂，尘土扬起，双角圣衣初现，背景群山。色调：青铜+山石灰+尘烟金。氛围：苦练、神力、豪迈。。台词："拳头要快，更要重。我这一拳下去，连山都要让路。"。动作帧（动图）：①蓄力 ②一拳轰出 ③山石崩裂 ④仰天大笑。诗词：力拔山河练角鸣，一吼震得山欲倾。拳落石裂尘飞扬，豪迈男儿胆气横。。主题句：豪迈的巨人，用宽厚的背脊为圣域撑起屏障。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+山石灰+尘烟金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 力拔山河, age 少年, 力拔山河·练巨型号角; scene: 对着山壁练吼，声震四野；以双臂撼动巨石，一吼山石俱裂。; 力拔山河练角鸣，一吼震得山欲倾。拳落石裂尘飞扬，豪迈男儿胆气横。; palette: 青铜+山石灰+尘烟金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 角力初成**
- 金牛座·阿鲁迪巴，道法初成阶段·巨型号角（少年，第七感·巨型号角）。形象：金牛座·阿鲁迪巴，金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。品性：第七感觉醒，巨型号角大成，一吼山岳为之震颤。。姿态：双角一抵，巨型号角冲撞，如蛮牛破阵；一吼震彻山谷。。服饰：青铜圣衣具金牛之形，双角如巨柱。。体型：身高约6头身，青铜金牛圣衣，魁梧如塔。。衣物细节：青铜金牛圣衣，双角如柱。。发型妆造：短发。。脸型五官：圆脸，粗眉，虎目如炬，咧嘴。。武器招式：巨型号角（冲撞）。。功法：巨型号角（冲撞绝技）；金牛威压。。功法表现：一角震山，吼断山谷。。画面：构图：山谷之中，青铜金牛圣衣少年低头一角，巨型号角冲击，山岳震颤，碎石四溅，背景群峰。色调：青铜+角金+山灰。氛围：第七感、初成、威震。。台词："角长出来了，我也该替圣域守门了。谁想过宫，先接我一角！"。动作帧（动图）：①低头蓄角 ②巨型号角冲击 ③山岳震颤 ④仰天怒吼。诗词：第七感觉角初横，巨型号角震岳鸣。金牛圣衣初盈月，一声吼断万山兵。。主题句：豪迈的巨人，用宽厚的背脊为圣域撑起屏障。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+角金+山灰。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 巨型号角, age 少年, 第七感·巨型号角; scene: 双角一抵，巨型号角冲撞，如蛮牛破阵；一吼震彻山谷。; 第七感觉角初横，巨型号角震岳鸣。金牛圣衣初盈月，一声吼断万山兵。; palette: 青铜+角金+山灰; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 巨型号角**
- 金牛座·阿鲁迪巴，大劫淬炼阶段·金牛宫战（青年，守金牛宫·力战不悔）。形象：金牛座·阿鲁迪巴，金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。品性：圣战爆发，镇守金牛宫。任谁来犯，他一角一吼，寸步不让。。姿态：金牛宫前，巨型号角连发；受了伤也不退，豪迈大笑"再来！"。服饰：白银金牛圣衣，战痕斑斑。。体型：身高约7头身，白银金牛圣衣，战痕斑斑。。衣物细节：白银金牛圣衣，战痕。。发型妆造：短发。。脸型五官：圆脸，粗眉，虎目豪迈，咧嘴大笑。。武器招式：巨型号角·金牛威压。。功法：巨型号角·金牛威压连击；铁壁身躯。。功法表现：一角撞散敌阵。。画面：构图：金牛宫前圣战，白银金牛圣斗士一角撞散敌阵，身上带伤仍大笑，背景战火宫阙。色调：银白+战火金+金牛角金。氛围：守护、力战、不悔。。台词："想过去？从我身上踩过去！我阿鲁迪巴，从不后退！"。动作帧（动图）：①立于宫前 ②巨型号角连发 ③受伤不退 ④豪迈大笑。诗词：金牛宫前战火燃，一角一吼破千关。受伤犹作豪迈笑，铁壁身躯立不弯。。主题句：豪迈的巨人，用宽厚的背脊为圣域撑起屏障。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：银白+战火金+金牛角金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 金牛宫战, age 青年, 守金牛宫·力战不悔; scene: 金牛宫前，巨型号角连发；受了伤也不退，豪迈大笑"再来！"; 金牛宫前战火燃，一角一吼破千关。受伤犹作豪迈笑，铁壁身躯立不弯。; palette: 银白+战火金+金牛角金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金牛宫主**
- 金牛座·阿鲁迪巴，封神登天阶段·金牛宫主（黄金圣斗士，黄金金牛·圣域屏障）。形象：金牛座·阿鲁迪巴，金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。品性：受封黄金圣斗士，金牛宫主。用宽厚的背脊，为圣域撑起屏障。。姿态：端坐金牛宫，来者先过巨型号角；战后为伤者让出臂弯，豪迈而宽厚。。服饰：黄金金牛圣衣，双角如巨柱，威光赫赫。。体型：身高约7头身，黄金金牛圣衣，魁梧如岳。。衣物细节：黄金金牛圣衣，双角如柱。。发型妆造：短发。。脸型五官：圆脸，粗眉，虎目豪迈，咧嘴。。武器招式：巨型号角（大成）。。功法：黄金金牛之力；巨型号角大成；铁壁身躯。。功法表现：黄金牛角威光。。画面：构图：金牛宫大殿，黄金金牛圣斗士负手而立，魁梧如山，身后金牛宫与星象，背景金色殿堂。色调：黄金+牛角金+殿堂暖。氛围：黄金、镇守、宽厚。。台词："这身黄金圣衣重得很，可我的背，背得起整个圣域。"。动作帧（动图）：①负手而立 ②来者一角 ③战后让臂弯 ④豪迈大笑。诗词：黄金金牛镇圣宫，巨型号角撼苍穹。宽厚背脊撑天起，豪迈一吼万山通。。主题句：豪迈的巨人，用宽厚的背脊为圣域撑起屏障。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+牛角金+殿堂暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 金牛宫主, age 黄金圣斗士, 黄金金牛·圣域屏障; scene: 端坐金牛宫，来者先过巨型号角；战后为伤者让出臂弯，豪迈而宽厚。; 黄金金牛镇圣宫，巨型号角撼苍穹。宽厚背脊撑天起，豪迈一吼万山通。; palette: 黄金+牛角金+殿堂暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 黄金之角**
- 金牛座·阿鲁迪巴，道果圆满阶段·大地共鸣（终极，巨型号角·大地共鸣）。形象：金牛座·阿鲁迪巴，金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。品性：第八感领悟，金牛之魂与大地共鸣。豪迈的巨人，最终与大地一体，永不言败。。姿态：一拳轰向大地，与大地共鸣之力化作巨型号角·大地共鸣；护住身后一切。。服饰：神圣衣·金牛（或黄金终极态），角如擎天柱。。体型：身高约7头身，神圣衣金牛，魁梧如岳。。衣物细节：神圣衣·金牛，角如擎天柱。。发型妆造：短发。。脸型五官：圆脸，粗眉，虎目含光，咧嘴。。武器招式：巨型号角·大地共鸣。。功法：大地共鸣；巨型号角终极；第八感。。功法表现：大地共鸣，金色波纹。。画面：构图：旷野大地，金牛圣斗士一拳轰地，大地共鸣之力化作金色波纹扩散，山岳呼应，身后圣域宫阙，背景天地浩大。色调：神光金+大地褐+圣域暖。氛围：终极、大地、不灭。。台词："大地从来不会败。它载着万物，也背着圣域——我阿鲁迪巴，就是大地的背脊。"。动作帧（动图）：①一拳轰地 ②大地共鸣波纹 ③山岳呼应 ④护住身后。诗词：大地共鸣金牛鸣，一拳撼岳万山倾。豪迈不灭魂犹在，背脊撑起圣域城。。主题句：豪迈的巨人，用宽厚的背脊为圣域撑起屏障。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光金+大地褐+圣域暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 大地共鸣, age 终极, 巨型号角·大地共鸣; scene: 一拳轰向大地，与大地共鸣之力化作巨型号角·大地共鸣；护住身后一切。; 大地共鸣金牛鸣，一拳撼岳万山倾。豪迈不灭魂犹在，背脊撑起圣域城。; palette: 神光金+大地褐+圣域暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 双子座·撒加（`gemini`） · 人生档案版

**灵胎初醒 · 星辉双子**
- 双子座·撒加，灵胎初醒阶段·善恶双生（童年，善恶双生的初醒）。形象：双子座·撒加，黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。品性：圣域双子宫幼徒，天生双面人格——一半是仁善的神之化身，一半是潜伏的恶念。。姿态：白日常行善举，救人济弱；夜深独坐，能听见心底另一道声音低语。。服饰：青铜学徒服，眉目俊美，却总有一瞬眼神幽深。。体型：身高约5头身，青铜学徒少年，俊美。。衣物细节：青铜学徒服。。发型妆造：束发，俊美。。脸型五官：俊美面容，眉目温柔，眼含仁善，鼻梁挺，嘴角温和；眼底偶现暗影。。武器招式：无兵器，小宇宙初醒。。功法：小宇宙初醒；双子宿命初现。。功法表现：善念之光与暗影微澜。。画面：构图：圣域夜空，青铜学徒少年立于宫前，月光下一个影子却悄悄长出另一张面孔，背景双子宫。色调：月光银+阴影紫+圣域白。氛围：双生、初醒、暗影潜伏。。台词："师父说我是神的化身。可为什么，我偶尔能听见心底有个声音在笑？"。动作帧（动图）：①白昼行善 ②夜深独坐 ③听心底低语 ④蹙眉压住暗影。诗词：双生善恶一胎生，白日仁心夜隐声。本是神躯承天命，谁料影里藏狰狞。。主题句：一半是神明，一半是恶鬼，一生都在与自己为敌。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：月光银+阴影紫+圣域白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 善恶双生, age 童年, 善恶双生的初醒; scene: 白日常行善举，救人济弱；夜深独坐，能听见心底另一道声音低语。; 双生善恶一胎生，白日仁心夜隐声。本是神躯承天命，谁料影里藏狰狞。; palette: 月光银+阴影紫+圣域白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 善恶交织**
- 双子座·撒加，凡尘砺心阶段·神之化身（少年，神之化身·光暗拉扯）。形象：双子座·撒加，黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。品性：被誉为"神的化身"，善面光芒万丈，恶面却在暗中疯长。光与暗在体内拉扯。。姿态：小宇宙浩瀚如神，一招惊动圣域；可那夜，恶面第一次夺过他的身体，浅笑。。服饰：青铜圣衣渐具双子之形，黑金交织。。体型：身高约6头身，青铜圣衣渐具，气质如神。。衣物细节：青铜圣衣，黑金交织。。发型妆造：束发，长发。。脸型五官：俊美面容，眉目时而圣洁时而幽深，眼神明暗交替。。武器招式：银河星爆初窥。。功法：银河星爆初窥；小宇宙如神。。功法表现：小宇宙如神，暗影随行。。画面：构图：圣域演武场，少年一招银河星爆惊动四座，然月光下其影渐生暗面，背景群星与双子宫。色调：星光银+暗影紫+圣域金。氛围：光暗拉扯、神之化身。。台词："我越是强大，那道声音就越清晰。它说——你本该是神，也可以是魔。"。动作帧（动图）：①小宇宙爆发 ②银河星爆 ③众目惊叹 ④暗影掠眉。诗词：神之化身耀圣域，银河星爆动天虚。谁料善光愈耀处，恶念随之愈粗。。主题句：一半是神明，一半是恶鬼，一生都在与自己为敌。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星光银+暗影紫+圣域金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 神之化身, age 少年, 神之化身·光暗拉扯; scene: 小宇宙浩瀚如神，一招惊动圣域；可那夜，恶面第一次夺过他的身体，浅笑。; 神之化身耀圣域，银河星爆动天虚。谁料善光愈耀处，恶念随之愈粗。; palette: 星光银+暗影紫+圣域金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 异次元微光**
- 双子座·撒加，道法初成阶段·伪教皇（青年，恶面夺舍·堕落伪教皇）。形象：双子座·撒加，黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。品性：恶面夺舍，善面沉睡。他戴上教皇面具，屠戮忠良，谋夺圣域——堕落为伪教皇。。姿态：暗中弑杀恩师，戴上面具号令圣域；善面偶尔挣扎，旋即被恶面压下。。服饰：教皇袍，面具，双面圣衣黑金狰狞。。体型：身高约7头身，教皇袍，戴面具。。衣物细节：教皇袍，面具，双面圣衣。。发型妆造：束发，面具半掩。。脸型五官：面具之下眼神阴鸷，偶现挣扎。。武器招式：银河星爆（大成）。。功法：银河星爆大成；伪教皇权术。。功法表现：暗影如蝠，恶念滔天。。画面：构图：教皇殿，戴面具的教皇立于宝座前，身后双面圣衣一善一狞，暗影如蝠，背景幽暗殿堂。色调：教皇袍白+暗影黑+面具金。氛围：堕落、恶面、伪教皇。。台词："善的那一半，睡着了。现在执掌这具身体的，是神，也是魔。"。动作帧（动图）：①戴上教皇面具 ②号令圣域 ③善面偶现挣扎 ④恶面压下。诗词：恶面夺舍换圣纲，教皇面具掩锋芒。善心沉睡深宫底，一半神明一半狂。。主题句：一半是神明，一半是恶鬼，一生都在与自己为敌。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：教皇袍白+暗影黑+面具金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 伪教皇, age 青年, 恶面夺舍·堕落伪教皇; scene: 暗中弑杀恩师，戴上面具号令圣域；善面偶尔挣扎，旋即被恶面压下。; 恶面夺舍换圣纲，教皇面具掩锋芒。善心沉睡深宫底，一半神明一半狂。; palette: 教皇袍白+暗影黑+面具金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 异次元空间**
- 双子座·撒加，大劫淬炼阶段·银河星爆（青年，银河星爆·群星破碎）。形象：双子座·撒加，黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。品性：十二宫之战，他以银河星爆迎敌，群星破碎。恶战之中，善面开始苏醒。。姿态：银河星爆一出，群星为之破碎；与来者死斗，生死一线间，善面激烈挣扎。。服饰：白银双子圣衣，黑金交织。。体型：身高约7头身，白银双子圣衣，战姿。。衣物细节：白银双子圣衣，黑金。。发型妆造：束发。。脸型五官：俊美面容，一半狰狞一半悲悯，明暗交替。。武器招式：银河星爆/异次元空间。。功法：银河星爆·群星破碎；异次元空间。。功法表现：银河破碎，善恶光暗交替。。画面：构图：双子宫前，白银双子圣斗士银河星爆击碎群星，光影交错间善恶两相争抢这具身体，背景破碎星穹。色调：银河银+暗影紫+战火红。氛围：星爆、恶战、宿命挣扎。。台词："够了！这具身体，本就不该只有恶魔……我，到底是谁？"。动作帧（动图）：①银河星爆 ②群星破碎 ③善恶争抢 ④跪地挣扎。诗词：银河星爆碎群星，双子宫前战血腥。生死关头善面醒，一半神魔问谁明。。主题句：一半是神明，一半是恶鬼，一生都在与自己为敌。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：银河银+暗影紫+战火红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 银河星爆, age 青年, 银河星爆·群星破碎; scene: 银河星爆一出，群星为之破碎；与来者死斗，生死一线间，善面激烈挣扎。; 银河星爆碎群星，双子宫前战血腥。生死关头善面醒，一半神魔问谁明。; palette: 银河银+暗影紫+战火红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 双子宫主**
- 双子座·撒加，封神登天阶段·回归本心（黄金圣斗士，回归本心·神之化身）。形象：双子座·撒加，黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。品性：善恶合一，善面重掌身体，以生命赎罪。回归本心的神之化身，坦然赴死。。姿态：向同袍托付圣域，卸下伪教皇之罪；最后一战，以善面之身挥出最强星爆。。服饰：黄金双子圣衣，善恶双面归于一体，圣洁与威严并具。。体型：身高约7头身，黄金双子圣衣，气质庄严。。衣物细节：黄金双子圣衣，双面合一。。发型妆造：束发。。脸型五官：俊美面容，眉目释然，眼含悲悯。。武器招式：银河星爆（圆满）。。功法：黄金双子之力；银河星爆圆满；赎罪之心。。功法表现：圣光与暗影交融。。画面：构图：黄金圣殿，黄金双子圣斗士负手而立，身后善恶双面圣衣归于一体，圣光与暗影交融，背景金色殿堂。色调：黄金+圣光白+暗影紫。氛围：赎罪、回归、善面。。台词："我不求原谅。只求这最后一战，让圣域知道——那个善的撒加，一直在。"。动作帧（动图）：①托付圣域 ②卸下罪孽 ③善面挥星爆 ④坦然而立。诗词：回归本心赎罪身，双面归一看分明。黄金圣衣承天命，善面终胜恶魔兵。。主题句：一半是神明，一半是恶鬼，一生都在与自己为敌。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+圣光白+暗影紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 回归本心, age 黄金圣斗士, 回归本心·神之化身; scene: 向同袍托付圣域，卸下伪教皇之罪；最后一战，以善面之身挥出最强星爆。; 回归本心赎罪身，双面归一看分明。黄金圣衣承天命，善面终胜恶魔兵。; palette: 黄金+圣光白+暗影紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 银河星爆**
- 双子座·撒加，道果圆满阶段·善恶归一（终极，银河星爆·双子归元）。形象：双子座·撒加，黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。品性：第八感领悟，善恶归一，双子之神。不再与自己为敌，而是拥抱了两面——善与恶，皆是"我"。。姿态：双手一合，善与恶化作双子归元之力；银河星爆·双子归元，护圣域亦度自身。。服饰：神圣衣·双子（或黄金终极态），黑白双子神相。。体型：身高约7头身，神圣衣双子，黑白双子神相。。衣物细节：神圣衣·双子，黑白。。发型妆造：束发。。脸型五官：俊美面容，眉目平和，善恶双相皆含于一眼。。武器招式：银河星爆·双子归元。。功法：双子归元；第八感；善恶合一。。功法表现：善恶归一，双子之光。。画面：构图：双子星座之下，双子圣斗士双手一合，善恶双相化作黑白双子归元，银河星爆·双子归元笼罩圣域，背景星河浩瀚。色调：神光白+双子黑+银河银。氛围：终极、归一、双子神。。台词："我曾与自己为敌一辈子。如今才懂——那个恶的我，也是我。拥抱了他，才成了神。"。动作帧（动图）：①双手一合 ②善恶双相归元 ③银河星爆·双子归元 ④护圣域。诗词：善恶归一双子神，银河星爆照乾坤。不再与己为敌久，一合双元立圣门。。主题句：一半是神明，一半是恶鬼，一生都在与自己为敌。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白+双子黑+银河银。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 善恶归一, age 终极, 银河星爆·双子归元; scene: 双手一合，善与恶化作双子归元之力；银河星爆·双子归元，护圣域亦度自身。; 善恶归一双子神，银河星爆照乾坤。不再与己为敌久，一合双元立圣门。; palette: 神光白+双子黑+银河银; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 巨蟹座·迪斯马斯克（`cancer`） · 人生档案版

**灵胎初醒 · 冥辉之种**
- 巨蟹座·迪斯马斯克，灵胎初醒阶段·黄泉回响（童年，黄泉比良坡的回声）。形象：巨蟹座·迪斯马斯克，青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。品性：圣域巨蟹宫幼徒，资质奇诡，小宇宙带着一股黄泉比良坡的阴冷气息。。姿态：旁人练拳他练影，指尖一缕阴气探向冥界；笑声里带着旁人听不出的狂妄。。服饰：青铜学徒服，眉目阴柔。。体型：身高约5头身，青铜学徒少年，身形瘦削。。衣物细节：青铜学徒服。。发型妆造：束发，阴柔。。脸型五官：阴柔面容，细眉，狭长眼含戾气，薄唇微勾，下巴尖。。武器招式：无兵器，积尸气初窥。。功法：小宇宙初醒（带冥界气）；积尸气初窥。。功法表现：冥界阴气。。画面：构图：巨蟹宫暗影处，阴柔少年指尖一缕青灰阴气探向冥界，背景比良坡虚影，月光惨淡。色调：冥界灰绿+月光惨白+巨蟹青。氛围：邪气、初醒、狂妄。。台词："你们怕死，我可不怕。黄泉比良坡的那头，我一直听得见。"。动作帧（动图）：①练影 ②指尖阴气探冥界 ③侧耳听 ④邪笑。诗词：黄泉回响巨蟹宫，少年眼底阴气浓。别人练拳我练影，冥界之声入耳中。。主题句：视杀戮为玩笑，却终被自己的邪念反噬。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：冥界灰绿+月光惨白+巨蟹青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 黄泉回响, age 童年, 黄泉比良坡的回声; scene: 旁人练拳他练影，指尖一缕阴气探向冥界；笑声里带着旁人听不出的狂妄。; 黄泉回响巨蟹宫，少年眼底阴气浓。别人练拳我练影，冥界之声入耳中。; palette: 冥界灰绿+月光惨白+巨蟹青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 冥气初生**
- 巨蟹座·迪斯马斯克，凡尘砺心阶段·狂妄渐浓（少年，狂妄笑意·视杀为戏）。形象：巨蟹座·迪斯马斯克，青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。品性：小宇宙渐强，狂妄渐浓。视杀戮为玩笑，以强凌弱为乐，师父的告诫充耳不闻。。姿态：与同门切磋，出手渐重，伤了人还大笑"真有趣"；俯视弱者的眼神如看蝼蚁。。服饰：青铜圣衣具巨蟹之形，甲壳泛冥光。。体型：身高约6头身，青铜巨蟹圣衣，身形阴鸷。。衣物细节：青铜巨蟹圣衣，甲壳冥光。。发型妆造：束发。。脸型五官：阴柔面容，狂笑，眼神轻蔑。。武器招式：积尸气冥界波初成。。功法：积尸气冥界波初成；狂妄之心。。功法表现：冥气渐浓。。画面：构图：演武场，巨蟹少年一掌击飞同门，笑声张狂，身后冥气渐浓，背景圣域宫阙。色调：巨蟹青灰+血战红+冥气灰绿。氛围：狂妄、杀戮、玩笑。。台词："命这种东西，在我这儿就是个笑话。想取，随时可取。"。动作帧（动图）：①一掌击飞 ②张狂大笑 ③俯视弱者 ④冥气渐浓。诗词：狂妄渐浓眉目狂，切磋出手竟伤伤。视人命作玩笑戏，师父叮咛全忘光。。主题句：视杀戮为玩笑，却终被自己的邪念反噬。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：巨蟹青灰+血战红+冥气灰绿。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 狂妄渐浓, age 少年, 狂妄笑意·视杀为戏; scene: 与同门切磋，出手渐重，伤了人还大笑"真有趣"；俯视弱者的眼神如看蝼蚁。; 狂妄渐浓眉目狂，切磋出手竟伤伤。视人命作玩笑戏，师父叮咛全忘光。; palette: 巨蟹青灰+血战红+冥气灰绿; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冥界门扉**
- 巨蟹座·迪斯马斯克，道法初成阶段·冥界之门（少年，积尸气冥界波·直通黄泉）。形象：巨蟹座·迪斯马斯克，青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。品性：积尸气冥界波大成，一手拉开冥界之门，将对手拖入黄泉。引路者之名渐起。。姿态：双掌一推，积尸气冥界波缠住对手，冥界之门在身后轰然洞开。。服饰：白银巨蟹圣衣渐具，甲壳冥光凛冽。。体型：身高约6头身，白银巨蟹圣衣，阴鸷。。衣物细节：白银巨蟹圣衣。。发型妆造：束发。。脸型五官：阴柔面容，眼含幽光，嘴角狞笑。。武器招式：积尸气冥界波/冥界之门。。功法：积尸气冥界波；冥界之门；黄泉引渡。。功法表现：冥界之门洞开，黄泉之气。。画面：构图：巨蟹宫前，白银巨蟹圣斗士双掌推出积尸气，身后冥界之门轰然洞开，黄泉之气涌出，背景幽暗。色调：蟹甲银+冥界灰绿+黄泉幽蓝。氛围：积尸气、冥界、引路。。台词："看见那道门了吗？进去的人，从来没有出来过。我送送你。"。动作帧（动图）：①双掌一推 ②积尸气缠敌 ③冥门洞开 ④送人入黄泉。诗词：积尸气波冥门开，黄泉引渡手中来。蟹甲凛冽幽光泛，笑送众生入夜台。。主题句：视杀戮为玩笑，却终被自己的邪念反噬。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：蟹甲银+冥界灰绿+黄泉幽蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 冥界之门, age 少年, 积尸气冥界波·直通黄泉; scene: 双掌一推，积尸气冥界波缠住对手，冥界之门在身后轰然洞开。; 积尸气波冥门开，黄泉引渡手中来。蟹甲凛冽幽光泛，笑送众生入夜台。; palette: 蟹甲银+冥界灰绿+黄泉幽蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 积尸气冥界波**
- 巨蟹座·迪斯马斯克，大劫淬炼阶段·冥河浸染（青年，被冥河浸染·邪念反噬）。形象：巨蟹座·迪斯马斯克，青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。品性：常年在冥界进进出出，终被冥河浸染。邪念反噬，他开始害怕——那扇门，似乎也要吞噬他自己。。姿态：又一次拉开冥门，却第一次感到门内的目光在看他；夜半惊醒，冷汗涔涔。。服饰：白银巨蟹圣衣，冥气浸入甲壳。。体型：身高约7头身，白银巨蟹圣衣，冥气浸染。。衣物细节：白银巨蟹圣衣，冥气。。发型妆造：束发。。脸型五官：阴柔面容，第一次浮现恐惧。。武器招式：积尸气·黄泉引渡。。功法：积尸气·黄泉引渡；被冥河浸染。。功法表现：冥门窥视，邪念反噬。。画面：构图：冥界入口，白银巨蟹圣斗士立于门边，门内一道目光窥视着他，蟹甲被冥气浸染，他第一次面露惧色，背景幽暗冥河。色调：冥河暗绿+蟹甲灰+恐惧黑。氛围：被浸染、邪念、反噬。。台词："这门……到底是我开的，还是它开的？我好像，也走不出去了。"。动作帧（动图）：①拉冥门 ②门内目光窥视 ③夜半惊醒 ④冷汗涔涔。诗词：冥河浸染蟹甲深，引渡千魂终噬心。夜半惊回冷汗透，门中目光正窥人。。主题句：视杀戮为玩笑，却终被自己的邪念反噬。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：冥河暗绿+蟹甲灰+恐惧黑。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 冥河浸染, age 青年, 被冥河浸染·邪念反噬; scene: 又一次拉开冥门，却第一次感到门内的目光在看他；夜半惊醒，冷汗涔涔。; 冥河浸染蟹甲深，引渡千魂终噬心。夜半惊回冷汗透，门中目光正窥人。; palette: 冥河暗绿+蟹甲灰+恐惧黑; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 巨蟹宫主**
- 巨蟹座·迪斯马斯克，封神登天阶段·巨蟹宫主（黄金圣斗士，黄金巨蟹·黄泉引渡）。形象：巨蟹座·迪斯马斯克，青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。品性：受封黄金圣斗士，巨蟹宫主。黄泉引渡者之名威震冥界，却是黄金中最孤冷的一个。。姿态：镇守巨蟹宫，来者要么过冥门，要么不过；无人敢近，他也习惯了独处。。服饰：黄金巨蟹圣衣，甲壳威光，冥气凝为徽记。。体型：身高约7头身，黄金巨蟹圣衣，孤冷。。衣物细节：黄金巨蟹圣衣，冥气徽记。。发型妆造：束发。。脸型五官：阴柔面容，眼含幽光，嘴角冷峭。。武器招式：积尸气·黄泉引渡。。功法：黄金巨蟹之力；积尸气·黄泉引渡大成。。功法表现：冥气凝徽，黄金威光。。画面：构图：巨蟹宫大殿，黄金巨蟹圣斗士端坐，冥气凝成蟹形徽记，甲壳威光凛冽，宫阙幽冷，背景金色殿堂与冥气交织。色调：黄金+冥气灰绿+殿堂幽冷。氛围：黄金、邪气、孤守。。台词："这身黄金圣衣，是圣域给的。可我这双手，早就习惯黄泉的凉了。"。动作帧（动图）：①端坐宫门 ②冥气凝徽 ③无人敢近 ④独处。诗词：黄金巨蟹镇宫门，黄泉引渡名震魂。黄金虽贵心孤冷，冥气凝徽独处身。。主题句：视杀戮为玩笑，却终被自己的邪念反噬。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+冥气灰绿+殿堂幽冷。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 巨蟹宫主, age 黄金圣斗士, 黄金巨蟹·黄泉引渡; scene: 镇守巨蟹宫，来者要么过冥门，要么不过；无人敢近，他也习惯了独处。; 黄金巨蟹镇宫门，黄泉引渡名震魂。黄金虽贵心孤冷，冥气凝徽独处身。; palette: 黄金+冥气灰绿+殿堂幽冷; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 冥界之王**
- 巨蟹座·迪斯马斯克，道果圆满阶段·黄泉归尽（终极，积尸气·黄泉归尽）。形象：巨蟹座·迪斯马斯克，青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。品性：第八感领悟前夜，他终被邪念反噬——自己也走进了那扇引渡过无数人的冥门。黄泉归尽，善恶终有报。。姿态：最后一战，积尸气反噬其身；他回头看了一眼冥门，忽然笑了——"原来门里这么凉。"。服饰：神圣衣·巨蟹（或黄金终极态），冥气与黄金交织。。体型：身高约7头身，神圣衣巨蟹，没入冥门。。衣物细节：神圣衣·巨蟹，冥气黄金。。发型妆造：束发。。脸型五官：阴柔面容，释然一笑。。武器招式：积尸气·黄泉归尽。。功法：积尸气·黄泉归尽；邪念终噬其身。。功法表现：冥门归尽，黄泉幽光。。画面：构图：冥界之门，黄金巨蟹圣斗士被冥气反噬，身形没入门中，回头一望释然一笑，背景黄泉幽光。色调：冥门灰绿+黄金残辉+黄泉幽蓝。氛围：终极、归尽、反噬。。台词："我送过那么多人进黄泉，今天，黄泉来接我了。报应，来得不冤。"。动作帧（动图）：①积尸气反噬 ②回首望冥门 ③释然一笑 ④没入黄泉。诗词：黄泉归尽蟹甲凉，引渡千魂终自尝。善恶有报天之道，冥门终向自身张。。主题句：视杀戮为玩笑，却终被自己的邪念反噬。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：冥门灰绿+黄金残辉+黄泉幽蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 黄泉归尽, age 终极, 积尸气·黄泉归尽; scene: 最后一战，积尸气反噬其身；他回头看了一眼冥门，忽然笑了——"原来门里这么凉。"; 黄泉归尽蟹甲凉，引渡千魂终自尝。善恶有报天之道，冥门终向自身张。; palette: 冥门灰绿+黄金残辉+黄泉幽蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 狮子座·艾奥里亚（`leo`） · 人生档案版

**灵胎初醒 · 星辉之种**
- 狮子座·艾奥里亚，灵胎初醒阶段·正义初燃（童年，狮子座·正义之火）。形象：狮子座·艾奥里亚，金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。品性：圣域狮子宫幼徒，热血正义，见不平事便按捺不住。。姿态：见强欺弱，冲上去护住弱者；白天练拳，夜里看星星想成为最亮的狮子。。服饰：青铜学徒服，眉目英武。。体型：身高约5头身，青铜学徒少年，英武。。衣物细节：青铜学徒服。。发型妆造：短发，英武。。脸型五官：少年面容，剑眉，虎目炯炯，鼻梁挺，嘴角坚毅。。武器招式：无兵器，拳脚。。功法：小宇宙初醒；正义之心。。功法表现：小宇宙初燃。。画面：构图：圣域狮子宫，青铜学徒少年护住被欺的伙伴，怒目而视，背景金色宫阙。色调：青铜+热血红+狮子宫金。氛围：正义、热血、初燃。。台词："等我长大了，要当圣域最亮的那头狮子，替所有被打的人出头！"。动作帧（动图）：①护住弱者 ②怒目而视 ③挥拳 ④夜里望星。诗词：狮子宫前小火苗，见欺弱小便挥拳。一腔热血藏胸底，誓做圣域最亮天。。主题句：以雷霆之拳，为世间不平而战。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+热血红+狮子宫金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 正义初燃, age 童年, 狮子座·正义之火; scene: 见强欺弱，冲上去护住弱者；白天练拳，夜里看星星想成为最亮的狮子。; 狮子宫前小火苗，见欺弱小便挥拳。一腔热血藏胸底，誓做圣域最亮天。; palette: 青铜+热血红+狮子宫金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 热血修行**
- 狮子座·艾奥里亚，凡尘砺心阶段·光速之拳（少年，练闪电光速拳·拳如雷霆）。形象：狮子座·艾奥里亚，金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。品性：日复一日练闪电光速拳，拳速渐近光速，拳影如雷霆。。姿态：对沙袋连击，拳影如电；练至兴起，一拳轰碎石柱。。服饰：青铜训练服，汗如雨下。。体型：身高约6头身，青铜训练服，身形矫健。。衣物细节：青铜训练服。。发型妆造：短发。。脸型五官：少年面容，剑眉，虎目如电，咧嘴。。武器招式：闪电光速拳初成。。功法：闪电光速拳初成；等离子光速拳初窥。。功法表现：拳影如雷。。画面：构图：训练场，青铜少年对沙袋连击，拳影如闪电，碎石四溅，背景星夜。色调：青铜+电光蓝+汗光金。氛围：苦练、光速、雷霆。。台词："再快一点！这拳头，要快到连光都追不上！"。动作帧（动图）：①蓄拳 ②闪电光速连击 ③拳影如电 ④一拳碎石。诗词：闪电光速练少年，拳影如雷破千岩。汗珠落尽星月转，光速之拳渐齐肩。。主题句：以雷霆之拳，为世间不平而战。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+电光蓝+汗光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 光速之拳, age 少年, 练闪电光速拳·拳如雷霆; scene: 对沙袋连击，拳影如电；练至兴起，一拳轰碎石柱。; 闪电光速练少年，拳影如雷破千岩。汗珠落尽星月转，光速之拳渐齐肩。; palette: 青铜+电光蓝+汗光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 拳光初现**
- 狮子座·艾奥里亚，道法初成阶段·第七感（少年，第七感·光子爆发）。形象：狮子座·艾奥里亚，金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。品性：第七感觉醒，闪电光速拳·光子爆发初成，正义之火燃得更旺。。姿态：一拳轰出，光子爆发划破夜空；狮牙突袭，快若惊雷。。服饰：青铜圣衣具狮子之形，金鬃威光。。体型：身高约6头身，青铜狮子圣衣，英武。。衣物细节：青铜狮子圣衣，金鬃。。发型妆造：短发。。脸型五官：少年面容，剑眉，虎目如炬。。武器招式：闪电光速拳·光子爆发。。功法：闪电光速拳·光子爆发；狮牙突袭。。功法表现：光子光柱冲天。。画面：构图：圣域演武场夜，青铜狮子圣衣少年一拳轰出光子爆发，光柱冲天，金鬃威光，背景夜穹。色调：青铜+光子金+夜穹蓝。氛围：第七感、光子、初成。。台词："师父，我这拳头的光，能不能照到所有不公的地方？"。动作帧（动图）：①蓄拳 ②光子爆发 ③光柱冲天 ④金鬃威光。诗词：第七感觉光子生，闪电光速照夜明。狮牙突袭惊雷起，正义之火渐渐横。。主题句：以雷霆之拳，为世间不平而战。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+光子金+夜穹蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 第七感, age 少年, 第七感·光子爆发; scene: 一拳轰出，光子爆发划破夜空；狮牙突袭，快若惊雷。; 第七感觉光子生，闪电光速照夜明。狮牙突袭惊雷起，正义之火渐渐横。; palette: 青铜+光子金+夜穹蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 闪电光速拳**
- 狮子座·艾奥里亚，大劫淬炼阶段·圣战死斗（青年，守狮子宫·燃烧小宇宙）。形象：狮子座·艾奥里亚，金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。品性：圣战爆发，镇守狮子宫。燃烧小宇宙，以雷霆之拳为世间不平而战。。姿态：狮子宫前，闪电光速拳连发如雨；受伤越重，小宇宙烧得越烈。。服饰：白银狮子圣衣，金鬃战意熊熊。。体型：身高约7头身，白银狮子圣衣，战意熊熊。。衣物细节：白银狮子圣衣。。发型妆造：短发。。脸型五官：青年面容，剑眉，虎目如电。。武器招式：闪电光速拳·光子爆发。。功法：闪电光速拳·光子爆发连击；小宇宙燃烧。。功法表现：拳影如雨，小宇宙燃烧。。画面：构图：狮子宫前圣战，白银狮子圣斗士拳影如雷霆连发，金鬃猎猎，小宇宙燃烧，背景战火宫阙。色调：银白+电光蓝+战火金。氛围：守护、燃烧、小宇宙。。台词："正义也许会迟到，但我的拳头不会！来者，报上名来！"。动作帧（动图）：①立于宫前 ②光速拳连发 ③小宇宙燃烧 ④金鬃猎猎。诗词：狮子宫前战意横，光速拳落如雷霆。燃烧小宇宙不灭，雷霆之拳护圣城。。主题句：以雷霆之拳，为世间不平而战。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：银白+电光蓝+战火金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 圣战死斗, age 青年, 守狮子宫·燃烧小宇宙; scene: 狮子宫前，闪电光速拳连发如雨；受伤越重，小宇宙烧得越烈。; 狮子宫前战意横，光速拳落如雷霆。燃烧小宇宙不灭，雷霆之拳护圣城。; palette: 银白+电光蓝+战火金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 狮子宫主**
- 狮子座·艾奥里亚，封神登天阶段·狮子宫主（黄金圣斗士，黄金狮子·守护女神）。形象：狮子座·艾奥里亚，金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。品性：受封黄金圣斗士，狮子宫主。守护女神的狮子，雷霆之拳守护正义。。姿态：端坐狮子宫，黄金鬃毛威光赫赫；来者过宫，先尝光速之拳。。服饰：黄金狮子圣衣，金鬃威光，雷霆绕拳。。体型：身高约7头身，黄金狮子圣衣，威仪。。衣物细节：黄金狮子圣衣。。发型妆造：短发。。脸型五官：青年面容，剑眉，虎目威光。。武器招式：闪电光速拳（大成）。。功法：黄金狮子之力；闪电光速拳大成。。功法表现：雷霆绕拳，金鬃威光。。画面：构图：狮子宫大殿，黄金狮子圣斗士端坐，金鬃威光赫赫，拳上雷霆流转，身后狮子宫星象，背景金色殿堂。色调：黄金+雷霆蓝+殿堂金。氛围：黄金、守护、正义。。台词："这身黄金圣衣，是正义的分量。我握着它，替女神守着圣域的光。"。动作帧（动图）：①端坐 ②金鬃威光 ③拳上雷霆 ④俯瞰狮子宫。诗词：黄金狮子镇圣宫，光速之拳照万重。守护女神心向正，雷霆一击破敌锋。。主题句：以雷霆之拳，为世间不平而战。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+雷霆蓝+殿堂金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 狮子宫主, age 黄金圣斗士, 黄金狮子·守护女神; scene: 端坐狮子宫，黄金鬃毛威光赫赫；来者过宫，先尝光速之拳。; 黄金狮子镇圣宫，光速之拳照万重。守护女神心向正，雷霆一击破敌锋。; palette: 黄金+雷霆蓝+殿堂金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 等离子光速拳**
- 狮子座·艾奥里亚，道果圆满阶段·雷霆守护（终极，闪电光速拳·光速巅峰）。形象：狮子座·艾奥里亚，金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。品性：第八感领悟，闪电光速拳臻至光速巅峰。雷霆之拳，守护圣域直到最后一刻。。姿态：一拳光速巅峰，雷霆万钧护圣域；小宇宙化作最亮的那头狮子，永远守护。。服饰：神圣衣·狮子（或黄金终极态），金鬃如日，雷霆为翼。。体型：身高约7头身，神圣衣狮子，雷霆为翼。。衣物细节：神圣衣·狮子。。发型妆造：短发。。脸型五官：青年面容，剑眉，虎目含光。。武器招式：闪电光速拳·光速巅峰。。功法：光速巅峰；第八感；雷霆守护。。功法表现：雷霆万钧，光速拳影。。画面：构图：狮子宫顶，狮子圣斗士立于雷霆之巅，光速拳影漫天，金鬃如日，雷霆为翼守护圣域，背景雷霆星河。色调：神光金+雷霆蓝+圣域紫。氛围：终极、光速、守护。。台词："光速有多快？快到所有不公，都来不及躲开我的拳头。"。动作帧（动图）：①立于雷霆之巅 ②光速拳影漫天 ③金鬃如日 ④守护圣域。诗词：雷霆守护狮子心，光速巅峰照夜明。化作圣域最亮宿，守护女神到天明。。主题句：以雷霆之拳，为世间不平而战。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光金+雷霆蓝+圣域紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 雷霆守护, age 终极, 闪电光速拳·光速巅峰; scene: 一拳光速巅峰，雷霆万钧护圣域；小宇宙化作最亮的那头狮子，永远守护。; 雷霆守护狮子心，光速巅峰照夜明。化作圣域最亮宿，守护女神到天明。; palette: 神光金+雷霆蓝+圣域紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 处女座·沙加（`virgo`） · 人生档案版

**灵胎初醒 · 星辉莲花**
- 处女座·沙加，灵胎初醒阶段·佛子降世（童年，佛陀转世·闭目冥想）。形象：处女座·沙加，白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。品性：圣域处女宫幼徒，佛陀转世之身，天生闭目，一念可窥因果。。姿态：众人嬉戏他冥想；闭目静坐，却能"看"见众生因果。。服饰：青铜学徒服，双手合十，双目紧闭。。体型：身高约5头身，青铜学徒小童，清瘦。。衣物细节：青铜学徒服。。发型妆造：短发，佛相。。脸型五官：面容清秀，双目紧闭，眉宇含慈悲，嘴角微阖。。武器招式：无兵器，一念。。功法：小宇宙初醒（佛性）；一念窥因果。。功法表现：佛光微晕。。画面：构图：沙罗双树下，青铜学徒小童合十闭目静坐，身周佛光微晕，背景沙罗双树与圣域。色调：青铜+佛光金+沙罗青。氛围：佛性、降世、闭目。。台词："我闭着眼，可你们每个人的命，我都看得见。"。动作帧（动图）：①合十闭目 ②静坐 ③佛光微晕 ④"看"见众生。诗词：佛子降世处女宫，合十闭目悟玲珑。一念窥得因果线，众生疾苦入眼中。。主题句：闭目是为了看得更远，悟道即是降魔。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+佛光金+沙罗青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 佛子降世, age 童年, 佛陀转世·闭目冥想; scene: 众人嬉戏他冥想；闭目静坐，却能"看"见众生因果。; 佛子降世处女宫，合十闭目悟玲珑。一念窥得因果线，众生疾苦入眼中。; palette: 青铜+佛光金+沙罗青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 禅心初定**
- 处女座·沙加，凡尘砺心阶段·沙罗冥想（少年，沙罗双树下·六道轮回）。形象：处女座·沙加，白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。品性：常在沙罗双树下冥想，六道轮回中看尽悲喜，越发接近神。。姿态：沙罗双树下入定，六道轮回虚影绕身；睁眼时，眼中无喜无悲。。服饰：青铜圣衣渐具，双手合十。。体型：身高约6头身，青铜少年，合十静坐。。衣物细节：青铜圣衣渐具。。发型妆造：短发，佛相。。脸型五官：面容清秀，双目紧闭，佛相渐成。。武器招式：无兵器，六道观想。。功法：六道轮回（观想）；天魔降伏初窥。。功法表现：六道轮回虚影。。画面：构图：沙罗双树下，青铜少年合十入定，六道轮回之相绕身流转，佛光渐盛，背景沙罗叶落。色调：青铜+轮回金光+沙罗青。氛围：冥想、轮回、看尽。。台词："六道轮回，我看了千遍。悲喜都看尽了，才知悲喜都是相。"。动作帧（动图）：①入定 ②六道轮回绕身 ③佛光渐盛 ④睁眼无悲喜。诗词：沙罗树下闭目深，六道轮回入禅心。看尽悲欢千般相，方知万象皆空音。。主题句：闭目是为了看得更远，悟道即是降魔。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+轮回金光+沙罗青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 沙罗冥想, age 少年, 沙罗双树下·六道轮回; scene: 沙罗双树下入定，六道轮回虚影绕身；睁眼时，眼中无喜无悲。; 沙罗树下闭目深，六道轮回入禅心。看尽悲欢千般相，方知万象皆空音。; palette: 青铜+轮回金光+沙罗青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 天眼微启**
- 处女座·沙加，道法初成阶段·天魔降伏（少年，第七感·天魔降伏）。形象：处女座·沙加，白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。品性：第七感觉醒，天魔降伏初成，闭目之间神光自生。最接近神的存在。。姿态：双掌一合，天魔降伏，闭目之间神光自生；一声喝断万邪。。服饰：青铜圣衣具处女之相，白金之光。。体型：身高约6头身，青铜处女圣衣，合十。。衣物细节：青铜处女圣衣。。发型妆造：短发。。脸型五官：面容清秀，闭目，眉宇慈悲含威。。武器招式：天魔降伏。。功法：天魔降伏；不动明王；第七感。。功法表现：神光自生，喝断万邪。。画面：构图：圣域处女宫前，青铜处女圣衣少年双掌一合，天魔降伏神光自生，一声喝断万邪，背景神光初照。色调：白金+神光金+圣域青。氛围：第七感、降魔、神光。。台词："闭目，是为了看得更远。这一掌，降的是魔，也是我自己的执念。"。动作帧（动图）：①双掌一合 ②天魔降伏 ③神光自生 ④一声喝断万邪。诗词：第七感觉降魔音，闭目神光自内生。天魔降伏一声喝，万邪俯首叩禅心。。主题句：闭目是为了看得更远，悟道即是降魔。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：白金+神光金+圣域青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 天魔降伏, age 少年, 第七感·天魔降伏; scene: 双掌一合，天魔降伏，闭目之间神光自生；一声喝断万邪。; 第七感觉降魔音，闭目神光自内生。天魔降伏一声喝，万邪俯首叩禅心。; palette: 白金+神光金+圣域青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 天魔降伏**
- 处女座·沙加，大劫淬炼阶段·沙罗之战（青年，沙罗双树下一战·舍身悟道）。形象：处女座·沙加，白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。品性：圣战之中，沙罗双树下一战，血染双树。舍身之际，悟得阿赖耶识——第八感，最接近神的他，终于成神。。姿态：沙罗双树下独战群敌，血染叶落；临死之际合十微笑，悟得第八感，佛光冲霄。。服饰：白银处女圣衣，染血。。体型：身高约7头身，白银处女圣衣，染血。。衣物细节：白银处女圣衣，染血。。发型妆造：短发。。脸型五官：面容清秀，闭目微笑，佛相圆满。。武器招式：天舞宝轮/天魔降伏。。功法：天舞宝轮；阿赖耶识（第八感）。。功法表现：佛光冲霄，第八感。。画面：构图：沙罗双树下，白银处女圣斗士血染圣衣立于落叶之中，合十微笑，佛光冲霄，六道轮回之相在光中流转，背景双树血染。色调：白银+佛光金+血绛+沙罗青。氛围：舍身、悟道、第八感。。台词："花开，花谢，皆是无常。这一战，我看到了第八感——那是神的目光。"。动作帧（动图）：①立于双树下 ②独战群敌 ③血染叶落 ④合十悟第八感。诗词：沙罗树下血染红，独战群魔气如虹。临死合十悟第八，佛光冲霄化作龙。。主题句：闭目是为了看得更远，悟道即是降魔。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：白银+佛光金+血绛+沙罗青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 沙罗之战, age 青年, 沙罗双树下一战·舍身悟道; scene: 沙罗双树下独战群敌，血染叶落；临死之际合十微笑，悟得第八感，佛光冲霄。; 沙罗树下血染红，独战群魔气如虹。临死合十悟第八，佛光冲霄化作龙。; palette: 白银+佛光金+血绛+沙罗青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 处女宫主**
- 处女座·沙加，封神登天阶段·处女宫主（黄金圣斗士，黄金处女·最接近神）。形象：处女座·沙加，白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。品性：受封黄金圣斗士，处女宫主。最接近神的人，镇守处女宫。。姿态：端坐处女宫，闭目含光；来者过宫，先过他的天舞宝轮。。服饰：黄金处女圣衣，白金之光，闭目禅相。。体型：身高约7头身，黄金处女圣衣，闭目端坐。。衣物细节：黄金处女圣衣。。发型妆造：短发。。脸型五官：面容清秀，闭目，佛相慈悲。。武器招式：天舞宝轮。。功法：黄金处女之力；天舞宝轮大成；第八感。。功法表现：白金之光，天舞宝轮。。画面：构图：处女宫大殿，黄金处女圣斗士闭目端坐，白金之光环绕，天舞宝轮虚影，身后处女宫星象，背景金色殿堂。色调：黄金+白金+佛光金。氛围：黄金、最接近神、镇守。。台词："我不是神，我只是看得比你们远一些。这处女宫，我来守。"。动作帧（动图）：①闭目端坐 ②白金之光 ③天舞宝轮 ④俯瞰处女宫。诗词：黄金处女镇圣宫，闭目禅相万法通。天舞宝轮随念起，最接近神守苍穹。。主题句：闭目是为了看得更远，悟道即是降魔。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+白金+佛光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 处女宫主, age 黄金圣斗士, 黄金处女·最接近神; scene: 端坐处女宫，闭目含光；来者过宫，先过他的天舞宝轮。; 黄金处女镇圣宫，闭目禅相万法通。天舞宝轮随念起，最接近神守苍穹。; palette: 黄金+白金+佛光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 第八感觉醒**
- 处女座·沙加，道果圆满阶段·万法皆空（终极，阿赖耶识·万法皆空）。形象：处女座·沙加，白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。品性：第八感圆满，阿赖耶识·万法皆空。最接近神的人，终于看透——一切皆空，唯有慈悲。。姿态：合十而立，六道轮回尽归空；天舞宝轮化作万法皆空之光，渡尽众生。。服饰：神圣衣·处女（或黄金终极态），佛光圆满，莲台为座。。体型：身高约7头身，神圣衣处女，莲台为座。。衣物细节：神圣衣·处女。。发型妆造：短发。。脸型五官：面容清秀，闭目，佛相圆满慈悲。。武器招式：天舞宝轮·阿赖耶识。。功法：阿赖耶识·万法皆空；第八感圆满；渡尽众生。。功法表现：万法皆空，佛光遍洒。。画面：构图：云端莲台，处女圣斗士合十而立，佛光圆满，六道轮回化作万法皆空之光遍洒，背景浩瀚星空与金莲。色调：神光白+佛光金+星海紫。氛围：终极、皆空、觉者。。台词："我看了千遍轮回，终于看透——色即是空，空即是色。慈悲，是唯一的真实。"。动作帧（动图）：①合十而立 ②六道归空 ③万法皆空之光 ④遍洒渡众。诗词：万法皆空处女尊，阿赖耶识照乾坤。轮回看尽归空处，唯有慈悲渡世人。。主题句：闭目是为了看得更远，悟道即是降魔。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白+佛光金+星海紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 万法皆空, age 终极, 阿赖耶识·万法皆空; scene: 合十而立，六道轮回尽归空；天舞宝轮化作万法皆空之光，渡尽众生。; 万法皆空处女尊，阿赖耶识照乾坤。轮回看尽归空处，唯有慈悲渡世人。; palette: 神光白+佛光金+星海紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 天秤座·童虎（`libra`） · 人生档案版

**灵胎初醒 · 星辉龙珠**
- 天秤座·童虎，灵胎初醒阶段·庐山修行（少年，庐山瀑布下修行）。形象：天秤座·童虎，白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。品性：圣域天秤宫少年，在庐山瀑布下修行，与龙为伴。。姿态：庐山瀑布下练拳，水花四溅；山中采药，与老龙对话。。服饰：青铜学徒服，少年英气。。体型：身高约5头身，青铜学徒少年，英武。。衣物细节：青铜学徒服。。发型妆造：短发。。脸型五官：少年面容，浓眉，虎目有神，鼻梁挺，嘴角坚毅。。武器招式：无兵器，拳。。功法：小宇宙初醒；庐山升龙霸初窥。。功法表现：水花随拳起。。画面：构图：庐山瀑布前，青铜学徒少年练拳，水花飞溅，背景飞瀑青山。色调：青铜+水花白+庐山青。氛围：苦修、瀑布、少年。。台词："瀑布日夜落，拳法日日新。这庐山的龙气，我都记在拳里了。"。动作帧（动图）：①瀑布下练拳 ②水花飞溅 ③与龙对话 ④望瀑布沉思。诗词：庐山瀑布练拳初，水花溅处龙气濡。少年不识愁滋味，只把山河入拳书。。主题句：两百年的守望，只为一战报圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+水花白+庐山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 庐山修行, age 少年, 庐山瀑布下修行; scene: 庐山瀑布下练拳，水花四溅；山中采药，与老龙对话。; 庐山瀑布练拳初，水花溅处龙气濡。少年不识愁滋味，只把山河入拳书。; palette: 青铜+水花白+庐山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 百裂拳初修**
- 天秤座·童虎，凡尘砺心阶段·返老还童（少年，返老还童秘法·上代圣战）。形象：天秤座·童虎，白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。品性：上代圣战在即，受赐返老还童秘法（每243年苏醒一次），为的是守候一场更久远的圣战。。姿态：习得返老还童秘法，肉身凝为少年，以换两百年后的重逢。。服饰：青铜圣衣，少年身形却含沧桑。。体型：身高约5头身，少年身形，眼中含沧桑。。衣物细节：青铜圣衣。。发型妆造：短发。。脸型五官：少年面容，眼神却苍老深沉。。武器招式：庐山升龙霸初成。。功法：返老还童秘法；庐山升龙霸初成。。功法表现：时光流转虚影。。画面：构图：庐山之巅，青铜少年习返老还童秘法，身周时光流转虚影（老/少/老），背景山川千年。色调：青铜+时光金+庐山青。氛围：秘法、返老、守望。。台词："返老还童，不是贪生，是要用两百年，等一场必须打的仗。"。动作帧（动图）：①习秘法 ②时光流转 ③少年凝形 ④眼中含沧桑。诗词：返老还童秘法修，少年之身藏千秋。不是贪生求寿久，只缘圣战待重头。。主题句：两百年的守望，只为一战报圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+时光金+庐山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 返老还童, age 少年, 返老还童秘法·上代圣战; scene: 习得返老还童秘法，肉身凝为少年，以换两百年后的重逢。; 返老还童秘法修，少年之身藏千秋。不是贪生求寿久，只缘圣战待重头。; palette: 青铜+时光金+庐山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 五老峰悟道**
- 天秤座·童虎，道法初成阶段·升龙霸（少年，第七感·庐山升龙霸）。形象：天秤座·童虎，白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。品性：第七感觉醒，庐山升龙霸初成，一拳升龙破云。。姿态：庐山瀑布前一拳轰出，升龙破云而上；龙吟声中，瀑布倒卷。。服饰：青铜圣衣具天秤之相。。体型：身高约6头身，青铜天秤圣衣，英武。。衣物细节：青铜天秤圣衣。。发型妆造：短发。。脸型五官：少年面容，浓眉，虎目如龙。。武器招式：庐山升龙霸。。功法：庐山升龙霸；龙飞翔。。功法表现：升龙破云，瀑布倒卷。。画面：构图：庐山瀑布，青铜天秤圣衣少年一拳升龙破云，瀑布倒卷，龙影冲天，背景青山。色调：青铜+升龙金+水光蓝。氛围：第七感、升龙、初成。。台词："龙从拳里生，云从拳里破。这一拳，庐山认得，我也认得。"。动作帧（动图）：①蓄拳 ②升龙破云 ③瀑布倒卷 ④龙影冲天。诗词：第七感觉升龙吟，一拳破云瀑布惊。龙从拳起云从散，少年初露霸主形。。主题句：两百年的守望，只为一战报圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+升龙金+水光蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 升龙霸, age 少年, 第七感·庐山升龙霸; scene: 庐山瀑布前一拳轰出，升龙破云而上；龙吟声中，瀑布倒卷。; 第七感觉升龙吟，一拳破云瀑布惊。龙从拳起云从散，少年初露霸主形。; palette: 青铜+升龙金+水光蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 庐山升龙霸**
- 天秤座·童虎，大劫淬炼阶段·上代圣战（上代，上代圣战·血战苍穹）。形象：天秤座·童虎，白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。品性：上代圣战爆发，他血战苍穹，重伤濒死。靠着秘法沉睡，等着下一场圣战。。姿态：庐山顶上，以一敌众，百龙齐出；重伤之际，望了一眼圣域方向，沉沉睡去。。服饰：白银天秤圣衣，战痕遍布。。体型：身高约7头身，白银天秤圣衣，战痕。。衣物细节：白银天秤圣衣。。发型妆造：短发。。脸型五官：少年面容染血，虎目不屈。。武器招式：庐山百龙霸。。功法：庐山百龙霸（初现）；以一敌众。。功法表现：百龙齐出。。画面：构图：庐山顶，白银天秤圣斗士百龙齐出，浴血而战，身后圣域方向一抹曙光，背景战火残云。色调：白银+战火红+庐山青。氛围：圣战、血战、两百年。。台词："这一战，我撑到最后一刻。剩下的，留给两百年后的童虎。"。动作帧（动图）：①以一敌众 ②百龙齐出 ③重伤 ④望圣域沉睡。诗词：上代圣战血苍穹，庐山顶上百龙冲。重伤沉眠蓄一诺，待得重生再交锋。。主题句：两百年的守望，只为一战报圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：白银+战火红+庐山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 上代圣战, age 上代, 上代圣战·血战苍穹; scene: 庐山顶上，以一敌众，百龙齐出；重伤之际，望了一眼圣域方向，沉沉睡去。; 上代圣战血苍穹，庐山顶上百龙冲。重伤沉眠蓄一诺，待得重生再交锋。; palette: 白银+战火红+庐山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 老师坐镇**
- 天秤座·童虎，封神登天阶段·五老峰守望（黄金圣斗士，独守五老峰两百年·紫龙之师）。形象：天秤座·童虎，白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。品性：受封黄金圣斗士，独守五老峰两百年，静候圣战。收紫龙为徒，将一身武学与守望之心尽授。。姿态：五老峰上负手而立，望了圣域两百年；授紫龙庐山升龙霸，也授他守望之心。。服饰：黄金天秤圣衣，背负十二件黄金武器。。体型：身高约7头身，黄金天秤圣衣，背负武器。。衣物细节：黄金天秤圣衣，十二件武器。。发型妆造：短发。。脸型五官：少年面容，浓眉，虎目含沧桑。。武器招式：庐山百龙霸。。功法：黄金天秤之力；庐山百龙霸大成；师者之心。。功法表现：黄金武器列阵。。画面：构图：五老峰巅，黄金天秤圣斗士负手而立，背后十二件黄金武器，紫龙少年在侧练拳，背景群山与圣域远影。色调：黄金+五峰青+圣域金。氛围：守望、师者、两百年。。台词："两百年的守望，不是怕死，是怕等不到这场仗。紫龙，为师等你长大。"。动作帧（动图）：①负手立五峰 ②望圣域 ③授徒升龙 ④静候圣战。诗词：黄金天秤立五峰，两百年守望苍穹。授徒升龙传道义，只待圣战再相逢。。主题句：两百年的守望，只为一战报圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+五峰青+圣域金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 五老峰守望, age 黄金圣斗士, 独守五老峰两百年·紫龙之师; scene: 五老峰上负手而立，望了圣域两百年；授紫龙庐山升龙霸，也授他守望之心。; 黄金天秤立五峰，两百年守望苍穹。授徒升龙传道义，只待圣战再相逢。; palette: 黄金+五峰青+圣域金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 返老还童**
- 天秤座·童虎，道果圆满阶段·升龙传说（终极，庐山百龙霸·升龙归元）。形象：天秤座·童虎，白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。品性：第八感领悟，庐山百龙霸·升龙归元。两百年守望，终成升龙传说，报圣域于最后一战。。姿态：庐山百龙霸全力轰出，百龙归元成龙神；一拳开天，护住圣域与弟子。。服饰：神圣衣·天秤（或黄金终极态），龙气缠身。。体型：身高约7头身，神圣衣天秤，龙气缠身。。衣物细节：神圣衣·天秤。。发型妆造：短发。。脸型五官：少年面容，虎目如龙，含光。。武器招式：庐山百龙霸·升龙归元。。功法：庐山百龙霸·升龙归元；第八感。。功法表现：百龙归元，龙神冲天。。画面：构图：五老峰巅，天秤圣斗士一拳轰出，百龙归元成龙神冲天，圣域与群山皆护，背景天地变色。色调：龙神金+升龙蓝+神光白。氛围：终极、升龙、传说。。台词："两百年，就为这一拳。紫龙，看好了——这就是庐山的传说！"。动作帧（动图）：①蓄拳 ②百龙齐出 ③归元成龙神 ④一拳开天。诗词：升龙归元震苍穹，两百年诺一朝倾。百龙化神开天去，庐山传说镇圣城。。主题句：两百年的守望，只为一战报圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：龙神金+升龙蓝+神光白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 升龙传说, age 终极, 庐山百龙霸·升龙归元; scene: 庐山百龙霸全力轰出，百龙归元成龙神；一拳开天，护住圣域与弟子。; 升龙归元震苍穹，两百年诺一朝倾。百龙化神开天去，庐山传说镇圣城。; palette: 龙神金+升龙蓝+神光白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 天蝎座·米罗（`scorpio`） · 人生档案版

**灵胎初醒 · 星辉之种**
- 天蝎座·米罗，灵胎初醒阶段·忠义初立（童年，天蝎座·忠义少年）。形象：天蝎座·米罗，深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。品性：圣域天蝎宫幼徒，重情重义的少年，认定了朋友就护到底。。姿态：为兄弟两肋插刀；练针之余，与伙伴并肩立誓守护圣域。。服饰：青铜学徒服，眉目英气。。体型：身高约5头身，青铜学徒少年，英气。。衣物细节：青铜学徒服。。发型妆造：短发。。脸型五官：少年面容，剑眉，虎目坚定，鼻梁挺，嘴角坚毅。。武器招式：无兵器，针意初窥。。功法：小宇宙初醒；忠义之心。。功法表现：小宇宙初燃。。画面：构图：圣域天蝎宫，青铜学徒少年与伙伴并肩而立，眼神坚定，背景金色宫阙。色调：青铜+天蝎红+宫阙金。氛围：忠义、少年、热血。。台词："我这辈子，认了兄弟就认到底。圣域的朋友，谁也不能动。"。动作帧（动图）：①与伙伴立誓 ②练针 ③护友 ④并肩望圣域。诗词：天蝎少年义气浓，为友能闯九重峰。一腔热血藏针底，认准情义便从容。。主题句：毒针虽烈，只刺敌人，义字当先。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+天蝎红+宫阙金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 忠义初立, age 童年, 天蝎座·忠义少年; scene: 为兄弟两肋插刀；练针之余，与伙伴并肩立誓守护圣域。; 天蝎少年义气浓，为友能闯九重峰。一腔热血藏针底，认准情义便从容。; palette: 青铜+天蝎红+宫阙金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 忠义修行**
- 天蝎座·米罗，凡尘砺心阶段·毒针初练（少年，深红毒针·只刺敌人）。形象：天蝎座·米罗，深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。品性：深红毒针初练，猩红星光流转。毒针虽烈，只刺敌人，义字当先。。姿态：对靶练针，猩红针影如雨；练罢，只对自己人露出少年憨笑。。服饰：青铜训练服，指尖猩红光。。体型：身高约6头身，青铜训练服，指尖星光。。衣物细节：青铜训练服。。发型妆造：短发。。脸型五官：少年面容，剑眉，虎目专注。。武器招式：深红毒针初成。。功法：深红毒针初成；毒针点穴。。功法表现：猩红针影。。画面：构图：训练场，青铜少年指尖猩红针影如雨点向靶心，猩红星光流转，背景天蝎宫。色调：青铜+深红针光+天蝎暗红。氛围：毒针、忠义、深红。。台词："这针很毒，可我的毒，只喂给敌人。朋友面前，我还是那个米罗。"。动作帧（动图）：①指尖蓄光 ②深红针影如雨 ③点中靶心 ④对友憨笑。诗词：深红毒针初练成，猩红星光指间生。针毒只向敌人刺，义字当先少年行。。主题句：毒针虽烈，只刺敌人，义字当先。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+深红针光+天蝎暗红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 毒针初练, age 少年, 深红毒针·只刺敌人; scene: 对靶练针，猩红针影如雨；练罢，只对自己人露出少年憨笑。; 深红毒针初练成，猩红星光指间生。针毒只向敌人刺，义字当先少年行。; palette: 青铜+深红针光+天蝎暗红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 毒针微芒**
- 天蝎座·米罗，道法初成阶段·猩红连刺（少年，第七感·猩红连刺）。形象：天蝎座·米罗，深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。品性：第七感觉醒，猩红连刺初成，深红毒针·心宿二之威渐显。。姿态：指尖连点，猩红连刺如暴雨；心宿二·终针一出，星光炸裂。。服饰：青铜圣衣具天蝎之形，尾针猩红。。体型：身高约6头身，青铜天蝎圣衣。。衣物细节：青铜天蝎圣衣。。发型妆造：短发。。脸型五官：少年面容，剑眉，虎目如电。。武器招式：猩红连刺/心宿二·终针。。功法：猩红连刺；心宿二·终针。。功法表现：猩红针雨。。画面：构图：圣域演武场，青铜天蝎圣衣少年指尖猩红连刺如暴雨，尾针猩红星芒，背景夜穹。色调：青铜+猩红针光+天蝎暗红。氛围：第七感、连刺、初成。。台词："十五针封喉，是我给敌人的礼数。朋友，我一根针都不舍得。"。动作帧（动图）：①指尖蓄光 ②猩红连刺 ③心宿二终针 ④尾针猩红。诗词：第七感觉针如虹，猩红连刺破长空。心宿二光终针出，天蝎之尾立苍穹。。主题句：毒针虽烈，只刺敌人，义字当先。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+猩红针光+天蝎暗红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 猩红连刺, age 少年, 第七感·猩红连刺; scene: 指尖连点，猩红连刺如暴雨；心宿二·终针一出，星光炸裂。; 第七感觉针如虹，猩红连刺破长空。心宿二光终针出，天蝎之尾立苍穹。; palette: 青铜+猩红针光+天蝎暗红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 深红毒针**
- 天蝎座·米罗，大劫淬炼阶段·天蝎宫战（青年，守天蝎宫·毒针义战）。形象：天蝎座·米罗，深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。品性：圣战爆发，镇守天蝎宫。毒针只刺敌人，义字当先，为守护圣域而战。。姿态：天蝎宫前，深红毒针连发；对来犯者毫不留情，对战友转身便笑。。服饰：白银天蝎圣衣，尾针猩红凛冽。。体型：身高约7头身，白银天蝎圣衣。。衣物细节：白银天蝎圣衣。。发型妆造：短发。。脸型五官：青年面容，剑眉，虎目坚定。。武器招式：深红毒针·猩红连刺。。功法：深红毒针·猩红连刺连击；忠义之心。。功法表现：毒针如雨。。画面：构图：天蝎宫前圣战，白银天蝎圣斗士毒针如雨破敌，尾针猩红，背后圣域宫阙，背景战火。色调：银白+猩红针光+战火金。氛围：守护、义战、毒针。。台词："我的针，认敌友。是敌，十五针送你；是友，我米罗护你到底！"。动作帧（动图）：①立于宫前 ②毒针连发 ③破敌 ④护战友。诗词：天蝎宫前战意横，深红毒针破敌营。义字当先针有眼，护得圣域到天明。。主题句：毒针虽烈，只刺敌人，义字当先。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：银白+猩红针光+战火金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 天蝎宫战, age 青年, 守天蝎宫·毒针义战; scene: 天蝎宫前，深红毒针连发；对来犯者毫不留情，对战友转身便笑。; 天蝎宫前战意横，深红毒针破敌营。义字当先针有眼，护得圣域到天明。; palette: 银白+猩红针光+战火金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 天蝎宫主**
- 天蝎座·米罗，封神登天阶段·天蝎宫主（黄金圣斗士，黄金天蝎·忠义战士）。形象：天蝎座·米罗，深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。品性：受封黄金圣斗士，天蝎宫主。重情重义的战士，毒针只向敌人。。姿态：端坐天蝎宫，深红毒针凛然；来者过宫，先尝十五针的义气。。服饰：黄金天蝎圣衣，尾针猩红星光。。体型：身高约7头身，黄金天蝎圣衣。。衣物细节：黄金天蝎圣衣。。发型妆造：短发。。脸型五官：青年面容，剑眉，虎目含光。。武器招式：深红毒针（大成）。。功法：黄金天蝎之力；深红毒针大成。。功法表现：猩红星光。。画面：构图：天蝎宫大殿，黄金天蝎圣斗士端坐，尾针猩红星光流转，身后天蝎宫星象，背景金色殿堂。色调：黄金+猩红针光+殿堂金。氛围：黄金、忠义、镇守。。台词："这身黄金圣衣，是义气的分量。谁与圣域为敌，我就给他十五针。"。动作帧（动图）：①端坐 ②尾针猩红 ③来者十五针 ④俯瞰天蝎宫。诗词：黄金天蝎镇圣宫，深红毒针指苍穹。忠义战士心如铁，十五针下定敌锋。。主题句：毒针虽烈，只刺敌人，义字当先。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+猩红针光+殿堂金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 天蝎宫主, age 黄金圣斗士, 黄金天蝎·忠义战士; scene: 端坐天蝎宫，深红毒针凛然；来者过宫，先尝十五针的义气。; 黄金天蝎镇圣宫，深红毒针指苍穹。忠义战士心如铁，十五针下定敌锋。; palette: 黄金+猩红针光+殿堂金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 天蝎之怒**
- 天蝎座·米罗，道果圆满阶段·心宿之刺（终极，深红毒针·心宿之刺）。形象：天蝎座·米罗，深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。品性：第八感领悟，深红毒针·心宿之刺。忠义之心化入针尖，为圣域刺出最后一记。。姿态：指尖一点，心宿之刺破空，猩红星光如天蝎之心；刺尽敌，亦护尽友。。服饰：神圣衣·天蝎（或黄金终极态），尾针如星，猩红为光。。体型：身高约7头身，神圣衣天蝎。。衣物细节：神圣衣·天蝎。。发型妆造：短发。。脸型五官：青年面容，剑眉，虎目含光，嘴角含笑。。武器招式：深红毒针·心宿之刺。。功法：深红毒针·心宿之刺；第八感。。功法表现：猩红星芒贯穿夜空。。画面：构图：天蝎宫顶，天蝎圣斗士指尖一点心宿之刺，猩红星光如天蝎之心贯穿夜空，护住圣域，背景星河天蝎座。色调：猩红+星芒金+圣域紫。氛围：终极、心宿、忠义。。台词："天蝎的心，就藏在这一针里。为敌，是一记心宿之刺；为友，是一生守护。"。动作帧（动图）：①指尖蓄光 ②心宿之刺 ③猩红星芒 ④护圣域。诗词：心宿之刺贯长空，猩红星光照圣宫。忠义化针深红里，天蝎之心永从容。。主题句：毒针虽烈，只刺敌人，义字当先。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：猩红+星芒金+圣域紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 心宿之刺, age 终极, 深红毒针·心宿之刺; scene: 指尖一点，心宿之刺破空，猩红星光如天蝎之心；刺尽敌，亦护尽友。; 心宿之刺贯长空，猩红星光照圣宫。忠义化针深红里，天蝎之心永从容。; palette: 猩红+星芒金+圣域紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 射手座·艾俄洛斯（`sagittarius`） · 人生档案版

**灵胎初醒 · 星辉之箭**
- 射手座·艾俄洛斯，灵胎初醒阶段·仁心少年（童年，仁之圣斗士·天性仁厚）。形象：射手座·艾俄洛斯，金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。品性：圣域射手宫幼徒，仁之圣斗士，天性仁厚，见不得任何人受难。。姿态：与伙伴玩耍，总把好的让给别人；夜里看射手座，立志守护弱者。。服饰：青铜学徒服，眉目温厚。。体型：身高约5头身，青铜学徒少年，温厚。。衣物细节：青铜学徒服。。发型妆造：短发。。脸型五官：少年面容，眉目温厚，眼含仁光，鼻梁挺，嘴角温和。。武器招式：无兵器。。功法：小宇宙初醒；仁心。。功法表现：小宇宙仁光。。画面：构图：圣域射手宫，青铜学徒少年抬头望射手座星，眼神温厚坚定，背景星空宫阙。色调：青铜+星月光+射手金。氛围：仁厚、少年、初醒。。台词："师父说，射手座是仁的星座。那我就用这一生，护住所有需要守护的人。"。动作帧（动图）：①让出好东西 ②夜里望射手座 ③立志守护 ④握拳。诗词：仁心少年望星穹，射手座光照心胸。让尽身前好与物，只愿人间少苦容。。主题句：肉身虽灭，星矢长存，永远守望圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+星月光+射手金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 仁心少年, age 童年, 仁之圣斗士·天性仁厚; scene: 与伙伴玩耍，总把好的让给别人；夜里看射手座，立志守护弱者。; 仁心少年望星穹，射手座光照心胸。让尽身前好与物，只愿人间少苦容。; palette: 青铜+星月光+射手金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 弓弦初张**
- 射手座·艾俄洛斯，凡尘砺心阶段·黄金之弓（少年，练黄金之箭·搭弦苍穹）。形象：射手座·艾俄洛斯，金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。品性：苦练黄金之箭，弓意渐成。箭随仁心，所向皆为护人。。姿态：对月练箭，一箭穿云；搭弦之间，仁心与箭意合一。。服饰：青铜训练服，背负黄金之弓。。体型：身高约6头身，青铜训练服，持弓。。衣物细节：青铜训练服，黄金之弓。。发型妆造：短发。。脸型五官：少年面容，眉目温厚，眼含坚毅。。武器招式：黄金之箭初练。。功法：黄金之箭初练；箭意。。功法表现：箭光穿云。。画面：构图：圣域月夜，青铜少年张弓搭箭，一箭穿云，金色箭光划破夜空，背景月宫星野。色调：青铜+箭光金+月白。氛围：练箭、仁厚、弓意。。台词："我的箭，不射人，只射向黑暗——护住箭后的每一个人。"。动作帧（动图）：①张弓 ②搭箭 ③一箭穿云 ④收弓望星。诗词：黄金之弓月下张，一箭穿云射天狼。仁心搭弦意随矢，只向黑暗不向乡。。主题句：肉身虽灭，星矢长存，永远守望圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+箭光金+月白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 黄金之弓, age 少年, 练黄金之箭·搭弦苍穹; scene: 对月练箭，一箭穿云；搭弦之间，仁心与箭意合一。; 黄金之弓月下张，一箭穿云射天狼。仁心搭弦意随矢，只向黑暗不向乡。; palette: 青铜+箭光金+月白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 闪电拳雏形**
- 射手座·艾俄洛斯，道法初成阶段·第七感（少年，第七感·正义之矢）。形象：射手座·艾俄洛斯，金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。品性：第七感觉醒，正义之矢初成。仁心化箭，为正义而射。。姿态：一箭射出，正义之矢破空；箭落处，黑暗尽散。。服饰：青铜圣衣具射手之形，背负黄金弓。。体型：身高约6头身，青铜射手圣衣。。衣物细节：青铜射手圣衣。。发型妆造：短发。。脸型五官：少年面容，眉目温厚，眼含正义。。武器招式：正义之矢。。功法：正义之矢；黄金之箭初成。。功法表现：正义金光破空。。画面：构图：圣域演武场，青铜射手圣衣少年一箭正义之矢破空，金光划开黑暗，背景夜穹。色调：青铜+正义金+夜蓝。氛围：第七感、正义、之矢。。台词："这支箭，是替所有不敢开口的人射的。"。动作帧（动图）：①张弓 ②正义之矢 ③破空 ④黑暗尽散。诗词：第七感觉正义生，仁心化矢破空行。箭落黑暗皆散尽，敢为无声发雷声。。主题句：肉身虽灭，星矢长存，永远守望圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+正义金+夜蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 第七感, age 少年, 第七感·正义之矢; scene: 一箭射出，正义之矢破空；箭落处，黑暗尽散。; 第七感觉正义生，仁心化矢破空行。箭落黑暗皆散尽，敢为无声发雷声。; palette: 青铜+正义金+夜蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 黄金之箭**
- 射手座·艾俄洛斯，大劫淬炼阶段·圣战牺牲（青年，为护雅典娜·中箭陨落）。形象：射手座·艾俄洛斯，金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。品性：圣战之中，为保护幼小的雅典娜，以身为盾，中箭陨落。仁之圣斗士，用生命守住了女神。。姿态：抱起初生的小雅典娜突围，背后中箭；临终将襁褓托付给同伴，含笑闭目。。服饰：白银射手圣衣，中箭染血。。体型：身高约7头身，白银射手圣衣，中箭。。衣物细节：白银射手圣衣，染血。。发型妆造：短发。。脸型五官：青年面容，含笑，眼神安详。。武器招式：以身护神。。功法：仁之守护；以身护神。。功法表现：仁心化光。。画面：构图：圣战乱局，白银射手圣斗士抱着襁褓突围，背后中箭，血染圣衣，仍紧护怀中婴儿，背景烽火圣域。色调：白银+血绛+圣火金。氛围：牺牲、护神、英魂。。台词："她……是圣域的光。我护住了她，就是护住了所有人的明天。"。动作帧（动图）：①抱婴突围 ②背后中箭 ③托付同伴 ④含笑陨落。诗词：圣战烽火护婴身，背后中箭血染尘。含笑托孤星陨落，仁心化光照圣门。。主题句：肉身虽灭，星矢长存，永远守望圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：白银+血绛+圣火金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 圣战牺牲, age 青年, 为护雅典娜·中箭陨落; scene: 抱起初生的小雅典娜突围，背后中箭；临终将襁褓托付给同伴，含笑闭目。; 圣战烽火护婴身，背后中箭血染尘。含笑托孤星陨落，仁心化光照圣门。; palette: 白银+血绛+圣火金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 英魂长存**
- 射手座·艾俄洛斯，封神登天阶段·星空守望（英魂，射手座英魂·永守圣域）。形象：射手座·艾俄洛斯，金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。品性：肉身虽灭，星矢长存。射手座英魂永在星空中守望圣域，护佑后辈。。姿态：化为射手座星辉，于夜穹注视圣域；危难之时，星光降下守护。。服饰：英魂星辉，黄金之弓虚影。。体型：英魂星辉虚影，持弓。。衣物细节：英魂星辉，黄金弓虚影。。发型妆造：短发。。脸型五官：青年面容，含笑，眼神温柔守望。。武器招式：星辉之箭。。功法：英魂守望；星辉护佑。。功法表现：星辉如雨护圣域。。画面：构图：夜穹射手座星辉流转，英魂虚影持弓立于星间，俯视下方圣域宫阙，星辉如雨洒落，背景浩瀚星空。色调：星辉银+射手金+夜穹蓝。氛围：英魂、守望、星空。。台词："我不在了，可这支箭还在。星光落处，就是我在护着你们。"。动作帧（动图）：①化星辉 ②立星空 ③俯视圣域 ④星辉洒落。诗词：肉身虽灭魂长存，射手星辉照圣门。夜夜望断云深处，星光落处是吾身。。主题句：肉身虽灭，星矢长存，永远守望圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星辉银+射手金+夜穹蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 星空守望, age 英魂, 射手座英魂·永守圣域; scene: 化为射手座星辉，于夜穹注视圣域；危难之时，星光降下守护。; 肉身虽灭魂长存，射手星辉照圣门。夜夜望断云深处，星光落处是吾身。; palette: 星辉银+射手金+夜穹蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 黄金之箭·觉醒**
- 射手座·艾俄洛斯，道果圆满阶段·穿云破晓（终极，黄金之箭·穿云破晓）。形象：射手座·艾俄洛斯，金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。品性：第八感领悟，黄金之箭·穿云破晓。星矢长存，永远守望——他的箭，终于射穿了黑夜，迎来破晓。。姿态：黄金之箭搭弦，一箭穿云破晓，金色曙光遍洒圣域；英魂化光，护住圣域到天明。。服饰：神圣衣·射手（或黄金终极态），星辉为翼。。体型：身高约7头身，神圣衣射手，星辉为翼。。衣物细节：神圣衣·射手。。发型妆造：短发。。脸型五官：青年面容，含笑，眼含破晓之光。。武器招式：黄金之箭·穿云破晓。。功法：黄金之箭·穿云破晓；第八感。。功法表现：破晓金光冲破黑夜。。画面：构图：夜穹之巅，射手圣斗士张弓搭箭，一箭穿云破晓，金色曙光冲破黑夜遍洒圣域，星辉为翼，背景破晓天地。色调：曙光金+破晓橙+星辉银。氛围：终极、破晓、星矢。。台词："我守了一辈子黑夜，终于等来这一箭——天亮了。"。动作帧（动图）：①张弓 ②搭箭 ③穿云破晓 ④曙光遍洒。诗词：穿云破晓射天光，仁心化箭照万方。星矢长存终有报，一箭迎得圣域阳。。主题句：肉身虽灭，星矢长存，永远守望圣域。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：曙光金+破晓橙+星辉银。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 穿云破晓, age 终极, 黄金之箭·穿云破晓; scene: 黄金之箭搭弦，一箭穿云破晓，金色曙光遍洒圣域；英魂化光，护住圣域到天明。; 穿云破晓射天光，仁心化箭照万方。星矢长存终有报，一箭迎得圣域阳。; palette: 曙光金+破晓橙+星辉银; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 摩羯座·修罗（`capricorn`） · 人生档案版

**灵胎初醒 · 星辉之剑**
- 摩羯座·修罗，灵胎初醒阶段·奉剑之志（童年，奉女神的剑心）。形象：摩羯座·修罗，银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。品性：圣域摩羯宫幼徒，奉女神的剑心，天生便是为剑而生。。姿态：捡一根树枝当剑，练到忘我；心里只有一个念头——为女神铸一柄忠诚的剑。。服饰：青铜学徒服，眉目锐利。。体型：身高约5头身，青铜学徒少年，锐利。。衣物细节：青铜学徒服。。发型妆造：短发。。脸型五官：少年面容，剑眉，眼含剑意，鼻梁挺，唇线坚毅。。武器招式：无兵器，折枝为剑。。功法：小宇宙初醒；剑心。。功法表现：剑意初凝。。画面：构图：摩羯宫前，青铜学徒少年折枝为剑，晨光中练剑忘我，背景金色宫阙。色调：青铜+剑意银+晨光金。氛围：剑心、忠诚、初志。。台词："我这一生，只做一件事——做一柄女神最锋利的剑。"。动作帧（动图）：①折枝为剑 ②晨光练剑 ③剑指苍穹 ④立誓。诗词：摩羯少年剑心纯，折枝为剑练晨昏。一生只奉一件事，为女铸得忠诚身。。主题句：一生为剑，剑之所向，即忠诚所往。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+剑意银+晨光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 奉剑之志, age 童年, 奉女神的剑心; scene: 捡一根树枝当剑，练到忘我；心里只有一个念头——为女神铸一柄忠诚的剑。; 摩羯少年剑心纯，折枝为剑练晨昏。一生只奉一件事，为女铸得忠诚身。; palette: 青铜+剑意银+晨光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 剑士修行**
- 摩羯座·修罗，凡尘砺心阶段·圣剑初芒（少年，Excalibur·圣剑初芒）。形象：摩羯座·修罗，银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。品性：Excalibur初芒，双臂凝剑，剑锋渐利。为剑而生的少年，锋芒毕露。。姿态：双臂凝作剑芒，一剑裂石；练剑房中，剑光如瀑。。服饰：青铜训练服，双臂剑芒。。体型：身高约6头身，青铜训练服，双臂剑芒。。衣物细节：青铜训练服。。发型妆造：短发。。脸型五官：少年面容，剑眉，眼含剑光。。武器招式：圣剑Excalibur初成。。功法：圣剑Excalibur初成；剑刃风暴初窥。。功法表现：剑芒裂石。。画面：构图：练剑房，青铜少年双臂凝作剑芒，一剑裂石，剑光如瀑，背景夜月。色调：青铜+剑芒银+夜月白。氛围：圣剑、初芒、剑士。。台词："我的双臂，就是剑。剑锋所指，唯忠诚二字。"。动作帧（动图）：①双臂凝剑 ②一剑裂石 ③剑光如瀑 ④收剑。诗词：圣剑初芒双臂凝，一剑裂石月夜惊。忠诚铸剑锋芒利，剑锋所指心自明。。主题句：一生为剑，剑之所向，即忠诚所往。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+剑芒银+夜月白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 圣剑初芒, age 少年, Excalibur·圣剑初芒; scene: 双臂凝作剑芒，一剑裂石；练剑房中，剑光如瀑。; 圣剑初芒双臂凝，一剑裂石月夜惊。忠诚铸剑锋芒利，剑锋所指心自明。; palette: 青铜+剑芒银+夜月白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 圣剑雏形**
- 摩羯座·修罗，道法初成阶段·圣剑裂空（少年，第七感·圣剑裂空）。形象：摩羯座·修罗，银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。品性：第七感觉醒，圣剑·裂空初成，一剑可断山河。忠诚之剑，越发锋锐。。姿态：圣剑裂空，一剑劈开山壁；忠诚之心，剑锋所指不偏不倚。。服饰：青铜圣衣具摩羯之相，双臂剑芒。。体型：身高约6头身，青铜摩羯圣衣。。衣物细节：青铜摩羯圣衣。。发型妆造：短发。。脸型五官：少年面容，剑眉，眼含剑意。。武器招式：圣剑·裂空。。功法：圣剑·裂空；剑刃风暴。。功法表现：剑光裂空。。画面：构图：山壁前，青铜摩羯圣衣少年圣剑裂空，一剑劈开山壁，剑光如虹，背景群山。色调：青铜+剑虹银+山青。氛围：第七感、裂空、忠诚。。台词："剑要快，心要正。我的剑，只为忠诚而挥。"。动作帧（动图）：①凝剑 ②圣剑裂空 ③劈开山壁 ④剑光如虹。诗词：第七感觉剑裂空，一剑劈开万山峰。忠诚为骨锋为刃，剑光所指正与忠。。主题句：一生为剑，剑之所向，即忠诚所往。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+剑虹银+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 圣剑裂空, age 少年, 第七感·圣剑裂空; scene: 圣剑裂空，一剑劈开山壁；忠诚之心，剑锋所指不偏不倚。; 第七感觉剑裂空，一剑劈开万山峰。忠诚为骨锋为刃，剑光所指正与忠。; palette: 青铜+剑虹银+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 剑影纵横**
- 摩羯座·修罗，大劫淬炼阶段·摩羯宫战（青年，以身为剑·守护圣域）。形象：摩羯座·修罗，银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。品性：圣战爆发，镇守摩羯宫。以身为剑的战士，剑锋所向，唯守护女神与圣域。。姿态：摩羯宫前，圣剑连斩；身受重创，犹以断剑之姿护住宫门。。服饰：白银摩羯圣衣，剑痕遍布。。体型：身高约7头身，白银摩羯圣衣，剑痕。。衣物细节：白银摩羯圣衣。。发型妆造：短发。。脸型五官：青年面容，剑眉，眼含剑意不屈。。武器招式：圣剑Excalibur·剑刃风暴。。功法：圣剑Excalibur·剑刃风暴；以身为剑。。功法表现：剑光如虹。。画面：构图：摩羯宫前圣战，白银摩羯圣斗士圣剑连斩，剑痕遍布，血染圣衣仍护宫门，背景战火。色调：白银+剑光银+血战红。氛围：以身为剑、守护、剑士。。台词："我的剑可以断，我的忠诚不会断。谁想过去，先问过我的剑。"。动作帧（动图）：①立于宫前 ②圣剑连斩 ③身受重创 ④剑断忠不断。诗词：摩羯宫前剑光寒，以身为剑镇雄关。剑断忠诚犹未断，血染宫门志不弯。。主题句：一生为剑，剑之所向，即忠诚所往。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：白银+剑光银+血战红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 摩羯宫战, age 青年, 以身为剑·守护圣域; scene: 摩羯宫前，圣剑连斩；身受重创，犹以断剑之姿护住宫门。; 摩羯宫前剑光寒，以身为剑镇雄关。剑断忠诚犹未断，血染宫门志不弯。; palette: 白银+剑光银+血战红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 摩羯宫主**
- 摩羯座·修罗，封神登天阶段·摩羯宫主（黄金圣斗士，黄金摩羯·忠诚剑士）。形象：摩羯座·修罗，银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。品性：受封黄金圣斗士，摩羯宫主。忠诚的剑士，圣剑Excalibur誓守女神。。姿态：端坐摩羯宫，双臂凝剑；来者过宫，先接圣剑之锋。。服饰：黄金摩羯圣衣，双臂剑芒如月。。体型：身高约7头身，黄金摩羯圣衣。。衣物细节：黄金摩羯圣衣。。发型妆造：短发。。脸型五官：青年面容，剑眉，眼含剑光。。武器招式：圣剑Excalibur（大成）。。功法：黄金摩羯之力；圣剑Excalibur大成。。功法表现：剑芒如月。。画面：构图：摩羯宫大殿，黄金摩羯圣斗士端坐，双臂剑芒如月，身后摩羯宫星象，背景金色殿堂。色调：黄金+剑芒银+殿堂金。氛围：黄金、忠诚、剑士。。台词："剑之所向，即忠诚所往。女神在，我的剑就在。"。动作帧（动图）：①端坐 ②双臂凝剑 ③来者接剑 ④俯瞰摩羯宫。诗词：黄金摩羯镇圣宫，圣剑一横万法空。忠诚为刃光如月，誓守女神立苍穹。。主题句：一生为剑，剑之所向，即忠诚所往。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+剑芒银+殿堂金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 摩羯宫主, age 黄金圣斗士, 黄金摩羯·忠诚剑士; scene: 端坐摩羯宫，双臂凝剑；来者过宫，先接圣剑之锋。; 黄金摩羯镇圣宫，圣剑一横万法空。忠诚为刃光如月，誓守女神立苍穹。; palette: 黄金+剑芒银+殿堂金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 觉醒之剑**
- 摩羯座·修罗，道果圆满阶段·以身归剑（终极，圣剑Excalibur·以身归剑）。形象：摩羯座·修罗，银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。品性：第八感领悟，圣剑Excalibur·以身归剑。一生为剑，终于与剑合一——人即剑，剑即忠。。姿态：双臂并作一柄通天圣剑，一剑化虹；以身为剑，护圣域万世。。服饰：神圣衣·摩羯（或黄金终极态），身化剑芒。。体型：身高约7头身，神圣衣摩羯，身化剑芒。。衣物细节：神圣衣·摩羯。。发型妆造：短发。。脸型五官：青年面容，剑眉，眼含剑魂。。武器招式：圣剑Excalibur·以身归剑。。功法：以身归剑；第八感；剑魂圆满。。功法表现：身化剑虹，贯穿天地。。画面：构图：摩羯宫顶，摩羯圣斗士双臂化作通天圣剑，一剑化虹贯穿天地，身化剑芒护住圣域，背景剑光星河。色调：剑芒银+神光金+圣域紫。氛围：终极、归剑、剑魂。。台词："这一生，我都在练剑。如今懂了——最好的剑，就是把自己铸成忠诚。"。动作帧（动图）：①双臂并剑 ②一剑化虹 ③以身归剑 ④护圣域。诗词：以身归剑摩羯魂，双臂化虹破乾坤。一生铸剑终成道，人剑合一护圣门。。主题句：一生为剑，剑之所向，即忠诚所往。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：剑芒银+神光金+圣域紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 以身归剑, age 终极, 圣剑Excalibur·以身归剑; scene: 双臂并作一柄通天圣剑，一剑化虹；以身为剑，护圣域万世。; 以身归剑摩羯魂，双臂化虹破乾坤。一生铸剑终成道，人剑合一护圣门。; palette: 剑芒银+神光金+圣域紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 水瓶座·卡妙（`aquarius`） · 人生档案版

**灵胎初醒 · 星辉冰晶**
- 水瓶座·卡妙，灵胎初醒阶段·绝对零度的冷静（童年，水瓶座·冷静少年）。形象：水瓶座·卡妙，冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。品性：圣域水瓶宫幼徒，天生冷静，心如冰晶，却藏着滚烫的心。。姿态：别人慌乱他沉静；冰晶在指尖凝结，他观察每一道裂纹。。服饰：青铜学徒服，眉目清冷。。体型：身高约5头身，青铜学徒少年，清冷。。衣物细节：青铜学徒服。。发型妆造：短发，清冷。。脸型五官：少年面容，剑眉，凤目清冷，鼻梁挺，唇线薄而静。。武器招式：无兵器，冰晶初凝。。功法：小宇宙初醒（寒性）；冰晶初凝。。功法表现：指尖冰晶。。画面：构图：水瓶宫水畔，青铜学徒少年指尖凝出冰晶，静看冰纹，背景清冷宫阙。色调：青铜+冰晶蓝+水光白。氛围：冷静、少年、初醒。。台词："冰要慢慢凝，心要时时静。我冷，是为了护得住热的东西。"。动作帧（动图）：①指尖凝冰 ②静看冰纹 ③凝神 ④微抬眸。诗词：水瓶少年性如冰，指尖晶凝照夜明。静看冰纹千条线，冷心藏着一腔情。。主题句：冰封外表的严师，藏着滚烫的守护之心。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+冰晶蓝+水光白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 绝对零度的冷静, age 童年, 水瓶座·冷静少年; scene: 别人慌乱他沉静；冰晶在指尖凝结，他观察每一道裂纹。; 水瓶少年性如冰，指尖晶凝照夜明。静看冰纹千条线，冷心藏着一腔情。; palette: 青铜+冰晶蓝+水光白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 冰道修行**
- 水瓶座·卡妙，凡尘砺心阶段·冰原严师（少年，冰原上的严师）。形象：水瓶座·卡妙，冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。品性：成为冰之战士的导师，在冰原上授徒。严苛如冰，心却为弟子滚烫。。姿态：冰原上教弟子练拳，毫不留情；弟子冻得发抖，他默默挡在风口。。服饰：青铜训练服，冰气绕身。。体型：身高约6头身，青铜训练服，冰气绕身。。衣物细节：青铜训练服。。发型妆造：短发。。脸型五官：少年面容，剑眉，凤目清冷，嘴角却有一丝温度。。武器招式：冰晶环。。功法：冰晶环；严师之道。。功法表现：冰气绕身。。画面：构图：冰原之上，青铜少年立于风雪中教弟子练拳，口鼻呵出白气，却默默挡在弟子风口，背景冰川。色调：青铜+冰原白+风雪蓝。氛围：严师、冰原、弟子。。台词："冰要承得住重，才护得住人。你们现在受的寒，将来都是护人的铠甲。"。动作帧（动图）：①冰原授徒 ②毫不留情 ③弟子发抖 ④默默挡风。诗词：冰原严师立风中，授徒拳法意从容。严苛如冰心似火，风口默默挡寒风。。主题句：冰封外表的严师，藏着滚烫的守护之心。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+冰原白+风雪蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 冰原严师, age 少年, 冰原上的严师; scene: 冰原上教弟子练拳，毫不留情；弟子冻得发抖，他默默挡在风口。; 冰原严师立风中，授徒拳法意从容。严苛如冰心似火，风口默默挡寒风。; palette: 青铜+冰原白+风雪蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冰晶环雏形**
- 水瓶座·卡妙，道法初成阶段·曙光女神（少年，第七感·曙光女神之宽恕）。形象：水瓶座·卡妙，冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。品性：第七感觉醒，曙光女神之宽恕初成，一拳可达绝对零度。。姿态：双拳并出，曙光女神之宽恕，冰封一片；绝对零度，万物凝霜。。服饰：青铜圣衣具水瓶之相，水光流转。。体型：身高约6头身，青铜水瓶圣衣。。衣物细节：青铜水瓶圣衣。。发型妆造：短发。。脸型五官：少年面容，剑眉，凤目如冰。。武器招式：曙光女神之宽恕。。功法：曙光女神之宽恕；绝对零度。。功法表现：冰封大地，绝对零度。。画面：构图：冰原夜，青铜水瓶圣衣少年双拳并出，曙光女神之宽恕冰封大地，绝对零度凝霜，背景寒星。色调：青铜+冰晶蓝+曙光银。氛围：第七感、曙光、初成。。台词："这一拳，是冰的宽恕——冻住罪恶，放过善良。"。动作帧（动图）：①双拳并出 ②冰封大地 ③绝对零度 ④凝霜而立。诗词：第七感觉曙光开，双拳并出冰封台。绝对零度凝万物，宽恕之心照雪来。。主题句：冰封外表的严师，藏着滚烫的守护之心。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+冰晶蓝+曙光银。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 曙光女神, age 少年, 第七感·曙光女神之宽恕; scene: 双拳并出，曙光女神之宽恕，冰封一片；绝对零度，万物凝霜。; 第七感觉曙光开，双拳并出冰封台。绝对零度凝万物，宽恕之心照雪来。; palette: 青铜+冰晶蓝+曙光银; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 绝对零度之道**
- 水瓶座·卡妙，大劫淬炼阶段·为弟子凝冰（青年，为弟子而凝的冰）。形象：水瓶座·卡妙，冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。品性：圣战之中，为护弟子，以身凝冰。冰封外表的严师，用冰守护了弟子滚烫的未来。。姿态：冰原上以身挡住致命一击，冰封自身；最后一拳，将弟子推向生路。。服饰：白银水瓶圣衣，凝霜。。体型：身高约7头身，白银水瓶圣衣，凝霜。。衣物细节：白银水瓶圣衣。。发型妆造：短发。。脸型五官：青年面容，凤目温柔，嘴角含笑。。武器招式：曙光女神之宽恕。。功法：曙光女神之宽恕；以身凝冰。。功法表现：以身凝冰，冰封自护。。画面：构图：冰原圣战，白银水瓶圣斗士以身凝冰挡住攻击，冰封自身，一手将弟子推向远方，背景风雪。色调：白银+冰晶蓝+血绛+风雪白。氛围：守护、凝冰、师徒。。台词："冰，是替你们挡在最前面的。我冻住自己，是为了你们还能暖。"。动作帧（动图）：①以身挡剑 ②冰封自身 ③推向弟子 ④含笑凝冰。诗词：圣战冰原凝此身，为徒挡剑血溅尘。冰封自己存弟子，严师之心热于春。。主题句：冰封外表的严师，藏着滚烫的守护之心。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：白银+冰晶蓝+血绛+风雪白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 为弟子凝冰, age 青年, 为弟子而凝的冰; scene: 冰原上以身挡住致命一击，冰封自身；最后一拳，将弟子推向生路。; 圣战冰原凝此身，为徒挡剑血溅尘。冰封自己存弟子，严师之心热于春。; palette: 白银+冰晶蓝+血绛+风雪白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 水瓶宫主**
- 水瓶座·卡妙，封神登天阶段·水瓶宫主（黄金圣斗士，黄金水瓶·冰之导师）。形象：水瓶座·卡妙，冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。品性：受封黄金圣斗士，水瓶宫主。冰之战士的导师，冷静严师，守护圣域。。姿态：端坐水瓶宫，水光流转；来者过宫，先尝绝对零度。。服饰：黄金水瓶圣衣，冰蓝水光。。体型：身高约7头身，黄金水瓶圣衣。。衣物细节：黄金水瓶圣衣。。发型妆造：短发。。脸型五官：青年面容，剑眉，凤目如冰含温。。武器招式：曙光女神之宽恕（大成）。。功法：黄金水瓶之力；曙光女神之宽恕大成。。功法表现：冰蓝水光。。画面：构图：水瓶宫大殿，黄金水瓶圣斗士端坐，冰蓝水光环绕，身后水瓶宫星象，背景金色殿堂。色调：黄金+冰蓝+殿堂金。氛围：黄金、导师、镇守。。台词："冰要守得住，才护得了。这水瓶宫，我用绝对零度来守。"。动作帧（动图）：①端坐 ②水光流转 ③绝对零度 ④俯瞰水瓶宫。诗词：黄金水瓶镇圣宫，水光流转冰气浓。严师冷面心如火，曙光之下守苍穹。。主题句：冰封外表的严师，藏着滚烫的守护之心。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+冰蓝+殿堂金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 水瓶宫主, age 黄金圣斗士, 黄金水瓶·冰之导师; scene: 端坐水瓶宫，水光流转；来者过宫，先尝绝对零度。; 黄金水瓶镇圣宫，水光流转冰气浓。严师冷面心如火，曙光之下守苍穹。; palette: 黄金+冰蓝+殿堂金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 冰河之师**
- 水瓶座·卡妙，道果圆满阶段·冰心无尘（终极，曙光女神·绝对零度）。形象：水瓶座·卡妙，冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。品性：第八感领悟，曙光女神·绝对零度圆满。冰封外表的严师，终成守护圣域之冰——冰心无尘，守护滚烫。。姿态：双拳并出，绝对零度化作曙光；冰封罪恶，暖护善良，圣域万物皆宁。。服饰：神圣衣·水瓶（或黄金终极态），冰晶为翼。。体型：身高约7头身，神圣衣水瓶，冰晶为翼。。衣物细节：神圣衣·水瓶。。发型妆造：短发。。脸型五官：青年面容，剑眉，凤目温柔含光。。武器招式：曙光女神·绝对零度。。功法：曙光女神·绝对零度；第八感。。功法表现：冰晶曙光，护佑圣域。。画面：构图：水瓶宫顶，水瓶圣斗士双拳并出，曙光女神·绝对零度化作漫天冰晶与曙光交织，护住圣域，背景极光冰川。色调：冰晶蓝+曙光银+极光青。氛围：终极、冰心、无尘。。台词："我冻了一辈子，才懂冰的温柔——冻住该冻的，放过该暖的。"。动作帧（动图）：①双拳并出 ②绝对零度 ③曙光交织 ④护圣域。诗词：冰心无尘水瓶尊，曙光绝对零度存。冻住罪恶放过暖，严师之魂护圣门。。主题句：冰封外表的严师，藏着滚烫的守护之心。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：冰晶蓝+曙光银+极光青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 冰心无尘, age 终极, 曙光女神·绝对零度; scene: 双拳并出，绝对零度化作曙光；冰封罪恶，暖护善良，圣域万物皆宁。; 冰心无尘水瓶尊，曙光绝对零度存。冻住罪恶放过暖，严师之魂护圣门。; palette: 冰晶蓝+曙光银+极光青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 双鱼座·阿布罗狄（`pisces`） · 人生档案版

**灵胎初醒 · 星辉花苞**
- 双鱼座·阿布罗狄，灵胎初醒阶段·玫瑰丛中（童年，双鱼座·玫瑰少年）。形象：双鱼座·阿布罗狄，绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。品性：圣域双鱼宫幼徒，玫瑰丛中的美战士，生的极美，却天生与玫瑰为伴。。姿态：在玫瑰园中长大，与玫瑰说话；指尖轻触花瓣，玫瑰应声绽放。。服饰：青铜学徒服，眉目如画。。体型：身高约5头身，青铜学徒少年，极美。。衣物细节：青铜学徒服。。发型妆造：长发披肩，美。。脸型五官：绝美面容，细眉，桃花眼含情，鼻梁挺，唇色如玫瑰。。武器招式：无兵器，玫瑰为伴。。功法：小宇宙初醒；玫瑰之缘。。功法表现：玫瑰随指绽放。。画面：构图：双鱼宫玫瑰园，青铜学徒美少年立于花丛，指尖触花，玫瑰绽放，背景绯红宫阙。色调：青铜+玫瑰绯红+园青。氛围：美、少年、初醒。。台词："玫瑰最美，也最护自己。我养玫瑰，玫瑰也护我。"。动作帧（动图）：①立花丛 ②指尖触花 ③玫瑰绽放 ④浅笑。诗词：双鱼少年貌如花，玫瑰园中伴朝霞。指尖轻触花应绽，美玉初成待芳华。。主题句：最美的玫瑰，藏着最致命的毒。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+玫瑰绯红+园青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 灵胎初醒 玫瑰丛中, age 童年, 双鱼座·玫瑰少年; scene: 在玫瑰园中长大，与玫瑰说话；指尖轻触花瓣，玫瑰应声绽放。; 双鱼少年貌如花，玫瑰园中伴朝霞。指尖轻触花应绽，美玉初成待芳华。; palette: 青铜+玫瑰绯红+园青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 花之修行**
- 双鱼座·阿布罗狄，凡尘砺心阶段·以玫瑰为兵（少年，以玫瑰为兵刃）。形象：双鱼座·阿布罗狄，绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。品性：以玫瑰为兵刃，毒与美集于一身。皇家恶魔玫瑰初练。。姿态：一挥手，漫天血玫飞射；毒与美并存，越美的花越致命。。服饰：青铜训练服，身绕玫瑰藤。。体型：身高约6头身，青铜训练服，玫瑰绕身。。衣物细节：青铜训练服。。发型妆造：长发。。脸型五官：绝美面容，桃花眼微眯。。武器招式：皇家恶魔玫瑰。。功法：皇家恶魔玫瑰初成；食人鱼玫瑰初窥。。功法表现：血玫漫天。。画面：构图：训练场，青铜少年一挥手漫天血玫飞射，身绕玫瑰藤，背景绯红。色调：青铜+玫瑰绯红+藤蔓绿。氛围：玫瑰、兵刃、美与毒。。台词："美，是武器；毒，是诚意。我这玫瑰，只送该收的人。"。动作帧（动图）：①抬手 ②漫天血玫 ③玫瑰飞射 ④浅笑。诗词：以玫瑰为刃少年行，漫天血玫化刀兵。毒与美集于一身，越美之花越无情。。主题句：最美的玫瑰，藏着最致命的毒。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+玫瑰绯红+藤蔓绿。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 凡尘砺心 以玫瑰为兵, age 少年, 以玫瑰为兵刃; scene: 一挥手，漫天血玫飞射；毒与美并存，越美的花越致命。; 以玫瑰为刃少年行，漫天血玫化刀兵。毒与美集于一身，越美之花越无情。; palette: 青铜+玫瑰绯红+藤蔓绿; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 恶魔玫瑰雏形**
- 双鱼座·阿布罗狄，道法初成阶段·皇家恶魔玫瑰（少年，第七感·皇家恶魔玫瑰）。形象：双鱼座·阿布罗狄，绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。品性：第七感觉醒，皇家恶魔玫瑰大成，毒与美，越发致命。。姿态：指尖捻起一朵恶魔玫瑰，轻轻吹出；玫瑰所至，黑暗尽染绯红。。服饰：青铜圣衣具双鱼之相，身绕绯红藤蔓。。体型：身高约6头身，青铜双鱼圣衣。。衣物细节：青铜双鱼圣衣。。发型妆造：长发。。脸型五官：绝美面容，桃花眼含媚。。武器招式：皇家恶魔玫瑰。。功法：皇家恶魔玫瑰；血腥玫瑰。。功法表现：花瓣染夜。。画面：构图：双鱼宫前，青铜双鱼圣衣少年指尖捻恶魔玫瑰轻吹，绯红花瓣飘散染暗夜，背景绯红月光。色调：青铜+恶魔玫红+夜紫。氛围：第七感、恶魔玫瑰、毒。。台词："恶魔玫瑰，是这世间最诚实的谎——它明明有毒，却美得让人想摘。"。动作帧（动图）：①捻玫瑰 ②轻吹 ③花瓣飘散 ④染暗夜。诗词：第七感觉恶魔玫，指尖轻捻绯红开。美到极处藏剧毒，双鱼之相立天台。。主题句：最美的玫瑰，藏着最致命的毒。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜+恶魔玫红+夜紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道法初成 皇家恶魔玫瑰, age 少年, 第七感·皇家恶魔玫瑰; scene: 指尖捻起一朵恶魔玫瑰，轻轻吹出；玫瑰所至，黑暗尽染绯红。; 第七感觉恶魔玫，指尖轻捻绯红开。美到极处藏剧毒，双鱼之相立天台。; palette: 青铜+恶魔玫红+夜紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 食人鱼玫瑰**
- 双鱼座·阿布罗狄，大劫淬炼阶段·血色终章（青年，双鱼宫战·血色玫瑰）。形象：双鱼座·阿布罗狄，绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。品性：圣战爆发，镇守双鱼宫。毒与美并存，血色玫瑰为守护圣域而绽。。姿态：双鱼宫前，黑玫瑰·吸魂齐飞；纵使血色染身，玫瑰依旧绝美。。服饰：白银双鱼圣衣，身绕血色藤蔓。。体型：身高约7头身，白银双鱼圣衣，血色染。。衣物细节：白银双鱼圣衣。。发型妆造：长发。。脸型五官：绝美面容，血色中依旧从容。。武器招式：黑玫瑰·吸魂。。功法：黑玫瑰·吸魂；血腥玫瑰连击。。功法表现：玫瑰漫天。。画面：构图：双鱼宫前圣战，白银双鱼圣斗士玫瑰漫天，血色染身仍绽笑，背景战火绯红宫阙。色调：白银+血色玫红+战火金。氛围：血色、守护、终章。。台词："我这一身血色，也是玫瑰的颜色。为了圣域，染红也无妨。"。动作帧（动图）：①立于宫前 ②玫瑰齐飞 ③血色染身 ④绽笑。诗词：双鱼宫前血色横，黑玫吸魂夜风惊。纵使染红身如玉，玫瑰依旧绽倾城。。主题句：最美的玫瑰，藏着最致命的毒。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：白银+血色玫红+战火金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 大劫淬炼 血色终章, age 青年, 双鱼宫战·血色玫瑰; scene: 双鱼宫前，黑玫瑰·吸魂齐飞；纵使血色染身，玫瑰依旧绝美。; 双鱼宫前血色横，黑玫吸魂夜风惊。纵使染红身如玉，玫瑰依旧绽倾城。; palette: 白银+血色玫红+战火金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 双鱼宫主**
- 双鱼座·阿布罗狄，封神登天阶段·双鱼宫主（黄金圣斗士，黄金双鱼·玫瑰宫主）。形象：双鱼座·阿布罗狄，绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。品性：受封黄金圣斗士，双鱼宫主。双鱼宫的芳香，藏着最锋利的毒。。姿态：端坐双鱼宫，满宫玫瑰盛开；来者过宫，先闻这致命的香。。服饰：黄金双鱼圣衣，身绕金色玫瑰藤。。体型：身高约7头身，黄金双鱼圣衣。。衣物细节：黄金双鱼圣衣。。发型妆造：长发。。脸型五官：绝美面容，桃花眼含笑。。武器招式：皇家恶魔玫瑰（大成）。。功法：黄金双鱼之力；皇家恶魔玫瑰大成。。功法表现：满宫玫瑰。。画面：构图：双鱼宫大殿，黄金双鱼圣斗士端坐，满宫玫瑰盛开，金色藤蔓绕身，身后双鱼宫星象，背景金色殿堂。色调：黄金+玫瑰金红+殿堂金。氛围：黄金、芳香、镇守。。台词："这满宫的香，是我给来客的礼数。香到深处，便是锋芒。"。动作帧（动图）：①端坐 ②玫瑰盛开 ③芳香弥漫 ④俯瞰双鱼宫。诗词：黄金双鱼镇圣宫，满殿玫瑰香正浓。芳香深处藏锋刃，美与毒并守苍穹。。主题句：最美的玫瑰，藏着最致命的毒。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金+玫瑰金红+殿堂金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 封神登天 双鱼宫主, age 黄金圣斗士, 黄金双鱼·玫瑰宫主; scene: 端坐双鱼宫，满宫玫瑰盛开；来者过宫，先闻这致命的香。; 黄金双鱼镇圣宫，满殿玫瑰香正浓。芳香深处藏锋刃，美与毒并守苍穹。; palette: 黄金+玫瑰金红+殿堂金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 凄美绽放**
- 双鱼座·阿布罗狄，道果圆满阶段·玫瑰之殇（终极，玫瑰之殇·血色终章）。形象：双鱼座·阿布罗狄，绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。品性：第八感领悟，玫瑰之殇·血色终章。最美的玫瑰，守护到最后一瓣凋落——毒与美，皆成守护。。姿态：一挥手，漫天玫瑰化作玫瑰之殇，血色终章笼罩；护住圣域，任花瓣零落成泥。。服饰：神圣衣·双鱼（或黄金终极态），玫瑰为翼。。体型：身高约7头身，神圣衣双鱼，玫瑰为翼。。衣物细节：神圣衣·双鱼。。发型妆造：长发。。脸型五官：绝美面容，桃花眼含慈悲。。武器招式：玫瑰之殇·血色终章。。功法：玫瑰之殇·血色终章；第八感。。功法表现：漫天玫瑰护圣域。。画面：构图：双鱼宫顶，双鱼圣斗士一挥手漫天玫瑰化作玫瑰之殇，血色终章笼罩圣域，花瓣零落成泥仍留香，背景绯红星夜。色调：玫瑰绯红+血色深红+神光金。氛围：终极、玫瑰、血色。。台词："玫瑰谢了，毒还在。我用最后一瓣，为圣域开出一片安全的花园。"。动作帧（动图）：①抬手 ②玫瑰之殇 ③血色终章 ④零落成泥。诗词：玫瑰之殇双鱼魂，血色终章照圣门。毒与美皆成守护，零落成泥香犹存。。主题句：最美的玫瑰，藏着最致命的毒。。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：玫瑰绯红+血色深红+神光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; 道果圆满 玫瑰之殇, age 终极, 玫瑰之殇·血色终章; scene: 一挥手，漫天玫瑰化作玫瑰之殇，血色终章笼罩；护住圣域，任花瓣零落成泥。; 玫瑰之殇双鱼魂，血色终章照圣门。毒与美皆成守护，零落成泥香犹存。; palette: 玫瑰绯红+血色深红+神光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 9. 传统节日（12 物种）

> **风格**：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。**阶段演绎**：
> - 灵胎初醒：原料初成，厨房温暖，节日的前奏与期待（原料暖米/糯粉）
> - 凡尘砺心：制作过程中的形态，蒸煮捏塑，巧手匠心（制作棕/喜庆红）
> - 道法初成：节庆登场，寓意成形，喜庆氛围渐浓（节日红/麦金）
> - 大劫淬炼：阖家团圆的时刻，热气腾腾，温暖满溢（团圆红/蒸金）
> - 封神登天：成为节日主角与名品，家家户户的吉祥符号（主角红/名品金）
> - 道果圆满：化为民俗文化图腾，代代相传，温暖永续（民俗暖白/传承金）

#### 粽子（`zongzi`）

**灵胎初醒 · 箬叶籽**
- 粽子，灵胎初醒阶段·箬叶籽。初始形态：一粒箬叶包裹的糯米灵种，青绿叶片系着绳结，米香与端午的艾草气萦绕。木属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 糯米团**
- 粽子，凡尘砺心阶段·糯米团。形象：翠绿粽形，系着绳结。 核心意象：箬叶、绳结、端午龙舟。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 蜜枣粽**
- 粽子，道法初成阶段·蜜枣粽。形象：翠绿粽形，系着绳结。 核心意象：箬叶、绳结、端午龙舟。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，剥开箬叶，露出饱满晶莹的糯米。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 龙舟粽**
- 粽子，大劫淬炼阶段·龙舟粽。形象：翠绿粽形，系着绳结。 核心意象：箬叶、绳结、端午龙舟。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 五月粽**
- 粽子，封神登天阶段·五月粽。形象：翠绿粽形，系着绳结。 核心意象：箬叶、绳结、端午龙舟。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 屈子之志**
- 粽子，道果圆满阶段·屈子之志。形象：翠绿粽形，系着绳结。 核心意象：箬叶、绳结、端午龙舟。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，剥开箬叶，露出饱满晶莹的糯米。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 汤圆（`tangyuan`）

**灵胎初醒 · 糯米团**
- 汤圆，灵胎初醒阶段·糯米团。初始形态：一团雪白糯米团，软糯圆润浮在甜汤灵光里，流心芝麻的甜意在团中酝酿。水属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小圆子**
- 汤圆，凡尘砺心阶段·小圆子。形象：雪白圆润，浮于甜汤。 核心意象：雪白圆身、甜汤、团圆之月。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 流心芝麻**
- 汤圆，道法初成阶段·流心芝麻。形象：雪白圆润，浮于甜汤。 核心意象：雪白圆身、甜汤、团圆之月。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，在甜汤里轻轻翻滚，一口咬下软糯流心。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 元宵**
- 汤圆，大劫淬炼阶段·元宵。形象：雪白圆润，浮于甜汤。 核心意象：雪白圆身、甜汤、团圆之月。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 元夕汤圆**
- 汤圆，封神登天阶段·元夕汤圆。形象：雪白圆润，浮于甜汤。 核心意象：雪白圆身、甜汤、团圆之月。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 玉润尊者**
- 汤圆，道果圆满阶段·玉润尊者。形象：雪白圆润，浮于甜汤。 核心意象：雪白圆身、甜汤、团圆之月。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，在甜汤里轻轻翻滚，一口咬下软糯流心。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 月饼（`mooncake`）

**灵胎初醒 · 面团胚**
- 月饼，灵胎初醒阶段·面团胚。初始形态：一块面团胚，金黄饼形压着花印纹，满月般的圆满寓意在胚中成形。金属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小月饼**
- 月饼，凡尘砺心阶段·小月饼。形象：棕金饼形，压花纹章。 核心意象：饼上花印、满月、玉兔。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 莲蓉月饼**
- 月饼，道法初成阶段·莲蓉月饼。形象：棕金饼形，压花纹章。 核心意象：饼上花印、满月、玉兔。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，刀锋落下，露出咸蛋黄流心。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 冰皮流心**
- 月饼，大劫淬炼阶段·冰皮流心。形象：棕金饼形，压花纹章。 核心意象：饼上花印、满月、玉兔。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 中秋月饼**
- 月饼，封神登天阶段·中秋月饼。形象：棕金饼形，压花纹章。 核心意象：饼上花印、满月、玉兔。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 明月尊者**
- 月饼，道果圆满阶段·明月尊者。形象：棕金饼形，压花纹章。 核心意象：饼上花印、满月、玉兔。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，刀锋落下，露出咸蛋黄流心。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 青团（`qingtuan`）

**灵胎初醒 · 艾草芽**
- 青团，灵胎初醒阶段·艾草芽。初始形态：一枚艾草芽，碧绿糯米团裹着青碧草色，清明踏青的清香在芽尖浮动。木属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小青团**
- 青团，凡尘砺心阶段·小青团。形象：碧绿糯米团，艾草清香。 核心意象：艾草青碧、糯香、清明踏青。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 艾香团**
- 青团，道法初成阶段·艾香团。形象：碧绿糯米团，艾草清香。 核心意象：艾草青碧、糯香、清明踏青。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，揭开蒸笼，艾草清香扑面而来。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 清明团**
- 青团，大劫淬炼阶段·清明团。形象：碧绿糯米团，艾草清香。 核心意象：艾草青碧、糯香、清明踏青。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 春意青团**
- 青团，封神登天阶段·春意青团。形象：碧绿糯米团，艾草清香。 核心意象：艾草青碧、糯香、清明踏青。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 青翠尊者**
- 青团，道果圆满阶段·青翠尊者。形象：碧绿糯米团，艾草清香。 核心意象：艾草青碧、糯香、清明踏青。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，揭开蒸笼，艾草清香扑面而来。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 重阳糕（`chongyang_cake`）

**灵胎初醒 · 米粉种**
- 重阳糕，灵胎初醒阶段·米粉种。初始形态：一粒米粉种，金黄米塔初具层级，重阳登高的茱萸香气在粉粒间萦绕。金属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小糕塔**
- 重阳糕，凡尘砺心阶段·小糕塔。形象：金色米糕，层叠如塔。 核心意象：层叠米糕、登高茱萸、金秋。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 枣泥糕**
- 重阳糕，道法初成阶段·枣泥糕。形象：金色米糕，层叠如塔。 核心意象：层叠米糕、登高茱萸、金秋。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，层层剥开如塔，撒上红绿果脯。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 步步糕**
- 重阳糕，大劫淬炼阶段·步步糕。形象：金色米糕，层叠如塔。 核心意象：层叠米糕、登高茱萸、金秋。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 福寿糕**
- 重阳糕，封神登天阶段·福寿糕。形象：金色米糕，层叠如塔。 核心意象：层叠米糕、登高茱萸、金秋。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 登高尊者**
- 重阳糕，道果圆满阶段·登高尊者。形象：金色米糕，层叠如塔。 核心意象：层叠米糕、登高茱萸、金秋。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，层层剥开如塔，撒上红绿果脯。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 年糕（`niangao`）

**灵胎初醒 · 糯米砖**
- 年糕，灵胎初醒阶段·糯米砖。初始形态：一块糯米砖，莹白软糯的方砖雏形，年年高升的红纸光晕在砖面浮动。金属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 软年糕**
- 年糕，凡尘砺心阶段·软年糕。形象：白糯方糕，软糯弹牙。 核心意象：白糯方糕、新春对联、年年高升。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 打年糕**
- 年糕，道法初成阶段·打年糕。形象：白糯方糕，软糯弹牙。 核心意象：白糯方糕、新春对联、年年高升。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，油锅煎至金黄，外脆里糯夹起拉丝。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 步步高**
- 年糕，大劫淬炼阶段·步步高。形象：白糯方糕，软糯弹牙。 核心意象：白糯方糕、新春对联、年年高升。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 福禄年糕**
- 年糕，封神登天阶段·福禄年糕。形象：白糯方糕，软糯弹牙。 核心意象：白糯方糕、新春对联、年年高升。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 年糕尊者**
- 年糕，道果圆满阶段·年糕尊者。形象：白糯方糕，软糯弹牙。 核心意象：白糯方糕、新春对联、年年高升。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，油锅煎至金黄，外脆里糯夹起拉丝。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 腊八粥（`laba_porridge`）

**灵胎初醒 · 谷粒籽**
- 腊八粥，灵胎初醒阶段·谷粒籽。初始形态：一把谷粒籽，红豆桂圆莲子红枣五谷杂陈，腊八的暖意在百味中酝酿。土属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 五谷粥**
- 腊八粥，凡尘砺心阶段·五谷粥。形象：五谷杂粮熬成的浓粥。 核心意象：五谷杂粮、腊八蒜、暖冬。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 佛粥**
- 腊八粥，道法初成阶段·佛粥。形象：五谷杂粮熬成的浓粥。 核心意象：五谷杂粮、腊八蒜、暖冬。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，舀起一勺，红枣桂圆莲子浮沉其间。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 百谷粥**
- 腊八粥，大劫淬炼阶段·百谷粥。形象：五谷杂粮熬成的浓粥。 核心意象：五谷杂粮、腊八蒜、暖冬。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 融和大师**
- 腊八粥，封神登天阶段·融和大师。形象：五谷杂粮熬成的浓粥。 核心意象：五谷杂粮、腊八蒜、暖冬。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 百味尊者**
- 腊八粥，道果圆满阶段·百味尊者。形象：五谷杂粮熬成的浓粥。 核心意象：五谷杂粮、腊八蒜、暖冬。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，舀起一勺，红枣桂圆莲子浮沉其间。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 春饼（`spring_pancake`）

**灵胎初醒 · 面糊团**
- 春饼，灵胎初醒阶段·面糊团。初始形态：一团面糊，薄嫩面胚泛着麦香，春蔬的嫩绿在饼皮中透出咬春之意。木属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 烙春饼**
- 春饼，凡尘砺心阶段·烙春饼。形象：薄嫩面饼，卷着春蔬。 核心意象：薄饼、春蔬、立春之绿。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 咬春饼**
- 春饼，道法初成阶段·咬春饼。形象：薄嫩面饼，卷着春蔬。 核心意象：薄饼、春蔬、立春之绿。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，摊开薄饼，卷入时令春蔬一卷而食。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 薄脆饼**
- 春饼，大劫淬炼阶段·薄脆饼。形象：薄嫩面饼，卷着春蔬。 核心意象：薄饼、春蔬、立春之绿。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 咬春大师**
- 春饼，封神登天阶段·咬春大师。形象：薄嫩面饼，卷着春蔬。 核心意象：薄饼、春蔬、立春之绿。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 春味尊者**
- 春饼，道果圆满阶段·春味尊者。形象：薄嫩面饼，卷着春蔬。 核心意象：薄饼、春蔬、立春之绿。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，摊开薄饼，卷入时令春蔬一卷而食。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 冰糖葫芦（`tanghulu`）

**灵胎初醒 · 山楂果**
- 冰糖葫芦，灵胎初醒阶段·山楂果。初始形态：一粒山楂果，红亮圆润裹着初凝的冰糖衣，街头的酸甜在脆壳下待绽。火属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小糖葫芦**
- 冰糖葫芦，凡尘砺心阶段·小糖葫芦。形象：冰糖山楂串，红亮晶莹。 核心意象：红亮山楂、冰糖衣、竹签串。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冰糖葫芦**
- 冰糖葫芦，道法初成阶段·冰糖葫芦。形象：冰糖山楂串，红亮晶莹。 核心意象：红亮山楂、冰糖衣、竹签串。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，举着竹签一咬，"咔嚓"咬开糖衣。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 糖葫芦串**
- 冰糖葫芦，大劫淬炼阶段·糖葫芦串。形象：冰糖山楂串，红亮晶莹。 核心意象：红亮山楂、冰糖衣、竹签串。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 甜蜜大师**
- 冰糖葫芦，封神登天阶段·甜蜜大师。形象：冰糖山楂串，红亮晶莹。 核心意象：红亮山楂、冰糖衣、竹签串。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 酸甜尊者**
- 冰糖葫芦，道果圆满阶段·酸甜尊者。形象：冰糖山楂串，红亮晶莹。 核心意象：红亮山楂、冰糖衣、竹签串。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，举着竹签一咬，"咔嚓"咬开糖衣。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 桂花糕（`osmanthus_cake`）

**灵胎初醒 · 桂子籽**
- 桂花糕，灵胎初醒阶段·桂子籽。初始形态：一粒桂子籽，米白糕胚缀着金桂花瓣，金秋的桂香在甜意中酝酿。金属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 桂花糕**
- 桂花糕，凡尘砺心阶段·桂花糕。形象：米白方糕，点缀桂花瓣。 核心意象：米白方糕、金桂花瓣、秋香。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 金桂糕**
- 桂花糕，道法初成阶段·金桂糕。形象：米白方糕，点缀桂花瓣。 核心意象：米白方糕、金桂花瓣、秋香。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，拈起一块，桂香在齿间化开。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 月桂糕**
- 桂花糕，大劫淬炼阶段·月桂糕。形象：米白方糕，点缀桂花瓣。 核心意象：米白方糕、金桂花瓣、秋香。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金桂之王**
- 桂花糕，封神登天阶段·金桂之王。形象：米白方糕，点缀桂花瓣。 核心意象：米白方糕、金桂花瓣、秋香。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 桂香尊者**
- 桂花糕，道果圆满阶段·桂香尊者。形象：米白方糕，点缀桂花瓣。 核心意象：米白方糕、金桂花瓣、秋香。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，拈起一块，桂香在齿间化开。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 馄饨（`wonton`）

**灵胎初醒 · 面皮种**
- 馄饨，灵胎初醒阶段·面皮种。初始形态：一张面皮种，薄如云翳的面胚裹着鲜肉初形，热汤的暖意在褶皱中藏起。水属性灵光微微环绕。神态：原料的宁静，蕴着期待。动作：静置案板，等待巧手。衣着：米面/食材原料，朴实。梳造：无，原料形态。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 原料的宁静，蕴着期待; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 云吞**
- 馄饨，凡尘砺心阶段·云吞。形象：薄皮小馄饨，汤鲜味美。 核心意象：薄皮、热汤、虾皮紫菜。神态：制作中的专注。动作：揉捏塑形，蒸汽初升。衣着：半成品，形态渐显。梳造：裹叶/印模初成。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 制作中的专注; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 元宝云吞**
- 馄饨，道法初成阶段·元宝云吞。形象：薄皮小馄饨，汤鲜味美。 核心意象：薄皮、热汤、虾皮紫菜。神态：成形时的欢喜。动作：摆上蒸笼，静待火候，一勺舀起，薄皮裹馅在热汤里打个转。衣着：形态完整，色泽初显。梳造：印花/装饰点缀。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 成形时的欢喜; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 鲜肉馄饨**
- 馄饨，大劫淬炼阶段·鲜肉馄饨。形象：薄皮小馄饨，汤鲜味美。 核心意象：薄皮、热汤、虾皮紫菜。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大厨馄饨**
- 馄饨，封神登天阶段·大厨馄饨。形象：薄皮小馄饨，汤鲜味美。 核心意象：薄皮、热汤、虾皮紫菜。神态：出锅封名，暖意融融。动作：热气腾腾，登堂亮相。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅封名，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 馄饨尊者**
- 馄饨，道果圆满阶段·馄饨尊者。形象：薄皮小馄饨，汤鲜味美。 核心意象：薄皮、热汤、虾皮紫菜。神态：名品镇席，香飘万里。动作：香溢满座，百年传承，一勺舀起，薄皮裹馅在热汤里打个转。衣着：华美名品，灵光透色。梳造：宝光流转，镇世之味。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品镇席，香飘万里; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 花灯（`festival_lantern`）

**灵胎初醒 · 竹骨灯**
- 花灯，灵胎初醒阶段·竹骨灯。初始形态：一盏竹骨灯，素绢灯面初糊成胚，烛火微光在骨架上跳动，灯谜的期待在绢面浮现。火属性灵光微微环绕。神态：沉睡的灵光，物灵未醒。动作：静置无声，灵光内蕴。衣着：素坯/原石，未成形。梳造：无，器物胚形。意境：原料初成，厨房温暖，节日的前奏与期待。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：原料暖米（#F5E6D3）主调 + 糯粉（#F8B195）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festive prelude; 沉睡的灵光，物灵未醒; palette #F5E6D3 with #F8B195 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 花灯仔**
- 花灯，凡尘砺心阶段·花灯仔。形象：红彩花灯，烛光摇曳。 核心意象：红彩绢纱、烛火、灯谜。神态：初生灵智的好奇。动作：微光颤动，器灵初醒。衣着：初雕成型，轮廓渐显。梳造：雕纹/铭文初刻。意境：制作过程中的形态，蒸煮捏塑，巧手匠心。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：制作棕（#E8A87C）主调 + 喜庆红（#F67280）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; crafting warmth; 初生灵智的好奇; palette #E8A87C with #F67280 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 上元灯**
- 花灯，道法初成阶段·上元灯。形象：红彩花灯，烛光摇曳。 核心意象：红彩绢纱、烛火、灯谜。神态：灵光渐盛，灵动自生。动作：器光初现，灵光流转，烛火摇曳，光影在绢纱上流转。衣着：成形精工，纹饰渐繁。梳造：纹饰/嵌饰增辉。意境：节庆登场，寓意成形，喜庆氛围渐浓。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：节日红（#F67280）主调 + 麦金（#FFD54F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; festival debut; 灵光渐盛，灵动自生; palette #F67280 with #FFD54F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 宫灯**
- 花灯，大劫淬炼阶段·宫灯。形象：红彩花灯，烛光摇曳。 核心意象：红彩绢纱、烛火、灯谜。神态：淬炼中的忍耐。动作：全力受炼，火炼淬洗。衣着：历经淬炼，温润内敛。梳造：包浆/裂纹，岁月痕。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 淬炼中的忍耐; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 灯会主角**
- 花灯，封神登天阶段·灯会主角。形象：红彩花灯，烛光摇曳。 核心意象：红彩绢纱、烛火、灯谜。神态：受封灵宝，灵光大成的庄严。动作：器灵大成，光华夺目。衣着：华彩流光，灵气充盈。梳造：祥纹满饰，宝光外放。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 受封灵宝，灵光大成的庄严; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 灯影尊者**
- 花灯，道果圆满阶段·灯影尊者。形象：红彩花灯，烛光摇曳。 核心意象：红彩绢纱、烛火、灯谜。神态：器灵化神，灵性通明。动作：器道通天，光华耀世，烛火摇曳，光影在绢纱上流转。衣着：神纹流转，镇世之宝。梳造：灵光化形，器魂永驻。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 器灵化神，灵性通明; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 10. 虹猫蓝兔七侠传（10 物种）

> **风格**：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。**阶段演绎**：
> - 灵胎初醒：剑意种子萌芽，江湖初闻，少年侠气的微光（剑意青/银灰）
> - 凡尘砺心：习武入门，剑招初练，磨砺锋芒（侠义蓝/热血红）
> - 道法初成：剑法初成，行走江湖，惩恶扬善（江湖橙/侠气蓝）
> - 大劫淬炼：七剑合璧，魔教大战，刀光剑影中淬炼（魔教紫/血战红）
> - 封神登天：七侠之名扬天下，快意恩仇，豪情万丈（七剑金/侠名白）
> - 道果圆满：剑道圆满，一剑开天，守护苍生的宗师气象（剑道白/圆满蓝）

#### 虹猫（`hongmao`） · 人生档案版

**灵胎初醒 · 剑意种子**
- 虹猫，灵胎初醒阶段·剑意萌芽（幼年，长虹剑主·家园太平）。形象：虹猫，红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。品性：长虹剑传人，七剑之首。幼年与长虹剑为伴，家园太平，一颗赤子之心，不知江湖险恶。。姿态：山林中习剑，长虹剑带起七色虹光；父亲含笑而立，告诉他——"剑，是护家的。"。服饰：红毛小猫，素色劲装。。体型：身高约4头身，红毛小猫，身形灵巧。。衣物细节：红布短打，袖口利落，腰间系红带，长虹剑佩于身侧。。发型妆造：一头红毛，发间系着红色发带，耳朵尖立。。脸型五官：猫脸，圆眼有神，眉间英气，嘴角含笑。。武器招式：长虹剑初练。。功法：长虹剑法初窥；剑意萌芽。。功法表现：剑光七色虹。。画面：构图：家园晨光，红毛小猫持长虹剑习练，剑光带起七色虹光，父亲在旁含笑，背景青山家园。色调：红毛+长虹七色+晨光金。氛围：太平、赤子、不谙世事。。台词："爹说，长虹剑是护家园的。那时候我还不懂，什么叫守护。"。动作帧（动图）：①持剑 ②剑光带虹 ③父亲含笑 ④望家园。诗词：红毛小猫抱剑眠，家园太平岁月甜。赤子不知江湖险，只道长虹护故园。。主题句：家破志不灭，长虹映月明——忍辱负重心中有光，赤子之心永不改。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：红毛+长虹七色+晨光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 剑意萌芽, age 幼年, 长虹剑主·家园太平; scene: 山林中习剑，长虹剑带起七色虹光；父亲含笑而立，告诉他——"剑，是护家的。"; 红毛小猫抱剑眠，家园太平岁月甜。赤子不知江湖险，只道长虹护故园。; palette: 红毛+长虹七色+晨光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 红衣小猫**
- 虹猫，凡尘砺心阶段·家园遭难（幼年，魔教来袭·父亲托付）。形象：虹猫，红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。品性：魔教突袭，家园被毁，火光冲天。父亲以身为盾护他逃生，临终将长虹剑与复兴七剑之任托付于他——"爹最后的目光是安详的，因为他知道，儿子会长大，会把这片家重新点亮。"。姿态：火光中父亲推他上马，塞给他长虹剑，低声道"爹把光留给你，你会替爹点亮它"，转身迎向魔教；他回首，家园成灰，却把父亲的剑紧紧抱在怀里。。服饰：红毛小猫，身上带伤，怀抱长虹剑。。体型：身高约4头身，红毛小猫，带伤抱剑。。衣物细节：红布短打染尘带伤，腰间红带，紧抱长虹剑。。发型妆造：红毛凌乱，沾着灰，发带松脱。。脸型五官：猫脸，圆眼含泪，眉间悲恸。。武器招式：长虹剑（护身）。。功法：长虹剑法初成；背负父亲托付的希望。。功法表现：火光映长虹。。画面：构图：火光家园，红毛小猫怀抱长虹剑立于废墟前，父亲身影没入火光却回头一笑，长虹剑泛着温暖的光，背景冲天烈焰。色调：火光红+废墟黑+长虹暖金。氛围：家破之痛、父亲的托付与希望。。台词："爹，我记下了——您说这剑，是留给我的光。我会好好长大，替您把它点亮。"。动作帧（动图）：①火光中回头 ②父亲推他上马 ③家园成灰 ④抱剑垂泪。诗词：魔火焚家夜正狂，父亲挡剑护儿郎。临别塞剑留一语：此光留你代我亮。。主题句：家破志不灭，长虹映月明——忍辱负重心中有光，赤子之心永不改。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：火光红+废墟黑+长虹暖金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 家园遭难, age 幼年, 魔教来袭·父亲托付; scene: 火光中父亲推他上马，塞给他长虹剑，低声道"爹把光留给你，你会替爹点亮它"，转身迎向魔教；他回首，家园成灰，却把父亲的剑紧紧抱在怀里。; 魔火焚家夜正狂，父亲挡剑护儿郎。临别塞剑留一语：此光留你代我亮。; palette: 火光红+废墟黑+长虹暖金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 长虹剑出**
- 虹猫，道法初成阶段·寻剑苦修（少年，独行寻剑·剑随心动）。形象：虹猫，红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。品性：为复兴七剑，独自踏上寻剑之路。长虹剑法在父亲手中出神入化，到他手里却迟迟不成。他背着父亲的长虹剑，一路被魔教追杀、被人嘲笑"黄毛小儿"，却始终把父亲的剑贴身带着——那是他唯一还能触碰到的、父亲留下的东西。。姿态：深山老林里，他一遍遍练长虹贯日，剑尖挑不起山泉，练到深夜；一个赶路的老者随口指点"剑随心动"，他记了一整夜；雪夜里，他把长虹剑放在膝上，低声对剑说："爹，您教我练剑那会儿，也是这么慢的吗？"。服饰：红毛小猫，粗布劲装，长虹剑时刻在手。。体型：身高约5头身，红毛小猫，粗布劲装。。衣物细节：粗布劲装，腰间红带磨旧，长虹剑时刻在手，腕上缠着练剑的布条。。发型妆造：红毛束起，发带褪色，沾尘。。脸型五官：猫脸，圆眼坚毅，眉间风霜。。武器招式：长虹剑（苦练中）。。功法：长虹剑法苦练渐精；寻剑之路。。功法表现：剑光微亮于月夜。。画面：构图：深山雪夜，红毛小猫独坐石上，长虹剑横膝，对剑低语，远处老者指点之影，背景苍茫群山与月光。色调：红毛+月白+山青。氛围：苦修、孤独、有寄托。。台词："剑，不是一两天能练成的。爹练了半辈子，我急什么？我且慢慢练，练到它认我为主——也练到，我能配得上它。"。动作帧（动图）：①练长虹贯日 ②剑尖挑泉屡败 ③记败招 ④夜中对着月光再练。诗词：独行寻剑少年孤，长虹未成苦练初。剑随心动有人点，雪夜对剑问阿爹。。主题句：家破志不灭，长虹映月明——忍辱负重心中有光，赤子之心永不改。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：红毛+月白+山青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 寻剑苦修, age 少年, 独行寻剑·剑随心动; scene: 深山老林里，他一遍遍练长虹贯日，剑尖挑不起山泉，练到深夜；一个赶路的老者随口指点"剑随心动"，他记了一整夜；雪夜里，他把长虹剑放在膝上，低声对剑说："爹，您教我练剑那会儿，也是这么慢的吗？"; 独行寻剑少年孤，长虹未成苦练初。剑随心动有人点，雪夜对剑问阿爹。; palette: 红毛+月白+山青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 侠义之心**
- 虹猫，大劫淬炼阶段·剑成初试（青年，剑法小成·情恨交织）。形象：虹猫，红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。品性：长虹剑法小成，一剑贯日初成。寻剑途中初遇蓝兔，并肩初战魔教——第一次，他用长虹剑护住了别人，也第一次，心里有了一个想守护的人。。姿态：与蓝兔并肩，长虹冰魄双剑合璧初显，退敌成功；他望着蓝兔，又望了一眼魔教方向——对家的恨、对眼前人的情，一起在心里亮起来。。服饰：红毛小猫，侠客装，长虹剑在手。。体型：身高约5头身，红毛小猫，侠客装。。衣物细节：红色侠客劲装，袖口挽起，腰间长虹剑，腕缠护腕。。发型妆造：红毛束起，发带鲜亮。。脸型五官：猫脸，圆眼含笑，剑意初亮。。武器招式：长虹贯日（小成）。。功法：长虹贯日小成；初与蓝兔并肩。。功法表现：长虹初亮，冰魄交映。。画面：构图：初战之地，红毛小猫长虹剑小成，与蓝兔并肩退敌，长虹与冰魄交映，一侧魔教远影，背景初战烽火。色调：红毛+长虹金+冰魄蓝+魔教暗。氛围：小成、情恨交织、看见希望。。台词："爹，我的剑开始能护人了。我会护住这片家，也护住眼前这个并肩的人——黑心虎欠的账，我也记着。"。动作帧（动图）：①与蓝兔并肩 ②长虹冰魄合璧 ③退敌 ④望剑而笑。诗词：剑成初试贯日横，与侠并肩退敌兵。家恨未销情初起，长虹与冰魄交明。。主题句：家破志不灭，长虹映月明——忍辱负重心中有光，赤子之心永不改。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：红毛+长虹金+冰魄蓝+魔教暗。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 剑成初试, age 青年, 剑法小成·情恨交织; scene: 与蓝兔并肩，长虹冰魄双剑合璧初显，退敌成功；他望着蓝兔，又望了一眼魔教方向——对家的恨、对眼前人的情，一起在心里亮起来。; 剑成初试贯日横，与侠并肩退敌兵。家恨未销情初起，长虹与冰魄交明。; palette: 红毛+长虹金+冰魄蓝+魔教暗; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 七剑之首**
- 虹猫，封神登天阶段·长虹贯日（青年，长虹大成·一剑扬威）。形象：虹猫，红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。品性：长虹剑法大成，一剑贯日。决战之前，他凭这一剑重创魔教先锋，威震江湖——父亲传下的剑法，终于在他手中圆满。。姿态：魔教先锋来犯，他一剑长虹贯日，剑光贯日而落，敌阵溃散；望着剑光，他轻声道："爹，您的长虹剑，我会了。"。服饰：红毛小猫，长虹剑在手，七侠之首气度。。体型：身高约6头身，红毛小猫，长虹剑在手。。衣物细节：红色侠客装，腰间长虹剑，护腕革带，七侠之首气度。。发型妆造：红毛束起，发带迎风。。脸型五官：猫脸，圆眼含光，剑意圆满。。武器招式：长虹贯日（圆满）。。功法：长虹贯日圆满；一剑扬威。。功法表现：剑光贯日，威震敌阵。。画面：构图：战场，红毛小猫一剑长虹贯日，剑光贯日而落重创敌阵，威震四方，背景决战前奏。色调：红毛+长虹金+战火红。氛围：大成、扬威、初显锋芒。。台词："爹，您传我的长虹贯日，今日终于圆满。这一剑，够魔教喝一壶了！"。动作帧（动图）：①蓄剑 ②长虹贯日 ③剑光贯日而落 ④重创敌阵。诗词：长虹贯日圆满成，一剑扬威震魔营。七侠之首初显威，长虹之光照天明。。主题句：家破志不灭，长虹映月明——忍辱负重心中有光，赤子之心永不改。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：红毛+长虹金+战火红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 长虹贯日, age 青年, 长虹大成·一剑扬威; scene: 魔教先锋来犯，他一剑长虹贯日，剑光贯日而落，敌阵溃散；望着剑光，他轻声道："爹，您的长虹剑，我会了。"; 长虹贯日圆满成，一剑扬威震魔营。七侠之首初显威，长虹之光照天明。; palette: 红毛+长虹金+战火红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 虹猫，道果圆满阶段·七剑合璧（终极，七剑合璧·长虹为引）。形象：虹猫，红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。品性：七剑合璧，长虹为引，与魔教决战。这是复兴七剑的最终一战，也是一路磨砺的终点——七色剑光合璧，守护天下。。姿态：七剑合璧，长虹为引，七色剑光冲天斩向魔教；挚友之剑也合入其中——"这剑光，有你我一半。"。服饰：红毛小猫，长虹剑为引，七侠之首。。体型：身高约6头身，红毛小猫，长虹剑为引。。衣物细节：红色侠客装，长虹剑为引，护腕革带，七侠之首。。发型妆造：红毛束起，发带在剑气中猎猎。。脸型五官：猫脸，圆眼含笑含光，剑意圆满。。武器招式：七剑合璧·长虹贯日。。功法：七剑合璧（长虹为引）；长虹贯日圆满。。功法表现：七色剑光冲天，护山河人间。。画面：构图：决战之巅，红毛小猫长虹为引，七剑合璧七色剑光冲天斩向魔教，护住山河人间，背景天地。色调：长虹七色+战火红+神光金。氛围：终极、合璧、复兴兑现。。台词："七剑合璧，复兴七剑——爹，我做到了。这剑光，护的是我们一路护过来的天下。"。动作帧（动图）：①长虹为引 ②七剑合璧 ③七色剑光斩魔 ④护山河人间。诗词：七剑合璧光耀天，长虹为引斩魔烟。复兴之任终兑现，七色虹光护人间。。主题句：家破志不灭，长虹映月明——忍辱负重心中有光，赤子之心永不改。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：长虹七色+战火红+神光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 七剑合璧, age 终极, 七剑合璧·长虹为引; scene: 七剑合璧，长虹为引，七色剑光冲天斩向魔教；挚友之剑也合入其中——"这剑光，有你我一半。"; 七剑合璧光耀天，长虹为引斩魔烟。复兴之任终兑现，七色虹光护人间。; palette: 长虹七色+战火红+神光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 蓝兔（`lantu`） · 人生档案版

**灵胎初醒 · 冰魄剑意**
- 蓝兔，灵胎初醒阶段·玉蟾宫雪（幼年，玉蟾宫主·冰雪聪明）。形象：蓝兔，蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。品性：玉蟾宫宫主，冰魄剑传人。冰雪聪明，在玉蟾宫雪中长大，宫中姐妹相伴，岁月静好。。姿态：玉蟾宫前看雪，与姐妹共舞；指尖凝霜，冰魄剑意随雪而生。。服饰：蓝毛白兔，素衣清冷。。体型：身高约4头身，蓝毛白兔，身形清雅。。衣物细节：玉蟾宫素白长裙，淡蓝披帛，腰间缀冰魄坠，冰魄剑侧佩。。发型妆造：蓝毛梳得齐整，簪蓝色珠花，兔耳微垂。。脸型五官：兔脸，眼眸清亮，眉目温柔，嘴角浅淡。。武器招式：冰魄剑意初窥。。功法：冰魄剑意初窥；凝霜。。功法表现：指尖凝霜。。画面：构图：玉蟾宫雪境，蓝毛白兔立于雪中，与姐妹共舞，指尖凝霜，背景宫阙飞雪。色调：蓝毛+雪白+玉蟾青。氛围：太平、冰雪聪明、岁月静好。。台词："玉蟾宫的雪，最懂冰魄剑的心——清冷，却护着底下的小草。"。动作帧（动图）：①立雪 ②与姐妹共舞 ③指尖凝霜 ④望雪。诗词：玉蟾宫雪落纷纷，白兔凝霜看晓昏。姐妹相伴岁月好，冰魄剑意藏温存。。主题句：玉蟾宫雪染劫火，冰魄护世历孤勇——冰雪聪明，磨难中愈显刚柔。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：蓝毛+雪白+玉蟾青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 玉蟾宫雪, age 幼年, 玉蟾宫主·冰雪聪明; scene: 玉蟾宫前看雪，与姐妹共舞；指尖凝霜，冰魄剑意随雪而生。; 玉蟾宫雪落纷纷，白兔凝霜看晓昏。姐妹相伴岁月好，冰魄剑意藏温存。; palette: 蓝毛+雪白+玉蟾青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 蓝衣兔**
- 蓝兔，凡尘砺心阶段·玉蟾宫围（少年，魔教围宫·姐妹受难）。形象：蓝兔，蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。品性：魔教围困玉蟾宫，要夺冰魄剑、折七侠之翼。宫难当头，冰雪聪明的她，第一次直面血与火。。姿态：玉蟾宫前浴血守宫，冰魄剑一剑封霜；眼见姐妹被掳受难，她悲愤却更冷静。。服饰：蓝毛白兔，素衣染血，冰魄剑在手。。体型：身高约5头身，蓝毛白兔，素衣染血。。衣物细节：素白长裙染血，淡蓝披帛断了一角，紧握冰魄剑。。发型妆造：蓝毛凌乱，珠花沾雪。。脸型五官：兔脸，眼眸含泪却冷静，眉间悲愤。。武器招式：冰魄剑（守宫）。。功法：冰魄剑初成；守宫孤勇。。功法表现：冰魄封霜护宫。。画面：构图：玉蟾宫前，魔教围困，蓝毛白兔浴血守宫，冰魄剑封霜，姐妹被掳的虚影，背景宫阙战火。色调：蓝毛+雪染红+冰魄蓝。氛围：宫难、悲愤、孤勇。。台词："这宫里的雪，我守了这么多年。今日，我用它守我的姐妹！"。动作帧（动图）：①浴血守宫 ②一剑封霜 ③姐妹被掳 ④悲愤冷静。诗词：魔教围宫雪染红，白兔持剑守玉蟾。姐妹受难悲且愤，冰魄封霜护旧宫。。主题句：玉蟾宫雪染劫火，冰魄护世历孤勇——冰雪聪明，磨难中愈显刚柔。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：蓝毛+雪染红+冰魄蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 玉蟾宫围, age 少年, 魔教围宫·姐妹受难; scene: 玉蟾宫前浴血守宫，冰魄剑一剑封霜；眼见姐妹被掳受难，她悲愤却更冷静。; 魔教围宫雪染红，白兔持剑守玉蟾。姐妹受难悲且愤，冰魄封霜护旧宫。; palette: 蓝毛+雪染红+冰魄蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冰魄剑出**
- 蓝兔，道法初成阶段·雪中苦修（少年，宫难之后·苦练冰魄）。形象：蓝兔，蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。品性：宫难之后，她深知"剑太慢，护不住人"。被救回的姐妹临别前握着她的手说："蓝兔，替我们守好玉蟾宫。"她把这句话记进骨头里，日日雪中苦练。。姿态：玉蟾宫雪地里，她一遍遍练冰魄剑，一剑封霜从只能封住一片叶，练到封住一整片湖；夜里看雪，把每一道冰纹都记成剑招——每一招，都是对姐妹们的承诺。。服饰：蓝毛白兔，素衣，练剑磨出的薄茧的手。。体型：身高约5头身，蓝毛白兔，素衣，手有薄茧。。衣物细节：素白练剑裙，袖口束起，腰间冰魄坠，手上缠着练剑磨出的薄茧。。发型妆造：蓝毛束起，珠花依旧，眼神坚定。。脸型五官：兔脸，眼眸坚毅，眉间专注。。武器招式：冰魄剑（苦练中）。。功法：冰魄剑法苦练渐精；剑心渐明。。功法表现：一剑封霜渐成林。。画面：构图：玉蟾宫雪地，蓝毛白兔一遍遍练冰魄剑，一剑封霜从叶到湖，手被冻红仍不停，远处姐妹送别之影，背景飞雪。色调：蓝毛+冰魄蓝+雪白。氛围：苦修、坚韧、守诺。。台词："宫难那日，我的剑太慢，没能护住姐妹们。从今往后，我要把剑练到比雪还快——替她们，守好玉蟾宫。"。动作帧（动图）：①练冰魄剑 ②一剑封叶 ③封住整片湖 ④记冰纹为招。诗词：雪中苦练冰魄勤，一剑封霜渐成林。姐妹一诺心头记，白兔练剑到夜深。。主题句：玉蟾宫雪染劫火，冰魄护世历孤勇——冰雪聪明，磨难中愈显刚柔。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：蓝毛+冰魄蓝+雪白。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 雪中苦修, age 少年, 宫难之后·苦练冰魄; scene: 玉蟾宫雪地里，她一遍遍练冰魄剑，一剑封霜从只能封住一片叶，练到封住一整片湖；夜里看雪，把每一道冰纹都记成剑招——每一招，都是对姐妹们的承诺。; 雪中苦练冰魄勤，一剑封霜渐成林。姐妹一诺心头记，白兔练剑到夜深。; palette: 蓝毛+冰魄蓝+雪白; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 寒光剑影**
- 蓝兔，大劫淬炼阶段·剑成护侠（少年，剑法小成·为情挡箭）。形象：蓝兔，蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。品性：冰魄剑法小成，魔教来袭时第一次护住了人——为护虹猫挡下暗箭。情之所系，剑也更有力；恨之所向，霜也更冷。。姿态：魔教暗箭射向虹猫，她一剑封霜挡下，自己中箭；虹猫拼死救回她，雪夜围着炉火替她裹伤——她望着那双眼睛，心里那点情，比火还暖。。服饰：蓝毛白兔，白衣染血，卧于雪中。。体型：身高约5头身，蓝毛白兔，中箭卧雪。。衣物细节：素白长裙染血，淡蓝披帛为同伴裹身，冰魄剑在手。。发型妆造：蓝毛散落，眉间坚毅。。脸型五官：兔脸，眼眸含笑，眼底有暖。。武器招式：冰魄剑法（小成）。。功法：冰魄剑法小成；为情而护，为恨而战。。功法表现：炉火映暖雪地。。画面：构图：雪夜炉火，蓝毛白兔中箭卧着，虹猫替她裹伤，炉火映暖雪地，一侧魔教远影，背景风雪。色调：蓝毛+炉火暖+雪白+魔教暗。氛围：小成、为情护人、被护的暖。。台词："这一箭，替虹猫挡了，替大家挡了。我不后悔——因为想护着的人，就值得我用命去护；想讨的账，也记在心里。"。动作帧（动图）：①一剑封霜挡箭 ②自己中箭 ③同伴救回 ④炉火裹伤。诗词：剑成护侠为情挡，暗芒一箭护同裳。雪夜炉火裹伤处，情比火暖恨亦藏。。主题句：玉蟾宫雪染劫火，冰魄护世历孤勇——冰雪聪明，磨难中愈显刚柔。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：蓝毛+炉火暖+雪白+魔教暗。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 剑成护侠, age 少年, 剑法小成·为情挡箭; scene: 魔教暗箭射向虹猫，她一剑封霜挡下，自己中箭；虹猫拼死救回她，雪夜围着炉火替她裹伤——她望着那双眼睛，心里那点情，比火还暖。; 剑成护侠为情挡，暗芒一箭护同裳。雪夜炉火裹伤处，情比火暖恨亦藏。; palette: 蓝毛+炉火暖+雪白+魔教暗; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 玉兔剑仙**
- 蓝兔，封神登天阶段·冰魄封霜（青年，冰魄大成·雪夜护盟）。形象：蓝兔，蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。品性：冰魄剑法大成，一剑封霜百里。决战之前，她凭这一剑护住七侠，退去魔教夜袭。。姿态：魔教夜袭七侠营地，她一剑封霜百里，冻住来敌；雪夜里，她守在营前，直到天明。。服饰：蓝毛白兔，冰魄剑在手，七侠巾帼。。体型：身高约6头身，蓝毛白兔，冰魄剑在手。。衣物细节：玉蟾宫素白长裙，淡蓝披帛，冰魄剑在手，雪夜守护。。发型妆造：蓝毛束起，珠花含雪。。脸型五官：兔脸，眼眸坚定，雪夜守望。。武器招式：冰魄剑法（圆满）。。功法：冰魄剑法圆满；一剑封霜百里。。功法表现：一剑封霜百里。。画面：构图：雪夜营地，蓝毛白兔一剑封霜百里，冻住夜袭之敌，守在营前，背景雪夜战火。色调：冰魄蓝+雪白+战火微红。氛围：大成、护盟、守夜。。台词："这一剑封霜，替七侠守了一夜。魔教想趁夜偷袭，先过我蓝兔的霜！"。动作帧（动图）：①魔教夜袭 ②一剑封霜百里 ③冻住来敌 ④守在营前到天明。诗词：冰魄封霜百里明，雪夜护盟守天明。一剑冻住来敌路，玉蟾之霜镇夜营。。主题句：玉蟾宫雪染劫火，冰魄护世历孤勇——冰雪聪明，磨难中愈显刚柔。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：冰魄蓝+雪白+战火微红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 冰魄封霜, age 青年, 冰魄大成·雪夜护盟; scene: 魔教夜袭七侠营地，她一剑封霜百里，冻住来敌；雪夜里，她守在营前，直到天明。; 冰魄封霜百里明，雪夜护盟守天明。一剑冻住来敌路，玉蟾之霜镇夜营。; palette: 冰魄蓝+雪白+战火微红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 蓝兔，道果圆满阶段·七剑合璧（终极，七剑合璧·冰魄护世）。形象：蓝兔，蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。品性：七剑合璧，冰魄为护，与魔教决战。这是历劫的终点——冰魄剑气化寒光，护住七侠与天下。。姿态：七剑合璧，冰魄剑气化寒光护世；大战之中，为七侠封霜御敌，护到最后一刻。。服饰：蓝毛白兔，冰魄剑在手，七侠巾帼。。体型：身高约6头身，蓝毛白兔，冰魄剑在手。。衣物细节：玉蟾宫素白长裙，淡蓝披帛飘动，冰魄剑在手，七侠巾帼。。发型妆造：蓝毛束起，珠花在冰魄光中流转。。脸型五官：兔脸，眼眸坚定含笑。。武器招式：七剑合璧·冰魄惊鸿。。功法：七剑合璧（冰魄）；冰魄惊鸿圆满。。功法表现：冰魄寒光护世，护到最后一刻。。画面：构图：决战之巅，蓝毛白兔冰魄剑气化寒光护世，与长虹交映，七剑合璧护住天下，背景大战。色调：冰魄蓝+长虹金+战火红。氛围：终极、合璧、护世。。台词："七剑合璧，冰火相济。我蓝兔的霜，为七侠铺出一条生路，护到最后一刻！"。动作帧（动图）：①冰魄为护 ②七剑合璧 ③寒光护世 ④护七侠到最后一刻。诗词：七剑合璧护苍生，冰魄寒光映长虹。历劫玉蟾终合璧，玉蟾之霜护万穹。。主题句：玉蟾宫雪染劫火，冰魄护世历孤勇——冰雪聪明，磨难中愈显刚柔。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：冰魄蓝+长虹金+战火红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 七剑合璧, age 终极, 七剑合璧·冰魄护世; scene: 七剑合璧，冰魄剑气化寒光护世；大战之中，为七侠封霜御敌，护到最后一刻。; 七剑合璧护苍生，冰魄寒光映长虹。历劫玉蟾终合璧，玉蟾之霜护万穹。; palette: 冰魄蓝+长虹金+战火红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 逗逗（`doudou`） · 人生档案版

**灵胎初醒 · 雨花剑意**
- 逗逗，灵胎初醒阶段·药囊少年（幼年，雨花剑主·小药囊）。形象：逗逗，白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。品性：雨花剑传人，七侠开心果。调皮捣蛋，却随身带着小药囊，一颗医者仁心。。姿态：山林中嬉闹采药，药囊渐满；受伤的小兽，他偷偷包扎。。服饰：白毛小狗，小药囊斜挎。。体型：身高约4头身，白毛小狗，圆润。。衣物细节：黄米色小短褂，斜挎布药囊，腰间挂小药葫芦，脚踩布鞋。。发型妆造：白毛软蓬蓬，头顶扎个小发揪。。脸型五官：狗脸，圆眼机灵，咧嘴笑，腮圆。。武器招式：无兵器，药囊。。功法：采药识药；雨花剑意初窥。。功法表现：无，医者之能。。画面：构图：山林溪边，白毛小狗嬉闹采药，药囊渐满，身旁一只包扎好伤口的小兽，背景青山。色调：白毛+药草绿+溪光。氛围：嬉闹、药囊、少年。。台词："别看我爱闹，我这药囊里，装的可是救命的方子。"。动作帧（动图）：①嬉闹 ②采药 ③药囊渐满 ④给小兽包扎。诗词：白毛小狗闹山林，药囊斜挎采药寻。看似贪玩心实细，医者仁心藏得深。。主题句：嬉闹皮囊下，藏着一颗救人的医者仁心。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：白毛+药草绿+溪光。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 药囊少年, age 幼年, 雨花剑主·小药囊; scene: 山林中嬉闹采药，药囊渐满；受伤的小兽，他偷偷包扎。; 白毛小狗闹山林，药囊斜挎采药寻。看似贪玩心实细，医者仁心藏得深。; palette: 白毛+药草绿+溪光; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 白绒小狗**
- 逗逗，凡尘砺心阶段·医者无奈（少年，救不回至亲·医者之痛）。形象：逗逗，白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。品性：至亲中毒，他翻遍医书、尝尽百草，却终究回天乏术。医者能救千万人，却救不回至亲——那是他第一次尝到医者的无力。。姿态：他守着至亲的榻前，一遍遍换药、喂药，手指抖得端不稳碗；至亲离去那夜，他把药囊摔在墙上，又哭着捡起来。。服饰：白毛小狗，布衣，药囊。。体型：身高约5头身，白毛小狗，布衣。。衣物细节：黄米色短褂，药囊摔过又缝上，手抖着端药。。发型妆造：白毛凌乱，发揪散开。。脸型五官：狗脸，圆眼含泪，手抖。。武器招式：无，药囊。。功法：医术初成；尝到医者无力。。功法表现：烛火昏黄。。画面：构图：昏黄屋中，白毛小狗守在榻前，手抖着喂药，药囊散落一地，背景烛火。色调：白毛+烛火昏黄+药草。氛围：医者无奈、无力、悲恸。。台词："我学了这么久的医，却救不回最想救的人……可我不能停下。这药囊，还得背着，去救更多的人。"。动作帧（动图）：①翻医书 ②尝百草 ③守榻喂药 ④哭后捡起药囊。诗词：至亲中毒医难回，尝遍百草泪满杯。医者能救千万众，却救不得眼前人。。主题句：嬉闹皮囊下，藏着一颗救人的医者仁心。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：白毛+烛火昏黄+药草。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 医者无奈, age 少年, 救不回至亲·医者之痛; scene: 他守着至亲的榻前，一遍遍换药、喂药，手指抖得端不稳碗；至亲离去那夜，他把药囊摔在墙上，又哭着捡起来。; 至亲中毒医难回，尝遍百草泪满杯。医者能救千万众，却救不得眼前人。; palette: 白毛+烛火昏黄+药草; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 雨花剑出**
- 逗逗，道法初成阶段·采药苦修（少年，翻遍群山·苦练医剑）。形象：逗逗，白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。品性：把救不回至亲的痛，化作苦修之力。翻遍群山采药，练医术也练雨花剑——途中每遇到一个伤病的人，他都想起至亲躺在榻上的样子，然后多救一个。。姿态：背着药囊翻山越岭，采遍奇草；途中救下坠崖的樵夫、退热的孩童，药囊空了又满；夜里在灯下研医书、练雨花剑，剑气如雨，药香满屋，偶尔恍惚，烛火里仿佛映着至亲的影子。。服饰：白毛小狗，粗布劲装，药囊。。体型：身高约5头身，白毛小狗，粗布劲装。。衣物细节：黄米色劲装，药囊背带磨旧，腰间药葫芦，袖口沾着药渍。。发型妆造：白毛束起，发揪系着褪色布条。。脸型五官：狗脸，圆眼坚毅。。武器招式：雨花剑（苦练中）。。功法：医术精进；雨花剑法苦练渐精。。功法表现：剑气如雨，药香满屋。。画面：构图：群山之中，白毛小狗背药囊翻山采药，途中救下坠崖樵夫，夜灯下研医书练剑，烛火里映着至亲虚影，背景群山夜色。色调：白毛+药草绿+夜灯暖。氛围：苦修、救人、思念化力。。台词："救不回一个人，就救十个、百个。我这药囊和剑，要练到能接住每一个垂死的人——也算，替至亲多活几个。"。动作帧（动图）：①翻山 ②采奇草 ③夜读医书 ④练雨花剑。诗词：翻山采药苦修行，救得樵夫退热童。烛火恍惚至亲影，医剑双修济苍生。。主题句：嬉闹皮囊下，藏着一颗救人的医者仁心。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：白毛+药草绿+夜灯暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 采药苦修, age 少年, 翻遍群山·苦练医剑; scene: 背着药囊翻山越岭，采遍奇草；途中救下坠崖的樵夫、退热的孩童，药囊空了又满；夜里在灯下研医书、练雨花剑，剑气如雨，药香满屋，偶尔恍惚，烛火里仿佛映着至亲的影子。; 翻山采药苦修行，救得樵夫退热童。烛火恍惚至亲影，医剑双修济苍生。; palette: 白毛+药草绿+夜灯暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 悬壶剑医**
- 逗逗，大劫淬炼阶段·妙手初成（青年，第一次救回垂死之人）。形象：逗逗，白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。品性：医术与剑法小成，第一次救回一个垂死之人。救回来的那一刻，他想起救不回的至亲——原来医者的路，是替救不回的人，多救几个。。姿态：一个中毒濒死的侠客被送到他面前，他施针喂药，熬了三日三夜，终于把人拉回来；那人握紧他的手，他咧嘴笑了，眼泪却掉下来——"要是当年，我也能这样救回你就好了。"。服饰：白毛小狗，劲装，药囊。。体型：身高约6头身，白毛小狗，劲装。。衣物细节：黄米色劲装，药囊敞开，手指沾着药粉与银针。。发型妆造：白毛汗湿，发揪歪着。。脸型五官：狗脸，圆眼含泪带笑。。武器招式：雨花剑法（小成）。。功法：妙手回春小成；把对至亲的思念化作救人之力。。功法表现：烛光暖，药香浓。。画面：构图：医馆，白毛小狗守在病榻前三日救醒垂死侠客，笑中带泪，烛火旁仿佛映着至亲的影子，背景烛光药香。色调：白毛+烛火暖+药香+至亲虚影。氛围：妙手、初成、思念化力。。台词："救回来了！爹娘，你们看到了吗——我救不回你们，可我能救回别人了。这药囊，替你们，也替大家。"。动作帧（动图）：①施针喂药 ②三日三夜 ③救醒垂危 ④咧嘴笑中带泪。诗词：妙手初成救垂危，三日三夜针与灰。至亲若在当含笑，医者仁心代你辉。。主题句：嬉闹皮囊下，藏着一颗救人的医者仁心。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：白毛+烛火暖+药香+至亲虚影。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 妙手初成, age 青年, 第一次救回垂死之人; scene: 一个中毒濒死的侠客被送到他面前，他施针喂药，熬了三日三夜，终于把人拉回来；那人握紧他的手，他咧嘴笑了，眼泪却掉下来——"要是当年，我也能这样救回你就好了。"; 妙手初成救垂危，三日三夜针与灰。至亲若在当含笑，医者仁心代你辉。; palette: 白毛+烛火暖+药香+至亲虚影; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 雨花医圣**
- 逗逗，封神登天阶段·妙手回春（青年，妙手大成·名动江湖）。形象：逗逗，白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。品性：医术大成，妙手回春，救回无数重伤之人，医者之名传遍江湖。决战之前，他成了七侠最可靠的后盾。。姿态：连日救治伤者，妙手回春；江湖人称"逗神医"，他却笑着说"我还是那个爱闹的逗逗"。。服饰：白毛小狗，药囊，雨花剑在手。。体型：身高约6头身，白毛小狗，药囊。。衣物细节：黄米色医者装，斜挎药囊，腰挂药葫芦，江湖人称逗神医。。发型妆造：白毛束起，发揪系着新布条。。脸型五官：狗脸，圆眼含笑。。武器招式：妙手回春（圆满）。。功法：妙手回春圆满；医者之名。。功法表现：药香救人，仁心名动。。画面：构图：医馆，白毛小狗连日救治伤者，妙手回春，江湖百姓感激，背景医馆灯火。色调：白毛+药草绿+人间暖。氛围：大成、名动、医者仁心。。台词："江湖叫我逗神医，可我还是那个逗逗——医者，就是让更多的人，还能继续笑。"。动作帧（动图）：①连日救治 ②妙手回春 ③江湖称神医 ④笑说还是逗逗。诗词：妙手回春名动江，救人无数药囊香。七侠后盾医者在，嬉闹皮囊仁心长。。主题句：嬉闹皮囊下，藏着一颗救人的医者仁心。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：白毛+药草绿+人间暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 妙手回春, age 青年, 妙手大成·名动江湖; scene: 连日救治伤者，妙手回春；江湖人称"逗神医"，他却笑着说"我还是那个爱闹的逗逗"。; 妙手回春名动江，救人无数药囊香。七侠后盾医者在，嬉闹皮囊仁心长。; palette: 白毛+药草绿+人间暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 逗逗，道果圆满阶段·七剑合璧（终极，七剑合璧·雨润天下）。形象：逗逗，白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。品性：七剑合璧，雨花润泽，与魔教决战。医者在战场救死扶伤，雨花剑气润泽——剑与药，皆为守护天下。。姿态：七剑合璧，雨花剑气润泽战场；大战之中，他穿梭救人，雨润天下，药囊见底也心甘。。服饰：白毛小狗，雨花剑在手，药囊。。体型：身高约6头身，白毛小狗，雨花剑在手。。衣物细节：黄米色侠客装，药囊与雨花剑齐备，医者侠客。。发型妆造：白毛束起，发揪迎风。。脸型五官：狗脸，圆眼含笑。。武器招式：七剑合璧·雨花剑。。功法：七剑合璧（雨花）；妙手回春圆满。。功法表现：雨润天下，妙手仁心。。画面：构图：决战之巅，白毛小狗雨花剑气润泽战场，穿梭救人，七剑合璧护住天下，背景大战。色调：雨花青+战火红+药草绿。氛围：终极、合璧、救世。。台词："七剑合璧，雨润天下。我逗逗，剑要合璧，人也要都救回来！"。动作帧（动图）：①雨花润泽 ②七剑合璧 ③穿梭救人 ④护天下。诗词：七剑合璧雨润天，医者穿梭救死间。雨花剑气护天下，妙手仁心照万川。。主题句：嬉闹皮囊下，藏着一颗救人的医者仁心。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：雨花青+战火红+药草绿。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 七剑合璧, age 终极, 七剑合璧·雨润天下; scene: 七剑合璧，雨花剑气润泽战场；大战之中，他穿梭救人，雨润天下，药囊见底也心甘。; 七剑合璧雨润天，医者穿梭救死间。雨花剑气护天下，妙手仁心照万川。; palette: 雨花青+战火红+药草绿; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 大奔（`dabeng`） · 人生档案版

**灵胎初醒 · 奔雷剑意**
- 大奔，灵胎初醒阶段·蛮力少年（幼年，奔雷剑主·天生神力）。形象：大奔，棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。品性：奔雷剑传人，天生神力的棕毛大熊。憨厚豪爽，一把子蛮力，帮人搬重物，以为力气大就能护住所有想护的人。。姿态：山间扛木涉水，力大无穷；帮人搬重物，憨笑不居功。。服饰：棕毛大熊，粗布短打。。体型：身高约4头身，棕毛大熊，魁梧。。衣物细节：粗布短打，右袖卷到肩头，露出一只粗壮的胳膊，腰束麻绳，脚蹬草鞋。。发型妆造：一头棕毛乱蓬蓬，汗珠常挂，憨厚相。。脸型五官：熊脸，圆眼憨厚，咧嘴笑。。武器招式：无兵器，神力。。功法：天生神力；奔雷剑意初窥。。功法表现：无，天生神力。。画面：构图：山林溪边，棕毛大熊扛木涉水，憨笑回头，背景青山。色调：棕毛+木褐+溪青。氛围：蛮力、憨厚、少年。。台词："力气大有什么用？能帮人，才有用！"。动作帧（动图）：①扛木 ②涉水 ③帮人搬重物 ④憨笑。诗词：棕熊少年力无穷，扛木涉水笑山中。憨厚不居半分功，蛮力只为帮人雄。。主题句：蛮力护不住至亲，重伤磨砺赤子心——奔雷一剑，护住想护的人。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：棕毛+木褐+溪青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 蛮力少年, age 幼年, 奔雷剑主·天生神力; scene: 山间扛木涉水，力大无穷；帮人搬重物，憨笑不居功。; 棕熊少年力无穷，扛木涉水笑山中。憨厚不居半分功，蛮力只为帮人雄。; palette: 棕毛+木褐+溪青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 棕色小熊**
- 大奔，凡尘砺心阶段·至亲之痛（少年，家园遭袭·蛮力难护）。形象：大奔，棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。品性：魔教为夺奔雷剑突袭他家，至亲为护他身受重伤。他一拳砸碎山石，却发现蛮力护不住想护的人——第一次，他恨自己"只有一把子蛮力"。。姿态：至亲重伤在榻，他端药的手直抖；守了三天三夜，至亲虽救回，却留下暗疾——他跪在屋外，一拳一拳砸地，恨自己不够强。。服饰：棕毛大熊，布衣，眼眶泛红。。体型：身高约5头身，棕毛大熊，布衣。。衣物细节：布衣沾尘，一只胳膊仍露在外，端药的手青筋凸起。。发型妆造：棕毛凌乱，眼眶泛红。。脸型五官：熊脸，圆眼含泪，眼底坚毅。。武器招式：无，蛮力难护。。功法：天生神力；恨自己护不住。。功法表现：烛火映夜。。画面：构图：屋外，棕毛大熊跪地砸拳，眼眶泛红，屋内至亲重伤的灯影，背景夜。色调：棕毛+夜蓝+烛火暖。氛围：变故、自责、立志。。台词："我这一身蛮力，连至亲都护不住……我要练剑，练到谁都伤不了我想护的人。"。动作帧（动图）：①至亲重伤 ②端药手抖 ③守三天三夜 ④跪地砸拳立志。诗词：魔教突袭伤至亲，蛮力护不住一尘。跪地砸拳恨力弱，立志练剑护亲人。。主题句：蛮力护不住至亲，重伤磨砺赤子心——奔雷一剑，护住想护的人。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：棕毛+夜蓝+烛火暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 至亲之痛, age 少年, 家园遭袭·蛮力难护; scene: 至亲重伤在榻，他端药的手直抖；守了三天三夜，至亲虽救回，却留下暗疾——他跪在屋外，一拳一拳砸地，恨自己不够强。; 魔教突袭伤至亲，蛮力护不住一尘。跪地砸拳恨力弱，立志练剑护亲人。; palette: 棕毛+夜蓝+烛火暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 奔雷剑出**
- 大奔，道法初成阶段·奔雷苦修（少年，蛮力化剑·苦练奔雷）。形象：大奔，棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。品性：把"护不住至亲"的痛化作苦修之力，日夜练奔雷剑，把一身蛮力练成剑力——从一剑劈不开石，练到一剑奔雷。。姿态：深山石壁前，他一剑剑劈石，虎口震裂也咬牙不停；夜里对着至亲送他的旧剑，一遍遍练，汗混着泪。。服饰：棕毛大熊，粗布劲装，缠着布条的手。。体型：身高约5头身，棕毛大熊，粗布劲装，缠布的手。。衣物细节：破袖劲装，右臂整只裸露缠着布条，虎口磨出老茧，腰束革带挂剑鞘。。发型妆造：棕毛束起高马尾，汗与灰混染。。脸型五官：熊脸，圆眼坚毅，咬牙。。武器招式：奔雷剑（苦练中）。。功法：奔雷剑苦练渐精；蛮力化剑。。功法表现：剑势渐成奔雷。。画面：构图：深山石壁前，棕毛大熊一剑剑劈石，虎口缠布，汗泪齐落，背景群山夜色。色调：棕毛+石灰+奔雷金。氛围：苦修、化痛为志。。台词："爹说，蛮力护不住人，剑能。我且把这一身蛮力，一寸寸练成剑。"。动作帧（动图）：①一剑剑劈石 ②虎口震裂 ③夜里对旧剑练 ④汗泪齐落。诗词：蛮力化剑苦练勤，一剑奔雷渐成真。虎口裂开犹不停，只待剑能护亲人。。主题句：蛮力护不住至亲，重伤磨砺赤子心——奔雷一剑，护住想护的人。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：棕毛+石灰+奔雷金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 奔雷苦修, age 少年, 蛮力化剑·苦练奔雷; scene: 深山石壁前，他一剑剑劈石，虎口震裂也咬牙不停；夜里对着至亲送他的旧剑，一遍遍练，汗混着泪。; 蛮力化剑苦练勤，一剑奔雷渐成真。虎口裂开犹不停，只待剑能护亲人。; palette: 棕毛+石灰+奔雷金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 力拔山兮**
- 大奔，大劫淬炼阶段·为义重伤（青年，第一次护人·为义重伤）。形象：大奔，棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。品性：奔雷剑初成，第一次在魔教手下护住一个素不相识的人——虽然身受重伤，但看到那人逃出生天，他咧嘴笑了。。姿态：魔教追捕一个无辜者，他一剑奔雷挡下追兵，身受重伤；同伴救回他，他躺在担架上还在憨笑："值了，我护住了一个人。"。服饰：棕毛大熊，劲装染血，躺在担架。。体型：身高约6头身，棕毛大熊，劲装染血。。衣物细节：劲装染血，露出的右臂缠着绷带，仍紧握奔雷剑。。发型妆造：棕毛汗湿，眼底憨笑。。脸型五官：熊脸，圆眼憨笑，眼底有光。。武器招式：奔雷剑（初成）。。功法：奔雷剑初成；第一次护住人。。功法表现：一剑挡追兵。。画面：构图：林间，棕毛大熊为护无辜者身受重伤，同伴救回，他躺在担架憨笑，背景追兵远影。色调：棕毛+血绛+林青。氛围：小成、为义、被救的暖。。台词："受了伤，可那人活着。我这一剑，没白劈——我大奔，终于能护住人了。"。动作帧（动图）：①一剑奔雷挡追兵 ②身受重伤 ③同伴救回 ④担架憨笑。诗词：为义重伤剑初成，挡下追兵护苍生。担架犹作憨厚笑，蛮力终能护一人。。主题句：蛮力护不住至亲，重伤磨砺赤子心——奔雷一剑，护住想护的人。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：棕毛+血绛+林青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 为义重伤, age 青年, 第一次护人·为义重伤; scene: 魔教追捕一个无辜者，他一剑奔雷挡下追兵，身受重伤；同伴救回他，他躺在担架上还在憨笑："值了，我护住了一个人。"; 为义重伤剑初成，挡下追兵护苍生。担架犹作憨厚笑，蛮力终能护一人。; palette: 棕毛+血绛+林青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 奔雷剑圣**
- 大奔，封神登天阶段·奔雷震山（青年，奔雷大成·一剑震山）。形象：大奔，棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。品性：奔雷剑大成，一剑奔雷震山，威震江湖。至亲望着他，眼里有泪也有光——"我的儿，终能护住家了。"。姿态：江湖扬名，一剑奔雷震山，敌阵溃散；他回身，向至亲行了一礼——"爹娘，我能护住你们了。"。服饰：棕毛大熊，侠客装，奔雷剑在手。。体型：身高约6头身，棕毛大熊，侠客装。。衣物细节：粗犷侠客装，右臂整只裸露，肌肉虬结，腰束宽革带，奔雷剑斜挎。。发型妆造：棕毛束起，眉目豪迈。。脸型五官：熊脸，圆眼豪迈，眼底含光。。武器招式：奔雷震山（大成）。。功法：奔雷震山大成；威震江湖。。功法表现：一剑震山，威震江湖。。画面：构图：战场，棕毛大熊一剑奔雷震山，敌阵溃散，至亲在旁含泪含笑，背景山河。色调：奔雷金+战火红+至亲暖。氛围：大成、扬名、护亲。。台词："这一剑，为我自己，也为爹娘——从今往后，谁也别想伤我想护的人。"。动作帧（动图）：①一剑奔雷震山 ②敌阵溃散 ③回身向至亲行礼 ④护亲。诗词：奔雷震山一剑成，威震江湖护亲名。至亲见儿终能守，泪里有光笑里明。。主题句：蛮力护不住至亲，重伤磨砺赤子心——奔雷一剑，护住想护的人。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：奔雷金+战火红+至亲暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 奔雷震山, age 青年, 奔雷大成·一剑震山; scene: 江湖扬名，一剑奔雷震山，敌阵溃散；他回身，向至亲行了一礼——"爹娘，我能护住你们了。"; 奔雷震山一剑成，威震江湖护亲名。至亲见儿终能守，泪里有光笑里明。; palette: 奔雷金+战火红+至亲暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 大奔，道果圆满阶段·七剑合璧（终极，七剑合璧·奔雷震山）。形象：大奔，棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。品性：七剑合璧，奔雷震山，与魔教决战。他凭一身蛮力练出的剑，为七侠开路——护住天下，也护住至亲。。姿态：七剑合璧，他奔雷震山开路；大战之中，一力当先，护住七侠与身后的人间。。服饰：棕毛大熊，奔雷剑在手，七侠。。体型：身高约6头身，棕毛大熊，奔雷剑在手。。衣物细节：七侠劲装，右臂仍露，奔雷剑在手，粗犷豪迈。。发型妆造：棕毛束起，豪迈含笑。。脸型五官：熊脸，圆眼豪迈含笑。。武器招式：七剑合璧·奔雷震山。。功法：七剑合璧（奔雷）；奔雷震山圆满。。功法表现：奔雷开路，护世无前。。画面：构图：决战之巅，棕毛大熊奔雷震山开路，七剑合璧护世，护住七侠与人间，背景大战。色调：奔雷金+战火红+神光金。氛围：终极、合璧、护世。。台词："七剑合璧，我大奔开路！这一剑奔雷，护的是想护的所有人！"。动作帧（动图）：①七剑合璧 ②奔雷震山开路 ③护住七侠 ④护人间。诗词：七剑合璧奔雷鸣，棕熊开路护苍生。蛮力终成护世剑，赤子之心永不移。。主题句：蛮力护不住至亲，重伤磨砺赤子心——奔雷一剑，护住想护的人。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：奔雷金+战火红+神光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 七剑合璧, age 终极, 七剑合璧·奔雷震山; scene: 七剑合璧，他奔雷震山开路；大战之中，一力当先，护住七侠与身后的人间。; 七剑合璧奔雷鸣，棕熊开路护苍生。蛮力终成护世剑，赤子之心永不移。; palette: 奔雷金+战火红+神光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 跳跳（`tiaotiao`） · 人生档案版

**灵胎初醒 · 青光剑意**
- 跳跳，灵胎初醒阶段·灵巧小猴（幼年，青光剑传人·灵巧少年）。形象：跳跳，绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。品性：青光剑传人，绿毛小猴，身轻如燕，古灵精怪。家中父母双全，不知愁滋味。。姿态：林间跳跃如飞，与父母嬉戏；青光剑意在指尖跳跃，以为天下皆如这般太平。。服饰：绿毛小猴，劲装利落。。体型：身高约4头身，绿毛小猴，身形轻盈。。衣物细节：绿色短打劲装，袖口利落，腰间系绿带，身轻如燕。。发型妆造：绿毛束起小马尾，机灵相。。脸型五官：猴脸，圆眼机灵，咧嘴笑。。武器招式：青光剑意初窥。。功法：青光剑意初窥；灵巧身法。。功法表现：指尖青光。。画面：构图：山林，绿毛小猴在林间跳跃如飞，父母在旁含笑，指尖青光剑意，背景青山家园。色调：绿毛+林青+青光。氛围：太平、灵巧、不知愁。。台词："娘说，练好轻功，天大的祸事也追不上我。我跳跳，就是要一辈子这么跳下去。"。动作帧（动图）：①跳跃 ②与父母嬉戏 ③指尖剑意 ④大笑。诗词：绿毛小猴跳林间，父母相伴岁月闲。青光剑意指尖跃，以为天下皆安然。。主题句：灭门之仇深似海，忍辱卧底魔教中——血仇与道义，终有澄明一日。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：绿毛+林青+青光。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 灵巧小猴, age 幼年, 青光剑传人·灵巧少年; scene: 林间跳跃如飞，与父母嬉戏；青光剑意在指尖跳跃，以为天下皆如这般太平。; 绿毛小猴跳林间，父母相伴岁月闲。青光剑意指尖跃，以为天下皆安然。; palette: 绿毛+林青+青光; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 绿衣小猴**
- 跳跳，凡尘砺心阶段·灭门血仇（少年，黑心虎灭门·亲人尽丧）。形象：跳跳，绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。品性：黑心虎为夺七剑、灭口，血洗他家，父母亲人尽丧。他侥幸逃生，从此背负灭门血仇。。姿态：火光中父母将他推出，转身迎敌；他滚落山崖，回头只看见家园成灰，亲人尽丧——那夜，他再没笑过。。服饰：绿毛小猴，带伤，目光如刀。。体型：身高约5头身，绿毛小猴，带伤。。衣物细节：绿色短打染尘带伤，绿带松散，紧握青光剑。。发型妆造：绿毛凌乱，眼底如刀。。脸型五官：猴脸，圆眼含泪却如刀，再不见笑。。武器招式：青光剑（护身）。。功法：青光剑（护身）；血仇入骨。。功法表现：青光冷冽。。画面：构图：火光家园，绿毛小猴带伤立于崖边，回望家园成灰，目光如刀，背景冲天烈焰。色调：火光红+废墟黑+青光冷。氛围：灭门、血仇、幸存。。台词："爹、娘……我记下了。黑心虎，你欠我家的血，我跳跳，记一辈子。"。动作帧（动图）：①双亲推他 ②滚落山崖 ③回望家园成灰 ④目光如刀。诗词：灭门血火夜正深，双亲推儿入山林。滚落崖头回望处，故园成灰泪满襟。。主题句：灭门之仇深似海，忍辱卧底魔教中——血仇与道义，终有澄明一日。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：火光红+废墟黑+青光冷。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 灭门血仇, age 少年, 黑心虎灭门·亲人尽丧; scene: 火光中父母将他推出，转身迎敌；他滚落山崖，回头只看见家园成灰，亲人尽丧——那夜，他再没笑过。; 灭门血火夜正深，双亲推儿入山林。滚落崖头回望处，故园成灰泪满襟。; palette: 火光红+废墟黑+青光冷; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 青光剑出**
- 跳跳，道法初成阶段·忍辱卧底（少年，潜入魔教·剑藏于心）。形象：跳跳，绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。品性：为报灭门血仇，他潜入魔教卧底。扮作忠犬替黑心虎办事，笑里藏刀；无人时，在暗处练剑，把青光剑意一点点练回来——仇恨是他活下去的光，也是压着他的重。。姿态：魔教中他低头哈腰，替黑心虎跑腿，学他眯眼笑的样子；夜里回屋，他掏出藏着的青光剑，一遍遍练，剑光不敢太亮，怕被人看见。。服饰：绿毛小猴，夜行黑衣蒙面，头戴斗笠压低，笑意盈盈眼底藏刀。。体型：身高约5头身，绿毛小猴，魔教装束。。衣物细节：夜行黑衣，蒙面遮住大半张脸，头戴黑斗笠压低帽檐，腰悬魔教信物，袖中暗藏青光软剑。。发型妆造：斗笠压低，蒙面下绿毛隐现，只露一双带刀的眼。。脸型五官：猴脸，笑意盈盈，眼底藏刀。。武器招式：青光剑（暗练）。。功法：青光剑（暗练）；卧底之术；忍辱之心。。功法表现：一缕青光，压得极暗。。画面：构图：魔教暗室，绿毛小猴夜里掏出青光剑，剑光压得极暗，独自练剑，窗外黑心虎身影，背景魔教暗殿。色调：绿毛+魔教暗+一缕青光。氛围：卧底、隐忍、剑藏于心。。台词："黑心虎，你越信我，我离手刃你的那日就越近。这口恶气，我咽了三年——可剑，我没敢放下。"。动作帧（动图）：①低头哈腰 ②夜里掏剑 ③暗处练剑 ④剑光不敢亮。诗词：忍辱卧底入魔营，笑里藏刀待时横。暗夜练剑不敢亮，只为一刀报灭门。。主题句：灭门之仇深似海，忍辱卧底魔教中——血仇与道义，终有澄明一日。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：绿毛+魔教暗+一缕青光。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 忍辱卧底, age 少年, 潜入魔教·剑藏于心; scene: 魔教中他低头哈腰，替黑心虎跑腿，学他眯眼笑的样子；夜里回屋，他掏出藏着的青光剑，一遍遍练，剑光不敢太亮，怕被人看见。; 忍辱卧底入魔营，笑里藏刀待时横。暗夜练剑不敢亮，只为一刀报灭门。; palette: 绿毛+魔教暗+一缕青光; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 风驰剑动**
- 跳跳，大劫淬炼阶段·复仇道义（青年，恨与义·首次护人）。形象：跳跳，绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。品性：卧底日久，对黑心虎的灭门之恨与心底那点良善撕扯。一次黑心虎令他屠戮无辜，他剑发抖——最终忍不住出手，第一次护下了一个素不相识的人。。姿态：黑心虎命他杀一个无辜者，他举剑的手抖了又抖，忽然青光一闪，他斩开的不是那人，而是绳索；放那人走时，他低声道："快走，别问我是谁。"。服饰：绿毛小猴，夜行黑衣染尘，蒙面半掩，头戴斗笠，眉间第一次有了光。。体型：身高约6头身，绿毛小猴，魔教装束。。衣物细节：夜行黑衣已沾尘土，蒙面半掩露出嘴角，斗笠斜戴，腕上缠着放人时割断的绳索余痕。。发型妆造：斗笠歪戴，蒙面松落一角，眼底第一次有了光。。脸型五官：猴脸，眉间第一次有光。。武器招式：青光剑（初亮）。。功法：青光剑（初亮）；恨与义的交战。。功法表现：青光一闪，斩绳放人。。画面：构图：魔教暗室，绿毛小猴青光一闪斩断绳索，放走无辜者，眉间第一次亮起光，背景暗室与烛火。色调：绿毛+青光+暗室+一缕暖。氛围：恨与义、首次护人、心渐明。。台词："为报家仇，我什么都忍。可黑心虎，我不当你的刀——我跳跳，纵使满身是恨，也想做一回人。"。动作帧（动图）：①举剑欲屠 ②手抖 ③青光斩绳索 ④放人走。诗词：卧底日久恨与义，屠戮无辜剑难提。青光一闪斩绳索，放走无辜心始明。。主题句：灭门之仇深似海，忍辱卧底魔教中——血仇与道义，终有澄明一日。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：绿毛+青光+暗室+一缕暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 复仇道义, age 青年, 恨与义·首次护人; scene: 黑心虎命他杀一个无辜者，他举剑的手抖了又抖，忽然青光一闪，他斩开的不是那人，而是绳索；放那人走时，他低声道："快走，别问我是谁。"; 卧底日久恨与义，屠戮无辜剑难提。青光一闪斩绳索，放走无辜心始明。; palette: 绿毛+青光+暗室+一缕暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 灵猴剑圣**
- 跳跳，封神登天阶段·青光掠影（青年，快剑大成·一朝扬名）。形象：跳跳，绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。品性：真相大白，他放下卧底之身，归附七侠。快剑大成，一剑青光掠影——江湖从此记住了这柄快到看不清的剑。。姿态：与七侠并肩，他青光掠影一剑破敌，快若闪电；收剑时，他轻声道："爹娘，我回家了。"。服饰：绿毛小猴，七侠装束，青光剑在手。。体型：身高约6头身，绿毛小猴，七侠装束。。衣物细节：绿色侠客装，青光剑在手，腰间绿带，快剑扬名。。发型妆造：绿毛束起，机灵含笑。。脸型五官：猴脸，圆眼含笑，剑意明亮。。武器招式：青光掠影（大成）。。功法：青光掠影大成；归附七侠。。功法表现：青光一掠，快若闪电。。画面：构图：战场，绿毛小猴青光掠影一剑破敌，快若闪电，与七侠并肩，背景大战前奏。色调：青光蓝+七侠辉+战火红。氛围：大成、扬名、归正。。台词："这三年，我的剑藏着掖着。今日起，我的青光，要堂堂正正地亮。"。动作帧（动图）：①归附七侠 ②青光掠影一剑 ③破敌 ④轻声道回家了。诗词：青光掠影破敌营，快剑一扬江湖惊。卧底三年剑藏匣，一朝归正耀光明。。主题句：灭门之仇深似海，忍辱卧底魔教中——血仇与道义，终有澄明一日。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：青光蓝+七侠辉+战火红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 青光掠影, age 青年, 快剑大成·一朝扬名; scene: 与七侠并肩，他青光掠影一剑破敌，快若闪电；收剑时，他轻声道："爹娘，我回家了。"; 青光掠影破敌营，快剑一扬江湖惊。卧底三年剑藏匣，一朝归正耀光明。; palette: 青光蓝+七侠辉+战火红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 跳跳，道果圆满阶段·七剑合璧（终极，七剑合璧·青光掠影）。形象：跳跳，绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。品性：七剑合璧，青光掠影，与魔教决战。他迎着黑心虎，一剑报灭门之仇，也归于守护——仇恨落幕，青光从此为守护而亮。。姿态：七剑合璧，他青光掠影直取黑心虎；剑落之际，泪与血齐涌——"爹娘，我替你们报仇了。这剑光，往后再不染仇。"。服饰：绿毛小猴，青光剑在手，七侠。。体型：身高约6头身，绿毛小猴，青光剑在手。。衣物细节：绿色侠客装，青光剑在手，绿带猎猎，七侠。。发型妆造：绿毛束起，眼底澄明。。脸型五官：猴脸，圆眼含泪含笑。。武器招式：七剑合璧·青光掠影。。功法：七剑合璧（青光）；青光掠影圆满；报灭门之仇。。功法表现：青光护世，仇怨随风。。画面：构图：决战之巅，绿毛小猴青光掠影直取黑心虎，七剑合璧青光护世，泪血齐涌，背景大战。色调：青光蓝+战火红+神光金。氛围：终极、合璧、报仇归守护。。台词："黑心虎，这一剑还我爹娘。往后的青光，只护人，不染仇。"。动作帧（动图）：①七剑合璧 ②青光直取黑心虎 ③剑落泪血齐涌 ④青光护世。诗词：七剑合璧青光横，一剑报得灭门明。仇怨从此随风去，青光护世照苍生。。主题句：灭门之仇深似海，忍辱卧底魔教中——血仇与道义，终有澄明一日。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：青光蓝+战火红+神光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 七剑合璧, age 终极, 七剑合璧·青光掠影; scene: 七剑合璧，他青光掠影直取黑心虎；剑落之际，泪与血齐涌——"爹娘，我替你们报仇了。这剑光，往后再不染仇。"; 七剑合璧青光横，一剑报得灭门明。仇怨从此随风去，青光护世照苍生。; palette: 青光蓝+战火红+神光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 莎丽（`shali`） · 人生档案版

**灵胎初醒 · 紫云剑意**
- 莎丽，灵胎初醒阶段·紫云剑主（幼年，紫云剑主·温柔坚韧）。形象：莎丽，紫毛松鼠，温柔坚韧。 核心意象：紫云剑、蓬松大尾、温柔剑意、七侠巾帼。品性：紫云剑传人，紫毛松鼠，外柔内刚，剑法细腻。紫云剑在手，岁月静好。。姿态：溪边练剑，剑影如紫云；师门同伴相伴，她以为这样的日子会一直下去。。服饰：紫毛松鼠，素衣。。体型：身高约4头身，紫毛松鼠，身形柔美。。衣物细节：紫色短打劲装，袖口利落，腰间系紫带，紫云剑佩于身侧。。发型妆造：紫毛梳得齐整，发间系紫色发带。。脸型五官：松鼠脸，圆眼温柔，颊边蓬毛，眉间却含坚韧。。武器招式：紫云剑意初窥。。功法：紫云剑意初窥；剑法细腻。。功法表现：剑影如云。。画面：构图：溪边，紫毛松鼠练剑，剑影如紫云，素衣飘然，师门同伴相伴，背景青山流水。色调：紫毛+紫云+溪青。氛围：太平、剑主、岁月静好。。台词："剑要柔，心要刚。这紫云剑，我握得稳稳的，谁都夺不走。"。动作帧（动图）：①立溪边 ②练剑 ③师门相伴 ④以为岁月长。诗词：紫毛松鼠立溪边，剑影如云自悠然。温柔外表刚心在，以为紫云常在手。。主题句：紫云剑被夺，武功被废——隐忍蛰伏磨剑心，终还紫云归其主。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫毛+紫云+溪青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 紫云剑主, age 幼年, 紫云剑主·温柔坚韧; scene: 溪边练剑，剑影如紫云；师门同伴相伴，她以为这样的日子会一直下去。; 紫毛松鼠立溪边，剑影如云自悠然。温柔外表刚心在，以为紫云常在手。; palette: 紫毛+紫云+溪青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 紫衣松鼠**
- 莎丽，凡尘砺心阶段·被夺剑（少年，马三娘夺剑·遭陷重伤）。形象：莎丽，紫毛松鼠，温柔坚韧。 核心意象：紫云剑、蓬松大尾、温柔剑意、七侠巾帼。品性：马三娘（魔教卧底）设局夺走紫云剑，她遭陷害，身负重伤，紫云剑离手——那是她第一次觉得天塌了。。姿态：马三娘假作亲近，突下杀手；她重伤倒地，眼睁睁看着紫云剑被夺走，嘶声："那是我的剑……"。服饰：紫毛松鼠，重伤，素衣染血。。体型：身高约5头身，紫毛松鼠，重伤。。衣物细节：紫色劲装染血，紫带断裂，紫云剑被夺的手还伸向空处。。发型妆造：紫毛凌乱，发带散开。。脸型五官：松鼠脸，圆眼含泪，颊边蓬毛沾血，手伸向空处。。武器招式：紫云剑被夺。。功法：紫云剑被夺；重伤。。功法表现：剑光黯淡。。画面：构图：暗夜，紫毛松鼠重伤倒地，眼睁睁看马三娘夺走紫云剑，手伸向空处，背景阴冷。色调：紫毛+血绛+夜色暗。氛围：被夺、受陷、失剑。。台词："紫云剑……还给我！那是师门传给我的，是我的剑！"。动作帧（动图）：①马三娘假作亲近 ②突下杀手 ③重伤倒地 ④眼睁睁看剑被夺。诗词：马娘设局夺紫云，重伤倒地剑离身。嘶声欲唤剑归处，从此剑主成失魂。。主题句：紫云剑被夺，武功被废——隐忍蛰伏磨剑心，终还紫云归其主。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫毛+血绛+夜色暗。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 被夺剑, age 少年, 马三娘夺剑·遭陷重伤; scene: 马三娘假作亲近，突下杀手；她重伤倒地，眼睁睁看着紫云剑被夺走，嘶声："那是我的剑……"; 马娘设局夺紫云，重伤倒地剑离身。嘶声欲唤剑归处，从此剑主成失魂。; palette: 紫毛+血绛+夜色暗; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 紫云剑出**
- 莎丽，道法初成阶段·武功被废（少年，武功被废·沦为废人）。形象：莎丽，紫毛松鼠，温柔坚韧。 核心意象：紫云剑、蓬松大尾、温柔剑意、七侠巾帼。品性：被废武功，从紫云剑主沦为废人。江湖只记得"紫云剑主已死"，无人知她苟活于世。。姿态：荒郊破庙中，她握不起剑，连筷子都端不稳；每日对着断剑，恨意与不甘在心里翻涌。。服饰：紫毛松鼠，布衣，落魄。。体型：身高约5头身，紫毛松鼠，布衣落魄。。衣物细节：破旧布衣，紫带不见，双手无力垂着，落魄。。发型妆造：紫毛披散，面无血色。。脸型五官：松鼠脸，圆眼不甘，颊边蓬毛凌乱，眉间恨意。。武器招式：无，握不起剑。。功法：武功被废；隐忍蛰伏。。功法表现：断剑寒光。。画面：构图：荒郊破庙，紫毛松鼠布衣独坐，对着一柄断剑，握不起的手在颤，神色不甘，背景残庙荒草。色调：紫毛+破庙灰+断剑锈。氛围：废武、蛰伏、不甘。。台词："我连剑都握不起来了……可我不甘心！紫云剑，我一定会亲手夺回来！"。动作帧（动图）：①握不起剑 ②端不稳筷 ③对断剑 ④恨意翻涌。诗词：废武之痛蚀骨深，破庙独对断剑痕。江湖皆道剑主死，谁知青灯隐恨心。。主题句：紫云剑被夺，武功被废——隐忍蛰伏磨剑心，终还紫云归其主。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫毛+破庙灰+断剑锈。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 武功被废, age 少年, 武功被废·沦为废人; scene: 荒郊破庙中，她握不起剑，连筷子都端不稳；每日对着断剑，恨意与不甘在心里翻涌。; 废武之痛蚀骨深，破庙独对断剑痕。江湖皆道剑主死，谁知青灯隐恨心。; palette: 紫毛+破庙灰+断剑锈; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 云剑逍遥**
- 莎丽，大劫淬炼阶段·忍辱重练（青年，暗中重练·剑心淬炼）。形象：莎丽，紫毛松鼠，温柔坚韧。 核心意象：紫云剑、蓬松大尾、温柔剑意、七侠巾帼。品性：暗中重练剑法，从握不起剑到再握紫云剑法。忍辱负重，剑心在痛苦中淬炼，越发坚韧。。姿态：深夜无人处，她一遍遍练剑，跌倒再起；手腕上的旧伤每每撕裂，她咬碎了牙也不吭声。。服饰：紫毛松鼠，布衣，缠着绷带的手。。体型：身高约6头身，紫毛松鼠，缠绷带的手。。衣物细节：粗布劲装，腕上缠着绷带，紫带重新系上，紫云剑影在练。。发型妆造：紫毛束起，发带系紧，眼神坚韧。。脸型五官：松鼠脸，圆眼坚韧，颊边蓬毛沾汗，咬牙。。武器招式：紫云剑法重练。。功法：紫云剑法重练；剑心淬炼。。功法表现：剑光微亮。。画面：构图：深夜荒野，紫毛松鼠缠着绷带练剑，跌倒再起，剑光微亮，背景夜月孤影。色调：紫毛+绷带白+夜蓝。氛围：隐忍、重练、磨砺。。台词："手断了，可以再练；剑没了，可以再夺。我莎丽，没那么容易认输。"。动作帧（动图）：①深夜练剑 ②跌倒 ③咬牙再起 ④旧伤撕心不吭声。诗词：忍辱重练夜夜勤，跌起咬牙剑欲寻。旧伤撕心浑不顾，只待紫云再归身。。主题句：紫云剑被夺，武功被废——隐忍蛰伏磨剑心，终还紫云归其主。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫毛+绷带白+夜蓝。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 忍辱重练, age 青年, 暗中重练·剑心淬炼; scene: 深夜无人处，她一遍遍练剑，跌倒再起；手腕上的旧伤每每撕裂，她咬碎了牙也不吭声。; 忍辱重练夜夜勤，跌起咬牙剑欲寻。旧伤撕心浑不顾，只待紫云再归身。; palette: 紫毛+绷带白+夜蓝; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 紫霞剑圣**
- 莎丽，封神登天阶段·夺回紫云（青年，与马三娘对决·夺剑雪耻）。形象：莎丽，紫毛松鼠，温柔坚韧。 核心意象：紫云剑、蓬松大尾、温柔剑意、七侠巾帼。品性：终于等到时机，与马三娘对决。旧伤之痛化作剑意，紫云万丈，一剑夺回本属于她的剑。。姿态：紫云剑相抵，她盯着马三娘，一字一句："这剑，是我的。"剑落，紫云归主，雪耻。。服饰：紫毛松鼠，紫云剑重回在手，战意凛然。。体型：身高约6头身，紫毛松鼠，紫云剑在手。。衣物细节：紫色女侠劲装，紫云剑重回在手，发带鲜亮，战意凛然。。发型妆造：紫毛束起，发带迎风。。脸型五官：松鼠脸，圆眼凛然，颊边蓬毛扬起，含雪耻之光。。武器招式：紫云剑法（大成）。。功法：紫云剑法大成；夺剑雪耻。。功法表现：紫云万丈。。画面：构图：对决之地，紫毛松鼠与马三娘紫云剑相抵，剑落紫云归主，神色雪耻，背景风云。色调：紫毛+紫云+对决冷。氛围：夺回、雪耻、对决。。台词："你夺我的剑，废我的武功，我忍了这些年。今日，紫云归主——这一剑，还你！"。动作帧（动图）：①紫云剑相抵 ②盯着马三娘 ③一剑归主 ④雪耻而立。诗词：夺剑雪耻紫云横，一剑归主旧仇平。隐忍数载终有报，紫毛猫影立天明。。主题句：紫云剑被夺，武功被废——隐忍蛰伏磨剑心，终还紫云归其主。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫毛+紫云+对决冷。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 夺回紫云, age 青年, 与马三娘对决·夺剑雪耻; scene: 紫云剑相抵，她盯着马三娘，一字一句："这剑，是我的。"剑落，紫云归主，雪耻。; 夺剑雪耻紫云横，一剑归主旧仇平。隐忍数载终有报，紫毛猫影立天明。; palette: 紫毛+紫云+对决冷; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 莎丽，道果圆满阶段·七剑合璧（终极，七剑合璧·紫云万丈）。形象：莎丽，紫毛松鼠，温柔坚韧。 核心意象：紫云剑、蓬松大尾、温柔剑意、七侠巾帼。品性：七剑合璧，紫云万丈，与魔教决战。历尽被夺、废武、隐忍、雪耻的她，在合璧中剑心圆满——紫云归主，护苍生。。姿态：七剑合璧，紫云万丈护世；大战之中，她以夺回的紫云剑，为七侠遮风挡雨——"这剑，失而复得，如今护的是大家。"。服饰：紫毛松鼠，紫云剑在手，七侠巾帼。。体型：身高约6头身，紫毛松鼠，紫云剑在手。。衣物细节：紫色女侠装，紫云剑在手，紫带飘动，七侠巾帼。。发型妆造：紫毛束起，发带在紫云剑气中流转。。脸型五官：松鼠脸，圆眼温柔明亮，颊边蓬毛轻拂，眼底历劫之光。。武器招式：七剑合璧·紫云万丈。。功法：七剑合璧（紫云）；紫云万丈圆满。。功法表现：紫云万丈护苍生。。画面：构图：决战之巅，紫毛松鼠紫云万丈护世，七剑合璧紫云护苍生，背景大战。色调：紫云+战火红+神光金。氛围：终极、合璧、护世。。台词："马三娘夺了我的剑，可夺不走我的心。如今紫云归主，合璧护世——剑失而复得，心从未丢过。"。动作帧（动图）：①七剑合璧 ②紫云万丈 ③护世 ④为七侠遮风挡雨。诗词：七剑合璧紫云横，万丈剑影护苍生。失而复得终合璧，巾帼之剑镇乾坤。。主题句：紫云剑被夺，武功被废——隐忍蛰伏磨剑心，终还紫云归其主。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫云+战火红+神光金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 七剑合璧, age 终极, 七剑合璧·紫云万丈; scene: 七剑合璧，紫云万丈护世；大战之中，她以夺回的紫云剑，为七侠遮风挡雨——"这剑，失而复得，如今护的是大家。"; 七剑合璧紫云横，万丈剑影护苍生。失而复得终合璧，巾帼之剑镇乾坤。; palette: 紫云+战火红+神光金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 达达（`dada`） · 人生档案版

**灵胎初醒 · 旋风剑意**
- 达达，灵胎初醒阶段·沉稳少年（幼年，旋风剑主·闲云野鹤）。形象：达达，绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。品性：旋风剑传人，绿面熊猫，文雅温和，沉稳大气。。姿态：竹林间闲坐品茶，观竹叶旋落；旋风剑意，随竹影而生。。服饰：绿面熊猫，青衫。。体型：身高约5头身，绿面熊猫，身形沉稳。。衣物细节：青衫书生装，袖口宽大，腰间系玉带，手持折扇，竹林闲坐。。发型妆造：绿毛梳得齐整，束发道髻，白净儒雅。。脸型五官：熊猫脸，圆眼温和，眉目沉稳。。武器招式：旋风剑意初窥。。功法：旋风剑意初窥；以静制动。。功法表现：剑意随竹。。画面：构图：竹林，绿面熊猫青衫闲坐，观竹叶旋落，旋风剑意随竹影，背景竹海。色调：绿面+竹青+青衫。氛围：沉稳、野鹤、旋风。。台词："竹叶能借风，剑也能。我达达，惯会借力打力。"。动作帧（动图）：①坐竹林 ②品茶 ③观竹叶旋 ④剑意随竹。诗词：绿面熊猫坐竹林，闲观竹叶随风吟。旋风剑意竹影起，沉稳之中藏剑心。。主题句：竹林家园被毁，隐忍磨砺借力剑——以柔克刚，守一方太平。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：绿面+竹青+青衫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 沉稳少年, age 幼年, 旋风剑主·闲云野鹤; scene: 竹林间闲坐品茶，观竹叶旋落；旋风剑意，随竹影而生。; 绿面熊猫坐竹林，闲观竹叶随风吟。旋风剑意竹影起，沉稳之中藏剑心。; palette: 绿面+竹青+青衫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 竹林熊猫**
- 达达，凡尘砺心阶段·竹林被毁（少年，竹林家园·魔教焚毁）。形象：达达，绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。品性：魔教为夺旋风剑，焚毁他的竹林家园。他一向隐忍退让，这一次却眼睁睁看着家园成灰——他握剑的手，第一次握出青筋。。姿态：魔教放火焚竹，他救出几户村民，回身看竹林成灰；他蹲在焦土前，捡起一片焦叶，低声道："我会让你们，再长回来。"。服饰：绿面熊猫，青衫染灰，神色隐忍。。体型：身高约5头身，绿面熊猫，青衫染灰。。衣物细节：青衫书生装染灰，玉带仍在，折扇已收起，握剑的手攥紧。。发型妆造：束发微乱，白净面容沾尘。。脸型五官：熊猫脸，圆眼隐忍，眼底坚毅。。武器招式：旋风剑（护人）。。功法：旋风剑（护人）；隐忍之痛。。功法表现：焦土余烬。。画面：构图：焦土竹林，绿面熊猫拾起焦叶，身后村民受护，背景余烬与灰烟。色调：竹灰+焦黑+一线青。氛围：变故、隐忍、立志。。台词："我惯会忍，可这竹林是我的家。家没了，我还忍什么？——我要练剑，练到没人再敢烧我的家。"。动作帧（动图）：①魔教放火 ②救出村民 ③看竹林成灰 ④拾焦叶立志。诗词：魔火焚竹家园残，救得村邻看烬寒。拾起焦叶誓重长，隐忍之心化剑澜。。主题句：竹林家园被毁，隐忍磨砺借力剑——以柔克刚，守一方太平。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：竹灰+焦黑+一线青。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 竹林被毁, age 少年, 竹林家园·魔教焚毁; scene: 魔教放火焚竹，他救出几户村民，回身看竹林成灰；他蹲在焦土前，捡起一片焦叶，低声道："我会让你们，再长回来。"; 魔火焚竹家园残，救得村邻看烬寒。拾起焦叶誓重长，隐忍之心化剑澜。; palette: 竹灰+焦黑+一线青; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 旋风剑出**
- 达达，道法初成阶段·借风苦修（少年，以柔克刚·苦练旋风）。形象：达达，绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。品性：把家园被毁的痛化作苦修之力，日夜练旋风剑，以柔克刚、借力打力——从借竹叶之力，练到借山河之力。。姿态：荒山风口，他一剑剑借风旋叶，以柔克刚；夜里对着竹林的方向，一遍遍练，剑势化柔为刚。。服饰：绿面熊猫，粗布劲装，青衫已旧。。体型：身高约5头身，绿面熊猫，粗布劲装。。衣物细节：青衫练剑装，袖口束起，腰间玉带，头戴斗笠，旋风剑在手。。发型妆造：束发道髻，斗笠压着，眼神沉静。。脸型五官：熊猫脸，圆眼坚毅。。武器招式：旋风剑（苦练中）。。功法：旋风剑苦练渐精；借力打力。。功法表现：借风旋叶。。画面：构图：荒山风口，绿面熊猫一剑剑借风旋叶，以柔克刚，背景荒山与远处焦竹。色调：绿面+风青+旋叶金。氛围：苦修、借力、化痛为志。。台词："竹被烧了，可风还在。我且借这风，练成护得住家园的剑。"。动作帧（动图）：①荒山练剑 ②借风旋叶 ③夜里对竹林方向练 ④剑势化柔为刚。诗词：借风苦练旋叶勤，以柔克刚渐成林。焦叶化作剑中力，只为护家再成阴。。主题句：竹林家园被毁，隐忍磨砺借力剑——以柔克刚，守一方太平。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：绿面+风青+旋叶金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 借风苦修, age 少年, 以柔克刚·苦练旋风; scene: 荒山风口，他一剑剑借风旋叶，以柔克刚；夜里对着竹林的方向，一遍遍练，剑势化柔为刚。; 借风苦练旋叶勤，以柔克刚渐成林。焦叶化作剑中力，只为护家再成阴。; palette: 绿面+风青+旋叶金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 竹林剑舞**
- 达达，大劫淬炼阶段·以柔护弱（青年，第一次护人·以柔克刚）。形象：达达，绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。品性：旋风剑初成，第一次在魔教手下以柔克刚护住村民——借对方的力，还对方一身。家园的痛，化作守护的力。。姿态：魔教欺凌村民，他以柔克刚，借力打力，四两拨千斤退敌；村民得救，他轻轻道："有我在，这太平，我守。"。服饰：绿面熊猫，青衫侠者，旋风剑在手。。体型：身高约6头身，绿面熊猫，青衫侠者。。衣物细节：青衫侠者装，头戴斗笠，旋风剑在手，以柔克刚。。发型妆造：束发，斗笠下目光温和坚定。。脸型五官：熊猫脸，圆眼沉稳温和。。武器招式：旋风剑（初成）。。功法：旋风剑初成；以柔克刚护人。。功法表现：四两拨千斤。。画面：构图：村落，绿面熊猫以柔克刚退敌，村民得救相护，背景家园与魔教远影。色调：绿面+竹青+家园暖。氛围：小成、以柔护弱、守太平。。台词："我借过竹叶之力，借过山河之力。如今，我借你的力，护我想护的人。"。动作帧（动图）：①魔教欺凌村民 ②以柔克刚退敌 ③村民得救 ④轻道守太平。诗词：以柔护弱借力成，四两拨千斤退兵。家园之痛化剑意，一柄旋风守太平。。主题句：竹林家园被毁，隐忍磨砺借力剑——以柔克刚，守一方太平。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：绿面+竹青+家园暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 以柔护弱, age 青年, 第一次护人·以柔克刚; scene: 魔教欺凌村民，他以柔克刚，借力打力，四两拨千斤退敌；村民得救，他轻轻道："有我在，这太平，我守。"; 以柔护弱借力成，四两拨千斤退兵。家园之痛化剑意，一柄旋风守太平。; palette: 绿面+竹青+家园暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 熊猫剑圣**
- 达达，封神登天阶段·旋风破竹（青年，旋风大成·一剑破竹）。形象：达达，绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。品性：旋风剑大成，一剑破竹，威震一方。他回到焦土竹林，一剑挥出，春风过处，新竹破土——家园，终是重新长起来了。。姿态：江湖扬名，旋风破竹退敌无数；他重回故地，一剑借风，春风过处，焦土里钻出新笋——他蹲下身，轻抚新竹。。服饰：绿面熊猫，侠客装，旋风剑在手。。体型：身高约6头身，绿面熊猫，侠客装。。衣物细节：青衫侠者装，头戴斗笠，旋风剑在手，立于春风新竹间。。发型妆造：束发，斗笠轻抬，眼含笑意。。脸型五官：熊猫脸，圆眼含笑，眼底有光。。武器招式：旋风破竹（大成）。。功法：旋风破竹大成；守一方太平。。功法表现：春风过处新竹生。。画面：构图：故地竹林，绿面熊猫一剑借风，焦土中新笋破土，他蹲身轻抚，背景春风。色调：绿面+竹青+新绿。氛围：大成、守太平、家园重立。。台词："竹烧了，会再长；家没了，会再立。这一剑，我护的不是一片竹，是一方太平。"。动作帧（动图）：①一剑破竹 ②焦土新笋 ③蹲身轻抚 ④守太平。诗词：旋风破竹一剑成，威震江湖守太平。焦土新笋春风起，家园重立绿又明。。主题句：竹林家园被毁，隐忍磨砺借力剑——以柔克刚，守一方太平。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：绿面+竹青+新绿。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 旋风破竹, age 青年, 旋风大成·一剑破竹; scene: 江湖扬名，旋风破竹退敌无数；他重回故地，一剑借风，春风过处，焦土里钻出新笋——他蹲下身，轻抚新竹。; 旋风破竹一剑成，威震江湖守太平。焦土新笋春风起，家园重立绿又明。; palette: 绿面+竹青+新绿; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 达达，道果圆满阶段·七剑合璧（终极，七剑合璧·旋风破竹）。形象：达达，绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。品性：七剑合璧，旋风破竹，与魔教决战。他以柔克刚，卸去强敌之力，守一方太平——护住天下，也护住重立的家园。。姿态：七剑合璧，他旋风破竹，借力打力，卸去魔教千钧之力；大战之中，稳如定海神针。。服饰：绿面熊猫，旋风剑在手，七侠。。体型：身高约6头身，绿面熊猫，旋风剑在手。。衣物细节：青衫侠客装，头戴斗笠，旋风剑在手，七侠长者。。发型妆造：束发，斗笠下目光沉稳。。脸型五官：熊猫脸，圆眼沉稳含笑。。武器招式：七剑合璧·旋风破竹。。功法：七剑合璧（旋风）；旋风破竹圆满。。功法表现：以柔克刚，卸力守天下。。画面：构图：决战之巅，绿面熊猫旋风破竹卸去敌力，七剑合璧守天下，背景大战。色调：竹青+旋风金+战火红。氛围：终极、合璧、守太平。。台词："七剑合璧，我达达借力——把敌人的力道，还给天下太平。"。动作帧（动图）：①七剑合璧 ②旋风破竹 ③卸去敌力 ④守天下。诗词：七剑合璧旋风横，破竹之势卸千钧。以柔克刚守天下，熊猫一剑镇乾坤。。主题句：竹林家园被毁，隐忍磨砺借力剑——以柔克刚，守一方太平。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：竹青+旋风金+战火红。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 七剑合璧, age 终极, 七剑合璧·旋风破竹; scene: 七剑合璧，他旋风破竹，借力打力，卸去魔教千钧之力；大战之中，稳如定海神针。; 七剑合璧旋风横，破竹之势卸千钧。以柔克刚守天下，熊猫一剑镇乾坤。; palette: 竹青+旋风金+战火红; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 黑心虎（`heixinhu`） · 人生档案版

**灵胎初醒 · 暗影虎意**
- 黑心虎，灵胎初醒阶段·黑虎幼崽（幼年，魔教之虎·天生争强）。形象：黑心虎，黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。品性：魔教的黑虎，天生争强好胜。虎亦有虎的尊严，亦正亦邪。。姿态：山林中称王，虎啸生风；见幼兽受欺，也会抬眼放它一马。。服饰：黑毛猛虎，双目如炬。。体型：身高约5头身，黑毛猛虎，魁梧。。衣物细节：通体墨黑虎纹，前额"王"纹初显，粗壮四肢，黑尾如鞭。。发型妆造：墨黑虎毛根根如针，双目如炬。。脸型五官：虎脸，圆眼如炬，眉间王纹。。武器招式：虎爪/虎啸。。功法：天生虎威；黑风初起。。功法表现：虎啸生风。。画面：构图：山林，黑毛猛虎幼崽立于高石，虎啸生风，双目如炬，背景青山。色调：黑毛+山青+虎目金。氛围：虎性、争强、亦正亦邪。。台词："虎有虎的规矩——强者为王，可我不欺负弱小的。"。动作帧（动图）：①立高石 ②虎啸 ③放走幼兽 ④回望。诗词：黑虎幼崽立山林，虎啸生风欲称尊。亦正亦邪天性里，放走幼兽留一分仁。。主题句：虎踞魔道亦有性，只叹一生总在争与执。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：黑毛+山青+虎目金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 黑虎幼崽, age 幼年, 魔教之虎·天生争强; scene: 山林中称王，虎啸生风；见幼兽受欺，也会抬眼放它一马。; 黑虎幼崽立山林，虎啸生风欲称尊。亦正亦邪天性里，放走幼兽留一分仁。; palette: 黑毛+山青+虎目金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小老虎**
- 黑心虎，凡尘砺心阶段·黑风虎王（少年，黑风虎王·威压渐盛）。形象：黑心虎，黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。品性：黑风渐盛，虎王威压。争强好胜，一步步踏向魔道。。姿态：以黑风震慑群兽；与正道交锋，威压初显。。服饰：黑毛猛虎，黑风缭绕。。体型：身高约6头身，黑毛猛虎，黑风绕身。。衣物细节：墨黑虎纹，肩胛渐披暗色甲片，黑风缭绕，威压四野。。发型妆造：墨黑虎毛，目露凶光。。脸型五官：虎脸，圆眼含威。。武器招式：黑风虎啸。。功法：黑风虎啸；威压。。功法表现：黑风缭绕。。画面：构图：幽暗山林，黑毛猛虎黑风缭绕，威压四野，背景暗色。色调：黑毛+黑风+暗林。氛围：威压、黑风、魔道。。台词："这山是我的，这天下，也该有我黑虎的位置！"。动作帧（动图）：①黑风起 ②虎啸 ③震慑群兽 ④威压四野。诗词：黑风渐起虎王横，威压山林万兽惊。争强好胜踏魔道，一啸风云暗地明。。主题句：虎踞魔道亦有性，只叹一生总在争与执。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：黑毛+黑风+暗林。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 黑风虎王, age 少年, 黑风虎王·威压渐盛; scene: 以黑风震慑群兽；与正道交锋，威压初显。; 黑风渐起虎王横，威压山林万兽惊。争强好胜踏魔道，一啸风云暗地明。; palette: 黑毛+黑风+暗林; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 暗影虎**
- 黑心虎，道法初成阶段·魔教霸主（青年，魔教霸主·争霸江湖）。形象：黑心虎，黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。品性：执掌魔教，争霸江湖。心中执念——一统天下，压过七侠。。姿态：魔教大殿发号施令，黑风裹挟威压；与七侠为敌，执念渐深。。服饰：黑毛猛虎，魔教霸主之姿。。体型：身高约7头身，黑毛猛虎，霸主之姿。。衣物细节：墨黑虎纹，身披暗色铁甲，肩披玄色披风，黑风裹挟，霸主之姿。。发型妆造：墨黑虎毛如焰，王纹深邃，威压赫赫。。脸型五官：虎脸，圆眼含威，王纹深。。武器招式：黑风威压。。功法：黑风威压；魔教霸业。。功法表现：黑风裹挟。。画面：构图：魔教大殿，黑毛猛虎霸主之姿发号施令，黑风裹挟，背景暗色殿堂。色调：黑毛+魔教暗+霸业金。氛围：霸主、魔教、执念。。台词："七侠算什么？我黑心虎，要这江湖，只有一个王！"。动作帧（动图）：①登殿 ②发号施令 ③黑风裹挟 ④执念渐深。诗词：魔教霸主黑虎王，争霸江湖威压张。执念一统压七侠，黑风裹挟战四方。。主题句：虎踞魔道亦有性，只叹一生总在争与执。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：黑毛+魔教暗+霸业金。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 魔教霸主, age 青年, 魔教霸主·争霸江湖; scene: 魔教大殿发号施令，黑风裹挟威压；与七侠为敌，执念渐深。; 魔教霸主黑虎王，争霸江湖威压张。执念一统压七侠，黑风裹挟战四方。; palette: 黑毛+魔教暗+霸业金; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 幽影虎王**
- 黑心虎，大劫淬炼阶段·霸业之争（中年，霸业之争·决战七侠）。形象：黑心虎，黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。品性：与七侠决战，霸业之争。一心想称霸，其实缺个拥抱。。姿态：七剑合璧前，他独战七侠；黑风尽出，心底却空了一块。。服饰：黑毛猛虎，战痕。。体型：身高约7头身，黑毛猛虎，战痕。。衣物细节：墨黑虎纹，暗甲战痕遍布，玄色披风残破，黑风尽出。。发型妆造：墨黑虎毛凌乱，眼底空落。。脸型五官：虎脸，圆眼含威又空落。。武器招式：黑风霸业。。功法：黑风霸业；以一敌七。。功法表现：黑风对决剑光。。画面：构图：战场，黑毛猛虎独对七剑合璧，黑风尽出，眼底却有一丝空落，背景大战。色调：黑毛+黑风+七剑彩光。氛围：霸业、死战、七侠。。台词："我争了一辈子霸业，打赢了这仗，却不知道要给谁看。"。动作帧（动图）：①独战 ②黑风尽出 ③七剑合璧 ④心底空落。诗词：霸业之争战七侠，黑风独对剑光华。一心想赢空落落，原来霸业冷如沙。。主题句：虎踞魔道亦有性，只叹一生总在争与执。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：黑毛+黑风+七剑彩光。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 霸业之争, age 中年, 霸业之争·决战七侠; scene: 七剑合璧前，他独战七侠；黑风尽出，心底却空了一块。; 霸业之争战七侠，黑风独对剑光华。一心想赢空落落，原来霸业冷如沙。; palette: 黑毛+黑风+七剑彩光; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 霸业之虎**
- 黑心虎，封神登天阶段·迷途之虎（中年，迷途之虎·开始怀疑）。形象：黑心虎，黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。品性：大战之后，开始怀疑自己的路。一心想称霸，夜里却会想家。。姿态：魔教大殿空荡，他独坐王座，想起幼时山林的月亮；争了一辈子，忽然不知为了什么。。服饰：黑毛猛虎，神色疲态。。体型：身高约7头身，黑毛猛虎，独坐王座。。衣物细节：墨黑虎纹，暗甲卸下一半，玄色披风垂地，独自望月。。发型妆造：墨黑虎毛，眼神迷惘。。脸型五官：虎脸，圆眼迷惘，望向月光。。武器招式：无，迷途。。功法：霸业动摇；迷途。。功法表现：月下黑影。。画面：构图：魔教大殿，黑毛猛虎独坐王座，月光从窗洒入，神色迷惘，背景空荡殿堂。色调：黑毛+月光白+殿堂暗。氛围：迷途、怀疑、想家。。台词："我打赢了天下，可这天下，怎么这么冷？我到底……在争什么？"。动作帧（动图）：①独坐 ②望月 ③想起山林 ④茫然。诗词：迷途之虎独坐堂，霸业成时空自伤。想起山林旧时月，争来争去一场凉。。主题句：虎踞魔道亦有性，只叹一生总在争与执。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：黑毛+月光白+殿堂暗。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 迷途之虎, age 中年, 迷途之虎·开始怀疑; scene: 魔教大殿空荡，他独坐王座，想起幼时山林的月亮；争了一辈子，忽然不知为了什么。; 迷途之虎独坐堂，霸业成时空自伤。想起山林旧时月，争来争去一场凉。; palette: 黑毛+月光白+殿堂暗; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 暗黑尊者**
- 黑心虎，道果圆满阶段·天地霸业（终极，天地霸业·孤独终局）。形象：黑心虎，黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。品性：魔教覆灭，霸业成空。争了一辈子的天下，最终才懂——反派的心，也想有人来接住。。姿态：最后一战落败，他望着七侠守护的家园，忽然低声道："原来我争的，不是天下，是有人陪我。"。服饰：黑毛猛虎，伤痕累累，神色释然。。体型：身高约7头身，黑毛猛虎，伤痕，立于暮色。。衣物细节：墨黑虎纹，暗甲残破，玄色披风落地，立于暮色。。发型妆造：墨黑虎毛，眼底释然含暖。。脸型五官：虎脸，圆眼释然含暖。。武器招式：无，归真。。功法：霸业尽散；虎性归真。。功法表现：暮色暖光。。画面：构图：暮色山巅，黑毛猛虎伤痕而立，望着七侠守护的家园灯火，神色释然，背景暮色天地。色调：黑毛+暮色橙+家园灯火暖。氛围：终极、霸业、孤独。。台词："虎踞魔道一辈子，最后才懂——我想赢的，从来不是天下，是有人愿意收留我这颗争强的心。"。动作帧（动图）：①立于暮色 ②望家园灯火 ③低头 ④释然。诗词：天地霸业终成空，黑虎伤痕立晚风。争尽天下方知冷，原来所求是相逢。。主题句：虎踞魔道亦有性，只叹一生总在争与执。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：黑毛+暮色橙+家园灯火暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 天地霸业, age 终极, 天地霸业·孤独终局; scene: 最后一战落败，他望着七侠守护的家园，忽然低声道："原来我争的，不是天下，是有人陪我。"; 天地霸业终成空，黑虎伤痕立晚风。争尽天下方知冷，原来所求是相逢。; palette: 黑毛+暮色橙+家园灯火暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 黑小虎（`heixiaohu`） · 人生档案版

**灵胎初醒 · 暗影虎卵**
- 黑小虎，灵胎初醒阶段·虎子初生（幼年，黑心虎之子·宿命起点）。形象：黑小虎，墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。品性：黑心虎之子，一身黑毛，学着父亲的样子，却藏不住少年心。。姿态：跟着父亲习武，进步神速；心里憋着一股劲——要证明给父亲看，他配得上虎王血脉。。服饰：墨黑小虎，爪牙如刃。。体型：身高约4头身，墨黑小虎，健壮。。衣物细节：墨黑幼虎纹，前额"王"纹初显，爪牙初利，少年意气。。发型妆造：墨黑虎毛，双耳立起。。脸型五官：虎脸，圆眼热切，眉间王纹初显。。武器招式：虎爪初利。。功法：天生虎威；爪牙初利。。功法表现：无，少年锐气。。画面：构图：魔教庭院，墨黑小虎学父习武，爪牙初利，眼神热切，背景暗色殿堂。色调：墨黑+魔教暗+少年光。氛围：争强、少年、证明。。台词："爹，你看我这一爪！我黑小虎，要让你骄傲！"。动作帧（动图）：①学父习武 ②一爪 ③热切望父 ④盼夸。诗词：黑虎之子初生横，学父习武少年行。一心想证血脉贵，爪牙初利盼父惊。。主题句：一心想证明自己，却将执念熬成了歧途。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：墨黑+魔教暗+少年光。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 虎子初生, age 幼年, 黑心虎之子·宿命起点; scene: 跟着父亲习武，进步神速；心里憋着一股劲——要证明给父亲看，他配得上虎王血脉。; 黑虎之子初生横，学父习武少年行。一心想证血脉贵，爪牙初利盼父惊。; palette: 墨黑+魔教暗+少年光; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 黑毛小虎**
- 黑小虎，凡尘砺心阶段·黑风小虎（少年，黑风小虎·初遇宿敌）。形象：黑小虎，墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。品性：虎啸带风，与七侠的初见已成宿敌。争强好胜，一心要赢过虹猫。。姿态：与虹猫初战，落败；咬牙立誓——终有一天要赢。。服饰：墨黑小虎，黑风渐起。。体型：身高约5头身，墨黑小虎，黑风绕。。衣物细节：墨黑幼虎纹，暗色短甲初披，黑风渐起。。发型妆造：墨黑虎毛，眼神执拗。。脸型五官：虎脸，圆眼执拗。。武器招式：黑风虎啸。。功法：黑风虎啸；爪风凌厉。。功法表现：黑风初起。。画面：构图：初战之地，墨黑小虎立于黑风中，望向虹猫离去的方向，眼神执拗，背景暮色。色调：墨黑+黑风+暮色。氛围：黑风、宿敌、初成。。台词："虹猫，你等着！总有一天，我黑小虎要堂堂正正赢你！"。动作帧（动图）：①初战 ②落败 ③咬牙 ④立誓。诗词：黑风小虎初遇雄，败于虹猫誓不平。少年执念心中种，一朝要证虎王名。。主题句：一心想证明自己，却将执念熬成了歧途。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：墨黑+黑风+暮色。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 黑风小虎, age 少年, 黑风小虎·初遇宿敌; scene: 与虹猫初战，落败；咬牙立誓——终有一天要赢。; 黑风小虎初遇雄，败于虹猫誓不平。少年执念心中种，一朝要证虎王名。; palette: 墨黑+黑风+暮色; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 黑风小虎**
- 黑小虎，道法初成阶段·少年黑虎（少年，少年黑虎·心中埋疑）。形象：黑小虎，墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。品性：跟着父亲习武，心中却埋着疑惑——父亲的霸业，真是对的吗？。姿态：习武之余，偷看七侠守护的村庄，那里有他没见过的人间烟火；心底的疑惑渐生。。服饰：墨黑小虎，劲装。。体型：身高约6头身，墨黑小虎，健壮。。衣物细节：墨黑虎纹，暗色劲装，爪刃藏锋，望向远处灯火。。发型妆造：墨黑虎毛，眼底疑惑。。脸型五官：虎脸，圆眼疑惑，望向灯火。。武器招式：黑风剑法。。功法：黑风剑法；爪刃。。功法表现：黑风剑光。。画面：构图：魔教高地，墨黑小虎望向远处村庄灯火，神色疑惑，背景冷暖对比。色调：墨黑+村庄暖灯。氛围：少年、疑惑、父影。。台词："爹说的天下，是冷的。可那个村子，好像……有点暖。"。动作帧（动图）：①习武 ②偷望村庄 ③心生疑惑 ④低头。诗词：少年黑虎习父刀，偷望人间烟火高。疑惑心中渐渐起，霸业真比温情豪？。主题句：一心想证明自己，却将执念熬成了歧途。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：墨黑+村庄暖灯。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 少年黑虎, age 少年, 少年黑虎·心中埋疑; scene: 习武之余，偷看七侠守护的村庄，那里有他没见过的人间烟火；心底的疑惑渐生。; 少年黑虎习父刀，偷望人间烟火高。疑惑心中渐渐起，霸业真比温情豪？; palette: 墨黑+村庄暖灯; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 黑小虎·剑**
- 黑小虎，大劫淬炼阶段·黑小虎·谋（青年，毒药暗器·越来越像父亲）。形象：黑小虎，墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。品性：为了证明自己，手段越来越像父亲——毒药、暗器、算计，渐渐迷失。。姿态：以毒药暗器设局；眼看就要赢，却觉得心里越来越空。。服饰：墨黑小虎，暗器在身，神色阴冷。。体型：身高约6头身，墨黑小虎，暗器在身。。衣物细节：墨黑虎纹，暗色劲装，腰间暗器皮囊，神色阴冷。。发型妆造：墨黑虎毛，眼神阴鸷。。脸型五官：虎脸，圆眼阴冷，眼底迷茫。。武器招式：毒药暗器。。功法：毒药暗器；黑风剑法大成。。功法表现：暗器毒光。。画面：构图：暗室，墨黑小虎摆弄毒药暗器，神色阴冷，阴影渐深，背景暗色。色调：墨黑+毒液暗绿+阴影。氛围：手段、迷失、像父。。台词："为了赢，我什么都用上了。可赢的手段越多，我怎么越不像我自己了？"。动作帧（动图）：①设局 ②用毒 ③眼看要赢 ④心底渐空。诗词：为证虎名手段多，毒药暗器设网罗。赢得一时心愈空，迷失来路竟蹉跎。。主题句：一心想证明自己，却将执念熬成了歧途。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：墨黑+毒液暗绿+阴影。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 黑小虎·谋, age 青年, 毒药暗器·越来越像父亲; scene: 以毒药暗器设局；眼看就要赢，却觉得心里越来越空。; 为证虎名手段多，毒药暗器设网罗。赢得一时心愈空，迷失来路竟蹉跎。; palette: 墨黑+毒液暗绿+阴影; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 黑小虎·狂**
- 黑小虎，封神登天阶段·黑小虎·狂（青年，一心要赢·执念成狂）。形象：黑小虎，墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。品性：一心要赢过虹猫，争强好胜到极致。执念熬成狂性，路越走越偏。。姿态：魔教之巅，黑风狂啸；与虹猫再战，赢了的念头已经压倒一切。。服饰：墨黑小虎，狂性尽显，黑风猎猎。。体型：身高约7头身，墨黑小虎，狂性尽显。。衣物细节：墨黑虎纹，暗色劲装，黑风狂啸，少主气势尽显。。发型妆造：墨黑虎毛如焰，王纹扭曲。。脸型五官：虎脸，圆眼狂热，王纹扭曲。。武器招式：黑风狂啸。。功法：黑风狂啸；执念成狂。。功法表现：黑风猎猎。。画面：构图：魔教之巅，墨黑小虎黑风狂啸，神色疯狂，背景狂风暗云。色调：墨黑+狂风暴+暗云。氛围：狂、执念、歧途。。台词："赢！我只要赢！虹猫，我要证明给爹看，我不比你差！"。动作帧（动图）：①登巅 ②黑风狂啸 ③再战虹猫 ④执念压倒。诗词：执念成狂黑小虎，一心要胜虹猫主。黑风狂啸歧途去，忘了他日少年初。。主题句：一心想证明自己，却将执念熬成了歧途。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：墨黑+狂风暴+暗云。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 黑小虎·狂, age 青年, 一心要赢·执念成狂; scene: 魔教之巅，黑风狂啸；与虹猫再战，赢了的念头已经压倒一切。; 执念成狂黑小虎，一心要胜虹猫主。黑风狂啸歧途去，忘了他日少年初。; palette: 墨黑+狂风暴+暗云; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 幽暗少主**
- 黑小虎，道果圆满阶段·呼唤本心（终极，魔王之影·呼唤本心）。形象：黑小虎，墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。品性：最后一战落败，他想起幼时那点少年心——他想证明的，其实不是赢过谁，是有人愿意唤他一声"小虎"。。姿态：落败之际，望着虹猫伸出的手，忽然怔住——原来他要的，从来不是赢，是有人愿意接住他。。服饰：墨黑小虎，伤痕，狂性散尽，露出少年眉眼。。体型：身高约7头身，墨黑小虎，伤痕，少年眉眼。。衣物细节：墨黑虎纹，暗色劲装破损，狂性散尽，露出少年眉眼。。发型妆造：墨黑虎毛，眼底含泪，少年柔软。。脸型五官：虎脸，圆眼含泪，少年柔软。。武器招式：无，本心。。功法：执念散去；本心回归。。功法表现：暮色暖光。。画面：构图：落败之地，墨黑小虎伤痕而立，望向虹猫伸出的手，狂性散尽露出少年眉眼，背景暮色渐暖。色调：墨黑+暮色暖+一缕光。氛围：终极、呼唤、本心。。台词："我想证明了一辈子……原来我只是想，有人愿意唤我一声小虎，说一句"你很好"。"。动作帧（动图）：①落败 ②望伸出的手 ③怔住 ④唤回本心。诗词：魔王之影终散尽，唤得本心少年真。争强半生求一诺，原是想要有人疼。。主题句：一心想证明自己，却将执念熬成了歧途。。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：墨黑+暮色暖+一缕光。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 呼唤本心, age 终极, 魔王之影·呼唤本心; scene: 落败之际，望着虹猫伸出的手，忽然怔住——原来他要的，从来不是赢，是有人愿意接住他。; 魔王之影终散尽，唤得本心少年真。争强半生求一诺，原是想要有人疼。; palette: 墨黑+暮色暖+一缕光; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 马三娘（`ma_sanniang`） · 人生档案版

**灵胎初醒 · 紫云剑意**
- 马三娘，灵胎初醒阶段·魔教剑手（幼年，紫影剑手·卧底之始）。形象：马三娘，紫衣女剑客，眉目清冷，佩剑如影。 核心意象：紫影剑、卧底身份、真假剑心。品性：魔教中善于用剑的紫衣女剑手，被选中执行卧底大计，潜入七侠身边。。姿态：魔教剑场中练紫影剑，剑光如影；受命之日，将真实身份藏进剑里。。服饰：紫衣劲装，眉目清冷，佩剑如影。。体型：身高约6头身，紫衣女剑手，身形清瘦。。衣物细节：紫黑劲装，袖口缀暗金云纹，腰悬紫影剑，一身肃杀。。发型妆造：长发束起高马尾，清冷如霜。。脸型五官：清冷面容，细眉，凤目含剑意，鼻梁挺，薄唇。。武器招式：紫影剑初成。。功法：紫影剑初成；藏锋之术。。功法表现：剑光如影。。画面：构图：魔教剑场，紫衣女剑手练剑，剑光如影，面容清冷，背景暗色殿堂。色调：紫衣+剑影紫+魔教暗。氛围：剑手、卧底之始、藏锋。。台词："这身紫衣，是我唯一的行头。剑，是我唯一的老实话。"。动作帧（动图）：①练剑 ②剑光如影 ③受命 ④藏锋。诗词：魔教剑场紫影横，剑光如影隐真名。受命卧底江湖去，一颗剑心暂藏锋。。主题句：卧底江湖半生，剑心在真假间摇摆——我到底是谁的剑？。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫衣+剑影紫+魔教暗。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 灵胎初醒 魔教剑手, age 幼年, 紫影剑手·卧底之始; scene: 魔教剑场中练紫影剑，剑光如影；受命之日，将真实身份藏进剑里。; 魔教剑场紫影横，剑光如影隐真名。受命卧底江湖去，一颗剑心暂藏锋。; palette: 紫衣+剑影紫+魔教暗; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 剑心初掩**
- 马三娘，凡尘砺心阶段·奉命潜伏（少年，假作侍女·剑藏于拙）。形象：马三娘，紫衣女剑客，眉目清冷，佩剑如影。 核心意象：紫影剑、卧底身份、真假剑心。品性：奉命潜伏，假扮紫云剑主身边的小侍女。剑法出众却刻意藏拙。。姿态：玉蟾宫外伺机接近，装作笨手笨脚的小侍女；剑藏于拙，暗观七侠。。服饰：侍女素裙，眉目温顺，眼底却压着剑光。。体型：身高约6头身，侍女素裙，低眉顺眼。。衣物细节：素裙为饰，腕间暗藏袖剑，发间一缕紫纱，似不经意。。发型妆造：挽作侍女髻，发间斜插紫簪，眼底微光。。脸型五官：温顺面容，眼底微光一闪。。武器招式：藏锋之剑。。功法：伪装术；以拙掩锋。。功法表现：眼底暗光。。画面：构图：玉蟾宫外，紫衣侍女素裙低眉，看似笨拙，眼底却压着一道剑光，背景宫阙。色调：素裙白+玉蟾青+眼底紫。氛围：潜伏、伪装、剑藏于拙。。台词："做侍女，手要笨，眼要亮。我的剑，藏在裙子下面。"。动作帧（动图）：①伺机接近 ②装作笨拙 ③暗观七侠 ④眼底剑光。诗词：奉命潜伏化钗裙，佯装笨拙掩剑心。玉蟾宫外观云动，紫影暗藏待时侵。。主题句：卧底江湖半生，剑心在真假间摇摆——我到底是谁的剑？。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：素裙白+玉蟾青+眼底紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 凡尘砺心 奉命潜伏, age 少年, 假作侍女·剑藏于拙; scene: 玉蟾宫外伺机接近，装作笨手笨脚的小侍女；剑藏于拙，暗观七侠。; 奉命潜伏化钗裙，佯装笨拙掩剑心。玉蟾宫外观云动，紫影暗藏待时侵。; palette: 素裙白+玉蟾青+眼底紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 潜伏剑影**
- 马三娘，道法初成阶段·潜伏玉蟾宫（少年，潜伏宫中·紫影暗度）。形象：马三娘，紫衣女剑客，眉目清冷，佩剑如影。 核心意象：紫影剑、卧底身份、真假剑心。品性：成功潜伏玉蟾宫，接近蓝兔，以待时机。紫影剑法在暗处渐长。。姿态：为蓝兔斟茶递水，暗记宫路径与七侠动静；夜半无人，紫影剑一舞。。服饰：侍女素裙，袖中藏着软剑。。体型：身高约6头身，侍女扮相，袖中藏剑。。衣物细节：素裙之上，斗篷一角压入腰带，袖中软剑贴身，夜深时紫影一闪。。发型妆造：发髻半散，夜半眼底露剑意。。脸型五官：温顺面容，夜半眼底露剑意。。武器招式：紫影剑（暗处）。。功法：潜伏术；紫影剑法渐长。。功法表现：月下紫影。。画面：构图：玉蟾宫月夜，侍女立于回廊，身后一道紫影剑光于暗处一舞，背景宫阙月色。色调：素裙白+玉蟾青+月下紫。氛围：潜伏、紫影、暗度陈仓。。台词："我在玉蟾宫当差的日子越久，越分不清——这宫里的暖，是真是假。"。动作帧（动图）：①斟茶 ②暗记动静 ③夜半舞剑 ④收剑藏袖。诗词：潜伏宫闱紫影深，斟茶暗记七侠音。夜半无人舞一剑，真假之间乱本心。。主题句：卧底江湖半生，剑心在真假间摇摆——我到底是谁的剑？。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：素裙白+玉蟾青+月下紫。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道法初成 潜伏玉蟾宫, age 少年, 潜伏宫中·紫影暗度; scene: 为蓝兔斟茶递水，暗记宫路径与七侠动静；夜半无人，紫影剑一舞。; 潜伏宫闱紫影深，斟茶暗记七侠音。夜半无人舞一剑，真假之间乱本心。; palette: 素裙白+玉蟾青+月下紫; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 卧底江湖**
- 马三娘，大劫淬炼阶段·身份动摇（青年，剑心在真假间摇摆）。形象：马三娘，紫衣女剑客，眉目清冷，佩剑如影。 核心意象：紫影剑、卧底身份、真假剑心。品性：与七侠相处日久，剑心在真假间摇摆。紫云剑主待她如亲，她开始犹豫——卧底还是归心？。姿态：紫云剑主遇险，她本能拔剑相护，才惊觉自己忘了"卧底"二字；夜里，剑放在枕边，一夜无眠。。服饰：紫衣，剑不离身，眉间多了一分纠结。。体型：身高约6头身，紫衣剑客，眉间纠结。。衣物细节：紫黑劲装，斗篷半披，剑不离身，月下剑影半紫半暖。。发型妆造：长发束起，眉间迷惘。。脸型五官：清冷面容，眼底多了一分迷惘。。武器招式：紫影剑法（大成）。。功法：紫影剑法大成；剑心摇摆。。功法表现：剑影半紫半暖。。画面：构图：玉蟾宫夜，紫衣剑客独立，剑横膝上，月光下剑影一半偏紫一半映暖，背景宫阙月色。色调：紫衣+月白+暖烛。氛围：动摇、真假、剑心。。台词："我这剑，本该刺向她们。可为什么，它自己会挡在她们前面？"。动作帧（动图）：①拔剑护主 ②惊觉忘卧底 ③夜半枕剑 ④剑影徘徊。诗词：真假之间剑心摇，护主一剑忘卧槽。夜半枕剑难成寐，紫影徘徊到明朝。。主题句：卧底江湖半生，剑心在真假间摇摆——我到底是谁的剑？。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫衣+月白+暖烛。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 大劫淬炼 身份动摇, age 青年, 剑心在真假间摇摆; scene: 紫云剑主遇险，她本能拔剑相护，才惊觉自己忘了"卧底"二字；夜里，剑放在枕边，一夜无眠。; 真假之间剑心摇，护主一剑忘卧槽。夜半枕剑难成寐，紫影徘徊到明朝。; palette: 紫衣+月白+暖烛; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 魔教剑锋**
- 马三娘，封神登天阶段·身份败露（青年，计谋败露·魔教本性）。形象：马三娘，紫衣女剑客，眉目清冷，佩剑如影。 核心意象：紫影剑、卧底身份、真假剑心。品性：计谋败露，魔教本性显露。她撕下伪装，执紫影剑与七侠反目——卧底半生，终究还是魔教的剑。。姿态：玉蟾宫前，紫影剑指向曾经护她的人；剑在抖，她却咬牙——"我是魔教的人。"。服饰：紫衣猎猎，魔教气息，剑锋尽露。。体型：身高约6头身，紫衣猎猎，剑锋尽露。。衣物细节：紫黑劲装尽展魔教气息，斗篷猎猎，面纱半落，紫影剑锋尽露。。发型妆造：长发披散，眼底含泪却狠戾。。脸型五官：清冷面容，眼底含泪却狠戾。。武器招式：魔教剑锋·紫云夺主。。功法：魔教剑锋；紫云夺主。。功法表现：剑锋紫光尽露。。画面：构图：玉蟾宫前，紫衣剑客撕下伪装，紫影剑指向七侠，剑锋尽露，眼底含泪却咬牙，背景宫阙对峙。色调：紫衣+剑锋紫+对峙冷。氛围：败露、魔教、剑锋。。台词："我演了这么久的好人，连我自己都差点信了。可剑骗不了人——我是魔教的剑。"。动作帧（动图）：①撕下伪装 ②紫影剑出 ③指向七侠 ④咬牙立誓。诗词：计败原形紫影寒，魔教剑锋指旧欢。半生卧底终露相，剑抖心碎强自安。。主题句：卧底江湖半生，剑心在真假间摇摆——我到底是谁的剑？。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫衣+剑锋紫+对峙冷。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 封神登天 身份败露, age 青年, 计谋败露·魔教本性; scene: 玉蟾宫前，紫影剑指向曾经护她的人；剑在抖，她却咬牙——"我是魔教的人。"; 计败原形紫影寒，魔教剑锋指旧欢。半生卧底终露相，剑抖心碎强自安。; palette: 紫衣+剑锋紫+对峙冷; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 真伪对决**
- 马三娘，道果圆满阶段·剑心归处（终极，真伪对决·剑心自问）。形象：马三娘，紫衣女剑客，眉目清冷，佩剑如影。 核心意象：紫影剑、卧底身份、真假剑心。品性：与紫云剑主真伪对决，剑心遭最大拷问。剑落的一刻，她终于明白——剑无正邪，人心才是归处。。姿态：双剑相抵，她望着对面那双真诚的眼睛，忽然收剑，低声道："这一剑，我不刺。"。服饰：紫衣，剑垂于地，神色释然。。体型：身高约6头身，紫衣，剑垂于地。。衣物细节：紫黑劲装，斗篷垂地，面纱已落，剑垂于地，神色释然。。发型妆造：长发披散，眉目释然，眼底含泪而笑。。脸型五官：清冷面容，眉目释然，眼底含泪而笑。。武器招式：收剑，不刺。。功法：剑心归处；真伪自分。。功法表现：剑光消融于月色。。画面：构图：玉蟾宫前月下，紫衣剑客收剑垂地，与紫云剑主相视一笑，剑光消融于月色，背景宫阙与月光。色调：紫衣+月色白+融光暖。氛围：归处、真伪、自问。。台词："我以为我是魔教的剑。可这剑，早就认了她们做归处。我不是谁的剑——我，只是马三娘。"。动作帧（动图）：①双剑相抵 ②望着对方 ③收剑 ④低头一笑。诗词：真伪对决剑心明，双剑相抵自分清。收剑低头一声笑，原来归处是曾经。。主题句：卧底江湖半生，剑心在真假间摇摆——我到底是谁的剑？。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：紫衣+月色白+融光暖。画布规格：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; 道果圆满 剑心归处, age 终极, 真伪对决·剑心自问; scene: 双剑相抵，她望着对面那双真诚的眼睛，忽然收剑，低声道："这一剑，我不刺。"; 真伪对决剑心明，双剑相抵自分清。收剑低头一声笑，原来归处是曾经。; palette: 紫衣+月色白+融光暖; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text
