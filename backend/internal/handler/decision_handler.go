package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"what-to-eat/internal/model"
	"what-to-eat/internal/service"
	"what-to-eat/pkg/logger"
	"what-to-eat/pkg/middleware"
)

type DecisionHandler struct {
	decisionService *service.DecisionService
}

func NewDecisionHandler(decisionService *service.DecisionService) *DecisionHandler {
	return &DecisionHandler{decisionService: decisionService}
}

// Decide 执行决策
// @Summary 执行随机决策
// @Tags 决策
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body model.DecideRequest true "决策请求"
// @Success 200 {object} model.Response{data=model.DecideResponse}
// @Router /api/decide [post]
func (h *DecisionHandler) Decide(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, model.Error(401, "用户未登录"))
		return
	}

	var req model.DecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 允许空body
		req = model.DecideRequest{}
	}

	resp, err := h.decisionService.Decide(userID, req.MenuIDs)
	if err != nil {
		logger.Error("Decision failed", zap.Int64("userID", userID), zap.Error(err))
		if errors.Is(err, service.ErrNoMenus) {
			c.JSON(http.StatusBadRequest, model.Error(400, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Error(500, "决策失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

// History 获取历史记录
// @Summary 获取最近5天的决策历史
// @Tags 决策
// @Security Bearer
// @Produce json
// @Success 200 {object} model.Response{data=model.HistoryResponse}
// @Router /api/history [get]
func (h *DecisionHandler) History(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, model.Error(401, "用户未登录"))
		return
	}

	resp, err := h.decisionService.GetHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "获取历史记录失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}
