package proxy

import (
	"context"
	"errors"
	"fmt"
	"testing"

	bannerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/banner/v1"
	eventv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/event/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
)

func Test_proxyRepository_ListHomeBanners(t *testing.T) {
	type fields struct {
		client *bannerv1.MockBannerAPIClient
		mapper *MockProxyMapper
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx       context.Context
		sessionID string
		portalURL string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityBanner.BannersList, error)
	}{
		{
			name: "success",
			args: args{
				ctx:       context.Background(),
				sessionID: "test-session-id",
			},
			want: func(a args, f fields) (*entityBanner.BannersList, error) {
				mockResponse := &bannerv1.ListHomeBannersResponse{
					Items: []*bannerv1.Banner{
						{
							ImageUrl: "http://example1.com/test.jpg",
							Url:      "http://example1.com/1",
							Order:    0,
						},
					},
					Interval: 1,
				}
				expectedEntity := &entityBanner.BannersList{
					Items: []*entityBanner.Banner{
						{
							ImageUrl: "http://example1.com/test.jpg",
							Url:      "http://example1.com/1",
							Order:    0,
						},
					},
					Interval: 1,
				}

				f.client.EXPECT().
					ListHomeBanners(a.ctx, &bannerv1.ListHomeBannersRequest{
						PortalSessionId: a.sessionID,
						PortalUrl:       a.portalURL,
					}).
					Return(mockResponse, nil)

				f.mapper.EXPECT().
					BannersListToEntity(mockResponse).
					Return(expectedEntity)

				return expectedEntity, nil
			},
		},
		{
			name: "error not found",
			args: args{
				ctx:       context.Background(),
				sessionID: "test-session-id",
			},
			want: func(a args, f fields) (*entityBanner.BannersList, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(diterrors.ErrNotFound))
				f.client.EXPECT().
					ListHomeBanners(a.ctx, &bannerv1.ListHomeBannersRequest{
						PortalSessionId: a.sessionID,
						PortalUrl:       a.portalURL,
					}).
					Return(nil, diterrors.ErrNotFound)

				return nil, ErrNotFound
			},
		},
		{
			name: "error validation",
			args: args{
				ctx:       context.Background(),
				sessionID: "test-session-id",
			},
			want: func(a args, f fields) (*entityBanner.BannersList, error) {
				validationErr := diterrors.NewValidationError(errors.New("validation error"))
				f.logger.EXPECT().Warn("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(validationErr))
				f.client.EXPECT().
					ListHomeBanners(a.ctx, &bannerv1.ListHomeBannersRequest{
						PortalSessionId: a.sessionID,
						PortalUrl:       a.portalURL,
					}).
					Return(nil, validationErr)

				return nil, fmt.Errorf("invalid request list home banners: %w", validationErr)
			},
		},
		{
			name: "error unauthorized",
			args: args{
				ctx:       context.Background(),
				sessionID: "test-session-id",
			},
			want: func(a args, f fields) (*entityBanner.BannersList, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(diterrors.ErrUnauthenticated))
				f.client.EXPECT().
					ListHomeBanners(a.ctx, &bannerv1.ListHomeBannersRequest{
						PortalSessionId: a.sessionID,
						PortalUrl:       a.portalURL,
					}).
					Return(nil, diterrors.ErrUnauthenticated)

				return nil, ErrUnauthorized
			},
		},
		{
			name: "error permission denied",
			args: args{
				ctx:       context.Background(),
				sessionID: "test-session-id",
			},
			want: func(a args, f fields) (*entityBanner.BannersList, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(diterrors.ErrPermissionDenied))
				f.client.EXPECT().
					ListHomeBanners(a.ctx, &bannerv1.ListHomeBannersRequest{
						PortalSessionId: a.sessionID,
						PortalUrl:       a.portalURL,
					}).
					Return(nil, diterrors.ErrPermissionDenied)

				return nil, ErrPermissionDenied
			},
		},
		{
			name: "error internal",
			args: args{
				ctx:       context.Background(),
				sessionID: "test-session-id",
			},
			want: func(a args, f fields) (*entityBanner.BannersList, error) {
				err := errors.New("some internal error")
				f.logger.EXPECT().Error("proxyRepository.ListHomeBanners: failed to list home banners", zap.Error(err))
				f.client.EXPECT().
					ListHomeBanners(a.ctx, &bannerv1.ListHomeBannersRequest{
						PortalSessionId: a.sessionID,
						PortalUrl:       a.portalURL,
					}).
					Return(nil, err)

				return nil, ErrInternal
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				client: bannerv1.NewMockBannerAPIClient(ctrl),
				mapper: NewMockProxyMapper(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			r := NewProxyRepository(f.client, nil, f.mapper, f.logger)
			got, err := r.ListHomeBanners(tt.args.ctx, tt.args.sessionID, tt.args.portalURL)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_proxyRepository_ListCalendarEvents(t *testing.T) {
	type fields struct {
		bannersAPIClient *bannerv1.MockBannerAPIClient
		eventsAPIClient  *eventv1.MockEventAPIClient
		logger           *ditzap.MockLogger
		mapper           *MockProxyMapper
	}
	type args struct {
		ctx context.Context
		req entityEvent.CalendarEventRequest
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityEvent.CalendarEventsList, error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					Year:      2024,
					Month:     11,
				},
			},
			want: func(a args, f fields) (*entityEvent.CalendarEventsList, error) {
				apiResponse := &eventv1.ListCalendarEventsResponse{
					// Populate with expected API response data
				}

				f.eventsAPIClient.EXPECT().
					ListCalendarEvents(a.ctx, &eventv1.ListCalendarEventsRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Year:            int32(a.req.Year),
						Month:           int32(a.req.Month),
					}).
					Return(apiResponse, nil)

				expectedList := &entityEvent.CalendarEventsList{
					// Populate with expected mapped data
				}

				f.mapper.EXPECT().
					CalendarEventsListToEntity(apiResponse).
					Return(expectedList)

				return expectedList, nil
			},
		},
		{
			name: "error not found",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					Year:      2024,
					Month:     11,
				},
			},
			want: func(a args, f fields) (*entityEvent.CalendarEventsList, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(diterrors.ErrNotFound))
				f.eventsAPIClient.EXPECT().
					ListCalendarEvents(a.ctx, &eventv1.ListCalendarEventsRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Year:            int32(a.req.Year),
						Month:           int32(a.req.Month),
					}).
					Return(nil, diterrors.ErrNotFound)

				return nil, ErrNotFound
			},
		},
		{
			name: "error validation",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					Year:      2024,
					Month:     11,
				},
			},
			want: func(a args, f fields) (*entityEvent.CalendarEventsList, error) {
				validationErr := diterrors.NewValidationError(errors.New("validation error"))
				f.logger.EXPECT().Warn("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(validationErr))
				f.eventsAPIClient.EXPECT().
					ListCalendarEvents(a.ctx, &eventv1.ListCalendarEventsRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Year:            int32(a.req.Year),
						Month:           int32(a.req.Month),
					}).
					Return(nil, validationErr)

				return nil, fmt.Errorf("invalid request list calendar events: %w", validationErr)
			},
		},
		{
			name: "error unauthenticated",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					Year:      2024,
					Month:     11,
				},
			},
			want: func(a args, f fields) (*entityEvent.CalendarEventsList, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(diterrors.ErrUnauthenticated))
				f.eventsAPIClient.EXPECT().
					ListCalendarEvents(a.ctx, &eventv1.ListCalendarEventsRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Year:            int32(a.req.Year),
						Month:           int32(a.req.Month),
					}).
					Return(nil, diterrors.ErrUnauthenticated)

				return nil, ErrUnauthorized
			},
		},
		{
			name: "error permission denied",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					Year:      2024,
					Month:     11,
				},
			},
			want: func(a args, f fields) (*entityEvent.CalendarEventsList, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(diterrors.ErrPermissionDenied))
				f.eventsAPIClient.EXPECT().
					ListCalendarEvents(a.ctx, &eventv1.ListCalendarEventsRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Year:            int32(a.req.Year),
						Month:           int32(a.req.Month),
					}).
					Return(nil, diterrors.ErrPermissionDenied)

				return nil, ErrPermissionDenied
			},
		},
		{
			name: "error internal",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					Year:      2024,
					Month:     11,
				},
			},
			want: func(a args, f fields) (*entityEvent.CalendarEventsList, error) {
				err := errors.New("internal error")
				f.logger.EXPECT().Error("proxyRepository.ListCalendarEvents: failed to list calendar events", zap.Error(err))
				f.eventsAPIClient.EXPECT().
					ListCalendarEvents(a.ctx, &eventv1.ListCalendarEventsRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Year:            int32(a.req.Year),
						Month:           int32(a.req.Month),
					}).
					Return(nil, err)

				return nil, ErrInternal
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				eventsAPIClient: eventv1.NewMockEventAPIClient(ctrl),
				logger:          ditzap.NewMockLogger(ctrl),
				mapper:          NewMockProxyMapper(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			r := NewProxyRepository(nil, f.eventsAPIClient, f.mapper, f.logger)
			got, err := r.ListCalendarEvents(tt.args.ctx, tt.args.req)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_proxyRepository_ListCalendarEventsLinks(t *testing.T) {
	type fields struct {
		bannersAPIClient *bannerv1.MockBannerAPIClient
		eventsAPIClient  *eventv1.MockEventAPIClient
		logger           *ditzap.MockLogger
		mapper           *MockProxyMapper
	}
	type args struct {
		ctx context.Context
		req entityEvent.CalendarEventLinksRequest
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*entityEvent.CalendarEventLink, error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventLinksRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					EventIDs:  []string{"event1", "event2"},
				},
			},
			want: func(a args, f fields) ([]*entityEvent.CalendarEventLink, error) {
				apiResponse := &eventv1.ListCalendarEventsLinksResponse{
					// Populate with expected API response data
				}

				f.eventsAPIClient.EXPECT().
					ListCalendarEventsLinks(a.ctx, &eventv1.ListCalendarEventsLinksRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Ids:             a.req.EventIDs,
					}).
					Return(apiResponse, nil)

				expectedLinks := []*entityEvent.CalendarEventLink{
					// Populate with expected mapped data
				}

				f.mapper.EXPECT().
					CalendarEventsLinksToEntity(apiResponse).
					Return(expectedLinks)

				return expectedLinks, nil
			},
		},
		{
			name: "error not found",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventLinksRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					EventIDs:  []string{"event1", "event2"},
				},
			},
			want: func(a args, f fields) ([]*entityEvent.CalendarEventLink, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(diterrors.ErrNotFound))
				f.eventsAPIClient.EXPECT().
					ListCalendarEventsLinks(a.ctx, &eventv1.ListCalendarEventsLinksRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Ids:             a.req.EventIDs,
					}).
					Return(nil, diterrors.ErrNotFound)

				return nil, ErrNotFound
			},
		},
		{
			name: "error validation",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventLinksRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					EventIDs:  []string{"event1", "event2"},
				},
			},
			want: func(a args, f fields) ([]*entityEvent.CalendarEventLink, error) {
				validationErr := diterrors.NewValidationError(errors.New("validation error"))
				f.logger.EXPECT().Warn("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(validationErr))
				f.eventsAPIClient.EXPECT().
					ListCalendarEventsLinks(a.ctx, &eventv1.ListCalendarEventsLinksRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Ids:             a.req.EventIDs,
					}).
					Return(nil, validationErr)

				return nil, fmt.Errorf("invalid request list calendar events links: %w", validationErr)
			},
		},
		{
			name: "error unauthenticated",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventLinksRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					EventIDs:  []string{"event1", "event2"},
				},
			},
			want: func(a args, f fields) ([]*entityEvent.CalendarEventLink, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(diterrors.ErrUnauthenticated))
				f.eventsAPIClient.EXPECT().
					ListCalendarEventsLinks(a.ctx, &eventv1.ListCalendarEventsLinksRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Ids:             a.req.EventIDs,
					}).
					Return(nil, diterrors.ErrUnauthenticated)

				return nil, ErrUnauthorized
			},
		},
		{
			name: "error permission denied",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventLinksRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					EventIDs:  []string{"event1", "event2"},
				},
			},
			want: func(a args, f fields) ([]*entityEvent.CalendarEventLink, error) {
				f.logger.EXPECT().Warn("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(diterrors.ErrPermissionDenied))
				f.eventsAPIClient.EXPECT().
					ListCalendarEventsLinks(a.ctx, &eventv1.ListCalendarEventsLinksRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Ids:             a.req.EventIDs,
					}).
					Return(nil, diterrors.ErrPermissionDenied)

				return nil, ErrPermissionDenied
			},
		},
		{
			name: "error internal",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventLinksRequest{
					SessionID: "test-session-id",
					PortalURL: "http://example.com",
					EventIDs:  []string{"event1", "event2"},
				},
			},
			want: func(a args, f fields) ([]*entityEvent.CalendarEventLink, error) {
				err := errors.New("internal error")
				f.logger.EXPECT().Error("proxyRepository.ListCalendarEventsLinks: failed to list calendar events links", zap.Error(err))
				f.eventsAPIClient.EXPECT().
					ListCalendarEventsLinks(a.ctx, &eventv1.ListCalendarEventsLinksRequest{
						PortalSessionId: a.req.SessionID,
						PortalUrl:       a.req.PortalURL,
						Ids:             a.req.EventIDs,
					}).
					Return(nil, err)

				return nil, ErrInternal
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				bannersAPIClient: bannerv1.NewMockBannerAPIClient(ctrl),
				eventsAPIClient:  eventv1.NewMockEventAPIClient(ctrl),
				logger:           ditzap.NewMockLogger(ctrl),
				mapper:           NewMockProxyMapper(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			r := NewProxyRepository(f.bannersAPIClient, f.eventsAPIClient, f.mapper, f.logger)
			got, err := r.ListCalendarEventsLinks(tt.args.ctx, tt.args.req)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
