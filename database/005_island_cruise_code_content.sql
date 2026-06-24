USE zhuhai_travel;

ALTER TABLE island_cruise_orders
  ADD COLUMN supplier_code_content VARCHAR(255) NULL AFTER supplier_ticket_no;

ALTER TABLE island_cruise_passengers
  ADD COLUMN supplier_code_content VARCHAR(255) NULL AFTER supplier_ticket_no;
