package services

import (
	"context"

	"city-of-lies/backend/internal/application/dto"
	"city-of-lies/backend/internal/game/simulation"
)

type StartSessionService struct {
	manager *simulation.SessionManager
}

func NewStartSessionService(manager *simulation.SessionManager) *StartSessionService {
	return &StartSessionService{manager: manager}
}

func (s *StartSessionService) Execute(ctx context.Context, req dto.StartSessionRequest) (*dto.SessionResponse, error) {
	var seed int64
	if req.Seed != nil {
		seed = *req.Seed
	}

	engine, err := s.manager.CreateSession(seed)
	if err != nil {
		return nil, err
	}

	snapshot := engine.BuildSnapshot()
	return &dto.SessionResponse{
		ID:                  snapshot.ID,
		ScenarioID:          snapshot.ScenarioID,
		Status:              string(snapshot.Status),
		GameSecond:          snapshot.GameSecond,
		FalseNarrativeRatio: snapshot.FalseNarrativeRatio,
		EvidenceStrength:    snapshot.EvidenceStrength,
		PlayerCredibility:   snapshot.PlayerCredibility,
		SimulationSeed:      snapshot.SimulationSeed,
		Revision:            snapshot.Revision,
		Sequence:            snapshot.Sequence,
	}, nil
}
