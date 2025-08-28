package proxy

import (
	viewBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/banner"
	viewEvents "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/events"
	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
)

type proxyPresenter struct {
}

func NewProxyPresenter() *proxyPresenter {
	return &proxyPresenter{}
}

func (b proxyPresenter) bannerToView(banner *entityBanner.Banner) *viewBanner.Banner {
	if banner == nil {
		return nil
	}
	return &viewBanner.Banner{
		ImageUrl: banner.ImageUrl,
		Url:      banner.Url,
		Order:    banner.Order,
	}
}

func (b proxyPresenter) BannersListToView(banners *entityBanner.BannersList) *viewBanner.BannersList {
	if banners == nil {
		return nil
	}

	items := make([]*viewBanner.Banner, 0, len(banners.Items))
	for _, banner := range banners.Items {
		items = append(items, b.bannerToView(banner))
	}

	return &viewBanner.BannersList{
		Items:    items,
		Interval: banners.Interval,
	}
}

func (b proxyPresenter) eventToView(event *entityEvent.CalendarEvent) *viewEvents.CalendarEvent {
	if event == nil {
		return nil
	}
	return &viewEvents.CalendarEvent{
		ID:    event.ID,
		Title: event.Title,
		Time:  event.Time,
		Date:  event.Date,
	}
}

func (b proxyPresenter) EventsListToView(events *entityEvent.CalendarEventsList) *viewEvents.CalendarEventsList {
	if events == nil {
		return nil
	}

	items := make([]*viewEvents.CalendarEvent, 0, len(events.Items))
	for _, event := range events.Items {
		items = append(items, b.eventToView(event))
	}

	return &viewEvents.CalendarEventsList{
		Items: items,
		Count: events.Count,
	}
}

func (b proxyPresenter) EventsLinksToView(events []*entityEvent.CalendarEventLink) []*viewEvents.CalendarEventLink {
	if events == nil {
		return nil
	}

	items := make([]*viewEvents.CalendarEventLink, 0, len(events))
	for _, event := range events {
		items = append(items, &viewEvents.CalendarEventLink{
			ID:          event.ID,
			RedirectURL: event.RedirectURL,
		})
	}

	return items
}