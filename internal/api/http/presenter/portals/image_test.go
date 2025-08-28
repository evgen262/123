package portals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

func Test_imagePresenter_ToNewEntity(t *testing.T) {
	type args struct {
		image *viewPortals.NewImage
	}
	tests := []struct {
		name string
		args args
		want *portal.Image
	}{
		{
			name: "from viewPortals to new portal",
			args: args{
				image: &viewPortals.NewImage{
					Name: "image_name.jpg",
					Data: "R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",
				},
			},
			want: &portal.Image{
				Name: "image_name.jpg",
				Data: portal.ImageData("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := NewImagePresenter()
			assert.Equalf(t, tt.want, ip.ToNewEntity(tt.args.image), "ToNewEntity(%v)", tt.args.image)
		})
	}
}

func Test_imagePresenter_ToEntity(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)

	type args struct {
		image *viewPortals.Image
	}
	tests := []struct {
		name string
		args args
		want *portal.Image
	}{
		{
			name: "from viewPortals to portal",
			args: args{
				image: &viewPortals.Image{
					Id:        25,
					Name:      "image_name.jpg",
					Path:      "/some/path/to/image_name.jpg",
					Data:      "R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",
					CreatedAt: &testTime,
					UpdatedAt: &testTime,
				},
			},
			want: &portal.Image{
				Id:        25,
				Name:      "image_name.jpg",
				Data:      []byte("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
				Path:      "/some/path/to/image_name.jpg",
				CreatedAt: &testTime,
				UpdatedAt: &testTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := NewImagePresenter()
			assert.Equalf(t, tt.want, ip.ToEntity(tt.args.image), "ToEntity(%v)", tt.args.image)
		})
	}
}

func Test_imagePresenter_ToEntities(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)

	type args struct {
		images []*viewPortals.Image
	}
	tests := []struct {
		name string
		args args
		want []*portal.Image
	}{
		{
			name: "from views to entities",
			args: args{
				images: []*viewPortals.Image{
					{
						Id:        1,
						Name:      "image_name_1.jpg",
						Path:      "/some/path/to/image_name_1.jpg",
						Data:      "R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",
						CreatedAt: &testTime,
						UpdatedAt: nil,
					},
					{
						Id:        2,
						Name:      "image_name_2.jpg",
						Path:      "/some/path/to/image_name_2.jpg",
						Data:      "R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",
						CreatedAt: &testTime,
						UpdatedAt: &testTime,
					},
				},
			},
			want: []*portal.Image{
				{
					Id:        1,
					Name:      "image_name_1.jpg",
					Path:      "/some/path/to/image_name_1.jpg",
					Data:      []byte("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
					CreatedAt: &testTime,
					UpdatedAt: nil,
				},
				{
					Id:        2,
					Name:      "image_name_2.jpg",
					Path:      "/some/path/to/image_name_2.jpg",
					Data:      []byte("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
					CreatedAt: &testTime,
					UpdatedAt: &testTime,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := NewImagePresenter()
			assert.Equalf(t, tt.want, ip.ToEntities(tt.args.images), "ToEntities(%v)", tt.args.images)
		})
	}
}

func Test_imagePresenter_ToShortView(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)

	type args struct {
		image *portal.Image
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.ImageInfo
	}{
		{
			name: "from portal to short viewPortals",
			args: args{
				&portal.Image{
					Id:        1,
					Name:      "image_name_1.jpg",
					Data:      portal.ImageData("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
					Path:      "/some/path/to/image_name_1.jpg",
					CreatedAt: &testTime,
					UpdatedAt: &testTime,
				},
			},
			want: &viewPortals.ImageInfo{
				Id:   1,
				Name: "image_name_1.jpg",
				Path: "/some/path/to/image_name_1.jpg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := NewImagePresenter()
			assert.Equalf(t, tt.want, ip.ToShortView(tt.args.image), "ToShortView(%v)", tt.args.image)
		})
	}
}

func Test_imagePresenter_ToShortViews(t *testing.T) {
	type args struct {
		images []*portal.Image
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.ImageInfo
	}{
		{
			name: "from entities to short views",
			args: args{
				images: []*portal.Image{
					{
						Id:   1,
						Name: "image_name_1.jpg",
						Path: "/some/path/to/image_name_1.jpg",
					},
					{
						Id:   2,
						Name: "image_name_2.jpg",
						Path: "/some/path/to/image_name_2.jpg",
					},
					{
						Id:   3,
						Name: "image_name_3.jpg",
						Path: "/some/path/to/image_name_3.jpg",
					},
				},
			},
			want: []*viewPortals.ImageInfo{
				{
					Id:   1,
					Name: "image_name_1.jpg",
					Path: "/some/path/to/image_name_1.jpg",
				},
				{
					Id:   2,
					Name: "image_name_2.jpg",
					Path: "/some/path/to/image_name_2.jpg",
				},
				{
					Id:   3,
					Name: "image_name_3.jpg",
					Path: "/some/path/to/image_name_3.jpg",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := NewImagePresenter()
			assert.Equalf(t, tt.want, ip.ToShortViews(tt.args.images), "ToShortViews(%v)", tt.args.images)
		})
	}
}

func Test_imagePresenter_ToView(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)

	type args struct {
		image *portal.Image
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.Image
	}{
		{
			name: "from portal to viewPortals",
			args: args{
				image: &portal.Image{
					Id:        1,
					Name:      "image_name_1.jpg",
					Data:      portal.ImageData("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
					Path:      "/some/path/to/image_name_1.jpg",
					CreatedAt: &testTime,
					UpdatedAt: &testTime,
				},
			},
			want: &viewPortals.Image{
				Id:        1,
				Name:      "image_name_1.jpg",
				Data:      "R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",
				Path:      "/some/path/to/image_name_1.jpg",
				CreatedAt: &testTime,
				UpdatedAt: &testTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := imagePresenter{}
			assert.Equalf(t, tt.want, ip.ToView(tt.args.image), "ToView(%v)", tt.args.image)
		})
	}
}

func Test_imagePresenter_ToViews(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)

	type args struct {
		images []*portal.Image
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.Image
	}{
		{
			name: "from entities to views",
			args: args{
				images: []*portal.Image{
					{
						Id:        1,
						Name:      "image_name_1.jpg",
						Data:      portal.ImageData("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
						Path:      "/some/path/to/image_name_1.jpg",
						CreatedAt: &testTime,
						UpdatedAt: &testTime,
					},
					{
						Id:        2,
						Name:      "image_name_2.jpg",
						Data:      portal.ImageData("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
						Path:      "/some/path/to/image_name_2.jpg",
						CreatedAt: &testTime,
						UpdatedAt: nil,
					},
					{
						Id:        3,
						Name:      "image_name_3.jpg",
						Data:      portal.ImageData("R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="),
						Path:      "/some/path/to/image_name_3.jpg",
						CreatedAt: &testTime,
						UpdatedAt: nil,
					},
				},
			},
			want: []*viewPortals.Image{
				{
					Id:        1,
					Name:      "image_name_1.jpg",
					Data:      "R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",
					Path:      "/some/path/to/image_name_1.jpg",
					CreatedAt: &testTime,
					UpdatedAt: &testTime,
				},
				{
					Id:        2,
					Name:      "image_name_2.jpg",
					Data:      "R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",
					Path:      "/some/path/to/image_name_2.jpg",
					CreatedAt: &testTime,
					UpdatedAt: nil,
				},
				{
					Id:        3,
					Name:      "image_name_3.jpg",
					Data:      "R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==",
					Path:      "/some/path/to/image_name_3.jpg",
					CreatedAt: &testTime,
					UpdatedAt: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := NewImagePresenter()
			assert.Equalf(t, tt.want, ip.ToViews(tt.args.images), "ToViews(%v)", tt.args.images)
		})
	}
}
