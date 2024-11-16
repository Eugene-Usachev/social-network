package logger

import "github.com/Eugune-Usachev/social-network/src/internal/config"

type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

func MustCreateLogger(isProduction bool, appConfig *config.AppConfig) Logger {
	if isProduction {
		return MustCreateElasticSearchLogger(appConfig.EsAddr(), appConfig.EsUser(), appConfig.EsPass())
	}

	return NewZeroLogger()
}
