package portalsv2

import (
	"context"
	"fmt"
	"testing"

	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCompexesInteractor_Filter(t *testing.T) {
	type fields struct {
		complexesRepository *MockComplexesRepository
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
			name: "repository error",
			args: args{
				ctx:     ctx,
				filters: &entityPortalsV2.FilterComplexesFilters{},
				options: &entityPortalsV2.FilterComplexesOptions{},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.Complex, error) {
				f.complexesRepository.EXPECT().
					Filter(a.ctx, a.filters, a.options).
					Return(nil, testErr)

				return nil, fmt.Errorf("complexesRepository.Filter: can't filter complexes: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				filters: &entityPortalsV2.FilterComplexesFilters{},
				options: &entityPortalsV2.FilterComplexesOptions{},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.Complex, error) {
				expectedComplexes := []*entityPortalsV2.Complex{
					{ID: 3, Name: "Complex C"},
					{ID: 4, Name: "Complex D"},
				}

				f.complexesRepository.EXPECT().
					Filter(a.ctx, a.filters, a.options).
					Return(expectedComplexes, nil)

				return expectedComplexes, nil
			},
		},
		{
			name: "correct with nil filters and options",
			args: args{
				ctx:     ctx,
				filters: nil,
				options: nil,
			},
			want: func(a args, f fields) ([]*entityPortalsV2.Complex, error) {
				expectedComplexes := []*entityPortalsV2.Complex{
					{ID: 1, Name: "Complex A"},
					{ID: 2, Name: "Complex B"},
				}

				f.complexesRepository.EXPECT().
					Filter(a.ctx, a.filters, a.options).
					Return(expectedComplexes, nil)

				return expectedComplexes, nil
			},
		},
		{
			name: "correct with empty result",
			args: args{
				ctx:     ctx,
				filters: &entityPortalsV2.FilterComplexesFilters{},
				options: &entityPortalsV2.FilterComplexesOptions{},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.Complex, error) {
				expectedComplexes := []*entityPortalsV2.Complex{}
				f.complexesRepository.EXPECT().
					Filter(a.ctx, a.filters, a.options).
					Return(expectedComplexes, nil)

				return expectedComplexes, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockRepo := NewMockComplexesRepository(ctrl)

			f := fields{
				complexesRepository: mockRepo,
			}

			want, wantErr := tt.want(tt.args, f)

			interactor := NewComplexesUseCase(f.complexesRepository)

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
