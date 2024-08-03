// отправка данных на клиент
package octoserver

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
)

var errSaveForUser = errors.New("нельзя сохранить файл для текущего юзера")

// GetArray Отправляет все сохраненные данные на клиент.
func (ts *OctoServer) GetArray(ctx context.Context, req *pb.GetArrayRequest) (*pb.GetArrayResponse, error) {
	resp := &pb.GetArrayResponse{}

	// получаем неоходимые id
	ids, err := getIDs(ctx, model.UserIDCtxKey)
	if err != nil {
		return nil, fmt.Errorf("read metadata: %w", err)
	}

	userID := ids[0]

	arr, err := ts.store.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("storeList: %w", err)
	}

	resp.Arr = buildProtoArr(arr)

	return resp, nil
}

// ServerStreamFile Отправляет сохраненные файлы юзера на клиент.
func (ts *OctoServer) ServerStreamFile(req *pb.FileStreamServerRequest, stream pb.Octo_ServerStreamFileServer) error {
	fmt.Println("start server stream out")
	defer fmt.Println("stop server stream out")

	ctx := stream.Context()

	ids, err := getIDs(ctx, model.UserIDCtxKey)
	if err != nil {
		return fmt.Errorf("get md: %w", err)
	}

	fmt.Printf("SendStream: IDs: %v\n", ids)

	userID := ids[0]
	fileID := req.GetIdFile()

	if !ts.checkData(ctx, userID, fileID) {
		return errSaveForUser
	}

	file, err := ts.fileStore.OpenReadeFile(fileID)
	if err != nil {
		return fmt.Errorf("fileStore open file: %w", err)
	}

	if err := sendFile(stream, file, ts.sendBufSize); err != nil {
		return fmt.Errorf("sendFile: %w", err)
	}

	return nil
}

// sendFile Читает файл и отправляет в стим.
func sendFile(stream pb.Octo_ServerStreamFileServer, file io.Reader, bufSize int) error {
	buf0 := make([]byte, bufSize)

	for {
		n, err := file.Read(buf0)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("file read: %w", err)
		}

		if err := stream.Send(&pb.FileStreamServerResponse{Data: buf0[:n]}); err != nil {
			return fmt.Errorf("send file: %w", err)
		}
	}

	return nil
}

// buildArray Собирает срез model.Data.
func buildArray(protoArr []*pb.Data) []*model.Data {
	arr := make([]*model.Data, len(protoArr))

	for i := range protoArr {
		arr[i] = model.NewData(
			protoArr[i].GetId(),
			protoArr[i].GetName(),
			model.TypeData(protoArr[i].GetRtype()),
			protoArr[i].GetMeta(),
			protoArr[i].GetData(),
			true,
		)
	}

	return arr
}

// buildProtoArr Собирает срез proto.Data.
func buildProtoArr(arr []*model.Data) []*pb.Data {
	protoArr := make([]*pb.Data, len(arr))

	for i := range arr {
		protoArr[i] = &pb.Data{
			Id:    arr[i].ID(),
			Name:  arr[i].Name(),
			Rtype: pb.DataType(arr[i].Type()),
			Meta:  arr[i].Meta(),
			Data:  arr[i].Data(),
		}
	}

	return protoArr
}
