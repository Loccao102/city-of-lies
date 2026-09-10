package domain

// Relationship models the directed epistemic trust from one agent toward another.
type Relationship struct {
	SessionID   string  `json:"session_id"`
	FromAgentID string  `json:"from_agent_id"`
	ToAgentID   string  `json:"to_agent_id"`
	Trust       float64 `json:"trust"` // [0.0, 1.0]
	Revision    int64   `json:"revision"`
}

// ClampTrust bounds trust strictly within [0.0, 1.0].
func ClampTrust(t float64) float64 {
	if t < 0.0 {
		return 0.0
	}
	if t > 1.0 {
		return 1.0
	}
	return t
}
