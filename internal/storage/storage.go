package storage

import (
	"fmt"
	"strconv"

	"github.com/denistakeda/alerting/internal/metric"
)

type Storage interface {
	StoreGauge(name string, value float64)
	StoreCounter(name string, value int64)
}

func Store(storage Storage, m metric.Metric) error {
	switch m.Type {
	case metric.Gauge:
		val, err := strconv.ParseFloat(m.Value, 64)
		if err != nil {
			return err
		}
		storage.StoreGauge(m.Name, val)
	case metric.Counter:
		val, err := strconv.ParseInt(m.Value, 10, 64)
		if err != nil {
			return err
		}
		storage.StoreCounter(m.Name, val)
	default:
		return fmt.Errorf("no such type \"%s\"", m.Type)
	}
	return nil
}
