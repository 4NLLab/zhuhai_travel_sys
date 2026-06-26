# Phase 2 执行计划对抗式评审

日期：2026-06-26

## Round 1

| 角色 | 发现 | 处理 |
|---|---|---|
| Builder reviewer | 计划聚焦首页 P0：hero 资产、顶部叠加结构和热门目的地，是最大视觉偏差来源。 | 保留首页单页修复范围，不扩散到其他页面。 |
| Skeptic reviewer | 只替换图片不足以通过，原版首屏还包含搜索、问候语、轮播文案、quick-sheet 负 margin 和底部导航。 | 将首屏结构、quick-sheet、promo 和 tabBar 纳入修复目标。 |
| Verifier reviewer | Phase 2 必须保存修复前和修复后截图；Phase 1 已有修复前，可在 Phase 2 引用并重新采集修复后。 | Phase 2 `screenshots.md` 必须列出 Phase 1 修复前截图和 Phase 2 修复后截图。 |

结论：无 P0/P1 计划缺陷，进入第二轮确认兼容风险。

## Round 2

| 角色 | 发现 | 处理 |
|---|---|---|
| Builder reviewer | 复制原版首张 hero 图即可对齐当前截图首屏，避免引入三张大图造成包体压力。 | 只复制 `wechat_2026-06-20_150339_912.png` 为本地首页 hero。 |
| Skeptic reviewer | `navigationStyle: custom` 会影响 H5/小程序顶部安全区，需要首页内部自行留白和布局。 | 首页 hero 顶部使用原版 top-bar 结构并保留安全区 padding。 |
| Verifier reviewer | 底部 tabBar 原生图标限制可能留下 P1 差异，修复后必须截图判定；不允许只靠代码判断。 | 修复后重新采集 `390x844` 和 `519x927` 截图，若仍明显不同继续迭代。 |

结论：第二轮后未发现 P0/P1 计划缺陷，可以实施。
