# 宠物 AI 生图提示词 · 全 125 物种 × 6 阶段（v2）

> 生成：2026-08-07 · 由 `scripts/generate-pet-prompts-v2.js` 自动生成，**不要手改**（改数据源后重跑）
> 取代 v1 `pet-prompts.csv` 通用模板。生成图片后按 `docs/pet-image-manifest.md` 流程登记（MANIFEST + 拷贝到 backend/public/pets/）。

## 一、使用说明

- **中文主句**（每段首行）：复制给即梦 / 通义万相 / Stable Diffusion（中文模型）。
- **EN 行**：复制给 Midjourney / DALL·E / Stable Diffusion（英文模型）。
- **六阶等级制**：灵胎初醒(Lv1) → 凡尘砺心(Lv3) → 道法初成(Lv5) → 大劫淬炼(Lv7) → 封神登天(Lv9) → 道果圆满(Lv11)。
- **精修四维度**：每阶段含 神态 / 动作 / 衣着 / 梳造 四维描述（重点角色手写，其余按大类基座），见 `frontend-vue/src/utils/petRefine.ts`。
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
| 虹猫蓝兔七侠传 | 9 | 2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇 | 2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop |

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
- 烛龙，大劫淬炼阶段·昼夜使者。形象：人面蛇身，通体赤红，口衔火烛。 核心意象：烛火、昼夜、钟山。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身，开眼为昼，闭眼为夜，呼吸之间风雷自生。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 烛龙巡天**
- 烛龙，封神登天阶段·烛龙巡天。形象：人面蛇身，通体赤红，口衔火烛。 核心意象：烛火、昼夜、钟山。神态：龙目洞彻，神威赫赫。动作：绝技大成，行云布雨，开眼为昼，闭眼为夜，呼吸之间风雷自生。衣着：金鳞覆身，祥光万道。梳造：龙角如珊瑚，须垂百丈。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 龙目洞彻，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 开目为昼**
- 烛龙，道果圆满阶段·开目为昼。形象：人面蛇身，通体赤红，口衔火烛。 核心意象：烛火、昼夜、钟山。神态：真身圆满，龙威盖世。动作：腾云驾雾，号令风雨，开眼为昼，闭眼为夜，呼吸之间风雷自生。衣着：金鳞神光，日月同辉。梳造：龙角如岳，须垂星河。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 真身圆满，龙威盖世; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 应龙，大劫淬炼阶段·应龙雏形。形象：身披金鳞，背生双翼的神龙。 核心意象：金鳞、双翼、风雷云雨。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身，双翼一展，云腾雨至，江河为之让路。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 战龙**
- 应龙，封神登天阶段·战龙。形象：身披金鳞，背生双翼的神龙。 核心意象：金鳞、双翼、风雷云雨。神态：龙目洞彻，神威赫赫。动作：绝技大成，行云布雨，双翼一展，云腾雨至，江河为之让路。衣着：金鳞覆身，祥光万道。梳造：龙角如珊瑚，须垂百丈。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 龙目洞彻，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 神龙降世**
- 应龙，道果圆满阶段·神龙降世。形象：身披金鳞，背生双翼的神龙。 核心意象：金鳞、双翼、风雷云雨。神态：真身圆满，龙威盖世。动作：腾云驾雾，号令风雨，双翼一展，云腾雨至，江河为之让路。衣着：金鳞神光，日月同辉。梳造：龙角如岳，须垂星河。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 真身圆满，龙威盖世; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 九尾狐，大劫淬炼阶段·五尾。形象：通体雪白，身后拖九条长尾。 核心意象：九尾、青丘山、月光。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，九尾摇曳如月华流泻，一顾倾人再顾倾国。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 七尾**
- 九尾狐，封神登天阶段·七尾。形象：通体雪白，身后拖九条长尾。 核心意象：九尾、青丘山、月光。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，九尾摇曳如月华流泻，一顾倾人再顾倾国。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 九尾天狐**
- 九尾狐，道果圆满阶段·九尾天狐。形象：通体雪白，身后拖九条长尾。 核心意象：九尾、青丘山、月光。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，九尾摇曳如月华流泻，一顾倾人再顾倾国。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 鲲鹏，大劫淬炼阶段·鹏雏。形象：北冥之巨鱼，可化身为鸟。 核心意象：北冥、垂天之云、九万里长空。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，振翅扶摇九万里，水击三千里，绝云气负青天。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 鲲鹏**
- 鲲鹏，封神登天阶段·鲲鹏。形象：北冥之巨鱼，可化身为鸟。 核心意象：北冥、垂天之云、九万里长空。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，振翅扶摇九万里，水击三千里，绝云气负青天。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 北冥之主**
- 鲲鹏，道果圆满阶段·北冥之主。形象：北冥之巨鱼，可化身为鸟。 核心意象：北冥、垂天之云、九万里长空。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，振翅扶摇九万里，水击三千里，绝云气负青天。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 凤凰，大劫淬炼阶段·凤雏。形象：五彩华羽，尾羽如虹。 核心意象：梧桐、竹实、醴泉。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，五色并举，振翅则百鸟来朝，一鸣则天下太平。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 神凤**
- 凤凰，封神登天阶段·神凤。形象：五彩华羽，尾羽如虹。 核心意象：梧桐、竹实、醴泉。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，五色并举，振翅则百鸟来朝，一鸣则天下太平。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 涅槃之凤**
- 凤凰，道果圆满阶段·涅槃之凤。形象：五彩华羽，尾羽如虹。 核心意象：梧桐、竹实、醴泉。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，五色并举，振翅则百鸟来朝，一鸣则天下太平。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 麒麟，大劫淬炼阶段·踏火麒麟。形象：麋身牛尾，马蹄而一角。 核心意象：独角、麋身、太平之兆。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，不踏生草不履生虫，步生祥云，角有瑞光。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 天麟**
- 麒麟，封神登天阶段·天麟。形象：麋身牛尾，马蹄而一角。 核心意象：独角、麋身、太平之兆。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，不踏生草不履生虫，步生祥云，角有瑞光。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 万瑞之祖**
- 麒麟，道果圆满阶段·万瑞之祖。形象：麋身牛尾，马蹄而一角。 核心意象：独角、麋身、太平之兆。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，不踏生草不履生虫，步生祥云，角有瑞光。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 穷奇，大劫淬炼阶段·穷奇·残暴。形象：状如虎而生双翼，遍体凶煞。 核心意象：虎身双翼、凶煞之气、颠倒黑白。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，振翅扑食，专噬忠善，风随其身。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 穷奇·噬恶**
- 穷奇，封神登天阶段·穷奇·噬恶。形象：状如虎而生双翼，遍体凶煞。 核心意象：虎身双翼、凶煞之气、颠倒黑白。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，振翅扑食，专噬忠善，风随其身。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 四凶之穷奇**
- 穷奇，道果圆满阶段·四凶之穷奇。形象：状如虎而生双翼，遍体凶煞。 核心意象：虎身双翼、凶煞之气、颠倒黑白。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，振翅扑食，专噬忠善，风随其身。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 毕方，大劫淬炼阶段·毕方·衔火。形象：青色单足之鸟，白喙赤足。 核心意象：单足、白喙、青色火焰。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，一足独立，长鸣如鹤，过处即有讹火。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 毕方·焚风**
- 毕方，封神登天阶段·毕方·焚风。形象：青色单足之鸟，白喙赤足。 核心意象：单足、白喙、青色火焰。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，一足独立，长鸣如鹤，过处即有讹火。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 毕方神鸟**
- 毕方，道果圆满阶段·毕方神鸟。形象：青色单足之鸟，白喙赤足。 核心意象：单足、白喙、青色火焰。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，一足独立，长鸣如鹤，过处即有讹火。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 貔貅，大劫淬炼阶段·貔貅·镇宅。形象：龙头马身麟脚，背生双翼。 核心意象：龙头麟脚、双翼、招财纳福。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，张口吸纳四方财气，只进不出，聚而不散。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 貔貅·吞金**
- 貔貅，封神登天阶段·貔貅·吞金。形象：龙头马身麟脚，背生双翼。 核心意象：龙头麟脚、双翼、招财纳福。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，张口吸纳四方财气，只进不出，聚而不散。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 貔貅至尊**
- 貔貅，道果圆满阶段·貔貅至尊。形象：龙头马身麟脚，背生双翼。 核心意象：龙头麟脚、双翼、招财纳福。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，张口吸纳四方财气，只进不出，聚而不散。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 精卫，大劫淬炼阶段·精卫·填海。形象：形如乌，赤首白喙赤足。 核心意象：木石、东海、不灭之志。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，衔西山之木石，日夜往复，以填东海。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 精卫·不屈**
- 精卫，封神登天阶段·精卫·不屈。形象：形如乌，赤首白喙赤足。 核心意象：木石、东海、不灭之志。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，衔西山之木石，日夜往复，以填东海。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 精卫神鸟**
- 精卫，道果圆满阶段·精卫神鸟。形象：形如乌，赤首白喙赤足。 核心意象：木石、东海、不灭之志。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，衔西山之木石，日夜往复，以填东海。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 相柳，大劫淬炼阶段·相柳。形象：九首蛇身，自环其尾。 核心意象：九首、蛇身、水泽之祸。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身，九首齐张，所过之地化为泽国，毒水横流。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 凶兽相柳**
- 相柳，封神登天阶段·凶兽相柳。形象：九首蛇身，自环其尾。 核心意象：九首、蛇身、水泽之祸。神态：龙目洞彻，神威赫赫。动作：绝技大成，行云布雨，九首齐张，所过之地化为泽国，毒水横流。衣着：金鳞覆身，祥光万道。梳造：龙角如珊瑚，须垂百丈。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 龙目洞彻，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 相柳·湮世**
- 相柳，道果圆满阶段·相柳·湮世。形象：九首蛇身，自环其尾。 核心意象：九首、蛇身、水泽之祸。神态：真身圆满，龙威盖世。动作：腾云驾雾，号令风雨，九首齐张，所过之地化为泽国，毒水横流。衣着：金鳞神光，日月同辉。梳造：龙角如岳，须垂星河。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 真身圆满，龙威盖世; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 獬豸，大劫淬炼阶段·獬豸·司法。形象：似羊而独角，目光清正。 核心意象：独角、公正、断狱之威。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，见人争斗，以独角触其不直者，刚正不阿。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 獬豸·执法**
- 獬豸，封神登天阶段·獬豸·执法。形象：似羊而独角，目光清正。 核心意象：独角、公正、断狱之威。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，见人争斗，以独角触其不直者，刚正不阿。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 司法神兽**
- 獬豸，道果圆满阶段·司法神兽。形象：似羊而独角，目光清正。 核心意象：独角、公正、断狱之威。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，见人争斗，以独角触其不直者，刚正不阿。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 青龙，大劫淬炼阶段·腾云龙。形象：青鳞神龙，鹿角长须。 核心意象：青鳞鹿角、东方之位、风调雨顺。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身，腾云驾雾而起，行云布雨泽润四方。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 星宿青龙**
- 青龙，封神登天阶段·星宿青龙。形象：青鳞神龙，鹿角长须。 核心意象：青鳞鹿角、东方之位、风调雨顺。神态：龙目洞彻，神威赫赫。动作：绝技大成，行云布雨，腾云驾雾而起，行云布雨泽润四方。衣着：金鳞覆身，祥光万道。梳造：龙角如珊瑚，须垂百丈。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 龙目洞彻，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 青龙尊者**
- 青龙，道果圆满阶段·青龙尊者。形象：青鳞神龙，鹿角长须。 核心意象：青鳞鹿角、东方之位、风调雨顺。神态：真身圆满，龙威盖世。动作：腾云驾雾，号令风雨，腾云驾雾而起，行云布雨泽润四方。衣着：金鳞神光，日月同辉。梳造：龙角如岳，须垂星河。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 真身圆满，龙威盖世; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 白虎，大劫淬炼阶段·啸林虎。形象：白毛巨虎，金纹隐现。 核心意象：白毛金纹、西方之位、战神之象。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，金纹一闪，虎啸裂空威震百兽。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 镇西神虎**
- 白虎，封神登天阶段·镇西神虎。形象：白毛巨虎，金纹隐现。 核心意象：白毛金纹、西方之位、战神之象。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，金纹一闪，虎啸裂空威震百兽。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 白虎尊者**
- 白虎，道果圆满阶段·白虎尊者。形象：白毛巨虎，金纹隐现。 核心意象：白毛金纹、西方之位、战神之象。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，金纹一闪，虎啸裂空威震百兽。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 朱雀，大劫淬炼阶段·焚风朱雀。形象：赤红神鸟，周身烈焰。 核心意象：赤羽烈焰、南方之位、浴火重生。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，振翅浴火而起，烈焰燎空化为凤凰。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 浴火朱雀**
- 朱雀，封神登天阶段·浴火朱雀。形象：赤红神鸟，周身烈焰。 核心意象：赤羽烈焰、南方之位、浴火重生。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，振翅浴火而起，烈焰燎空化为凤凰。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 朱雀尊者**
- 朱雀，道果圆满阶段·朱雀尊者。形象：赤红神鸟，周身烈焰。 核心意象：赤羽烈焰、南方之位、浴火重生。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，振翅浴火而起，烈焰燎空化为凤凰。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 玄武，大劫淬炼阶段·驮山龟。形象：龟蛇合体，背负硬壳。 核心意象：龟蛇合体、北冥之水、玄武星宿。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身，龟蛇盘踞，身周四象水光流转，一尾横扫千军。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 龙威炽烈，怒目电光; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 镇水玄武**
- 玄武，封神登天阶段·镇水玄武。形象：龟蛇合体，背负硬壳。 核心意象：龟蛇合体、北冥之水、玄武星宿。神态：龙目洞彻，神威赫赫。动作：绝技大成，行云布雨，龟蛇盘踞，身周四象水光流转，一尾横扫千军。衣着：金鳞覆身，祥光万道。梳造：龙角如珊瑚，须垂百丈。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 龙目洞彻，神威赫赫; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 玄武尊者**
- 玄武，道果圆满阶段·玄武尊者。形象：龟蛇合体，背负硬壳。 核心意象：龟蛇合体、北冥之水、玄武星宿。神态：真身圆满，龙威盖世。动作：腾云驾雾，号令风雨，龟蛇盘踞，身周四象水光流转，一尾横扫千军。衣着：金鳞神光，日月同辉。梳造：龙角如岳，须垂星河。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 真身圆满，龙威盖世; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 饕餮，大劫淬炼阶段·无餍之兽。形象：巨口凶兽，双角狰狞。 核心意象：巨口、狰狞双角、贪食之纹。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，巨口一张，吞尽万物永无餍足。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 贪噬之相**
- 饕餮，封神登天阶段·贪噬之相。形象：巨口凶兽，双角狰狞。 核心意象：巨口、狰狞双角、贪食之纹。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，巨口一张，吞尽万物永无餍足。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 饕餮尊者**
- 饕餮，道果圆满阶段·饕餮尊者。形象：巨口凶兽，双角狰狞。 核心意象：巨口、狰狞双角、贪食之纹。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，巨口一张，吞尽万物永无餍足。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 白泽，大劫淬炼阶段·万卷兽。形象：白身独角，通体智慧纹。 核心意象：白身独角、智慧纹、通晓万物。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，开口便能道出天下妖魅鬼怪之名。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：上古大劫降临，水火风雷交加，异兽浴劫淬炼，直面生死蜕变。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：玄墨黑（#1A0A0A）主调 + 血赤（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; trial by elements; 眸光如电，威严中带着坚韧; palette #1A0A0A with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 万物之书**
- 白泽，封神登天阶段·万物之书。形象：白身独角，通体智慧纹。 核心意象：白身独角、智慧纹、通晓万物。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，开口便能道出天下妖魅鬼怪之名。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：神通大成，受万民祭祀封神，身绕鎏金神光，威仪震慑四方。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：鎏金（#D4AF37）主调 + 朱紫（#7B2E8B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; divine coronation; 神兽威严，目光洞彻九幽; palette #D4AF37 with #7B2E8B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 白泽尊者**
- 白泽，道果圆满阶段·白泽尊者。形象：白身独角，通体智慧纹。 核心意象：白身独角、智慧纹、通晓万物。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，开口便能道出天下妖魅鬼怪之名。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：图腾显圣，山海为印，洪荒万古唯此一尊。风格：新中式上古神话国风插画，融合青铜器纹样与敦煌岩彩质感，洪荒神秘、庄严恢弘。色彩：云白（#F5F0E8）主调 + 淡金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese ancient mythology illustration, xinzhongshi, bronze ritual patterns, Dunhuang fresco texture; eternal harmony; 瑞气圆满，祥光普照; palette #F5F0E8 with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 2. 东方神话（20 物种）

> **风格**：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。**阶段演绎**：
> - 灵胎初醒：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘（玄青/仙金）
> - 凡尘砺心：离开洞府行走人间，红尘历练，磨砺道心（黛青/灰白）
> - 道法初成：道法初成，法宝显威，法相庄严初现（朱砂/明黄）
> - 大劫淬炼：天劫淬炼、心魔试炼，仙体历经大劫而不灭（墨褐/血绛）
> - 封神登天：功行圆满封神登天，位列仙班，金光紫气环绕（紫金/御金）
> - 道果圆满：道果圆满，紫气东来，万法皆通证道果（素白/淡紫）

#### 姜子牙（`jiang_ziya`）

**灵胎初醒 · 灵光种子**
- 姜子牙，灵胎初醒阶段·灵光种子。初始形态：一缕灵光凝成垂钓老者的剪影，渭水之畔愿者上钩，封神之榜在其身侧浮沉。水属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：灵光种子中的垂钓老影。梳造：白发，斗笠未戴。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 垂钓渭水**
- 姜子牙，凡尘砺心阶段·垂钓渭水。形象：白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：渭水边布衣老翁，直钩垂钓。梳造：白发苍苍，束发布巾。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 悟道下山**
- 姜子牙，道法初成阶段·悟道下山。形象：白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，渭水直钩垂钓，愿者上钩。衣着：杏黄旗初展，道袍素净。梳造：道髻高束，眉目慈和。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 杏黄旗展**
- 姜子牙，大劫淬炼阶段·杏黄旗展。形象：白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，渭水直钩垂钓，愿者上钩。衣着：封神榜启，杏黄旗挥。梳造：道冠峨立，仙风道骨。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 封神榜启**
- 姜子牙，封神登天阶段·封神榜启。形象：白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，渭水直钩垂钓，愿者上钩。衣着：姜太公临，金甲法衣。梳造：仙髻法冠，打神鞭在手。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 姜太公临**
- 姜子牙，道果圆满阶段·姜太公临。形象：白发道袍，手持杏黄旗，背负封神榜。 核心意象：封神榜、杏黄旗、打神鞭、四不像。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，渭水直钩垂钓，愿者上钩。衣着：封神执榜人，坐于云端。梳造：白发如银，慈眉垂目。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 哪吒（`nezha`）

**灵胎初醒 · 灵珠种子**
- 哪吒，灵胎初醒阶段·灵珠种子。初始形态：一团赤红灵珠之光，火尖枪的锋芒在其中跳跃，莲花初生的轮廓若隐若现。火属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：灵珠之光中的孩童轮廓，红肚兜若隐若现。梳造：总角发髻，呆毛微翘。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 陈塘哪吒**
- 哪吒，凡尘砺心阶段·陈塘哪吒。形象：莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：红肚兜赤足孩童，混天绫缠在腕间。梳造：双总角发髻，额前齐刘海。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 闹海屠龙**
- 哪吒，道法初成阶段·闹海屠龙。形象：莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，脚踏风火轮，手掷乾坤圈，一声"变"。衣着：莲花战甲初成，乾坤圈套在颈间。梳造：发髻渐紧，束发带飞扬。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 莲花化身**
- 哪吒，大劫淬炼阶段·莲花化身。形象：莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，脚踏风火轮，手掷乾坤圈，一声"变"。衣着：三头六臂初现，火尖枪在手，战甲燃焰。梳造：六臂之姿，发丝燃着火光。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 三太子**
- 哪吒，封神登天阶段·三太子。形象：莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，脚踏风火轮，手掷乾坤圈，一声"变"。衣着：三头六臂全开，莲花金甲，风火轮踏焰。梳造：金冠束发，英姿勃发。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 肉身成圣**
- 哪吒，道果圆满阶段·肉身成圣。形象：莲花化身灵童子，火尖枪挑乾坤。 核心意象：风火轮、乾坤圈、混天绫、火尖枪、莲花。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，脚踏风火轮，手掷乾坤圈，一声"变"。衣着：三头六臂莲花战神，金甲全开。梳造：金冠束发，六臂威仪。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 杨戬（`yang_jian`）

