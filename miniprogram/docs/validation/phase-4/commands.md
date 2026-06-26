# Phase 4 命令记录

采集时间：2026-06-26

| 命令 | 结果 | 关键输出 |
|---|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 | `vue-tsc --noEmit` |
| `corepack pnpm -C miniprogram test` | 通过 | 6 个测试文件、17 个测试通过 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | `Platform lint passed.` |
| `corepack pnpm -C miniprogram build:h5` | 通过 | `DONE Build complete.` |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | `DONE Build complete.`；产物在 `dist/build/mp-weixin` |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 401.5 KB，司机分包 14.4 KB，总包 415.8 KB |
| `npx playwright screenshot --browser=firefox ...` | 通过 | 生成司机 active、待审核、空钱包、余额不足状态截图 |

## 说明

- H5 dev server 使用 `corepack pnpm -C miniprogram dev:h5 --host 127.0.0.1 --port 5173` 启动，截图完成后已关闭。
- active 与余额不足移动截图使用 `--full-page`，覆盖钱包、佣金、提现申请和提现记录。
