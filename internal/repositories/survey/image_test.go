package survey

import (
	"context"
	"errors"
	"fmt"
	"testing"

	imagev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/image/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_imageRepository_Get(t *testing.T) {
	type fields struct {
		client *imagev1.MockImageAPIClient
		mapper *MockImageMapper
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
				imageName: "test",
			},
			want: func(a args, f fields) ([]byte, error) {
				f.client.EXPECT().Get(a.ctx, &imagev1.GetRequest{
					Id: a.imageName,
				}).Return(&imagev1.GetResponse{
					Image: []byte("test"),
				}, nil)
				return []byte("test"), nil
			},
		},
		{
			name: "default err",
			args: args{
				ctx:       ctx,
				imageName: "test",
			},
			want: func(a args, f fields) ([]byte, error) {
				f.client.EXPECT().Get(a.ctx, &imagev1.GetRequest{
					Id: a.imageName,
				}).Return(nil, status.Error(codes.Internal, testErr.Error()))
				return nil, fmt.Errorf("can't get image: %s", testErr.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: imagev1.NewMockImageAPIClient(ctrl),
				mapper: NewMockImageMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			sr := NewImageRepository(f.client, f.mapper)
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
