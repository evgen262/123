package portal

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	questionsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/questions/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_portalQuestionsRepository_All(t *testing.T) {
	type fields struct {
		client *questionsv1.MockQuestionsAPIClient
		mapper *MockQuestionsMapper
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
		want func(a args, f fields) (*portal.Questions, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Questions, error) {
				questionsPb := []*questionsv1.Question{
					{
						Id:          1,
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						DeletedTime: time,
						IsDeleted:   false,
						CreatedTime: time,
						UpdatedTime: time,
					},
				}
				entityQuestions := []*portal.Question{
					{
						Id:          1,
						Name:        "Name Test",
						Description: "Desciption Test",
						Sort:        1,
						DeletedAt:   &testT,
						IsDeleted:   false,
						CreatedAt:   &testT,
						UpdatedAt:   &testT,
					},
				}

				f.client.EXPECT().All(a.ctx, &questionsv1.AllRequest{
					WithDeleted: a.withDeleted,
				}).Return(&questionsv1.AllResponse{
					SupportEmail: "SupportEmail Test",
					Questions:    questionsPb,
				}, nil)
				f.mapper.EXPECT().QuestionsToEntity(questionsPb).Return(entityQuestions)

				return &portal.Questions{
					SupportEmail: "SupportEmail Test",
					Questions:    entityQuestions,
				}, nil
			},
		},
		{
			name: "get all questions from portal service error",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) (*portal.Questions, error) {
				f.client.EXPECT().All(a.ctx, &questionsv1.AllRequest{
					WithDeleted: a.withDeleted,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't get all questions: %w", testErr)
			},
		},
		{
			name: "get all question invalid argument",
			args: args{
				ctx:         ctx,
				withDeleted: true,
			},
			want: func(a args, f fields) (*portal.Questions, error) {
				f.client.EXPECT().All(a.ctx, &questionsv1.AllRequest{
					WithDeleted: a.withDeleted,
				}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "get all question not found",
			args: args{
				ctx:         ctx,
				withDeleted: true,
			},
			want: func(a args, f fields) (*portal.Questions, error) {
				f.client.EXPECT().All(a.ctx, &questionsv1.AllRequest{
					WithDeleted: a.withDeleted,
				}).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return nil, repositories.ErrNotFound
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: questionsv1.NewMockQuestionsAPIClient(ctrl),
				mapper: NewMockQuestionsMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewQuestionsRepository(f.client, f.mapper)
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

func Test_portalQuestionsRepository_Get(t *testing.T) {
	type fields struct {
		client *questionsv1.MockQuestionsAPIClient
		mapper *MockQuestionsMapper
	}
	type args struct {
		ctx         context.Context
		questionId  portal.QuestionId
		withDeleted bool
	}
	ctx := context.TODO()
	testT := time.Now()
	time := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Question, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				questionId:  1,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Question, error) {
				question := &questionsv1.Question{
					Id:          1,
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					DeletedTime: time,
					IsDeleted:   false,
					CreatedTime: time,
					UpdatedTime: time,
				}
				entityQuestion := &portal.Question{
					Id:          1,
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					DeletedAt:   &testT,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					IsDeleted:   false,
				}

				f.client.EXPECT().Get(a.ctx, &questionsv1.GetRequest{
					Id:          int32(a.questionId),
					WithDeleted: a.withDeleted,
				}).Return(&questionsv1.GetResponse{
					SupportEmail: "SupportEmail Test",
					Question:     question,
				}, nil)
				f.mapper.EXPECT().QuestionToEntity(question).Return(entityQuestion)
				return entityQuestion, nil
			},
		},
		{
			name: "get question from portal service error",
			args: args{
				ctx:         ctx,
				questionId:  1,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.client.EXPECT().Get(a.ctx, &questionsv1.GetRequest{
					Id:          int32(a.questionId),
					WithDeleted: a.withDeleted,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't get question: %w", testErr)
			},
		},
		{
			name: "get question invalid argument",
			args: args{
				ctx:         ctx,
				questionId:  1,
				withDeleted: true,
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.client.EXPECT().Get(a.ctx, &questionsv1.GetRequest{
					Id:          int32(a.questionId),
					WithDeleted: a.withDeleted,
				}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "get question not found",
			args: args{
				ctx:         ctx,
				questionId:  1,
				withDeleted: true,
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.client.EXPECT().Get(a.ctx, &questionsv1.GetRequest{
					Id:          int32(a.questionId),
					WithDeleted: a.withDeleted,
				}).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return nil, repositories.ErrNotFound
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: questionsv1.NewMockQuestionsAPIClient(ctrl),
				mapper: NewMockQuestionsMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewQuestionsRepository(f.client, f.mapper)
			got, err := repo.Get(tt.args.ctx, tt.args.questionId, tt.args.withDeleted)
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

func Test_portalQuestionsRepository_Add(t *testing.T) {
	type fields struct {
		client *questionsv1.MockQuestionsAPIClient
		mapper *MockQuestionsMapper
	}
	type args struct {
		ctx       context.Context
		questions []*portal.Question
	}
	ctx := context.TODO()
	testT := time.Now()
	time := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Question, error)
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
				questions: []*portal.Question{
					{
						Id:          1,
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						CreatedAt:   &testT,
						UpdatedAt:   &testT,
						DeletedAt:   &testT,
						IsDeleted:   false,
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Question, error) {
				getQuestion := []*questionsv1.Question{
					{
						Id:          1,
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						DeletedTime: time,
						IsDeleted:   false,
						CreatedTime: time,
						UpdatedTime: time,
					},
				}
				questionAddRequest := []*questionsv1.AddRequest_Question{
					{
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						IsDeleted:   false,
					},
				}
				entityQuestion := []*portal.Question{
					{
						Id:          1,
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						DeletedAt:   &testT,
						CreatedAt:   &testT,
						UpdatedAt:   &testT,
						IsDeleted:   false,
					},
				}

				f.mapper.EXPECT().NewQuestionsToPb(a.questions).Return(questionAddRequest)
				f.client.EXPECT().Add(a.ctx, &questionsv1.AddRequest{
					Questions: questionAddRequest,
				}).Return(&questionsv1.AddResponse{
					Questions: getQuestion,
				}, nil)
				f.mapper.EXPECT().QuestionsToEntity(getQuestion).Return(entityQuestion)
				return entityQuestion, nil
			},
		},
		{
			name: "add questions service error",
			args: args{
				ctx: ctx,
				questions: []*portal.Question{
					{
						Id:          1,
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						CreatedAt:   &testT,
						UpdatedAt:   &testT,
						DeletedAt:   &testT,
						IsDeleted:   false,
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Question, error) {
				questionAddRequest := []*questionsv1.AddRequest_Question{
					{
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						IsDeleted:   false,
					},
				}
				f.mapper.EXPECT().NewQuestionsToPb(a.questions).Return(questionAddRequest)
				f.client.EXPECT().Add(a.ctx, &questionsv1.AddRequest{
					Questions: questionAddRequest,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't add questions: %w", testErr)
			},
		},
		{
			name: "add question invalid argument",
			args: args{
				ctx: ctx,
				questions: []*portal.Question{
					{
						Id:          1,
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						CreatedAt:   &testT,
						UpdatedAt:   &testT,
						DeletedAt:   &testT,
						IsDeleted:   false,
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Question, error) {
				questionAddRequest := []*questionsv1.AddRequest_Question{
					{
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						IsDeleted:   false,
					},
				}
				f.mapper.EXPECT().NewQuestionsToPb(a.questions).Return(questionAddRequest)
				f.client.EXPECT().Add(a.ctx, &questionsv1.AddRequest{
					Questions: questionAddRequest,
				}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "add question not found",
			args: args{
				ctx: ctx,
				questions: []*portal.Question{
					{
						Id:          1,
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						CreatedAt:   &testT,
						UpdatedAt:   &testT,
						DeletedAt:   &testT,
						IsDeleted:   false,
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Question, error) {
				questionAddRequest := []*questionsv1.AddRequest_Question{
					{
						Name:        "Name Test",
						Description: "Description Test",
						Sort:        1,
						IsDeleted:   false,
					},
				}
				f.mapper.EXPECT().NewQuestionsToPb(a.questions).Return(questionAddRequest)
				f.client.EXPECT().Add(a.ctx, &questionsv1.AddRequest{
					Questions: questionAddRequest,
				}).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return nil, repositories.ErrNotFound
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: questionsv1.NewMockQuestionsAPIClient(ctrl),
				mapper: NewMockQuestionsMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewQuestionsRepository(f.client, f.mapper)
			got, err := repo.Add(tt.args.ctx, tt.args.questions)
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

func Test_portalQuestionsRepository_Update(t *testing.T) {
	type fields struct {
		client *questionsv1.MockQuestionsAPIClient
		mapper *MockQuestionsMapper
	}
	type args struct {
		ctx      context.Context
		question *portal.Question
	}
	ctx := context.TODO()
	testT := time.Now()
	time := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Question, error)
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
				question: &portal.Question{
					Id:          1,
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					DeletedAt:   &testT,
					IsDeleted:   false,
				},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				getQuestion := &questionsv1.Question{
					Id:          1,
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					DeletedTime: time,
					IsDeleted:   false,
					CreatedTime: time,
					UpdatedTime: time,
				}
				entityQuestion := &portal.Question{
					Id:          1,
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					DeletedAt:   &testT,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					IsDeleted:   false,
				}

				f.mapper.EXPECT().QuestionToPb(a.question).Return(getQuestion)
				f.client.EXPECT().Update(a.ctx, &questionsv1.UpdateRequest{
					Question: getQuestion,
				}).Return(&questionsv1.UpdateResponse{
					Question: getQuestion,
				}, nil)
				f.mapper.EXPECT().QuestionToEntity(getQuestion).Return(entityQuestion)
				return entityQuestion, nil
			},
		},
		{
			name: "update question service error",
			args: args{
				ctx: ctx,
				question: &portal.Question{
					Id:          1,
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					DeletedAt:   &testT,
					IsDeleted:   false,
				},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				questionUpdateRequest := &questionsv1.Question{
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					IsDeleted:   false,
					Id:          1,
					CreatedTime: time,
					UpdatedTime: time,
					DeletedTime: time,
				}
				f.mapper.EXPECT().QuestionToPb(a.question).Return(questionUpdateRequest)
				f.client.EXPECT().Update(a.ctx, &questionsv1.UpdateRequest{
					Question: questionUpdateRequest,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't update question: %w", testErr)
			},
		},
		{
			name: "update question invalid argument",
			args: args{
				ctx: ctx,
				question: &portal.Question{
					Id:          1,
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					DeletedAt:   &testT,
					IsDeleted:   false,
				},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				questionUpdateRequest := &questionsv1.Question{
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					IsDeleted:   false,
					Id:          1,
					CreatedTime: time,
					UpdatedTime: time,
					DeletedTime: time,
				}
				f.mapper.EXPECT().QuestionToPb(a.question).Return(questionUpdateRequest)
				f.client.EXPECT().Update(a.ctx, &questionsv1.UpdateRequest{
					Question: questionUpdateRequest,
				}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "update question not found",
			args: args{
				ctx: ctx,
				question: &portal.Question{
					Id:          1,
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					DeletedAt:   &testT,
					IsDeleted:   false,
				},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				questionUpdateRequest := &questionsv1.Question{
					Name:        "Name Test",
					Description: "Description Test",
					Sort:        1,
					IsDeleted:   false,
					Id:          1,
					CreatedTime: time,
					UpdatedTime: time,
					DeletedTime: time,
				}
				f.mapper.EXPECT().QuestionToPb(a.question).Return(questionUpdateRequest)
				f.client.EXPECT().Update(a.ctx, &questionsv1.UpdateRequest{
					Question: questionUpdateRequest,
				}).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return nil, repositories.ErrNotFound
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: questionsv1.NewMockQuestionsAPIClient(ctrl),
				mapper: NewMockQuestionsMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewQuestionsRepository(f.client, f.mapper)
			got, err := repo.Update(tt.args.ctx, tt.args.question)
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

func Test_portalQuestionsRepository_Delete(t *testing.T) {
	type fields struct {
		client *questionsv1.MockQuestionsAPIClient
		mapper *MockQuestionsMapper
	}
	type args struct {
		ctx        context.Context
		questionId portal.QuestionId
	}
	ctx := context.TODO()
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "correct",
			args: args{
				ctx:        ctx,
				questionId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &questionsv1.DeleteRequest{
					Id: int32(a.questionId),
				}).Return(nil, nil)
				return nil
			},
		},
		{
			name: "delete question service error",
			args: args{
				ctx:        ctx,
				questionId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &questionsv1.DeleteRequest{
					Id: int32(a.questionId),
				}).Return(nil, testErr)
				return fmt.Errorf("can't delete question: %w", diterrors.NewLocalizedError(diterrors.LocalizeLocale, testErr))
			},
		},
		{
			name: "delete question invalid argument",
			args: args{
				ctx:        ctx,
				questionId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &questionsv1.DeleteRequest{
					Id: int32(a.questionId),
				}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return diterrors.NewValidationError(diterrors.NewLocalizedError(diterrors.LocalizeLocale, status.Error(codes.InvalidArgument, testErr.Error())))
			},
		},
		{
			name: "delete question not found",
			args: args{
				ctx:        ctx,
				questionId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &questionsv1.DeleteRequest{
					Id: int32(a.questionId),
				}).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return repositories.ErrNotFound
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: questionsv1.NewMockQuestionsAPIClient(ctrl),
				mapper: NewMockQuestionsMapper(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			repo := NewQuestionsRepository(f.client, f.mapper)
			err := repo.Delete(tt.args.ctx, tt.args.questionId)
			assert.Equal(t, wantErr, err)
		})
	}
}
