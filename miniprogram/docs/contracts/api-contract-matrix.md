# 接口契约矩阵

本矩阵从 Phase 1 建立骨架，后续 Phase 2-5 按业务切片补齐。后端路由以 `backend/routes/router.go` 为准，`backend/docs/ARCHITECTURE.md` 仅作为背景参考。

## 环境与 Base URL

| 运行环境 | API 模式 | Base URL | 配置方式 | 备注 |
|---|---|---|---|---|
| H5 Mock | `mock` | 不访问网络 | `VITE_API_MODE=mock`、`VITE_MOCK_SCENARIO=phase1-success` | 默认模式 |
| H5 本地后端 | `local` | `http://127.0.0.1:8080/api/v1` 或 `/api/v1` proxy | `VITE_API_MODE=local`、`VITE_API_BASE_URL=...` | devServer 可代理 `/api` |
| 微信 DevTools | `local` | 局域网 IP 或本地代理完整 URL | `VITE_API_BASE_URL=http://<LAN-IP>:8080/api/v1` | 开发期可关闭合法域名校验；需记录设置 |
| 真机/体验版 | `test/prod` | HTTPS 合法域名 | `VITE_API_MODE=test|prod` | 无域名时记录 blocker |

## Phase 1 骨架

| 页面动作 | endpoint | 方法 | auth 类型 | 字段映射 | token/账号/seed | 当前能否本地联调 | allowFallback | dataSource | 错误态 |
|---|---|---|---|---|---|---|---|---|---|
| 平台健康占位 | `/health` | GET | public | 后端 `{status}` -> Phase 1 状态占位 | 无 | 可在后端启动后验证 | false | local | 非 2xx 显示平台状态读取失败 |
| Mock 平台状态 | 无 | N/A | none | scenario -> `PhaseOneStatus` | `VITE_MOCK_SCENARIO` | 不访问后端 | N/A | mock | unknown scenario 回到 success |
| 用户订单入口占位 | `/orders` | GET | user_token | Phase 2 补齐 | 缺测试 user token/seed | Phase 2/5 验证 | false | mock/local | 未登录、401/403、空列表 |
| 司机端入口占位 | `/driver/me` | GET | driver_active_token | Phase 4 补齐 | 缺 active driver token/seed | Phase 4/5 验证 | false | mock/local | 未登录、待审核、401/403 |

## 路由边界快照

| 类型 | 路由示例 | Phase 1 处理 |
|---|---|---|
| 公开接口 | `/health`、`/auth/phone-login`、`/products`、`/products/:id`、`/products/schedules/query`、`/categories`、`/banners` | 只记录，不在 Phase 1 接真实业务 |
| 当前公开但应加固 | `/island-cruise/lock`、`/island-cruise/sale`、`/island-cruise/unlock`、`/island-cruise/refund`、`/island-cruise/change-*`、`/island-cruise/verify-notify` | 标记 `auth_type=current_public_but_should_be_hardened`，Phase 3/5 细化 |
| 用户 Bearer token | `/users/profile`、`/travelers`、`/orders`、`/tickets/*` | Phase 2/5 补齐 |
| 司机 active token | `/driver/me`、`/driver/wallet`、`/driver/commissions`、`/driver/withdraw`、`/driver/withdrawals` | Phase 4/5 补齐 |
| 后台 super_admin token | `/admin/*`、`/drivers`、`/commissions/*`、`/tickets/verify` | 不迁入小程序用户端 |
