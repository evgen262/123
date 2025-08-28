package viewBanners

import "time"

type SetBanners struct {
	Banners []*SetBanner `json:"banners"`
}

type SetBanner struct {
	ID           *string    `json:"id,omitempty"`
	Type         BannerType `json:"type"`
	Position     *int       `json:"position,omitempty"`
	ImageDesktop string     `json:"imageWeb"`
	ImageMobile  *string    `json:"imageMobile,omitempty"`
	Url          string     `json:"url"`
}

type BannerInfo struct {
	ID       string     `json:"id"`
	Type     BannerType `json:"type"`
	Position *int       `json:"position,omitempty"`
}

type Banner struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Content   Content    `json:"content"`
}

type Content struct {
	Position       *int    `json:"position,omitempty"`
	DesktopImageID string  `json:"imageWeb"`
	MobileImageID  *string `json:"imageMobile,omitempty"`
	Url            string  `json:"url"`
}

type BannerType string

const (
	BannerTypeInvalid BannerType = ""
	BannerTypePromo   BannerType = "promo"
	BannerTypeSlider  BannerType = "slider"
	BannerTypeBottom  BannerType = "bottom"
)
