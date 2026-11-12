package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ngPhDuy/bookmark-management/internal/model"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// Mock redis client
type mockRedisClient struct {
	//Store simulation
	store map[string]string

	//Errors simulation
	setNXErr error
	getErr   error
}

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{
		store: map[string]string{},
	}
}

func (m *mockRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.BoolCmd {

	// error connection simulation
	if m.setNXErr != nil {
		return goredis.NewBoolResult(false, m.setNXErr)
	}

	// check if key already exist
	if _, exist := m.store[key]; exist {
		return goredis.NewBoolResult(false, m.setNXErr)
	}

	m.store[key] = value.(string)
	return goredis.NewBoolResult(true, nil)
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *goredis.StringCmd {
	if m.getErr != nil {
		return goredis.NewStringResult("", m.getErr)
	}

	value, exist := m.store[key]

	// value not found
	if !exist {
		return goredis.NewStringResult("", goredis.Nil)
	}

	return goredis.NewStringResult(value, nil)
}

// mock repository
func newMockRepository(client redisClient) model.Repository {
	return &redisRepository{
		rdbClient: client,
	}
}

// Test for Mock
func TestMockGet(t *testing.T) {
	client := newMockRedisClient()
	client.store["key2"] = "val2"

	val, err := client.Get(context.Background(), "key2").Result()

	assert.NoError(t, err)
	assert.Equal(t, "val2", val)
}

// Test for SET
func TestSet_Success(t *testing.T) {
	//mock redis client
	client := newMockRedisClient()
	//create mock repo
	rep := newMockRepository(client)

	//set
	_, err := rep.SetIfNotExist(context.Background(), "key1", "https://www.google.com", 10)

	//no error, store: key1-"https://www.google.com"
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"key1": "https://www.google.com"}, client.store)
}

func TestSet_NoSlot(t *testing.T) {
	client := newMockRedisClient()
	client.store["key1"] = "this is my slot"
	client.store["key2"] = "this is 2 slot"

	rep := newMockRepository(client)

	_, err := rep.SetIfNotExist(context.Background(), "key2", "https://www.google.com", 10)

	//error -> ErrKeyExists, store: key2-"this is 2 slot"
	assert.ErrorIs(t, err, ErrKeyExists)
	assert.Equal(t, "this is 2 slot", client.store["key2"])
}

func TestGetFunc(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		store     map[string]string
		key       string
		getErr    error
		wantValue string
		wantErr   error
	}{
		{
			name: "success",
			store: map[string]string{
				"key1": "val1", "key2": "val2", "key3": "val3",
			},
			getErr:    nil,
			key:       "key2",
			wantValue: "val2",
			wantErr:   nil,
		},
		{
			name:      "error_connection",
			store:     map[string]string{},
			key:       "key2",
			getErr:    errors.New("error"),
			wantValue: "",
			wantErr:   errors.New("error"),
		},
		{
			name: "key_not_found",
			store: map[string]string{
				"key1": "val1", "key2": "val2", "key3": "val3",
			},
			key:       "key4",
			wantValue: "",
			wantErr:   ErrKeyNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newMockRedisClient()
			client.store = tc.store
			client.getErr = tc.getErr

			rep := newMockRepository(client)

			gotVal, gotErr := rep.Get(context.Background(), tc.key)

			assert.Equal(t, tc.wantValue, gotVal)
			assert.Equal(t, tc.wantErr, gotErr)
		})
	}
}
