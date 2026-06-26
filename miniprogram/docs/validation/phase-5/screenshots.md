# Phase 5 截图记录

采集时间：2026-06-26

截图命令：`npx playwright screenshot --browser=firefox --wait-for-timeout=1500 --viewport-size=<size> <url> <file>`

| 页面 / 模式 | 桌面截图 | 移动截图 | 结果 |
|---|---|---|---|
| 首页 local | `screenshots/local-home-desktop.png` | `screenshots/local-home-mobile.png` | 通过；公开产品/分类来自本地后端 seed，用户订单/票券提示需 token/seed |
| 环岛游 local | N/A | `screenshots/local-island-mobile.png` | 通过；显示本地接口来源和暂无真实班次状态 |
| 司机端 local fallback | N/A | `screenshots/local-driver-mobile.png` | 通过；显示 `/driver/me` HTTP 401 后 fallback 到司机 Mock 工作台 |
| 首页 mock | N/A | `screenshots/mock-home-mobile.png` | 通过；mock 模式无需改页面代码 |

## 截图路由

- `/#/pages/home/index`
- `/#/pages/island-cruise/index`
- `/#/subpackages/driver/pages/home/index`
