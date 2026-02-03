package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"what-to-eat/internal/model"
	"what-to-eat/internal/service"
)

type RestaurantHandler struct {
	restaurantService *service.RestaurantService
}

func NewRestaurantHandler(restaurantService *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{restaurantService: restaurantService}
}

// List 获取餐厅列表
// @Summary 获取餐厅列表
// @Tags 餐厅
// @Security Bearer
// @Produce json
// @Success 200 {object} model.Response{data=[]model.Restaurant}
// @Router /api/restaurants [get]
func (h *RestaurantHandler) List(c *gin.Context) {
	restaurants, err := h.restaurantService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "获取餐厅列表失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(restaurants))
}

// Create 创建餐厅
// @Summary 创建餐厅
// @Tags 餐厅
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body model.CreateRestaurantRequest true "餐厅信息"
// @Success 200 {object} model.Response{data=model.Restaurant}
// @Router /api/restaurants [post]
func (h *RestaurantHandler) Create(c *gin.Context) {
	var req model.CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "参数错误: "+err.Error()))
		return
	}

	restaurant, err := h.restaurantService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "创建餐厅失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(restaurant))
}

// Delete 删除餐厅
// @Summary 删除餐厅
// @Tags 餐厅
// @Security Bearer
// @Param id path int true "餐厅ID"
// @Success 200 {object} model.Response
// @Router /api/restaurants/{id} [delete]
func (h *RestaurantHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "无效的餐厅ID"))
		return
	}

	if err := h.restaurantService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "删除餐厅失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}
