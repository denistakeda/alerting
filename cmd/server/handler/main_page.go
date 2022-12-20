package handler

import (
	"net/http"

	"github.com/denistakeda/alerting/internal/metric"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/gin-gonic/gin"
)

type met struct {
	Type  string
	Name  string
	Value string
}

func metricsToRepresentation(metrics []metric.Metric) []met {
	ms := make([]met, len(metrics))
	for i, m := range metrics {
		ms[i] = met{
			Type:  m.Type(),
			Name:  m.Name(),
			Value: m.StrValue(),
		}
	}
	return ms
}

func MainPageHandler(storage s.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Metrics": metricsToRepresentation(storage.All()),
		})
	}
}
