# Phase 2 执行级计划：首页视觉还原

日期：2026-06-26

## 目标与边界

- 目标：以 Phase 1 的 `original-home-390x844.png`、`original-home-519x927.png` 为真相源，修复 uni-app 首页首屏和主体视觉差异。
- 原版 URL：`http://127.0.0.1:8000/index.html`
- 迁移版 URL：`http://127.0.0.1:5174/#/pages/home/index`
- 视口：`390x844`、`519x927`
- 不修改后端、旧静态 HTML、订单/我的/票券/环岛游/司机端页面。

## 修复范围

- `miniprogram/src/pages/home/index.vue`
- `miniprogram/src/pages.json`
- `miniprogram/src/static/phase2/home-hero-watercolor.png`
- 必要时补充 `miniprogram/docs/validation/visual-parity/phase-2/` 证据文档和截图。

## 差异关闭目标

- P0：首页 hero 从夜景邮轮卡片改为原版水彩珠海湾区插画。
- P0：去除 H5 默认白色导航栏，恢复品牌、搜索、问候语、轮播文案叠加在 hero 内。
- P1：热门目的地恢复三列浅色彩块卡片，只显示前 6 项。
- P1：恢复首屏下缘露出的“九洲港至蛇口船票”促销卡和“立即购票”按钮。
- P1：商品分类 tabs、商品卡、价格按钮贴近原版视觉语言。
- P1：底部 tabBar 文案和色彩尽量贴近原版“首页 / 客服 / 我的”。

## 验证命令

- 采集修复后首页截图：`390x844`、`519x927`
- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 回滚边界

- 回滚本 Phase 只需要还原首页 Vue、pages.json、复制进 `src/static/phase2/` 的首页 hero 资产和 Phase 2 证据目录。
