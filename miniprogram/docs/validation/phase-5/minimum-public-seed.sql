-- Phase 5 local-only seed for public product/category contract checks.
-- Run against the Docker Compose MySQL database with utf8mb4 enabled.

INSERT INTO product_categories (id, name, slug, sort_order, status, created_at, updated_at)
VALUES (910001, '船票', 'ship', 1, 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE name = VALUES(name), status = VALUES(status), updated_at = NOW();

INSERT INTO products (id, category_id, title, subtitle, product_type, cover_url, location_name, status, sort_order, created_at, updated_at)
VALUES (
  910001,
  910001,
  '澳门环岛游夜景船票',
  '湾仔旅游码头出发，凭电子票码检票登船',
  'ship_ticket',
  '/static/phase2/macau-cruise-night-banner-web.jpg',
  '湾仔旅游码头',
  'active',
  1,
  NOW(),
  NOW()
)
ON DUPLICATE KEY UPDATE title = VALUES(title), status = VALUES(status), updated_at = NOW();

INSERT INTO product_skus (id, product_id, sku_name, sale_price, settlement_price, stock_mode, status, created_at, updated_at)
VALUES (910001, 910001, '成人票', 88.00, 70.00, 'schedule', 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE sale_price = VALUES(sale_price), status = VALUES(status), updated_at = NOW();

INSERT INTO product_images (id, product_id, image_url, alt_text, sort_order, created_at)
VALUES (910001, 910001, '/static/phase2/macau-cruise-night-banner-web.jpg', '澳门环岛游夜景船票', 1, NOW())
ON DUPLICATE KEY UPDATE image_url = VALUES(image_url);

INSERT INTO product_schedules (
  id,
  product_id,
  sku_id,
  travel_date,
  start_time,
  end_time,
  venue,
  total_stock,
  locked_stock,
  sold_stock,
  status,
  created_at,
  updated_at
)
VALUES (910001, 910001, 910001, '2026-07-15', '19:30:00', '21:00:00', '湾仔旅游码头', 50, 0, 0, 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE total_stock = VALUES(total_stock), status = VALUES(status), updated_at = NOW();
