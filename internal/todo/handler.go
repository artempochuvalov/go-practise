package todo

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type CreateTodoResponse struct {
	ID int `json:"id"`
}

type UpdateTodoTitleRequest struct {
	Title string `json:"title"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) CreateTodo(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	id, err := handler.service.CreateTodo(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := CreateTodoResponse{
		ID: id,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	todos, err := handler.service.GetAllTodos()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, err := handler.service.GetTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = handler.service.DeleteTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) MarkAsDone(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = handler.service.MarkAsDone(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) UpdateTitle(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateTodoTitleRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = handler.service.UpdateTitle(id, req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
