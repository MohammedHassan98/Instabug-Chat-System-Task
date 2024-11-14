package handlers

import (
	"chat-system/internal/pkg/httputil"
	"chat-system/internal/pkg/validation"
	"chat-system/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ApplicationHandler struct {
	service *service.ApplicationService
}

func NewApplicationHandler(service *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

type createApplicationRequest struct {
	Name string `json:"name" validate:"required,app_name"`
}

func (h *ApplicationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createApplicationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the request
	if errors := validation.ValidateStruct(req); len(errors) > 0 {
		httputil.WriteValidationErrors(w, errors)
		return
	}

	app, err := h.service.CreateApplication(r.Context(), req.Name)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Token string `json:"token"`
		Name  string `json:"name"`
	}{
		Token: app.Token,
		Name:  app.Name,
	}

	httputil.WriteJSON(w, http.StatusCreated, response)
}

func (h *ApplicationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	apps, err := h.service.GetAllApplications(r.Context())

	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]struct {
		Token string `json:"token"`
		Name  string `json:"name"`
	}, len(apps))

	for i, app := range apps {
		response[i] = struct {
			Token string `json:"token"`
			Name  string `json:"name"`
		}{
			Token: app.Token,
			Name:  app.Name,
		}
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *ApplicationHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	var req createApplicationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the request
	if errors := validation.ValidateStruct(req); len(errors) > 0 {
		httputil.WriteValidationErrors(w, errors)
		return
	}

	app, err := h.service.UpdateApplication(r.Context(), token, req.Name)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Token string `json:"Token"`
		Name  string `json:"Name"`
	}{
		Token: app.Token,
		Name:  app.Name,
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}
