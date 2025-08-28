package portal

import (
	"time"
)

//go:generate ditgen -source=question.go

type QuestionId int

type Question struct {
	Id          QuestionId
	Name        string
	Description string
	Sort        int
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	IsDeleted   bool
}

type Questions struct {
	SupportEmail string
	Questions    []*Question
}
