package interceptor

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-metrics/internal/securety/signature"
)

const headerHash = "HashSHA256"

func SignatureClientInterceptor(sign *signature.SignManager) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if sign.Enable() {
			jsonb, err := json.Marshal(req)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to marshal request: %s", err.Error())
			}

			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.New(nil)
			}

			md.Set(headerHash, sign.Sing(jsonb))

			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func CheckSignatureServerInterceptor(sign *signature.SignManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if sign.Enable() {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.FailedPrecondition, "don't have signature")
			}

			signHash := md.Get(headerHash)
			if len(signHash) == 0 {
				return nil, status.Error(codes.FailedPrecondition, "don't have signature")
			}

			jsonb, err := json.Marshal(req)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to marshal request: %s", err.Error())
			}

			isValid, _ := sign.Verify(jsonb, signHash[0])
			if !isValid {
				return nil, status.Error(codes.FailedPrecondition, "invalid signature")
			}
		}

		return handler(ctx, req)
	}
}
