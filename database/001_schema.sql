-- Zhuhai travel platform schema for MySQL 8.4+
-- Import with: mysql -u root < database/001_schema.sql

CREATE DATABASE IF NOT EXISTS zhuhai_travel
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_0900_ai_ci;

USE zhuhai_travel;

SET NAMES utf8mb4;
SET time_zone = '+08:00';

CREATE TABLE users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  openid VARCHAR(80) NULL,
  unionid VARCHAR(80) NULL,
  phone VARCHAR(32) NULL,
  nickname VARCHAR(80) NULL,
  avatar_url VARCHAR(500) NULL,
  real_name VARCHAR(80) NULL,
  id_card_no VARBINARY(255) NULL,
  realname_status VARCHAR(24) NOT NULL DEFAULT 'unverified',
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  last_login_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_users_openid (openid),
  UNIQUE KEY uk_users_unionid (unionid),
  KEY idx_users_phone (phone),
  KEY idx_users_status (status)
) ENGINE=InnoDB;

CREATE TABLE travelers (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(80) NOT NULL,
  phone VARCHAR(32) NULL,
  id_type VARCHAR(24) NOT NULL DEFAULT 'id_card',
  id_no VARBINARY(255) NOT NULL,
  is_default TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_travelers_user (user_id),
  CONSTRAINT fk_travelers_user FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=InnoDB;

CREATE TABLE admin_users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(80) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  display_name VARCHAR(80) NOT NULL,
  role VARCHAR(32) NOT NULL DEFAULT 'operator',
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  last_login_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_admin_username (username),
  KEY idx_admin_role_status (role, status)
) ENGINE=InnoDB;

CREATE TABLE drivers (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  driver_no VARCHAR(40) NOT NULL,
  name VARCHAR(80) NOT NULL,
  phone VARCHAR(32) NOT NULL,
  password_hash VARCHAR(255) NOT NULL DEFAULT '',
  id_card_no VARBINARY(255) NULL,
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  id_card_front_url VARCHAR(500) NOT NULL DEFAULT '',
  driver_license_url VARCHAR(500) NOT NULL DEFAULT '',
  vehicle_license_url VARCHAR(500) NOT NULL DEFAULT '',
  vehicle_photo_url VARCHAR(500) NOT NULL DEFAULT '',
  review_remark VARCHAR(255) NOT NULL DEFAULT '',
  commission_rate DECIMAL(6,4) NOT NULL DEFAULT 0.0800,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_drivers_no (driver_no),
  KEY idx_drivers_phone (phone),
  KEY idx_drivers_status (status)
) ENGINE=InnoDB;

CREATE TABLE vehicles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  driver_id BIGINT UNSIGNED NOT NULL,
  plate_no VARCHAR(32) NOT NULL,
  model VARCHAR(80) NULL,
  seats INT UNSIGNED NULL,
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_vehicles_plate (plate_no),
  KEY idx_vehicles_driver (driver_id),
  CONSTRAINT fk_vehicles_driver FOREIGN KEY (driver_id) REFERENCES drivers (id)
) ENGINE=InnoDB;

CREATE TABLE driver_qr_codes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  driver_id BIGINT UNSIGNED NOT NULL,
  vehicle_id BIGINT UNSIGNED NULL,
  code VARCHAR(80) NOT NULL,
  scene VARCHAR(80) NOT NULL DEFAULT 'seat',
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_driver_qr_code (code),
  KEY idx_driver_qr_driver (driver_id),
  KEY idx_driver_qr_vehicle (vehicle_id),
  CONSTRAINT fk_driver_qr_driver FOREIGN KEY (driver_id) REFERENCES drivers (id),
  CONSTRAINT fk_driver_qr_vehicle FOREIGN KEY (vehicle_id) REFERENCES vehicles (id)
) ENGINE=InnoDB;

CREATE TABLE product_categories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  parent_id BIGINT UNSIGNED NULL,
  name VARCHAR(80) NOT NULL,
  slug VARCHAR(80) NOT NULL,
  icon VARCHAR(80) NULL,
  sort_order INT NOT NULL DEFAULT 0,
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_categories_slug (slug),
  KEY idx_categories_parent (parent_id),
  CONSTRAINT fk_categories_parent FOREIGN KEY (parent_id) REFERENCES product_categories (id)
) ENGINE=InnoDB;

