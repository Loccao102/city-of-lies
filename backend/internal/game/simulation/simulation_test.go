package simulation

import (
	"path/filepath"
	"testing"

	"city-of-lies/backend/internal/config"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/scenario"
)

func TestSimulationEngineProgression(t *testing.T) {
	scenarioDir := filepath.Join("..", "..", "..", "..", "data", "scenarios", "riverside-factory")
	bundle, err := scenario.LoadScenario(scenarioDir)
	if err != nil {
		t.Fatalf("failed to load scenario: %v", err)
	}

	cfg, _ := config.Load()
	cfg.Simulation.MaxNewConversationsPerTick = 4
	cfg.Simulation.PairCooldownSeconds = 10

	engine := NewSimulationEngine(bundle, 12345, cfg)
	initialRatio := engine.Session.FalseNarrativeRatio

	// Simulate 120 ticks
	for i := 0; i < 120; i++ {
		engine.Tick(1)
		if !engine.Session.IsRunning() {
			break
		}
	}

	if engine.Session.GameSecond != 120 && engine.Session.Status == domain.SessionStatusRunning {
		t.Errorf("expected 120 game seconds, got %d", engine.Session.GameSecond)
	}

	finalRatio := engine.Session.FalseNarrativeRatio
	if finalRatio < initialRatio {
		t.Errorf("expected false narrative ratio to increase or stay steady, from %f to %f", initialRatio, finalRatio)
	}

	snapshot := engine.BuildSnapshot()
	if len(snapshot.Agents) != 20 {
		t.Errorf("expected 20 agents in snapshot, got %d", len(snapshot.Agents))
	}
}
