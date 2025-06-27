package sink

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ximura/teleprobe/internal/metric"
	"golang.org/x/time/rate"
)

var (
	ErrRateLimit = errors.New("rate limit exceeded")
)

type Marshaller interface {
	Marshal(*metric.Measurement) ([]byte, error)
}

type AppendableBuffer interface {
	Append(line string) error
}

// MetricService handles incoming telemetry metrics by rate-limiting,
// formatting, and forwarding them to a buffered writer. It encapsulates
// the core business logic for the telemetry sink.
type MetricService struct {
	buffer      AppendableBuffer
	marshal     Marshaller
	rateLimiter *rate.Limiter
}

func New(buffer AppendableBuffer, marshal Marshaller, ratelimiter int) MetricService {
	limiter := rate.NewLimiter(rate.Limit(ratelimiter), ratelimiter)
	return MetricService{
		buffer:      buffer,
		marshal:     marshal,
		rateLimiter: limiter,
	}
}

func (s *MetricService) Handle(ctx context.Context, m *metric.Measurement) error {
	data, err := s.marshal.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal metric: %w", err)
	}

	if !s.rateLimiter.AllowN(time.Now(), len(data)) {
		return ErrRateLimit
	}

	line := string(data) + "\n"
	return s.buffer.Append(line)
}
