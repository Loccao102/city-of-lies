package services

import (
	"context"

	"city-of-lies/backend/internal/application/dto"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/simulation"
)

type SubmitTruthService struct {
	manager *simulation.SessionManager
}

func NewSubmitTruthService(manager *simulation.SessionManager) *SubmitTruthService {
	return &SubmitTruthService{manager: manager}
}

func (s *SubmitTruthService) Execute(ctx context.Context, sessionID string, req dto.SubmitTruthRequest) (*domain.TruthSubmissionResult, error) {
	engine, err := s.manager.GetEngine(sessionID)
	if err != nil {
		return nil, err
	}

	// Cost: 15 seconds of game time
	engine.Tick(15)

	sub := domain.GroundTruthSubmission{
		EventType:       req.EventType,
		CauseCategory:   req.CauseCategory,
		MajorExplosion:  req.MajorExplosion,
		Fatalities:      req.Fatalities,
		ChemicalRelease: req.ChemicalRelease,
		ManagementIssue: req.ManagementIssue,
	}

	result, err := engine.SubmitGroundTruth(sub)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
