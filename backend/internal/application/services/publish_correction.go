package services

import (
	"context"

	"city-of-lies/backend/internal/application/dto"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/simulation"
)

type PublishCorrectionService struct {
	manager *simulation.SessionManager
}

func NewPublishCorrectionService(manager *simulation.SessionManager) *PublishCorrectionService {
	return &PublishCorrectionService{manager: manager}
}

func (s *PublishCorrectionService) Execute(ctx context.Context, sessionID string, req dto.PublishCorrectionRequest) (*domain.Publication, error) {
	engine, err := s.manager.GetEngine(sessionID)
	if err != nil {
		return nil, err
	}

	// Cost: 45 seconds of game time
	engine.Tick(45)

	return engine.PublishCorrection(req.ChallengedClaimID, req.EvidenceIDs, req.Message)
}
