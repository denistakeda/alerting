package gauge

import (
	"reflect"
	"testing"

	"github.com/denistakeda/alerting/internal/metric"
)

func Test_gauge_ImplementsMetric(t *testing.T) {
	var _ metric.Metric = (*gauge)(nil)
}

func TestFromStr(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *gauge
		wantErr bool
	}{
		{
			name: "simple test",
			args: args{
				name:  "test",
				value: "3.14",
			},
			want:    &gauge{"test", 3.14},
			wantErr: false,
		},
		{
			name: "error test",
			args: args{
				name:  "test",
				value: "test",
			},
			want:    &gauge{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromStr(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gauge_UpdateValue(t *testing.T) {
	type fields struct {
		name  string
		value float64
	}
	tests := []struct {
		name      string
		fields    fields
		newVal    any
		wantGauge *gauge
		wantErr   bool
	}{
		{
			name:      "simple test",
			fields:    fields{"test", 3.14},
			newVal:    5.16,
			wantGauge: &gauge{"test", 5.16},
			wantErr:   false,
		},
		{
			name:      "update with string",
			fields:    fields{"test", 3.14},
			newVal:    "5.16",
			wantGauge: &gauge{"test", 5.16},
			wantErr:   false,
		},
		{
			name:      "update with incorrect string",
			fields:    fields{"test", 3.14},
			newVal:    "invalid",
			wantGauge: &gauge{"test", 3.14},
			wantErr:   true,
		},
		{
			name:      "update with unknown type",
			fields:    fields{"test", 3.14},
			newVal:    5,
			wantGauge: &gauge{"test", 3.14},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{tt.fields.name, tt.fields.value}
			err := g.UpdateValue(tt.newVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("gauge.UpdateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(g, tt.wantGauge) {
				t.Errorf("UpdateValue() = %v, want %v", g, tt.wantGauge)
			}
		})
	}
}
