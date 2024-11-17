// cmd/server/main.go
package main

import (
	"chat-system/internal/api/handlers"
	"chat-system/internal/db"
	"chat-system/internal/logger"
	"chat-system/internal/middleware"
	"chat-system/internal/queue"
	"chat-system/internal/service"
	"chat-system/internal/worker"
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

	// Initialize Queue
	messageQueue := queue.NewMessageQueue(db.Redis)

	// Initialize Services
	appService := service.NewApplicationService(db.GormDB)
	chatService := service.NewChatService(db.GormDB, db.Redis, messageQueue)
	messageService := service.NewMessageService(db.GormDB, db.Redis, messageQueue, db.ES)

	// Initialize Handlers
	appHandler := handlers.NewApplicationHandler(appService)
	chatHandler := handlers.NewChatHandler(chatService)
	messageHandler := handlers.NewMessageHandler(messageService)

	// Initialize Worker
	worker := worker.NewWorker(messageQueue, db.ES)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	worker.Start(ctx)

	// Initialize middlewares
	rateLimiter := middleware.NewRateLimiter(100, 200)

	router := mux.NewRouter()
	router.Use(middleware.RequestLogger)
	router.Use(middleware.ErrorHandler)
	router.Use(rateLimiter.RateLimit)

	// Application routes
	router.HandleFunc("/applications", appHandler.GetAll).Methods("GET")
	router.HandleFunc("/applications/{token}/chats", appHandler.GetChats).Methods("GET")
	router.HandleFunc("/applications", appHandler.Create).Methods("POST")
	router.HandleFunc("/applications/{token}", appHandler.Update).Methods("PUT")

	// Chat routes
	router.HandleFunc("/chats/{token}", chatHandler.Create).Methods("POST")

	// Message routes
	router.HandleFunc("/messages/{chatNumber}", messageHandler.Create).Methods("POST")
	router.HandleFunc("/applications/{token}/chats/{chatNumber}/messages", messageHandler.GetMessages).Methods("GET")
	router.HandleFunc("/chats/{chatNumber}/messages/search", messageHandler.Search).Methods("GET")

	// Create server with timeouts
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
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