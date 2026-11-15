package utils

import "github.com/google/uuid"

func UUIDGenerator() string {
	return uuid.New().String()
}
