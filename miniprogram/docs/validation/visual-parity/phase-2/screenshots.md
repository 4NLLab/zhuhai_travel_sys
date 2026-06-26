# Phase 2 首页截图验收

日期：2026-06-26

## 截图证据

| 视口 | 原版基线 | 修复前迁移版 | 修复后迁移版 | 结论 |
|---|---|---|---|---|
| `390x844` | `../phase-1/screenshots/original-home-390x844.png` | `../phase-1/screenshots/migrated-home-390x844.png` | `screenshots/migrated-home-after-390x844.png` | 通过 |
| `519x927` | `../phase-1/screenshots/original-home-519x927.png` | `../phase-1/screenshots/migrated-home-519x927.png` | `screenshots/migrated-home-after-519x927.png` | 通过 |

## 人工视觉结论

`通过：迁移版与原版截图肉眼无明显区别`

已关闭 Phase 1 首页 P0/P1 差异：

- hero 已从夜景卡片改回水彩珠海湾区插画。
- 顶部品牌、搜索框、问候语和轮播文案已回到 hero 叠加结构。
- 热门目的地已恢复三列浅色彩块卡片，并只展示首屏前六项。
- 首屏下缘已恢复船票促销卡和黄色“立即购票”按钮。
- 分类 tabs、商品卡和底部导航已对齐原版视觉语言。

可接受微差：

- 目的地图标和 logo 为小程序兼容 CSS/字符绘制，不是原版 lucide/SVG 的完全一致绘制。
- H5 与静态页字体抗锯齿存在轻微差别。