**灵胎初醒 · 天眼灵光**
- 杨戬，灵胎初醒阶段·天眼灵光。初始形态：一道天眼灵光，银甲虚影与额间竖瞳的星芒若隐若现，二郎神威初显。金属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：天眼灵光中的银甲虚影。梳造：束发未冠，额间天目微闭。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年二郎**
- 杨戬，凡尘砺心阶段·幼年二郎。形象：银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：灌江口少年，布衣素服。梳造：束发利落，英气初显。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 八九玄功**
- 杨戬，道法初成阶段·八九玄功。形象：银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，额间天目一睁，照破妖邪原形。衣着：银甲初着，三尖两刃刀在手。梳造：银冠束发，天目渐睁。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 天眼洞开**
- 杨戬，大劫淬炼阶段·天眼洞开。形象：银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，额间天目一睁，照破妖邪原形。衣着：八九玄功激荡，战袍猎猎。梳造：发丝凌乱，天目怒张。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 二郎真君**
- 杨戬，封神登天阶段·二郎真君。形象：银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，额间天目一睁，照破妖邪原形。衣着：二郎真君银甲金冠，哮天犬随侧。梳造：银盔高束，天目如电。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 玄功大成**
- 杨戬，道果圆满阶段·玄功大成。形象：银甲神将，额生天目，手持三尖两刃刀。 核心意象：三尖两刃刀、哮天犬、天眼、银甲。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，额间天目一睁，照破妖邪原形。衣着：清源妙道真君，玄功大成。梳造：束发道冠，天目内敛。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 雷震子（`lei_zhenzi`）

**灵胎初醒 · 雷卵种子**
- 雷震子，灵胎初醒阶段·雷卵种子。初始形态：一枚雷纹灵卵，紫电纹路在蛋壳上蜿蜒，风雷双翅的雏形在雷光中孕育。雷属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：雷卵灵光中的雏鸟轮廓。梳造：无，蛋中初形。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年雷仔**
- 雷震子，凡尘砺心阶段·幼年雷仔。形象：鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：幼年孩童，青布短衣。梳造：总角发髻，天真。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 食杏生翅**
- 雷震子，道法初成阶段·食杏生翅。形象：鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，双翅一振扶摇千里，黄金棍引下雷霆万钧。衣着：食杏果后鸟面初显，身生细羽。梳造：发间生出细羽，青面初现。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 黄金棍**
- 雷震子，大劫淬炼阶段·黄金棍。形象：鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，双翅一振扶摇千里，黄金棍引下雷霆万钧。衣着：鸟面雷公之相，背生双翅，黄金棍在手。梳造：风雷双翅，发如闪电。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 救父出关**
- 雷震子，封神登天阶段·救父出关。形象：鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，双翅一振扶摇千里，黄金棍引下雷霆万钧。衣着：雷公真身，风雷战甲，金棍引雷。梳造：雷霆之发，翼展千里。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 风雷显圣**
- 雷震子，道果圆满阶段·风雷显圣。形象：鸟面人身，青面獠牙，背生双翼。 核心意象：黄金棍、风雷双翅、雷声。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，双翅一振扶摇千里，黄金棍引下雷霆万钧。衣着：风雷显圣，忠义神将之姿。梳造：雷冠神相，威而不怒。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 黄天化（`huang_tianhua`）

**灵胎初醒 · 将门灵种**
- 黄天化，灵胎初醒阶段·将门灵种。初始形态：一粒将门灵种，玉麒麟的蹄影与莫邪剑的寒光在土色灵光中交叠。土属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：将门灵种，虎子之气的灵光。梳造：总角发髻，稚气。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年天化**
- 黄天化，凡尘砺心阶段·幼年天化。形象：青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：少年将军，素色劲装。梳造：束发带冠，少年英气。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 火龙标**
- 黄天化，道法初成阶段·火龙标。形象：青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，胯下玉麒麟冲阵，莫邪剑出鞘如虹。衣着：青布战袍，莫邪剑初佩。梳造：束发簪冠，眉目英武。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 下山辅周**
- 黄天化，大劫淬炼阶段·下山辅周。形象：青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，胯下玉麒麟冲阵，莫邪剑出鞘如虹。衣着：玉麒麟踏阵，莫邪剑出鞘。梳造：战盔斜戴，血战之姿。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 将门虎子**
- 黄天化，封神登天阶段·将门虎子。形象：青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，胯下玉麒麟冲阵，莫邪剑出鞘如虹。衣着：少年将军金甲，玉麒麟为骑。梳造：金冠束发，英姿飒爽。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 莫邪封神**
- 黄天化，道果圆满阶段·莫邪封神。形象：青面红发，手持莫邪宝剑，胯下玉麒麟。 核心意象：莫邪宝剑、玉麒麟、攒心钉。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，胯下玉麒麟冲阵，莫邪剑出鞘如虹。衣着：封神之姿，少年英雄的永恒。梳造：神光化冠，风华长存。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 土行孙（`tu_xingsun`）

**灵胎初醒 · 地灵种子**
- 土行孙，灵胎初醒阶段·地灵种子。初始形态：一粒地灵种子，黄土灵光凝成遁形之影，镔铁棍的轮廓埋在地脉之中。土属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：地灵种子，土色灵光。梳造：矮小身形，发丝粗短。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年土行**
- 土行孙，凡尘砺心阶段·幼年土行。形象：矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：矮小少年，粗布短衣。梳造：平头短发，机灵。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 土遁神通**
- 土行孙，道法初成阶段·土遁神通。形象：矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，就地一滚没入黄土，遁地千里倏忽即至。衣着：镔铁棍在手，土遁初成。梳造：发间沾土，狡黠。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 地行如飞**
- 土行孙，大劫淬炼阶段·地行如飞。形象：矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，就地一滚没入黄土，遁地千里倏忽即至。衣着：土遁如飞，捆仙绳藏袖。梳造：灰头土脸，笑意狡黠。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 入地为营**
- 土行孙，封神登天阶段·入地为营。形象：矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，就地一滚没入黄土，遁地千里倏忽即至。衣着：地行仙之姿，镔铁棍如臂使指。梳造：短须微蓄，仙气内藏。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 地行仙归**
- 土行孙，道果圆满阶段·地行仙归。形象：矮小身形，手持镔铁棍，擅遁地术。 核心意象：镔铁棍、捆仙绳、土遁。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，就地一滚没入黄土，遁地千里倏忽即至。衣着：地行仙圆满，镔铁棍破土裂山。梳造：短须如戟，仙威内藏。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 杨任（`yang_ren`）

**灵胎初醒 · 忠臣灵光**
- 杨任，灵胎初醒阶段·忠臣灵光。初始形态：一缕忠臣灵光，双掌之目在光影中睁开，掌中星芒是剖目重见天日的光。木属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：忠臣灵光，朝服虚影。梳造：束发官帽，端肃。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年杨任**
- 杨任，凡尘砺心阶段·幼年杨任。形象：双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：商朝上大夫，绯色官袍。梳造：冠冕齐整，正直。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 被挖双目**
- 杨任，道法初成阶段·被挖双目。形象：双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，双手一抬，掌中之目扫视天地阴阳。衣着：被剜双目，白巾覆眼。梳造：发丝散乱，悲怆。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 眼中生手**
- 杨任，大劫淬炼阶段·眼中生手。形象：双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，双手一抬，掌中之目扫视天地阴阳。衣着：掌中双目初睁，飞电枪在手。梳造：道袍布冠，掌目生辉。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 飞电枪出**
- 杨任，封神登天阶段·飞电枪出。形象：双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，双手一抬，掌中之目扫视天地阴阳。衣着：掌中目仙，云霞兽为骑。梳造：道冠仙髻，目光如炬。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 道德真君**
- 杨任，道果圆满阶段·道德真君。形象：双目长于掌中，手持飞电枪，坐骑云霞兽。 核心意象：掌中目、飞电枪、云霞兽。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，双手一抬，掌中之目扫视天地阴阳。衣着：道德真君之姿，掌目洞彻天地。梳造：仙光化发，慈严并具。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 韦护（`wei_hu`）

**灵胎初醒 · 金刚灵种**
- 韦护，灵胎初醒阶段·金刚灵种。初始形态：一粒金刚灵种，金甲虚影与降魔杵的宝光交错，护法威严初凝。金属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：金刚灵种，佛光初凝。梳造：短发，僧俗未定。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年韦护**
- 韦护，凡尘砺心阶段·幼年韦护。形象：身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：护法少年，素衣。梳造：束发，端正。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 降魔杵**
- 韦护，道法初成阶段·降魔杵。形象：身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，降魔杵往下一杵，万邪辟易。衣着：金甲初着，降魔杵在手。梳造：发冠，威仪渐显。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 金刚怒目**
- 韦护，大劫淬炼阶段·金刚怒目。形象：身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，降魔杵往下一杵，万邪辟易。衣着：护法金刚相，金甲怒目。梳造：怒目圆睁，法相庄严。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 护法大神**
- 韦护，封神登天阶段·护法大神。形象：身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，降魔杵往下一杵，万邪辟易。衣着：韦陀护法，降魔杵镇邪。梳造：金冠宝相，慈悲怒目。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 韦护尊者**
- 韦护，道果圆满阶段·韦护尊者。形象：身披金甲，手持降魔杵，法相庄严。 核心意象：降魔杵、金甲、韦陀像。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，降魔杵往下一杵，万邪辟易。衣着：三教护法，金刚不坏。梳造：神光化髻，万邪辟易。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 妲己（`daji`）

**灵胎初醒 · 妖狐之卵**
- 妲己，灵胎初醒阶段·妖狐之卵。初始形态：一枚妖狐之卵，绯红蛋壳泛着魅惑幽光，九尾虚影在狐火中摇曳。暗属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：妖狐之卵，绯红幽光。梳造：狐尾虚影蜷绕。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼年妖狐**
- 妲己，凡尘砺心阶段·幼年妖狐。形象：千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：幼年狐妖，素衫。梳造：青丝如墨，狐尾初藏。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 入宫惑主**
- 妲己，道法初成阶段·入宫惑主。形象：千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，长袖轻舞，回眸一笑百媚生，九尾一现狐火冲天。衣着：入宫华服，宫装初着。梳造：云髻初盘，珠钗点缀。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 炮烙忠臣**
- 妲己，大劫淬炼阶段·炮烙忠臣。形象：千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，长袖轻舞，回眸一笑百媚生，九尾一现狐火冲天。衣着：九尾华裳，狐火隐现。梳造：高髻金步摇，媚眼如丝。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 九尾祸世**
- 妲己，封神登天阶段·九尾祸世。形象：千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，长袖轻舞，回眸一笑百媚生，九尾一现狐火冲天。衣着：九尾祸世，长裙曳地。梳造：云髻高绾，狐尾尽展。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 妖狐末路**
- 妲己，道果圆满阶段·妖狐末路。形象：千年九尾狐化身绝色美人，狐尾隐现。 核心意象：九尾、狐火、摘星楼、酒池肉林。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，长袖轻舞，回眸一笑百媚生，九尾一现狐火冲天。衣着：九尾齐天，祸乱天下。梳造：九尾尽展，魅惑万世。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 申公豹（`shen_gongbao`）

**灵胎初醒 · 邪道灵光**
- 申公豹，灵胎初醒阶段·邪道灵光。初始形态：一缕邪道灵光，双剑虚影与黑虎之息交织，三寸舌的巧言在风中回响。风属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：邪道灵光，黑袍虚影。梳造：束发，面带谄笑。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 阐教入门**
- 申公豹，凡尘砺心阶段·阐教入门。形象：黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：阐教弟子，青白道袍。梳造：道髻初束，眼神闪烁。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 叛教下山**
- 申公豹，道法初成阶段·叛教下山。形象：黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，一句"道友请留步"，三寸舌说动天下英雄。衣着：叛教下山，玄黑道袍。梳造：道髻斜歪，黑虎随行。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 请道友下山**
- 申公豹，大劫淬炼阶段·请道友下山。形象：黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，一句"道友请留步"，三寸舌说动天下英雄。衣着：黑虎为骑，双剑在手。梳造：乱发披肩，三寸舌如刀。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 五岳三山**
- 申公豹，封神登天阶段·五岳三山。形象：黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，一句"道友请留步"，三寸舌说动天下英雄。衣着：说客之相，五岳三山皆访。梳造：长须垂胸，眼神阴鸷。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 封神台末**
- 申公豹，道果圆满阶段·封神台末。形象：黑袍道人，身负双剑，面容阴鸷。 核心意象：双剑、黑虎、开天珠、三寸舌。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，一句"道友请留步"，三寸舌说动天下英雄。衣着：翻江倒海，三寸舌搅动乾坤。梳造：长须飞扬，舌灿莲花。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 孙悟空（`sun_wukong`）

**灵胎初醒 · 灵根育孕**
- 孙悟空，灵胎初醒阶段·灵根育孕。初始形态：一粒灵石灵根，花果山仙气凝成石卵，石中孕育着金毛猴王的轮廓。土属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：灵石中的石猴轮廓，身裹混沌仙气。梳造：石纹未褪，通体金毛初生。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 花果山王**
- 孙悟空，凡尘砺心阶段·花果山王。形象：通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：花果山的小石猴，藤叶为衣。梳造：金色猴毛，红冠未成。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 筋斗云**
- 孙悟空，道法初成阶段·筋斗云。形象：通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，一个筋斗云，或吹毛化猴。衣着：虎皮裙初穿，金箍棒在手。梳造：金毛如焰，凤翅紫金冠初戴。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 大闹天宫**
- 孙悟空，大劫淬炼阶段·大闹天宫。形象：通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，一个筋斗云，或吹毛化猴。衣着：大闹天宫的战甲，锁子黄金甲。梳造：紫金冠斜戴，战意凛然。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 五百年劫**
- 孙悟空，封神登天阶段·五百年劫。形象：通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，一个筋斗云，或吹毛化猴。衣着：齐天大圣冠冕，金甲红袍。梳造：凤翅紫金冠，金毛猎猎。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 斗战胜佛**
- 孙悟空，道果圆满阶段·斗战胜佛。形象：通体金毛石猴，火眼金睛，金箍棒在手。 核心意象：金箍棒、筋斗云、虎皮裙、花果山。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，一个筋斗云，或吹毛化猴。衣着：斗战胜佛袈裟，佛光内蕴。梳造：金毛归于平静，项间念珠。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 吕洞宾（`lv_dongbin`）

**灵胎初醒 · 黄粱书生**
- 吕洞宾，灵胎初醒阶段·黄粱书生。初始形态：一位黄粱书生剪影，青衫负剑立于道袍灵光中，纯阳剑气如露未凝。金属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：黄粱书生，青衫负剑。梳造：书生巾，鬓发如裁。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 弃儒从道**
- 吕洞宾，凡尘砺心阶段·弃儒从道。形象：白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：儒生装束，弃儒从道。梳造：束发，眉目清俊。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 御剑飞行**
- 吕洞宾，道法初成阶段·御剑飞行。形象：白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，酒葫芦一倾，纯阳剑出鞘，飞剑千里斩妖。衣着：灰白道袍，纯阳剑在手。梳造：道髻，逍遥巾束发。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 五雷天心**
- 吕洞宾，大劫淬炼阶段·五雷天心。形象：白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，酒葫芦一倾，纯阳剑出鞘，飞剑千里斩妖。衣着：五雷天心，道袍猎猎。梳造：长发飞扬，剑眉入鬓。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 点石成金**
- 吕洞宾，封神登天阶段·点石成金。形象：白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，酒葫芦一倾，纯阳剑出鞘，飞剑千里斩妖。衣着：纯阳帝君，白衣佩剑。梳造：仙髻金冠，剑仙风姿。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 纯阳帝君**
- 吕洞宾，道果圆满阶段·纯阳帝君。形象：白衣佩剑的俊逸剑仙，逍遥巾束发，腰悬酒葫芦。 核心意象：纯阳剑、酒葫芦、黄粱梦、岳阳楼。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，酒葫芦一倾，纯阳剑出鞘，飞剑千里斩妖。衣着：纯阳剑仙，飞剑千里斩妖。梳造：仙髻金冠，剑意冲霄。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 何仙姑（`he_xiangu`）

**灵胎初醒 · 岭南采药女**
- 何仙姑，灵胎初醒阶段·岭南采药女。初始形态：一位采药女虚影，竹篓药香凝成青雾，莲花在足下缓缓绽放。木属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：采药女虚影，竹篓在背。梳造：青丝，荆钗布裙。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 仙桃轻身**
- 何仙姑，凡尘砺心阶段·仙桃轻身。形象：素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：岭南采药女，素裙。梳造：发辫垂肩，灵秀。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 八仙之女**
- 何仙姑，道法初成阶段·八仙之女。形象：素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，玉手轻拈一朵莲花，抛起便是漫天仙光。衣着：八仙之女，云母为饰。梳造：云髻初盘，莲花簪。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 百花阵**
- 何仙姑，大劫淬炼阶段·百花阵。形象：素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，玉手轻拈一朵莲花，抛起便是漫天仙光。衣着：百花阵启，仙裙生辉。梳造：云髻高绾，花钿点额。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 八仙过海**
- 何仙姑，封神登天阶段·八仙过海。形象：素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，玉手轻拈一朵莲花，抛起便是漫天仙光。衣着：八仙过海，素衣仙袍。梳造：仙髻凤钗，容光摄人。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 瑶池仙姑**
- 何仙姑，道果圆满阶段·瑶池仙姑。形象：素衣仙姑，手持如意莲花，长发及腰。 核心意象：如意莲花、荷叶、云母、仙桃。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，玉手轻拈一朵莲花，抛起便是漫天仙光。衣着：瑶池仙姑，莲花在手。梳造：素发仙髻，不染凡尘。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 张果老（`zhang_guolao`）

**灵胎初醒 · 恒山樵夫**
- 张果老，灵胎初醒阶段·恒山樵夫。初始形态：一位恒山樵夫虚影，背负柴薪倒骑白驴，鱼鼓之声在灵光中轻响。土属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：樵夫虚影，柴担在肩。梳造：白发，粗布巾。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 轻身之能**
- 张果老，凡尘砺心阶段·轻身之能。形象：灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：恒山樵夫，麻布短褐。梳造：白发苍然，腰间酒葫芦。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 鱼鼓神通**
- 张果老，道法初成阶段·鱼鼓神通。形象：灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，倒骑白驴，鱼鼓一拍，声定风波。衣着：倒骑白驴，鱼鼓在手。梳造：白发束起，仙翁之气。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 呼风唤雨**
- 张果老，大劫淬炼阶段·呼风唤雨。形象：灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，倒骑白驴，鱼鼓一拍，声定风波。衣着：呼风唤雨，道袍翻飞。梳造：白发飘逸，仙骨道姿。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 鱼鼓天音**
- 张果老，封神登天阶段·鱼鼓天音。形象：灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，倒骑白驴，鱼鼓一拍，声定风波。衣着：鱼鼓天音，混元仙翁。梳造：仙髻道冠，银须飘飘。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 中条洞主**
- 张果老，道果圆满阶段·中条洞主。形象：灰白道袍老翁，倒骑白驴，手持鱼鼓。 核心意象：白驴、鱼鼓、蝙蝠、恒山。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，倒骑白驴，鱼鼓一拍，声定风波。衣着：混元仙翁，鱼鼓天音震三界。梳造：银发仙髻，倒骑驴巡天。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 铁拐李（`tie_guaili`）

