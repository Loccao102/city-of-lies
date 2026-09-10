package config

import (
	"testing"
	"time"
)

func TestConfigValidation(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected valid default config, got error: %v", err)
	}

	if cfg.HTTP.Addr != DefaultHTTPAddr {
		t.Errorf("expected %s, got %s", DefaultHTTPAddr, cfg.HTTP.Addr)
	}

	// Test invalid belief adoption threshold
	cfg.Game.BeliefAdoptionThreshold = 1.5
	if err := cfg.Validate(); err == nil {
		t.Errorf("expected validation error for threshold > 1.0, got nil")
	}

	// Test invalid tick duration
	cfg.Game.BeliefAdoptionThreshold = 0.65
	cfg.Simulation.TickDuration = -10 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Errorf("expected validation error for negative tick duration, got nil")
	}
}
