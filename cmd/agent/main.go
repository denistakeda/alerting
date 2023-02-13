package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/denistakeda/alerting/internal/config/agentcfg"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	conf, err := agentcfg.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	logService := loggerservice.New()
	logger := logService.ComponentLogger("Agent")

	logger.Info().Msgf("configuration: %v", conf)

	memStorage := memstorage.NewMemStorage(conf.Key, logService)

	go readStats(conf.PollInterval, memStorage, logger)
	sendStats(conf.ReportInterval, conf.RateLimit, logger, memStorage, conf.Address)
}

func readStats(pollInterval time.Duration, store storage.Storage, logger zerolog.Logger) {
	pollTicker := time.NewTicker(pollInterval)

	for range pollTicker.C {
		go func() {
			if err := registerRuntimeMetrics(store, logger); err != nil {
				logger.Error().Err(err).Msg("failed to register runtime metrics")
			}
		}()

		go func() {
			if err := registerGoOpsMetrics(store, logger); err != nil {
				logger.Error().Err(err).Msg("failed to register goops metrics")
			}
		}()
	}
}

func sendStats(
	reportInterval time.Duration,
	rateLimit int,
	logger zerolog.Logger,
	store storage.Storage,
	address string,
) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    20,
			MaxConnsPerHost: 20,
		},
	}

	bus := make(chan []*metric.Metric, 100)

	// Initiate workers
	for i := 0; i < rateLimit; i++ {
		go func(workerId int) {
			for metrics := range bus {
				if err := sendMetrics(client, metrics, address); err != nil {
					logger.Error().Err(err).Int("workerId", workerId).Msg("failed to send metrics")
					continue
				}

				logger.Info().Int("workerId", workerId).Msgf("successfully sent %d metrics", len(metrics))
			}
		}(i)
	}

	// Task publisher
	reportTicker := time.NewTicker(reportInterval)
	for range reportTicker.C {
		bus <- store.All(context.Background())
	}
}

func registerRuntimeMetrics(store storage.Storage, logger zerolog.Logger) error {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

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

	return nil
}

func registerGoOpsMetrics(store storage.Storage, logger zerolog.Logger) error {
	gopsutilMemory, err := mem.VirtualMemory()
	if err != nil {
		return errors.Wrap(err, "failed to read virtual memory stats")
	}

	registerMetric(store, logger, metric.NewGauge("TotalMemory", float64(gopsutilMemory.Total)))
	registerMetric(store, logger, metric.NewGauge("FreeMemory", float64(gopsutilMemory.Free)))

	cpus, err := cpu.Percent(0, true)
	if err != nil {
		return errors.Wrap(err, "failed to read get the number of cores")
	}

	for idx, cpuUsage := range cpus {
		registerMetric(store, logger, metric.NewGauge(fmt.Sprintf("CPUutilization%d", idx), cpuUsage))
	}

	return nil
}

func registerMetric(store storage.Storage, logger zerolog.Logger, m *metric.Metric) {
	_, err := store.Update(context.Background(), m)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to update metric %v\n", m)
	}
}

func sendMetrics(client *http.Client, metrics []*metric.Metric, server string) error {
	url := fmt.Sprintf("%s/updates/", server)
	m, err := json.Marshal(metrics)
	if err != nil {
		return errors.Wrap(err, "failed to marshal metrics")
	}
	body := bytes.NewBuffer(m)
	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		return errors.Wrapf(err, "unable to file a request to URL: %s", url)
	}
	if err := resp.Body.Close(); err != nil {
		return errors.Wrap(err, "unable to close a body")
	}

	return nil
}
