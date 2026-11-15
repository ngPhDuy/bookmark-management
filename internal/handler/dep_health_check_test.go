package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/service"
	"github.com/stretchr/testify/assert"
)

type stubDepHealthCheckService struct {
	err error
}

func (s stubDepHealthCheckService) HealthCheck(ctx context.Context) error { return s.err }

func TestDepHealthCheckHandler(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name   string
		svcErr error

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "Happy case",
			svcErr:             nil,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"OK"}`,
		},
		{
			name:               "Dependency unhealthy",
			svcErr:             service.ErrDepUnhealthy,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"dependency unhealthy"}`,
		},
		{
			name:               "Unknown error is masked",
			svcErr:             errors.New("dial tcp 10.0.0.1:6379: AUTH failed password=hunter2"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"health check failed"}`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()

			handler := NewDepHealthCheckHandler(stubDepHealthCheckService{err: tc.svcErr})

			router := gin.New()
			router.GET("/health-check/redis", handler.HealthCheck)
			req := httptest.NewRequest(http.MethodGet, "/health-check/redis", nil)

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
