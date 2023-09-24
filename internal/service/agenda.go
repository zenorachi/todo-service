package service

import (
	"context"
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

}

func (a *AgendaService) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {

}

func (a *AgendaService) SetTaskStatus(ctx context.Context, id int, status string) error {

}

func (a *AgendaService) DeleteTaskByID(ctx context.Context, id int) error {

}

func (a *AgendaService) DeleteUserTasks(ctx context.Context, userId int) error {

}

func (a *AgendaService) GetUserTasks(ctx context.Context, userId int) ([]entity.Task, error) {

}
