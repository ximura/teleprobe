package main

import (
	"context"
	"log"

	"github.com/ximura/teleprobe/api"
	"github.com/ximura/teleprobe/internal/grpc"
	"github.com/ximura/teleprobe/internal/sink"
	"github.com/ximura/teleprobe/pkg/async"
	googleGrpc "google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("sink initing")

	cfg, err := sink.LoadConfig("data/sink.json")
	if err != nil {
		log.Fatalf("failed to read config, %v", err)
	}

	buffer, err := sink.NewBuffer(cfg.LogFilePath, cfg.BufferSize)
	if err != nil {
		log.Fatalf("failed to read config, %v", err)
	}

	service := sink.New(buffer, &sink.JSONFormatter{}, cfg.RateLimit)
	flusher := sink.NewFlusher(buffer, cfg.FlushInterval)

	server := grpc.NewTelemetrySinkServer(&service)
	grpcService := grpc.NewGRPCService(cfg.BindAddr)
	grpcService.Register(func(serviceRegister *googleGrpc.Server) {
		api.RegisterTelemetrySinkServiceServer(serviceRegister, server)
	})

	acts := []async.Runner{
		async.NewShutdown(),
		flusher,
		grpcService,
	}

	log.Println("sink starting")
	if err := async.RunGroup(acts).Run(ctx); err != nil {
		log.Println("sink stopped")
		return
	}
}
