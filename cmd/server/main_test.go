package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	"github.com/denistakeda/alerting/mocks"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_updateMetric(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		met      *metric.Metric
		wantCode int
	}{
		{
			name:     "gauge success case",
			request:  "/update/gauge/metric_name/100",
			met:      metric.NewGauge("metric_name", 100),
			wantCode: http.StatusOK,
		},
		{
			name:     "counter success case",
			request:  "/update/counter/metric_name/100",
			met:      metric.NewCounter("metric_name", 100),
			wantCode: http.StatusOK,
		},
		{
			name:     "gauge without name and type",
			request:  "/update/gauge/",
			met:      nil,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "counter without name and type",
			request:  "/update/counter/",
			met:      nil,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "gauge invalid value",
			request:  "/update/gauge/test_counter/none",
			met:      nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "counter invalid value",
			request:  "/update/counter/test_counter/none",
			met:      nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "unknown type",
			request:  "/update/unknown/testCounter/100",
			met:      nil,
			wantCode: http.StatusNotImplemented,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStorage(ctrl)
			s.EXPECT().Update(gomock.Any(), tt.met).Return(tt.met, nil).AnyTimes()

			router := setupRouter(s, "", loggerservice.New())

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
	type storageMock struct {
		reqType   metric.Type
		reqName   string
		retMetric *metric.Metric
		retOk     bool
	}
	tests := []struct {
		name        string
		request     string
		storageMock storageMock
		want        want
	}{
		{
			name:    "gauge success case",
			request: fmt.Sprintf("/value/%s/%s", m1.StrType(), m1.Name()),
			storageMock: storageMock{
				reqType:   m1.Type(),
				reqName:   m1.Name(),
				retMetric: m1,
				retOk:     true,
			},
			want: want{
				code: http.StatusOK,
				body: m1.StrValue(),
			},
		},
		{
			name:    "counter success case",
			request: fmt.Sprintf("/value/%s/%s", m3.StrType(), m3.Name()),
			storageMock: storageMock{
				reqType:   m3.Type(),
				reqName:   m3.Name(),
				retMetric: m3,
				retOk:     true,
			},
			want: want{
				code: http.StatusOK,
				body: m3.StrValue(),
			},
		},
		{
			name:    "request not existing metric",
			request: fmt.Sprintf("/value/%s/%s", m2.StrType(), m2.Name()),
			storageMock: storageMock{
				reqType:   m2.Type(),
				reqName:   m2.Name(),
				retMetric: nil,
				retOk:     false,
			},
			want: want{
				code: http.StatusNotFound,
				body: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStorage(ctrl)
			s.EXPECT().
				Get(gomock.Any(), tt.storageMock.reqType, tt.storageMock.reqName).
				Return(tt.storageMock.retMetric, tt.storageMock.retOk).
				AnyTimes()

			router := setupRouter(s, "", loggerservice.New())

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
	g2 := metric.NewGauge("gauge2", 5.18)
	c3 := metric.NewCounter("counter1", 7)
	c4 := metric.NewCounter(c3.Name(), 14)

	type want struct {
		code int
		body []byte
	}
	type storageMock struct {
		reqMetric *metric.Metric
		resMetric *metric.Metric
		resError  error
	}
	tests := []struct {
		name        string
		requestBody []byte
		storageMock storageMock
		want        want
	}{
		{
			name:        "not existing gauge",
			requestBody: marshal(t, g1),
			storageMock: storageMock{
				reqMetric: g1,
				resMetric: g1,
				resError:  nil,
			},
			want: want{
				code: http.StatusOK,
				body: marshal(t, g1),
			},
		},
		{
			name:        "existing gauge",
			requestBody: marshal(t, metric.NewGauge(g2.Name(), 5.18)),
			storageMock: storageMock{
				reqMetric: g2,
				resMetric: g2,
				resError:  nil,
			},
			want: want{
				code: http.StatusOK,
				body: marshal(t, g2),
			},
		},
		{
			name:        "not existing counter",
			requestBody: marshal(t, c3),
			storageMock: storageMock{
				reqMetric: c3,
				resMetric: c3,
				resError:  nil,
			},
			want: want{
				code: http.StatusOK,
				body: marshal(t, c3),
			},
		},
		{
			name:        "existing counter",
			requestBody: marshal(t, c3),
			storageMock: storageMock{
				reqMetric: c3,
				resMetric: c4,
				resError:  nil,
			},
			want: want{
				code: http.StatusOK,
				body: marshal(t, metric.NewCounter(c3.Name(), 14)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStorage(ctrl)
			s.EXPECT().
				Update(gomock.Any(), tt.storageMock.reqMetric).
				Return(tt.storageMock.resMetric, tt.storageMock.resError).
				AnyTimes()

			router := setupRouter(s, "", loggerservice.New())

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
