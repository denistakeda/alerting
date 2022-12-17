package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denistakeda/alerting/internal/metric"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct{}

func (m *mockStorage) Get(metricType string, metricName string) (metric.Metric, bool) {
	return nil, false
}
func (m *mockStorage) Update(metric metric.Metric) error { return nil }
func (m *mockStorage) All() []metric.Metric              { return []metric.Metric{} }

func Test_setupRouter(t *testing.T) {
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
