package scenario

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"city-of-lies/backend/internal/domain"
)

// LoadScenario reads and deserializes all scenario files from the given directory.
func LoadScenario(dir string) (*ScenarioBundle, error) {
	bundle := &ScenarioBundle{}

	// 1. scenario.json
	if err := readJSONFile(filepath.Join(dir, "scenario.json"), &bundle.Metadata); err != nil {
		return nil, fmt.Errorf("load scenario.json: %w", err)
	}

	// 2. locations.json
	if err := readJSONFile(filepath.Join(dir, "locations.json"), &bundle.Locations); err != nil {
		return nil, fmt.Errorf("load locations.json: %w", err)
	}

	// 3. agents.json
	if err := readJSONFile(filepath.Join(dir, "agents.json"), &bundle.Agents); err != nil {
		return nil, fmt.Errorf("load agents.json: %w", err)
	}

	// 4. claims.json
	type rawClaim struct {
		ID             string  `json:"id"`
		Text           string  `json:"text"`
		Classification string  `json:"classification"`
		Role           string  `json:"role"`
		Weight         float64 `json:"weight"`
	}
	var rawClaims []rawClaim
	if err := readJSONFile(filepath.Join(dir, "claims.json"), &rawClaims); err != nil {
		return nil, fmt.Errorf("load claims.json: %w", err)
	}
	for _, rc := range rawClaims {
		bundle.Claims = append(bundle.Claims, domain.Claim{
			ScenarioClaimID: rc.ID,
			DisplayText:     rc.Text,
			Classification:  rc.Classification,
			NarrativeRole:   domain.NarrativeRole(rc.Role),
			NarrativeWeight: rc.Weight,
		})
	}

	// 5. observations.json
	if err := readJSONFile(filepath.Join(dir, "observations.json"), &bundle.Observations); err != nil {
		return nil, fmt.Errorf("load observations.json: %w", err)
	}

	// 6. relationships.json
	if err := readJSONFile(filepath.Join(dir, "relationships.json"), &bundle.Relationships); err != nil {
		return nil, fmt.Errorf("load relationships.json: %w", err)
	}

	// 7. evidence.json
	if err := readJSONFile(filepath.Join(dir, "evidence.json"), &bundle.Evidence); err != nil {
		return nil, fmt.Errorf("load evidence.json: %w", err)
	}

	// 8. public-events.json
	if err := readJSONFile(filepath.Join(dir, "public-events.json"), &bundle.PublicEvents); err != nil {
		return nil, fmt.Errorf("load public-events.json: %w", err)
	}

	// 9. truth-form.json
	if err := readJSONFile(filepath.Join(dir, "truth-form.json"), &bundle.TruthForm); err != nil {
		return nil, fmt.Errorf("load truth-form.json: %w", err)
	}

	return bundle, nil
}

func readJSONFile(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
