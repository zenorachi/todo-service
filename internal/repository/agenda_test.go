package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zenorachi/todo-service/internal/entity"
)

func TestAgendaRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewAgenda(db)

	testTask := entity.Task{
		UserID:      1,
		Title:       "Test Task",
		Description: "Test Description",
		Date:        time.Now().Round(time.Second),
		Status:      entity.StatusNotDone,
	}

	type args struct {
		task entity.Task
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
				task: testTask,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO agenda (user_id, title, description, date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.task.UserID, args.task.Title, args.task.Description, args.task.Date, args.task.Status).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				task: entity.Task{},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO agenda (user_id, title, description, date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.task.UserID, args.task.Title, args.task.Description, args.task.Date, args.task.Status).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			_, err := repo.Create(context.Background(), tt.args.task)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAgendaRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewAgenda(db)

	testTaskID := 1
	testTask := entity.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Date:        time.Now().Round(time.Second),
		Status:      entity.StatusNotDone,
	}

	type args struct {
		id     int
		userId int
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantTask      entity.Task
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				id: testTaskID,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"title", "description", "date", "status"}).
					AddRow(testTask.Title, testTask.Description, testTask.Date, testTask.Status)

				expectedQuery := "SELECT title, description, date, status FROM agenda WHERE id = $1 AND user_id = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id, args.userId).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantTask: testTask,
		},
		{
			name: "ERROR",
			args: args{
				id: testTaskID,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT title, description, date, status FROM agenda WHERE id = $1 AND user_id = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id, args.userId).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantTask: entity.Task{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			task, err := repo.GetByID(context.Background(), tt.args.id, tt.args.userId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTask, task)
			}
		})
	}
}

func TestAgendaRepository_GetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewAgenda(db)

	testTaskTitle := "Test Task"
	testTask := entity.Task{
		Title:       testTaskTitle,
		Description: "Test Description",
		Date:        time.Now().Round(time.Second),
		Status:      entity.StatusNotDone,
	}

	type args struct {
		title  string
		userId int
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantTask      entity.Task
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				title: testTaskTitle,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"title", "description", "date", "status"}).
					AddRow(testTask.Title, testTask.Description, testTask.Date, testTask.Status)

				expectedQuery := "SELECT title, description, date, status FROM agenda WHERE title = $1 AND user_id = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.title, args.userId).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantTask: testTask,
		},
		{
			name: "ERROR",
			args: args{
				title: testTaskTitle,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT title, description, date, status FROM agenda WHERE title = $1 AND user_id = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.title, args.userId).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantTask: entity.Task{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			task, err := repo.GetByTitleAndUserID(context.Background(), tt.args.title, tt.args.userId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTask, task)
			}
		})
	}
}

func TestAgendaRepository_SetStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewAgenda(db)

	type args struct {
		id     int
		userId int
		status string
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
				id:     1,
				status: entity.StatusDone,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "UPDATE agenda SET status = $1 WHERE id = $2 AND user_id = $3"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.status, args.id, args.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				id:     2,
				status: entity.StatusNotDone,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "UPDATE agenda SET status = $1 WHERE id = $2 AND user_id = $3"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.status, args.id, args.userId).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			err := repo.SetStatus(context.Background(), tt.args.id, tt.args.userId, tt.args.status)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAgendaRepository_DeleteByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewAgenda(db)

	type args struct {
		id     int
		userId int
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
				id: 1,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM agenda WHERE id = $1 AND user_id = $2"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.id, args.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				id: 2,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM agenda WHERE id = $1 AND user_id = $2"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.id, args.userId).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			err := repo.DeleteByID(context.Background(), tt.args.id, tt.args.userId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAgendaRepository_DeleteByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewAgenda(db)

	type args struct {
		userID int
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
				userID: 1,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM agenda WHERE user_id = $1"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ERROR",
			args: args{
				userID: 2,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM agenda WHERE user_id = $1"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			err := repo.DeleteByUserID(context.Background(), tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAgendaRepository_GetByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewAgenda(db)

	type args struct {
		userID int
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name           string
		args           args
		mockBehaviour  mockBehaviour
		wantErr        bool
		expectedResult []entity.Task
	}{
		{
			name: "OK",
			args: args{
				userID: 1,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, title, description, date, status FROM agenda WHERE user_id = $1"
				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "status"}).
					AddRow(0, "Task 1", "Description 1", time.Now().Round(time.Second), "done").
					AddRow(0, "Task 2", "Description 2", time.Now().Round(time.Second), "not done")

				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			expectedResult: []entity.Task{
				{
					Title:       "Task 1",
					Description: "Description 1",
					Date:        time.Now().Round(time.Second),
					Status:      "done",
				},
				{
					Title:       "Task 2",
					Description: "Description 2",
					Date:        time.Now().Round(time.Second),
					Status:      "not done",
				},
			},
		},
		{
			name: "ERROR",
			args: args{
				userID: 2,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, title, description, date, status FROM agenda WHERE user_id = $1"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			tasks, err := repo.GetByUserID(context.Background(), tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, tasks)
			}
		})
	}
}

