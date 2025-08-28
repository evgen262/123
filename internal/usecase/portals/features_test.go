package portals

import (
	"context"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_featuresUseCase_Add(t *testing.T) {
	type fields struct {
		repo   *MockFeaturesRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx     context.Context
		feature *portal.Feature
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Feature, error)
	}{
		{
			name: "err",
			args: args{
				ctx:     ctx,
				feature: &portal.Feature{Name: "test"},
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				f.repo.EXPECT().Add(a.ctx, []*portal.Feature{a.feature}).Return(nil, testErr)
				f.logger.EXPECT().Error("can't add feature into repo", zap.Error(testErr))
				return nil, fmt.Errorf("can't add feature: %w", testErr)
			},
		},
		{
			name: "correct with empty",
			args: args{
				ctx:     ctx,
				feature: &portal.Feature{Name: "test"},
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				f.repo.EXPECT().Add(a.ctx, []*portal.Feature{a.feature}).Return([]*portal.Feature{}, nil)
				return nil, nil
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				feature: &portal.Feature{Name: "test"},
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				feature := &portal.Feature{Id: 1}
				f.repo.EXPECT().Add(a.ctx, []*portal.Feature{a.feature}).Return([]*portal.Feature{feature}, nil)
				return feature, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockFeaturesRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			fuc := NewFeatureUseCase(f.repo, f.logger)
			got, err := fuc.Add(tt.args.ctx, tt.args.feature)
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

func Test_featuresUseCase_MultipleAdd(t *testing.T) {
	type fields struct {
		repo   *MockFeaturesRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		features []*portal.Feature
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Feature, error)
	}{
		{
			name: "err",
			args: args{
				ctx:      ctx,
				features: []*portal.Feature{{Name: "test"}},
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				f.repo.EXPECT().Add(a.ctx, a.features).Return(nil, testErr)
				f.logger.EXPECT().Error("can't add features into repo", zap.Error(testErr))
				return nil, fmt.Errorf("can't add features: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:      ctx,
				features: []*portal.Feature{{Name: "test"}},
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				features := []*portal.Feature{
					{Id: 1},
					{Id: 2},
				}
				f.repo.EXPECT().Add(a.ctx, a.features).Return(features, nil)
				return features, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockFeaturesRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			fuc := NewFeatureUseCase(f.repo, f.logger)
			got, err := fuc.MultipleAdd(tt.args.ctx, tt.args.features)
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

func Test_featuresUseCase_Delete(t *testing.T) {
	type fields struct {
		repo   *MockFeaturesRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx       context.Context
		featureId int
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "err",
			args: args{
				ctx:       ctx,
				featureId: 1,
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Delete(a.ctx, portal.FeatureId(a.featureId)).Return(testErr)
				f.logger.EXPECT().Error("can't delete feature from repo", zap.Int("feature_id", a.featureId), zap.Error(testErr))
				return fmt.Errorf("can't delete feature: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:       ctx,
				featureId: 1,
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Delete(a.ctx, portal.FeatureId(a.featureId)).Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockFeaturesRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			fuc := NewFeatureUseCase(f.repo, f.logger)
			err := fuc.Delete(tt.args.ctx, tt.args.featureId)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_featuresUseCase_GetAllFeatures(t *testing.T) {
	type fields struct {
		repo   *MockFeaturesRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx          context.Context
		withDisabled bool
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Feature, error)
	}{
		{
			name: "err",
			args: args{
				ctx:          ctx,
				withDisabled: true,
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				f.repo.EXPECT().All(a.ctx, a.withDisabled).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get all features from repo", zap.Error(testErr))
				return nil, fmt.Errorf("can't get all features: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:          ctx,
				withDisabled: true,
			},
			want: func(a args, f fields) ([]*portal.Feature, error) {
				features := []*portal.Feature{
					{Name: "test"},
				}
				f.repo.EXPECT().All(a.ctx, a.withDisabled).Return(features, nil)
				return features, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockFeaturesRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			fuc := NewFeatureUseCase(f.repo, f.logger)
			got, err := fuc.All(tt.args.ctx, tt.args.withDisabled)
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

func Test_featuresUseCase_GetFeature(t *testing.T) {
	type fields struct {
		repo   *MockFeaturesRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx          context.Context
		featureId    int
		withDisabled bool
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Feature, error)
	}{
		{
			name: "err",
			args: args{
				ctx:          ctx,
				featureId:    1,
				withDisabled: true,
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				f.repo.EXPECT().Get(a.ctx, portal.FeatureId(a.featureId), a.withDisabled).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get feature from repo", zap.Error(testErr))
				return nil, fmt.Errorf("can't get feature: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:          ctx,
				featureId:    1,
				withDisabled: true,
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				feature := &portal.Feature{
					Name: "test",
				}
				f.repo.EXPECT().Get(a.ctx, portal.FeatureId(a.featureId), a.withDisabled).Return(feature, nil)
				return feature, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockFeaturesRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			fuc := NewFeatureUseCase(f.repo, f.logger)
			got, err := fuc.Get(tt.args.ctx, tt.args.featureId, tt.args.withDisabled)
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

func Test_featuresUseCase_Update(t *testing.T) {
	type fields struct {
		repo   *MockFeaturesRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx     context.Context
		feature *portal.Feature
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Feature, error)
	}{
		{
			name: "err",
			args: args{
				ctx: ctx,
				feature: &portal.Feature{
					Name: "test",
				},
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				f.repo.EXPECT().Update(a.ctx, a.feature).Return(nil, testErr)
				f.logger.EXPECT().Error(
					"can't update feature into repo",
					zap.String("feature_name", a.feature.Name),
					zap.Error(testErr),
				)
				return nil, fmt.Errorf("can't update feature: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				feature: &portal.Feature{
					Name: "test",
				},
			},
			want: func(a args, f fields) (*portal.Feature, error) {
				f.repo.EXPECT().Update(a.ctx, a.feature).Return(&portal.Feature{Name: "test"}, nil)
				return &portal.Feature{Name: "test"}, nil
			},
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		f := fields{
			repo:   NewMockFeaturesRepository(ctrl),
			logger: ditzap.NewMockLogger(ctrl),
		}
		want, wantErr := tt.want(tt.args, f)
		fuc := NewFeatureUseCase(f.repo, f.logger)
		got, err := fuc.Update(tt.args.ctx, tt.args.feature)
		if wantErr != nil {
			assert.EqualError(t, err, wantErr.Error())
			assert.Nil(t, got)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}
}
