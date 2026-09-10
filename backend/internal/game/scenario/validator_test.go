package scenario

import (
	"path/filepath"
	"testing"
)

func TestValidateRiversideScenario(t *testing.T) {
	scenarioDir := filepath.Join("..", "..", "..", "..", "data", "scenarios", "riverside-factory")
	bundle, err := LoadScenario(scenarioDir)
	if err != nil {
		t.Fatalf("failed to load scenario bundle from %s: %v", scenarioDir, err)
	}

	if err := ValidateScenarioBundle(bundle); err != nil {
		t.Errorf("scenario bundle validation failed: %v", err)
	}

	if len(bundle.Agents) != 20 {
		t.Errorf("expected exactly 20 agents in Riverside scenario, got %d", len(bundle.Agents))
	}
}
