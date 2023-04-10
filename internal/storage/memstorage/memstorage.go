package memstorage

import (
	"context"
	"sync"

	"github.com/denistakeda/alerting/internal/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
)

var _ storage.Storage = (*Memstorage)(nil)

// Memstorage is a memory storage.
type Memstorage struct {
	gauges   map[string]*metric.Metric
	counters map[string]*metric.Metric
	hashKey  string
	mx       sync.Mutex
	logger   zerolog.Logger
}

// NewMemStorage instantiates a new MemStorage instance.
func NewMemStorage(hashKey string, logService *loggerservice.LoggerService) *Memstorage {
	return &Memstorage{
		gauges:   make(map[string]*metric.Metric),
		counters: make(map[string]*metric.Metric),
		hashKey:  hashKey,
		logger:   logService.ComponentLogger("Memstorage"),
	}
}

// Get returns a metric if exists.
func (m *Memstorage) Get(_ context.Context, metricType metric.Type, metricName string) (*metric.Metric, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if metricType == metric.Gauge {
		met, ok := m.gauges[metricName]
		return met, ok
	}
	if metricType == metric.Counter {
		met, ok := m.counters[metricName]
		return met, ok
	}

	return nil, false
}

// Update updates a metric if exists.
func (m *Memstorage) Update(_ context.Context, updatedMetric *metric.Metric) (*metric.Metric, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if updatedMetric.Type() == metric.Gauge {
		m.gauges[updatedMetric.Name()] = updatedMetric
		updatedMetric.FillHash(m.hashKey)
		return updatedMetric, nil
	}

	if updatedMetric.Type() == metric.Counter {
		res, ok := m.counters[updatedMetric.Name()]
		if !ok {
			m.counters[updatedMetric.Name()] = updatedMetric
			updatedMetric.FillHash(m.hashKey)
			return updatedMetric, nil
		}

		res = metric.Update(res, updatedMetric)
		res.FillHash(m.hashKey)
		m.counters[updatedMetric.Name()] = res
		return res, nil
	}

	return nil, errors.New("unknown metric type")
}

// UpdateAll updates all the metrics in list.
func (m *Memstorage) UpdateAll(ctx context.Context, metrics []*metric.Metric) error {
	for _, met := range metrics {
		_, _ = m.Update(ctx, met)
	}

	return nil
}

// Replace replaces metric with another one.
func (m *Memstorage) Replace(_ context.Context, met *metric.Metric) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if met.Type() == metric.Gauge {
		m.gauges[met.Name()] = met
	} else {
		m.counters[met.Name()] = met
	}
	met.FillHash(m.hashKey)
}

// All returns all the metrics.
func (m *Memstorage) All(_ context.Context) []*metric.Metric {
	m.mx.Lock()
	defer m.mx.Unlock()

	res := make([]*metric.Metric, 0, len(m.gauges)+len(m.counters))
	for _, c := range m.counters {
		res = append(res, c)
	}
	for _, g := range m.gauges {
		res = append(res, g)
	}

	return res
}

func (m *Memstorage) Close(_ context.Context) error {
	// For memory storage there is no need to do anything on close
	return nil
}

func (m *Memstorage) Ping(_ context.Context) error {
	// For memory storage there is no need to do anything on ping
	return nil
}
