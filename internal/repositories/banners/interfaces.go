package repositoryBanners

import (
	bannersv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/banners/banners/v1"
	dtoBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/banners"
	entityBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banners"
)

//go:generate mockgen -source=interfaces.go -destination=./banners_mock.go -package=repositoryBanners

type BannersMapper interface {
	SetBannersToPb(authorID string, portalID int, banners []*dtoBanners.SetBanner) *bannersv1.SetRequest
	BannerToPb(banners *dtoBanners.SetBanner) *bannersv1.SetRequest_Banner

	BannerInfoToEntity(info *bannersv1.SetResponse_Banner) *entityBanners.BannerInfo
	BannerInfosToEntities(info []*bannersv1.SetResponse_Banner) []*entityBanners.BannerInfo

	BannerToEntity(banner *bannersv1.Banner) *entityBanners.Banner
	BannersToEntities(banners []*bannersv1.Banner) []*entityBanners.Banner
	ContentToEntity(content *bannersv1.Banner_Content) entityBanners.Content
}
