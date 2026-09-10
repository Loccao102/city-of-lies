package domain

import (
	"encoding/json"
	"time"
)

type WorldEventType string

const (
	WorldEventTypeSessionStarted      WorldEventType = "session.started"
	WorldEventTypeSessionStatusChange WorldEventType = "session.status_changed"
	WorldEventTypeAgentMoved          WorldEventType = "agent.moved"
	WorldEventTypeRumorShared         WorldEventType = "rumor.shared"
	WorldEventTypeRumorVisualized     WorldEventType = "rumor.visualized"
	WorldEventTypeBeliefChanged       WorldEventType = "belief.changed"
	WorldEventTypeBeliefStatsUpdated  WorldEventType = "belief_stats.updated"
	WorldEventTypeEvidenceDiscovered  WorldEventType = "evidence.discovered"
	WorldEventTypeCorrectionPublished WorldEventType = "correction.published"
	WorldEventTypeAgentMessage        WorldEventType = "conversation.agent_message"
	WorldEventTypeGameWon             WorldEventType = "game.won"
	WorldEventTypeGameLost            WorldEventType = "game.lost"
)

// WorldEvent represents an immutable timeline event envelope.
type WorldEvent struct {
	ID           string          `json:"id"`
	SessionID    string          `json:"session_id"`
	Sequence     int64           `json:"sequence"`
	EventType    WorldEventType  `json:"type"`
	GameSecond   int64           `json:"game_second"`
	ActorAgentID *string         `json:"actor_agent_id,omitempty"`
	TargetAgentID *string        `json:"target_agent_id,omitempty"`
	Payload      json.RawMessage `json:"payload"`
	OccurredAt   time.Time       `json:"occurred_at"`
}
