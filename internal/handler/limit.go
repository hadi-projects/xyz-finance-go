package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	services "github.com/hadi-projects/xyz-finance-go/internal/service"
)

type LimitHandler struct {
	limitService services.LimitService
}

// NewLimitHandler creates a new limit handler instance
func NewLimitHandler(limitService services.LimitService) *LimitHandler {
	return &LimitHandler{
		limitService: limitService,
	}
}

func (h *LimitHandler) GetLimits(c *gin.Context) {
	// get user id from context
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limits, err := h.limitService.GetLimits(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get User Limits",
		"data":    limits,
	})
}

func (h *LimitHandler) CreateLimit(c *gin.Context) {
	var req dto.CreateLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.limitService.CreateLimit(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Limit created successfully",
		"data": dto.LimitResponse{
			UserID:      req.TargetUserID,
			TenorMonth:  req.TenorMonth,
			LimitAmount: req.LimitAmount,
		},
	})
}

func (h *LimitHandler) UpdateLimit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	var req dto.UpdateLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.limitService.UpdateLimit(uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Limit updated successfully"})
}

func (h *LimitHandler) DeleteLimit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit ID"})
		return
	}

	if err := h.limitService.DeleteLimit(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Limit deleted successfully"})
}
