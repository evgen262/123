package analytics

import (
	"strings"

	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

const (
	keyDeviceID   = "DeviceID"
	keyDeviceType = "DeviceType"
)

type XCFCUserAgentHeader string

type CFCHeaders struct {
	Header        XCFCUserAgentHeader // X-CFC-UserAgent передаем как есть парсится внутри analytics
	Portal        string              // Из сессии
	Authorization string
}

func (h XCFCUserAgentHeader) IsValid() bool {
	if h == "" {
		return false
	}
	cfcHeaders := strings.Split(string(h), ";")
	var isDeviceIDValid, isDeviceTypeValid bool
	for _, header := range cfcHeaders {
		key, value, ok := strings.Cut(header, "=")
		if !ok {
			return false
		}
		switch key {
		case keyDeviceID:
			if err := uuid.Validate(value); err != nil {
				return false
			}
			isDeviceIDValid = true
		case keyDeviceType:
			if !entity.DeviceType(value).IsValid() {
				return false
			}
			isDeviceTypeValid = true
		}
	}

	return isDeviceIDValid && isDeviceTypeValid
}
