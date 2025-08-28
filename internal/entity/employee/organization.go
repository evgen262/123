package employee

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID              uuid.UUID
	GlobalID        *string
	Name            string
	ShortName       string
	FullName        *string
	INN             string
	IconID          string // TODO: пока непонятно откуда брать
	IsActive        bool
	IsLiquidated    bool
	CreatedTime     *time.Time
	LiquidationTime *time.Time
}
