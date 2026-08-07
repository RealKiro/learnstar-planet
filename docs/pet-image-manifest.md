# 宠物 AI 图片登记表（Pet Image Manifest）

> 与 `frontend-vue/src/utils/petImage.ts` 的 `MANIFEST` 数组同步：**生成图片后**把物种 id 登记进该数组，`PetSprite` 会自动优先加载图片、未登记则回退程序化 SVG。

## 生成流程
1. 用 `frontend-vue/scripts/generate-pet-prompts.js` 生成全部物种 × 6 阶段的提示词 CSV
2. 用任意 AI 生图工具（Midjourney / DALL·E / Stable Diffusion / 即梦 等）按提示词生成图片
3. 图片放入 `frontend-vue/public/pets/{seriesId}/{speciesId}-{stage}.webp`
4. 在 `utils/petImage.ts` 的 `MANIFEST` 登记 speciesId，并把 `public/pets/` 拷贝到 `backend/public/pets/`（部署时）

## 已登记物种（MANIFEST）

| 物种 id | 系列 | 状态 |
|---------|------|------|
| （空） | — | 尚未生成任何图片 |

## 图片规格
- 目录：`public/pets/{seriesId}/{speciesId}-{stage}.webp`
- 阶段：egg / baby / growing / mature / legendary / transcendent（6 阶段）
- 画布：800×1000 竖版；角色主体占 60% 居中偏下；六阶视觉基调见 `宠物系统全案.md` 6.4 节

## 推荐生图顺序（先验收风格）
1. 东方神话 20 角色（孙悟空 / 哪吒 / 杨戬 / 八仙 / 太上老君 / 钟馗 等）
2. 星座圣斗士 12（白羊座·穆 等）
3. 宝可梦御三家 + 虹猫蓝兔 9
4. 其余系列全量
