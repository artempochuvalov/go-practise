package main

import (
	"context"
	"net/http"
	"os"

	"go-service/internal/config"
	"go-service/internal/database"
	"go-service/internal/logger"
	"go-service/internal/middlewares"
	"go-service/internal/todo"
)

func main() {
	log := logger.NewLogger()

	config, err := config.Load()
	if err != nil {
		log.Error("failed to get project envirenments", "eror", err.Error())
		os.Exit(1)
	}

	conn, err := database.NewDB(config.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	log.Info("application started")

	defer conn.Close(context.Background())

	r := todo.NewRepository(conn)
	service := todo.NewService(r)
	handler := todo.NewHandler(service, log)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /todos", handler.CreateTodo)
	mux.HandleFunc("GET /todos", handler.GetAll)
	mux.HandleFunc("GET /todos/{id}", handler.GetByID)
	mux.HandleFunc("DELETE /todos/{id}", handler.Delete)
	mux.HandleFunc("PUT /todos/{id}", handler.UpdateTitle)
	mux.HandleFunc("PATCH /todos/{id}/done", handler.MarkAsDone)

	recoveredMux := middlewares.RecoverMiddleware(log)(mux)
	loggedMux := middlewares.LoggingMiddleware(log)(recoveredMux)

	log.Info("starting http server", "address", config.Port)

	if err = http.ListenAndServe(config.Port, loggedMux); err != nil {
		log.Error("http server stopped", "error", err)
		os.Exit(1)
	}
}
