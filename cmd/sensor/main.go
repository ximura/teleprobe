package main

import (
	"context"
	"log"

	"github.com/ximura/teleprobe/internal/async"
	"github.com/ximura/teleprobe/internal/metric"
)

func main() {
	ctx := context.Background()
	log.Println("sensor")

	cfg, err := LoadConfig("../../data/sensor.json")
	if err != nil {
		log.Fatalf("failed to read config, %v", err)
	}

	manager := metric.NewManager(10)
	for _, m := range cfg.Metrics {
		manager.Register(m.Name, m.Rate)
	}

	reporter := metric.NewReporter("sensor_1", manager.Data())

	acts := []async.Runner{
		async.NewShutdown(),
		&manager,
		&reporter,
	}

	if err := async.RunGroup(acts).Run(ctx); err != nil {
		return
	}
}
