package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/internal/middleware"
)

func (r *Router) setupPrivateRoutes(api *gin.Engine) {

	protected := api.Group("/api")
	protected.Use(middleware.APIKeyMiddleware(r.Config.Security.APIKey))
	protected.Use(middleware.JWTAuth(r.Config.JWT.Secret))
	{
		user := protected.Group("/user")
		{
			user.GET("/profile", r.UserHandler.GetProfile)
		}

		limit := protected.Group("/limit")
		{
			limit.GET("/", middleware.PermissionMiddleware(r.UserRepo, r.PermCache, "get-limit"), r.LimitHandler.GetLimits)
			limit.POST("/", middleware.PermissionMiddleware(r.UserRepo, r.PermCache, "create-limit"), r.LimitHandler.CreateLimit)
			limit.PUT("/:id", middleware.PermissionMiddleware(r.UserRepo, r.PermCache, "edit-limit"), r.LimitHandler.UpdateLimit)
			limit.DELETE("/:id", middleware.PermissionMiddleware(r.UserRepo, r.PermCache, "delete-limit"), r.LimitHandler.DeleteLimit)
		}

		transaction := protected.Group("/transaction")
		{
			transaction.POST("/", middleware.PermissionMiddleware(r.UserRepo, r.PermCache, "create-transaction"), r.TransactionHandler.CreateTransaction)
			transaction.GET("/", middleware.PermissionMiddleware(r.UserRepo, r.PermCache, "get-transactions"), r.TransactionHandler.GetTransactions)
		}

		logs := protected.Group("/logs")
		{
			logs.GET("/audit", middleware.PermissionMiddleware(r.UserRepo, r.PermCache, "get-audit-log"), r.LogHandler.GetAuditLog)
			logs.GET("/auth", middleware.PermissionMiddleware(r.UserRepo, r.PermCache, "get-auth-log"), r.LogHandler.GetAuthLog)
		}
	}
}
