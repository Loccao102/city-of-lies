package handlers

import (
	"net/http"

	"city-of-lies/backend/internal/application/services"
	"city-of-lies/backend/internal/http/respond"
)

type ScenariosHandler struct {
	listScenarios *services.ListScenariosService
}

func NewScenariosHandler(listScenarios *services.ListScenariosService) *ScenariosHandler {
	return &ScenariosHandler{
		listScenarios: listScenarios,
	}
}

func (h *ScenariosHandler) ListScenarios(w http.ResponseWriter, r *http.Request) {
	scenarios, err := h.listScenarios.Execute(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, scenarios)
}
