package dto

type StartSessionRequest struct {
	ScenarioID string `json:"scenario_id"`
	Seed       *int64 `json:"seed,omitempty"`
}

type InterviewRequest struct {
	Message string `json:"message"`
}

type PublishCorrectionRequest struct {
	ChallengedClaimID string   `json:"challenged_claim_id"`
	EvidenceIDs       []string `json:"evidence_ids"`
	Message           string   `json:"message"`
}

type SubmitTruthRequest struct {
	EventType       string `json:"event_type"`
	CauseCategory   string `json:"cause_category"`
	MajorExplosion  bool   `json:"major_explosion"`
	Fatalities      int    `json:"fatalities"`
	ChemicalRelease *bool  `json:"chemical_release,omitempty"`
	ManagementIssue string `json:"management_issue,omitempty"`
}
