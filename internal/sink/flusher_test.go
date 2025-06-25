package sink_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ximura/teleprobe/internal/sink"
)

// mockFlushable is a thread-safe mock implementing Flushable.
type mockFlushable struct {
	count int32
}

func (m *mockFlushable) Flush() error {
	atomic.AddInt32(&m.count, 1)
	return nil
}

func (m *mockFlushable) Calls() int {
	return int(atomic.LoadInt32(&m.count))
}

func TestFlusher_PeriodicFlush(t *testing.T) {
	mock := &mockFlushable{}
	flusher := sink.NewFlusher(mock, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 220*time.Millisecond)
	defer cancel()

	go flusher.Run(ctx)

	// Wait for flushes to happen
	time.Sleep(250 * time.Millisecond)

	require.GreaterOrEqual(t, mock.Calls(), 4, "should flush ~4-5 times")
}
