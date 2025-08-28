package repositoryBanners

import (
	"context"
	"fmt"

	bannersv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/banners/banners/v1"
	dtoBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/banners"
	entityBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banners"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
)

type bannersRepository struct {
	bannersApi    bannersv1.BannersAPIClient
	bannersMapper BannersMapper
	logger        ditzap.Logger
}

func NewBannersRepository(bannersApi bannersv1.BannersAPIClient, bannersMapper BannersMapper, logger ditzap.Logger) *bannersRepository {
	return &bannersRepository{
		bannersApi:    bannersApi,
		bannersMapper: bannersMapper,
		logger:        logger,
	}
}

func (r *bannersRepository) List(ctx context.Context, portalID int) ([]*entityBanners.Banner, error) {
	resp, err := r.bannersApi.List(ctx, &bannersv1.ListRequest{PortalId: int32(portalID)})
	if err != nil {
		return nil, fmt.Errorf("bannersRepository.List: can't list banners: %w", diterrors.GrpcErrorToError(err))
	}

	return r.bannersMapper.BannersToEntities(resp.GetBanners()), nil
}

func (r *bannersRepository) Set(ctx context.Context, authorID string,
	portalID int, banners []*dtoBanners.SetBanner) ([]*entityBanners.BannerInfo, error) {
	resp, err := r.bannersApi.Set(ctx, r.bannersMapper.SetBannersToPb(authorID, portalID, banners))
	if err != nil {
		return nil, fmt.Errorf("bannersRepository.Set: can't set banners: %w", diterrors.GrpcErrorToError(err))
	}

	return r.bannersMapper.BannerInfosToEntities(resp.GetBanners()), nil
}
