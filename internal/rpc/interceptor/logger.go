package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-metrics/internal/logger"
)

func LoggerClientInterceptor(log logger.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		log.WithField("method", method).
			Info("Send RPC request")
		return err
	}
}

func LoggerServerInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		log.
			WithField("method", info.FullMethod).
			WithField("code", status.Code(err)).
			Info("Got RPC request")

		return resp, err
	}
}
