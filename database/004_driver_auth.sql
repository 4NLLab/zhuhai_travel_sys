-- 司机端登录密码字段
-- 已有司机需通过运营后台或数据脚本写入 bcrypt password_hash 后才能使用密码登录。

ALTER TABLE drivers
  ADD COLUMN password_hash VARCHAR(255) NOT NULL AFTER phone;
