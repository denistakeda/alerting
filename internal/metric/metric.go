package metric

import "fmt"

type Type string

const (
	Gauge   Type = "gauge"
	Counter Type = "counter"
)

type Metric struct {
	Type  Type
	Name  string
	Value string
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
