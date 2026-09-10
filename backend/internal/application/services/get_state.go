package services

import (
	"context"

	"city-of-lies/backend/internal/game/simulation"
)

type GetStateService struct {
	manager *simulation.SessionManager
}

func NewGetStateService(manager *simulation.SessionManager) *GetStateService {
	return &GetStateService{manager: manager}
}

func (s *GetStateService) Execute(ctx context.Context, sessionID string) (*simulation.SessionSnapshot, error) {
	engine, err := s.manager.GetEngine(sessionID)
	if err != nil {
		return nil, err
	}

	snapshot := engine.BuildSnapshot()
	return &snapshot, nil
}
