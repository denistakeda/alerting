package loggerservice

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggerService stores the global logger.
type LoggerService struct{}

// New instantiates a LoggerService.
func New() *LoggerService {
	return &LoggerService{}
}

// ComponentLogger creates a new logger for a component.
func (ls *LoggerService) ComponentLogger(component string) zerolog.Logger {
	return log.Logger.With().Str("component", component).Logger()
}
