package domain

import (
	"time"
)

type SessionStatus string

const (
	SessionStatusCreated         SessionStatus = "created"
	SessionStatusRunning         SessionStatus = "running"
	SessionStatusWon             SessionStatus = "won"
	SessionStatusLostFalseBelief SessionStatus = "lost_false_belief"
	SessionStatusLostTimeout     SessionStatus = "lost_timeout"
	SessionStatusAborted         SessionStatus = "aborted"
)

// GameSession represents an active or concluded game session aggregate.
type GameSession struct {
	ID                  string        `json:"id"`
	ScenarioID          string        `json:"scenario_id"`
	Status              SessionStatus `json:"status"`
	SimulationSeed      int64         `json:"simulation_seed"`
	GameSecond          int64         `json:"game_second"`
	PlayerCredibility   float64       `json:"player_credibility"`
	FalseNarrativeRatio float64       `json:"false_narrative_ratio"`
	EvidenceStrength    float64       `json:"evidence_strength"`
	Revision            int64         `json:"revision"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	EndedAt             *time.Time    `json:"ended_at,omitempty"`
}

func (s *GameSession) IsRunning() bool {
	return s.Status == SessionStatusRunning
}

func (s *GameSession) IsTerminal() bool {
	return s.Status == SessionStatusWon ||
		s.Status == SessionStatusLostFalseBelief ||
		s.Status == SessionStatusLostTimeout ||
		s.Status == SessionStatusAborted
}
