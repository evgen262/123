package portalsv2

import (
	"context"
	"fmt"
	"testing"

	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPortalsInteractor_Filter(t *testing.T) {
	type fields struct {
		portalsRepository *MockPortalsRepository
	}
	type args struct {
		ctx     context.Context
		filters *entityPortalsV2.FilterPortalsFilters
		options *entityPortalsV2.FilterPortalsOptions
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("simulated repository error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*entityPortalsV2.PortalWithCounts, error)
	}{
		{
			name: "repository error",
			args: args{
				ctx:     ctx,
				filters: &entityPortalsV2.FilterPortalsFilters{},
				options: &entityPortalsV2.FilterPortalsOptions{},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.PortalWithCounts, error) {
				f.portalsRepository.EXPECT().
					Filter(a.ctx, a.filters, a.options).
					Return(nil, testErr)

				return nil, fmt.Errorf("portalsRepository.Filter: can't filter portals: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				filters: &entityPortalsV2.FilterPortalsFilters{},
				options: &entityPortalsV2.FilterPortalsOptions{},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.PortalWithCounts, error) {
				expectedPortals := []*entityPortalsV2.PortalWithCounts{
					{
						Portal: &entityPortalsV2.Portal{
							ID:   1,
							Name: "Test Portal 1",
						},
					},
					{
						Portal: &entityPortalsV2.Portal{
							ID:   2,
							Name: "Test Portal 2",
						},
					},
				}

				f.portalsRepository.EXPECT().
					Filter(a.ctx, a.filters, a.options).
					Return(expectedPortals, nil)

				return expectedPortals, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				portalsRepository: NewMockPortalsRepository(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			interactor := NewPortalsUseCase(f.portalsRepository)

			got, err := interactor.Filter(tt.args.ctx, tt.args.filters, tt.args.options)

			if wantErr != nil {
				assert.Nil(t, got)
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}
		})
	}
}
