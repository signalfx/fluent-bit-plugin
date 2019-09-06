package util

import (
	"C"
	"fmt"
	"github.com/fluent/fluent-bit-go/output"
	"github.com/signalfx/golib/datapoint"
	"time"
)

func ValueOrDefault(value string, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

func MapValueOrDefault(record map[interface{}]interface{}, key string, defaultValue string) string {
	if value, exist := record[key]; exist {
		return fmt.Sprintf("%s", value)
	}
	return defaultValue
}

func FLBTimestampAsTime(ts interface{}) time.Time {
	switch val := ts.(type) {
	case output.FLBTime:
		return val.Time
	case uint64:
		return time.Unix(int64(val), 0)
	default:
		return time.Now()
	}
}

func MetricTypeAsString(mt datapoint.MetricType) string {
	switch mt {
	case datapoint.Count:
		return "counter"
	case datapoint.Counter:
		return "cumulative counter"
	case datapoint.Gauge:
		return "gauge"
	default:
		return fmt.Sprintf("MetricType(%d)", int(mt))
	}
}
