package portals

import (
	"testing"
	"time"

	featuresv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/features/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

func TestPortalsMapper_NewFeatureToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		feature *portal.Feature
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *featuresv1.AddRequest_Feature
	}{
		{
			name: "correct",
			args: args{
				feature: &portal.Feature{
					Id:        1,
					Name:      "Test ImageID Name",
					Version:   "Test Feature Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				},
			},
			want: func(a args, f fields) *featuresv1.AddRequest_Feature {
				return &featuresv1.AddRequest_Feature{
					Name:    "Test ImageID Name",
					Version: "Test Feature Version",
					Enabled: false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewFeaturesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.NewFeatureToPb(tt.args.feature)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_NewFeaturesToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		feature []*portal.Feature
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*featuresv1.AddRequest_Feature
	}{
		{
			name: "correct",
			args: args{
				feature: []*portal.Feature{{
					Id:        1,
					Name:      "Test ImageID Name",
					Version:   "Test Feature Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				}},
			},
			want: func(a args, f fields) []*featuresv1.AddRequest_Feature {
				return []*featuresv1.AddRequest_Feature{{
					Name:    "Test ImageID Name",
					Version: "Test Feature Version",
					Enabled: false,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewFeaturesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.NewFeaturesToPb(tt.args.feature)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_FeatureToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		feature *portal.Feature
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *featuresv1.Feature
	}{
		{
			name: "correct",
			args: args{
				feature: &portal.Feature{
					Id:        1,
					Name:      "Test ImageID Name",
					Version:   "Test Feature Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				},
			},
			want: func(a args, f fields) *featuresv1.Feature {
				t := timestamppb.New(testT)
				f.timeUtils.EXPECT().TimeToTimestamp(a.feature.CreatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.feature.UpdatedAt).Return(t)
				return &featuresv1.Feature{
					Id:          1,
					Name:        "Test ImageID Name",
					Version:     "Test Feature Version",
					CreatedTime: t,
					UpdatedTime: t,
					Enabled:     false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewFeaturesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.FeatureToPb(tt.args.feature)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_FeaturesToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		features []*portal.Feature
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*featuresv1.Feature
	}{
		{
			name: "correct",
			args: args{
				features: []*portal.Feature{{
					Id:        1,
					Name:      "Test ImageID Name",
					Version:   "Test Feature Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				}},
			},
			want: func(a args, f fields) []*featuresv1.Feature {
				t := timestamppb.New(testT)
				f.timeUtils.EXPECT().TimeToTimestamp(a.features[0].CreatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.features[0].UpdatedAt).Return(t)
				return []*featuresv1.Feature{{
					Id:          1,
					Name:        "Test ImageID Name",
					Version:     "Test Feature Version",
					CreatedTime: t,
					UpdatedTime: t,
					Enabled:     false,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewFeaturesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.FeaturesToPb(tt.args.features)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_FeatureToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		feature *featuresv1.Feature
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)
	testTime := timestamppb.New(testT)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portal.Feature
	}{
		{
			name: "correct",
			args: args{
				feature: &featuresv1.Feature{
					Id:          1,
					Name:        "Test ImageID Name",
					Version:     "Test Feature Version",
					CreatedTime: testTime,
					UpdatedTime: testTime,
					Enabled:     false,
				},
			},
			want: func(a args, f fields) *portal.Feature {
				f.timeUtils.EXPECT().TimestampToTime(a.feature.CreatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.feature.UpdatedTime).Return(&testT)
				return &portal.Feature{
					Id:        1,
					Name:      "Test ImageID Name",
					Version:   "Test Feature Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewFeaturesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.FeatureToEntity(tt.args.feature)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_FeaturesToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		featuresPb []*featuresv1.Feature
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)
	testTime := timestamppb.New(testT)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*portal.Feature
	}{
		{
			name: "correct",
			args: args{
				featuresPb: []*featuresv1.Feature{{
					Id:          1,
					Name:        "Test ImageID Name",
					Version:     "Test Feature Version",
					CreatedTime: testTime,
					UpdatedTime: testTime,
					Enabled:     false,
				}},
			},
			want: func(a args, f fields) []*portal.Feature {
				f.timeUtils.EXPECT().TimestampToTime(a.featuresPb[0].CreatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.featuresPb[0].UpdatedTime).Return(&testT)
				return []*portal.Feature{{
					Id:        1,
					Name:      "Test ImageID Name",
					Version:   "Test Feature Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewFeaturesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.FeaturesToEntity(tt.args.featuresPb)

			assert.Equal(t, want, got)
		})
	}
}
