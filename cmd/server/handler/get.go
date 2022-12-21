package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getMetricURI struct {
	MetricType string `uri:"metric_type" binding:"required"`
	MetricName string `uri:"metric_name" binding:"required"`
}

func (h *handler) GetMetricHandler(c *gin.Context) {
	var uri getMetricURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	m, ok := h.storage.Get(uri.MetricType, uri.MetricName)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.String(http.StatusOK, m.StrValue())
}
