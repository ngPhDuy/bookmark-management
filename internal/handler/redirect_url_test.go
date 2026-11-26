package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/domain"
	"github.com/ngPhDuy/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testURL = "https://www.google.com"

func TestRedirectURLHandler(t *testing.T) {
	t.Parallel()

	testcase := []struct {
		name string

		//setup
		code string
		// nil means the service must never be reached
		setUpSvc func(svc *mocks.MockRedirectURLService)

		//expected
		expectedStatusCode int
		expectedLocation   string
		expectedResponse   string
	}{
		{
			name: "happy case",
			code: "abracadabra",
			setUpSvc: func(svc *mocks.MockRedirectURLService) {
				svc.EXPECT().RedirectURL(mock.Anything, "abracadabra").
					Return(testURL, nil).Once()
			},
			expectedStatusCode: http.StatusFound,
			expectedLocation:   testURL,
		},
		{
			name: "key not found",
			code: "unknown-key",
			setUpSvc: func(svc *mocks.MockRedirectURLService) {
				svc.EXPECT().RedirectURL(mock.Anything, "unknown-key").
					Return("", domain.ErrKeyNotFound).Once()
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"message":"Key not found"}`,
		},
		{
			name: "internal error",
			code: "abracadabra",
			setUpSvc: func(svc *mocks.MockRedirectURLService) {
				svc.EXPECT().RedirectURL(mock.Anything, "abracadabra").
					Return("", domain.ErrInternalErr).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Redirect URL fail"}`,
		},
		{
			name: "unknown error",
			code: "abracadabra",
			setUpSvc: func(svc *mocks.MockRedirectURLService) {
				svc.EXPECT().RedirectURL(mock.Anything, "abracadabra").
					Return("", errors.New("boom")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Redirect URL fail"}`,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewMockRedirectURLService(t)
			if tc.setUpSvc != nil {
				tc.setUpSvc(svc)
			}

			h := NewRedirectURLHandler(svc)

			rec := httptest.NewRecorder()
			router := gin.New()
			router.GET("/v1/links/redirect/:code", h.RedirectURL)

			req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/"+tc.code, nil)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)

			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, rec.Header().Get("Location"))
				return
			}

			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
