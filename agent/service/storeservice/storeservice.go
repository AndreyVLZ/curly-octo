package storeservice

import (
	"context"
	"fmt"
	"io"

	"github.com/AndreyVLZ/curly-octo/internal/model"
)

type storager interface {
	List(ctx context.Context) ([]*model.Data, error)
	SaveArray(ctx context.Context, arr []*model.Data) error
}

type fileStorager interface {
	SaveFiles(arr []*model.Data) ([]*model.File, error)
	GetFiles(arr []*model.Data) ([]*model.File, error)
}

type crypter interface {
	NewStreamDecWriter(wc io.WriteCloser) io.WriteCloser
	NewStreamEncReader(rc io.ReadCloser) io.ReadCloser
	Encode(msg []byte) ([]byte, error)
	Decode(encMsg []byte) ([]byte, error)
}

// StoreService Отвечает за шифрование исходящих и расшифрование входящих данных.
type StoreService struct {
	crypto    crypter
	store     storager
	fileStore fileStorager
}

func NewStoreService(crypter crypter, store storager, fileStore fileStorager) *StoreService {
	return &StoreService{
		crypto:    crypter,
		store:     store,
		fileStore: fileStore,
	}
}

// GetAll ...
func (srv *StoreService) GetAll(ctx context.Context) ([]*model.Data, []*model.EncFile, error) {
	// получаем данные из хранилища
	array, err := srv.store.List(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("store list: %w", err)
	}

	// получаем файлы
	files, err := srv.fileStore.GetFiles(array)
	if err != nil {
		return nil, nil, fmt.Errorf("filestore GetFiles: %w", err)
	}

	// зашифровываем данные
	if err := encryptArray(srv.crypto.Encode, array); err != nil {
		return nil, nil, err
	}

	// шифратор для стрима
	encFiles := encryptFiles(srv.crypto, files)

	return array, encFiles, nil
}

// SaveArray ...
func (srv *StoreService) SaveArray(ctx context.Context, arr []*model.Data) ([]*model.DecFile, error) {
	fmt.Printf("ss данные получены[%d]\n", len(arr))

	// расшифровываем полученные данные
	if err := decryptArray(srv.crypto.Decode, arr); err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	// сохранаем данные
	if err := srv.store.SaveArray(ctx, arr); err != nil {
		return nil, fmt.Errorf("store saveArr: %w", err)
	}

	// сохраняем файлы
	files, err := srv.fileStore.SaveFiles(arr)
	if err != nil {
		return nil, fmt.Errorf("fileStore SaveFiles: %w", err)
	}

	// дешифратор для стрима
	decFiles := decryptFiles(srv.crypto, files)

	return decFiles, nil
}

// encryptArray Зашифровывает данные.
func encryptArray(fnEnc func([]byte) ([]byte, error), arr []*model.Data) error {
	for i := range arr {
		if err := arr[i].Encrypt(fnEnc); err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
	}

	return nil
}

// decryptArray Расшифровывает данные.
func decryptArray(fnDec func([]byte) ([]byte, error), arr []*model.Data) error {
	for i := range arr {
		if err := arr[i].Decrypt(fnDec); err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	}

	return nil
}

// encryptFiles ...
func encryptFiles(crypt crypter, files []*model.File) []*model.EncFile {
	encFiles := make([]*model.EncFile, len(files))

	for i := range files {
		encFiles[i] = model.NewEncFile(
			files[i].ID(),
			crypt.NewStreamEncReader(files[i]),
		)
	}

	return encFiles
}

// decryptFiles ...
func decryptFiles(crypt crypter, files []*model.File) []*model.DecFile {
	decFiles := make([]*model.DecFile, len(files))

	for i := range files {
		decFiles[i] = model.NewDecFile(
			files[i].ID(),
			crypt.NewStreamDecWriter(files[i]),
		)
	}

	return decFiles
}
