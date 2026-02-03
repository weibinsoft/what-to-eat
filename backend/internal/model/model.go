package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"type:varchar(50);not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Restaurant 餐厅模型
type Restaurant struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Menus     []Menu    `json:"menus,omitempty" gorm:"foreignKey:RestaurantID;constraint:false"`
}

// Menu 菜单模型（菜品）- 支持软删除
type Menu struct {
	ID           int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	RestaurantID int64          `json:"restaurant_id" gorm:"not null;index:idx_restaurant"`
	DishName     string         `json:"dish_name" gorm:"type:varchar(100);not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"` // 软删除字段
	Restaurant   Restaurant     `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID;constraint:false"`
}

// DecisionRecord 决策记录模型
type DecisionRecord struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null;index:idx_user_decided"`
	MenuID    int64     `json:"menu_id" gorm:"not null"`
	DecidedAt time.Time `json:"decided_at" gorm:"index:idx_user_decided"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:false"`
	Menu      Menu      `json:"menu,omitempty" gorm:"foreignKey:MenuID;constraint:false"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (Restaurant) TableName() string {
	return "restaurants"
}

func (Menu) TableName() string {
	return "menus"
}

func (DecisionRecord) TableName() string {
	return "decision_records"
}
