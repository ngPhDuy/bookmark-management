package redis

import (
	"context"
	"time"

	"fmt"

	"github.com/ngPhDuy/bookmark-management/internal/utils"
	goredis "github.com/redis/go-redis/v9"
)

const (
	MAX_ATTEMPTS   = 5
	MAX_KEY_LENGTH = 7
)

type RedisRepository interface {
	Set(ctx context.Context, value string, exp int) (string, error)
	Get(ctx context.Context, key string) (string, error)
}

type redisRepository struct {
	rdbClient          *goredis.Client
	randomKeyGenerator func(int) (string, error)
}

func NewRedisRepository(rdbClient *goredis.Client) RedisRepository {
	return &redisRepository{
		rdbClient:          rdbClient,
		randomKeyGenerator: utils.RandomKey,
	}
}

// Set value for a random key with expiration is exp (seconds), retry 5 times to generate a unique key
func (r *redisRepository) Set(ctx context.Context, value string, exp int) (string, error) {
	for i := 0; i < MAX_ATTEMPTS; i++ {
		key, err := r.randomKeyGenerator(MAX_KEY_LENGTH)

		if err != nil {
			return "", err
		}

		ok, err := r.rdbClient.SetNX(ctx, key, value, time.Duration(exp)*time.Second).Result()

		if err != nil {
			return "", err
		}

		if ok {
			return key, nil
		}
	}

	return "", fmt.Errorf("redis: failed to set key after %d attempts", MAX_ATTEMPTS)
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	value, err := r.rdbClient.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return value, nil
}
