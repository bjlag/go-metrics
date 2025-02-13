package interceptor

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	clientIP "github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/logger"
)

const RealIPMeta = "real-ip"

func RealIPInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	md.Set(RealIPMeta, clientIP.GetOutboundIP().String())

	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}

func CheckRealIPMiddleware(trustedSubnet *net.IPNet, logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		realIP := md.Get(RealIPMeta)
		if len(realIP) == 0 {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		ip := net.ParseIP(realIP[0])
		if ip == nil {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		if !trustedSubnet.Contains(ip) {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		return handler(ctx, req)
	}
}
