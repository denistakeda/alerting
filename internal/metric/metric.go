package metrick

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
)

type metricType string

const (
	Gauge   metricType = "gauge"
	Counter metricType = "counter"
)

type metricName string

const (
	Alloc         metricName = "Alloc"
	BuckHashSys   metricName = "BuckHashSys"
	Frees         metricName = "Frees"
	GCCPUFraction metricName = "GCCPUFraction"
	GCSys         metricName = "GCSys"
	HeapAlloc     metricName = "HeapAlloc"
	HeapIdle      metricName = "HeapIdle"
	HeapInuse     metricName = "HeapInuse"
	HeapObjects   metricName = "HeapObjects"
	HeapReleased  metricName = "HeapReleased"
	HeapSys       metricName = "HeapSys"
	LastGC        metricName = "LastGC"
	Lookups       metricName = "Lookups"
	MCacheInuse   metricName = "MCacheInuse"
	MCacheSys     metricName = "MCacheSys"
	MSpanInuse    metricName = "MSpanInuse"
	MSpanSys      metricName = "MSpanSys"
	Mallocs       metricName = "Mallocs"
	NextGC        metricName = "NextGC"
	NumForcedGC   metricName = "NumForcedGC"
	NumGC         metricName = "NumGC"
	OtherSys      metricName = "OtherSys"
	PauseTotalNs  metricName = "PauseTotalNs"
	StackInuse    metricName = "StackInuse"
	StackSys      metricName = "StackSys"
	Sys           metricName = "Sys"
	TotalAlloc    metricName = "TotalAlloc"
	PollCount     metricName = "PollCount"
	RandomValue   metricName = "RandomValue"
)

type metric struct {
	Type  metricType
	Value string
}

type metrics struct {
	ms map[metricName]*metric
}

func NewMetrics() *metrics {
	return &metrics{
		ms: map[metricName]*metric{
			Alloc:         {Gauge, ""},
			BuckHashSys:   {Gauge, ""},
			Frees:         {Gauge, ""},
			GCCPUFraction: {Gauge, ""},
			GCSys:         {Gauge, ""},
			HeapAlloc:     {Gauge, ""},
			HeapIdle:      {Gauge, ""},
			HeapInuse:     {Gauge, ""},
			HeapObjects:   {Gauge, ""},
			HeapReleased:  {Gauge, ""},
			HeapSys:       {Gauge, ""},
			LastGC:        {Gauge, ""},
			Lookups:       {Gauge, ""},
			MCacheInuse:   {Gauge, ""},
			MCacheSys:     {Gauge, ""},
			MSpanInuse:    {Gauge, ""},
			MSpanSys:      {Gauge, ""},
			Mallocs:       {Gauge, ""},
			NextGC:        {Gauge, ""},
			NumForcedGC:   {Gauge, ""},
			NumGC:         {Gauge, ""},
			OtherSys:      {Gauge, ""},
			PauseTotalNs:  {Gauge, ""},
			StackInuse:    {Gauge, ""},
			StackSys:      {Gauge, ""},
			Sys:           {Gauge, ""},
			TotalAlloc:    {Gauge, ""},
			PollCount:     {Counter, "0"},
			RandomValue:   {Gauge, ""},
		}}
}

func (m *metrics) Fill(stats *runtime.MemStats) {
	m.setMetric(Alloc, strconv.FormatUint(stats.Alloc, 10))
	m.setMetric(BuckHashSys, strconv.FormatUint(stats.BuckHashSys, 10))
	m.setMetric(Frees, strconv.FormatUint(stats.Frees, 10))
	m.setMetric(GCCPUFraction, strconv.FormatUint(uint64(stats.GCCPUFraction), 10))
	m.setMetric(GCSys, strconv.FormatUint(stats.GCSys, 10))
	m.setMetric(HeapAlloc, strconv.FormatUint(stats.HeapAlloc, 10))
	m.setMetric(HeapIdle, strconv.FormatUint(stats.HeapIdle, 10))
	m.setMetric(HeapInuse, strconv.FormatUint(stats.HeapInuse, 10))
	m.setMetric(HeapObjects, strconv.FormatUint(stats.HeapObjects, 10))
	m.setMetric(HeapReleased, strconv.FormatUint(stats.HeapReleased, 10))
	m.setMetric(HeapSys, strconv.FormatUint(stats.HeapSys, 10))
	m.setMetric(LastGC, strconv.FormatUint(stats.LastGC, 10))
	m.setMetric(Lookups, strconv.FormatUint(stats.Lookups, 10))
	m.setMetric(MCacheInuse, strconv.FormatUint(stats.MCacheInuse, 10))
	m.setMetric(MCacheSys, strconv.FormatUint(stats.MCacheSys, 10))
	m.setMetric(MSpanInuse, strconv.FormatUint(stats.MSpanInuse, 10))
	m.setMetric(MSpanSys, strconv.FormatUint(stats.MSpanSys, 10))
	m.setMetric(Mallocs, strconv.FormatUint(stats.Mallocs, 10))
	m.setMetric(NextGC, strconv.FormatUint(stats.NextGC, 10))
	m.setMetric(NumForcedGC, strconv.FormatUint(uint64(stats.NumForcedGC), 10))
	m.setMetric(NumGC, strconv.FormatUint(uint64(stats.NumGC), 10))
	m.setMetric(OtherSys, strconv.FormatUint(stats.OtherSys, 10))
	m.setMetric(PauseTotalNs, strconv.FormatUint(stats.PauseTotalNs, 10))
	m.setMetric(StackInuse, strconv.FormatUint(stats.StackInuse, 10))
	m.setMetric(StackSys, strconv.FormatUint(stats.StackSys, 10))
	m.setMetric(Sys, strconv.FormatUint(stats.Sys, 10))
	m.setMetric(TotalAlloc, strconv.FormatUint(stats.TotalAlloc, 10))

	// Special cases
	oldPollCount, err := strconv.Atoi(m.ms[PollCount].Value)
	if err != nil {
		oldPollCount = 0
	}
	m.setMetric(PollCount, strconv.Itoa(oldPollCount+1))
	m.setMetric(RandomValue, strconv.Itoa(rand.Int()))
}

func (m *metrics) Send(url string) {
	for n, m := range m.ms {
		name, metric := n, m
		go func() {
			url := fmt.Sprintf("%s/update/%s/%s/%s", url, metric.Type, name, metric.Value)
			resp, err := http.Post(url, "text/plain", nil)

			if err == nil {
				resp.Body.Close()
			} else {
				fmt.Printf("Unable to file a request to URL: %s, error: %v\n", url, err)
			}
		}()
	}
}

func (m *metrics) ToString() string {
	b := new(bytes.Buffer)
	for key, value := range m.ms {
		fmt.Fprintf(b, "%s=%v\n", key, value)
	}
	return b.String()
}

func (m *metrics) setMetric(name metricName, val string) {
	metric, ok := m.ms[name]
	if ok {
		metric.Value = val
	}
}
