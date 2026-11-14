package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/domain"
	"github.com/ngPhDuy/bookmark-management/internal/model/request"
	service_mock "github.com/ngPhDuy/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURLHandler(t *testing.T) {
	t.Parallel()

	testcase := []struct {
		name string

		//setup
		req     interface{}
		rawBody string
		err     error
		key     string

		//expected
		expectedStatusCode int
		expectedResponse   string
		expectedSvcURL     string
		expectedSvcExp     int
	}{
		{
			name: "happy case",
			req: request.ShortenURLRequest{
				URL: "https://google.com",
				Exp: 60,
			},
			err:                nil,
			key:                "abracadabra",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"code":"abracadabra", "message":""}`,
			expectedSvcURL:     "https://google.com",
			expectedSvcExp:     60,
		},
		{
			name:               "missing exp field",
			req:                struct{ URL string }{URL: "https://google.com"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"code":"", "message":"invalid request format"}`,
		},
		{
			name: "url is not a valid url",
			req: request.ShortenURLRequest{
				URL: "not-a-url",
				Exp: 60,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"code":"", "message":"invalid request format"}`,
		},
		{
			name: "exp is not positive",
			req: request.ShortenURLRequest{
				URL: "https://google.com",
				Exp: 0,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"code":"", "message":"invalid request format"}`,
		},
		{
			name:               "malformed json body",
			rawBody:            `{"url":`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"code":"", "message":"invalid request format"}`,
		},
		{
			name: "internal server error",
			req: request.ShortenURLRequest{
				URL: "https://google.com",
				Exp: 60,
			},
			err:                errors.New("internal server error"),
			key:                "abracadabra",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"code":"", "message":"could not generate a unique code"}`,
			expectedSvcURL:     "https://google.com",
			expectedSvcExp:     60,
		},
		{
			name: "max retry",
			req: request.ShortenURLRequest{
				URL: "https://google.com",
				Exp: 60,
			},
			err:                domain.ErrShortenUrlFail,
			key:                "abracadabra",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"code":"", "message":"could not generate a unique code"}`,
			expectedSvcURL:     "https://google.com",
			expectedSvcExp:     60,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := service_mock.NewMockShortenURLService(tc.key, tc.err)
			h := NewShortenURLHandler(svc)

			rec := httptest.NewRecorder()
			router := gin.New()
			router.POST("/v1/links/shorten", h.ShortenURL)

			//parse body req to bytes
			body := []byte(tc.rawBody)

			if tc.rawBody == "" {
				marshaled, err := json.Marshal(tc.req)
				require.NoError(t, err)
				body = marshaled
			}

			//mock req with bytes body
			req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(body))
			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())

			//zero values also assert the service was never reached
			assert.Equal(t, tc.expectedSvcURL, svc.GotURL)
			assert.Equal(t, tc.expectedSvcExp, svc.GotExp)
		})
	}
}
