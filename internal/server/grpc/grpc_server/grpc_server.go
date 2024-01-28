package grpcserver

import (
	"context"
	"errors"
	"net"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	api "github.com/ilyakharev/url-short/proto"
)

type GrpcServer struct {
	port   string
	server *grpc.Server
	logger *zap.Logger
}
type GRPCHandlers interface {
	api.GrpcHandlerServer
}

func (server *GrpcServer) ListenAndServe(_ context.Context) error {
	grpcListener, err := net.Listen("tcp", ":"+server.port)
	if err != nil {
		return err
	}
	return server.server.Serve(grpcListener)
}

func (server *GrpcServer) Run(ctx context.Context) error {
	ch := make(chan error)
	go func(ch chan error) {
		if err := server.ListenAndServe(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.logger.Error("Failed to listen", zap.Error(err))
			ch <- err
		}
	}(ch)
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():

		server.logger.Info("Shutting down server gracefully")

		server.server.GracefulStop()

		return nil
	}
}

func New(port string, grpcHandlers GRPCHandlers,
	logger *zap.Logger,
) *GrpcServer {
	grpcServ := grpc.NewServer()
	api.RegisterGrpcHandlerServer(grpcServ, grpcHandlers)

	return &GrpcServer{
		port:   port,
		server: grpcServ,
		logger: logger,
	}
}
