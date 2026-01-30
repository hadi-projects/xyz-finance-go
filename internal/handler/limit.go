package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *LimitHandler) GetConsumerLimit(c *gin.Context) {
	limit, err := h.limitService.GetConsumerLimit(c.Param("limit"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Get Consumer Limit",
		"data":    limit,
	})
}
