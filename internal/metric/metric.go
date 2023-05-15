package metric

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/denistakeda/alerting/proto"
	"github.com/pkg/errors"
)

// Type represents a type of the metric.
type Type string

const (
	// Gauge metric type.
	Gauge Type = "gauge"
	// Counter metric type.
	Counter Type = "counter"
)

// Metric is a record to store a metric of one of two types: gauge or counter.
type Metric struct {
	ID    string   `json:"id" db:"id"`
	MType Type     `json:"type" db:"mtype"`
	Value *float64 `json:"value,omitempty" db:"value"`
	Delta *int64   `json:"delta,omitempty" db:"delta"`
	Hash  string   `json:"hash,omitempty" db:"-"`
}

// NewGauge instantiates a new metric of type Gauge.
func NewGauge(name string, value float64) *Metric {
	return &Metric{
		MType: Gauge,
		ID:    name,
		Value: &value,
	}
}

// NewCounter instantiates a new metric of type Counter.
func NewCounter(name string, delta int64) *Metric {
	return &Metric{
		MType: Counter,
		ID:    name,
		Delta: &delta,
	}
}

func FromProto(p *proto.Metric) *Metric {
	mtype := Gauge
	if p.Mtype == proto.Metric_COUNTER {
		mtype = Counter
	}
	return &Metric{
		ID:    p.Id,
		MType: mtype,
		Value: &p.Value,
		Delta: &p.Delta,
		Hash:  p.Hash,
	}
}

func (m *Metric) ToProto() *proto.Metric {
	res := &proto.Metric{
		Id:    m.ID,
		Mtype: proto.Metric_UNSPECIFIED,
		Hash:  m.Hash,
	}

	switch m.MType {
	case Gauge:
		res.Mtype = proto.Metric_GAUGE
	case Counter:
		res.Mtype = proto.Metric_COUNTER
	}

	if m.Value != nil {
		res.Value = *m.Value
	}
	if m.Delta != nil {
		res.Delta = *m.Delta
	}
	return res
}

// Type returns type of the metric.
func (m *Metric) Type() Type {
	return m.MType
}

// Name returns a name of a metric.
func (m *Metric) Name() string {
	return m.ID
}

// StrValue returns the string representation of metric value.
func (m *Metric) StrValue() string {
	switch m.MType {
	case Gauge:
		return strconv.FormatFloat(*m.Value, 'f', 3, 64)
	case Counter:
		return strconv.FormatInt(*m.Delta, 10)
	default:
		return ""
	}
}

// StrType returns the string representation of metric type.
func (m *Metric) StrType() string {
	switch m.MType {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	default:
		return ""
	}
}

// Validate validates the metric.
func (m *Metric) Validate() error {
	switch m.MType {
	case Gauge:
		if m.Value == nil {
			return fmt.Errorf("metric should have a 'value' field for type 'gauge'")
		}
	case Counter:
		if m.Delta == nil {
			return fmt.Errorf("metric should have a 'delta' field for type 'counter'")
		}
	default:
		return fmt.Errorf("unknown metric type: '%s'", m.MType)
	}

	return nil
}

// String representation of a metric.
func (m *Metric) String() string {
	res, err := json.Marshal(m)
	if err != nil {
		return "error" // should never happen
	}
	return string(res)
}

// FillHash fills hash of a metric.
func (m *Metric) FillHash(hashKey string) {
	if hashKey == "" {
		return
	}

	switch m.MType {
	case Gauge:
		m.Hash = getGaugeHash(m.ID, *m.Value, hashKey)
	case Counter:
		m.Hash = getCounterHash(m.ID, *m.Delta, hashKey)
	}
}

// VerifyHash verifies the hash of the metric.
func (m *Metric) VerifyHash(hashKey string) error {
	if hashKey == "" {
		return nil
	}

	var isValid bool

	switch m.MType {
	case Gauge:
		isValid = m.Hash == getGaugeHash(m.ID, *m.Value, hashKey)
	case Counter:
		isValid = m.Hash == getCounterHash(m.ID, *m.Delta, hashKey)
	}

	if !isValid {
		return errors.New("invalid hash")
	}

	return nil
}

// Update updates the metric with new data.
func Update(old *Metric, new *Metric) *Metric {
	if old == nil {
		return new
	}
	if new == nil {
		return old
	}
	if old.MType != new.MType {
		return old
	}
	switch old.MType {
	case Gauge:
		return new
	case Counter:
		return NewCounter(old.ID, *old.Delta+*new.Delta)
	default:
		// Should never happen
		return old
	}
}

// TypeFromString converts a string into a metric type.
func TypeFromString(str string) (Type, error) {
	switch str {
	case "gauge":
		return Gauge, nil
	case "counter":
		return Counter, nil
	default:
		return "", fmt.Errorf("unknown metric type: '%s'", str)
	}
}

func getGaugeHash(name string, value float64, hashKey string) string {
	return hash(fmt.Sprintf("%s:gauge:%f", name, value), hashKey)
}

func getCounterHash(name string, delta int64, hashKey string) string {
	return hash(fmt.Sprintf("%s:counter:%d", name, delta), hashKey)
}

func hash(src string, key string) string {
	if key == "" {
		return ""
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))

	return hex.EncodeToString(h.Sum(nil))
}
