package portalsv2

import (
	"testing"
	"time"

	complexesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/complexes/v1"
	entityPortalsv2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestComplexesMapper_ComplexesFiltersToPb(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		filters *entityPortalsv2.FilterComplexesFilters
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *complexesv1.FilterRequest_Filters
	}{
		{
			name: "nil filters",
			args: args{
				filters: nil,
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Filters {
				return nil
			},
		},
		{
			name: "empty filters",
			args: args{
				filters: &entityPortalsv2.FilterComplexesFilters{
					IDs:       []int{},
					PortalIDs: []int{},
				},
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Filters {
				return &complexesv1.FilterRequest_Filters{
					Ids:       []int32{},
					PortalIds: []int32{},
				}
			},
		},
		{
			name: "filters with data",
			args: args{
				filters: &entityPortalsv2.FilterComplexesFilters{
					IDs:       []int{1, 5, 10},
					PortalIDs: []int{100, 200},
				},
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Filters {
				return &complexesv1.FilterRequest_Filters{
					Ids:       []int32{1, 5, 10},
					PortalIds: []int32{100, 200},
				}
			},
		},
		{
			name: "filters with only IDs",
			args: args{
				filters: &entityPortalsv2.FilterComplexesFilters{
					IDs:       []int{42},
					PortalIDs: []int{},
				},
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Filters {
				return &complexesv1.FilterRequest_Filters{
					Ids:       []int32{42},
					PortalIds: []int32{},
				}
			},
		},
		{
			name: "filters with only PortalIDs",
			args: args{
				filters: &entityPortalsv2.FilterComplexesFilters{
					IDs:       []int{},
					PortalIDs: []int{99, 88, 77},
				},
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Filters {
				return &complexesv1.FilterRequest_Filters{
					Ids:       []int32{},
					PortalIds: []int32{99, 88, 77},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{tu: timeUtils.NewMockTimeUtils(ctrl)}
			mapper := NewComplexesMapper(f.tu)
			expected := tt.want(tt.args, f)

			got := mapper.ComplexesFiltersToPb(tt.args.filters)
			assert.Equal(t, expected, got)
		})
	}
}

func TestComplexesMapper_ComplexesOptionsToPb(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		options *entityPortalsv2.FilterComplexesOptions
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *complexesv1.FilterRequest_Options
	}{
		{
			name: "nil options",
			args: args{
				options: nil,
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Options {
				return nil
			},
		},
		{
			name: "options with both true",
			args: args{
				options: &entityPortalsv2.FilterComplexesOptions{
					WithDeleted:  true,
					WithDisabled: true,
				},
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Options {
				return &complexesv1.FilterRequest_Options{
					WithDeleted:  true,
					WithDisabled: true,
				}
			},
		},
		{
			name: "options with WithDeleted true",
			args: args{
				options: &entityPortalsv2.FilterComplexesOptions{
					WithDeleted:  true,
					WithDisabled: false,
				},
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Options {
				return &complexesv1.FilterRequest_Options{
					WithDeleted:  true,
					WithDisabled: false,
				}
			},
		},
		{
			name: "options with WithDisabled true",
			args: args{
				options: &entityPortalsv2.FilterComplexesOptions{
					WithDeleted:  false,
					WithDisabled: true,
				},
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Options {
				return &complexesv1.FilterRequest_Options{
					WithDeleted:  false,
					WithDisabled: true,
				}
			},
		},
		{
			name: "options with both false",
			args: args{
				options: &entityPortalsv2.FilterComplexesOptions{
					WithDeleted:  false,
					WithDisabled: false,
				},
			},
			want: func(a args, f fields) *complexesv1.FilterRequest_Options {
				return &complexesv1.FilterRequest_Options{
					WithDeleted:  false,
					WithDisabled: false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{tu: timeUtils.NewMockTimeUtils(ctrl)}
			cm := NewComplexesMapper(f.tu)

			want := tt.want(tt.args, f)

			got := cm.ComplexesOptionsToPb(tt.args.options)
			assert.Equal(t, want, got)
		})
	}
}

func TestComplexesMapper_ComplexesToEntity(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		complexesPb []*complexesv1.Complex
	}

	testTime := time.Date(2024, 5, 21, 10, 30, 0, 0, time.UTC)
	testTimestamp := timestamppb.New(testTime)
	stringPtr := func(s string) *string { return &s }

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*entityPortalsv2.Complex
	}{
		{
			name: "nil input",
			args: args{
				complexesPb: nil,
			},
			want: func(a args, f fields) []*entityPortalsv2.Complex {
				return nil
			},
		},
		{
			name: "empty slice",
			args: args{
				complexesPb: []*complexesv1.Complex{},
			},
			want: func(a args, f fields) []*entityPortalsv2.Complex {
				return []*entityPortalsv2.Complex{}
			},
		},
		{
			name: "correct conversion with various data",
			args: args{
				complexesPb: []*complexesv1.Complex{
					{
						Id:           1,
						Name:         "Test Complex 1",
						Description:  wrapperspb.String("Complex Description 1"),
						ImageId:      wrapperspb.String("img-001"),
						ComplexGroup: 1,
						Sort:         10,
						IsDisabled:   false,
						Responsible: &complexesv1.Complex_Responsible{
							FirstName:   "John",
							LastName:    "Doe",
							MiddleName:  wrapperspb.String("M."),
							ImageId:     wrapperspb.String("resp-img-001"),
							Description: "Responsible person for complex 1",
						},
						Portals: []*complexesv1.Complex_Portal{
							{Id: 101, Sort: 1},
							{Id: 102, Sort: 2},
						},
						CreatedTime: testTimestamp,
						UpdatedTime: testTimestamp,
						DeletedTime: testTimestamp,
					},
					nil,
					{
						Id:           2,
						Name:         "Test Complex 2",
						Description:  nil,
						ImageId:      nil,
						ComplexGroup: 2,
						Sort:         20,
						IsDisabled:   true,
						Responsible:  nil,
						Portals:      nil, 
						CreatedTime:  testTimestamp,
						UpdatedTime:  testTimestamp,
						DeletedTime:  nil,
					},
				},
			},
			want: func(a args, f fields) []*entityPortalsv2.Complex {
				f.tu.EXPECT().TimestampToTime(a.complexesPb[0].GetUpdatedTime()).Return(&testTime)
				f.tu.EXPECT().TimestampToTime(a.complexesPb[0].GetDeletedTime()).Return(&testTime)

				f.tu.EXPECT().TimestampToTime(a.complexesPb[2].GetUpdatedTime()).Return(&testTime)
				f.tu.EXPECT().TimestampToTime(a.complexesPb[2].GetDeletedTime()).Return(nil)

				return []*entityPortalsv2.Complex{
					{
						ID:           1,
						Name:         "Test Complex 1",
						Description:  stringPtr("Complex Description 1"),
						ImageID:      stringPtr("img-001"),
						ComplexGroup: 1,
						Sort:         10,
						IsDisabled:   false,
						Responsible: &entityPortalsv2.ComplexResponsible{
							FirstName:   "John",
							LastName:    "Doe",
							MiddleName:  stringPtr("M."),
							ImageID:     stringPtr("resp-img-001"),
							Description: "Responsible person for complex 1",
						},
						Portals: []*entityPortalsv2.ComplexPortal{
							{ID: 101, Sort: 1},
							{ID: 102, Sort: 2},
						},
						CreatedAt: testTime,
						UpdatedAt: &testTime,
						DeletedAt: &testTime,
					},
					{
						ID:           2,
						Name:         "Test Complex 2",
						Description:  nil,
						ImageID:      nil,
						ComplexGroup: 2,
						Sort:         20,
						IsDisabled:   true,
						Responsible:  nil,
						Portals:      nil,
						CreatedAt:    testTime,
						UpdatedAt:    &testTime,
						DeletedAt:    nil,
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{tu: timeUtils.NewMockTimeUtils(ctrl)}
			cm := NewComplexesMapper(f.tu)
			want := tt.want(tt.args, f)
			got := cm.ComplexesToEntity(tt.args.complexesPb)

			assert.Equal(t, want, got)
		})
	}
}
