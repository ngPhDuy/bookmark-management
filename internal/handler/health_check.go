package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/service"
)

type HealthCheckHandler interface {
	HealthCheck(c *gin.Context)
}

type healthCheckHandler struct {
	healthCheckService service.HealthCheckService
}

func NewHealthCheckHandler(healthCheckService service.HealthCheckService) HealthCheckHandler {
	return &healthCheckHandler{
		healthCheckService: healthCheckService,
	}
}

func (h *healthCheckHandler) HealthCheck(c *gin.Context) {
	response, err := h.healthCheckService.HealthCheck()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
