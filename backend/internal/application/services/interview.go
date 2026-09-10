package services

import (
	"context"

	"city-of-lies/backend/internal/application/ports"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/dialogue"
	"city-of-lies/backend/internal/game/simulation"
)

type InterviewService struct {
	manager  *simulation.SessionManager
	provider ports.DialogueProvider
}

func NewInterviewService(manager *simulation.SessionManager, provider ports.DialogueProvider) *InterviewService {
	return &InterviewService{
		manager:  manager,
		provider: provider,
	}
}

func (s *InterviewService) Execute(ctx context.Context, sessionID, agentID, query string) (*domain.DialogueResponse, error) {
	engine, err := s.manager.GetEngine(sessionID)
	if err != nil {
		return nil, err
	}

	if !engine.Session.IsRunning() {
		return nil, domain.ErrSessionNotRunning
	}

	// Advance game time by 30 seconds for an interview exchange
	engine.Tick(30)

	// Locate agent
	var targetAgent *domain.Agent
	for _, a := range engine.Agents {
		if a.ID == agentID || a.ScenarioAgentID == agentID {
			targetAgent = a
			break
		}
	}
	if targetAgent == nil {
		return nil, domain.ErrAgentNotFound
	}

	// Build safe context
	knownClaims := make([]domain.AgentKnownClaim, 0)
	for claimID, b := range engine.Beliefs[targetAgent.ID] {
		knownClaims = append(knownClaims, domain.AgentKnownClaim{
			ClaimID:    claimID,
			Confidence: b.Confidence,
			Basis:      "memory",
		})
	}

	// Determine allowed reveal evidence IDs
	allowedEvidence := make([]string, 0)
	switch targetAgent.ScenarioAgentID {
	case "agent_security_guard":
		allowedEvidence = append(allowedEvidence, "ev_cctv_gate")
	case "agent_firefighter":
		allowedEvidence = append(allowedEvidence, "ev_fire_report")
	case "agent_paramedic", "agent_nurse":
		allowedEvidence = append(allowedEvidence, "ev_hospital_summary")
	case "agent_technician":
		allowedEvidence = append(allowedEvidence, "ev_maintenance_ticket")
	case "agent_photographer":
		allowedEvidence = append(allowedEvidence, "ev_photo_sequence")
	case "agent_factory_worker_a":
		allowedEvidence = append(allowedEvidence, "ev_worker_testimony")
	}

	dCtx := dialogue.BuildDialogueContext(dialogue.AgentContextInput{
		Agent:                    *targetAgent,
		KnownClaims:              knownClaims,
		Memories:                 engine.Memories[targetAgent.ID],
		RelationshipToPlayer:     engine.Session.PlayerCredibility,
		AllowedRevealEvidenceIDs: allowedEvidence,
	})

	resp, err := s.provider.GenerateDialogue(ctx, dCtx, query)
	if err != nil {
		return nil, err
	}

	// If the NPC revealed an evidence item, discover it in the engine
	for _, evID := range resp.RevealedEvidenceIDs {
		_, _ = engine.DiscoverEvidence(evID)
	}

	// Append World Event for dialogue response
	engine.AppendEvent(
		domain.WorldEventTypeAgentMessage,
		&targetAgent.ID, nil,
		map[string]interface{}{
			"agent_id":   targetAgent.ID,
			"name":       targetAgent.Name,
			"utterance":  resp.Utterance,
			"emotion":    resp.Emotion,
			"certainty":  resp.Certainty,
		},
	)

	return &resp, nil
}
