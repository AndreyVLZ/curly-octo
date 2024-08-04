package octoserver

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"google.golang.org/grpc/metadata"
)

type fileStorager interface {
	OpenReadeFile(filePath string) (io.ReadCloser, error)
	OpenWiteFile(filePath string) (io.WriteCloser, error)
}

type storager interface {
	GetData(ctx context.Context, userID string, dataID string) (model.Data, error)
	SaveArray(ctx context.Context, userID string, arr []*model.Data) error
	List(ctx context.Context, userID string) ([]*model.Data, error) // sync
}

// OctoServer Отвечает за приём/отправку данных и стрим файлов.
type OctoServer struct {
	pb.UnimplementedOctoServer
	store       storager
	fileStore   fileStorager
	sendBufSize int
}

func NewOctoServer(store storager, fileStore fileStorager, sendBufSize int) *OctoServer {
	return &OctoServer{
		fileStore:   fileStore,
		store:       store,
		sendBufSize: sendBufSize,
	}
}

// checkData Проверяет, что юзер сохраняет свой файл.
func (ts *OctoServer) checkData(ctx context.Context, userID, dataID string) bool {
	_, err := ts.store.GetData(ctx, userID, dataID)

	return err == nil
}

// getIDs Получает из контекста значения по ключам.
func getIDs(ctx context.Context, keys ...string) ([]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("not matadata")
	}

	fmt.Printf("SERVER MD: %v\n", md)

	arrRes := make([]string, len(keys))

	for i := range keys {
		vals := md.Get(keys[i])
		if len(vals) == 0 {
			return nil, errors.New("incorrect ctxKey")
		}

		arrRes[i] = vals[0]
	}

	return arrRes, nil
}
