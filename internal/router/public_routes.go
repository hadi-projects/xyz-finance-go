package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) setupPublicRoutes(router *gin.Engine) {

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"app":     "XYZ Multifinance",
			"version": "1.0.0",
		})
	})

}
