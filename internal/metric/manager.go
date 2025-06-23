package metric

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/ximura/teleprobe/internal/async"
)

var (
	ErrMetricDuplicated = errors.New("metric already registred")
)

type Manager struct {
	metrics map[string]*Metric
	out     chan MetricValue
	wg      sync.WaitGroup
}

func NewManager(bufferSize int) Manager {
	return Manager{
		metrics: make(map[string]*Metric),
		out:     make(chan MetricValue, bufferSize),
	}
}

func (m *Manager) Register(name string, rate int) error {
	_, ok := m.metrics[name]
	if ok {
		return ErrMetricDuplicated
	}

	metric := NewMetric(name, rate)
	m.metrics[name] = &metric

	return nil
}

func (m *Manager) Run(ctx context.Context) error {
	acts := make([]async.Runner, 0, 2*len(m.metrics))
	for i := range m.metrics {
		metric := m.metrics[i]
		m.wg.Add(1)
		f := async.Func(func(ctx context.Context) error {
			defer m.wg.Done()

			for value := range metric.Data() {
				select {
				case m.out <- value:
				case <-ctx.Done():
					log.Println("Stoping manager")
					return ctx.Err()
				}
			}

			return nil
		})

		acts = append(acts, f, m.metrics[i])
	}

	return async.RunGroup(acts).Run(ctx)
}

func (m *Manager) Close() error {
	m.wg.Wait()
	close(m.out)
	return nil
}

func (m *Manager) Data() <-chan MetricValue {
	return m.out
}
