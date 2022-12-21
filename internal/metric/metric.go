package metric

type Metric interface {
	Type() string
	Name() string
	Value() any
	StrValue() string
	UpdateValue(val any) error
}
