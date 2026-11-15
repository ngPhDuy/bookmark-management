package service

import (
	"context"
	"testing"

	"github.com/ngPhDuy/bookmark-management/internal/domain"
	"github.com/ngPhDuy/bookmark-management/internal/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testURL = "https://www.google.com"

func TestShortenURL(t *testing.T) {
	t.Parallel()

	testCase := []struct {
		name string

		// Setup
		setUpRepo  func(rep *mocks.MockRepository)
		genKeyFunc func(int) (string, error)
		exp        int

		// Expected
		expectedKey string
		expectedErr error
	}{
		{
			name: "happy case 1",
			setUpRepo: func(rep *mocks.MockRepository) {
				rep.EXPECT().SetIfNotExist(mock.Anything, "mock-key", testURL, 60).
					Return("mock-key", nil).Once()
			},
			genKeyFunc: func(int) (string, error) {
				return "mock-key", nil
			},
			exp:         60,
			expectedKey: "mock-key",
			expectedErr: nil,
		},
		{
			name: "retry 1 time",
			setUpRepo: func(rep *mocks.MockRepository) {
				rep.EXPECT().SetIfNotExist(mock.Anything, "key1", testURL, 10).
					Return("", domain.ErrKeyExists).Once()
				rep.EXPECT().SetIfNotExist(mock.Anything, "key2", testURL, 10).
					Return("key2", nil).Once()
			},
			genKeyFunc:  sequenceKeys("key1", "key2"),
			exp:         10,
			expectedKey: "key2",
			expectedErr: nil,
		},
		{
			name: "retry max times, fail",
			setUpRepo: func(rep *mocks.MockRepository) {
				rep.EXPECT().SetIfNotExist(mock.Anything, "key1", testURL, 5).
					Return("", domain.ErrKeyExists).Times(MAX_ATTEMPTS)
			},
			genKeyFunc:  sequenceKeys("key1"),
			exp:         5,
			expectedKey: "",
			expectedErr: domain.ErrShortenUrlFail,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			rep := mocks.NewMockRepository(t)
			tc.setUpRepo(rep)

			// Test
			svc := NewShortenURLService(rep, tc.genKeyFunc)
			key, err := svc.ShortenURL(context.Background(), testURL, tc.exp)

			// Assert
			assert.Equal(t, tc.expectedKey, key)
			assert.ErrorIs(t, err, tc.expectedErr)
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
