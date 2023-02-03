package handler

import (
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
		h.logger.Warn().Err(err).Msg("failed to bind uri")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	metricType, err := metric.TypeFromString(uri.MetricType)
	if err != nil {
		h.logger.Warn().Err(err).Msgf("wrong metric type '%s'", uri.MetricType)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	m, ok := h.storage.Get(c, metricType, uri.MetricName)
	if !ok {
		h.logger.Warn().Err(err).Msgf("no such metric with type '%s' and name '%s'", metricType, uri.MetricName)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.String(http.StatusOK, m.StrValue())
}

func (h *Handler) GetMetricHandler2(c *gin.Context) {
	var requestMetric *metric.Metric
	if err := c.ShouldBind(&requestMetric); err != nil {
		h.logger.Warn().Err(err).Msg("unable to bind uri")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	m, ok := h.storage.Get(c, requestMetric.Type(), requestMetric.Name())
	if !ok {
		h.logger.Warn().Msgf("metric not found %v", m)
		c.AbortWithStatusJSON(http.StatusNotFound, requestMetric)
		return
	}
	c.JSON(http.StatusOK, m)
}
