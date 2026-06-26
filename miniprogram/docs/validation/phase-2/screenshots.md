# Phase 2 截图记录

采集时间：2026-06-26

| 页面 | 桌面截图 | 移动截图 | 结果 |
|---|---|---|---|
| 首页 | `screenshots/home-desktop.png` | `screenshots/home-mobile.png` | 通过；首屏非空白，目的地两列移动布局可读 |
| 我的 | `screenshots/profile-desktop.png` | `screenshots/profile-mobile.png` | 通过；用户摘要、票券/订单入口和常用服务可读 |
| 订单 | `screenshots/orders-desktop.png` | `screenshots/orders-mobile.png` | 通过；筛选、订单状态、票券入口可读 |
| 票券详情 | `screenshots/ticket-desktop.png` | `screenshots/ticket-mobile.png` | 通过；票券视觉、券码脱敏、须知和订单信息可读 |

## 截图命令

使用 Playwright Firefox CLI，分别以 `1440x900` 和 `390x844` viewport 访问：

- `/#/pages/home/index`
- `/#/pages/profile/index`
- `/#/pages/orders/index`
- `/#/pages/ticket/index`
