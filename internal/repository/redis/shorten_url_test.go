package redis

import (
	"context"
	"errors"
	"testing"
	"time"

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

	//To test usecase in mock repository
	countSetNXCall int
	lastExp        time.Duration
}

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{
		store: map[string]string{},
	}
}

func (m *mockRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.BoolCmd {
	m.countSetNXCall++
	m.lastExp = expiration

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

// mock random key generator
func sequenceKeys(keys ...string) func(int) (string, error) {
	i := 0
	return func(int) (string, error) {
		k := keys[min(i, len(keys)-1)]
		i++
		return k, nil
	}
}

// mock repository
func newMockRepository(client redisClient, gen func(int) (string, error)) RedisRepository {
	return &redisRepository{
		rdbClient:          client,
		randomKeyGenerator: gen,
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
	rep := newMockRepository(client, sequenceKeys("key1"))

	//set
	_, err := rep.Set(context.Background(), "https://www.google.com", 10)

	//no error, store: key1-"https://www.google.com", count 1, lastExp 10
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"key1": "https://www.google.com"}, client.store)
	assert.Equal(t, 1, client.countSetNXCall)
	assert.Equal(t, time.Duration(10)*time.Second, client.lastExp)
}

func TestSet_OneCollision(t *testing.T) {
	client := newMockRedisClient()
	client.store["key1"] = "this is my slot"

	rep := newMockRepository(client, sequenceKeys("key1", "key2"))

	_, err := rep.Set(context.Background(), "https://www.google.com", 10)

	//no error, store: key2-"https://www.google.com", count 2, lastExp 10
	assert.NoError(t, err)
	assert.Equal(t, "https://www.google.com", client.store["key2"])
	assert.Equal(t, 2, client.countSetNXCall)
	assert.Equal(t, time.Duration(10)*time.Second, client.lastExp)
}

func TestSet_NoSlot(t *testing.T) {
	client := newMockRedisClient()
	client.store["key1"] = "this is my slot"

	rep := newMockRepository(client, sequenceKeys("key1"))

	_, err := rep.Set(context.Background(), "https://www.google.com", 10)

	//error, store: key1-"this is my slot", count 1, lastExp 10
	assert.Error(t, err)
	assert.Equal(t, "this is my slot", client.store["key1"])
	assert.Equal(t, MAX_ATTEMPTS, client.countSetNXCall)
	assert.Equal(t, time.Duration(10)*time.Second, client.lastExp)
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

			rep := newMockRepository(client, sequenceKeys())

			gotVal, gotErr := rep.Get(context.Background(), tc.key)

			assert.Equal(t, tc.wantValue, gotVal)
			assert.Equal(t, tc.wantErr, gotErr)
		})
	}
}
