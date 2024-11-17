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

// @title Message API
// @version 1.0
// @description Message handler manages message operations

type MessageHandler struct {
	service *service.MessageService
}

type createMessageRequest struct {
	Body string `json:"body" validate:"required,min=1"`
}

func NewMessageHandler(service *service.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}


// @Summary Create a new message
// @Description Creates a new message in a chat
// @Tags Messages
// @Accept json
// @Produce json
// @Param chatNumber path int true "Chat Number"
// @Param message body createMessageRequest true "Message creation request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /chats/{chatNumber}/messages [post]
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

// @Summary Get messages
// @Description Retrieves all messages for a chat
// @Tags Messages
// @Accept json
// @Produce json
// @Param token path string true "Application Token"
// @Param chatNumber path int true "Chat Number"
// @Success 200 {array} MessageListResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /applications/{token}/chats/{chatNumber}/messages [get]
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"] // Get the application token from the URL
	chatNumber, err := strconv.Atoi(vars["chatNumber"])
	if err != nil {
		http.Error(w, "Invalid chat number", http.StatusBadRequest)
		return
	}

	messages, err := h.service.GetMessagesByChatNumberAndToken(r.Context(), token, uint(chatNumber))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]struct {
		MessageNumber int    `json:"Message Number"`
		Body          string `json:"body"`
	}, len(messages))

	for i, message := range messages {
		response[i] = struct {
			MessageNumber int    `json:"Message Number"`
			Body          string `json:"body"`
		}{
			MessageNumber: message.MessageNumber,
			Body:          message.Body,
		}
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

// @Summary Search messages
// @Description Search messages in a chat
// @Tags Messages
// @Accept json
// @Produce json
// @Param chatNumber path int true "Chat Number"
// @Param q query string true "Search query"
// @Success 200 {object} MessageListResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /chats/{chatNumber}/messages/search [get]
func (h *MessageHandler) Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatNumber, err := strconv.Atoi(vars["chatNumber"])
	if err != nil {
		http.Error(w, "Invalid chat number", http.StatusBadRequest)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	messages, err := h.service.SearchMessages(r.Context(), uint(chatNumber), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Messages []struct {
			MessageNumber int    `json:"Message Number"`
			Body          string `json:"Body"`
		} `json:"messages"`
	}{
		Messages: make([]struct {
			MessageNumber int    `json:"Message Number"`
			Body          string `json:"Body"`
		}, len(messages)),
	}

	for i, msg := range messages {
		response.Messages[i].MessageNumber = msg.MessageNumber
		response.Messages[i].Body = msg.Body
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update the Create handler annotation
// @Success 200 {object} MessageResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse

// Update the GetMessages handler annotation
// @Success 200 {array} MessageListResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse

// Update the Search handler annotation
// @Success 200 {object} MessageListResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
