package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/service"
)

type RedisHealthCheckHandler interface {
	HealthCheck(c *gin.Context)
}

type redisHealthCheckHandler struct {
	svc service.RedisHealthCheckService
}

func NewRedisHealthCheckHandler(healthCheckService service.RedisHealthCheckService) RedisHealthCheckHandler {
	return &redisHealthCheckHandler{
		svc: healthCheckService,
	}
}

// HealthCheck handles health check requests.
//
// @Summary      Redis Health check
// @Description  Check whether the redis service is running.
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /health-check/redis [get]
func (h *redisHealthCheckHandler) HealthCheck(c *gin.Context) {
	if err := h.svc.HealthCheck(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
