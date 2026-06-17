package main

import (
	"context"
	"net/http"
	"os"

	"go-service/internal/database"
	"go-service/internal/todo"
)

func main() {
	conn, err := database.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	r := todo.NewRepository(conn)
	service := todo.NewService(r)
	handler := todo.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /todos", handler.CreateTodo)
	mux.HandleFunc("GET /todos", handler.GetAll)
	mux.HandleFunc("GET /todos/{id}", handler.GetByID)
	mux.HandleFunc("DELETE /todos/{id}", handler.Delete)
	mux.HandleFunc("PUT /todos/{id}", handler.UpdateTitle)
	mux.HandleFunc("PATCH /todos/{id}/done", handler.MarkAsDone)

	if err = http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