CREATE TABLE products (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  category_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(160) NOT NULL,
  subtitle VARCHAR(255) NULL,
  product_type VARCHAR(32) NOT NULL,
  cover_url VARCHAR(500) NULL,
  location_name VARCHAR(120) NULL,
  address VARCHAR(255) NULL,
  notice TEXT NULL,
  refund_policy TEXT NULL,
  status VARCHAR(24) NOT NULL DEFAULT 'draft',
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_products_category (category_id),
  KEY idx_products_type_status (product_type, status),
  CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES product_categories (id)
) ENGINE=InnoDB;

CREATE TABLE product_images (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  product_id BIGINT UNSIGNED NOT NULL,
  image_url VARCHAR(500) NOT NULL,
  alt_text VARCHAR(160) NULL,
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_product_images_product (product_id),
  CONSTRAINT fk_product_images_product FOREIGN KEY (product_id) REFERENCES products (id)
) ENGINE=InnoDB;

CREATE TABLE product_skus (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  product_id BIGINT UNSIGNED NOT NULL,
  sku_name VARCHAR(120) NOT NULL,
  market_price DECIMAL(10,2) NULL,
  sale_price DECIMAL(10,2) NOT NULL,
  settlement_price DECIMAL(10,2) NULL,
  stock_mode VARCHAR(24) NOT NULL DEFAULT 'schedule',
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_skus_product (product_id),
  KEY idx_skus_status (status),
  CONSTRAINT fk_skus_product FOREIGN KEY (product_id) REFERENCES products (id)
) ENGINE=InnoDB;

CREATE TABLE product_schedules (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  product_id BIGINT UNSIGNED NOT NULL,
  sku_id BIGINT UNSIGNED NOT NULL,
  travel_date DATE NOT NULL,
  start_time TIME NULL,
  end_time TIME NULL,
  venue VARCHAR(120) NULL,
  total_stock INT NOT NULL DEFAULT 0,
  locked_stock INT NOT NULL DEFAULT 0,
  sold_stock INT NOT NULL DEFAULT 0,
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_schedule_sku_time (sku_id, travel_date, start_time),
  KEY idx_schedules_product_date (product_id, travel_date),
  CONSTRAINT fk_schedules_product FOREIGN KEY (product_id) REFERENCES products (id),
  CONSTRAINT fk_schedules_sku FOREIGN KEY (sku_id) REFERENCES product_skus (id)
) ENGINE=InnoDB;

CREATE TABLE banners (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title VARCHAR(120) NOT NULL,
  subtitle VARCHAR(255) NULL,
  image_url VARCHAR(500) NOT NULL,
  link_type VARCHAR(32) NOT NULL DEFAULT 'product',
  link_target VARCHAR(160) NULL,
  sort_order INT NOT NULL DEFAULT 0,
  status VARCHAR(24) NOT NULL DEFAULT 'active',
  starts_at DATETIME NULL,
  ends_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_banners_status_sort (status, sort_order)
) ENGINE=InnoDB;

CREATE TABLE orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_no VARCHAR(40) NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  source VARCHAR(32) NOT NULL DEFAULT 'miniapp',
  driver_id BIGINT UNSIGNED NULL,
  driver_qr_code_id BIGINT UNSIGNED NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'pending_payment',
  total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  payable_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  paid_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  contact_name VARCHAR(80) NULL,
  contact_phone VARCHAR(32) NULL,
  remark VARCHAR(255) NULL,
  paid_at DATETIME NULL,
  cancelled_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_orders_no (order_no),
  KEY idx_orders_user_created (user_id, created_at),
  KEY idx_orders_status (status),
  KEY idx_orders_driver (driver_id),
  CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_orders_driver FOREIGN KEY (driver_id) REFERENCES drivers (id),
  CONSTRAINT fk_orders_driver_qr FOREIGN KEY (driver_qr_code_id) REFERENCES driver_qr_codes (id)
) ENGINE=InnoDB;

