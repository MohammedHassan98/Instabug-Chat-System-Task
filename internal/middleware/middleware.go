package middleware

import (
	"chat-system/internal/errors"
	"chat-system/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				appErr := errors.ErrInternalServer(fmt.Errorf("%v", err))
				logger.Error(r.Context(), "Panic recovered", fmt.Errorf("%v", err))
				respondWithError(w, appErr)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()

		// Create a custom response writer to capture status code
		rw := &ResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Add request ID to context
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Process request
		next.ServeHTTP(rw, r)

		// Log request details
		logger.Info(ctx, "Request completed",
			map[string]interface{}{
				"request_id":  requestID,
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      rw.status,
				"size":        rw.size,
				"duration_ms": time.Since(start).Milliseconds(),
				"user_agent":  r.UserAgent(),
				"remote_addr": r.RemoteAddr,
			})
	})
}

func respondWithError(w http.ResponseWriter, err *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	response := struct {
		Error string `json:"error"`
	}{
		Error: err.Message,
	}

	json.NewEncoder(w).Encode(response)
}
