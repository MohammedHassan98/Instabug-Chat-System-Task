package handlers

import (
	"chat-system/internal/pkg/httputil"
	"chat-system/internal/pkg/validation"
	"chat-system/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MessageHandler struct {
	service *service.MessageService
}

type createMessageRequest struct {
	Body string `json:"body" validate:"required,min=1"`
}

func NewMessageHandler(service *service.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createMessageRequest

	vars := mux.Vars(r)
	chatNumber, err := strconv.Atoi(vars["chatNumber"])

	if err != nil {
		http.Error(w, "Invalid chat number", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if Body is empty after decoding
	if req.Body == "" {
		http.Error(w, "Body cannot be empty", http.StatusBadRequest)
		return
	}

	// Validate the request
	if errors := validation.ValidateStruct(req); len(errors) > 0 {
		httputil.WriteValidationErrors(w, errors)
		return
	}

	message, err := h.service.CreateMessage(r.Context(), uint(chatNumber), req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		MessageNumber int `json:"Message Number"`
	}{
		MessageNumber: message.MessageNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
