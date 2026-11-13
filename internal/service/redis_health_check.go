package service

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type RedisHealthCheckService interface {
	HealthCheck(ctx context.Context) error
}

type redisHealthCheckService struct {
	redisClient *goredis.Client
}

func NewRedisHealthCheckService(redisClient *goredis.Client) RedisHealthCheckService {
	return &redisHealthCheckService{
		redisClient: redisClient,
	}
}

// Check redis connection, return error if failed
func (s *redisHealthCheckService) HealthCheck(ctx context.Context) error {
	return s.redisClient.Ping(ctx).Err()
}
