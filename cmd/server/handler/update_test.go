package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct{}

func (m *mockStorage) StoreGauge(name string, value float64) {}
func (m *mockStorage) StoreCounter(name string, value int64) {}

func TestUpdateHandler(t *testing.T) {
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
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(UpdateHandler(tt.storage))
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
