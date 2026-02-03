package service

import (
	"what-to-eat/internal/model"
	"what-to-eat/internal/repository"
)

type RestaurantService struct {
	restaurantRepo *repository.RestaurantRepository
}

func NewRestaurantService(restaurantRepo *repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{
		restaurantRepo: restaurantRepo,
	}
}

// Create 创建餐厅
func (s *RestaurantService) Create(req *model.CreateRestaurantRequest) (*model.Restaurant, error) {
	restaurant := &model.Restaurant{
		Name: req.Name,
	}

	if err := s.restaurantRepo.Create(restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

// GetAll 获取所有餐厅
func (s *RestaurantService) GetAll() ([]model.Restaurant, error) {
	return s.restaurantRepo.GetAll()
}

// GetByID 根据ID获取餐厅
func (s *RestaurantService) GetByID(id int64) (*model.Restaurant, error) {
	return s.restaurantRepo.GetByID(id)
}

// Delete 删除餐厅
func (s *RestaurantService) Delete(id int64) error {
	return s.restaurantRepo.Delete(id)
}
