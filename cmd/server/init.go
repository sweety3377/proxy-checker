package main

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/sweety3377/proxy-checker/internal/config"
	"github.com/sweety3377/proxy-checker/internal/model"
	"net/http"
	"os"
	"strings"
)

func initProxyListFromFile(ctx context.Context, cfg config.Proxy) ([]string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing proxy list from file")

	data, err := os.ReadFile(cfg.File)
	if err != nil {
		return nil, err
	}

	dataStr := strings.ReplaceAll(string(data), "\r", "")
	dataStr = strings.TrimSpace(dataStr)
	list := strings.Split(dataStr, "\n")

	return list, nil
}

func initProxyListFromURL(ctx context.Context, cfg config.Proxy) ([]string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing proxy list from url")

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

	var response model.CheckedProxyResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	proxiesList := make([]string, 0, len(response))
	for _, resp := range response {
		proxiesList = append(proxiesList, resp.Addr)
	}

	return proxiesList, nil
}
