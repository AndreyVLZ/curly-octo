package agent

/*
func TestEnc(t *testing.T) {
	key := []byte("MY SECRET KEY")
	hashKey := crypto.CreateHash(key)

	fnEnc, err := crypto.Encrypt(hashKey)
	if err != nil {
		t.Errorf("new fnEnc: %v\n", err)

		return
	}

	userData := []byte("Very secret in file 1 Very secret in file 2 Very secret in file 3 Very secret in file 4 Very secret in file 5 Very secret in file 6 Very secret in file 7 Very secret in file 8 Very secret in file 9 Very secret in file 10 ")
	fmt.Printf("lenUserData: %d\n", len(userData))
	dat, err := model.NewLogPassData(userData)
	if err != nil {
		t.Errorf("new logPass: %v\n", err)

		return
	}

	fn1 := func(d []byte) ([]byte, error) {
		return fnEnc(d), nil
	}

	encData := NewCryptoData(fn1, dat, bytes.NewReader(dat.Data()), 100)

	resEnc, err := io.ReadAll(encData)
	if err != nil {
		t.Errorf("readAll: %v\n", err)

		return
	}

	t.Logf("len resEnc: %d\n", len(resEnc))
}
*/

/*
func TwestEncFile(t *testing.T) {
	key := []byte("MY SECRET KEY")
	hashKey := crypto.CreateHash(key)

	fnEnc, err := crypto.Encrypt(hashKey)
	if err != nil {
		t.Errorf("new fnEnc: %v\n", err)

		return
	}

	dataFile, err := model.NewBinaryData([]byte("secretFile.txt"))
	if err != nil {
		t.Errorf("new binary data: %v\n", err)

		return
	}
	t.Logf("lenFile: %d\n", dataFile.Len())

	list := []*model.Data{dataFile}

	if err := buildSecretArray1(fnEnc, list); err != nil {
		t.Errorf("build secret array: %v\n", err)

		return
	}

	var resBuf bytes.Buffer

	for _, el := range list {
		buf := make([]byte, 100)

		for {
			t.Logf("len_ost: %d\n", el.Len())
			n, err := el.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Errorf("read err: %v\n", err)

				return
			}
			t.Logf("lenBytes: %d\n", n)

			if _, err := resBuf.Write(buf[:n]); err != nil {
				t.Errorf("write res: %v\n", err)
			}
		}
	}

	t.Logf("Res: %d\n", len(resBuf.Bytes()))

	//listEnc := []*model.Data{}

}
func T1estEncrypt(t *testing.T) {
	key := []byte("MY SECRET KEY")
	hashKey := crypto.CreateHash(key)

	fnEnc, err := crypto.Encrypt(hashKey)
	if err != nil {
		t.Errorf("new fnEnc: %v\n", err)

		return
	}

	mData1 := model.NewLogPassData([]byte("myLogPass"))

	list := []*model.Data{mData1}

	if err := buildSecretArray1(fnEnc, list); err != nil {
		t.Errorf("build: %v\n", err)

		return
	}

	var buf01 bytes.Buffer

	for _, el := range list {
		t.Logf("len enc data0: %d\n", el.Len())
		buf := make([]byte, el.Len())

		n, err := io.ReadFull(el, buf)
		if err != nil {
			t.Errorf("read full: %v\n", err)

			return
		}

		t.Logf("n: %d\n", n)
		t.Logf("res: %s\n", hex.EncodeToString(buf))
		t.Logf("len enc data: %d\n", el.Len())
		buf01.Write(buf)
	}

	fnDec, err := crypto.Decrypt(hashKey)
	if err != nil {
		t.Errorf("new fnDec: %v\n", err)

		return
	}

	mddd := model.NewData(
		"id",
		"name",
		model.LogPassData,
		bytes.NewBuffer(buf01.Bytes()),
		[]byte("meta"),
	)

	listEnc := []*model.Data{mddd}

	if err := decryptSecretArray(fnDec, listEnc); err != nil {
		t.Errorf("dec array: %v\n", err)

		return
	}
	bud := make([]byte, listEnc[0].Len())
	listEnc[0].Read(bud)
	fmt.Printf("ress: %s\n", string(bud))
}

func TEncrypt(t *testing.T) {
	key := []byte("MY SECRET KEY")
	hashKey := crypto.CreateHash(key)

	mData1 := model.NewLogPassData([]byte("myLogPass"))
	/*
			buf0 := make([]byte, mData1.SecretData.Len())
			n, err := mData1.SecretData.Read(buf0)
			if err != nil {
				t.Errorf("read: %v\n", err)

				return
			}

			t.Logf("n: %d\n", n)
			t.Logf("buf: %s\n", string(buf0))

			return

		for _, el := range list {
			buf := make([]byte, el.Len())
			n, err := io.ReadFull(el, buf)
			if err != nil {
				t.Errorf("read encData: %v\n", err)

				return
			}

			t.Logf("n: %d\n", n)
			t.Logf("res: %s\n", string(buf))
			t.Logf("len data: %d\n", el.Len())
		}

	list := []model.Dater{mData1}
	encList, err := buildSecretArray(hashKey, list)
	if err != nil {
		t.Errorf("build: %v\n", err)

		return
	}

	/*
		// ReadAll
		for _, el := range encList {
			t.Logf("len enc data0: %d\n", el.Len())
			buf, err := io.ReadAll(el)
			if err != nil {
				t.Errorf("read all: %v\n", err)

				return
			}

			t.Logf("res: %s\n", hex.EncodeToString(buf))
			t.Logf("len enc data: %d\n", el.Len())
		}

	// ReadFull
	for _, el := range encList {
		t.Logf("len enc data0: %d\n", el.Len())

		buf := make([]byte, el.Len())
		n, err := io.ReadFull(el, buf)
		if err != nil {
			t.Errorf("read full: %v\n", err)

			return
		}

		t.Logf("n: %d\n", n)
		t.Logf("res: %s\n", hex.EncodeToString(buf))
		t.Logf("len enc data: %d\n", el.Len())
	}
}
*/
