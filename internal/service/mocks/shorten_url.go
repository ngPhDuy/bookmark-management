package service_mock

import "context"

type MockShortenURLService struct {
	key string
	err error
}

func NewMockShortenURLService(key string, err error) *MockShortenURLService {
	return &MockShortenURLService{
		key: key,
		err: err,
	}
}

func (m *MockShortenURLService) ShortenURL(ctx context.Context, url string, exp int) (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.key, nil
}
