package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/model/request"
	"github.com/ngPhDuy/bookmark-management/internal/model/response"
	"github.com/ngPhDuy/bookmark-management/internal/service"
)

type ShortenURLHandler interface {
	ShortenURL(c *gin.Context)
}

type shortenURLHandler struct {
	shortenURLService service.ShortenURLService
}

func NewShortenURLHandler(svc service.ShortenURLService) ShortenURLHandler {
	return &shortenURLHandler{
		shortenURLService: svc,
	}
}

// ShortenURL handles shorten url requests.
//
// @Summary      Shorten url
// @Description  Gen a random key and save it to redis (with retry if key already exist). Return key if success.
// @Tags         Links
// @Accept       json
// @Produce      json
// @Param        request  body      request.ShortenURLRequest  true  "URL to shorten and its lifetime in seconds"
// @Success      200      {object}  response.ShortenURLResponse
// @Failure      400      {object}  response.ShortenURLResponse  "Invalid request body"
// @Failure      500      {object}  response.ShortenURLResponse  "Could not generate a unique code"
// @Router       /v1/links/shorten [post]
func (h *shortenURLHandler) ShortenURL(c *gin.Context) {
	var req request.ShortenURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ShortenURLResponse{
			Code:    "",
			Message: "invalid request format",
		})
		return
	}

	key, err := h.shortenURLService.ShortenURL(c.Request.Context(), req.URL, req.Exp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ShortenURLResponse{
			Code:    "",
			Message: "could not generate a unique code",
		})
		return
	}

	c.JSON(http.StatusOK, response.ShortenURLResponse{
		Code:    key,
		Message: "",
	})
}
