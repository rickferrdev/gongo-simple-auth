package platform

import (
	"errors"

	"github.com/rickferrdev/gongo-simple-auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	cost int
}

func NewHasher(cost int) Hasher {
	return Hasher{
		cost: cost,
	}
}

func (u *Hasher) capture(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return domain.ErrInvalidCredentials
	default:
		return domain.ErrInternal
	}
}

func (u *Hasher) GenerateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), u.cost)
	if err != nil {
		return "", u.capture(err)
	}

	return string(bytes), nil
}

func (u *Hasher) ValidateHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return u.capture(err)
}
