package model

import "testing"

func TestUser(t *testing.T) {
	login := "LOGIN"
	pass := "PASS"

	user, err := NewUser(login, pass)
	if err != nil {
		t.Errorf("new user: %v\n", err)

		return
	}

	if user.ID() == "" || user.Login() != login {
		t.Error("no eq")

		return
	}

	if !user.CheckPass(pass) {
		t.Error("checkPass")

		return
	}
}