**灵胎初醒 · 巴国隐士**
- 铁拐李，灵胎初醒阶段·巴国隐士。初始形态：一位巴国隐士虚影，铁拐顿地破衣飘摇，葫芦灵光纳着一方天地。水属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：巴国隐士虚影，铁拐顿地。梳造：乱发，破帽。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 弃儒从道**
- 铁拐李，凡尘砺心阶段·弃儒从道。形象：跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：弃儒从道，布衣。梳造：束发，清癯。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 附身乞丐**
- 铁拐李，道法初成阶段·附身乞丐。形象：跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，铁拐一顿，宝葫芦口一开，乾坤皆入其中。衣着：附身乞丐，破衣跛足。梳造：乱发蓬面，铁拐为杖。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 混元丹术**
- 铁拐李，大劫淬炼阶段·混元丹术。形象：跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，铁拐一顿，宝葫芦口一开，乾坤皆入其中。衣着：混元丹术，破袍仙光。梳造：发须皆乱，目藏精光。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 葫芦纳天地**
- 铁拐李，封神登天阶段·葫芦纳天地。形象：跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，铁拐一顿，宝葫芦口一开，乾坤皆入其中。衣着：铁拐仙人，宝葫芦纳天地。梳造：乱发仙髻，破中藏道。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 铁拐大仙**
- 铁拐李，道果圆满阶段·铁拐大仙。形象：跛足乞丐形象，铁拐为杖，腰悬宝葫芦。 核心意象：铁拐、宝葫芦、破衣、跛足。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，铁拐一顿，宝葫芦口一开，乾坤皆入其中。衣着：铁拐大仙，破衣而入道。梳造：形残道全，宝葫芦常挂。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 汉钟离（`han_zhongli`）

**灵胎初醒 · 燕台武将**
- 汉钟离，灵胎初醒阶段·燕台武将。初始形态：一位燕台武将虚影，红袍赤面战意未消，芭蕉扇上焰光初燃。火属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：武将虚影，红袍未褪。梳造：束发，鬓角霜。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 弃武修道**
- 汉钟离，凡尘砺心阶段·弃武修道。形象：赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：燕台武将，玄甲红袍。梳造：束发盔缨，战意未消。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 芭蕉扇**
- 汉钟离，道法初成阶段·芭蕉扇。形象：赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，芭蕉扇一扇，火海滔天，再一扇清风送人归。衣着：弃武修道，道袍初着。梳造：束发，芭蕉扇藏袖。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 炎龙召唤**
- 汉钟离，大劫淬炼阶段·炎龙召唤。形象：赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，芭蕉扇一扇，火海滔天，再一扇清风送人归。衣着：赤阳神功，红袍大氅。梳造：赤发虬髯，目若朗星。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 赤阳神功**
- 汉钟离，封神登天阶段·赤阳神功。形象：赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，芭蕉扇一扇，火海滔天，再一扇清风送人归。衣着：正阳帝君，红袍金冠。梳造：金冠束发，赤面威严。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 正阳帝君**
- 汉钟离，道果圆满阶段·正阳帝君。形象：赤面红袍，手持芭蕉扇，腰系宝葫芦。 核心意象：芭蕉扇、宝葫芦、红袍、大日。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，芭蕉扇一扇，火海滔天，再一扇清风送人归。衣着：正阳帝君，大日金轮镇山河。梳造：赤发金冠，仙光圆满。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 蓝采和（`lan_caihe`）

**灵胎初醒 · 濠梁乐童**
- 蓝采和，灵胎初醒阶段·濠梁乐童。初始形态：一位踏歌乐童虚影，破绿衫赤足迎风，玉板轻拍溅起歌谣般的灵光。风属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：乐童虚影，破绿衫。梳造：赤足，散发。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 歌中有道**
- 蓝采和，凡尘砺心阶段·歌中有道。形象：破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：市井歌者，绿衫赤足。梳造：散发，笑意洒脱。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 踏歌仙人**
- 蓝采和，道法初成阶段·踏歌仙人。形象：破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，玉板一拍，踏歌而行，一步一莲花。衣着：踏歌而行，玉板在手。梳造：发丝飞扬，踏歌自若。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 踏歌成阵**
- 蓝采和，大劫淬炼阶段·踏歌成阵。形象：破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，玉板一拍，踏歌而行，一步一莲花。衣着：踏歌成阵，歌声破空。梳造：散发狂歌，仙气外溢。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 九天玄音**
- 蓝采和，封神登天阶段·九天玄音。形象：破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，玉板一拍，踏歌而行，一步一莲花。衣着：九天玄音，歌仙之姿。梳造：发髻松散，逍遥不羁。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 濠梁仙翁**
- 蓝采和，道果圆满阶段·濠梁仙翁。形象：破绿衫赤足的清秀少年，手持玉板，腰挂铜钱。 核心意象：玉板、破绿衫、铜钱、歌谣。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，玉板一拍，踏歌而行，一步一莲花。衣着：濠梁仙翁，赤足高歌。梳造：白发作歌，笑看红尘。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 曹国舅（`cao_guojiu`）

**灵胎初醒 · 锦绣公子**
- 曹国舅，灵胎初醒阶段·锦绣公子。初始形态：一位锦绣公子虚影，锦袍玉带随灵光浮动，玉磬之音如金玉轻鸣。金属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：公子虚影，锦袍玉带。梳造：束发金冠，华贵。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 青布道袍**
- 曹国舅，凡尘砺心阶段·青布道袍。形象：白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：锦绣公子，锦袍。梳造：束发玉冠，温润。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 玉磬仙人**
- 曹国舅，道法初成阶段·玉磬仙人。形象：白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，玉磬一击，清音绕梁，天地皆静。衣着：青布道袍，弃官入道。梳造：道髻，眉目清雅。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 礼乐大道**
- 曹国舅，大劫淬炼阶段·礼乐大道。形象：白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，玉磬一击，清音绕梁，天地皆静。衣着：礼乐大道，玉磬在手。梳造：道冠，音律为魂。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 九天神音**
- 曹国舅，封神登天阶段·九天神音。形象：白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，玉磬一击，清音绕梁，天地皆静。衣着：九天神音，鹤氅白袍。梳造：莲花冠，玉磬清响。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 国舅真仙**
- 曹国舅，道果圆满阶段·国舅真仙。形象：白衣鹤氅的雅士，头戴莲花冠，手持莹白玉磬。 核心意象：玉磬、莲花冠、白鹤、朝服。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，玉磬一击，清音绕梁，天地皆静。衣着：国舅真仙，白鹤为伴。梳造：素发仙髻，超然物外。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 太上老君（`taishang_laojun`）

**灵胎初醒 · 上古真人**
- 太上老君，灵胎初醒阶段·上古真人。初始形态：一缕上古真人灵光，紫气自东而来，拂尘虚影与八卦炉的焰光交织。火属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：上古真人灵光，紫气东来。梳造：白发童子面，未成相。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 紫气东来**
- 太上老君，凡尘砺心阶段·紫气东来。形象：白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：青牛老者，紫金道袍初现。梳造：白发白须，面如童子。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 八卦炉**
- 太上老君，道法初成阶段·八卦炉。形象：白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，拂尘一挥，紫气东来三万里。衣着：八卦炉前，道袍烟云。梳造：发髻高束，拂尘在手。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 西游护法**
- 太上老君，大劫淬炼阶段·西游护法。形象：白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，拂尘一挥，紫气东来三万里。衣着：老君炉炼，紫气万丈。梳造：白发怒张，炉火映面。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 道德天尊**
- 太上老君，封神登天阶段·道德天尊。形象：白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，拂尘一挥，紫气东来三万里。衣着：道德天尊，紫金仙袍。梳造：仙髻金冠，紫气绕身。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 三清合一**
- 太上老君，道果圆满阶段·三清合一。形象：白发白须面如童子，紫金道袍，手持拂尘，身绕紫气。 核心意象：拂尘、八卦炉、金丹、青牛、紫气。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，拂尘一挥，紫气东来三万里。衣着：道祖至尊，三清合一紫气东来。梳造：紫金道冠，万道之源。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 钟馗（`zhong_kui`）

**灵胎初醒 · 终南寒窗**
- 钟馗，灵胎初醒阶段·终南寒窗。初始形态：一缕书生灵光伏案夜读，寒窗孤灯下字迹成链，隐约凝成赤面虬髯的轮廓。暗属性灵光微微环绕。神态：灵光中沉睡，仙缘初定的宁静。动作：灵光化形，气息未定。衣着：寒窗书生灵光，布巾束发。梳造：布巾束发，书生清癯。意境：仙山福地灵气孕育，灵光化形，渺小却蕴含仙缘。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：玄青（#3A4A5A）主调 + 仙金（#C9B37E）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; spiritual genesis; 灵光中沉睡，仙缘初定的宁静; palette #3A4A5A with #C9B37E accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 触阶而亡**
- 钟馗，凡尘砺心阶段·触阶而亡。形象：赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。神态：入世初见的纯澈。动作：初踏红尘，好奇四顾。衣着：青衫落拓，腰悬旧笔。梳造：发丝散乱，眉宇郁结。意境：离开洞府行走人间，红尘历练，磨砺道心。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：黛青（#5B6B7A）主调 + 灰白（#8A9BA8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; mortal wandering; 入世初见的纯澈; palette #5B6B7A with #8A9BA8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 地府扬名**
- 钟馗，道法初成阶段·地府扬名。形象：赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。神态：悟道中的专注，眸光渐亮。动作：功法初显，招式渐成，拔剑前先合书，判官笔一落定善恶。衣着：大红官袍初着，乌纱帽。梳造：乌纱帽微歪，髭须初张。意境：道法初成，法宝显威，法相庄严初现。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：朱砂（#B8502E）主调 + 明黄（#E8B84B）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; dao awakening; 悟道中的专注，眸光渐亮; palette #B8502E with #E8B84B accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 驱鬼之王**
- 钟馗，大劫淬炼阶段·驱鬼之王。形象：赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。神态：渡劫时的坚毅，眼神无畏。动作：全力抗劫，天劫加身，拔剑前先合书，判官笔一落定善恶。衣着：驱鬼之王，红袍与剑。梳造：虬髯怒张，豹眼圆睁。意境：天劫淬炼、心魔试炼，仙体历经大劫而不灭。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：墨褐（#241A14）主调 + 血绛（#8B1A1A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; heavenly tribulation; 渡劫时的坚毅，眼神无畏; palette #241A14 with #8B1A1A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 五方鬼卒**
- 钟馗，封神登天阶段·五方鬼卒。形象：赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。神态：仙光内蕴，目光深邃。动作：功法大成，威临天下，拔剑前先合书，判官笔一落定善恶。衣着：玄黑蟒袍金纹，紫金冠。梳造：紫金冠，虬髯如戟。意境：功行圆满封神登天，位列仙班，金光紫气环绕。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：紫金（#5B2A8A）主调 + 御金（#F5D742）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; ascension to immortals; 仙光内蕴，目光深邃; palette #5B2A8A with #F5D742 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 驱魔真君**
- 钟馗，道果圆满阶段·驱魔真君。形象：赤面虬髯，豹头环眼，额间竖纹，判官笔与青锋剑。 核心意象：判官笔、青锋剑、红袍、葫芦、终南山。神态：法相圆满，神光内蕴的威仪。动作：功行圆满，法相全开，拔剑前先合书，判官笔一落定善恶。衣着：镇宅捉鬼，判官笔指鬼青锋剑镇邪。梳造：虬髯怒张，紫金冠生威。意境：道果圆满，紫气东来，万法皆通证道果。风格：新国风水墨仙侠插画，工笔重彩、金箔点缀，仙气缥缈、道韵悠长。色彩：素白（#F5F0E8）主调 + 淡紫（#9A8FB8）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：xianxia Chinese ink-wash painting, gongbi detail, gold leaf accents, ethereal immortal aura; union with dao; 法相圆满，神光内蕴的威仪; palette #F5F0E8 with #9A8FB8 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 3. 宝可梦（12 物种）

> **风格**：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。**阶段演绎**：
> - 灵胎初醒：蛋壳微光，新生命初醒，温暖亲近的孵化氛围（奶油/嫩粉）
> - 凡尘砺心：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮（晴蓝/暖橙）
> - 道法初成：初次对战或训练，火花四溅，能力初显的高光时刻（活力橙/电光蓝）
> - 大劫淬炼：进化之光笼罩，身体在光中蜕变，挣扎又坚定（深海蓝/战斗红）
> - 封神登天：完全进化形态，可靠的战斗伙伴，自信昂扬（伙伴金/烈焰橙）
> - 道果圆满：超进化或极致形态，圣光守护，散发传说级气场（圣光白/彩虹金）

#### 小火龙（`charmander`）

**灵胎初醒 · 火纹蛋**
- 小火龙，灵胎初醒阶段·火纹蛋。初始形态：一枚火纹蛋，橙色蛋壳燃着一点尾焰，暖光从壳缝漏出，火之生命初醒。火属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 蹒跚火苗**
- 小火龙，凡尘砺心阶段·蹒跚火苗。形象：橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 烈火小龙**
- 小火龙，道法初成阶段·烈火小龙。形象：橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，鼓起腮帮，一记"火花"喷出炽热火焰。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 烈火恐龙**
- 小火龙，大劫淬炼阶段·烈火恐龙。形象：橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，鼓起腮帮，一记"火花"喷出炽热火焰。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 喷火龙**
- 小火龙，封神登天阶段·喷火龙。形象：橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，鼓起腮帮，一记"火花"喷出炽热火焰。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 超级喷火龙X**
- 小火龙，道果圆满阶段·超级喷火龙X。形象：橙色小蜥蜴，尾尖燃着火焰。 核心意象：尾尖火焰、橙红鳞片、进化之焰。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，鼓起腮帮，一记"火花"喷出炽热火焰。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 妙蛙种子（`bulbasaur`）

**灵胎初醒 · 种子蛋**
- 妙蛙种子，灵胎初醒阶段·种子蛋。初始形态：一粒种子蛋，青白种壳顶着一枚嫩芽，翠绿叶脉蔓延壳面，像会发芽的小世界。木属性灵光微微环绕。神态：沉睡的种子里，藏着破土的心跳。动作：静卧泥土，待春雨而萌。衣着：种子/种壳，纹理细腻。梳造：顶端嫩芽微顶。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 沉睡的种子里，藏着破土的心跳; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼种子**
- 妙蛙种子，凡尘砺心阶段·幼种子。形象：蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。神态：初生的好奇，向着光的方向。动作：嫩芽破土，努力伸展。衣着：幼芽新叶，娇嫩欲滴。梳造：顶芽鲜嫩，两片子叶。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 初生的好奇，向着光的方向; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 叶芽成长**
- 妙蛙种子，道法初成阶段·叶芽成长。形象：蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。神态：生机勃发，神采奕奕。动作：生机初现，枝叶舒展，藤鞭从背苞探出，一记"飞叶快刀"斩向敌手。衣着：茎叶渐盛，花苞初成。梳造：花苞/嫩叶环生。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 生机勃发，神采奕奕; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 花苞绽放**
- 妙蛙种子，大劫淬炼阶段·花苞绽放。形象：蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。神态：风雨中的坚韧，眸光不折。动作：全力生长，根系深扎，藤鞭从背苞探出，一记"飞叶快刀"斩向敌手。衣着：枝叶繁茂，带风霜痕迹。梳造：花叶翻卷，仍自挺立。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 风雨中的坚韧，眸光不折; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 繁花之姿**
- 妙蛙种子，封神登天阶段·繁花之姿。形象：蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。神态：华彩绽放，生机盎然的威仪。动作：生机大成，繁花满枝，藤鞭从背苞探出，一记"飞叶快刀"斩向敌手。衣着：花开满枝/树冠如云。梳造：花冠/树冠，光华流转。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 华彩绽放，生机盎然的威仪; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 花神**
- 妙蛙种子，道果圆满阶段·花神。形象：蓝绿蛙状小兽，背驮种子苞。 核心意象：背上的种子苞、藤鞭、暖阳。神态：生机圆满，枝叶参天的宁静力量。动作：根深叶茂，繁花结成果实，藤鞭从背苞探出，一记"飞叶快刀"斩向敌手。衣着：古木参天，树冠如盖。梳造：万叶花冠，生机无限。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 生机圆满，枝叶参天的宁静力量; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 杰尼龟（`squirtle`）

**灵胎初醒 · 水纹蛋**
- 杰尼龟，灵胎初醒阶段·水纹蛋。初始形态：一枚水纹蛋，蓝色蛋壳泛着水波鳞光，浪纹在壳面流动，是水之灵胎。水属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 幼龟**
- 杰尼龟，凡尘砺心阶段·幼龟。形象：蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 水枪龟**
- 杰尼龟，道法初成阶段·水枪龟。形象：蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，深吸一口气，一记"水枪"喷射而出。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 水炮龟**
- 杰尼龟，大劫淬炼阶段·水炮龟。形象：蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，深吸一口气，一记"水枪"喷射而出。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 水箭龟**
- 杰尼龟，封神登天阶段·水箭龟。形象：蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，深吸一口气，一记"水枪"喷射而出。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 海神龟**
- 杰尼龟，道果圆满阶段·海神龟。形象：蓝色小龟，背驮硬壳。 核心意象：坚固龟壳、水枪水泡、机灵眼神。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，深吸一口气，一记"水枪"喷射而出。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 伊布，大劫淬炼阶段·元素显现。形象：棕毛狐狸样小兽，颈毛蓬松。 核心意象：蓬松颈毛、进化之力、多形态可能。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，蹭蹭主人的手，一记"撞击"轻快而出。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 水伊布·潮汐**
- 伊布，封神登天阶段·水伊布·潮汐。形象：棕毛狐狸样小兽，颈毛蓬松。 核心意象：蓬松颈毛、进化之力、多形态可能。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，蹭蹭主人的手，一记"撞击"轻快而出。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 水伊布·海神**
- 伊布，道果圆满阶段·水伊布·海神。形象：棕毛狐狸样小兽，颈毛蓬松。 核心意象：蓬松颈毛、进化之力、多形态可能。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，蹭蹭主人的手，一记"撞击"轻快而出。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 皮卡丘（`pikachu`）

