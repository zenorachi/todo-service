package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zenorachi/todo-service/internal/entity"
)

type AgendaRepository struct {
	db *sql.DB
}

func NewAgenda(db *sql.DB) *AgendaRepository {
	return &AgendaRepository{db: db}
}

func (a *AgendaRepository) Create(ctx context.Context, task entity.Task) (int, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var (
		id    int
		query = fmt.Sprintf("INSERT INTO %s (user_id, title, description, date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			collectionAgenda)
	)

	err = tx.QueryRowContext(ctx, query, task.UserID, task.Title, task.Description, task.Date, task.Status).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (a *AgendaRepository) GetByID(ctx context.Context, id int, userId int) (entity.Task, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Task{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var (
		task  entity.Task
		query = fmt.Sprintf("SELECT title, description, date, status FROM %s WHERE id = $1 AND user_id = $2",
			collectionAgenda)
	)

	err = tx.QueryRowContext(ctx, query, id, userId).
		Scan(&task.Title, &task.Description, &task.Date, &task.Status)
	if err != nil {
		return entity.Task{}, err
	}

	return task, tx.Commit()
}

func (a *AgendaRepository) GetByTitleAndUserID(ctx context.Context, title string, userId int) (entity.Task, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Task{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var (
		task  entity.Task
		query = fmt.Sprintf("SELECT title, description, date, status FROM %s WHERE title = $1 AND user_id = $2",
			collectionAgenda)
	)

	err = tx.QueryRowContext(ctx, query, title, userId).
		Scan(&task.Title, &task.Description, &task.Date, &task.Status)
	if err != nil {
		return entity.Task{}, err
	}

	return task, tx.Commit()
}

func (a *AgendaRepository) SetStatus(ctx context.Context, id int, userId int, status string) error {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2 AND user_id = $3",
		collectionAgenda)

	_, err = tx.ExecContext(ctx, query, status, id, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (a *AgendaRepository) DeleteByID(ctx context.Context, id int, userId int) error {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var query = fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", collectionAgenda)

	_, err = tx.ExecContext(ctx, query, id, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (a *AgendaRepository) DeleteByUserID(ctx context.Context, userId int) error {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var query = fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", collectionAgenda)

	_, err = tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (a *AgendaRepository) GetByUserID(ctx context.Context, userId int) ([]entity.Task, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var (
		tasks []entity.Task
		query = fmt.Sprintf("SELECT id, title, description, date, status FROM %s WHERE user_id = $1",
			collectionAgenda)
	)

	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var task entity.Task
		if err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, tx.Commit()
}

func (a *AgendaRepository) GetByDateAndStatus(ctx context.Context, userId int, status string, date time.Time, limit, offset int) ([]entity.Task, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var (
		tasks []entity.Task
		query string
		rows  *sql.Rows
	)

	if date.Equal(time.Time{}) {
		query = fmt.Sprintf(
			"SELECT id, title, description, date, status FROM %s WHERE user_id = $1 AND status = $2 LIMIT $3 OFFSET $4",
			collectionAgenda)
		rows, err = tx.QueryContext(ctx, query, userId, status, limit, offset)
	} else {
		query = fmt.Sprintf(
			"SELECT id, title, description, date, status FROM %s WHERE user_id = $1 AND status = $2 AND DATE(date) = $3 LIMIT $4 OFFSET $5",
			collectionAgenda)
		rows, err = tx.QueryContext(ctx, query, userId, status, date, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var task entity.Task
		if err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, tx.Commit()
}
