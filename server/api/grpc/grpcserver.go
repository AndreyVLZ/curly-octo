package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"github.com/AndreyVLZ/curly-octo/server/api/grpc/authserver"
	"github.com/AndreyVLZ/curly-octo/server/api/grpc/interceptor"
	"github.com/AndreyVLZ/curly-octo/server/api/grpc/octoserver"
	"github.com/AndreyVLZ/curly-octo/server/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type fileStorager interface {
	OpenReadeFile(filePath string) (io.ReadCloser, error)
	OpenWiteFile(filePath string) (io.WriteCloser, error)
}

type jwtGenerator interface {
	Generate(user *model.User) (string, error)
}

type storager interface {
	// main store
	GetData(ctx context.Context, userID string, dataID string) (model.Data, error)
	SaveArray(ctx context.Context, userID string, arr []*model.Data) error
	List(ctx context.Context, userID string) ([]*model.Data, error) // sync
	// user store
	FindByName(name string) (*model.User, error)
	Add(ctx context.Context, user model.User) error
}

// GRPCServer ...
type GRPCServer struct {
	addr       string
	server     *grpc.Server
	octoServer pb.OctoServer
	authServer pb.AuthServer
	store      storager
}

func NewGRPCServere(addr string, tlsConf *tls.Config, jwt *jwt.JWT, store storager, fileStore fileStorager, sendBufSize int) *GRPCServer {
	tlsCred := credentials.NewTLS(tlsConf)

	inter := interceptor.New(jwt)

	opts := []grpc.ServerOption{
		grpc.Creds(tlsCred),
		grpc.UnaryInterceptor(inter.Unary()),
		grpc.StreamInterceptor(inter.Stream()),
	}

	// для регистрации и входа
	authServer := authserver.NewAuthService(store, jwt)

	// для работы с данными
	octoServer := octoserver.NewOctoServer(store, fileStore, sendBufSize)

	return &GRPCServer{
		addr:       addr,
		server:     grpc.NewServer(opts...),
		authServer: authServer,
		octoServer: octoServer,
		store:      store,
	}
}

// Start Запуск сервера.
func (s *GRPCServer) Start() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	pb.RegisterOctoServer(s.server, s.octoServer)
	pb.RegisterAuthServer(s.server, s.authServer)

	fmt.Println("start server main")

	return s.server.Serve(listen)
}

// Stop Остановка сервера.
func (s *GRPCServer) Stop() {
	// s.store.Stop() пока нет
	s.server.GracefulStop()
}