**灵胎初醒 · 电气蛋**
- 皮卡丘，灵胎初醒阶段·电气蛋。初始形态：一枚电气蛋，暖黄蛋壳泛着微光，细碎电流纹在壳面爬动，仿佛蛰伏的雷种。雷属性灵光微微环绕。神态：眼帘低垂，呼吸轻浅，沉在初生的微光里。动作：蜷卧于光中，静默沉睡。衣着：细嫩绒毛/初生软鳞，未成形。梳造：新生短毛，耳未立/角未出。意境：蛋壳微光，新生命初醒，温暖亲近的孵化氛围。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：奶油（#F5E6D3）主调 + 嫩粉（#FFB6C1）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; hatching glow; 眼帘低垂，呼吸轻浅，沉在初生的微光里; palette #F5E6D3 with #FFB6C1 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 皮丘萌动**
- 皮卡丘，凡尘砺心阶段·皮丘萌动。形象：黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。神态：眸光清澈，好奇打量，怯生生又藏不住欢喜。动作：蹒跚学步，试探着触碰四周。衣着：绒毛渐丰，隐约透出本色。梳造：耳尾渐立，茸毛蓬松。意境：蹒跚学步的幼体，在森林或水边初遇世界，好奇明亮。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：晴蓝（#87CEEB）主调 + 暖橙（#FFA07A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first steps; 眸光清澈，好奇打量，怯生生又藏不住欢喜; palette #87CEEB with #FFA07A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 蓄电皮卡丘**
- 皮卡丘，道法初成阶段·蓄电皮卡丘。形象：黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。神态：眸光晶亮，意气初显，跃跃欲试。动作：初展身手，初露锋芒，双颊电光噼啪，一记"十万伏特"放射而出。衣着：毛色分明，姿态挺拔。梳造：鬃毛/羽角初生，精神抖擞。意境：初次对战或训练，火花四溅，能力初显的高光时刻。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：活力橙（#FF8C42）主调 + 电光蓝（#4FC3F7）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; first battle spark; 眸光晶亮，意气初显，跃跃欲试; palette #FF8C42 with #4FC3F7 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 雷丘**
- 皮卡丘，大劫淬炼阶段·雷丘。形象：黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，双颊电光噼啪，一记"十万伏特"放射而出。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 飞空皮卡丘**
- 皮卡丘，封神登天阶段·飞空皮卡丘。形象：黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，双颊电光噼啪，一记"十万伏特"放射而出。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 极巨化皮卡丘**
- 皮卡丘，道果圆满阶段·极巨化皮卡丘。形象：黄色小鼠，双颊红晕。 核心意象：红色脸颊、闪电尾巴、十万伏特。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，双颊电光噼啪，一记"十万伏特"放射而出。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 利欧路，大劫淬炼阶段·路卡利欧。形象：蓝黑相间的小犬，胸前黑纹。 核心意象：胸前黑纹、波导之力、格斗之魂。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，双掌凝聚波导弹，一发轰出破空而去。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 波导大师**
- 利欧路，封神登天阶段·波导大师。形象：蓝黑相间的小犬，胸前黑纹。 核心意象：胸前黑纹、波导之力、格斗之魂。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，双掌凝聚波导弹，一发轰出破空而去。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 传说守护者**
- 利欧路，道果圆满阶段·传说守护者。形象：蓝黑相间的小犬，胸前黑纹。 核心意象：胸前黑纹、波导之力、格斗之魂。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，双掌凝聚波导弹，一发轰出破空而去。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 冰·霜翎狐，大劫淬炼阶段·凛冬狐。形象：冰蓝狐狸，尾带霜华。 核心意象：冰蓝毛皮、霜华、极寒雪原。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，雪原上一纵而过，尾尖带起一路霜尘。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 永冬狐**
- 冰·霜翎狐，封神登天阶段·永冬狐。形象：冰蓝狐狸，尾带霜华。 核心意象：冰蓝毛皮、霜华、极寒雪原。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，雪原上一纵而过，尾尖带起一路霜尘。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 极地尊者**
- 冰·霜翎狐，道果圆满阶段·极地尊者。形象：冰蓝狐狸，尾带霜华。 核心意象：冰蓝毛皮、霜华、极寒雪原。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，雪原上一纵而过，尾尖带起一路霜尘。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 岩·岩甲犀，大劫淬炼阶段·铁岩犀。形象：岩石甲身，额生尖角。 核心意象：岩石甲身、额尖角、磐石之重。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，低头一顶，尖角撞碎前方巨岩。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大地犀**
- 岩·岩甲犀，封神登天阶段·大地犀。形象：岩石甲身，额生尖角。 核心意象：岩石甲身、额尖角、磐石之重。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，低头一顶，尖角撞碎前方巨岩。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 岩之尊者**
- 岩·岩甲犀，道果圆满阶段·岩之尊者。形象：岩石甲身，额生尖角。 核心意象：岩石甲身、额尖角、磐石之重。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，低头一顶，尖角撞碎前方巨岩。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 风·掠风隼，大劫淬炼阶段·破空隼。形象：翠绿猛隼，翼疾如风。 核心意象：翠绿翼羽、疾风、万里长空。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，高空敛翅，一瞬俯冲如风掠地。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 王者风隼**
- 风·掠风隼，封神登天阶段·王者风隼。形象：翠绿猛隼，翼疾如风。 核心意象：翠绿翼羽、疾风、万里长空。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，高空敛翅，一瞬俯冲如风掠地。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 天空尊者**
- 风·掠风隼，道果圆满阶段·天空尊者。形象：翠绿猛隼，翼疾如风。 核心意象：翠绿翼羽、疾风、万里长空。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，高空敛翅，一瞬俯冲如风掠地。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 光·晨曦鹿，大劫淬炼阶段·破晓鹿。形象：金光神鹿，鹿角如炬。 核心意象：金鹿角、光束、驱暗之光。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，鹿角燃起金光，踏光而行驱散阴霾。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 辉明鹿**
- 光·晨曦鹿，封神登天阶段·辉明鹿。形象：金光神鹿，鹿角如炬。 核心意象：金鹿角、光束、驱暗之光。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，鹿角燃起金光，踏光而行驱散阴霾。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 光辉尊者**
- 光·晨曦鹿，道果圆满阶段·光辉尊者。形象：金光神鹿，鹿角如炬。 核心意象：金鹿角、光束、驱暗之光。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，鹿角燃起金光，踏光而行驱散阴霾。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 暗·夜隐豹，大劫淬炼阶段·幽冥豹。形象：墨黑豹影，身姿流线。 核心意象：墨黑皮毛、流线身姿、暗夜利爪。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，身影没入夜色，骤然扑出快如闪电。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 暗夜猎豹**
- 暗·夜隐豹，封神登天阶段·暗夜猎豹。形象：墨黑豹影，身姿流线。 核心意象：墨黑皮毛、流线身姿、暗夜利爪。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，身影没入夜色，骤然扑出快如闪电。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 暗夜尊者**
- 暗·夜隐豹，道果圆满阶段·暗夜尊者。形象：墨黑豹影，身姿流线。 核心意象：墨黑皮毛、流线身姿、暗夜利爪。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，身影没入夜色，骤然扑出快如闪电。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 钢·磁甲兽，大劫淬炼阶段·磁暴兽。形象：钢甲犰狳，甲带如鳞。 核心意象：钢甲鳞带、铁球之姿、坚不可摧。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，蜷身一滚，化作铁球轰然撞出。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：进化之光笼罩，身体在光中蜕变，挣扎又坚定。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：深海蓝（#283593）主调 + 战斗红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; evolution light; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #283593 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 合金兽**
- 钢·磁甲兽，封神登天阶段·合金兽。形象：钢甲犰狳，甲带如鳞。 核心意象：钢甲鳞带、铁球之姿、坚不可摧。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，蜷身一滚，化作铁球轰然撞出。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：完全进化形态，可靠的战斗伙伴，自信昂扬。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：伙伴金（#FFD54F）主调 + 烈焰橙（#FF7043）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; fully evolved; 目光如炬，不怒自威，威仪自生; palette #FFD54F with #FF7043 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 钢之尊者**
- 钢·磁甲兽，道果圆满阶段·钢之尊者。形象：钢甲犰狳，甲带如鳞。 核心意象：钢甲鳞带、铁球之姿、坚不可摧。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，蜷身一滚，化作铁球轰然撞出。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：超进化或极致形态，圣光守护，散发传说级气场。风格：经典日系动画冒险风，赛璐璐质感、高饱和明快，热血治愈。色彩：圣光白（#FFFFFF）主调 + 彩虹金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：classic Japanese anime adventure style, cel-shaded, vibrant and wholesome; legendary aura; 圆满自足，神光内蕴的从容; palette #FFFFFF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 亚古兽，大劫淬炼阶段·暴龙兽·强攻。形象：橙色小恐龙，背甲坚挺。 核心意象：勇气徽章、橙色背甲、数码蛋。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙，口中聚焰，一记"小型火焰"喷向敌人。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 机械暴龙兽·强化**
- 亚古兽，封神登天阶段·机械暴龙兽·强化。形象：橙色小恐龙，背甲坚挺。 核心意象：勇气徽章、橙色背甲、数码蛋。神态：究极体的神圣威仪。动作：绝技大成，数据凝成圣甲，口中聚焰，一记"小型火焰"喷向敌人。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体的神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 战斗暴龙兽·勇气**
- 亚古兽，道果圆满阶段·战斗暴龙兽·勇气。形象：橙色小恐龙，背甲坚挺。 核心意象：勇气徽章、橙色背甲、数码蛋。神态：究极圆满，守护之光如日。动作：绝技巅峰，数码洪流归于一念，口中聚焰，一记"小型火焰"喷向敌人。衣着：究极神圣甲，数据凝光。梳造：光之圣盔，金色辉光。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 究极圆满，守护之光如日; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 迪路兽，大劫淬炼阶段·猫人兽·雷霆。形象：白色圣猫，耳戴神圣之环。 核心意象：神圣之环、光明徽章、圣洁之光。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙，猫拳连击如电，一记"猫猫拳"撕裂黑暗。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 神圣天女兽·圣光**
- 迪路兽，封神登天阶段·神圣天女兽·圣光。形象：白色圣猫，耳戴神圣之环。 核心意象：神圣之环、光明徽章、圣洁之光。神态：究极体的神圣威仪。动作：绝技大成，数据凝成圣甲，猫拳连击如电，一记"猫猫拳"撕裂黑暗。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体的神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 圣龙兽·光明**
- 迪路兽，道果圆满阶段·圣龙兽·光明。形象：白色圣猫，耳戴神圣之环。 核心意象：神圣之环、光明徽章、圣洁之光。神态：究极圆满，守护之光如日。动作：绝技巅峰，数码洪流归于一念，猫拳连击如电，一记"猫猫拳"撕裂黑暗。衣着：究极神圣甲，数据凝光。梳造：光之圣盔，金色辉光。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 究极圆满，守护之光如日; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 巴达兽，大劫淬炼阶段·天使兽·圣杖。形象：奶油色圆身，扑腾小翅膀。 核心意象：希望徽章、奶油圆身、神圣羽翼。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙，扑扇小翅盘旋，一记"空气炮"轰向敌人。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 神圣天使兽·审判**
- 巴达兽，封神登天阶段·神圣天使兽·审判。形象：奶油色圆身，扑腾小翅膀。 核心意象：希望徽章、奶油圆身、神圣羽翼。神态：究极体的神圣威仪。动作：绝技大成，数据凝成圣甲，扑扇小翅盘旋，一记"空气炮"轰向敌人。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体的神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 炽天使兽·希望**
- 巴达兽，道果圆满阶段·炽天使兽·希望。形象：奶油色圆身，扑腾小翅膀。 核心意象：希望徽章、奶油圆身、神圣羽翼。神态：究极圆满，守护之光如日。动作：绝技巅峰，数码洪流归于一念，扑扇小翅盘旋，一记"空气炮"轰向敌人。衣着：究极神圣甲，数据凝光。梳造：光之圣盔，金色辉光。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 究极圆满，守护之光如日; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 加布兽，大劫淬炼阶段·加鲁鲁兽·狂岚。形象：身披蓝毛的小狼，额生独角。 核心意象：友情徽章、独角、蓝色毛皮。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙，额角凝聚蓝光，一记"小狼爪"快如闪电。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 兽人加鲁鲁·月狼**
- 加布兽，封神登天阶段·兽人加鲁鲁·月狼。形象：身披蓝毛的小狼，额生独角。 核心意象：友情徽章、独角、蓝色毛皮。神态：究极体的神圣威仪。动作：绝技大成，数据凝成圣甲，额角凝聚蓝光，一记"小狼爪"快如闪电。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体的神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 钢铁加鲁鲁兽·友情**
- 加布兽，道果圆满阶段·钢铁加鲁鲁兽·友情。形象：身披蓝毛的小狼，额生独角。 核心意象：友情徽章、独角、蓝色毛皮。神态：究极圆满，守护之光如日。动作：绝技巅峰，数码洪流归于一念，额角凝聚蓝光，一记"小狼爪"快如闪电。衣着：究极神圣甲，数据凝光。梳造：光之圣盔，金色辉光。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 究极圆满，守护之光如日; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 比丘兽，大劫淬炼阶段·巴多拉兽·烈焰。形象：粉色小鸟，头顶凤冠。 核心意象：纯真徽章、凤冠、藤蔓之花。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙，凤冠一闪，一记"毒藤蔓"甩出缠绕敌手。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 伽楼达兽·疾风**
- 比丘兽，封神登天阶段·伽楼达兽·疾风。形象：粉色小鸟，头顶凤冠。 核心意象：纯真徽章、凤冠、藤蔓之花。神态：究极体的神圣威仪。动作：绝技大成，数据凝成圣甲，凤冠一闪，一记"毒藤蔓"甩出缠绕敌手。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体的神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 凤凰兽·涅槃**
- 比丘兽，道果圆满阶段·凤凰兽·涅槃。形象：粉色小鸟，头顶凤冠。 核心意象：纯真徽章、凤冠、藤蔓之花。神态：究极圆满，守护之光如日。动作：绝技巅峰，数码洪流归于一念，凤冠一闪，一记"毒藤蔓"甩出缠绕敌手。衣着：究极神圣甲，数据凝光。梳造：光之圣盔，金色辉光。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 究极圆满，守护之光如日; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 哥玛兽，大劫淬炼阶段·海狮兽·破浪。形象：白色海狮，圆滚滚胖乎乎。 核心意象：诚实徽章、圆滚滚、鱼群伙伴。神态：直面黑暗的勇毅。动作：全力施为，与深渊对峙，尾巴一甩，召唤"鱼群大进军"扑向敌阵。衣着：完全体铠甲，伤痕累累。梳造：战甲破损，眸光如焰。意境：面对黑暗深渊的试炼，进化之光与黑暗对峙。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：深渊紫（#2C2C54）主调 + 战斗红（#FF3B30）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; dark trial; 直面黑暗的勇毅; palette #2C2C54 with #FF3B30 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 祖顿兽·重锤**
- 哥玛兽，封神登天阶段·祖顿兽·重锤。形象：白色海狮，圆滚滚胖乎乎。 核心意象：诚实徽章、圆滚滚、鱼群伙伴。神态：究极体的神圣威仪。动作：绝技大成，数据凝成圣甲，尾巴一甩，召唤"鱼群大进军"扑向敌阵。衣着：究极体圣甲，辉光万丈。梳造：圣盔金角，威仪堂堂。意境：究极进化，圣光加冕，数据粒子凝成铠甲。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣金（#FFD700）主调 + 纯白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; ultimate evolution; 究极体的神圣威仪; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 维京兽·咆哮**
- 哥玛兽，道果圆满阶段·维京兽·咆哮。形象：白色海狮，圆滚滚胖乎乎。 核心意象：诚实徽章、圆滚滚、鱼群伙伴。神态：究极圆满，守护之光如日。动作：绝技巅峰，数码洪流归于一念，尾巴一甩，召唤"鱼群大进军"扑向敌阵。衣着：究极神圣甲，数据凝光。梳造：光之圣盔，金色辉光。意境：神圣形态，希望之光普照，数据洪流归于一念。风格：日系数码进化动画风，赛璐璐质感配数码科技光效，热血进化、高光时刻。色彩：圣光白（#FFFFFF）主调 + 虹彩（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Japanese digital evolution anime, cel shading with digital particles, techno glow; holy transcendence; 究极圆满，守护之光如日; palette #FFFFFF with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 大熊猫，大劫淬炼阶段·青涩少年。形象：黑白圆滚，黑耳黑眼圈。 核心意象：黑白皮毛、竹笋、团团圆圆。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，抱起竹笋慢悠悠地啃，一个翻身滚下坡。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 成年熊猫**
- 大熊猫，封神登天阶段·成年熊猫。形象：黑白圆滚，黑耳黑眼圈。 核心意象：黑白皮毛、竹笋、团团圆圆。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，抱起竹笋慢悠悠地啃，一个翻身滚下坡。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 太极熊猫**
- 大熊猫，道果圆满阶段·太极熊猫。形象：黑白圆滚，黑耳黑眼圈。 核心意象：黑白皮毛、竹笋、团团圆圆。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，抱起竹笋慢悠悠地啃，一个翻身滚下坡。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 金丝猴，大劫淬炼阶段·长尾猴。形象：金色长毛，蓝脸朝天鼻。 核心意象：金色长毛、蓝脸、高山密林。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，抓住藤蔓荡起，一纵数米穿行林间。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金猴王**
- 金丝猴，封神登天阶段·金猴王。形象：金色长毛，蓝脸朝天鼻。 核心意象：金色长毛、蓝脸、高山密林。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，抓住藤蔓荡起，一纵数米穿行林间。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 金丝大圣**
- 金丝猴，道果圆满阶段·金丝大圣。形象：金色长毛，蓝脸朝天鼻。 核心意象：金色长毛、蓝脸、高山密林。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，抓住藤蔓荡起，一纵数米穿行林间。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 丹顶鹤，大劫淬炼阶段·长寿鹤。形象：白羽黑颈，头顶一点朱红。 核心意象：头顶朱红、白羽黑颈、松鹤延年。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，引颈长鸣，展翅起舞，声闻九天。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 祥云伴鹤**
- 丹顶鹤，封神登天阶段·祥云伴鹤。形象：白羽黑颈，头顶一点朱红。 核心意象：头顶朱红、白羽黑颈、松鹤延年。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，引颈长鸣，展翅起舞，声闻九天。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 羽化鹤**
- 丹顶鹤，道果圆满阶段·羽化鹤。形象：白羽黑颈，头顶一点朱红。 核心意象：头顶朱红、白羽黑颈、松鹤延年。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，引颈长鸣，展翅起舞，声闻九天。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 华南虎，大劫淬炼阶段·百兽之王。形象：橙底黑纹，身形矫健。 核心意象：橙底黑纹、额间王字、密林之王。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，夜幕中悄然潜行，一声虎啸震慑山林。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金虎**
- 华南虎，封神登天阶段·金虎。形象：橙底黑纹，身形矫健。 核心意象：橙底黑纹、额间王字、密林之王。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，夜幕中悄然潜行，一声虎啸震慑山林。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 神虎**
- 华南虎，道果圆满阶段·神虎。形象：橙底黑纹，身形矫健。 核心意象：橙底黑纹、额间王字、密林之王。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，夜幕中悄然潜行，一声虎啸震慑山林。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 扬子鳄，大劫淬炼阶段·铁甲鳄。形象：低伏的宽吻鳄龙，身披硬鳞。 核心意象：硬鳞、宽吻、江河泥潭。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身，水中悄然潜行，猛然一扑咬定猎物。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 龙威炽烈，怒目电光; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大鳄**
- 扬子鳄，封神登天阶段·大鳄。形象：低伏的宽吻鳄龙，身披硬鳞。 核心意象：硬鳞、宽吻、江河泥潭。神态：龙目洞彻，神威赫赫。动作：绝技大成，行云布雨，水中悄然潜行，猛然一扑咬定猎物。衣着：金鳞覆身，祥光万道。梳造：龙角如珊瑚，须垂百丈。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 龙目洞彻，神威赫赫; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 蛟鳄**
- 扬子鳄，道果圆满阶段·蛟鳄。形象：低伏的宽吻鳄龙，身披硬鳞。 核心意象：硬鳞、宽吻、江河泥潭。神态：真身圆满，龙威盖世。动作：腾云驾雾，号令风雨，水中悄然潜行，猛然一扑咬定猎物。衣着：金鳞神光，日月同辉。梳造：龙角如岳，须垂星河。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 真身圆满，龙威盖世; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 朱鹮，大劫淬炼阶段·祥瑞之鸟。形象：白羽粉翼，头冠蓬松。 核心意象：粉红羽翼、长喙、水田浅滩。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，长喙探入浅水，轻巧地叼起鱼虾。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 吉祥朱鹮**
- 朱鹮，封神登天阶段·吉祥朱鹮。形象：白羽粉翼，头冠蓬松。 核心意象：粉红羽翼、长喙、水田浅滩。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，长喙探入浅水，轻巧地叼起鱼虾。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 凤凰近亲**
- 朱鹮，道果圆满阶段·凤凰近亲。形象：白羽粉翼，头冠蓬松。 核心意象：粉红羽翼、长喙、水田浅滩。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，长喙探入浅水，轻巧地叼起鱼虾。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 藏羚羊，大劫淬炼阶段·奔风羚。形象：白褐相间，长角如剑。 核心意象：剑形长角、白褐毛皮、可可西里。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，四蹄翻飞，在高原上一路狂奔如飞。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 雪原之王**
- 藏羚羊，封神登天阶段·雪原之王。形象：白褐相间，长角如剑。 核心意象：剑形长角、白褐毛皮、可可西里。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，四蹄翻飞，在高原上一路狂奔如飞。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 雪原尊者**
- 藏羚羊，道果圆满阶段·雪原尊者。形象：白褐相间，长角如剑。 核心意象：剑形长角、白褐毛皮、可可西里。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，四蹄翻飞，在高原上一路狂奔如飞。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 雪豹，大劫淬炼阶段·雪原豹。形象：灰白底黑斑，长尾粗壮。 核心意象：灰白斑纹、粗壮长尾、雪山之巅。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，长尾平衡，纵身跃过悬崖峭壁。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 高山雪豹**
- 雪豹，封神登天阶段·高山雪豹。形象：灰白底黑斑，长尾粗壮。 核心意象：灰白斑纹、粗壮长尾、雪山之巅。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，长尾平衡，纵身跃过悬崖峭壁。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 雪山尊者**
- 雪豹，道果圆满阶段·雪山尊者。形象：灰白底黑斑，长尾粗壮。 核心意象：灰白斑纹、粗壮长尾、雪山之巅。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，长尾平衡，纵身跃过悬崖峭壁。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 麋鹿，大劫淬炼阶段·泽畔鹿。形象：角似鹿非鹿，蹄似牛非牛。 核心意象：四不像之形、水泽湿地、角蹄。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，踏水而行，慢悠悠地在湿地踱步。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 瑞泽鹿**
- 麋鹿，封神登天阶段·瑞泽鹿。形象：角似鹿非鹿，蹄似牛非牛。 核心意象：四不像之形、水泽湿地、角蹄。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，踏水而行，慢悠悠地在湿地踱步。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 泽畔尊者**
- 麋鹿，道果圆满阶段·泽畔尊者。形象：角似鹿非鹿，蹄似牛非牛。 核心意象：四不像之形、水泽湿地、角蹄。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，踏水而行，慢悠悠地在湿地踱步。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 东北虎，大劫淬炼阶段·啸雪虎。形象：硕大虎躯，厚毛如披风。 核心意象：厚毛披风、硕大虎躯、林海雪原。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，一声长啸震彻林海，威压百兽。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 雪原巨虎**
- 东北虎，封神登天阶段·雪原巨虎。形象：硕大虎躯，厚毛如披风。 核心意象：厚毛披风、硕大虎躯、林海雪原。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，一声长啸震彻林海，威压百兽。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 雪原尊者**
- 东北虎，道果圆满阶段·雪原尊者。形象：硕大虎躯，厚毛如披风。 核心意象：厚毛披风、硕大虎躯、林海雪原。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，一声长啸震彻林海，威压百兽。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 小熊猫，大劫淬炼阶段·红锦猫熊。形象：红褐毛皮，环纹大尾。 核心意象：红褐毛皮、环纹大尾、竹叶。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，竖起环纹大尾，抱着一捧竹叶细嚼。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 祥瑞小熊**
- 小熊猫，封神登天阶段·祥瑞小熊。形象：红褐毛皮，环纹大尾。 核心意象：红褐毛皮、环纹大尾、竹叶。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，竖起环纹大尾，抱着一捧竹叶细嚼。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 竹林尊者**
- 小熊猫，道果圆满阶段·竹林尊者。形象：红褐毛皮，环纹大尾。 核心意象：红褐毛皮、环纹大尾、竹叶。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，竖起环纹大尾，抱着一捧竹叶细嚼。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 江豚，大劫淬炼阶段·微笑天使。形象：圆头无背鳍，灰蓝光滑。 核心意象：无背鳍圆头、江涛、微笑。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，跃出江面，嘴角的弧度像一抹微笑。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：成年期沉稳担当，守护领地与族群，历经风雨。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：深林绿（#3E5C4B）主调 + 熟褐（#C45A3C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian prime; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E5C4B with #C45A3C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 江中守护**
- 江豚，封神登天阶段·江中守护。形象：圆头无背鳍，灰蓝光滑。 核心意象：无背鳍圆头、江涛、微笑。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，跃出江面，嘴角的弧度像一抹微笑。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为族群的传奇，祥瑞光环，被守护与敬仰。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：祥瑞金（#D4AF37）主调 + 王绿（#2F4F4F）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; legend of the species; 目光如炬，不怒自威，威仪自生; palette #D4AF37 with #2F4F4F accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 江流尊者**
- 江豚，道果圆满阶段·江流尊者。形象：圆头无背鳍，灰蓝光滑。 核心意象：无背鳍圆头、江涛、微笑。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，跃出江面，嘴角的弧度像一抹微笑。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：生态图腾显圣，祥瑞护世，生生不息。风格：皮克斯式写实萌化卡通，柔软毛发质感，自然纪录片背景，治愈萌暖。色彩：云白（#F5F0E8）主调 + 生态金绿（#A3C1AD）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Pixar-style realistic cutesy cartoon, soft fur, nature documentary backdrop; guardian spirit; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #A3C1AD accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 独角兽，大劫淬炼阶段·彩虹鬃毛。形象：纯白骏马，额生螺旋独角。 核心意象：螺旋独角、纯白毛皮、圣光。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，踏云而来，独角泛起一圈圣洁光环。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 眸光如电，威严中带着坚韧; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 星夜独角兽**
- 独角兽，封神登天阶段·星夜独角兽。形象：纯白骏马，额生螺旋独角。 核心意象：螺旋独角、纯白毛皮、圣光。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，踏云而来，独角泛起一圈圣洁光环。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 神兽威严，目光洞彻九幽; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 永恒独角兽**
- 独角兽，道果圆满阶段·永恒独角兽。形象：纯白骏马，额生螺旋独角。 核心意象：螺旋独角、纯白毛皮、圣光。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，踏云而来，独角泛起一圈圣洁光环。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 瑞气圆满，祥光普照; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 飞龙，大劫淬炼阶段·雷龙。形象：赤色飞龙，蝠翼如膜，尾带尖刺。 核心意象：蝠翼、尖刺长尾、火山烈焰。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身，双翼一振腾空，喷吐烈焰焚尽大地。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 龙威炽烈，怒目电光; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 古龙**
- 飞龙，封神登天阶段·古龙。形象：赤色飞龙，蝠翼如膜，尾带尖刺。 核心意象：蝠翼、尖刺长尾、火山烈焰。神态：龙目洞彻，神威赫赫。动作：绝技大成，行云布雨，双翼一振腾空，喷吐烈焰焚尽大地。衣着：金鳞覆身，祥光万道。梳造：龙角如珊瑚，须垂百丈。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 龙目洞彻，神威赫赫; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 太古龙**
- 飞龙，道果圆满阶段·太古龙。形象：赤色飞龙，蝠翼如膜，尾带尖刺。 核心意象：蝠翼、尖刺长尾、火山烈焰。神态：真身圆满，龙威盖世。动作：腾云驾雾，号令风雨，双翼一振腾空，喷吐烈焰焚尽大地。衣着：金鳞神光，日月同辉。梳造：龙角如岳，须垂星河。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 真身圆满，龙威盖世; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 精灵，大劫淬炼阶段·火精灵。形象：金翼小仙子，身绕光点。 核心意象：金翼、光点、花间飞舞。神态：风雨中的坚韧，眸光不折。动作：全力生长，根系深扎，金翼轻振，洒下点点星尘与花粉。衣着：枝叶繁茂，带风霜痕迹。梳造：花叶翻卷，仍自挺立。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 风雨中的坚韧，眸光不折; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大地精灵**
- 精灵，封神登天阶段·大地精灵。形象：金翼小仙子，身绕光点。 核心意象：金翼、光点、花间飞舞。神态：华彩绽放，生机盎然的威仪。动作：生机大成，繁花满枝，金翼轻振，洒下点点星尘与花粉。衣着：花开满枝/树冠如云。梳造：花冠/树冠，光华流转。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 华彩绽放，生机盎然的威仪; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 森林之子**
- 精灵，道果圆满阶段·森林之子。形象：金翼小仙子，身绕光点。 核心意象：金翼、光点、花间飞舞。神态：生机圆满，枝叶参天的宁静力量。动作：根深叶茂，繁花结成果实，金翼轻振，洒下点点星尘与花粉。衣着：古木参天，树冠如盖。梳造：万叶花冠，生机无限。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 生机圆满，枝叶参天的宁静力量; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 树人，大劫淬炼阶段·古树人。形象：深绿树人，枝干为臂。 核心意象：深绿树皮、年轮、森林之魂。神态：风雨中的坚韧，眸光不折。动作：全力生长，根系深扎，大地震颤着迈步，枝条如臂挥展。衣着：枝叶繁茂，带风霜痕迹。梳造：花叶翻卷，仍自挺立。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 风雨中的坚韧，眸光不折; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 森林卫士**
- 树人，封神登天阶段·森林卫士。形象：深绿树人，枝干为臂。 核心意象：深绿树皮、年轮、森林之魂。神态：华彩绽放，生机盎然的威仪。动作：生机大成，繁花满枝，大地震颤着迈步，枝条如臂挥展。衣着：花开满枝/树冠如云。梳造：花冠/树冠，光华流转。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 华彩绽放，生机盎然的威仪; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 世界树幼苗**
- 树人，道果圆满阶段·世界树幼苗。形象：深绿树人，枝干为臂。 核心意象：深绿树皮、年轮、森林之魂。神态：生机圆满，枝叶参天的宁静力量。动作：根深叶茂，繁花结成果实，大地震颤着迈步，枝条如臂挥展。衣着：古木参天，树冠如盖。梳造：万叶花冠，生机无限。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 生机圆满，枝叶参天的宁静力量; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 狮鹫，大劫淬炼阶段·皇家狮鹫。形象：鹰首狮身，金色双翼。 核心意象：鹰首狮身、金色双翼、忠勇。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，双翼一振冲天，俯冲时迅如流星。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 眸光如电，威严中带着坚韧; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 雷光狮鹫**
- 狮鹫，封神登天阶段·雷光狮鹫。形象：鹰首狮身，金色双翼。 核心意象：鹰首狮身、金色双翼、忠勇。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，双翼一振冲天，俯冲时迅如流星。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 神兽威严，目光洞彻九幽; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 神圣狮鹫**
- 狮鹫，道果圆满阶段·神圣狮鹫。形象：鹰首狮身，金色双翼。 核心意象：鹰首狮身、金色双翼、忠勇。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，双翼一振冲天，俯冲时迅如流星。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 瑞气圆满，祥光普照; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 美人鱼，大劫淬炼阶段·珊瑚公主。形象：碧鳞人鱼，长发如波。 核心意象：碧鳞鱼尾、如波长发、深海歌声。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，鱼尾轻摆破浪，一曲歌声飘向海面。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 眸光如电，威严中带着坚韧; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 人鱼女王**
- 美人鱼，封神登天阶段·人鱼女王。形象：碧鳞人鱼，长发如波。 核心意象：碧鳞鱼尾、如波长发、深海歌声。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，鱼尾轻摆破浪，一曲歌声飘向海面。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 神兽威严，目光洞彻九幽; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 大洋之魂**
- 美人鱼，道果圆满阶段·大洋之魂。形象：碧鳞人鱼，长发如波。 核心意象：碧鳞鱼尾、如波长发、深海歌声。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，鱼尾轻摆破浪，一曲歌声飘向海面。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 瑞气圆满，祥光普照; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 灰袍智者，大劫淬炼阶段·魔典持有者。形象：灰袍巫师，手持法杖。 核心意象：灰袍、法杖、奥术符文。神态：暗夜中坚持，眼神倔强。动作：全力施法，魔力风暴，法杖顿地，口中吟唱，星光自天穹应召而下。衣着：法袍破损，魔力燃烧。梳造：发丝凌乱，眸中带焰。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 暗夜中坚持，眼神倔强; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 白袍法师**
- 灰袍智者，封神登天阶段·白袍法师。形象：灰袍巫师，手持法杖。 核心意象：灰袍、法杖、奥术符文。神态：威能内蕴，光晕庄严。动作：魔法大成，法阵漫天，法杖顿地，口中吟唱，星光自天穹应召而下。衣着：贤者白袍/秘袍，权杖在手。梳造：王冠/法冠，威仪自生。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 威能内蕴，光晕庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 白袍尊者**
- 灰袍智者，道果圆满阶段·白袍尊者。形象：灰袍巫师，手持法杖。 核心意象：灰袍、法杖、奥术符文。神态：魔力圆满，与元素共鸣的从容。动作：大法圆满，法阵漫天归于一，法杖顿地，口中吟唱，星光自天穹应召而下。衣着：元素法袍，权杖圣光。梳造：法冠星冕，威能内蕴。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 魔力圆满，与元素共鸣的从容; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 魔杖猫，大劫淬炼阶段·魔咒猫。形象：紫毛小猫，爪握魔法杖。 核心意象：紫色毛皮、魔法杖、闪烁咒光。神态：暗夜中坚持，眼神倔强。动作：全力施法，魔力风暴，挥舞魔法杖，星光自杖尖倾泻而下。衣着：法袍破损，魔力燃烧。梳造：发丝凌乱，眸中带焰。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 暗夜中坚持，眼神倔强; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大贤者猫**
- 魔杖猫，封神登天阶段·大贤者猫。形象：紫毛小猫，爪握魔法杖。 核心意象：紫色毛皮、魔法杖、闪烁咒光。神态：威能内蕴，光晕庄严。动作：魔法大成，法阵漫天，挥舞魔法杖，星光自杖尖倾泻而下。衣着：贤者白袍/秘袍，权杖在手。梳造：王冠/法冠，威仪自生。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 威能内蕴，光晕庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 魔法尊者**
- 魔杖猫，道果圆满阶段·魔法尊者。形象：紫毛小猫，爪握魔法杖。 核心意象：紫色毛皮、魔法杖、闪烁咒光。神态：魔力圆满，与元素共鸣的从容。动作：大法圆满，法阵漫天归于一，挥舞魔法杖，星光自杖尖倾泻而下。衣着：元素法袍，权杖圣光。梳造：法冠星冕，威能内蕴。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 魔力圆满，与元素共鸣的从容; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 龙骑士，大劫淬炼阶段·骑士之誓。形象：红铠骑士，身侧有龙影。 核心意象：红铠、骑士长枪、龙之影。神态：暗夜中坚持，眼神倔强。动作：全力施法，魔力风暴，举枪跃上龙背，巨龙振翼腾空而起。衣着：法袍破损，魔力燃烧。梳造：发丝凌乱，眸中带焰。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 暗夜中坚持，眼神倔强; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 天空骑士**
- 龙骑士，封神登天阶段·天空骑士。形象：红铠骑士，身侧有龙影。 核心意象：红铠、骑士长枪、龙之影。神态：威能内蕴，光晕庄严。动作：魔法大成，法阵漫天，举枪跃上龙背，巨龙振翼腾空而起。衣着：贤者白袍/秘袍，权杖在手。梳造：王冠/法冠，威仪自生。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 威能内蕴，光晕庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 龙骑尊者**
- 龙骑士，道果圆满阶段·龙骑尊者。形象：红铠骑士，身侧有龙影。 核心意象：红铠、骑士长枪、龙之影。神态：魔力圆满，与元素共鸣的从容。动作：大法圆满，法阵漫天归于一，举枪跃上龙背，巨龙振翼腾空而起。衣着：元素法袍，权杖圣光。梳造：法冠星冕，威能内蕴。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 魔力圆满，与元素共鸣的从容; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 炼金魔像，大劫淬炼阶段·点金手。形象：金甲魔像，符文刻身。 核心意象：金甲、符文刻痕、炼金秘术。神态：淬炼中的忍耐。动作：全力受炼，火炼淬洗，符文亮起，双拳砸落震裂大地。衣着：历经淬炼，温润内敛。梳造：包浆/裂纹，岁月痕。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 淬炼中的忍耐; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 黄金石像**
- 炼金魔像，封神登天阶段·黄金石像。形象：金甲魔像，符文刻身。 核心意象：金甲、符文刻痕、炼金秘术。神态：灵光大成的庄严。动作：器灵大成，光华夺目，符文亮起，双拳砸落震裂大地。衣着：华彩流光，灵气充盈。梳造：祥纹满饰，宝光外放。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 灵光大成的庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 炼金尊者**
- 炼金魔像，道果圆满阶段·炼金尊者。形象：金甲魔像，符文刻身。 核心意象：金甲、符文刻痕、炼金秘术。神态：器灵圆满，灵性通明。动作：器道大成，光华满室，符文亮起，双拳砸落震裂大地。衣着：宝光流转，灵性通神。梳造：神纹满饰，镇世之宝。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 器灵圆满，灵性通明; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 梦魇夜马，大劫淬炼阶段·月光马。形象：紫焰骏马，四蹄踏火。 核心意象：紫焰、踏火四蹄、梦魇之力。神态：眸光如电，威严中带着坚韧。动作：全力施为，祥光与煞气并舞，四蹄踏火奔腾，紫焰在身后拖成长练。衣着：神纹满身，羽鳞如铠。梳造：羽冠/灵角，光华流转。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 眸光如电，威严中带着坚韧; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 夜之守护**
- 梦魇夜马，封神登天阶段·夜之守护。形象：紫焰骏马，四蹄踏火。 核心意象：紫焰、踏火四蹄、梦魇之力。神态：神兽威严，目光洞彻九幽。动作：绝技大成，百兽来朝，四蹄踏火奔腾，紫焰在身后拖成长练。衣着：神光加身，五色祥云。梳造：圣羽垂天，瑞角冲霄。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 神兽威严，目光洞彻九幽; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 夜月尊者**
- 梦魇夜马，道果圆满阶段·夜月尊者。形象：紫焰骏马，四蹄踏火。 核心意象：紫焰、踏火四蹄、梦魇之力。神态：瑞气圆满，祥光普照。动作：瑞兽真身，百瑞齐鸣，四蹄踏火奔腾，紫焰在身后拖成长练。衣着：五色祥光，神纹满身。梳造：圣羽垂天，瑞角生辉。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 瑞气圆满，祥光普照; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 元素灯灵，大劫淬炼阶段·流光灯灵。形象：金灯灵体，灯火摇曳。 核心意象：神灯、摇曳灯芯、许愿之光。神态：淬炼中的忍耐。动作：全力受炼，火炼淬洗，灯芯一燃，应召而出，金光漫室。衣着：历经淬炼，温润内敛。梳造：包浆/裂纹，岁月痕。意境：黑暗试炼、魔力对决，在深渊边缘淬炼。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：暗夜紫（#1A0A2E）主调 + 火焰红（#FF4500）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; abyssal trial; 淬炼中的忍耐; palette #1A0A2E with #FF4500 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 千愿灯灵**
- 元素灯灵，封神登天阶段·千愿灯灵。形象：金灯灵体，灯火摇曳。 核心意象：神灯、摇曳灯芯、许愿之光。神态：灵光大成的庄严。动作：器灵大成，光华夺目，灯芯一燃，应召而出，金光漫室。衣着：华彩流光，灵气充盈。梳造：祥纹满饰，宝光外放。意境：大法师或神话生物完全体，威震大陆，法阵与冠冕加身。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：传奇紫（#8B00FF）主调 + 王冠金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; archmage glory; 灵光大成的庄严; palette #8B00FF with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 灯灵尊者**
- 元素灯灵，道果圆满阶段·灯灵尊者。形象：金灯灵体，灯火摇曳。 核心意象：神灯、摇曳灯芯、许愿之光。神态：器灵圆满，灵性通明。动作：器道大成，光华满室，灯芯一燃，应召而出，金光漫室。衣着：宝光流转，灵性通神。梳造：神纹满饰，镇世之宝。意境：元素法相圆满，光与秘术铸就神格，威震大陆。风格：欧美奇幻概念艺术风，史诗光影、法阵与魔法光效，神秘史诗。色彩：永恒白（#F5F0E8）主调 + 圣光紫（#C8A2FF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：western fantasy concept art, epic lighting, glowing magic circles, arcane atmosphere; elemental oneness; 器灵圆满，灵性通明; palette #F5F0E8 with #C8A2FF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 霸王龙，大劫淬炼阶段·领地之王。形象：直立霸王龙，头大颚壮，前肢短小。 核心意象：巨颚利齿、粗壮后肢、白垩纪。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食，巨颚一口咬碎骨骼，怒吼震彻山谷。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 白垩纪领主**
- 霸王龙，封神登天阶段·白垩纪领主。形象：直立霸王龙，头大颚壮，前肢短小。 核心意象：巨颚利齿、粗壮后肢、白垩纪。神态：霸主之威，眼神睥睨。动作：猎技大成，万兽辟易，巨颚一口咬碎骨骼，怒吼震彻山谷。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 霸主之威，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 终极掠食者**
- 霸王龙，道果圆满阶段·终极掠食者。形象：直立霸王龙，头大颚壮，前肢短小。 核心意象：巨颚利齿、粗壮后肢、白垩纪。神态：霸主圆满，天地臣服。动作：万兽辟易，踏碎山河，巨颚一口咬碎骨骼，怒吼震彻山谷。衣着：远古霸躯，王纹生辉。梳造：巨角王冠，威压众生。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 霸主圆满，天地臣服; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 三角龙，大劫淬炼阶段·铁头龙。形象：大颈盾配三尖角，身强体壮。 核心意象：三尖角、大颈盾、群居护群。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食，低头亮出三尖角，怒吼着顶向来敌。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 群首领**
- 三角龙，封神登天阶段·群首领。形象：大颈盾配三尖角，身强体壮。 核心意象：三尖角、大颈盾、群居护群。神态：霸主之威，眼神睥睨。动作：猎技大成，万兽辟易，低头亮出三尖角，怒吼着顶向来敌。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 霸主之威，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 化石之王**
- 三角龙，道果圆满阶段·化石之王。形象：大颈盾配三尖角，身强体壮。 核心意象：三尖角、大颈盾、群居护群。神态：霸主圆满，天地臣服。动作：万兽辟易，踏碎山河，低头亮出三尖角，怒吼着顶向来敌。衣着：远古霸躯，王纹生辉。梳造：巨角王冠，威压众生。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 霸主圆满，天地臣服; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 翼龙，大劫淬炼阶段·火焰翼龙。形象：展翼无齿，颈长身轻。 核心意象：翼膜、长颈、远古天空。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食，翼膜一展，借风滑翔，俯冲掠水捕鱼。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 王者翼龙**
- 翼龙，封神登天阶段·王者翼龙。形象：展翼无齿，颈长身轻。 核心意象：翼膜、长颈、远古天空。神态：霸主之威，眼神睥睨。动作：猎技大成，万兽辟易，翼膜一展，借风滑翔，俯冲掠水捕鱼。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 霸主之威，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 太古翼龙**
- 翼龙，道果圆满阶段·太古翼龙。形象：展翼无齿，颈长身轻。 核心意象：翼膜、长颈、远古天空。神态：霸主圆满，天地臣服。动作：万兽辟易，踏碎山河，翼膜一展，借风滑翔，俯冲掠水捕鱼。衣着：远古霸躯，王纹生辉。梳造：巨角王冠，威压众生。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 霸主圆满，天地臣服; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 猛犸象，大劫淬炼阶段·象牙王。形象：长毛巨象，弯长象牙。 核心意象：弯长象牙、厚毛、冰河雪原。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，长鼻卷起枯枝，弯长象牙挑开积雪。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 冰川之主**
- 猛犸象，封神登天阶段·冰川之主。形象：长毛巨象，弯长象牙。 核心意象：弯长象牙、厚毛、冰河雪原。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，长鼻卷起枯枝，弯长象牙挑开积雪。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 目光如炬，不怒自威，威仪自生; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 远古巨兽**
- 猛犸象，道果圆满阶段·远古巨兽。形象：长毛巨象，弯长象牙。 核心意象：弯长象牙、厚毛、冰河雪原。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，长鼻卷起枯枝，弯长象牙挑开积雪。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 剑齿虎，大劫淬炼阶段·冰河猎手。形象：上颚獠牙如剑，肌肉虬结。 核心意象：剑形獠牙、虬结肌肉、一击必杀。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，低伏潜行，骤然跃起一剑封喉。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金剑齿**
- 剑齿虎，封神登天阶段·金剑齿。形象：上颚獠牙如剑，肌肉虬结。 核心意象：剑形獠牙、虬结肌肉、一击必杀。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，低伏潜行，骤然跃起一剑封喉。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 目光如炬，不怒自威，威仪自生; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 太古剑齿**
- 剑齿虎，道果圆满阶段·太古剑齿。形象：上颚獠牙如剑，肌肉虬结。 核心意象：剑形獠牙、虬结肌肉、一击必杀。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，低伏潜行，骤然跃起一剑封喉。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 沧龙，大劫淬炼阶段·王沧。形象：深海巨蜥龙，鳍状四肢。 核心意象：鳍状四肢、巨颚、苍茫汪洋。神态：龙威炽烈，怒目电光。动作：全力施为，风雷随身，鳍肢一摆，巨颚咬碎深海猎物。衣着：战损鳞甲，雷火纹显。梳造：角芒凌厉，须张如戟。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 龙威炽烈，怒目电光; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 古沧**
- 沧龙，封神登天阶段·古沧。形象：深海巨蜥龙，鳍状四肢。 核心意象：鳍状四肢、巨颚、苍茫汪洋。神态：龙目洞彻，神威赫赫。动作：绝技大成，行云布雨，鳍肢一摆，巨颚咬碎深海猎物。衣着：金鳞覆身，祥光万道。梳造：龙角如珊瑚，须垂百丈。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 龙目洞彻，神威赫赫; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 神沧**
- 沧龙，道果圆满阶段·神沧。形象：深海巨蜥龙，鳍状四肢。 核心意象：鳍状四肢、巨颚、苍茫汪洋。神态：真身圆满，龙威盖世。动作：腾云驾雾，号令风雨，鳍肢一摆，巨颚咬碎深海猎物。衣着：金鳞神光，日月同辉。梳造：龙角如岳，须垂星河。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 真身圆满，龙威盖世; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 棘龙，大劫淬炼阶段·河畔霸主。形象：背生巨大帆状脊，形似鳄龙。 核心意象：巨大帆脊、长吻利齿、河畔。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食，帆脊竖起示威，长吻探入水中叼起巨鱼。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 河畔之王**
- 棘龙，封神登天阶段·河畔之王。形象：背生巨大帆状脊，形似鳄龙。 核心意象：巨大帆脊、长吻利齿、河畔。神态：霸主之威，眼神睥睨。动作：猎技大成，万兽辟易，帆脊竖起示威，长吻探入水中叼起巨鱼。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 霸主之威，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 棘龙尊者**
- 棘龙，道果圆满阶段·棘龙尊者。形象：背生巨大帆状脊，形似鳄龙。 核心意象：巨大帆脊、长吻利齿、河畔。神态：霸主圆满，天地臣服。动作：万兽辟易，踏碎山河，帆脊竖起示威，长吻探入水中叼起巨鱼。衣着：远古霸躯，王纹生辉。梳造：巨角王冠，威压众生。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 霸主圆满，天地臣服; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 甲龙，大劫淬炼阶段·堡垒龙。形象：身披重甲，尾端生骨锤。 核心意象：重甲骨板、尾端骨锤、不动如山。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食，骨锤横扫，一击足以震退掠食者。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 要塞甲龙**
- 甲龙，封神登天阶段·要塞甲龙。形象：身披重甲，尾端生骨锤。 核心意象：重甲骨板、尾端骨锤、不动如山。神态：霸主之威，眼神睥睨。动作：猎技大成，万兽辟易，骨锤横扫，一击足以震退掠食者。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 霸主之威，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 甲龙尊者**
- 甲龙，道果圆满阶段·甲龙尊者。形象：身披重甲，尾端生骨锤。 核心意象：重甲骨板、尾端骨锤、不动如山。神态：霸主圆满，天地臣服。动作：万兽辟易，踏碎山河，骨锤横扫，一击足以震退掠食者。衣着：远古霸躯，王纹生辉。梳造：巨角王冠，威压众生。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 霸主圆满，天地臣服; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 梁龙，大劫淬炼阶段·群居龙。形象：长颈细尾的巨型蜥脚龙。 核心意象：超长脖颈、鞭状长尾、成群迁徙。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食，长颈探向高树，细尾如鞭甩动护身。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 温柔巨兽**
- 梁龙，封神登天阶段·温柔巨兽。形象：长颈细尾的巨型蜥脚龙。 核心意象：超长脖颈、鞭状长尾、成群迁徙。神态：霸主之威，眼神睥睨。动作：猎技大成，万兽辟易，长颈探向高树，细尾如鞭甩动护身。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 霸主之威，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 梁龙尊者**
- 梁龙，道果圆满阶段·梁龙尊者。形象：长颈细尾的巨型蜥脚龙。 核心意象：超长脖颈、鞭状长尾、成群迁徙。神态：霸主圆满，天地臣服。动作：万兽辟易，踏碎山河，长颈探向高树，细尾如鞭甩动护身。衣着：远古霸躯，王纹生辉。梳造：巨角王冠，威压众生。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 霸主圆满，天地臣服; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 巨齿鲨，大劫淬炼阶段·深蓝鲨。形象：巨牙利齿的远古巨鲨。 核心意象：巨牙、庞大身影、深海之渊。神态：生存竞争中的冷酷坚毅。动作：全力搏杀，与天地争食，巨口张开，一口吞下整条海兽。衣着：成体鳞甲，伤痕累累。梳造：角/帆/鬃威猛。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 生存竞争中的冷酷坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 巨齿霸主**
- 巨齿鲨，封神登天阶段·巨齿霸主。形象：巨牙利齿的远古巨鲨。 核心意象：巨牙、庞大身影、深海之渊。神态：霸主之威，眼神睥睨。动作：猎技大成，万兽辟易，巨口张开，一口吞下整条海兽。衣着：霸主之躯，王纹隐现。梳造：顶冠/长角，王者相。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 霸主之威，眼神睥睨; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 巨齿尊者**
- 巨齿鲨，道果圆满阶段·巨齿尊者。形象：巨牙利齿的远古巨鲨。 核心意象：巨牙、庞大身影、深海之渊。神态：霸主圆满，天地臣服。动作：万兽辟易，踏碎山河，巨口张开，一口吞下整条海兽。衣着：远古霸躯，王纹生辉。梳造：巨角王冠，威压众生。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 霸主圆满，天地臣服; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 大地懒，大劫淬炼阶段·耐心巨兽。形象：巨爪大懒兽，毛长体壮。 核心意象：巨型爪钩、长毛、缓慢而强大。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，巨爪刨开泥土，缓缓挪动庞然之躯。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大地巨人**
- 大地懒，封神登天阶段·大地巨人。形象：巨爪大懒兽，毛长体壮。 核心意象：巨型爪钩、长毛、缓慢而强大。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，巨爪刨开泥土，缓缓挪动庞然之躯。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 目光如炬，不怒自威，威仪自生; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 大地尊者**
- 大地懒，道果圆满阶段·大地尊者。形象：巨爪大懒兽，毛长体壮。 核心意象：巨型爪钩、长毛、缓慢而强大。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，巨爪刨开泥土，缓缓挪动庞然之躯。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 冰河犀牛，大劫淬炼阶段·破冰犀。形象：长毛犀牛，一双弯角。 核心意象：双弯角、厚长毛、冰原冻土。神态：眸光深沉，眉峰微蹙，带着淬炼的坚毅。动作：豁尽全力，势如破竹，低头拱开积雪，弯角翻出苔草为食。衣着：成体毛皮/鳞甲，或有伤痕。梳造：角/鬃/尾冠成型，凌乱而烈。意境：严酷的冰河或火山环境下生存试炼，与天争命。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：火山黑（#3E2723）主调 + 熔岩红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; survival trial; 眸光深沉，眉峰微蹙，带着淬炼的坚毅; palette #3E2723 with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 冰河之主**
- 冰河犀牛，封神登天阶段·冰河之主。形象：长毛犀牛，一双弯角。 核心意象：双弯角、厚长毛、冰原冻土。神态：目光如炬，不怒自威，威仪自生。动作：绝技大成，昂首傲立山巅，低头拱开积雪，弯角翻出苔草为食。衣着：金色祥纹缀身，王者的气象。梳造：金鬃/长羽/圣角，威风凛凛。意境：成为领地霸主，王者的姿态与威严。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：霸主棕（#4E342E）主调 + 王金（#FFB300）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; apex predator; 目光如炬，不怒自威，威仪自生; palette #4E342E with #FFB300 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 冰河尊者**
- 冰河犀牛，道果圆满阶段·冰河尊者。形象：长毛犀牛，一双弯角。 核心意象：双弯角、厚长毛、冰原冻土。神态：圆满自足，神光内蕴的从容。动作：神形合一，万兽来朝，低头拱开积雪，弯角翻出苔草为食。衣着：圣光缠身，王兽之威。梳造：神鬃/圣冠，威仪圆满。意境：远古霸主图腾显圣，冰川之巅威压万古。风格：BBC 纪录片古生物复原结合卡通拟人，鳞甲皮肤纹理真实，远古苍茫。色彩：冰川白（#F5F0E8）主调 + 远古蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：paleoart documentary style with cartoon charm, realistic scale/skin texture; ancient totem; 圆满自足，神光内蕴的从容; palette #F5F0E8 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 8. 星座守护（12 物种）

