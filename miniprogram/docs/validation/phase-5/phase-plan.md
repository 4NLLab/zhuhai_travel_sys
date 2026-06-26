# Phase 5 执行计划

日期：2026-06-26

## 范围

- 将主包首页公开产品/分类和环岛游公开查询接入 `local` adapter。
- 司机端在 `local` 下探测 `/driver/me` 鉴权边界，缺 active driver token 时按 `allowFallback` 回落 Mock。
- 以 `backend/routes/router.go` 为准更新接口契约矩阵。
- 使用 Docker Compose 后端执行 curl 探针，记录公开接口、鉴权接口、供应商依赖和交易边界。
- 采集 H5 mock/local 模式截图；`mp-weixin` 继续以构建和包体证据为准。

## 不做

- 不修复环岛游交易接口当前公开、支付幂等、供应商正式联调等生产问题。
- 不伪造 user/driver token、微信登录、微信支付或供应商出票成功。
- 不把 Docker 本地临时 seed 当作正式测试数据来源。
