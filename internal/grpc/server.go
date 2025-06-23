package grpc

import (
	"context"
	"log"
	"time"

	"github.com/ximura/teleprobe/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	api.UnimplementedTelemetrySinkServiceServer
}

func NewTelemetrySinkServer() api.TelemetrySinkServiceServer {
	return server{}
}

func (s server) Report(ctx context.Context, m *api.Metric) (*emptypb.Empty, error) {
	log.Printf("Report: %s %d %s\n", m.Name, m.Value, m.CreatedAt.AsTime().Format(time.RFC3339))
	return &emptypb.Empty{}, nil
}
