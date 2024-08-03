package jwt

import (
	"fmt"
	"time"

	"github.com/AndreyVLZ/curly-octo/internal/model"
	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	key       string
	expiresAt time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	UserID string `json:"userid"`
}

func New(key string, expiresAt time.Duration) *JWT {
	return &JWT{
		key:       key,
		expiresAt: expiresAt,
	}
}

// Generate Создание токена.
func (j *JWT) Generate(user *model.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.expiresAt).Unix(),
		},
		UserID: user.ID(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.key))
}

// Verify Проверка токена.
func (j *JWT) Verify(accToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accToken, &UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(j.key), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
