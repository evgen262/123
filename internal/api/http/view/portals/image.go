package portals

import (
	"time"
)

type Image struct {
	Id        int        `json:"id,omitempty"`
	Name      string     `json:"name"`
	Path      string     `json:"path,omitempty"`
	Data      string     `json:"data,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
} // @name Image

type NewImage struct {
	// Наименование файла
	Name string `json:"name" binding:"required"`
	// Base64 содержимое файла
	Data string `json:"data" format:"base64" binding:"required"`
} // @name NewImage

type ImageInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path,omitempty"`
} // @name ImageInfo
