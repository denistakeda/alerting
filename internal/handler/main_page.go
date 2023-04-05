package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/denistakeda/alerting/internal/metric"
)

type met struct {
	Type  string
	Name  string
	Value string
}

func metricsToRepresentation(metrics []*metric.Metric) []met {
	ms := make([]met, len(metrics))
	for i, m := range metrics {
		ms[i] = met{
			Type:  m.StrType(),
			Name:  m.Name(),
			Value: m.StrValue(),
		}
	}
	return ms
}

func (h *Handler) MainPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Metrics": metricsToRepresentation(h.storage.All(c)),
	})
}
