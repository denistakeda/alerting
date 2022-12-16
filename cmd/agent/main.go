package main

import (
	"runtime"
	"time"

	metric "github.com/denistakeda/alerting/internal/metric"
)

const (
	PollInterval   = 2 * time.Second
	ReportInterval = 10 * time.Second
)

func main() {
	mem := &runtime.MemStats{}
	metrics := metric.NewMetrics()

	// Update metrics
	pollTicker := time.NewTicker(PollInterval)
	go func() {
		for range pollTicker.C {
			runtime.ReadMemStats(mem)
			metrics.Fill(mem)
		}
	}()

	// Send metrics
	reportTicker := time.NewTicker(ReportInterval)
	for range reportTicker.C {
		metrics.Send("http://127.0.0.1:8080")
	}
}
