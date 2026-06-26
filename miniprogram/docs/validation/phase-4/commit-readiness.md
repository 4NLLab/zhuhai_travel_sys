# Phase 4 提交准备

采集时间：2026-06-26

## 准备状态

- [x] 司机端保持在 `subpackages/driver`。
- [x] 登录、注册、待审核、active 工作台、推广码、钱包、佣金、提现申请、提现记录完成。
- [x] 司机类型、mock、API、mapper 和 adapter 测试完成。
- [x] API 契约矩阵补 Phase 4。
- [x] H5 桌面/移动截图完成。
- [x] `mp-weixin` 构建和包体记录完成。
- [x] 五类提交前 review 完成。

## 待提交文件类别

- `miniprogram/src/subpackages/driver/pages/home/index.vue`
- `miniprogram/src/types/driver.ts`
- `miniprogram/src/mock/driver.ts`
- `miniprogram/src/api/driver.ts`
- `miniprogram/src/utils/driver-mappers.ts`
- `miniprogram/tests/driver-adapter.test.ts`
- `miniprogram/src/pages/profile/index.vue`
- `miniprogram/docs/contracts/api-contract-matrix.md`
- `miniprogram/docs/validation/phase-4/*`

## 阻塞说明

微信开发者工具运行时验收、真实 active driver token、后台审核和提现打款链路缺环境，只能作为 blocker 记录；不阻塞 Phase 4 build-validated 提交。
