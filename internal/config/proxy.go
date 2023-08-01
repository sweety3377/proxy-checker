package config

import "time"

type Proxy struct {
	URL       string        `env:"URL"`
	File      string        `env:"FILE"`
	InputType string        `env:"INPUT_TYPE" default:"URL"`
	Timeout   time.Duration `env:"TIMEOUT" default:"1m"`
}
