package model

import "context"

type Repository interface {
	// Set value with key if key does not exist, return key if success
	SetIfNotExist(ctx context.Context, key string, value string, exp int) (string, error)

	// Get value by key, return ErrKeyNotFound if key does not exist
	Get(ctx context.Context, key string) (string, error)
}