> **风格**：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。**阶段演绎**：
> - 灵胎初醒：星屑微光中孕育，小宇宙初醒，命定的星芒（星屑银/白银）
> - 凡尘砺心：圣斗士修行岁月，拳法初练，热血汗水（青铜金/银灰）
> - 道法初成：第七感觉醒，星座之光闪耀，必杀技初现（黄金/星座蓝）
> - 大劫淬炼：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙（暗紫/圣战红）
> - 封神登天：身着黄金圣衣，镇守宫门，星座图腾全开（黄金圣衣/圣光白）
> - 道果圆满：第八感领悟，神之领域，小宇宙化作无尽星光（神光白/小宇宙金）

#### 白羊座·穆（`aries`）

**灵胎初醒 · 星屑之种**
- 白羊座·穆，灵胎初醒阶段·星屑之种。初始形态：一团星屑灵光，白羊命星在其中闪烁，暖金色小宇宙初醒，羊角般的星光微凝。火属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星屑虚影，白羊命星微光。梳造：束发，眉心一点星芒。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 念动初觉**
- 白羊座·穆，凡尘砺心阶段·念动初觉。形象：金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：素白练功服，圣衣未成。梳造：束发，汗珠映星。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 水晶微光**
- 白羊座·穆，道法初成阶段·水晶微光。形象：金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，水晶墙一筑，星屑旋转功。衣着：青铜圣衣初覆肩甲。梳造：羊角头冠微启。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 水晶墙**
- 白羊座·穆，大劫淬炼阶段·水晶墙。形象：金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，水晶墙一筑，星屑旋转功。衣着：白银圣衣，水晶墙光。梳造：银盔，眉宇沉静。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 白羊宫主**
- 白羊座·穆，封神登天阶段·白羊宫主。形象：金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，水晶墙一筑，星屑旋转功。衣着：黄金圣衣全装，白羊宫主。梳造：白羊金盔，宝相庄严。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 星光灭绝**
- 白羊座·穆，道果圆满阶段·星光灭绝。形象：金色羊首圣衣星灵，角如新月。 核心意象：白羊宫、星光灭绝。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，水晶墙一筑，星屑旋转功。衣着：神圣衣星辉，星光灭绝之姿。梳造：星光化冠，超脱凡尘。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 金牛座·阿鲁迪巴（`taurus`）

