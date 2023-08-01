package config

type Runtime struct {
	UseCPUs        int `env:"USE_CPUS" default:"4"`
	MaxThreads     int `env:"MAX_THREADS" default:"1000"`
	CurrentThreads int `env:"CURRENT_THREADS" default:"10"`
}
