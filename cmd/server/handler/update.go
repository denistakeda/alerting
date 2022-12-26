package handler

import (
	"errors"
	"net/http"
	"strconv"

	errextra "github.com/pkg/errors"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/gin-gonic/gin"
)

var ErrUnknownMetricType = errors.New("unknown metric type")
var ErrIncorrectValue = errors.New("incorrect metric value")

type updateMetricURI struct {
	MetricType  string `uri:"metric_type" binding:"required"`
	MetricName  string `uri:"metric_name" binding:"required"`
	MetricValue string `uri:"metric_value" binding:"required"`
}

func (h *handler) UpdateMetricHandler(c *gin.Context) {
	var uri updateMetricURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	m, err := createMetric(uri)
	if errors.Is(err, ErrUnknownMetricType) {
		c.AbortWithError(http.StatusNotImplemented, err)
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.storage.Update(m); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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
