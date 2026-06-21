# 珠海文旅数据库设计文档

> 数据库：`zhuhai_travel` | 引擎：MariaDB 10.5 | 字符集：`utf8mb4_general_ci` | 表数：25

---

## 一、ER 关系总览

```
users ──< travelers
  │
  ├──< orders ──< order_items ──< tickets ──< ticket_verifications
  │       │           │
  │       │           ├──< order_travelers ──< travelers (ref)
  │       │           └──< driver_commissions ──< drivers
  │       │
  │       ├──< payments
  │       ├──< refunds ──< order_items (opt)
  │       ├──< invoices ──< invoice_titles
  │       ├──< driver_commissions
  │       └──< support_tickets
  │
  ├──< user_favorites ──< products
  └──< audit_logs (actor_id)

drivers ──< vehicles
  └──< driver_qr_codes ──< vehicles (opt)
       └──< orders (driver_qr_code_id)

product_categories ──< product_categories (parent_id)
  └──< products ──< product_images
       └──< product_skus ──< product_schedules
```

---

## 二、表结构详解

### 2.1 用户端（4 张表）

#### `users` — 用户

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 用户 ID |
| `openid` | VARCHAR(80) | UNIQUE, NULLABLE | 微信 OpenID |
| `unionid` | VARCHAR(80) | UNIQUE, NULLABLE | 微信 UnionID |
| `phone` | VARCHAR(32) | INDEX | 手机号 |
| `nickname` | VARCHAR(80) | | 昵称 |
| `avatar_url` | VARCHAR(500) | | 头像 URL |
| `real_name` | VARCHAR(80) | | 真实姓名 |
| `id_card_no` | VARBINARY(255) | | ⚠ 身份证号（加密存储） |
| `realname_status` | VARCHAR(24) | DEFAULT 'unverified' | 实名状态：unverified/verified |
| `status` | VARCHAR(24) | INDEX, 'active' | 账号状态 |
| `last_login_at` | DATETIME | | 最后登录 |
| `created_at` | DATETIME | CURRENT_TIMESTAMP | |
| `updated_at` | DATETIME | ON UPDATE | |

**索引**：`uk_users_openid`, `uk_users_unionid`, `idx_users_phone`, `idx_users_status`

---

#### `travelers` — 出游人

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `user_id` | BIGINT UNSIGNED | FK→users(id), INDEX | 所属用户 |
| `name` | VARCHAR(80) | NOT NULL | 姓名 |
| `phone` | VARCHAR(32) | | 手机 |
| `id_type` | VARCHAR(24) | DEFAULT 'id_card' | 证件类型 |
| `id_no` | VARBINARY(255) | NOT NULL | ⚠ 证件号（加密存储） |
| `is_default` | TINYINT(1) | DEFAULT 0 | 是否默认出游人 |

---

#### `user_favorites` — 用户收藏

| 字段 | 类型 | 约束 |
|---|---|---|
| `id` | BIGINT UNSIGNED | PK |
| `user_id` | BIGINT UNSIGNED | FK→users(id), UNIQUE 联合 |
| `product_id` | BIGINT UNSIGNED | FK→products(id), UNIQUE 联合, INDEX |

**唯一约束**：`uk_favorites_user_product (user_id, product_id)`

---

#### `support_tickets` — 客服工单

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `user_id` | BIGINT UNSIGNED | FK→users(id), INDEX | |
| `order_id` | BIGINT UNSIGNED | FK→orders(id), INDEX, NULLABLE | 关联订单 |
| `type` | VARCHAR(32) | NOT NULL | 工单类型 |
| `title` | VARCHAR(160) | NOT NULL | |
| `content` | TEXT | | |
| `status` | VARCHAR(32) | INDEX, DEFAULT 'open' | |
| `handled_by` | BIGINT UNSIGNED | FK→admin_users(id), NULLABLE | 处理人 |

---

### 2.2 订单模块（5 张表）

