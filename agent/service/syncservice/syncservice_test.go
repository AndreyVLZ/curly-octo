package syncservice

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	"github.com/stretchr/testify/assert"
)

type client struct {
	dataSend []*model.Data
	fileSend []*model.EncFile
	fileDec  []*model.DecFile
}

func (c *client) SendData(ctx context.Context, arr []*model.Data) error {
	c.dataSend = arr

	return nil
}

func (c *client) SendFiles(ctx context.Context, files []*model.EncFile) error {
	c.fileSend = files

	return nil
}

func (c *client) RecvData(ctx context.Context) ([]*model.Data, error) {
	return c.dataSend, nil
}

func (c *client) RecvFiles(ctx context.Context, files []*model.DecFile) error {
	c.fileDec = files
	return nil
}

type storeSrv struct {
	myData    []*model.Data
	myFile    []*model.EncFile
	myDecFile []*model.DecFile
}

func (ss *storeSrv) GetAll(ctx context.Context) ([]*model.Data, []*model.EncFile, error) {
	return ss.myData, ss.myFile, nil
}

func (ss *storeSrv) SaveArray(ctx context.Context, arr []*model.Data) ([]*model.DecFile, error) {
	ss.myData = arr

	return ss.myDecFile, nil
}

type writeCloser struct {
	*bytes.Buffer
}

func (wc *writeCloser) Close() error { return nil }

func TestRecv(t *testing.T) {
	ctx := context.Background()
	cl := &client{}
	storeS := &storeSrv{}

	buf1 := bytes.NewBuffer([]byte("МОЙ-ФАЙЛ-1"))
	buf2 := bytes.NewBuffer([]byte("МОЙ-ФАЙЛ-2"))

	storeS.myDecFile = []*model.DecFile{
		model.NewDecFile("F-1", &writeCloser{Buffer: buf1}),
		model.NewDecFile("F-2", &writeCloser{Buffer: buf2}),
	}

	cl.dataSend = []*model.Data{
		model.NewData("D-1", "NameD1", model.LogPassData, []byte("meta-1"), []byte("data-1"), false),
		model.NewData("D-2", "NameD1", model.LogPassData, []byte("meta-2"), []byte("data-2"), false),
	}

	syncSrv := NewSyncService(cl, storeS)

	if err := syncSrv.Recv(ctx); err != nil {
		t.Errorf("syncService send: %v\n", err)
	}

	assert.Equal(t, storeS.myDecFile, cl.fileDec)
	assert.Equal(t, storeS.myData, cl.dataSend)
}

func TestSend(t *testing.T) {
	ctx := context.Background()
	cl := &client{}

	storeS := &storeSrv{
		myFile: []*model.EncFile{
			model.NewEncFile("F-1", io.NopCloser(strings.NewReader("МОЙ-ФАЙЛ-1"))),
			model.NewEncFile("F-2", io.NopCloser(strings.NewReader("МОЙ-ФАЙЛ-2"))),
		},
		myData: []*model.Data{
			model.NewData("D-1", "NameD1", model.LogPassData, []byte("meta-1"), []byte("data-1"), false),
			model.NewData("D-2", "NameD1", model.LogPassData, []byte("meta-2"), []byte("data-2"), false),
		},
	}

	syncSrv := NewSyncService(cl, storeS)

	if err := syncSrv.Send(ctx); err != nil {
		t.Errorf("syncService send: %v\n", err)
	}

	assert.Equal(t, storeS.myData, cl.dataSend)
	assert.Equal(t, storeS.myFile, cl.fileSend)
}
