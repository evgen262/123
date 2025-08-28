package portals

import (
	"testing"
	"time"

	questionsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/questions/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

func TestPortalsMapper_QuestionToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		request *questionsv1.Question
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)
	time := timestamppb.New(testT)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portal.Question
	}{
		{
			name: "correct",
			args: args{
				request: &questionsv1.Question{
					Id:          1,
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					CreatedTime: time,
					UpdatedTime: time,
					DeletedTime: time,
					IsDeleted:   false,
				},
			},
			want: func(a args, f fields) *portal.Question {
				f.timeUtils.EXPECT().TimestampToTime(a.request.CreatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.request.UpdatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.request.DeletedTime).Return(&testT)

				return &portal.Question{
					Id:          1,
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					DeletedAt:   &testT,
					IsDeleted:   false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewQuestionsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.QuestionToEntity(tt.args.request)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_QuestionToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		request *portal.Question
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *questionsv1.Question
	}{
		{
			name: "correct",
			args: args{
				request: &portal.Question{
					Id:          1,
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					DeletedAt:   &testT,
					IsDeleted:   false,
				},
			},
			want: func(a args, f fields) *questionsv1.Question {
				t := timestamppb.New(testT)
				f.timeUtils.EXPECT().TimeToTimestamp(a.request.CreatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.request.UpdatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.request.DeletedAt).Return(t)

				return &questionsv1.Question{
					Id:          1,
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					CreatedTime: t,
					UpdatedTime: t,
					DeletedTime: t,
					IsDeleted:   false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewQuestionsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.QuestionToPb(tt.args.request)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_NewQuestionToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		request *portal.Question
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *questionsv1.AddRequest_Question
	}{
		{
			name: "correct",
			args: args{
				request: &portal.Question{
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					IsDeleted:   false,
				},
			},
			want: func(a args, f fields) *questionsv1.AddRequest_Question {
				return &questionsv1.AddRequest_Question{
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					IsDeleted:   false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewQuestionsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.NewQuestionToPb(tt.args.request)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_NewQuestionsToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		request []*portal.Question
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*questionsv1.AddRequest_Question
	}{
		{
			name: "correct",
			args: args{
				request: []*portal.Question{
					{
						Id:          1,
						Name:        "Name 1",
						Description: "Description 1",
						Sort:        1,
						CreatedAt:   &testT,
						UpdatedAt:   &testT,
						DeletedAt:   &testT,
						IsDeleted:   false,
					},
				},
			},
			want: func(a args, f fields) []*questionsv1.AddRequest_Question {
				return []*questionsv1.AddRequest_Question{{
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					IsDeleted:   false,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewQuestionsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.NewQuestionsToPb(tt.args.request)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_QuestionsToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		request []*portal.Question
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*questionsv1.Question
	}{
		{
			name: "correct",
			args: args{
				request: []*portal.Question{{
					Id:          1,
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					DeletedAt:   &testT,
					IsDeleted:   false,
				}},
			},
			want: func(a args, f fields) []*questionsv1.Question {
				t := timestamppb.New(testT)
				f.timeUtils.EXPECT().TimeToTimestamp(a.request[0].CreatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.request[0].UpdatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.request[0].DeletedAt).Return(t)

				return []*questionsv1.Question{{
					Id:          1,
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					CreatedTime: t,
					UpdatedTime: t,
					DeletedTime: t,
					IsDeleted:   false,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewQuestionsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.QuestionsToPb(tt.args.request)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_QuestionsToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		request []*questionsv1.Question
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)
	time := timestamppb.New(testT)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*portal.Question
	}{
		{
			name: "correct",
			args: args{
				request: []*questionsv1.Question{{
					Id:          1,
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					CreatedTime: time,
					UpdatedTime: time,
					DeletedTime: time,
					IsDeleted:   false,
				}},
			},
			want: func(a args, f fields) []*portal.Question {
				f.timeUtils.EXPECT().TimestampToTime(a.request[0].CreatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.request[0].UpdatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.request[0].DeletedTime).Return(&testT)

				return []*portal.Question{{
					Id:          1,
					Name:        "Name 1",
					Description: "Description 1",
					Sort:        1,
					CreatedAt:   &testT,
					UpdatedAt:   &testT,
					DeletedAt:   &testT,
					IsDeleted:   false,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewQuestionsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.QuestionsToEntity(tt.args.request)

			assert.Equal(t, want, got)
		})
	}
}
