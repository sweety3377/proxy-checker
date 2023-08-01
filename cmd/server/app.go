package main

import (
	"context"
	"encoding/csv"
	"github.com/rs/zerolog"
	"github.com/sweety3377/proxy-checker/internal/config"
	proxy_repository "github.com/sweety3377/proxy-checker/internal/repository"
	proxy_service "github.com/sweety3377/proxy-checker/internal/service"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func runApp(ctx context.Context, cfg *config.Config) {
	// Create channel with signals for gracefully shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Kill, syscall.SIGKILL, syscall.SIGTERM)

	logger := zerolog.Ctx(ctx)

	// Get proxies list be loaded depending on which file_input
	var (
		proxiesList []string
		err         error
	)

	switch cfg.InputType {
	// Get proxies from url
	case "URL":
		proxiesList, err = initProxyListFromURL(ctx, cfg.Proxy)
	// Get proxies from file
	case "FILE":
		proxiesList, err = initProxyListFromFile(ctx, cfg.Proxy)
	}
	if err != nil {
		logger.Fatal().Err(err).Str("mode", cfg.Proxy.InputType).Msg("error initializing proxy")
	}

	logger.Info().Int("len", len(proxiesList)).Msg("successfully loaded proxies list")

	proxyRepository := proxy_repository.New(ctx, cfg.Proxy, cfg.Runtime.MaxThreads)
	proxyService := proxy_service.New(proxyRepository)

	// Start check
	go func() {
		results := proxyService.StartChecker(proxiesList)

		timestamp := strconv.Itoa(int(time.Now().Unix()))

		fileName := filepath.Join("results", "results"+timestamp+".csv")
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			logger.Error().Err(err).Msg("error creating file for results")
		}

		w := csv.NewWriter(file)

		err = w.WriteAll(results)
		if err != nil {
			logger.Error().Err(err).Msg("error saving results in csv file")
		}
	}()

	// Wait shutdown signal
	sig := <-shutdownCh
	logger.Info().Str("signal", sig.String()).Msg("shutdown signal receive")
}
