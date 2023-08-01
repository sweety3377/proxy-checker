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
	"sync"
	"sync/atomic"
	"time"
)

type ProxiesResults struct {
	// Results slices for csv
	data [][]string

	// Locker for structure
	mx *sync.Mutex
}

type ProxiesStorage struct {
	// WaitGroup, it's counter for goroutines
	wg *sync.WaitGroup

	// Console logger
	logger *zerolog.Logger

	// Channel for workers (for threads balancing)
	// Buffered channel when len = len(proxies_list)
	workersCh chan struct{}

	// Protocols (ex: []string{"http", "socks5", "socks4a"} and etc)
	protocols []string

	// Checks results
	results ProxiesResults

	// Proxies config
	cfg config.Proxy
}

func New(ctx context.Context, cfg config.Proxy, maxThreads int) *ProxiesStorage {
	return &ProxiesStorage{
		results: ProxiesResults{
			data: make([][]string, 0),
			mx:   new(sync.Mutex),
		},
		protocols: []string{
			"http",
			"socks5",
		},
		workersCh: make(chan struct{}, maxThreads),
		wg:        new(sync.WaitGroup),
		logger:    zerolog.Ctx(ctx),
		cfg:       cfg,
	}
}

func (p *ProxiesStorage) StartChecker(proxiesList []string) [][]string {
	// Increment wait group
	p.wg.Add(len(proxiesList))

	start := time.Now().Local()

	var successfullyCount atomic.Uint64
	for _, proxyAddress := range proxiesList {
		// Add worker in channel
		p.workersCh <- struct{}{}

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
			p.results.mx.Lock()
			p.results.data = append(p.results.data, records...)
			p.results.mx.Unlock()

			// Remove worker from channel
			<-p.workersCh
		}(proxyAddress)

		time.Sleep(time.Millisecond * 50)
	}

	// Wait all checks
	p.wg.Wait()

	// Close workers channel
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

	return p.results.data
}

func (p *ProxiesStorage) checkProxy(ctx context.Context, proxyAddress, scheme string) ([]string, error) {
	// Parse proxy url
	proxyURL, err := url.Parse(scheme + "://" + proxyAddress)
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
		scheme,
		response.RegionName,
		sub.String(),
	}, nil
}
