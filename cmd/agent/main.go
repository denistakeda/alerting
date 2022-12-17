package main

import (
	"bytes"
	"fmt"
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
	store.Update(gauge.New("Alloc", float64(memStats.Alloc)))
	store.Update(gauge.New("BuckHashSys", float64(memStats.BuckHashSys)))
	store.Update(gauge.New("Frees", float64(memStats.Frees)))
	store.Update(gauge.New("GCCPUFraction", float64(memStats.GCCPUFraction)))
	store.Update(gauge.New("GCSys", float64(memStats.GCSys)))
	store.Update(gauge.New("HeapAlloc", float64(memStats.HeapAlloc)))
	store.Update(gauge.New("HeapIdle", float64(memStats.HeapIdle)))
	store.Update(gauge.New("HeapInuse", float64(memStats.HeapInuse)))
	store.Update(gauge.New("HeapObjects", float64(memStats.HeapObjects)))
	store.Update(gauge.New("HeapReleased", float64(memStats.HeapReleased)))
	store.Update(gauge.New("HeapSys", float64(memStats.HeapSys)))
	store.Update(gauge.New("LastGC", float64(memStats.LastGC)))
	store.Update(gauge.New("Lookups", float64(memStats.Lookups)))
	store.Update(gauge.New("MCacheInuse", float64(memStats.MCacheInuse)))
	store.Update(gauge.New("MCacheSys", float64(memStats.MCacheSys)))
	store.Update(gauge.New("MSpanInUse", float64(memStats.MSpanInuse)))
	store.Update(gauge.New("MSpanSys", float64(memStats.MSpanSys)))
	store.Update(gauge.New("Mallocs", float64(memStats.Mallocs)))
	store.Update(gauge.New("NextGC", float64(memStats.NextGC)))
	store.Update(gauge.New("NumForcedGC", float64(memStats.NumForcedGC)))
	store.Update(gauge.New("NumGC", float64(memStats.NumGC)))
	store.Update(gauge.New("OtherSys", float64(memStats.OtherSys)))
	store.Update(gauge.New("PauseTotalNs", float64(memStats.PauseTotalNs)))
	store.Update(gauge.New("StackInuse", float64(memStats.StackInuse)))
	store.Update(gauge.New("StackSys", float64(memStats.StackSys)))
	store.Update(gauge.New("Sys", float64(memStats.Sys)))
	store.Update(gauge.New("TotalAlloc", float64(memStats.TotalAlloc)))

	store.Update(counter.New("PollCount", 1))
	store.Update(gauge.New("RandomValue", float64(rand.Int())))
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
