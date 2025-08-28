package portals

import (
	"time"
)

type Questions struct {
	Email     string          `json:"helpdeskEmail"`
	Questions []*QuestionInfo `json:"questions"`
} // @name Questions

type Question struct {
	Id          int        `json:"id,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	Name        string     `json:"title"`
	Description string     `json:"description"`
	Sort        int        `json:"sort,omitempty"`
	IsDeleted   bool       `json:"isDeleted"`
} // @name Question

type NewQuestion struct {
	Name        string `json:"title"`
	Description string `json:"description"`
	Sort        int    `json:"sort,omitempty"`
	IsDeleted   bool   `json:"isDeleted"`
} // @name NewQuestion

type UpdateQuestion struct {
	Id          int    `json:"id,omitempty" swaggerignore:"true"`
	Name        string `json:"title"`
	Description string `json:"description"`
	Sort        int    `json:"sort,omitempty"`
	IsDeleted   bool   `json:"isDeleted"`
} // @name UpdateQuestion

type QuestionInfo struct {
	Name        string `json:"title"`
	Description string `json:"description"`
} // @name QuestionInfo
