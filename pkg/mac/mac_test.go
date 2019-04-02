package mac

import (
	"reflect"
	"testing"
)

func TestTransformPrometheusline(t *testing.T) {
	tests := []struct {
		name     string
		ipmac    map[string]string
		line     string
		expected string
	}{
		{
			"Missing mac",
			map[string]string{},
			`fping_response_duration_seconds_bucket{ip="192.168.1.153",le="0.4096"} 7652`,
			`fping_response_duration_seconds_bucket{mac="UNKNOWN",ip="192.168.1.153",le="0.4096"} 7652`,
		},
		{
			"Simple mac exists",
			map[string]string{"192.168.1.153": "00:11:22:33:44:55:66"},
			`fping_response_duration_seconds_bucket{ip="192.168.1.153",le="0.4096"} 7652`,
			`fping_response_duration_seconds_bucket{mac="00:11:22:33:44:55:66",ip="192.168.1.153",le="0.4096"} 7652`,
		},
		{
			"Unrelated line",
			map[string]string{},
			`# TYPE go_memstats_next_gc_bytes gauge`,
			`# TYPE go_memstats_next_gc_bytes gauge`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transformer := &macTransformer{ipmac: test.ipmac}
			parsed := transformer.Transform(test.line)
			if !reflect.DeepEqual(test.expected, parsed) {
				t.Errorf("\n%+v\n!=\n%+v", test.expected, parsed)
			}
		})
	}
}