#### `orders` — 订单

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `order_no` | VARCHAR(40) | UNIQUE, NOT NULL | 订单号（如 ZH0618-1028） |
| `user_id` | BIGINT UNSIGNED | FK→users(id), INDEX | |
| `source` | VARCHAR(32) | DEFAULT 'miniapp' | 来源：miniapp / driver_qr |
| `driver_id` | BIGINT UNSIGNED | FK→drivers(id), INDEX, NULLABLE | 推广司机 |
| `driver_qr_code_id` | BIGINT UNSIGNED | FK→driver_qr_codes(id), NULLABLE | 扫码二维码 |
| `status` | VARCHAR(32) | INDEX, 'pending_payment' | pending_payment/paid/cancelled/completed |
| `total_amount` | DECIMAL(10,2) | DEFAULT 0 | 原价 |
| `discount_amount` | DECIMAL(10,2) | DEFAULT 0 | 优惠金额 |
| `payable_amount` | DECIMAL(10,2) | DEFAULT 0 | 应付金额 |
| `paid_amount` | DECIMAL(10,2) | DEFAULT 0 | 实付金额 |
| `contact_name` | VARCHAR(80) | | 联系人 |
| `contact_phone` | VARCHAR(32) | | 联系手机 |
| `remark` | VARCHAR(255) | | 备注 |
| `paid_at` | DATETIME | | 支付时间 |
| `cancelled_at` | DATETIME | | 取消时间 |

**索引**：`uk_orders_no`, `idx_orders_user_created(user_id, created_at)`, `idx_orders_status`, `idx_orders_driver`

---

#### `order_items` — 订单项

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `order_id` | BIGINT UNSIGNED | FK→orders(id), INDEX | |
| `product_id` | BIGINT UNSIGNED | FK→products(id), INDEX | |
| `sku_id` | BIGINT UNSIGNED | FK→product_skus(id) | |
| `schedule_id` | BIGINT UNSIGNED | FK→product_schedules(id), INDEX, NULLABLE | 关联排期 |
| `product_title` | VARCHAR(160) | NOT NULL | 产品标题快照 |
| `sku_name` | VARCHAR(120) | NOT NULL | SKU 名称快照 |
| `travel_date` | DATE | | 出行日期 |
| `start_time` | TIME | | 出发时间 |
| `quantity` | INT UNSIGNED | NOT NULL | 数量 |
| `unit_price` | DECIMAL(10,2) | NOT NULL | 单价 |
| `total_price` | DECIMAL(10,2) | NOT NULL | 小计 |
| `status` | VARCHAR(32) | DEFAULT 'active' | |

> 订单项存储下单时的产品和 SKU 信息快照，防止后续产品变更影响已产生的订单。

---

#### `order_travelers` — 订单出游人

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `order_item_id` | BIGINT UNSIGNED | FK→order_items(id), INDEX | |
| `traveler_id` | BIGINT UNSIGNED | FK→travelers(id), INDEX, NULLABLE | 关联已保存的出游人 |
| `name` | VARCHAR(80) | NOT NULL | |
| `phone` | VARCHAR(32) | | |
| `id_type` | VARCHAR(24) | DEFAULT 'id_card' | |
| `id_no` | VARBINARY(255) | NOT NULL | ⚠ 证件号加密 |

---

#### `payments` — 支付记录

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `order_id` | BIGINT UNSIGNED | FK→orders(id), INDEX, UNIQUE | 一个订单一条支付记录 |
| `payment_no` | VARCHAR(64) | UNIQUE, NOT NULL | 支付流水号 |
| `channel` | VARCHAR(32) | DEFAULT 'wechat' | 支付渠道 |
| `amount` | DECIMAL(10,2) | NOT NULL | 支付金额 |
| `status` | VARCHAR(32) | INDEX, 'pending' | pending/success/failed/refunded |
| `transaction_id` | VARCHAR(128) | | 微信/支付宝交易号 |
| `paid_at` | DATETIME | | 实际付款时间 |
| `raw_payload` | JSON | | ⚠ 支付回调原始数据（JSON 校验） |

---

#### `refunds` — 退款记录

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `order_id` | BIGINT UNSIGNED | FK→orders(id), INDEX | |
| `order_item_id` | BIGINT UNSIGNED | FK→order_items(id), NULLABLE | 部分退款时指定 |
| `refund_no` | VARCHAR(64) | UNIQUE, NOT NULL | 退款单号 |
| `amount` | DECIMAL(10,2) | NOT NULL | 退款金额 |
| `reason` | VARCHAR(255) | | 退款原因 |
| `status` | VARCHAR(32) | INDEX, 'requested' | requested/approved/rejected/refunded |
| `requested_at` | DATETIME | CURRENT_TIMESTAMP | 申请时间 |
| `processed_at` | DATETIME | | 处理时间 |
| `raw_payload` | JSON | | 退款回调原始数据 |

