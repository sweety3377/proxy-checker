package proxy_repository

import (
	"context"
	"github.com/rs/zerolog"
	"net"
	"proxy-checker/internal/config"
	"sync"
)

type ProxiesRepository struct {
	wg     *sync.WaitGroup
	logger *zerolog.Logger

	dialer *net.Dialer
}

func New(ctx context.Context, cfg config.Proxy) *ProxiesRepository {
	return &ProxiesRepository{
		wg:     new(sync.WaitGroup),
		logger: zerolog.Ctx(ctx),
		dialer: &net.Dialer{
			Timeout: cfg.Timeout,
		},
	}
}

func (p *ProxiesRepository) StartChecker(proxiesList []string) {
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

	p.logger.Info().Msg("")
}

func (p *ProxiesRepository) checkProxy(ctx context.Context, proxyAddress string) error {
	conn, err := p.dialer.DialContext(ctx, "tcp", proxyAddress)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			p.logger.Error().Err(err).Msg("error closing connection")
		}
	}()

}
