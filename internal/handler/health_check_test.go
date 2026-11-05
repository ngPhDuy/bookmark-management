package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/model/response"
	service_mock "github.com/ngPhDuy/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name         string
		setUpMockSvc func() service_mock.MockHealthService

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Happy case",
			setUpMockSvc: func() service_mock.MockHealthService {
				return service_mock.NewMockHealthService(response.HealthCheckResponse{
					Message:     "OK",
					ServiceName: "bookmark-management",
					InstanceID:  "mock-instance-id",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"OK","service_name":"bookmark-management","instance_id":"mock-instance-id"}`,
		},
		{
			name: "Error case",
			setUpMockSvc: func() service_mock.MockHealthService {
				return service_mock.NewMockHealthService(response.HealthCheckResponse{}, errors.New("internal server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"internal server error"}`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()

			mockSvc := tc.setUpMockSvc()
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
