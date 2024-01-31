package lib

import (
	"sync"
	"time"
)

// PerformanceMetrics holds the metrics for performance evaluation.
type PerformanceMetrics struct {
    Mu               sync.Mutex
    TotalRequests    int64
    TotalResponses   int64
    TotalLatency     time.Duration
    MaxLatency       time.Duration
    MinLatency       time.Duration
    ResponseCounters map[int]int
}

type RequestError struct {
	Verb string
	URL  string
	Err  error
}
