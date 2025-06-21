package main

import (
	"context"
	"log"
	"time"

	"github.com/ximura/teleprobe/internal/async"
	"github.com/ximura/teleprobe/internal/teleprobe"
)

func main() {
	ctx := context.Background()
	log.Println("sink")
	metric1 := teleprobe.NewMetric("metric_1", time.Second)
	metric2 := teleprobe.NewMetric("metric_2", 2*time.Second)

	reporter := teleprobe.NewReporter("sensor_1", metric1.Data())

	acts := []async.Runner{
		async.NewShutdown(),
		&metric1,
		&metric2,
		&reporter,
	}

	if err := async.RunGroup(acts).Run(ctx); err != nil {
		return
	}
}
