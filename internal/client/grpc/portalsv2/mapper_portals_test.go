package portalsv2

import (
	"testing"
	"time"

	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/portals/v1"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestPortalsMapper_PortalsWithCountsToEntity(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		portalsWithCountsPb []*portalsv1.PortalWithCounts
	}

	testTime := time.Date(2024, 5, 21, 10, 0, 0, 0, time.UTC)
	testTimestamp := timestamppb.New(testTime)
	testTimePtr := &testTime

	imageIDStr := "image-uuid-123"
	middleNameStr := "Middle"

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*entityPortalsV2.PortalWithCounts
	}{
		{
			name: "nil input",
			args: args{
				portalsWithCountsPb: nil,
			},
			want: func(a args, f fields) []*entityPortalsV2.PortalWithCounts {
				return nil
			},
		},
		{
			name: "empty input",
			args: args{
				portalsWithCountsPb: []*portalsv1.PortalWithCounts{},
			},
			want: func(a args, f fields) []*entityPortalsV2.PortalWithCounts {
				return []*entityPortalsV2.PortalWithCounts{}
			},
		},
		{
			name: "correct conversion with nil items",
			args: args{
				portalsWithCountsPb: []*portalsv1.PortalWithCounts{
					{
						Portal: &portalsv1.Portal{
							Id:                1,
							Name:              "Portal One",
							ShortName:         "P1",
							Url:               "https://portal1.mos.ru",
							ImageId:           wrapperspb.String(imageIDStr),
							Status:            portalsv1.PortalStatus_PORTAL_STATUS_ACTIVE,
							StatusUpdatedTime: testTimestamp,
							IsDisabled:        false,
							Sort:              1,
							Manager: &portalsv1.PortalManager{
								FirstName:  "John",
								LastName:   "Doe",
								MiddleName: wrapperspb.String(middleNameStr),
								ImageId:    wrapperspb.String(imageIDStr),
							},
							CreatedTime: testTimestamp,
							UpdatedTime: testTimestamp,
							DeletedTime: testTimestamp,
						},
						OrgsCount:      10,
						EmployeesCount: 100,
					},
					nil,
					{
						Portal:         nil,
						OrgsCount:      5,
						EmployeesCount: 20,
					},
					{
						Portal: &portalsv1.Portal{
							Id:                2,
							Name:              "Portal Two",
							ShortName:         "P2",
							Url:               "https://portal2.mos.ru",
							ImageId:           nil,
							Status:            portalsv1.PortalStatus_PORTAL_STATUS_INACTIVE,
							StatusUpdatedTime: nil,
							IsDisabled:        true,
							Sort:              2,
							Manager:           nil,
							CreatedTime:       testTimestamp,
							UpdatedTime:       nil,
							DeletedTime:       nil,
						},
						OrgsCount:      0,
						EmployeesCount: 0,
					},
				},
			},
			want: func(a args, f fields) []*entityPortalsV2.PortalWithCounts {
				f.tu.EXPECT().TimestampToTime(a.portalsWithCountsPb[0].GetPortal().GetStatusUpdatedTime()).Return(testTimePtr)
				f.tu.EXPECT().TimestampToTime(a.portalsWithCountsPb[0].GetPortal().GetUpdatedTime()).Return(testTimePtr)
				f.tu.EXPECT().TimestampToTime(a.portalsWithCountsPb[0].GetPortal().GetDeletedTime()).Return(testTimePtr)

				f.tu.EXPECT().TimestampToTime(a.portalsWithCountsPb[3].GetPortal().GetStatusUpdatedTime()).Return(nil)
				f.tu.EXPECT().TimestampToTime(a.portalsWithCountsPb[3].GetPortal().GetUpdatedTime()).Return(nil)
				f.tu.EXPECT().TimestampToTime(a.portalsWithCountsPb[3].GetPortal().GetDeletedTime()).Return(nil)

				return []*entityPortalsV2.PortalWithCounts{
					{
						Portal: &entityPortalsV2.Portal{
							ID:                1,
							Name:              "Portal One",
							ShortName:         "P1",
							Url:               "https://portal1.mos.ru",
							ImageID:           &imageIDStr,
							Status:            entityPortalsV2.PortalStatusActive,
							StatusUpdatedTime: testTimePtr,
							IsDisabled:        false,
							Sort:              1,
							Manager: &entityPortalsV2.PortalManager{
								FirstName:  "John",
								LastName:   "Doe",
								MiddleName: &middleNameStr,
								ImageID:    &imageIDStr,
							},
							CreatedAt: testTime,
							UpdatedAt: testTimePtr,
							DeletedAt: testTimePtr,
						},
						OrgsCount:      10,
						EmployeesCount: 100,
					},
					{
						Portal: &entityPortalsV2.Portal{
							ID:                2,
							Name:              "Portal Two",
							ShortName:         "P2",
							Url:               "https://portal2.mos.ru",
							ImageID:           nil,
							Status:            entityPortalsV2.PortalStatusInactive,
							StatusUpdatedTime: nil,
							IsDisabled:        true,
							Sort:              2,
							Manager:           nil,
							CreatedAt:         testTime,
							UpdatedAt:         nil,
							DeletedAt:         nil,
						},
						OrgsCount:      0,
						EmployeesCount: 0,
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{tu: timeUtils.NewMockTimeUtils(ctrl)}
			m := NewPortalsMapper(f.tu)

			want := tt.want(tt.args, f)
			got := m.PortalsWithCountsToEntity(tt.args.portalsWithCountsPb)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_PortalsFiltersToPb(t *testing.T) {
	type fields struct {
		tu timeUtils.TimeUtils
	}
	type args struct {
		filters *entityPortalsV2.FilterPortalsFilters
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portalsv1.FilterRequest_Filters
	}{
		{
			name: "correct",
			args: args{
				filters: &entityPortalsV2.FilterPortalsFilters{
					IDs: []int{1, 2, 3},
				},
			},
			want: func(a args, f fields) *portalsv1.FilterRequest_Filters {
				return &portalsv1.FilterRequest_Filters{
					Ids: []int32{1, 2, 3},
				}
			},
		},
		{
			name: "empty ids",
			args: args{
				filters: &entityPortalsV2.FilterPortalsFilters{
					IDs: []int{},
				},
			},
			want: func(a args, f fields) *portalsv1.FilterRequest_Filters {
				return &portalsv1.FilterRequest_Filters{
					Ids: []int32{},
				}
			},
		},
		{
			name: "nil filters",
			args: args{
				filters: nil,
			},
			want: func(a args, f fields) *portalsv1.FilterRequest_Filters {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				tu: timeUtils.NewMockTimeUtils(ctrl),
			}
			m := NewPortalsMapper(f.tu)
			want := tt.want(tt.args, f)
			got := m.PortalsFiltersToPb(tt.args.filters)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_PortalsOptionsToPb(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		options *entityPortalsV2.FilterPortalsOptions
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portalsv1.FilterRequest_Options
	}{
		{
			name: "nil options",
			args: args{
				options: nil,
			},
			want: func(a args, f fields) *portalsv1.FilterRequest_Options {
				return nil
			},
		},
		{
			name: "correct with employees count",
			args: args{
				options: &entityPortalsV2.FilterPortalsOptions{
					WithEmployeesCount: true,
				},
			},
			want: func(a args, f fields) *portalsv1.FilterRequest_Options {
				return &portalsv1.FilterRequest_Options{
					WithEmployeesCount: true,
				}
			},
		},
		{
			name: "correct without employees count",
			args: args{
				options: &entityPortalsV2.FilterPortalsOptions{
					WithEmployeesCount: false,
				},
			},
			want: func(a args, f fields) *portalsv1.FilterRequest_Options {
				return &portalsv1.FilterRequest_Options{
					WithEmployeesCount: false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{tu: timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewPortalsMapper(f.tu)
			want := tt.want(tt.args, f)
			got := pfm.PortalsOptionsToPb(tt.args.options)

			assert.Equal(t, want, got)
		})
	}
}
