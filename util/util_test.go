package util

import (
	"github.com/fluent/fluent-bit-go/output"
	"github.com/signalfx/golib/datapoint"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestUtils(t *testing.T) {

	Convey("ValueOrDefault", t, func() {

		Convey("non empty value shall be returned", func() {
			So(ValueOrDefault("abc", "def"), ShouldEqual, "abc")
		})

		Convey("default shall be returned when value is empty", func() {
			So(ValueOrDefault("", "def"), ShouldEqual, "def")
		})
	})

	Convey("MapValueOrDefault", t, func() {

		mapOf := func(key interface{}, value interface{}) map[interface{}]interface{} { return map[interface{}]interface{}{key: value} }

		Convey("map value shall be returned when present", func() {
			So(MapValueOrDefault(mapOf("a", "b"), "a", "def"), ShouldEqual, "b")
		})

		Convey("map []uint8 value shall be returned as string", func() {
			So(MapValueOrDefault(mapOf("a", []uint8{88, 89, 90}), "a", "def"), ShouldEqual, "XYZ")
		})

		Convey("default value shall be returned when map does not contain requested key", func() {
			So(MapValueOrDefault(mapOf("a", "b"), "x", "def"), ShouldEqual, "def")
		})
	})

	Convey("FLBTimestampAsTime", t, func() {

		t := time.Now()

		Convey("shall convert uint64", func() {
			So(FLBTimestampAsTime(uint64(t.Unix())), ShouldEqual, time.Unix(t.Unix(), 0))
		})

		Convey("shall convert FLBTime struct", func() {
			flbTime := output.FLBTime{Time: t}
			So(FLBTimestampAsTime(flbTime), ShouldEqual, t)
		})

		Convey("unknown type shall yield current time", func() {
			t := time.Now()
			So(FLBTimestampAsTime(struct{ unknown int }{1}), ShouldHappenOnOrAfter, t)
		})
	})

	Convey("MetricTypeAsString", t, func() {

		Convey("shall convert all supported metric types", func() {
			testCases := map[datapoint.MetricType]string{
				datapoint.Gauge:   "gauge",
				datapoint.Count:   "counter",
				datapoint.Counter: "cumulative counter",
			}
			for metricType, expectedText := range testCases {
				So(MetricTypeAsString(metricType), ShouldEqual, expectedText)
			}
		})

		Convey("shall handle unknown type", func() {
			So(MetricTypeAsString(datapoint.MetricType(1234)), ShouldEqual, "MetricType(1234)")

		})
	})

}
