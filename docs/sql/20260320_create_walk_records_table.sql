-- 创建遛狗轨迹表
-- 创建日期：2026-03-20

CREATE TABLE IF NOT EXISTS walk_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '轨迹ID',
    order_id BIGINT NOT NULL COMMENT '主订单ID',
    sub_order_id BIGINT NOT NULL COMMENT '子订单ID',
    path JSON COMMENT '路径数据 LatLng数组',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_order_id (order_id),
    INDEX idx_sub_order_id (sub_order_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='遛狗轨迹记录表';

-- 添加外键约束（可选，根据实际情况）
-- ALTER TABLE walk_records ADD CONSTRAINT fk_walk_records_order_id FOREIGN KEY (order_id) REFERENCES orders(id);
-- ALTER TABLE walk_records ADD CONSTRAINT fk_walk_records_sub_order_id FOREIGN KEY (sub_order_id) REFERENCES sub_orders(id);
