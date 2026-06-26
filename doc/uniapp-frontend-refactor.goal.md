---
name: uniapp-frontend-refactor
overview: 用 uni-app + Vue 3 + TypeScript 将现有静态 HTML 前端重构为可持续开发的小程序前端，按业务纵向切片同步完成页面迁移与 Mock/API 适配，真实接入与新增后端能力后置。
todos:
  - id: phase-1-platform-foundation
    content: Phase 1：搭建 uni-app 工程与小程序平台护栏
    status: completed
  - id: phase-2-main-shell-orders
    content: Phase 2：迁移主包首页、我的、订单/票券基础切片
    status: completed
    depends_on:
      - phase-1-platform-foundation
  - id: phase-3-island-cruise-slice
    content: Phase 3：迁移环岛游订票纵向切片
    status: completed
    depends_on:
      - phase-2-main-shell-orders
  - id: phase-4-driver-slice
    content: Phase 4：迁移司机端纵向切片
    status: completed
    depends_on:
      - phase-2-main-shell-orders
  - id: phase-5-local-contract
    content: Phase 5：本地后端契约联调与接口矩阵
    status: completed
    depends_on:
      - phase-3-island-cruise-slice
      - phase-4-driver-slice
  - id: phase-6-cross-platform-hardening
    content: Phase 6：跨端质量、包体治理、文档与交付收口
    status: pending
    depends_on:
      - phase-5-local-contract
isProject: false
---

# uni-app 小程序前端重构分阶段计划

## /goal 提示词：uni-app 小程序前端重构

将下面整份 Markdown 作为 `/goal` 输入使用。目标是在不先接入真实微信登录、微信支付、短信和供应商正式账号的前提下，用 `uni-app + Vue 3 + TypeScript` 新建小程序前端工程，并把当前静态 HTML 原型中首版需要的页面按业务纵向切片迁移为可运行、可验收、可继续开发的小程序前端。每个业务切片必须同时完成页面、状态、类型、Mock 场景、API adapter、测试和截图验收，避免把页面迁移做成不可维护的静态壳。

## 背景

- [已验证] 当前根目录存在 12 个静态 HTML 页面：`index.html`、`ticket.html`、`admin.html`、`driver.html`、`flow.html`、`island-cruise-booking.html`、`island-cruise.html`、`hotel.html`、`hotel-list.html`、`macau.html`、`hongkong.html`、`car.html`。
- [已验证] 当前仓库没有 `package.json`、`pages.json`、`manifest.json`、`miniprogram/` 或已有 uni-app/Taro 小程序工程。
- [已验证] `assets/` 目录有 69 个文件，约 32 MB，包含 `preview-*`、`qa-*` 截图和多个 1 MB 以上大图；本次小程序重构必须筛选和压缩首版实际需要的静态资源，不能整体搬入小程序包。
- [已验证] `island-cruise-booking.html` 直接设置 `API_BASE = "/api/v1/island-cruise"`，并调用 ports、cert-types、smart-search、voyage-calendar、voyages、lock、unlock、sale、order、refund、change 等接口。
- [已验证] `driver.html` 和 `admin.html` 使用 `window.ZHUHAI_API_BASE || "/api/v1"`；司机 token 当前依赖浏览器 `localStorage`，迁移时必须改为小程序可用的 storage adapter。
- [已验证] 现有 HTML 大量依赖浏览器能力：`document.querySelector*`、`window.*`、`fetch`、`localStorage`、`history`、`window.confirm`、外链 `<script>`、lucide CDN、qrcode CDN、canvas 二维码等；这些不能直接搬入 `mp-weixin`。
- [已验证] 后端实际路由以 `backend/routes/router.go` 为准；`backend/docs/ARCHITECTURE.md` 的 API 清单已落后，未完整覆盖当前环岛游和新版司机工作台路由。
- [已验证] 后端接口分为公开接口、用户 Bearer token 接口、司机 active token 接口、后台 super_admin token 接口；不能把订单、司机钱包、司机佣金误写成公开接口。
- [已验证] 当前公开路由示例包括 `/health`、公开登录、`/products`、`/products/:id`、`/products/schedules/query`、`/categories`、`/banners`、当前全部 `/island-cruise/*` 和 `/payments/callback`；其中环岛游锁票/出票/退改签/核销通知/支付回调虽然当前挂在公开路由下，前端计划中必须标注为 `auth_type=current_public_but_should_be_hardened` 或独立签名回调，不得当作普通公开查询接口处理。
- [已验证] 页面来源映射：`index.html` 的首页区域对应品牌、分类、产品和环岛游入口；`index.html` 的 mine/orders 区域对应“我的”和订单列表；`ticket.html` 与 `island-cruise-booking.html` 的 ticket screen 对应票券详情视觉；`driver.html` 对应司机登录、注册、钱包、佣金和提现申请；`flow.html` 只作为后续船票业务参考，不进入本轮目标。
- [已验证] 仓库提供 Docker Compose 本地运行方式：后端默认暴露 `8080`，前端 nginx 默认暴露 `8000`，`docker/nginx/default.conf` 将 `/api/` 代理到 `backend:8080`；小程序端不能依赖浏览器相对路径代理，`mp-weixin` 必须配置完整 base URL，并记录 DevTools 或真机所用的局域网 IP、HTTPS 域名、代理或“不校验合法域名”开发设置。
- [用户声称] 当前目标是第一步先用 `uni-app + Vue 3 + TypeScript` 重构前端，后面再开发其它功能和真实接入。
- [用户声称] 原“页面迁移”和“接口适配/Mock 闭环”应放在一起完成。本计划按业务纵向切片落实该要求：迁移哪个业务切片，就同步完成该切片的 Mock/API 适配。

