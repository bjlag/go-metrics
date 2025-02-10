package rpc

import (
	"context"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/bjlag/go-metrics/internal/generated/rpc"
	"github.com/bjlag/go-metrics/internal/logger"
)

type Server struct {
	rpc.UnimplementedMetricServiceServer

	log logger.Logger
}

func NewServer(log logger.Logger) *Server {
	return &Server{
		log: log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
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
