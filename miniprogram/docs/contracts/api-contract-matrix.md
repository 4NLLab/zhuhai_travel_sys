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

## Phase 2 主包首页 / 我的 / 订单 / 票券基础壳

Phase 2 只迁移主包展示壳与通用 view model。当前实现默认 `VITE_API_MODE=mock`，真实接口仅记录为 Phase 5 联调候选；禁止在 local/test/prod 静默回退到 mock。

| 页面动作 | endpoint | 方法 | auth 类型 | 字段映射 | token/账号/seed | 当前能否本地联调 | allowFallback | dataSource | 错误态 |
|---|---|---|---|---|---|---|---|---|---|
| 首页 Banner / 分类入口 | `/banners`、`/categories` | GET | public | `Banner.image_url/title/link`、`ProductCategory.name/type` -> `DestinationEntry` / 首页头图；Phase 2 以 `index.html` 静态切片和压缩图片占位 | 无 | 可联调但尚未接入小程序请求层 | false | mock | 接口失败时显示首页错误提示，不回退到伪数据 |
| 首页精选商品 | `/products?keyword=&category_id=&type=&page=&size=` | GET | public | `Product.name/summary/type/min_price/images` -> `ProductEntry.title/subtitle/category/priceLabel/imageUrl` | 无 | 可联调但字段需 Phase 5 校准 | false | mock | 空数组显示“暂无商品”；非 2xx 显示错误态 |
| 我的用户摘要 | `/users/profile` | GET | user_token | `User.nickname/mobile/realname_status` -> `UserSummary.displayName/maskedMobile/realnameStatus`；订单/票券计数暂由订单票券聚合 | 缺稳定测试用户 token/seed | Phase 5 联调 | false | mock | 401/403 显示未登录；非 2xx 显示场景提示 |
| 常用出行人入口 | `/travelers` | GET | user_token | `Traveler.name/id_type/id_no/phone/is_default` -> 常用服务入口后续列表；Phase 2 只放入口 | 缺稳定测试用户 token/seed | Phase 5 联调 | false | mock | 401/403 显示未登录；空列表不阻断“我的”页 |
| 订单列表 | `/orders?status=&page=&size=` | GET | user_token | `Order.id/order_no/status/product_name/travel_date/quantity/paid_amount` -> `OrderSummary`；覆盖 `pending_use/pending_pay/reserved/completed/refunded/refunding` | `VITE_MOCK_SCENARIO=phase2-success|phase2-empty|phase2-unauthorized|phase2-failure` | Phase 5 联调 | false | mock | 空态显示“暂无订单”；401/403 显示未登录；非 2xx 显示订单服务不可用 |
| 票券详情 | `/tickets/:id` | GET | user_token | `Ticket.id/status/product_name/valid_date/valid_time/code/verify_location/notices/order_no` -> `TicketSummary`；券码必须通过 `maskSensitive` 脱敏展示 | `VITE_MOCK_SCENARIO=phase2-success` 含 `available/unavailable` | Phase 5 联调 | false | mock | 404 显示“暂无票券”；401/403 显示未登录；非 2xx 显示读取失败 |

## Phase 3 环岛游订票纵向切片

Phase 3 当前只接 mock view model。以下真实 endpoint 均来自 `backend/routes/router.go`，其中锁票、出票、退票、改签、核销通知当前虽挂公开路由，但涉及库存/资金/供应商状态，前端契约中标注为 `current_public_but_should_be_hardened`，Phase 5 只能按受控测试环境联调，不能在生产语义下视为普通公开接口。

