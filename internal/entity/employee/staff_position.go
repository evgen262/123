package employee

import (
	"time"

	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
)

type Position struct {
	ID   uuid.UUID
	Name string
}

func (p *Position) GetID() uuid.UUID {
	if p == nil {
		return uuid.Nil
	}
	return p.ID
}

func (p *Position) GetName() string {
	if p == nil {
		return ""
	}
	return p.Name
}

type StaffPosition struct {
	ID            uuid.UUID
	Portal        Portal
	Organization  Organization
	Subdivision   *Subdivision
	Position      *Position
	ResponsibleID string
	Name          string
	RateNumbers   int
	CreateDate    *timeUtils.Date
	CloseDate     *timeUtils.Date
	CreatedTime   *time.Time
	UpdatedTime   *time.Time
}

func (sp *StaffPosition) GetPosition() *Position {
	if sp == nil {
		return nil
	}
	return sp.Position
}

func (sp *StaffPosition) GetSubdivision() *Subdivision {
	if sp == nil {
		return nil
	}
	return sp.Subdivision
}

type Subdivision struct {
	ID        uuid.UUID
	Name      string
	ParentID  string
	Sort      int
	IsDeleted bool
}

func (s *Subdivision) GetID() uuid.UUID {
	if s == nil {
		return uuid.Nil
	}
	return s.ID
}

func (s *Subdivision) GetName() string {
	if s == nil {
		return ""
	}
	return s.Name
}

type SubdivisionTree struct {
	ID                   uuid.UUID
	Name                 string
	ManagementPositionID string
	ParentID             string
	Children             []*SubdivisionTree
	Sort                 int
	IsDeleted            bool
}
