package main

import (
	"context"
	"log"

	"github.com/ximura/teleprobe/api"
	"github.com/ximura/teleprobe/internal/async"
	"github.com/ximura/teleprobe/internal/grpc"
	googleGrpc "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	log.Println("sink")

	server := grpc.NewTelemetrySinkServer()
	grpcService := grpc.NewGRPCService(50051)
	grpcService.Register(func(serviceRegister *googleGrpc.Server) {
		api.RegisterTelemetrySinkServiceServer(serviceRegister, server)
	})

	acts := []async.Runner{
		grpcService,
	}

	if err := async.RunGroup(acts).Run(ctx); err != nil {
		return
	}
}