| 页面动作 | endpoint | 方法 | auth 类型 | 字段映射 | token/账号/seed | 当前能否本地联调 | allowFallback | dataSource | 错误态 |
|---|---|---|---|---|---|---|---|---|---|
| 码头与证件类型 | `/island-cruise/ports`、`/island-cruise/cert-types` | GET | public | `portId/portName` -> `IslandPort`；`certTypeId/certTypeName` -> cert option | 无 | Phase 5 联调 | false | mock | 失败时使用页面错误态，不静默伪造 local 成功 |
| 近期推荐 | `/island-cruise/smart-search?start_date=&days=&people_num=` | GET | public | `recommended.date/first_time/count/min_price` -> `recommended` | 无 | Phase 5 联调 | false | mock | 无推荐显示“暂无推荐班次”；非 2xx 显示服务失败 |
| 班次日历 | `/island-cruise/voyage-calendar?start_date=&days=&up_port_id=&down_port_id=&people_num=` | GET | public | `days[].date/count` -> 日期/班次数提示 | 无 | Phase 5 联调 | false | mock | 查询失败保留手动日期选择和错误提示 |
| 可售班次 | `/island-cruise/voyages?departure_date=&up_port_id=&down_port_id=&people_num=` | GET | public | `voyageId/voyageName/voyageNo/shipName/departureTime/cabinPriceList[].fareTypeList[]` -> `IslandVoyage`/`IslandFare` | `phase3-success|phase3-no-voyage|phase3-failure|phase3-insufficient-stock` | Phase 5 联调 | false | mock | 无班次、库存不足、接口失败均有可解释状态 |
| 保留座位 | `/island-cruise/lock` | POST | current_public_but_should_be_hardened | `local_order_no/orderNo/ticketNo/lock_expire_at` -> `IslandLockedOrder`；请求体由 `IslandOrderDraft` 生成 | `phase3-lock-expired|phase3-passenger-invalid` | Phase 5 需受控 seed | false | mock | 乘客缺失、库存不足、锁票失败、锁过期提示用户重新选择 |
| 取消保留 | `/island-cruise/unlock` | POST | current_public_but_should_be_hardened | `local_order_no/status` -> 取消结果 | 需可取消测试订单 | Phase 5 需受控 seed | false | mock | 已出票不可取消；失败保留当前订单状态 |
| 模拟支付出票 | `/island-cruise/sale` | POST | current_public_but_should_be_hardened | `ticketNo/codeContent/paid_at` -> `IslandTicketResult`；`codeContent` 必须专用脱敏 | `phase3-sale-failed` | Phase 5 需受控 seed | false | mock | 出票失败显示待处理，不展示完整核销码 |
| 查询订单 | `/island-cruise/order?local_order_no=` | GET | current_public_but_should_be_hardened | `status/passengers/ticket_no/code_content/pay_amount/go_time` -> 票券详情与售后状态 | 需已出票测试订单 | Phase 5 需受控 seed | false | mock | 404 显示订单不存在；非 2xx 显示刷新失败 |
| 退票费用 / 提交退票 | `/island-cruise/refund-fee`、`/island-cruise/refund` | GET/POST | current_public_but_should_be_hardened | Phase 3 只保留入口和说明，不执行真实资金流 | 需已出票且可退 seed | 不在 Phase 3 联调 | false | mock | 页面明确“后续接入”，不宣称退票闭环 |
| 改签费用 / 班次 / 锁定 / 取消 | `/island-cruise/change-fee`、`/island-cruise/change-voyages`、`/island-cruise/change-lock`、`/island-cruise/change-unlock` | GET/POST | current_public_but_should_be_hardened | Phase 3 只保留入口和说明，不执行供应商改签闭环 | 需已出票且可改签 seed | 不在 Phase 3 联调 | false | mock | 页面明确“后续接入”，不宣称改签闭环 |

## Phase 4 司机端纵向切片

Phase 4 当前只接 mock view model。司机端工作台真实接口均需要 `driver_active_token`；注册/登录是公开入口，但登录后必须由后端中间件校验 driver token 和 active 状态。Phase 5 联调前需要准备 active、pending、余额不足、提现成功等司机 seed。

