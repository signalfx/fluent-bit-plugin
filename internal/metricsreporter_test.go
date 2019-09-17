package internal

import (
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/signalfx-fluent-bit-plugin/internal/config"
	"github.com/signalfx/signalfx-fluent-bit-plugin/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestMetricsReporter(t *testing.T) {

	Convey("MetricsReporter", t, func() {

		Convey("shall setup scheduler", func() {
			pluginConfig := config.Config{
				BufferSize:      2,
				ReportingRate:   time.Hour,
				MetricType:      datapoint.Gauge,
				MetricName:      "config.metric",
				Dimensions:      []string{"cluster", "realm"},
			}
			logger := getMockLogger()
			scheduler := getMockScheduler()
			reporter := NewMetricsReporter(logger, pluginConfig, scheduler)

			scheduler.WaitForScheduleToBeCalled()
			scheduler.AssertExpectations(t)

			ts, record := getFluentBitRecord()

			Convey("shall return empty slice when internal buffer is empty", func() {
				points := reporter.Datapoints()
				So(points, ShouldBeEmpty)
			})

			Convey("shall add metric to the buffer using record's metric name", func() {
				reporter.AddMetric(ts, record)
				datapoints := reporter.Datapoints()

				So(datapoints, ShouldHaveLength, 1)
				So(datapoints[0].Timestamp, ShouldEqual, time.Unix(int64(ts), 0))
				So(datapoints[0].Metric, ShouldEqual, "record.metric")
				So(datapoints[0].MetricType, ShouldEqual, datapoint.Gauge)
				So(datapoints[0].Dimensions, ShouldResemble, map[string]string{
					"cluster": "record.cluster",
					"realm":   "record.realm",
				})

				Convey("shall return empty slice when all datapoints were reported", func() {
					points := reporter.Datapoints()
					So(points, ShouldBeEmpty)
				})
			})

			Convey("shall use config's metric name", func() {
				delete(record, "MetricName")
				reporter.AddMetric(ts, record)
				datapoints := reporter.Datapoints()

				So(datapoints, ShouldHaveLength, 1)
				So(datapoints[0].Metric, ShouldEqual, "config.metric")
			})

			Convey("shall ignore missing dimensions", func() {
				delete(record, "cluster")
				reporter.AddMetric(ts, record)
				datapoints := reporter.Datapoints()

				So(datapoints, ShouldHaveLength, 1)
				So(datapoints[0].Dimensions, ShouldResemble, map[string]string{
					"realm": "record.realm",
				})
			})

			Convey("shall flush metrics when buffer is full", func() {
				scheduler.On("ReportOnce", mock.AnythingOfType("*context.cancelCtx")).Return(nil)
				reporter.AddMetric(ts, record)
				reporter.AddMetric(ts, record)
				scheduler.AssertExpectations(t)
			})

			Convey("shall log record", func() {
				reporter.config.DebugLogEnabled = true
				logger.On("Debug", "Processing record:").Return(nil)
				logger.On("Debugf", "\tmetric: %s", mock.AnythingOfType("string")).Return(nil)
				logger.On("Debugf", "\ttime: %s", mock.AnythingOfType("time.Time")).Return(nil)
				logger.On("Debugf", "\t%s: %s", "MetricName", "record.metric").Return(nil)
				logger.On("Debugf", "\t%s: %s", "cluster", "record.cluster").Return(nil)
				logger.On("Debugf", "\t%s: %s", "realm", "record.realm").Return(nil)
				reporter.AddMetric(ts, record)
				logger.AssertExpectations(t)
			})

			Convey("shall cancel scheduler's context when closed", func() {
				So(reporter.context.Err(), ShouldBeNil)
				reporter.Close()
				So(reporter.context.Err(), ShouldNotBeNil)
			})
		})
	})
}

func getMockLogger() *mocks.Logger {
	return new(mocks.Logger)
}

func getMockScheduler() *mocks.Scheduler {
	scheduler := new(mocks.Scheduler)
	scheduler.On("AddCallback", mock.AnythingOfType("*internal.MetricsReporter")).Return(nil)
	scheduler.On("Schedule", mock.AnythingOfType("*context.cancelCtx")).Return(nil)
	scheduler.ExpectScheduleToBeCalled()
	return scheduler
}

func getFluentBitRecord() (uint64, map[interface{}]interface{}) {
	m := make(map[interface{}]interface{})
	m["MetricName"] = "record.metric"
	m["cluster"] = "record.cluster"
	m["realm"] = "record.realm"
	return uint64(time.Now().Unix()), m
}
