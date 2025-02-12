package rpc

import (
	"context"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-metrics/internal/generated/rpc"
	"github.com/bjlag/go-metrics/internal/logger"
)

const (
	UpdatesMethodName = "updates"
)

type Server struct {
	rpc.UnimplementedMetricServiceServer

	methods map[string]any
	log     logger.Logger
}

func NewServer(log logger.Logger) *Server {
	return &Server{
		methods: make(map[string]any),
		log:     log,
	}
}

func (s *Server) AddMethod(name string, method any) {
	s.methods[name] = method
}

func (s *Server) Start(ctx context.Context) error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	rpc.RegisterMetricServiceServer(grpcServer, s)

	s.log.Info("Starting gRPC server")

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return grpcServer.Serve(listen)
	})

	g.Go(func() error {
		<-gCtx.Done()

		s.log.Info("Shutting down gRPC server")
		grpcServer.GracefulStop()

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Updates(ctx context.Context, in *rpc.UpdatesIn) (*rpc.UpdatesOut, error) {
	method, ok := s.methods[UpdatesMethodName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown method: %s", UpdatesMethodName)
	}

	return method.(func(context.Context, *rpc.UpdatesIn) (*rpc.UpdatesOut, error))(ctx, in)
}
