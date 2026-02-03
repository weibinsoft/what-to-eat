package repository

import (
	"what-to-eat/internal/model"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

// GetAll 获取所有菜单（包含餐厅信息）
func (r *MenuRepository) GetAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Preload("Restaurant").Order("id ASC").Find(&menus).Error
	return menus, err
}

// GetByID 根据ID查询菜单
func (r *MenuRepository) GetByID(id int64) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Preload("Restaurant").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetByIDs 根据ID列表查询菜单
func (r *MenuRepository) GetByIDs(ids []int64) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Preload("Restaurant").Where("id IN ?", ids).Find(&menus).Error
	return menus, err
}

// GetByRestaurantID 根据餐厅ID查询菜单
func (r *MenuRepository) GetByRestaurantID(restaurantID int64) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Preload("Restaurant").Where("restaurant_id = ?", restaurantID).Find(&menus).Error
	return menus, err
}

// Delete 删除菜单
func (r *MenuRepository) Delete(id int64) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

// Count 统计菜单数量
func (r *MenuRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).Count(&count).Error
	return count, err
}

// ExistsByRestaurantAndDish 检查餐厅下是否已有该菜品
func (r *MenuRepository) ExistsByRestaurantAndDish(restaurantID int64, dishName string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).
		Where("restaurant_id = ? AND dish_name = ?", restaurantID, dishName).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
