package main

import (
	"context"
	"github.com/sweety3377/proxy-checker/internal/config"
	pkgLogger "github.com/sweety3377/proxy-checker/pkg/logger"
	"log"
	"os"
	"runtime"
)

func init() {
	os.Mkdir("data", os.ModePerm)
	os.Mkdir("results", os.ModePerm)
}

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

	// Run app
	runApp(logger.WithContext(ctx), cfg)
}