| 页面动作 | endpoint | 方法 | auth 类型 | 字段映射 | token/账号/seed | 当前能否本地联调 | allowFallback | dataSource | 错误态 |
|---|---|---|---|---|---|---|---|---|---|
| 司机注册申请 | `/driver/register` | POST | public | `driver_no/status/message` -> 注册结果与待审核状态 | `phase4-pending-review`；需测试手机号/车辆 seed | Phase 5 联调 | false | mock | 参数错误显示注册提示；成功后显示待审核 |
| 司机登录 | `/driver/login` | POST | public | `access_token/driver_id/driver_no/name/status/commission_rate` -> driver session/profile | `phase4-active|phase4-login-failed` | Phase 5 联调 | false | mock | 账号密码错误显示登录失败；非 active 后续接口不可进入 |
| 司机资料刷新 | `/driver/me` | GET | driver_active_token | `driver_no/name/phone/status/car_plate_no/vehicle_type/commission_rate/qr_code` -> `DriverProfileSummary` | 缺 active driver token/seed | Phase 5 联调 | false | mock | 401/403 显示未登录或待审核 |
| 钱包余额 | `/driver/wallet` | GET | driver_active_token | `available/pending_total/settled_total/withdrawn_total` -> `DriverWalletSummary` | `phase4-active|phase4-empty-wallet` | Phase 5 联调 | false | mock | 空钱包显示 0；非 2xx 显示服务失败 |
| 佣金列表 | `/driver/commissions?page=&size=&status=` | GET | driver_active_token | `order_no/commission_amount/status/created_at` -> `DriverCommissionSummary[]` | `phase4-active|phase4-empty-wallet` | Phase 5 联调 | false | mock | 空列表显示“暂无佣金记录” |
| 提现申请 | `/driver/withdraw` | POST | driver_active_token | `withdrawal_no/amount/status/message` -> 提现申请结果 | `phase4-insufficient-balance|phase4-withdraw-success` | Phase 5 联调 | false | mock | 余额不足显示可提现余额；成功显示待管理员审核打款 |
| 提现记录 | `/driver/withdrawals?page=&size=` | GET | driver_active_token | `withdrawal_no/amount/status/account/created_at` -> `DriverWithdrawalSummary[]`，账号必须脱敏 | `phase4-active|phase4-empty-wallet` | Phase 5 联调 | false | mock | 空列表显示“暂无提现记录” |

## Phase 5 本地后端契约联调

Phase 5 直接以 `backend/routes/router.go` 为准。当前前端支持 `VITE_API_MODE=mock|local` 和 `VITE_ALLOW_FALLBACK=true|false`，页面代码不需要因 mock/local 切换而改动。H5 local 使用 `VITE_API_BASE_URL=http://127.0.0.1:8080/api/v1`；`mp-weixin` local 需要改成局域网 IP 或代理域名，真机/体验版必须使用 HTTPS 合法域名。

