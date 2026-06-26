# Phase 5 mp-weixin 包体记录

日期：2026-06-26

`corepack pnpm -C miniprogram check:size` 输出：

```json
{
  "distRoot": "dist/build/mp-weixin",
  "mainPackage": "407.0 KB",
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

结论：主包、司机分包和总包体积均低于 Phase 1 预算；未发现大体积首版无关资源。
