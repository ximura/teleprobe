package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MetricConfig struct {
	Name string `json:"name,omitempty"`
	Rate int    `json:"rate,omitempty"`
}

type Config struct {
	Metrics []MetricConfig `json:"metrics,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}
