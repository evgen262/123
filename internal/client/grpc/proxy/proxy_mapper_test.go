package proxy

import (
	"testing"

	bannerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/banner/v1"
	"github.com/stretchr/testify/assert"

	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
)

func Test_proxyMapper_BannersListToEntity(t *testing.T) {
	type args struct {
		banners *bannerv1.ListHomeBannersResponse
	}
	tests := []struct {
		name string
		args args
		want *entityBanner.BannersList
	}{
		{
			name: "nil input",
			args: args{
				banners: nil,
			},
			want: nil,
		},
		{
			name: "empty list",
			args: args{
				banners: &bannerv1.ListHomeBannersResponse{
					Interval: 0,
					Items:    nil,
				},
			},
			want: &entityBanner.BannersList{
				Interval: 0,
				Items:    []*entityBanner.Banner{},
			},
		},
		{
			name: "single banner",
			args: args{
				banners: &bannerv1.ListHomeBannersResponse{
					Interval: 10,
					Items: []*bannerv1.Banner{
						{
							ImageUrl: "http://image.com",
							Url:      "http://example.com",
						},
					},
				},
			},
			want: &entityBanner.BannersList{
				Interval: 10,
				Items: []*entityBanner.Banner{
					{
						ImageUrl: "http://image.com",
						Url:      "http://example.com",
					},
				},
			},
		},
		{
			name: "multiple banners",
			args: args{
				banners: &bannerv1.ListHomeBannersResponse{
					Interval: 15,
					Items: []*bannerv1.Banner{
						{
							ImageUrl: "http://image1.com",
							Url:      "http://example1.com",
						},
						{
							ImageUrl: "http://image2.com",
							Url:      "http://example2.com",
						},
					},
				},
			},
			want: &entityBanner.BannersList{
				Interval: 15,
				Items: []*entityBanner.Banner{
					{
						ImageUrl: "http://image1.com",
						Url:      "http://example1.com",
					},
					{
						ImageUrl: "http://image2.com",
						Url:      "http://example2.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := proxyMapper{}
			got := m.BannersListToEntity(tt.args.banners)
			assert.Equal(t, tt.want, got)
		})
	}
}
