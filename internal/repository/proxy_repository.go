package proxy_repository

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/sweety3377/proxy-checker/internal/config"
	"github.com/sweety3377/proxy-checker/internal/model"
	httpTransport "github.com/sweety3377/proxy-checker/internal/transport/http"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
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
	//	p.wg.Add(len(proxiesList))

	var successfullyCount atomic.Uint64
	for ind, proxyAddress := range proxiesList {
		if ind == 10 {
			break
		}

		p.wg.Add(1)

		go func(proxyAddress string) {
			defer p.wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), p.cfg.Timeout)
			defer cancel()

			err := p.checkProxy(ctx, proxyAddress)
			if err != nil {
				p.logger.Info().Str("proxy", proxyAddress).Msg("proxy is not active")
			} else {
				successfullyCount.Add(1)
			}
		}(proxyAddress)

		time.Sleep(time.Millisecond * 10)
	}

	p.wg.Wait()

	successfullyCountUint := successfullyCount.Load()
	unsuccessfullyCountUint := uint64(len(proxiesList)) - successfullyCountUint

	p.logger.Info().
		Uint64("successfully", successfullyCountUint).
		Uint64("unsuccessfully", unsuccessfullyCountUint).
		Msg("successfully checked selected proxies")
}

func (p *ProxiesStorage) checkProxy(ctx context.Context, proxyAddress string) error {
	proxyURL, err := url.Parse("http://" + proxyAddress)
	if err != nil {
		return nil
	}

	httpClient, err := httpTransport.NewHttpClient(proxyURL, p.cfg.Timeout)
	if err != nil {
		return err
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://ip-api.com/json/?fields=61439",
		nil,
	)

	start := time.Now().Local()

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	sub := time.Now().Local().Sub(start)

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
		Msg("proxy is active")

	return nil
}
