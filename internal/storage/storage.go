package storage

import (
	"context"

	"github.com/denistakeda/alerting/internal/metric"
)

type Storage interface {
	// Get returns a metric if exists.
	Get(ctx context.Context, metricType metric.Type, metricName string) (*metric.Metric, bool)
	// Update updates a metric if exists.
	Update(ctx context.Context, metric *metric.Metric) (*metric.Metric, error)
	// UpdateAll updates all the metrics in list.
	UpdateAll(ctx context.Context, metrics []*metric.Metric) error
	// All returns all the metrics.
	All(ctx context.Context) []*metric.Metric
	// Close closes the connection to db.
	Close(ctx context.Context) error
	// Ping pings the database.
	Ping(ctx context.Context) error
}
