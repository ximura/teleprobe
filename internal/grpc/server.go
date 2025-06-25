package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/ximura/teleprobe/api"
	"github.com/ximura/teleprobe/internal/metric"
	"github.com/ximura/teleprobe/internal/sink"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetricHandler interface {
	Handle(ctx context.Context, m *metric.Measurement) error
}

type server struct {
	handler MetricHandler
	api.UnimplementedTelemetrySinkServiceServer
}

func NewTelemetrySinkServer(h MetricHandler) api.TelemetrySinkServiceServer {
	return server{
		handler: h,
	}
}

func (s server) Report(ctx context.Context, m *api.Metric) (*emptypb.Empty, error) {
	err := s.handler.Handle(ctx, &metric.Measurement{
		Name:      m.Name,
		Value:     int(m.Value),
		Timestamp: m.CreatedAt.AsTime(),
	})
	if err != nil {
		if errors.Is(err, sink.ErrRateLimit) {
			log.Printf("dropped metric %s: %v", m.Name, err)
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		}
		log.Printf("handle error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, err
}
