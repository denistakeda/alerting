package storage

import (
	"github.com/denistakeda/alerting/internal/metric"
)

type Storage interface {
	Get(metricType string, metricName string) (metric.Metric, bool)
	Update(metric metric.Metric) error
	All() []metric.Metric
}
