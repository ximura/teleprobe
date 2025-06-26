package sensor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ximura/teleprobe/internal/metric"
	"github.com/ximura/teleprobe/internal/sensor"
)

// mockTransport implements the Transport interface for testing
type mockTransport struct {
	sendCalls []metric.Measurement
	failTimes int // how many times to fail before success
	calls     int
}

func (m *mockTransport) Send(ctx context.Context, data metric.Measurement) error {
	m.calls++
	m.sendCalls = append(m.sendCalls, data)

	if m.calls <= m.failTimes {
		return errors.New("mock failure")
	}
	return nil
}

func TestSensor_Run_SuccessfulSend(t *testing.T) {
	in := make(chan metric.Measurement, 1)
	in <- metric.Measurement{
		Name:      "test",
		Value:     123,
		Timestamp: time.Now(),
	}
	close(in)

	mock := &mockTransport{}
	s := sensor.New("test-sensor", mock, in)

	ctx := context.Background()
	err := s.Run(ctx)
	assert.Error(t, err, "should return when input channel is closed")

	require.Len(t, mock.sendCalls, 1)
	assert.Equal(t, "test", mock.sendCalls[0].Name)
}

func TestSensor_Run_RetryOnFailure(t *testing.T) {
	in := make(chan metric.Measurement, 1)
	in <- metric.Measurement{Name: "retry-me", Value: 42, Timestamp: time.Now()}
	close(in)

	mock := &mockTransport{failTimes: 2} // fail twice before succeeding
	s := sensor.New("retry-sensor", mock, in)

	ctx := context.Background()
	err := s.Run(ctx)
	assert.Error(t, err)
	assert.GreaterOrEqual(t, mock.calls, 3, "should retry at least 3 times")
}

func TestSensor_Run_ContextCancel(t *testing.T) {
	in := make(chan metric.Measurement)

	mock := &mockTransport{}
	s := sensor.New("cancel-sensor", mock, in)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := s.Run(ctx)
	assert.Error(t, err)
	assert.Less(t, time.Since(start), time.Second, "should exit quickly on context cancel")
}
