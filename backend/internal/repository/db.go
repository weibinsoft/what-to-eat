package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"what-to-eat/config"
	"what-to-eat/internal/model"
	"what-to-eat/pkg/logger"
)

var DB *gorm.DB

// ============================================================================
// 数据库初始化
// ============================================================================

// InitDB 初始化数据库连接
// 执行顺序：1. 确保数据库存在 -> 2. 连接数据库 -> 3. 自动迁移表结构 -> 4. 初始化默认数据
func InitDB(cfg *config.DatabaseConfig) error {
	// 1. 确保数据库存在（不存在则创建）
	if err := ensureDatabase(cfg); err != nil {
		return fmt.Errorf("failed to ensure database: %w", err)
	}

	// 2. 连接到指定数据库
	if err := connectDatabase(cfg); err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 3. 自动迁移表结构
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 4. 初始化默认数据（仅首次启动）
	if err := initDefaultData(); err != nil {
		return fmt.Errorf("failed to init default data: %w", err)
	}

	logger.Info("Database initialized successfully",
		zap.String("host", cfg.Host),
		zap.String("database", cfg.DBName),
	)
	return nil
}

// ensureDatabase 确保数据库存在，不存在则创建
func ensureDatabase(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Discard,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// 创建数据库（如果不存在）
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	if err := db.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	logger.Info("Database ensured", zap.String("database", cfg.DBName))

	// 关闭临时连接
	sqlDB, _ := db.DB()
	sqlDB.Close()

	return nil
}

// connectDatabase 连接到指定数据库
func connectDatabase(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	gormLogger := &zapGormLogger{SlowThreshold: time.Second}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})
	if err != nil {
		return err
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return nil
}

// autoMigrate 自动迁移表结构
// 表创建顺序：users -> restaurants -> menus -> decision_records
func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Restaurant{},
		&model.Menu{},
		&model.DecisionRecord{},
	)
}

// ============================================================================
// 默认数据初始化
// ============================================================================

// initDefaultData 初始化默认数据（仅当数据库为空时执行）
func initDefaultData() error {
	// 先初始化游客用户
	if err := initGuestUser(); err != nil {
		return fmt.Errorf("failed to init guest user: %w", err)
	}

	// 检查是否已有菜单数据
	var menuCount int64
	if err := DB.Model(&model.Menu{}).Count(&menuCount).Error; err != nil {
		return err
	}

	if menuCount > 0 {
		logger.Debug("Default data already exists, skipping initialization")
		return nil
	}

	// 默认餐厅和菜品数据
	defaultData := []struct {
		Restaurant string
		Dishes     []string
	}{
		{"麦当劳", []string{"巨无霸", "麦辣鸡腿堡", "薯条", "麦旋风"}},
		{"肯德基", []string{"原味鸡", "香辣鸡腿堡", "上校鸡块", "蛋挞"}},
		{"沙县小吃", []string{"拌面", "蒸饺", "炖罐", "扁肉"}},
		{"兰州拉面", []string{"牛肉面", "刀削面", "炒面", "凉面"}},
		{"黄焖鸡米饭", []string{"黄焖鸡米饭", "黄焖排骨", "黄焖豆腐"}},
		{"海底捞", []string{"麻辣锅底", "番茄锅底", "菌汤锅底"}},
	}

	// 使用事务批量插入
	return DB.Transaction(func(tx *gorm.DB) error {
		for _, data := range defaultData {
			// 创建餐厅
			restaurant := model.Restaurant{Name: data.Restaurant}
			if err := tx.Create(&restaurant).Error; err != nil {
				return fmt.Errorf("failed to create restaurant %s: %w", data.Restaurant, err)
			}

			// 批量创建菜品
			menus := make([]model.Menu, len(data.Dishes))
			for i, dish := range data.Dishes {
				menus[i] = model.Menu{
					RestaurantID: restaurant.ID,
					DishName:     dish,
				}
			}
			if err := tx.Create(&menus).Error; err != nil {
				return fmt.Errorf("failed to create menus for %s: %w", data.Restaurant, err)
			}
		}

		logger.Info("Initialized default data",
			zap.Int("restaurants", len(defaultData)),
		)
		return nil
	})
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// ============================================================================
// GORM 日志适配器（使用 zap）
// ============================================================================

type zapGormLogger struct {
	SlowThreshold time.Duration
}

func (l *zapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.Info(fmt.Sprintf(msg, data...))
}

func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.Warn(fmt.Sprintf(msg, data...))
}

func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.Error(fmt.Sprintf(msg, data...))
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("GORM SQL Error", append(fields, zap.Error(err))...)
		return
	}

	if elapsed > l.SlowThreshold && l.SlowThreshold > 0 {
		logger.Warn("GORM Slow SQL", fields...)
		return
	}

	logger.Debug("GORM SQL", fields...)
}

// ============================================================================
// 游客用户初始化
// ============================================================================

// GuestUserID 游客用户ID
const GuestUserID int64 = 1

// GuestUsername 游客用户名
const GuestUsername = "guest"

// initGuestUser 初始化游客用户（id=1, username=guest）
func initGuestUser() error {
	var user model.User
	err := DB.First(&user, GuestUserID).Error
	if err == nil {
		// 游客用户已存在
		logger.Debug("Guest user already exists", zap.Int64("id", GuestUserID))
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建游客用户
	// 密码哈希为 "guest123" 使用 bcrypt 生成
	guestUser := model.User{
		ID:           GuestUserID,
		Username:     GuestUsername,
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
	}

	if err := DB.Create(&guestUser).Error; err != nil {
		return fmt.Errorf("failed to create guest user: %w", err)
	}

	logger.Info("Guest user initialized",
		zap.Int64("id", GuestUserID),
		zap.String("username", GuestUsername),
	)
	return nil
}
