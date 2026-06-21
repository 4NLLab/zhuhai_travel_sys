# 珠海文旅后端架构文档

## 技术栈

| 层 | 技术 | 版本 |
|---|---|---|
| 语言 | Go | 1.25.10 |
| Web 框架 | Gin | v1.12.0 |
| ORM | GORM | v1.31.1 |
| 数据库 | MariaDB | 10.5.29 |
| 数据库驱动 | go-sql-driver/mysql | v1.8.1 |

## 目录结构

```
wwwroot/backend/
├── main.go              # 应用入口
├── .env                 # 环境变量配置
├── .gitignore
├── go.mod / go.sum      # Go 模块依赖
├── server               # 编译产物（静态二进制）
│
├── config/              # 配置加载
│   └── config.go        # 读取环境变量，生成 Config 实例
│
├── database/            # 数据库层
│   └── mysql.go         # GORM 连接初始化，连接池配置
│
├── models/              # 数据模型（GORM 映射）
│   ├── user.go          # 用户/出游人/管理员/产品/SKU/排期/分类/轮播图
│   ├── order.go         # 订单/订单项/出游人/支付/票/核销/退款
│   └── driver.go        # 司机/车辆/二维码/佣金/发票/收藏/工单/审计
│
├── handlers/            # HTTP 处理器（Controller 层）
│   ├── auth.go          # 管理员登录
│   ├── user.go          # 用户/实名/出游人/收藏/发票
│   ├── product.go       # 产品列表(搜索+分页)/详情/排期/分类
│   ├── order.go         # 下单(防超卖事务)/支付回调/订单查询
│   ├── ticket.go        # 票详情/二维码验票/核销
│   ├── driver.go        # 司机/佣金
│   └── admin.go         # 看板/订单管理/轮播图CRUD/佣金结算/参数
│
├── middleware/           # 中间件
│   └── middleware.go     # 请求日志 + CORS 跨域
│
├── routes/              # 路由注册
│   └── router.go        # 统一路由入口
│
├── dto/                 # 数据传输对象
│   └── response.go      # 统一 API 响应格式
│
├── services/            # 业务逻辑层（预留扩展）
├── utils/               # 工具函数（预留扩展）
└── docs/                # 项目文档
    ├── ARCHITECTURE.md  # 本文档
    └── DATABASE.md      # 数据库设计文档
```

## 核心设计

### 配置管理

