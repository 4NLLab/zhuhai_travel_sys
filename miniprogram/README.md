# 珠海文旅小程序前端

`miniprogram/` 是 uni-app + Vue 3 + TypeScript 小程序前端主工程。现阶段为 Phase 1 平台护栏：H5 与 mp-weixin 构建、主包 tabBar、司机分包占位、平台 adapter、Mock 场景、禁用浏览器 API 扫描和包体预算。

## 命令

```bash
corepack pnpm -C miniprogram install
corepack pnpm -C miniprogram dev:h5 --host 127.0.0.1 --port 5173
corepack pnpm -C miniprogram build:h5
corepack pnpm -C miniprogram build:mp-weixin
corepack pnpm -C miniprogram typecheck
corepack pnpm -C miniprogram test
corepack pnpm -C miniprogram lint:platform
corepack pnpm -C miniprogram check:size
```

## 环境切换

复制 `.env.example` 为本地 `.env` 后调整：

- `VITE_API_MODE=mock`：默认，不访问后端。
- `VITE_API_MODE=local`：访问 `VITE_API_BASE_URL` 指向的本地后端。
- `VITE_MOCK_SCENARIO=phase1-success|phase1-empty|phase1-failure`：切换 Phase 1 场景。
- `VITE_ALLOW_FALLBACK=false`：local 模式默认不静默降级。

## 目录

- `src/pages/`：主包 tabBar 页面。
- `src/subpackages/driver/`：司机端分包占位。
- `src/api/`：API adapter 入口，页面不得直接拼 URL。
- `src/adapters/`：request、storage、logger、modal、navigation。
- `src/mock/`：Mock 场景。
- `src/utils/`：配置和敏感信息脱敏。
- `docs/contracts/`：接口契约矩阵。
- `docs/validation/`：阶段验证证据。

## 包体策略

首版小程序不整包复制根目录 `assets/`。`check:size` 读取 `dist/build/mp-weixin`，门禁为主包 `< 2 MB`、单分包 `< 2 MB`、总包 `< 8 MB`，主包 `>= 1.8 MB` 会预警。
