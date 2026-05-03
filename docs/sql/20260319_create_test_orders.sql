-- 测试订单数据
-- 创建日期：2026-03-19
-- 包含5个订单：3个遛狗订单，2个其他类型
-- 状态：OrderPayed(1) 已支付，sitter未接受

-- 清除旧数据（可选）
-- DELETE FROM sub_orders WHERE order_id IN (SELECT id FROM orders WHERE order_number LIKE 'ORDTEST%');
-- DELETE FROM order_pet WHERE order_id IN (SELECT id FROM orders WHERE order_number LIKE 'ORDTEST%');
-- DELETE FROM orders WHERE order_number LIKE 'ORDTEST%';

-- ==========================================
-- 订单1：遛狗订单（3天）
-- ==========================================
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
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    3, -- Walking（遛狗）
    '2026-03-20',
    '2026-03-22',
    75.00,
    75.00,
    1, -- OrderPayed（已支付）
    '张三',
    '李四',
    'ORDTEST202603190001',
    NOW(),
    NOW()
);
SET @order1_id = LAST_INSERT_ID();

-- 订单1的宠物
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
    @order1_id,
    'p_test001',
    'Dog',
    2, -- 中型
    'Lucky',
    'Corgi',
    75.00,
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    NOW(),
    NOW()
);

-- 订单1的子订单（3天）
INSERT INTO sub_orders (
    order_id,
    date,
    state,
    start_code,
    end_code,
    created_at,
    updated_at
) VALUES 
(
    @order1_id,
    '2026-03-20',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    @order1_id,
    '2026-03-21',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    @order1_id,
    '2026-03-22',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
);

-- ==========================================
-- 订单2：遛狗订单（5天）
-- ==========================================
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
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    3, -- Walking（遛狗）
    '2026-03-25',
    '2026-03-29',
    125.00,
    125.00,
    1, -- OrderPayed（已支付）
    '王五',
    '赵六',
    'ORDTEST202603190002',
    NOW(),
    NOW()
);
SET @order2_id = LAST_INSERT_ID();

-- 订单2的宠物
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
    @order2_id,
    'p_test002',
    'Dog',
    1, -- 小型
    'Buddy',
    'Poodle',
    125.00,
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    NOW(),
    NOW()
);

-- 订单2的子订单（5天）
INSERT INTO sub_orders (
    order_id,
    date,
    state,
    start_code,
    end_code,
    created_at,
    updated_at
) VALUES 
(
    @order2_id,
    '2026-03-25',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    @order2_id,
    '2026-03-26',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    @order2_id,
    '2026-03-27',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    @order2_id,
    '2026-03-28',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    @order2_id,
    '2026-03-29',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
);

-- ==========================================
-- 订单3：遛狗订单（2天）
-- ==========================================
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
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    3, -- Walking（遛狗）
    '2026-04-01',
    '2026-04-02',
    50.00,
    50.00,
    1, -- OrderPayed（已支付）
    '钱七',
    '孙八',
    'ORDTEST202603190003',
    NOW(),
    NOW()
);
SET @order3_id = LAST_INSERT_ID();

-- 订单3的宠物
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
    @order3_id,
    'p_test003',
    'Dog',
    3, -- 大型
    'Max',
    'Golden Retriever',
    50.00,
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    NOW(),
    NOW()
);

-- 订单3的子订单（2天）
INSERT INTO sub_orders (
    order_id,
    date,
    state,
    start_code,
    end_code,
    created_at,
    updated_at
) VALUES 
(
    @order3_id,
    '2026-04-01',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
),
(
    @order3_id,
    '2026-04-02',
    1, -- OrderPayed
    NULL,
    NULL,
    NOW(),
    NOW()
);

-- ==========================================
-- 订单4：寄养订单（非遛狗）
-- ==========================================
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
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    1, -- Boarding（寄养）
    '2026-03-20',
    '2026-03-25',
    300.00,
    300.00,
    1, -- OrderPayed（已支付）
    '周九',
    '吴十',
    'ORDTEST202603190004',
    NOW(),
    NOW()
);
SET @order4_id = LAST_INSERT_ID();

-- 订单4的宠物
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
    @order4_id,
    'p_test004',
    'Cat',
    1, -- 小型
    'Mimi',
    'British Shorthair',
    300.00,
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    NOW(),
    NOW()
);

-- ==========================================
-- 订单5：日托订单（非遛狗）
-- ==========================================
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
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    2, -- Daycare（日托）
    '2026-03-21',
    '2026-03-23',
    105.00,
    105.00,
    1, -- OrderPayed（已支付）
    '郑十一',
    '王十二',
    'ORDTEST202603190005',
    NOW(),
    NOW()
);
SET @order5_id = LAST_INSERT_ID();

-- 订单5的宠物
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
    @order5_id,
    'p_test005',
    'Dog',
    2, -- 中型
    'Charlie',
    'Beagle',
    105.00,
    'bcdd55b8-f081-7043-b733-d39cb17afe05',
    '5c8dd518-0031-70a5-dcf6-05574edfc787',
    NOW(),
    NOW()
);

-- ==========================================
-- 查询验证数据
-- ==========================================
SELECT '========== 主订单信息 ==========' AS info;
SELECT 
    id,
    order_number,
    type,
    CASE type
        WHEN 1 THEN 'Boarding（寄养）'
        WHEN 2 THEN 'Daycare（日托）'
        WHEN 3 THEN 'Walking（遛狗）'
        WHEN 4 THEN 'DropIn（上门）'
        ELSE 'Unknown'
    END AS type_name,
    from_date,
    to_date,
    state,
    CASE state
        WHEN 0 THEN 'OrderInitialized（初始化）'
        WHEN 1 THEN 'OrderPayed（已支付）'
        WHEN 2 THEN 'OrderAccepted（已接受）'
        ELSE 'Other'
    END AS state_name,
    owner_name,
    sitter_name,
    total_price
FROM orders WHERE order_number LIKE 'ORDTEST%' ORDER BY id;

SELECT '========== 子订单信息（仅遛狗订单）==========' AS info;
SELECT 
    so.id,
    so.order_id,
    o.order_number,
    so.date,
    so.state,
    CASE so.state
        WHEN 0 THEN 'OrderInitialized（初始化）'
        WHEN 1 THEN 'OrderPayed（已支付）'
        WHEN 2 THEN 'OrderAccepted（已接受）'
        ELSE 'Other'
    END AS state_name
FROM sub_orders so
JOIN orders o ON so.order_id = o.id
WHERE o.order_number LIKE 'ORDTEST%'
ORDER BY so.order_id, so.date;

SELECT '========== 订单宠物信息 ==========' AS info;
SELECT 
    op.id,
    op.order_id,
    o.order_number,
    op.pet_id,
    op.pet_name,
    op.pet_type,
    op.pet_price
FROM order_pet op
JOIN orders o ON op.order_id = o.id
WHERE o.order_number LIKE 'ORDTEST%'
ORDER BY op.order_id;
