package octo

import (
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"google.golang.org/grpc"
)

// Client Отвечает за отправку/приём данных и стрима файлов.
type Client struct {
	octoClient  pb.OctoClient
	sendBufSize int
}

func New(conn *grpc.ClientConn, sendBufSize int) *Client {
	return &Client{
		octoClient:  pb.NewOctoClient(conn),
		sendBufSize: sendBufSize,
	}
}
