package domain

// Publication represents an official fact-checking correction issued by the player.
type Publication struct {
	ID                      string   `json:"id"`
	SessionID               string   `json:"session_id"`
	ChallengedClaimID       string   `json:"challenged_claim_id"`
	Message                 string   `json:"message"`
	Strength                float64  `json:"strength"`
	PlayerCredibilityBefore float64  `json:"player_credibility_before"`
	PlayerCredibilityAfter  float64  `json:"player_credibility_after"`
	PublishedAtGameSecond   int64    `json:"published_at_game_second"`
	EvidenceIDs             []string `json:"evidence_ids"`
}

// GroundTruthSubmission represents the final report submitted by the player to resolve the investigation.
type GroundTruthSubmission struct {
	EventType       string `json:"event_type"`
	CauseCategory   string `json:"cause_category"`
	MajorExplosion  bool   `json:"major_explosion"`
	Fatalities      int    `json:"fatalities"`
	ChemicalRelease *bool  `json:"chemical_release,omitempty"`
	ManagementIssue string `json:"management_issue,omitempty"`
}

// TruthSubmissionResult contains the evaluation feedback for the player's submission.
type TruthSubmissionResult struct {
	Won              bool              `json:"won"`
	Terminal         bool              `json:"terminal"`
	EvidenceStrength float64           `json:"evidence_strength"`
	RequiredEvidence float64           `json:"required_evidence"`
	FalseBeliefRatio float64           `json:"false_belief_ratio"`
	FieldStatus      map[string]bool   `json:"field_status"`
	FeedbackMessage  string            `json:"feedback_message"`
}
