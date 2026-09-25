package services

import (
	"context"
	"sort"

	"city-of-lies/backend/internal/application/dto"
	"city-of-lies/backend/internal/game/simulation"
)

type ListScenariosService struct {
	manager *simulation.SessionManager
}

func NewListScenariosService(manager *simulation.SessionManager) *ListScenariosService {
	return &ListScenariosService{manager: manager}
}

func (s *ListScenariosService) Execute(ctx context.Context) ([]dto.ScenarioSummaryDTO, error) {
	bundles := s.manager.GetBundles()
	result := make([]dto.ScenarioSummaryDTO, 0, len(bundles))

	for id, b := range bundles {
		primaryClaims := make([]string, 0, len(b.Metadata.PrimaryFalseNarrative.Claims))
		for _, c := range b.Metadata.PrimaryFalseNarrative.Claims {
			primaryClaims = append(primaryClaims, c.ClaimID)
		}

		result = append(result, dto.ScenarioSummaryDTO{
			ID:            id,
			Title:         b.Metadata.Title,
			Version:       b.Metadata.Version,
			Description:   b.Metadata.Description,
			AgentCount:    len(b.Agents),
			LocationCount: len(b.Locations),
			EvidenceCount: len(b.Evidence),
			DefeatRatio:   b.Metadata.PrimaryFalseNarrative.DefeatRatio,
			PrimaryClaims: primaryClaims,
		})
	}

	// Sort alphabetically by ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}
