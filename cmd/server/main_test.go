package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denistakeda/alerting/internal/metric"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct{}

func (m *mockStorage) Get(_metricType metric.Type, _metricName string) (*metric.Metric, bool) {
	return nil, false
}
func (m *mockStorage) Update(_metric *metric.Metric) (*metric.Metric, error) { return nil, nil }
func (m *mockStorage) All() []*metric.Metric                                 { return []*metric.Metric{} }

func Test_updateMetric(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		storage  s.Storage
		wantCode int
	}{
		{
			name:     "gauge success case",
			request:  "/update/gauge/metric_name/100",
			storage:  &mockStorage{},
			wantCode: http.StatusOK,
		},
		{
			name:     "counter success case",
			request:  "/update/counter/metric_name/100",
			storage:  &mockStorage{},
			wantCode: http.StatusOK,
		},
		{
			name:     "gauge without name and type",
			request:  "/update/gauge/",
			storage:  &mockStorage{},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "counter without name and type",
			request:  "/update/counter/",
			storage:  &mockStorage{},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "gauge invalid value",
			request:  "/update/gauge/test_counter/none",
			storage:  &mockStorage{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "counter invalid value",
			request:  "/update/counter/test_counter/none",
			storage:  &mockStorage{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "unknown type",
			request:  "/update/unknown/testCounter/100",
			storage:  &mockStorage{},
			wantCode: http.StatusNotImplemented,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.storage)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", tt.request, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func Test_getMetric(t *testing.T) {
	m1 := metric.NewGauge("gauge1", 3.14)
	m2 := metric.NewGauge("gauge2", 5.18)
	m3 := metric.NewCounter("counter1", 7)
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name    string
		request string
		storage s.Storage
		want    want
	}{
		{
			name:    "gauge success case",
			request: fmt.Sprintf("/value/%s/%s", m1.StrType(), m1.Name()),
			storage: createStorage(t, []*metric.Metric{m1, m2, m3}),
			want: want{
				code: http.StatusOK,
				body: m1.StrValue(),
			},
		},
		{
			name:    "counter success case",
			request: fmt.Sprintf("/value/%s/%s", m3.StrType(), m3.Name()),
			storage: createStorage(t, []*metric.Metric{m1, m2, m3}),
			want: want{
				code: http.StatusOK,
				body: m3.StrValue(),
			},
		},
		{
			name:    "request not existing metric",
			request: fmt.Sprintf("/value/%s/%s", m2.StrType(), m2.Name()),
			storage: createStorage(t, []*metric.Metric{m1, m3}),
			want: want{
				code: http.StatusNotFound,
				body: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.storage)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.request, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func Test_update(t *testing.T) {

	g1 := metric.NewGauge("gauge1", 3.14)
	//g2 := metric.NewGauge("gauge2", 5.18)
	c3 := metric.NewCounter("counter1", 7)

	type want struct {
		code int
		body []byte
	}
	tests := []struct {
		name        string
		requestBody []byte
		storage     s.Storage
		want        want
	}{
		{
			name:        "not existing gauge",
			requestBody: marshal(t, g1),
			storage:     createStorage(t, make([]*metric.Metric, 0)),
			want: want{
				code: http.StatusOK,
				body: marshal(t, g1),
			},
		},
		{
			name:        "existing gauge",
			requestBody: marshal(t, metric.NewGauge(g1.Name(), 5.18)),
			storage:     createStorage(t, []*metric.Metric{g1}),
			want: want{
				code: http.StatusOK,
				body: marshal(t, metric.NewGauge(g1.Name(), 5.18)),
			},
		},
		{
			name:        "not existing counter",
			requestBody: marshal(t, c3),
			storage:     createStorage(t, make([]*metric.Metric, 0)),
			want: want{
				code: http.StatusOK,
				body: marshal(t, c3),
			},
		},
		{
			name:        "existing counter",
			requestBody: marshal(t, c3),
			storage:     createStorage(t, []*metric.Metric{c3}),
			want: want{
				code: http.StatusOK,
				body: marshal(t, metric.NewCounter(c3.Name(), 14)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.storage)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/update/", bytes.NewBuffer(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.JSONEq(t, string(tt.want.body), w.Body.String())
		})
	}
}

func marshal(t *testing.T, v any) []byte {
	res, err := json.Marshal(v)
	require.NoError(t, err)
	return res
}

func createStorage(t *testing.T, metrics []*metric.Metric) s.Storage {
	ms := memstorage.New()
	for _, m := range metrics {
		_, err := ms.Update(m)
		require.NoError(t, err)
	}
	return ms
}
