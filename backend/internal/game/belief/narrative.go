package belief

type NarrativeClaimWeight struct {
	ClaimID string
	Weight  float64
}

// CalculateAgentNarrativeScore computes the weighted sum of an agent's confidence across primary false claims.
func CalculateAgentNarrativeScore(claims []NarrativeClaimWeight, beliefs map[string]float64) float64 {
	score := 0.0
	for _, cw := range claims {
		if conf, exists := beliefs[cw.ClaimID]; exists {
			score += conf * cw.Weight
		}
	}
	return score
}

// IsNarrativeAdopted checks if an agent's narrative score crosses the adoption threshold (default 0.65).
func IsNarrativeAdopted(score float64, threshold float64) bool {
	if threshold <= 0 {
		threshold = 0.65
	}
	return score >= threshold
}

// CalculateFalseNarrativeRatio returns the proportion of agents who have adopted the false narrative.
func CalculateFalseNarrativeRatio(adoptedCount int, totalAgents int) float64 {
	if totalAgents <= 0 {
		return 0.0
	}
	return float64(adoptedCount) / float64(totalAgents)
}
