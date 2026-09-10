package ports

import (
	"context"

	"city-of-lies/backend/internal/domain"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, s *domain.GameSession) error
	GetSession(ctx context.Context, id string) (*domain.GameSession, error)
	UpdateSession(ctx context.Context, s *domain.GameSession) error
}

type AgentRepository interface {
	SaveAgents(ctx context.Context, sessionID string, agents []*domain.Agent) error
	ListAgents(ctx context.Context, sessionID string) ([]*domain.Agent, error)
	GetAgent(ctx context.Context, sessionID string, agentID string) (*domain.Agent, error)
}

type BeliefRepository interface {
	SaveBeliefs(ctx context.Context, beliefs []*domain.AgentBelief) error
	ListBeliefs(ctx context.Context, sessionID string, agentID string) ([]*domain.AgentBelief, error)
}

type EvidenceRepository interface {
	SaveEvidence(ctx context.Context, items []*domain.EvidenceItem) error
	ListDiscoveredEvidence(ctx context.Context, sessionID string) ([]*domain.EvidenceItem, error)
	MarkDiscovered(ctx context.Context, sessionID, evidenceID string, gameSecond int64) error
}

type EventRepository interface {
	AppendEvent(ctx context.Context, evt domain.WorldEvent) error
	ListEventsSince(ctx context.Context, sessionID string, sequence int64, limit int) ([]domain.WorldEvent, error)
}
