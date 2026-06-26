# Phase 6 包体与资源复核

日期：2026-06-26

## 首版静态资源

| 文件 | 大小 | 用途 |
|---|---:|---|
| `src/static/phase2/macau-cruise-night-banner-web.jpg` | 116351 bytes | 首页/环岛游 hero |
| `src/static/phase2/taxi-scan-illustration-web.jpg` | 118019 bytes | 首页接送/司机引导 |
| `src/static/phase2/ticket-wallet-illustration-web.jpg` | 30599 bytes | 首页套餐/票券钱包 |

未发现 `preview-*`、`qa-*` 或非首版大图进入 `src/static`。

## mp-weixin 包体

最终 `check:size`：

```json
{
  "distRoot": "dist/build/mp-weixin",
  "mainPackage": "406.9 KB",
  "subPackages": {
    "subpackages/driver": "14.6 KB"
  },
  "total": "421.5 KB",
  "largestFile": {
    "path": "static/phase2/taxi-scan-illustration-web.jpg",
    "size": "115.3 KB"
  },
  "largeAssets": []
}
```

结论：主包、分包和总包均低于门禁；首版资源归属清晰。
