# Phase 1 Blockers

- 未发现可用微信开发者工具 CLI 或 `miniprogram-ci`。本阶段只能声明 `mp-weixin` build-validated，不能声明小程序运行时验收完成。
- Python Playwright 库未安装；已改用 Playwright Firefox CLI 完成 H5 截图。
- Playwright Chromium CLI 启动缺 `libnspr4.so`，`npx playwright install-deps chromium` 需要 sudo 交互认证；已改用 Firefox 截图，不影响 H5 截图证据。
