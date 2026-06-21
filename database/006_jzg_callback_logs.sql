CREATE TABLE IF NOT EXISTS jzg_callback_logs (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
