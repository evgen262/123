package portals

import (
	"testing"
	"time"

	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"

	"github.com/stretchr/testify/assert"
)

func Test_featurePresenter_ToEntity(t *testing.T) {
	type args struct {
		feature *viewPortals.UpdateFeature
	}
	tests := []struct {
		name string
		args args
		want *portal.Feature
	}{
		{
			name: "from feature to portal",
			args: args{
				feature: &viewPortals.UpdateFeature{
					Id:      5,
					Name:    "Test Feature Name",
					Version: "Version 0.0.1",
					Enabled: true,
				},
			},
			want: &portal.Feature{
				Id:      5,
				Name:    "Test Feature Name",
				Version: "Version 0.0.1",
				Enabled: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFeaturePresenter()
			assert.Equalf(t, tt.want, fp.ToEntity(tt.args.feature), "ToEntity(%v)", tt.args.feature)
		})
	}
}

func Test_featurePresenter_ToEntities(t *testing.T) {
	type args struct {
		features []*viewPortals.UpdateFeature
	}
	tests := []struct {
		name string
		args args
		want []*portal.Feature
	}{
		{
			name: "from features to entities",
			args: args{
				features: []*viewPortals.UpdateFeature{
					{
						Id:      1,
						Name:    "Test Feature 1",
						Version: "Version 0.0.1",
						Enabled: true,
					},
					{
						Id:      2,
						Name:    "Test Feature 2",
						Version: "Version 0.0.1",
						Enabled: false,
					},
				},
			},
			want: []*portal.Feature{
				{
					Id:      1,
					Name:    "Test Feature 1",
					Version: "Version 0.0.1",
					Enabled: true,
				},
				{
					Id:      2,
					Name:    "Test Feature 2",
					Version: "Version 0.0.1",
					Enabled: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFeaturePresenter()
			assert.Equalf(t, tt.want, fp.ToEntities(tt.args.features), "ToEntities(%v)", tt.args.features)
		})
	}
}

func Test_featurePresenter_ToNewEntity(t *testing.T) {
	type args struct {
		feature *viewPortals.NewFeature
	}
	tests := []struct {
		name string
		args args
		want *portal.Feature
	}{
		{
			name: "from feature to new portal",
			args: args{
				feature: &viewPortals.NewFeature{
					Name:    "Test Feature Name",
					Version: "Version 0.0.1",
					Enabled: true,
				},
			},
			want: &portal.Feature{
				Name:    "Test Feature Name",
				Version: "Version 0.0.1",
				Enabled: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFeaturePresenter()
			assert.Equalf(t, tt.want, fp.ToNewEntity(tt.args.feature), "ToNewEntity(%v)", tt.args.feature)
		})
	}
}

func Test_featurePresenter_ToNewEntities(t *testing.T) {
	type args struct {
		features []*viewPortals.NewFeature
	}
	tests := []struct {
		name string
		args args
		want []*portal.Feature
	}{
		{
			name: "from features to new entities",
			args: args{
				features: []*viewPortals.NewFeature{
					{
						Name:    "Test Feature 1",
						Version: "Version 0.0.1",
						Enabled: true,
					},
					{
						Name:    "Test Feature 2",
						Version: "Version 0.0.1",
						Enabled: false,
					},
				},
			},
			want: []*portal.Feature{
				{
					Name:    "Test Feature 1",
					Version: "Version 0.0.1",
					Enabled: true,
				},
				{
					Name:    "Test Feature 2",
					Version: "Version 0.0.1",
					Enabled: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFeaturePresenter()
			assert.Equalf(t, tt.want, fp.ToNewEntities(tt.args.features), "ToNewEntities(%v)", tt.args.features)
		})
	}
}

func Test_featurePresenter_ToShortView(t *testing.T) {
	type args struct {
		feature *portal.Feature
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.FeatureInfo
	}{
		{
			name: "from portal to short viewPortals feature",
			args: args{
				feature: &portal.Feature{
					Id:        7,
					Name:      "Feature Name",
					Version:   "Version 0.0.1",
					CreatedAt: nil,
					UpdatedAt: nil,
					Enabled:   true,
				},
			},
			want: &viewPortals.FeatureInfo{
				Id:      7,
				Version: "Version 0.0.1",
				Enabled: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFeaturePresenter()
			assert.Equalf(t, tt.want, fp.ToShortView(tt.args.feature), "ToShortView(%v)", tt.args.feature)
		})
	}
}

func Test_featurePresenter_ToShortViews(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)
	type args struct {
		features []*portal.Feature
	}
	tests := []struct {
		name string
		args args
		want viewPortals.Features
	}{
		{
			name: "from entities to short views feature",
			args: args{
				features: []*portal.Feature{
					{
						Id:        1,
						Name:      "Test Feature 1",
						Version:   "Version 0.0.1",
						CreatedAt: &testTime,
						UpdatedAt: nil,
						Enabled:   true,
					},
					{
						Id:        9,
						Name:      "Test Feature 9",
						CreatedAt: &testTime,
						UpdatedAt: nil,
						Version:   "Version 0.0.1",
						Enabled:   false,
					},
				},
			},
			want: viewPortals.Features{
				"Test Feature 1": {
					Id:      1,
					Version: "Version 0.0.1",
					Enabled: true,
				},
				"Test Feature 9": {
					Id:      9,
					Version: "Version 0.0.1",
					Enabled: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFeaturePresenter()
			assert.Equalf(t, tt.want, fp.ToShortViews(tt.args.features), "ToShortViews(%v)", tt.args.features)
		})
	}
}

func Test_featurePresenter_ToView(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)
	type args struct {
		feature *portal.Feature
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.Feature
	}{
		{
			name: "from portal to viewPortals",
			args: args{
				feature: &portal.Feature{
					Id:        55,
					Name:      "Test Feature Name 55",
					Version:   "0.0.2",
					CreatedAt: &testTime,
					UpdatedAt: nil,
					Enabled:   true,
				},
			},
			want: &viewPortals.Feature{
				Id:        55,
				Name:      "Test Feature Name 55",
				Version:   "0.0.2",
				CreatedAt: &testTime,
				Enabled:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFeaturePresenter()
			assert.Equalf(t, tt.want, fp.ToView(tt.args.feature), "ToView(%v)", tt.args.feature)
		})
	}
}

func Test_featurePresenter_ToViews(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)
	type args struct {
		features []*portal.Feature
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.Feature
	}{
		{
			name: "from entities to views",
			args: args{
				features: []*portal.Feature{
					{
						Id:        1,
						Name:      "Test Feature Name 1",
						Version:   "0.0.1",
						CreatedAt: &testTime,
						UpdatedAt: &testTime,
						Enabled:   true,
					},
					{
						Id:        2,
						Name:      "Test Feature Name 2",
						Version:   "0.0.2",
						CreatedAt: &testTime,
						UpdatedAt: nil,
						Enabled:   false,
					},
				},
			},
			want: []*viewPortals.Feature{
				{
					Id:        1,
					Name:      "Test Feature Name 1",
					Version:   "0.0.1",
					CreatedAt: &testTime,
					UpdatedAt: &testTime,
					Enabled:   true,
				},
				{
					Id:        2,
					Name:      "Test Feature Name 2",
					Version:   "0.0.2",
					CreatedAt: &testTime,
					UpdatedAt: nil,
					Enabled:   false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFeaturePresenter()
			assert.Equalf(t, tt.want, fp.ToViews(tt.args.features), "ToViews(%v)", tt.args.features)
		})
	}
}
