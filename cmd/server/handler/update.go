package handler

import (
	"net/http"

	"github.com/denistakeda/alerting/internal/metric"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/gin-gonic/gin"
)

func UpdateMetricHandler(storage s.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m metric.Metric
		if err := c.ShouldBindUri(&m); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if _, err := metric.ParseType(string(m.Type)); err != nil {
			c.AbortWithError(http.StatusNotImplemented, err)
		}
		if err := s.Store(storage, m); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		c.Status(http.StatusOK)
	}
}
