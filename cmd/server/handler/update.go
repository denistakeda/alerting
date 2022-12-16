package handler

import (
	"net/http"
	"strings"

	"github.com/denistakeda/alerting/internal/metric"
	s "github.com/denistakeda/alerting/internal/storage"
)

func UpdateHandler(storage s.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		segments := strings.Split(r.URL.Path, "/")
		// The slice is expected to have 5 segments:
		// [ "", "update", "metricType", "metricName", "metricValue"]
		if len(segments) != 5 || segments[1] != "update" {
			http.Error(w, "method not found", http.StatusNotFound)
			return
		}
		metricType, err := metric.ParseType(segments[2])
		if err != nil {
			http.Error(w, "method not allowed", http.StatusNotImplemented)
			return
		}

		metricName := segments[3]
		metricValue := segments[4]

		if err = s.Store(storage, metricType, metricName, metricValue); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
