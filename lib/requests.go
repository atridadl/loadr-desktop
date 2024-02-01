package lib

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Global HTTP client used for making requests.
var client = &http.Client{}

func (e *RequestError) Error() string {
	return fmt.Sprintf("error making %s request to %s: %v", e.Verb, e.URL, e.Err)
}

// makeRequest sends an HTTP request and updates performance metrics.
func makeRequest(verb, url, token string, second int) error {
	startTime := time.Now()

	// Create a new request with the provided verb and URL
	var req *http.Request
	var err error
	req, err = http.NewRequest(verb, url, nil)

	if err != nil {
		return &RequestError{Verb: verb, URL: url, Err: err}
	}

	// Add the bearer token to the request's Authorization header if provided.
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Send the request.
	resp, err := client.Do(req)
	if err != nil {
		return &RequestError{Verb: verb, URL: url, Err: err}
	}
	defer resp.Body.Close()

	// Calculate the duration of the request.
	duration := time.Since(startTime)

	// Update the performance metrics in a separate goroutine.
	go UpdateMetrics(duration, resp, second)

	return nil
}

// SendRequests sends requests at the specified rate.
func SendRequests(ctx context.Context, url, bearerToken, requestType string, maxRequests int, requestsPerSecond int) PerformanceMetrics {
	// Calculate the rate limit based on the requests per second.
	rateLimit := time.Second / time.Duration(requestsPerSecond)
	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()

	// Initialize the request count.
	var requestCount int32 = 0

	// Wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Log beginning of requests
	fmt.Println("Starting Loadr Requests...")

	// Start sending requests at the specified rate.
	startTime := time.Now()
	for {
		select {
		case <-ticker.C:
			second := int(time.Since(startTime).Seconds())
			if int(requestCount) >= maxRequests {
				return *GetMetrics()
			}
			wg.Add(1)
			go func(u, t, verb string, sec int) {
				defer wg.Done()
				err := makeRequest(verb, u, t, sec)
				if err != nil {
					fmt.Println(err)
					return
				}

				atomic.AddInt32(&requestCount, 1)
			}(url, bearerToken, strings.ToUpper(requestType), second)
		case <-ctx.Done():
			// The context has been cancelled, wait for all requests to finish and return the metrics collected so far
			wg.Wait()
			return *GetMetrics()
		}
	}
}
