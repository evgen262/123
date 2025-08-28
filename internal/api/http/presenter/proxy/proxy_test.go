package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	viewBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/banner"
	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
)

func Test_proxyPresenter_BannersListToView(t *testing.T) {
	type args struct {
		banners *entityBanner.BannersList
	}
	tests := []struct {
		name string
		args args
		want *viewBanner.BannersList
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
				banners: &entityBanner.BannersList{
					Items:    []*entityBanner.Banner{},
					Interval: 10,
				},
			},
			want: &viewBanner.BannersList{
				Items:    []*viewBanner.Banner{},
				Interval: 10,
			},
		},
		{
			name: "non-empty list",
			args: args{
				banners: &entityBanner.BannersList{
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
					Interval: 15,
				},
			},
			want: &viewBanner.BannersList{
				Items: []*viewBanner.Banner{
					{
						ImageUrl: "http://image1.com",
						Url:      "http://example1.com",
					},
					{
						ImageUrl: "http://image2.com",
						Url:      "http://example2.com",
					},
				},
				Interval: 15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewProxyPresenter()
			got := b.BannersListToView(tt.args.banners)
			assert.Equal(t, tt.want, got)
		})
	}
}