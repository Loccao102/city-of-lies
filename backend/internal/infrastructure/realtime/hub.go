package realtime

import (
	"encoding/json"
	"sync"

	"city-of-lies/backend/internal/domain"
)

type Hub struct {
	mu           sync.RWMutex
	clients      map[string]map[*Client]bool // sessionID -> set of Clients
	eventHistory map[string][]domain.WorldEvent // sessionID -> recent event cache for replay
}

func NewHub() *Hub {
	return &Hub{
		clients:      make(map[string]map[*Client]bool),
		eventHistory: make(map[string][]domain.WorldEvent),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[c.sessionID]; !exists {
		h.clients[c.sessionID] = make(map[*Client]bool)
	}
	h.clients[c.sessionID][c] = true
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if group, exists := h.clients[c.sessionID]; exists {
		delete(group, c)
		if len(group) == 0 {
			delete(h.clients, c.sessionID)
		}
	}
}

func (h *Hub) Broadcast(sessionID string, evt domain.WorldEvent) {
	h.mu.Lock()
	// Store in recent history (last 100 events)
	history := h.eventHistory[sessionID]
	history = append(history, evt)
	if len(history) > 100 {
		history = history[len(history)-100:]
	}
	h.eventHistory[sessionID] = history

	group := make([]*Client, 0)
	if clients, exists := h.clients[sessionID]; exists {
		for c := range clients {
			group = append(group, c)
		}
	}
	h.mu.Unlock()

	envelope := WSMessageEnvelope{
		Sequence:   evt.Sequence,
		Type:       string(evt.EventType),
		SessionID:  sessionID,
		GameSecond: evt.GameSecond,
		OccurredAt: evt.OccurredAt,
		Payload:    evt.Payload,
	}
	msgBytes, err := json.Marshal(envelope)
	if err != nil {
		return
	}

	for _, c := range group {
		c.Send(msgBytes)
	}
}

func (h *Hub) HandleResume(c *Client, lastSequence int64) {
	h.mu.RLock()
	history, exists := h.eventHistory[c.sessionID]
	if !exists {
		h.mu.RUnlock()
		return
	}
	replay := make([]domain.WorldEvent, 0)
	for _, ev := range history {
		if ev.Sequence > lastSequence {
			replay = append(replay, ev)
		}
	}
	h.mu.RUnlock()

	for _, ev := range replay {
		envelope := WSMessageEnvelope{
			Sequence:   ev.Sequence,
			Type:       string(ev.EventType),
			SessionID:  c.sessionID,
			GameSecond: ev.GameSecond,
			OccurredAt: ev.OccurredAt,
			Payload:    ev.Payload,
		}
		bytes, err := json.Marshal(envelope)
		if err == nil {
			c.Send(bytes)
		}
	}
}