## 核心目标

- 建立 `uni-app + Vue 3 + TypeScript` 小程序前端工程，作为后续小程序开发的唯一主线。
- 前置小程序平台护栏：`mp-weixin` 构建、H5 预览、分包策略、包体预算、request/storage/logger adapter、Mock 场景切换、禁用浏览器 API 扫描。
- 按纵向业务切片迁移页面：首页/我的/订单票券基础、环岛游订票、司机端；每个切片同时完成页面、状态、类型、Mock、API adapter 和截图验收。
- 明确现有 UI 迁移与新增/补齐 UI 的边界：例如首页司机入口、司机推广二维码、提现记录属于小程序信息架构补齐，不是从 `index.html` 或 `driver.html` 原样迁移。
- 保留现有静态 HTML 作为视觉和流程参考，不把它继续扩展成新功能承载层。
- 第一轮只完成前端重构可验收闭环，不以真实微信登录、微信支付、短信、供应商正式账号、生产上线为验收条件。

## 非目标

- 不在本目标内接入真实微信登录、手机号授权、短信验证码、微信支付、退款原路退回、提现真实打款。
- 不在本目标内重构 Go 后端业务逻辑、支付逻辑、环岛游供应商正式接入或生产安全一致性问题。
- 不把 Web 管理后台迁入小程序；`admin.html` 继续作为 Web 管理后台参考或保留页面。
- 不迁移酒店、香港、澳门、用车等非首版交易路径，除非它们只作为首页入口占位或静态参考。
- 不迁移 `flow.html` 九洲港至蛇口船票交易路径；本轮只可把它作为船票业务后续开发参考。
- 不删除现有静态 HTML 页面，除非后续用户明确要求。

## 总体原则

- `miniprogram/` 是小程序重构主目录，默认采用 `src/pages`、`src/components`、`src/api`、`src/adapters`、`src/stores`、`src/types`、`src/utils`、`src/mock`、`src/static` 结构；静态 HTML 只作为参考输入。
- 小程序验收以 `mp-weixin` 构建和小程序端约束为主；H5 预览用于浏览器截图和快速视觉验证，不能替代 `mp-weixin` 构建证据。最低证据阶梯为：`build:mp-weixin` 成功、检查 `dist/mp-weixin/app.json` 的 pages/tabBar/subPackages、记录主包/分包/总包体积；若微信开发者工具或 `miniprogram-ci` 可用，再打开核心页面并截图。不可用时不得声称已完成小程序运行时验证，只能记录 runtime blocker；整个目标最多标记为 build-validated，正式提测前必须补小程序运行时验证。
- 页面迁移必须组件化：表单、价格、票券卡片、订单状态、底部操作栏、空态/错误态不能散落复制。
- 前端 API 必须通过 `api/` 适配层访问，页面不能直接拼临时 URL、不能绑定后端原始字段、不能直接调用 `fetch`。
- 接口契约矩阵从 Phase 1 建立骨架，Phase 2-4 每迁移一个业务切片就补齐该切片的 endpoint、auth、字段映射、Mock 场景和 fallback 规则；Phase 5 只做本地联调校验和差异修正，不能等到 Phase 5 才第一次梳理接口。
- Mock 必须通过场景 ID 切换，并覆盖与当前切片相关的 success、empty、failure、unauthorized、expired、insufficient-stock、insufficient-balance 等状态。API adapter 返回或记录 `dataSource: mock | local | fallback`；`local` 模式默认不得静默 fallback，只有接口矩阵中显式 `allowFallback=true` 的场景才可降级。
- 静态资源必须按阶段筛选、压缩和命名；任何阶段都不得整包复制 `assets/`。视频不得进入小程序包；如必须保留演示视频，只能用 HTTPS 远程资源并提供静态封面 fallback。
- 请求层必须区分 H5、本地后端、`mp-weixin`、未来测试/生产环境；H5 可使用 dev proxy，小程序 DevTools 可使用局域网 IP/代理并记录域名校验设置，真机或体验版必须使用 HTTPS 合法域名或明确记录隧道/域名 blocker。
- storage、logger、confirm/modal、navigation、QR code/canvas 必须走小程序兼容 adapter；禁止在页面中直接使用 `window`、`document`、`localStorage`、`history`、外链 CDN 脚本。
- 敏感字段展示、Mock fixture、storage、测试快照和日志从 Phase 1 起就必须遵守脱敏规则，证件号、手机号、票码、收款账号、token 不得完整写入普通日志；`maskSensitive()`、`redactLogPayload()` 或等价能力必须有测试。
- 包体预算是硬门禁：主包 `< 2 MB`、每个分包 `< 2 MB`、总包 `< 8 MB`；主包 `< 1.8 MB` 为预警线，单个本地图片建议 `< 200 KB`，超过 200 KB 必须进入包体报告并说明不能继续压缩或远程化的原因。
- 每个 Phase 的验证证据统一记录在 `miniprogram/docs/validation/phase-N/` 或等价目录，至少包含 `commands.md`、`screenshots.md`、`mp-weixin-size.md`、`blockers.md` 或同等信息。
- 每个阶段保持一个清晰主提交；若 post-commit review 发现问题，按仓库规则创建后续修复 commit，阶段仍不得标记完成，直到 post-commit review 无问题。
- Phase 依赖以 frontmatter 的 DAG 为准；Phase 3 和 Phase 4 在 Phase 2 完成后可独立细化和执行，Phase 5 必须等待 Phase 3 与 Phase 4 都完成。
- 涉及代码 symbol 编辑前按仓库规则执行影响分析；提交前执行变更检测。

