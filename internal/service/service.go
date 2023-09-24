package service

import (
	"context"
	"time"

	"github.com/zenorachi/todo-service/internal/repository"
	"github.com/zenorachi/todo-service/pkg/auth"
	"github.com/zenorachi/todo-service/pkg/hash"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type (
	Users interface {
		SignUp(ctx context.Context, login, email, password string) (int, error)
		SignIn(ctx context.Context, login, password string) (Tokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	}
)

type Services struct {
	Users
}

type Deps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func New(deps Deps) *Services {
	return &Services{
		Users: NewUsers(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}
