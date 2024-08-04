package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/AndreyVLZ/curly-octo/agent/client/grpc/auth"
	"github.com/AndreyVLZ/curly-octo/agent/client/grpc/interceptor"
	"github.com/AndreyVLZ/curly-octo/agent/client/grpc/octo"
	"github.com/AndreyVLZ/curly-octo/agent/service/syncservice"
	"github.com/AndreyVLZ/curly-octo/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type iAuthClient interface {
	Login(ctx context.Context, login, pass string) error        // вход
	Registration(ctx context.Context, login, pass string) error // регистрация
}

type iSyncService interface {
	Send(ctx context.Context) error // отправка сохрённеных данных
	Recv(ctx context.Context) error // получение данных с сервера
}

type iStoreService interface {
	GetAll(ctx context.Context) ([]*model.Data, []*model.EncFile, error)        // получение
	SaveArray(ctx context.Context, arr []*model.Data) ([]*model.DecFile, error) // сохранение
}

// Client GRPC-клиент.
type Client struct {
	authClient  iAuthClient  // регистрация/авторизация
	syncService iSyncService // отправка/получение шифрованных данных
	conn        *grpc.ClientConn
	authConn    *grpc.ClientConn
}

func NewClient(addr string, tlsConf *tls.Config, storeSrv iStoreService, sendBufSize int) (*Client, error) {
	tlsCrend := credentials.NewTLS(tlsConf)

	// опции для auth-клиента
	authOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCrend),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// соединение для auth-клиента
	authConn, err := grpc.NewClient(addr, authOpts...)
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	// auth-клиент
	authClient := auth.NewAuthClient(authConn)

	// middle
	inter := interceptor.NewAuthInterceptor(authClient)

	// опции для основного клиента
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCrend),
		grpc.WithUnaryInterceptor(inter.Unary()),
		grpc.WithStreamInterceptor(inter.Stream()),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	// основной клиент
	octoClient := octo.New(conn, sendBufSize)

	// сервис синхронизации данных
	syncSrv := syncservice.NewSyncService(octoClient, storeSrv)

	return &Client{
		syncService: syncSrv,
		authClient:  authClient,
		conn:        conn,
		authConn:    authConn,
	}, nil
}

// Close Закрывает все соединения клиента.
func (c *Client) Close(ctx context.Context) error {
	return errors.Join(c.authConn.Close(), c.conn.Close())
}

// Login Вход для юзера.
func (c *Client) Login(ctx context.Context, login, pass string) error {
	return c.authClient.Login(ctx, login, pass)
}

// Registration Регистрация юзера.
func (c *Client) Registration(ctx context.Context, login, pass string) error {
	return c.authClient.Registration(ctx, login, pass)
}

// Recv Получение данных.
func (c *Client) Recv(ctx context.Context) error {
	return c.syncService.Recv(ctx)
}

// Send Отправка данных.
func (c *Client) Send(ctx context.Context) error {
	return c.syncService.Send(ctx)
}
