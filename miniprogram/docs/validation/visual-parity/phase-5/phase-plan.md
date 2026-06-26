# Phase 5 执行级计划：最终视觉验收、文档与交付收口

日期：2026-06-26

## 目标与边界

- 目标：重新采集所有原版与迁移版最终截图，建立最终对比矩阵，并运行最终全量质量命令。
- 原版端口：`http://127.0.0.1:8000/`
- 迁移版端口：`http://127.0.0.1:5174/`
- 页面范围：首页、我的、订单、票券、环岛游详情、环岛游订票、司机端。
- 不再改动旧静态 HTML、后端、数据库或 Docker。

## 截图范围

- 原版：
  - `index.html` 首页
  - `index.html` 我的状态
  - `index.html` 订单状态
  - `ticket.html`
  - `island-cruise.html`
  - `island-cruise-booking.html`
  - `driver.html`
- 迁移版：
  - `#/pages/home/index`
  - `#/pages/profile/index`
  - `#/pages/orders/index`
  - `#/pages/ticket/index`
  - `#/pages/island-cruise/index`
  - `#/pages/island-cruise/index?step=traveler`
  - `#/subpackages/driver/pages/home/index`

## 验证命令

- 最终截图：每页 `390x844`、`519x927`。
- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 输出文档

- `miniprogram/docs/validation/visual-parity/phase-5/screenshots.md`
- `miniprogram/docs/validation/visual-parity/phase-5/final-matrix.md`
- `miniprogram/docs/validation/visual-parity/final-report.md`
