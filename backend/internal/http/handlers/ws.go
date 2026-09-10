package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coder/websocket"

	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/simulation"
	"city-of-lies/backend/internal/infrastructure/realtime"
	"city-of-lies/backend/internal/http/respond"
)

type WSHandler struct {
	hub     *realtime.Hub
	manager *simulation.SessionManager
}

func NewWSHandler(hub *realtime.Hub, manager *simulation.SessionManager) *WSHandler {
	return &WSHandler{
		hub:     hub,
		manager: manager,
	}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		respond.Error(w, http.StatusBadRequest, "session_id query param required")
		return
	}

	engine, err := h.manager.GetEngine(sessionID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, "session not found")
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Local dev & frontend cross-origin
	})
	if err != nil {
		return
	}

	client := realtime.NewClient(h.hub, conn, sessionID, 256)
	h.hub.Register(client)

	// Send initial snapshot on connect
	snapshot := engine.BuildSnapshot()
	snapshotBytes, _ := json.Marshal(snapshot)
	initialEnvelope := realtime.WSMessageEnvelope{
		Sequence:   snapshot.Sequence,
		Type:       "session.snapshot",
		SessionID:  sessionID,
		GameSecond: snapshot.GameSecond,
		OccurredAt: time.Now(),
		Payload:    snapshotBytes,
	}
	envBytes, _ := json.Marshal(initialEnvelope)
	client.Send(envBytes)

	// Wire engine events directly to hub broadcast
	engine.OnEventEmitted = func(evt domain.WorldEvent) {
		h.hub.Broadcast(sessionID, evt)
	}

	go client.WritePump(5*time.Second, 20*time.Second)
	client.ReadPump()
}
