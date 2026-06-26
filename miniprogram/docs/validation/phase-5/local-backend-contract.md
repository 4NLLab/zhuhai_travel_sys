# Phase 5 本地后端契约结果

日期：2026-06-26

## 环境

- Docker Compose 服务：`mysql`、`backend`
- 后端 Base URL：`http://127.0.0.1:8080/api/v1`
- H5 local 配置：`VITE_API_MODE=local VITE_API_BASE_URL=http://127.0.0.1:8080/api/v1 VITE_ALLOW_FALLBACK=true`
- `mp-weixin` local 配置要求：`VITE_API_BASE_URL=http://<LAN-IP>:8080/api/v1` 或本地代理域名；微信开发者工具可临时关闭合法域名校验，真机/体验版必须使用 HTTPS 合法域名。

## 结果摘要

| 分组 | 结果 |
|---|---|
| 进程健康 | `/health` 200 |
| 公开产品/分类 | 最小 seed 后 `/products`、`/categories`、`/products/schedules/query` 均 200 且有数据 |
| 用户鉴权 | `/orders` 无 token 返回 401 |
| 司机鉴权 | `/driver/me` 无 token 返回 401；前端 local 模式按 `allowFallback=true` 回落 Mock |
| 后台鉴权 | `/admin/me` 无 token 返回 401 |
| 支付回调 | 无签名 `/payments/callback` 返回 401 |
| 环岛游查询 | `/smart-search` 200 但无推荐；`/ports`、`/cert-types`、`/voyages`、`/price` 返回 502，原因是供应商账号未配置 |
| 环岛游交易 | `/lock` 空请求返回 400；未做锁票/出票成功断言 |

## 不能判定通过的项

- 环岛游供应商真实查询、锁票、出票、退票、改签：缺供应商账号与受控订单。
- 用户订单/票券真实读取：缺 user token、订单 seed、票券 seed。
- 司机钱包/佣金/提现真实读取：缺 active driver token、佣金/提现 seed。
- 微信 DevTools/真机 runtime：当前环境无微信开发者工具或 `miniprogram-ci`，只能标记为 build-validated。
