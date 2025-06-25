package sink_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ximura/teleprobe/internal/metric"
	"github.com/ximura/teleprobe/internal/sink"
	"golang.org/x/time/rate"
)

// mockBuffer implements the minimal Append method for testing.
type mockBuffer struct {
	lines []string
}

func (b *mockBuffer) Append(line string) error {
	b.lines = append(b.lines, line)
	return nil
}

type fakeFormatter struct{}

func (f *fakeFormatter) Marshal(m *metric.Measurement) ([]byte, error) {
	return []byte(`1`), nil
}

func TestMetricService_Handle_Success(t *testing.T) {
	ctx := context.Background()
	buf := &mockBuffer{}
	limiter := rate.NewLimiter(rate.Inf, 0) // no rate limit
	service := sink.New(buf, &sink.JSONFormatter{}, limiter)

	m := &metric.Measurement{
		Name:      "cpu",
		Value:     42,
		Timestamp: time.Now(),
	}

	err := service.Handle(ctx, m)
	require.NoError(t, err)

	require.Len(t, buf.lines, 1)
	var result metric.Measurement
	err = json.Unmarshal([]byte(buf.lines[0]), &result)
	require.NoError(t, err)
	require.Equal(t, "cpu", result.Name)
	require.Equal(t, 42, result.Value)
}

func TestMetricService_Handle_TooFast(t *testing.T) {
	ctx := context.Background()
	buf := &mockBuffer{}
	limiter := rate.NewLimiter(rate.Every(10*time.Second), 1) // allow 1 token per 10s
	service := sink.New(buf, &fakeFormatter{}, limiter)

	m := &metric.Measurement{
		Name:      "mem",
		Value:     99,
		Timestamp: time.Now(),
	}

	// 1st call should succeed
	err := service.Handle(ctx, m)
	require.NoError(t, err)

	// 2nd call should be rate-limited
	err = service.Handle(ctx, m)
	require.Error(t, err)
	require.ErrorIs(t, err, sink.ErrRateLimit)
}
