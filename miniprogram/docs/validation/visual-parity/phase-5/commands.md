# Phase 5 最终命令记录

日期：2026-06-26

## 服务确认

- 原版静态页：`http://127.0.0.1:8000/`
- 迁移版 H5：`http://127.0.0.1:5174/`

## 截图结果

- Phase 5 共保存 `28` 张最终截图。
- 原版截图：7 页 × 2 视口。
- 迁移版截图：7 页 × 2 视口。

## 自动化验证

| 命令 | 结果 |
|---|---|
| `corepack pnpm -C miniprogram typecheck` | 通过 |
| `corepack pnpm -C miniprogram test` | 通过，7 个测试文件、20 个测试用例通过 |
| `corepack pnpm -C miniprogram build:h5` | 通过 |
| `corepack pnpm -C miniprogram build:mp-weixin` | 通过 |
| `corepack pnpm -C miniprogram lint:platform` | 通过 |
| `corepack pnpm -C miniprogram check:size` | 通过，主包 `1815.7 KB`、司机分包 `16.3 KB`、总包 `1832.0 KB` |

## mp-weixin 构建产物

- pages：
  - `pages/home/index`
  - `pages/orders/index`
  - `pages/profile/index`
  - `pages/ticket/index`
  - `pages/island-cruise/index`
- tabBar：
  - `pages/home/index`：首页
  - `pages/orders/index`：客服
  - `pages/profile/index`：我的
- subPackages：
  - root `subpackages/driver`
  - page `pages/home/index`
- 包体：
  - mainPackage：`1815.7 KB`
  - `subpackages/driver`：`16.3 KB`
  - total：`1832.0 KB`
