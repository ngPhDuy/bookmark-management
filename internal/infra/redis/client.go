package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

// New a redis client and check connection. Caller need to close the client
func NewRedisClient(ctx context.Context, config *RedisConfig) (*goredis.Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
