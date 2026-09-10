package outcome

import (
	"city-of-lies/backend/internal/domain"
)

type EvaluatorConfig struct {
	TruthEvidenceThreshold float64 // default: 0.70
	DefeatFalseBeliefRatio float64 // default: 0.75
}

// CheckLossCondition evaluates if the false narrative has reached or crossed the defeat threshold.
func CheckLossCondition(falseBeliefRatio float64, threshold float64) bool {
	if threshold <= 0 {
		threshold = 0.75
	}
	return falseBeliefRatio >= threshold
}

type GroundTruthAnswers struct {
	EventType      string
	CauseCategory  string
	MajorExplosion bool
	Fatalities     int
}

// EvaluateSubmission verifies player answers, evidence strength, and current rumor ratio.
func EvaluateSubmission(
	sub domain.GroundTruthSubmission,
	expected GroundTruthAnswers,
	evidenceStrength float64,
	falseBeliefRatio float64,
	cfg EvaluatorConfig,
) domain.TruthSubmissionResult {
	if cfg.TruthEvidenceThreshold <= 0 {
		cfg.TruthEvidenceThreshold = 0.70
	}
	if cfg.DefeatFalseBeliefRatio <= 0 {
		cfg.DefeatFalseBeliefRatio = 0.75
	}

	fieldStatus := make(map[string]bool)
	fieldStatus["event_type"] = (sub.EventType == expected.EventType)
	fieldStatus["cause_category"] = (sub.CauseCategory == expected.CauseCategory)
	fieldStatus["major_explosion"] = (sub.MajorExplosion == expected.MajorExplosion)
	fieldStatus["fatalities"] = (sub.Fatalities == expected.Fatalities)

	allFieldsCorrect := fieldStatus["event_type"] &&
		fieldStatus["cause_category"] &&
		fieldStatus["major_explosion"] &&
		fieldStatus["fatalities"]

	// Check immediate loss precedence
	if falseBeliefRatio >= cfg.DefeatFalseBeliefRatio {
		return domain.TruthSubmissionResult{
			Won:              false,
			Terminal:         true,
			EvidenceStrength: evidenceStrength,
			RequiredEvidence: cfg.TruthEvidenceThreshold,
			FalseBeliefRatio: falseBeliefRatio,
			FieldStatus:      fieldStatus,
			FeedbackMessage:  "The false narrative has already dominated public consensus. Investigation failed.",
		}
	}

	if !allFieldsCorrect {
		return domain.TruthSubmissionResult{
			Won:              false,
			Terminal:         false,
			EvidenceStrength: evidenceStrength,
			RequiredEvidence: cfg.TruthEvidenceThreshold,
			FalseBeliefRatio: falseBeliefRatio,
			FieldStatus:      fieldStatus,
			FeedbackMessage:  "Some elements of your factual report are inaccurate or unverified.",
		}
	}

	if evidenceStrength < cfg.TruthEvidenceThreshold {
		return domain.TruthSubmissionResult{
			Won:              false,
			Terminal:         false,
			EvidenceStrength: evidenceStrength,
			RequiredEvidence: cfg.TruthEvidenceThreshold,
			FalseBeliefRatio: falseBeliefRatio,
			FieldStatus:      fieldStatus,
			FeedbackMessage:  "Your claims are true, but your corroborating evidence is not yet sufficient to convince the public.",
		}
	}

	// Victorious conclusion
	return domain.TruthSubmissionResult{
		Won:              true,
		Terminal:         true,
		EvidenceStrength: evidenceStrength,
		RequiredEvidence: cfg.TruthEvidenceThreshold,
		FalseBeliefRatio: falseBeliefRatio,
		FieldStatus:      fieldStatus,
		FeedbackMessage:  "Truth established! The evidence successfully dismantled the false narrative.",
	}
}