---

### 2.3 票务模块（2 张表）

#### `tickets` — 电子票

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `order_item_id` | BIGINT UNSIGNED | FK→order_items(id), INDEX | |
| `ticket_no` | VARCHAR(64) | UNIQUE, NOT NULL | 票号 |
| `qr_token_hash` | CHAR(64) | UNIQUE, NOT NULL | ⚠ 二维码哈希（单向加密） |
| `status` | VARCHAR(32) | INDEX, DEFAULT 'valid' | valid/used/refunded/cancelled |
| `valid_from` | DATETIME | | 有效期始 |
| `valid_to` | DATETIME | | 有效期止 |
| `used_at` | DATETIME | | 核销时间 |

> 每个 order_item 按 quantity 生成对应数量的 tickets。

---

#### `ticket_verifications` — 核销记录

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `ticket_id` | BIGINT UNSIGNED | FK→tickets(id), INDEX | |
| `verifier_admin_id` | BIGINT UNSIGNED | FK→admin_users(id), NULLABLE | 核销员 |
| `verify_location` | VARCHAR(120) | | 核销地点 |
| `verify_result` | VARCHAR(32) | NOT NULL | success/fail |
| `verify_note` | VARCHAR(255) | | 备注 |
| `created_at` | DATETIME | INDEX | 核销时间 |

---

### 2.4 产品模块（6 张表）

#### `product_categories` — 产品分类

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `parent_id` | BIGINT UNSIGNED | FK→product_categories(id), INDEX | ⚠ 自引用，支持无限层级 |
| `name` | VARCHAR(80) | NOT NULL | |
| `slug` | VARCHAR(80) | UNIQUE, NOT NULL | URL 友好标识 |
| `icon` | VARCHAR(80) | | 图标 |
| `sort_order` | INT | DEFAULT 0 | 排序 |
| `status` | VARCHAR(24) | DEFAULT 'active' | |

---

#### `products` — 产品

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `category_id` | BIGINT UNSIGNED | FK→product_categories(id), INDEX | |
| `title` | VARCHAR(160) | NOT NULL | 产品标题 |
| `subtitle` | VARCHAR(255) | | 副标题 |
| `product_type` | VARCHAR(32) | NOT NULL | 类型：cruise/hotel/scenic/entertainment |
| `cover_url` | VARCHAR(500) | | 封面图 |
| `location_name` | VARCHAR(120) | | 地点名 |
| `address` | VARCHAR(255) | | 地址 |
| `notice` | TEXT | | 注意事项 |
| `refund_policy` | TEXT | | 退款政策 |
| `status` | VARCHAR(24) | DEFAULT 'draft' | draft/active/offline |
| `sort_order` | INT | DEFAULT 0 | |

**索引**：`idx_products_type_status(product_type, status)`

---

#### `product_images` — 产品图片

| 字段 | 类型 | 约束 |
|---|---|---|
| `id` | BIGINT UNSIGNED | PK |
| `product_id` | BIGINT UNSIGNED | FK→products(id), INDEX |
| `image_url` | VARCHAR(500) | NOT NULL |
| `alt_text` | VARCHAR(160) | |
| `sort_order` | INT | DEFAULT 0 |

---

#### `product_skus` — SKU（价格规格）

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `product_id` | BIGINT UNSIGNED | FK→products(id), INDEX | |
| `sku_name` | VARCHAR(120) | NOT NULL | 如"一等座"/"儿童票" |
| `market_price` | DECIMAL(10,2) | | 市场价 |
| `sale_price` | DECIMAL(10,2) | NOT NULL | 售价 |
| `settlement_price` | DECIMAL(10,2) | | 结算价（给渠道） |
| `stock_mode` | VARCHAR(24) | DEFAULT 'schedule' | schedule/unlimited |
| `status` | VARCHAR(24) | INDEX, DEFAULT 'active' | |

---