| 页面动作 | endpoint | 方法 | auth 类型 | token/账号/seed | local 结果（2026-06-26） | allowFallback | dataSource | 错误态 / 后续要求 |
|---|---|---|---|---|---|---|---|---|
| 健康检查 | `/health` | GET | public | Docker Compose 后端 | 200 `{"status":"ok"}` | false | local | 只证明进程存活，不替代业务联调 |
| 首页分类 | `/categories` | GET | public | Phase 5 最小 seed：`船票/ship` | 200，返回 1 条 active 分类 | true | local | 无 seed 时为空数组；页面展示空态或 fallback |
| 首页商品 | `/products?size=3` | GET | public | Phase 5 最小 seed：`澳门环岛游夜景船票` + SKU + 图片 + 排期 | 200，返回 1 条 active 商品，`sale_price=88` | true | local | 无 seed 时为空数组；adapter 测试覆盖字段映射 |
| 商品排期 | `/products/schedules/query?product_id=910001&date=2026-07-15` | GET | public | Phase 5 最小排期 seed | 200，返回 1 条 `19:30` 排期 | true | local | 交易下单仍需 user token，不在本阶段闭环 |
| 环岛游推荐 | `/island-cruise/smart-search?start_date=2026-07-15&days=1&people_num=1` | GET | public wrapper + supplier dependency | 缺 `ISLAND_CRUISE_DISTRIBUTOR_CODE` / `ISLAND_CRUISE_ACCESS_TOKEN` | 200，但 `recommended=null`、各 route `count=0`；后端吞掉供应商错误作为空推荐 | true | local | 页面显示暂无班次；不能判定供应商真实班次联调通过 |
| 环岛游码头/证件/班次/价格 | `/island-cruise/ports`、`/cert-types`、`/voyages`、`/price` | GET | public wrapper + supplier dependency | 缺环岛游供应商账号 | 502 `环岛游接口账号未配置` | true | fallback | 保留 Mock/fallback；正式接入前必须配置供应商沙箱/正式凭据 |
| 环岛游锁票 | `/island-cruise/lock` | POST | current_public_but_should_be_hardened | 需航班、乘客、供应商、订单 seed | 空请求 400 `缺少航班或港口参数` | true | fallback | 当前仍公开，应在后续目标中加鉴权/幂等/风控 |
| 支付回调 | `/payments/callback` | POST | independent_signature | 需支付平台签名头 | 无签名 401 `支付回调签名无效` | false | local | 不是小程序公开接口；继续保持独立签名校验 |
| 用户订单 | `/orders` | GET | user_token | 缺测试 user token 和订单/票券 seed | 无 token 401 `请先登录` | true | fallback | 前端不得当公开接口；需准备 user、order、ticket seed |
| 司机工作台 | `/driver/me` | GET | driver_active_token | 缺 active driver token/driver seed | 无 token 401 `请先登录` | true | fallback | 前端可回落 Mock 展示；真实钱包/佣金需 active driver seed |
| 后台接口 | `/admin/me`、`/drivers`、`/commissions/*`、`/tickets/verify` | GET/POST | admin_super_admin_token | 缺 admin token | 无 token 401 `请先登录` | false | N/A | 不迁入小程序用户端 |

### 最小准备清单

| 缺口 | 最小准备项 | 期望状态 | 验证 endpoint |
|---|---|---|---|
| 用户订单/票券 | active user、Bearer token、至少 1 个 paid/valid 订单、1 张 ticket、1 条 verification 可选记录 | `/orders` 200 且订单能映射到 `OrderSummary`；`/tickets/:id` 200 且券码脱敏 | `/orders`、`/tickets/:id` |
| 司机钱包/佣金 | active driver、登录密码、vehicle、driver_qr_code、pending/settled commission、withdrawal 记录 | `/driver/me`、`/driver/wallet`、`/driver/commissions`、`/driver/withdrawals` 200 | `/driver/login`、`/driver/me`、`/driver/wallet` |
| 环岛游供应商 | `ISLAND_CRUISE_DISTRIBUTOR_CODE`、`ISLAND_CRUISE_ACCESS_TOKEN`、供应商测试可售航班、可锁/可出票账号余额 | 查询类返回真实 ports/cert/voyages；锁票/出票在受控订单上返回供应商订单号和票号 | `/island-cruise/ports`、`/voyages`、`/lock`、`/sale` |
| 支付回调 | `PAYMENT_WEBHOOK_SECRET`、测试 payment、合法 timestamp/signature | 无签名 401；合法签名且金额匹配时更新 payment/order | `/payments/callback` |

### Phase 5 本地 seed SQL

本阶段本地 Docker MySQL 使用了最小公开业务 seed，完整 SQL 记录在 `miniprogram/docs/validation/phase-5/minimum-public-seed.sql`。正式测试环境应通过迁移/seed 管道写入等价数据，不应依赖手工容器状态。
