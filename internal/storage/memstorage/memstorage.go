package memstorage

import (
	"github.com/denistakeda/alerting/internal/metric"
	"sync"
)

type Memstorage struct {
	types   map[metric.Type]map[string]*metric.Metric
	hashKey string
	mx      sync.Mutex
}

func New(hashKey string) *Memstorage {
	return &Memstorage{
		types:   make(map[metric.Type]map[string]*metric.Metric),
		hashKey: hashKey,
	}
}

func (m *Memstorage) Get(metricType metric.Type, metricName string) (*metric.Metric, bool) {
	m.mx.Lock()
	defer m.mx.Unlock()

	group, ok := m.types[metricType]
	if !ok {
		return nil, false
	}

	met, ok := group[metricName]
	return met, ok
}

func (m *Memstorage) Update(updatedMetric *metric.Metric) (*metric.Metric, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	group, ok := m.types[updatedMetric.Type()]
	if !ok {
		group = make(map[string]*metric.Metric)
		m.types[updatedMetric.Type()] = group
	}

	res := metric.Update(group[updatedMetric.Name()], updatedMetric, m.hashKey)
	group[updatedMetric.Name()] = res
	return res, nil
}

func (m *Memstorage) Replace(met *metric.Metric) {
	m.mx.Lock()
	defer m.mx.Unlock()

	group, ok := m.types[met.Type()]
	if !ok {
		group = make(map[string]*metric.Metric)
		m.types[met.Type()] = group
	}
	group[met.Name()] = met
}

func (m *Memstorage) All() []*metric.Metric {
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

func (m *Memstorage) Close() error {
	// For memory storage there is no need to do anything on close
	return nil
}
