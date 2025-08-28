package employee

import (
	"github.com/google/uuid"
)

type Product struct {
	ID            string
	PortalID      int32
	ClusterID     string
	ClusterName   string
	FullName      string
	ShortName     string
	Type          string
	IconID        uuid.UUID
	ResponsibleID string
	Tutor         string
	IsMain        bool
}

func (p *Product) GetFullName() string {
	if p == nil {
		return ""
	}
	return p.FullName
}
