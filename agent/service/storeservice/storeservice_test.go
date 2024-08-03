package storeservice

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	"github.com/stretchr/testify/assert"
)

type writeCloser struct {
	*bytes.Buffer
}

func (wc *writeCloser) Close() error { return nil }

type crypto struct{}

func (c *crypto) NewStreamDecWriter(wc io.WriteCloser) io.WriteCloser { return wc }
func (c *crypto) NewStreamEncReader(rc io.ReadCloser) io.ReadCloser {
	return rc
}
func (c *crypto) Encode(msg []byte) ([]byte, error)    { return msg, nil }
func (c *crypto) Decode(encMsg []byte) ([]byte, error) { return encMsg, nil }

type store struct {
	mDatas []*model.Data
}

func (s *store) List(ctx context.Context) ([]*model.Data, error) {
	return s.mDatas, nil
}

func (s *store) SaveArray(ctx context.Context, arr []*model.Data) error {
	s.mDatas = arr
	return nil
}

type fileStore struct {
	mFiles   []*model.File
	mDataRes []*model.Data
}

func (fs *fileStore) SaveFiles(arr []*model.Data) ([]*model.File, error) {
	fs.mDataRes = arr
	return fs.mFiles, nil
}

func (fs *fileStore) GetFiles(arr []*model.Data) ([]*model.File, error) {
	return fs.mFiles, nil
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()

	crypto := &crypto{}
	store := &store{}

	fStore := &fileStore{}

	storeSrv := NewStoreService(crypto, store, fStore)

	encData, encFiles, err := storeSrv.GetAll(ctx)
	if err != nil {
		t.Errorf("getAll: %v\n", err)

		return
	}

	_ = encData
	_ = encFiles
}

func TestSaveArr(t *testing.T) {
	ctx := context.Background()

	arr := []*model.Data{
		model.NewData("D-1", "NameD1", model.LogPassData, []byte("meta-1"), []byte("data-1"), false),
		model.NewData("D-2", "NameD1", model.LogPassData, []byte("meta-2"), []byte("data-2"), false),
	}

	buf1 := bytes.NewBuffer([]byte("FILE-111"))
	f1 := model.NewFile("F-1", &writeCloser{Buffer: buf1})
	files := []*model.File{f1}

	crypto := &crypto{}
	store := &store{}
	store.mDatas = arr

	fStore := &fileStore{}
	fStore.mFiles = files

	storeSrv := NewStoreService(crypto, store, fStore)

	resEncDatas, resEncFiles, err := storeSrv.GetAll(ctx)
	if err != nil {
		t.Errorf("getAll: %v\n", err)

		return
	}

	assert.Equal(t, arr, resEncDatas)

	if !equalsFiles(files, resEncFiles) {
		t.Errorf("eq: %v\n", err)

		return
	}
}

func equalsFiles(files []*model.File, encFiles []*model.EncFile) bool {
	if len(files) != len(encFiles) {
		return false
	}

	for i := range files {
		if files[i].ID() != encFiles[i].ID() {
			return false
		}

		if _, err := io.ReadAll(files[i]); err != nil {
			return false
		}

		b2, err := io.ReadAll(encFiles[i])
		if err != nil {
			return false
		}
		// данные уже прочитаны значит OK
		if len(b2) != 0 {
			return false
		}
	}

	return true
}
