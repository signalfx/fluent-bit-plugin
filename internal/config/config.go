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

var gitVersion string // initialized by -ldflags build option

func (c Config) String() string {
	builder := strings.Builder{}
	addLine := func(format string, a...interface{}) { builder.WriteString(fmt.Sprintf(format + "\n", a...)) }

	addLine("SignalFx output plugin (%s)", gitVersion)
	addLine("Id            = %s", c.Id)
	addLine("Ingest URL    = %s", c.IngestURL)
	addLine("Token         = %s", obfuscatedToken(c.Token))
	addLine("Metric Name   = %s", util.ValueOrDefault(c.MetricName, "<no default value>"))
	addLine("Metric Type   = %s", c.MetricType)
	addLine("Dimensions    = %s", strings.Join(c.Dimensions, ", "))
	addLine("BufferSize    = %d", c.BufferSize)
	addLine("ReportingRate = %s", c.ReportingRate)
	addLine("LogLevel      = %s", c.LogLevel)

	return builder.String()
}

func obfuscatedToken(token string) string {
	if len(token) < minTokenLength {
		return "<invalid token>"
	}
	return fmt.Sprintf("%s...%s", token[0:2], token[len(token)-2:])
}
