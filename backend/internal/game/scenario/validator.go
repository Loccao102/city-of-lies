package scenario

import (
	"fmt"
	"math"
)

// ValidateScenarioBundle enforces referential integrity and invariant checks across all scenario files.
func ValidateScenarioBundle(b *ScenarioBundle) error {
	// 1. Locations
	locationSet := make(map[string]bool)
	for _, loc := range b.Locations {
		if loc.ID == "" {
			return fmt.Errorf("empty location id found")
		}
		if locationSet[loc.ID] {
			return fmt.Errorf("duplicate location id: %s", loc.ID)
		}
		locationSet[loc.ID] = true
	}

	// 2. Agents
	agentSet := make(map[string]bool)
	for _, a := range b.Agents {
		if a.ID == "" {
			return fmt.Errorf("empty agent id found")
		}
		if agentSet[a.ID] {
			return fmt.Errorf("duplicate agent id: %s", a.ID)
		}
		agentSet[a.ID] = true

		if !locationSet[a.LocationID] {
			return fmt.Errorf("agent %s references unknown location: %s", a.ID, a.LocationID)
		}

		// Probability checks
		traits := []struct {
			name string
			val  float64
		}{
			{"influence", a.Influence},
			{"credibility", a.Credibility},
			{"skepticism", a.Skepticism},
			{"sociality", a.Sociality},
			{"deception_tendency", a.DeceptionTendency},
			{"receptiveness", a.Receptiveness},
		}
		for _, t := range traits {
			if t.val < 0.0 || t.val > 1.0 {
				return fmt.Errorf("agent %s trait %s has invalid value %f", a.ID, t.name, t.val)
			}
		}
	}

	// 3. Claims
	claimSet := make(map[string]bool)
	for _, c := range b.Claims {
		if c.ScenarioClaimID == "" {
			return fmt.Errorf("empty claim id found")
		}
		if claimSet[c.ScenarioClaimID] {
			return fmt.Errorf("duplicate claim id: %s", c.ScenarioClaimID)
		}
		claimSet[c.ScenarioClaimID] = true
	}

	// 4. Primary False Narrative Weights Sum near 1.0
	narrativeWeightSum := 0.0
	for _, c := range b.Metadata.PrimaryFalseNarrative.Claims {
		if !claimSet[c.ClaimID] {
			return fmt.Errorf("narrative references unknown claim: %s", c.ClaimID)
		}
		narrativeWeightSum += c.Weight
	}
	if math.Abs(narrativeWeightSum-1.0) > 0.01 {
		return fmt.Errorf("narrative claim weights sum to %f, must sum near 1.0", narrativeWeightSum)
	}

	// 5. Observations
	for _, obs := range b.Observations {
		if !agentSet[obs.AgentID] {
			return fmt.Errorf("observation references unknown agent: %s", obs.AgentID)
		}
		if !claimSet[obs.ClaimID] {
			return fmt.Errorf("observation references unknown claim: %s", obs.ClaimID)
		}
		if obs.Confidence < 0.0 || obs.Confidence > 1.0 {
			return fmt.Errorf("observation confidence %f outside [0, 1]", obs.Confidence)
		}
	}

	// 6. Relationships
	for _, rel := range b.Relationships {
		if !agentSet[rel.FromAgentID] {
			return fmt.Errorf("relationship from unknown agent: %s", rel.FromAgentID)
		}
		if !agentSet[rel.ToAgentID] {
			return fmt.Errorf("relationship to unknown agent: %s", rel.ToAgentID)
		}
		if rel.FromAgentID == rel.ToAgentID {
			return fmt.Errorf("self-referential relationship loop on agent: %s", rel.FromAgentID)
		}
		if rel.Trust < 0.0 || rel.Trust > 1.0 {
			return fmt.Errorf("relationship trust %f outside [0, 1]", rel.Trust)
		}
	}

	// 7. Evidence
	evidenceSet := make(map[string]bool)
	for _, ev := range b.Evidence {
		if ev.ID == "" {
			return fmt.Errorf("empty evidence id found")
		}
		if evidenceSet[ev.ID] {
			return fmt.Errorf("duplicate evidence id: %s", ev.ID)
		}
		evidenceSet[ev.ID] = true

		if !locationSet[ev.LocationID] {
			return fmt.Errorf("evidence %s references unknown location: %s", ev.ID, ev.LocationID)
		}
		if ev.Reliability < 0.0 || ev.Reliability > 1.0 {
			return fmt.Errorf("evidence reliability %f outside [0, 1]", ev.Reliability)
		}
		for _, sc := range ev.Supports {
			if !claimSet[sc] {
				return fmt.Errorf("evidence %s supports unknown claim: %s", ev.ID, sc)
			}
		}
		for _, cc := range ev.Contradicts {
			if !claimSet[cc] {
				return fmt.Errorf("evidence %s contradicts unknown claim: %s", ev.ID, cc)
			}
		}
	}

	// 8. Required Truth Fields
	for _, rf := range b.Metadata.RequiredTruthFields {
		if _, exists := b.TruthForm[rf]; !exists {
			return fmt.Errorf("required truth field %q missing from truth-form.json", rf)
		}
	}

	return nil
}