**灵胎初醒 · 星辉之种**
- 金牛座·阿鲁迪巴，灵胎初醒阶段·星辉之种。初始形态：一团星辉灵光，金牛星芒沉稳厚重，大地般的小宇宙在星辉中缓缓脉动。土属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉灵光，金牛星芒沉稳。梳造：粗发，魁伟身形。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 蛮力渐长**
- 金牛座·阿鲁迪巴，凡尘砺心阶段·蛮力渐长。形象：金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，赤膊练力。梳造：短发如刺，汗如雨下。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 角力初成**
- 金牛座·阿鲁迪巴，道法初成阶段·角力初成。形象：金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，巨型号角一吼，山岳为之震颤。衣着：青铜圣衣，力士之姿。梳造：牛角头冠初现。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 巨型号角**
- 金牛座·阿鲁迪巴，大劫淬炼阶段·巨型号角。形象：金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，巨型号角一吼，山岳为之震颤。衣着：白银圣衣，巨型号角。梳造：双角盔，怒目圆睁。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金牛宫主**
- 金牛座·阿鲁迪巴，封神登天阶段·金牛宫主。形象：金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，巨型号角一吼，山岳为之震颤。衣着：黄金圣衣，金牛宫主。梳造：金牛金盔，力撼山岳。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 黄金之角**
- 金牛座·阿鲁迪巴，道果圆满阶段·黄金之角。形象：金色牛首圣衣星灵，体格魁伟，双角如巨柱。 核心意象：金牛宫、巨型号角、黄金双角。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，巨型号角一吼，山岳为之震颤。衣着：黄金之角，仁厚巨人。梳造：金盔生辉，憨厚威仪。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 双子座·撒加（`gemini`）

**灵胎初醒 · 星辉双子**
- 双子座·撒加，灵胎初醒阶段·星辉双子。初始形态：两团交缠的星辉灵光，一明一暗宛如双子，善恶之争在风中无声展开。风属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：双子星辉，一明一暗。梳造：双生虚影，善恶交织。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 善恶交织**
- 双子座·撒加，凡尘砺心阶段·善恶交织。形象：黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，温雅少年。梳造：束发，眼神却晦暗。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 异次元微光**
- 双子座·撒加，道法初成阶段·异次元微光。形象：黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，银河星爆一击，群星为之破碎。衣着：青铜圣衣，善恶分形。梳造：半面面具，半面光明。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 异次元空间**
- 双子座·撒加，大劫淬炼阶段·异次元空间。形象：黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，银河星爆一击，群星为之破碎。衣着：异次元空间，白银圣衣。梳造：黑发狂舞，邪气外露。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 双子宫主**
- 双子座·撒加，封神登天阶段·双子宫主。形象：黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，银河星爆一击，群星为之破碎。衣着：黄金圣衣，双子宫主。梳造：双面金盔，善恶并存。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 银河星爆**
- 双子座·撒加，道果圆满阶段·银河星爆。形象：黑金双面圣衣星灵，一体双生，善恶并存。 核心意象：双子宫、银河星爆、善恶双面。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，银河星爆一击，群星为之破碎。衣着：银河星爆，神的化身。梳造：神光化发，一念神魔。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 巨蟹座·迪斯马斯克（`cancer`）

**灵胎初醒 · 冥辉之种**
- 巨蟹座·迪斯马斯克，灵胎初醒阶段·冥辉之种。初始形态：一团冥辉灵光，蟹钳般的紫芒时隐时现，冥界幽光在蛋形光晕中流转。水属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：冥辉灵光，幽光沉沉。梳造：发丝带冥火。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 冥气初生**
- 巨蟹座·迪斯马斯克，凡尘砺心阶段·冥气初生。形象：青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，邪气未显。梳造：束发，嘴角冷笑。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冥界门扉**
- 巨蟹座·迪斯马斯克，道法初成阶段·冥界门扉。形象：青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，积尸气冥界波，将对手拖入黄泉。衣着：青铜圣衣，冥气缠绕。梳造：蟹甲头冠，诡异。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 积尸气冥界波**
- 巨蟹座·迪斯马斯克，大劫淬炼阶段·积尸气冥界波。形象：青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，积尸气冥界波，将对手拖入黄泉。衣着：积尸气冥界波，冥光冲霄。梳造：冥火为发，阴笑森然。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 巨蟹宫主**
- 巨蟹座·迪斯马斯克，封神登天阶段·巨蟹宫主。形象：青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，积尸气冥界波，将对手拖入黄泉。衣着：黄金圣衣，巨蟹宫主。梳造：金盔蟹钳，邪中带狂。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 冥界之王**
- 巨蟹座·迪斯马斯克，道果圆满阶段·冥界之王。形象：青灰蟹形圣衣星灵，甲壳泛着冥界幽光。 核心意象：巨蟹宫、积尸气、冥界之门。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，积尸气冥界波，将对手拖入黄泉。衣着：冥界之王，身陷幽暗。梳造：冥光吞没，亦正亦邪。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 狮子座·艾奥里亚（`leo`）

