package proxy

import (
	"context"

	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
)

//go:generate mockgen -source=interfaces.go -destination=./proxy_mock.go -package=proxy

type ProxyRepository interface {
	ListHomeBanners(ctx context.Context, sessionID, portalURL string) (*entityBanner.BannersList, error)
	ListCalendarEvents(ctx context.Context, req entityEvent.CalendarEventRequest) (*entityEvent.CalendarEventsList, error)
	ListCalendarEventsLinks(ctx context.Context, req entityEvent.CalendarEventLinksRequest) ([]*entityEvent.CalendarEventLink, error)
}
