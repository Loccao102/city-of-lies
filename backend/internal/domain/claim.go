package domain

type NarrativeRole string

const (
	NarrativeRoleFact           NarrativeRole = "fact"
	NarrativeRolePrimaryFalse   NarrativeRole = "primary_false"
	NarrativeRoleSecondaryFalse NarrativeRole = "secondary_false"
	NarrativeRoleMeta           NarrativeRole = "meta"
)

// Claim represents an informational statement circulating in the district.
type Claim struct {
	ID              string        `json:"id"`
	SessionID       string        `json:"session_id"`
	ScenarioClaimID string        `json:"scenario_claim_id"`
	DisplayText     string        `json:"display_text"`
	Classification  string        `json:"-"` // SERVER-ONLY: never serialize to player DTO
	NarrativeRole   NarrativeRole `json:"narrative_role"`
	NarrativeWeight float64       `json:"narrative_weight"`
}

// IsPrimaryFalse returns true if this claim contributes to the defeat false narrative score.
func (c *Claim) IsPrimaryFalse() bool {
	return c.NarrativeRole == NarrativeRolePrimaryFalse
}
