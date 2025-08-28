package survey

import (
	"context"
	"errors"
	"fmt"
	"testing"

	answerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answer/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mapper "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/surveys"
)

func Test_answerRepository_Add(t *testing.T) {
	type fields struct {
		client *answerv1.MockAnswerAPIClient
		mapper *MockAnswerMapper
	}
	type args struct {
		ctx     context.Context
		answers []*surveys.RespondentAnswer
	}

	ctx := context.TODO()
	testErr := errors.New("some error")
	testUUID := uuid.New()
	respondentID := surveys.RespondentID(testUUID)
	localMapper := mapper.NewAnswerMapper()
	testContent := "content"

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]uuid.UUID, error)
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
				answers: []*surveys.RespondentAnswer{
					{
						ChosenVariant: surveys.AnswerID(testUUID),
						Content:       testContent,
						RespondentId:  &respondentID,
					},
				},
			},
			want: func(a args, f fields) ([]uuid.UUID, error) {
				answerIDs := []string{testUUID.String()}
				answersPb := localMapper.NewAnswersToPb(a.answers)
				f.client.EXPECT().Add(a.ctx, &answerv1.AddRequest{
					Answers: answersPb,
					Respondent: &answerv1.Respondent{
						Id: a.answers[0].RespondentId.String(),
					},
				}).Return(&answerv1.AddResponse{
					AnswerIds: answerIDs,
				}, nil)
				f.mapper.EXPECT().NewAnswersToPb(a.answers).Return(answersPb)
				IDs, _ := localMapper.IDsToUUIDs(answerIDs)
				f.mapper.EXPECT().IDsToUUIDs(answerIDs).Return(IDs, nil)
				return IDs, nil
			},
		},
		{
			name: "correct 2",
			args: args{
				ctx: ctx,
				answers: []*surveys.RespondentAnswer{
					{
						ChosenVariant: surveys.AnswerID(testUUID),
					},
				},
			},
			want: func(a args, f fields) ([]uuid.UUID, error) {
				answerIDs := []string{testUUID.String()}
				answersPb := localMapper.NewAnswersToPb(a.answers)
				f.client.EXPECT().Add(a.ctx, &answerv1.AddRequest{
					Answers: answersPb,
				}).Return(&answerv1.AddResponse{
					AnswerIds: answerIDs,
				}, nil)
				f.mapper.EXPECT().NewAnswersToPb(a.answers).Return(answersPb)
				IDs, _ := localMapper.IDsToUUIDs(answerIDs)
				f.mapper.EXPECT().IDsToUUIDs(answerIDs).Return(IDs, nil)
				return IDs, nil
			},
		},
		{
			name: "answer ids to uuids error",
			args: args{
				ctx: ctx,
				answers: []*surveys.RespondentAnswer{
					{
						ChosenVariant: surveys.AnswerID(testUUID),
						Content:       testContent,
						RespondentId:  &respondentID,
					},
				},
			},
			want: func(a args, f fields) ([]uuid.UUID, error) {
				answerIDs := []string{testUUID.String()}
				answersPb := localMapper.NewAnswersToPb(a.answers)
				f.client.EXPECT().Add(a.ctx, &answerv1.AddRequest{
					Answers: answersPb,
					Respondent: &answerv1.Respondent{
						Id: a.answers[0].RespondentId.String(),
					},
				}).Return(&answerv1.AddResponse{
					AnswerIds: answerIDs,
				}, nil)
				f.mapper.EXPECT().NewAnswersToPb(a.answers).Return(answersPb)
				f.mapper.EXPECT().IDsToUUIDs(answerIDs).Return(nil, testErr)
				return nil, fmt.Errorf("can't convert answer ids to uuids: %w", testErr)
			},
		},
		{
			name: "add answers invalid argument error",
			args: args{
				ctx: ctx,
				answers: []*surveys.RespondentAnswer{
					{
						ChosenVariant: surveys.AnswerID(testUUID),
						Content:       testContent,
						RespondentId:  &respondentID,
					},
				},
			},
			want: func(a args, f fields) ([]uuid.UUID, error) {
				answersPb := localMapper.NewAnswersToPb(a.answers)
				f.client.EXPECT().Add(a.ctx, &answerv1.AddRequest{
					Answers: answersPb,
					Respondent: &answerv1.Respondent{
						Id: a.answers[0].RespondentId.String(),
					},
				}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				f.mapper.EXPECT().NewAnswersToPb(a.answers).Return(answersPb)
				return nil, diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "add answers default error",
			args: args{
				ctx: ctx,
				answers: []*surveys.RespondentAnswer{
					{
						ChosenVariant: surveys.AnswerID(testUUID),
						Content:       testContent,
						RespondentId:  &respondentID,
					},
				},
			},
			want: func(a args, f fields) ([]uuid.UUID, error) {
				answersPb := localMapper.NewAnswersToPb(a.answers)
				f.client.EXPECT().Add(a.ctx, &answerv1.AddRequest{
					Answers: answersPb,
					Respondent: &answerv1.Respondent{
						Id: a.answers[0].RespondentId.String(),
					},
				}).Return(nil, status.Error(codes.Internal, testErr.Error()))
				f.mapper.EXPECT().NewAnswersToPb(a.answers).Return(answersPb)
				return nil, fmt.Errorf("can't add answers: %s", testErr.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: answerv1.NewMockAnswerAPIClient(ctrl),
				mapper: NewMockAnswerMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			sr := NewAnswerRepository(f.client, f.mapper)
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
