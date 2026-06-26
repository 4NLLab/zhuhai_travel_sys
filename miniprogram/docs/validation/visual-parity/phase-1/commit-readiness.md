# Phase 1 提交准备检查

日期：2026-06-26

## 变更范围

- 新增 `miniprogram/docs/validation/visual-parity/phase-1/` 视觉基线证据。
- 新增 `miniprogram/scripts/capture-visual-parity.mjs` 作为后续可复跑截图脚本。
- 新增 `playwright` devDependency，用于本地截图脚本。

## 检查结论

- 凭据检查：未新增 `.env`、密钥、证书、service account JSON 或私密配置。
- 生成物检查：新增 PNG 均为本 Phase 必需的截图证据。
- 无关文件检查：未修改原静态 HTML、uni-app 页面实现、后端、Docker 或数据库。
- 风险边界：本阶段不修 UI，只建立基线；所有页面当前视觉结论均为“不通过”。
