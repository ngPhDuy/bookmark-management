package domain

import "errors"

var (
	ErrKeyExists      = errors.New("key already exists")
	ErrKeyNotFound    = errors.New("key not found")
	ErrShortenUrlFail = errors.New("shorten url fail")
)
