package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"city-of-lies/backend/internal/application/dto"
	"city-of-lies/backend/internal/application/services"
	"city-of-lies/backend/internal/game/simulation"
	"city-of-lies/backend/internal/http/respond"
)

type AgentsHandler struct {
	manager   *simulation.SessionManager
	interview *services.InterviewService
}

func NewAgentsHandler(manager *simulation.SessionManager, interview *services.InterviewService) *AgentsHandler {
	return &AgentsHandler{
		manager:   manager,
		interview: interview,
	}
}

func (h *AgentsHandler) ListAgents(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	engine, err := h.manager.GetEngine(sessionID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}

	snapshot := engine.BuildSnapshot()
	respond.JSON(w, http.StatusOK, snapshot.Agents)
}

func (h *AgentsHandler) GetAgent(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	agentID := chi.URLParam(r, "agentId")

	engine, err := h.manager.GetEngine(sessionID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}

	for _, a := range engine.Agents {
		if a.ID == agentID || a.ScenarioAgentID == agentID {
			status := "available"
			if a.BusyUntilGameSecond > engine.Session.GameSecond {
				status = "busy"
			}
			detail := dto.AgentDetailDTO{
				AgentSummaryDTO: dto.AgentSummaryDTO{
					ID:              a.ID,
					ScenarioAgentID: a.ScenarioAgentID,
					Name:            a.Name,
					Role:            a.Role,
					LocationID:      a.CurrentLocationID,
					PublicStatus:    status,
					Interviewed:     false,
				},
				Personality: "",
			}
			respond.JSON(w, http.StatusOK, detail)
			return
		}
	}

	respond.Error(w, http.StatusNotFound, "agent not found")
}

func (h *AgentsHandler) Interview(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	agentID := chi.URLParam(r, "agentId")

	var req dto.InterviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.interview.Execute(r.Context(), sessionID, agentID, req.Message)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, resp)
}
