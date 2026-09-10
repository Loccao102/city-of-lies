package ports

import (
	"context"

	"city-of-lies/backend/internal/domain"
)

// DialogueProvider generates NPC utterances constrained by the supplied knowledge context.
type DialogueProvider interface {
	GenerateDialogue(ctx context.Context, dCtx domain.DialogueContext, playerQuery string) (domain.DialogueResponse, error)
}
