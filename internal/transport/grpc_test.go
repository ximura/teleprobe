package transport_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ximura/teleprobe/api"
	"github.com/ximura/teleprobe/internal/metric"
	"github.com/ximura/teleprobe/internal/sink"
	"github.com/ximura/teleprobe/internal/transport"
)

const bufSize = 1024 * 1024

type mockHandler struct {
	received *metric.Measurement
	err      error
}

func (h *mockHandler) Handle(_ context.Context, m *metric.Measurement) error {
	h.received = m
	return h.err
}

func setupServer(t *testing.T, h transport.MetricHandler) (api.TelemetrySinkServiceClient, func()) {
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()

	api.RegisterTelemetrySinkServiceServer(server, transport.NewTelemetrySinkServer(h))

	go func() {
		_ = server.Serve(listener)
	}()

	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		return listener.Dial()
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	client := api.NewTelemetrySinkServiceClient(conn)
	return client, func() {
		conn.Close()
		server.Stop()
	}
}

func TestReport_ValidationAndHandlerErrors(t *testing.T) {
	now := timestamppb.Now()
	validMsg := &api.Metric{
		Name:      "test.valid",
		Value:     100,
		CreatedAt: now,
	}

	tests := []struct {
		name     string
		input    *api.Metric
		handler  func() transport.MetricHandler
		wantCode codes.Code
		wantMsg  string
	}{
		{
			name:  "missing name",
			input: &api.Metric{Name: "", Value: 1, CreatedAt: now},
			handler: func() transport.MetricHandler {
				return &mockHandler{}
			},
			wantCode: codes.InvalidArgument,
			wantMsg:  "invalid telemetry message",
		},
		{
			name:  "nil timestamp",
			input: &api.Metric{Name: "x", Value: 1, CreatedAt: nil},
			handler: func() transport.MetricHandler {
				return &mockHandler{}
			},
			wantCode: codes.InvalidArgument,
			wantMsg:  "invalid telemetry message",
		},
		{
			name:  "zero timestamp",
			input: &api.Metric{Name: "x", Value: 1, CreatedAt: &timestamppb.Timestamp{}},
			handler: func() transport.MetricHandler {
				return &mockHandler{}
			},
			wantCode: codes.InvalidArgument,
			wantMsg:  "invalid telemetry message",
		},
		{
			name:  "handler returns rate limit error",
			input: validMsg,
			handler: func() transport.MetricHandler {
				return &mockHandler{err: sink.ErrRateLimit}
			},
			wantCode: codes.ResourceExhausted,
			wantMsg:  "rate limit",
		},
		{
			name:  "handler returns internal error",
			input: validMsg,
			handler: func() transport.MetricHandler {
				return &mockHandler{err: errors.New("unexpected failure")}
			},
			wantCode: codes.Internal,
			wantMsg:  "unexpected failure",
		},
		{
			name:  "valid input succeeds",
			input: validMsg,
			handler: func() transport.MetricHandler {
				return &mockHandler{}
			},
			wantCode: codes.OK,
			wantMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, teardown := setupServer(t, tt.handler())
			defer teardown()

			_, err := client.Report(context.Background(), tt.input)

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)

				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantCode, st.Code())
				assert.Contains(t, st.Message(), tt.wantMsg)
			}
		})
	}
}
