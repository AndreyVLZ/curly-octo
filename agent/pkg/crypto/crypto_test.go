package crypto

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randBytes(n int) ([]byte, error) {
	arr := make([]byte, n)

	if _, err := io.ReadFull(rand.Reader, arr); err != nil {
		return nil, fmt.Errorf("readFull: %w", err)
	}

	return arr, nil
}

type writeCloser struct {
	*bytes.Buffer
}

func (wc *writeCloser) Close() error { return nil }

func TestReadWrite(t *testing.T) {
	secretKey := []byte("SECRET-KEY")
	sizeChunk := 1000

	myCrypto, err := NewCrypto(secretKey, sizeChunk)
	if err != nil {
		t.Errorf("new crypto: %v\n", err)

		return
	}

	bufRead, err := randBytes(100000)
	if err != nil {
		t.Errorf("rand bytes: %v\n", err)

		return
	}
	reader := bytes.NewReader(bufRead)

	cryptoRead := myCrypto.NewStreamEncReader(io.NopCloser(reader))

	bufEnc, err := io.ReadAll(cryptoRead)
	if err != nil {
		t.Errorf("readAll: %v\n", err)

		return
	}

	var bufWrite bytes.Buffer

	decWrite := myCrypto.NewStreamDecWriter(&writeCloser{Buffer: &bufWrite})

	if _, err := decWrite.Write(bufEnc); err != nil {
		t.Errorf("dec write: %v\n", err)

		return
	}

	assert.Equal(t, len(bufRead), len(bufWrite.Bytes()))
	assert.Equal(t, bufRead, bufWrite.Bytes())
}

func TestEncDec(t *testing.T) {
	secretKey := []byte("SECRET-KEY")
	sizeChunk := 1000

	myCrypto, err := NewCrypto(secretKey, sizeChunk)
	if err != nil {
		t.Errorf("new crypto: %v\n", err)

		return
	}

	msg := []byte("MESSAGE")

	encMsg, err := myCrypto.Encode(msg)
	if err != nil {
		t.Errorf("encMsg: %v\n", err)

		return
	}

	resMsg, err := myCrypto.Decode(encMsg)
	if err != nil {
		t.Errorf("encMsg: %v\n", err)

		return
	}

	assert.Equal(t, msg, resMsg)
}