## 强制阶段门禁

1. 执行任何 Phase 前，先只阅读本 goal 文档和当前 Phase，确定阶段边界；随后在 Phase plan 中按需读取本阶段相关 HTML、后端路由、资源和配置证据。
2. 为当前 Phase 创建执行级 plan，明确文件、模块、数据流、测试、浏览器验证、`mp-weixin` 验证、风险、证据文件和回滚边界。
3. 对执行级 plan 进行至少两轮三角色对抗式优化：Builder reviewer、Skeptic reviewer、Verifier reviewer。
4. 若第二轮后仍有 P0/P1 计划缺陷，继续优化直到无 P0/P1 缺陷，或向用户报告 blocker。
5. 只有优化后的 Phase plan 可以进入实现。
6. 实现时优先把可拆分任务交给聚焦子代理；主代理负责阶段控制、集成、冲突处理和最终验证。
7. 完成当前 Phase 实现后，运行 Phase 指定的自动化测试、类型检查、语法检查、H5 构建、`mp-weixin` 构建和 smoke check。
8. 任何触及前端代码的 Phase，都必须启动真实本地 H5 页面做浏览器截图验收，并保存或记录 `mp-weixin` 构建/模拟器验证证据。
9. 逐条验证当前 Phase 的 Acceptance Criteria，并记录通过/失败结果。
10. 提交前并行执行五类专项 review：Security、Logic bug、Test coverage、Maintainability、Performance / concurrency。
11. 任一专项 review 发现问题后，必须修复、重跑相关测试和浏览器/小程序验证，然后重跑全部五类专项 review。
12. 只有五类提交前 review 都明确 `未发现问题` 后，才允许执行提交准备检查。
13. 提交准备检查必须包括 `git status`、完整 `git diff`、仓库要求的影响/变更分析、凭据/生成物/缓存/无关文件检查。
14. 每个 Phase 的主实现提交必须遵守 `.claude/rules/git-submit.mdc`。
15. 提交后必须复查最新 commit diff；若发现问题，创建新的修复 commit 并重复 post-commit review，直到无问题。
16. 只有 post-commit review 无问题后，才允许把当前 Phase todo 标记为 completed 并进入下一 Phase。

## Review 门禁细则

提交前五类专项 review 必须使用自包含 packet：用户目标、当前 Phase、相关 diff、测试结果、浏览器截图证据、`mp-weixin` 构建/模拟器证据、仓库约束。

- Security reviewer：检查鉴权绕过、路径穿越、HTML 注入、XSS、敏感数据泄露、权限边界、危险日志、凭据读取。
- Logic bug reviewer：检查状态流、错误分支、前后端契约、边界条件、数据兼容、回归风险。
- Test coverage reviewer：检查缺失的单元/集成/前端/浏览器覆盖，高风险路径是否有断言。
- Maintainability reviewer：检查模块边界、命名、重复、抽象层级、配置设计、文档和后续扩展成本。
- Performance/concurrency reviewer：检查大文件读取、包体压力、缓存、并发请求、锁、频控、超时、N+1 调用。

