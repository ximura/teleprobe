package metric

import (
	"context"
	"log"
	"math/rand/v2"
	"time"
)

type Measurement struct {
	Name      string
	Value     int
	Timestamp time.Time
}

type Metric struct {
	name     string
	duration time.Duration
	out      chan Measurement
}

func NewMetric(name string, rate int) Metric {
	interval := time.Duration(float64(time.Second) / float64(rate))
	return Metric{
		name:     name,
		duration: interval,
		out:      make(chan Measurement),
	}
}

func (m *Metric) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stoping metric %s\n", m.name)
			return ctx.Err()
		case <-ticker.C:
			value := rand.Int()
			select {
			case m.out <- Measurement{
				Name:      m.name,
				Value:     value,
				Timestamp: time.Now().UTC(),
			}:
			default:
				log.Printf("Dropping metric values for %s\n", m.name)
			}
		}
	}
}

func (m *Metric) Close() error {
	close(m.out)
	return nil
}

func (m *Metric) Data() <-chan Measurement {
	return m.out
}
