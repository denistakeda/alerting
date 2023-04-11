package memstorage

import (
	"context"
	"testing"

	"github.com/denistakeda/alerting/internal/services/loggerservice"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/storage"
)

func Test_memstorage_ImplementsStorage(t *testing.T) {
	var _ storage.Storage = (*Memstorage)(nil)
}

func Test_memstorage_Get(t *testing.T) {
	m1 := metric.NewGauge("m1_name", 3.14)
	type args struct {
		metricType metric.Type
		metricName string
	}
	type want struct {
		metric *metric.Metric
		ok     bool
	}
	tests := []struct {
		name    string
		metrics []*metric.Metric
		args    args
		want    want
	}{
		{
			name:    "should return existing value",
			metrics: []*metric.Metric{m1},
			args: args{
				metricType: m1.Type(),
				metricName: m1.Name(),
			},
			want: want{
				metric: m1,
				ok:     true,
			},
		},
		{
			name:    "should return nil for non-existing value",
			metrics: []*metric.Metric{},
			args: args{
				metricType: m1.Type(),
				metricName: m1.Name(),
			},
			want: want{
				metric: nil,
				ok:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := create(t, tt.metrics)
			got, ok := m.Get(context.Background(), tt.args.metricType, tt.args.metricName)
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.metric, got)
		})
	}
}

func Test_memstorage_Update(t *testing.T) {
	m1 := metric.NewGauge("m1_name", 3.14)
	m2 := metric.NewGauge("m2_name", 5.16)
	type args struct {
		updatedMetric *metric.Metric
	}
	tests := []struct {
		name    string
		metrics []*metric.Metric
		args    args
		wantErr bool
	}{
		{
			name:    "should return put if not present",
			metrics: []*metric.Metric{},
			args: args{
				updatedMetric: m1,
			},
			wantErr: false,
		},
		{
			name:    "should return update if present",
			metrics: []*metric.Metric{m1, m2},
			args: args{
				updatedMetric: m1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := create(t, tt.metrics)
			if _, err := m.Update(context.Background(), tt.args.updatedMetric); (err != nil) != tt.wantErr {
				t.Errorf("Memstorage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func create(t *testing.T, metrics []*metric.Metric) *Memstorage {
	ms := NewMemStorage("", loggerservice.New())
	for _, m := range metrics {
		_, err := ms.Update(context.Background(), m)
		require.NoError(t, err)
	}
	return ms
}
