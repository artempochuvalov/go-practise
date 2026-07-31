package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-service/internal/config"
	"go-service/internal/database"
	"go-service/internal/logger"
	"go-service/internal/middlewares"
	"go-service/internal/todo"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to get project environments", "error", err)
		os.Exit(1)
	}

	conn, err := database.NewDB(cfg.DatabaseURL)
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

	log.Info("starting http server", "address", cfg.Port)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: loggedMux,
	}

	serverErr := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
	}()

	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-signalContext.Done():
		log.Info("shutdown signal received")

		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Info("shutting down server")
		shutdownStart := time.Now()

		if err := server.Shutdown(shutdownContext); err != nil {
			log.Error("unable to shutdown", "error", err)

			if closeErr := server.Close(); closeErr != nil {
				log.Error("unable to force close server", "error", closeErr)
			}
		}

		shutdownDuration := time.Since(shutdownStart)
		log.Info("server shutdown completed", "duration", shutdownDuration)

	case err := <-serverErr:
		log.Error("server failed", "error", err)
	}
}
