USE zhuhai_travel;

ALTER TABLE island_cruise_orders
  ADD COLUMN balance_amount DECIMAL(12,2) NULL AFTER pay_evidence_no,
  ADD COLUMN balance_checked_at DATETIME NULL AFTER balance_amount,
  ADD COLUMN refund_status VARCHAR(32) NULL AFTER balance_checked_at,
  ADD COLUMN refund_flow_no VARCHAR(80) NULL AFTER refund_status,
  ADD COLUMN refund_fee DECIMAL(10,2) NOT NULL DEFAULT 0 AFTER refund_flow_no,
  ADD COLUMN refund_amount DECIMAL(10,2) NOT NULL DEFAULT 0 AFTER refund_fee,
  ADD COLUMN change_status VARCHAR(32) NULL AFTER refund_amount,
  ADD COLUMN change_order_no VARCHAR(80) NULL AFTER change_status,
  ADD COLUMN change_fee DECIMAL(10,2) NOT NULL DEFAULT 0 AFTER change_order_no,
  ADD COLUMN change_price_diff DECIMAL(10,2) NOT NULL DEFAULT 0 AFTER change_fee,
  ADD COLUMN verified_at DATETIME NULL AFTER change_price_diff,
  ADD COLUMN verify_status VARCHAR(32) NULL AFTER verified_at,
  ADD COLUMN balance_response JSON NULL AFTER order_response,
  ADD COLUMN refund_fee_response JSON NULL AFTER balance_response,
  ADD COLUMN refund_response JSON NULL AFTER refund_fee_response,
  ADD COLUMN change_fee_response JSON NULL AFTER refund_response,
  ADD COLUMN change_voyage_response JSON NULL AFTER change_fee_response,
  ADD COLUMN change_lock_response JSON NULL AFTER change_voyage_response,
  ADD COLUMN change_unlock_response JSON NULL AFTER change_lock_response,
  ADD COLUMN verify_response JSON NULL AFTER change_unlock_response;
