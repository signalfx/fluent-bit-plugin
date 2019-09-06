package config

import (
	"github.com/signalfx/golib/datapoint"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {

	Convey("Config", t, func() {

		config := Config{
			Id:            "abc",
			IngestURL:     "http://a.b.c",
			Token:         "abc123",
			MetricName:    "com.example.abc",
			MetricType:    datapoint.Gauge,
			Dimensions:    []string{"foo", "bar"},
			BufferSize:    2000,
			ReportingRate: time.Minute,
			LogLevel:      zapcore.WarnLevel,
		}

		Convey("shall use stringer", func() {
			So(config.String(), ShouldEqual, `SignalFx output plugin configuration:
Id            = abc
Ingest URL    = http://a.b.c
Token         = ab...23
Metric Name   = com.example.abc
Metric Type   = gauge
Dimensions    = foo, bar
BufferSize    = 2000
ReportingRate = 1m0s
LogLevel      = warn
`)
		})
	})
}
