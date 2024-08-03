package interceptor

import (
	"context"
	"os"
	"strings"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	"github.com/AndreyVLZ/curly-octo/server/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwt    *jwt.JWT
	metods map[string]struct{}
}

type wrapServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *wrapServerStream) Context() context.Context { return s.ctx }

func New(jwt *jwt.JWT) *AuthInterceptor {
	return &AuthInterceptor{
		jwt: jwt,
		metods: map[string]struct{}{
			"Registration": struct{}{},
			"Login":        struct{}{},
		},
	}
}

func (inter *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := inter.metods[getMethod(info.FullMethod)]; ok {
			return handler(ctx, req)
		}

		md, err := inter.check(ctx)
		if err != nil {
			return nil, err
		}

		newCtx := metadata.NewIncomingContext(ctx, md)

		return handler(newCtx, req)
	}
}

func (inter *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := stream.Context()

		md, err := inter.check(ctx)
		if err != nil {
			return err
		}

		newStream := &wrapServerStream{
			ServerStream: stream,
			ctx:          metadata.NewIncomingContext(ctx, md),
		}

		return handler(srv, newStream)
	}
}

func (inter *AuthInterceptor) check(ctx context.Context) (metadata.MD, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "отсутствуют метаданные к запросе")
	}

	vals := md["token"]
	if len(vals) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "токен отстутствует")
	}

	token := vals[0]

	claims, err := inter.jwt.Verify(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "токен не валидный: %v", err)
	}

	md.Delete("token")
	md.Set(model.UserIDCtxKey, claims.UserID)

	return md, nil
}

func getMethod(method string) string {
	arr := strings.Split(method, string(os.PathSeparator))

	return arr[len(arr)-1]
}
