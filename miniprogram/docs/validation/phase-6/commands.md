# Phase 6 命令记录

日期：2026-06-26

| 命令 | 结果 | 说明 |
|---|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 | Vue/TS 类型检查 |
| `corepack pnpm -C miniprogram test` | 通过 | 7 个测试文件，20 个测试 |
| `corepack pnpm -C miniprogram build:h5` | 通过 | H5 构建 |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | 微信小程序构建 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | 禁用平台 API 扫描 |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 406.9 KB，司机分包 14.6 KB，总包 421.5 KB |
| `cd backend && go test ./...` | 通过 | 确认后端未被破坏 |
| `npx playwright screenshot --browser=firefox ...` | 通过 | Phase 6 H5 截图矩阵已生成 |
