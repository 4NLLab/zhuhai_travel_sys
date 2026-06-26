# Phase 6 平台禁用 API 扫描

日期：2026-06-26

命令：`corepack pnpm -C miniprogram lint:platform`

结果：通过。

扫描范围：

- `src/pages`
- `src/components`
- `src/api`
- `src/adapters`
- `src/stores`
- `src/utils`
- `dist/build/mp-weixin`
- `dist/mp-weixin`

规则覆盖：`window`、`document`、`localStorage`、`fetch`、`history`、外链 `<script src="http(s)://...">`。

说明：正则人工复核时出现的 `<script setup>` 是 Vue 单文件组件语法，不是外链脚本。
