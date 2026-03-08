package auth

import (
	"context"
	"errors"

	"github.com/rickferrdev/gongo-simple-auth/internal/domain"
	"github.com/rickferrdev/gongo-simple-auth/internal/platform"
)

var (
	ErrHashingFailed   = errors.New("failed to process password")
	ErrTokenGeneration = errors.New("could not generate access token")
)

type AuthService struct {
	storage   AuthStorage
	tokenizer platform.Tokenizer
	hasher    platform.Hasher
}

type LoginInput struct {
	Email    string
	Password string
}

type RegisterInput struct {
	Email    string
	Username string
	Password string
}

type RegisterOutput struct {
	ID string
}

type LoginOutput struct {
	ID    string
	Token string
}

func NewAuthService(storage AuthStorage, tokenizer platform.Tokenizer, hasher platform.Hasher) AuthService {
	return AuthService{
		storage:   storage,
		tokenizer: tokenizer,
		hasher:    hasher,
	}
}

func (u *AuthService) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	user, err := u.storage.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := u.hasher.ValidateHash(user.Password, input.Password); err != nil {
		return nil, err
	}

	token, err := u.tokenizer.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{ID: user.ID, Token: token}, nil
}

func (u *AuthService) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	hash, err := u.hasher.GenerateHash(input.Password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email:    input.Email,
		Username: input.Username,
		Password: hash,
	}

	created, err := u.storage.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{ID: created.ID}, nil
}
