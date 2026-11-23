package service

import (
	"context"
	"errors"

	"github.com/ngPhDuy/bookmark-management/internal/domain"
	"github.com/ngPhDuy/bookmark-management/internal/model"
)

type RedirectURLService interface {
	// Get <key, url> within redis. Return url if it present, else return err
	RedirectURL(ctx context.Context, key string) (string, error)
}

type redirectURLService struct {
	rep model.Repository
}

func NewRedirectURLService(rep model.Repository) RedirectURLService {
	return &redirectURLService{
		rep: rep,
	}
}

func (s *redirectURLService) RedirectURL(ctx context.Context, key string) (string, error) {
	url, err := s.rep.Get(ctx, key)

	if errors.Is(err, domain.ErrKeyNotFound) {
		//log "get url with <key> not found"
		return "", err
	}

	if err != nil {
		//log error
		return "", domain.ErrInternalErr
	}

	return url, nil
}
