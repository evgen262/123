package proxy

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
)

func Test_proxyInteractor_ListHomeBanners(t *testing.T) {
	type fields struct {
		proxyRepository *MockProxyRepository
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
				expectedList := &entityBanner.BannersList{
					Interval: 10,
					Items: []*entityBanner.Banner{
						{
							ImageUrl: "http://image.com",
							Url:      "http://example.com",
						},
					},
				}

				f.proxyRepository.EXPECT().
					ListHomeBanners(a.ctx, a.sessionID, a.portalURL).
					Return(expectedList, nil)

				return expectedList, nil
			},
		},
		{
			name: "error fetching banners",
			args: args{
				ctx:       context.Background(),
				sessionID: "test-session-id",
			},
			want: func(a args, f fields) (*entityBanner.BannersList, error) {
				f.proxyRepository.EXPECT().
					ListHomeBanners(a.ctx, a.sessionID, a.portalURL).
					Return(nil, fmt.Errorf("repository error"))

				return nil, fmt.Errorf("can't get list home banners: %w", fmt.Errorf("repository error"))
			},
		},
		// Add more test cases for different error scenarios if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				proxyRepository: NewMockProxyRepository(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			i := NewProxyInteractor(f.proxyRepository)
			got, err := i.ListHomeBanners(tt.args.ctx, tt.args.sessionID, tt.args.portalURL)
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

func Test_proxyInteractor_ListCalendarEvents(t *testing.T) {
	type fields struct {
		proxyRepository *MockProxyRepository
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
					// заполните необходимые поля запроса
				},
			},
			want: func(a args, f fields) (*entityEvent.CalendarEventsList, error) {
				expectedList := &entityEvent.CalendarEventsList{
					// заполните ожидаемый список событий
				}

				f.proxyRepository.EXPECT().
					ListCalendarEvents(a.ctx, a.req).
					Return(expectedList, nil)

				return expectedList, nil
			},
		},
		{
			name: "error fetching events",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventRequest{
					// заполните необходимые поля запроса
				},
			},
			want: func(a args, f fields) (*entityEvent.CalendarEventsList, error) {
				f.proxyRepository.EXPECT().
					ListCalendarEvents(a.ctx, a.req).
					Return(nil, fmt.Errorf("repository error"))

				return nil, fmt.Errorf("can't get list calendar events: %w", fmt.Errorf("repository error"))
			},
		},
		// Добавьте дополнительные тестовые случаи для других ошибок, если необходимо
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				proxyRepository: NewMockProxyRepository(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			i := NewProxyInteractor(f.proxyRepository)
			got, err := i.ListCalendarEvents(tt.args.ctx, tt.args.req)
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

func Test_proxyInteractor_ListCalendarEventsLinks(t *testing.T) {
	type fields struct {
		proxyRepository *MockProxyRepository
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
				req: entityEvent.CalendarEventLinksRequest{},
			},
			want: func(a args, f fields) ([]*entityEvent.CalendarEventLink, error) {
				expectedLinks := []*entityEvent.CalendarEventLink{
					{},
				}

				f.proxyRepository.EXPECT().
					ListCalendarEventsLinks(a.ctx, a.req).
					Return(expectedLinks, nil)

				return expectedLinks, nil
			},
		},
		{
			name: "error fetching links",
			args: args{
				ctx: context.Background(),
				req: entityEvent.CalendarEventLinksRequest{},
			},
			want: func(a args, f fields) ([]*entityEvent.CalendarEventLink, error) {
				f.proxyRepository.EXPECT().
					ListCalendarEventsLinks(a.ctx, a.req).
					Return(nil, fmt.Errorf("repository error"))

				return nil, fmt.Errorf("can't get list calendar events links: %w", fmt.Errorf("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				proxyRepository: NewMockProxyRepository(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			i := NewProxyInteractor(f.proxyRepository)
			got, err := i.ListCalendarEventsLinks(tt.args.ctx, tt.args.req)
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
