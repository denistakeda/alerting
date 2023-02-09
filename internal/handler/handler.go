package handler

import (
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/gin-gonic/gin"
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
	logService *loggerservice.LoggerService,
) *Handler {
	return &Handler{
		storage: storage,
		hashKey: hashKey,
		logger:  logService.ComponentLogger("Handler"),
	}
}

func (h *Handler) RegisterHandlers(engine *gin.Engine) {
	engine.POST("/update/", h.UpdateMetricHandler2)
	engine.POST("/update/:metric_type/:metric_name/:metric_value", h.UpdateMetricHandler)
	engine.POST("/updates/", h.UpdateMetricsHandler)
	engine.POST("/value/", h.GetMetricHandler2)
	engine.GET("/value/:metric_type/:metric_name", h.GetMetricHandler)
	engine.GET("/ping", h.PingHandler)
	engine.GET("/", h.MainPageHandler)
}
