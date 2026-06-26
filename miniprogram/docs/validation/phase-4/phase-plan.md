# Phase 4 执行级 Plan

## 阶段边界

- 迁移司机端纵向切片：登录/注册参考、审核状态、推广二维码 fallback、钱包、佣金、提现申请和提现记录。
- 真实管理员审核、真实提现打款、真实二维码生成和生产司机账号不进入本阶段。
- 司机端保持在 `subpackages/driver` 分包，不进入主包。
- 不搬 `driver.html` 中的 DOM、`fetch`、`localStorage`、全局 API base 或外链图标脚本。

## 文件与模块

| 模块 | 文件 | 目的 |
|---|---|---|
| 页面 | `src/subpackages/driver/pages/home/index.vue` | 司机端单页状态机 |
| API | `src/api/driver.ts` | mock/local 分支入口 |
| Mock | `src/mock/driver.ts` | 未登录、登录失败、待审核、active、钱包/佣金/提现场景 |
| 类型 | `src/types/driver.ts` | 司机 domain types 和 view model |
| Mapper | `src/utils/driver-mappers.ts` | 后端字段到司机 view model 映射 |
| 测试 | `tests/driver-adapter.test.ts` | profile、wallet、commission、withdrawal 映射 |
| 文档 | `docs/contracts/api-contract-matrix.md`、`docs/validation/phase-4` | Phase 4 契约、验证和审查证据 |

## Mock 场景

- `phase4-active`：active 司机，钱包、佣金、提现记录都有数据。
- `phase4-unauthorized`：未登录。
- `phase4-login-failed`：登录失败。
- `phase4-pending-review`：注册后待审核。
- `phase4-empty-wallet`：钱包/佣金/提现记录为空。
- `phase4-insufficient-balance`：提现金额超过余额。
- `phase4-withdraw-success`：提现提交成功。
- `phase4-failure`：司机服务失败。

## 验证命令

- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 截图路由

- `/#/subpackages/driver/pages/home/index?scenario=phase4-active`
- `/#/subpackages/driver/pages/home/index?scenario=phase4-pending-review`
- `/#/subpackages/driver/pages/home/index?scenario=phase4-empty-wallet`
- `/#/subpackages/driver/pages/home/index?scenario=phase4-insufficient-balance`

## 两轮三角色审查

### 第一轮

| 角色 | 发现 | 调整 |
|---|---|---|
| Builder reviewer | 司机工作台不应进入主包。 | 保持在 `subpackages/driver`，只从“我的”页提供入口。 |
| Skeptic reviewer | token 存储如果沿用 Web `localStorage` 会破坏小程序兼容。 | 只通过 storage adapter 暴露登录态占位，不在页面直接访问存储。 |
| Verifier reviewer | 真实 active token 缺失，不能声明已联调。 | 契约矩阵标明 Phase 5 需要 active driver token/seed。 |

### 第二轮

| 角色 | 发现 | 调整 |
|---|---|---|
| Builder reviewer | 只做 active 工作台不足以覆盖审核链路。 | 增加待审核、未登录、登录失败和注册结果状态。 |
| Skeptic reviewer | 提现入口容易被误认为真实打款。 | 文案明确“提交申请，待管理员审核打款”，不宣称转账完成。 |
| Verifier reviewer | 佣金/提现 mapping 必须可测。 | 增加 driver adapter 单测覆盖四类 view model。 |

第二轮后无 P0/P1 计划缺陷，进入实现。
