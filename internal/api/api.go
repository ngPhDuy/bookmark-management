package api

import (
	"fmt"

	_ "github.com/ngPhDuy/bookmark-management/cmd/api/docs"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/config"
	"github.com/ngPhDuy/bookmark-management/internal/handler"
	"github.com/ngPhDuy/bookmark-management/internal/repository/redis"
	"github.com/ngPhDuy/bookmark-management/internal/service"
	"github.com/ngPhDuy/bookmark-management/internal/utils"
	goredis "github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start() error
}

type engine struct {
	app         *gin.Engine
	cfg         *config.AppConfig
	redisClient *goredis.Client
}

func NewEngine(cfg *config.AppConfig, redisClient *goredis.Client) Engine {
	engine := engine{
		app:         gin.Default(),
		cfg:         cfg,
		redisClient: redisClient,
	}
	engine.initRoutes()

	return &engine
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.AppPort))
}

func (e *engine) initRoutes() {
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	redisRep := redis.NewRedisRepository(e.redisClient)

	healthCheckService := service.NewHealthCheckService(utils.UUIDGenerator, e.cfg.ServiceName, e.cfg.InstanceID)
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheckService)

	shortenURLService := service.NewShortenURLService(redisRep, utils.RandomKey)
	shortenURLHandler := handler.NewShortenURLHandler(shortenURLService)

	redisHealthCheckService := service.NewRedisHealthCheckService(e.redisClient)
	redisHealthCheckHandler := handler.NewRedisHealthCheckHandler(redisHealthCheckService)

	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
	e.app.POST("/v1/links/shorten", shortenURLHandler.ShortenURL)
	e.app.GET("/health-check/redis", redisHealthCheckHandler.HealthCheck)
}
