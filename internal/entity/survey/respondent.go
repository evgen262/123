package entitySurveys

import (
	"strings"

	"github.com/google/uuid"
)

//go:generate ditgen -source=respondent.go

type RespondentID uuid.UUID

func (r RespondentID) String() string {
	return uuid.UUID(r).String()
}

type RespondentIDs []RespondentID

func (ri RespondentIDs) ToStringSlice() []string {
	stringIds := make([]string, 0, len(ri))
	for _, id := range ri {
		stringIds = append(stringIds, uuid.UUID(id).String())
	}
	return stringIds
}

func (ri RespondentIDs) ToString() string {
	return strings.Join(ri.ToStringSlice(), ",")
}
