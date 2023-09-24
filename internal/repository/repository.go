package repository

import (
	"context"
	"database/sql"
	"time"

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

	Agenda interface {
		Create(ctx context.Context, task entity.Task) (int, error)
		GetByID(ctx context.Context, id, userId int) (entity.Task, error)
		GetByTitleAndUserID(ctx context.Context, title string, userId int) (entity.Task, error)
		SetStatus(ctx context.Context, id, userId int, status string) error
		DeleteByID(ctx context.Context, id, userId int) error
		DeleteByUserID(ctx context.Context, userId int) error
		GetByUserID(ctx context.Context, userId int) ([]entity.Task, error)
		GetByDateAndStatus(ctx context.Context, userId int, status string, date time.Time, limit, offset int) ([]entity.Task, error)
	}
)

type Repositories struct {
	Users
	Agenda
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		Users:  NewUsers(db),
		Agenda: NewAgenda(db),
	}
}
