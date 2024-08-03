package authserver

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type mockUserStore struct {
	user *model.User
}

func (ms *mockUserStore) FindByName(name string) (*model.User, error) {
	return ms.user, nil
}

func (ms *mockUserStore) Add(ctx context.Context, user model.User) error {
	return nil
}

type mockJWT struct {
	token string
}

func (mj *mockJWT) Generate(user *model.User) (string, error) {
	return mj.token, nil
}

func dialer(srv pb.AuthServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterAuthServer(server, srv)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestAuthServiceReg(t *testing.T) {
	ctx := context.Background()

	login, pass := "LOGIN", "PASS"
	token := "AUTH-TOKEN"

	user, err := model.NewUser(login, pass)
	if err != nil {
		t.Errorf("new user: %v\n", err)

		return
	}

	auth := NewAuthService(&mockUserStore{user: user}, &mockJWT{token: token})

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(auth)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewAuthClient(conn)

	req := &pb.RegRequest{Name: login, Password: pass}

	resp, err := client.Registration(ctx, req)
	if err != nil {
		t.Errorf("reg: %v\n", err)

		return
	}

	assert.Equal(t, token, resp.GetToken())

	respL, err := client.Login(ctx, &pb.LoginRequest{Name: login, Password: pass})
	if err != nil {
		t.Errorf("login: %v\n", err)

		return
	}

	assert.Equal(t, token, respL.GetToken())
}
