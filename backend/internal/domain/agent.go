package domain

import (
	"fmt"
)

// Agent represents an autonomous non-player character in a game session.
type Agent struct {
	ID                  string  `json:"id"`
	SessionID           string  `json:"session_id"`
	ScenarioAgentID     string  `json:"scenario_agent_id"`
	Name                string  `json:"name"`
	Role                string  `json:"role"`
	CurrentLocationID   string  `json:"current_location_id"`
	Influence           float64 `json:"influence"`
	Credibility         float64 `json:"credibility"`
	Skepticism          float64 `json:"skepticism"`
	Sociality           float64 `json:"sociality"`
	DeceptionTendency   float64 `json:"deception_tendency"`
	Receptiveness       float64 `json:"receptiveness"`
	BusyUntilGameSecond int64   `json:"busy_until_game_second"`
	Revision            int64   `json:"revision"`
}

// IsBusy returns true if the agent is engaged in an activity at the given game second.
func (a *Agent) IsBusy(currentGameSecond int64) bool {
	return a.BusyUntilGameSecond > currentGameSecond
}

// ValidateTraits ensures all numeric attributes are clamped within [0, 1].
func (a *Agent) ValidateTraits() error {
	traits := map[string]float64{
		"influence":          a.Influence,
		"credibility":        a.Credibility,
		"skepticism":         a.Skepticism,
		"sociality":          a.Sociality,
		"deception_tendency": a.DeceptionTendency,
		"receptiveness":      a.Receptiveness,
	}
	for name, val := range traits {
		if val < 0.0 || val > 1.0 {
			return fmt.Errorf("agent %s: trait %s is %.2f, must be within [0, 1]", a.ScenarioAgentID, name, val)
		}
	}
	return nil
}
