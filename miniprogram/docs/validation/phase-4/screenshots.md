# Phase 4 截图记录

采集时间：2026-06-26

| 页面 / 状态 | 桌面截图 | 移动截图 | 结果 |
|---|---|---|---|
| active 司机工作台 | `screenshots/driver-active-desktop.png` | `screenshots/driver-active-mobile.png` | 通过；资料、推广码、钱包、佣金、提现和记录可读 |
| 待审核司机 | N/A | `screenshots/driver-pending-mobile.png` | 通过；待审核状态和推广码生成说明可读 |
| 空钱包 / 空记录 | N/A | `screenshots/driver-empty-mobile.png` | 通过；钱包、佣金、提现记录空态可解释 |
| 余额不足提现 | N/A | `screenshots/driver-insufficient-mobile.png` | 通过；提现表单可见，余额不足场景由 mock 覆盖 |

## 截图路由

- `/#/subpackages/driver/pages/home/index?scenario=phase4-active`
- `/#/subpackages/driver/pages/home/index?scenario=phase4-pending-review`
- `/#/subpackages/driver/pages/home/index?scenario=phase4-empty-wallet`
- `/#/subpackages/driver/pages/home/index?scenario=phase4-insufficient-balance`
