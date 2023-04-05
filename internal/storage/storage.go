package storage

import (
	"context"

	"github.com/denistakeda/alerting/internal/metric"
)

type Storage interface {
	Get(ctx context.Context, metricType metric.Type, metricName string) (*metric.Metric, bool)
	Update(ctx context.Context, metric *metric.Metric) (*metric.Metric, error)
	UpdateAll(ctx context.Context, metrics []*metric.Metric) error
	All(ctx context.Context) []*metric.Metric
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
}
