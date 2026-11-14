package service

import (
	"context"
	"testing"

	"github.com/ngPhDuy/bookmark-management/internal/domain"
	mock "github.com/ngPhDuy/bookmark-management/internal/repository/redis/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		// Setup
		store      map[string]string
		repoErr    error
		genKeyFunc func(int) (string, error)
		exp        int

		// Expected
		expectedKey     string
		expectedErr     error
		expectedCalls   int
		expectedLastExp int
	}{
		{
			name: "happy case 1",
			store: map[string]string{
				"key1": "val1", "key2": "val2", "key3": "val3",
			},
			repoErr: nil,
			genKeyFunc: func(int) (string, error) {
				return "mock-key", nil
			},
			exp:             60,
			expectedKey:     "mock-key",
			expectedErr:     nil,
			expectedCalls:   1,
			expectedLastExp: 60,
		},
		{
			name: "retry 1 time",
			store: map[string]string{
				"key1": "val1",
			},
			repoErr:         nil,
			genKeyFunc:      sequenceKeys("key1", "key2"),
			exp:             10,
			expectedKey:     "key2",
			expectedErr:     nil,
			expectedCalls:   2,
			expectedLastExp: 10,
		},
		{
			name: "retry max times, fail",
			store: map[string]string{
				"key1": "val1",
			},
			repoErr:         nil,
			genKeyFunc:      sequenceKeys("key1"),
			exp:             5,
			expectedKey:     "",
			expectedErr:     domain.ErrShortenUrlFail,
			expectedCalls:   MAX_ATTEMPTS,
			expectedLastExp: 0,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			rep := mock.NewMockRedisRepository()
			rep.SetErr = tc.repoErr
			for k, v := range tc.store {
				rep.Store[k] = v
			}

			// Test
			svc := NewShortenURLService(rep, tc.genKeyFunc)
			key, err := svc.ShortenURL(context.Background(), "https://www.google.com", tc.exp)

			// Assert
			assert.Equal(t, tc.expectedKey, key)
			assert.ErrorIs(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedCalls, rep.CountSetCall)
			assert.Equal(t, tc.expectedLastExp, rep.LastExp)
		})
	}
}

// mock genKeyFunc
func sequenceKeys(keys ...string) func(int) (string, error) {
	i := 0
	return func(int) (string, error) {
		k := keys[min(i, len(keys)-1)]
		i++
		return k, nil
	}
}
