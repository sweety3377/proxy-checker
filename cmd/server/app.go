package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/sweety3377/proxy-checker/internal/config"
	proxy_repository "github.com/sweety3377/proxy-checker/internal/repository"
	proxy_service "github.com/sweety3377/proxy-checker/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func runApp(ctx context.Context, cfg *config.Config) {
	// Create channel with signals for gracefully shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Kill, syscall.SIGKILL, syscall.SIGTERM)

	logger := zerolog.Ctx(ctx)

	var (
		proxiesList []string
		err         error
	)

	switch cfg.InputType {
	case "URL":
		proxiesList, err = initProxyListFromURL(ctx, cfg.Proxy)
	case "FILE":
		proxiesList, err = initProxyListFromFile(ctx, cfg.Proxy)
	}
	if err != nil {
		logger.Fatal().Err(err).Str("mode", cfg.Proxy.InputType).Msg("error initializing proxy")
	}

	logger.Info().Int("len", len(proxiesList)).Msg("successfully loaded proxies list")

	proxyRepository := proxy_repository.New(ctx, cfg.Proxy)
	proxyService := proxy_service.New(proxyRepository)

	// Start check
	go proxyService.StartChecker(proxiesList)

	sig := <-shutdownCh
	logger.Info().Str("signal", sig.String()).Msg("shutdown signal receive")
}
