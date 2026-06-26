# Phase 4 执行级计划：环岛游订票与司机端视觉还原

日期：2026-06-26

## 目标与边界

- 目标：以 `island-cruise.html`、`island-cruise-booking.html`、`driver.html` 的 Phase 1 原版截图为真相源，还原环岛游详情/订票和司机端核心状态视觉。
- 迁移版入口：
  - `http://127.0.0.1:5174/#/pages/island-cruise/index`
  - `http://127.0.0.1:5174/#/pages/island-cruise/index?step=traveler`
  - `http://127.0.0.1:5174/#/subpackages/driver/pages/home/index`
- 视口：`390x844`、`519x927`
- 不修改后端、旧静态 HTML、数据库或 Docker。

## 修复范围

- `miniprogram/src/pages/island-cruise/index.vue`
- `miniprogram/src/subpackages/driver/pages/home/index.vue`
- `miniprogram/docs/validation/visual-parity/phase-4/`

## 差异关闭目标

- 环岛游详情：恢复深色夜游背景、媒体 hero、顶部 chip、深色信息卡、金色咨询/购买卡和固定底部栏。
- 环岛游订票：恢复深色 sticky 顶栏、视频卡、hero copy、三列指标、近期推荐、班次查询和固定价格栏。
- 司机端：默认回到司机登录/注册首屏，恢复车内背景 hero、返回首页/司机入驻 chip、三步入驻卡、双 tab 和白色表单卡。
- 工作台、钱包、佣金、提现和推广码继续可通过登录按钮进入，保持 Phase 4 业务状态可用。

## 验证命令

- 修复后截图：环岛游详情、环岛游订票、司机端各 `390x844`、`519x927`。
- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 回滚边界

- 回滚本 Phase 只需还原两个页面和 Phase 4 证据目录。