#### `product_schedules` — 排期/库存

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `product_id` | BIGINT UNSIGNED | FK→products(id), INDEX | |
| `sku_id` | BIGINT UNSIGNED | FK→product_skus(id), UNIQUE 联合 | |
| `travel_date` | DATE | NOT NULL, UNIQUE 联合 | 出行日期 |
| `start_time` | TIME | UNIQUE 联合 | 场次时间 |
| `end_time` | TIME | | |
| `venue` | VARCHAR(120) | | 场馆/码头 |
| `total_stock` | INT | DEFAULT 0 | 总库存 |
| `locked_stock` | INT | DEFAULT 0 | ⚠ 锁定库存（下单未支付） |
| `sold_stock` | INT | DEFAULT 0 | 已售库存（已支付） |
| `status` | VARCHAR(24) | DEFAULT 'active' | |

**唯一约束**：`uk_schedule_sku_time (sku_id, travel_date, start_time)`

> ⚠ **库存防超卖公式**：可用库存 = `total_stock - locked_stock - sold_stock`
> 下单时先对 locked_stock +N，支付成功后再移动到 sold_stock。

**索引**：`idx_schedules_product_date(product_id, travel_date)`

---

#### `banners` — 首页轮播图

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `title` | VARCHAR(120) | NOT NULL | |
| `subtitle` | VARCHAR(255) | | |
| `image_url` | VARCHAR(500) | NOT NULL | |
| `link_type` | VARCHAR(32) | DEFAULT 'product' | product/category/url |
| `link_target` | VARCHAR(160) | | 跳转目标 ID 或 URL |
| `sort_order` | INT | DEFAULT 0 | |
| `status` | VARCHAR(24) | DEFAULT 'active' | |
| `starts_at` | DATETIME | | 展示开始 |
| `ends_at` | DATETIME | | 展示结束 |

**索引**：`idx_banners_status_sort(status, sort_order)`

---

### 2.5 司机模块（4 张表）

#### `drivers` — 司机

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `driver_no` | VARCHAR(40) | UNIQUE, NOT NULL | 司机编号 |
| `name` | VARCHAR(80) | NOT NULL | |
| `phone` | VARCHAR(32) | NOT NULL, INDEX | |
| `id_card_no` | VARBINARY(255) | | ⚠ 身份证号加密 |
| `status` | VARCHAR(24) | INDEX, DEFAULT 'active' | |
| `commission_rate` | DECIMAL(6,4) | DEFAULT 0.0800 | ⚠ 提成比例（如 0.08 = 8%） |

---

#### `vehicles` — 车辆

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `driver_id` | BIGINT UNSIGNED | FK→drivers(id), INDEX | |
| `plate_no` | VARCHAR(32) | UNIQUE, NOT NULL | 车牌号 |
| `model` | VARCHAR(80) | | 车型 |
| `seats` | INT UNSIGNED | | 座位数 |
| `status` | VARCHAR(24) | DEFAULT 'active' | |

---

#### `driver_qr_codes` — 司机推广二维码

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `driver_id` | BIGINT UNSIGNED | FK→drivers(id), INDEX | |
| `vehicle_id` | BIGINT UNSIGNED | FK→vehicles(id), INDEX, NULLABLE | |
| `code` | VARCHAR(80) | UNIQUE, NOT NULL | 二维码标识码 |
| `scene` | VARCHAR(80) | DEFAULT 'seat' | 场景：seat/window/entrance |
| `status` | VARCHAR(24) | DEFAULT 'active' | |

> 乘客扫车上二维码 → 后端通过 code 定位司机和车辆 → 订单自动关联。

---

#### `driver_commissions` — 司机佣金

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `driver_id` | BIGINT UNSIGNED | FK→drivers(id), 联合索引 | |
| `order_id` | BIGINT UNSIGNED | FK→orders(id), INDEX | |
| `order_item_id` | BIGINT UNSIGNED | FK→order_items(id) | |
| `commission_no` | VARCHAR(64) | UNIQUE, NOT NULL | 佣金单号 |
| `base_amount` | DECIMAL(10,2) | NOT NULL | 计算基数 |
| `rate` | DECIMAL(6,4) | NOT NULL | 提成率 |
| `commission_amount` | DECIMAL(10,2) | NOT NULL | ⚠ 佣金金额 |
| `status` | VARCHAR(32) | DEFAULT 'pending' | pending/settled/cancelled |
| `settled_at` | DATETIME | | 结算时间 |

**索引**：`idx_commissions_driver_status(driver_id, status)`, `idx_commissions_order`

