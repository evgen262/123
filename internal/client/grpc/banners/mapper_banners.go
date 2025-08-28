package mapperBanners

import (
	bannersv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/banners/banners/v1"
	dtoBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/banners"
	entityBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banners"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
)

type bannersMapper struct {
	tu           timeUtils.TimeUtils
	sharedMapper repositories.SharedMapper
}

func NewBannersMapper(tu timeUtils.TimeUtils, sharedMapper repositories.SharedMapper) *bannersMapper {
	return &bannersMapper{
		tu:           tu,
		sharedMapper: sharedMapper,
	}
}

func (m bannersMapper) SetBannersToPb(authorID string, portalID int, banners []*dtoBanners.SetBanner) *bannersv1.SetRequest {
	pbs := make([]*bannersv1.SetRequest_Banner, 0, len(banners))
	for _, v := range banners {
		pb := m.BannerToPb(v)
		if pb != nil {
			pbs = append(pbs, pb)
		}
	}
	return &bannersv1.SetRequest{
		Banners:  pbs,
		AuthorId: authorID,
		PortalId: int32(portalID),
	}
}

func (m bannersMapper) BannerToPb(banner *dtoBanners.SetBanner) *bannersv1.SetRequest_Banner {
	if banner == nil {
		return nil
	}

	return &bannersv1.SetRequest_Banner{
		Id:             m.sharedMapper.StringValue(banner.ID),
		Type:           bannersv1.BannerType(banner.Type),
		Position:       m.sharedMapper.Int32Value(banner.Position),
		DesktopImageId: banner.DesktopImageID,
		MobileImageId:  m.sharedMapper.StringValue(banner.MobileImageID),
		Url:            banner.Url,
	}
}

func (m bannersMapper) BannerInfoToEntity(info *bannersv1.SetResponse_Banner) *entityBanners.BannerInfo {
	if info == nil {
		return nil
	}

	var pos *int
	if info.Position != nil {
		p := int(info.Position.GetValue())
		pos = &p
	}

	return &entityBanners.BannerInfo{
		ID:       info.Id,
		Type:     entityBanners.BannerType(info.Type),
		Position: pos,
	}
}

func (m bannersMapper) BannerInfosToEntities(infos []*bannersv1.SetResponse_Banner) []*entityBanners.BannerInfo {
	if len(infos) == 0 {
		return nil
	}

	banners := make([]*entityBanners.BannerInfo, 0, len(infos))
	for _, v := range infos {
		b := m.BannerInfoToEntity(v)
		if b != nil {
			banners = append(banners, b)
		}
	}
	return banners
}

func (m bannersMapper) BannerToEntity(banner *bannersv1.Banner) *entityBanners.Banner {
	if banner == nil {
		return nil
	}

	return &entityBanners.Banner{
		ID:         banner.Id,
		BannerType: entityBanners.BannerType(banner.Type),
		CreatedAt:  *m.tu.TimestampToTime(banner.CreatedTime),
		UpdatedAt:  m.tu.TimestampToTime(banner.UpdatedTime),
		Content:    m.ContentToEntity(banner.Content),
	}
}

func (m bannersMapper) BannersToEntities(banners []*bannersv1.Banner) []*entityBanners.Banner {
	if len(banners) == 0 {
		return nil
	}

	entites := make([]*entityBanners.Banner, 0, len(banners))
	for _, v := range banners {
		b := m.BannerToEntity(v)
		if b != nil {
			entites = append(entites, b)
		}
	}
	return entites
}

func (m bannersMapper) ContentToEntity(content *bannersv1.Banner_Content) entityBanners.Content {
	if content == nil {
		return entityBanners.Content{}
	}

	var pos *int
	if content.Position != nil {
		p := int(content.GetPosition().GetValue())
		pos = &p
	}

	var mobileImageID *string
	if content.MobileImageId != nil {
		mobileImageID = &content.MobileImageId.Value
	}

	return entityBanners.Content{
		Position:       pos,
		DesktopImageID: content.DesktopImageId,
		MobileImageID:  mobileImageID,
		Url:            content.Url,
	}
}
