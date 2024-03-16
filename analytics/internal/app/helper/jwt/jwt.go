package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

type Manager struct {
	secretKey []byte
}

func New() *Manager {
	return &Manager{secretKey: []byte(os.Getenv("JWT_SECRET_KEY"))}
}

type UserClaims struct {
	jwt.RegisteredClaims
	PublicID uuid.UUID `json:"public_id"`
}

func NewUserClaims(publicID uuid.UUID) UserClaims {
	return UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		PublicID: publicID,
	}
}

func (j *Manager) NewAccessToken(claims UserClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secretKey)
}

func (j *Manager) ParseAccessToken(token string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedAccessToken.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
