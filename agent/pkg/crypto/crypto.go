package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

// nonceSizeStandart Стандартный размер блока nonce.
// https://cs.opensource.google/go/go/+/master:src/crypto/cipher/gcm.go;l=157
const (
	nonceSizeStandart = 12
	blockSize         = 16
)

func createHash256(key []byte) ([]byte, error) {
	hasher := sha256.New()
	if _, err := hasher.Write(key); err != nil {
		return nil, fmt.Errorf("hasher write: %w", err)
	}

	return hasher.Sum(nil), nil
}

func genNonce() ([]byte, error) {
	randomBytes := make([]byte, nonceSizeStandart)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("crypto rand read: %w", err)
	}

	return randomBytes, nil
}

type Crypto struct {
	aead      cipher.AEAD
	sizeChunk int
}

func NewCrypto(key []byte, sizeChunk int) (*Crypto, error) {
	hashKey, err := createHash256(key)
	if err != nil {
		return nil, fmt.Errorf("new hash: %w", err)
	}

	cBlock, err := aes.NewCipher(hashKey)
	if err != nil {
		return nil, fmt.Errorf("newCipher: %w", err)
	}

	cipherAEAD, err := cipher.NewGCM(cBlock)
	if err != nil {
		return nil, fmt.Errorf("newGCM: %w", err)
	}

	return &Crypto{
		aead:      cipherAEAD,
		sizeChunk: sizeChunk,
	}, nil
}

func (c *Crypto) NewStreamEncReader(rc io.ReadCloser) io.ReadCloser {
	return &StreaEncReader{
		aead:      c.aead,
		rc:        rc,
		buf:       make([]byte, 0, c.sizeChunk+blockSize+nonceSizeStandart-1),
		sizeChunk: c.sizeChunk,
	}
}

// StreaEncReader Попытка реализации потокового шифратора для AEAD-GCM.
// Зашифровывает байты прочитанные из файла [rc].
// Байты накапливаются во внутренний буфер [buf].
// Когда размер буфер подходит для шифрования,
// байты из буфера зашифровываются.
type StreaEncReader struct {
	aead      cipher.AEAD
	rc        io.ReadCloser
	buf       []byte
	sizeChunk int
}

func (se *StreaEncReader) Close() error { return se.rc.Close() }

// Read Копирует байты из буфера в [dst].
// Если буфер пуст - шифрует новую порцию байт.
func (se *StreaEncReader) Read(dst []byte) (int, error) {
	for len(se.buf) > 0 {
		n := copy(dst, se.buf)
		se.buf = se.buf[n:]

		return n, nil
	}

	nn, err := se.read()
	if err != nil {
		return 0, err
	}

	n := copy(dst, se.buf[:nn])
	se.buf = se.buf[n:]

	return n, nil
}

func (se *StreaEncReader) read() (int, error) {
	chunk := make([]byte, se.sizeChunk)

	nonce, err := genNonce()
	if err != nil {
		return 0, fmt.Errorf("gen nonce: %w", err)
	}

	nn, err := io.ReadFull(se.rc, chunk)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		if !errors.Is(err, io.EOF) {
			return 0, fmt.Errorf("readAll from rc: %w", err)
		}

		return 0, err
	}

	encData := se.aead.Seal(nonce /*dst*/, nonce /*nonce*/, chunk[:nn] /*plainText*/, nil)
	se.buf = append(se.buf, encData...)

	return nn, nil
}

func (c *Crypto) NewStreamDecWriter(wc io.WriteCloser) io.WriteCloser {
	return &StreaDecWriter{
		aead:      c.aead,
		wc:        wc,
		buf:       make([]byte, 0, c.sizeChunk+blockSize+nonceSizeStandart-1),
		sizeChunk: c.sizeChunk,
	}
}

// StreaDecWriter Попытка реализации потокового дешифратора для AEAD-GCM.
// Расшифровывает входящие байты.
// Входящие байты накапливаются во внутренний буфер [buf].
// Когда размер буфер подходит для расшифрования,
// байты из буфера расшифровываются.
// Результат сбрасывается в файл [wc].
// Обязательно вызов Close по завершению.
type StreaDecWriter struct {
	aead      cipher.AEAD
	wc        io.WriteCloser
	buf       []byte
	sizeChunk int
}

// Close Расшифровывает оставшиеся байты в буфере.
// Закрывает файл.
func (sdw *StreaDecWriter) Close() error {
	errs := make([]error, 0, 2)

	if len(sdw.buf) > 0 {
		if _, err := sdw.write(len(sdw.buf)); err != nil {
			errs = append(errs, err)
		}
	}

	sdw.buf = sdw.buf[:0]

	if err := sdw.wc.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// Write Расшифровывает входящие байты.
// Может вернуть 0, и нулевую ошибку -
// это значит что байты сохранены в буфер, но еще не расшифрованы.
func (sdw *StreaDecWriter) Write(src []byte) (int, error) {
	var n int

	sdw.buf = append(sdw.buf, src...)
	t := sdw.sizeChunk + blockSize + nonceSizeStandart

	for len(sdw.buf) >= t {
		nn, err := sdw.write(t)
		if err != nil {
			return 0, err
		}

		n += nn
	}

	return n, nil
}

func (sdw *StreaDecWriter) write(t int) (int, error) {
	res, err := sdw.aead.Open(nil /*dst*/, sdw.buf[:nonceSizeStandart] /*nonce*/, sdw.buf[nonceSizeStandart:t] /*cipherText*/, nil)
	if err != nil {
		return 0, fmt.Errorf("decrypt msg: %w", err)
	}

	n, err := sdw.wc.Write(res)
	if err != nil {
		return 0, fmt.Errorf("write to wc: %w", err)
	}

	sdw.buf = sdw.buf[t:]

	return n, nil
}

func (c *Crypto) Encode(msg []byte) ([]byte, error) {
	nonce, err := genNonce()
	if err != nil {
		return nil, fmt.Errorf("gen nonce: %w", err)
	}

	res := c.aead.Seal(nonce /*dst*/, nonce /*nonce*/, msg /*plain*/, nil)

	return res, nil
}

func (c *Crypto) Decode(encMsg []byte) ([]byte, error) {
	res, err := c.aead.Open(nil /*dst*/, encMsg[:nonceSizeStandart] /*nonce*/, encMsg[nonceSizeStandart:] /*cipherText*/, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt msg: %w", err)
	}

	return res, nil
}
