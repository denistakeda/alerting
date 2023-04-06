package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler godoc
// @Summary health of the service
// @Accept  json
// @Produce text/plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func (h *Handler) PingHandler(c *gin.Context) {
	if err := h.storage.Ping(c); err != nil {
		h.logger.Error().Err(err).Msg("failed to ping database")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, "pong")
}
