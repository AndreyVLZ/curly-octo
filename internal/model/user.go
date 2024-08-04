package model

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ID uuid.UUID

func NewID() ID { return ID(uuid.New()) }

func (id ID) String() string { return uuid.UUID(id).String() }

type User struct {
	id    ID
	login string
	hash  string
}

func (u User) ID() string    { return u.id.String() }
func (u User) Login() string { return u.login }

func NewUser(login string, pass string) (*User, error) {
	hashe, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	user := &User{
		id:    NewID(),
		login: login,
		hash:  string(hashe),
	}

	return user, nil
}

func (u *User) CheckPass(pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.hash), []byte(pass))

	return err == nil
}
