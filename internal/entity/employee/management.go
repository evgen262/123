package employee

import (
	"time"

	"github.com/google/uuid"
)

type Management struct {
	ID          uuid.UUID
	Portal      Portal
	RoleID      string
	RoleName    string
	RoleType    string
	EmployeeID  string
	ParentID    string
	ProductID   string
	IsMain      bool
	IsDeleted   bool
	CreatedTime *time.Time
	UpdatedTime *time.Time
}

type ManagementTree struct {
	ID        uuid.UUID
	Name      string
	RoleName  string
	RoleType  string
	ParentID  string
	Children  []*ManagementTree
	Sort      int
	IsDeleted bool
}
