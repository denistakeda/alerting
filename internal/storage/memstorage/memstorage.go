package memstorage

import (
	"github.com/denistakeda/alerting/internal/metric"
)

type memstorage struct {
	types map[string]map[string]metric.Metric
	size  int
}

func New() *memstorage {
	return &memstorage{
		types: make(map[string]map[string]metric.Metric),
		size:  0,
	}
}

func (m *memstorage) Get(metricType string, metricName string) (metric.Metric, bool) {
	group, ok := m.types[metricType]
	if !ok {
		return nil, false
	}

	metric, ok := group[metricName]
	return metric, ok
}

func (m *memstorage) Update(updatedMetric metric.Metric) error {
	group, ok := m.types[updatedMetric.Type()]
	if !ok {
		group = make(map[string]metric.Metric)
		m.types[updatedMetric.Type()] = group
	}

	oldMetric, ok := group[updatedMetric.Name()]
	if !ok {
		group[updatedMetric.Name()] = updatedMetric
		m.size++
		return nil
	}
	return oldMetric.UpdateValue(updatedMetric.Value())
}

func (m *memstorage) All() []metric.Metric {
	res := make([]metric.Metric, m.size)
	i := 0
	for _, group := range m.types {
		for _, met := range group {
			res[i] = met
			i++
		}
	}
	return res
}
