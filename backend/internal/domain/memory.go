package domain

type MemoryType string

const (
	MemoryTypeDirectObservation MemoryType = "direct_observation"
	MemoryTypeHearsay           MemoryType = "hearsay"
	MemoryTypeOfficialStatement MemoryType = "official_statement"
	MemoryTypeSocialMedia       MemoryType = "social_media"
	MemoryTypeCorrection        MemoryType = "correction"
)

// AgentMemory represents an atomic recollection or observation stored by an agent.
type AgentMemory struct {
	ID                   string     `json:"id"`
	SessionID            string     `json:"session_id"`
	AgentID              string     `json:"agent_id"`
	MemoryType           MemoryType `json:"memory_type"`
	Content              string     `json:"content"`
	Reliability          float64    `json:"reliability"`
	Importance           float64    `json:"importance"`
	SourceAgentID        *string    `json:"source_agent_id,omitempty"`
	RootSourceAgentID    *string    `json:"root_source_agent_id,omitempty"`
	RelatedClaimIDs      []string   `json:"related_claim_ids"`
	CreatedAtGameSecond  int64      `json:"created_at_game_second"`
}
