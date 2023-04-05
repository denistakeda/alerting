package handler

import (
	"errors"
	"net/http"
	"strconv"

	errextra "github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/denistakeda/alerting/internal/metric"
)

var (
	ErrUnknownMetricType = errors.New("unknown metric type")
	ErrIncorrectValue    = errors.New("incorrect metric value")
)

type updateMetricURI struct {
	MetricType  string `uri:"metric_type" binding:"required"`
	MetricName  string `uri:"metric_name" binding:"required"`
	MetricValue string `uri:"metric_value" binding:"required"`
}

func (h *Handler) UpdateMetricHandler(c *gin.Context) {
	var uri updateMetricURI
	if err := c.ShouldBindUri(&uri); err != nil {
		h.logger.Warn().Err(err).Msg("failed to bind uri")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	m, err := createMetric(uri)
	if errors.Is(err, ErrUnknownMetricType) {
		h.logger.Warn().Err(err).Msg("unknown metric type")
		c.AbortWithStatus(http.StatusNotImplemented)
		return
	}
	if err != nil {
		h.logger.Warn().Err(err).Msg("failed to create a metric")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := h.storage.Update(c, m); err != nil {
		h.logger.Warn().Err(err).Msgf("failed to update a metric %v", m)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.Status(http.StatusOK)
}

func (h *Handler) UpdateMetricHandler2(c *gin.Context) {
	var m *metric.Metric
	if err := c.ShouldBind(&m); err != nil {
		h.logger.Warn().Err(err).Msg("failed to bind uri")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := m.Validate(); err != nil {
		h.logger.Warn().Err(err).Msgf("incorrect metric %v", m)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := m.VerifyHash(h.hashKey); err != nil {
		h.logger.Warn().Err(err).Msgf("incorrect metric hash %v", m.Hash)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	m, err := h.storage.Update(c, m)
	if err != nil {
		h.logger.Warn().Err(err).Msgf("failed to update a metric %v", m)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *Handler) UpdateMetricsHandler(c *gin.Context) {
	var metrics []*metric.Metric
	if err := c.ShouldBind(&metrics); err != nil {
		h.logger.Warn().Err(err).Msg("failed to bind uri")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	for _, m := range metrics {
		if err := m.Validate(); err != nil {
			h.logger.Warn().Err(err).Msgf("incorrect metric %v", m)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := m.VerifyHash(h.hashKey); err != nil {
			h.logger.Warn().Err(err).Msgf("incorrect metric hash %v", m.Hash)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	if err := h.storage.UpdateAll(c, metrics); err != nil {
		h.logger.Warn().Err(err).Msg("failed to update a metrics")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

func createMetric(uri updateMetricURI) (*metric.Metric, error) {
	switch uri.MetricType {
	case "gauge":
		val, err := strconv.ParseFloat(uri.MetricValue, 64)
		if err != nil {
			return nil, errextra.Wrapf(ErrIncorrectValue, "expected to be float64, got \"%s\"", uri.MetricValue)
		}
		return metric.NewGauge(uri.MetricName, val), nil
	case "counter":
		val, err := strconv.ParseInt(uri.MetricValue, 10, 64)
		if err != nil {
			return nil, errextra.Wrapf(ErrIncorrectValue, "expected to be int64, got \"%s\"", uri.MetricValue)
		}
		return metric.NewCounter(uri.MetricName, val), nil
	default:
		return nil, errextra.Wrapf(ErrUnknownMetricType, "expected \"gauge\" or \"counter\", got \"%s\"", uri.MetricType)
	}
}