每个 reviewer 必须返回以下两类结果之一：

- `未发现问题`
- 或按严重程度排序的问题清单，包含具体文件/路径引用和建议修复方式。

## 浏览器截图验收门禁

- 前端 Phase 必须使用真实本地页面验证，不接受用户截图或口头确认替代。
- uni-app 小程序页面必须提供 H5 本地预览作为浏览器截图目标；可用命令应在 Phase 1 明确，例如 `pnpm -C miniprogram dev:h5` 或等价脚本。
- 每个前端 Phase 至少截取一个桌面视口和一个窄屏视口截图，建议覆盖 `1440x900` 和 `390x844`。
- `mp-weixin` 验收必须至少包含构建成功证据、`dist/mp-weixin/app.json` pages/tabBar/subPackages 检查、包体报告；如微信开发者工具或 `miniprogram-ci` 可用，还应包含小程序运行时截图。不可用时必须记录 runtime blocker，不能用 H5 截图冒充小程序运行时验证，也不能宣称已完成小程序运行时验收。
- 验证内容必须覆盖：路由跳转、导航高亮、加载态、空态、错误态、按钮可用性、文本截断、内容重叠、滚动行为、底部操作栏遮挡。
- 迁移环岛游订票、订单/票券和司机端页面时，必须分别截取关键页面截图，不能只截首页。
- 截图验收或 `mp-weixin` 构建不通过不得提交。

## Phase 1：搭建 uni-app 工程与小程序平台护栏

### Codex Goal

在仓库中新增 `miniprogram/` 小程序前端工程，固定为 `uni-app + Vue 3 + TypeScript`，并先建立小程序平台护栏：H5 与 `mp-weixin` 构建、主包/分包策略、包体预算、request/storage/logger adapter、Mock 场景切换、禁用浏览器 API 扫描和基础样式系统。

### Scope

- 新增 `miniprogram/` 工程目录和必要配置文件：`package.json`、`pnpm-lock.yaml`、TypeScript 配置、uni-app 配置、`pages.json`、`manifest.json` 或对应框架配置；`manifest` 使用安全占位 appid，本地真实 appid 放入未跟踪配置。
- 固定包管理器为 `pnpm`，至少提供 `dev:h5`、`build:h5`、`build:mp-weixin`、`typecheck`、`test`、`lint:platform`、`check:size`。
- 建立基础目录：`src/pages/`、`src/components/`、`src/api/`、`src/adapters/`、`src/stores/`、`src/types/`、`src/utils/`、`src/mock/`、`src/static/`、`docs/validation/`。
- 配置主包 tabBar：`首页`、`订单`、`我的`；司机端和低频页面按分包策略预留。
- 建立项目包体预算和 `check:size`：记录主包、每个分包、总体输出体积、最大单资源；硬门禁主包 `< 2 MB`、每个分包 `< 2 MB`、总包 `< 8 MB`，主包 `< 1.8 MB` 为预警线。
- 封装 request adapter、storage adapter、logger adapter、modal/confirm adapter、navigation adapter，禁止页面直接使用浏览器 API。
- 建立 Mock scenario switch，支持通过配置切换 `mock`、`local`、未来 `test/prod` 环境；`local` 模式默认不静默 fallback。
- 建立接口契约矩阵骨架，字段至少包含页面动作、endpoint、HTTP 方法、auth 类型、字段映射、需要的 token/测试账号/seed、当前能否本地联调、`allowFallback`、错误态。
- 建立全局样式变量和基础布局约束：颜色、字号、间距、按钮、表单、价格、卡片、底部操作栏、安全区。
- 新增最小占位页面，证明路由、tabBar、样式、request/storage/logger adapter 和 Mock 开关可以运行。
- 不迁移完整业务页面，不接真实后端接口，不处理真实登录和支付。

### Acceptance Criteria

