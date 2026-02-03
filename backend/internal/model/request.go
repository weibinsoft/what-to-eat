package model

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateRestaurantRequest 创建餐厅请求
type CreateRestaurantRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// CreateMenuRequest 创建菜单请求（同时包含餐厅信息）
type CreateMenuRequest struct {
	RestaurantID   int64  `json:"restaurant_id"`                                    // 餐厅ID（可选，如果提供则使用现有餐厅）
	RestaurantName string `json:"restaurant_name" binding:"required,min=1,max=100"` // 餐厅名称（必填，用于查找或创建）
	DishName       string `json:"dish_name" binding:"required,min=1,max=100"`       // 菜品名称（必填）
}

// DecideRequest 决策请求
type DecideRequest struct {
	// 可选：指定参与决策的菜单ID列表，为空则使用全部菜单
	MenuIDs []int64 `json:"menu_ids"`
}