**灵胎初醒 · 星辉之种**
- 狮子座·艾奥里亚，灵胎初醒阶段·星辉之种。初始形态：一团星辉灵光，狮鬃般的金色芒焰喷张，战意如火的小宇宙炽热燃烧。火属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉灵光，狮鬃般的金焰。梳造：金发，烈如狮鬃。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 热血修行**
- 狮子座·艾奥里亚，凡尘砺心阶段·热血修行。形象：金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，热血修行。梳造：束发，汗意蒸腾。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 拳光初现**
- 狮子座·艾奥里亚，道法初成阶段·拳光初现。形象：金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，闪电光速拳连发，拳影如雷霆万钧。衣着：青铜圣衣，拳光初现。梳造：狮盔微启，眸光如电。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 闪电光速拳**
- 狮子座·艾奥里亚，大劫淬炼阶段·闪电光速拳。形象：金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，闪电光速拳连发，拳影如雷霆万钧。衣着：闪电光速拳，战意熊熊。梳造：金发飞扬，怒目如狮。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 狮子宫主**
- 狮子座·艾奥里亚，封神登天阶段·狮子宫主。形象：金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，闪电光速拳连发，拳影如雷霆万钧。衣着：黄金圣衣，狮子宫主。梳造：黄金狮盔，威风凛凛。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 等离子光速拳**
- 狮子座·艾奥里亚，道果圆满阶段·等离子光速拳。形象：金鬃狮首圣衣星灵，威光赫赫，战意熊熊。 核心意象：狮子宫、闪电光速拳、黄金鬃毛。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，闪电光速拳连发，拳影如雷霆万钧。衣着：等离子光速拳，正义之光。梳造：金光化冠，浩然正气。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 处女座·沙加（`virgo`）

**灵胎初醒 · 星辉莲花**
- 处女座·沙加，灵胎初醒阶段·星辉莲花。初始形态：一朵星辉莲花，白金花瓣层层合拢，闭目禅相在莲心流转，最接近神的光。土属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉莲花，禅意初凝。梳造：束发，眉目安详。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 禅心初定**
- 处女座·沙加，凡尘砺心阶段·禅心初定。形象：白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，静坐冥想。梳造：束发，闭目禅坐。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 天眼微启**
- 处女座·沙加，道法初成阶段·天眼微启。形象：白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，天魔降伏，闭目之间神光自生。衣着：青铜圣衣，禅相庄严。梳造：发髻微束，宝相初显。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 天魔降伏**
- 处女座·沙加，大劫淬炼阶段·天魔降伏。形象：白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，天魔降伏，闭目之间神光自生。衣着：天魔降伏，神光内敛。梳造：闭目怒相，禅威并具。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 处女宫主**
- 处女座·沙加，封神登天阶段·处女宫主。形象：白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，天魔降伏，闭目之间神光自生。衣着：黄金圣衣，处女宫主。梳造：金莲花冠，佛陀相。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 第八感觉醒**
- 处女座·沙加，道果圆满阶段·第八感觉醒。形象：白金圣衣星灵，双手合十，双目紧闭。 核心意象：处女宫、天舞宝轮、闭目禅相。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，天魔降伏，闭目之间神光自生。衣着：第八感觉醒，最接近神。梳造：圣光化莲，眉目悲悯。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 天秤座·童虎（`libra`）

**灵胎初醒 · 星辉龙珠**
- 天秤座·童虎，灵胎初醒阶段·星辉龙珠。初始形态：一粒星辉龙珠，天秤虚影与十二道金光浮沉，守望者的龙息在珠中沉眠。风属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉龙珠，天秤虚影。梳造：束发，苍老之相。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 百裂拳初修**
- 天秤座·童虎，凡尘砺心阶段·百裂拳初修。形象：白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：五老峰练功，布衣。梳造：须发皆白，佝偻。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 五老峰悟道**
- 天秤座·童虎，道法初成阶段·五老峰悟道。形象：白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，庐山百龙霸，龙吟声中百龙齐出。衣着：青铜圣衣，庐山升龙霸。梳造：白发，龙吟初响。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 庐山升龙霸**
- 天秤座·童虎，大劫淬炼阶段·庐山升龙霸。形象：白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，庐山百龙霸，龙吟声中百龙齐出。衣着：返老还童，白银圣衣。梳造：银发转黑，青年之相。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 老师坐镇**
- 天秤座·童虎，封神登天阶段·老师坐镇。形象：白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，庐山百龙霸，龙吟声中百龙齐出。衣着：黄金圣衣，天秤宫主。梳造：金盔，十二件武器负背。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 返老还童**
- 天秤座·童虎，道果圆满阶段·返老还童。形象：白金圣衣星灵，背负天秤与十二件黄金武器。 核心意象：天秤宫、十二件黄金武器、五老峰。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，庐山百龙霸，龙吟声中百龙齐出。衣着：老师坐镇，两百年守望。梳造：白发金盔，目光如渊。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 天蝎座·米罗（`scorpio`）

**灵胎初醒 · 星辉之种**
- 天蝎座·米罗，灵胎初醒阶段·星辉之种。初始形态：一团星辉灵光，蝎尾般的猩红针芒微露，深红毒针的锐意在星辉中蛰伏。水属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉灵光，猩红针芒微露。梳造：束发，眼神犀利。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 忠义修行**
- 天蝎座·米罗，凡尘砺心阶段·忠义修行。形象：深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，忠义修行。梳造：束发，沉稳。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 毒针微芒**
- 天蝎座·米罗，道法初成阶段·毒针微芒。形象：深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，猩红毒针疾点，十五针封喉。衣着：青铜圣衣，毒针初芒。梳造：发间一点猩红。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 深红毒针**
- 天蝎座·米罗，大劫淬炼阶段·深红毒针。形象：深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，猩红毒针疾点，十五针封喉。衣着：深红毒针，蝎尾流光。梳造：红发如焰，凛然。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 天蝎宫主**
- 天蝎座·米罗，封神登天阶段·天蝎宫主。形象：深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，猩红毒针疾点，十五针封喉。衣着：黄金圣衣，天蝎宫主。梳造：蝎尾金盔，锐利。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 天蝎之怒**
- 天蝎座·米罗，道果圆满阶段·天蝎之怒。形象：深红蝎形圣衣星灵，尾针流转猩红星光。 核心意象：天蝎宫、深红毒针、蝎尾之光。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，猩红毒针疾点，十五针封喉。衣着：天蝎之怒，深红针下留情。梳造：金盔生辉，义字当先。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 射手座·艾俄洛斯（`sagittarius`）

**灵胎初醒 · 星辉之箭**
- 射手座·艾俄洛斯，灵胎初醒阶段·星辉之箭。初始形态：一支星辉之箭，黄金弓芒凝成箭形，射穿苍穹的意志在星光中引而不发。火属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉之箭，弓芒初凝。梳造：束发，俊朗。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 弓弦初张**
- 射手座·艾俄洛斯，凡尘砺心阶段·弓弦初张。形象：金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，弓弦初张。梳造：束发，目光如鹰。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 闪电拳雏形**
- 射手座·艾俄洛斯，道法初成阶段·闪电拳雏形。形象：金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，黄金之箭搭弦，一箭射穿苍穹。衣着：青铜圣衣，黄金之箭雏形。梳造：发丝随风，弓在手。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 黄金之箭**
- 射手座·艾俄洛斯，大劫淬炼阶段·黄金之箭。形象：金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，黄金之箭搭弦，一箭射穿苍穹。衣着：黄金之箭，守护而战。梳造：金发飞扬，英烈之气。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 英魂长存**
- 射手座·艾俄洛斯，封神登天阶段·英魂长存。形象：金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，黄金之箭搭弦，一箭射穿苍穹。衣着：黄金圣衣，射手宫主。梳造：金盔，弓满月。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 黄金之箭·觉醒**
- 射手座·艾俄洛斯，道果圆满阶段·黄金之箭·觉醒。形象：金光圣衣星灵，背负黄金之弓，箭指苍穹。 核心意象：射手宫、黄金之弓、群星之矢。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，黄金之箭搭弦，一箭射穿苍穹。衣着：英魂长存，星空中守望。梳造：星光化发，永恒之姿。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 摩羯座·修罗（`capricorn`）

**灵胎初醒 · 星辉之剑**
- 摩羯座·修罗，灵胎初醒阶段·星辉之剑。初始形态：一柄星辉之剑，银白剑芒凝成圣剑雏形，Excalibur的锋锐在土色星光中成形。土属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉之剑，剑芒初凝。梳造：束发，冷峻。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 剑士修行**
- 摩羯座·修罗，凡尘砺心阶段·剑士修行。形象：银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，剑士修行。梳造：束发，沉静如渊。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 圣剑雏形**
- 摩羯座·修罗，道法初成阶段·圣剑雏形。形象：银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，双臂化剑，圣剑一挥可斩星辰。衣着：青铜圣衣，圣剑雏形。梳造：发间剑光隐现。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 剑影纵横**
- 摩羯座·修罗，大劫淬炼阶段·剑影纵横。形象：银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，双臂化剑，圣剑一挥可斩星辰。衣着：圣剑Excalibur，剑影纵横。梳造：黑发如墨，眸光如剑。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 摩羯宫主**
- 摩羯座·修罗，封神登天阶段·摩羯宫主。形象：银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，双臂化剑，圣剑一挥可斩星辰。衣着：黄金圣衣，摩羯宫主。梳造：金盔，双臂为剑。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 觉醒之剑**
- 摩羯座·修罗，道果圆满阶段·觉醒之剑。形象：银白圣衣星灵，双臂凝作圣剑锋芒。 核心意象：摩羯宫、圣剑Excalibur、双臂剑芒。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，双臂化剑，圣剑一挥可斩星辰。衣着：觉醒之剑，忠诚所向。梳造：剑光化发，一生为剑。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 水瓶座·卡妙（`aquarius`）

**灵胎初醒 · 星辉冰晶**
- 水瓶座·卡妙，灵胎初醒阶段·星辉冰晶。初始形态：一粒星辉冰晶，冰蓝灵光凝霜成晶，绝对零度的气息在寒光中沉淀。冰属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉冰晶，寒气初凝。梳造：束发，清冷。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 冰道修行**
- 水瓶座·卡妙，凡尘砺心阶段·冰道修行。形象：冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，冰道修行。梳造：束发，凝霜为汗。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冰晶环雏形**
- 水瓶座·卡妙，道法初成阶段·冰晶环雏形。形象：冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，曙光女神之宽恕，一拳挥出绝对零度。衣着：青铜圣衣，冰晶环雏形。梳造：发间凝霜，眼如冰湖。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 绝对零度之道**
- 水瓶座·卡妙，大劫淬炼阶段·绝对零度之道。形象：冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，曙光女神之宽恕，一拳挥出绝对零度。衣着：绝对零度之道，冰封千里。梳造：白发冰霜，无情之冷。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 水瓶宫主**
- 水瓶座·卡妙，封神登天阶段·水瓶宫主。形象：冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，曙光女神之宽恕，一拳挥出绝对零度。衣着：黄金圣衣，水瓶宫主。梳造：冰蓝金盔，严师之姿。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 冰河之师**
- 水瓶座·卡妙，道果圆满阶段·冰河之师。形象：冰蓝圣衣星灵，周身凝霜，水光流转。 核心意象：水瓶宫、绝对零度、冰蓝水光。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，曙光女神之宽恕，一拳挥出绝对零度。衣着：冰河之师，严中带暖。梳造：冰发微融，目光温柔。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 双鱼座·阿布罗狄（`pisces`）

**灵胎初醒 · 星辉花苞**
- 双鱼座·阿布罗狄，灵胎初醒阶段·星辉花苞。初始形态：一朵星辉花苞，绯红花瓣裹着双鱼游影，玫瑰之毒的暗香在灵光中浮动。水属性灵光微微环绕。神态：命星微光中沉睡，小宇宙初醒而不自知。动作：星屑环绕，静立若定。衣着：星辉花苞，绯红暗香。梳造：束发，俊美如画。意境：星屑微光中孕育，小宇宙初醒，命定的星芒。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：星屑银（#E8E8F0）主调 + 白银（#C0C0C0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; cosmic birth; 命星微光中沉睡，小宇宙初醒而不自知; palette #E8E8F0 with #C0C0C0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 花之修行**
- 双鱼座·阿布罗狄，凡尘砺心阶段·花之修行。形象：绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。神态：修行中的专注，汗珠映着星辉。动作：挥拳踢腿，苦练招式。衣着：练功服，花之修行。梳造：束发，眉目如诗。意境：圣斗士修行岁月，拳法初练，热血汗水。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：青铜金（#D4AF37）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; saint training; 修行中的专注，汗珠映着星辉; palette #D4AF37 with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 恶魔玫瑰雏形**
- 双鱼座·阿布罗狄，道法初成阶段·恶魔玫瑰雏形。形象：绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。神态：为守护而战的炽热，第七感觉醒。动作：绝技初现，拳风呼啸，皇家魔宫玫瑰，漫天血玫飞射。衣着：青铜圣衣，恶魔玫瑰雏形。梳造：发间玫瑰初绽。意境：第七感觉醒，星座之光闪耀，必杀技初现。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金（#FFD700）主调 + 星座蓝（#87CEEB）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; seventh sense; 为守护而战的炽热，第七感觉醒; palette #FFD700 with #87CEEB accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 食人鱼玫瑰**
- 双鱼座·阿布罗狄，大劫淬炼阶段·食人鱼玫瑰。形象：绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。神态：绝境中燃烧小宇宙，眸光不屈。动作：全力施展，血战不退，皇家魔宫玫瑰，漫天血玫飞射。衣着：食人鱼玫瑰，血玫飞射。梳造：金发染血，凄美。意境：圣战死斗，伤痕与信念，在绝境中燃烧小宇宙。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：暗紫（#4B0082）主调 + 圣战红（#DC143C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; holy war trial; 绝境中燃烧小宇宙，眸光不屈; palette #4B0082 with #DC143C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 双鱼宫主**
- 双鱼座·阿布罗狄，封神登天阶段·双鱼宫主。形象：绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。神态：金眸光华内敛，不怒自威。动作：绝技大成，镇守宫门，皇家魔宫玫瑰，漫天血玫飞射。衣着：黄金圣衣，双鱼宫主。梳造：玫瑰金盔，美中藏毒。意境：身着黄金圣衣，镇守宫门，星座图腾全开。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：黄金圣衣（#FFD700）主调 + 圣光白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; golden cloth; 金眸光华内敛，不怒自威; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 凄美绽放**
- 双鱼座·阿布罗狄，道果圆满阶段·凄美绽放。形象：绯红圣衣星灵，周身缠绕玫瑰藤蔓。 核心意象：双鱼宫、魔宫玫瑰、血色藤蔓。神态：第八感圆满，战意与神性合一。动作：绝技巅峰，星碎神裂，皇家魔宫玫瑰，漫天血玫飞射。衣着：凄美绽放，玫瑰成刃。梳造：花瓣化发，绝美之姿。意境：第八感领悟，神之领域，小宇宙化作无尽星光。风格：华丽圣衣战士风，金属圣衣光泽，希腊神话加黄道星座图腾，神圣崇高。色彩：神光白（#FFFFFF）主调 + 小宇宙金（#E8C87A）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：golden armor saint warrior style, polished metal sheen, Greek myth with zodiac motifs; eighth sense transcendence; 第八感圆满，战意与神性合一; palette #FFFFFF with #E8C87A accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 粽子，大劫淬炼阶段·龙舟粽。形象：翠绿粽形，系着绳结。 核心意象：箬叶、绳结、端午龙舟。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，剥开箬叶，露出饱满晶莹的糯米。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 五月粽**
- 粽子，封神登天阶段·五月粽。形象：翠绿粽形，系着绳结。 核心意象：箬叶、绳结、端午龙舟。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，剥开箬叶，露出饱满晶莹的糯米。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 屈子之志**
- 粽子，道果圆满阶段·屈子之志。形象：翠绿粽形，系着绳结。 核心意象：箬叶、绳结、端午龙舟。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，剥开箬叶，露出饱满晶莹的糯米。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 汤圆，大劫淬炼阶段·元宵。形象：雪白圆润，浮于甜汤。 核心意象：雪白圆身、甜汤、团圆之月。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，在甜汤里轻轻翻滚，一口咬下软糯流心。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 元夕汤圆**
- 汤圆，封神登天阶段·元夕汤圆。形象：雪白圆润，浮于甜汤。 核心意象：雪白圆身、甜汤、团圆之月。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，在甜汤里轻轻翻滚，一口咬下软糯流心。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 玉润尊者**
- 汤圆，道果圆满阶段·玉润尊者。形象：雪白圆润，浮于甜汤。 核心意象：雪白圆身、甜汤、团圆之月。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，在甜汤里轻轻翻滚，一口咬下软糯流心。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 月饼，大劫淬炼阶段·冰皮流心。形象：棕金饼形，压花纹章。 核心意象：饼上花印、满月、玉兔。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，刀锋落下，露出咸蛋黄流心。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 中秋月饼**
- 月饼，封神登天阶段·中秋月饼。形象：棕金饼形，压花纹章。 核心意象：饼上花印、满月、玉兔。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，刀锋落下，露出咸蛋黄流心。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 明月尊者**
- 月饼，道果圆满阶段·明月尊者。形象：棕金饼形，压花纹章。 核心意象：饼上花印、满月、玉兔。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，刀锋落下，露出咸蛋黄流心。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 青团，大劫淬炼阶段·清明团。形象：碧绿糯米团，艾草清香。 核心意象：艾草青碧、糯香、清明踏青。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，揭开蒸笼，艾草清香扑面而来。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 春意青团**
- 青团，封神登天阶段·春意青团。形象：碧绿糯米团，艾草清香。 核心意象：艾草青碧、糯香、清明踏青。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，揭开蒸笼，艾草清香扑面而来。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 青翠尊者**
- 青团，道果圆满阶段·青翠尊者。形象：碧绿糯米团，艾草清香。 核心意象：艾草青碧、糯香、清明踏青。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，揭开蒸笼，艾草清香扑面而来。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 重阳糕，大劫淬炼阶段·步步糕。形象：金色米糕，层叠如塔。 核心意象：层叠米糕、登高茱萸、金秋。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，层层剥开如塔，撒上红绿果脯。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 福寿糕**
- 重阳糕，封神登天阶段·福寿糕。形象：金色米糕，层叠如塔。 核心意象：层叠米糕、登高茱萸、金秋。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，层层剥开如塔，撒上红绿果脯。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 登高尊者**
- 重阳糕，道果圆满阶段·登高尊者。形象：金色米糕，层叠如塔。 核心意象：层叠米糕、登高茱萸、金秋。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，层层剥开如塔，撒上红绿果脯。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 年糕，大劫淬炼阶段·步步高。形象：白糯方糕，软糯弹牙。 核心意象：白糯方糕、新春对联、年年高升。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，油锅煎至金黄，外脆里糯夹起拉丝。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 福禄年糕**
- 年糕，封神登天阶段·福禄年糕。形象：白糯方糕，软糯弹牙。 核心意象：白糯方糕、新春对联、年年高升。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，油锅煎至金黄，外脆里糯夹起拉丝。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 年糕尊者**
- 年糕，道果圆满阶段·年糕尊者。形象：白糯方糕，软糯弹牙。 核心意象：白糯方糕、新春对联、年年高升。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，油锅煎至金黄，外脆里糯夹起拉丝。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 腊八粥，大劫淬炼阶段·百谷粥。形象：五谷杂粮熬成的浓粥。 核心意象：五谷杂粮、腊八蒜、暖冬。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，舀起一勺，红枣桂圆莲子浮沉其间。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 融和大师**
- 腊八粥，封神登天阶段·融和大师。形象：五谷杂粮熬成的浓粥。 核心意象：五谷杂粮、腊八蒜、暖冬。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，舀起一勺，红枣桂圆莲子浮沉其间。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 百味尊者**
- 腊八粥，道果圆满阶段·百味尊者。形象：五谷杂粮熬成的浓粥。 核心意象：五谷杂粮、腊八蒜、暖冬。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，舀起一勺，红枣桂圆莲子浮沉其间。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 春饼，大劫淬炼阶段·薄脆饼。形象：薄嫩面饼，卷着春蔬。 核心意象：薄饼、春蔬、立春之绿。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，摊开薄饼，卷入时令春蔬一卷而食。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 咬春大师**
- 春饼，封神登天阶段·咬春大师。形象：薄嫩面饼，卷着春蔬。 核心意象：薄饼、春蔬、立春之绿。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，摊开薄饼，卷入时令春蔬一卷而食。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 春味尊者**
- 春饼，道果圆满阶段·春味尊者。形象：薄嫩面饼，卷着春蔬。 核心意象：薄饼、春蔬、立春之绿。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，摊开薄饼，卷入时令春蔬一卷而食。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 冰糖葫芦，大劫淬炼阶段·糖葫芦串。形象：冰糖山楂串，红亮晶莹。 核心意象：红亮山楂、冰糖衣、竹签串。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，举着竹签一咬，"咔嚓"咬开糖衣。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 甜蜜大师**
- 冰糖葫芦，封神登天阶段·甜蜜大师。形象：冰糖山楂串，红亮晶莹。 核心意象：红亮山楂、冰糖衣、竹签串。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，举着竹签一咬，"咔嚓"咬开糖衣。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 酸甜尊者**
- 冰糖葫芦，道果圆满阶段·酸甜尊者。形象：冰糖山楂串，红亮晶莹。 核心意象：红亮山楂、冰糖衣、竹签串。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，举着竹签一咬，"咔嚓"咬开糖衣。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 桂花糕，大劫淬炼阶段·月桂糕。形象：米白方糕，点缀桂花瓣。 核心意象：米白方糕、金桂花瓣、秋香。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，拈起一块，桂香在齿间化开。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 金桂之王**
- 桂花糕，封神登天阶段·金桂之王。形象：米白方糕，点缀桂花瓣。 核心意象：米白方糕、金桂花瓣、秋香。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，拈起一块，桂香在齿间化开。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 桂香尊者**
- 桂花糕，道果圆满阶段·桂香尊者。形象：米白方糕，点缀桂花瓣。 核心意象：米白方糕、金桂花瓣、秋香。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，拈起一块，桂香在齿间化开。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 馄饨，大劫淬炼阶段·鲜肉馄饨。形象：薄皮小馄饨，汤鲜味美。 核心意象：薄皮、热汤、虾皮紫菜。神态：火候淬炼的专注。动作：经受蒸煮/煎炸，坚韧定型，一勺舀起，薄皮裹馅在热汤里打个转。衣着：色泽加深，香气酝酿。梳造：糖衣/油光初亮。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 火候淬炼的专注; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 大厨馄饨**
- 馄饨，封神登天阶段·大厨馄饨。形象：薄皮小馄饨，汤鲜味美。 核心意象：薄皮、热汤、虾皮紫菜。神态：出锅的骄傲，暖意融融。动作：热气腾腾，登堂亮相，一勺舀起，薄皮裹馅在热汤里打个转。衣着：成品佳肴，色香味全。梳造：装饰华美，名品之姿。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 出锅的骄傲，暖意融融; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 馄饨尊者**
- 馄饨，道果圆满阶段·馄饨尊者。形象：薄皮小馄饨，汤鲜味美。 核心意象：薄皮、热汤、虾皮紫菜。神态：名品圆满，色香俱全的骄傲。动作：登堂入室，香飘满座，一勺舀起，薄皮裹馅在热汤里打个转。衣着：华美名品，灵光透色。梳造：宝光流转，镇席之品。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 名品圆满，色香俱全的骄傲; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

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
- 花灯，大劫淬炼阶段·宫灯。形象：红彩花灯，烛光摇曳。 核心意象：红彩绢纱、烛火、灯谜。神态：淬炼中的忍耐。动作：全力受炼，火炼淬洗，烛火摇曳，光影在绢纱上流转。衣着：历经淬炼，温润内敛。梳造：包浆/裂纹，岁月痕。意境：阖家团圆的时刻，热气腾腾，温暖满溢。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：团圆红（#C0392B）主调 + 蒸金（#F39C12）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; family reunion; 淬炼中的忍耐; palette #C0392B with #F39C12 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 灯会主角**
- 花灯，封神登天阶段·灯会主角。形象：红彩花灯，烛光摇曳。 核心意象：红彩绢纱、烛火、灯谜。神态：灵光大成的庄严。动作：器灵大成，光华夺目，烛火摇曳，光影在绢纱上流转。衣着：华彩流光，灵气充盈。梳造：祥纹满饰，宝光外放。意境：成为节日主角与名品，家家户户的吉祥符号。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：主角红（#E74C3C）主调 + 名品金（#FFD700）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; iconic festive food; 灵光大成的庄严; palette #E74C3C with #FFD700 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 灯影尊者**
- 花灯，道果圆满阶段·灯影尊者。形象：红彩花灯，烛光摇曳。 核心意象：红彩绢纱、烛火、灯谜。神态：器灵圆满，灵性通明。动作：器道大成，光华满室，烛火摇曳，光影在绢纱上流转。衣着：宝光流转，灵性通神。梳造：神纹满饰，镇世之宝。意境：化为民俗文化图腾，代代相传，温暖永续。风格：民俗吉祥物拟人插画风，年画配色，圆润喜庆，温暖治愈。色彩：民俗暖白（#FFF8E7）主调 + 传承金（#FFC93C）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：Chinese folk festival mascot illustration, nianhua palette, round and festive; cultural totem; 器灵圆满，灵性通明; palette #FFF8E7 with #FFC93C accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