- `miniprogram/` 下存在可运行的 `uni-app + Vue 3 + TypeScript` 工程，且不破坏根目录现有静态 HTML 和后端代码。
- `pnpm -C miniprogram install` 或项目实际安装命令可执行；lockfile 被提交；若环境缺失导致无法执行，必须记录 blocker。
- `pnpm -C miniprogram typecheck`、`pnpm -C miniprogram build:h5`、`pnpm -C miniprogram build:mp-weixin`、`pnpm -C miniprogram test`、`pnpm -C miniprogram lint:platform`、`pnpm -C miniprogram check:size` 或等价脚本有明确结果。
- Phase 1 验证证据记录当前仓库基线命令，例如根目录 HTML 清单、`assets/` 体积、当前无 `miniprogram/` 工程、Docker/nginx 本地代理配置位置。
- 首页、订单、我的三个基础页面可打开、可切换，司机端分包或入口存在明确占位。
- request/storage/logger/modal/navigation adapter 有最小测试或脚本校验；`maskSensitive()`、`redactLogPayload()` 或等价脱敏能力有测试，storage 只保存最小 token/profile，退出登录会清理。
- 页面不直接使用 `fetch`、`window`、`document`、`localStorage`、`history`、外链 `<script>`；`lint:platform` 扫描范围至少覆盖 `src/pages`、`src/components`、`src/api`、`src/adapters`、`src/stores`、`src/utils` 和 `dist/mp-weixin`。
- Mock fixture、测试快照和日志样例不得包含真实手机号、证件号、票码、收款账号或 token。
- 接口契约矩阵骨架已创建，并记录 H5、本地后端、DevTools、真机/体验版分别如何配置 base URL。
- 记录 H5 与 `mp-weixin` 构建产物体积、`dist/mp-weixin/app.json` pages/tabBar/subPackages 检查结果，且没有整包复制 `assets/`。
- 浏览器截图验收覆盖 H5 首页、订单页、我的页的桌面和窄屏视口；`mp-weixin` 构建、包体和可用时的小程序运行时证据被记录在 `miniprogram/docs/validation/phase-1/`。
- 提交前完成五类专项 review；全部 `未发现问题` 后才允许提交。
- 提交后再次 review 最新 commit；review 无问题后才允许进入 Phase 2。

## Phase 2：迁移主包首页、我的、订单/票券基础切片

### Codex Goal

迁移主包基础体验切片：从 `index.html` 抽取首页、我的、订单列表基础结构，从 `ticket.html` 抽取票券详情视觉参考，同时建立通用 domain types、Mock scenario 机制、订单/票券 adapter 约定和主包组件体系。

### Scope

- 迁移首页核心内容：城市文旅品牌、搜索/分类入口、产品入口、环岛游入口；订单入口主要来自 `index.html` 的 mine/orders 区域，若放到首页必须标注为新增小程序信息架构。
- 可在首页或“我的”露出司机入口，但必须标注为小程序新增信息架构入口；司机端完整流程仍属于 Phase 4。
- 迁移“我的”基础页：用户信息占位、常用出行人入口、订单入口、客服/协议入口；来源参考 `index.html` 的 mine 区域。
- 迁移订单/票券基础壳：订单列表、订单详情入口、票券卡片、票券详情视觉；来源参考 `index.html` 的 orders/mine 区域、`ticket.html` 和环岛游出票页。
- 建立通用组件：应用壳、tabBar 适配、产品入口卡片、订单卡片、票券卡片、状态徽标、空态、错误态、底部操作栏。
- 建立通用 domain types 和 view model：产品入口、最小 `OrderSummary`、最小 `TicketSummary`、用户摘要；只放通用 `source/type/status` 等字段，环岛游专属字段留到 Phase 3。
- 建立 Mock scenario 机制和基础场景：订单为空、待支付、待使用、已退款、票券不可用、未登录。
- 建立 adapter 目录约定和测试样例：Mock/后端原始字段 -> 前端 view model；同步补齐接口契约矩阵中主包基础切片的字段和 fallback 规则。
- 不迁移环岛游订票完整流程，不迁移司机端完整流程，不接真实用户订单接口。

### Acceptance Criteria

- H5 与 `mp-weixin` 构建均通过；主包 tabBar 页面可打开、可切换、可返回。
- 首页、我的、订单列表、票券详情基础页可在 Mock 模式下展示 success、empty、unauthorized、failure 场景。
- 首页司机入口被标注为新增小程序入口；订单/票券基础壳区分通用订单与环岛游订单，不混用字段。
- 页面不依赖原 HTML 文件、浏览器 DOM 查询、`window.ZHUHAI_API_BASE` 或页面内联脚本。
- API adapter 测试覆盖至少一种订单列表和一种票券详情 view model 映射。
- 接口契约矩阵已补齐 Phase 2 页面动作，用户订单接口标注为需 user token，Mock/local fallback 规则明确。
- 本阶段只复制或压缩实际使用的资源；提交中不得包含 `preview-*`、`qa-*`、非首版大图或整包 assets 复制。
- 浏览器截图验收覆盖首页、我的、订单列表、票券详情的桌面和窄屏视口；`mp-weixin` 构建、包体和可用时的小程序运行时证据被记录在 `miniprogram/docs/validation/phase-2/`。
- 提交前完成五类专项 review；全部 `未发现问题` 后才允许提交。
- 提交后再次 review 最新 commit；review 无问题后才允许启动依赖 Phase 2 的后续阶段（Phase 3 和 Phase 4）。

## Phase 3：迁移环岛游订票纵向切片

### Codex Goal

按纵向切片迁移环岛游订票流程：页面、状态、Mock、types、adapter、测试和截图一起完成。首轮只做可演示的开发期闭环，不接真实微信支付和供应商正式账号。

