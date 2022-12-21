package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/metric/counter"
	"github.com/denistakeda/alerting/internal/metric/gauge"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct{}

func (m *mockStorage) Get(metricType string, metricName string) (metric.Metric, bool) {
	return nil, false
}
func (m *mockStorage) Update(metric metric.Metric) error { return nil }
func (m *mockStorage) All() []metric.Metric              { return []metric.Metric{} }

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
	m1 := gauge.New("gauge1", 3.14)
	m2 := gauge.New("gauge2", 5.18)
	m3 := counter.New("counter1", 7)
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
			request: fmt.Sprintf("/value/%s/%s", m1.Type(), m1.Name()),
			storage: createStorage(t, []metric.Metric{m1, m2, m3}),
			want: want{
				code: http.StatusOK,
				body: m1.StrValue(),
			},
		},
		{
			name:    "counter success case",
			request: fmt.Sprintf("/value/%s/%s", m3.Type(), m3.Name()),
			storage: createStorage(t, []metric.Metric{m1, m2, m3}),
			want: want{
				code: http.StatusOK,
				body: m3.StrValue(),
			},
		},
		{
			name:    "request not existing metric",
			request: fmt.Sprintf("/value/%s/%s", m2.Type(), m2.Name()),
			storage: createStorage(t, []metric.Metric{m1, m3}),
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

func createStorage(t *testing.T, metrics []metric.Metric) s.Storage {
	ms := memstorage.New()
	for _, m := range metrics {
		err := ms.Update(m)
		require.NoError(t, err)
	}
	return ms
}
