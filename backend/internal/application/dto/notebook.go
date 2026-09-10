package dto

import (
	"city-of-lies/backend/internal/domain"
)

type NotebookDTO struct {
	DiscoveredEvidence []domain.EvidenceItem `json:"discovered_evidence"`
	KnownClaims        []ClaimSummaryDTO     `json:"known_claims"`
	InterviewedAgents  []string              `json:"interviewed_agents"`
	TimelineEvents     []TimelineEventDTO    `json:"timeline_events"`
}

type ClaimSummaryDTO struct {
	ID          string `json:"id"`
	DisplayText string `json:"display_text"`
	Status      string `json:"status"` // "Unverified", "Supported", "Contradicted"
}

type TimelineEventDTO struct {
	GameSecond  int64  `json:"game_second"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
}
