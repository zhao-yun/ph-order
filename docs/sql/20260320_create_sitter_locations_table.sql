-- 创建Sitter实时位置表
-- 创建日期：2026-03-20

CREATE TABLE IF NOT EXISTS sitter_locations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '位置记录ID',
    sitter_id VARCHAR(50) NOT NULL COMMENT 'Sitter ID',
    order_id BIGINT NOT NULL COMMENT '主订单ID',
    sub_order_id BIGINT NOT NULL COMMENT '子订单ID',
    lat DECIMAL(10,6) NOT NULL COMMENT '纬度',
    lng DECIMAL(10,6) NOT NULL COMMENT '经度',
    timestamp BIGINT NOT NULL COMMENT '时间戳',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX idx_sitter_id (sitter_id),
    INDEX idx_order_id (order_id),
    INDEX idx_sub_order_id (sub_order_id),
    INDEX idx_timestamp (timestamp),
    INDEX idx_sitter_suborder (sitter_id, sub_order_id),
    INDEX idx_order_suborder (order_id, sub_order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Sitter实时位置记录表';

-- 可选：添加外键约束
-- ALTER TABLE sitter_locations ADD CONSTRAINT fk_sitter_locations_order_id FOREIGN KEY (order_id) REFERENCES orders(id);
-- ALTER TABLE sitter_locations ADD CONSTRAINT fk_sitter_locations_sub_order_id FOREIGN KEY (sub_order_id) REFERENCES sub_orders(id);

-- 查看表结构
-- DESC sitter_locations;

-- 清理旧数据（谨慎使用）
-- DELETE FROM sitter_locations WHERE created_at < DATE_SUB(NOW(), INTERVAL 1 DAY);
