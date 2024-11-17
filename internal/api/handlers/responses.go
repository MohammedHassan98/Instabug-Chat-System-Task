package handlers

// Error response structures
type ErrorResponse struct {
    Error string `json:"error"`
}

type ValidationErrorResponse struct {
    Errors map[string]string `json:"errors"`
}

// Message response structures
type MessageResponse struct {
    MessageNumber int    `json:"Message Number"`
    Body          string `json:"body"`
}

type MessageListResponse []MessageResponse

type MessageSearchResponse struct {
    Messages []MessageResponse `json:"messages"`
}

// Chat response structures
type ChatResponse struct {
    ChatNumber int `json:"Chat Number"`
}

type ChatWithMessagesResponse struct {
    ChatNumber int `json:"Chat Number"`
    Messages   int `json:"Messages"`
}

type ChatListResponse []ChatWithMessagesResponse

// ApplicationChatResponse represents the chat information for an application
type ApplicationChatResponse struct {
    ChatNumber int `json:"Chat Number"`
    Messages   int `json:"Messages"`
}

// ApplicationListResponse represents a list of applications
type ApplicationListResponse []ApplicationResponse
