# Phase 6 命令记录

日期：2026-06-26

| 命令 | 结果 | 说明 |
|---|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 | Vue/TS 类型检查 |
| `corepack pnpm -C miniprogram test` | 通过 | 7 个测试文件，20 个测试 |
| `corepack pnpm -C miniprogram build:h5` | 通过 | H5 构建 |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | 微信小程序构建 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | 禁用平台 API 扫描 |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 1986.7 KB，司机分包 16.3 KB，总包 2003.0 KB；主包达到 1.8 MB 预警线但低于 2 MB 门禁 |
| `cd backend && go test ./...` | 通过 | 确认后端未被破坏 |
| `npx playwright screenshot --browser=firefox ...` | 通过 | Phase 6 H5 截图矩阵已生成 |
| `docker compose config` | 通过 | Compose 展开成功，前端构建参数为 `VITE_API_MODE=local`、`VITE_API_BASE_URL=/api/v1` |
| `docker compose build frontend` | 通过 | Docker 内执行 `pnpm build:h5` 成功，镜像包含重构后的 H5 产物 |
| `docker compose up -d --no-deps --build frontend` | 通过 | 已重建并启动 `zhuhai-travel-frontend`，端口 `8000 -> 80` |
| `curl -fsS http://127.0.0.1:8000/` | 通过 | 返回 uni-app H5 入口 `index.html` |
| `curl -fsS http://127.0.0.1:8000/api/v1/health` | 通过 | 通过前端 Nginx 反代后端，返回 `{"status":"ok"}` |
| `curl -fsSI http://127.0.0.1:8000/static/phase2/home-hero-sunset-web.jpg` | 通过 | `/static/` 图片资源返回 200，并带 7 天缓存头 |
| `curl -fsSI http://127.0.0.1:8000/assets/pages-home-index.Be59ZmWA.js` | 通过 | `/assets/` 编译资源返回 200，并带 7 天缓存头 |
| Playwright Chromium `http://127.0.0.1:8000/` | 通过 | 移动与桌面视口首页渲染正常，无 warning/error；`/api/v1/products`、`/api/v1/categories` 均为 200 |
