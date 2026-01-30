package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/config"
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

	r.setupPublicRoutes(router)

	return router
}
