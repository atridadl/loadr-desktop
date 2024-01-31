package main

import (
	"changeme/lib"
	"context"
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
	ctx context.Context
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
	metrics := lib.SendRequests(config.URL, config.BearerToken, config.RequestType, config.MaxRequests, config.RequestsPerSec)

	return &metrics, nil
}
