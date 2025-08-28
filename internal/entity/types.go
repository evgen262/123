package entity

import (
	"time"

	"github.com/google/uuid"
)

//go:generate ditgen -source=types.go

const (
	DeviceTypeWeb     = "web"
	DeviceTypeAndroid = "Android"
	DeviceTypeIOS     = "iOS"
)

type DeviceType string

func (dt DeviceType) IsValid() bool {
	switch dt {
	case DeviceTypeWeb:
		return true
	case DeviceTypeAndroid:
		return true
	case DeviceTypeIOS:
		return true
	default:
		return false
	}
}

type ID interface {
	GetId() uuid.UUID
	String() string
}

type Pagination struct {
	LastId          ID
	LastCreatedTime *time.Time
	Limit           *int
	Total           *int
}

type StringID interface {
	String() string
}

type StringPagination struct {
	LastId          StringID
	LastCreatedTime *time.Time
	Limit           *int
	Total           *int
}
