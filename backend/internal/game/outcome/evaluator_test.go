package outcome

import (
	"testing"

	"city-of-lies/backend/internal/domain"
)

func TestOutcomeEvaluations(t *testing.T) {
	cfg := EvaluatorConfig{
		TruthEvidenceThreshold: 0.70,
		DefeatFalseBeliefRatio: 0.75,
	}

	answers := GroundTruthAnswers{
		EventType:      "fire",
		CauseCategory:  "electrical_fault",
		MajorExplosion: false,
		Fatalities:     0,
	}

	correctSub := domain.GroundTruthSubmission{
		EventType:      "fire",
		CauseCategory:  "electrical_fault",
		MajorExplosion: false,
		Fatalities:     0,
	}

	// 1. Correct answers + 0.69 evidence => Not won
	resLowEv := EvaluateSubmission(correctSub, answers, 0.69, 0.40, cfg)
	if resLowEv.Won || resLowEv.Terminal {
		t.Errorf("expected not won due to low evidence (0.69), got won=%v", resLowEv.Won)
	}

	// 2. Correct answers + 0.70 evidence + false ratio 0.74 => Win
	resWin := EvaluateSubmission(correctSub, answers, 0.70, 0.74, cfg)
	if !resWin.Won || !resWin.Terminal {
		t.Errorf("expected win with 0.70 evidence and 0.74 false ratio, got won=%v", resWin.Won)
	}

	// 3. Correct answers + 1.0 evidence + false ratio 0.75 => Loss
	resLost := EvaluateSubmission(correctSub, answers, 1.0, 0.75, cfg)
	if resLost.Won || !resLost.Terminal {
		t.Errorf("expected loss when false ratio >= 0.75, got won=%v terminal=%v", resLost.Won, resLost.Terminal)
	}

	// 4. Incorrect answer => Never wins
	wrongSub := correctSub
	wrongSub.EventType = "chemical_explosion"
	resWrong := EvaluateSubmission(wrongSub, answers, 1.0, 0.20, cfg)
	if resWrong.Won {
		t.Errorf("expected failure on wrong event type, got won=true")
	}
}
