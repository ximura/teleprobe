package grpc

import (
	"context"

	"github.com/ximura/teleprobe/api"
	"github.com/ximura/teleprobe/internal/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TelemetryClient struct {
	client api.TelemetrySinkServiceClient
	conn   *grpc.ClientConn
}

func NewTelemetryClient(addr string) (*TelemetryClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := api.NewTelemetrySinkServiceClient(conn)
	return &TelemetryClient{client: client, conn: conn}, nil
}

func (tc *TelemetryClient) Report(ctx context.Context, m metric.MetricValue) error {
	_, err := tc.client.Report(ctx, &api.Metric{
		Name:      m.Name,
		Value:     int32(m.Value),
		CreatedAt: timestamppb.New(m.Timestamp),
	})
	return err
}

func (tc *TelemetryClient) Close() error {
	return tc.conn.Close()
}
