package entityBanners

import "time"

//go:generate ditgen -source=./banners.go -all-fields

type BannerInfo struct {
	ID       string
	Type     BannerType
	Position *int
}

type Banner struct {
	ID         string
	BannerType BannerType
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	Content    Content
}

func (b *Banner) Type() BannerType {
	switch b.BannerType {
	case BannerTypeBottom, BannerTypePromo, BannerTypeSlider:
		return b.BannerType
	default:
		return BannerTypeInvalid
	}
}

type Content struct {
	Position       *int
	DesktopImageID string
	MobileImageID  *string
	Url            string
}

type BannerType uint8

const (
	BannerTypeInvalid BannerType = iota
	BannerTypePromo
	BannerTypeSlider
	BannerTypeBottom
)
