package model

import (
	"errors"
	"fmt"
	"os"
)

const (
	UserIDCtxKey = "userID"
	FileIDCtxKey = "fileID"
)

var (
	errEncode = errors.New("данные уже расшифрованы")
	errDecode = errors.New("данные уже зашифрованны")
)

type TypeData uint8

var namesType = []string{
	"unknown",
	"logPass",
	"text",
	"binary",
	"card",
}

func (td TypeData) String() string { return namesType[td] }

const (
	UnknownData TypeData = iota
	LogPassData
	TextData
	BinaryData
	CardData
)

// Data Хранит данные от пользователя.
type Data struct {
	id          string   // id записи
	name        string   // имя записи
	dType       TypeData // тип записи
	meta        []byte   // описание
	data        []byte   // данные [для файла: путь]
	isEncrypted bool
}

func NewData(id, name string, dType TypeData, meta, data []byte, isEncrypted bool) *Data {
	return &Data{
		id:          id,
		name:        name,
		dType:       dType,
		meta:        meta,
		data:        data,
		isEncrypted: isEncrypted,
	}
}

func (d *Data) Decrypt(fnDec func([]byte) ([]byte, error)) error {
	if !d.isEncrypted {
		return errEncode
	}

	decData, err := fnDec(d.data)
	if err != nil {
		return fmt.Errorf("data decrypr: %w", err)
	}

	// так не работает, надо менять базовый массив
	// d.data = append(d.data[:0], decData...)

	d.data = decData
	d.isEncrypted = false

	return nil
}

func (d *Data) Encrypt(fnEnc func([]byte) ([]byte, error)) error {
	if d.isEncrypted {
		return errDecode
	}

	encData, err := fnEnc(d.data)
	if err != nil {
		return fmt.Errorf("data encrypt: %w", err)
	}

	d.data = append(d.data[:0], encData...)
	d.isEncrypted = true

	return nil
}

func (d *Data) ID() string     { return d.id }
func (d *Data) Name() string   { return d.name }
func (d *Data) Type() TypeData { return d.dType }
func (d *Data) Meta() []byte   { return d.meta }
func (d *Data) Data() []byte   { return d.data } // удалить

func NewLogPassData(name string, meta []byte, userData []byte) (*Data, error) {
	data := userData // сделать

	return &Data{
		id:          NewID().String(),
		name:        name,
		dType:       LogPassData,
		meta:        meta,
		data:        data,
		isEncrypted: false,
	}, nil
}

func NewBinaryData(name string, meta []byte, userData []byte) (*Data, error) {
	filePath := userData // сделать

	// проверяем существование файла
	if _, err := os.Stat(string(filePath)); err != nil {
		return nil, fmt.Errorf("check file: %w", err)
	}

	return &Data{
		id:          NewID().String(),
		name:        name,
		dType:       BinaryData,
		meta:        meta,
		data:        filePath,
		isEncrypted: false,
	}, nil
}
