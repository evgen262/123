package entitySurveys

import (
	"strings"

	"github.com/google/uuid"
)

//go:generate ditgen -source=image.go

type ImageID string
type ImageData []byte

// Image изображение.
type Image struct {
	// id изображения.
	ID ImageID
	// base64 кодированное изображние.
	Data string
	// Информация об изображении во внешнем хранилище.
	// 	Используется для сохранения в S3.
	ExternalImageInfo *ExternalProperties
}

type ExternalProperties struct {
	// Внешний uuid изображения.
	ID uuid.UUID
	// Наименование изображения.
	FileName string
	// uri изображения во внешнем хранилище.
	URL string
	// Размер изображения.
	Size int64
}

type ImageIDs []ImageID

func (ii ImageIDs) ToStringSlice() []string {
	stringIDs := make([]string, 0, len(ii))
	for _, id := range ii {
		stringIDs = append(stringIDs, string(id))
	}
	return stringIDs
}

func (ii ImageIDs) ToString() string {
	return strings.Join(ii.ToStringSlice(), ",")
}
