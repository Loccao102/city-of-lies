package dialogue

import (
	"fmt"
	"strings"

	"city-of-lies/backend/internal/domain"
)

type KnowledgeGuard struct{}

func NewKnowledgeGuard() *KnowledgeGuard {
	return &KnowledgeGuard{}
}

// ValidateResponse checks if an agent utterance conforms to safety and knowledge boundaries.
func (g *KnowledgeGuard) ValidateResponse(resp domain.DialogueResponse, ctx domain.DialogueContext) error {
	// 1. Check intent
	switch resp.Intent {
	case domain.DialogueIntentAnswer, domain.DialogueIntentDeflect,
		domain.DialogueIntentLie, domain.DialogueIntentRefuse,
		domain.DialogueIntentAskQuestion:
		// Valid
	default:
		return fmt.Errorf("disallowed intent: %q", resp.Intent)
	}

	// 2. Check emotion
	switch resp.Emotion {
	case domain.DialogueEmotionNeutral, domain.DialogueEmotionCalm,
		domain.DialogueEmotionUncertain, domain.DialogueEmotionAfraid,
		domain.DialogueEmotionAngry, domain.DialogueEmotionExcited,
		domain.DialogueEmotionDefensive, domain.DialogueEmotionSad:
		// Valid
	default:
		return fmt.Errorf("disallowed emotion: %q", resp.Emotion)
	}

	// 3. Utterance non-empty and bounded
	trimmed := strings.TrimSpace(resp.Utterance)
	if trimmed == "" {
		return fmt.Errorf("utterance cannot be empty")
	}
	if len(trimmed) > 1200 {
		return fmt.Errorf("utterance length (%d) exceeds limit (1200)", len(trimmed))
	}

	// 4. Ensure all referenced claims are actually known by this agent
	knownClaimMap := make(map[string]bool)
	for _, kc := range ctx.KnownClaims {
		knownClaimMap[kc.ClaimID] = true
	}
	for _, rc := range resp.ReferencedClaimIDs {
		if !knownClaimMap[rc] {
			return fmt.Errorf("agent referenced unknown claim: %q", rc)
		}
	}

	// 5. Ensure revealed evidence is explicitly allowed
	allowedEvidenceMap := make(map[string]bool)
	for _, evID := range ctx.AllowedRevealEvidenceIDs {
		allowedEvidenceMap[evID] = true
	}
	for _, rev := range resp.RevealedEvidenceIDs {
		if !allowedEvidenceMap[rev] {
			return fmt.Errorf("agent attempted to reveal unpermitted evidence: %q", rev)
		}
	}

	return nil
}
