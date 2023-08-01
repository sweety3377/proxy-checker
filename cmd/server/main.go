package main

import (
	"context"
	"github.com/sweety3377/proxy-checker/internal/config"
	pkgLogger "github.com/sweety3377/proxy-checker/pkg/logger"
	"log"
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
