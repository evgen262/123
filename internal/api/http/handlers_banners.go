package http

import (
	"context"
	"errors"
	"net/http"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	viewBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/banners"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bannersHandlers struct {
	bannersInteractor BannersInteractor
	bannersPresenter  BannersPresenter
	logger            ditzap.Logger
}

func NewBannersHandlers(
	bannersInteractor BannersInteractor,
	bannersPresenter BannersPresenter,
	logger ditzap.Logger,
) *bannersHandlers {
	return &bannersHandlers{
		bannersInteractor: bannersInteractor,
		bannersPresenter:  bannersPresenter,
		logger:            logger,
	}
}

func (h *bannersHandlers) list(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	promo, slider, bottom, err := h.bannersInteractor.List(ctx)
	if err != nil {
		var (
			code     int
			response *view.Response
		)
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			code = http.StatusBadRequest
			response = view.NewErrorResponse(view.ErrMessageInvalidRequest)
		case errors.Is(err, usecase.ErrGetSessionFromContext):
			code = http.StatusUnauthorized
			response = view.NewErrorResponse(viewAuth.ErrMessageFrontUnauthenticated)
		default:
			code = http.StatusInternalServerError
			response = view.NewErrorResponse(view.ErrMessageInternalError)
		}
		c.JSON(code, response)
		return
	}

	type response struct {
		Promo  *viewBanners.Banner   `json:"promo,omitempty"`
		Slider []*viewBanners.Banner `json:"slider"`
		Bottom *viewBanners.Banner   `json:"bottom,omitempty"`
	}

	resp := &response{}

	if len(promo) > 0 {
		resp.Promo = h.bannersPresenter.BannerToView(promo[0])
	}

	if len(bottom) > 0 {
		resp.Bottom = h.bannersPresenter.BannerToView(bottom[0])
	}

	resp.Slider = h.bannersPresenter.BannersToViews(slider)

	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}

func (h *bannersHandlers) set(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	var setBanners viewBanners.SetBanners
	if err := c.BindJSON(&setBanners); err != nil {
		h.logger.Debug("can't bind banners into struct", zap.Error(err))
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	dto := h.bannersPresenter.SetBannersToDTOs(&setBanners)

	banners, err := h.bannersInteractor.Set(ctx, dto)
	if err != nil {
		var (
			code     int
			response *view.Response
		)
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			code = http.StatusBadRequest
			response = view.NewErrorResponse(view.ErrMessageInvalidRequest)
		case errors.Is(err, usecase.ErrGetSessionFromContext):
			code = http.StatusUnauthorized
			response = view.NewErrorResponse(viewAuth.ErrMessageFrontUnauthenticated)
		default:
			code = http.StatusInternalServerError
			response = view.NewErrorResponse(view.ErrMessageInternalError)
		}
		c.JSON(code, response)
		return
	}

	type response struct {
		Banners []*viewBanners.BannerInfo `json:"banners"`
	}

	resp := &response{
		Banners: h.bannersPresenter.BannerInfosToViews(banners),
	}
	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}
