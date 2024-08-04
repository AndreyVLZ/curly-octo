package interceptor

import (
	"context"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type iAuthClient interface {
	Refresh(ctx context.Context) error // не реализовано
	Token() string
}

type AuthInterceptor struct {
	authClient iAuthClient
}

func NewAuthInterceptor(authClient iAuthClient) *AuthInterceptor {
	return &AuthInterceptor{
		authClient: authClient,
	}
}

func (inter *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "token", inter.authClient.Token())

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (inter *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "token", inter.authClient.Token())

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func getMethod(method string) string {
	arr := strings.Split(method, string(os.PathSeparator))

	return arr[len(arr)-1]
}
