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
	AverageLatency:   0,
	MaxLatency:       0,
	MinLatency:       time.Duration(math.MaxInt64),
	ResponseCounters: make(map[int]int),
}

func UpdateMetrics(duration time.Duration, resp *http.Response, second int) {
	metrics.Mu.Lock()
	defer metrics.Mu.Unlock()

	convertedDuration := duration / time.Millisecond

	metrics.TotalRequests++
	if resp.StatusCode == http.StatusOK {
		metrics.TotalResponses++
		metrics.ResponseCounters[second]++
	}

	// Calculate the average latency
	metrics.AverageLatency = (metrics.AverageLatency*time.Duration(metrics.TotalRequests-1) + convertedDuration) / time.Duration(metrics.TotalRequests)

	if convertedDuration > metrics.MaxLatency {
		metrics.MaxLatency = convertedDuration
	}
	if convertedDuration < metrics.MinLatency {
		metrics.MinLatency = convertedDuration
	}
}

// calculateAndPrintMetrics calculates and prints the performance metrics.
func GetMetrics() *PerformanceMetrics {
	return &metrics
}

// ResetMetrics resets all the performance metrics to their initial values.
func ResetMetrics() {
	metrics := GetMetrics()
	metrics.TotalRequests = 0
	metrics.TotalResponses = 0
	metrics.AverageLatency = 0
	metrics.MaxLatency = 0
	metrics.MinLatency = math.MaxInt64
}