CREATE TABLE order_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_id BIGINT UNSIGNED NOT NULL,
  product_id BIGINT UNSIGNED NOT NULL,
  sku_id BIGINT UNSIGNED NOT NULL,
  schedule_id BIGINT UNSIGNED NULL,
  product_title VARCHAR(160) NOT NULL,
  sku_name VARCHAR(120) NOT NULL,
  travel_date DATE NULL,
  start_time TIME NULL,
  quantity INT UNSIGNED NOT NULL,
  unit_price DECIMAL(10,2) NOT NULL,
  total_price DECIMAL(10,2) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_order_items_order (order_id),
  KEY idx_order_items_product (product_id),
  KEY idx_order_items_schedule (schedule_id),
  CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders (id),
  CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products (id),
  CONSTRAINT fk_order_items_sku FOREIGN KEY (sku_id) REFERENCES product_skus (id),
  CONSTRAINT fk_order_items_schedule FOREIGN KEY (schedule_id) REFERENCES product_schedules (id)
) ENGINE=InnoDB;

CREATE TABLE order_travelers (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_item_id BIGINT UNSIGNED NOT NULL,
  traveler_id BIGINT UNSIGNED NULL,
  name VARCHAR(80) NOT NULL,
  phone VARCHAR(32) NULL,
  id_type VARCHAR(24) NOT NULL DEFAULT 'id_card',
  id_no VARBINARY(255) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_order_travelers_item (order_item_id),
  KEY idx_order_travelers_traveler (traveler_id),
  CONSTRAINT fk_order_travelers_item FOREIGN KEY (order_item_id) REFERENCES order_items (id),
  CONSTRAINT fk_order_travelers_traveler FOREIGN KEY (traveler_id) REFERENCES travelers (id)
) ENGINE=InnoDB;

CREATE TABLE payments (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_id BIGINT UNSIGNED NOT NULL,
  payment_no VARCHAR(64) NOT NULL,
  channel VARCHAR(32) NOT NULL DEFAULT 'wechat',
  amount DECIMAL(10,2) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'pending',
  transaction_id VARCHAR(128) NULL,
  paid_at DATETIME NULL,
  raw_payload JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_payments_no (payment_no),
  KEY idx_payments_order (order_id),
  KEY idx_payments_status (status),
  CONSTRAINT fk_payments_order FOREIGN KEY (order_id) REFERENCES orders (id)
) ENGINE=InnoDB;

CREATE TABLE tickets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_item_id BIGINT UNSIGNED NOT NULL,
  ticket_no VARCHAR(64) NOT NULL,
  qr_token_hash CHAR(64) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'valid',
  valid_from DATETIME NULL,
  valid_to DATETIME NULL,
  used_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_tickets_no (ticket_no),
  UNIQUE KEY uk_tickets_qr_hash (qr_token_hash),
  KEY idx_tickets_item (order_item_id),
  KEY idx_tickets_status (status),
  CONSTRAINT fk_tickets_order_item FOREIGN KEY (order_item_id) REFERENCES order_items (id)
) ENGINE=InnoDB;

CREATE TABLE ticket_verifications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  ticket_id BIGINT UNSIGNED NOT NULL,
  verifier_admin_id BIGINT UNSIGNED NULL,
  verify_location VARCHAR(120) NULL,
  verify_result VARCHAR(32) NOT NULL,
  verify_note VARCHAR(255) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_verifications_ticket (ticket_id),
  KEY idx_verifications_created (created_at),
  CONSTRAINT fk_verifications_ticket FOREIGN KEY (ticket_id) REFERENCES tickets (id),
  CONSTRAINT fk_verifications_admin FOREIGN KEY (verifier_admin_id) REFERENCES admin_users (id)
) ENGINE=InnoDB;

CREATE TABLE refunds (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_id BIGINT UNSIGNED NOT NULL,
  order_item_id BIGINT UNSIGNED NULL,
  refund_no VARCHAR(64) NOT NULL,
  amount DECIMAL(10,2) NOT NULL,
  reason VARCHAR(255) NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'requested',
  requested_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  processed_at DATETIME NULL,
  raw_payload JSON NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_refunds_no (refund_no),
  KEY idx_refunds_order (order_id),
  KEY idx_refunds_status (status),
  CONSTRAINT fk_refunds_order FOREIGN KEY (order_id) REFERENCES orders (id),
  CONSTRAINT fk_refunds_order_item FOREIGN KEY (order_item_id) REFERENCES order_items (id)
) ENGINE=InnoDB;

