package auth

import "go.uber.org/zap"

var (
	LogModuleUE = zap.String("module", "Универсальный вход")
	LogCode     = func(code string) zap.Field {
		return zap.String("code", code)
	}
)
