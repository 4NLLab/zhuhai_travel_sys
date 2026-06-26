# Phase 3 执行计划对抗式评审

日期：2026-06-26

## Round 1

| 角色 | 发现 | 处理 |
|---|---|---|
| Builder reviewer | 三页均已有 Mock 数据，可直接重排模板和样式，不需要改 adapter。 | 保持数据层不动，只做页面视觉结构。 |
| Skeptic reviewer | 票券页差异最大，若只沿用 `TicketSummaryCard` 仍会缺大二维码和底部操作栏。 | 票券页改为页面内专用布局，不复用现有摘要卡。 |
| Verifier reviewer | 我的/订单原版是 `index.html` 内部状态，不是独立 URL；截图索引必须引用 Phase 1 对应状态截图。 | Phase 3 `screenshots.md` 明确引用 `original-profile-*` 和 `original-orders-*`。 |

结论：无 P0/P1 计划缺陷，进入第二轮。

## Round 2

| 角色 | 发现 | 处理 |
|---|---|---|
| Builder reviewer | H5 原生 tabBar 已在 Phase 2 隐藏，基础页需要自绘底部导航，否则我的/订单缺导航。 | 我的、订单页补自绘底部导航；票券页按 `ticket.html` 使用底部操作栏。 |
| Skeptic reviewer | 为了票券页大图复制 `verify-hero` 会增加包体，需要复核 2 MB 主包预算。 | 只复制一张 `verify-hero-wide-clean-web.png`，Phase 3 后运行 `check:size`。 |
| Verifier reviewer | 图标细节可有小差，但二维码、状态、筛选、订单卡和底部操作必须截图可见。 | 修复后逐页截图并人工判定。 |

结论：第二轮后未发现 P0/P1 计划缺陷，可以实施。
