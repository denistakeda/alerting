package memstorage

import (
	"github.com/denistakeda/alerting/internal/metric"
)

type Memstorage struct {
	types map[metric.Type]map[string]*metric.Metric
}

func New() *Memstorage {
	return &Memstorage{
		types: make(map[metric.Type]map[string]*metric.Metric),
	}
}

func (m *Memstorage) Get(metricType metric.Type, metricName string) (*metric.Metric, bool) {
	group, ok := m.types[metricType]
	if !ok {
		return nil, false
	}

	met, ok := group[metricName]
	return met, ok
}

func (m *Memstorage) Update(updatedMetric *metric.Metric) error {
	group, ok := m.types[updatedMetric.Type()]
	if !ok {
		group = make(map[string]*metric.Metric)
		m.types[updatedMetric.Type()] = group
	}

	group[updatedMetric.Name()] = metric.Update(group[updatedMetric.Name()], updatedMetric)
	return nil
}

func (m *Memstorage) All() []*metric.Metric {
	res := []*metric.Metric{}
	for _, group := range m.types {
		for _, met := range group {
			res = append(res, met)
		}
	}
	return res
}
