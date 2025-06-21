package teleprobe

import (
	"context"
	"fmt"
	"log"
)

type Reporter struct {
	sensorName string
	in         <-chan MetricValue
}

func NewReporter(sensorName string, in <-chan MetricValue) Reporter {
	return Reporter{
		sensorName: sensorName,
		in:         in,
	}
}

func (r *Reporter) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stoping reporter %s\n", r.sensorName)
			return ctx.Err()
		case data, ok := <-r.in:
			if !ok {
				return fmt.Errorf("input channel closed")
			}
			r.report(data)
		}
	}
}

func (r *Reporter) Close() error {
	return nil
}

func (r *Reporter) report(data MetricValue) {
	log.Printf("%s: %.4f\n", data.Name, data.Value)
}
