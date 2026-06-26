# Phase 5 命令记录

日期：2026-06-26

| 命令 | 结果 | 说明 |
|---|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 | Vue/TS 类型检查通过 |
| `corepack pnpm -C miniprogram test` | 通过 | 7 个测试文件，20 个测试；新增 local contract adapter 测试 |
| `cd backend && go test ./...` | 通过 | 后端全部 package 编译测试通过 |
| `cd backend && go build -o /tmp/zhuhai_travel_backend_server .` | 通过 | 后端二进制可构建 |
| `timeout 8s /tmp/zhuhai_travel_backend_server` | 未通过 | 本机未启动 MySQL 时返回 `dial tcp 127.0.0.1:3306: connect: connection refused` |
| `docker compose up -d mysql backend` | 通过 | `zhuhai-travel-mysql` healthy，`zhuhai-travel-backend` running |
| `curl http://127.0.0.1:8080/api/v1/health` | 通过 | 200 `{"status":"ok"}` |
| `curl http://127.0.0.1:8080/api/v1/products?size=3` | 通过 | 最小 seed 后返回 1 条商品 |
| `curl http://127.0.0.1:8080/api/v1/categories` | 通过 | 最小 seed 后返回 1 条分类 |
| `curl /api/v1/island-cruise/smart-search?...` | 部分通过 | 200，但 `recommended=null`，真实供应商未配置 |
| `curl /api/v1/island-cruise/ports`、`/cert-types`、`/voyages`、`/price` | 未通过 | 502 `环岛游接口账号未配置` |
| `curl /api/v1/orders`、`/api/v1/driver/me`、`/api/v1/admin/me` | 通过 | 无 token 均返回 401，鉴权边界成立 |
| `curl -X POST /api/v1/payments/callback ...` | 通过 | 无签名返回 401，独立签名边界成立 |
| `corepack pnpm -C miniprogram build:h5` | 通过 | mock 和 local+fallback 模式均构建通过 |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 | 微信小程序构建通过 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 | 禁用平台 API 扫描通过 |
| `corepack pnpm -C miniprogram check:size` | 通过 | 主包 407.0 KB，司机分包 14.6 KB，总包 421.5 KB |
| `npx playwright screenshot --browser=firefox ...` | 通过 | 生成 mock/local H5 截图 |
