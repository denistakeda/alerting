package gauge

import (
	"errors"
	"strconv"
)

type gauge struct {
	name  string
	value float64
}

func New(name string, value float64) *gauge {
	return &gauge{
		name:  name,
		value: value,
	}
}

func FromStr(name, value string) (*gauge, error) {
	val, err := strToValue(value)
	if err != nil {
		return nil, err
	}
	return &gauge{name, val}, nil
}

func (g *gauge) Type() string {
	return "gauge"
}

func (g *gauge) Name() string {
	return g.name
}

func (g *gauge) Value() any {
	return g.value
}

func (g *gauge) StrValue() string {
	return strconv.FormatFloat(g.value, 'f', 3, 64)
}

func (g *gauge) UpdateValue(val any) error {
	switch v := val.(type) {
	case string:
		vFloat64, err := strToValue(v)
		if err != nil {
			return err
		}
		g.value = vFloat64
	case float64:
		g.value = v
	default:
		return errors.New("unknown type")
	}
	return nil
}

func strToValue(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}
