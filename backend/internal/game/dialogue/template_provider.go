package dialogue

import (
	"context"
	"fmt"
	"strings"

	"city-of-lies/backend/internal/domain"
)

type TemplateProvider struct{}

func NewTemplateProvider() *TemplateProvider {
	return &TemplateProvider{}
}

// GenerateDialogue produces a deterministic response using the character's direct knowledge.
func (p *TemplateProvider) GenerateDialogue(ctx context.Context, dCtx domain.DialogueContext, playerQuery string) (domain.DialogueResponse, error) {
	name := dCtx.Identity.Name
	role := dCtx.Identity.Role

	var claims []string
	for _, c := range dCtx.KnownClaims {
		claims = append(claims, c.ClaimID)
	}

	var utterance string
	intent := domain.DialogueIntentAnswer
	emotion := domain.DialogueEmotionNeutral
	certainty := 0.80
	var revealedEv []string

	// Check if this character has allowed evidence to reveal
	if len(dCtx.AllowedRevealEvidenceIDs) > 0 {
		revealedEv = append(revealedEv, dCtx.AllowedRevealEvidenceIDs[0])
	}

	// Tailor utterance based on memories or known claims
	if len(dCtx.RelevantMemories) > 0 {
		mem := dCtx.RelevantMemories[0].Content
		utterance = fmt.Sprintf("Tôi là %s (%s). Về việc đó: %s", name, role, mem)
	} else if len(dCtx.KnownClaims) > 0 {
		utterance = fmt.Sprintf("Tôi chỉ biết những gì nghe được gần đây quanh khu vực. Thông tin vẫn còn khá hỗn loạn.", )
		emotion = domain.DialogueEmotionUncertain
		certainty = 0.50
	} else {
		utterance = fmt.Sprintf("Tôi là %s. Lúc xảy ra sự cố tôi không có mặt trực tiếp nên không rõ chi tiết cụ thể.", name)
		intent = domain.DialogueIntentDeflect
		emotion = domain.DialogueEmotionUncertain
		certainty = 0.30
	}

	if strings.Contains(strings.ToLower(playerQuery), "ai chết") || strings.Contains(strings.ToLower(playerQuery), "tử vong") {
		hasDeathClaim := false
		for _, c := range dCtx.KnownClaims {
			if c.ClaimID == "claim_zero_fatalities" {
				utterance = fmt.Sprintf("Theo những gì tôi nắm được thì không có ca tử vong nào cả.")
				hasDeathClaim = true
				emotion = domain.DialogueEmotionCalm
				certainty = 0.90
				break
			} else if c.ClaimID == "claim_multiple_deaths" {
				utterance = fmt.Sprintf("Tôi có nghe người ta đồn có người chết, nhưng tôi chưa tận mắt thấy bằng chứng xác thực.")
				hasDeathClaim = true
				emotion = domain.DialogueEmotionAfraid
				certainty = 0.60
				break
			}
		}
		if !hasDeathClaim {
			utterance = fmt.Sprintf("Về thương vong, tôi không dám khẳng định vì chưa có báo cáo chính thức.")
			intent = domain.DialogueIntentRefuse
		}
	}

	return domain.DialogueResponse{
		Intent:             intent,
		Utterance:          utterance,
		ReferencedClaimIDs: claims,
		RevealedEvidenceIDs: revealedEv,
		Emotion:            emotion,
		Certainty:          certainty,
	}, nil
}
