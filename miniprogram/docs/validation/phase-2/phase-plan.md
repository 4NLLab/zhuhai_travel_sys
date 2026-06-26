# Phase 2 执行级 Plan

## 阶段边界

- 迁移主包首页、我的、订单列表和票券详情基础视觉。
- 来源限定为 `index.html` 的首页/我的/订单信息架构，以及 `ticket.html` 的票券视觉。
- 不迁移环岛游完整下单链路、司机端完整链路、支付、供应商锁票/售票/退改逻辑。
- 不复制 preview/qa 资源，不搬运根目录大资源；仅使用压缩后的 Phase 2 小图。

## 文件与模块

| 模块 | 文件 | 目的 |
|---|---|---|
| 页面 | `src/pages/home`、`src/pages/profile`、`src/pages/orders`、`src/pages/ticket` | 主包首页、我的、订单、票券基础壳 |
| 组件 | `EmptyState`、`ProductEntryCard`、`OrderSummaryCard`、`TicketSummaryCard`、`StatusPill` | 复用商品卡、订单卡、票券卡、状态与空态 |
| 类型 | `src/types/main-shell.ts` | `ProductEntry`、`OrderSummary`、`TicketSummary`、`UserSummary` 等通用 view model |
| Mock | `src/mock/main-shell.ts` | 覆盖订单空态、待使用、待支付、已预约、已退款、票券不可用、未登录、失败场景 |
| Adapter | `src/api/main-shell.ts`、`src/utils/main-shell-mappers.ts` | 页面 API 入口和后端字段到 view model 的纯映射 |
| 资源 | `src/static/phase2/*.jpg` | 压缩后的首页/票券/接送视觉图 |
| 文档 | `docs/contracts/api-contract-matrix.md`、`docs/validation/phase-2` | Phase 2 契约与验收证据 |

## 数据流

页面只调用 `loadMainShellViewModel()` 或 `loadTicketDetail()`。API 层读取 runtime config，当前 `mock` 模式返回 `getMainShellMock()`；真实用户接口留到 Phase 5 接入。字段映射集中在 `main-shell-mappers.ts`，mock 和未来 API response mapper 共用同一转换逻辑。

## 验证命令

- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 浏览器与小程序验证

- H5：启动 `corepack pnpm -C miniprogram dev:h5 --host 127.0.0.1 --port 5173`，用 Playwright Firefox CLI 截取首页、我的、订单、票券详情桌面和移动截图。
- mp-weixin：通过 uni CLI 生成 `dist/build/mp-weixin`，记录包体。当前环境没有微信开发者工具或 `miniprogram-ci`，只能声明 build-validated。

## 风险与回滚边界

- Phase 2 没有接真实 token 和网络请求，不能声明真实订单/票券接口已联调。
- 票券二维码是小程序兼容的静态视觉 fallback，不是可核销真实码。
- 所有实现限制在 `miniprogram/`，不修改根静态 HTML 和后端业务逻辑。

## 两轮三角色审查

### 第一轮

| 角色 | 发现 | 调整 |
|---|---|---|
| Builder reviewer | 订单和票券 mapper 不能和 mock 形成循环依赖。 | 将 mapper 移到 `src/utils/main-shell-mappers.ts`。 |
| Skeptic reviewer | 只展示成功订单会掩盖空态、未登录和不可用票券。 | Mock 覆盖 `phase2-empty`、`phase2-unauthorized`、`phase2-failure`，成功场景含待支付、已预约、已退款和不可用票券。 |
| Verifier reviewer | 页面迁移必须留下截图证据，而不是只靠构建。 | 增加 8 张 H5 桌面/移动截图。 |

### 第二轮

| 角色 | 发现 | 调整 |
|---|---|---|
| Builder reviewer | H5 标签如 `span/strong` 会增加小程序差异。 | 改成 `view/text` 基础组件。 |
| Skeptic reviewer | 首页目的地卡片在移动端受 button 默认样式影响，网格视觉不稳定。 | 补 `width/box-sizing/margin` 和移动端两列断点，并重新截图。 |
| Verifier reviewer | 订单/票券字段映射需要单元测试，不应只靠页面渲染。 | 增加 `main-shell-adapter.test.ts`，覆盖订单列表和票券详情映射。 |

第二轮后无 P0/P1 计划缺陷，进入提交准备。
