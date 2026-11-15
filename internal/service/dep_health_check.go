package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// ErrDepUnhealthy là error duy nhất service này trả ra ngoài: nội dung của nó
// an toàn để lộ cho client, err gốc chỉ được log.
var ErrDepUnhealthy = errors.New("dependency unhealthy")

type DepHealthCheckService interface {
	HealthCheck(ctx context.Context) error
}

type depHealthCheckService struct {
	name string
	ping func(ctx context.Context) error
}

func NewDepHealthCheckService(name string, ping func(ctx context.Context) error) DepHealthCheckService {
	return &depHealthCheckService{
		name: name,
		ping: ping,
	}
}

func (s *depHealthCheckService) HealthCheck(ctx context.Context) error {
	if err := s.ping(ctx); err != nil {
		slog.Error("dependency health check failed", "dep", s.name, "err", err)
		return fmt.Errorf("%w: %s", ErrDepUnhealthy, s.name)
	}
	return nil
}
