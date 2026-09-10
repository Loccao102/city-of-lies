package services

import (
	"context"

	"city-of-lies/backend/internal/application/dto"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/simulation"
)

type NotebookService struct {
	manager *simulation.SessionManager
}

func NewNotebookService(manager *simulation.SessionManager) *NotebookService {
	return &NotebookService{manager: manager}
}

func (s *NotebookService) Execute(ctx context.Context, sessionID string) (*dto.NotebookDTO, error) {
	engine, err := s.manager.GetEngine(sessionID)
	if err != nil {
		return nil, err
	}

	discoveredEv := make([]domain.EvidenceItem, 0)
	for _, ev := range engine.Evidence {
		if ev.Discovered {
			discoveredEv = append(discoveredEv, *ev)
		}
	}

	knownClaims := make([]dto.ClaimSummaryDTO, 0)
	for _, c := range engine.Claims {
		status := "Unverified"
		for _, ev := range discoveredEv {
			for _, sc := range engine.ScenarioBundle.Evidence {
				if sc.ID == ev.ScenarioEvidenceID {
					for _, sup := range sc.Supports {
						if sup == c.ScenarioClaimID {
							status = "Supported"
						}
					}
					for _, con := range sc.Contradicts {
						if con == c.ScenarioClaimID {
							status = "Contradicted"
						}
					}
				}
			}
		}
		knownClaims = append(knownClaims, dto.ClaimSummaryDTO{
			ID:          c.ScenarioClaimID,
			DisplayText: c.DisplayText,
			Status:      status,
		})
	}

	timeline := make([]dto.TimelineEventDTO, 0)
	for _, pe := range engine.ScenarioBundle.PublicEvents {
		if pe.GameSecond <= engine.Session.GameSecond {
			timeline = append(timeline, dto.TimelineEventDTO{
				GameSecond:  pe.GameSecond,
				Headline:    pe.Headline,
				Description: pe.Description,
			})
		}
	}

	return &dto.NotebookDTO{
		DiscoveredEvidence: discoveredEv,
		KnownClaims:        knownClaims,
		InterviewedAgents:  []string{},
		TimelineEvents:     timeline,
	}, nil
}
