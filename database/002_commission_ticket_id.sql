USE zhuhai_travel;

ALTER TABLE driver_commissions
  ADD COLUMN IF NOT EXISTS ticket_id BIGINT UNSIGNED NULL AFTER order_item_id;

ALTER TABLE driver_commissions
  ADD UNIQUE KEY IF NOT EXISTS uk_commissions_ticket (ticket_id);
