package realtime

import (
	"encoding/json"
	"time"
)

type WSMessageEnvelope struct {
	Sequence    int64           `json:"sequence"`
	Type        string          `json:"type"`
	SessionID   string          `json:"session_id"`
	GameSecond  int64           `json:"game_second"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Payload     json.RawMessage `json:"payload"`
}

type ClientResumeMessage struct {
	Type         string `json:"type"` // "client.resume"
	LastSequence int64  `json:"last_sequence"`
}
