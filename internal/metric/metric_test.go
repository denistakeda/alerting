package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric_StrValue(t *testing.T) {
	type fields struct {
		metricType   MetricType
		name         string
		gaugeValue   float64
		counterValue int64
	}
	tests := []struct {
		name   string
		metric *Metric
		want   string
	}{
		{
			name:   "gauge",
			metric: NewGauge("gauge", 3.14159265),
			want:   "3.142",
		},
		{
			name:   "counter",
			metric: NewCounter("gauge", 5),
			want:   "5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.metric.StrValue())
		})
	}
}

func TestMetric_UpdateValue(t *testing.T) {
	g1, g2 := NewGauge("g1", 3.14), NewGauge("g2", 5.17)
	c1, c2 := NewCounter("c1", 5), NewCounter("c2", 7)

	tests := []struct {
		name string
		old  *Metric
		new  *Metric
		want *Metric
	}{
		{
			name: "gauge to gauge",
			old:  g1,
			new:  g2,
			want: g2,
		},
		{
			name: "counter to counter",
			old:  c1,
			new:  c2,
			want: NewCounter("c1", 5+7),
		},
		{
			name: "old is nil",
			old:  nil,
			new:  g2,
			want: g2,
		},
		{
			name: "new is nil",
			old:  g1,
			new:  nil,
			want: g1,
		},
		{
			name: "incompatible types",
			old:  g1,
			new:  c1,
			want: g1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Update(tt.old, tt.new))
		})
	}
}
