package scenario

import (
	"city-of-lies/backend/internal/domain"
)

type ScenarioMetadata struct {
	ID                    string `json:"id"`
	Title                 string `json:"title"`
	Version               int    `json:"version"`
	Description           string `json:"description"`
	PrimaryFalseNarrative struct {
		AdoptionThreshold float64 `json:"adoption_threshold"`
		DefeatRatio       float64 `json:"defeat_ratio"`
		Claims            []struct {
			ClaimID string  `json:"claim_id"`
			Weight  float64 `json:"weight"`
		} `json:"claims"`
	} `json:"primary_false_narrative"`
	RequiredTruthFields []string `json:"required_truth_fields"`
	BonusTruthFields    []string `json:"bonus_truth_fields"`
}

type AgentSeed struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Role              string  `json:"role"`
	LocationID        string  `json:"location_id"`
	Influence         float64 `json:"influence"`
	Credibility       float64 `json:"credibility"`
	Skepticism        float64 `json:"skepticism"`
	Sociality         float64 `json:"sociality"`
	DeceptionTendency float64 `json:"deception_tendency"`
	Receptiveness     float64 `json:"receptiveness"`
	Personality       string  `json:"personality"`
	PrivateMotive     string  `json:"private_motive"`
	DirectKnowledge   string  `json:"direct_knowledge"`
	BlindSpots        string  `json:"blind_spots"`
}

type ObservationSeed struct {
	AgentID    string  `json:"agent_id"`
	ClaimID    string  `json:"claim_id"`
	Confidence float64 `json:"confidence"`
	Basis      string  `json:"basis"`
}

type RelationshipSeed struct {
	FromAgentID string  `json:"from_agent_id"`
	ToAgentID   string  `json:"to_agent_id"`
	Trust       float64 `json:"trust"`
	Reason      string  `json:"reason"`
}

type EvidenceSeed struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	LocationID      string   `json:"location_id"`
	Reliability     float64  `json:"reliability"`
	SourceCategory  string   `json:"source_category"`
	Supports        []string `json:"supports"`
	Contradicts     []string `json:"contradicts"`
	UnlockCondition string   `json:"unlock_condition"`
}

type PublicEventSeed struct {
	GameSecond  int64  `json:"game_second"`
	EventType   string `json:"event_type"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
	LocationID  string `json:"location_id"`
}

type TruthFormField struct {
	Type    string      `json:"type"`
	Label   string      `json:"label"`
	Answer  interface{} `json:"answer"`
	Options []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"options,omitempty"`
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

// ScenarioBundle holds all parsed assets for a single scenario.
type ScenarioBundle struct {
	Metadata      ScenarioMetadata
	Locations     []domain.Location
	Agents        []AgentSeed
	Claims        []domain.Claim
	Observations  []ObservationSeed
	Relationships []RelationshipSeed
	Evidence      []EvidenceSeed
	PublicEvents  []PublicEventSeed
	TruthForm     map[string]TruthFormField
}
