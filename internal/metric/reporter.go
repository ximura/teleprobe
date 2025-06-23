package metric

import (
	"context"
	"fmt"
	"log"
)

type Reporter interface {
	Report(ctx context.Context, m MetricValue) error
}

type MetricReporter struct {
	sensorName string
	reporter   Reporter
	in         <-chan MetricValue
}

func NewReporter(sensorName string, r Reporter, in <-chan MetricValue) MetricReporter {
	return MetricReporter{
		sensorName: sensorName,
		reporter:   r,
		in:         in,
	}
}

func (r *MetricReporter) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stoping reporter %s\n", r.sensorName)
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

func (r *MetricReporter) Close() error {
	return nil
}

func (r *MetricReporter) report(ctx context.Context, data MetricValue) error {
	log.Printf("%s: %.4f\n", data.Name, data.Value)
	return r.reporter.Report(ctx, data)
}
