-- ============================================================================
-- What To Eat 数据库初始化脚本
-- ============================================================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS what_to_eat 
  DEFAULT CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;

USE what_to_eat;

-- ============================================================================
-- 表结构定义
-- ============================================================================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 餐厅表
CREATE TABLE IF NOT EXISTS restaurants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE COMMENT '餐厅名称',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='餐厅表';

-- 菜单表（菜品）- 支持软删除
CREATE TABLE IF NOT EXISTS menus (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    restaurant_id BIGINT NOT NULL COMMENT '所属餐厅ID',
    dish_name VARCHAR(100) NOT NULL COMMENT '菜品名称',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL COMMENT '软删除时间',
    INDEX idx_restaurant (restaurant_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- 决策记录表
CREATE TABLE IF NOT EXISTS decision_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    menu_id BIGINT NOT NULL COMMENT '菜单ID',
    decided_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '决策时间',
    INDEX idx_user_decided (user_id, decided_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='决策记录表';

-- ============================================================================
-- 默认数据（可选，后端启动时会自动初始化）
-- ============================================================================

-- 插入游客用户（id=1）
-- 密码哈希为 "guest123" 使用 bcrypt 生成
INSERT INTO users (id, username, password_hash, created_at, updated_at) 
VALUES (1, 'guest', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE username = username;

-- 插入默认餐厅
INSERT IGNORE INTO restaurants (name) VALUES 
    ('麦当劳'),
    ('肯德基'),
    ('沙县小吃'),
    ('兰州拉面'),
    ('黄焖鸡米饭'),
    ('海底捞');

-- 插入默认菜品
INSERT IGNORE INTO menus (restaurant_id, dish_name) 
SELECT r.id, d.dish_name FROM restaurants r
JOIN (
    SELECT '麦当劳' as restaurant, '巨无霸' as dish_name UNION ALL
    SELECT '麦当劳', '麦辣鸡腿堡' UNION ALL
    SELECT '麦当劳', '薯条' UNION ALL
    SELECT '麦当劳', '麦旋风' UNION ALL
    SELECT '肯德基', '原味鸡' UNION ALL
    SELECT '肯德基', '香辣鸡腿堡' UNION ALL
    SELECT '肯德基', '上校鸡块' UNION ALL
    SELECT '肯德基', '蛋挞' UNION ALL
    SELECT '沙县小吃', '拌面' UNION ALL
    SELECT '沙县小吃', '蒸饺' UNION ALL
    SELECT '沙县小吃', '炖罐' UNION ALL
    SELECT '沙县小吃', '扁肉' UNION ALL
    SELECT '兰州拉面', '牛肉面' UNION ALL
    SELECT '兰州拉面', '刀削面' UNION ALL
    SELECT '兰州拉面', '炒面' UNION ALL
    SELECT '兰州拉面', '凉面' UNION ALL
    SELECT '黄焖鸡米饭', '黄焖鸡米饭' UNION ALL
    SELECT '黄焖鸡米饭', '黄焖排骨' UNION ALL
    SELECT '黄焖鸡米饭', '黄焖豆腐' UNION ALL
    SELECT '海底捞', '麻辣锅底' UNION ALL
    SELECT '海底捞', '番茄锅底' UNION ALL
    SELECT '海底捞', '菌汤锅底'
) d ON r.name = d.restaurant
WHERE NOT EXISTS (
    SELECT 1 FROM menus m WHERE m.restaurant_id = r.id AND m.dish_name = d.dish_name
);
