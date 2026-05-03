/*
 Navicat MySQL Data Transfer

 Source Server         : 阿里云RDS
 Source Server Type    : MySQL
 Source Server Version : 80036
 Source Host           : rm-bp1acok94fvb6jmxzio.mysql.rds.aliyuncs.com:3306
 Source Schema         : ph_orders

 Target Server Type    : MySQL
 Target Server Version : 80036
 File Encoding         : 65001

 Date: 04/05/2026 01:02:47
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for bank_card
-- ----------------------------
DROP TABLE IF EXISTS `bank_card`;
CREATE TABLE `bank_card` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
  `user_id` varchar(50) NOT NULL COMMENT 'User ID',
  `card_number` varchar(20) NOT NULL COMMENT 'Bank card number',
  `card_type` tinyint(1) NOT NULL COMMENT 'Card type (e.g., 1: debit, 2: credit)',
  `bank_code` varchar(10) NOT NULL COMMENT 'Bank code',
  `bank_name` varchar(50) NOT NULL COMMENT 'Bank name',
  `interbank_transfer_code` varchar(20) DEFAULT NULL COMMENT 'Interbank transfer code (e.g., SWIFT/IFSC)',
  `account_holder` varchar(50) NOT NULL COMMENT 'Account holder name',
  `is_default` tinyint(1) DEFAULT '0' COMMENT 'Is default card (0: no, 1: yes)',
  `status` tinyint(1) DEFAULT '1' COMMENT 'Status (0: inactive, 1: active)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_card_number` (`card_number`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Bank card information table';

-- ----------------------------
-- Table structure for interview_record
-- ----------------------------
DROP TABLE IF EXISTS `interview_record`;
CREATE TABLE `interview_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '面试预约记录唯一ID',
  `order_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联的订单ID，用于绑定业务订单',
  `initiator_type` tinyint NOT NULL COMMENT '预约发起方类型：1-用户 2-Sitter',
  `interview_type` tinyint NOT NULL COMMENT '面试类型：1-线上 2-离线会议',
  `appointment_time` datetime DEFAULT NULL COMMENT '预约面试时间',
  `location` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '面试地点（线上则为空）',
  `message` text COLLATE utf8mb4_unicode_ci COMMENT '预约时填写的消息内容',
  `status` tinyint NOT NULL DEFAULT (0) COMMENT '预约状态：1-待确认（刚发起） 2-接受方修改待确认 3-已接受 4-已取消 5-已拒绝',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '预约创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  `user_result` tinyint DEFAULT NULL,
  `sitter_result` tinyint DEFAULT NULL,
  `user_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sitter_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=171 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for order_modification_log
-- ----------------------------
DROP TABLE IF EXISTS `order_modification_log`;
CREATE TABLE `order_modification_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `order_id` bigint NOT NULL COMMENT '关联的订单ID',
  `owner_id` varchar(50) NOT NULL COMMENT '用户ID',
  `sitter_id` varchar(50) NOT NULL COMMENT 'Sitter ID',
  `previous_date` date DEFAULT NULL COMMENT '修改前的日期',
  `new_date` date DEFAULT NULL COMMENT '修改后的日期',
  `previous_pet_list` text COMMENT '修改前的宠物列表(JSON格式)',
  `new_pet_list` text COMMENT '修改后的宠物列表(JSON格式)',
  `previous_price` decimal(10,2) DEFAULT NULL COMMENT '修改前的价格',
  `new_price` decimal(10,2) DEFAULT NULL COMMENT '修改后的价格',
  `state` int NOT NULL DEFAULT '0' COMMENT '修改状态',
  `type` int NOT NULL COMMENT '修改类型',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_owner_id` (`owner_id`),
  KEY `idx_sitter_id` (`sitter_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_order_modification_log_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单修改日志表';

-- ----------------------------
-- Table structure for order_pet
-- ----------------------------
DROP TABLE IF EXISTS `order_pet`;
CREATE TABLE `order_pet` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_id` bigint NOT NULL COMMENT '订单ID',
  `pet_id` varchar(255) NOT NULL COMMENT '宠物ID',
  `pet_type` varchar(255) NOT NULL COMMENT '宠物类型，如猫、狗',
  `pet_shape` int NOT NULL COMMENT '宠物体型 1为小型犬 2为中型犬 3为大型犬',
  `pet_name` varchar(255) DEFAULT NULL COMMENT '宠物名字',
  `breed` varchar(255) DEFAULT NULL COMMENT '宠物品种',
  `pet_price` decimal(10,2) NOT NULL COMMENT '宠物费用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `owner_id` varchar(255) DEFAULT NULL,
  `sitter_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_pet_id` (`pet_id`),
  KEY `idx_pet_type` (`pet_type`),
  KEY `idx_pet_shape` (`pet_shape`),
  CONSTRAINT `fk_order_pet_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=60 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单宠物关联表';

-- ----------------------------
-- Table structure for orders
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `owner_id` varchar(50) NOT NULL COMMENT '用户ID',
  `sitter_id` varchar(50) NOT NULL COMMENT 'Sitter ID',
  `type` bigint NOT NULL COMMENT '订单类型',
  `from_date` date NOT NULL COMMENT '开始时间',
  `to_date` date NOT NULL COMMENT '结束时间',
  `tips_price` decimal(10,2) DEFAULT NULL COMMENT '小费',
  `total_price` decimal(10,2) DEFAULT NULL COMMENT '总费用',
  `sub_total_price` decimal(10,2) DEFAULT NULL COMMENT '宠物总费用',
  `service_fee` decimal(10,2) DEFAULT NULL COMMENT '服务费',
  `taxes` decimal(10,2) DEFAULT NULL COMMENT '税',
  `state` int NOT NULL COMMENT '订单状态',
  `sitter_handle_at` timestamp NULL DEFAULT NULL COMMENT 'Sitter 处理时间',
  `sitter_finish_at` timestamp NULL DEFAULT NULL COMMENT 'Sitter 完成时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `owner_name` varchar(255) DEFAULT NULL,
  `sitter_name` varchar(255) DEFAULT NULL,
  `contact` varchar(255) DEFAULT NULL,
  `alternative_contact` varchar(255) DEFAULT NULL,
  `note` text,
  `cancel_at` timestamp NULL DEFAULT NULL,
  `cancel_reason` varchar(255) DEFAULT NULL,
  `refund_price` decimal(10,2) DEFAULT NULL,
  `order_number` varchar(255) DEFAULT NULL,
  `user_deleted` tinyint NOT NULL DEFAULT '0',
  `sitter_deleted` tinyint NOT NULL DEFAULT '0',
  `user_rating_state` tinyint DEFAULT '0',
  `sitter_rating_state` tinyint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_owner_id` (`owner_id`),
  KEY `idx_sitter_id` (`sitter_id`),
  KEY `idx_state` (`state`),
  KEY `idx_dates` (`from_date`,`to_date`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单表';

-- ----------------------------
-- Table structure for owner_rating
-- ----------------------------
DROP TABLE IF EXISTS `owner_rating`;
CREATE TABLE `owner_rating` (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` varchar(50) DEFAULT NULL,
  `owner_id` varchar(50) DEFAULT NULL,
  `sitter_id` varchar(50) DEFAULT NULL,
  `owner_name` varchar(50) DEFAULT NULL,
  `sitter_name` varchar(50) DEFAULT NULL,
  `score` tinyint DEFAULT NULL,
  `satisfaction_level` tinyint DEFAULT NULL,
  `instructions_clarity` tinyint DEFAULT NULL,
  `communication` tinyint DEFAULT NULL,
  `supplies_preparation` tinyint DEFAULT NULL,
  `respect_courtesy` tinyint DEFAULT NULL,
  `suggestions` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Table structure for pet_rating
-- ----------------------------
DROP TABLE IF EXISTS `pet_rating`;
CREATE TABLE `pet_rating` (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` varchar(50) DEFAULT NULL,
  `pet_id` varchar(50) DEFAULT NULL,
  `pet_name` varchar(50) DEFAULT NULL,
  `score` tinyint DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Table structure for report
-- ----------------------------
DROP TABLE IF EXISTS `report`;
CREATE TABLE `report` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '举报ID',
  `reporter_id` varchar(50) NOT NULL COMMENT '举报人ID',
  `reported_id` varchar(50) NOT NULL COMMENT '被举报人ID',
  `reported_name` varchar(50) DEFAULT NULL COMMENT '被举报人姓名',
  `report_type` varchar(50) DEFAULT NULL COMMENT '举报类型：1-欺诈，2-骚扰，3-不当内容，4-其他',
  `report_reason` varchar(50) DEFAULT NULL COMMENT '举报原因',
  `report_desc` text NOT NULL COMMENT '举报内容',
  `report_img` text COMMENT '证据(图片/视频等)',
  `status` tinyint(1) DEFAULT '0' COMMENT '处理状态：0-待处理，1-处理中，2-已处理',
  `handler_id` varchar(50) DEFAULT NULL COMMENT '处理人ID',
  `handle_result` text COMMENT '处理结果',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_reporter_id` (`reporter_id`) COMMENT '举报人索引',
  KEY `idx_reported_id` (`reported_id`) COMMENT '被举报人索引',
  KEY `idx_status` (`status`) COMMENT '状态索引'
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户举报记录表';

-- ----------------------------
-- Table structure for sitter_locations
-- ----------------------------
DROP TABLE IF EXISTS `sitter_locations`;
CREATE TABLE `sitter_locations` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '位置记录ID',
  `sitter_id` varchar(50) NOT NULL COMMENT 'Sitter ID',
  `order_id` bigint NOT NULL COMMENT '主订单ID',
  `sub_order_id` bigint NOT NULL COMMENT '子订单ID',
  `lat` decimal(10,6) NOT NULL COMMENT '纬度',
  `lng` decimal(10,6) NOT NULL COMMENT '经度',
  `timestamp` bigint NOT NULL COMMENT '时间戳',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_sitter_id` (`sitter_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_sub_order_id` (`sub_order_id`),
  KEY `idx_timestamp` (`timestamp`),
  KEY `idx_sitter_suborder` (`sitter_id`,`sub_order_id`),
  KEY `idx_order_suborder` (`order_id`,`sub_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Sitter实时位置记录表';

-- ----------------------------
-- Table structure for sitter_rating
-- ----------------------------
DROP TABLE IF EXISTS `sitter_rating`;
CREATE TABLE `sitter_rating` (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` int NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `sitter_id` varchar(255) NOT NULL,
  `punctuality` tinyint DEFAULT NULL,
  `responsibility` tinyint DEFAULT NULL,
  `communication` tinyint DEFAULT NULL,
  `pet_care_skills` tinyint DEFAULT NULL,
  `cleanliness` tinyint DEFAULT NULL,
  `suggestions` text,
  `score` tinyint DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `order_id` (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Table structure for sub_orders
-- ----------------------------
DROP TABLE IF EXISTS `sub_orders`;
CREATE TABLE `sub_orders` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '子订单ID',
  `order_id` bigint NOT NULL COMMENT '主订单ID',
  `date` date NOT NULL COMMENT '日期',
  `state` int NOT NULL COMMENT '子订单状态',
  `start_code` varchar(10) DEFAULT NULL COMMENT '开始验证码',
  `end_code` varchar(10) DEFAULT NULL COMMENT '结束验证码',
  `sitter_handle_at` timestamp NULL DEFAULT NULL COMMENT 'Sitter处理时间',
  `sitter_finish_at` timestamp NULL DEFAULT NULL COMMENT 'Sitter完成时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `walk_thumbnail_url` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_date` (`date`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='子订单表';

-- ----------------------------
-- Table structure for user_wallet
-- ----------------------------
DROP TABLE IF EXISTS `user_wallet`;
CREATE TABLE `user_wallet` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` varchar(255) NOT NULL COMMENT '用户ID（关联用户表）',
  `balance` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '账户余额（元）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '钱包状态：1-正常，2-冻结',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`) COMMENT '用户ID唯一索引'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户钱包主表';

-- ----------------------------
-- Table structure for walk_records
-- ----------------------------
DROP TABLE IF EXISTS `walk_records`;
CREATE TABLE `walk_records` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '轨迹ID',
  `order_id` bigint NOT NULL COMMENT '主订单ID',
  `sub_order_id` bigint NOT NULL COMMENT '子订单ID',
  `path` json DEFAULT NULL COMMENT '路径数据 LatLng数组',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_sub_order_id` (`sub_order_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='遛狗轨迹记录表';

-- ----------------------------
-- Table structure for wallet_transaction
-- ----------------------------
DROP TABLE IF EXISTS `wallet_transaction`;
CREATE TABLE `wallet_transaction` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` varchar(255) NOT NULL COMMENT '用户ID',
  `transaction_type` tinyint NOT NULL COMMENT '交易类型：1-收入，2-支出',
  `order_type` int NOT NULL COMMENT '订单类型（如：咨询服务、会员订阅等）',
  `order_id` bigint DEFAULT NULL COMMENT '关联订单ID（可选）',
  `order_created_at` datetime DEFAULT NULL COMMENT '订单开始时间',
  `amount` decimal(12,2) NOT NULL COMMENT '交易金额（元，正数）',
  `balance_after` decimal(12,2) NOT NULL COMMENT '交易后余额（元）',
  `transaction_time` datetime NOT NULL COMMENT '交易时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '交易备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_type_time` (`user_id`,`transaction_type`,`transaction_time`) COMMENT '用户+交易类型+时间索引，优化列表查询',
  KEY `idx_transaction_time` (`transaction_time`) COMMENT '交易时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='钱包收支记录明细表';

SET FOREIGN_KEY_CHECKS = 1;
