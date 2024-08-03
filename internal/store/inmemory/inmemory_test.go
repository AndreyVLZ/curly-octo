package inmemory

import (
	"context"
	"testing"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	ctx := context.Background()
	store := New()
	login := "LOGIN"
	pass := "PASS"
	user, err := model.NewUser(login, pass)
	if err != nil {
		t.Errorf("new user: %v\n", err)

		return
	}

	if err := store.Add(ctx, *user); err != nil {
		t.Errorf("add: %v\n", err)

		return
	}

	userInDB, err := store.FindByName(login)
	if err != nil {
		t.Errorf("find user: %v\n", err)

		return
	}

	assert.Equal(t, user, userInDB)

	logPass := model.NewData(
		"ID-1",
		"LogPass-1",
		model.LogPassData,
		[]byte("meta"),
		[]byte("user data"),
		false,
	)

	if err := store.SaveData(ctx, user.ID(), *logPass); err != nil {
		t.Errorf("save data: %v\n", err)

		return
	}

	logPassInDB, err := store.GetData(ctx, user.ID(), logPass.ID())
	if err != nil {
		t.Errorf("get data: %v\n", err)

		return
	}

	assert.Equal(t, logPass, &logPassInDB)

	binary := model.NewData(
		"ID-2",
		"binary-1",
		model.BinaryData,
		[]byte("meta-1"),
		[]byte("user data"),
		true,
	)

	arr := []*model.Data{binary}

	if err := store.SaveArray(ctx, user.ID(), arr); err != nil {
		t.Errorf("save array: %v\n", err)

		return
	}

	listInDB, err := store.List(ctx, user.ID())
	if err != nil {
		t.Errorf("list: %v\n", err)

		return
	}

	if len(listInDB) != 2 {
		t.Error("no eq")
	}
}
