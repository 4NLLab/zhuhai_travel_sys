# Phase 1 执行级计划：原版视觉基线与逐页差异清单

日期：2026-06-26

## 目标与边界

- 目标：为原静态 HTML 与当前 uni-app H5 建立同视口截图基线，输出逐页视觉差异清单。
- 真相源：`http://127.0.0.1:8000/` 下的原静态 HTML。
- 迁移版：`http://127.0.0.1:5174/` 下的 uni-app H5，默认 `VITE_API_MODE=mock`。
- 本 Phase 不做大规模 UI 修复；仅允许新增视觉验收证据、截图记录、执行计划和必要的截图辅助资料。

## 页面映射

| 范围 | 原版 URL / 状态 | 迁移版 URL / 状态 | 视口 | 截图文件 |
|---|---|---|---|---|
| 首页 | `/index.html` 默认首页 | `/#/pages/home/index` | `390x844`, `519x927` | `original-home-*`, `migrated-home-*` |
| 我的 | `/index.html` 点击底部“我的” | `/#/pages/profile/index` | `390x844`, `519x927` | `original-profile-*`, `migrated-profile-*` |
| 订单 | `/index.html` 从“我的订单”进入订单状态 | `/#/pages/orders/index` | `390x844`, `519x927` | `original-orders-*`, `migrated-orders-*` |
| 票券 | `/ticket.html` 默认票券 | `/#/pages/ticket/index` | `390x844`, `519x927` | `original-ticket-*`, `migrated-ticket-*` |
| 环岛游详情 | `/island-cruise.html` 默认详情 | `/#/pages/island-cruise/index?step=detail` | `390x844`, `519x927` | `original-island-detail-*`, `migrated-island-detail-*` |
| 环岛游订票 | `/island-cruise-booking.html` 默认详情 | `/#/pages/island-cruise/index?step=detail` | `390x844`, `519x927` | `original-island-booking-*`, `migrated-island-booking-*` |
| 司机端 | `/driver.html` 默认登录/注册页 | `/#/subpackages/driver/pages/home/index` | `390x844`, `519x927` | `original-driver-*`, `migrated-driver-*` |

## 截图与验收方法

- 等待页面 `networkidle` 后再额外等待 1000ms，降低异步数据和首屏动画造成的波动。
- 首屏截图使用相同 viewport、相同 scrollY=0、相同 mock 场景。
- 原版和迁移版截图均保存到 `miniprogram/docs/validation/visual-parity/phase-1/screenshots/`。
- Phase 1 验收结论只能用于建立基线；如截图显示明显不同，结论必须记录为“不通过”，后续 Phase 继续修复。

## 预计修改文件

- `miniprogram/docs/validation/visual-parity/phase-1/phase-plan.md`
- `miniprogram/docs/validation/visual-parity/phase-1/plan-reviews.md`
- `miniprogram/docs/validation/visual-parity/phase-1/commands.md`
- `miniprogram/docs/validation/visual-parity/phase-1/differences.md`
- `miniprogram/docs/validation/visual-parity/phase-1/visual-standard.md`
- `miniprogram/docs/validation/visual-parity/phase-1/screenshots/*`

## 测试命令

- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 回滚边界

- Phase 1 只新增视觉验收证据和文档，回滚时可整体删除 `miniprogram/docs/validation/visual-parity/phase-1/`。
- 不修改原静态 HTML、uni-app 页面实现、后端代码、Docker 或数据库文件。
