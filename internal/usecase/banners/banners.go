package usecaseBanners

import (
	"context"
	"fmt"

	dtoBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/banners"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banners"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
)

type bannersInteractor struct {
	bannersRepo BannersRepository
	logger      ditzap.Logger
}

func NewBannersInteractor(bannersRepository BannersRepository, logger ditzap.Logger) *bannersInteractor {
	return &bannersInteractor{
		bannersRepo: bannersRepository,
		logger:      logger,
	}
}

func (i *bannersInteractor) List(ctx context.Context) ([]*entityBanners.Banner, []*entityBanners.Banner, []*entityBanners.Banner, error) {
	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		return nil, nil, nil, usecase.ErrGetSessionFromContext
	}
	portalID := session.GetActivePortal().GetPortal().ID

	result, err := i.bannersRepo.List(ctx, portalID)
	if err != nil {
		i.logger.Error("can't list banners", zap.Error(err))
		return nil, nil, nil, fmt.Errorf("bannersInteractor.List: can't list banners: %w", err)
	}

	promo, slider, bottom := []*entityBanners.Banner{}, []*entityBanners.Banner{}, []*entityBanners.Banner{}
	for _, v := range result {
		switch v.Type() {
		case entityBanners.BannerTypePromo:
			promo = append(promo, v)
		case entityBanners.BannerTypeSlider:
			slider = append(slider, v)
		case entityBanners.BannerTypeBottom:
			bottom = append(bottom, v)
		default:
			i.logger.Warn("got banner with invalid type")
		}
	}

	return promo, slider, bottom, nil
}

func (i *bannersInteractor) Set(ctx context.Context, banners []*dtoBanners.SetBanner) ([]*entityBanners.BannerInfo, error) {
	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		return nil, usecase.ErrGetSessionFromContext
	}

	portalID := session.GetActivePortal().GetPortal().ID
	authorID := session.GetUser().GetEmployee().GetExtID()

	result, err := i.bannersRepo.Set(ctx, authorID, portalID, banners)
	if err != nil {
		i.logger.Error("can't set banners", zap.Error(err))
		return nil, fmt.Errorf("bannersInteractor.Set: can't set banners: %w", err)
	}
	return result, nil
}
