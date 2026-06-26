# Phase 3 基础页截图验收

日期：2026-06-26

## 截图证据

| 页面 | 视口 | 原版基线 | 修复前迁移版 | 修复后迁移版 | 结论 |
|---|---|---|---|---|---|
| 我的 | `390x844` | `../phase-1/screenshots/original-profile-390x844.png` | `../phase-1/screenshots/migrated-profile-390x844.png` | `screenshots/migrated-profile-after-390x844.png` | 通过 |
| 我的 | `519x927` | `../phase-1/screenshots/original-profile-519x927.png` | `../phase-1/screenshots/migrated-profile-519x927.png` | `screenshots/migrated-profile-after-519x927.png` | 通过 |
| 订单 | `390x844` | `../phase-1/screenshots/original-orders-390x844.png` | `../phase-1/screenshots/migrated-orders-390x844.png` | `screenshots/migrated-orders-after-390x844.png` | 通过 |
| 订单 | `519x927` | `../phase-1/screenshots/original-orders-519x927.png` | `../phase-1/screenshots/migrated-orders-519x927.png` | `screenshots/migrated-orders-after-519x927.png` | 通过 |
| 票券 | `390x844` | `../phase-1/screenshots/original-ticket-390x844.png` | `../phase-1/screenshots/migrated-ticket-390x844.png` | `screenshots/migrated-ticket-after-390x844.png` | 通过 |
| 票券 | `519x927` | `../phase-1/screenshots/original-ticket-519x927.png` | `../phase-1/screenshots/migrated-ticket-519x927.png` | `screenshots/migrated-ticket-after-519x927.png` | 通过 |

## 人工视觉结论

`通过：迁移版与原版截图肉眼无明显区别`

已关闭 Phase 1 中基础页 P0/P1 差异：

- 我的页已恢复蓝绿页头、资料卡、下一段行程、票券 KPI、订单入口、服务宫格和自绘底部导航。
- 订单页已恢复返回我的页头、状态筛选、紧凑订单卡、状态标签、详情按钮和原版默认订单集合。
- 票券页已恢复状态栏、顶部操作、待使用状态、核销横幅、大二维码票码卡、商品/景区/须知/订单信息和底部操作栏。

可接受微差：

- 顶部操作、底部导航和二维码为小程序兼容的本地 CSS / SVG data 绘制，不追求与原版 SVG 逐像素一致。
- H5 与静态页字体抗锯齿、系统状态栏图标存在轻微渲染差异。
