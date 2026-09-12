# 宠物图片 AI 生成 · 统一约束与交付规范

> 目的：让开发者用 AI（Gemini / 即梦 / Midjourney 等）批量生成宠物各阶段形象，产出图可直接替换进系统，无需二次处理。
> 提示词来源：`docs/stage-ai-prompts/{画风}.md`（由 `frontend-vue/scripts/generate-stage-prompts.mjs` 自动生成，共 756 条）。

## 一、统一约束（每张图必须满足）

| 约束项 | 要求 |
|--------|------|
| 画布 | **1024×1024 正方形**（固定，不得变化） |
| 格式 | **PNG**，透明背景（模型不支持时用纯白背景，交付前抠图转透明） |
| 构图 | 角色单主体居中，正面 3/4 视角，占画面 **60-70%**，全身完整不出画框 |
| 禁止 | 不得出现文字、水印、边框、背景杂物、多角色 |
| 一致性 | **同一角色的 6 个阶段 = 同一画风 + 同一角色特征递进**（配色/标志元素逐阶段延续） |
| 命名 | `{speciesId}-{stage}.png`，stage ∈ `L1 / L3 / L5 / L7 / L9 / L12`（如 `zhulong-L9.png`） |
| 存放 | `frontend-vue/public/pets/{画风}/{系列Id}/{speciesId}-{stage}.png` |

## 二、系列 → 画风绑定（已定，生成时不得混用）

| 系列 | 画风 | 提示词文件 |
|------|------|-----------|
| 山海经 / 东方神话 | 中国风水墨 | `docs/stage-ai-prompts/inkcn.md`（228 条） |
| 宝可梦 / 数码宝贝 / 虹猫蓝兔七侠传 | 动漫 | `docs/stage-ai-prompts/anime.md`（168 条） |
| 国宝 | 卡通 | `docs/stage-ai-prompts/cartoon.md`（72 条） |
| 魔法奇幻 / 史前生物 | 3D 渲染 CG | `docs/stage-ai-prompts/cg3d.md`（144 条） |
| 星座守护 | 手办 | `docs/stage-ai-prompts/figure.md`（72 条） |
| 传统节日 | 水彩画 | `docs/stage-ai-prompts/watercolor.md`（72 条） |

> 画风按系列固定，**同一角色的六个阶段严格同一画风**——保证角色随阶段成长的视觉连续性。

## 三、使用流程

1. 打开对应画风的提示词文件（如 `anime.md`），按物种分节，每节 6 条（L1/L3/L5/L7/L9/L12）。
2. 每条提示词已拼好：画风行 + 角色/阶段 + 叙事关键词 + 神态/动作/服饰/能力 + 画面描述 + 统一约束。直接整条复制投喂即可。
3. **一致性技巧**：同一物种先生成 L1，之后每阶把上一阶成图作为参考图附给模型，描述只改增量（体型变大 / 新增部件 / 光效升级）。
4. 生成后按「命名」规则存放到对应目录。

## 四、人工精写与覆盖

- 档案数据 `petLifeStories.ts` 的每阶段已支持两个可选字段：
  - `interaction`：本阶段与其他角色的互动/冲突（剧情推演展示用，不进提示词）
  - `aiPrompt`：人工精写的提示词——**填写后优先生效**（自动组装的作为兜底）
- 修改数据源后重跑 `node scripts/generate-stage-prompts.mjs` 即可刷新全部提示词文件。

## 五、图片接入（开发侧）

1. 图片放入 `frontend-vue/public/pets/{画风}/{系列Id}/` 后，在 `docs/pet-image-manifest.md` 登记。
2. `frontend-vue/src/utils/petImage.ts` 的 `PET_IMAGE_MODE` 从 `'emoji'` 改为 `'image'`，`getPetImageUrl` 恢复按 MANIFEST 取图。
3. 全站（课堂卡片 / 图鉴 / 详情 / 进化树）自动切换为生成图。