### Scope

- 来源参考 `island-cruise.html` 和 `island-cruise-booking.html`，迁移环岛游详情/入口、班次选择、乘客填写、订单确认、模拟支付、出票结果、票券展示。
- 建立 `island-cruise` domain types、view model、adapter 和 Mock fixtures。
- Mock 场景覆盖：班次查询成功、无班次、查询失败、库存不足、锁票过期、乘客校验失败、模拟支付成功、出票失败。
- 二维码/票码展示必须使用小程序兼容方案，不直接迁移 CDN qrcode、DOM canvas 或 `window.QRCode` 逻辑。
- `window.confirm`、`history.back`、`fetch`、`document.querySelector*`、外链 lucide/qrcode 脚本都必须替换为 uni-app/Vue/adapter 方案。
- 退款、改签只做入口、状态占位和“后续接入说明”；不在本 Phase 要求完整资金流、库存回滚或真实供应商退改签闭环，不能在文案中暗示已闭环。
- 原型中的视频不得进入小程序包；如保留导览视频入口，只能使用 HTTPS 远程视频并提供静态封面 fallback。
- 同步补齐接口契约矩阵中的环岛游查询、锁票、出票、退票、改签相关 endpoint、auth 类型、字段映射、`allowFallback` 和错误态。
- 只可引用当前切片实际需要且压缩后的图片资源。

### Acceptance Criteria

- Mock 模式下可完整演示：首页环岛游入口 -> 班次选择 -> 乘客填写 -> 订单确认 -> 模拟支付 -> 出票结果 -> 票券展示。
- 环岛游 adapter 测试覆盖班次、乘客、订单确认、票券结果四类 view model 映射。
- 表单校验覆盖联系人、手机号、乘客姓名、证件类型、证件号、乘客数量不一致等错误。
- 锁票倒计时、重复提交保护、库存不足、锁票过期和接口失败有明确页面反馈。
- 敏感字段展示层默认脱敏，证件号、手机号、票码、token 不得完整输出到普通日志。
- 二维码在 H5 和可用的小程序运行时证据中可见；如小程序端二维码生成能力受限，必须展示可复制/可核验的码内容 fallback。
- H5 与 `mp-weixin` 构建均通过；不得引入外链脚本或浏览器专用 API。
- 浏览器截图验收覆盖班次页、乘客填写页、订单确认/模拟支付页、出票结果/票券页，以及至少一个失败/过期状态；`mp-weixin` 构建、包体和可用时的小程序运行时证据被记录在 `miniprogram/docs/validation/phase-3/`。
- 提交前完成五类专项 review；全部 `未发现问题` 后才允许提交。
- 提交后再次 review 最新 commit；review 无问题后才算满足 Phase 5 的环岛游依赖。

## Phase 4：迁移司机端纵向切片

### Codex Goal

按纵向切片迁移司机端：登录/注册参考、状态、二维码、钱包、佣金、提现申请、提现记录、types、adapter、Mock、测试和截图一起完成。真实审核、真实提现打款和真实后台运营后置。

### Scope

- 来源参考 `driver.html`，迁移司机登录、注册申请、钱包、佣金、提现申请等已有 UI。
- 首页司机入口、推广二维码、提现记录属于小程序端补齐能力；后端存在相关路由/模型，但 `driver.html` 并未完整展示，计划中必须标注为新增/补齐 UI。
- 建立 driver domain types、view model、adapter 和 Mock fixtures。
- Mock 场景覆盖：未登录、登录失败、待审核、active、钱包为空、佣金为空、余额不足、提现提交成功、提现记录为空、接口失败。
- token 存储必须通过 storage adapter，不得使用 `localStorage`。
- 推广二维码/码内容展示必须使用小程序兼容方案，不直接依赖浏览器 DOM 或外链二维码库。
- 同步补齐接口契约矩阵中的司机注册、登录、profile、wallet、commission、withdrawal 相关 endpoint、auth 类型、字段映射、`allowFallback` 和错误态。
- 不要求本 Phase 使用真实司机账号、管理员审核、真实提现打款或真实后端司机接口。

### Acceptance Criteria

- Mock 模式下可完整演示：司机入口 -> 登录/注册占位 -> 待审核/active 状态 -> 推广二维码 -> 钱包 -> 佣金明细 -> 提现申请 -> 提现记录。
- driver adapter 测试覆盖 profile、wallet、commission、withdrawal 四类 view model 映射。
- 未登录、待审核、active、余额不足、提现成功、提现记录为空都有明确页面状态。
- token、手机号、收款账号默认脱敏，不写入普通日志。
- 推广二维码在 H5 和可用的小程序运行时证据中可见；如小程序端二维码生成能力受限，必须展示推广码内容 fallback。
- H5 与 `mp-weixin` 构建均通过；司机端页面可按分包策略运行。
- 浏览器截图验收覆盖司机登录/注册、司机首页、钱包、佣金、提现申请、提现记录的桌面和窄屏视口；`mp-weixin` 构建、包体和可用时的小程序运行时证据被记录在 `miniprogram/docs/validation/phase-4/`。
- 提交前完成五类专项 review；全部 `未发现问题` 后才允许提交。
- 提交后再次 review 最新 commit；review 无问题后才算满足 Phase 5 的司机端依赖；Phase 5 必须同时等待 Phase 3 和 Phase 4 完成。

