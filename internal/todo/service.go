package todo

import (
	"errors"
	"strings"
)

type Service struct {
	repo TodoRepository
}

func (s *Service) CreateTodo(title string) (int, error) {
	if strings.TrimSpace(title) == "" {
		return 0, errors.New("title cannot be empty")
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
		return errors.New("title cannot be empty")
	}

	return s.repo.UpdateTitle(id, title)
}
func (s *Service) MarkAsDone(id int) error {
	return s.repo.MarkAsDone(id)
}
func (s *Service) DeleteTodo(id int) error {
	return s.repo.Delete(id)
}

func NewService(repo TodoRepository) *Service {
	return &Service{repo}
}
