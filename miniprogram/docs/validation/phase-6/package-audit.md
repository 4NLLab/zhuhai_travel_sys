# Phase 6 包体与资源复核

日期：2026-06-26

## 首版静态资源

| 文件 | 大小 | 用途 |
|---|---:|---|
| `src/static/phase2/home-hero-moon-web.jpg` | 84295 bytes | 首页轮播 |
| `src/static/phase2/home-hero-sunset-web.jpg` | 82279 bytes | 首页轮播 |
| `src/static/phase2/home-hero-watercolor.png` | 875061 bytes | 首页轮播主视觉 |
| `src/static/phase2/macau-cruise-night-banner-web.jpg` | 116351 bytes | 首页/环岛游 hero |
| `src/static/phase2/taxi-scan-illustration-web.jpg` | 118019 bytes | 首页接送/司机引导 |
| `src/static/phase2/ticket-wallet-illustration-web.jpg` | 30599 bytes | 首页套餐/票券钱包 |
| `src/static/phase2/zhuhai-bay-home-hero-web.jpg` | 214282 bytes | 首页首版视觉备选资源 |
| `src/static/phase3/verify-hero-wide-clean-web.png` | 329779 bytes | 环岛游验证页首版视觉资源 |

未发现 `preview-*`、`qa-*` 或非首版大图进入 `src/static`。

超过 200 KB 的图片仍处于首版页面视觉范围内；其中 `home-hero-watercolor.png` 是当前最大资源，后续若继续新增图片，应优先远程化或压缩该首屏资源，避免主包越过 2 MB 硬门禁。

## mp-weixin 包体

最终 `check:size`：

```json
{
  "distRoot": "dist/build/mp-weixin",
  "mainPackage": "1986.7 KB",
  "subPackages": {
    "subpackages/driver": "16.3 KB"
  },
  "total": "2003.0 KB",
  "largestFile": {
    "path": "static/phase2/home-hero-watercolor.png",
    "size": "854.6 KB"
  },
  "largeAssets": [
    {
      "path": "static/phase2/home-hero-watercolor.png",
      "size": "854.6 KB"
    },
    {
      "path": "static/phase2/zhuhai-bay-home-hero-web.jpg",
      "size": "209.3 KB"
    },
    {
      "path": "static/phase3/verify-hero-wide-clean-web.png",
      "size": "322.0 KB"
    }
  ]
}
```

结论：主包、分包和总包均低于门禁，但主包已达到 1.8 MB 预警线；首版资源归属清晰，继续新增本地大图前必须先压缩或远程化现有首屏资源。