## Phase 5：本地后端契约联调与接口矩阵

### Codex Goal

将重构后的小程序前端切换到本地 Docker/Go 后端能提供的现有接口上，明确公开接口、用户鉴权接口、司机鉴权接口和后台接口边界，验证 Mock/local 切换并输出后续真实接入清单。

### Scope

- 使用当前 Docker Compose 或本地 Go 后端启动方式，验证小程序 API base URL 能切换到本地后端；H5 可走 dev proxy，`mp-weixin` 必须记录局域网 IP/代理/域名校验设置，真机或体验版必须使用 HTTPS 合法域名或记录 blocker。
- 建立接口契约矩阵，直接以 `backend/routes/router.go` 为准，不以过期架构文档为唯一依据。
- 接口矩阵至少包含：页面动作、endpoint、HTTP 方法、auth 类型、需要的 token/测试账号/seed、当前能否本地联调、`allowFallback`、`dataSource`、错误态。
- 公开接口清单以 `backend/routes/router.go` 为准并区分风险：`/health`、公开登录、`/products`、`/products/:id`、`/products/schedules/query`、`/categories`、`/banners`、环岛游查询类接口、当前公开但应加固的环岛游交易类接口、`/payments/callback` 独立签名回调。
- 公开接口验证不能只用 `/health` 代替业务联调；`/health` 只证明后端进程存活。
- 用户订单接口只在可取得测试 user token/测试数据时验证；否则保留 Mock/fallback 并记录缺口。
- 司机钱包/佣金/提现接口只在可取得 active driver token/测试数据时验证；否则保留 Mock/fallback 并记录缺口。
- 对缺失的 user/driver token、seed、账号和外部依赖，必须输出最小准备清单，包含所需角色、最小 SQL/seed、测试账号字段、期望状态和对应验证 endpoint。
- 环岛游锁票、出票、退票、改签涉及供应商依赖；本阶段不得把供应商失败伪装成已联调通过，必须记录真实接口、Mock fallback 和后续接入要求。
- 不在本 Phase 修复后端交易鉴权、支付幂等、供应商正式联调等生产问题；只记录为后续目标。

### Acceptance Criteria

- H5 与 `mp-weixin` 构建均通过；`mock` 与 `local` 环境可通过配置切换，不需要改页面代码。
- `/health` 成功、至少一个产品/分类公开业务接口真实读取成功、至少一个环岛游公开查询接口真实读取成功，并在页面或 adapter 测试中验证；若 DB/seed 缺失导致业务接口不可用，必须记录具体 SQL/seed/account 缺口，不能判定通过。
- 至少一次未授权访问用户或司机鉴权接口返回 401/403，证明权限边界未被前端误当公开接口。
- 用户鉴权接口和司机鉴权接口的验证结果真实记录：成功、缺测试账号、缺 seed、外部依赖失败或保留 Mock，不能用公开接口通过替代。
- 环岛游供应商相关接口的 local 结果被明确记录；若返回 502 或外部失败，必须说明 fallback 行为。
- 接口契约矩阵被写入前端 README 或专门文档，后续真实登录/支付/供应商接入可直接使用。
- 浏览器截图验收覆盖 Mock 模式和 local 模式下的关键页面；`mp-weixin` 构建、包体和可用时的小程序运行时证据被记录在 `miniprogram/docs/validation/phase-5/`。
- 提交前完成五类专项 review；全部 `未发现问题` 后才允许提交。
- 提交后再次 review 最新 commit；review 无问题后才允许进入 Phase 6。

### Completion Record

- 主提交：`95d1c176c567a5a0edd0d0006c9b20e3cc21cfe8`（`feat(miniprogram): 接入本地后端契约联调`）。
- 完成内容：`local` 模式已接入公开产品/分类、环岛游 `smart-search` 查询和司机鉴权边界 fallback；接口契约矩阵已更新到 Phase 5。
- 验证材料：`miniprogram/docs/validation/phase-5/`，包含 Docker 后端 curl 探针、最小公开 seed、mock/local 截图、包体记录、blocker、提交前五类 review 和提交后 review。
- 真实结论：`/health`、`/products`、`/categories`、`/products/schedules/query` 已在本地 Docker 后端验证；`/orders`、`/driver/me`、`/admin/me` 无 token 返回 401；支付回调无签名返回 401。
- 保留缺口：环岛游供应商账号未配置，ports/cert-types/voyages/price 返回 502；用户订单/票券和司机钱包/佣金/提现缺 token 与 seed，继续保留 Mock/fallback。

