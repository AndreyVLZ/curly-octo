package auth

import (
	"context"
	"log"
	"net"
	"testing"

	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type mockAuthServer struct {
	pb.UnimplementedAuthServer
}

func (ms *mockAuthServer) Registration(ctx context.Context, req *pb.RegRequest) (*pb.RegResponce, error) {
	if req.GetName() == "" || req.GetPassword() == "" {
		return nil, nil
	}

	resp := pb.RegResponce{Token: req.GetName() + req.GetPassword()}

	return &resp, nil
}

func (ms *mockAuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponce, error) {
	if req.GetName() == "" || req.GetPassword() == "" {
		return nil, nil
	}

	resp := pb.LoginResponce{Token: req.GetName() + req.GetPassword()}

	return &resp, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterAuthServer(server, &mockAuthServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestAuth(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	authClient := NewAuthClient(conn)

	if err := authClient.Registration(ctx, "LOGIN", "PASS"); err != nil {
		t.Errorf("reg: %v\n", err)

		return
	}

	if authClient.Token() == "" {
		t.Error("token empty")
	}

	if err := authClient.Login(ctx, "LOGIN", "PASS"); err != nil {
		t.Errorf("login: %v\n", err)

		return
	}

	if authClient.Token() == "" {
		t.Error("token empty")

		return
	}

	if err := authClient.Refresh(ctx); err != nil {
		t.Errorf("refresh: %v\n", err)

		return
	}
}
