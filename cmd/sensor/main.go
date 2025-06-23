package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ximura/teleprobe/internal/async"
	"github.com/ximura/teleprobe/internal/grpc"
	"github.com/ximura/teleprobe/internal/metric"
)

type MetricConfig struct {
	Name string `json:"name,omitempty"`
	Rate int    `json:"rate,omitempty"`
}

type Config struct {
	Addr    string         `json:"addr,omitempty"`
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

func main() {
	ctx := context.Background()
	log.Println("sensor")

	cfg, err := LoadConfig("data/sensor.json")
	if err != nil {
		log.Fatalf("failed to read config, %v", err)
	}
	client, err := grpc.NewTelemetryClient(cfg.Addr)
	if err != nil {
		log.Fatalf("failed to create telemetry client, %v", err)
	}

	manager := metric.NewManager(10)
	for _, m := range cfg.Metrics {
		manager.Register(m.Name, m.Rate)
	}

	reporter := metric.NewReporter("sensor_1", client, manager.Data())

	acts := []async.Runner{
		async.NewShutdown(),
		&manager,
		&reporter,
	}

	if err := async.RunGroup(acts).Run(ctx); err != nil {
		return
	}
}
