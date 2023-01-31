package storage

import (
	"context"
	"github.com/denistakeda/alerting/internal/metric"
)

type Storage interface {
	Get(metricType metric.Type, metricName string) (*metric.Metric, bool)
	Update(metric *metric.Metric) (*metric.Metric, error)
	All() []*metric.Metric
	Close() error
	Ping(ctx context.Context) error
}
