package repository

import (
	"what-to-eat/internal/model"

	"gorm.io/gorm"
)

type RestaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

// Create 创建餐厅
func (r *RestaurantRepository) Create(restaurant *model.Restaurant) error {
	return r.db.Create(restaurant).Error
}

// GetAll 获取所有餐厅
func (r *RestaurantRepository) GetAll() ([]model.Restaurant, error) {
	var restaurants []model.Restaurant
	err := r.db.Order("id ASC").Find(&restaurants).Error
	return restaurants, err
}

// GetByID 根据ID查询餐厅
func (r *RestaurantRepository) GetByID(id int64) (*model.Restaurant, error) {
	var restaurant model.Restaurant
	err := r.db.First(&restaurant, id).Error
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// GetByIDs 根据ID列表查询餐厅
func (r *RestaurantRepository) GetByIDs(ids []int64) ([]model.Restaurant, error) {
	var restaurants []model.Restaurant
	err := r.db.Where("id IN ?", ids).Find(&restaurants).Error
	return restaurants, err
}

// GetByName 根据名称查询餐厅
func (r *RestaurantRepository) GetByName(name string) (*model.Restaurant, error) {
	var restaurant model.Restaurant
	err := r.db.Where("name = ?", name).First(&restaurant).Error
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// GetOrCreate 获取或创建餐厅
func (r *RestaurantRepository) GetOrCreate(name string) (*model.Restaurant, bool, error) {
	var restaurant model.Restaurant
	err := r.db.Where("name = ?", name).First(&restaurant).Error
	if err == nil {
		return &restaurant, false, nil // 已存在
	}
	if err != gorm.ErrRecordNotFound {
		return nil, false, err
	}

	// 创建新餐厅
	restaurant = model.Restaurant{Name: name}
	if err := r.db.Create(&restaurant).Error; err != nil {
		return nil, false, err
	}
	return &restaurant, true, nil // 新创建
}

// Delete 删除餐厅
func (r *RestaurantRepository) Delete(id int64) error {
	return r.db.Delete(&model.Restaurant{}, id).Error
}

// Count 统计餐厅数量
func (r *RestaurantRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Restaurant{}).Count(&count).Error
	return count, err
}
