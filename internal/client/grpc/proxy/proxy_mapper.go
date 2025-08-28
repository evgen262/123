package proxy

import (
	bannerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/banner/v1"
	eventv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/event/v1"

	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
)

type proxyMapper struct{}

func NewProxyMapper() *proxyMapper {
	return &proxyMapper{}
}

func (m proxyMapper) bannerToEntity(banner *bannerv1.Banner) *entityBanner.Banner {
	if banner == nil {
		return nil
	}
	return &entityBanner.Banner{
		ImageUrl: banner.GetImageUrl(),
		Url:      banner.GetUrl(),
		Order:    int(banner.GetOrder()),
	}
}

func (m proxyMapper) BannersListToEntity(banners *bannerv1.ListHomeBannersResponse) *entityBanner.BannersList {
	if banners == nil {
		return nil
	}

	items := make([]*entityBanner.Banner, 0, len(banners.GetItems()))
	for _, banner := range banners.GetItems() {
		items = append(items, m.bannerToEntity(banner))
	}

	return &entityBanner.BannersList{
		Interval: int(banners.GetInterval()),
		Items:    items,
	}
}

func (m proxyMapper) calendarEventToEntity(event *eventv1.CalendarEvent) *entityEvent.CalendarEvent {
	if event == nil {
		return nil
	}
	return &entityEvent.CalendarEvent{
		ID:    event.GetId(),
		Title: event.GetTitle(),
		Time:  event.GetTime(),
		Date:  event.GetDate(),
	}
}

func (m proxyMapper) CalendarEventsListToEntity(list *eventv1.ListCalendarEventsResponse) *entityEvent.CalendarEventsList {
	if list == nil {
		return nil
	}

	items := make([]*entityEvent.CalendarEvent, 0, len(list.GetEvents()))
	for _, event := range list.GetEvents() {
		items = append(items, m.calendarEventToEntity(event))
	}

	return &entityEvent.CalendarEventsList{
		Items: items,
		Count: int(list.GetTotalCount()),
	}
}

func (m proxyMapper) CalendarEventsLinksToEntity(list *eventv1.ListCalendarEventsLinksResponse) []*entityEvent.CalendarEventLink {
	if list == nil {
		return nil
	}

	items := make([]*entityEvent.CalendarEventLink, 0, len(list.GetLinks()))
	for _, link := range list.GetLinks() {
		items = append(items, &entityEvent.CalendarEventLink{
			ID:          link.GetId(),
			RedirectURL: link.GetRedirectUrl(),
		})
	}

	return items
}
