package service

import (
	"context"
	"github.com/zenorachi/todo-service/internal/entity"
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

	Agenda interface {
		CreateTask(ctx context.Context, task entity.Task) (int, error)
		GetTaskByID(ctx context.Context, id int) (entity.Task, error)
		SetTaskStatus(ctx context.Context, id int, status string) error
		DeleteTaskByID(ctx context.Context, id int) error
		DeleteUserTasks(ctx context.Context, userId int) error
		GetUserTasks(ctx context.Context, userId int) ([]entity.Task, error)
	}
)

type Services struct {
	Users
	Agenda
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
		Users:  NewUsers(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Agenda: NewAgenda(deps.Repos.Agenda),
	}
}
