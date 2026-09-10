package evidence

import (
	"math"
)

type EvidenceEvaluationItem struct {
	ID             string
	Reliability    float64
	SourceCategory string
	CoveredFields  []string
}

// CalculateEvidenceStrength calculates the composite strength score of discovered evidence.
func CalculateEvidenceStrength(items []EvidenceEvaluationItem, requiredFields []string) (float64, float64, float64, float64) {
	if len(requiredFields) == 0 {
		return 0, 0, 0, 0
	}

	// 1. Coverage
	coveredMap := make(map[string]bool)
	sumReliability := 0.0
	categoryMap := make(map[string]bool)

	for _, item := range items {
		for _, f := range item.CoveredFields {
			coveredMap[f] = true
		}
		sumReliability += item.Reliability
		categoryMap[item.SourceCategory] = true
	}

	coveredCount := 0
	for _, rf := range requiredFields {
		if coveredMap[rf] {
			coveredCount++
		}
	}
	coverageScore := float64(coveredCount) / float64(len(requiredFields))

	// 2. Reliability
	reliabilityScore := 0.0
	if len(items) > 0 {
		reliabilityScore = sumReliability / float64(len(items))
	}

	// 3. Diversity (4 or more distinct categories = 1.0)
	diversityScore := math.Min(1.0, float64(len(categoryMap))/4.0)

	// Composite
	totalStrength := (coverageScore * 0.50) + (reliabilityScore * 0.30) + (diversityScore * 0.20)
	totalStrength = math.Min(1.0, math.Max(0.0, totalStrength))

	return totalStrength, coverageScore, reliabilityScore, diversityScore
}

// MapEvidenceToFields determines which ground truth fields an evidence item informs.
func MapEvidenceToFields(scenarioEvidenceID string) []string {
	switch scenarioEvidenceID {
	case "ev_cctv_gate":
		return []string{"event_type", "major_explosion"}
	case "ev_fire_report":
		return []string{"event_type", "cause_category"}
	case "ev_hospital_summary":
		return []string{"fatalities"}
	case "ev_maintenance_ticket":
		return []string{"cause_category"}
	case "ev_panel_photo":
		return []string{"cause_category"}
	case "ev_photo_sequence":
		return []string{"event_type", "major_explosion"}
	case "ev_worker_testimony":
		return []string{"cause_category", "fatalities"}
	case "ev_dispatch_log":
		return []string{"fatalities"}
	default:
		return nil
	}
}
