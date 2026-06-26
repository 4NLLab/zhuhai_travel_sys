# Phase 3 提交准备

采集时间：2026-06-26

## 准备状态

- [x] 范围限制在 `miniprogram/`。
- [x] 环岛游详情、班次、实名、锁座支付、出票票券基础闭环完成。
- [x] 退改签仅保留入口和后续接入说明。
- [x] 环岛游类型、mock、API、mapper 和 adapter 测试完成。
- [x] API 契约矩阵补 Phase 3。
- [x] H5 桌面/移动截图完成。
- [x] `mp-weixin` 构建和包体记录完成。
- [x] 五类提交前 review 完成。

## 待提交文件类别

- `miniprogram/src/pages/island-cruise/index.vue`
- `miniprogram/src/types/island-cruise.ts`
- `miniprogram/src/mock/island-cruise.ts`
- `miniprogram/src/api/island-cruise.ts`
- `miniprogram/src/utils/island-cruise-mappers.ts`
- `miniprogram/tests/island-cruise-adapter.test.ts`
- `miniprogram/src/pages.json`
- `miniprogram/src/pages/home/index.vue`
- `miniprogram/src/mock/main-shell.ts`
- `miniprogram/docs/contracts/api-contract-matrix.md`
- `miniprogram/docs/validation/phase-3/*`

## 阻塞说明

微信开发者工具运行时验收、真实供应商/支付联调和退改签 seed 缺环境，只能作为 blocker 记录；不阻塞 Phase 3 build-validated 提交。
