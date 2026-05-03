-- 遛狗订单和子订单示例数据
-- 创建日期：2026-03-14

-- 1. 先创建一个主订单（遛狗订单）
INSERT INTO orders (
    owner_id,
    sitter_id,
    type,
    from_date,
    to_date,
    total_price,
    sub_total_price,
    state,
    owner_name,
    sitter_name,
    order_number,
    created_at,
    updated_at
) VALUES (
    'u_123',
    's_456',
    3, -- Walking（遛狗）
    '2026-03-14',
    '2026-03-16',
    75.00,
    75.00,
    2, -- OrderAccepted（已接受）
    '张三',
    '李四',
    'ORD20260314120000abc12345',
    NOW(),
    NOW()
);

-- 获取刚插入的主订单ID（假设是100），实际执行时需要使用真实的ID
SET @order_id = LAST_INSERT_ID();

-- 2. 为这个订单添加宠物
INSERT INTO order_pet (
    order_id,
    pet_id,
    pet_type,
    pet_shape,
    pet_name,
    breed,
    pet_price,
    owner_id,
    sitter_id,
    created_at,
    updated_at
) VALUES 
(
    @order_id,
    'p_001',
    'Dog',
    2, -- 中型
    'Lucky',
    'Corgi',
    75.00,
    'u_123',
    's_456',
    NOW(),
    NOW()
);

-- 3. 创建对应的子订单（3天：2026-03-14, 2026-03-15, 2026-03-16）
INSERT INTO sub_orders (
    order_id,
    date,
    state,
    start_code,
    end_code,
    created_at,
    updated_at
) VALUES 
-- 第1天：2026-03-14（已完成）
(
    @order_id,
    '2026-03-14',
    5, -- OrderCompleted（已完成）
    '1234',
    '5678',
    NOW(),
    NOW()
),
-- 第2天：2026-03-15（订单开始）
(
    @order_id,
    '2026-03-15',
    3, -- OrderEstablished（订单开始）
    '2345',
    '6789',
    NOW(),
    NOW()
),
-- 第3天：2026-03-16（未开始/已接受）
(
    @order_id,
    '2026-03-16',
    2, -- OrderAccepted（未开始/已接受）
    '3456',
    NULL,
    NOW(),
    NOW()
);

-- 4. 查询插入的结果，验证数据是否正确
SELECT '主订单信息:' AS info;
SELECT 
    id,
    order_number,
    type,
    from_date,
    to_date,
    state,
    owner_name,
    sitter_name
FROM orders WHERE id = @order_id;

SELECT '订单宠物信息:' AS info;
SELECT 
    id,
    order_id,
    pet_id,
    pet_name,
    pet_type
FROM order_pet WHERE order_id = @order_id;

SELECT '子订单信息:' AS info;
SELECT 
    id,
    order_id,
    date,
    state,
    start_code,
    end_code
FROM sub_orders WHERE order_id = @order_id ORDER BY date ASC;
