package dtoBanners

import entityBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banners"

//go:generate ditgen -source=./banners.go -all-fields

type SetBanner struct {
	ID             *string
	Type           entityBanners.BannerType
	Position       *int
	DesktopImageID string
	MobileImageID  *string
	Url            string
}

