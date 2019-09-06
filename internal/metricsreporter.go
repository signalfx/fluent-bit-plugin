package internal

import (
	"C"
	"context"
	"fmt"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/signalfx-fluent-bit-plugin/internal/config"
	"github.com/signalfx/signalfx-fluent-bit-plugin/util"
	"sync"
	"time"
)

type MetricsReporter struct {
	logger    util.Logger
	config    config.Config
	scheduler Scheduler
	buffer    []*datapoint.Datapoint
	mutex     sync.Mutex
	context   context.Context
	cancelFn  context.CancelFunc
}

var oneAsIntValue = datapoint.NewIntValue(1)

func NewMetricsReporter(logger util.Logger, config config.Config, scheduler Scheduler) *MetricsReporter {
	buffer := make([]*datapoint.Datapoint, 0, config.BufferSize)
	ctx, cancelFn := context.WithCancel(context.Background())
	reporter := &MetricsReporter{logger, config, scheduler, buffer, sync.Mutex{}, ctx, cancelFn}

	scheduler.AddCallback(reporter)
	go scheduler.Schedule(ctx)

	return reporter
}

func (r *MetricsReporter) AddMetric(ts interface{}, record map[interface{}]interface{}) {
	timestamp := util.FLBTimestampAsTime(ts)
	metric := util.MapValueOrDefault(record, "MetricName", r.config.MetricName)
	metricType := r.config.MetricType
	dimensions := r.getDimensions(record)

	dp := datapoint.New(metric, dimensions, oneAsIntValue, metricType, timestamp)

	if r.config.DebugLogEnabled {
		r.printRecord(timestamp, record)
	}

	r.mutex.Lock()
	r.buffer = append(r.buffer, dp)
	if len(r.buffer) == r.config.BufferSize {
		defer r.scheduler.ReportOnce(r.context) // flush entire buffer
	}
	r.mutex.Unlock()
}

func (r *MetricsReporter) Datapoints() []*datapoint.Datapoint {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	result := make([]*datapoint.Datapoint, len(r.buffer))
	copy(result, r.buffer)
	r.buffer = r.buffer[:0]
	return result
}

func (r *MetricsReporter) Close() {
	r.cancelFn()
}

func (r *MetricsReporter) getDimensions(record map[interface{}]interface{}) map[string]string {
	result := make(map[string]string)
	for _, dimension := range r.config.Dimensions {
		if value, exist := record[dimension]; exist {
			result[dimension] = fmt.Sprintf("%s", value)
		}
	}
	return result
}

func (r *MetricsReporter) printRecord(t time.Time, record map[interface{}]interface{}) {
	r.logger.Debug("Processing record:")
	r.logger.Debugf("\ttime: %s", t)

	for k, v := range record {
		r.logger.Debugf("\t%s: %s", k, v)
	}
}
