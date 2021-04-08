# Deprecation Notice

:warning: **Be advised this project is deprecated.** :warning:

Use the official [Fluent Bit
forwarder](https://docs.fluentbit.io/manual/installation/getting-started-with-fluent-bit) with the [Splunk
output](https://docs.fluentbit.io/manual/pipeline/outputs/splunk).

# SignalFx Output Plugin for Fluent Bit

The SignalFx output plugin for [Fluent Bit](https://docs.fluentbit.io) sends log-based metrics to [SignalFx](https://www.SignalFx.com).

This enables you to filter your logs for specific phrases like "error", "exception", etc. and have the plugin report metrics whenever any of those phrases is present in a log stream.

## Docker Images

There are two Docker images available with the SignalFx output plugin for Fluent Bit:
   * [quay.io/signalfx/fluent-bit](https://quay.io/repository/signalfx/fluent-bit?tab=tags) ([Dockerfile](build/package/Dockerfile)) -- this image is based on the official [Fluent Bit image](https://hub.docker.com/r/fluent/fluent-bit/tags). This image contains Fluent Bit binaries and the SignalFx output plugin.

  * [quay.io/signalfx/fluent-bit-aws](https://quay.io/repository/signalfx/fluent-bit-aws?tab=tags) ([Dockerfile](build/package/Dockerfile.aws)) -- this image is based on the the official [Amazon Fluent Bit image](https://hub.docker.com/r/amazon/aws-for-fluent-bit). This image contains Fluent Bit binaries, additional plugins for AWS Firehose and AWS CloudWatch provided by Amazon, and the SignalFx output plugin.

## Configuration Parameters

| Key | Description | Default | Example |
| :--- | :--- | :--- | :--- |
| IngestURL | Specifies the data ingest address of SignalFx. This address is on your SignalFx profile page. | https://ingest.signalfx.com | https://ingest.eu0.signalfx.com |
| Token | Specifies the SignalFx [access token](https://docs.signalfx.com/en/latest/admin-guide/tokens.html#working-with-access-tokens). This access token is on your profile page in SignalFx. | | abcdefgh12345678 |
| MetricName | Specifies the [metric name](https://docs.signalfx.com/en/latest/reference/glossary/glossary.html#term-metric). You can override metric name for each Fluent Bit record using `modify` filter. Refer to `com.example.app.error` and `com.example.app.exception` in the [sample config](example/fluent-bit.conf) to see how metric names are defined. |  | com.example.app.requests |
| MetricType | Specifies [metric type](https://docs.signalfx.com/en/latest/metrics-metadata/metric-types.html#metric-types). | gauge | "gauge", "counter" or "cumulative counter" (without quotes) |
| Dimensions | Specifies a list of dimensions attached to reported metric. For instance if your Fluent Bit record contains "ecs_cluster" and "container_name" fields you can use them as dimensions. If you want to add an additional dimension that is not available in the Fluent Bit record you may use [Fluent Bit filters](https://docs.fluentbit.io/manual/filter) to add extra fields to a record. Refer to the [sample config](example/fluent-bit.conf) to see how the additional `env` dimension is configured. | empty list | ecs_cluster, container_name, realm |
| BufferSize | Specifies maximum number of metrics to buffer before they are sent to SignalFx. Minimum value is 100. | 10000 | any value >= 100 |
| ReportingRate | Specifies how often buffered metrics are sent to SignalFx. Minimum value is 1s. | 5s | 1s, 5s, 3m, etc. |
| LogLevel | Specifies log level for plugin diagnostic messages. | info | debug, info, warning, error |

## Getting Started

You can run the plugin from the command line or through the configuration file (recommended).

### Command Line

The SignalFx plugin can read the parameters from the command line through the `-p` argument \(property\) as shown below. We recommend using the configuration file instead. See **Configuration File** below.

`$ docker run -it --rm fluent-bit-signalfx /fluent-bit/bin/fluent-bit -i cpu -t cpu -e /fluent-bit/signalfx.so -o SignalFx -p IngestURL=https://ingest.corp.signalfx.com -p Token=<ACCESS TOKEN> -p MetricName=com.example.app.requests -p LogLevel=debug -m '*'`

### Configuration File

Refer to the [fluent-bit.conf](example/fluent-bit.conf) file for an example of how to configure the SignalFx plugin.

### Reporting Rate vs. Buffer Size 

Metrics are reported at a rate specified by the `ReportingRate` configuration parameter. They are also reported whenever buffer gets full (see `BufferSize` above).

### Test Data and Metric Timestamps

Please note the SignalFx Ingest service silently rejects datapoints it receives if they have a timestamp that’s earlier than previous datapoints for the same MTS. See [Sending Metrics and Events](https://developers.signalfx.com/metrics/data_ingest_overview.html) documentation for details.

This silent rejection becomes important when you test the SignalFx plugin using captured log data. In such a case you need to update timestamps in the captured log file for each test run. Another option is to use a different metric name for each test run.

For test purposes you may use the [gen-log.js](example/gen-log.js) node.js script to generate sample log data. See **Engineering Notes** below.

### Engineering Notes

You may find the following snippets useful when working with the SignalFx output plugin for Fluent Bit.

`make images TAG=x.y.z` - builds Docker images with the SignalFx plugin and tags them (refer to **Docker Images** for details).

`make demo TAG=x.y.z INGEST_URL=https://ingest.signalfx.com TOKEN=<ACCESS TOKEN>` -- reports sample metric `com.example.app.requests` using default `ReportingRate`.

`node example/gen-log.js 10 > example/fluent-bit-sample.log` -- generates sample log file with 10 log records using timestamps close to current time.

`docker run -it --rm -v $(PWD)/example/fluent-bit-sample.log:/fluent-bit-sample.log -v $(PWD)/example/fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf quay.io/signalfx/fluent-bit:x.y.z` - runs the Docker container with [fluent-bit-sample.log](example/fluent-bit-sample.log) as an input and [fluent-bit.conf](example/fluent-bit.conf) as a config file.
