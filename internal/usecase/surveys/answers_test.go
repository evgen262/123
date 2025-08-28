package surveys

import (
	"context"
	"errors"
	"fmt"
	"testing"

	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_answersUseCase_Add(t *testing.T) {
	type fields struct {
		repo *MockAnswersRepository
	}
	type args struct {
		ctx     context.Context
		answers []*surveys.RespondentAnswer
	}

	ctx := context.TODO()
	testErr := errors.New("some error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]uuid.UUID, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				answers: []*surveys.RespondentAnswer{},
			},
			want: func(a args, f fields) ([]uuid.UUID, error) {
				f.repo.EXPECT().Add(a.ctx, a.answers).Return([]uuid.UUID{}, nil)
				return []uuid.UUID{}, nil
			},
		},
		{
			name: "err",
			args: args{
				ctx:     ctx,
				answers: []*surveys.RespondentAnswer{},
			},
			want: func(a args, f fields) ([]uuid.UUID, error) {
				f.repo.EXPECT().Add(a.ctx, a.answers).Return(nil, testErr)
				return nil, fmt.Errorf("can't add answers to repository: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo: NewMockAnswersRepository(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			sr := NewAnswersUseCase(f.repo)
			got, err := sr.Add(tt.args.ctx, tt.args.answers)
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
