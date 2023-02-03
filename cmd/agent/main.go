package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/denistakeda/alerting/internal/config/agentcfg"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	"github.com/rs/zerolog"
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

	logService := loggerservice.New()
	logger := logService.ComponentLogger("Agent")

	logger.Info().Msgf("configuration: %v", conf)

	mem := &runtime.MemStats{}
	memStorage := memstorage.New(conf.Key, logService)

	// Update metrics
	pollTicker := time.NewTicker(conf.PollInterval)
	go func() {
		for range pollTicker.C {
			runtime.ReadMemStats(mem)
			registerMetrics(mem, memStorage, logger)
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
		sendMetrics(client, logger, memStorage.All(context.Background()), conf.Address)
	}
}

func registerMetrics(memStats *runtime.MemStats, store storage.Storage, logger zerolog.Logger) {
	registerMetric(store, logger, metric.NewGauge("Alloc", float64(memStats.Alloc)))
	registerMetric(store, logger, metric.NewGauge("BuckHashSys", float64(memStats.BuckHashSys)))
	registerMetric(store, logger, metric.NewGauge("Frees", float64(memStats.Frees)))
	registerMetric(store, logger, metric.NewGauge("GCCPUFraction", float64(memStats.GCCPUFraction)))
	registerMetric(store, logger, metric.NewGauge("GCSys", float64(memStats.GCSys)))
	registerMetric(store, logger, metric.NewGauge("HeapAlloc", float64(memStats.HeapAlloc)))
	registerMetric(store, logger, metric.NewGauge("HeapIdle", float64(memStats.HeapIdle)))
	registerMetric(store, logger, metric.NewGauge("HeapInuse", float64(memStats.HeapInuse)))
	registerMetric(store, logger, metric.NewGauge("HeapObjects", float64(memStats.HeapObjects)))
	registerMetric(store, logger, metric.NewGauge("HeapReleased", float64(memStats.HeapReleased)))
	registerMetric(store, logger, metric.NewGauge("HeapSys", float64(memStats.HeapSys)))
	registerMetric(store, logger, metric.NewGauge("LastGC", float64(memStats.LastGC)))
	registerMetric(store, logger, metric.NewGauge("Lookups", float64(memStats.Lookups)))
	registerMetric(store, logger, metric.NewGauge("MCacheInuse", float64(memStats.MCacheInuse)))
	registerMetric(store, logger, metric.NewGauge("MCacheSys", float64(memStats.MCacheSys)))
	registerMetric(store, logger, metric.NewGauge("MSpanInuse", float64(memStats.MSpanInuse)))
	registerMetric(store, logger, metric.NewGauge("MSpanSys", float64(memStats.MSpanSys)))
	registerMetric(store, logger, metric.NewGauge("Mallocs", float64(memStats.Mallocs)))
	registerMetric(store, logger, metric.NewGauge("NextGC", float64(memStats.NextGC)))
	registerMetric(store, logger, metric.NewGauge("NumForcedGC", float64(memStats.NumForcedGC)))
	registerMetric(store, logger, metric.NewGauge("NumGC", float64(memStats.NumGC)))
	registerMetric(store, logger, metric.NewGauge("OtherSys", float64(memStats.OtherSys)))
	registerMetric(store, logger, metric.NewGauge("PauseTotalNs", float64(memStats.PauseTotalNs)))
	registerMetric(store, logger, metric.NewGauge("StackInuse", float64(memStats.StackInuse)))
	registerMetric(store, logger, metric.NewGauge("StackSys", float64(memStats.StackSys)))
	registerMetric(store, logger, metric.NewGauge("Sys", float64(memStats.Sys)))
	registerMetric(store, logger, metric.NewGauge("TotalAlloc", float64(memStats.TotalAlloc)))

	registerMetric(store, logger, metric.NewCounter("PollCount", 1))
	registerMetric(store, logger, metric.NewGauge("RandomValue", float64(rand.Int())))
}

func registerMetric(store storage.Storage, logger zerolog.Logger, m *metric.Metric) {
	_, err := store.Update(context.Background(), m)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to update metric %v\n", m)
	}
}

func sendMetrics(client *http.Client, logger zerolog.Logger, metrics []*metric.Metric, server string) {
	startTime := time.Now()

	url := fmt.Sprintf("%s/updates/", server)
	m, err := json.Marshal(metrics)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal metrics")
	}
	body := bytes.NewBuffer(m)
	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		logger.Error().Err(err).Msgf("unable to file a request to URL: %s, error: %v, metric: %v\n", url, err, string(m))
		return
	}
	if err := resp.Body.Close(); err != nil {
		logger.Error().Err(err).Msg("unable to close a body")
		return
	}
	now := time.Now()
	logger.Info().Msgf("successfully updated %d metrics in %f seconds", len(metrics), now.Sub(startTime).Seconds())
}
