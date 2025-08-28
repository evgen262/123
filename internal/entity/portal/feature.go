package portal

import (
	"time"
)

//go:generate ditgen -source=feature.go

type FeatureId int

type Feature struct {
	Id        FeatureId
	Name      string
	Version   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Enabled   bool
}
