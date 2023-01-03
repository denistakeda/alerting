package metric

import (
	"fmt"
	"strconv"
)

type Type string

const (
	Gauge   Type = "gauge"
	Counter Type = "counter"
)

type Metric struct {
	ID    string   `json:"name"`
	MType Type     `json:"type"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
}

func NewGauge(name string, value float64) *Metric {
	return &Metric{
		MType: Gauge,
		ID:    name,
		Value: &value,
	}
}

func NewCounter(name string, value int64) *Metric {
	return &Metric{
		MType: Counter,
		ID:    name,
		Delta: &value,
	}
}

func (m *Metric) Type() Type {
	return m.MType
}

func (m *Metric) Name() string {
	return m.ID
}

func (m *Metric) StrValue() string {
	switch m.MType {
	case Gauge:
		return strconv.FormatFloat(*m.Value, 'f', 3, 64)
	case Counter:
		return strconv.FormatInt(*m.Delta, 10)
	default:
		return ""
	}
}

func (m *Metric) StrType() string {
	switch m.MType {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	default:
		return ""
	}
}

func Update(old *Metric, new *Metric) *Metric {
	if old == nil {
		return new
	}
	if new == nil {
		return old
	}
	if old.MType != new.MType {
		return old
	}
	switch old.MType {
	case Gauge:
		return new
	case Counter:
		return NewCounter(old.ID, *old.Delta+*new.Delta)
	default:
		// Should never happen
		return old
	}
}

func TypeFromString(str string) (Type, error) {
	switch str {
	case "gauge":
		return Gauge, nil
	case "counter":
		return Counter, nil
	default:
		return "", fmt.Errorf("unknown metric type: '%s'", str)
	}
}
