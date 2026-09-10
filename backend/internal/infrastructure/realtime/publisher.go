package realtime

import (
	"context"

	"city-of-lies/backend/internal/domain"
)

type HubPublisher struct {
	hub *Hub
}

func NewHubPublisher(hub *Hub) *HubPublisher {
	return &HubPublisher{hub: hub}
}

func (p *HubPublisher) PublishEvent(ctx context.Context, sessionID string, event domain.WorldEvent) error {
	p.hub.Broadcast(sessionID, event)
	return nil
}