func TestAgendaRepository_GetByDateAndStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating database connection: %v\n", err)
	}
	defer db.Close()

	repo := NewAgenda(db)

	type args struct {
		userID int
		status string
		date   time.Time
		limit  int
		offset int
	}
	type mockBehaviour func(args args)

	tests := []struct {
		name           string
		args           args
		mockBehaviour  mockBehaviour
		wantErr        bool
		expectedResult []entity.Task
	}{
		{
			name: "OK with date",
			args: args{
				userID: 1,
				status: "completed",
				date:   time.Now().Round(time.Second),
				limit:  10,
				offset: 0,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, title, description, date, status FROM agenda WHERE user_id = $1 AND status = $2 AND DATE(date) = $3 LIMIT $4 OFFSET $5"
				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "status"}).
					AddRow(0, "Task 1", "Description 1", time.Now().Round(time.Second), "done").
					AddRow(0, "Task 2", "Description 2", time.Now().Round(time.Second), "done")

				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID, args.status, args.date, args.limit, args.offset).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			expectedResult: []entity.Task{
				{
					Title:       "Task 1",
					Description: "Description 1",
					Date:        time.Now().Round(time.Second),
					Status:      "done",
				},
				{
					Title:       "Task 2",
					Description: "Description 2",
					Date:        time.Now().Round(time.Second),
					Status:      "done",
				},
			},
		},
		{
			name: "OK without date",
			args: args{
				userID: 1,
				status: "completed",
				date:   time.Time{},
				limit:  10,
				offset: 0,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, title, description, date, status FROM agenda WHERE user_id = $1 AND status = $2 LIMIT $3 OFFSET $4"
				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "status"}).
					AddRow(0, "Task 1", "Description 1", time.Time{}, "done").
					AddRow(0, "Task 2", "Description 2", time.Time{}, "done")

				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID, args.status, args.limit, args.offset).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			expectedResult: []entity.Task{
				{
					Title:       "Task 1",
					Description: "Description 1",
					Date:        time.Time{},
					Status:      "done",
				},
				{
					Title:       "Task 2",
					Description: "Description 2",
					Date:        time.Time{},
					Status:      "done",
				},
			},
		},
		{
			name: "ERROR",
			args: args{
				userID: 2,
				status: "completed",
				date:   time.Now().Round(time.Second),
				limit:  10,
				offset: 0,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "SELECT id, title, description, date, status FROM agenda WHERE user_id = $1 AND status = $2 AND DATE(date) = $3 LIMIT $4 OFFSET $5"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userID, args.status, args.date, args.limit, args.offset).
					WillReturnError(errors.New("test error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)
			tasks, err := repo.GetByDateAndStatus(context.Background(), tt.args.userID, tt.args.status, tt.args.date, tt.args.limit, tt.args.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, tasks)
			}
		})
	}
}
