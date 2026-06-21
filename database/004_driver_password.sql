USE zhuhai_travel;

ALTER TABLE drivers
  ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255) NOT NULL DEFAULT '' AFTER phone;
