package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/denistakeda/alerting/internal/config/agentcfg"
	"github.com/denistakeda/alerting/internal/grpcclient"
	"github.com/denistakeda/alerting/internal/httpclient"
	"github.com/denistakeda/alerting/internal/ports"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/v3/cpu"

	"github.com/shirou/gopsutil/v3/mem"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printInfo()

	conf, err := agentcfg.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	logService := loggerservice.New()
	logger := logService.ComponentLogger("Agent")

	logger.Info().Msgf("configuration: %v", conf)

	memStorage := memstorage.NewMemStorage(conf.Key, logService)

	client, err := makeClient(conf)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to initiate a client")
	}

	defer client.Stop()

	go readStats(conf.PollInterval, memStorage, logger)
	go sendStats(client, conf.ReportInterval, logger, memStorage)

	<-handleInterrupt()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.SendMetrics(memStorage.All(ctx)); err != nil {
		logger.Fatal().Err(err).Msg("unable to send metrics before stop")
	}
}

func makeClient(conf agentcfg.Config) (ports.Client, error) {
	if conf.GRPCAddress == "" {
		return httpclient.New(conf.RateLimit, conf.CryptoKey, conf.Address)
	} else {
		return grpcclient.NewGRPCClient(conf.GRPCAddress)
	}
}

func printInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func handleInterrupt() <-chan os.Signal {
	out := make(chan os.Signal, 2)
	signal.Notify(out, os.Interrupt)
	signal.Notify(out, syscall.SIGTERM)
	return out
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
	client ports.Client,
	reportInterval time.Duration,
	logger zerolog.Logger,
	store storage.Storage,
) {
	// Task publisher
	reportTicker := time.NewTicker(reportInterval)
	for range reportTicker.C {
		metrics := store.All(context.Background())
		if err := client.SendMetrics(metrics); err != nil {
			logger.Error().Err(err).Msg("failed to send metrics")
			continue
		}
		logger.Info().Msgf("successfully sent %d metrics", len(metrics))
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
