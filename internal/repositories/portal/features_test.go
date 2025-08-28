package portal

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	featuresv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/features/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_portalFeaturesRepository_All(t *testing.T) {
	type fields struct {
		client *featuresv1.MockFeaturesAPIClient
		mapper *MockFeaturesMapper
	}
	type args struct {
		ctx         context.Context
		withDeleted bool
	}
	ctx := context.TODO()
	testT := time.Now()
	time := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Feature, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				withDeleted: false,
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				features := []*featuresv1.Feature{
					{
						Id:          1,
						Name:        "Test Feature",
						Version:     "Test Version",
						CreatedTime: time,
						UpdatedTime: time,
						Enabled:     false,
					},
				}
				portalFeats := []*portal.Feature{
					{
						Id:        1,
						Name:      "Test Feature",
						Version:   "Test Version",
						CreatedAt: &testT,
						UpdatedAt: &testT,
						Enabled:   false,
					},
				}
				f.client.EXPECT().All(a.ctx, &featuresv1.AllRequest{
					WithDeleted: a.withDeleted,
				}).Return(&featuresv1.AllResponse{
					Features: features,
				}, nil)
				f.mapper.EXPECT().FeaturesToEntity(features).Return(portalFeats)
				return portalFeats, nil
			},
		},
		{
			name: "get all features from portal service error",
			args: args{
				ctx:         ctx,
				withDeleted: true,
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				f.client.EXPECT().All(a.ctx, &featuresv1.AllRequest{
					WithDeleted: a.withDeleted,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't get all features: %w", testErr)
			},
		},
		{
			name: "features not found",
			args: args{
				ctx:         ctx,
				withDeleted: true,
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				f.client.EXPECT().All(a.ctx, &featuresv1.AllRequest{
					WithDeleted: a.withDeleted,
				}).Return(nil, diterrors.NewApiError(codes.NotFound, testErr.Error(), testErr))
				return nil, repositories.ErrNotFound
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: featuresv1.NewMockFeaturesAPIClient(ctrl),
				mapper: NewMockFeaturesMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewFeaturesRepository(f.client, f.mapper)
			got, err := repo.All(tt.args.ctx, tt.args.withDeleted)
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

func Test_portalFeaturesRepository_Get(t *testing.T) {
	type fields struct {
		client *featuresv1.MockFeaturesAPIClient
		mapper *MockFeaturesMapper
	}
	type args struct {
		ctx         context.Context
		featureId   portal.FeatureId
		withDeleted bool
	}
	ctx := context.TODO()
	testT := time.Now()
	time := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Feature, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				feature := &featuresv1.Feature{
					Id:          1,
					Name:        "Test Feature",
					Version:     "Test Version",
					CreatedTime: time,
					UpdatedTime: time,
					Enabled:     false,
				}
				entityFeat := &portal.Feature{
					Id:        1,
					Name:      "Test Feature",
					Version:   "Test Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				}
				f.client.EXPECT().Get(a.ctx, &featuresv1.GetRequest{
					Id:          int32(a.featureId),
					WithDeleted: a.withDeleted,
				}).Return(&featuresv1.GetResponse{
					Feature: feature,
				}, nil)
				f.mapper.EXPECT().FeatureToEntity(feature).Return(entityFeat)
				return entityFeat, nil
			},
		},
		{
			name: "get feature from portal service error",
			args: args{
				ctx:         ctx,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				f.client.EXPECT().Get(a.ctx, &featuresv1.GetRequest{
					Id:          int32(a.featureId),
					WithDeleted: a.withDeleted,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't get feature: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: featuresv1.NewMockFeaturesAPIClient(ctrl),
				mapper: NewMockFeaturesMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewFeaturesRepository(f.client, f.mapper)
			got, err := repo.Get(tt.args.ctx, tt.args.featureId, tt.args.withDeleted)
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

func Test_portalFeaturesRepository_Add(t *testing.T) {
	type fields struct {
		client *featuresv1.MockFeaturesAPIClient
		mapper *MockFeaturesMapper
	}
	type args struct {
		ctx      context.Context
		features []*portal.Feature
	}
	ctx := context.TODO()
	testT := time.Now()
	time := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Feature, error)
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
				features: []*portal.Feature{
					{
						Id:        1,
						Name:      "Test Feature",
						Version:   "Test Version",
						CreatedAt: &testT,
						UpdatedAt: &testT,
						Enabled:   false,
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				featuresPb := []*featuresv1.Feature{
					{
						Id:          1,
						Name:        "Test Feature",
						Version:     "Test Version",
						CreatedTime: time,
						UpdatedTime: time,
						Enabled:     false,
					},
				}
				featureToPb := []*featuresv1.AddRequest_Feature{
					{
						Name:    "Test Feature",
						Version: "Test Version",
						Enabled: false,
					},
				}
				entitiesFeats := []*portal.Feature{
					{
						Id:        1,
						Name:      "Test Feature",
						Version:   "Test Version",
						CreatedAt: &testT,
						UpdatedAt: &testT,
						Enabled:   false,
					},
				}
				f.mapper.EXPECT().NewFeaturesToPb(a.features).Return(featureToPb)
				f.client.EXPECT().Add(a.ctx, &featuresv1.AddRequest{
					Features: featureToPb,
				}).Return(&featuresv1.AddResponse{
					Features: featuresPb,
				}, nil)
				f.mapper.EXPECT().FeaturesToEntity(featuresPb).Return(entitiesFeats)
				return entitiesFeats, nil
			},
		},
		{
			name: "add features service error",
			args: args{
				ctx: ctx,
				features: []*portal.Feature{
					{
						Id:        1,
						Name:      "Test Feature",
						Version:   "Test Version",
						CreatedAt: &testT,
						UpdatedAt: &testT,
						Enabled:   false,
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				featureToPb := []*featuresv1.AddRequest_Feature{
					{
						Name:    "Test Feature",
						Version: "Test Version",
						Enabled: false,
					},
				}
				f.mapper.EXPECT().NewFeaturesToPb(a.features).Return(featureToPb)
				f.client.EXPECT().Add(a.ctx, &featuresv1.AddRequest{
					Features: featureToPb,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't add features: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: featuresv1.NewMockFeaturesAPIClient(ctrl),
				mapper: NewMockFeaturesMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewFeaturesRepository(f.client, f.mapper)
			got, err := repo.Add(tt.args.ctx, tt.args.features)
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

func Test_portalFeaturesRepository_Update(t *testing.T) {
	type fields struct {
		client *featuresv1.MockFeaturesAPIClient
		mapper *MockFeaturesMapper
	}
	type args struct {
		ctx     context.Context
		feature *portal.Feature
	}
	ctx := context.TODO()
	testT := time.Now()
	time := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Feature, error)
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
				feature: &portal.Feature{
					Id:        1,
					Name:      "Test Feature",
					Version:   "Test Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				},
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				feature := &featuresv1.Feature{
					Id:          1,
					Name:        "Test Feature",
					Version:     "Test Version",
					CreatedTime: time,
					UpdatedTime: time,
					Enabled:     false,
				}
				entityFeat := &portal.Feature{
					Id:        1,
					Name:      "Test Feature",
					Version:   "Test Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				}
				f.mapper.EXPECT().FeatureToPb(a.feature).Return(feature)
				f.client.EXPECT().Update(a.ctx, &featuresv1.UpdateRequest{
					Feature: feature,
				}).Return(&featuresv1.UpdateResponse{
					Feature: feature,
				}, nil)
				f.mapper.EXPECT().FeatureToEntity(feature).Return(entityFeat)
				return entityFeat, nil
			},
		},
		{
			name: "update feature service error",
			args: args{
				ctx: ctx,
				feature: &portal.Feature{
					Id:        1,
					Name:      "Test Feature",
					Version:   "Test Version",
					CreatedAt: &testT,
					UpdatedAt: &testT,
					Enabled:   false,
				},
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				featureToPb := &featuresv1.Feature{
					Id:          1,
					Name:        "Test Feature",
					Version:     "Test Version",
					CreatedTime: time,
					UpdatedTime: time,
					Enabled:     false,
				}
				f.mapper.EXPECT().FeatureToPb(a.feature).Return(featureToPb)
				f.client.EXPECT().Update(a.ctx, &featuresv1.UpdateRequest{
					Feature: featureToPb,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't update feature: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: featuresv1.NewMockFeaturesAPIClient(ctrl),
				mapper: NewMockFeaturesMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewFeaturesRepository(f.client, f.mapper)
			got, err := repo.Update(tt.args.ctx, tt.args.feature)
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

func Test_portalFeaturesRepository_Delete(t *testing.T) {
	type fields struct {
		client *featuresv1.MockFeaturesAPIClient
		mapper *MockFeaturesMapper
	}
	type args struct {
		ctx       context.Context
		featureId portal.FeatureId
	}
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "correct",
			args: args{
				ctx:       context.TODO(),
				featureId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &featuresv1.DeleteRequest{
					Id: int32(a.featureId),
				}).Return(nil, nil)
				return nil
			},
		},
		{
			name: "delete feature service error",
			args: args{
				ctx:       context.TODO(),
				featureId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &featuresv1.DeleteRequest{
					Id: int32(a.featureId),
				}).Return(nil, testErr)
				return fmt.Errorf("can't delete feature: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: featuresv1.NewMockFeaturesAPIClient(ctrl),
				mapper: NewMockFeaturesMapper(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			repo := NewFeaturesRepository(f.client, f.mapper)
			err := repo.Delete(tt.args.ctx, tt.args.featureId)
			assert.Equal(t, wantErr, err)
		})
	}
}
