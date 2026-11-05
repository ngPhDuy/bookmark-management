package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/config"
	"github.com/ngPhDuy/bookmark-management/internal/handler"
	"github.com/ngPhDuy/bookmark-management/internal/service"
	"github.com/ngPhDuy/bookmark-management/internal/utils"
)

type Engine interface {
	Start() error
}

type engine struct {
	app *gin.Engine
	cfg *config.AppConfig
}

func NewEngine(cfg *config.AppConfig) Engine {
	engine := engine{
		app: gin.Default(),
		cfg: cfg,
	}
	engine.initRoutes()

	return &engine
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.AppPort))
}

func (e *engine) initRoutes() {
	healthCheckService := service.NewHealthCheckService(utils.UUIDGenerator)
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheckService)

	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
}
