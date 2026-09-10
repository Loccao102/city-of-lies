package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"city-of-lies/backend/internal/application/dto"
	"city-of-lies/backend/internal/application/services"
	"city-of-lies/backend/internal/http/respond"
)

type InvestigationHandler struct {
	inspect    *services.InspectService
	notebook   *services.NotebookService
	correction *services.PublishCorrectionService
	truth      *services.SubmitTruthService
}

func NewInvestigationHandler(
	inspect *services.InspectService,
	notebook *services.NotebookService,
	correction *services.PublishCorrectionService,
	truth *services.SubmitTruthService,
) *InvestigationHandler {
	return &InvestigationHandler{
		inspect:    inspect,
		notebook:   notebook,
		correction: correction,
		truth:      truth,
	}
}

func (h *InvestigationHandler) InspectLocation(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	locationID := chi.URLParam(r, "locationId")

	resp, err := h.inspect.Execute(r.Context(), sessionID, locationID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, resp)
}

func (h *InvestigationHandler) GetNotebook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")

	notebook, err := h.notebook.Execute(r.Context(), sessionID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, notebook)
}

func (h *InvestigationHandler) PublishCorrection(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")

	var req dto.PublishCorrectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pub, err := h.correction.Execute(r.Context(), sessionID, req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.JSON(w, http.StatusCreated, pub)
}

func (h *InvestigationHandler) SubmitTruth(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")

	var req dto.SubmitTruthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.truth.Execute(r.Context(), sessionID, req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, res)
}