CREATE TABLE driver_commissions (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  driver_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL,
  order_item_id BIGINT UNSIGNED NOT NULL,
  ticket_id BIGINT UNSIGNED NULL,
  commission_no VARCHAR(64) NOT NULL,
  base_amount DECIMAL(10,2) NOT NULL,
  rate DECIMAL(6,4) NOT NULL,
  commission_amount DECIMAL(10,2) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'pending',
  settled_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_commissions_no (commission_no),
  UNIQUE KEY uk_commissions_ticket (ticket_id),
  KEY idx_commissions_driver_status (driver_id, status),
  KEY idx_commissions_order (order_id),
  CONSTRAINT fk_commissions_driver FOREIGN KEY (driver_id) REFERENCES drivers (id),
  CONSTRAINT fk_commissions_order FOREIGN KEY (order_id) REFERENCES orders (id),
  CONSTRAINT fk_commissions_order_item FOREIGN KEY (order_item_id) REFERENCES order_items (id),
  CONSTRAINT fk_commissions_ticket FOREIGN KEY (ticket_id) REFERENCES tickets (id)
) ENGINE=InnoDB;

CREATE TABLE support_tickets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NULL,
  type VARCHAR(32) NOT NULL,
  title VARCHAR(160) NOT NULL,
  content TEXT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'open',
  handled_by BIGINT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_support_user (user_id),
  KEY idx_support_order (order_id),
  KEY idx_support_status (status),
  CONSTRAINT fk_support_user FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_support_order FOREIGN KEY (order_id) REFERENCES orders (id),
  CONSTRAINT fk_support_admin FOREIGN KEY (handled_by) REFERENCES admin_users (id)
) ENGINE=InnoDB;

CREATE TABLE invoice_titles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  title_type VARCHAR(24) NOT NULL DEFAULT 'company',
  title_name VARCHAR(160) NOT NULL,
  tax_no VARCHAR(80) NULL,
  email VARCHAR(120) NULL,
  is_default TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_invoice_titles_user (user_id),
  CONSTRAINT fk_invoice_titles_user FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=InnoDB;

CREATE TABLE invoices (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL,
  invoice_title_id BIGINT UNSIGNED NULL,
  invoice_no VARCHAR(80) NULL,
  amount DECIMAL(10,2) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'requested',
  issued_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_invoices_user (user_id),
  KEY idx_invoices_order (order_id),
  CONSTRAINT fk_invoices_user FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_invoices_order FOREIGN KEY (order_id) REFERENCES orders (id),
  CONSTRAINT fk_invoices_title FOREIGN KEY (invoice_title_id) REFERENCES invoice_titles (id)
) ENGINE=InnoDB;

CREATE TABLE user_favorites (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  product_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_favorites_user_product (user_id, product_id),
  KEY idx_favorites_product (product_id),
  CONSTRAINT fk_favorites_user FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_favorites_product FOREIGN KEY (product_id) REFERENCES products (id)
) ENGINE=InnoDB;

CREATE TABLE audit_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  actor_type VARCHAR(24) NOT NULL,
  actor_id BIGINT UNSIGNED NULL,
  actor_name VARCHAR(120) NULL,
  action VARCHAR(80) NOT NULL,
  target_type VARCHAR(80) NULL,
  target_id BIGINT UNSIGNED NULL,
  method VARCHAR(12) NOT NULL DEFAULT '',
  path VARCHAR(255) NOT NULL DEFAULT '',
  status_code INT NOT NULL DEFAULT 0,
  ip VARCHAR(64) NULL,
  user_agent VARCHAR(500) NULL,
  payload JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_audit_actor (actor_type, actor_id),
  KEY idx_audit_target (target_type, target_id),
  KEY idx_audit_created (created_at)
) ENGINE=InnoDB;

CREATE TABLE jzg_callback_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  event_type VARCHAR(40) NOT NULL,
  order_no VARCHAR(64) NULL,
  verified TINYINT(1) NOT NULL DEFAULT 0,
  request_ip VARCHAR(64) NULL,
  payload JSON NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_jzg_callback_event_time (event_type, created_at),
  KEY idx_jzg_callback_order (order_no),
  KEY idx_jzg_callback_verified (verified)
) ENGINE=InnoDB;
