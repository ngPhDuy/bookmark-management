package service

import (
	"testing"

	"github.com/ngPhDuy/bookmark-management/internal/model/response"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {

	t.Parallel()

	testCase := []struct {
		name            string
		generator       func() string
		serviceName     string
		instanceID      string
		expectedMessage response.HealthCheckResponse
		expectedError   error
	}{
		{
			name: "Happy case 1",
			generator: func() string {
				return "mock-instance-id"
			},
			serviceName: "bookmark-management",
			instanceID:  "",
			expectedMessage: response.HealthCheckResponse{
				Message:     "OK",
				ServiceName: "bookmark-management",
				InstanceID:  "mock-instance-id",
			},
			expectedError: nil,
		},
		{
			name: "Happy case 2",
			generator: func() string {
				return "mock-instance-id"
			},
			serviceName: "bookmark-management",
			instanceID:  "real-instance-id",
			expectedMessage: response.HealthCheckResponse{
				Message:     "OK",
				ServiceName: "bookmark-management",
				InstanceID:  "real-instance-id",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			healthCheckService := NewHealthCheckService(tc.generator, tc.serviceName, tc.instanceID)

			response, err := healthCheckService.HealthCheck()

			assert.Equal(t, tc.expectedMessage, response)
			assert.ErrorIs(t, tc.expectedError, err)
		})
	}
}
