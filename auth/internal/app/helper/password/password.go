package password

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	salt string
}

func New() *Hasher {
	return &Hasher{salt: os.Getenv("USER_PASSWORD_SALT")}
}

func (p *Hasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+p.salt), 14)
	return string(bytes), err
}

func (p *Hasher) CheckHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+p.salt))
	return err == nil
}
