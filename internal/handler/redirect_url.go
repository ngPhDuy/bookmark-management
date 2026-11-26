package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngPhDuy/bookmark-management/internal/domain"
	"github.com/ngPhDuy/bookmark-management/internal/model/response"
	"github.com/ngPhDuy/bookmark-management/internal/service"
)

type RedirectURLHandler interface {
	RedirectURL(ctx *gin.Context)
}

type redirectURLHandler struct {
	svc service.RedirectURLService
}

func NewRedirectURLHandler(svc service.RedirectURLService) RedirectURLHandler {
	return &redirectURLHandler{
		svc: svc,
	}
}

// RedirectURL handles short link redirect requests.
//
// @Summary      Redirect to the original url
// @Description  Look up the original url by its short code and redirect to it.
// @Tags         Links
// @Produce      json
// @Param        code  path      string  true  "Short code returned by the shorten endpoint"
// @Success      302   {string}  string                  "Redirect to the original url"
// @Failure      404   {object}  response.ErrorResponse  "Key not found"
// @Failure      500   {object}  response.ErrorResponse  "Redirect url fail"
// @Router       /v1/links/redirect/{code} [get]
func (h *redirectURLHandler) RedirectURL(ctx *gin.Context) {
	code := ctx.Param("code")

	url, err := h.svc.RedirectURL(ctx, code)

	if errors.Is(err, domain.ErrKeyNotFound) {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "Key not found",
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: "Redirect URL fail",
		})
		return
	}

	ctx.Redirect(http.StatusFound, url)
}
