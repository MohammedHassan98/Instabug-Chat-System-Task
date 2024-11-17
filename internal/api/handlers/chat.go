// Package handlers provides HTTP request handlers
package handlers

import (
	"chat-system/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// @title Chat API
// @version 1.0
// @description Chat handler manages chat operations

type ChatHandler struct {
	service *service.ChatService
}

func NewChatHandler(service *service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

// @Summary Create a new chat
// @Description Creates a new chat for the given application token
// @Tags Chats
// @Accept json
// @Produce json
// @Param token path string true "Application Token"
// @Success 200 {object} handlers.ChatResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /applications/{token}/chats [post]
func (h *ChatHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	chat, err := h.service.CreateChat(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ChatNumber int `json:"Chat Number"`
	}{
		ChatNumber: chat.ChatNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}