package handler

import (
	"net/http"

	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/gin-gonic/gin"
)

type getMetricURI struct {
	MetricType string `uri:"metric_type" binding:"required"`
	MetricName string `uri:"metric_name" binding:"required"`
}

func GetMetricHandler(storage s.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri getMetricURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		m, ok := storage.Get(uri.MetricType, uri.MetricName)
		if !ok {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, m.StrValue())
	}
}
