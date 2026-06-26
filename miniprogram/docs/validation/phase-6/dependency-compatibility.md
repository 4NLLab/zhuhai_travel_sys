# Phase 6 依赖兼容性

日期：2026-06-26

| 依赖 | 类型 | H5 结论 | mp-weixin 结论 | 风险 |
|---|---|---|---|---|
| `@dcloudio/uni-app` | runtime | 构建通过 | 构建通过 | alpha 版本，正式发版前建议锁定并复测 |
| `@dcloudio/uni-components` | runtime | 构建通过 | 构建通过 | 使用 uni 组件体系 |
| `@dcloudio/uni-h5` | runtime | 构建通过 | N/A | 仅 H5 端 |
| `@dcloudio/uni-mp-weixin` | runtime | N/A | 构建通过 | 需微信开发者工具补运行时验证 |
| `vue` | runtime | 构建通过 | 构建通过 | 无额外平台限制 |
| `@dcloudio/vite-plugin-uni` | build | 构建通过 | 构建通过 | 不进入运行时包 |
| `vite` | build | 构建通过 | 构建通过 | 不进入运行时包 |
| `typescript`、`vue-tsc` | dev | 类型检查通过 | N/A | 不进入运行时包 |
| `vitest` | dev | 测试通过 | N/A | 不进入运行时包 |

未新增二维码、canvas、地图、支付或 DOM 依赖库；平台禁用 API 扫描通过。
