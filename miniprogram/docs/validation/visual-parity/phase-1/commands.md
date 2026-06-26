# Phase 1 命令记录

日期：2026-06-26

| 命令 | 结果 | 说明 |
|---|---|---|
| `python3 -m http.server 8000 --bind 127.0.0.1` | 未启动新进程 | `8000` 已被现有服务占用。 |
| `curl -I --max-time 5 http://127.0.0.1:8000/index.html` | 通过 | 原版静态页返回 `200 OK`。 |
| `VITE_API_MODE=mock corepack pnpm -C miniprogram dev:h5 -- --host 127.0.0.1 --port 5174` | 运行中 | uni-app H5 mock 服务。 |
| `curl -I --max-time 5 http://127.0.0.1:5174/#/pages/home/index` | 通过 | 迁移版 H5 返回 `200 OK`。 |
| `Playwright MCP + 本地 8765 截图接收器` | 通过 | 已采集 28 张原版/迁移版截图，保存到 `screenshots/`。 |
| `file miniprogram/docs/validation/visual-parity/phase-1/screenshots/*.png` | 通过 | 28 张截图均为 `390x844` 或 `519x927` PNG。 |
| `corepack pnpm -C miniprogram add -D playwright@1.61.1` | 通过 | 新增本地截图脚本运行依赖。 |
| `corepack pnpm -C miniprogram exec node scripts/capture-visual-parity.mjs` | 未通过 | 当前系统缺 `libnspr4.so`，Playwright headless Chromium 无法启动；截图已由 MCP 完成，脚本需安装系统依赖后复跑。 |
| `corepack pnpm -C miniprogram exec playwright install-deps chromium` | 未完成 | 需要 sudo 密码，已中止；不影响本次 MCP 截图证据。 |
| `corepack pnpm -C miniprogram typecheck` | 通过 | Vue/TypeScript 类型检查通过。 |
| `corepack pnpm -C miniprogram test` | 通过 | 7 个测试文件、20 个测试通过。 |
| `corepack pnpm -C miniprogram build:h5` | 通过 | H5 构建通过。 |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | 微信小程序构建通过。 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | 平台 API 扫描通过。 |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 406.9 KB，司机分包 14.6 KB，总包 421.5 KB。 |
