# Phase 3 命令记录

采集时间：2026-06-26

| 命令 | 结果 | 关键输出 |
|---|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 | `vue-tsc --noEmit` |
| `corepack pnpm -C miniprogram test` | 通过 | 5 个测试文件、13 个测试通过 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | `Platform lint passed.` |
| `corepack pnpm -C miniprogram build:h5` | 通过 | `DONE Build complete.` |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | `DONE Build complete.`；产物在 `dist/build/mp-weixin` |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 397.1 KB，司机分包 0.9 KB，总包 398.1 KB |
| `npx playwright screenshot --browser=firefox ...` | 通过 | 生成环岛游详情、出行人、支付、票券和无班次状态截图 |

## 说明

- H5 dev server 使用 `corepack pnpm -C miniprogram dev:h5 --host 127.0.0.1 --port 5173` 启动，截图完成后已关闭。
- 本阶段继续使用 Firefox CLI 截图；Chromium 系统依赖 blocker 沿用 Phase 1/2 记录。
