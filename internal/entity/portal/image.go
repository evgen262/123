package portal

import (
	"time"
)

//go:generate ditgen -source=image.go

type ImageId int

type ImageData []byte

type Image struct {
	Id        ImageId
	Name      string
	Path      string
	Data      ImageData
	Type      ImageType
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type ImageType int

const (
	ImageTypeUnknown ImageType = iota
	ImageTypeJpeg
	ImageTypePng
	ImageTypeSvg
	ImageTypeGif
)
