package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/service"
)

type DepHealthCheckHandler interface {
	HealthCheck(c *gin.Context)
}

type depHealthCheckHandler struct {
	svc service.DepHealthCheckService
}

func NewDepHealthCheckHandler(healthCheckService service.DepHealthCheckService) DepHealthCheckHandler {
	return &depHealthCheckHandler{
		svc: healthCheckService,
	}
}

// HealthCheck handles dependency health check requests.
//
// @Summary      Dependency health check
// @Description  Ping an external dependency and report whether it is reachable.
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} map[string]string  "Dependency is unreachable"
// @Router       /health-check/redis [get]
func (h *depHealthCheckHandler) HealthCheck(c *gin.Context) {
	if err := h.svc.HealthCheck(c); err != nil {
		// Chỉ error của service mới an toàn để trả ra ngoài; error lạ có thể
		// chứa host/credential nên phải thay bằng message chung.
		msg := "health check failed"
		if errors.Is(err, service.ErrDepUnhealthy) {
			msg = err.Error()
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
