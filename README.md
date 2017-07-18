# DataDog

[![Build Status](https://travis-ci.org/syntaqx/go-metrics-datadog.svg?branch=master)](https://travis-ci.org/syntaqx/go-metrics-datadog)
[![GoDoc](https://godoc.org/github.com/syntaqx/go-metrics-datadog?status.svg)](https://godoc.org/github.com/syntaqx/go-metrics-datadog)
![license](https://img.shields.io/github/license/syntaqx/go-metrics-datadog.svg)

[go-metrics]: https://github.com/rcrowley/go-metrics
[datadog-go]: https://github.com/DataDog/datadog-go
[license]:    ./LICENSE

This package provides a reporter for the [go-metrics][] library that will post
the metrics to [datadog-go][].

## Installation

```sh
go get -u github.com/syntaqx/go-metrics-datadog
go get -u github.com/DataDog/datadog-go/...
```

## Usage

Simply check out the [example](./example/main.go) file for implementation
details.

## License

Distributed under the MIT license. See [LICENSE][] file for details.
