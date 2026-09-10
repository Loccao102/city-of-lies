package belief

import (
	"math"
)

type PropagationInput struct {
	SpeakerConfidence    float64
	RelationshipTrust    float64
	SpeakerInfluence     float64
	SpeakerCredibility   float64
	ListenerReceptiveness float64
	ListenerSkepticism   float64
	HasDirectContradiction bool
	HasSameRootSource      bool
	OldConfidence        float64
	BeliefUpdateRate     float64 // default: 0.38
	SkepticismEffect     float64 // default: 0.55
	DirectEvidenceMultiplier float64 // default: 0.35
	SameRootSourceMultiplier float64 // default: 0.45
}

// CalculatePropagation computes the new confidence level following a positive claim assertion.
func CalculatePropagation(in PropagationInput) float64 {
	// 1. Base propagation strength (weights must sum to 1.0)
	base := in.SpeakerConfidence*0.30 +
		in.RelationshipTrust*0.25 +
		in.SpeakerInfluence*0.15 +
		in.SpeakerCredibility*0.15 +
		in.ListenerReceptiveness*0.15

	// 2. Skepticism multiplier
	skepEffect := in.SkepticismEffect
	if skepEffect == 0 {
		skepEffect = 0.55
	}
	skepticismMultiplier := 1.0 - (in.ListenerSkepticism * skepEffect)
	if skepticismMultiplier < 0 {
		skepticismMultiplier = 0
	}

	// 3. Direct evidence contradiction modifier
	evidenceMod := 1.0
	if in.HasDirectContradiction {
		if in.DirectEvidenceMultiplier > 0 {
			evidenceMod = in.DirectEvidenceMultiplier
		} else {
			evidenceMod = 0.35
		}
	}

	// 4. Repeated root source attenuation
	provenanceMod := 1.0
	if in.HasSameRootSource {
		if in.SameRootSourceMultiplier > 0 {
			provenanceMod = in.SameRootSourceMultiplier
		} else {
			provenanceMod = 0.45
		}
	}

	// 5. Update rate
	updateRate := in.BeliefUpdateRate
	if updateRate == 0 {
		updateRate = 0.38
	}

	// 6. Compute delta
	delta := (1.0 - in.OldConfidence) *
		base *
		updateRate *
		skepticismMultiplier *
		evidenceMod *
		provenanceMod

	newConfidence := in.OldConfidence + delta
	return math.Min(1.0, math.Max(0.0, newConfidence))
}

// EnsureWeightsSumValid validates that default propagation weights sum to 1.0.
func EnsureWeightsSumValid() bool {
	sum := 0.30 + 0.25 + 0.15 + 0.15 + 0.15
	return math.Abs(sum-1.0) < 1e-6
}
