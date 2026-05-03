-- ----------------------------
-- PostgreSQL 兼容版本 ph_orders 建表语句
-- ----------------------------

-- ----------------------------
-- 自动更新 updated_at 的通用触发器函数
-- ----------------------------
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ----------------------------
-- Table structure for bank_card
-- ----------------------------
DROP TABLE IF EXISTS bank_card CASCADE;
CREATE TABLE bank_card (
  id BIGSERIAL PRIMARY KEY,
  user_id VARCHAR(50) NOT NULL,
  card_number VARCHAR(20) NOT NULL,
  card_type SMALLINT NOT NULL,
  bank_code VARCHAR(10) NOT NULL,
  bank_name VARCHAR(50) NOT NULL,
  interbank_transfer_code VARCHAR(20) DEFAULT NULL,
  account_holder VARCHAR(50) NOT NULL,
  is_default SMALLINT DEFAULT 0,
  status SMALLINT DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE bank_card IS 'Bank card information table';
COMMENT ON COLUMN bank_card.id IS 'Primary key';
COMMENT ON COLUMN bank_card.user_id IS 'User ID';
COMMENT ON COLUMN bank_card.card_number IS 'Bank card number';
COMMENT ON COLUMN bank_card.card_type IS 'Card type (e.g., 1: debit, 2: credit)';
COMMENT ON COLUMN bank_card.bank_code IS 'Bank code';
COMMENT ON COLUMN bank_card.bank_name IS 'Bank name';
COMMENT ON COLUMN bank_card.interbank_transfer_code IS 'Interbank transfer code (e.g., SWIFT/IFSC)';
COMMENT ON COLUMN bank_card.account_holder IS 'Account holder name';
COMMENT ON COLUMN bank_card.is_default IS 'Is default card (0: no, 1: yes)';
COMMENT ON COLUMN bank_card.status IS 'Status (0: inactive, 1: active)';
COMMENT ON COLUMN bank_card.created_at IS 'Creation time';
COMMENT ON COLUMN bank_card.updated_at IS 'Update time';

CREATE INDEX idx_bank_card_user_id ON bank_card (user_id);
CREATE INDEX idx_bank_card_card_number ON bank_card (card_number);

CREATE TRIGGER update_bank_card_modtime BEFORE UPDATE ON bank_card FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for interview_record
-- ----------------------------
DROP TABLE IF EXISTS interview_record CASCADE;
CREATE TABLE interview_record (
  id BIGSERIAL PRIMARY KEY,
  order_id VARCHAR(64) NOT NULL,
  initiator_type SMALLINT NOT NULL,
  interview_type SMALLINT NOT NULL,
  appointment_time TIMESTAMP DEFAULT NULL,
  location VARCHAR(255) DEFAULT NULL,
  message TEXT,
  status SMALLINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_result SMALLINT DEFAULT NULL,
  sitter_result SMALLINT DEFAULT NULL,
  user_reason VARCHAR(255) DEFAULT NULL,
  sitter_reason VARCHAR(255) DEFAULT NULL
);
COMMENT ON COLUMN interview_record.id IS '面试预约记录唯一ID';
COMMENT ON COLUMN interview_record.order_id IS '关联的订单ID，用于绑定业务订单';
COMMENT ON COLUMN interview_record.initiator_type IS '预约发起方类型：1-用户 2-Sitter';
COMMENT ON COLUMN interview_record.interview_type IS '面试类型：1-线上 2-离线会议';
COMMENT ON COLUMN interview_record.appointment_time IS '预约面试时间';
COMMENT ON COLUMN interview_record.location IS '面试地点（线上则为空）';
COMMENT ON COLUMN interview_record.message IS '预约时填写的消息内容';
COMMENT ON COLUMN interview_record.status IS '预约状态：1-待确认（刚发起） 2-接受方修改待确认 3-已接受 4-已取消 5-已拒绝';
COMMENT ON COLUMN interview_record.created_at IS '预约创建时间';
COMMENT ON COLUMN interview_record.updated_at IS '记录更新时间';

CREATE TRIGGER update_interview_record_modtime BEFORE UPDATE ON interview_record FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for orders
-- ----------------------------
DROP TABLE IF EXISTS orders CASCADE;
CREATE TABLE orders (
  id BIGSERIAL PRIMARY KEY,
  owner_id VARCHAR(50) NOT NULL,
  sitter_id VARCHAR(50) NOT NULL,
  type BIGINT NOT NULL,
  from_date DATE NOT NULL,
  to_date DATE NOT NULL,
  tips_price DECIMAL(10,2) DEFAULT NULL,
  total_price DECIMAL(10,2) DEFAULT NULL,
  sub_total_price DECIMAL(10,2) DEFAULT NULL,
  service_fee DECIMAL(10,2) DEFAULT NULL,
  taxes DECIMAL(10,2) DEFAULT NULL,
  state INTEGER NOT NULL,
  sitter_handle_at TIMESTAMP NULL DEFAULT NULL,
  sitter_finish_at TIMESTAMP NULL DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  owner_name VARCHAR(255) DEFAULT NULL,
  sitter_name VARCHAR(255) DEFAULT NULL,
  contact VARCHAR(255) DEFAULT NULL,
  alternative_contact VARCHAR(255) DEFAULT NULL,
  note TEXT,
  cancel_at TIMESTAMP NULL DEFAULT NULL,
  cancel_reason VARCHAR(255) DEFAULT NULL,
  refund_price DECIMAL(10,2) DEFAULT NULL,
  order_number VARCHAR(255) DEFAULT NULL,
  user_deleted SMALLINT NOT NULL DEFAULT 0,
  sitter_deleted SMALLINT NOT NULL DEFAULT 0,
  user_rating_state SMALLINT DEFAULT 0,
  sitter_rating_state SMALLINT DEFAULT 0
);
COMMENT ON TABLE orders IS '订单表';
COMMENT ON COLUMN orders.id IS '订单ID';
COMMENT ON COLUMN orders.owner_id IS '用户ID';
COMMENT ON COLUMN orders.sitter_id IS 'Sitter ID';
COMMENT ON COLUMN orders.type IS '订单类型';
COMMENT ON COLUMN orders.from_date IS '开始时间';
COMMENT ON COLUMN orders.to_date IS '结束时间';
COMMENT ON COLUMN orders.tips_price IS '小费';
COMMENT ON COLUMN orders.total_price IS '总费用';
COMMENT ON COLUMN orders.sub_total_price IS '宠物总费用';
COMMENT ON COLUMN orders.service_fee IS '服务费';
COMMENT ON COLUMN orders.taxes IS '税';
COMMENT ON COLUMN orders.state IS '订单状态';
COMMENT ON COLUMN orders.sitter_handle_at IS 'Sitter 处理时间';
COMMENT ON COLUMN orders.sitter_finish_at IS 'Sitter 完成时间';
COMMENT ON COLUMN orders.created_at IS '创建时间';
COMMENT ON COLUMN orders.updated_at IS '更新时间';

CREATE INDEX idx_orders_owner_id ON orders (owner_id);
CREATE INDEX idx_orders_sitter_id ON orders (sitter_id);
CREATE INDEX idx_orders_state ON orders (state);
CREATE INDEX idx_orders_dates ON orders (from_date, to_date);

CREATE TRIGGER update_orders_modtime BEFORE UPDATE ON orders FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for order_modification_log
-- ----------------------------
DROP TABLE IF EXISTS order_modification_log CASCADE;
CREATE TABLE order_modification_log (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL,
  owner_id VARCHAR(50) NOT NULL,
  sitter_id VARCHAR(50) NOT NULL,
  previous_date DATE DEFAULT NULL,
  new_date DATE DEFAULT NULL,
  previous_pet_list TEXT,
  new_pet_list TEXT,
  previous_price DECIMAL(10,2) DEFAULT NULL,
  new_price DECIMAL(10,2) DEFAULT NULL,
  state INTEGER NOT NULL DEFAULT 0,
  type INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_order_modification_log_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);
COMMENT ON TABLE order_modification_log IS '订单修改日志表';
COMMENT ON COLUMN order_modification_log.id IS '日志ID';
COMMENT ON COLUMN order_modification_log.order_id IS '关联的订单ID';
COMMENT ON COLUMN order_modification_log.owner_id IS '用户ID';
COMMENT ON COLUMN order_modification_log.sitter_id IS 'Sitter ID';
COMMENT ON COLUMN order_modification_log.previous_date IS '修改前的日期';
COMMENT ON COLUMN order_modification_log.new_date IS '修改后的日期';
COMMENT ON COLUMN order_modification_log.previous_pet_list IS '修改前的宠物列表(JSON格式)';
COMMENT ON COLUMN order_modification_log.new_pet_list IS '修改后的宠物列表(JSON格式)';
COMMENT ON COLUMN order_modification_log.previous_price IS '修改前的价格';
COMMENT ON COLUMN order_modification_log.new_price IS '修改后的价格';
COMMENT ON COLUMN order_modification_log.state IS '修改状态';
COMMENT ON COLUMN order_modification_log.type IS '修改类型';
COMMENT ON COLUMN order_modification_log.created_at IS '创建时间';
COMMENT ON COLUMN order_modification_log.updated_at IS '更新时间';

CREATE INDEX idx_oml_order_id ON order_modification_log (order_id);
CREATE INDEX idx_oml_owner_id ON order_modification_log (owner_id);
CREATE INDEX idx_oml_sitter_id ON order_modification_log (sitter_id);
CREATE INDEX idx_oml_created_at ON order_modification_log (created_at);

CREATE TRIGGER update_order_modification_log_modtime BEFORE UPDATE ON order_modification_log FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for order_pet
-- ----------------------------
DROP TABLE IF EXISTS order_pet CASCADE;
CREATE TABLE order_pet (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL,
  pet_id VARCHAR(255) NOT NULL,
  pet_type VARCHAR(255) NOT NULL,
  pet_shape INTEGER NOT NULL,
  pet_name VARCHAR(255) DEFAULT NULL,
  breed VARCHAR(255) DEFAULT NULL,
  pet_price DECIMAL(10,2) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  owner_id VARCHAR(255) DEFAULT NULL,
  sitter_id VARCHAR(255) DEFAULT NULL,
  CONSTRAINT fk_order_pet_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);
COMMENT ON TABLE order_pet IS '订单宠物关联表';
COMMENT ON COLUMN order_pet.id IS '主键ID';
COMMENT ON COLUMN order_pet.order_id IS '订单ID';
COMMENT ON COLUMN order_pet.pet_id IS '宠物ID';
COMMENT ON COLUMN order_pet.pet_type IS '宠物类型，如猫、狗';
COMMENT ON COLUMN order_pet.pet_shape IS '宠物体型 1为小型犬 2为中型犬 3为大型犬';
COMMENT ON COLUMN order_pet.pet_name IS '宠物名字';
COMMENT ON COLUMN order_pet.breed IS '宠物品种';
COMMENT ON COLUMN order_pet.pet_price IS '宠物费用';
COMMENT ON COLUMN order_pet.created_at IS '创建时间';
COMMENT ON COLUMN order_pet.updated_at IS '更新时间';

CREATE INDEX idx_order_pet_order_id ON order_pet (order_id);
CREATE INDEX idx_order_pet_pet_id ON order_pet (pet_id);
CREATE INDEX idx_order_pet_pet_type ON order_pet (pet_type);
CREATE INDEX idx_order_pet_pet_shape ON order_pet (pet_shape);

CREATE TRIGGER update_order_pet_modtime BEFORE UPDATE ON order_pet FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for owner_rating
-- ----------------------------
DROP TABLE IF EXISTS owner_rating CASCADE;
CREATE TABLE owner_rating (
  id SERIAL PRIMARY KEY,
  order_id VARCHAR(50) DEFAULT NULL,
  owner_id VARCHAR(50) DEFAULT NULL,
  sitter_id VARCHAR(50) DEFAULT NULL,
  owner_name VARCHAR(50) DEFAULT NULL,
  sitter_name VARCHAR(50) DEFAULT NULL,
  score SMALLINT DEFAULT NULL,
  satisfaction_level SMALLINT DEFAULT NULL,
  instructions_clarity SMALLINT DEFAULT NULL,
  communication SMALLINT DEFAULT NULL,
  supplies_preparation SMALLINT DEFAULT NULL,
  respect_courtesy SMALLINT DEFAULT NULL,
  suggestions TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_owner_rating_modtime BEFORE UPDATE ON owner_rating FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for pet_rating
-- ----------------------------
DROP TABLE IF EXISTS pet_rating CASCADE;
CREATE TABLE pet_rating (
  id SERIAL PRIMARY KEY,
  order_id VARCHAR(50) DEFAULT NULL,
  pet_id VARCHAR(50) DEFAULT NULL,
  pet_name VARCHAR(50) DEFAULT NULL,
  score SMALLINT DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_pet_rating_modtime BEFORE UPDATE ON pet_rating FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for report
-- ----------------------------
DROP TABLE IF EXISTS report CASCADE;
CREATE TABLE report (
  id BIGSERIAL PRIMARY KEY,
  reporter_id VARCHAR(50) NOT NULL,
  reported_id VARCHAR(50) NOT NULL,
  reported_name VARCHAR(50) DEFAULT NULL,
  report_type VARCHAR(50) DEFAULT NULL,
  report_reason VARCHAR(50) DEFAULT NULL,
  report_desc TEXT NOT NULL,
  report_img TEXT,
  status SMALLINT DEFAULT 0,
  handler_id VARCHAR(50) DEFAULT NULL,
  handle_result TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE report IS '用户举报记录表';
COMMENT ON COLUMN report.id IS '举报ID';
COMMENT ON COLUMN report.reporter_id IS '举报人ID';
COMMENT ON COLUMN report.reported_id IS '被举报人ID';
COMMENT ON COLUMN report.reported_name IS '被举报人姓名';
COMMENT ON COLUMN report.report_type IS '举报类型：1-欺诈，2-骚扰，3-不当内容，4-其他';
COMMENT ON COLUMN report.report_reason IS '举报原因';
COMMENT ON COLUMN report.report_desc IS '举报内容';
COMMENT ON COLUMN report.report_img IS '证据(图片/视频等)';
COMMENT ON COLUMN report.status IS '处理状态：0-待处理，1-处理中，2-已处理';
COMMENT ON COLUMN report.handler_id IS '处理人ID';
COMMENT ON COLUMN report.handle_result IS '处理结果';
COMMENT ON COLUMN report.created_at IS '创建时间';
COMMENT ON COLUMN report.updated_at IS '更新时间';

CREATE INDEX idx_report_reporter_id ON report (reporter_id);
CREATE INDEX idx_report_reported_id ON report (reported_id);
CREATE INDEX idx_report_status ON report (status);

CREATE TRIGGER update_report_modtime BEFORE UPDATE ON report FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for sitter_locations
-- ----------------------------
DROP TABLE IF EXISTS sitter_locations CASCADE;
CREATE TABLE sitter_locations (
  id BIGSERIAL PRIMARY KEY,
  sitter_id VARCHAR(50) NOT NULL,
  order_id BIGINT NOT NULL,
  sub_order_id BIGINT NOT NULL,
  lat DECIMAL(10,6) NOT NULL,
  lng DECIMAL(10,6) NOT NULL,
  timestamp BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE sitter_locations IS 'Sitter实时位置记录表';
COMMENT ON COLUMN sitter_locations.id IS '位置记录ID';
COMMENT ON COLUMN sitter_locations.sitter_id IS 'Sitter ID';
COMMENT ON COLUMN sitter_locations.order_id IS '主订单ID';
COMMENT ON COLUMN sitter_locations.sub_order_id IS '子订单ID';
COMMENT ON COLUMN sitter_locations.lat IS '纬度';
COMMENT ON COLUMN sitter_locations.lng IS '经度';
COMMENT ON COLUMN sitter_locations.timestamp IS '时间戳';
COMMENT ON COLUMN sitter_locations.created_at IS '创建时间';

CREATE INDEX idx_sl_sitter_id ON sitter_locations (sitter_id);
CREATE INDEX idx_sl_order_id ON sitter_locations (order_id);
CREATE INDEX idx_sl_sub_order_id ON sitter_locations (sub_order_id);
CREATE INDEX idx_sl_timestamp ON sitter_locations (timestamp);
CREATE INDEX idx_sl_sitter_suborder ON sitter_locations (sitter_id, sub_order_id);
CREATE INDEX idx_sl_order_suborder ON sitter_locations (order_id, sub_order_id);

-- ----------------------------
-- Table structure for sitter_rating
-- ----------------------------
DROP TABLE IF EXISTS sitter_rating CASCADE;
CREATE TABLE sitter_rating (
  id SERIAL PRIMARY KEY,
  order_id INTEGER NOT NULL UNIQUE,
  user_id VARCHAR(255) NOT NULL,
  sitter_id VARCHAR(255) NOT NULL,
  punctuality SMALLINT DEFAULT NULL,
  responsibility SMALLINT DEFAULT NULL,
  communication SMALLINT DEFAULT NULL,
  pet_care_skills SMALLINT DEFAULT NULL,
  cleanliness SMALLINT DEFAULT NULL,
  suggestions TEXT,
  score SMALLINT DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_sitter_rating_modtime BEFORE UPDATE ON sitter_rating FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for sub_orders
-- ----------------------------
DROP TABLE IF EXISTS sub_orders CASCADE;
CREATE TABLE sub_orders (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL,
  date DATE NOT NULL,
  state INTEGER NOT NULL,
  start_code VARCHAR(10) DEFAULT NULL,
  end_code VARCHAR(10) DEFAULT NULL,
  sitter_handle_at TIMESTAMP NULL DEFAULT NULL,
  sitter_finish_at TIMESTAMP NULL DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  walk_thumbnail_url VARCHAR(255) DEFAULT NULL
);
COMMENT ON TABLE sub_orders IS '子订单表';
COMMENT ON COLUMN sub_orders.id IS '子订单ID';
COMMENT ON COLUMN sub_orders.order_id IS '主订单ID';
COMMENT ON COLUMN sub_orders.date IS '日期';
COMMENT ON COLUMN sub_orders.state IS '子订单状态';
COMMENT ON COLUMN sub_orders.start_code IS '开始验证码';
COMMENT ON COLUMN sub_orders.end_code IS '结束验证码';
COMMENT ON COLUMN sub_orders.sitter_handle_at IS 'Sitter处理时间';
COMMENT ON COLUMN sub_orders.sitter_finish_at IS 'Sitter完成时间';
COMMENT ON COLUMN sub_orders.created_at IS '创建时间';
COMMENT ON COLUMN sub_orders.updated_at IS '更新时间';

CREATE INDEX idx_sub_orders_order_id ON sub_orders (order_id);
CREATE INDEX idx_sub_orders_date ON sub_orders (date);

CREATE TRIGGER update_sub_orders_modtime BEFORE UPDATE ON sub_orders FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for user_wallet
-- ----------------------------
DROP TABLE IF EXISTS user_wallet CASCADE;
CREATE TABLE user_wallet (
  id BIGSERIAL PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL UNIQUE,
  balance DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  status SMALLINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE user_wallet IS '用户钱包主表';
COMMENT ON COLUMN user_wallet.id IS '主键ID';
COMMENT ON COLUMN user_wallet.user_id IS '用户ID（关联用户表）';
COMMENT ON COLUMN user_wallet.balance IS '账户余额（元）';
COMMENT ON COLUMN user_wallet.status IS '钱包状态：1-正常，2-冻结';
COMMENT ON COLUMN user_wallet.created_at IS '创建时间';
COMMENT ON COLUMN user_wallet.updated_at IS '更新时间';

CREATE TRIGGER update_user_wallet_modtime BEFORE UPDATE ON user_wallet FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for walk_records
-- ----------------------------
DROP TABLE IF EXISTS walk_records CASCADE;
CREATE TABLE walk_records (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL,
  sub_order_id BIGINT NOT NULL,
  path JSONB DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE walk_records IS '遛狗轨迹记录表';
COMMENT ON COLUMN walk_records.id IS '轨迹ID';
COMMENT ON COLUMN walk_records.order_id IS '主订单ID';
COMMENT ON COLUMN walk_records.sub_order_id IS '子订单ID';
COMMENT ON COLUMN walk_records.path IS '路径数据 LatLng数组';
COMMENT ON COLUMN walk_records.created_at IS '创建时间';
COMMENT ON COLUMN walk_records.updated_at IS '更新时间';

CREATE INDEX idx_walk_records_order_id ON walk_records (order_id);
CREATE INDEX idx_walk_records_sub_order_id ON walk_records (sub_order_id);
CREATE INDEX idx_walk_records_created_at ON walk_records (created_at);

CREATE TRIGGER update_walk_records_modtime BEFORE UPDATE ON walk_records FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- ----------------------------
-- Table structure for wallet_transaction
-- ----------------------------
DROP TABLE IF EXISTS wallet_transaction CASCADE;
CREATE TABLE wallet_transaction (
  id BIGSERIAL PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  transaction_type SMALLINT NOT NULL,
  order_type INTEGER NOT NULL,
  order_id BIGINT DEFAULT NULL,
  order_created_at TIMESTAMP DEFAULT NULL,
  amount DECIMAL(12,2) NOT NULL,
  balance_after DECIMAL(12,2) NOT NULL,
  transaction_time TIMESTAMP NOT NULL,
  remark VARCHAR(255) DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE wallet_transaction IS '钱包收支记录明细表';
COMMENT ON COLUMN wallet_transaction.id IS '主键ID';
COMMENT ON COLUMN wallet_transaction.user_id IS '用户ID';
COMMENT ON COLUMN wallet_transaction.transaction_type IS '交易类型：1-收入，2-支出';
COMMENT ON COLUMN wallet_transaction.order_type IS '订单类型（如：咨询服务、会员订阅等）';
COMMENT ON COLUMN wallet_transaction.order_id IS '关联订单ID（可选）';
COMMENT ON COLUMN wallet_transaction.order_created_at IS '订单开始时间';
COMMENT ON COLUMN wallet_transaction.amount IS '交易金额（元，正数）';
COMMENT ON COLUMN wallet_transaction.balance_after IS '交易后余额（元）';
COMMENT ON COLUMN wallet_transaction.transaction_time IS '交易时间';
COMMENT ON COLUMN wallet_transaction.remark IS '交易备注';
COMMENT ON COLUMN wallet_transaction.created_at IS '记录创建时间';

CREATE INDEX idx_wt_user_type_time ON wallet_transaction (user_id, transaction_type, transaction_time);
CREATE INDEX idx_wt_transaction_time ON wallet_transaction (transaction_time);
