package services

import (
	"context"

	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/simulation"
)

type InspectService struct {
	manager *simulation.SessionManager
}

func NewInspectService(manager *simulation.SessionManager) *InspectService {
	return &InspectService{manager: manager}
}

type InspectResponse struct {
	LocationID         string                `json:"location_id"`
	DiscoveredEvidence []domain.EvidenceItem `json:"discovered_evidence"`
	Message            string                `json:"message"`
}

func (s *InspectService) Execute(ctx context.Context, sessionID, locationID string) (*InspectResponse, error) {
	engine, err := s.manager.GetEngine(sessionID)
	if err != nil {
		return nil, err
	}

	if !engine.Session.IsRunning() {
		return nil, domain.ErrSessionNotRunning
	}

	// Cost: 45 seconds of game time
	engine.Tick(45)

	var discovered []domain.EvidenceItem
	for _, ev := range engine.Evidence {
		if ev.LocationID == locationID && !ev.Discovered {
			dEv, err := engine.DiscoverEvidence(ev.ScenarioEvidenceID)
			if err == nil && dEv != nil {
				discovered = append(discovered, *dEv)
			}
		}
	}

	msg := "Không tìm thấy manh mối mới tại khu vực này."
	if len(discovered) > 0 {
		msg = "Bạn đã tìm thấy bằng chứng mới tại hiện trường!"
	}

	return &InspectResponse{
		LocationID:         locationID,
		DiscoveredEvidence: discovered,
		Message:            msg,
	}, nil
}
