package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AndreyVLZ/curly-octo/internal/store/filestore"
	"github.com/AndreyVLZ/curly-octo/internal/store/inmemory"
	"github.com/AndreyVLZ/curly-octo/server/api/grpc"
	"github.com/AndreyVLZ/curly-octo/server/pkg/jwt"
)

const (
	caCertFile     = "cert/ca-cert.pem"
	serverCertFile = "cert/server-cert.pem"
	serverKeyFile  = "cert/server-key.pem"
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

func New(cfg Config) (*Server, error) {
	// для работы с ключами
	jwt := jwt.New(cfg.SecretKey, cfg.ExpireAt)

	// хранилище
	store := inmemory.New()

	// для работы с файлами
	fStore := filestore.NewFileStore(cfg.TmpDir)

	tlsConf, err := initTLS()
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewGRPCServere(cfg.Addr, tlsConf, jwt, store, fStore, cfg.SendBufSize)

	return &Server{api: grpcServer}, nil
}

func (s *Server) Start() error { return s.api.Start() }
func (s *Server) Stop()        { s.api.Stop() }

func initTLS() (*tls.Config, error) {
	// сертификат 'удостоверяющего центра'
	caPEM, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("open CA cert: %w", err)
	}

	certPool := x509.NewCertPool()

	// добавляем сертификат
	if !certPool.AppendCertsFromPEM(caPEM) {
		return nil, errors.New("add to cert poll")
	}

	serverTLSCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, fmt.Errorf("load x509 cert: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // сервер должен доверять сертификату клиента
		ClientCAs:    certPool,
	}, nil
}
