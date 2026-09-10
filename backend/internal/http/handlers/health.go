package handlers

import (
	"net/http"

	"city-of-lies/backend/internal/http/respond"
)

type HealthHandler struct {
	isReady func() bool
}

func NewHealthHandler(isReady func() bool) *HealthHandler {
	return &HealthHandler{isReady: isReady}
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusOK, map[string]string{
		"status": "alive",
	})
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.isReady != nil && !h.isReady() {
		respond.Error(w, http.StatusServiceUnavailable, "service not ready")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}
