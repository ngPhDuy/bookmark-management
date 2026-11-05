package service

import (
	"testing"

	"github.com/ngPhDuy/bookmark-management/internal/model/response"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {

	testCase := []struct {
		name            string
		expectedMessage response.HealthCheckResponse
		expectedError   error
	}{
		{
			name: "Happy case",
			expectedMessage: response.HealthCheckResponse{
				Message:     "OK",
				ServiceName: "bookmark-management",
				InstanceID:  "mock-instance-id",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			healthCheckService := NewHealthCheckService(func() string {
				return "mock-instance-id"
			})

			response, err := healthCheckService.HealthCheck()

			assert.Equal(t, tc.expectedMessage, response)
			assert.ErrorIs(t, tc.expectedError, err)
		})
	}
}
