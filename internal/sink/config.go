package sink

import (
	"time"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	BindAddr      string        `env:"BIND_ADDR"        envDefault:":50051"`
	LogFilePath   string        `env:"LOG_FILE"         envDefault:"telemetry.log"`
	BufferSize    int           `env:"BUFFER_SIZE"      envDefault:"10240"` // in bytes
	FlushInterval time.Duration `env:"FLUSH_INTERVAL"   envDefault:"10s"`
	RateLimit     int           `env:"RATE_LIMIT"       envDefault:"1048576"` // bytes/sec
}

func LoadConfig(path string) (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
