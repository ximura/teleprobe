package sensor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ximura/teleprobe/internal/metric"
)

type Transport interface {
	Send(ctx context.Context, m metric.Measurement) error
}

type Sensor struct {
	name      string
	transport Transport
	in        <-chan metric.Measurement
}

func New(name string, t Transport, in <-chan metric.Measurement) Sensor {
	return Sensor{
		name:      name,
		transport: t,
		in:        in,
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
			if err := sendWithRetry(ctx, r.transport, data); err != nil {
				log.Printf("failed to send metric, %v", err)
			}
		}
	}
}

func (r *Sensor) Close() error {
	return nil
}

func sendWithRetry(ctx context.Context, transport Transport, data metric.Measurement) error {
	operation := func() error {
		return transport.Send(ctx, data)
	}

	// Customize backoff strategy
	expBackoff := backoff.NewExponentialBackOff(
		backoff.WithMaxElapsedTime(3*time.Second),
		backoff.WithMaxInterval(time.Second),
		backoff.WithInitialInterval(100*time.Millisecond),
	)

	return backoff.Retry(operation, backoff.WithContext(expBackoff, ctx))
}
