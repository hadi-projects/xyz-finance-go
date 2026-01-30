package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/config"
	"github.com/hadi-projects/xyz-finance-go/internal/delivery/http/middleware"
)

type Router struct {
	Config *config.Config
}

func NewRouter(
	cfg *config.Config,
) *Router {
	return &Router{
		Config: cfg,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {

	if r.Config.AppPort == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.CORS(r.Config.Security.CORSAllowedOrigins, r.Config.Security.CORSAllowCredentials))
	router.Use(middleware.RateLimiter(r.Config.Security.RateLimitRPS, r.Config.Security.RateLimitBurst))

	r.setupPublicRoutes(router)

	return router
}
