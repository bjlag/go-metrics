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
	"github.com/bjlag/go-metrics/internal/rpc/interceptor"
	"github.com/bjlag/go-metrics/internal/securety/signature"
)

const (
	UpdatesMethodName = "updates"
)

type Server struct {
	rpc.UnimplementedMetricServiceServer

	methods       map[string]any
	addr          string
	trustedSubnet *net.IPNet
	singManager   *signature.SignManager
	log           logger.Logger
}

func NewServer(addr string, trustedSubnet *net.IPNet, singManager *signature.SignManager, log logger.Logger) *Server {
	return &Server{
		methods: make(map[string]any),

		addr:          addr,
		trustedSubnet: trustedSubnet,
		singManager:   singManager,
		log:           log,
	}
}

func (s *Server) AddMethod(name string, method any) {
	s.methods[name] = method
}

func (s *Server) Start(ctx context.Context) error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerServerInterceptor(s.log),
			interceptor.CheckRealIPServerMiddleware(s.trustedSubnet),
			interceptor.CheckSignatureServerInterceptor(s.singManager),
		),
	)
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
