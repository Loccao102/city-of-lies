package belief

import (
	"testing"
)

func TestPropagationRules(t *testing.T) {
	if !EnsureWeightsSumValid() {
		t.Fatalf("propagation weights do not sum to 1.0")
	}

	baseInput := PropagationInput{
		SpeakerConfidence:    0.80,
		RelationshipTrust:    0.70,
		SpeakerInfluence:     0.60,
		SpeakerCredibility:   0.60,
		ListenerReceptiveness: 0.50,
		ListenerSkepticism:   0.30,
		OldConfidence:        0.20,
	}

	resNormal := CalculatePropagation(baseInput)
	if resNormal <= baseInput.OldConfidence {
		t.Errorf("expected belief to increase, got %f", resNormal)
	}
	if resNormal > 1.0 {
		t.Errorf("confidence exceeded 1.0: %f", resNormal)
	}

	// High trust should yield greater confidence than low trust
	lowTrustInput := baseInput
	lowTrustInput.RelationshipTrust = 0.10
	resLowTrust := CalculatePropagation(lowTrustInput)
	if resLowTrust >= resNormal {
		t.Errorf("expected high trust (%f) > low trust (%f)", resNormal, resLowTrust)
	}

	// Skeptical listener should receive smaller delta
	skepticalInput := baseInput
	skepticalInput.ListenerSkepticism = 0.90
	resSkeptical := CalculatePropagation(skepticalInput)
	if resSkeptical >= resNormal {
		t.Errorf("expected skeptical listener (%f) < normal (%f)", resSkeptical, resNormal)
	}

	// Contradictory direct observation reduces delta
	contradictoryInput := baseInput
	contradictoryInput.HasDirectContradiction = true
	resContradictory := CalculatePropagation(contradictoryInput)
	if resContradictory >= resNormal {
		t.Errorf("expected contradictory observation (%f) < normal (%f)", resContradictory, resNormal)
	}

	// Repeated root source attenuates propagation
	sameRootInput := baseInput
	sameRootInput.HasSameRootSource = true
	resSameRoot := CalculatePropagation(sameRootInput)
	if resSameRoot >= resNormal {
		t.Errorf("expected same root source (%f) < normal (%f)", resSameRoot, resNormal)
	}
}
