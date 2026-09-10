package belief

import (
	"testing"
)

func TestNarrativeScoreAndAdoption(t *testing.T) {
	weights := []NarrativeClaimWeight{
		{ClaimID: "claim_chemical_explosion", Weight: 0.35},
		{ClaimID: "claim_multiple_deaths", Weight: 0.40},
		{ClaimID: "claim_company_coverup", Weight: 0.25},
	}

	beliefsBelow := map[string]float64{
		"claim_chemical_explosion": 0.60, // 0.21
		"claim_multiple_deaths":    0.60, // 0.24
		"claim_company_coverup":    0.60, // 0.15 => total 0.60
	}
	scoreBelow := CalculateAgentNarrativeScore(weights, beliefsBelow)
	if IsNarrativeAdopted(scoreBelow, 0.65) {
		t.Errorf("expected score %f to not be adopted at 0.65", scoreBelow)
	}

	beliefsAbove := map[string]float64{
		"claim_chemical_explosion": 0.70, // 0.245
		"claim_multiple_deaths":    0.70, // 0.280
		"claim_company_coverup":    0.60, // 0.150 => total 0.675
	}
	scoreAbove := CalculateAgentNarrativeScore(weights, beliefsAbove)
	if !IsNarrativeAdopted(scoreAbove, 0.65) {
		t.Errorf("expected score %f to be adopted at 0.65", scoreAbove)
	}

	// 14/20 = 70% (no defeat)
	ratio14 := CalculateFalseNarrativeRatio(14, 20)
	if ratio14 >= 0.75 {
		t.Errorf("expected 14/20 to be < 0.75, got %f", ratio14)
	}

	// 15/20 = 75% (defeat)
	ratio15 := CalculateFalseNarrativeRatio(15, 20)
	if ratio15 < 0.75 {
		t.Errorf("expected 15/20 to be >= 0.75, got %f", ratio15)
	}
}