### 10. 虹猫蓝兔七侠传（9 物种）

> **风格**：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。**阶段演绎**：
> - 灵胎初醒：剑意种子萌芽，江湖初闻，少年侠气的微光（剑意青/银灰）
> - 凡尘砺心：习武入门，剑招初练，磨砺锋芒（侠义蓝/热血红）
> - 道法初成：剑法初成，行走江湖，惩恶扬善（江湖橙/侠气蓝）
> - 大劫淬炼：七剑合璧，魔教大战，刀光剑影中淬炼（魔教紫/血战红）
> - 封神登天：七侠之名扬天下，快意恩仇，豪情万丈（七剑金/侠名白）
> - 道果圆满：剑道圆满，一剑开天，守护苍生的宗师气象（剑道白/圆满蓝）

#### 虹猫（`hongmao`）

**灵胎初醒 · 剑意种子**
- 虹猫，灵胎初醒阶段·剑意种子。初始形态：一缕赤虹剑意凝成的种子，剑光如丝缠绕，隐约映出红衣小猫的虚影。火属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 红衣小猫**
- 虹猫，凡尘砺心阶段·红衣小猫。形象：红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 长虹剑出**
- 虹猫，道法初成阶段·长虹剑出。形象：红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，长虹贯日，一剑劈出七色虹光。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 侠义之心**
- 虹猫，大劫淬炼阶段·侠义之心。形象：红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，长虹贯日，一剑劈出七色虹光。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 七剑之首**
- 虹猫，封神登天阶段·七剑之首。形象：红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，长虹贯日，一剑劈出七色虹光。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 虹猫，道果圆满阶段·剑之尊者。形象：红毛小猫，手持长虹剑。 核心意象：长虹剑、七侠之首、家园守护。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，长虹贯日，一剑劈出七色虹光。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 蓝兔（`lantu`）

**灵胎初醒 · 冰魄剑意**
- 蓝兔，灵胎初醒阶段·冰魄剑意。初始形态：一缕冰魄剑意凝成的种子，寒光凝霜成种，蓝衣玉兔的轮廓在冰魄中浮现。冰属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 蓝衣兔**
- 蓝兔，凡尘砺心阶段·蓝衣兔。形象：蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 冰魄剑出**
- 蓝兔，道法初成阶段·冰魄剑出。形象：蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，冰魄寒光一闪，剑锋凝霜封百里。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 寒光剑影**
- 蓝兔，大劫淬炼阶段·寒光剑影。形象：蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，冰魄寒光一闪，剑锋凝霜封百里。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 玉兔剑仙**
- 蓝兔，封神登天阶段·玉兔剑仙。形象：蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，冰魄寒光一闪，剑锋凝霜封百里。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 蓝兔，道果圆满阶段·剑之尊者。形象：蓝毛白兔，冰魄剑在手。 核心意象：冰魄剑、玉蟾宫、七剑之灵。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，冰魄寒光一闪，剑锋凝霜封百里。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 逗逗（`doudou`）

**灵胎初醒 · 雨花剑意**
- 逗逗，灵胎初醒阶段·雨花剑意。初始形态：一缕雨花剑意凝成的种子，剑光如雨丝纷扬，白绒小狗与药囊的虚影若隐若现。水属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 白绒小狗**
- 逗逗，凡尘砺心阶段·白绒小狗。形象：白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 雨花剑出**
- 逗逗，道法初成阶段·雨花剑出。形象：白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，雨花剑挥洒，剑气如雨点纷飞。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 悬壶剑医**
- 逗逗，大劫淬炼阶段·悬壶剑医。形象：白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，雨花剑挥洒，剑气如雨点纷飞。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 雨花医圣**
- 逗逗，封神登天阶段·雨花医圣。形象：白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，雨花剑挥洒，剑气如雨点纷飞。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 逗逗，道果圆满阶段·剑之尊者。形象：白色小狗，爱玩爱闹。 核心意象：雨花剑、小药囊、七侠开心果。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，雨花剑挥洒，剑气如雨点纷飞。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 大奔（`dabeng`）

**灵胎初醒 · 奔雷剑意**
- 大奔，灵胎初醒阶段·奔雷剑意。初始形态：一缕奔雷剑意凝成的种子，雷光在剑种中轰鸣，棕熊神力的轮廓在雷芒中显现。雷属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 棕色小熊**
- 大奔，凡尘砺心阶段·棕色小熊。形象：棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 奔雷剑出**
- 大奔，道法初成阶段·奔雷剑出。形象：棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，奔雷剑劈下，剑势如奔雷滚滚。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 力拔山兮**
- 大奔，大劫淬炼阶段·力拔山兮。形象：棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，奔雷剑劈下，剑势如奔雷滚滚。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 奔雷剑圣**
- 大奔，封神登天阶段·奔雷剑圣。形象：棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，奔雷剑劈下，剑势如奔雷滚滚。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 大奔，道果圆满阶段·剑之尊者。形象：棕毛大熊，天生神力。 核心意象：奔雷剑、蛮力、七侠豪杰。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，奔雷剑劈下，剑势如奔雷滚滚。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 跳跳（`tiaotiao`）

**灵胎初醒 · 青光剑意**
- 跳跳，灵胎初醒阶段·青光剑意。初始形态：一缕青光剑意凝成的种子，剑光快如飞燕掠影，绿衣小猴的虚影在风中跳跃。风属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 绿衣小猴**
- 跳跳，凡尘砺心阶段·绿衣小猴。形象：绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 青光剑出**
- 跳跳，道法初成阶段·青光剑出。形象：绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，青光剑出鞘，剑光一闪快若闪电。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 风驰剑动**
- 跳跳，大劫淬炼阶段·风驰剑动。形象：绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，青光剑出鞘，剑光一闪快若闪电。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 灵猴剑圣**
- 跳跳，封神登天阶段·灵猴剑圣。形象：绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，青光剑出鞘，剑光一闪快若闪电。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 跳跳，道果圆满阶段·剑之尊者。形象：绿毛小猴，身轻如燕。 核心意象：青光剑、灵巧身法、七侠奇才。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，青光剑出鞘，剑光一闪快若闪电。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 莎丽（`shali`）

**灵胎初醒 · 紫云剑意**
- 莎丽，灵胎初醒阶段·紫云剑意。初始形态：一缕紫云剑意凝成的种子，紫云翻涌如雾，紫衣剑客的温柔轮廓在其中隐现。木属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 紫衣少年**
- 莎丽，凡尘砺心阶段·紫衣少年。形象：紫毛小猫，温柔坚韧。 核心意象：紫云剑、温柔剑意、七侠巾帼。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 紫云剑出**
- 莎丽，道法初成阶段·紫云剑出。形象：紫毛小猫，温柔坚韧。 核心意象：紫云剑、温柔剑意、七侠巾帼。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，紫云剑一舞，剑影如紫云翻涌。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 云剑逍遥**
- 莎丽，大劫淬炼阶段·云剑逍遥。形象：紫毛小猫，温柔坚韧。 核心意象：紫云剑、温柔剑意、七侠巾帼。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，紫云剑一舞，剑影如紫云翻涌。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 紫霞剑圣**
- 莎丽，封神登天阶段·紫霞剑圣。形象：紫毛小猫，温柔坚韧。 核心意象：紫云剑、温柔剑意、七侠巾帼。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，紫云剑一舞，剑影如紫云翻涌。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 莎丽，道果圆满阶段·剑之尊者。形象：紫毛小猫，温柔坚韧。 核心意象：紫云剑、温柔剑意、七侠巾帼。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，紫云剑一舞，剑影如紫云翻涌。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 达达（`dada`）

**灵胎初醒 · 旋风剑意**
- 达达，灵胎初醒阶段·旋风剑意。初始形态：一缕旋风剑意凝成的种子，剑气旋成环舞，竹林熊猫的沉稳虚影在风中端坐。风属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 竹林熊猫**
- 达达，凡尘砺心阶段·竹林熊猫。形象：绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 旋风剑出**
- 达达，道法初成阶段·旋风剑出。形象：绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，旋风剑轻转，剑势化柔为刚四两拨千斤。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 竹林剑舞**
- 达达，大劫淬炼阶段·竹林剑舞。形象：绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，旋风剑轻转，剑势化柔为刚四两拨千斤。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 熊猫剑圣**
- 达达，封神登天阶段·熊猫剑圣。形象：绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，旋风剑轻转，剑势化柔为刚四两拨千斤。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 剑之尊者**
- 达达，道果圆满阶段·剑之尊者。形象：绿面熊猫，沉稳大气。 核心意象：旋风剑、沉稳气度、七侠长者。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，旋风剑轻转，剑势化柔为刚四两拨千斤。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 黑心虎（`heixinhu`）

**灵胎初醒 · 暗影虎意**
- 黑心虎，灵胎初醒阶段·暗影虎意。初始形态：一缕暗影虎意凝成的种子，黑风裹着威压，猛虎如炬的双目在黑暗中睁开。暗属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 小老虎**
- 黑心虎，凡尘砺心阶段·小老虎。形象：黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 暗影虎**
- 黑心虎，道法初成阶段·暗影虎。形象：黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，虎爪一探，黑风裹挟着威压扑来。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 幽影虎王**
- 黑心虎，大劫淬炼阶段·幽影虎王。形象：黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，虎爪一探，黑风裹挟着威压扑来。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 霸业之虎**
- 黑心虎，封神登天阶段·霸业之虎。形象：黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，虎爪一探，黑风裹挟着威压扑来。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 暗黑尊者**
- 黑心虎，道果圆满阶段·暗黑尊者。形象：黑毛猛虎，双目如炬。 核心意象：黑虎纹、威压、亦正亦邪。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，虎爪一探，黑风裹挟着威压扑来。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

#### 黑小虎（`heixiaohu`）

**灵胎初醒 · 暗影虎卵**
- 黑小虎，灵胎初醒阶段·暗影虎卵。初始形态：一枚暗影虎卵，墨黑灵光裹着虎影，利爪的锋芒在执念般的暗光中初现。暗属性灵光微微环绕。神态：剑意初凝，稚气未脱的认真。动作：剑意化种，静候萌发。衣着：剑意虚影，素衣未备。梳造：发丝初束，剑穗微晃。意境：剑意种子萌芽，江湖初闻，少年侠气的微光。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑意青（#A8D8EA）主调 + 银灰（#E0E0E0）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; sword spirit; 剑意初凝，稚气未脱的认真; palette #A8D8EA with #E0E0E0 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**凡尘砺心 · 黑毛小虎**
- 黑小虎，凡尘砺心阶段·黑毛小虎。形象：墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。神态：初握剑的紧张与欢喜。动作：笨拙挥剑，破绽百出却认真。衣着：素色短打，布带缠手。梳造：束发利落，额前碎发。意境：习武入门，剑招初练，磨砺锋芒。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：侠义蓝（#4FC3F7）主调 + 热血红（#FF8A80）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; martial training; 初握剑的紧张与欢喜; palette #4FC3F7 with #FF8A80 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道法初成 · 黑风小虎**
- 黑小虎，道法初成阶段·黑风小虎。形象：墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。神态：剑气初成，眸光锐利。动作：剑意初现，剑光如虹，利爪如刃划下，黑风裂空。衣着：练功劲装，护腕已备。梳造：发带束起，英气渐显。意境：剑法初成，行走江湖，惩恶扬善。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：江湖橙（#FF7043）主调 + 侠气蓝（#29B6F6）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; journeying hero; 剑气初成，眸光锐利; palette #FF7043 with #29B6F6 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**大劫淬炼 · 黑小虎·剑**
- 黑小虎，大劫淬炼阶段·黑小虎·剑。形象：墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。神态：浴血后的沉静，杀气内敛。动作：全力挥剑，剑影纵横，利爪如刃划下，黑风裂空。衣着：侠客劲装，血迹未干。梳造：发乱而不颓，眼神如刃。意境：七剑合璧，魔教大战，刀光剑影中淬炼。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：魔教紫（#4A148C）主调 + 血战红（#E53935）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; seven swords battle; 浴血后的沉静，杀气内敛; palette #4A148C with #E53935 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**封神登天 · 黑小虎·狂**
- 黑小虎，封神登天阶段·黑小虎·狂。形象：墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。神态：目光如剑，豪气干云。动作：剑法大成，一剑定乾坤，利爪如刃划下，黑风裂空。衣着：锦袍披风，名剑在腰。梳造：侠客冠束，风采卓然。意境：七侠之名扬天下，快意恩仇，豪情万丈。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：七剑金（#FFD700）主调 + 侠名白（#FFFFFF）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; legendary hero; 目光如剑，豪气干云; palette #FFD700 with #FFFFFF accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text

**道果圆满 · 幽暗少主**
- 黑小虎，道果圆满阶段·幽暗少主。形象：墨黑猛虎，爪牙如刃。 核心意象：虎爪利刃、执念、黑虎血脉。神态：剑道圆满，锋芒内敛却慑人。动作：人剑合一，剑光裂天，利爪如刃划下，黑风裂空。衣着：剑者长袍，剑气绕身。梳造：高冠束发，剑眉入鬓。意境：剑道圆满，一剑开天，守护苍生的宗师气象。风格：2D 国漫武侠动画风，干净线条，水墨武侠背景，快意恩仇。色彩：剑道白（#F5F0E8）主调 + 圆满蓝（#90CAF9）点缀。构图：800×1000 竖版，居中全身像，正面 3/4 视角，角色主体占画面约 60%。画质：高细节，柔和渐变光，背景干净，无文字、无水印。
- EN：2D Chinese wuxia animation, clean linework, ink-wash martial arts backdrop; master of the blade; 剑道圆满，锋芒内敛却慑人; palette #F5F0E8 with #90CAF9 accents; centered full-body 3/4 view, 800x1000 vertical, soft gradient light, high detail, no text
