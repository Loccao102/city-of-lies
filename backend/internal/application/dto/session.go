package dto

import (
	"city-of-lies/backend/internal/domain"
)

type SessionResponse struct {
	ID                  string            `json:"id"`
	ScenarioID          string            `json:"scenario_id"`
	ScenarioTitle       string            `json:"scenario_title,omitempty"`
	ScenarioDescription string            `json:"scenario_description,omitempty"`
	Status              string            `json:"status"`
	GameSecond          int64             `json:"game_second"`
	FalseNarrativeRatio float64           `json:"primary_false_narrative_ratio"`
	EvidenceStrength    float64           `json:"evidence_strength"`
	PlayerCredibility   float64           `json:"player_credibility"`
	SimulationSeed      int64             `json:"simulation_seed"`
	Revision            int64             `json:"revision"`
	Sequence            int64             `json:"sequence"`
	Locations           []domain.Location `json:"locations,omitempty"`
}

type ScenarioSummaryDTO struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Version       int      `json:"version"`
	Description   string   `json:"description"`
	AgentCount    int      `json:"agent_count"`
	LocationCount int      `json:"location_count"`
	EvidenceCount int      `json:"evidence_count"`
	DefeatRatio   float64  `json:"defeat_ratio"`
	PrimaryClaims []string `json:"primary_claims"`
}

