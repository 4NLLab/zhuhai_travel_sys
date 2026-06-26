# Phase 3 执行级计划：我的、订单、票券基础页视觉还原

日期：2026-06-26

## 目标与边界

- 目标：以 Phase 1 的原版截图为真相源，还原我的、订单、票券详情三个基础页面视觉。
- 原版来源：`index.html` 内 mine/orders 状态、`ticket.html` 默认票券页。
- 迁移版入口：
  - `http://127.0.0.1:5174/#/pages/profile/index`
  - `http://127.0.0.1:5174/#/pages/orders/index`
  - `http://127.0.0.1:5174/#/pages/ticket/index`
- 视口：`390x844`、`519x927`
- 不修改后端、旧静态 HTML、环岛游、司机端。

## 修复范围

- `miniprogram/src/pages/profile/index.vue`
- `miniprogram/src/pages/orders/index.vue`
- `miniprogram/src/pages/ticket/index.vue`
- `miniprogram/src/static/phase3/verify-hero-wide-clean-web.png`
- `miniprogram/docs/validation/visual-parity/phase-3/`

## 差异关闭目标

- 我的页：恢复蓝绿渐变页头、资料卡、下一段行程、票券 KPI、我的订单入口、常用服务宫格和底部导航 active=我的。
- 订单页：恢复返回我的页头、状态筛选、紧凑订单卡、状态标签和详情按钮。
- 票券页：恢复顶部导航、待使用状态、大二维码核销卡、票券信息、使用须知、订单信息和底部操作栏。

## 验证命令

- 修复后截图：三页各 `390x844`、`519x927`。
- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 回滚边界

- 回滚本 Phase 只需还原三个页面、Phase 3 静态资产和 Phase 3 证据目录。
