package simulation

import (
	"city-of-lies/backend/internal/domain"
)

type AgentPublicView struct {
	ID              string `json:"id"`
	ScenarioAgentID string `json:"scenario_agent_id"`
	Name            string `json:"name"`
	Role            string `json:"role"`
	LocationID      string `json:"location_id"`
	PublicStatus    string `json:"public_status"`
	Interviewed     bool   `json:"interviewed"`
}

type SessionSnapshot struct {
	ID                  string               `json:"id"`
	ScenarioID          string               `json:"scenario_id"`
	Status              domain.SessionStatus `json:"status"`
	GameSecond          int64                `json:"game_second"`
	FalseNarrativeRatio float64              `json:"primary_false_narrative_ratio"`
	EvidenceStrength    float64              `json:"evidence_strength"`
	PlayerCredibility   float64              `json:"player_credibility"`
	SimulationSeed      int64                `json:"simulation_seed"`
	Revision            int64                `json:"revision"`
	Sequence            int64                `json:"sequence"`
	Agents              []AgentPublicView    `json:"agents"`
	DiscoveredEvidence  []domain.EvidenceItem `json:"discovered_evidence"`
	RecentEvents        []domain.WorldEvent  `json:"recent_events"`
}

// BuildSnapshot creates a player-safe state projection with zero leak of Ground Truth or server-only tags.
func (e *SimulationEngine) BuildSnapshot() SessionSnapshot {
	e.mu.RLock()
	defer e.mu.RUnlock()

	agentViews := make([]AgentPublicView, 0, len(e.Agents))
	for _, a := range e.Agents {
		status := "available"
		if a.BusyUntilGameSecond > e.Session.GameSecond {
			status = "busy"
		}
		agentViews = append(agentViews, AgentPublicView{
			ID:              a.ID,
			ScenarioAgentID: a.ScenarioAgentID,
			Name:            a.Name,
			Role:            a.Role,
			LocationID:      a.CurrentLocationID,
			PublicStatus:    status,
			Interviewed:     false,
		})
	}

	discoveredEv := make([]domain.EvidenceItem, 0)
	for _, ev := range e.Evidence {
		if ev.Discovered {
			discoveredEv = append(discoveredEv, *ev)
		}
	}

	// Last 50 events
	eventCount := len(e.EventLog)
	start := 0
	if eventCount > 50 {
		start = eventCount - 50
	}
	recentEvents := make([]domain.WorldEvent, eventCount-start)
	copy(recentEvents, e.EventLog[start:])

	return SessionSnapshot{
		ID:                  e.Session.ID,
		ScenarioID:          e.Session.ScenarioID,
		Status:              e.Session.Status,
		GameSecond:          e.Session.GameSecond,
		FalseNarrativeRatio: e.Session.FalseNarrativeRatio,
		EvidenceStrength:    e.Session.EvidenceStrength,
		PlayerCredibility:   e.Session.PlayerCredibility,
		SimulationSeed:      e.Session.SimulationSeed,
		Revision:            e.Session.Revision,
		Sequence:            e.LatestSequence,
		Agents:              agentViews,
		DiscoveredEvidence:  discoveredEv,
		RecentEvents:        recentEvents,
	}
}
