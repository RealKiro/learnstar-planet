# 宠物进化体系 · 5 阶段 → 6 阶段改造 TODO

> 创建：2026-08-07 · 状态：**✅ 已完成**
> 新 6 阶段：`egg`(Lv1-2) / `baby`(Lv3-4) / `growing`(Lv5-6) / `mature`(Lv7-8) / `legendary`(Lv9-10) / `transcendent`(Lv11-12，归真级)

## ✅ 已完成
- [x] `petData.ts`：`PetLevel.stage` 类型加 `'transcendent'`；`getLevelStage` 改 6 桶
- [x] `petHandbookData.ts`：`stageIndexForLevel` 改 6 桶（对齐诗文/台词 6 槽）
- [x] 全量等级数据 stage 字段重映射（809 个，petData.ts + petDataExtended.ts，脚本验证 12 级每级 1 种 stage）
- [x] `petArtParts.ts` computePetArt：`scaleByStage`(transcendent 1.55) / `glowByStage` / `eyesByStage` / `stageShade` / `stageAllowed`(transcendent 全部件)；`stageVariants`/`stageMain`/`PetArtRender.stage` 类型加 transcendent
- [x] `Pet.php` `currentStage()`（6 桶 + 归真级）+ `currentStageEmoji()`（6 桶）
- [x] **UI 阶段标签/徽章**：`PetEvolutionTree.vue` stageLabels 加 `transcendent: '归真级'`；`PetHandbook.vue` 徽章 label 加 `transcendent: '归真'` + `.stage--transcendent` 样式；`PetDisplay.vue` 加 `.level--transcendent` 样式
- [x] **类型**：`types/index.ts` PetDetail.stage 加 'transcendent'
- [x] **文档**：`docs/宠物系统全案.md` 1.2 节"5 阶段"与 3.1 节"12 级 → 6 阶段"表更新
- [x] **验证**：`npm run build` 全绿
- [x] **提交**：见 git log
