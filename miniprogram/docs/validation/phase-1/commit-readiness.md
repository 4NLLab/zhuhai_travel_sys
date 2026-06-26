# Phase 1 提交准备检查

采集时间：2026-06-26

## 工作区状态

`git status --short --ignored=matching miniprogram`：

```text
?? miniprogram/
!! miniprogram/dist/
!! miniprogram/node_modules/
```

结论：

- 本阶段变更只新增 `miniprogram/`。
- `node_modules/` 和 `dist/` 已被仓库 `.gitignore` 忽略，不属于提交候选。
- 当前未执行 `git add`、`git commit`、push 或 PR 操作。

## 提交候选文件

命令：`git ls-files --others --exclude-standard`

- 候选文件数：49
- 候选文件总体积：约 640K
- 类型分布：17 个 `.ts`、10 个 `.md`、6 个 `.vue`、6 个 `.png`、4 个 `.json`、2 个 `.mjs`、1 个 `.yaml`、1 个 `.html`、1 个 `.example`、1 个 `.css`

候选范围：

- 工程与构建配置：`package.json`、`pnpm-lock.yaml`、`tsconfig.json`、`vite.config.ts`、`vitest.config.ts`、`src/pages.json`、`src/manifest.json`
- 源码：`src/pages`、`src/subpackages`、`src/api`、`src/adapters`、`src/mock`、`src/utils`、`src/types`、`src/components`
- 脚本：`scripts/lint-platform.mjs`、`scripts/check-size.mjs`
- 测试：`tests/*.test.ts`
- 文档和证据：`README.md`、`docs/contracts`、`docs/validation/phase-1`

## Diff 策略说明

由于 `miniprogram/` 是全新未跟踪目录，普通 `git diff` 不会展示内容。使用空目录与 `miniprogram/` 做 `git diff --no-index` 会把被 `.gitignore` 忽略的 `node_modules/` 和 `dist/` 也纳入统计，产生 2 万多个文件的噪声。因此提交准备检查以 `git ls-files --others --exclude-standard` 作为提交候选的权威范围。

若用户授权提交，提交前应先按 `.claude/rules/git-submit.mdc` 执行：

1. `git status`
2. 将上述 49 个候选文件显式加入暂存区
3. `git diff --cached --stat`
4. `git diff --cached`
5. 复核凭据、生成物、截图和 lockfile

## 凭据与敏感信息扫描

命令：

```bash
rg -n --hidden --glob '!miniprogram/node_modules/**' --glob '!miniprogram/dist/**' --glob '!miniprogram/pnpm-lock.yaml' --glob '!miniprogram/docs/validation/phase-1/screenshots/*.png' "(AKIA[0-9A-Z]{16}|-----BEGIN [A-Z ]*PRIVATE KEY-----|JWT_SECRET|PAYMENT_WEBHOOK_SECRET|password\\s*[:=]|secret\\s*[:=]|token\\s*[:=]|accessToken\\s*[:=]|api[_-]?key\\s*[:=])" miniprogram
```

命中项：

- `src/types/platform.ts` 的 `accessToken` 类型字段
- `src/adapters/storage.ts` 的最小 session 存储字段
- `tests/storage.test.ts` 的假 token `dev-token`
- `tests/sensitive.test.ts` 的脱敏测试假 token

结论：未发现真实凭据、私钥、生产 token 或 `.env` 实值文件。`.env.example` 是模板文件，符合仓库规则。

## 生成物与资源扫描

命令：

```bash
git ls-files --others --exclude-standard | rg '(^|/)(node_modules|dist)/|\\.env$|\\.env\\.'
git ls-files --others --exclude-standard | rg '(^|/)(preview-|qa-)|\\.(mp4|mov|avi)$'
find miniprogram -path '*/node_modules/*' -prune -o -path '*/dist/*' -prune -o -type f \\( -name 'preview-*' -o -name 'qa-*' -o -size +200k \\) -print
```

结果：

- 候选中没有 `node_modules/` 或 `dist/`。
- 候选中没有真实 `.env`。
- 候选中没有 `preview-*`、`qa-*` 或视频文件。
- 超过 200K 的候选仅 `pnpm-lock.yaml`，不是小程序静态资源。

## 影响范围

- 未修改根目录静态 HTML。
- 未修改 Go 后端。
- 未复制根目录 `assets/`。
- Phase 1 仅建立小程序工程和平台护栏，不改变现有线上/本地 Web 入口行为。

## 剩余门禁

- 缺少用户明确授权，不能执行 commit。
- 未发现可用微信开发者工具 CLI 或 `miniprogram-ci`，小程序运行时验收仍是 blocker；当前只能声明 `mp-weixin` build-validated。
