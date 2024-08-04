// приём данных от клиента.
package octoserver

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
)

// SendArray Сохраняет все данные юзера.
func (ts *OctoServer) SendArray(ctx context.Context, req *pb.SendArrayRequest) (*pb.SendArrayResponse, error) {
	resp := pb.SendArrayResponse{}

	// получаем неоходимые id
	ids, err := getIDs(ctx, model.UserIDCtxKey)
	if err != nil {
		resp.Error = "read metadata: " + err.Error() // не критичная ошибка

		return &resp, err
	}

	userID := ids[0]

	protoArr := req.GetArr()
	if err := ts.store.SaveArray(ctx, userID, buildArray(protoArr)); err != nil {
		return nil, fmt.Errorf("save arr: %w", err)
	}

	return &resp, nil
}

// recv
// StreamFile Записывает в файлы данные из стрима.
func (ts *OctoServer) StreamFile(stream pb.Octo_StreamFileServer) error {
	var resp pb.FileResponse

	defer stream.SendAndClose(&resp)

	ctx := stream.Context()
	// получаем неоходимые id
	ids, err := getIDs(ctx, model.UserIDCtxKey, model.FileIDCtxKey)
	if err != nil {
		resp.Error = "read metadata: " + err.Error() // не критичная ошибка

		return nil
	}

	userID, fileID := ids[0], ids[1]
	// проверяем, что данные которые мы хотим сохранить принадлежат конкретному юзеру
	if !ts.checkData(ctx, userID, fileID) {
		resp.Error = "incorrect userID or fileID"

		return nil
	}

	// открываем файл
	file, err := ts.fileStore.OpenWiteFile(fileID)
	if err != nil {
		return fmt.Errorf("fileStore open file: %w", err)
	}

	// читаем данные из стрима
	if err := fileRecv(stream, file); err != nil {
		return fmt.Errorf("fileRecv: %w", err)
	}

	return nil
}

// fileRecv Записывает в файл данные из стрима.
func fileRecv(stream pb.Octo_StreamFileServer, file io.Writer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("streamFile read: %w\n", err)
		}

		n, err := file.Write(req.GetData())
		if err != nil {
			return fmt.Errorf("write: %v\n", err)
		}

		fmt.Printf("steam len: %d\n", n)
	}

	return nil
}
