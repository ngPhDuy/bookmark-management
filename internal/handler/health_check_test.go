package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/model/response"
	"github.com/ngPhDuy/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		svcResponse response.HealthCheckResponse
		svcErr      error

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Happy case",
			svcResponse: response.HealthCheckResponse{
				Message:     "OK",
				ServiceName: "bookmark-management",
				InstanceID:  "mock-instance-id",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"OK","service_name":"bookmark-management","instance_id":"mock-instance-id"}`,
		},
		{
			name:               "Error case",
			svcErr:             errors.New("internal server error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"internal server error"}`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()

			mockSvc := mocks.NewMockHealthCheckService(t)
			mockSvc.EXPECT().HealthCheck().Return(tc.svcResponse, tc.svcErr).Once()

			handler := NewHealthCheckHandler(mockSvc)

			router := gin.New()
			router.GET("/health-check", handler.HealthCheck)
			req := httptest.NewRequest(http.MethodGet, "/health-check", nil)

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
