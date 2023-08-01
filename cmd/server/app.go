package main

import (
	"context"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"proxy-checker/internal/config"
	"syscall"
)

func runApp(ctx context.Context, cfg *config.Config) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Kill, syscall.SIGKILL, syscall.SIGTERM)

	logger := zerolog.Ctx(ctx)

	proxiesList, err := initProxyList(ctx, cfg.Proxy)
	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing proxies list")
	}

	logger.Info().Int("len", len(proxiesList)).Msg("successfully loaded proxies list")

}
