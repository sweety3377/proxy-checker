package proxy_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/sweety3377/proxy-checker/internal/config"
	"github.com/sweety3377/proxy-checker/internal/model"
	httpTransport "github.com/sweety3377/proxy-checker/internal/transport/http"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ProxiesStorage struct {
	wg        *sync.WaitGroup
	mx        *sync.Mutex
	logger    *zerolog.Logger
	workersCh chan struct{}
	protocols []string
	results   [][]string
	cfg       config.Proxy
}

func New(ctx context.Context, cfg config.Proxy, maxThreads int) *ProxiesStorage {
	return &ProxiesStorage{
		workersCh: make(chan struct{}, maxThreads),
		results:   make([][]string, 0),
		protocols: []string{"http://", "socks5://"},
		wg:        new(sync.WaitGroup),
		mx:        new(sync.Mutex),
		logger:    zerolog.Ctx(ctx),
		cfg:       cfg,
	}
}

func (p *ProxiesStorage) StartChecker(proxiesList []string) [][]string {
	//p.wg.Add(len(proxiesList))

	start := time.Now().Local()

	var successfullyCount atomic.Uint64
	for ind, proxyAddress := range proxiesList {
		if ind == 50 {
			break
		}

		// Add worker in channel
		p.workersCh <- struct{}{}

		// Increment wait group
		p.wg.Add(1)

		// Start goroutine for check proxy
		go func(proxyAddress string) {
			defer p.wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), p.cfg.Timeout)
			defer cancel()

			records := make([][]string, 0, 2)

			var (
				record []string
				err    error
			)
			for _, protocol := range p.protocols {
				record, err = p.checkProxy(ctx, proxyAddress, protocol)
				if err != nil {
					p.logger.Info().Str("proxy", proxyAddress).Msg("proxy is not active")
				} else {
					successfullyCount.Add(1)
					records = append(records, record)
				}

			}

			// Add record in result
			p.mx.Lock()
			p.results = append(p.results, records...)
			p.mx.Unlock()

			// Remove worker from channel
			<-p.workersCh
		}(proxyAddress)
	}

	// Wait all checks
	p.wg.Wait()
	close(p.workersCh)

	sub := time.Now().Local().Sub(start)

	// Get successfully count
	successfullyCountUint := successfullyCount.Load()

	// Get unsuccessfully count
	unsuccessfullyCountUint := uint64(len(proxiesList)) - successfullyCountUint

	p.logger.Info().
		Dur("dur", sub).
		Uint64("successfully", successfullyCountUint).
		Uint64("unsuccessfully", unsuccessfullyCountUint).
		Msg("successfully checked selected proxies")

	return p.results
}

func (p *ProxiesStorage) checkProxy(ctx context.Context, proxyAddress, scheme string) ([]string, error) {
	// Parse proxy url
	proxyURL, err := url.Parse(scheme + proxyAddress)
	if err != nil {
		return nil, err
	}

	// Get http client with transport on selected proxy url
	httpClient, err := httpTransport.NewHttpClient(proxyURL, p.cfg.Timeout)
	if err != nil {
		return nil, err
	}

	// Create request
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://ip-api.com/json/?fields=61439",
		nil,
	)

	start := time.Now().Local()

	// Do request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	sub := time.Now().Local().Sub(start)

	defer resp.Body.Close()

	// Check on 200 status ok
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is wrong: %v", resp.StatusCode)
	}

	// Get proxy data
	var response model.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	p.logger.Info().
		Str("address", proxyAddress).
		Str("protocol", scheme).
		Str("country", response.RegionName).
		Dur("duration", sub).
		Msg("proxy is active")

	return []string{
		proxyAddress,
		strings.Trim(scheme, "://"),
		response.RegionName,
		sub.String(),
	}, nil
}
