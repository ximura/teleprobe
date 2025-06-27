package sink_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ximura/teleprobe/internal/metric"
	"github.com/ximura/teleprobe/internal/sink"
)

// mockBuffer implements the minimal Append method for testing.
type mockBuffer struct {
	lines []string
	err   error
}

func (b *mockBuffer) Append(line string) error {
	if b.err != nil {
		return b.err
	}
	b.lines = append(b.lines, line)
	return nil
}

type mockMarshaller struct {
	output string
	err    error
}

func (m *mockMarshaller) Marshal(_ *metric.Measurement) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []byte(m.output), nil
}

func TestMetricService_Handle_Success(t *testing.T) {
	buf := &mockBuffer{}
	marsh := &mockMarshaller{output: `{"name":"a","value":1}`}
	s := sink.New(buf, marsh, 1024)

	err := s.Handle(context.Background(), &metric.Measurement{Name: "a", Value: 1})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(buf.lines))
	assert.True(t, strings.HasSuffix(buf.lines[0], "\n"))
}

func TestMetricService_Handle_RateLimitExceeded(t *testing.T) {
	buf := &mockBuffer{}
	marsh := &mockMarshaller{output: strings.Repeat("x", 100)}
	s := sink.New(buf, marsh, 10) // only 10 bytes/sec allowed

	err := s.Handle(context.Background(), &metric.Measurement{Name: "too-big"})
	assert.ErrorIs(t, err, sink.ErrRateLimit)
	assert.Equal(t, 0, len(buf.lines))
}

func TestMetricService_Handle_MarshalError(t *testing.T) {
	buf := &mockBuffer{}
	marsh := &mockMarshaller{err: errors.New("marshal failed")}
	s := sink.New(buf, marsh, 100)

	err := s.Handle(context.Background(), &metric.Measurement{Name: "bad"})
	assert.ErrorContains(t, err, "marshal metric")
}

func TestMetricService_Handle_AppendError(t *testing.T) {
	buf := &mockBuffer{err: errors.New("append fail")}
	marsh := &mockMarshaller{output: `{"name":"a","value":1}`}
	s := sink.New(buf, marsh, 1000)

	err := s.Handle(context.Background(), &metric.Measurement{Name: "x"})
	assert.ErrorContains(t, err, "append fail")
}

func TestMetricService_Handle_BurstWindow(t *testing.T) {
	buf := &mockBuffer{}
	marsh := &mockMarshaller{output: strings.Repeat("x", 5)}
	s := sink.New(buf, marsh, 10) // 10 bytes/sec

	// within burst allowance
	_ = s.Handle(context.Background(), &metric.Measurement{})
	_ = s.Handle(context.Background(), &metric.Measurement{}) // burst allows this

	// now should exceed
	err := s.Handle(context.Background(), &metric.Measurement{})
	assert.ErrorIs(t, err, sink.ErrRateLimit)
}
