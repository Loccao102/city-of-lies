package domain

type EvidenceRelation string

const (
	EvidenceRelationSupports   EvidenceRelation = "supports"
	EvidenceRelationContradicts EvidenceRelation = "contradicts"
)

// EvidenceItem represents a collectible clue, document, or record.
type EvidenceItem struct {
	ID                     string   `json:"id"`
	SessionID              string   `json:"session_id"`
	ScenarioEvidenceID     string   `json:"scenario_evidence_id"`
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	LocationID             string   `json:"location_id"`
	Reliability            float64  `json:"reliability"`
	SourceCategory         string   `json:"source_category"`
	Discovered             bool     `json:"discovered"`
	DiscoveredAtGameSecond *int64   `json:"discovered_at_game_second,omitempty"`
}

// EvidenceClaimLink defines whether an evidence item supports or contradicts a claim.
type EvidenceClaimLink struct {
	EvidenceID string           `json:"evidence_id"`
	ClaimID    string           `json:"claim_id"`
	Relation   EvidenceRelation `json:"relation"`
	Weight     float64          `json:"weight"`
}
