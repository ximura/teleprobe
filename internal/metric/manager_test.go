package metric_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ximura/teleprobe/internal/metric"
)

func TestManager_Register(t *testing.T) {
	mgr := metric.NewManager(10)

	err := mgr.Register("cpu_usage", 100)
	require.NoError(t, err)

	err = mgr.Register("cpu_usage", 100)
	assert.ErrorIs(t, err, metric.ErrMetricDuplicated)
}

func TestManager_RunAndReceive(t *testing.T) {
	mgr := metric.NewManager(10)
	err := mgr.Register("test", 10)
	require.NoError(t, err)

	// simulate a short-lived context
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	go func() {
		_ = mgr.Run(ctx) // don't assert, context will cancel
	}()

	// Wait for at least one metric to arrive
	select {
	case <-mgr.Data():
		// received OK
	case <-time.After(300 * time.Millisecond):
		t.Fatal("expected to receive a metric, but timed out")
	}

	err = mgr.Close()
	assert.NoError(t, err)
}

func TestManager_DataChannelClosedAfterClose(t *testing.T) {
	mgr := metric.NewManager(10)
	err := mgr.Register("test", 1)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go func() {
		_ = mgr.Run(ctx)
	}()

	time.Sleep(200 * time.Millisecond) // allow metric loop to emit & exit
	_ = mgr.Close()

	// Ensure channel is closed
	_, ok := <-mgr.Data()
	assert.False(t, ok, "Data channel should be closed after Close()")
}
