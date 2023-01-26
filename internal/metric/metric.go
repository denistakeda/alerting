package metric

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
)

type Type string

const (
	Gauge   Type = "gauge"
	Counter Type = "counter"
)

type Metric struct {
	ID    string   `json:"id"`
	MType Type     `json:"type"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func NewGauge(name string, value float64, hashKey string) *Metric {
	return &Metric{
		MType: Gauge,
		ID:    name,
		Value: &value,
		Hash:  hash(fmt.Sprintf("%s:gauge:%f", name, value), hashKey),
	}
}

func NewCounter(name string, delta int64, hashKey string) *Metric {
	return &Metric{
		MType: Counter,
		ID:    name,
		Delta: &delta,
		Hash:  hash(fmt.Sprintf("%s:counter:%d", name, delta), hashKey),
	}
}

func (m *Metric) Type() Type {
	return m.MType
}

func (m *Metric) Name() string {
	return m.ID
}

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

func (m *Metric) String() string {
	res, err := json.Marshal(m)
	if err != nil {
		return "error" // should never happen
	}
	return string(res)
}

func Update(old *Metric, new *Metric, hashKey string) *Metric {
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
		return NewCounter(old.ID, *old.Delta+*new.Delta, hashKey)
	default:
		// Should never happen
		return old
	}
}

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

func hash(src string, key string) string {
	if key == "" {
		return ""
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))

	return string(h.Sum(nil))
}
