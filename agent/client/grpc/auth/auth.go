package auth

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"google.golang.org/grpc"
)

// AuthClient Отвечает за авторизацию.
type AuthClient struct {
	client pb.AuthClient
	uLogin string
	uPass  string
	token  string
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		client: pb.NewAuthClient(conn),
		uLogin: "",
		uPass:  "",
		token:  "",
	}
}

// Token Возвращает сохранённый токен.
func (ac *AuthClient) Token() string { return ac.token }

// Refresh Оюновление токена.
func (ac *AuthClient) Refresh(ctx context.Context) error {
	if ac.uLogin == "" || ac.uPass == "" {
		return errors.New("authClient: user is empty")
	}

	if err := ac.Login(ctx, ac.uLogin, ac.uPass); err != nil {
		return fmt.Errorf("authClient login: %w", err)
	}

	return nil
}

// Registration Отправляет запрос на регистрацию.
func (ac *AuthClient) Registration(ctx context.Context, login, pass string) error {
	ac.uLogin, ac.uPass = login, pass

	req := pb.RegRequest{Name: login, Password: pass}

	resp, err := ac.client.Registration(ctx, &req)
	if err != nil {
		return fmt.Errorf("client Registration: %w", err)
	}

	ac.token = resp.GetToken()

	return nil
}

// Login Отправляет запрос на вход.
func (ac *AuthClient) Login(ctx context.Context, login, pass string) error {
	ac.uLogin, ac.uPass = login, pass
	req := pb.LoginRequest{Name: login, Password: pass}

	resp, err := ac.client.Login(ctx, &req)
	if err != nil {
		return fmt.Errorf("client Login: %w", err)
	}

	ac.token = resp.GetToken()

	return nil
}
