package server

import (
	"time"

	"github.com/AndreyVLZ/curly-octo/internal/store/filestore"
	"github.com/AndreyVLZ/curly-octo/internal/store/inmemory"
	"github.com/AndreyVLZ/curly-octo/server/api/grpc"
	"github.com/AndreyVLZ/curly-octo/server/pkg/jwt"
)

const (
	AddresDefault      = ":3200"
	TmpDirDefault      = "tmpDir/server"
	ExpiresAtDefault   = time.Hour
	SendBufSizeDefault = 30000
)

type iAPI interface {
	Start() error
	Stop()
}

type Server struct {
	api iAPI
}

type Config struct {
	Addr        string        // адрес подключения
	TmpDir      string        // куда сохранять принятые файлы
	SecretKey   string        // ключ
	ExpireAt    time.Duration // время жизни токена
	SendBufSize int           // размер буфера чтения/записи файла
}

func NewConfig(addr, tmpDir, secretKey string, expiresAt time.Duration, sendBufSize int) Config {
	return Config{
		Addr:        addr,
		TmpDir:      tmpDir,
		SecretKey:   secretKey,
		ExpireAt:    expiresAt,
		SendBufSize: sendBufSize,
	}
}

func ConfigDefault() Config {
	return NewConfig(AddresDefault, TmpDirDefault, "", ExpiresAtDefault, SendBufSizeDefault)
}

func New(cfg Config) *Server {
	// для работы с ключами
	jwt := jwt.New(cfg.SecretKey, cfg.ExpireAt)

	// хранилище
	store := inmemory.New()

	// для работы с файлами
	fStore := filestore.NewFileStore(cfg.TmpDir)

	grpcServer := grpc.NewGRPCServere(cfg.Addr, jwt, store, fStore, cfg.SendBufSize)

	return &Server{api: grpcServer}
}

func (s *Server) Start() error { return s.api.Start() }
func (s *Server) Stop()        { s.api.Stop() }
