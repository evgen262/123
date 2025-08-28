package portalsv2

import (
	"context"
	"fmt"
	"testing"

	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/portals/v1"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPortalsRepository_Filter(t *testing.T) {
	type fields struct {
		portalsClient *portalsv1.MockPortalsAPIClient
		portalsMapper *MockPortalsMapper // Use the mock generated for this package's mapper interface
	}
	type args struct {
		ctx     context.Context
		filters *entityPortalsV2.FilterPortalsFilters
		options *entityPortalsV2.FilterPortalsOptions
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*entityPortalsV2.PortalWithCounts, error) // Adjusted return type
	}{
		{
			name: "portalsClient error",
			args: args{
				ctx: ctx,
				filters: &entityPortalsV2.FilterPortalsFilters{
					IDs: []int{1, 2},
				},
				options: &entityPortalsV2.FilterPortalsOptions{
					WithEmployeesCount: true,
				},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.PortalWithCounts, error) {
				// Prepare protobuf request parts
				pbFilters := &portalsv1.FilterRequest_Filters{Ids: []int32{1, 2}}
				pbOptions := &portalsv1.FilterRequest_Options{WithEmployeesCount: true}

				// Set up mapper expectations for entity -> pb conversion
				f.portalsMapper.EXPECT().PortalsFiltersToPb(a.filters).Return(pbFilters)
				f.portalsMapper.EXPECT().PortalsOptionsToPb(a.options).Return(pbOptions)

				// Set up gRPC client expectation for error
				f.portalsClient.EXPECT().
					Filter(a.ctx, &portalsv1.FilterRequest{Filters: pbFilters, Options: pbOptions}).
					Return(nil, testErr)

				// Define expected error matching the repository's error wrapping
				return nil, fmt.Errorf("portalsClient.Filter: can't get portals: %w", diterrors.GrpcErrorToError(testErr))
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				filters: &entityPortalsV2.FilterPortalsFilters{
					IDs: []int{1, 2},
				},
				options: &entityPortalsV2.FilterPortalsOptions{
					WithEmployeesCount: true,
				},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.PortalWithCounts, error) {
				// Prepare protobuf request parts
				pbFilters := &portalsv1.FilterRequest_Filters{Ids: []int32{1, 2}}
				pbOptions := &portalsv1.FilterRequest_Options{WithEmployeesCount: true}

				// Prepare protobuf response data (slice of PortalWithCounts)
				pbResponsePortalsWithCounts := []*portalsv1.PortalWithCounts{
					{
						Portal:         &portalsv1.Portal{Id: 1, Name: "Test portal 1"},
						OrgsCount:      1,
						EmployeesCount: 100,
					},
					{
						Portal:         &portalsv1.Portal{Id: 2, Name: "Test portal 2"},
						OrgsCount:      2,
						EmployeesCount: 200,
					},
				}
				pbResponse := &portalsv1.FilterResponse{Portals: pbResponsePortalsWithCounts} // FilterResponse contains []*PortalWithCounts

				// Prepare expected entity output (slice of PortalWithCounts)
				entityResult := []*entityPortalsV2.PortalWithCounts{
					{
						Portal:         &entityPortalsV2.Portal{ID: 1, Name: "Test portal 1"},
						OrgsCount:      1,
						EmployeesCount: 100,
					},
					{
						Portal:         &entityPortalsV2.Portal{ID: 2, Name: "Test portal 2"},
						OrgsCount:      2,
						EmployeesCount: 200,
					},
				}

				// Set up mapper expectations for entity -> pb conversion
				f.portalsMapper.EXPECT().PortalsFiltersToPb(a.filters).Return(pbFilters)
				f.portalsMapper.EXPECT().PortalsOptionsToPb(a.options).Return(pbOptions)

				// Set up gRPC client expectation for successful call
				f.portalsClient.EXPECT().
					Filter(a.ctx, &portalsv1.FilterRequest{Filters: pbFilters, Options: pbOptions}).
					Return(pbResponse, nil)

				// Set up mapper expectation for pb -> entity conversion
				// Mapper expects []*portalsv1.PortalWithCounts and returns []*entityPortalsV2.PortalWithCounts
				f.portalsMapper.EXPECT().PortalsWithCountsToEntity(pbResponsePortalsWithCounts).Return(entityResult)

				// Return expected entity result and no error
				return entityResult, nil
			},
		},
		{
			name: "correct - empty result",
			args: args{
				ctx: ctx,
				filters: &entityPortalsV2.FilterPortalsFilters{
					IDs: []int{999}, // Filter for non-existent IDs
				},
				options: &entityPortalsV2.FilterPortalsOptions{},
			},
			want: func(a args, f fields) ([]*entityPortalsV2.PortalWithCounts, error) {
				// Prepare protobuf request parts
				pbFilters := &portalsv1.FilterRequest_Filters{Ids: []int32{999}}
				pbOptions := &portalsv1.FilterRequest_Options{}

				// Prepare empty protobuf response
				pbResponsePortalsWithCounts := []*portalsv1.PortalWithCounts{}
				pbResponse := &portalsv1.FilterResponse{Portals: pbResponsePortalsWithCounts}

				// Prepare empty expected entity output
				entityResult := []*entityPortalsV2.PortalWithCounts{}

				// Set up mapper expectations for entity -> pb conversion
				f.portalsMapper.EXPECT().PortalsFiltersToPb(a.filters).Return(pbFilters)
				f.portalsMapper.EXPECT().PortalsOptionsToPb(a.options).Return(pbOptions)

				// Set up gRPC client expectation for successful call returning empty list
				f.portalsClient.EXPECT().
					Filter(a.ctx, &portalsv1.FilterRequest{Filters: pbFilters, Options: pbOptions}).
					Return(pbResponse, nil)

				// Set up mapper expectation for pb -> entity conversion (empty list)
				f.portalsMapper.EXPECT().PortalsWithCountsToEntity(pbResponsePortalsWithCounts).Return(entityResult)

				// Return empty entity result and no error
				return entityResult, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Ensure controller is finished

			// Create mocks
			portalsClient := portalsv1.NewMockPortalsAPIClient(ctrl)
			portalsMapper := NewMockPortalsMapper(ctrl) // Use the mock for portalsv2 mapper

			// Populate fields struct with mocks
			f := fields{
				portalsClient: portalsClient,
				portalsMapper: portalsMapper,
			}

			// Call the want function to set up expectations on mocks and get expected results
			want, wantErr := tt.want(tt.args, f)

			// Create the repository instance with mocks
			repo := NewPortalsRepository(f.portalsClient, f.portalsMapper) // Use portalsMapper

			// Act: Call the method under test
			got, err := repo.Filter(tt.args.ctx, tt.args.filters, tt.args.options) // Pass filters and options

			// Assert: Compare actual results with expected results
			if wantErr != nil {
				assert.Nil(t, got) // Expect nil slice on error
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}
		})
	}
}