配置通过环境变量注入（[config/config.go](file:///www/wwwroot/backend/config/config.go)），默认值定义于 `.env` 文件：

| 变量 | 默认值 | 说明 |
|---|---|---|
| `SERVER_PORT` | `8080` | 服务监听端口 |
| `DB_HOST` | `127.0.0.1` | 数据库地址 |
| `DB_PORT` | `3306` | 数据库端口 |
| `DB_USER` | `root` | 数据库用户 |
| `DB_PASSWORD` | `wuyuanjian0` | 数据库密码 |
| `DB_NAME` | `zhuhai_travel` | 数据库名 |
| `JWT_SECRET` | - | JWT 签名密钥 |

### 数据库连接池

```go
sqlDB.SetMaxIdleConns(10)     // 最大空闲连接
sqlDB.SetMaxOpenConns(50)     // 最大打开连接
sqlDB.SetConnMaxLifetime(time.Hour)  // 连接最大存活时间
```

### 库存防超卖

下单时通过 **GORM 事务 + 行级锁** 保证库存安全：

```
1. 查询排期可用库存 = total_stock - locked_stock - sold_stock
2. 库存不足 → 返回错误
3. 开启事务 → UPDATE locked_stock += quantity
4. 创建 order / order_item / order_traveler
5. 提交事务
```

### 统一响应格式

```json
// 成功
{ "code": 200, "message": "ok", "data": {...} }

// 失败
{ "code": 400, "message": "参数错误" }

// 分页
{ "code": 200, "message": "ok", "data": [...], "total": 100, "page": 1, "size": 10 }
```

## API 接口列表（共 40 个）

### 系统

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/health` | 健康检查 |

### 用户端

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| GET | `/api/v1/users/profile?user_id=` | 用户信息 | `UserProfile` |
| POST | `/api/v1/users/realname` | 提交实名认证 | `UserRealnameSubmit` |

### 出游人

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| GET | `/api/v1/travelers?user_id=` | 列表 | `TravelerList` |
| POST | `/api/v1/travelers` | 新增 | `TravelerCreate` |
| PUT | `/api/v1/travelers/:id` | 编辑 | `TravelerUpdate` |
| DELETE | `/api/v1/travelers/:id` | 删除 | `TravelerDelete` |

### 收藏

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| GET | `/api/v1/favorites?user_id=` | 收藏列表 | `FavoriteList` |
| POST | `/api/v1/favorites/toggle` | 切换收藏 | `FavoriteToggle` |

### 发票

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| GET | `/api/v1/invoice-titles?user_id=` | 发票抬头列表 | `InvoiceTitleList` |
| POST | `/api/v1/invoice-titles` | 新增抬头 | `InvoiceTitleCreate` |
| DELETE | `/api/v1/invoice-titles/:id` | 删除抬头 | `InvoiceTitleDelete` |
| POST | `/api/v1/invoices` | 申请开票 | `InvoiceCreate` |

### 产品模块

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| GET | `/api/v1/products?keyword=&category_id=&type=&page=&size=` | 产品列表（搜索+分类+分页） | `ProductList` |
| GET | `/api/v1/products/schedules/query?product_id=&date=` | 排期查询 | `ProductSchedulesQuery` |
| GET | `/api/v1/products/:id` | 产品详情（SKU/排期/图片） | `ProductDetail` |
| GET | `/api/v1/categories` | 分类树（含子分类） | `CategoryList` |
| GET | `/api/v1/banners` | 轮播图（仅启用） | `BannerList` |

### 订单模块

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| POST | `/api/v1/orders` | 创建订单（库存锁+防超卖+司机扫码） | `CreateOrder` |
| GET | `/api/v1/orders?user_id=&status=&page=&size=` | 用户订单列表 | `OrderList` |
| GET | `/api/v1/orders/:id` | 订单详情（子项/票/出游人/支付） | `OrderDetail` |
| POST | `/api/v1/payments/callback` | 支付回调（生成票+库存转移+佣金计算） | `PaymentCallback` |

### 票务模块

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| GET | `/api/v1/tickets/:id` | 票详情 | `TicketDetail` |
| GET | `/api/v1/tickets/by-qr?qr_hash=` | 二维码验票查询 | `TicketByQR` |
| POST | `/api/v1/tickets/verify` | 核销电子票 | `TicketVerify` |
| GET | `/api/v1/tickets/verifications?ticket_id=` | 核销历史 | `VerificationHistory` |

### 司机模块

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| GET | `/api/v1/drivers` | 司机列表 | `DriverList` |
| GET | `/api/v1/commissions?driver_id=` | 佣金记录 | `DriverCommissionList` |
| GET | `/api/v1/commissions/summary` | 佣金汇总 | `CommissionSummary` |

### 管理后台

| 方法 | 路径 | 说明 | Handler |
|---|---|---|---|
| POST | `/api/v1/admin/login` | 登录 | `AdminLogin` |
| GET | `/api/v1/admin/dashboard` | 运营看板 | `AdminDashboard` |
| GET | `/api/v1/admin/trend` | 7 日交易趋势 | `AdminTrend` |
| GET | `/api/v1/admin/orders?status=&page=&size=` | 订单管理 | `AdminOrderList` |
| POST | `/api/v1/admin/orders/refund` | 处理退款 | `AdminOrderRefund` |
| GET | `/api/v1/admin/banners` | 轮播图全列表 | `AdminBannerList` |
| POST | `/api/v1/admin/banners` | 新增轮播图 | `AdminBannerCreate` |
| PUT | `/api/v1/admin/banners/:id` | 编辑轮播图 | `AdminBannerUpdate` |
| DELETE | `/api/v1/admin/banners/:id` | 删除轮播图 | `AdminBannerDelete` |
| GET | `/api/v1/admin/commission-batches` | 佣金批次 | `AdminCommissionBatches` |
| POST | `/api/v1/admin/commissions/settle` | 手动结算佣金 | `AdminCommissionSettle` |
| GET | `/api/v1/admin/params` | 系统参数配置 | `AdminParams` |

## 启动方式

```bash
# 编译
cd /www/wwwroot/backend
GOPROXY=https://goproxy.cn,direct go build -o server .

# 前台运行
./server

# 后台运行
nohup ./server > server.log 2>&1 &
```

## 扩展指南

- 新增业务接口 → 在 `handlers/` 下新建对应文件，在 `routes/router.go` 注册路由
- 新增中间件（如 JWT 鉴权）→ 在 `middleware/` 下编写，在 `router.go` 中 `Use()`
- 新增数据表 → 在 `models/` 下新增结构体，执行 SQL 或 AutoMigrate
- 复杂业务逻辑 → 提取到 `services/` 层，handler 调用 service
- 微信支付集成 → 引入 [gopay](https://github.com/go-pay/gopay)，在 `services/` 下封装
