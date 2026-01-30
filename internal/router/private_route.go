package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/internal/middleware"
)

func (r *Router) setupPrivateRoutes(api *gin.Engine) {

	protected := api.Group("/api")
	protected.Use(middleware.JWTAuth(r.Config.JWT.Secret))
	{
		user := protected.Group("/user")
		{
			user.GET("/profile", r.UserHandler.GetProfile)
			user.GET("/limit", middleware.PermissionMiddleware(r.UserRepo, "get-limit"), r.UserHandler.GetLimitsByUserID)
		}
	}
}
