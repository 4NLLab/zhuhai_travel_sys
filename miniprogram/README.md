# 珠海文旅小程序前端

`miniprogram/` 是新的小程序前端主工程，使用 uni-app + Vue 3 + TypeScript。首版已迁移首页、我的、订单/票券基础壳、澳门环岛游订票流程和司机端分包，并保留 Mock/local 两套数据源切换。

## 命令

```bash
corepack pnpm -C miniprogram install
corepack pnpm -C miniprogram dev:h5 -- --host 127.0.0.1 --port 5173
corepack pnpm -C miniprogram build:h5
corepack pnpm -C miniprogram build:mp-weixin
corepack pnpm -C miniprogram typecheck
corepack pnpm -C miniprogram test
corepack pnpm -C miniprogram lint:platform
corepack pnpm -C miniprogram check:size
```

## 本地预览

H5 快速预览：

```bash
VITE_API_MODE=mock VITE_MOCK_SCENARIO=phase2-success corepack pnpm -C miniprogram dev:h5 -- --host 127.0.0.1 --port 5173
```

微信小程序构建：

```bash
corepack pnpm -C miniprogram build:mp-weixin
```

构建后使用微信开发者工具导入 `miniprogram/dist/build/mp-weixin`。当前仓库环境没有微信开发者工具或 `miniprogram-ci`，因此已完成的是 build-validated；正式提测前需要补小程序运行时截图。

## 环境切换

复制 `.env.example` 为本地 `.env` 后调整：

- `VITE_API_MODE=mock`：默认，不访问后端。
- `VITE_API_MODE=local`：访问 `VITE_API_BASE_URL` 指向的本地后端。
- `VITE_MOCK_SCENARIO`：切换 Mock 场景，常用值见下表。
- `VITE_ALLOW_FALLBACK=false`：local 模式默认不静默降级。

| 切片 | 常用场景 |
|---|---|
| 主包首页/我的/订单/票券 | `phase2-success`、`phase2-empty`、`phase2-unauthorized`、`phase2-failure` |
| 环岛游 | `phase3-success`、`phase3-no-voyage`、`phase3-insufficient-stock`、`phase3-lock-expired`、`phase3-passenger-invalid`、`phase3-sale-failed` |
| 司机端 | `phase4-active`、`phase4-pending-review`、`phase4-empty-wallet`、`phase4-insufficient-balance`、`phase4-withdraw-success` |

local 模式示例：

```bash
docker compose up -d mysql backend
VITE_API_MODE=local VITE_API_BASE_URL=http://127.0.0.1:8080/api/v1 VITE_ALLOW_FALLBACK=true corepack pnpm -C miniprogram dev:h5 -- --host 127.0.0.1 --port 5173
```

Docker Compose 部署重构后的 H5 前端：

```bash
docker compose up -d --build frontend
```

该命令会使用 `docker/frontend/Dockerfile` 构建 `miniprogram/dist/build/h5`，由 Nginx 在 `http://127.0.0.1:8000` 提供访问，并将浏览器里的 `/api/v1` 请求反向代理到 Compose 内部的 `backend:8080`。

`mp-weixin` local 调试不能使用模拟器里的 `127.0.0.1` 访问宿主机后端，需要改为局域网 IP、代理域名或 HTTPS 测试域名。真机/体验版必须配置 HTTPS 合法域名。

## 目录

- `src/pages/`：主包页面：首页、订单、我的、票券详情、环岛游订票。
- `src/subpackages/driver/`：司机端分包。
- `src/api/`：API adapter 入口，页面不得直接拼 URL。
- `src/adapters/`：request、storage、logger、modal、navigation。
- `src/mock/`：Mock 场景。
- `src/utils/`：配置、字段映射和敏感信息脱敏。
- `src/types/`：页面 view model 和后端契约类型。
- `docs/contracts/`：接口契约矩阵。
- `docs/validation/`：阶段验证证据。

## 截图验收

H5 截图使用 Playwright Firefox：

```bash
python3 -m http.server 4175 --bind 127.0.0.1 -d miniprogram/dist/build/h5
npx playwright screenshot --browser=firefox --wait-for-timeout=1500 --viewport-size=390,844 'http://127.0.0.1:4175/#/pages/home/index' miniprogram/docs/validation/phase-6/screenshots/home-mobile.png
```

阶段截图和命令记录在 `docs/validation/phase-*/`。Chromium 在当前环境缺系统依赖，使用 Firefox 作为浏览器截图工具。

## 包体策略

首版小程序不整包复制根目录 `assets/`。`check:size` 读取 `dist/build/mp-weixin`，门禁为主包 `< 2 MB`、单分包 `< 2 MB`、总包 `< 8 MB`，主包 `>= 1.8 MB` 会预警。

当前首版资源只包含：

- `src/static/phase2/home-hero-moon-web.jpg`
- `src/static/phase2/home-hero-sunset-web.jpg`
- `src/static/phase2/home-hero-watercolor.png`
- `src/static/phase2/macau-cruise-night-banner-web.jpg`
- `src/static/phase2/taxi-scan-illustration-web.jpg`
- `src/static/phase2/ticket-wallet-illustration-web.jpg`
- `src/static/phase2/zhuhai-bay-home-hero-web.jpg`
- `src/static/phase3/verify-hero-wide-clean-web.png`

当前主包为 `1986.7 KB`，已接近 2 MB 硬门禁；后续新增首页大图应优先远程化或继续压缩。

## 依赖兼容性

| 依赖 | 用途 | H5 | mp-weixin | 备注 |
|---|---|---|---|---|
| `@dcloudio/uni-app` / `uni-components` / `uni-h5` / `uni-mp-weixin` | uni-app 跨端运行与构建 | 通过 | 通过构建 | 使用 uni API，不直接使用浏览器 API |
| `vue` | 页面和组件 | 通过 | 通过构建 | Vue 3 `<script setup>` |
| `vite` / `@dcloudio/vite-plugin-uni` | 构建 | 通过 | 通过 | 无运行时依赖 |
| `typescript` / `vue-tsc` | 类型检查 | 通过 | N/A | `typecheck` 已纳入门禁 |
| `vitest` | adapter 和工具测试 | 通过 | N/A | 测试中 mock `uni` |
