package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	services "github.com/hadi-projects/xyz-finance-go/internal/service"
)

type LogHandler struct {
	logService services.LogService
}

func NewLogHandler(logService services.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

func (h *LogHandler) GetAuditLog(c *gin.Context) {
	logs, err := h.logService.GetAuditLog()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func (h *LogHandler) GetAuthLog(c *gin.Context) {
	logs, err := h.logService.GetAuthLog()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs})
}
