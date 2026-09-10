package scheduler

import (
	"testing"
)

func TestSchedulerInteractions(t *testing.T) {
	rng := NewSeededRNG(42)
	sched := NewScheduler(rng, 90, 2)

	agents := []AgentSchedulingView{
		{
			ID:              "a1",
			ScenarioAgentID: "agent_a",
			LocationID:      "cafe",
			Sociality:       0.9,
			Beliefs:         map[string]float64{"claim_chemical_explosion": 0.8},
		},
		{
			ID:              "a2",
			ScenarioAgentID: "agent_b",
			LocationID:      "cafe",
			Sociality:       0.8,
			Beliefs:         map[string]float64{},
		},
		{
			ID:              "a3",
			ScenarioAgentID: "agent_c",
			LocationID:      "market",
			Sociality:       0.5,
			Beliefs:         map[string]float64{},
		},
	}

	rels := map[string]map[string]float64{
		"a1": {"a2": 0.8},
	}

	interactions := sched.PlanTickInteractions(agents, rels, 10)
	if len(interactions) == 0 {
		t.Fatalf("expected at least 1 scheduled interaction, got 0")
	}

	first := interactions[0]
	if first.SpeakerAgentID != "a1" || first.ListenerAgentID != "a2" {
		t.Errorf("expected a1 -> a2 interaction, got %s -> %s", first.SpeakerAgentID, first.ListenerAgentID)
	}

	// Immediately calling at second 11 should respect pair cooldown
	nextInteractions := sched.PlanTickInteractions(agents, rels, 11)
	for _, in := range nextInteractions {
		if (in.SpeakerAgentID == "a1" && in.ListenerAgentID == "a2") || (in.SpeakerAgentID == "a2" && in.ListenerAgentID == "a1") {
			t.Errorf("expected pair cooldown to prevent a1 and a2 from speaking again at second 11")
		}
	}
}
