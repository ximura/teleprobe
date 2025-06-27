package transport

import (
	"context"

	"github.com/ximura/teleprobe/api"
	"github.com/ximura/teleprobe/internal/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCTransport struct {
	client api.TelemetrySinkServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCTransport(addr string) (*GRPCTransport, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := api.NewTelemetrySinkServiceClient(conn)
	return &GRPCTransport{client: client, conn: conn}, nil
}

func (t *GRPCTransport) Send(ctx context.Context, data metric.Measurement) error {
	_, err := t.client.Report(ctx, &api.Metric{
		Name:      data.Name,
		Value:     int32(data.Value),
		CreatedAt: timestamppb.New(data.Timestamp),
	})
	return err
}
