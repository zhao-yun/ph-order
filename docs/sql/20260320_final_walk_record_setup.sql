-- ============================================
-- 遛狗轨迹模块 - 最终SQL脚本
-- 创建日期：2026-03-20
-- ============================================

-- ============================================
-- 1. 创建遛狗轨迹表
-- ============================================
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

-- ============================================
-- 2. 给子订单表添加轨迹缩略图URL字段
-- ============================================
DELIMITER //
CREATE PROCEDURE add_walk_thumbnail_column_if_not_exists()
BEGIN
    IF NOT EXISTS (
        SELECT * FROM information_schema.COLUMNS 
        WHERE TABLE_SCHEMA = DATABASE() 
        AND TABLE_NAME = 'sub_orders' 
        AND COLUMN_NAME = 'walk_thumbnail_url'
    ) THEN
        ALTER TABLE sub_orders 
        ADD COLUMN walk_thumbnail_url VARCHAR(500) 
        COMMENT '轨迹缩略图URL' 
        AFTER sitter_finish_at;
    END IF;
END //
DELIMITER ;

CALL add_walk_thumbnail_column_if_not_exists();
DROP PROCEDURE add_walk_thumbnail_column_if_not_exists;

-- ============================================
-- 验证表结构
-- ============================================
-- 查看轨迹表结构
-- DESC walk_records;

-- 查看子订单表结构
-- DESC sub_orders;

-- ============================================
-- 可选：添加外键约束
-- ============================================
-- ALTER TABLE walk_records ADD CONSTRAINT fk_walk_records_order_id FOREIGN KEY (order_id) REFERENCES orders(id);
-- ALTER TABLE walk_records ADD CONSTRAINT fk_walk_records_sub_order_id FOREIGN KEY (sub_order_id) REFERENCES sub_orders(id);

-- ============================================
-- 完成！
-- ============================================
-- 执行完此脚本后，遛狗轨迹模块的数据库表结构就准备好了
-- 接下来可以使用相关的接口：
-- - POST /WalkRecord/AppendWalkPath - 追加轨迹点
-- - GET /WalkRecord/GetWalkRecordBySubOrderId - 查询轨迹
-- - POST /SubOrder/UpdateWalkThumbnail - 更新缩略图
