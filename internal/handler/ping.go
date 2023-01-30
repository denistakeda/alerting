package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) PingHandler(c *gin.Context) {
	if err := h.dbStorage.Ping(c); err != nil {
		log.Println(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}
	c.String(http.StatusOK, "pong")
}
