package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepHealthCheckService(t *testing.T) {
	t.Parallel()

	secret := "dial tcp 10.0.0.1:6379: AUTH failed password=hunter2"

	testCase := []struct {
		name    string
		pingErr error

		expectedErr bool
	}{
		{
			name:        "Happy case",
			pingErr:     nil,
			expectedErr: false,
		},
		{
			name:        "Ping error",
			pingErr:     errors.New(secret),
			expectedErr: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := NewDepHealthCheckService("redis", func(ctx context.Context) error {
				return tc.pingErr
			})

			err := svc.HealthCheck(context.Background())

			if !tc.expectedErr {
				assert.NoError(t, err)
				return
			}

			assert.ErrorIs(t, err, ErrDepUnhealthy)
			assert.Equal(t, "dependency unhealthy: redis", err.Error())
			assert.False(t, strings.Contains(err.Error(), "hunter2"))
		})
	}
}
