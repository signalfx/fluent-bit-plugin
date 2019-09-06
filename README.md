# SignalFx Output Plugin for Fluent Bit

SignalFx output plugin for [Fluent Bit](https://docs.fluentbit.io) sends log based metrics to [SignalFx](https://www.SignalFx.com) service.

This means you can monitor your logs for specific phrases like "error", "exception", etc. and have the plugin report metrics whenever any of those phrases is present in a log stream.

## Configuration Parameters

| Key | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| IngestURL | Specifies data ingest address of the SignalFx service. You may find this address on your SignalFx profile page. | https://ingest.signalfx.com | https://ingest.eu0.signalfx.com |
| Token | Specify the SignalFx [access token](https://docs.signalfx.com/en/latest/admin-guide/tokens.html#working-with-access-tokens). | | abcdefgh12345678 |
| MetricName | Specify [metric name](https://docs.signalfx.com/en/latest/reference/glossary/glossary.html#term-metric). Please note you can override metric name for each Fluent Bit record using `modify` filter. See how `com.example.app.error` and `com.example.app.exception` metrics are defined in the [sample config](examples/fluent-bit.conf). |   | com.example.app.requests |
| MetricType | Specify [metric type](https://docs.signalfx.com/en/latest/metrics-metadata/metric-types.html#metric-types). | gauge | gauge, counter or cumulative counter |
| Dimensions | Specifies a list of dimensions attached to reported metric. For instance if your Fluent Bit record contains "ecs_cluster" and "container_name" fields you can use them as dimensions. If you want to add additional dimension that is not available in the Fluent Bit record you may use [Fluent Bit's filters](https://docs.fluentbit.io/manual/filter) to add extra fields to a record. See [sample config](examples/fluent-bit.conf) to see how additional `env` dimension is configured there. | | ecs_cluster, container_name, realm |
| BufferSize | Specifies maximum number of metrics to buffer before they are sent to SignalFx. Minimum value is 100. | 10000 | any value >= 100 |
| ReportingRate | Specifies how often buffered metrics are sent to SignalFx. Minimum value is 1s. | 5s | 1s, 5s, 3m, etc. |
| LogLevel | Specifies log level for plugin diagnostic messages. | info | debug, info, warning, error |

## Getting Started

You can run the plugin from the command line or through the configuration file (recommended).

### Command Line

The SignalFx plugin can read the parameters from the command line through the `-p` argument \(property\) as shown below. Although doable it is recommended to use configuration file instead.

`$ docker run -it --rm fluent-bit-signalfx /fluent-bit/bin/fluent-bit -i cpu -t cpu -e /fluent-bit/signalfx.so -o SignalFx -p IngestURL=https://ingest.corp.signalfx.com -p Token=<ACCESS TOKEN> -p MetricName=com.example.app.requests -p LogLevel=debug -m '*'`

### Configuration File

See sample config file: [fluent-bit.conf](examples/fluent-bit.conf) 

### Reporting Rate vs. Buffer Size 

Metrics are reported at a rate specified by the `ReportingRate` config param. They are also reported whenever buffer gets full (see `BufferSize` above).

### Test Data and Metric Timestamps

Please note SignalFx Ingest service silently rejects datapoints it receives if they have a timestamp that’s earlier than previous datapoints for the same MTS. See [Sending Metrics and Events](https://developers.signalfx.com/metrics/data_ingest_overview.html) documentation for details.

This is important when you test the SignalFx plugin using captured log data: in such a case you either need to update timestamps in the captured log file or use different metric name whenever you try to ingest metrics to SignalFx to ensure ingested timestamps are newer than those ingested earlier.

For test purposes you may use the `examples/gen-log.js` node.js script to generate sample log data. See below for details.

### Engineering Notes

You may find the following snippets useful when working with SignalFx plugin.

`docker build . -t fluent-bit-signalfx -f ./build/package/Dockerfile` - builds new Docker image with the SignalFx plugin on top of the official [Fluent Bit image](https://hub.docker.com/r/fluent/fluent-bit/tags). If you want to report log based metrics to SignalFx only most likely this is the image you should use.

`docker build . -t fluent-bit-signalfx -f ./build/package/Dockerfile.aws` - builds new Docker image with the SignalFx plugin on top of the official [Amazon Fluent Bit image](https://hub.docker.com/r/fluent/fluent-bit/tags). If you want to report log based metrics to SignalFx **and** you also want to use Amazon services like CloudWatch or Firehose this is the image you should use.

`node examples/gen-log.js 10 > examples/fluent-bit-sample.log` -- generates sample log file with timestamps set to current time.

`docker run -it --rm -v $(PWD)/examples/fluent-bit-sample.log:/fluent-bit-sample.log -v $(PWD)/examples/fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf fluent-bit-signalfx` - runs the Docker container with [fluent-bit-sample.log](examples/fluent-bit-sample.log) as an input and [fluent-bit.conf](examples/fluent-bit.conf) as a config file.
