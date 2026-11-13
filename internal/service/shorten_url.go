package service

import (
	"context"

	"github.com/ngPhDuy/bookmark-management/internal/domain"
	"github.com/ngPhDuy/bookmark-management/internal/model"
)

const (
	MAX_ATTEMPTS   = 5
	MAX_KEY_LENGTH = 7
)

type ShortenURLService interface {
	//Gen a random key and save it [key, url] to repository
	//Return key if success, domain.ErrShortenUrlFail if fail
	ShortenURL(ctx context.Context, url string, exp int) (string, error)
}

type shortenURLService struct {
	rep        model.Repository
	genKeyFunc func(int) (string, error)
}

func NewShortenURLService(rep model.Repository, genKeyFunc func(int) (string, error)) ShortenURLService {
	return &shortenURLService{
		rep:        rep,
		genKeyFunc: genKeyFunc,
	}
}

func (s *shortenURLService) ShortenURL(ctx context.Context, url string, exp int) (string, error) {
	for i := 0; i < MAX_ATTEMPTS; i++ {
		key, err := s.genKeyFunc(MAX_KEY_LENGTH)

		if err != nil {
			return "", err
		}

		_, err = s.rep.SetIfNotExist(ctx, key, url, exp)

		if err == domain.ErrKeyExists {
			continue
		}

		if err != nil {
			return "", err
		}

		return key, nil
	}

	return "", domain.ErrShortenUrlFail
}
