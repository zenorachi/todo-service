package repository

import (
	"context"
	"database/sql"
	"fmt"
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
		query = fmt.Sprintf("INSERT INTO %s (title, description, data, status) VALUES ($1, $2, $3, $4) RETURNING id",
			collectionAgenda)
	)

	err = tx.QueryRowContext(ctx, query, task.Title, task.Description, task.Data, task.Status).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (a *AgendaRepository) GetByID(ctx context.Context, id int) (entity.Task, error) {
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
		query = fmt.Sprintf("SELECT title, description, data, status FROM %s WHERE id = $1",
			collectionAgenda)
	)

	err = tx.QueryRowContext(ctx, query, id).
		Scan(&task.Title, &task.Description, &task.Data, &task.Status)
	if err != nil {
		return entity.Task{}, err
	}

	return task, tx.Commit()
}

func (a *AgendaRepository) SetStatus(ctx context.Context, id int, status string) error {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2",
		collectionAgenda)

	_, err = tx.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (a *AgendaRepository) DeleteByID(ctx context.Context, id int) error {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var query = fmt.Sprintf("DELETE FROM %s WHERE id = $1", collectionAgenda)

	_, err = tx.ExecContext(ctx, query, id)
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
		query = fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1",
			collectionAgenda)
	)

	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var task entity.Task
		if err = rows.Scan(&task.Title, &task.Description, &task.Data, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, tx.Commit()
}
