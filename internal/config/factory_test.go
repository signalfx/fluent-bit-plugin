package config

import (
	"github.com/signalfx/golib/datapoint"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestFactory(t *testing.T) {

	Convey("Factory", t, func() {

		rawConfig := map[string]string{
			"Id":            "abc",
			"IngestURL":     "http://a.b.c",
			"Token":         "def",
			"MetricName":    "com.example.abc",
			"MetricType":    "gauge",
			"Dimensions":    "foo, bar, baz",
			"BufferSize":    "987",
			"ReportingRate": "10s",
			"LogLevel":      "info",
		}

		valueGetter := func(key string) string {
			return rawConfig[key]
		}

		Convey("shall return non nil instance", func() {
			So(NewFactory(func(string) string { return "" }), ShouldNotBeNil)
		})

		Convey("shall parse valid config", func() {
			config := NewFactory(valueGetter).GetConfig()
			So(config.Id, ShouldEqual, "abc")
			So(config.IngestURL, ShouldEqual, "http://a.b.c")
			So(config.Token, ShouldEqual, "def")
			So(config.MetricName, ShouldEqual, "com.example.abc")
			So(config.MetricType, ShouldEqual, datapoint.Gauge)
			So(config.Dimensions, ShouldResemble, []string{"foo", "bar", "baz"})
			So(config.BufferSize, ShouldEqual, 987)
			So(config.ReportingRate, ShouldEqual, time.Second*10)
			So(config.LogLevel, ShouldEqual, zapcore.InfoLevel)
		})

		Convey("shall generate Id", func() {
			delete(rawConfig, "Id")
			factory := NewFactory(valueGetter)
			So(factory.GetConfig().Id, ShouldEqual, "SignalFx.2")
			So(factory.GetConfig().Id, ShouldEqual, "SignalFx.3")
		})

		Convey("shall use default IngestURL", func() {
			delete(rawConfig, "IngestURL")
			config := NewFactory(valueGetter).GetConfig()
			So(config.IngestURL, ShouldEqual, "https://ingest.signalfx.com")
		})

		Convey("shall handle 'counter' metric type", func() {
			rawConfig["MetricType"] = "counter"
			factory := NewFactory(valueGetter)
			So(factory.GetConfig().MetricType, ShouldEqual, datapoint.Count)
		})

		Convey("shall handle 'cumulative counter' metric type", func() {
			rawConfig["MetricType"] = "cumulative counter"
			factory := NewFactory(valueGetter)
			So(factory.GetConfig().MetricType, ShouldEqual, datapoint.Counter)
		})

		Convey("shall parse dimensions", func() {
			rawConfig["Dimensions"] = ",, ,host,cluster , realm, env "
			config := NewFactory(valueGetter).GetConfig()
			So(config.Dimensions, ShouldResemble, []string{"host", "cluster", "realm", "env"})
		})

		Convey("shall handle missing dimensions", func() {
			delete(rawConfig, "Dimensions")
			config := NewFactory(valueGetter).GetConfig()
			So(config.Dimensions, ShouldResemble, []string{})
		})

		Convey("shall panic when metric type is invalid", func() {
			rawConfig["MetricType"] = "blah"
			So(func() { NewFactory(valueGetter).GetConfig() }, ShouldPanicWith, "Invalid value for \"MetricType\": \"blah\". Supported values: \"gauge\", \"counter\", \"cumulative counter\"")
		})

		Convey("shall panic when buffer size is not a number", func() {
			rawConfig["BufferSize"] = "abc"
			So(func() { NewFactory(valueGetter).GetConfig() }, ShouldPanicWith, "Invalid value for \"BufferSize\": cannot parse \"abc\" as int")
		})

		Convey("shall panic when buffer size is out of range", func() {
			rawConfig["BufferSize"] = "1"
			So(func() { NewFactory(valueGetter).GetConfig() }, ShouldPanicWith, "Invalid value for \"BufferSize\": value has to be greater than or equal to 100")
		})

		Convey("shall panic when reporting rate is not valid", func() {
			rawConfig["ReportingRate"] = "qwe"
			So(func() { NewFactory(valueGetter).GetConfig() }, ShouldPanicWith, "Invalid value for \"ReportingRate\": cannot parse \"qwe\" as duration")
		})

		Convey("shall panic when reporting rate is out of range", func() {
			rawConfig["ReportingRate"] = "100ms"
			So(func() { NewFactory(valueGetter).GetConfig() }, ShouldPanicWith, "Invalid value for \"ReportingRate\": value has to be greater than or equal to 1s")
		})

		Convey("shall panic when log level is not valid", func() {
			rawConfig["LogLevel"] = "asd"
			So(func() { NewFactory(valueGetter).GetConfig() }, ShouldPanicWith, "Invalid value for \"LogLevel\": cannot parse \"asd\" as log level")
		})
	})
}
