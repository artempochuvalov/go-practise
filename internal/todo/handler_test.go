package todo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	receivedTitle string

	createErr  error
	receivedId int

	getTodo    Todo
	getTodoErr error

	getTodos    []Todo
	getTodosErr error

	updateTitleErr error

	markAsDoneErr error

	deleteErr error
}

func NewMockService() *MockService {
	return &MockService{}
}

func (service *MockService) CreateTodo(ctx context.Context, title string) (int, error) {
	service.receivedTitle = title
	return 1, service.createErr
}

func (service *MockService) GetTodo(ctx context.Context, id int) (Todo, error) {
	service.receivedId = id
	return service.getTodo, service.getTodoErr
}

func (service *MockService) GetAllTodos(ctx context.Context) ([]Todo, error) {
	return service.getTodos, service.getTodosErr
}

func (service *MockService) UpdateTitle(ctx context.Context, id int, title string) error {
	service.receivedId = id
	service.receivedTitle = title
	return service.updateTitleErr
}

func (service *MockService) MarkAsDone(ctx context.Context, id int) error {
	service.receivedId = id
	return service.markAsDoneErr
}

func (service *MockService) DeleteTodo(ctx context.Context, id int) error {
	service.receivedId = id
	return service.deleteErr
}

func setupHandler() (*Handler, *MockService) {
	service := NewMockService()
	var buf bytes.Buffer
	log := slog.New(
		slog.NewTextHandler(&buf, nil),
	)

	return NewHandler(service, log), service
}

func TestHandler_CreateTodo_Success(t *testing.T) {
	handler, service := setupHandler()

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{"title": "some title"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateTodo(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if service.receivedTitle != "some title" {
		t.Fatalf("wrong title passed to service")
	}

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusCreated, res.StatusCode)
	}

	var body struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		t.Fatalf("json parsing error: %v", err)
	}

	if body.ID != 1 {
		t.Fatalf("unexpected id: expected %d, received %d", 1, body.ID)
	}
}

