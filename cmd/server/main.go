package main

import (
	"context"
	"log"
	"proxy-checker/internal/config"
	pkgLogger "proxy-checker/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &config.Config{}
	err := config.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger := pkgLogger.New().With().
		Str("project", "proxy-checker").
		Logger()

	runApp(logger.WithContext(ctx), cfg)
}
