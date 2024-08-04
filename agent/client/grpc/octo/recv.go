// приём данных.
package octo

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
)

// RecvData Приём данных.
func (c *Client) RecvData(ctx context.Context) ([]*model.Data, error) {
	req := pb.GetArrayRequest{}

	resp, err := c.octoClient.GetArray(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("octo getArr: %w", err)
	}

	if resp.GetError() != "" {
		return nil, fmt.Errorf("resp err: %w\n", err)
	}

	return buildArr(resp.GetArr()), nil
}

// RecvFiles Прием файлов.
func (c *Client) RecvFiles(ctx context.Context, files []*model.DecFile) error {
	for i := range files {
		if err := c.recvFile(ctx, files[i]); err != nil {
			return err
		}
	}

	return nil
}

// recvFile Записывает в файл данные из стрима.
func (c *Client) recvFile(ctx context.Context, file *model.DecFile) error {
	defer file.Close()

	req := pb.FileStreamServerRequest{
		IdFile: file.ID(),
	}

	stream, err := c.octoClient.ServerStreamFile(ctx, &req)
	if err != nil {
		return fmt.Errorf("init stream: %w", err)
	}

	defer stream.CloseSend()

	for {
		// получаем данные
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("stream recv: %w", err)
		}
		// записываем в файл
		n, err := file.Write(resp.GetData())
		if err != nil {
			return fmt.Errorf("file write: %w", err)
		}

		fmt.Printf("NN: %d\n", n)
	}

	return nil
}

// buildArr Собирает срез model.Data из proto.Data.
func buildArr(protoArr []*pb.Data) []*model.Data {
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
