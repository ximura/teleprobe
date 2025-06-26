package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCService struct {
	listAddr string

	server *grpc.Server
}

// NewGRPCService creates a new GRPCService instance.
func NewGRPCService(addr string) *GRPCService {
	return &GRPCService{
		listAddr: addr,
		server:   grpc.NewServer(),
	}
}

func (gs *GRPCService) Run(ctx context.Context) error {
	lc := &net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", gs.listAddr)
	if err != nil {
		return err
	}
	defer lis.Close()
	reflection.Register(gs.server)

	// Shutdown the server when the context is canceled
	go func() {
		<-ctx.Done()
		gs.Close()
	}()

	log.Println("Serving GRPC Service on " + gs.listAddr)
	if err := gs.server.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (gs *GRPCService) Close() error {
	gs.server.GracefulStop()
	return nil
}

func (gs *GRPCService) Register(serviceRegister func(server *grpc.Server)) {
	serviceRegister(gs.server)
}
