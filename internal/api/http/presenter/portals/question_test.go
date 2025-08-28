package portals

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

func Test_questionPresenter_ToEntity(t *testing.T) {
	type args struct {
		question *viewPortals.UpdateQuestion
	}
	tests := []struct {
		name string
		args args
		want *portal.Question
	}{
		{
			name: "from viewPortals to portal",
			args: args{
				question: &viewPortals.UpdateQuestion{
					Id:          5,
					Name:        "Question Name 5",
					Description: "Some question 5 description",
					Sort:        12,
					IsDeleted:   false,
				},
			},
			want: &portal.Question{
				Id:          5,
				Name:        "Question Name 5",
				Description: "Some question 5 description",
				Sort:        12,
				IsDeleted:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewQuestionPresenter()
			assert.Equalf(t, tt.want, qp.ToEntity(tt.args.question), "ToEntity(%v)", tt.args.question)
		})
	}
}

func Test_questionPresenter_ToEntities(t *testing.T) {
	type args struct {
		questions []*viewPortals.UpdateQuestion
	}
	tests := []struct {
		name string
		args args
		want []*portal.Question
	}{
		{
			name: "from views to entities",
			args: args{
				questions: []*viewPortals.UpdateQuestion{
					{
						Id:          1,
						Name:        "Question Name 1",
						Description: "Some question 1 description",
						Sort:        10,
						IsDeleted:   false,
					},
					{
						Id:          2,
						Name:        "Question Name 2",
						Description: "Some question 2 description",
						Sort:        12,
						IsDeleted:   true,
					},
					{
						Id:          3,
						Name:        "Question Name 3",
						Description: "Some question 3 description",
						Sort:        12,
						IsDeleted:   false,
					},
				},
			},
			want: []*portal.Question{
				{
					Id:          1,
					Name:        "Question Name 1",
					Description: "Some question 1 description",
					Sort:        10,
					IsDeleted:   false,
				},
				{
					Id:          2,
					Name:        "Question Name 2",
					Description: "Some question 2 description",
					Sort:        12,
					IsDeleted:   true,
				},
				{
					Id:          3,
					Name:        "Question Name 3",
					Description: "Some question 3 description",
					Sort:        12,
					IsDeleted:   false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewQuestionPresenter()
			assert.Equalf(t, tt.want, qp.ToEntities(tt.args.questions), "ToEntities(%v)", tt.args.questions)
		})
	}
}

func Test_questionPresenter_ToNewEntity(t *testing.T) {
	type args struct {
		question *viewPortals.NewQuestion
	}
	tests := []struct {
		name string
		args args
		want *portal.Question
	}{
		{
			name: "from viewPortals to new portal",
			args: args{
				question: &viewPortals.NewQuestion{
					Name:        "New question name",
					Description: "Some new question description",
					Sort:        9,
					IsDeleted:   false,
				},
			},
			want: &portal.Question{
				Name:        "New question name",
				Description: "Some new question description",
				Sort:        9,
				IsDeleted:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewQuestionPresenter()
			assert.Equalf(t, tt.want, qp.ToNewEntity(tt.args.question), "ToNewEntity(%v)", tt.args.question)
		})
	}
}

func Test_questionPresenter_ToNewEntities(t *testing.T) {
	type args struct {
		questions []*viewPortals.NewQuestion
	}
	tests := []struct {
		name string
		args args
		want []*portal.Question
	}{
		{
			name: "from views to new entities",
			args: args{
				questions: []*viewPortals.NewQuestion{
					{
						Name:        "New question 1 name",
						Description: "Some new question 1 description",
						Sort:        1,
						IsDeleted:   false,
					},
					{
						Name:        "New question 2 name",
						Description: "Some new question 2 description",
						Sort:        1,
						IsDeleted:   false,
					},
				},
			},
			want: []*portal.Question{
				{
					Name:        "New question 1 name",
					Description: "Some new question 1 description",
					Sort:        1,
					IsDeleted:   false,
				},
				{
					Name:        "New question 2 name",
					Description: "Some new question 2 description",
					Sort:        1,
					IsDeleted:   false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewQuestionPresenter()
			assert.Equalf(t, tt.want, qp.ToNewEntities(tt.args.questions), "ToNewEntities(%v)", tt.args.questions)
		})
	}
}

func Test_questionPresenter_ToShortView(t *testing.T) {
	type args struct {
		question *portal.Question
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.QuestionInfo
	}{
		{
			name: "from portal to short viewPortals",
			args: args{
				question: &portal.Question{
					Id:          10,
					Name:        "The Question",
					Description: "Question description",
					Sort:        15,
					IsDeleted:   false,
				},
			},
			want: &viewPortals.QuestionInfo{
				Name:        "The Question",
				Description: "Question description",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewQuestionPresenter()
			assert.Equalf(t, tt.want, qp.ToShortView(tt.args.question), "ToShortView(%v)", tt.args.question)
		})
	}
}

func Test_questionPresenter_ToShortViews(t *testing.T) {
	type args struct {
		questions []*portal.Question
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.QuestionInfo
	}{
		{
			name: "from entities to short views",
			args: args{
				questions: []*portal.Question{
					{
						Id:          1,
						Name:        "Question 1",
						Description: "Question 1 description",
						Sort:        12,
						IsDeleted:   false,
					},
					{
						Id:          2,
						Name:        "Question 2",
						Description: "Question 2 description",
						Sort:        13,
						IsDeleted:   false,
					},
				},
			},
			want: []*viewPortals.QuestionInfo{
				{
					Name:        "Question 1",
					Description: "Question 1 description",
				},
				{
					Name:        "Question 2",
					Description: "Question 2 description",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewQuestionPresenter()
			assert.Equalf(t, tt.want, qp.ToShortViews(tt.args.questions), "ToShortViews(%v)", tt.args.questions)
		})
	}
}

func Test_questionPresenter_ToView(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)

	type args struct {
		question *portal.Question
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.Question
	}{
		{
			name: "from portal to viewPortals",
			args: args{
				question: &portal.Question{
					Id:          1,
					Name:        "Question 1",
					Description: "Question 1 description",
					Sort:        5,
					CreatedAt:   &testTime,
					UpdatedAt:   &testTime,
					DeletedAt:   &testTime,
					IsDeleted:   true,
				},
			},
			want: &viewPortals.Question{
				Id:          1,
				Name:        "Question 1",
				Description: "Question 1 description",
				Sort:        5,
				CreatedAt:   &testTime,
				UpdatedAt:   &testTime,
				DeletedAt:   &testTime,
				IsDeleted:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewQuestionPresenter()
			assert.Equalf(t, tt.want, qp.ToView(tt.args.question), "ToView(%v)", tt.args.question)
		})
	}
}

func Test_questionPresenter_ToViews(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 10, 10, 10, 0, time.UTC)

	type args struct {
		questions []*portal.Question
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.Question
	}{
		{
			name: "from entities to views",
			args: args{
				questions: []*portal.Question{
					{
						Id:          1,
						Name:        "Question 1",
						Description: "Question 1 description",
						Sort:        5,
						CreatedAt:   &testTime,
						UpdatedAt:   &testTime,
						DeletedAt:   &testTime,
						IsDeleted:   true,
					},
					{
						Id:          2,
						Name:        "Question 2",
						Description: "Question 2 description",
						Sort:        7,
						CreatedAt:   &testTime,
						UpdatedAt:   &testTime,
						DeletedAt:   nil,
						IsDeleted:   false,
					},
					{
						Id:          3,
						Name:        "Question 3",
						Description: "Question 3 description",
						Sort:        9,
						CreatedAt:   &testTime,
						UpdatedAt:   nil,
						DeletedAt:   nil,
						IsDeleted:   false,
					},
				},
			},
			want: []*viewPortals.Question{
				{
					Id:          1,
					Name:        "Question 1",
					Description: "Question 1 description",
					Sort:        5,
					CreatedAt:   &testTime,
					UpdatedAt:   &testTime,
					DeletedAt:   &testTime,
					IsDeleted:   true,
				},
				{
					Id:          2,
					Name:        "Question 2",
					Description: "Question 2 description",
					Sort:        7,
					CreatedAt:   &testTime,
					UpdatedAt:   &testTime,
					DeletedAt:   nil,
					IsDeleted:   false,
				},
				{
					Id:          3,
					Name:        "Question 3",
					Description: "Question 3 description",
					Sort:        9,
					CreatedAt:   &testTime,
					UpdatedAt:   nil,
					DeletedAt:   nil,
					IsDeleted:   false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewQuestionPresenter()
			assert.Equalf(t, tt.want, qp.ToViews(tt.args.questions), "ToViews(%v)", tt.args.questions)
		})
	}
}