## Phase 6：跨端质量、包体治理、文档与交付收口

### Codex Goal

完成前端重构收口：跨端 QA、包体和资源复核、平台禁用 API 扫描、README、重构完成报告和后续开发清单。

### Scope

- 复核主包/分包输出体积和资源归属，删除未使用资源，记录包体预算执行情况。
- 按 Phase 2-5 的页面/场景矩阵，复核所有首版页面在 H5 与 `mp-weixin` 构建下的主要路径、空态、错误态、窄屏布局和底部操作栏。
- 运行禁用平台 API 扫描，覆盖 `window`、`document`、`localStorage`、`fetch`、`history`、外链 `<script>`、不兼容 canvas/二维码调用。
- 补充前端 README：安装、启动、H5 预览、`mp-weixin` 构建、Mock 场景切换、本地后端联调、截图验收方式、目录说明、包体策略。
- 更新重构完成报告：已迁移页面、已接入/Mock 的接口、仍需后端或真实接入补齐的能力、后续业务开发建议。
- 列出所有新增前端依赖及其 H5、`mp-weixin` 兼容性结论；第三方库内部平台限制必须写入风险或替代方案。
- 保持现有静态 HTML 和 Go 后端不被无关重构破坏。

### Acceptance Criteria

- `typecheck`、`test`、`build:h5`、`build:mp-weixin`、`lint:platform`、`check:size` 全部通过；无法执行的项目必须有具体环境 blocker。
- 禁用平台 API 扫描无违规；若第三方库内部不可避免，必须有隔离说明和小程序端验证证据。
- 主包、分包、总包体积被记录，且未包含首版无关的 `preview-*`、`qa-*` 或非首版大图。
- README 能让新开发者完成安装、启动、H5 预览、小程序构建、Mock/local 切换和截图验收。
- 重构完成报告列出首版页面、接口适配状态、Mock 留存点、后续真实接入和业务开发事项。
- README 或重构完成报告包含新增前端依赖清单和小程序兼容性结论。
- 浏览器截图验收覆盖 Phase 2-5 页面/场景矩阵中的所有首版主路径和至少一个失败/空态场景；`mp-weixin` 构建/模拟器验证证据被记录在 `miniprogram/docs/validation/phase-6/`。
- 提交前完成五类专项 review；全部 `未发现问题` 后才允许提交。
- 提交后再次 review 最新 commit；review 无问题后才允许把整个 goal 标记为完成。

## 总体最终验收标准

- `miniprogram/` 成为新的小程序前端主工程，使用 `uni-app + Vue 3 + TypeScript`。
- H5 与 `mp-weixin` 构建、类型检查、测试、平台禁用 API 扫描均有可复现命令和结果。
- 首页、我的、订单/票券基础、环岛游订票、司机端完成首轮迁移，并能通过 Mock 场景完整演示。
- 每个业务切片都同步完成页面、状态、types、Mock、adapter、测试和截图验收。
- API 访问集中在适配层，页面不直接绑定临时 URL 或后端原始字段。
- 接口契约矩阵明确公开接口、用户鉴权接口、司机鉴权接口、后台接口和供应商依赖，不误把需鉴权接口当公开接口。
- 小程序资源经过筛选和压缩，首版无关图片不进入主包。
- H5 浏览器截图验收和 `mp-weixin` 构建/模拟器验收均有记录。
- 若执行环境缺少微信开发者工具或 `miniprogram-ci`，最终状态只能标记为 build-validated，并在正式提测前补齐小程序运行时验收。
- 每个 Phase 有清晰主提交，提交前五类专项 review 和提交后 review 均完成。
- 现有静态 HTML 和 Go 后端未被无关重构破坏。
- 真实微信登录、微信支付、短信、供应商正式接入、后端安全一致性问题已作为后续目标记录，而不是混入本次前端重构。

## 当前执行建议

下一步只细化 Phase 1 的执行 plan：初始化 `miniprogram/` 工程、固定脚本和 lockfile、配置 H5 与 `mp-weixin` 构建、建立主包 tabBar、分包策略、包体预算、request/storage/logger/modal/navigation adapter、Mock 场景切换、禁用浏览器 API 扫描和全局样式。不要在 Phase 1 迁移完整页面，也不要接入真实登录、支付或供应商接口。用户要求的“页面迁移和 API/Mock 适配一起做”已落实为 Phase 2-4 的纵向切片原则：迁移哪个业务切片，就同步完成该切片的 Mock/API 适配。
