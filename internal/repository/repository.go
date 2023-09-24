package repository

import (
	"context"
	"database/sql"

	"github.com/zenorachi/todo-service/internal/entity"
)

type (
	Users interface {
		Create(ctx context.Context, user entity.User) (int, error)
		GetByID(ctx context.Context, id int) (entity.User, error)
		GetByLogin(ctx context.Context, login string) (entity.User, error)
		GetByCredentials(ctx context.Context, login, password string) (entity.User, error)
		GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
		SetSession(ctx context.Context, userId int, session entity.Session) error
	}
)

type Repositories struct {
	Users
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		Users: NewUsers(db),
	}
}
