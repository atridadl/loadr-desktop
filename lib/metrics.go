package lib

import (
	"math"
	"net/http"
	"time"
)

// Initialize the metrics with default values.
var metrics = PerformanceMetrics{
	TotalRequests:    0,
	TotalResponses:   0,
	TotalLatency:     0,
	MaxLatency:       0,
	MinLatency:       time.Duration(math.MaxInt64),
	ResponseCounters: make(map[int]int),
}

// updateMetrics updates the performance metrics.
func UpdateMetrics(duration time.Duration, resp *http.Response, second int) {
	metrics.Mu.Lock()
	defer metrics.Mu.Unlock()

	convertedDuration := duration / time.Millisecond

	metrics.TotalRequests++
	metrics.TotalLatency += convertedDuration
	if convertedDuration > metrics.MaxLatency {
		metrics.MaxLatency = convertedDuration
	}
	if convertedDuration < metrics.MinLatency {
		metrics.MinLatency = convertedDuration
	}
	if resp.StatusCode == http.StatusOK {
		metrics.TotalResponses++
		metrics.ResponseCounters[second]++
	}
}

// calculateAndPrintMetrics calculates and prints the performance metrics.
func GetMetrics() *PerformanceMetrics {
	return &metrics
}
