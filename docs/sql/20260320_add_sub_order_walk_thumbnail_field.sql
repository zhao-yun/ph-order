-- 给子订单表添加轨迹缩略图URL字段
-- 创建日期：2026-03-20

-- 检查字段是否存在的存储过程
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

-- 如果需要设置默认值，可以使用：
-- ALTER TABLE sub_orders ALTER COLUMN walk_thumbnail_url SET DEFAULT '';

-- 查看表结构验证
-- DESC sub_orders;
