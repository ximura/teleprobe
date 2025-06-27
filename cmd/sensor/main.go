package main

import (
	"context"
	"log"

	"github.com/ximura/teleprobe/internal/metric"
	"github.com/ximura/teleprobe/internal/sensor"
	"github.com/ximura/teleprobe/internal/transport"
	"github.com/ximura/teleprobe/pkg/async"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("sensor initing")

	cfg, err := sensor.LoadConfig("data/sensor.json")
	if err != nil {
		log.Fatalf("failed to read config, %v", err)
	}

	grpcTransport, err := transport.NewGRPCTransport(cfg.SynAddr)
	if err != nil {
		log.Fatalf("failed to create telemetry client, %v", err)
	}

	manager := metric.NewManager(10)
	for _, m := range cfg.Metrics {
		manager.Register(m.Name, m.Rate)
	}

	reporter := sensor.New(cfg.Name, grpcTransport, manager.Data())

	acts := []async.Runner{
		async.NewShutdown(),
		&manager,
		&reporter,
	}
	log.Println("sensor starting")
	if err := async.RunGroup(acts).Run(ctx); err != nil {
		log.Println("sensor stopped")
		return
	}
}
