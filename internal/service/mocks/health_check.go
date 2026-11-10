package service_mock

import (
	"github.com/ngPhDuy/bookmark-management/internal/model/response"
)

type MockHealthService interface {
	HealthCheck() (response.HealthCheckResponse, error)
}

type mockHealthService struct {
	response response.HealthCheckResponse
	err      error
}

func NewMockHealthService(response response.HealthCheckResponse, err error) MockHealthService {
	return &mockHealthService{
		response: response,
		err:      err,
	}
}

func (m *mockHealthService) HealthCheck() (response.HealthCheckResponse, error) {
	return m.response, m.err
}
