package belief

import (
	"testing"
)

func TestCorrectionCalculations(t *testing.T) {
	params := CorrectionParams{
		EvidenceAverageReliability: 0.95,
		SourceDiversityRatio:       0.80,
		PlayerCredibility:          0.50,
		CorroborationScore:         0.70,
	}

	strength := CalculateCorrectionStrength(params)
	if strength <= 0.5 || strength > 1.0 {
		t.Errorf("expected high strength between 0.5 and 1.0, got %f", strength)
	}

	oldConf := 0.80
	newConf := ApplyCorrectionImpact(oldConf, strength, 0.70, 0.10)
	if newConf >= oldConf {
		t.Errorf("expected correction to lower confidence from %f, got %f", oldConf, newConf)
	}

	// Test credibility update
	newCred := UpdatePlayerCredibility(0.50, strength, false)
	if newCred <= 0.50 {
		t.Errorf("expected credibility to increase from 0.50, got %f", newCred)
	}

	penalizedCred := UpdatePlayerCredibility(0.50, 0.20, true)
	if penalizedCred >= 0.50 {
		t.Errorf("expected credibility to decrease from 0.50, got %f", penalizedCred)
	}
}
