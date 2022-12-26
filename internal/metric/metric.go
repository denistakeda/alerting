package metric

import (
	"errors"
	"fmt"
	"strconv"
)

type MetricType int

const (
	Gauge = iota
	Counter
)

var ErrorIncompatibleTypes = errors.New("incompatible types")

type Metric struct {
	metricType   MetricType
	name         string
	gaugeValue   float64
	counterValue int64
}

func NewGauge(name string, value float64) *Metric {
	return &Metric{
		metricType: Gauge,
		name:       name,
		gaugeValue: value,
	}
}

func NewCounter(name string, value int64) *Metric {
	return &Metric{
		metricType:   Counter,
		name:         name,
		counterValue: value,
	}
}

func (m *Metric) Type() MetricType {
	return m.metricType
}

func (m *Metric) Name() string {
	return m.name
}

func (m *Metric) StrValue() string {
	switch m.metricType {
	case Gauge:
		return strconv.FormatFloat(m.gaugeValue, 'f', 3, 64)
	case Counter:
		return strconv.FormatInt(m.counterValue, 10)
	default:
		return ""
	}
}

func (m *Metric) StrType() string {
	switch m.metricType {
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
	if old.metricType != new.metricType {
		return old
	}
	switch old.metricType {
	case Gauge:
		return new
	case Counter:
		return NewCounter(old.name, old.counterValue+new.counterValue)
	default:
		// Should never happen
		return old
	}
}

func TypeFromString(str string) (MetricType, error) {
	switch str {
	case "gauge":
		return Gauge, nil
	case "counter":
		return Counter, nil
	default:
		return -1, fmt.Errorf("unknown metric type: '%s'", str)
	}
}
