package simulation

import (
	"time"

	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/belief"
	"city-of-lies/backend/internal/game/evidence"
	"city-of-lies/backend/internal/game/outcome"
	"city-of-lies/backend/internal/game/scenario"
)

func getStringField(tf map[string]scenario.TruthFormField, key, fallback string) string {
	if tf == nil {
		return fallback
	}
	if f, ok := tf[key]; ok {
		if s, ok := f.Answer.(string); ok {
			return s
		}
	}
	return fallback
}

func getBoolField(tf map[string]scenario.TruthFormField, key string, fallback bool) bool {
	if tf == nil {
		return fallback
	}
	if f, ok := tf[key]; ok {
		if b, ok := f.Answer.(bool); ok {
			return b
		}
	}
	return fallback
}

func getIntField(tf map[string]scenario.TruthFormField, key string, fallback int) int {
	if tf == nil {
		return fallback
	}
	if f, ok := tf[key]; ok {
		switch v := f.Answer.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case int64:
			return int(v)
		}
	}
	return fallback
}

// DiscoverEvidence marks an evidence item as found and recalculates evidence strength.
func (e *SimulationEngine) DiscoverEvidence(scenarioEvidenceID string) (*domain.EvidenceItem, error) {

	e.mu.Lock()
	defer e.mu.Unlock()

	ev, ok := e.Evidence[scenarioEvidenceID]
	if !ok {
		return nil, domain.ErrEvidenceNotFound
	}

	if !ev.Discovered {
		ev.Discovered = true
		sec := e.Session.GameSecond
		ev.DiscoveredAtGameSecond = &sec
		e.recalculateEvidenceStrength()

		e.AppendEvent(
			domain.WorldEventTypeEvidenceDiscovered,
			nil, nil,
			map[string]interface{}{
				"evidence_id": ev.ScenarioEvidenceID,
				"name":        ev.Name,
				"location_id": ev.LocationID,
				"reliability": ev.Reliability,
			},
		)
	}

	return ev, nil
}

func (e *SimulationEngine) recalculateEvidenceStrength() {
	var evalItems []evidence.EvidenceEvaluationItem
	for _, ev := range e.Evidence {
		if ev.Discovered {
			evalItems = append(evalItems, evidence.EvidenceEvaluationItem{
				ID:             ev.ScenarioEvidenceID,
				Reliability:    ev.Reliability,
				SourceCategory: ev.SourceCategory,
				CoveredFields:  evidence.MapEvidenceToFields(ev.ScenarioEvidenceID),
			})
		}
	}
	strength, _, _, _ := evidence.CalculateEvidenceStrength(evalItems, e.ScenarioBundle.Metadata.RequiredTruthFields)
	e.Session.EvidenceStrength = strength
}

// PublishCorrection processes a player's fact-check correction against a specific claim.
func (e *SimulationEngine) PublishCorrection(challengedClaimID string, evidenceIDs []string, message string) (*domain.Publication, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.Session.IsRunning() {
		return nil, domain.ErrSessionNotRunning
	}

	claim, ok := e.Claims[challengedClaimID]
	if !ok {
		return nil, domain.ErrClaimNotFound
	}

	// Validate evidence is discovered
	totalRel := 0.0
	categories := make(map[string]bool)
	for _, evID := range evidenceIDs {
		ev, ok := e.Evidence[evID]
		if !ok || !ev.Discovered {
			return nil, domain.ErrEvidenceNotDiscovered
		}
		totalRel += ev.Reliability
		categories[ev.SourceCategory] = true
	}

	avgRel := 0.0
	if len(evidenceIDs) > 0 {
		avgRel = totalRel / float64(len(evidenceIDs))
	}
	diversityRatio := float64(len(categories)) / 4.0

	// Calculate correction strength
	strength := belief.CalculateCorrectionStrength(belief.CorrectionParams{
		EvidenceAverageReliability: avgRel,
		SourceDiversityRatio:       diversityRatio,
		PlayerCredibility:          e.Session.PlayerCredibility,
		CorroborationScore:         avgRel,
	})

	// Adjust player credibility
	oldCred := e.Session.PlayerCredibility
	newCred := belief.UpdatePlayerCredibility(oldCred, strength, len(evidenceIDs) == 0)
	e.Session.PlayerCredibility = newCred

	// Apply correction impact to agents holding the challenged claim
	for _, a := range e.Agents {
		if b, exists := e.Beliefs[a.ID][claim.ScenarioClaimID]; exists && b.Confidence > 0.10 {
			b.Confidence = belief.ApplyCorrectionImpact(b.Confidence, strength, a.Receptiveness, b.EvidenceResistance)
		}
	}

	// Recalculate false belief ratio
	e.recalculateNarrativeRatio()

	pub := &domain.Publication{
		ID:                      domain.NewUUID(),
		SessionID:               e.Session.ID,
		ChallengedClaimID:       challengedClaimID,
		Message:                 message,
		Strength:                strength,
		PlayerCredibilityBefore: oldCred,
		PlayerCredibilityAfter:  newCred,
		PublishedAtGameSecond:   e.Session.GameSecond,
		EvidenceIDs:             evidenceIDs,
	}

	e.AppendEvent(
		domain.WorldEventTypeCorrectionPublished,
		nil, nil,
		map[string]interface{}{
			"challenged_claim": claim.DisplayText,
			"strength":         strength,
			"player_credibility": newCred,
		},
	)

	return pub, nil
}

// SubmitGroundTruth evaluates the player's final investigative submission.
func (e *SimulationEngine) SubmitGroundTruth(sub domain.GroundTruthSubmission) (domain.TruthSubmissionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	expectedAnswers := outcome.GroundTruthAnswers{
		EventType:      getStringField(e.ScenarioBundle.TruthForm, "event_type", "fire"),
		CauseCategory:  getStringField(e.ScenarioBundle.TruthForm, "cause_category", "electrical_fault"),
		MajorExplosion: getBoolField(e.ScenarioBundle.TruthForm, "major_explosion", false),
		Fatalities:     getIntField(e.ScenarioBundle.TruthForm, "fatalities", 0),
	}

	cfg := outcome.EvaluatorConfig{
		TruthEvidenceThreshold: e.Config.Game.TruthEvidenceThreshold,
		DefeatFalseBeliefRatio: e.Config.Game.DefeatFalseBeliefRatio,
	}

	res := outcome.EvaluateSubmission(sub, expectedAnswers, e.Session.EvidenceStrength, e.Session.FalseNarrativeRatio, cfg)

	if res.Won {
		e.Session.Status = domain.SessionStatusWon
		now := time.Now()
		e.Session.EndedAt = &now
		e.AppendEvent(
			domain.WorldEventTypeGameWon,
			nil, nil,
			map[string]interface{}{
				"feedback":          res.FeedbackMessage,
				"evidence_strength": res.EvidenceStrength,
				"false_belief_ratio": res.FalseBeliefRatio,
			},
		)
	}

	return res, nil
}
