package main

import (
	"context"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"proxy-checker/internal/config"
	"strings"
)

func initProxyList(ctx context.Context, cfg config.Proxy) ([]string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing proxy list")

	req, _ := http.NewRequest(http.MethodGet, cfg.URL, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("error closing response body")
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	list := strings.Split(string(data), "\n")

	return list, nil
}
