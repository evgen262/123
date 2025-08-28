package portals

import (
	"context"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_organizationsUseCase_Filter(t *testing.T) {
	type fields struct {
		repo   *MockOrganizationsRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx        context.Context
		filters    portal.OrganizationsFilters
		pagination *entity.StringPagination
		options    portal.OrganizationsFilterOptions
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.OrganizationsWithPagination, error)
	}{
		{
			name: "repo err",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				f.repo.EXPECT().Filter(a.ctx, a.filters, a.pagination, a.options).Return(nil, testErr)
				f.logger.EXPECT().Error("can't filter organizations", zap.Error(testErr))
				return nil, fmt.Errorf("can't filter organizations: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				result := &portal.OrganizationsWithPagination{}
				f.repo.EXPECT().Filter(a.ctx, a.filters, a.pagination, a.options).Return(result, nil)
				return result, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockOrganizationsRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			ouc := NewOrganizationsUseCase(f.repo, f.logger)
			got, err := ouc.Filter(tt.args.ctx, tt.args.filters, tt.args.pagination, tt.args.options)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_organizationsUseCase_LinkOrganizationsToPortal(t *testing.T) {
	type fields struct {
		repo   *MockOrganizationsRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		portalId portal.PortalID
		orgIds   portal.OrganizationIDs
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "repo link error",
			args: args{
				ctx:      ctx,
				portalId: 1,
				orgIds:   portal.OrganizationIDs{"1", "2"},
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().LinkOrganizationsToPortal(a.ctx, a.portalId, a.orgIds).Return(testErr)
				f.logger.EXPECT().Error(
					"can't link organizations",
					zap.Int("portal_id", int(a.portalId)),
					zap.String("org_ids", a.orgIds.ToString()),
					zap.Error(testErr),
				)
				return fmt.Errorf("can't link organizations: %w", testErr)
			},
		},
		{
			name: "link correct",
			args: args{
				ctx:      ctx,
				portalId: 1,
				orgIds:   portal.OrganizationIDs{"1", "2"},
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().LinkOrganizationsToPortal(a.ctx, a.portalId, a.orgIds).Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockOrganizationsRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			ouc := NewOrganizationsUseCase(f.repo, f.logger)
			gotErr := ouc.Link(tt.args.ctx, tt.args.portalId, tt.args.orgIds)
			if wantErr != nil {
				assert.EqualError(t, gotErr, wantErr.Error())
			} else {
				assert.Nil(t, gotErr)
			}
		})
	}
}

func Test_organizationsUseCase_UnlinkOrganizations(t *testing.T) {
	type fields struct {
		repo   *MockOrganizationsRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx    context.Context
		orgIds portal.OrganizationIDs
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "repo unlink error",
			args: args{
				ctx:    ctx,
				orgIds: portal.OrganizationIDs{"1", "2"},
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().UnlinkOrganizations(a.ctx, a.orgIds).Return(testErr)
				f.logger.EXPECT().Error(
					"can't unlink organizations",
					zap.String("org_ids", a.orgIds.ToString()),
					zap.Error(testErr),
				)
				return fmt.Errorf("can't unlink organizations: %w", testErr)
			},
		},
		{
			name: "unlink correct",
			args: args{
				ctx:    ctx,
				orgIds: portal.OrganizationIDs{"1", "2"},
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().UnlinkOrganizations(a.ctx, a.orgIds).Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockOrganizationsRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			ouc := NewOrganizationsUseCase(f.repo, f.logger)
			gotErr := ouc.Unlink(tt.args.ctx, tt.args.orgIds)
			if wantErr != nil {
				assert.EqualError(t, gotErr, wantErr.Error())
			} else {
				assert.Nil(t, gotErr)
			}
		})
	}
}
