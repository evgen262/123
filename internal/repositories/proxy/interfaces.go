package proxy

import (
	bannerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/banner/v1"
	eventv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/event/v1"

	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
)

//go:generate mockgen -source=interfaces.go -destination=./proxy_mock.go -package=proxy

type ProxyMapper interface {
	BannersListToEntity(banners *bannerv1.ListHomeBannersResponse) *entityBanner.BannersList
	CalendarEventsListToEntity(list *eventv1.ListCalendarEventsResponse) *entityEvent.CalendarEventsList
	CalendarEventsLinksToEntity(list *eventv1.ListCalendarEventsLinksResponse) []*entityEvent.CalendarEventLink
}
