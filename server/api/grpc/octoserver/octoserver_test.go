package octoserver

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

type mockFileStore struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (mfs *mockFileStore) OpenReadeFile(filePath string) (io.ReadCloser, error) {
	return mfs.r, nil
}
func (mfs *mockFileStore) OpenWiteFile(filePath string) (io.WriteCloser, error) {
	return mfs.w, nil
}

type mockStore struct {
	list []*model.Data
}

func (ms *mockStore) GetData(ctx context.Context, userID string, dataID string) (model.Data, error) {
	return model.Data{}, nil
}

func (ms *mockStore) SaveArray(ctx context.Context, userID string, arr []*model.Data) error {
	return nil
}

func (ms *mockStore) List(ctx context.Context, userID string) ([]*model.Data, error) {
	return ms.list, nil
}

func dialer(srv pb.OctoServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterOctoServer(server, srv)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestOctoServer(t *testing.T) {
	ctx := context.Background()

	logPass := model.NewData(
		"ID-1",
		"LogPass-1",
		model.LogPassData,
		[]byte("meta"),
		[]byte("user data"),
		true,
	)

	arr := []*model.Data{logPass}

	octo := NewOctoServer(
		&mockStore{list: arr},
		&mockFileStore{},
		1000,
	)

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(octo)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewOctoClient(conn)

	md := metadata.New(map[string]string{
		model.UserIDCtxKey: "USER-ID",
	})

	nCtx := metadata.NewOutgoingContext(ctx, md)

	resp, err := client.GetArray(nCtx, &pb.GetArrayRequest{})
	if err != nil {
		t.Errorf("getArray: %v\n", err)

		return
	}

	t.Logf("respArr: %v\n", resp.GetArr())

	arrRes := buildArray(resp.GetArr())
	assert.Equal(t, arr, arrRes)
}
