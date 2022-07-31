/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

// Package metrics provides abstractions for registering which metrics
// to record.
package metrics

import (
	"net/url"
	"sync"
	"time"
)

var registerMetrics sync.Once

// LatencyMetric observes client latency partitioned by verb and url.
type LatencyMetric interface {
	Observe(verb string, u url.URL, latency time.Duration)
}

// ResultMetric counts response codes partitioned by method and host.
type ResultMetric interface {
	Increment(code string, method string, host string)
}

var (
	// RequestLatency is the latency metric that rest clients will update.
	RequestLatency LatencyMetric = noopLatency{}
	// RequestResult is the result metric that rest clients will update.
	RequestResult ResultMetric = noopResult{}
)

// Register registers metrics for the rest client to use. This can
// only be called once.
func Register(lm LatencyMetric, rm ResultMetric) {
	registerMetrics.Do(func() {
		RequestLatency = lm
		RequestResult = rm
	})
}

type noopLatency struct{}

func (noopLatency) Observe(string, url.URL, time.Duration) {}

type noopResult struct{}

func (noopResult) Increment(string, string, string) {}
