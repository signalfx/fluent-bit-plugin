package config

import (
	"C"
	"fmt"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/signalfx-fluent-bit-plugin/util"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type Config struct {
	Id              string
	MetricType      datapoint.MetricType
	MetricName      string
	Token           string
	IngestURL       string
	Dimensions      []string
	BufferSize      int
	ReportingRate   time.Duration
	LogLevel        zapcore.Level
	DebugLogEnabled bool
}

func (c Config) String() string {
	builder := strings.Builder{}
	addLine := func(format string, a...interface{}) { builder.WriteString(fmt.Sprintf(format + "\n", a...)) }

	addLine("SignalFx output plugin configuration:")
	addLine("Id            = %s", c.Id)
	addLine("Ingest URL    = %s", c.IngestURL)
	addLine("Token         = %s...%s", c.Token[0:2], c.Token[len(c.Token)-2:])
	addLine("Metric Name   = %s", util.ValueOrDefault(c.MetricName, "<no default value>"))
	addLine("Metric Type   = %s", util.MetricTypeAsString(c.MetricType))
	addLine("Dimensions    = %s", strings.Join(c.Dimensions, ", "))
	addLine("BufferSize    = %d", c.BufferSize)
	addLine("ReportingRate = %s", c.ReportingRate)
	addLine("LogLevel      = %s", c.LogLevel)

	return builder.String()
}