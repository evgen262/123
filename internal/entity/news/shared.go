package news

import (
	"github.com/google/uuid"
)

//go:generate ditgen -source=./shared.go -zero=true -all-fields=true

type Visitor struct {
	EmployeeID *uuid.UUID
	PortalID   int
}
