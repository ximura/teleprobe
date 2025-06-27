package sensor

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/env/v9"
)

type MetricConfig struct {
	Name string `json:"name,omitempty"`
	Rate int    `json:"rate,omitempty"`
}

type Config struct {
	SynAddr    string         `env:"SINK_ADDR"        envDefault:"localhost:50051"`
	ConfigFile string         `env:"CONFIG_FILE"      envDefault:"cfg/sensor.json"`
	Name       string         `env:"SENSOR_NAME"      envDefault:"sensor"`
	Metrics    []MetricConfig `json:"metrics,omitempty"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	f, err := os.Open(cfg.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}
