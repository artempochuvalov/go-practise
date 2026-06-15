package todo

type Service struct {
	repo *Repository
}

func (s *Service) CreateTodo(title string) (int, error) {
	return s.repo.Create(title)
}
func (s *Service) GetTodo(id int) (Todo, error) {
	return s.repo.GetByID(id)
}
func (s *Service) GetAllTodos() ([]Todo, error) {
	return s.repo.GetAll()
}
func (s *Service) UpdateTitle(id int, title string) error {
	return s.repo.UpdateTitle(id, title)
}
func (s *Service) MarkAsDone(id int) error {
	return s.repo.MarkAsDone(id)
}
func (s *Service) DeleteTodo(id int) error {
	return s.repo.Delete(id)
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}
