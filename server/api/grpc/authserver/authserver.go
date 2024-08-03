package authserver

import (
	"context"
	"fmt"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	pb "github.com/AndreyVLZ/curly-octo/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type iUserStore interface {
	FindByName(name string) (*model.User, error)
	Add(ctx context.Context, user model.User) error
}

type jwtGenerator interface {
	Generate(user *model.User) (string, error)
}

// AuthService Отвечает за регистрацию и аутентификацию юзера.
type AuthService struct {
	pb.UnimplementedAuthServer
	store iUserStore
	jwt   jwtGenerator
}

func NewAuthService(store iUserStore, jwt jwtGenerator) pb.AuthServer {
	return &AuthService{
		store: store,
		jwt:   jwt,
	}
}

// Registration Регистрация юзера.
func (s *AuthService) Registration(ctx context.Context, req *pb.RegRequest) (*pb.RegResponce, error) {
	uName := req.GetName()
	uPass := req.GetPassword()

	user, err := model.NewUser(uName, uPass)
	if err != nil {
		return nil, fmt.Errorf("model newUser: %w", err)
	}

	if err := s.store.Add(ctx, *user); err != nil {
		return nil, fmt.Errorf("store add: %w", err)
	}

	token, err := s.jwt.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка создания ключа")
	}

	return &pb.RegResponce{Token: token}, nil
}

// Login Вход для юзера.
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponce, error) {
	user, err := s.store.FindByName(req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "юзер не найден: %v", err)
	}

	if user == nil || !user.CheckPass(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "неверный пароль")
	}

	token, err := s.jwt.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка создания ключа")
	}

	return &pb.LoginResponce{Token: token}, nil
}
