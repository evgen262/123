package portals

import (
	"time"
)

// Features Список фича-флагов
//
//	где ключом мапы является наименование фичи
type Features map[string]*FeatureInfo // @name Features

type Feature struct {
	Id        int        `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Version   string     `json:"version"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Enabled   bool       `json:"enabled"`
} // @name Feature

type NewFeature struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
} // @name NewFeature

type UpdateFeature struct {
	Id      int    `json:"id,omitempty" swaggerignore:"true"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
} // @name UpdateFeature

type FeatureInfo struct {
	Id      int    `json:"id"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
} // @name FeatureInfo
