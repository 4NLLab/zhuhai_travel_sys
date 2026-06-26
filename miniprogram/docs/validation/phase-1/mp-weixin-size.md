# Phase 1 mp-weixin 构建与包体

采集时间：2026-06-26

## 构建

- 命令：`corepack pnpm -C miniprogram build:mp-weixin`
- 结果：通过
- 输出：`DONE Build complete.`
- 产物：`miniprogram/dist/build/mp-weixin`

## app.json 检查

`dist/build/mp-weixin/app.json` 已包含：

- `pages`: `pages/home/index`、`pages/orders/index`、`pages/profile/index`
- `tabBar`: 首页、订单、我的
- `subPackages`: `subpackages/driver/pages/home/index`

## 包体

`corepack pnpm -C miniprogram check:size`：

```json
{
  "distRoot": "dist/build/mp-weixin",
  "mainPackage": "78.7 KB",
  "subPackages": {
    "subpackages/driver": "0.9 KB"
  },
  "total": "79.6 KB",
  "largestFile": {
    "path": "common/vendor.js",
    "size": "64.7 KB"
  },
  "largeAssets": []
}
```

## 资源结论

- 未复制根目录 `assets/`。
- 无单文件超过 200 KB。
- 主包、司机分包、总包均低于 Phase 1 预算。
