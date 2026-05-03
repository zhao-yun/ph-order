-- 创建子订单表
-- 创建日期：2026-03-14

-- 如果表已存在则删除
DROP TABLE IF EXISTS `sub_orders`;

-- 创建子订单表
CREATE TABLE `sub_orders` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '子订单ID',
    `order_id` BIGINT NOT NULL COMMENT '主订单ID',
    `date` DATE NOT NULL COMMENT '日期',
    `state` INT NOT NULL COMMENT '子订单状态',
    `start_code` VARCHAR(10) DEFAULT NULL COMMENT '开始验证码',
    `end_code` VARCHAR(10) DEFAULT NULL COMMENT '结束验证码',
    `sitter_handle_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Sitter处理时间',
    `sitter_finish_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Sitter完成时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_order_id` (`order_id`),
    INDEX `idx_date` (`date`),
    INDEX `idx_state` (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='子订单表';

-- 验证表是否创建成功
SHOW CREATE TABLE `sub_orders`;
