package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/service"
	"github.com/ngPhDuy/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

			mockSvc := mocks.NewMockDepHealthCheckService(t)
			mockSvc.EXPECT().HealthCheck(mock.Anything).Return(tc.svcErr).Once()

			handler := NewDepHealthCheckHandler(mockSvc)

			router := gin.New()
			router.GET("/health-check/redis", handler.HealthCheck)
			req := httptest.NewRequest(http.MethodGet, "/health-check/redis", nil)

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
