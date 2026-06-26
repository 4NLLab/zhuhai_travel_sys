# Phase 1 执行级 Plan

## 阶段边界

- 只新增 `miniprogram/` uni-app Vue 3 TypeScript 工程与平台护栏。
- 不迁移完整业务页面，不接真实登录、支付、供应商接口，不删除现有静态 HTML。
- 司机端只创建分包占位，完整司机切片留到 Phase 4。

## 文件与模块

| 模块 | 文件 | 目的 |
|---|---|---|
| 工程配置 | `package.json`、`vite.config.ts`、`tsconfig.json`、`pages.json`、`manifest.json` | 固定 pnpm、构建 H5/mp-weixin、主包 tabBar、司机分包 |
| 基础页面 | `src/pages/home`、`src/pages/orders`、`src/pages/profile`、`src/subpackages/driver/pages/home` | 验证路由、tabBar、分包、状态占位 |
| 平台 adapter | `src/adapters/*` | request/storage/logger/modal/navigation 统一入口 |
| Mock 与配置 | `src/mock/scenario.ts`、`src/utils/config.ts` | scenario ID 与 mock/local 切换 |
| 安全与校验 | `src/utils/sensitive.ts`、`scripts/lint-platform.mjs`、`scripts/check-size.mjs` | 脱敏、禁用平台 API 扫描、包体门禁 |
| 测试 | `tests/*.test.ts` | adapter、Mock、脱敏最小覆盖 |
| 文档证据 | `docs/contracts`、`docs/validation/phase-1` | 接口矩阵骨架与验证记录 |

## 数据流

页面只调用 `src/api/*`；API 层读取 `getRuntimeConfig()` 决定 mock/local。Mock 返回 `dataSource=mock`，local 通过 `requestJson()` 使用 `uni.request` 并返回 `dataSource=local`。页面不直接访问 URL、不直接使用浏览器 API。

## 验证命令

- `corepack pnpm -C miniprogram install`
- `corepack pnpm -C miniprogram typecheck`
- `corepack pnpm -C miniprogram test`
- `corepack pnpm -C miniprogram build:h5`
- `corepack pnpm -C miniprogram build:mp-weixin`
- `corepack pnpm -C miniprogram lint:platform`
- `corepack pnpm -C miniprogram check:size`

## 浏览器与小程序验证

- H5：启动 `corepack pnpm -C miniprogram dev:h5 --host 127.0.0.1`，用 Playwright 截取首页、订单、我的桌面与窄屏截图。
- mp-weixin：检查 uni CLI 实际产物 `dist/build/mp-weixin/app.json` 的 pages/tabBar/subPackages，记录包体。若本机无微信开发者工具或 `miniprogram-ci`，记录 runtime blocker，不能声明小程序运行时已验收。

## 风险与回滚边界

- DCloud 包版本采用同一 `vue3` dist-tag 系列；若安装或构建失败，优先调整版本 pin，不扩大到业务迁移。
- 不复制根 `assets/`，避免包体超标。
- 所有改动限制在 `miniprogram/`，回滚边界清晰。

## 两轮三角色审查

### 第一轮

| 角色 | 发现 | 调整 |
|---|---|---|
| Builder reviewer | 需要同时证明 H5 和 mp-weixin 构建，不应只创建文件。 | plan 加入完整命令和产物检查。 |
| Skeptic reviewer | `local` 模式若静默 fallback 会掩盖接口问题。 | `requestJson()` 不提供默认 fallback，矩阵显式 `allowFallback=false`。 |
| Verifier reviewer | 仅有脚本不足以证明禁用 API 覆盖 dist。 | `lint:platform` 扫描 `src` 关键目录和 `dist/mp-weixin`。 |

### 第二轮

| 角色 | 发现 | 调整 |
|---|---|---|
| Builder reviewer | 分包策略不能停留在文档。 | `pages.json` 增加司机端 subPackage 占位页。 |
| Skeptic reviewer | 日志脱敏必须有测试，避免样例泄漏手机号/证件号。 | 增加 `sensitive.test.ts` 和 storage 最小会话测试。 |
| Verifier reviewer | 包体预算需要机器可执行门禁。 | 增加 `check:size`，输出主包/分包/总包和最大文件。 |

第二轮后无 P0/P1 计划缺陷，进入实现。
