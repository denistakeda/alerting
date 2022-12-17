package metric

import "fmt"

type Type string

const (
	Gauge   Type = "gauge"
	Counter Type = "counter"
)

type Metric struct {
	Type  Type   `uri:"metric_type" binding:"required"`
	Name  string `uri:"metric_name" binding:"required"`
	Value string `uri:"metric_value" binding:"required"`
}

func ParseType(metricType string) (Type, error) {
	switch metricType {
	case string(Gauge):
		return Gauge, nil
	case string(Counter):
		return Counter, nil
	default:
		return "", fmt.Errorf("no such type \"%s\"", metricType)
	}
}
