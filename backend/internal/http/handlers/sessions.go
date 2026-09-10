package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"city-of-lies/backend/internal/application/dto"
	"city-of-lies/backend/internal/application/services"
	"city-of-lies/backend/internal/http/respond"
)

type SessionsHandler struct {
	startSession *services.StartSessionService
	getState     *services.GetStateService
}

func NewSessionsHandler(startSession *services.StartSessionService, getState *services.GetStateService) *SessionsHandler {
	return &SessionsHandler{
		startSession: startSession,
		getState:     getState,
	}
}

func (h *SessionsHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req dto.StartSessionRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	res, err := h.startSession.Execute(r.Context(), req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusCreated, res)
}

func (h *SessionsHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		respond.Error(w, http.StatusBadRequest, "session id required")
		return
	}

	snapshot, err := h.getState.Execute(r.Context(), sessionID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, snapshot)
}
