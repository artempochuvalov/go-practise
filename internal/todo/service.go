package todo

import (
	"strings"
)

type TodoService interface {
	CreateTodo(title string) (int, error)
	GetTodo(id int) (Todo, error)
	GetAllTodos() ([]Todo, error)
	UpdateTitle(id int, title string) error
	MarkAsDone(id int) error
	DeleteTodo(id int) error
}

type Service struct {
	repo TodoRepository
}

func (s *Service) CreateTodo(title string) (int, error) {
	if strings.TrimSpace(title) == "" {
		return 0, EmptyTitleError
	}

	return s.repo.Create(title)
}
func (s *Service) GetTodo(id int) (Todo, error) {
	return s.repo.GetByID(id)
}
func (s *Service) GetAllTodos() ([]Todo, error) {
	return s.repo.GetAll()
}
func (s *Service) UpdateTitle(id int, title string) error {
	if strings.TrimSpace(title) == "" {
		return EmptyTitleError
	}

	return s.repo.UpdateTitle(id, title)
}
func (s *Service) MarkAsDone(id int) error {
	return s.repo.MarkAsDone(id)
}
func (s *Service) DeleteTodo(id int) error {
	return s.repo.Delete(id)
}

func NewService(repo TodoRepository) TodoService {
	return &Service{repo}
}
