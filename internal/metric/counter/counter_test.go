package counter

import (
	"reflect"
	"testing"

	"github.com/denistakeda/alerting/internal/metric"
)

func Test_counter_ImplementsMetric(t *testing.T) {
	var _ metric.Metric = (*counter)(nil)
}

func TestFromStr(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *counter
		wantErr bool
	}{
		{
			name: "simple test",
			args: args{
				name:  "test",
				value: "3",
			},
			want:    &counter{"test", 3},
			wantErr: false,
		},
		{
			name: "error test",
			args: args{
				name:  "test",
				value: "test",
			},
			want:    nil,
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

func Test_counter_UpdateValue(t *testing.T) {
	type fields struct {
		name  string
		value int64
	}
	tests := []struct {
		name        string
		fields      fields
		newVal      any
		wantCounter *counter
		wantErr     bool
	}{
		{
			name:        "simple test",
			fields:      fields{"test", 3},
			newVal:      int64(5),
			wantCounter: &counter{"test", 8},
			wantErr:     false,
		},
		{
			name:        "update with string",
			fields:      fields{"test", 3},
			newVal:      "5",
			wantCounter: &counter{"test", 8},
			wantErr:     false,
		},
		{
			name:        "update with incorrect string",
			fields:      fields{"test", 3},
			newVal:      "invalid",
			wantCounter: &counter{"test", 3},
			wantErr:     true,
		},
		{
			name:        "update with unknown type",
			fields:      fields{"test", 3},
			newVal:      5.18,
			wantCounter: &counter{"test", 3},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{tt.fields.name, tt.fields.value}
			err := c.UpdateValue(tt.newVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("counter.UpdateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(c, tt.wantCounter) {
				t.Errorf("UpdateValue() = %v, want %v", c, tt.wantCounter)
			}
		})
	}
}
