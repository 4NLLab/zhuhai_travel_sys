USE zhuhai_travel;

ALTER TABLE drivers
  ADD COLUMN IF NOT EXISTS id_card_front_url VARCHAR(500) NOT NULL DEFAULT '' AFTER status,
  ADD COLUMN IF NOT EXISTS driver_license_url VARCHAR(500) NOT NULL DEFAULT '' AFTER id_card_front_url,
  ADD COLUMN IF NOT EXISTS vehicle_license_url VARCHAR(500) NOT NULL DEFAULT '' AFTER driver_license_url,
  ADD COLUMN IF NOT EXISTS vehicle_photo_url VARCHAR(500) NOT NULL DEFAULT '' AFTER vehicle_license_url,
  ADD COLUMN IF NOT EXISTS review_remark VARCHAR(255) NOT NULL DEFAULT '' AFTER vehicle_photo_url;
