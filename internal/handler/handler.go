package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/denistakeda/alerting/internal/services/loggerservice"
	s "github.com/denistakeda/alerting/internal/storage"
)

type Handler struct {
	hashKey    string
	logger     zerolog.Logger
	cert       string
	privateKey string

	engine  *gin.Engine
	storage s.Storage

	server *http.Server
}

type Params struct {
	Addr       string
	HashKey    string
	Cert       string
	PrivateKey string

	Engine     *gin.Engine
	Storage    s.Storage
	LogService *loggerservice.LoggerService
}

func New(params Params) *Handler {

	handler := &Handler{
		engine:     params.Engine,
		storage:    params.Storage,
		hashKey:    params.HashKey,
		cert:       params.Cert,
		privateKey: params.PrivateKey,
		logger:     params.LogService.ComponentLogger("Handler"),

		server: &http.Server{
			Addr:    params.Addr,
			Handler: params.Engine,
		},
	}

	handler.registerHandlers(params.Engine)

	return handler
}

func (h *Handler) Start() <-chan error {
	out := make(chan error)

	go func() {
		var err error
		if h.cert != "" && h.privateKey != "" {
			err = h.server.ListenAndServeTLS(h.cert, h.privateKey)
		} else {
			err = h.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			out <- errors.Wrap(err, "critical api failure")
		}
	}()

	return out
}

func (h *Handler) Stop() {
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Fatal().Err(err).Msg("failed to stop server")
	}

	h.logger.Info().Msg("server exiting")
}

func (h *Handler) registerHandlers(engine *gin.Engine) {
	engine.POST("/update/", h.UpdateMetricHandler2)
	engine.POST("/update/:metric_type/:metric_name/:metric_value", h.UpdateMetricHandler)
	engine.POST("/updates/", h.UpdateMetricsHandler)
	engine.POST("/value/", h.GetMetricHandler2)
	engine.GET("/value/:metric_type/:metric_name", h.GetMetricHandler)
	engine.GET("/ping", h.PingHandler)
	engine.GET("/", h.MainPageHandler)
}
