package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/config"
	"github.com/hadi-projects/xyz-finance-go/internal/handler"
	"github.com/hadi-projects/xyz-finance-go/internal/middleware"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
)

type Router struct {
	Config             *config.AppConfig
	AuthHandler        *handler.AuthHandler
	LimitHandler       *handler.LimitHandler
	UserHandler        *handler.UserHandler
	TransactionHandler *handler.TransactionHandler
	LogHandler         *handler.LogHandler
	UserRepo           repository.UserRepository
}

func NewRouter(
	cfg *config.AppConfig,
	authHandler *handler.AuthHandler,
	limitHandler *handler.LimitHandler,
	userHandler *handler.UserHandler,
	transactionHandler *handler.TransactionHandler,
	logHandler *handler.LogHandler,
	userRepo repository.UserRepository,
) *Router {
	return &Router{
		Config:             cfg,
		AuthHandler:        authHandler,
		LimitHandler:       limitHandler,
		UserHandler:        userHandler,
		TransactionHandler: transactionHandler,
		LogHandler:         logHandler,
		UserRepo:           userRepo,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {

	if r.Config.AppPort == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS(r.Config.Security.CORSAllowedOrigins, r.Config.Security.CORSAllowCredentials))
	router.Use(middleware.RateLimiter(r.Config.Security.RateLimitRPS, r.Config.Security.RateLimitBurst))
	router.Use(middleware.RequestCancellation(time.Duration(r.Config.Security.RequestTimeout) * time.Second))
	router.Use(middleware.XSSProtection())
	router.Use(middleware.SecureHeaders())

	r.setupPublicRoutes(router)
	r.setupPrivateRoutes(router)

	return router
}
