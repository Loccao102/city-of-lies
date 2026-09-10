package dialogue

import (
	"city-of-lies/backend/internal/domain"
)

type AgentContextInput struct {
	Agent                    domain.Agent
	KnownClaims              []domain.AgentKnownClaim
	Memories                 []domain.AgentMemory
	RelationshipToPlayer     float64
	History                  []domain.ChatMessage
	AllowedRevealEvidenceIDs []string
}

// BuildDialogueContext constructs a privacy-safe, knowledge-guarded context payload.
func BuildDialogueContext(in AgentContextInput) domain.DialogueContext {
	var ctx domain.DialogueContext
	ctx.Identity.Name = in.Agent.Name
	ctx.Identity.Role = in.Agent.Role
	ctx.Traits.Skepticism = in.Agent.Skepticism
	ctx.Traits.DeceptionTendency = in.Agent.DeceptionTendency
	ctx.CurrentLocation = in.Agent.CurrentLocationID
	ctx.RelationshipToPlayer = in.RelationshipToPlayer

	ctx.KnownClaims = in.KnownClaims
	if ctx.KnownClaims == nil {
		ctx.KnownClaims = []domain.AgentKnownClaim{}
	}

	ctx.RelevantMemories = make([]domain.MemoryBrief, 0, len(in.Memories))
	for _, m := range in.Memories {
		ctx.RelevantMemories = append(ctx.RelevantMemories, domain.MemoryBrief{
			Type:    string(m.MemoryType),
			Content: m.Content,
		})
	}

	ctx.ConversationHistory = in.History
	if ctx.ConversationHistory == nil {
		ctx.ConversationHistory = []domain.ChatMessage{}
	}

	ctx.AllowedRevealEvidenceIDs = in.AllowedRevealEvidenceIDs
	if ctx.AllowedRevealEvidenceIDs == nil {
		ctx.AllowedRevealEvidenceIDs = []string{}
	}

	return ctx
}
