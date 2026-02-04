package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"what-to-eat/internal/model"
	"what-to-eat/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "注册信息"
// @Success 200 {object} model.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "参数错误: "+err.Error()))
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			c.JSON(http.StatusConflict, model.Error(409, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Error(500, "注册失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"id":       user.ID,
		"username": user.Username,
	}))
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "登录信息"
// @Success 200 {object} model.Response{data=model.LoginResponse}
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "参数错误: "+err.Error()))
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, model.Error(401, "用户名或密码错误"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Error(500, "登录失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

// GuestLogin 游客登录
// @Summary 游客登录
// @Tags 认证
// @Produce json
// @Success 200 {object} model.Response{data=model.LoginResponse}
// @Router /api/auth/guest [post]
func (h *AuthHandler) GuestLogin(c *gin.Context) {
	resp, err := h.authService.GuestLogin()
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusInternalServerError, model.Error(500, "游客账号未初始化"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Error(500, "游客登录失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}
