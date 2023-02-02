package handler

import (
	"log"
	"net/http"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/gin-gonic/gin"
)

type getMetricURI struct {
	MetricType string `uri:"metric_type" binding:"required"`
	MetricName string `uri:"metric_name" binding:"required"`
}

func (h *Handler) GetMetricHandler(c *gin.Context) {
	var uri getMetricURI
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Println(c.AbortWithError(http.StatusBadRequest, err))
		return
	}

	metricType, err := metric.TypeFromString(uri.MetricType)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	m, ok := h.storage.Get(c, metricType, uri.MetricName)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.String(http.StatusOK, m.StrValue())
}

func (h *Handler) GetMetricHandler2(c *gin.Context) {
	var requestMetric *metric.Metric
	if err := c.ShouldBind(&requestMetric); err != nil {
		log.Println(c.AbortWithError(http.StatusBadRequest, err))
		return
	}
	m, ok := h.storage.Get(c, requestMetric.Type(), requestMetric.Name())
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, requestMetric)
		return
	}
	c.JSON(http.StatusOK, m)
}
