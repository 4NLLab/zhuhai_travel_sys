# Phase 1 截图验收

H5 本地预览：`http://127.0.0.1:5173/`

截图命令：`npx playwright screenshot --browser=firefox --wait-for-timeout=1000 --viewport-size=<size> <url> <file>`

| 页面 | 桌面 1440x900 | 窄屏 390x844 |
|---|---|---|
| 首页 | `screenshots/home-desktop.png` | `screenshots/home-mobile.png` |
| 订单 | `screenshots/orders-desktop.png` | `screenshots/orders-mobile.png` |
| 我的 | `screenshots/profile-desktop.png` | `screenshots/profile-mobile.png` |

## 验收结论

- 首页、订单、我的三条 H5 路由均可直接打开。
- tabBar 可见，当前页面高亮正常。
- 首页展示 Mock 模式、scenario、platform、adapter 状态。
- 移动端截图未见明显文本重叠、内容遮挡或底部 tabBar 覆盖正文。
- Phase 1 只做占位页，业务空态和失败态将在 Phase 2-4 按切片补齐。
