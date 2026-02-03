package service

import (
	"errors"

	"what-to-eat/internal/model"
	"what-to-eat/internal/repository"
)

var (
	ErrMenuExists = errors.New("该餐厅已有此菜品")
)

type MenuService struct {
	menuRepo       *repository.MenuRepository
	restaurantRepo *repository.RestaurantRepository
}

func NewMenuService(menuRepo *repository.MenuRepository, restaurantRepo *repository.RestaurantRepository) *MenuService {
	return &MenuService{
		menuRepo:       menuRepo,
		restaurantRepo: restaurantRepo,
	}
}

// Create 创建菜单（同时处理餐厅）
func (s *MenuService) Create(req *model.CreateMenuRequest) (*model.Menu, bool, error) {
	// 获取或创建餐厅
	restaurant, isNewRestaurant, err := s.restaurantRepo.GetOrCreate(req.RestaurantName)
	if err != nil {
		return nil, false, err
	}

	// 检查该餐厅是否已有此菜品
	exists, err := s.menuRepo.ExistsByRestaurantAndDish(restaurant.ID, req.DishName)
	if err != nil {
		return nil, false, err
	}
	if exists {
		return nil, false, ErrMenuExists
	}

	// 创建菜单
	menu := &model.Menu{
		RestaurantID: restaurant.ID,
		DishName:     req.DishName,
	}

	if err := s.menuRepo.Create(menu); err != nil {
		return nil, false, err
	}

	// 加载餐厅信息
	menu.Restaurant = *restaurant

	return menu, isNewRestaurant, nil
}

// GetAll 获取所有菜单
func (s *MenuService) GetAll() ([]model.Menu, error) {
	return s.menuRepo.GetAll()
}

// GetByID 根据ID获取菜单
func (s *MenuService) GetByID(id int64) (*model.Menu, error) {
	return s.menuRepo.GetByID(id)
}

// Delete 删除菜单
func (s *MenuService) Delete(id int64) error {
	return s.menuRepo.Delete(id)
}

// GetAllRestaurants 获取所有餐厅（用于下拉选择）
func (s *MenuService) GetAllRestaurants() ([]model.Restaurant, error) {
	return s.restaurantRepo.GetAll()
}
