package config

import (
	"C"
	"fmt"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/signalfx-fluent-bit-plugin/util"
	"go.uber.org/zap/zapcore"
	"strconv"
	"strings"
	"time"
)

const defaultIngestURL = "https://ingest.signalfx.com"
const defaultMetricType = "gauge"
const defaultBufferSize = "10000"
const defaultReportingRate = "5s"
const defaultLogLevel = "info"

const minTokenLength = 10 // SFx Access Tokens are 22 chars long in 2019 but accept 10 or more chars just in case
const minBufferSize = 100
const minReportingRate = time.Second

type ValueGetter = func(key string) string

type Factory struct {
	configValue ValueGetter
	logger      util.Logger
}

var signalFxPluginInstanceCounter = 0

func NewFactory(configValue ValueGetter) *Factory {
	return &Factory{configValue, util.GetLogger(zapcore.InfoLevel)}
}

func (f *Factory) GetConfig() Config {
	signalFxPluginInstanceCounter++
	defaultPluginId := fmt.Sprintf("SignalFx.%v", signalFxPluginInstanceCounter)
	logLevel := f.configValueOrDefaultAsLogLevel("LogLevel", defaultLogLevel)

	return Config{
		Id:              util.ValueOrDefault(f.configValue("Id"), defaultPluginId),
		IngestURL:       util.ValueOrDefault(f.configValue("IngestURL"), defaultIngestURL),
		Token:           f.configValueAsToken("Token"),
		MetricName:      f.configValue("MetricName"),
		MetricType:      f.configValueOrDefaultAsMetricType("MetricType", defaultMetricType),
		Dimensions:      f.configValueAsSliceOfStrings("Dimensions", ","),
		BufferSize:      f.configValueOrDefaultAsInt("BufferSize", defaultBufferSize, minBufferSize),
		ReportingRate:   f.configValueOrDefaultAsDuration("ReportingRate", defaultReportingRate, minReportingRate),
		LogLevel:        logLevel,
		DebugLogEnabled: logLevel == zapcore.DebugLevel,
	}
}

func (f *Factory) configValueAsToken(configKey string) string {
	token := f.configValue(configKey)
	if len(token) < minTokenLength {
		f.logger.Panicf("Invalid value for %q: value shall be at least %d characters long", configKey, minTokenLength)
	}
	return token
}

func (f *Factory) configValueOrDefaultAsMetricType(configKey string, defaultValue string) datapoint.MetricType {
	value := util.ValueOrDefault(f.configValue(configKey), defaultValue)
	var mt datapoint.MetricType
	if err := mt.UnmarshalText([]byte(value)); err != nil {
		f.logger.Panicf("Invalid value for %q: %s", configKey, err)
	}
	return mt
}

func (f *Factory) configValueOrDefaultAsDuration(configKey string, defaultValue string, minDuration time.Duration) time.Duration {
	value := util.ValueOrDefault(f.configValue(configKey), defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		f.logger.Panicf("Invalid value for %q: cannot parse %q as duration", configKey, value)
	}
	if duration < minDuration {
		f.logger.Panicf("Invalid value for %q: value has to be greater than or equal to %s", configKey, minDuration)
	}
	return duration
}

func (f *Factory) configValueOrDefaultAsInt(configKey string, defaultValue string, minValue int) int {
	value := util.ValueOrDefault(f.configValue(configKey), defaultValue)
	i, err := strconv.Atoi(value)
	if err != nil {
		f.logger.Panicf("Invalid value for %q: cannot parse %q as int", configKey, value)
	}
	if i < minValue {
		f.logger.Panicf("Invalid value for %q: value has to be greater than or equal to %d", configKey, minValue)
	}
	return i
}

func (f *Factory) configValueOrDefaultAsLogLevel(configKey string, defaultValue string) zapcore.Level {
	value := util.ValueOrDefault(f.configValue(configKey), defaultValue)
	var level zapcore.Level
	err := level.UnmarshalText([]byte(value))
	if err != nil {
		f.logger.Panicf("Invalid value for %q: cannot parse %q as log level", configKey, value)
	}
	return level
}

func (f *Factory) configValueAsSliceOfStrings(configKey string, wordsSeparator string) []string {
	value := f.configValue(configKey)
	words := strings.Split(value, wordsSeparator)
	result := make([]string, 0)
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 0 {
			result = append(result, word)
		}
	}
	return result
}
