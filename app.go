package main

import (
	"changeme/lib"
	"context"
	"sync"
)

// LoadTestConfig struct
type LoadTestConfig struct {
	URL            string
	RequestsPerSec int
	MaxRequests    int
	RequestType    string
	BearerToken    string
}

// App struct
type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	mux    sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// LoadTest performs a load test based on the given configuration
func (a *App) LoadTest(config LoadTestConfig) (*lib.PerformanceMetrics, error) {
	a.mux.Lock()
	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.mux.Unlock()

	// Reset the metrics at the start of each test
	lib.ResetMetrics()

	// Pass the context to the SendRequests function
	metrics := lib.SendRequests(a.ctx, config.URL, "", config.RequestType, config.MaxRequests, config.RequestsPerSec)

	return &metrics, nil
}

// Cancel cancels the ongoing load test
func (a *App) Cancel() {
	a.mux.Lock()
	if a.cancel != nil {
		a.cancel()
	}
	a.mux.Unlock()
}
