package domain

// AgentBelief represents an agent's confidence level and provenance for a specific claim.
type AgentBelief struct {
	SessionID           string  `json:"session_id"`
	AgentID             string  `json:"agent_id"`
	ClaimID             string  `json:"claim_id"`
	Confidence          float64 `json:"confidence"` // [0.0, 1.0]
	FirstSourceAgentID  *string `json:"first_source_agent_id,omitempty"`
	LastSourceAgentID   *string `json:"last_source_agent_id,omitempty"`
	RootSourceAgentID   *string `json:"root_source_agent_id,omitempty"`
	EvidenceResistance  float64 `json:"evidence_resistance"`
	UpdatedAtGameSecond int64   `json:"updated_at_game_second"`
	Revision            int64   `json:"revision"`
}

// ClampConfidence ensures confidence is bounded strictly to [0.0, 1.0].
func ClampConfidence(c float64) float64 {
	if c < 0.0 {
		return 0.0
	}
	if c > 1.0 {
		return 1.0
	}
	return c
}
