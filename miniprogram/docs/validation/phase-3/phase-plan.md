# Phase 3 执行级 Plan

## 阶段边界

- 迁移环岛游订票纵向切片：详情/班次、出行人、锁座待支付、模拟支付出票、电子票展示。
- 退票和改签只做入口、状态占位和后续接入说明，不声明真实资金流、库存回滚或供应商退改签闭环。
- 不接真实微信支付、短信、供应商正式账号；真实接口只进入契约矩阵和 Phase 5 联调候选。
- 不复制视频，不引入外链二维码库，不使用 `window`、`document`、`fetch`、`history`、`localStorage`。

## 文件与模块

| 模块 | 文件 | 目的 |
|---|---|---|
| 页面 | `src/pages/island-cruise/index.vue` | 环岛游多步骤状态机页面，支持 query step 用于截图 |
| API | `src/api/island-cruise.ts` | mock/local 分支入口，local 暂不静默 fallback |
| Mock | `src/mock/island-cruise.ts` | 航线、班次、乘客、锁座、出票和失败场景 |
| 类型 | `src/types/island-cruise.ts` | 环岛游 domain types 和 view model |
| Mapper | `src/utils/island-cruise-mappers.ts` | 后端/供应商字段到前端 view model 映射 |
| 测试 | `tests/island-cruise-adapter.test.ts` | 班次、乘客、锁座、出票映射和错误场景 |
| 文档 | `docs/contracts/api-contract-matrix.md`、`docs/validation/phase-3` | Phase 3 契约、验证和审查证据 |

## 数据流

`island-cruise/index.vue` 调用 `loadIslandCruiseFlow()` 读取 view model。页面内的锁座、取消、模拟支付调用 `lockIslandCruiseOrder()`、`unlockIslandCruiseOrder()`、`saleIslandCruiseOrder()`；当前 mock 模式返回固定订单和票券结果。字段映射集中在 mapper，避免页面绑定后端原始字段。

## Mock 场景

- `phase3-success`：有推荐、可售班次、实名乘客、锁座成功、模拟支付出票成功。
- `phase3-no-voyage`：无可售班次。
- `phase3-failure`：查询服务失败。
- `phase3-insufficient-stock`：班次余票不足。
- `phase3-lock-expired`：锁座后过期，支付页提示重新选择。
- `phase3-passenger-invalid`：实名信息校验失败。
- `phase3-sale-failed`：模拟支付后出票失败。

## 验证命令

- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 浏览器截图

H5 截图覆盖：

- 详情/班次：`/#/pages/island-cruise/index?step=detail`
- 出行人：`/#/pages/island-cruise/index?step=traveler`
- 锁座/模拟支付：`/#/pages/island-cruise/index?step=pay`
- 出票/票券：`/#/pages/island-cruise/index?step=ticket`
- 失败/无班次：`/#/pages/island-cruise/index?scenario=phase3-no-voyage`

每个关键页至少保留桌面或移动截图；详情、出行人、支付、票券保留桌面和移动截图。

## 两轮三角色审查

### 第一轮

| 角色 | 发现 | 调整 |
|---|---|---|
| Builder reviewer | 如果拆成多个页面，阶段成本和导航状态会膨胀。 | 使用单页状态机，query step 用于截图直达。 |
| Skeptic reviewer | 退改入口容易被误认为已完成真实闭环。 | 文案和契约矩阵明确只做占位，真实退改留后续。 |
| Verifier reviewer | 原 HTML 依赖 DOM/canvas/QRCode，不能搬代码。 | 票券码采用静态兼容视觉和可复制码内容 fallback。 |

### 第二轮

| 角色 | 发现 | 调整 |
|---|---|---|
| Builder reviewer | Mock 只成功不够证明状态机。 | 增加无班次、库存不足、锁过期、乘客校验失败、出票失败场景。 |
| Skeptic reviewer | 页面若直接用供应商字段会影响 Phase 5 改接口。 | 建立 `IslandCruiseViewModel` 和 mapper 测试。 |
| Verifier reviewer | 必须证明入口从 Phase 2 主包可达。 | 更新首页环岛游入口到 `/pages/island-cruise/index` 并截图。 |

第二轮后无 P0/P1 计划缺陷，进入实现。
