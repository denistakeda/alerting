package loggerservice

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerService struct{}

func New() *LoggerService {
	return &LoggerService{}
}

func (ls *LoggerService) ComponentLogger(component string) zerolog.Logger {
	return log.Logger.With().Str("component", component).Logger()
}
