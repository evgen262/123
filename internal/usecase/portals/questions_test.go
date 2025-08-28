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

func Test_questionsUseCase_AddQuestion(t *testing.T) {
	type fields struct {
		repo   *MockQuestionsRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		question *portal.Question
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Question, error)
	}{
		{
			name: "err",
			args: args{
				ctx:      ctx,
				question: &portal.Question{Name: "test"},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.repo.EXPECT().Add(a.ctx, []*portal.Question{a.question}).Return(nil, testErr)
				f.logger.EXPECT().Error("can't add question", zap.Error(testErr))
				return nil, testErr
			},
		},
		{
			name: "correct with empty",
			args: args{
				ctx:      ctx,
				question: &portal.Question{Name: "test"},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.repo.EXPECT().Add(a.ctx, []*portal.Question{a.question}).Return([]*portal.Question{}, nil)
				return nil, nil
			},
		},
		{
			name: "correct",
			args: args{
				ctx:      ctx,
				question: &portal.Question{Name: "test"},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.repo.EXPECT().Add(a.ctx, []*portal.Question{a.question}).Return([]*portal.Question{{Id: 1}}, nil)
				return &portal.Question{Id: 1}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockQuestionsRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			quc := NewQuestionUseCase(f.repo, f.logger)
			got, err := quc.AddQuestion(tt.args.ctx, tt.args.question)
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

func Test_questionsUseCase_AddQuestions(t *testing.T) {
	type fields struct {
		repo   *MockQuestionsRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx       context.Context
		questions []*portal.Question
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Question, error)
	}{
		{
			name: "err",
			args: args{
				ctx:       ctx,
				questions: []*portal.Question{{Name: "test"}},
			},
			want: func(a args, f fields) ([]*portal.Question, error) {
				f.repo.EXPECT().Add(a.ctx, a.questions).Return(nil, testErr)
				f.logger.EXPECT().Error("can't add questions", zap.Error(testErr))
				return nil, testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx:       ctx,
				questions: []*portal.Question{{Name: "test"}},
			},
			want: func(a args, f fields) ([]*portal.Question, error) {
				q := []*portal.Question{{Id: 1}, {Id: 2}}
				f.repo.EXPECT().Add(a.ctx, a.questions).Return(q, nil)
				return q, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockQuestionsRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			quc := NewQuestionUseCase(f.repo, f.logger)
			got, err := quc.AddQuestions(tt.args.ctx, tt.args.questions)
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

func Test_questionsUseCase_DeleteQuestion(t *testing.T) {
	type fields struct {
		repo   *MockQuestionsRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx        context.Context
		questionId int
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
				ctx:        ctx,
				questionId: 1,
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Delete(a.ctx, portal.QuestionId(a.questionId)).Return(testErr)
				f.logger.EXPECT().Error("can't delete question", zap.Int("question_id", a.questionId), zap.Error(testErr))
				return testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx:        ctx,
				questionId: 1,
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Delete(a.ctx, portal.QuestionId(a.questionId)).Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockQuestionsRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			quc := NewQuestionUseCase(f.repo, f.logger)
			err := quc.DeleteQuestion(tt.args.ctx, tt.args.questionId)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_questionsUseCase_GetAllQuestions(t *testing.T) {
	type fields struct {
		repo   *MockQuestionsRepository
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
		want func(a args, f fields) (*portal.Questions, error)
	}{
		{
			name: "err",
			args: args{
				ctx:          ctx,
				withDisabled: true,
			},
			want: func(a args, f fields) (*portal.Questions, error) {
				f.repo.EXPECT().All(a.ctx, a.withDisabled).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get questions", zap.Error(testErr))
				return nil, testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx:          ctx,
				withDisabled: true,
			},
			want: func(a args, f fields) (*portal.Questions, error) {
				questions := []*portal.Question{
					{Name: "test"},
				}
				f.repo.EXPECT().All(a.ctx, a.withDisabled).
					Return(&portal.Questions{
						SupportEmail: "a@a.ru",
						Questions:    questions,
					}, nil)
				return &portal.Questions{
					SupportEmail: "a@a.ru",
					Questions:    questions,
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockQuestionsRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			quc := NewQuestionUseCase(f.repo, f.logger)
			got, err := quc.GetAllQuestions(tt.args.ctx, tt.args.withDisabled)
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

func Test_questionsUseCase_GetQuestion(t *testing.T) {
	type fields struct {
		repo   *MockQuestionsRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx          context.Context
		questionId   int
		withDisabled bool
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Question, error)
	}{
		{
			name: "err",
			args: args{
				ctx:          ctx,
				questionId:   1,
				withDisabled: true,
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.repo.EXPECT().Get(a.ctx, portal.QuestionId(a.questionId), a.withDisabled).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get question", zap.Error(testErr))
				return nil, testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx:          ctx,
				questionId:   1,
				withDisabled: true,
			},
			want: func(a args, f fields) (*portal.Question, error) {
				question := &portal.Question{
					Name: "test",
				}
				f.repo.EXPECT().Get(a.ctx, portal.QuestionId(a.questionId), a.withDisabled).Return(question, nil)
				return question, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockQuestionsRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			quc := NewQuestionUseCase(f.repo, f.logger)
			got, err := quc.GetQuestion(tt.args.ctx, tt.args.questionId, tt.args.withDisabled)
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

func Test_questionsUseCase_UpdateQuestion(t *testing.T) {
	type fields struct {
		repo   *MockQuestionsRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		question *portal.Question
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Question, error)
	}{
		{
			name: "err",
			args: args{
				ctx: ctx,
				question: &portal.Question{
					Name: "test",
				},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.repo.EXPECT().Update(a.ctx, a.question).Return(nil, testErr)
				f.logger.EXPECT().Error("can't update question", zap.Int("question_id", int(a.question.Id)), zap.Error(testErr))
				return nil, testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				question: &portal.Question{
					Name: "test",
				},
			},
			want: func(a args, f fields) (*portal.Question, error) {
				f.repo.EXPECT().Update(a.ctx, a.question).Return(&portal.Question{Name: "test"}, nil)
				return &portal.Question{Name: "test"}, nil
			},
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		f := fields{
			repo:   NewMockQuestionsRepository(ctrl),
			logger: ditzap.NewMockLogger(ctrl),
		}
		want, wantErr := tt.want(tt.args, f)
		quc := NewQuestionUseCase(f.repo, f.logger)
		got, err := quc.UpdateQuestion(tt.args.ctx, tt.args.question)
		if wantErr != nil {
			assert.EqualError(t, err, wantErr.Error())
			assert.Nil(t, got)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}
}
