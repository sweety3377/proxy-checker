package config

import "time"

type Proxy struct {
	URL     string        `env:"URL"`
	Timeout time.Duration `env:"TIMEOUT"`
}
