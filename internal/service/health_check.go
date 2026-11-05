package service

import (
	"github.com/ngPhDuy/bookmark-management/internal/model/response"
)

const (
	OK           = "OK"
	SERVICE_NAME = "bookmark-management"
	INSTANCE_ID  = ""
)

type HealthCheckService interface {
	HealthCheck() (response.HealthCheckResponse, error) //* hoặc không *
}

type healthCheckService struct {
	generator func() string
}

func NewHealthCheckService(generator func() string) HealthCheckService {
	return &healthCheckService{
		generator: generator,
	}
}

func (s *healthCheckService) HealthCheck() (response.HealthCheckResponse, error) {
	instanceId := INSTANCE_ID

	if instanceId == "" {
		instanceId = s.generator()
	}

	return response.HealthCheckResponse{
		Message:     OK,
		ServiceName: SERVICE_NAME,
		InstanceID:  instanceId,
	}, nil
}