func TestHandler_CreateTodo_InvalidJSON(t *testing.T) {
	handler, _ := setupHandler()

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{"title": "some title"`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateTodo(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusBadRequest, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "invalid json") {
		t.Fatalf("unexpected error message")
	}
}

func TestHandler_CreateTodo_EmptyTitle(t *testing.T) {
	handler, service := setupHandler()
	service.createErr = EmptyTitleError

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{"title": ""}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateTodo(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusBadRequest, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), EmptyTitleError.Error()) {
		t.Fatalf("unexpected error message")
	}
}

func TestHandler_GetTodo_Success(t *testing.T) {
	handler, service := setupHandler()
	service.getTodo = Todo{Title: "some title", ID: 1, Completed: false}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos/{id}", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if service.receivedId != 1 {
		t.Fatalf("unexpected id passed to service: expected %d, received %d", 1, service.receivedId)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusOK, res.StatusCode)
	}

	var todo Todo
	err := json.NewDecoder(res.Body).Decode(&todo)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if todo != service.getTodo {
		t.Fatalf("unexpected todo: expected %v, received %v", service.getTodo, todo)
	}
}

func TestHandler_GetTodo_InvalidID(t *testing.T) {
	handler, service := setupHandler()
	service.getTodo = Todo{Title: "some title", ID: 1, Completed: false}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos/{id}", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/todos/abc", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if service.receivedId != 0 {
		t.Fatalf("service should not have been called")
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestHandler_GetTodo_NotFound(t *testing.T) {
	handler, service := setupHandler()
	service.getTodoErr = TodoNotFound

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos/{id}", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusNotFound, res.StatusCode)
	}
}

func TestHandler_GetAll_Success(t *testing.T) {
	handler, service := setupHandler()
	service.getTodos = []Todo{
		{Title: "first", ID: 1, Completed: false},
		{Title: "second", ID: 2, Completed: false},
	}

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.GetAll(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusOK, res.StatusCode)
	}

	var todos []Todo
	err := json.NewDecoder(res.Body).Decode(&todos)
	if err != nil {
		t.Fatalf("error while reading response body: %v", err)
	}

	if len(todos) != len(service.getTodos) {
		t.Fatalf("unexpected todo list")
	}

	for i, todo := range todos {
		if todo != service.getTodos[i] {
			t.Fatalf("unexpected todo list")
		}
	}
}

func TestHandler_GetAll_Error(t *testing.T) {
	handler, service := setupHandler()
	service.getTodosErr = errors.New("some server error")

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.GetAll(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusInternalServerError, res.StatusCode)
	}
}

func TestHandler_MarkAsDone_Success(t *testing.T) {
	handler, service := setupHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /todos/{id}/done", handler.MarkAsDone)

	req := httptest.NewRequest(http.MethodPatch, "/todos/1/done", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if service.receivedId != 1 {
		t.Fatalf("unexpected todo id: expected %d, received %d", 1, service.receivedId)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusNoContent, res.StatusCode)
	}
}

func TestHandler_MarkAsDone_InvalidID(t *testing.T) {
	handler, service := setupHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /todos/{id}/done", handler.MarkAsDone)

	req := httptest.NewRequest(http.MethodPatch, "/todos/abc/done", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if service.receivedId != 0 {
		t.Fatalf("service should not have been called")
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestHandler_MarkAsDone_Error(t *testing.T) {
	handler, service := setupHandler()
	service.markAsDoneErr = errors.New("some error")

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /todos/{id}/done", handler.MarkAsDone)

	req := httptest.NewRequest(http.MethodPatch, "/todos/123/done", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusNotFound, res.StatusCode)
	}
}

func TestHandler_UpdateTitle_Success(t *testing.T) {
	handler, service := setupHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /todos/{id}", handler.UpdateTitle)

	req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewBufferString(`{"title": "new title"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if service.receivedId != 1 {
		t.Fatalf("unexpected todo id: expected %d, received %d", 1, service.receivedId)
	}

	if service.receivedTitle != "new title" {
		t.Fatalf("unexpected todo title: expected %s, received %s", "new title", service.receivedTitle)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusNoContent, res.StatusCode)
	}
}

func TestHandler_UpdateTitle_InvalidJSON(t *testing.T) {
	handler, _ := setupHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /todos/{id}", handler.UpdateTitle)

	req := httptest.NewRequest(http.MethodPut, "/todos/123", bytes.NewBufferString(`{"title": "some title"`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusBadRequest, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "invalid json") {
		t.Fatalf("unexpected error message")
	}
}

func TestHandler_UpdateTitle_InvalidID(t *testing.T) {
	handler, service := setupHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /todos/{id}", handler.UpdateTitle)

	req := httptest.NewRequest(http.MethodPut, "/todos/abc", bytes.NewBufferString(`{"title": "some title"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if service.receivedId != 0 {
		t.Fatalf("service should not have been called")
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestHandler_UpdateTitle_Error(t *testing.T) {
	handler, service := setupHandler()
	service.updateTitleErr = errors.New("some error")

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /todos/{id}", handler.UpdateTitle)

	req := httptest.NewRequest(http.MethodPut, "/todos/123", bytes.NewBufferString(`{"title": "some title"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusNotFound, res.StatusCode)
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	handler, service := setupHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /todos/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if service.receivedId != 1 {
		t.Fatalf("unexpected todo id: expected %d, received %d", 1, service.receivedId)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusNoContent, res.StatusCode)
	}
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	handler, service := setupHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /todos/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/todos/abc", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if service.receivedId != 0 {
		t.Fatalf("service should not have been called")
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestHandler_Delete_Error(t *testing.T) {
	handler, service := setupHandler()
	service.deleteErr = errors.New("some error")

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /todos/{id}", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/todos/123", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected http code: expected %d, received %d", http.StatusNotFound, res.StatusCode)
	}
}
