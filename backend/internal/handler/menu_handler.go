package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"what-to-eat/internal/model"
	"what-to-eat/internal/service"
)

type MenuHandler struct {
	menuService *service.MenuService
}

func NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

// List 获取菜单列表
// @Summary 获取菜单列表
// @Tags 菜单
// @Security Bearer
// @Produce json
// @Success 200 {object} model.Response{data=[]model.Menu}
// @Router /api/menus [get]
func (h *MenuHandler) List(c *gin.Context) {
	menus, err := h.menuService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "获取菜单列表失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(menus))
}

// Create 创建菜单（同时处理餐厅）
// @Summary 创建菜单
// @Tags 菜单
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body model.CreateMenuRequest true "菜单信息"
// @Success 200 {object} model.Response{data=model.Menu}
// @Router /api/menus [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req model.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "参数错误: "+err.Error()))
		return
	}

	menu, isNewRestaurant, err := h.menuService.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrMenuExists) {
			c.JSON(http.StatusConflict, model.Error(409, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Error(500, "创建菜单失败"))
		return
	}

	// 返回结果，包含是否创建了新餐厅的信息
	c.JSON(http.StatusOK, model.Success(gin.H{
		"menu":              menu,
		"is_new_restaurant": isNewRestaurant,
	}))
}

// Delete 删除菜单
// @Summary 删除菜单
// @Tags 菜单
// @Security Bearer
// @Param id path int true "菜单ID"
// @Success 200 {object} model.Response
// @Router /api/menus/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "无效的菜单ID"))
		return
	}

	if err := h.menuService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "删除菜单失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

// ListRestaurants 获取餐厅列表（用于下拉选择）
// @Summary 获取餐厅列表
// @Tags 菜单
// @Security Bearer
// @Produce json
// @Success 200 {object} model.Response{data=[]model.Restaurant}
// @Router /api/restaurants [get]
func (h *MenuHandler) ListRestaurants(c *gin.Context) {
	restaurants, err := h.menuService.GetAllRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "获取餐厅列表失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(restaurants))
}
