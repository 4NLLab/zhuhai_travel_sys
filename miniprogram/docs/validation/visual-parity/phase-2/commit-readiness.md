# Phase 2 提交准备检查

日期：2026-06-26

## 变更范围

- 首页视觉重构：`miniprogram/src/pages/home/index.vue`
- H5 首页导航与端口：`miniprogram/src/pages.json`、`miniprogram/src/manifest.json`
- H5 原生 tabBar 隐藏：`miniprogram/src/App.vue`
- 首页 Mock 内容：`miniprogram/src/mock/main-shell.ts`
- 首页视觉资产：`miniprogram/src/static/phase2/home-hero-watercolor.png`、`zhuhai-bay-home-hero-web.jpg`
- Phase 2 验证证据：`miniprogram/docs/validation/visual-parity/phase-2/`

## 检查结论

- 凭据检查：未新增 `.env`、密钥、证书或账号 JSON。
- 生成物检查：新增 PNG 均为本 Phase 必需截图证据；`.local/` 为本机 Playwright 运行库，已写入 `.git/info/exclude`，不会提交。
- 无关文件检查：未修改后端、旧静态 HTML、Docker 或数据库。
- 视觉结论：`通过：迁移版与原版截图肉眼无明显区别`。
