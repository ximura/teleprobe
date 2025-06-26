package metric_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ximura/teleprobe/internal/metric"
)

func TestMetric_EmitsValues(t *testing.T) {
	m := metric.NewMetric("test", 10) // 10 per second = every 100ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go func() {
		_ = m.Run(ctx)
	}()

	select {
	case val := <-m.Data():
		assert.Equal(t, "test", val.Name)
		assert.NotZero(t, val.Value)
		assert.WithinDuration(t, time.Now().UTC(), val.Timestamp, time.Second)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected to receive a metric within 200ms")
	}
}

func TestMetric_StopsOnContextCancel(t *testing.T) {
	m := metric.NewMetric("cancel-test", 100)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := m.Run(ctx)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.LessOrEqual(t, time.Since(start), 200*time.Millisecond)
}

func TestMetric_CloseClosesChannel(t *testing.T) {
	m := metric.NewMetric("close-test", 1)
	err := m.Close()
	assert.NoError(t, err)

	_, ok := <-m.Data()
	assert.False(t, ok, "expected data channel to be closed after Close()")
}
