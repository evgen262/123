package proxy

import (
	"context"
	"errors"
	"fmt"

	bannerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/banner/v1"
	eventv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/event/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"

	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
)

type proxyRepository struct {
	bannersAPIClient bannerv1.BannerAPIClient
	eventsAPIClient  eventv1.EventAPIClient

	mapper ProxyMapper
	logger ditzap.Logger
}

func NewProxyRepository(
	bannersAPIClient bannerv1.BannerAPIClient,
	eventsAPIClient eventv1.EventAPIClient,
	mapper ProxyMapper,
	logger ditzap.Logger,
) *proxyRepository {
	return &proxyRepository{
		bannersAPIClient: bannersAPIClient,
		eventsAPIClient:  eventsAPIClient,
		mapper:           mapper,
		logger:           logger,
	}
}

func (r proxyRepository) ListHomeBanners(ctx context.Context, sessionID, portalURL string) (*entityBanner.BannersList, error) {
	bannersList, err := r.bannersAPIClient.ListHomeBanners(ctx, &bannerv1.ListHomeBannersRequest{
		PortalSessionId: sessionID,
		PortalUrl:       portalURL,
	})
	if err != nil {
		switch {
		case errors.Is(err, diterrors.ErrNotFound):
			r.logger.Warn("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(err))
			return nil, ErrNotFound
		case errors.As(err, new(diterrors.ValidationError)):
			r.logger.Warn("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(err))
			return nil, fmt.Errorf("invalid request list home banners: %w", err)
		case errors.Is(err, diterrors.ErrUnauthenticated):
			r.logger.Warn("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(err))
			return nil, ErrUnauthorized
		case errors.Is(err, diterrors.ErrPermissionDenied):
			r.logger.Warn("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(err))
			return nil, ErrPermissionDenied
		default:
			r.logger.Error("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(err))
			return nil, ErrInternal
		}
	}

	return r.mapper.BannersListToEntity(bannersList), nil
}

func (r proxyRepository) ListCalendarEvents(ctx context.Context, req entityEvent.CalendarEventRequest) (*entityEvent.CalendarEventsList, error) {
	list, err := r.eventsAPIClient.ListCalendarEvents(ctx, &eventv1.ListCalendarEventsRequest{
		PortalSessionId: req.SessionID,
		PortalUrl:       req.PortalURL,
		Year:            int32(req.Year),
		Month:           int32(req.Month),
	})
	if err != nil {
		switch {
		case errors.Is(err, diterrors.ErrNotFound):
			r.logger.Warn("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(err))
			return nil, ErrNotFound
		case errors.As(err, new(diterrors.ValidationError)):
			r.logger.Warn("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(err))
			return nil, fmt.Errorf("invalid request list calendar events: %w", err)
		case errors.Is(err, diterrors.ErrUnauthenticated):
			r.logger.Warn("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(err))
			return nil, ErrUnauthorized
		case errors.Is(err, diterrors.ErrPermissionDenied):
			r.logger.Warn("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(err))
			return nil, ErrPermissionDenied
		default:
			r.logger.Error("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(err))
			return nil, ErrInternal
		}
	}

	return r.mapper.CalendarEventsListToEntity(list), nil
}

func (r proxyRepository) ListCalendarEventsLinks(ctx context.Context, req entityEvent.CalendarEventLinksRequest) ([]*entityEvent.CalendarEventLink, error) {
	list, err := r.eventsAPIClient.ListCalendarEventsLinks(ctx, &eventv1.ListCalendarEventsLinksRequest{
		PortalSessionId: req.SessionID,
		PortalUrl:       req.PortalURL,
		Ids:             req.EventIDs,
	})
	if err != nil {
		switch {
		case errors.Is(err, diterrors.ErrNotFound):
			r.logger.Warn("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(err))
			return nil, ErrNotFound
		case errors.As(err, new(diterrors.ValidationError)):
			r.logger.Warn("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(err))
			return nil, fmt.Errorf("invalid request list calendar events links: %w", err)
		case errors.Is(err, diterrors.ErrUnauthenticated):
			r.logger.Warn("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(err))
			return nil, ErrUnauthorized
		case errors.Is(err, diterrors.ErrPermissionDenied):
			r.logger.Warn("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(err))
			return nil, ErrPermissionDenied
		default:
			r.logger.Error("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(err))
			return nil, ErrInternal
		}
	}

	return r.mapper.CalendarEventsLinksToEntity(list), nil
}
