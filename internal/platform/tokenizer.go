package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rickferrdev/gongo-simple-auth/internal/domain"
)

var (
	ErrInvalidSignatureMethod = errors.New("unexpected signing method")
)

type Tokenizer struct {
	secret []byte
}

func NewTokenizer(secret string) Tokenizer {
	return Tokenizer{
		secret: []byte(secret),
	}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (u *Tokenizer) capture(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return domain.ErrTokenMalformed
	case errors.Is(err, ErrInvalidSignatureMethod):
	case errors.Is(err, jwt.ErrSignatureInvalid):
		return domain.ErrTokenInvalid
	case errors.Is(err, jwt.ErrTokenExpired):
		return domain.ErrTokenExpired
	default:
		return domain.ErrInternal
	}

	return domain.ErrInternal
}

func (u Tokenizer) GenerateToken(id string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rickferrdev/gongo-simple-auth",
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(u.secret)
}

func (u Tokenizer) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		switch t.Method.(type) {
		case *jwt.SigningMethodHMAC:
			return u.secret, nil
		default:
			return nil, fmt.Errorf("%w: %v", ErrInvalidSignatureMethod, t.Header["alg"])
		}
	})

	if err != nil {
		return nil, u.capture(err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, u.capture(domain.ErrTokenInvalid)
}
