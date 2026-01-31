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
		}

		limit := protected.Group("/limit")
		{
			limit.GET("/", middleware.PermissionMiddleware(r.UserRepo, "get-limit"), r.LimitHandler.GetLimits)
			limit.POST("/", middleware.PermissionMiddleware(r.UserRepo, "create-limit"), r.LimitHandler.CreateLimit)
			limit.DELETE("/:id", middleware.PermissionMiddleware(r.UserRepo, "delete-limit"), r.LimitHandler.DeleteLimit)
		}

		transaction := protected.Group("/transaction")
		{
			transaction.POST("/", middleware.PermissionMiddleware(r.UserRepo, "create-transaction"), r.TransactionHandler.CreateTransaction)
		}
	}
}
