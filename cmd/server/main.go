// cmd/server/main.go
package main

import (
	"chat-system/internal/api/handlers"
	"chat-system/internal/db"
	"chat-system/internal/logger"
	"chat-system/internal/middleware"
	"chat-system/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Database and Services
	db.Connect()

	// Initialize Services
	appService := service.NewApplicationService(db.GormDB)

	// Initialize Handlers
	appHandler := handlers.NewApplicationHandler(appService)

	// Set up router with middleware
	router := mux.NewRouter()
	router.Use(middleware.RequestLogger)
	router.Use(middleware.ErrorHandler)

	// Application routes
	router.HandleFunc("/applications", appHandler.GetAll).Methods("GET")
	router.HandleFunc("/applications", appHandler.Create).Methods("POST")
	router.HandleFunc("/applications/{token}", appHandler.Update).Methods("PUT")

	// Create server with timeouts
	srv := &http.Server{
		Addr:         os.Getenv("SERVER_ADDRESS"),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(context.Background(), "Starting server", map[string]string{"address": srv.Addr})
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "Server is shutting down...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(context.Background(), "Server forced to shutdown", err)
	}

	logger.Info(context.Background(), "Server stopped")
}
