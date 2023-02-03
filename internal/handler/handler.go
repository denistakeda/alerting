package handler

import (
	"github.com/denistakeda/alerting/internal/services/logger_service"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/rs/zerolog"
)

type Handler struct {
	storage s.Storage
	hashKey string
	logger  zerolog.Logger
}

func New(
	storage s.Storage,
	hashKey string,
	logService *logger_service.LoggerService,
) *Handler {
	return &Handler{
		storage: storage,
		hashKey: hashKey,
		logger:  logService.ComponentLogger("Handler"),
	}
}
