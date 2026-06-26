# Phase 2 命令记录

日期：2026-06-26

| 命令 | 结果 | 说明 |
|---|---|---|
| `VITE_API_MODE=mock corepack pnpm -C miniprogram dev:h5 -- --host 127.0.0.1` | 通过 | manifest 已固定 H5 devServer 端口为 `5174`。 |
| `LD_LIBRARY_PATH=.local/playwright-libs corepack pnpm -C miniprogram exec node ...` | 通过 | 使用本地解包 NSS/NSPR 库运行 Playwright Chromium，采集 Phase 2 首页截图。 |
| `file miniprogram/docs/validation/visual-parity/phase-2/screenshots/*.png` | 通过 | 修复后截图为 `390x844` 和 `519x927` PNG。 |
| `corepack pnpm -C miniprogram typecheck` | 通过 | Vue/TypeScript 类型检查通过。 |
| `corepack pnpm -C miniprogram test` | 通过 | 7 个测试文件、20 个测试通过。 |
| `corepack pnpm -C miniprogram build:h5` | 通过 | H5 构建通过。 |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | 微信小程序构建通过。 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | 平台 API 扫描通过。 |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 1479.9 KB，司机分包 14.6 KB，总包 1494.5 KB。 |
