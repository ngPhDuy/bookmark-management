package redis

import (
	"context"
	"time"

	"github.com/ngPhDuy/bookmark-management/internal/utils"
	goredis "github.com/redis/go-redis/v9"
	"fmt"
)

const (
	MAX_ATTEMPTS   = 5
	MAX_KEY_LENGTH = 7
)

type RedisRepository interface {
	Set(key string, value string, exp int) error
	Get(key string) (string, error)
}

type redisRepository struct {
	ctx       context.Context
	rdbClient *goredis.Client
}

func NewRedisRepository(ctx context.Context, rdbClient *goredis.Client) RedisRepository {
	return &redisRepository{
		ctx:       ctx,
		rdbClient: rdbClient,
	}
}

// Set key and value with expiration is exp (seconds), retry 5 times to generate a unique key
func (r *redisRepository) Set(key string, value string, exp int) error {
	for i := 0; i < MAX_ATTEMPTS; i++ {
		key, err := utils.RandomKey(MAX_KEY_LENGTH)

		if err != nil {
			return err
		}

		ok, err := r.rdbClient.SetNX(r.ctx, key, value, time.Duration(exp)*time.Second).Result()

		if err != nil {
			return err
		}

		if ok {
			return nil
		}
	}

	return fmt.Errorf("redis: failed to set key after %d attempts", MAX_ATTEMPTS)
}

func (r *redisRepository) Get(key string) (string, error) {
	value, err := r.rdbClient.Get(r.ctx, key).Result()

	if err != nil {
		return "", err
	}

	return value, nil
}
