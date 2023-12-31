package repository

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zenorachi/todo-service/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestUsersRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewUsers(db)

	type args struct {
		user entity.User
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				user: entity.User{
					ID:           1,
					Login:        "user",
					Email:        "email@go.dev",
					Password:     "password",
					RegisteredAt: time.Now().Round(time.Second),
				},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO users (login, email, password) VALUES ($1, $2, $3) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.user.Login, args.user.Email, args.user.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				user: entity.User{},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO users (login, email, password) VALUES ($1, $2, $3)"
				mock.ExpectExec(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.user.Login, args.user.Email, args.user.Password).
					WillReturnError(fmt.Errorf("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			_, err := repo.Create(context.Background(), tt.args.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUsersRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewUsers(db)

	type args struct {
		id int
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		args          args
		user          entity.User
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			user: entity.User{
				ID:           1,
				Login:        "login",
				Email:        "user-email",
				Password:     "password",
				RegisteredAt: time.Now().Round(time.Second),
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "login", "email", "password", "registered_at"}).
					AddRow(1, "login", "user-email", "password", time.Now().Round(time.Second))

				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id).WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				id: 1,
			},
			user: entity.User{},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			user, err := repo.GetByID(context.Background(), tt.args.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, user)
			}
		})
	}
}

func TestUsersRepository_GetByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewUsers(db)

	type args struct {
		login string
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		args          args
		user          entity.User
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				login: "login",
			},
			user: entity.User{
				ID:           1,
				Login:        "login",
				Email:        "user-email",
				Password:     "password",
				RegisteredAt: time.Now().Round(time.Second),
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "login", "email", "password", "registered_at"}).
					AddRow(1, "login", "user-email", "password", time.Now().Round(time.Second))

				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE login = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.login).WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				login: "login",
			},
			user: entity.User{},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.login).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			user, err := repo.GetByLogin(context.Background(), tt.args.login)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, user)
			}
		})
	}
}

func TestUsersRepository_GetByCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewUsers(db)

	type args struct {
		login    string
		password string
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		args          args
		user          entity.User
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				login:    "login",
				password: "password",
			},
			user: entity.User{
				ID:           1,
				Login:        "login",
				Email:        "user-email",
				Password:     "password",
				RegisteredAt: time.Now().Round(time.Second),
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "login", "email", "password", "registered_at"}).
					AddRow(1, "login", "user-email", "password", time.Now().Round(time.Second))

				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE login = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.login, args.password).WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				login:    "hello",
				password: "world",
			},
			user: entity.User{},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, login, email, password, registered_at FROM users WHERE login = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.login, args.password).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			user, err := repo.GetByCredentials(context.Background(), tt.args.login, tt.args.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, user)
			}
		})
	}
}
