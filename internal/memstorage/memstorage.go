package memstorage

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/denistakeda/alerting/internal/metric"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (ms *MemStorage) StoreGauge(name string, value float64) {
	ms.gauges[name] = value
}

func (ms *MemStorage) StoreCounter(name string, value int64) {
	oldValue, ok := ms.counters[name]
	if !ok {
		oldValue = 0
	}
	ms.counters[name] = oldValue + value
}

func (ms *MemStorage) Metrics() []metric.Metric {
	res := make([]metric.Metric, len(ms.gauges)+len(ms.counters))
	i := 0
	for name, value := range ms.gauges {
		res[i] = metric.Metric{
			Type:  metric.Gauge,
			Name:  name,
			Value: strconv.FormatFloat(value, 'E', -1, 64),
		}
		i++
	}
	for name, value := range ms.counters {
		res[i] = metric.Metric{
			Type:  metric.Counter,
			Name:  name,
			Value: strconv.FormatInt(value, 10),
		}
		i++
	}
	return res
}

func (ms *MemStorage) ToString() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "Gauges:\n")
	for name, value := range ms.gauges {
		fmt.Fprintf(b, "%s=%v\n", name, value)
	}
	fmt.Fprintf(b, "Counters:\n")
	for name, value := range ms.counters {
		fmt.Fprintf(b, "%s=%v\n", name, value)
	}
	return b.String()
}
