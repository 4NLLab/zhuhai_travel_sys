# Phase 4 命令记录

日期：2026-06-26

## 服务确认

- 原版静态页：`http://127.0.0.1:8000/`
- 迁移版 H5：`http://127.0.0.1:5174/`
- 迁移版启动命令：`VITE_API_MODE=mock corepack pnpm -C miniprogram dev:h5 -- --host 127.0.0.1`

## 截图命令

使用本机 Playwright Chromium，在 `390x844`、`519x927` 两个视口打开以下页面并保存截图：

- `http://127.0.0.1:5174/#/pages/island-cruise/index`
- `http://127.0.0.1:5174/#/pages/island-cruise/index?step=traveler`
- `http://127.0.0.1:5174/#/subpackages/driver/pages/home/index`

截图输出：

- `screenshots/migrated-island-detail-after-390x844.png`
- `screenshots/migrated-island-detail-after-519x927.png`
- `screenshots/migrated-island-booking-after-390x844.png`
- `screenshots/migrated-island-booking-after-519x927.png`
- `screenshots/migrated-driver-after-390x844.png`
- `screenshots/migrated-driver-after-519x927.png`

## 自动化验证

| 命令 | 结果 |
|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 |
| `corepack pnpm -C miniprogram test` | 通过，7 个测试文件、20 个测试用例通过 |
| `corepack pnpm -C miniprogram build:h5` | 通过 |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 |
| `corepack pnpm -C miniprogram check:size` | 通过，主包 `1815.7 KB`、司机分包 `16.3 KB`、总包 `1832.0 KB` |

## 构建备注

- `mp-weixin` 构建产物位于 `miniprogram/dist/build/mp-weixin`。
- Phase 4 未新增大图资产，复用 Phase 2 已入包的夜游和司机背景图。
