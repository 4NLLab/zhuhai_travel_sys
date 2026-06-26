# Phase 3 Blockers

- 未发现可用微信开发者工具 CLI 或 `miniprogram-ci`。本阶段只能声明 `mp-weixin` build-validated，不能声明小程序运行时验收完成。
- Playwright Chromium CLI 启动缺 `libnspr4.so`，`npx playwright install-deps chromium` 需要 sudo 交互认证；本阶段继续使用 Firefox CLI 完成 H5 截图。
- 真实供应商账号、锁票/出票测试订单、微信支付测试能力、退改签 seed 尚未提供；Phase 3 只验证 mock 闭环和接口契约，真实接口联调留到 Phase 5。
