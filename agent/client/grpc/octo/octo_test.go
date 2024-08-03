package octo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type writeCloser struct {
	*bytes.Buffer
}

func (wc *writeCloser) Close() error { return nil }

type mockOctoServer struct {
	pb.UnimplementedOctoServer
	files  int
	mDatas []*model.Data
}

func (os *mockOctoServer) SendArray(ctx context.Context, req *pb.SendArrayRequest) (*pb.SendArrayResponse, error) {
	n := len(req.GetArr())
	fmt.Printf("len: %d\n", n)

	return &pb.SendArrayResponse{}, nil
}

func (os *mockOctoServer) StreamFile(stream pb.Octo_StreamFileServer) error {
	os.files += 1
	stream.SendAndClose(&pb.FileResponse{})

	return nil
}

func (os *mockOctoServer) GetArray(ctx context.Context, req *pb.GetArrayRequest) (*pb.GetArrayResponse, error) {
	resp := &pb.GetArrayResponse{
		Arr: buildProtoArr(os.mDatas),
	}

	return resp, nil
}

func (os *mockOctoServer) ServerStreamFile(req *pb.FileStreamServerRequest, stream pb.Octo_ServerStreamFileServer) error {
	return nil
}

func dialer(ser pb.OctoServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterOctoServer(server, ser)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestOcto(t *testing.T) {
	ctx := context.Background()

	mockServer := &mockOctoServer{}

	//conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	conn, err := grpc.NewClient(":3300",
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(mockServer)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	octoServer := New(conn, 100)

	logPass := model.NewData(
		"ID-1",
		"LogPass-1",
		model.LogPassData,
		[]byte("meta"),
		[]byte("user data"),
		true,
	)

	arr := []*model.Data{logPass}

	if err := octoServer.SendData(ctx, arr); err != nil {
		t.Errorf("send Data: %v\n", err)

		return
	}

	myFile := []*model.EncFile{
		model.NewEncFile("F-1", io.NopCloser(strings.NewReader("МОЙ-ФАЙЛ-1"))),
		model.NewEncFile("F-2", io.NopCloser(strings.NewReader("МОЙ-ФАЙЛ-2"))),
	}

	if err := octoServer.SendFiles(ctx, myFile); err != nil {
		t.Errorf("send files: %v\n", err)

		return
	}

	if len(myFile) != mockServer.files {
		t.Error("no eq files")

		return
	}

	mockServer.mDatas = arr

	resData, err := octoServer.RecvData(ctx)
	if err != nil {
		t.Errorf("recv files: %v\n", err)

		return
	}

	assert.Equal(t, arr, resData)
}
