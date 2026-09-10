package ports

import (
	"context"

	"city-of-lies/backend/internal/domain"
)

// RealtimePublisher broadcasts state events to connected WebSocket clients.
type RealtimePublisher interface {
	PublishEvent(ctx context.Context, sessionID string, event domain.WorldEvent) error
}
