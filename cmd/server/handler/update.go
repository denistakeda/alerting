package handler

import (
	"errors"
	"net/http"

	errextra "github.com/pkg/errors"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/metric/counter"
	"github.com/denistakeda/alerting/internal/metric/gauge"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/gin-gonic/gin"
)

var ErrUnknownMetricType = errors.New("unknown metric type")

type updateMetricURI struct {
	MetricType  string `uri:"metric_type" binding:"required"`
	MetricName  string `uri:"metric_name" binding:"required"`
	MetricValue string `uri:"metric_value" binding:"required"`
}

func UpdateMetricHandler(storage s.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		if err := storage.Update(m); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		c.Status(http.StatusOK)
	}
}

func createMetric(uri updateMetricURI) (metric.Metric, error) {
	switch uri.MetricType {
	case "gauge":
		return gauge.FromStr(uri.MetricName, uri.MetricValue)
	case "counter":
		return counter.FromStr(uri.MetricName, uri.MetricValue)
	default:
		return nil, errextra.Wrapf(ErrUnknownMetricType, "expected \"gauge\" or \"counter\", got \"%s\"", uri.MetricType)
	}
}
