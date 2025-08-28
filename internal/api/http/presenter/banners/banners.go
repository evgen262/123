package presenterBanners

import (
	viewBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/banners"
	dtoBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/banners"
	entityBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banners"
)

type bannersPresenter struct{}

func NewBannersPresenter() *bannersPresenter {
	return &bannersPresenter{}
}

func (p bannersPresenter) SetBannerToDTO(view *viewBanners.SetBanner) *dtoBanners.SetBanner {
	if view == nil {
		return nil
	}

	return &dtoBanners.SetBanner{
		ID:             view.ID,
		Type:           p.BannerTypeToEntity(view.Type),
		Position:       view.Position,
		DesktopImageID: view.ImageDesktop,
		MobileImageID:  view.ImageMobile,
		Url:            view.Url,
	}
}

func (p bannersPresenter) SetBannersToDTOs(view *viewBanners.SetBanners) []*dtoBanners.SetBanner {
	if view == nil {
		return nil
	}

	result := make([]*dtoBanners.SetBanner, 0, len(view.Banners))
	for _, v := range view.Banners {
		dto := p.SetBannerToDTO(v)
		if dto != nil {
			result = append(result, dto)
		}
	}
	return result
}

func (p bannersPresenter) BannerTypeToEntity(t viewBanners.BannerType) entityBanners.BannerType {
	switch t {
	case viewBanners.BannerTypeBottom:
		return entityBanners.BannerTypeBottom

	case viewBanners.BannerTypePromo:
		return entityBanners.BannerTypePromo

	case viewBanners.BannerTypeSlider:
		return entityBanners.BannerTypeSlider

	default:
		return entityBanners.BannerTypeInvalid
	}
}

func (p bannersPresenter) BannerTypeToView(t entityBanners.BannerType) viewBanners.BannerType {
	switch t {
	case entityBanners.BannerTypeBottom:
		return viewBanners.BannerTypeBottom

	case entityBanners.BannerTypePromo:
		return viewBanners.BannerTypePromo

	case entityBanners.BannerTypeSlider:
		return viewBanners.BannerTypeSlider

	default:
		return viewBanners.BannerTypeInvalid
	}
}

func (p bannersPresenter) BannerInfoToView(banner *entityBanners.BannerInfo) *viewBanners.BannerInfo {
	if banner == nil {
		return nil
	}

	return &viewBanners.BannerInfo{
		ID:       banner.ID,
		Type:     p.BannerTypeToView(banner.Type),
		Position: banner.Position,
	}
}

func (p bannersPresenter) BannerInfosToViews(banners []*entityBanners.BannerInfo) []*viewBanners.BannerInfo {
	if len(banners) == 0 {
		return nil
	}

	result := make([]*viewBanners.BannerInfo, 0, len(banners))
	for _, v := range banners {
		dto := p.BannerInfoToView(v)
		if dto != nil {
			result = append(result, dto)
		}
	}
	return result
}

func (p bannersPresenter) BannerToView(banner *entityBanners.Banner) *viewBanners.Banner {
	if banner == nil {
		return nil
	}

	return &viewBanners.Banner{
		ID:        banner.ID,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
		Content:   p.ContentToView(banner.Content),
	}
}

func (p bannersPresenter) ContentToView(content entityBanners.Content) viewBanners.Content {
	return viewBanners.Content{
		Position:       content.Position,
		DesktopImageID: content.DesktopImageID,
		MobileImageID:  content.MobileImageID,
		Url:            content.Url,
	}
}

func (p bannersPresenter) BannersToViews(banners []*entityBanners.Banner) []*viewBanners.Banner {
	result := make([]*viewBanners.Banner, 0, len(banners))

	if len(banners) == 0 {
		return result
	}

	for _, v := range banners {
		dto := p.BannerToView(v)
		if dto != nil {
			result = append(result, dto)
		}
	}
	return result
}
