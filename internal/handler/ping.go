package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PingHandler(c *gin.Context) {
	if err := h.storage.Ping(c); err != nil {
		h.logger.Error().Err(err).Msg("failed to ping database")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, "pong")
}
