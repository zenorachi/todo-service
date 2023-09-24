package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zenorachi/todo-service/internal/entity"
	"github.com/zenorachi/todo-service/internal/repository"
)

type AgendaService struct {
	repo repository.Agenda
}

func NewAgenda(repo repository.Agenda) *AgendaService {
	return &AgendaService{repo: repo}
}

func (a *AgendaService) CreateTask(ctx context.Context, task entity.Task) (int, error) {
	if a.isTaskExists(ctx, task.Title) {
		return 0, entity.ErrTaskAlreadyExist
	}

	if len(task.Status) != 0 && task.Status != entity.StatusDone && task.Status != entity.StatusNotDone {
		return 0, entity.ErrInvalidStatus
	}

	if len(task.Status) == 0 {
		task.Status = entity.StatusNotDone
	}

	return a.repo.Create(ctx, task)
}

func (a *AgendaService) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	task, err := a.repo.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Task{}, entity.ErrTaskDoesNotExist
	}
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (a *AgendaService) SetTaskStatus(ctx context.Context, id int, status string) error {
	if !a.isTaskExists(ctx, id) {
		return entity.ErrTaskDoesNotExist
	}

	if status != entity.StatusDone && status != entity.StatusNotDone {
		return entity.ErrInvalidStatus
	}

	return a.repo.SetStatus(ctx, id, status)
}

func (a *AgendaService) DeleteTaskByID(ctx context.Context, id int) error {
	if !a.isTaskExists(ctx, id) {
		return entity.ErrTaskDoesNotExist
	}

	return a.repo.DeleteByID(ctx, id)
}

func (a *AgendaService) DeleteUserTasks(ctx context.Context, userId int) error {
	return a.repo.DeleteByUserID(ctx, userId)
}

func (a *AgendaService) GetUserTasks(ctx context.Context, userId int) ([]entity.Task, error) {
	return a.repo.GetByUserID(ctx, userId)
}

func (a *AgendaService) GetByDateAndStatus(ctx context.Context, userId int, status string, date time.Time, limit, offset int) ([]entity.Task, error) {
	if status == "not_done" {
		status = entity.StatusNotDone
	}

	if status != entity.StatusDone && status != entity.StatusNotDone {
		return nil, entity.ErrInvalidStatus
	}

	if limit < 0 || offset < 0 {
		return nil, entity.ErrInvalidPaginationSizes
	}

	return a.repo.GetByDateAndStatus(ctx, userId, status, date, limit, offset)
}

func (a *AgendaService) isTaskExists(ctx context.Context, data any) bool {
	var task entity.Task

	switch data.(type) {
	case string:
		task, _ = a.repo.GetByTitle(ctx, data.(string))
	case int:
		task, _ = a.repo.GetByID(ctx, data.(int))
	}

	return len(task.Title) != 0
}
