package todo

import (
	"context"
	"testing"
)

type MockRepository struct{}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (r *MockRepository) Create(ctx context.Context, title string) (int, error) {
	return 1, nil
}

func (r *MockRepository) GetByID(ctx context.Context, id int) (Todo, error) {
	return Todo{1, "title", false}, nil
}

func (r *MockRepository) GetAll(ctx context.Context) ([]Todo, error) {
	return []Todo{}, nil
}

func (r *MockRepository) UpdateTitle(ctx context.Context, id int, title string) error {
	return nil
}

func (r *MockRepository) MarkAsDone(ctx context.Context, id int) error {
	return nil
}

func (r *MockRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func setupService() TodoService {
	r := NewMockRepository()
	service := NewService(r)
	return service
}

func TestService_CreateTodo_Fail(t *testing.T) {
	service := setupService()

	tests := []struct {
		name  string
		title string
	}{
		{"empty", ""},
		{"spaces", "  "},
		{"tabs", "\t"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := service.CreateTodo(t.Context(), test.title)
			if err == nil {
				t.Fatalf("expected create to fail")
			}
		})
	}
}

func TestService_CreateTodo_Success(t *testing.T) {
	service := setupService()

	id, err := service.CreateTodo(t.Context(), "new todo")
	if err != nil {
		t.Fatalf("create todo failed: %v", err)
	}

	if id != 1 {
		t.Fatalf("unexpected todo id returned: %v", id)
	}
}
