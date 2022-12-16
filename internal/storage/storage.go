package storage

import (
	"strconv"

	"github.com/denistakeda/alerting/internal/metric"
)

type Storage interface {
	StoreGauge(name string, value float64)
	StoreCounter(name string, value int64)
}

func Store(storage Storage, metricType metric.Type, metricName, metricValue string) error {
	switch metricType {
	case metric.Gauge:
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return err
		}
		storage.StoreGauge(metricName, val)
	case metric.Counter:
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return err
		}
		storage.StoreCounter(metricName, val)
	}
	return nil
}
