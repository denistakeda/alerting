package ports

import "github.com/denistakeda/alerting/internal/metric"

type Client interface {
	SendMetrics([]*metric.Metric) error
}
