package sensor

import (
	"context"
	"fmt"
	"log"

	"github.com/ximura/teleprobe/internal/metric"
)

type Reporter interface {
	Report(ctx context.Context, m metric.Measurement) error
}

type Sensor struct {
	name     string
	reporter Reporter
	in       <-chan metric.Measurement
}

func New(name string, r Reporter, in <-chan metric.Measurement) Sensor {
	return Sensor{
		name:     name,
		reporter: r,
		in:       in,
	}
}

func (r *Sensor) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stoping reporter %s\n", r.name)
			return ctx.Err()
		case data, ok := <-r.in:
			if !ok {
				return fmt.Errorf("input channel closed")
			}
			if err := r.report(ctx, data); err != nil {
				log.Printf("failed to report metric, %v", err)
			}
		}
	}
}

func (r *Sensor) Close() error {
	return nil
}

func (r *Sensor) report(ctx context.Context, data metric.Measurement) error {
	log.Printf("%s: %d\n", data.Name, data.Value)
	return r.reporter.Report(ctx, data)
}
