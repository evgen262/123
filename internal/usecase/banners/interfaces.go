package usecaseBanners

import (
	"context"

	dtoBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/banners"
	entityBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banners"
)

//go:generate mockgen -source=interfaces.go -destination=./banners_mock.go -package=usecaseBanners

type BannersRepository interface {
	List(ctx context.Context, portalID int) ([]*entityBanners.Banner, error)
	Set(ctx context.Context, authorID string, portalID int, banners []*dtoBanners.SetBanner) ([]*entityBanners.BannerInfo, error)
}
