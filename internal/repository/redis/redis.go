package redis

import (
	"context"
	"errors"
	"time"

	"github.com/ngPhDuy/bookmark-management/internal/domain"
	"github.com/ngPhDuy/bookmark-management/internal/model"
	goredis "github.com/redis/go-redis/v9"
)

type redisClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.BoolCmd
	Get(ctx context.Context, key string) *goredis.StringCmd
}

type redisRepository struct {
	rdbClient redisClient
}

func NewRedisRepository(rdbClient redisClient) model.Repository {
	return &redisRepository{
		rdbClient: rdbClient,
	}
}

func (r *redisRepository) SetIfNotExist(ctx context.Context, key string, value string, exp int) (string, error) {
	ok, err := r.rdbClient.SetNX(ctx, key, value, time.Duration(exp)*time.Second).Result()

	if err != nil {
		return "", err
	}

	if !ok {
		return "", domain.ErrKeyExists
	}

	return key, nil
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	value, err := r.rdbClient.Get(ctx, key).Result()

	//err is goredis.Nil -> not found
	if errors.Is(err, goredis.Nil) {
		return "", domain.ErrKeyNotFound
	}

	if err != nil {
		return "", err
	}

	return value, nil
}
