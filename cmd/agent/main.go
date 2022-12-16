package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/denistakeda/alerting/internal/memstorage"
	"github.com/denistakeda/alerting/internal/metric"
)

const (
	PollInterval   = 2 * time.Second
	ReportInterval = 10 * time.Second
)

func main() {
	mem := &runtime.MemStats{}
	memStorage := memstorage.NewMemStorage()

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
		sendMetrics(memStorage.Metrics(), "http://127.0.0.1:8080")
	}
}

func registerMetrics(memStats *runtime.MemStats, memStorage *memstorage.MemStorage) {
	memStorage.StoreGauge("Alloc", float64(memStats.Alloc))
	memStorage.StoreGauge("BuckHashSys", float64(memStats.BuckHashSys))
	memStorage.StoreGauge("Frees", float64(memStats.Frees))
	memStorage.StoreGauge("GCCPUFraction", float64(memStats.GCCPUFraction))
	memStorage.StoreGauge("GCSys", float64(memStats.GCSys))
	memStorage.StoreGauge("HeapAlloc", float64(memStats.HeapAlloc))
	memStorage.StoreGauge("HeapIdle", float64(memStats.HeapIdle))
	memStorage.StoreGauge("HeapInuse", float64(memStats.HeapInuse))
	memStorage.StoreGauge("HeapObjects", float64(memStats.HeapObjects))
	memStorage.StoreGauge("HeapReleased", float64(memStats.HeapReleased))
	memStorage.StoreGauge("HeapSys", float64(memStats.HeapSys))
	memStorage.StoreGauge("LastGC", float64(memStats.LastGC))
	memStorage.StoreGauge("Lookups", float64(memStats.Lookups))
	memStorage.StoreGauge("MCacheInuse", float64(memStats.MCacheInuse))
	memStorage.StoreGauge("MCacheSys", float64(memStats.MCacheSys))
	memStorage.StoreGauge("MSpanInUse", float64(memStats.MSpanInuse))
	memStorage.StoreGauge("MSpanSys", float64(memStats.MSpanSys))
	memStorage.StoreGauge("Mallocs", float64(memStats.Mallocs))
	memStorage.StoreGauge("NextGC", float64(memStats.NextGC))
	memStorage.StoreGauge("NumForcedGC", float64(memStats.NumForcedGC))
	memStorage.StoreGauge("NumGC", float64(memStats.NumGC))
	memStorage.StoreGauge("OtherSys", float64(memStats.OtherSys))
	memStorage.StoreGauge("PauseTotalNs", float64(memStats.PauseTotalNs))
	memStorage.StoreGauge("StackInuse", float64(memStats.StackInuse))
	memStorage.StoreGauge("StackSys", float64(memStats.StackSys))
	memStorage.StoreGauge("Sys", float64(memStats.Sys))
	memStorage.StoreGauge("TotalAlloc", float64(memStats.TotalAlloc))

	memStorage.StoreCounter("PollCount", 1)
	memStorage.StoreGauge("RandomValue", float64(rand.Int()))
}

func sendMetrics(metrics []metric.Metric, server string) {
	for _, m := range metrics {
		met := m
		go func() {
			url := fmt.Sprintf("%s/update/%s/%s", server, met.Type, met.Name)
			body := bytes.NewBuffer([]byte(met.Value))
			resp, err := http.Post(url, "text/plain", body)
			if err != nil {
				fmt.Printf("Unable to file a request to URL: %s, error: %v\n", url, err)
				return
			}
			defer resp.Body.Close()
		}()
	}
}
