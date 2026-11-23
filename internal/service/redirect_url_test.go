package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ngPhDuy/bookmark-management/internal/domain"
	"github.com/ngPhDuy/bookmark-management/internal/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testKey = "test-key"

func TestRedirectURL(t *testing.T) {
	t.Parallel()

	testcase := []struct {
		name      string
		setUpRepo func(rep *mocks.MockRepository)

		expectedURL string
		expectedErr error
	}{
		{
			name: "happy case",
			setUpRepo: func(rep *mocks.MockRepository) {
				rep.EXPECT().Get(mock.Anything, testKey).Return(testURL, nil).Once()
			},
			expectedURL: testURL,
			expectedErr: nil,
		},
		{
			name: "key not found",
			setUpRepo: func(rep *mocks.MockRepository) {
				rep.EXPECT().Get(mock.Anything, testKey).Return("", domain.ErrKeyNotFound).Once()
			},
			expectedURL: "",
			expectedErr: domain.ErrKeyNotFound,
		},
		{
			name: "internal error",
			setUpRepo: func(rep *mocks.MockRepository) {
				rep.EXPECT().Get(mock.Anything, testKey).Return("", errors.New("redis err")).Once()
			},
			expectedURL: "",
			expectedErr: domain.ErrInternalErr,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rep := mocks.NewMockRepository(t)
			tc.setUpRepo(rep)

			svc := NewRedirectURLService(rep)
			url, err := svc.RedirectURL(context.Background(), testKey)

			assert.Equal(t, tc.expectedURL, url)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
