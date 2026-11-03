# DataDog

[![Test](https://github.com/syntaqx/go-metrics-datadog/actions/workflows/test.yml/badge.svg)](https://github.com/syntaqx/go-metrics-datadog/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/go-metrics-datadog.svg)](https://pkg.go.dev/github.com/syntaqx/go-metrics-datadog)

> [!WARNING]
> **This project is deprecated and archived.** It is a reporter for
> [`rcrowley/go-metrics`](https://github.com/rcrowley/go-metrics), which is itself
> archived and no longer maintained. It will continue to work as-is, but will not
> receive updates. For new code, consider:
>
> - [OpenTelemetry for Go](https://opentelemetry.io/docs/languages/go/): the Datadog Agent accepts OTLP natively.
> - [`hashicorp/go-metrics`](https://github.com/hashicorp/go-metrics): a maintained metrics registry with a built-in Datadog sink.
> - [`DataDog/datadog-go/v5`](https://github.com/DataDog/datadog-go): the official DogStatsD client, used directly.

`go-datadog-metrics` defines a DataDog reporter for [`rcrowley/go-metrics`](https://github.com/rcrowley/go-metrics).

## Getting started

```sh
go get github.com/syntaqx/go-metrics-datadog
```

## Usage

Check out the [example implementation](./example) for the base implementation.

## License

This project is open source available under the [MIT license](./LICENSE)
