package storage

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	"github.com/denistakeda/alerting/internal/storage/filestorage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/stretchr/testify/require"
)

func BenchmarkStorages(b *testing.B) {
	metricsCount := 1000
	logService := loggerservice.New()
	memStore := memstorage.NewMemStorage("", logService)
	fileStore, err := filestorage.NewFileStorage(
		context.Background(),
		"/tmp/store",
		500*time.Millisecond,
		false, "",
		logService,
	)
	require.NoError(b, err)

	metrics := generateMetrics(metricsCount)
	memStore.UpdateAll(context.Background(), metrics)
	fileStore.UpdateAll(context.Background(), metrics)

	b.Run("memstorage.Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			m := metrics[rand.Intn(metricsCount)]
			b.StartTimer()
			memStore.Get(context.Background(), m.Type(), m.Name())
		}
	})
	b.Run("filestorage.Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			m := metrics[rand.Intn(metricsCount)]
			b.StartTimer()
			fileStore.Get(context.Background(), m.Type(), m.Name())
		}
	})
	b.Run("memstorage.Update", func(b *testing.B) {
		for i := 0; i < metricsCount; i++ {
			b.StopTimer()
			m := metrics[rand.Intn(metricsCount)]
			if m.Type() == metric.Gauge {
				*m.Value += float64(1)
			} else {
				*m.Delta += int64(1)
			}
			b.StartTimer()
			memStore.Update(context.Background(), m)
		}
	})
	b.Run("filestorage.Update", func(b *testing.B) {
		for i := 0; i < metricsCount; i++ {
			b.StopTimer()
			m := metrics[rand.Intn(metricsCount)]
			if m.Type() == metric.Gauge {
				*m.Value += float64(1)
			} else {
				*m.Delta += int64(1)
			}
			b.StartTimer()
			fileStore.Update(context.Background(), m)
		}
	})
}

func generateMetrics(count int) []*metric.Metric {
	metrics := make([]*metric.Metric, 0, 1000)
	for i := 0; i < count; i++ {
		var m *metric.Metric
		if i%2 == 0 {
			m = metric.NewGauge(fmt.Sprintf("gauge-%d", i), float64(i))
		} else {
			m = metric.NewCounter(fmt.Sprintf("counter-%d", i), int64(i))
		}
		metrics = append(metrics, m)
	}

	return metrics
}
