package evidence

import (
	"testing"
)

func TestEvidenceStrengthScoring(t *testing.T) {
	reqFields := []string{"event_type", "cause_category", "major_explosion", "fatalities"}

	// Test empty evidence
	strengthEmpty, cov, rel, div := CalculateEvidenceStrength(nil, reqFields)
	if strengthEmpty != 0 || cov != 0 || rel != 0 || div != 0 {
		t.Errorf("expected 0 for empty evidence, got strength %f", strengthEmpty)
	}

	// Test comprehensive evidence items
	items := []EvidenceEvaluationItem{
		{ID: "ev_cctv_gate", Reliability: 0.94, SourceCategory: "visual_record", CoveredFields: []string{"event_type", "major_explosion"}},
		{ID: "ev_fire_report", Reliability: 0.98, SourceCategory: "official_response", CoveredFields: []string{"event_type", "cause_category"}},
		{ID: "ev_hospital_summary", Reliability: 0.99, SourceCategory: "medical", CoveredFields: []string{"fatalities"}},
		{ID: "ev_maintenance_ticket", Reliability: 0.91, SourceCategory: "technical_record", CoveredFields: []string{"cause_category"}},
	}

	totalStrength, coverage, reliability, diversity := CalculateEvidenceStrength(items, reqFields)
	if coverage != 1.0 {
		t.Errorf("expected 100%% coverage, got %f", coverage)
	}
	if reliability <= 0 {
		t.Errorf("expected reliability > 0, got %f", reliability)
	}
	if diversity != 1.0 {
		t.Errorf("expected 100%% diversity with 4 categories, got %f", diversity)
	}
	if totalStrength < 0.70 {
		t.Errorf("expected total strength >= 0.70, got %f", totalStrength)
	}
}
