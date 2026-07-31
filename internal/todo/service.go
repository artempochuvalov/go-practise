package todo

import (
	"context"
	"strings"
)

type TodoService interface {
	CreateTodo(ctx context.Context, title string) (int, error)
	GetTodo(ctx context.Context, id int) (Todo, error)
	GetAllTodos(ctx context.Context) ([]Todo, error)
	UpdateTitle(ctx context.Context, id int, title string) error
	MarkAsDone(ctx context.Context, id int) error
	DeleteTodo(ctx context.Context, id int) error
}

type Service struct {
	repo TodoRepository
}

func (s *Service) CreateTodo(ctx context.Context, title string) (int, error) {
	if strings.TrimSpace(title) == "" {
		return 0, EmptyTitleError
	}

	return s.repo.Create(ctx, title)
}
func (s *Service) GetTodo(ctx context.Context, id int) (Todo, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) GetAllTodos(ctx context.Context) ([]Todo, error) {
	return s.repo.GetAll(ctx)
}
func (s *Service) UpdateTitle(ctx context.Context, id int, title string) error {
	if strings.TrimSpace(title) == "" {
		return EmptyTitleError
	}

	return s.repo.UpdateTitle(ctx, id, title)
}
func (s *Service) MarkAsDone(ctx context.Context, id int) error {
	return s.repo.MarkAsDone(ctx, id)
}
func (s *Service) DeleteTodo(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func NewService(repo TodoRepository) TodoService {
	return &Service{repo}
}
