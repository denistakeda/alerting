package memstorage

import (
	"context"
	"sync"

	"github.com/denistakeda/alerting/internal/storage"
	"github.com/rs/zerolog"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
)

var _ storage.Storage = (*Memstorage)(nil)

// Memstorage is a memory storage.
type Memstorage struct {
	types   map[metric.Type]map[string]*metric.Metric
	hashKey string
	mx      sync.Mutex
	logger  zerolog.Logger
}

// NewMemStorage instantiates a new MemStorage instance.
func NewMemStorage(hashKey string, logService *loggerservice.LoggerService) *Memstorage {
	return &Memstorage{
		types:   make(map[metric.Type]map[string]*metric.Metric),
		hashKey: hashKey,
		logger:  logService.ComponentLogger("Memstorage"),
	}
}

// Get returns a metric if exists.
func (m *Memstorage) Get(_ context.Context, metricType metric.Type, metricName string) (*metric.Metric, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()

	group, ok := m.types[metricType]
	if !ok {
		return nil, false
	}

	met, ok := group[metricName]
	return met, ok
}

// Update updates a metric if exists.
func (m *Memstorage) Update(_ context.Context, updatedMetric *metric.Metric) (*metric.Metric, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	group, ok := m.types[updatedMetric.Type()]
	if !ok {
		group = make(map[string]*metric.Metric)
		m.types[updatedMetric.Type()] = group
	}

	res := metric.Update(group[updatedMetric.Name()], updatedMetric)
	group[updatedMetric.Name()] = res
	res.FillHash(m.hashKey)
	return res, nil
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

	group, ok := m.types[met.Type()]
	if !ok {
		group = make(map[string]*metric.Metric)
		m.types[met.Type()] = group
	}
	group[met.Name()] = met
	met.FillHash(m.hashKey)
}

// All returns all the metrics.
func (m *Memstorage) All(_ context.Context) []*metric.Metric {
	m.mx.Lock()
	defer m.mx.Unlock()

	var res []*metric.Metric
	for _, group := range m.types {
		for _, met := range group {
			res = append(res, met)
		}
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
