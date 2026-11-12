package model

import "context"

type Repository interface {
	SetIfNotExist(ctx context.Context, key string, value string, exp int) (string, error)
	Get(ctx context.Context, key string) (string, error)
}
