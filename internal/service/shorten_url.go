package service

import "github.com/ngPhDuy/bookmark-management/internal/model"

type ShortenURLService interface {
	ShortenURL(url string) (string, error)
}

type shortenURLService struct {
	rep        model.Repository
	genKeyFunc func(int) (string, error)
}
