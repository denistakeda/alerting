package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/metric/counter"
	"github.com/denistakeda/alerting/internal/metric/gauge"
	"github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
)

const (
	PollInterval   = 2 * time.Second
	ReportInterval = 10 * time.Second
)

func main() {
	mem := &runtime.MemStats{}
	memStorage := memstorage.New()

	// Update metrics
	pollTicker := time.NewTicker(PollInterval)
	go func() {
		for range pollTicker.C {
			runtime.ReadMemStats(mem)
			registerMetrics(mem, memStorage)
		}
	}()

	// Send metrics
	reportTicker := time.NewTicker(ReportInterval)
	for range reportTicker.C {
		sendMetrics(memStorage.All(), "http://127.0.0.1:8080")
	}
}

func registerMetrics(memStats *runtime.MemStats, store storage.Storage) {
	registerMetric(store, gauge.New("Alloc", float64(memStats.Alloc)))
	registerMetric(store, gauge.New("BuckHashSys", float64(memStats.BuckHashSys)))
	registerMetric(store, gauge.New("Frees", float64(memStats.Frees)))
	registerMetric(store, gauge.New("GCCPUFraction", float64(memStats.GCCPUFraction)))
	registerMetric(store, gauge.New("GCSys", float64(memStats.GCSys)))
	registerMetric(store, gauge.New("HeapAlloc", float64(memStats.HeapAlloc)))
	registerMetric(store, gauge.New("HeapIdle", float64(memStats.HeapIdle)))
	registerMetric(store, gauge.New("HeapInuse", float64(memStats.HeapInuse)))
	registerMetric(store, gauge.New("HeapObjects", float64(memStats.HeapObjects)))
	registerMetric(store, gauge.New("HeapReleased", float64(memStats.HeapReleased)))
	registerMetric(store, gauge.New("HeapSys", float64(memStats.HeapSys)))
	registerMetric(store, gauge.New("LastGC", float64(memStats.LastGC)))
	registerMetric(store, gauge.New("Lookups", float64(memStats.Lookups)))
	registerMetric(store, gauge.New("MCacheInuse", float64(memStats.MCacheInuse)))
	registerMetric(store, gauge.New("MCacheSys", float64(memStats.MCacheSys)))
	registerMetric(store, gauge.New("MSpanInUse", float64(memStats.MSpanInuse)))
	registerMetric(store, gauge.New("MSpanSys", float64(memStats.MSpanSys)))
	registerMetric(store, gauge.New("Mallocs", float64(memStats.Mallocs)))
	registerMetric(store, gauge.New("NextGC", float64(memStats.NextGC)))
	registerMetric(store, gauge.New("NumForcedGC", float64(memStats.NumForcedGC)))
	registerMetric(store, gauge.New("NumGC", float64(memStats.NumGC)))
	registerMetric(store, gauge.New("OtherSys", float64(memStats.OtherSys)))
	registerMetric(store, gauge.New("PauseTotalNs", float64(memStats.PauseTotalNs)))
	registerMetric(store, gauge.New("StackInuse", float64(memStats.StackInuse)))
	registerMetric(store, gauge.New("StackSys", float64(memStats.StackSys)))
	registerMetric(store, gauge.New("Sys", float64(memStats.Sys)))
	registerMetric(store, gauge.New("TotalAlloc", float64(memStats.TotalAlloc)))

	registerMetric(store, counter.New("PollCount", 1))
	registerMetric(store, gauge.New("RandomValue", float64(rand.Int())))
}

func registerMetric(store storage.Storage, m metric.Metric) {
	err := store.Update(m)
	if err != nil {
		log.Printf("Failed to update metric %v\n", m)
	}
}

func sendMetrics(metrics []metric.Metric, server string) {
	for _, m := range metrics {
		met := m
		go func() {
			url := fmt.Sprintf("%s/update/%s/%s", server, met.Type(), met.Name())
			body := bytes.NewBuffer([]byte(met.StrValue()))
			resp, err := http.Post(url, "text/plain", body)
			if err != nil {
				fmt.Printf("Unable to file a request to URL: %s, error: %v\n", url, err)
				return
			}
			defer resp.Body.Close()
		}()
	}
}
