# Phase 5 Blockers

日期：2026-06-26

- 环岛游供应商账号未配置：`/island-cruise/ports`、`/cert-types`、`/voyages`、`/price` 返回 502 `环岛游接口账号未配置`；`smart-search` 返回 200 空推荐，不能视为真实供应商班次通过。
- 缺 user token 和订单/票券 seed：`/orders` 无 token 返回 401；用户订单、票券详情仍保留 Mock/fallback。
- 缺 active driver token 和佣金/提现 seed：`/driver/me` 无 token 返回 401；司机钱包、佣金、提现仍保留 Mock/fallback。
- 本机直接运行后端二进制需要本地 MySQL 3306；当前通过 Docker Compose 后端完成联调。
- 当前环境无微信开发者工具或 `miniprogram-ci`；`mp-weixin` 状态仍为 build-validated，正式提测前需补小程序运行时验证。
- Playwright Chromium CLI 仍缺系统依赖；本阶段继续使用 Firefox CLI 完成 H5 截图。
