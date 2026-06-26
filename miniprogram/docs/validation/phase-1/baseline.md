# Phase 1 基线证据

采集时间：2026-06-26

## 仓库前端基线

- 根目录 HTML 页面共 12 个：`admin.html`、`car.html`、`driver.html`、`flow.html`、`hongkong.html`、`hotel-list.html`、`hotel.html`、`index.html`、`island-cruise-booking.html`、`island-cruise.html`、`macau.html`、`ticket.html`。
- Phase 1 执行前不存在 `miniprogram/` 工程文件。
- 根目录未发现既有 `package.json`、`pnpm-lock.yaml`、`pages.json`、`manifest.json`。

## 资源基线

- `assets/`：69 个文件，约 32M。
- 超过 1M 的图片包括 `zhuhai-bay-home-hero.png`、`macau-cruise-night-banner.png`、`taxi-scan-illustration.png`、`realname-illustration.png` 等。
- Phase 1 不复制 `assets/` 到小程序包。

## 后端与代理基线

- 后端入口：`backend/routes/router.go`。
- 架构参考：`backend/docs/ARCHITECTURE.md`，其 API 清单可能落后于实际路由。
- Docker nginx：`docker/nginx/default.conf` 将 `/api/` 代理到 `backend:8080`。
- 小程序端不能依赖浏览器相对路径代理，DevTools/local/真机 base URL 记录在接口矩阵。

## 本机工具

- Node：`v22.22.3`
- npm：`10.9.8`
- 全局 `pnpm`：未安装；使用 `corepack pnpm` 或 packageManager 固定的 pnpm。
