package filestore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/AndreyVLZ/curly-octo/internal/model"
)

type FileStore struct {
	tmpDir string
}

func NewFileStore(tmpDir string) *FileStore {
	return &FileStore{
		tmpDir: tmpDir,
	}
}

func (fs *FileStore) OpenReadeFile(filePath string) (io.ReadCloser, error) {
	osFile, err := os.OpenFile(filepath.Join(fs.tmpDir, filePath), os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("new file: %w", err)
	}

	return osFile, nil
}

func (fs *FileStore) OpenWiteFile(fileName string) (io.WriteCloser, error) {
	filePath := filepath.Join(fs.tmpDir, fileName)
	if isExistsFile(filePath) {
		return nil, fmt.Errorf("файла [%s] уже существует", filePath)
	}

	osFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("new file: %w", err)
	}

	return osFile, nil
}

// GetFiles Получает.
func (fs *FileStore) GetFiles(arr []*model.Data) ([]*model.File, error) {
	files := make([]*model.File, 0, len(arr))

	for i := range arr {
		if arr[i].Type() != model.BinaryData {
			continue
		}

		file, err := os.OpenFile(string(arr[i].Data()), os.O_RDONLY, 0666)
		if err != nil {
			return nil, fmt.Errorf("new model file: %w", err)
		}

		files = append(files, model.NewFile(arr[i].ID(), file))
	}

	return files, nil
}

func (fs *FileStore) SaveFiles(arr []*model.Data) ([]*model.File, error) {
	files := make([]*model.File, 0, len(arr))

	for i := range arr {
		if arr[i].Type() != model.BinaryData {
			continue
		}

		filePath := filepath.Join(fs.tmpDir, string(arr[i].Data()))
		if isExistsFile(filePath) {
			return nil, fmt.Errorf("файла [%s] уже существует", filePath)
		}

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, fmt.Errorf("new model file: %w", err)
		}

		files = append(files, model.NewFile(arr[i].ID(), file))
	}

	return files, nil
}

func isExistsFile(filePath string) bool {
	_, err := os.Stat(filePath)

	return err == nil
}
