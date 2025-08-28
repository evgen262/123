package surveys

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_imagesUseCase_Get(t *testing.T) {
	type fields struct {
		repo *MockImagesRepository
	}
	type args struct {
		ctx       context.Context
		imageName string
	}

	ctx := context.TODO()
	testErr := errors.New("some error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]byte, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:       ctx,
				imageName: "name",
			},
			want: func(a args, f fields) ([]byte, error) {
				f.repo.EXPECT().Get(a.ctx, a.imageName).Return([]byte{}, nil)
				return []byte{}, nil
			},
		},
		{
			name: "err",
			args: args{
				ctx:       ctx,
				imageName: "name",
			},
			want: func(a args, f fields) ([]byte, error) {
				f.repo.EXPECT().Get(a.ctx, a.imageName).Return(nil, testErr)
				return nil, fmt.Errorf("can't get image from repository: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo: NewMockImagesRepository(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			sr := NewImagesUseCase(f.repo)
			got, err := sr.Get(tt.args.ctx, tt.args.imageName)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
