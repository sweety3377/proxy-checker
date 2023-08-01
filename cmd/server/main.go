package main

import (
	"context"
	"github.com/sweety3377/proxy-checker/internal/config"
	pkgLogger "github.com/sweety3377/proxy-checker/pkg/logger"
	"log"
	"runtime"
	"runtime/debug"
)

func main() {
	// Base context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config [.env]
	cfg := &config.Config{}
	err := config.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create app logger
	logger := pkgLogger.New().With().
		Str("project", "proxy-checker").
		Logger()

	// Set max using cpus on instance
	if cfg.Runtime.UseCPUs != 0 {
		runtime.GOMAXPROCS(cfg.Runtime.UseCPUs)
	}

	// Set max threads for instance
	// 5 = it's default threads number for instance working
	debug.SetMaxThreads(5 + cfg.Runtime.MaxThreads)

	// Run app
	runApp(logger.WithContext(ctx), cfg)
}
