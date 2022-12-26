package memstorage

import (
	"github.com/denistakeda/alerting/internal/metric"
)

type memstorage struct {
	types map[metric.MetricType]map[string]*metric.Metric
}

func New() *memstorage {
	return &memstorage{
		types: make(map[metric.MetricType]map[string]*metric.Metric),
	}
}

func (m *memstorage) Get(metricType metric.MetricType, metricName string) (*metric.Metric, bool) {
	group, ok := m.types[metricType]
	if !ok {
		return nil, false
	}

	metric, ok := group[metricName]
	return metric, ok
}

func (m *memstorage) Update(updatedMetric *metric.Metric) error {
	group, ok := m.types[updatedMetric.Type()]
	if !ok {
		group = make(map[string]*metric.Metric)
		m.types[updatedMetric.Type()] = group
	}

	group[updatedMetric.Name()] = metric.Update(group[updatedMetric.Name()], updatedMetric)
	return nil
}

func (m *memstorage) All() []*metric.Metric {
	res := []*metric.Metric{}
	for _, group := range m.types {
		for _, met := range group {
			res = append(res, met)
		}
	}
	return res
}
