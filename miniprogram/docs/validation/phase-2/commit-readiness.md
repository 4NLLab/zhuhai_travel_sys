# Phase 2 提交准备

采集时间：2026-06-26

## 准备状态

- [x] 范围限制在 `miniprogram/`。
- [x] 首页、我的、订单列表、票券详情基础视觉完成。
- [x] 组件、类型、mock、adapter 测试完成。
- [x] API 契约矩阵补 Phase 2。
- [x] H5 桌面/移动截图完成。
- [x] `mp-weixin` 构建和包体记录完成。
- [x] 五类提交前 review 完成。

## 待提交文件类别

- `miniprogram/src/pages/*`
- `miniprogram/src/components/*`
- `miniprogram/src/types/main-shell.ts`
- `miniprogram/src/mock/main-shell.ts`
- `miniprogram/src/api/main-shell.ts`
- `miniprogram/src/utils/main-shell-mappers.ts`
- `miniprogram/tests/main-shell-adapter.test.ts`
- `miniprogram/docs/contracts/api-contract-matrix.md`
- `miniprogram/docs/validation/phase-2/*`
- `miniprogram/src/static/phase2/*`

## 阻塞说明

微信开发者工具运行时验收缺环境，只能作为 blocker 记录；不阻塞 Phase 2 build-validated 提交。
