package proxy

import (
	"context"
	"fmt"

	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
)

type proxyInteractor struct {
	proxyRepository ProxyRepository
}

func NewProxyInteractor(proxyRepository ProxyRepository) *proxyInteractor {
	return &proxyInteractor{
		proxyRepository: proxyRepository,
	}
}

func (i proxyInteractor) ListHomeBanners(ctx context.Context, sessionID, portalURL string) (*entityBanner.BannersList, error) {
	bannersList, err := i.proxyRepository.ListHomeBanners(ctx, sessionID, portalURL)
	if err != nil {
		return nil, fmt.Errorf("can't get list home banners: %w", err)
	}

	return bannersList, nil
}

func (i proxyInteractor) ListCalendarEvents(ctx context.Context, req entityEvent.CalendarEventRequest) (*entityEvent.CalendarEventsList, error) {
	eventsList, err := i.proxyRepository.ListCalendarEvents(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("can't get list calendar events: %w", err)
	}

	return eventsList, nil
}

func (i proxyInteractor) ListCalendarEventsLinks(ctx context.Context, req entityEvent.CalendarEventLinksRequest) ([]*entityEvent.CalendarEventLink, error) {
	eventsLinks, err := i.proxyRepository.ListCalendarEventsLinks(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("can't get list calendar events links: %w", err)
	}

	return eventsLinks, nil
}
