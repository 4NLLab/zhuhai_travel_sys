# Database setup

This project should use MySQL 8.4 LTS locally and on the future server.

## Why MySQL

- The current product is a transaction-heavy travel platform: products, schedules, orders, payments, tickets, verification, refunds, and driver commissions.
- MySQL is easy to run locally, common on Chinese cloud servers, and straightforward to migrate from local development to a remote server.
- MySQL 8.4 LTS is safer than starting with MySQL 9.x because many hosting environments and managed database services still standardize around 8.x.

## Local install with Homebrew

```bash
brew install mysql@8.4
brew services start mysql@8.4
```

If `mysql` is not on your PATH after install, use:

```bash
/opt/homebrew/opt/mysql@8.4/bin/mysql -u root
```

## Create the database

```bash
/opt/homebrew/opt/mysql@8.4/bin/mysql -u root < database/001_schema.sql
```

The script creates a database named `zhuhai_travel`.

## Suggested local connection

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=zhuhai_travel
DB_USER=root
DB_PASSWORD=
```

Before using a real server, create a dedicated app user instead of using `root`.

```sql
CREATE USER 'zhuhai_app'@'%' IDENTIFIED BY 'replace-with-a-strong-password';
GRANT SELECT, INSERT, UPDATE, DELETE ON zhuhai_travel.* TO 'zhuhai_app'@'%';
FLUSH PRIVILEGES;
```

## Main table groups

- User side: `users`, `travelers`, `orders`, `order_items`, `order_travelers`, `payments`, `tickets`, `refunds`, `invoice_titles`, `invoices`, `user_favorites`.
- Product and inventory: `product_categories`, `products`, `product_images`, `product_skus`, `product_schedules`, `banners`.
- Driver channel: `drivers`, `vehicles`, `driver_qr_codes`, `driver_commissions`.
- Admin and operations: `admin_users`, `ticket_verifications`, `support_tickets`, `audit_logs`.
