package jwt

import (
	"testing"
	"time"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	key := "KEY"
	exp := time.Minute

	user, err := model.NewUser("", "")
	if err != nil {
		t.Errorf("new user: %v\n", err)

		return
	}

	jwt := New(key, exp)

	token, err := jwt.Generate(user)
	if err != nil {
		t.Errorf("jwt gen: %v\n", err)

		return
	}

	userClaims, err := jwt.Verify(token)
	if err != nil {
		t.Errorf("jwt verify: %v\n", err)

		return
	}

	assert.Equal(t, userClaims.UserID, user.ID())
}