---

### 2.6 管理端（2 张表）

#### `admin_users` — 后台管理员

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `username` | VARCHAR(80) | UNIQUE, NOT NULL | 登录名 |
| `password_hash` | VARCHAR(255) | NOT NULL | ⚠ 密码哈希 |
| `display_name` | VARCHAR(80) | NOT NULL | 显示名 |
| `role` | VARCHAR(32) | DEFAULT 'operator' | super_admin/operator/viewer |
| `status` | VARCHAR(24) | DEFAULT 'active' | |
| `last_login_at` | DATETIME | | |

**索引**：`idx_admin_role_status(role, status)`

---

#### `audit_logs` — 操作审计

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `actor_type` | VARCHAR(24) | 联合索引, NOT NULL | admin/user/system |
| `actor_id` | BIGINT UNSIGNED | 联合索引, NULLABLE | |
| `action` | VARCHAR(80) | NOT NULL | 操作类型（如 order.refund） |
| `target_type` | VARCHAR(80) | 联合索引 | 目标类型 |
| `target_id` | BIGINT UNSIGNED | 联合索引 | 目标 ID |
| `ip` | VARCHAR(64) | | 操作 IP |
| `user_agent` | VARCHAR(500) | | |
| `payload` | JSON | | 操作详情 |
| `created_at` | DATETIME | INDEX | |

**索引**：`idx_audit_actor(actor_type, actor_id)`, `idx_audit_target(target_type, target_id)`, `idx_audit_created`

---

### 2.7 发票模块（2 张表）

#### `invoice_titles` — 发票抬头

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `user_id` | BIGINT UNSIGNED | FK→users(id), INDEX | |
| `title_type` | VARCHAR(24) | DEFAULT 'company' | company/person |
| `title_name` | VARCHAR(160) | NOT NULL | 抬头名称 |
| `tax_no` | VARCHAR(80) | | 税号 |
| `email` | VARCHAR(120) | | 接收邮箱 |
| `is_default` | TINYINT(1) | DEFAULT 0 | |

---

#### `invoices` — 发票

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `id` | BIGINT UNSIGNED | PK | |
| `user_id` | BIGINT UNSIGNED | FK→users(id), INDEX | |
| `order_id` | BIGINT UNSIGNED | FK→orders(id), INDEX | |
| `invoice_title_id` | BIGINT UNSIGNED | FK→invoice_titles(id), NULLABLE | |
| `invoice_no` | VARCHAR(80) | | 发票号码 |
| `amount` | DECIMAL(10,2) | NOT NULL | |
| `status` | VARCHAR(32) | DEFAULT 'requested' | requested/issued |
| `issued_at` | DATETIME | | 开票时间 |

---

## 三、安全设计要点

| 措施 | 应用表 | 说明 |
|---|---|---|
| VARBINARY 加密 | `users.id_card_no`, `travelers.id_no`, `drivers.id_card_no`, `order_travelers.id_no` | 身份证号不存明文 |
| 密码哈希 | `admin_users.password_hash` | 管理员密码单向加密 |
| 二维码哈希 | `tickets.qr_token_hash` | 二维码内容存 SHA-256 哈希，不存明文 token |
| 行锁防超卖 | `product_schedules` (locked_stock + sold_stock) | 事务 + 行级锁防止并发超售 |
| 操作审计 | `audit_logs` | 全量记录管理后台操作 |
| JSON 校验 | `payments.raw_payload`, `refunds.raw_payload` | MariaDB CHECK 约束确保 JSON 合法 |

## 四、索引策略

- **高频查询字段**：`user_id`, `order_id`, `status`, `driver_id`, `product_id` 均建索引
- **联合唯一约束**：排期 `(sku_id, travel_date, start_time)` / 收藏 `(user_id, product_id)`
- **订单列表**：`(user_id, created_at)` 复合索引加速我的订单翻页
- **佣金查询**：`(driver_id, status)` 复合索引加速司机佣金筛选

## 五、初始化脚本

SQL 文件位置：[database/001_schema.sql](file:///www/wwwroot/frontend/database/001_schema.sql)

```bash
# 创建数据库（已执行）
mysql -u root -p < database/001_schema.sql

# 验证
mysql -u root -p -e "USE zhuhai_travel; SHOW TABLES;"
```
