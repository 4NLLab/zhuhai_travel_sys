# Phase 1 命令记录

采集时间：2026-06-26

| 命令 | 结果 | 关键输出 |
|---|---|---|
| `node -v` | 通过 | `v22.22.3` |
| `npm -v` | 通过 | `10.9.8` |
| `pnpm -v` | 未通过 | 全局 `pnpm` 不存在；改用 `corepack pnpm` |
| `corepack pnpm -C miniprogram install` | 通过 | 生成 `pnpm-lock.yaml`；DCloud 依赖有 peer warning，已 pin 到 Vite 5.2.8 / Vue 3.4.21 |
| `corepack pnpm -C miniprogram typecheck` | 通过 | `vue-tsc --noEmit` |
| `corepack pnpm -C miniprogram test` | 通过 | 3 个测试文件、5 个测试通过 |
| `corepack pnpm -C miniprogram build:h5` | 通过 | `DONE Build complete.` |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | `DONE Build complete.`；产物在 `dist/build/mp-weixin` |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | 扫描 `src/pages`、`src/components`、`src/api`、`src/adapters`、`src/stores`、`src/utils`、`dist/build/mp-weixin` |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 78.7 KB，司机分包 0.9 KB，总包 79.6 KB |
| `python3 /home/zhouq/.codex/skills/webapp-testing/scripts/with_server.py --help` | 通过 | helper 可用；`python` 命令不存在 |
| `npx playwright install chromium` | 部分通过 | Chromium 下载成功，但 CLI 启动缺 `libnspr4.so` |
| `npx playwright install-deps chromium` | 未通过 | 需要 sudo 交互认证 |
| `npx playwright install firefox` | 通过 | Firefox 下载成功，并用于截图 |
| `npx playwright screenshot --browser=firefox ...` | 通过 | 生成首页、订单、我的桌面和窄屏截图 |

## 说明

- H5 dev server 使用 `corepack pnpm -C miniprogram dev:h5 --host 127.0.0.1 --port 5173` 启动，截图完成后已关闭。
- Playwright MCP 可渲染页面并读取文本；仓库内截图最终由 Playwright Firefox CLI 生成。
