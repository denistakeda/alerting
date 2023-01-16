package handler

import (
	"errors"
	"log"
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

func (h *Handler) UpdateMetricHandler(c *gin.Context) {
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

	if _, err := h.storage.Update(m); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.Status(http.StatusOK)
}

func (h *Handler) UpdateMetricHandler2(c *gin.Context) {
	var m *metric.Metric
	if err := c.ShouldBind(&m); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Printf("UpdateMetricHandler2: request: %v, response: ", m)
	if err := m.Validate(); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	m, err := h.storage.Update(m)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	log.Printf("%v\n", m)
	c.JSON(http.StatusOK, m)
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
