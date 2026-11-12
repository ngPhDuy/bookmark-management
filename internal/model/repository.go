package model

import "context"

type Repository interface {
	Set(ctx context.Context, value string, exp int) (string, error)
	Get(ctx context.Context, key string) (string, error)
}
