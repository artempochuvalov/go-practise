package todo

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type Handler struct {
	service TodoService
	logger  *slog.Logger
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

func NewHandler(service TodoService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (handler *Handler) CreateTodo(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		handler.logger.Warn(
			"invalid request body",
			"error", err,
		)
		return
	}

	id, err := handler.service.CreateTodo(r.Context(), req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		handler.logger.Warn(
			"failed to create todo",
			"error", err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := CreateTodoResponse{
		ID: id,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		handler.logger.Error(
			"unexpected response encoding",
			"error", err,
		)
		return
	}
}

func (handler *Handler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	todos, err := handler.service.GetAllTodos(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		handler.logger.Warn(
			"failed to get todos",
			"error", err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		handler.logger.Error(
			"unexpected response encoding",
			"error", err,
		)
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
		handler.logger.Warn(
			"unexpected todo's id",
			"error", err,
		)
		return
	}

	todo, err := handler.service.GetTodo(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		handler.logger.Warn(
			"failed to get todo",
			"error", err,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		handler.logger.Error(
			"unexpected response encoding",
			"error", err,
		)
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
		handler.logger.Warn(
			"unexpected todo's id",
			"error", err,
		)
		return
	}

	err = handler.service.DeleteTodo(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		handler.logger.Warn(
			"failed to delete todo",
			"error", err,
		)
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
		handler.logger.Warn(
			"unexpected todo's id",
			"error", err,
		)
		return
	}

	err = handler.service.MarkAsDone(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		handler.logger.Warn(
			"failed to mark todo as done",
			"error", err,
		)
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
		handler.logger.Warn(
			"unexpected todo's id",
			"error", err,
		)
		return
	}

	var req UpdateTodoTitleRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		handler.logger.Warn(
			"unexpected request body",
			"error", err,
		)
		return
	}

	err = handler.service.UpdateTitle(r.Context(), id, req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		handler.logger.Warn(
			"failed to update todo's title",
			"error", err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
