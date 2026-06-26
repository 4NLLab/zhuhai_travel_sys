# Phase 2 命令记录

采集时间：2026-06-26

| 命令 | 结果 | 关键输出 |
|---|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 | `vue-tsc --noEmit` |
| `corepack pnpm -C miniprogram test` | 通过 | 4 个测试文件、9 个测试通过 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | `Platform lint passed.` |
| `corepack pnpm -C miniprogram build:h5` | 通过 | `DONE Build complete.` |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | `DONE Build complete.`；产物在 `dist/build/mp-weixin` |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 368.0 KB，司机分包 0.9 KB，总包 369.0 KB |
| `npx playwright screenshot --browser=firefox ...` | 通过 | 生成首页、我的、订单、票券详情桌面和移动截图 |

## 说明

- H5 dev server 使用 `corepack pnpm -C miniprogram dev:h5 --host 127.0.0.1 --port 5173` 启动，截图完成后已关闭。
- CSS 修复后重新执行完整验证链，以上记录对应当前最新 Phase 2 代码。
