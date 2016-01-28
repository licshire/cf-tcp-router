package metrics_reporter

import "github.com/cloudfoundry/dropsonde/metrics"

type PrefixedValue string

func (name PrefixedValue) Send(prefix string, value uint64) {
	metrics.SendValue(prefix+"."+string(name), float64(value), "Metric")
}

type PrefixedDurationMs string

func (name PrefixedDurationMs) Send(prefix string, duration uint64) {
	metrics.SendValue(prefix+"."+string(name), float64(duration), "ms")
}

type PrefixedValueByPort string

func (name PrefixedValueByPort) Send(prefix string, port string, value uint64) {
	metrics.SendValue(prefix+"."+port+"."+string(name), float64(value), "Metric")
}

type PrefixedDurationMsByPort string

func (name PrefixedDurationMsByPort) Send(prefix string, port string, duration uint64) {
	metrics.SendValue(prefix+"."+port+"."+string(name), float64(duration), "ms")
}
