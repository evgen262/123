package portalsv2

import (
	"context"
	"fmt"
	"testing"

	complexesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/complexes/v1"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestComplexesRepository_Filter(t *testing.T) {
	type fields struct {
		complexesClient *complexesv1.MockComplexesAPIClient
		mapper          *MockComplexesMapper
	}
	type args struct {
		ctx     context.Context
		filters *entityPortalsV2.FilterComplexesFilters
		options *entityPortalsV2.FilterComplexesOptions
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*entityPortalsV2.Complex, error)
	}{
		{
			name: "complexesAPIClient error",
			args: args{
				ctx: ctx,
				filters: &entityPortalsV2.FilterComplexesFilters{
					IDs: []int{1, 2},
				},
				options: &entityPortalsV2.FilterComplexesOptions{
					WithDeleted: true,
				},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.Complex, error) {
				ids := make([]int32, 0, len(a.filters.IDs))
				for _, id := range a.filters.IDs {
					ids = append(ids, int32(id))
				}
				filtersPb := &complexesv1.FilterRequest_Filters{Ids: ids}
				optionsPb := &complexesv1.FilterRequest_Options{WithDeleted: a.options.WithDeleted}
				req := &complexesv1.FilterRequest{Filters: filtersPb, Options: optionsPb}

				f.mapper.EXPECT().ComplexesFiltersToPb(a.filters).Return(filtersPb)
				f.mapper.EXPECT().ComplexesOptionsToPb(a.options).Return(optionsPb)
				f.complexesClient.EXPECT().
					Filter(a.ctx, req).
					Return(nil, testErr)

				return nil, fmt.Errorf("complexesClient.Filter: can't get complexes: %w", diterrors.GrpcErrorToError(testErr))
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				filters: &entityPortalsV2.FilterComplexesFilters{
					IDs: []int{1, 2},
				},
				options: &entityPortalsV2.FilterComplexesOptions{
					WithDeleted: true,
				},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.Complex, error) {
				ids := make([]int32, 0, len(a.filters.IDs))
				for _, id := range a.filters.IDs {
					ids = append(ids, int32(id))
				}
				filtersPb := &complexesv1.FilterRequest_Filters{Ids: ids}
				optionsPb := &complexesv1.FilterRequest_Options{WithDeleted: a.options.WithDeleted}
				req := &complexesv1.FilterRequest{Filters: filtersPb, Options: optionsPb}

				complexesPb := []*complexesv1.Complex{
					{
						Id:          1,
						Name:        "Complex 1",
						Description: wrapperspb.String("Desc 1"),
					},
					{
						Id:          2,
						Name:        "Complex 2",
						Description: wrapperspb.String("Desc 2"),
					},
				}
				resp := &complexesv1.FilterResponse{Complexes: complexesPb}

				complexesEntity := []*entityPortalsV2.Complex{
					{
						ID:          1,
						Name:        "Complex 1",
						Description: testPtr("Desc 1"),
					},
					{
						ID:          2,
						Name:        "Complex 2",
						Description: testPtr("Desc 2"),
					},
				}

				f.mapper.EXPECT().ComplexesFiltersToPb(a.filters).Return(filtersPb)
				f.mapper.EXPECT().ComplexesOptionsToPb(a.options).Return(optionsPb)
				f.complexesClient.EXPECT().
					Filter(a.ctx, req).
					Return(resp, nil)
				f.mapper.EXPECT().ComplexesToEntity(complexesPb).Return(complexesEntity)

				return complexesEntity, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			complexesClient := complexesv1.NewMockComplexesAPIClient(ctrl)
			mapper := NewMockComplexesMapper(ctrl) // Use the mock mapper

			f := fields{
				complexesClient: complexesClient,
				mapper:          mapper,
			}

			want, wantErr := tt.want(tt.args, f)
			repo := NewComplexesRepository(f.complexesClient, f.mapper)

			got, err := repo.Filter(tt.args.ctx, tt.args.filters, tt.args.options)

			if wantErr != nil {
				assert.Empty(t, got)
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}
		})
	}
}

func testPtr(i string) *string {
	return &i
}
