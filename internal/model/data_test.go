package model

import (
	"bytes"
	"testing"
)

func TestDataType(t *testing.T) {
	for i, nameType := range namesType {
		if TypeData(i).String() != nameType {
			t.Error("no eq")
		}
	}
}
func TestData(t *testing.T) {
	id := "ID"
	name := "name"
	dType := LogPassData
	meta := []byte("meta")
	userData := []byte("user data")

	myData := NewData(
		id,
		name,
		dType,
		meta,
		userData,
		false,
	)

	if myData.ID() != id || myData.Name() != name || myData.Type() != dType || !bytes.Equal(myData.Meta(), meta) || !bytes.Equal(myData.Data(), userData) {
		t.Error("no eq")

		return
	}

	encMsg := []byte("testEnc")

	fnEnc := func([]byte) ([]byte, error) {
		return encMsg, nil
	}

	myData.Encrypt(fnEnc)

	if !bytes.Equal(myData.Data(), encMsg) {
		t.Error("enc data")

		return
	}

	decMsg := []byte("testDec")

	fnDec := func([]byte) ([]byte, error) {
		return decMsg, nil
	}

	myData.Decrypt(fnDec)

	if !bytes.Equal(myData.Data(), decMsg) {
		t.Error("dec data")

		return
	}
}
