package service

import (
	"github.com/ngPhDuy/bookmark-management/internal/model/response"
)

const (
	OK = "OK"
)

type HealthCheckService interface {
	HealthCheck() (response.HealthCheckResponse, error) //* hoặc không *
}

type healthCheckService struct {
	generator   func() string
	serviceName string
	instanceID  string
}

func NewHealthCheckService(generator func() string, serviceName, instanceID string) HealthCheckService {
	if instanceID == "" {
		instanceID = generator()
	}

	return &healthCheckService{
		generator:   generator,
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

func (s *healthCheckService) HealthCheck() (response.HealthCheckResponse, error) {

	return response.HealthCheckResponse{
		Message:     OK,
		ServiceName: s.serviceName,
		InstanceID:  s.instanceID,
	}, nil
}
