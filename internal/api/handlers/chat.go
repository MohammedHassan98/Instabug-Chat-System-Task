package handlers

import (
	"chat-system/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ChatHandler struct {
	service *service.ChatService
}

func NewChatHandler(service *service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	chat, err := h.service.CreateChat(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ChatNumber int `json:"chat_number"`
	}{
		ChatNumber: chat.ChatNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
