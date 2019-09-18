package main

import (
	"C"
	"fmt"
	"github.com/signalfx/golib/sfxclient"
	"github.com/signalfx/signalfx-fluent-bit-plugin/internal/config"
	"github.com/signalfx/signalfx-fluent-bit-plugin/util"
	"log"
	"net/http"
	"strings"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/signalfx/signalfx-fluent-bit-plugin/internal"
)

type PluginInstance struct {
	Logger   util.Logger
	Config   config.Config
	Reporter *internal.MetricsReporter
}

var instances []PluginInstance

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(def, "SignalFx", "SignalFx output plugin.")
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	configFactory := config.NewFactory(func(key string) string { return output.FLBPluginConfigKey(plugin, key) })
	pluginConfig := configFactory.GetConfig()

	logger := util.GetLogger(pluginConfig.LogLevel)
	logger.Info(pluginConfig)

	scheduler := getScheduler(logger, pluginConfig)

	reporter := internal.NewMetricsReporter(logger, pluginConfig, scheduler)

	instance := PluginInstance{logger, pluginConfig, reporter}
	instances = append(instances, instance)

	output.FLBPluginSetContext(plugin, instance)

	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	instance := output.FLBPluginGetContext(ctx).(PluginInstance)
	instance.Logger.Debugf("Flush called for instance %s", instance.Config.Id)

	dec := output.NewDecoder(data, int(length))

	count := 0
	for {
		ret, ts, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}
		instance.Reporter.AddMetric(ts, record)
		count++
	}
	instance.Logger.Infof("Buffered %d metric(s)", count)
	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	log.Print("Flush called for unknown instance")
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	for _, instance := range instances {
		instance.Logger.Infof("Closing instance %q", instance.Config.Id)
		instance.Reporter.Close()
	}
	return output.FLB_OK
}

func getScheduler(logger util.Logger, config config.Config) internal.Scheduler {
	scheduler := sfxclient.NewScheduler()
	scheduler.ReportingDelay(config.ReportingRate)

	sink := scheduler.Sink.(*sfxclient.HTTPSink)
	sink.DatapointEndpoint = fmt.Sprintf("%s/v2/datapoint", strings.TrimRight(config.IngestURL, "/"))
	sink.AuthToken = config.Token

	if config.DebugLogEnabled {
		sink.ResponseCallback = func(resp *http.Response, responseBody []byte) {
			logger.Debugf("Buffered metrics ingest status: %s, body: %s", resp.Status, responseBody)
		}
	}
	return scheduler
}

func main() {
}
