package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/denistakeda/alerting/internal/config/agentcfg"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
)

func main() {
	conf, err := agentcfg.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("configuration: %v", conf)
	mem := &runtime.MemStats{}
	memStorage := memstorage.New()

	// Update metrics
	pollTicker := time.NewTicker(conf.PollInterval)
	go func() {
		for range pollTicker.C {
			runtime.ReadMemStats(mem)
			registerMetrics(mem, memStorage)
		}
	}()

	// Send metrics
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    20,
			MaxConnsPerHost: 20,
		},
	}
	reportTicker := time.NewTicker(conf.ReportInterval)
	for range reportTicker.C {
		sendMetrics(client, memStorage.All(), conf.Address)
	}
}

func registerMetrics(memStats *runtime.MemStats, store storage.Storage) {
	registerMetric(store, metric.NewGauge("Alloc", float64(memStats.Alloc)))
	registerMetric(store, metric.NewGauge("BuckHashSys", float64(memStats.BuckHashSys)))
	registerMetric(store, metric.NewGauge("Frees", float64(memStats.Frees)))
	registerMetric(store, metric.NewGauge("GCCPUFraction", float64(memStats.GCCPUFraction)))
	registerMetric(store, metric.NewGauge("GCSys", float64(memStats.GCSys)))
	registerMetric(store, metric.NewGauge("HeapAlloc", float64(memStats.HeapAlloc)))
	registerMetric(store, metric.NewGauge("HeapIdle", float64(memStats.HeapIdle)))
	registerMetric(store, metric.NewGauge("HeapInuse", float64(memStats.HeapInuse)))
	registerMetric(store, metric.NewGauge("HeapObjects", float64(memStats.HeapObjects)))
	registerMetric(store, metric.NewGauge("HeapReleased", float64(memStats.HeapReleased)))
	registerMetric(store, metric.NewGauge("HeapSys", float64(memStats.HeapSys)))
	registerMetric(store, metric.NewGauge("LastGC", float64(memStats.LastGC)))
	registerMetric(store, metric.NewGauge("Lookups", float64(memStats.Lookups)))
	registerMetric(store, metric.NewGauge("MCacheInuse", float64(memStats.MCacheInuse)))
	registerMetric(store, metric.NewGauge("MCacheSys", float64(memStats.MCacheSys)))
	registerMetric(store, metric.NewGauge("MSpanInuse", float64(memStats.MSpanInuse)))
	registerMetric(store, metric.NewGauge("MSpanSys", float64(memStats.MSpanSys)))
	registerMetric(store, metric.NewGauge("Mallocs", float64(memStats.Mallocs)))
	registerMetric(store, metric.NewGauge("NextGC", float64(memStats.NextGC)))
	registerMetric(store, metric.NewGauge("NumForcedGC", float64(memStats.NumForcedGC)))
	registerMetric(store, metric.NewGauge("NumGC", float64(memStats.NumGC)))
	registerMetric(store, metric.NewGauge("OtherSys", float64(memStats.OtherSys)))
	registerMetric(store, metric.NewGauge("PauseTotalNs", float64(memStats.PauseTotalNs)))
	registerMetric(store, metric.NewGauge("StackInuse", float64(memStats.StackInuse)))
	registerMetric(store, metric.NewGauge("StackSys", float64(memStats.StackSys)))
	registerMetric(store, metric.NewGauge("Sys", float64(memStats.Sys)))
	registerMetric(store, metric.NewGauge("TotalAlloc", float64(memStats.TotalAlloc)))

	registerMetric(store, metric.NewCounter("PollCount", 1))
	registerMetric(store, metric.NewGauge("RandomValue", float64(rand.Int())))
}

func registerMetric(store storage.Storage, m *metric.Metric) {
	_, err := store.Update(m)
	if err != nil {
		log.Printf("Failed to update metric %v\n", m)
	}
}

func sendMetrics(client *http.Client, metrics []*metric.Metric, server string) {
	startTime := time.Now()
	for _, m := range metrics {
		m := m
		func() {
			url := fmt.Sprintf("%s/update/", server)
			m, err := json.Marshal(m)
			if err != nil {
				log.Printf("failed to marshal metric: %v\n", m)
			}
			body := bytes.NewBuffer(m)
			resp, err := client.Post(url, "application/json", body)
			if err != nil {
				log.Printf("unable to file a request to URL: %s, error: %v, metric: %v\n", url, err, string(m))
				return
			}
			if err := resp.Body.Close(); err != nil {
				log.Print("unable to close a body")
				return
			}
		}()
	}
	now := time.Now()
	log.Printf("successfully updated %d metrics in %f seconds", len(metrics), now.Sub(startTime).Seconds())
}
