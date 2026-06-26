# Phase 3 命令记录

日期：2026-06-26

## 服务确认

- 原版静态页：`http://127.0.0.1:8000/`
- 迁移版 H5：`http://127.0.0.1:5174/`
- 迁移版启动命令：`VITE_API_MODE=mock corepack pnpm -C miniprogram dev:h5 -- --host 127.0.0.1`

## 截图命令

使用本机 Playwright Chromium，在 `390x844`、`519x927` 两个视口打开以下页面并保存截图：

- `http://127.0.0.1:5174/#/pages/profile/index`
- `http://127.0.0.1:5174/#/pages/orders/index`
- `http://127.0.0.1:5174/#/pages/ticket/index`

截图输出：

- `screenshots/migrated-profile-after-390x844.png`
- `screenshots/migrated-profile-after-519x927.png`
- `screenshots/migrated-orders-after-390x844.png`
- `screenshots/migrated-orders-after-519x927.png`
- `screenshots/migrated-ticket-after-390x844.png`
- `screenshots/migrated-ticket-after-519x927.png`

## 自动化验证

| 命令 | 结果 |
|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 |
| `corepack pnpm -C miniprogram test` | 通过，7 个测试文件、20 个测试用例通过 |
| `corepack pnpm -C miniprogram build:h5` | 通过 |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 |
| `corepack pnpm -C miniprogram check:size` | 通过，主包 `1816.2 KB`、司机分包 `14.6 KB`、总包 `1830.8 KB` |

## 构建备注

- `mp-weixin` 构建产物位于 `miniprogram/dist/build/mp-weixin`。
- 最大资源仍为 `static/phase2/home-hero-watercolor.png`，大小 `854.6 KB`。
- Phase 3 新增票券横幅资源 `static/phase3/verify-hero-wide-clean-web.png`，大小 `322.0 KB`。
