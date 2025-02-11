package rpc

import (
	"context"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/bjlag/go-metrics/internal/generated/rpc"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
)

type Server struct {
	rpc.UnimplementedMetricServiceServer

	repo   repo
	backup backup
	log    log
}

func NewServer(repo repo, backup backup, log log) *Server {
	return &Server{
		repo:   repo,
		backup: backup,
		log:    log,
	}
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
	s.log.Info("Received Updates")

	if len(in.Metrics) == 0 {
		return nil, nil
	}

	gauges := make([]storage.Gauge, 0, len(in.Metrics))
	counters := make([]storage.Counter, 0, len(in.Metrics))

	for _, m := range in.Metrics {
		switch m.Type {
		case model.TypeGauge:
			if m.Value == nil {
				s.log.Info("Invalid value")
				continue
			}

			gauges = append(gauges, storage.Gauge{
				ID:    m.Id,
				Value: *m.Value,
			})
		case model.TypeCounter:
			if m.Delta == nil {
				s.log.Info("Invalid value")
				continue
			}

			counters = append(counters, storage.Counter{
				ID:    m.Id,
				Value: *m.Delta,
			})
		}
	}

	err := s.repo.SetGauges(ctx, gauges)
	if err != nil {
		s.log.WithError(err).Error("Failed to save gauges")
		return nil, err
	}

	err = s.repo.AddCounters(ctx, counters)
	if err != nil {
		s.log.WithError(err).Error("Failed to save counters")
		return nil, err
	}

	err = s.backup.Create(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to backup data")
	}

	return &rpc.UpdatesOut{}, nil
}
