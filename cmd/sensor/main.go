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

	manager := metric.NewManager(10)
	manager.Register("metric_1", 10)
	manager.Register("metric_2", 5)

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
