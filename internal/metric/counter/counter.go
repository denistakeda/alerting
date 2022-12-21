package counter

import (
	"errors"
	"strconv"
)

type counter struct {
	name  string
	value int64
}

func New(name string, value int64) *counter {
	return &counter{
		name:  name,
		value: value,
	}
}

func FromStr(name, value string) (*counter, error) {
	val, err := strToValue(value)
	if err != nil {
		return nil, err
	}
	return &counter{name, val}, nil
}

func (c *counter) Type() string {
	return "counter"
}

func (c *counter) Name() string {
	return c.name
}

func (c *counter) Value() any {
	return c.value
}

func (c *counter) StrValue() string {
	return strconv.FormatInt(c.value, 10)
}

func (c *counter) UpdateValue(val any) error {
	switch v := val.(type) {
	case string:
		vInt64, err := strToValue(v)
		if err != nil {
			return err
		}
		c.value += vInt64
	case int64:
		c.value += v
	default:
		return errors.New("unknown type")
	}
	return nil
}

func strToValue(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
