package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/AndreyVLZ/curly-octo/agent/api/cli"
	"github.com/AndreyVLZ/curly-octo/agent/client/grpc"
	"github.com/AndreyVLZ/curly-octo/agent/pkg/crypto"
	"github.com/AndreyVLZ/curly-octo/agent/pkg/store/localstore"
	"github.com/AndreyVLZ/curly-octo/agent/service/storeservice"
	"github.com/AndreyVLZ/curly-octo/internal/store/filestore"
	"github.com/AndreyVLZ/curly-octo/internal/store/inmemory"
)

const (
	caCertFile     = "cert/ca-cert.pem"
	clientCertFile = "cert/client-cert.pem"
	clientKeyFile  = "cert/client-key.pem"
)

const (
	AddresDefault       = ":3200"
	TmpDirDefault       = "tmpDir/agent"
	SendBufSizeDefault  = 30000
	ChunkBufSizeDefault = 10000
)

type iAPI interface {
	Start(ctx context.Context) error
}

type Agent struct {
	api iAPI
	// store storager // локальное хранение данных
}

type Config struct {
	Addr            string // адрес подключения
	TmpDir          string // куда сохранять принятые файлы
	SecretKey       string // ключ
	CryptoChunkSize int    // размер буфера для шифрования/расшифрования
	SendBufSize     int    // размер буфера чтения/записи файла
}

func NewConfig(addr, tmpDir, secretKey string, cryptoChunkSize, sendBufSize int) Config {
	return Config{
		Addr:            addr,
		TmpDir:          tmpDir,
		SecretKey:       secretKey,
		CryptoChunkSize: cryptoChunkSize,
		SendBufSize:     sendBufSize,
	}
}

func ConfigDefault() Config {
	return NewConfig(AddresDefault, TmpDirDefault, "", ChunkBufSizeDefault, SendBufSizeDefault)
}

func New(cfg Config) (*Agent, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// AEAD-GCM
	myCrypto, err := crypto.NewCrypto([]byte(cfg.SecretKey), cfg.CryptoChunkSize)
	if err != nil {
		return nil, fmt.Errorf("new myCrypto: %w", err)
	}

	store := inmemory.New()

	// используем только локальное хранилище на клиенте
	localStore := localstore.NewLocalStore("-LOCAL-USER-", store)

	// для работы с файлами
	fStore := filestore.NewFileStore(cfg.TmpDir)

	// для шифрования/расшифрования
	storeSrv := storeservice.NewStoreService(myCrypto, localStore, fStore)

	tlsConf, err := initTLS()
	if err != nil {
		return nil, fmt.Errorf("init tls: %w", err)
	}

	// основной клиент
	client, err := grpc.NewClient(cfg.Addr, tlsConf, storeSrv, cfg.SendBufSize)
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	// для работы из консоли
	api := cli.New(client, localStore)

	return &Agent{
		api: api,
	}, nil
}

func (a *Agent) Stop(ctx context.Context) error {
	// a.api.Stop() // пока нет

	return nil
}

func (a *Agent) Start(ctx context.Context) error { return a.api.Start(ctx) }

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

	clientTLSCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, fmt.Errorf("load x509 cert: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{clientTLSCert},
		RootCAs:      certPool,
	}, nil
}
