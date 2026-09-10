package belief

import (
	"math"

	"city-of-lies/backend/internal/domain"
)

type CorrectionParams struct {
	EvidenceAverageReliability float64
	SourceDiversityRatio       float64
	PlayerCredibility          float64
	CorroborationScore         float64
}

// CalculateCorrectionStrength evaluates the overall strength of a published correction.
func CalculateCorrectionStrength(p CorrectionParams) float64 {
	strength := p.EvidenceAverageReliability*0.45 +
		p.SourceDiversityRatio*0.20 +
		p.PlayerCredibility*0.20 +
		p.CorroborationScore*0.15
	return math.Min(1.0, math.Max(0.0, strength))
}

// ApplyCorrectionImpact computes the reduction in confidence for an agent who receives the correction.
func ApplyCorrectionImpact(oldConfidence float64, correctionStrength float64, listenerReceptiveness float64, evidenceResistance float64) float64 {
	// Higher receptiveness and lower resistance produce a greater reduction in false belief
	impactMultiplier := 0.55 + (listenerReceptiveness * 0.30) - (evidenceResistance * 0.25)
	if impactMultiplier < 0.1 {
		impactMultiplier = 0.1
	}

	reduction := oldConfidence * correctionStrength * impactMultiplier
	newConf := oldConfidence - reduction
	return domain.ClampConfidence(newConf)
}

// UpdatePlayerCredibility adjusts credibility based on evidence quality.
func UpdatePlayerCredibility(currentCredibility float64, correctionStrength float64, hasContradictoryEvidence bool) float64 {
	var delta float64
	if hasContradictoryEvidence || correctionStrength < 0.30 {
		delta = -0.12 // penalty for weak or unsubstantiated correction
	} else if correctionStrength >= 0.70 {
		delta = 0.08 // bonus for high quality evidence
	} else {
		delta = 0.03 // moderate improvement
	}

	return domain.ClampConfidence(currentCredibility + delta)
}
