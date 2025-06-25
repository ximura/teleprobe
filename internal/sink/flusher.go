package sink

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Flushable interface {
	Flush() error
}

// Flusher periodically calls Flush on a Flushable target at a fixed interval.
// It stops and performs a final flush when the provided context is canceled.
type Flusher struct {
	target   Flushable
	interval time.Duration
}

func NewFlusher(target Flushable, interval time.Duration) *Flusher {
	return &Flusher{
		target:   target,
		interval: interval,
	}
}

func (f *Flusher) Run(ctx context.Context) error {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("flusher time to flush")
			if err := f.target.Flush(); err != nil {
				return fmt.Errorf("flush failed: %w", err)
			}
		case <-ctx.Done():
			// final flush
			if err := f.target.Flush(); err != nil {
				return fmt.Errorf("flush failed: %w", err)
			}
			return ctx.Err()
		}
	}
}

func (f *Flusher) Close() error {
	return nil
}
