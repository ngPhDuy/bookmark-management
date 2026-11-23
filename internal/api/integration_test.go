package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/api"
	"github.com/ngPhDuy/bookmark-management/internal/config"
	"github.com/ngPhDuy/bookmark-management/internal/model/request"
	"github.com/ngPhDuy/bookmark-management/internal/model/response"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const originalURL = "https://www.google.com"

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T) (*httptest.Server, *miniredis.Miniredis, *http.Client) {
	t.Helper()

	mr := miniredis.RunT(t)

	rdb := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { rdb.Close() })

	srv := httptest.NewServer(api.NewEngine(&config.AppConfig{ServiceName: "test"}, rdb).Handler())
	t.Cleanup(srv.Close)

	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return srv, mr, client
}

func shorten(t *testing.T, srv *httptest.Server, client *http.Client, url string, exp int) string {
	t.Helper()

	body, err := json.Marshal(request.ShortenURLRequest{URL: url, Exp: exp})
	require.NoError(t, err)

	res, err := client.Post(srv.URL+"/v1/links/shorten", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	var parsed response.ShortenURLResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&parsed))
	require.NotEmpty(t, parsed.Code)

	return parsed.Code
}

func TestShortenThenRedirect(t *testing.T) {
	t.Parallel()

	srv, mr, client := newTestServer(t)

	code := shorten(t, srv, client, originalURL, 60)

	res, err := client.Get(srv.URL + "/v1/links/redirect/" + code)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusFound, res.StatusCode)
	assert.Equal(t, originalURL, res.Header.Get("Location"))

	stored, err := mr.Get(code)
	require.NoError(t, err)
	assert.Equal(t, originalURL, stored)
}

func TestRedirectAfterKeyExpired(t *testing.T) {
	t.Parallel()

	srv, mr, client := newTestServer(t)

	code := shorten(t, srv, client, originalURL, 60)

	mr.FastForward(61 * time.Second)

	res, err := client.Get(srv.URL + "/v1/links/redirect/" + code)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assertErrorMessage(t, res, "Key not found")
}

func TestRedirectUnknownCode(t *testing.T) {
	t.Parallel()

	srv, _, client := newTestServer(t)

	res, err := client.Get(srv.URL + "/v1/links/redirect/khong-ton-tai")
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assertErrorMessage(t, res, "Key not found")
}

func assertErrorMessage(t *testing.T, res *http.Response, expected string) {
	t.Helper()

	var parsed response.ErrorResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&parsed))
	assert.Equal(t, expected, parsed.Message)
}
