# Phase 6 截图记录

采集时间：2026-06-26

截图命令：`npx playwright screenshot --browser=firefox --wait-for-timeout=1500 --viewport-size=390,844 <url> <file>`

| 页面 / 场景 | 截图 | 结果 |
|---|---|---|
| 首页 mock 主路径 | `screenshots/home-mobile.png` | 通过；目的地、商品和底部 tabBar 可读 |
| 我的 mock 主路径 | `screenshots/profile-mobile.png` | 通过；用户摘要、订单/票券入口和司机入口可读 |
| 订单 mock 主路径 | `screenshots/orders-mobile.png` | 通过；筛选和订单卡片可读 |
| 票券详情 mock 主路径 | `screenshots/ticket-mobile.png` | 通过；票码脱敏和须知可读 |
| 环岛游 mock 选班次 | `screenshots/island-detail-mobile.png` | 通过；班次、票种和底部 CTA 可读 |
| 环岛游 mock 支付态 | `screenshots/island-pay-mobile.png` | 通过；订单确认和支付 CTA 可读 |
| 环岛游 mock 空态 | `screenshots/island-no-voyage-mobile.png` | 通过；无班次状态可解释，底部 CTA 不遮挡文本 |
| 司机端 mock active | `screenshots/driver-active-mobile.png` | 通过；推广码、钱包和佣金入口可读 |
| 司机端 mock 空钱包 | `screenshots/driver-empty-mobile.png` | 通过；钱包与记录空态可读 |
| 首页 local | `screenshots/local-home-mobile.png` | 通过；本地产品/分类 seed 可展示 |
| 环岛游 local | `screenshots/local-island-mobile.png` | 通过；供应商未配置时展示暂无班次 |
| 司机端 local fallback | `screenshots/local-driver-mobile.png` | 通过；无 driver token 时显示 HTTP 401 并 fallback |

## 截图路由

- `/#/pages/home/index`
- `/#/pages/profile/index`
- `/#/pages/orders/index`
- `/#/pages/ticket/index`
- `/#/pages/island-cruise/index`
- `/#/pages/island-cruise/index?step=pay`
- `/#/pages/island-cruise/index?scenario=phase3-no-voyage`
- `/#/subpackages/driver/pages/home/index?scenario=phase4-active`
- `/#/subpackages/driver/pages/home/index?scenario=phase4-empty-wallet`
