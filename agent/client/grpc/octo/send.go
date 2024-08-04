package octo

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"google.golang.org/grpc/metadata"
)

// SendData Оправляет данные [arr] на сервер.
func (c *Client) SendData(ctx context.Context, arr []*model.Data) error {
	req := &pb.SendArrayRequest{
		Arr: buildProtoArr(arr),
	}

	resp, err := c.octoClient.SendArray(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc sendArr: %w", err)
	}

	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}

	return nil
}

// SendData Оправляет файлы [files] на сервер.
func (c *Client) SendFiles(ctx context.Context, files []*model.EncFile) error {
	for i := range files {
		if err := c.sendFile(ctx, files[i]); err != nil {
			return fmt.Errorf("send files: %w", err)
		}
	}

	return nil
}

// sendFile Читает файл частями и отправляет в стрим.
func (c *Client) sendFile(ctx context.Context, file *model.EncFile) error {
	// метаданные с id передаваемого файла
	metaData := metadata.New(
		map[string]string{
			model.FileIDCtxKey: file.ID(),
		},
	)

	gCtx := metadata.NewOutgoingContext(ctx, metaData)

	stream, err := c.octoClient.StreamFile(gCtx)
	if err != nil {
		return fmt.Errorf("grpc StreamFile: %w", err)
	}

	buf0 := make([]byte, c.sendBufSize)

	for {
		n, err := file.Read(buf0)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("file read: %w", err)
		}

		fmt.Printf("send [%s] len[%d]\n", file.ID(), n)

		if err := stream.Send(&pb.FileRequest{Data: buf0[:n]}); err != nil {
			fmt.Printf("strem send: %v\n", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close and recv: %w", err)
	}

	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}

	return nil
}

// buildProtoArr Собирает срез proto.Data из model.Data.
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
