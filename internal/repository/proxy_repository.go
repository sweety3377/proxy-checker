package proxy_repository

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/sweety3377/proxy-checker/internal/config"
	httpTransport "github.com/sweety3377/proxy-checker/internal/transport/http"
	"github.com/sweety3377/proxy-checker/pkg/model"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type ProxiesStorage struct {
	wg     *sync.WaitGroup
	logger *zerolog.Logger
	cfg    config.Proxy
}

func New(ctx context.Context, cfg config.Proxy) *ProxiesStorage {
	return &ProxiesStorage{
		wg:     new(sync.WaitGroup),
		logger: zerolog.Ctx(ctx),
		cfg:    cfg,
	}
}

func (p *ProxiesStorage) StartChecker(proxiesList []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(proxiesList))

	for _, proxyAddress := range proxiesList {
		go func(proxyAddress string) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), p.cfg.Timeout)
			defer cancel()

			err := p.checkProxy(ctx, proxyAddress)
			if err != nil {
				p.logger.Error().Str("proxy", proxyAddress).Err(err).Msg("error checking proxy")
			}

		}(proxyAddress)
	}

	p.wg.Wait()

	p.logger.Info().Int("len", len(proxiesList)).Msg("successfully checked selected proxies")
}

func (p *ProxiesStorage) checkProxy(ctx context.Context, proxyAddress string) error {
	proxyURL, err := url.Parse(proxyAddress)
	if err != nil {
		return nil
	}

	httpClient, err := httpTransport.NewHttpClient(proxyURL, p.cfg.Timeout)
	if err != nil {
		return err
	}

	requestURL := "http://ip-api.com/json/?fields=61439"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)

	start := time.Now().Local()

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	sub := start.Sub(time.Now().Local())

	defer resp.Body.Close()

	var response model.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	p.logger.Info().
		Str("address", proxyAddress).
		Str("protocol", proxyURL.Scheme).
		Str("country", response.RegionName).
		Dur("duration", sub).
		Msg("[-] Proxy is active")

	return nil
}
